export interface UserInfo {
    name: string
    status: number
    isAdmin: number
}

export interface UserState {
    token: string
    refresh: string
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
