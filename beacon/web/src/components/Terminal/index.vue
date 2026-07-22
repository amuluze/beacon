<script setup lang="ts">
import { onBeforeUnmount, onMounted, shallowRef, useTemplateRef, watch } from 'vue'
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
const emit = defineEmits<{
  statusChange: [status: TerminalStatus]
  resize: [size: { rows: number, cols: number }]
}>()
const terminalRef = useTemplateRef<HTMLDivElement>('terminal')

const status = shallowRef<TerminalStatus>('idle')
const cols = shallowRef(0)
const rows = shallowRef(0)

let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null
// 连接代际：递增后旧连接的回调一律丢弃，避免 newSession/agentId 切换时
// 旧 ws 的异步 onclose/onerror 覆盖新连接的 connecting/connected 状态。
let wsGeneration = 0

function buildURL(agentId: string): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const params = new URLSearchParams({
    agent_id: agentId,
    rows: String(term?.rows || rows.value || 24),
    cols: String(term?.cols || cols.value || 80),
  })
  // 浏览器 ws 握手无法携带 Authorization header，把 access token 写入 query 通过服务端鉴权。
  const token = useUserStore().token
  if (token)
    params.set('token', token)
  return `${protocol}//${window.location.host}/ws/terminal?${params.toString()}`
}

function setStatus(nextStatus: TerminalStatus): void {
  if (status.value === nextStatus)
    return
  status.value = nextStatus
  emit('statusChange', nextStatus)
}

function sendInput(data: string): void {
  if (!ws || ws.readyState !== WebSocket.OPEN || status.value !== 'connected')
    return
  ws.send(JSON.stringify({ type: 'input', data }))
}

function sendResize(rows: number, cols: number): void {
  if (!ws || ws.readyState !== WebSocket.OPEN || status.value !== 'connected')
    return
  ws.send(JSON.stringify({ type: 'resize', rows, cols }))
}

function syncDims(): void {
  if (!term)
    return
  cols.value = term.cols
  rows.value = term.rows
  emit('resize', { rows: rows.value, cols: cols.value })
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
    sendInput(encoded)
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

  setStatus('connecting')
  const url = buildURL(props.agentId)
  ws = new WebSocket(url)

  // HTTP upgrade only means the Server accepted the socket. The component
  // remains connecting until the Server confirms the Agent PTY is ready.
  ws.onopen = guard(() => {})

  ws.onmessage = guard((event: MessageEvent) => {
    let msg
    try {
      msg = JSON.parse(event.data)
    }
    catch {
      term?.write(event.data)
      return
    }

    if (msg.type === 'ready') {
      setStatus('connected')
      term?.writeln('\r\n\x1B[32mConnected to agent terminal.\x1B[0m')
      if (term)
        sendResize(term.rows, term.cols)
    }
    else if (msg.type === 'output' && msg.data) {
      try {
        const decoded = decodeURIComponent(escape(atob(msg.data)))
        term?.write(decoded)
      }
      catch {
        term?.write(msg.data)
      }
    }
    else if (msg.type === 'error') {
      setStatus('error')
      term?.writeln(`\r\n\x1B[31mError: ${msg.msg}\x1B[0m`)
    }
  })

  ws.onclose = guard(() => {
    if (status.value !== 'error')
      setStatus('closed')
    term?.writeln('\r\n\x1B[31mConnection closed.\x1B[0m')
  })

  ws.onerror = guard(() => {
    setStatus('error')
    term?.writeln('\r\n\x1B[31mWebSocket error.\x1B[0m')
  })
}

function clear(): void {
  if (!term)
    return
  term.clear()
  term.write('\x1B[2J\x1B[H')
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
  setStatus('closed')
})

watch(() => props.agentId, () => {
  ws?.close()
  term?.clear()
  connect()
})

defineExpose({ status, cols, rows, clear, newSession })
</script>

<template>
    <div ref="terminal" class="am-terminal" />
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
