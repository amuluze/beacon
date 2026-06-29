import { describe, expect, it } from 'vitest'
import { useUserStore } from '@/store/modules/user'

describe('useUserStore', () => {
    it('initial state is empty', () => {
        const store = useUserStore()
        expect(store.token).toBe('')
        expect(store.refresh).toBe('')
        expect(store.userInfo).toEqual({})
    })

    it('setToken updates token and refresh', () => {
        const store = useUserStore()
        store.setToken('access-123', 'refresh-456')
        expect(store.token).toBe('access-123')
        expect(store.refresh).toBe('refresh-456')
    })

    it('setToken can clear tokens', () => {
        const store = useUserStore()
        store.setToken('access-123', 'refresh-456')
        store.setToken('', '')
        expect(store.token).toBe('')
        expect(store.refresh).toBe('')
    })

    it('setUserInfo updates user info fields', () => {
        const store = useUserStore()
        store.setUserInfo('admin', 1, 1)
        expect(store.userInfo.name).toBe('admin')
        expect(store.userInfo.status).toBe(1)
        expect(store.userInfo.isAdmin).toBe(1)
    })
})
