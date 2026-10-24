import { describe, expect, it } from 'vitest'
import { useUserStore } from '@/store/modules/user'

describe('useUserStore', () => {
    it('initial state is empty', () => {
        const store = useUserStore()
        expect(store.authenticated).toBe(false)
        expect(store.sessionChecked).toBe(false)
        expect(store.userInfo).toEqual({})
    })

    it('clearSession removes in-memory user state', () => {
        const store = useUserStore()
        store.setUserInfo('admin', 1, 1, false)
        store.clearSession()
        expect(store.authenticated).toBe(false)
        expect(store.sessionChecked).toBe(true)
        expect(store.userInfo).toEqual({})
    })

    it('setUserInfo updates user info fields', () => {
        const store = useUserStore()
        store.setUserInfo('admin', 1, 1, true)
        expect(store.userInfo.name).toBe('admin')
        expect(store.userInfo.status).toBe(1)
        expect(store.userInfo.isAdmin).toBe(1)
        expect(store.userInfo.mustChangePassword).toBe(true)
        expect(store.authenticated).toBe(true)
        expect(store.sessionChecked).toBe(true)
    })
})
