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
interface StubStore {
    user: { token: string, refresh: string, setToken: ReturnType<typeof vi.fn> }
    agent: { selectedAgentID: string, clear: ReturnType<typeof vi.fn> }
}
const stubStore: StubStore = {
    user: { token: '', refresh: '', setToken: vi.fn() },
    agent: { selectedAgentID: '', clear: vi.fn() },
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
    const inst: any = vi.fn().mockResolvedValue({ data: {} })
    inst.interceptors = {
        request: { use: (ok: ReqHandler) => { requestHandlers.push(ok) } },
        response: {
            use: (_ok: any, err: ResErrHandler) => { responseErrHandlers.push(err) },
        },
    }
    inst.get = vi.fn()
    inst.post = vi.fn().mockResolvedValue({ data: { access_token: 'new', refresh_token: 'new-r' } })
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
    stubStore.agent.clear = vi.fn()
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

    it('injects X-Agent-ID for audit queries', async () => {
        stubStore.agent.selectedAgentID = 'agent-audit'

        await loadRequest()
        const out = await requestHandlers[0]({ headers: {}, url: '/api/v1/audit/query?page=1&size=10' })

        expect(out.headers['X-Agent-ID']).toBe('agent-audit')
    })

    it('injects Authorization from stubStore.user.token', async () => {
        stubStore.user.token = 'abc123'
        stubStore.agent.selectedAgentID = ''

        await loadRequest()
        const out = await requestHandlers[0]({ headers: {}, url: '/api/v1/host/info' })
        expect(out.headers.Authorization).toBe('Bearer abc123')
    })

    it('preserves the refresh token Authorization on token_update', async () => {
        stubStore.user.token = 'expired-access'

        await loadRequest()
        const out = await requestHandlers[0]({
            headers: { Authorization: 'Bearer valid-refresh' },
            url: '/api/v1/auth/token_update',
        })

        expect(out.headers.Authorization).toBe('Bearer valid-refresh')
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
    it('token_update failures are rejected to the refresh owner without a warning', async () => {
        await loadRequest()
        expect(responseErrHandlers.length).toBe(1)

        const error = {
            response: {
                status: 400,
                data: { msg: 'invalid' },
                config: { url: '/api/v1/auth/token_update' },
            },
        }
        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).not.toHaveBeenCalled()
        expect(stubStore.user.setToken).not.toHaveBeenCalled()
        expect(stubStore.agent.clear).not.toHaveBeenCalled()
        expect(lastHref).toBeNull()
    })

    it('network error without response warns and rejects without logout', async () => {
        await loadRequest()
        // 模拟无 HTTP 响应（断网/超时/CORS），error.response 为 undefined
        await expect(responseErrHandlers[0](new Error('Network Error'))).rejects.toBeTruthy()
        expect(warning).toHaveBeenCalledWith('网络异常，请检查连接后重试')
        expect(stubStore.user.setToken).not.toHaveBeenCalled()
        expect(lastHref).toBeNull()
    })

    it('unknown status code warns instead of silent redirect', async () => {
        await loadRequest()
        const error = {
            response: {
                status: 502,
                data: { msg: 'bad gateway' },
                config: { url: '/api/v1/host/info' },
            },
        }
        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).toHaveBeenCalledWith('bad gateway')
        // 未清 token、未跳转，避免不清状态卡在首页
        expect(stubStore.user.setToken).not.toHaveBeenCalled()
        expect(lastHref).toBeNull()
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
        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).toHaveBeenCalledWith('invalid input')
    })

    it('400 prefers the service error detail for password failures', async () => {
        await loadRequest()
        const error = {
            response: {
                status: 400,
                data: { err: 'invalid password', msg: 'bad request' },
                config: { url: '/api/v1/auth/pass_update' },
            },
        }

        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).toHaveBeenCalledWith('invalid password')
    })

    it('403 calls warning() with permission message', async () => {
        await loadRequest()
        const error = {
            response: { status: 403, data: {}, config: { url: '/api/v1/host/info' } },
        }
        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).toHaveBeenCalledWith('您目前没有权限执行该操作，请联系管理员')
    })

    it('500 calls warning() with server error message', async () => {
        await loadRequest()
        const error = {
            response: { status: 500, data: {}, config: { url: '/api/v1/host/info' } },
        }
        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).toHaveBeenCalledWith('服务器错误，请稍后再试')
    })

    it('401 on non-token_update refreshes the token and resolves the retried request', async () => {
        stubStore.user.refresh = 'old-refresh'
        stubStore.user.token = 'old'

        await loadRequest()
        const retriedResponse = { data: { hostname: 'node-01' } }
        latestFakeInstance.mockResolvedValueOnce(retriedResponse)
        const error = {
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/host/info' },
            },
        }
        const request = responseErrHandlers[0](error)

        await expect(request).resolves.toBe(retriedResponse)
        expect(latestFakeInstance.post).toHaveBeenCalledWith(
            '/api/v1/auth/token_update',
            {},
            expect.objectContaining({
                headers: expect.objectContaining({ Authorization: 'Bearer old-refresh' }),
            }),
        )
        expect(latestFakeInstance).toHaveBeenCalledWith(error.response.config)
    })

    it('rejects every queued request when refresh fails', async () => {
        stubStore.user.refresh = 'expired-refresh'
        stubStore.user.token = 'expired-access'

        await loadRequest()
        const refreshError = {
            response: {
                status: 500,
                data: { err: 'key not found' },
                config: { url: '/api/v1/auth/token_update' },
            },
        }
        latestFakeInstance.post.mockImplementationOnce(async () => responseErrHandlers[0](refreshError))
        const firstError = {
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/agent/list' },
            },
        }
        const secondError = {
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/host/info' },
            },
        }

        const firstRequest = responseErrHandlers[0](firstError)
        const secondRequest = responseErrHandlers[0](secondError)

        await expect(firstRequest).rejects.toBe(refreshError.response)
        await expect(secondRequest).rejects.toBe(refreshError.response)
        expect(latestFakeInstance.post).toHaveBeenCalledTimes(1)
        expect(stubStore.user.setToken).toHaveBeenCalledWith('', '')
        expect(stubStore.agent.clear).toHaveBeenCalledOnce()
        expect(lastHref).toBe('/#/login')
        expect(warning).not.toHaveBeenCalled()
    })

    it('allows a new Agent request to settle after refresh failure and re-login', async () => {
        stubStore.user.refresh = 'expired-refresh'
        stubStore.user.token = 'expired-access'

        await loadRequest()
        const refreshError = {
            response: {
                status: 500,
                data: { err: 'key not found' },
                config: { url: '/api/v1/auth/token_update' },
            },
        }
        latestFakeInstance.post.mockImplementationOnce(async () => responseErrHandlers[0](refreshError))
        const expiredRequest = responseErrHandlers[0]({
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/agent/list' },
            },
        })
        await expect(expiredRequest).rejects.toBe(refreshError.response)

        stubStore.user.refresh = 'fresh-refresh'
        stubStore.user.token = 'fresh-access'
        latestFakeInstance.post.mockResolvedValueOnce({
            data: { access_token: 'renewed-access', refresh_token: 'renewed-refresh' },
        })
        const agentResponse = { data: [{ agent_id: 'node-01' }] }
        latestFakeInstance.mockResolvedValueOnce(agentResponse)

        const newAgentRequest = responseErrHandlers[0]({
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/agent/list' },
            },
        })

        await expect(newAgentRequest).resolves.toBe(agentResponse)
        expect(latestFakeInstance.post).toHaveBeenCalledTimes(2)
        expect(stubStore.user.setToken).toHaveBeenLastCalledWith('renewed-access', 'renewed-refresh')
    })
})
