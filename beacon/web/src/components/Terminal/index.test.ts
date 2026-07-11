import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import type { ComponentPublicInstance } from 'vue'

// --- xterm mocks ---------------------------------------------------------
const termClear = vi.fn()
const termWriteln = vi.fn()
const termWrite = vi.fn()
const termOnData = vi.fn()
const termLoadAddon = vi.fn()
const termOpen = vi.fn()
const termDispose = vi.fn()
const termReset = vi.fn()
const termInstance = {
  rows: 40,
  cols: 120,
  onData: termOnData,
  loadAddon: termLoadAddon,
  open: termOpen,
  clear: termClear,
  writeln: termWriteln,
  write: termWrite,
  reset: termReset,
  dispose: termDispose,
}

vi.mock('@xterm/xterm', () => ({
  Terminal: vi.fn(() => termInstance),
}))
vi.mock('@xterm/addon-fit', () => ({
  FitAddon: vi.fn(() => ({ fit: vi.fn(), dispose: vi.fn() })),
}))
vi.mock('@xterm/addon-web-links', () => ({
  WebLinksAddon: vi.fn(() => ({})),
}))
vi.mock('@xterm/xterm/css/xterm.css', () => ({}))

vi.mock('@/store/modules/user', () => ({
  useUserStore: () => ({ token: 'tok' }),
}))

// --- WebSocket mock ------------------------------------------------------
interface FakeWS {
  url: string
  readyState: number
  onopen: ((ev: Event) => void) | null
  onmessage: ((ev: MessageEvent) => void) | null
  onclose: ((ev: CloseEvent) => void) | null
  onerror: ((ev: Event) => void) | null
  send: ReturnType<typeof vi.fn>
  close: ReturnType<typeof vi.fn>
}
const instances: FakeWS[] = []
let autoOpen = true
class FakeWebSocket {
  static OPEN = 1
  static CLOSED = 3
  onopen: FakeWS['onopen'] = null
  onmessage: FakeWS['onmessage'] = null
  onclose: FakeWS['onclose'] = null
  onerror: FakeWS['onerror'] = null
  url: string
  readyState = 0
  send = vi.fn()
  close = vi.fn(() => {
    this.readyState = FakeWebSocket.CLOSED
    // 模拟真实浏览器：close() 异步触发 onclose，用以暴露
    // newSession/agentId 切换时旧连接回调污染新状态的竞态。
    Promise.resolve().then(() => this.onclose?.(new CloseEvent('close')))
  })

  constructor(url: string) {
    this.url = url
    instances.push(this as unknown as FakeWS)
    if (autoOpen) {
      Promise.resolve().then(() => {
        this.readyState = FakeWebSocket.OPEN
        this.onopen?.(new Event('open'))
      })
    }
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  instances.length = 0
  autoOpen = true
  ;(globalThis as unknown as { WebSocket: typeof WebSocket }).WebSocket
    = FakeWebSocket as unknown as typeof WebSocket
})

afterEach(() => {
  vi.useRealTimers()
})

interface TerminalVM extends ComponentPublicInstance {
  status: string
  cols: number
  rows: number
  clear: () => void
  newSession: () => void
}

async function mountTerminal(agentId = 'agent-a') {
  const wrapper = mount((await import('./index.vue')).default, {
    props: { agentId },
    global: { stubs: { teleport: true } },
  })
  await vi.dynamicImportSettled()
  return wrapper
}

async function flush() {
  await new Promise(r => setTimeout(r, 0))
}

describe('Terminal component', () => {
  it('transitions status connecting -> connected when ws opens', async () => {
    const wrapper = await mountTerminal()
    await flush()
    expect((wrapper.vm as TerminalVM).status).toBe('connected')
    expect(instances).toHaveLength(1)
  })

  it('transitions to closed when ws closes', async () => {
    const wrapper = await mountTerminal()
    await flush()
    const ws = instances[0]!
    ws.onclose?.(new CloseEvent('close'))
    expect((wrapper.vm as TerminalVM).status).toBe('closed')
  })

  it('transitions to error when ws errors', async () => {
    const wrapper = await mountTerminal()
    await flush()
    const ws = instances[0]!
    ws.onerror?.(new Event('error'))
    expect((wrapper.vm as TerminalVM).status).toBe('error')
  })

  it('clear() delegates to xterm clear', async () => {
    const wrapper = await mountTerminal()
    await flush()
    ;(wrapper.vm as TerminalVM).clear()
    expect(termClear).toHaveBeenCalled()
  })

  it('newSession() closes the existing ws and opens a new connection', async () => {
    const wrapper = await mountTerminal()
    await flush()
    const first = instances[0]!
    expect(instances).toHaveLength(1)

    ;(wrapper.vm as TerminalVM).newSession()
    await flush()

    expect(first.close).toHaveBeenCalled()
    expect(instances.length).toBeGreaterThanOrEqual(2)
  })

  it('newSession() keeps status connecting/connected despite the old ws onclose firing', async () => {
    const wrapper = await mountTerminal()
    await flush()

    ;(wrapper.vm as TerminalVM).newSession()
    await flush()

    // 旧连接的 onclose 已异步触发，但代际守卫使其不得覆盖新连接状态。
    expect((wrapper.vm as TerminalVM).status).not.toBe('closed')
  })

  it('switching agentId does not let the old ws onclose clobber the new status', async () => {
    const wrapper = await mountTerminal('agent-a')
    await flush()

    await wrapper.setProps({ agentId: 'agent-b' })
    await flush()

    expect(instances.length).toBeGreaterThanOrEqual(2)
    expect((wrapper.vm as TerminalVM).status).not.toBe('closed')
  })

  it('exposes cols/rows derived from the xterm instance', async () => {
    const wrapper = await mountTerminal()
    await flush()
    expect((wrapper.vm as TerminalVM).cols).toBe(120)
    expect((wrapper.vm as TerminalVM).rows).toBe(40)
  })
})
