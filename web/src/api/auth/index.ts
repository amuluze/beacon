/**
 * @Author     : Amu
 * @Date       : 2024/03/28 00:15:21
 * @Description:
 */

import request from '@/api'

import { loginFormData, LoginResponseData, UpdatePassword } from '@/interface/auth.ts'

/**
 * 登录
 * @param data
 */
export function login(data: loginFormData) {
    return request.post<LoginResponseData>('/api/v1/auth/login', data)
}

/**
 * 更新 token
 */
export function updateToken() {
    return request.post<LoginResponseData>('/api/v1/auth/token_update')
}

/**
 * 登出
 */
export function logout() {
    return request.post('/api/v1/auth/logout')
}

/**
 * 更新密码
 */
export function updatePassword(params: UpdatePassword) {
    return request.post('/api/v1/auth/pass_update', params)
}
