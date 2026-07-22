/**
 * @Author     : Amu
 * @Date       : 2024/1/8 14:18
 * @Description:
 */
/**
 * @Author     : Amu
 * @Date       : 2024/10/12 10:10
 * @Description:
 */

import { warning } from '@/components/Message/message.ts'

import type { ResultData } from '@/interface/result.ts'
import type { AxiosError, AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'

import useStore from '@/store'
import axios from 'axios'

const agentScopedPrefixes = [
    '/api/v1/host/',
    '/api/v1/container/',
    '/api/v1/audit/',
]

const agentScopedExcludes = [
    '/api/v1/host/install',
    '/api/v1/host/report',
    '/api/v1/host/get_install_token',
]

function requestPath(url?: string): string {
    if (!url)
        return ''
    try {
        return new URL(url, window.location.origin).pathname
    }
    catch {
        return url.split('?')[0]
    }
}

function requiresAgent(url?: string): boolean {
    const path = requestPath(url)
    if (agentScopedExcludes.some(prefix => path.startsWith(prefix)))
        return false
    return agentScopedPrefixes.some(prefix => path.startsWith(prefix))
}

const config = {
    // 默认地址请求地址，可在 .env.*** 文件中修改
    baseURL: '/',
    // 设置超时时间
    timeout: 600000,
    // 设置默认请求头
    headers: { 'Content-Type': 'application/json;charset=utf-8' },
    // 跨域时候允许携带凭证
    withCredentials: true,
}

interface QueuedRequest {
    retry: () => Promise<unknown>
    resolve: (value: unknown) => void
    reject: (reason?: unknown) => void
}

class Request {
    service: AxiosInstance
    isRefreshing: boolean

    requestQueue: QueuedRequest[]

    /** 存储因 token 过期而导致发送失败的请求 */
    private saveErrorRequest = (expiredRequest: QueuedRequest): void => {
        this.requestQueue.push(expiredRequest)
    }

    /** 清理当前存储的过期请求 */
    private clearExpiredRequest = (): void => {
        this.requestQueue = []
    }

    /** 执行当前存储的由于过期导致失败的请求 */
    private againRequest = (): void => {
        const expiredRequests = this.requestQueue
        this.clearExpiredRequest()
        expiredRequests.forEach((request): void => {
            Promise.resolve()
                .then(request.retry)
                .then(request.resolve, request.reject)
        })
    }

    /** refresh 失败时拒绝全部排队请求，确保调用方 finally 能够执行 */
    private rejectExpiredRequest = (reason: unknown): void => {
        const expiredRequests = this.requestQueue
        this.clearExpiredRequest()
        expiredRequests.forEach((request): void => {
            request.reject(reason)
        })
    }

    /** 强制登出：清 token 与排队请求，整页跳转到登录页（hash 路由） */
    private forceLogout = (reason: unknown = new Error('登录状态已失效')): void => {
        const store = useStore()
        store.user.setToken('', '')
        store.agent.clear()
        this.rejectExpiredRequest(reason)
        this.isRefreshing = false
        window.location.href = '/#/login'
    }

    /** 利用 refreshToken 更新 accessToken */
    private updateAccessTokenByRefreshToken = (): void => {
        const store = useStore()
        this.service.post('/api/v1/auth/token_update', {}, {
            headers: { Authorization: `Bearer ${store.user.refresh}` },
        }).then((res) => {
            // 更新本地 token
            store.user.setToken(res.data.access_token, res.data.refresh_token)
            // 更新 token 后，重放之前失败的请求
            this.againRequest()
        }).catch((error) => {
            // 此时 refreshToken 也失效了（或刷新请求本身异常），强制登出返回登录页
            this.forceLogout(error)
        }).finally(() => {
            this.isRefreshing = false
        })
    }

    private refreshToken = (expiredRequest: QueuedRequest): void => {
        this.saveErrorRequest(expiredRequest)
        if (this.isRefreshing) {
            // 已有刷新进行中，当前请求排队等待完成后由 againRequest 重放
            return
        }
        this.isRefreshing = true
        this.updateAccessTokenByRefreshToken()
    }

    public constructor(config: AxiosRequestConfig) {
        // instantiation
        this.service = axios.create(config)
        this.isRefreshing = false
        this.requestQueue = []

        /**
         * 请求拦截器
         * 客户端发送请求 -> [请求拦截器] -> 服务器
         * token 校验(JWT)：接收服务器返回的 token，存储到 pinia 本地存储当中
         */
        this.service.interceptors.request.use(
            async (config: InternalAxiosRequestConfig) => {
                const store = useStore()
                if (store.user.token !== '' && config.url !== '/api/v1/auth/token_update') {
                    config.headers.Authorization = `Bearer ${store.user.token}`
                }
                if (requiresAgent(config.url) && store.agent.selectedAgentID) {
                    config.headers['X-Agent-ID'] = store.agent.selectedAgentID
                }
                if (config.url?.endsWith('login')) {
                    config.headers['Content-Type'] = 'multipart/form-data;charset=UTF-8'
                    return config
                }
                return config
            },
            async (error: AxiosError) => {
                return Promise.reject(error)
            },
        )

        /**
         * 响应拦截器
         */
        this.service.interceptors.response.use(
            (response) => {
                if (response.headers['content-disposition']) {
                    const downLoadMark = response.headers['content-disposition'].split(';')
                    if (downLoadMark[0] === 'attachment') {
                        // 执行下载
                        let fileName = downLoadMark[1].split('filename=')[1]
                        if (fileName) {
                            fileName = decodeURI(fileName)
                            const content = response.data
                            const url = window.URL.createObjectURL(new Blob([content], { type: 'application/octet-stream' }))
                            const link = document.createElement('a')
                            link.style.display = 'none'
                            link.href = url
                            link.download = fileName
                            document.body.appendChild(link)
                            link.click()
                            link.remove()
                            window.URL.revokeObjectURL(url)
                        }
                        else {
                            return response
                        }
                    }
                }
                return response
            },
            async (error) => {
                if (!error.response) {
                    // 网络错误/超时：无 HTTP 响应，提示后拒绝，不登出
                    warning('网络异常，请检查连接后重试')
                    return Promise.reject(error)
                }
                const { data, config, status } = error.response
                if (config.url === '/api/v1/auth/token_update') {
                    // 刷新请求由 updateAccessTokenByRefreshToken 统一处理失败，
                    // 避免响应拦截器重复强退或展示无意义的 token 服务错误。
                    return Promise.reject(error.response)
                }
                if (status === 400) {
                    warning(data.err || data.msg)
                }
                else if (status === 403) {
                    warning('您目前没有权限执行该操作，请联系管理员')
                }
                else if (status === 500) {
                    warning(data.err || data.msg || '服务器错误，请稍后再试')
                }
                else if (status === 401 && config.url !== '/api/v1/auth/token_update') {
                    return new Promise((resolve, reject) => {
                        this.refreshToken({
                            retry: async () => this.service(config),
                            resolve,
                            reject,
                        })
                    })
                }
                else {
                    // 未识别状态码：提示而非静默跳转，避免不清 token 卡在首页
                    warning(data?.err || data?.msg || `请求失败（${status}）`)
                }
                return Promise.reject(error.response)
            },
        )
    }

    /**
     * 常用请求方法封装
     */
    async get<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.get(url, { params, ...config })
    }

    async post<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.post(url, params, config)
    }
}

const request = new Request(config)
export default request
