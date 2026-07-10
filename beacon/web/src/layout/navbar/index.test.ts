import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

import Navbar from './index.vue'

const push = vi.fn()

vi.mock('@/store', () => ({
    default: () => ({
        user: {
            userInfo: { name: 'admin' },
        },
    }),
}))

vi.mock('vue-router', () => ({
    useRoute: () => ({ path: '/monitor' }),
    useRouter: () => ({ push }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

describe('navbar', () => {
    it('keeps user settings out of primary navigation', () => {
        const wrapper = mount(Navbar, {
            global: {
                stubs: {
                    'AgentSelect': { template: '<div data-testid="agent-select" />' },
                    'Avatar': true,
                    'Language': true,
                    'ThemeChange': true,
                    'svg-icon': true,
                    'el-dropdown': { template: '<div class="dropdown-stub"><slot /><slot name="dropdown" /></div>' },
                    'el-dropdown-menu': { template: '<div><slot /></div>' },
                    'el-dropdown-item': { template: '<div><slot /></div>' },
                },
            },
        })

        const menuLabels = wrapper
            .findAll('.am-navbar__menu > .am-navbar__menu-item span')
            .map(item => item.text())

        expect(menuLabels).toEqual([
            'menu.monitor',
            'menu.container',
            'menu.setting',
            'menu.terminal',
        ])
        expect(wrapper.find('.dropdown-stub').exists()).toBe(false)
        expect(wrapper.text()).not.toContain('container.more')
    })

    it('shows the global agent selector before navbar utilities', () => {
        const wrapper = mount(Navbar, {
            global: {
                stubs: {
                    'AgentSelect': { template: '<div data-testid="agent-select" />' },
                    'Avatar': { template: '<div data-testid="avatar" />' },
                    'Language': { template: '<div data-testid="language" />' },
                    'ThemeChange': { template: '<div data-testid="theme-change" />' },
                    'svg-icon': true,
                },
            },
        })

        const rightItems = wrapper
            .find('.am-navbar__right')
            .findAll('[data-testid]')
            .map(item => item.attributes('data-testid'))

        expect(rightItems).toEqual([
            'agent-select',
            'language',
            'theme-change',
            'avatar',
        ])
    })
})
