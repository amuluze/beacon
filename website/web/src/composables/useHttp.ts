import type { FetchResponse } from 'ofetch'

export interface RequestOptions<T> {
    data: T
    code: number
    message: string
    success: boolean
}

function handleError<T>(response: FetchResponse<RequestOptions<T>> & FetchResponse<ResponseType>) {
    const err = (text: string) => {
        ElMessage.error({ message: response?._data?.message ?? text })
    }
    if (!response._data) {
        err('请求超时，服务器无响应！')
        return
    }
    // const store = useStore()
    const handleMap: { [key: number]: () => void } = {
        404: () => err('服务器资源不存在'),
        500: () => err('服务器内部错误'),
        403: () => err('没有权限访问该资源'),
        401: () => {
            err('登录状态已过期，需要重新登录')
            // store.user.setToken('', '')
            // 跳转实际登录页
            navigateTo('/')
        },
    }
    const handler = handleMap[response.status]
    if (handler)
        handler()
    else
        err('未知错误！')
}

const fetch = $fetch.create({
    // 请求拦截器
    onRequest({ options }) {
        const { public: { baseUrl } } = useRuntimeConfig()
        // 仅在显式配置时设置 baseURL；留空则走相对路径，
        // 由同源 nitro routeRules 代理到后端，避免 SPA 中 localhost 误指向用户本机。
        if (baseUrl) {
            options.baseURL = baseUrl as string
        }
    },
    onResponse({ response }) {
        const contentType = response.headers.get('Content-Type')
        if (!response.ok) {
            handleError(response)
            return Promise.resolve(response._data)
        }
        if (!contentType) {
            response._data = { code: -1, data: '返回数据不符合预期' }
            return
        }

        if (contentType === 'application/json') {
            response._data.headers = response.headers
        }
        else {
            const disposition = response.headers.get('content-disposition')
            if (!disposition) {
                response._data = { code: -2, data: '返回数据不符合预期' }
                return
            }
            // 切割文件名
            const blob = new Blob([response._data], { type: contentType })
            const blobURL = URL.createObjectURL(blob)
            response._data = { code: 1, data: blobURL, headers: response.headers }
        }
    },
    onResponseError({ response }) {
        handleError(response)
        return Promise.resolve(response?._data ?? null)
    },
})

// 自动导出
export const useHttp = {
    get: <T>(url: string, params?: any) => {
        return fetch<RequestOptions<T>>(url, { method: 'get', params })
    },
    post: <T>(url: string, body?: any) => {
        return fetch<RequestOptions<T>>(url, { method: 'post', body })
    },
}
