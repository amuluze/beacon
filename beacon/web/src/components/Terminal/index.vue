<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'
import { useUserStore } from '@/store/modules/user'

export type TerminalStatus = 'idle' | 'connecting' | 'connected' | 'closed' | 'error'

interface Props {
  agentId: string
}

const props = defineProps<Props>()
const terminalRef = ref<HTMLDivElement>()

const status = ref<TerminalStatus>('idle')
const cols = ref(0)
const rows = ref(0)

let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null
// 连接代际：递增后旧连接的回调一律丢弃，避免 newSession/agentId 切换时
// 旧 ws 的异步 onclose/onerror 覆盖新连接的 connecting/connected 状态。
let wsGeneration = 0

function buildURL(agentId: string): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // 浏览器 ws 握手无法携带 Authorization header，把 access token 写入 query 通过服务端鉴权。
  const token = useUserStore().token
  const tokenQuery = token ? `&token=${encodeURIComponent(token)}` : ''
  return `${protocol}//${window.location.host}/ws/terminal?agent_id=${encodeURIComponent(agentId)}${tokenQuery}`
}

function sendMessage(type: string, data?: unknown): void {
  if (!ws || ws.readyState !== WebSocket.OPEN)
    return
  ws.send(JSON.stringify({ type, data }))
}

function sendResize(rows: number, cols: number): void {
  if (!ws || ws.readyState !== WebSocket.OPEN)
    return
  ws.send(JSON.stringify({ type: 'resize', rows, cols }))
}

function syncDims(): void {
  if (!term)
    return
  cols.value = term.cols
  rows.value = term.rows
}

function initTerminal(): void {
  term = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
    },
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())

  if (terminalRef.value)
    term.open(terminalRef.value)

  term.onData((input: string) => {
    const encoded = btoa(unescape(encodeURIComponent(input)))
    sendMessage('input', encoded)
  })

  fitAddon.fit()
  syncDims()
  sendResize(term.rows, term.cols)
}

function connect(): void {
  if (!term)
    return

  const gen = ++wsGeneration
  function guard<T extends (...args: any[]) => void>(fn: T): T {
    return ((...args: any[]) => {
      if (gen !== wsGeneration)
        return
      fn(...args)
    }) as T
  }

  status.value = 'connecting'
  const url = buildURL(props.agentId)
  ws = new WebSocket(url)

  ws.onopen = guard(() => {
    status.value = 'connected'
    term?.writeln('\r\n\x1B[32mConnected to agent terminal.\x1B[0m')
  })

  ws.onmessage = guard((event: MessageEvent) => {
    let msg
    try {
      msg = JSON.parse(event.data)
    }
    catch {
      term?.write(event.data)
      return
    }

    if (msg.type === 'output' && msg.data) {
      try {
        const decoded = decodeURIComponent(escape(atob(msg.data)))
        term?.write(decoded)
      }
      catch {
        term?.write(msg.data)
      }
    }
    else if (msg.type === 'error') {
      term?.writeln(`\r\n\x1B[31mError: ${msg.msg}\x1B[0m`)
    }
  })

  ws.onclose = guard(() => {
    status.value = 'closed'
    term?.writeln('\r\n\x1B[31mConnection closed.\x1B[0m')
  })

  ws.onerror = guard(() => {
    status.value = 'error'
    term?.writeln('\r\n\x1B[31mWebSocket error.\x1B[0m')
  })
}

function clear(): void {
  term?.clear()
}

function newSession(): void {
  if (ws && ws.readyState !== WebSocket.CLOSED)
    ws.close()
  term?.clear()
  connect()
}

function onResize(): void {
  if (!term || !fitAddon)
    return
  fitAddon.fit()
  syncDims()
  sendResize(term.rows, term.cols)
}

onMounted(() => {
  initTerminal()
  connect()
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  wsGeneration++
  window.removeEventListener('resize', onResize)
  ws?.close()
  term?.dispose()
})

watch(() => props.agentId, () => {
  ws?.close()
  term?.clear()
  connect()
})

defineExpose({ status, cols, rows, clear, newSession })
</script>

<template>
    <div ref="terminalRef" class="am-terminal" />
</template>

<style scoped lang="scss">
.am-terminal {
  width: 100%;
  height: 100%;
  min-height: 0;
  padding: 8px 12px;
  background: #1e1e1e;
}
</style>
