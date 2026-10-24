import type { UserInfo, UserState } from '@/interface/store'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
    state: (): UserState => ({
        authenticated: false,
        sessionChecked: false,
        userInfo: {} as UserInfo,
    }),
    actions: {
        setUserInfo(name: string, status: number, isAdmin: number, mustChangePassword = false) {
            this.userInfo = { name, status, isAdmin, mustChangePassword }
            this.authenticated = true
            this.sessionChecked = true
        },
        clearSession() {
            this.authenticated = false
            this.sessionChecked = true
            this.userInfo = {} as UserInfo
        },
        markSessionChecked() {
            this.sessionChecked = true
        },
    },
})
