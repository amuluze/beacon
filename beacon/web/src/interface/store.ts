export interface UserInfo {
    name: string
    status: number
    isAdmin: number
    mustChangePassword: boolean
    readOnly: boolean
    noAuth: boolean
}

export interface UserState {
    authenticated: boolean
    sessionChecked: boolean
    userInfo: UserInfo
}

export interface AppState {
    language: string
}

export interface themeState {
    dark: boolean
}
