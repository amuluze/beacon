import { beforeEach, describe, expect, it, vi } from 'vitest'
import { login } from './index'

const mocks = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn(),
}))

vi.mock('@/api', () => ({
    default: {
        get: mocks.get,
        post: mocks.post,
    },
}))

describe('auth API', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.get.mockResolvedValue({ status: 204 })
        mocks.post.mockResolvedValue({ data: { must_change_password: false } })
    })

    it('provisions a CSRF Cookie before submitting credentials', async () => {
        await login({ username: 'admin', password: 'password' })

        expect(mocks.get).toHaveBeenCalledWith('/api/v1/auth/csrf')
        expect(mocks.post).toHaveBeenCalledWith('/api/v1/auth/login', {
            username: 'admin',
            password: 'password',
        })
        expect(mocks.get.mock.invocationCallOrder[0]).toBeLessThan(mocks.post.mock.invocationCallOrder[0])
    })
})
