import { UserState } from '@/interface/store'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
    state: (): UserState => ({
        token: '',
        refresh: ''
    }),
    actions: {
        setToken(token: string, refresh: string) {
            this.token = token
            this.refresh = refresh
        }
    },
    persist: true
})
