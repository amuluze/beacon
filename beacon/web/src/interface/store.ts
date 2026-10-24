export interface UserInfo {
    name: string
    status: number
    isAdmin: number
    mustChangePassword: boolean
}

export interface UserState {
    authenticated: boolean
    sessionChecked: boolean
    userInfo: UserInfo
}

export interface AppState {
    isCollapse: boolean
    language: string
}

export interface themeState {
    dark: boolean
}

export interface echartsThemeState {
    currentColorArray: string[]
}
