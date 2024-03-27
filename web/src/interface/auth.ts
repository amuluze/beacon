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
    token: string
    refresh: string
}

export interface UpdatePassword {
    username: string
    old_password: string
    net_password: string
}
