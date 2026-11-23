/**
 * @Author     : Amu
 * @Date       : 2024/3/24 17:18
 * @Description:
 */

import type { AppState } from '@/interface/store'
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
    state: (): AppState => ({
        language: 'zh',
    }),
    getters: {},
    actions: {
        setLanguage(language: string) {
            this.language = language
        },
    },
    persist: true,
})
