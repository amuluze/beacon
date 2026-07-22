import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import Avatar from './Avatar.vue'

const { clearAgent, logoutMock, openUpdatePassword, replace, setToken } = vi.hoisted(() => ({
    clearAgent: vi.fn(),
    logoutMock: vi.fn(),
    openUpdatePassword: vi.fn(),
    replace: vi.fn(),
    setToken: vi.fn(),
}))

vi.mock('@/api/auth', () => ({
    logout: logoutMock,
}))

vi.mock('@/hooks/useCommandComponent.ts', () => ({
    default: () => openUpdatePassword,
}))

vi.mock('@/store', () => ({
    default: () => ({
        user: {
            userInfo: { name: 'admin' },
            setToken,
        },
        agent: { clear: clearAgent },
    }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn(),
        replace,
    }),
}))

describe('navbar avatar', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        logoutMock.mockResolvedValue(undefined)
    })

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

    it('passes the logged-in username when opening the password drawer', async () => {
        const wrapper = mount(Avatar, {
            global: {
                stubs: {
                    'el-dropdown': { template: '<div><slot /><slot name="dropdown" /></div>' },
                    'el-dropdown-menu': { template: '<div data-testid="avatar-menu"><slot /></div>' },
                    'el-dropdown-item': {
                        emits: ['click'],
                        template: '<button @click="$emit(\'click\', $event)"><slot /></button>',
                    },
                    'el-divider': true,
                    'svg-icon': true,
                },
            },
        })

        await wrapper.get('[data-testid="avatar-menu"] button').trigger('click')

        expect(openUpdatePassword).toHaveBeenCalledWith({
            title: '更新密码',
            username: 'admin',
        })
    })

    it('clears the authenticated Agent state before returning to login', async () => {
        const wrapper = mount(Avatar, {
            global: {
                stubs: {
                    'el-dropdown': { template: '<div><slot /><slot name="dropdown" /></div>' },
                    'el-dropdown-menu': { template: '<div data-testid="avatar-menu"><slot /></div>' },
                    'el-dropdown-item': {
                        emits: ['click'],
                        template: '<button @click="$emit(\'click\', $event)"><slot /></button>',
                    },
                    'el-divider': true,
                    'svg-icon': true,
                },
            },
        })

        await wrapper.findAll('[data-testid="avatar-menu"] button')[1].trigger('click')
        await flushPromises()

        expect(logoutMock).toHaveBeenCalledOnce()
        expect(setToken).toHaveBeenCalledWith('', '')
        expect(clearAgent).toHaveBeenCalledOnce()
        expect(replace).toHaveBeenCalledWith('/login')
    })
})
