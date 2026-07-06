<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'
import { useUserStore } from '@/store/modules/user'

interface Props {
  agentId: string
}

const props = defineProps<Props>()
const terminalRef = ref<HTMLDivElement>()

let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null

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
  const dims = term
  sendResize(dims.rows, dims.cols)
}

function connect(): void {
  if (!term)
    return

  const url = buildURL(props.agentId)
  ws = new WebSocket(url)

  ws.onopen = () => {
    term?.writeln('\r\n\x1B[32mConnected to agent terminal.\x1B[0m')
  }

  ws.onmessage = (event: MessageEvent) => {
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
  }

  ws.onclose = () => {
    term?.writeln('\r\n\x1B[31mConnection closed.\x1B[0m')
  }

  ws.onerror = (err: Event) => {
    console.error('terminal websocket error', err)
    term?.writeln('\r\n\x1B[31mWebSocket error.\x1B[0m')
  }
}

function onResize(): void {
  if (!term || !fitAddon)
    return
  fitAddon.fit()
  sendResize(term.rows, term.cols)
}

onMounted(() => {
  initTerminal()
  connect()
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  ws?.close()
  term?.dispose()
})

watch(() => props.agentId, () => {
  ws?.close()
  connect()
})
</script>

<template>
    <div ref="terminalRef" class="am-terminal" />
</template>

<style scoped lang="scss">
.am-terminal {
  width: 100%;
  height: 100%;
  min-height: 400px;
  background: #1e1e1e;
  padding: 8px;
  border-radius: 4px;
}
</style>
