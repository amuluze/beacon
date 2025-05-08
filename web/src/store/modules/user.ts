import type { UserInfo, UserState } from '@/interface/store'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
    state: (): UserState => <UserState>({
        token: '',
        refresh: '',
        userInfo: {} as UserInfo,
    }),
    actions: {
        setToken(token: string, refresh: string) {
            this.token = token
            this.refresh = refresh
        },
        setUserInfo(name: string, status: number, isAmin: number) {
            this.userInfo.name = name
            this.userInfo.status = status
            this.userInfo.isAdmin = isAmin
        },
    },
    persist: true,
})
