import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

interface StubStore {
    user: { clearSession: ReturnType<typeof vi.fn> }
    agent: { selectedAgentID: string, clear: ReturnType<typeof vi.fn> }
}

const stubStore: StubStore = {
    user: { clearSession: vi.fn() },
    agent: { selectedAgentID: '', clear: vi.fn() },
}

vi.mock('@/store', () => ({
    default: () => stubStore,
}))

type ReqHandler = (cfg: any) => any
type ResErrHandler = (err: any) => any
const requestHandlers: ReqHandler[] = []
const responseErrHandlers: ResErrHandler[] = []

let latestFakeInstance: any = null
function makeFakeAxiosInstance() {
    const instance: any = vi.fn().mockResolvedValue({ data: {} })
    instance.interceptors = {
        request: { use: (ok: ReqHandler) => requestHandlers.push(ok) },
        response: { use: (_ok: any, err: ResErrHandler) => responseErrHandlers.push(err) },
    }
    instance.get = vi.fn()
    instance.post = vi.fn().mockResolvedValue({ status: 204 })
    latestFakeInstance = instance
    return instance
}

vi.mock('axios', () => ({
    default: { create: () => makeFakeAxiosInstance() },
}))

const warning = vi.fn()
vi.mock('@/components/Message/message.ts', () => ({
    warning: (...args: any[]) => warning(...args),
}))

let lastHref: string | null = null
Object.defineProperty(window, 'location', {
    configurable: true,
    get: () => ({
        set href(value: string) { lastHref = value },
        get href() { return lastHref ?? '' },
        origin: 'https://beacon.example.com',
    }) as any,
})

beforeEach(() => {
    stubStore.user.clearSession = vi.fn()
    stubStore.agent.selectedAgentID = ''
    stubStore.agent.clear = vi.fn()
    requestHandlers.length = 0
    responseErrHandlers.length = 0
    lastHref = null
    warning.mockClear()
    document.cookie = 'beacon_csrf_token=; Max-Age=0; path=/'
    vi.clearAllMocks()
    vi.resetModules()
})

afterEach(() => {
    vi.useRealTimers()
})

async function loadRequest() {
    return (await import('@/api/index')).default
}

describe('api/index request interceptor', () => {
    it('injects Agent ID without exposing an Authorization token', async () => {
        stubStore.agent.selectedAgentID = 'agent-1'
        await loadRequest()

        const output = await requestHandlers[0]({ headers: {}, method: 'get', url: '/api/v1/host/info' })
        expect(output.headers['X-Agent-ID']).toBe('agent-1')
        expect(output.headers.Authorization).toBeUndefined()
    })

    it('adds the double-submit CSRF header to unsafe requests', async () => {
        document.cookie = 'beacon_csrf_token=csrf-value; path=/'
        await loadRequest()

        const output = await requestHandlers[0]({ headers: {}, method: 'post', url: '/api/v1/dingtalk/update' })
        expect(output.headers['X-CSRF-Token']).toBe('csrf-value')
    })

    it('does not add CSRF headers to safe requests', async () => {
        document.cookie = 'beacon_csrf_token=csrf-value; path=/'
        await loadRequest()

        const output = await requestHandlers[0]({ headers: {}, method: 'get', url: '/api/v1/auth/user_info' })
        expect(output.headers['X-CSRF-Token']).toBeUndefined()
    })

    it('keeps the login multipart content type without Agent or Authorization headers', async () => {
        await loadRequest()

        const output = await requestHandlers[0]({ headers: {}, method: 'post', url: '/api/v1/auth/login' })
        expect(output.headers['Content-Type']).toContain('multipart/form-data')
        expect(output.headers['X-Agent-ID']).toBeUndefined()
        expect(output.headers.Authorization).toBeUndefined()
    })
})

describe('api/index response interceptor', () => {
    it('returns refresh failures to the refresh owner without duplicate warnings', async () => {
        await loadRequest()
        const error = {
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/auth/token_update' },
            },
        }

        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).not.toHaveBeenCalled()
        expect(stubStore.user.clearSession).not.toHaveBeenCalled()
    })

    it('warns on a network error without clearing the browser session', async () => {
        await loadRequest()

        await expect(responseErrHandlers[0](new Error('Network Error'))).rejects.toBeTruthy()
        expect(warning).toHaveBeenCalledWith('网络异常，请检查连接后重试')
        expect(stubStore.user.clearSession).not.toHaveBeenCalled()
    })

    it('does not try to refresh invalid login credentials', async () => {
        await loadRequest()
        const error = {
            response: {
                status: 401,
                data: { err: 'invalid username or password' },
                config: { url: '/api/v1/auth/login' },
            },
        }

        await expect(responseErrHandlers[0](error)).rejects.toEqual(error.response)
        expect(warning).toHaveBeenCalledWith('invalid username or password')
        expect(latestFakeInstance.post).not.toHaveBeenCalled()
    })

    it('refreshes with HttpOnly Cookies and replays the failed request', async () => {
        await loadRequest()
        const retriedResponse = { data: { hostname: 'node-01' } }
        latestFakeInstance.mockResolvedValueOnce(retriedResponse)
        const error = {
            response: {
                status: 401,
                data: {},
                config: { headers: {}, method: 'get', url: '/api/v1/host/info' },
            },
        }

        await expect(responseErrHandlers[0](error)).resolves.toBe(retriedResponse)
        expect(latestFakeInstance.post).toHaveBeenCalledWith('/api/v1/auth/token_update', {})
        expect(latestFakeInstance).toHaveBeenCalledWith(expect.objectContaining({ _retry: true }))
    })

    it('rejects all queued requests and clears user state when Cookie refresh fails', async () => {
        await loadRequest()
        const refreshError = {
            response: {
                status: 401,
                data: {},
                config: { url: '/api/v1/auth/token_update' },
            },
        }
        latestFakeInstance.post.mockImplementationOnce(async () => responseErrHandlers[0](refreshError))
        const first = responseErrHandlers[0]({
            response: {
                status: 401,
                data: {},
                config: { headers: {}, method: 'get', url: '/api/v1/agent/list' },
            },
        })
        const second = responseErrHandlers[0]({
            response: {
                status: 401,
                data: {},
                config: { headers: {}, method: 'get', url: '/api/v1/host/info' },
            },
        })

        await expect(first).rejects.toBe(refreshError.response)
        await expect(second).rejects.toBe(refreshError.response)
        expect(latestFakeInstance.post).toHaveBeenCalledTimes(1)
        expect(stubStore.user.clearSession).toHaveBeenCalledOnce()
        expect(stubStore.agent.clear).toHaveBeenCalledOnce()
        expect(lastHref).toBe('/#/login')
    })

    it('fails closed if the retried request is still unauthorized', async () => {
        await loadRequest()
        const response = {
            status: 401,
            data: {},
            config: { _retry: true, headers: {}, method: 'get', url: '/api/v1/host/info' },
        }

        await expect(responseErrHandlers[0]({ response })).rejects.toEqual(response)
        expect(stubStore.user.clearSession).toHaveBeenCalledOnce()
        expect(lastHref).toBe('/#/login')
    })
})
