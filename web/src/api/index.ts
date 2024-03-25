/**
 * @Author     : Amu
 * @Date       : 2024/1/8 14:18
 * @Description:
 */
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'
import { ResultData } from '@/interface/result.ts'

const config = {
    // 默认地址请求地址，可在 .env.*** 文件中修改
    baseURL: '/api',
    // baseURL: import.meta.env.VITE_API_URL,
    // 设置超时时间
    timeout: 30000,
    // 设置默认请求头
    headers: { 'Content-Type': 'application/json;charset=utf-8' },
    // 跨域时候允许携带凭证
    withCredentials: true
}

class Request {
    service: AxiosInstance

    public constructor(config: AxiosRequestConfig) {
        // instantiation
        this.service = axios.create(config)

        /**
         * 请求拦截器
         * 客户端发送请求 -> [请求拦截器] -> 服务器
         * token 校验(JWT)：接收服务器返回的 token，存储到 pinia 本地存储当中
         */
        this.service.interceptors.request.use(
            async (config: InternalAxiosRequestConfig) => {
                if (config.url?.endsWith('login')) {
                    config.headers['Content-Type'] = 'multipart/form-data;charset=UTF-8'
                    return config
                }
                return config
            },
            (error: AxiosError) => {
                return Promise.reject(error)
            }
        )
    }

    /**
     * 常用请求方法封装
     */
    get<T>(url: string, params?: object | string): Promise<ResultData<T>> {
        return this.service.get(url, { params, ...config })
    }

    post<T>(url: string, params?: object | string): Promise<ResultData<T>> {
        return this.service.post(url, params, config)
    }
}

export default new Request(config)
