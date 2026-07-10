import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import UserSettings from './UserSettings.vue'

describe('user settings', () => {
    it('groups user, role, and API management in one settings section', () => {
        const wrapper = mount(UserSettings, {
            global: {
                mocks: {
                    $t: (key: string) => key,
                },
                stubs: {
                    'el-tabs': { template: '<div data-testid="user-tabs"><slot /></div>' },
                    'el-tab-pane': {
                        props: ['label', 'name'],
                        template: '<section :data-tab="name"><span>{{ label }}</span><slot /></section>',
                    },
                    'UserManager': { template: '<div data-testid="user-manager" />' },
                    'RoleManager': { template: '<div data-testid="role-manager" />' },
                    'ApiManager': { template: '<div data-testid="api-manager" />' },
                },
            },
        })

        expect(wrapper.findAll('[data-tab]').map(tab => tab.attributes('data-tab'))).toEqual([
            'users',
            'roles',
            'apis',
        ])
        expect(wrapper.text()).toContain('menu.userManager')
        expect(wrapper.text()).toContain('menu.roleManager')
        expect(wrapper.text()).toContain('menu.apiManager')
        expect(wrapper.find('[data-testid="user-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="api-manager"]').exists()).toBe(true)
    })
})
