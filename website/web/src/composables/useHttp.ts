import type { FetchResponse } from 'ofetch'

import { showToast } from '~/composables/useToast'

// useHttp 的泛型 T 对应接口的真实响应体；是否存在 data 信封由各 API 边界显式处理，
// 避免在通用请求层猜测并改写后端契约。

function handleError(response: FetchResponse<any>) {
    // Toast 仅在浏览器上下文展示；SSR 阶段跳过，避免跨请求共享用户可见状态。
    if (import.meta.server)
        return
    const err = (text: string) => showToast(response?._data?.msg ?? text, 'error')
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

export interface HttpRequestOptions {
    silent?: boolean
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
    onResponseError({ response, options }) {
        if (!(options as HttpRequestOptions).silent)
            handleError(response)
        return Promise.reject(response?._data ?? null)
    },
})

type CreatedFetchOptions = NonNullable<Parameters<typeof fetch>[1]>

// 自动导出
export const useHttp = {
    get: <T>(url: string, params?: Record<string, unknown>, requestOptions: HttpRequestOptions = {}) => {
        const options = { method: 'get' as const, params, ...requestOptions } as CreatedFetchOptions & HttpRequestOptions
        return fetch<T>(url, options)
    },
    post: <T>(url: string, body?: BodyInit | Record<string, any> | null, requestOptions: HttpRequestOptions = {}) => {
        const options = { method: 'post' as const, body, ...requestOptions } as CreatedFetchOptions & HttpRequestOptions
        return fetch<T>(url, options)
    },
}
