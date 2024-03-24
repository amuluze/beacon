/**
 * @Author     : Amu
 * @Date       : 2024/3/24 17:18
 * @Description:
 */

import { AppState } from '@/interface/store';
import { defineStore } from 'pinia';

export const useAppStore = defineStore('app', {
    state: (): AppState => ({
        isCollapse: false,
    }),
    getters: {},
    actions: {
        setCollapse(collapse: boolean) {
            this.isCollapse = collapse;
        },
    },
    persist: true,
});
