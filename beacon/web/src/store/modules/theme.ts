/**
 * @Author     : Amu
 * @Date       : 2024/11/10 23:54
 * @Description:
 */
import type { themeState } from '@/interface/store.ts'

export const useThemeStore = defineStore('theme', {
    state: (): themeState => ({
        dark: false,
    }),
    getters: {},
    actions: {
        setDark(val: boolean) {
            console.log('set dark value: ', val)
            this.dark = val
        },
    },
    persist: true,
})
