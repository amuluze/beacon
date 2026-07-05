/**
 * @Author     : Amu
 * @Date       : 2024/3/24 17:18
 * @Description:
 */

import type { Image, Network } from '@/interface/container'
import type { AppState } from '@/interface/store'
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
    state: (): AppState => ({
        isCollapse: false,
        language: 'zh',
        networks: [],
        images: [],
    }),
    getters: {},
    actions: {
        setCollapse(collapse: boolean) {
            this.isCollapse = collapse
        },
        setLanguage(language: string) {
            this.language = language
        },
        setNetworks(networks: Network[]) {
            this.networks = networks
        },
        setImages(images: Image[]) {
            this.images = images
        },
    },
    persist: true,
})
