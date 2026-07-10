import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

import Avatar from './Avatar.vue'

vi.mock('@/api/auth', () => ({
    logout: vi.fn(),
}))

vi.mock('@/hooks/useCommandComponent.ts', () => ({
    default: () => vi.fn(),
}))

vi.mock('@/store', () => ({
    default: () => ({
        user: {
            userInfo: { name: 'admin' },
            setToken: vi.fn(),
        },
    }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn(),
        replace: vi.fn(),
    }),
}))

describe('navbar avatar', () => {
    it('keeps account actions without exposing the removed profile page', () => {
        const wrapper = mount(Avatar, {
            global: {
                stubs: {
                    'el-dropdown': { template: '<div><slot /><slot name="dropdown" /></div>' },
                    'el-dropdown-menu': { template: '<div data-testid="avatar-menu"><slot /></div>' },
                    'el-dropdown-item': { template: '<button><slot /></button>' },
                    'el-divider': true,
                    'svg-icon': true,
                },
            },
        })

        const menu = wrapper.get('[data-testid="avatar-menu"]')

        expect(menu.text()).not.toContain('avatar.profile')
        expect(menu.text()).toContain('avatar.updatePassword')
        expect(menu.text()).toContain('avatar.logout')
    })
})
