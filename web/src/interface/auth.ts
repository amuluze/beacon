/**
 * @Author     : Amu
 * @Date       : 2024/03/28 00:17:07
 * @Description:
 */

export interface loginFormData {
    username: string
    password: string
}

export interface LoginResponseData {
    access_token: string
    refresh_token: string
}

export interface UpdatePassword {
    username: string
    old_password: string
    new_password: string
}

export interface UserInfo {
    id: string
    username: string
    status: string
    is_admin: number
}
