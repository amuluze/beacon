/**
 * @Author     : Amu
 * @Date       : 2024/03/28 00:15:21
 * @Description:
 */

import request from '@/api'

import { loginFormData, LoginResponseData } from '@/interface/auth.ts'

/**
 * 登录
 * @param data
 */
export function login(data: loginFormData) {
    return request.post<LoginResponseData>('/api/v1/auth/login', data)
}

/**
 * 登出
 */
export function logout() {
    return request.post('/api/v1/auth/logout')
}
