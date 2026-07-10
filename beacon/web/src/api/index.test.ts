/**
 * Tests for the Axios request / response interceptors in api/index.ts.
 *
 * Strategy:
 *  - Stub `@/store` so the whole pinia chain is bypassed (avoids pulling in
 *    unrelated store modules and the `defineStore` global registration).
 *  - Stub `axios` so we control the axios instance returned by `axios.create()`
 *    and can introspect its interceptor arrays.
 *  - Stub `@/components/Message/message.ts` so the warning() helper is a
 *    vi.fn() and not noisy / harmful in jsdom.
 *  - Drive the registered interceptor callbacks manually against a stub store
 *    object; never issue real HTTP traffic.
 */
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Stub-pinia store handle exposed via the shared store.
type StubStore = {
    user: { token: string; refresh: string; setToken: ReturnType<typeof vi.fn> }
    agent: { selectedAgentID: string }
}
const stubStore: StubStore = {
    user: { token: '', refresh: '', setToken: vi.fn() },
    agent: { selectedAgentID: '' },
}

vi.mock('@/store', () => ({
    default: () => stubStore,
}))

type ReqHandler = (cfg: any) => any
type ResErrHandler = (err: any) => any
const requestHandlers: ReqHandler[] = []
const responseErrHandlers: ResErrHandler[] = []

// Each call to axios.create() returns a fresh stub instance whose interceptor
// callbacks are tracked in module-level arrays.
let latestFakeInstance: any = null
function makeFakeAxiosInstance() {
    const inst: any = {
        interceptors: {
            request: { use: (ok: ReqHandler) => { requestHandlers.push(ok) } },
            response: {
                use: (_ok: any, err: ResErrHandler) => { responseErrHandlers.push(err) },
            },
        },
        get: vi.fn(),
        post: vi.fn().mockResolvedValue({ data: { access_token: 'new', refresh_token: 'new-r' } }),
    }
    latestFakeInstance = inst
    return inst
}

vi.mock('axios', () => ({
    default: {
        create: () => makeFakeAxiosInstance(),
    },
}))

const warning = vi.fn()
vi.mock('@/components/Message/message.ts', () => ({
    warning: (...args: any[]) => warning(...args),
}))

let lastHref: string | null = null
Object.defineProperty(window, 'location', {
    configurable: true,
    get: () => ({
        set href(v: string) { lastHref = v },
        get href() { return lastHref ?? '' },
    }) as any,
})

window.URL.createObjectURL = vi.fn(() => 'blob:fake')
window.URL.revokeObjectURL = vi.fn()
HTMLAnchorElement.prototype.click = vi.fn()

beforeEach(() => {
    stubStore.user.token = ''
    stubStore.user.refresh = ''
    stubStore.user.setToken = vi.fn()
    stubStore.agent.selectedAgentID = ''
    requestHandlers.length = 0
    responseErrHandlers.length = 0
    lastHref = null
    warning.mockClear()
    vi.clearAllMocks()
    // Re-evaluate @/api/index so a fresh Request singleton is constructed
    // and its interceptors are pushed into our arrays via the mocked axios.create.
    vi.resetModules()
})

afterEach(() => {
    vi.useRealTimers()
})

async function loadRequest() {
    const mod = await import('@/api/index')
    return mod.default
}

describe('api/index — request interceptor', () => {
    it('injects X-Agent-ID from stubStore.agent.selectedAgentID', async () => {
        stubStore.agent.selectedAgentID = 'agent-1'
        stubStore.user.token = ''

        await loadRequest()
        expect(requestHandlers.length).toBe(1)

        const out = await requestHandlers[0]({ headers: {}, url: '/api/v1/host/info' })
        expect(out.headers['X-Agent-ID']).toBe('agent-1')
    })

    it('injects Authorization from stubStore.user.token', async () => {
        stubStore.user.token = 'abc123'
        stubStore.agent.selectedAgentID = ''

        await loadRequest()
        const out = await requestHandlers[0]({ headers: {}, url: '/api/v1/host/info' })
        expect(out.headers.Authorization).toBe('Bearer abc123')
    })

    it('rewrites Content-Type for the login endpoint', async () => {
        stubStore.user.token = 'should-stay'
        stubStore.agent.selectedAgentID = 'agent-x'

        await loadRequest()
        const out = await requestHandlers[0]({ headers: {}, url: '/api/v1/auth/login' })
        expect(out.headers['Content-Type']).toContain('multipart/form-data')
        expect(out.headers['X-Agent-ID']).toBeUndefined()
        expect(out.headers.Authorization).toBe('Bearer should-stay')
    })
})

describe('api/index — response error interceptor', () => {
    it('400 on /api/v1/auth/token_update clears tokens and navigates to /', async () => {
        await loadRequest()
        expect(responseErrHandlers.length).toBe(1)

        const error = {
            response: {
                status: 400,
                data: { msg: 'invalid' },
                config: { url: '/api/v1/auth/token_update' },
            },
        }
        void responseErrHandlers[0](error)
        expect(lastHref).toBe('/')
        expect(stubStore.user.setToken).toHaveBeenCalledWith('', '')
    })

    it('400 on a non-token_update URL forwards to warning()', async () => {
        vi.useFakeTimers()
        await loadRequest()
        const error = {
            response: {
                status: 400,
                data: { msg: 'invalid input' },
                config: { url: '/api/v1/host/report' },
            },
        }
        void responseErrHandlers[0](error)
        expect(warning).toHaveBeenCalledWith('invalid input')
    })

    it('403 calls warning() with permission message', async () => {
        await loadRequest()
        const error = {
            response: { status: 403, data: {}, config: { url: '/api/v1/host/info' } },
        }
        void responseErrHandlers[0](error)
        expect(warning).toHaveBeenCalledWith('您目前没有权限执行该操作，请联系管理员')
    })

    it('500 calls warning() with server error message', async () => {
        await loadRequest()
        const error = {
            response: { status: 500, data: {}, config: { url: '/api/v1/host/info' } },
        }
        void responseErrHandlers[0](error)
        expect(warning).toHaveBeenCalledWith('服务器错误，请稍后再试')
    })

    it('401 on non-token_update schedules a token refresh', async () => {
        vi.useFakeTimers()
        stubStore.user.refresh = 'old-refresh'
        stubStore.user.token = 'old'

        await loadRequest()
        const error = {
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/host/info' },
            },
        }
        void responseErrHandlers[0](error)

        // The implementation uses setTimeout(..., 500); advance the fake timer.
        vi.advanceTimersByTime(1000)

        expect(latestFakeInstance.post).toHaveBeenCalledWith(
            '/api/v1/auth/token_update',
            {},
            expect.objectContaining({
                headers: expect.objectContaining({ Authorization: 'Bearer old-refresh' }),
            }),
        )
    })
})
