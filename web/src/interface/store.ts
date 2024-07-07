import { Image, Network } from './container'

export interface UserState {
    token: string
    refresh: string
}

export interface AppState {
    isCollapse: boolean
    networks: Network[]
    images: Image[]
}
