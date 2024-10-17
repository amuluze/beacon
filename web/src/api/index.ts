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
import useStore from '@/store'
import type { AxiosError, AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'
import axios from 'axios'

const config = {
    // 默认地址请求地址，可在 .env.*** 文件中修改
    baseURL: '/app/',
    // baseURL: import.meta.env.VITE_API_URL,
    // 设置超时时间
    timeout: 600000,
    // 设置默认请求头
    headers: { 'Content-Type': 'application/json;charset=utf-8' },
    // 跨域时候允许携带凭证
    withCredentials: true,
}

class Request {
    service: AxiosInstance
    isRefreshing: boolean

    requestQueue: Function[]
    /** 存储因 token 过期而导致发送失败的请求 */
    private saveErrorRequest = (expiredRequest: () => any) => {
        this.requestQueue.push(expiredRequest)
    }

    /** 清理当前存储的过期请求 */
    private clearExpiredRequest = () => {
        this.requestQueue = []
    }

    /** 执行当前存储的由于过期导致失败的请求 */
    private againRequest = () => {
        this.requestQueue.forEach((request) => {
            request()
        })
        this.clearExpiredRequest()
    }

    /** 避免频繁发送更新 */
    firstRequest: boolean
    /** 利用 refreshToken 更新 accessToken */
    private updateAccessTokenByRefreshToken = () => {
        this.firstRequest = false
        const store = useStore()
        this.service.post('/api/v1/auth/token_update', {}, {
            headers: { Authorization: `Bearer ${store.user.refresh}` },
        }).then((res) => {
            // 更新本地 token
            store.user.setToken(res.data.access_token, res.data.refresh_token)
            // 更新 token 后，重启发起之前失败的请求
            this.againRequest()
        }).catch(() => {
            // 此时 refreshToken 也失效了，返回登录页
            window.location.href = '/app/app/login'
        })
    }

    private refreshToken = (expiredRequest: () => any) => {
        this.saveErrorRequest(expiredRequest)
        // 保证再发起更新时，已经没有过期请求要进行存储了

        setTimeout(() => {
            this.updateAccessTokenByRefreshToken()
        }, 500)
    }

    public constructor(config: AxiosRequestConfig) {
        // instantiation
        this.service = axios.create(config)
        this.isRefreshing = false
        this.firstRequest = false
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
                    config.headers.Authorization = `Bearer ${store.user.token}`
                }
                if (config.url?.endsWith('login')) {
                    config.headers['Content-Type'] = 'multipart/form-data;charset=UTF-8'
                    return config
                }
                return config
            },
            (error: AxiosError) => {
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
                            // fileName = decodeURIComponent(filename);//对filename进行转码
                            fileName = decodeURI(fileName)
                            const content = response.data
                            const url = window.URL.createObjectURL(new Blob([content], { type: 'application/octet-stream' }))
                            const link = document.createElement('a')
                            link.style.display = 'none'
                            link.href = url
                            link.download = fileName
                            // link.setAttribute('download', fileName)
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
                const { data, config, status } = error.response
                return new Promise((resolve, reject) => {
                    /** 判断当前请求失败的原因 */
                    if (status === 400) {
                        warning(data.msg)
                    }
                    else if (status === 403) {
                        warning('您目前没有权限执行该操作，请联系管理员')
                    }
                    else if (status === 500) {
                        warning('服务器错误，请稍后再试')
                    }
                    else if (status === 401 && config.url === '/api/v1/auth/token_update') {
                        // refreshToken 也失效了，返回登录页
                        this.clearExpiredRequest()
                        window.location.href = '/app/app/login'
                    }
                    else if (status === 401 && config.url !== '/api/v1/auth/token_update') {
                        this.refreshToken(() => {
                            resolve(this.service(config))
                        })
                    }
                    else {
                        window.location.href = '/app/app/login'
                        reject(error.response)
                    }
                })
            },
        )
    }

    /**
     * 常用请求方法封装
     */
    get<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.get(url, { params, ...config })
    }

    post<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.post(url, params, config)
    }
}

export default new Request(config)
