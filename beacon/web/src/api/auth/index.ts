/**
 * @Author     : Amu
 * @Date       : 2024/03/28 00:15:21
 * @Description:
 */

import request from '@/api'

import type { LoginForm, LoginResponse, UpdatePassword, UserInfo } from '@/interface/auth.ts'

/**
 * 登录
 * @param data
 */
export async function login(data: LoginForm) {
    return request.post<LoginResponse>('/api/v1/auth/login', data)
}

/**
 * 登出
 */
export async function logout() {
    return request.post<any>('/api/v1/auth/logout', {})
}

/**
 * 更新密码
 */
export async function updatePassword(data: UpdatePassword) {
    return request.post<any>('/api/v1/auth/pass_update', data)
}

/**
 * 获取用户信息
 */
export async function getUserInfo() {
    return request.get<UserInfo>('/api/v1/auth/user_info')
}
