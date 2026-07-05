import type { Image, Network } from './container'

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

export interface AgentInfo {
    agent_id: string
    hostname: string
    status: string
    last_seen: string
}

export interface AgentState {
    currentAgentID: string
    agents: AgentInfo[]
}

export interface AppState {
    isCollapse: boolean
    language: string
    networks: Network[]
    images: Image[]
}

export interface themeState {
    dark: boolean
}

export interface echartsThemeState {
    currentColorArray: string[]
}

export interface PermissionsState {
    permissions: string[]
}
