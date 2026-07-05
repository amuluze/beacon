/**
 * @Author     : Amu
 * @Date       : 2025/1/19 00:06
 * @Description:
 */

import type { echartsThemeState } from '@/interface/store.ts'

export const useEChartsStore = defineStore('echarts', {
    state: (): echartsThemeState => ({
        currentColorArray: [],
    }),
    getters: {},
    actions: {
        setCurrentColorArray(colorArray: string[]) {
            this.currentColorArray = colorArray
        },
    },
})
