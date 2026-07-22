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
let termOnDataHandler: ((input: string) => void) | undefined
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
        void Promise.resolve().then(() => this.onclose?.(new CloseEvent('close')))
    })

    constructor(url: string) {
        this.url = url
        instances.push(this)
        if (autoOpen) {
            void Promise.resolve().then(() => {
                this.readyState = FakeWebSocket.OPEN
                this.onopen?.(new Event('open'))
            })
        }
    }
}

beforeEach(() => {
    vi.clearAllMocks()
    termOnDataHandler = undefined
    termOnData.mockImplementation((handler: (input: string) => void) => {
        termOnDataHandler = handler
        return { dispose: vi.fn() }
    })
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

describe('terminal component', () => {
    it('stays connecting after ws opens and only connects after the Agent PTY ready message', async () => {
        const wrapper = await mountTerminal()
        await flush()

        expect((wrapper.vm as TerminalVM).status).toBe('connecting')
        instances[0].onmessage?.(new MessageEvent('message', { data: JSON.stringify({ type: 'ready' }) }))

        expect((wrapper.vm as TerminalVM).status).toBe('connected')
        expect(termWriteln).toHaveBeenCalledWith(expect.stringContaining('Connected to agent terminal.'))
        expect(instances).toHaveLength(1)
    })

    it('does not forward input before ready and sends base64 input after ready', async () => {
        await mountTerminal()
        await flush()
        const ws = instances[0]

        termOnDataHandler?.('echo blocked\r')
        expect(ws.send).not.toHaveBeenCalled()

        ws.onmessage?.(new MessageEvent('message', { data: JSON.stringify({ type: 'ready' }) }))
        ws.send.mockClear()
        termOnDataHandler?.('echo CODEX_TERMINAL_OK\r')

        expect(ws.send).toHaveBeenCalledTimes(1)
        expect(JSON.parse(ws.send.mock.calls[0][0] as string)).toEqual({
            type: 'input',
            data: btoa('echo CODEX_TERMINAL_OK\r'),
        })
    })

    it('includes agent, token, and initial dimensions in the WebSocket URL', async () => {
        await mountTerminal('node-01')
        await flush()

        const url = new URL(instances[0].url)
        expect(url.pathname).toBe('/ws/terminal')
        expect(url.searchParams.get('agent_id')).toBe('node-01')
        expect(url.searchParams.get('token')).toBe('tok')
        expect(url.searchParams.get('rows')).toBe('40')
        expect(url.searchParams.get('cols')).toBe('120')
    })

    it('transitions to closed when ws closes', async () => {
        const wrapper = await mountTerminal()
        await flush()
        const ws = instances[0]
        ws.onclose?.(new CloseEvent('close'))
        expect((wrapper.vm as TerminalVM).status).toBe('closed')
    })

    it('transitions to error when ws errors', async () => {
        const wrapper = await mountTerminal()
        await flush()
        const ws = instances[0]
        ws.onerror?.(new Event('error'))
        expect((wrapper.vm as TerminalVM).status).toBe('error')
    })

    it('clear() delegates to xterm clear', async () => {
        const wrapper = await mountTerminal()
        await flush()
        ;(wrapper.vm as TerminalVM).clear()
        expect(termClear).toHaveBeenCalled()
        expect(termWrite).toHaveBeenCalledWith('\x1B[2J\x1B[H')
    })

    it('newSession() closes the existing ws and opens a new connection', async () => {
        const wrapper = await mountTerminal()
        await flush()
        const first = instances[0]
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
