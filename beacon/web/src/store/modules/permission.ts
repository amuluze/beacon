/**
 * @Author     : Amu
 * @Date       : 2024/11/10 23:31
 * @Description:
 */
import type { PermissionsState } from '@/interface/store.ts'
import { defineStore } from 'pinia'

export const usePermissionStore = defineStore('permission', {
    state: (): PermissionsState => ({
        permissions: [],
    }),
    getters: {},
    actions: {
        setPermissions(permissions: string[]) {
            this.permissions = permissions
        },
    },
    persist: true,
})
