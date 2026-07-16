import type { FetchResponse } from 'ofetch'

// 官网后端直接返回业务对象（如 { id, times }），不使用统一信封。
// 因此 useHttp 的泛型 T 即为响应体类型，调用方直接读取业务字段。

function handleError(response: FetchResponse<any>) {
    // ElMessage 仅在浏览器上下文可用；SSR 阶段跳过，避免访问 window/document 崩溃
    if (import.meta.server)
        return
    const err = (text: string) => {
        ElMessage.error({ message: response?._data?.msg ?? text })
    }
    if (!response._data) {
        err('请求超时，服务器无响应！')
        return
    }
    const handleMap: Record<number, () => void> = {
        400: () => err(response._data?.msg ?? '请求参数有误'),
        404: () => err('服务器资源不存在'),
        403: () => err('没有权限访问该资源'),
        500: () => err('服务器内部错误'),
    }
    const handler = handleMap[response.status]
    if (handler)
        handler()
    else
        err('请求失败，请稍后再试')
}

const fetch = $fetch.create({
    onRequest({ options }) {
        const { public: { baseUrl } } = useRuntimeConfig()
        // 仅在显式配置时设置 baseURL；留空则走相对路径，
        // 由同源 nitro routeRules 代理到后端，避免 SPA 中 localhost 误指向用户本机。
        if (baseUrl)
            options.baseURL = baseUrl as string
    },
    // 成功响应直接透传业务对象，不在 _data 上附加 headers 等运行时字段（保持类型与运行时一致）；
    // 错误响应统一交给 onResponseError 处理，避免重复弹窗。
    onResponseError({ response }) {
        handleError(response)
        return Promise.reject(response?._data ?? null)
    },
})

// 自动导出
export const useHttp = {
    get: <T>(url: string, params?: Record<string, unknown>) => {
        return fetch<T>(url, { method: 'get', params })
    },
    post: <T>(url: string, body?: BodyInit | Record<string, any> | null) => {
        return fetch<T>(url, { method: 'post', body })
    },
}
