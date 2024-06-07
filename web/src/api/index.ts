/**
 * @Author     : Amu
 * @Date       : 2024/1/8 14:18
 * @Description:
 */
import { warning } from '@/components/Message/message.ts'
import { ResultData } from '@/interface/result.ts'
import useStore from '@/store'
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'

const config = {
    // 默认地址请求地址，可在 .env.*** 文件中修改
    baseURL: '/app/',
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
    isRefreshing: boolean
    requestQueue: Function[]

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
                if (store.user.token !== '') {
                    config.headers.Authorization = 'Bearer ' + store.user.token
                }
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

        /**
         * 响应拦截器
         */
        this.service.interceptors.response.use(
            (response: AxiosResponse) => {
                return response
            },
            async (error: AxiosError) => {
                // https://zhuanlan.zhihu.com/p/653595209
                const originalRequest: InternalAxiosRequestConfig<any> = error.config as InternalAxiosRequestConfig<any>
                const store = useStore()
                if (error.response?.status === 400) {
                    warning('请检查您的请求参数')
                }
                if (error.response?.status === 401) {
                    if (!this.isRefreshing) {
                        this.isRefreshing = true

                        // 发起刷新请求
                        originalRequest.headers.Authorization = store.user.refresh
                        this.service
                            .post('/api/v1/auth/token_update', originalRequest)
                            .then((res) => {
                                store.user.setToken(res.data.access_token, res.data.refresh_token)
                                // 刷新 token 完成后，重启发送之前失败的请求
                                this.requestQueue.forEach((request) => request(store.user.token))

                                originalRequest.headers.Authorization = `Bearer ${store.user.token}`
                                return this.service(originalRequest)
                            })
                            .catch((err) => {
                                console.log('update token error: ', err)
                                // 刷新 token 失败，跳转到登录页面
                                store.user.setToken('', '')
                            })
                            .finally(() => {
                                this.isRefreshing = false
                            })
                        this.requestQueue = []
                        this.isRefreshing = false
                        window.location.href = '/login'
                    } else {
                        // 正在刷新 token，将当前请求加入队列，等待刷新完成后再重新发送
                        return new Promise((resolve) => {
                            this.requestQueue.push((token: string) => {
                                originalRequest.headers.Authorization = `Bearer ${token}`
                                resolve(this.service(originalRequest as InternalAxiosRequestConfig))
                            })
                        })
                    }
                }
                if (error.response?.status === 403) {
                    warning('您目前没有权限执行该操作，请联系管理员')
                }
                if (error.response?.status === 500) {
                    console.log('500 error: ', error.response.data)
                    warning('服务器错误，请稍后再试')
                }
                return Promise.reject(error)
            }
        )
    }

    /**
     * 常用请求方法封装
     */
    get<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.get(url, { params: params, ...config })
    }

    post<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.post(url, params, config)
    }
}

export default new Request(config)
