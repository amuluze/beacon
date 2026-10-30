import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

import UserSettings from './UserSettings.vue'

vi.mock('vue-router', async (importOriginal) => {
    const actual = await importOriginal<typeof import('vue-router')>()
    return {
        ...actual,
        useRoute: () => ({ query: {} }),
        useRouter: () => ({ replace: vi.fn() }),
    }
})

function mountUserSettings() {
    return mount(UserSettings, {
        global: {
            mocks: {
                $t: (key: string) => key,
            },
            stubs: {
                'el-button': {
                    emits: ['click'],
                    template: '<button type="button" v-bind="$attrs" @click="$emit(\'click\')"><slot /></button>',
                },
                'UserManager': {
                    props: ['createTicket'],
                    template: '<div data-testid="user-manager" :data-ticket="createTicket" />',
                },
                'RoleManager': { template: '<div data-testid="role-manager" />' },
                'ApiManager': { template: '<div data-testid="api-manager" />' },
            },
        },
    })
}

describe('user settings', () => {
    it('groups user, role, and API management behind one view switcher', () => {
        const wrapper = mountUserSettings()

        expect(wrapper.findAll('.user-settings__tab').map(tab => tab.text())).toEqual([
            'menu.userManager',
            'menu.roleManager',
            'menu.apiManager',
        ])
        expect(wrapper.find('[data-testid="user-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="create-user"]').exists()).toBe(true)
    })

    it('keeps a visited view mounted and hides the primary action outside the user view', async () => {
        const wrapper = mountUserSettings()

        await wrapper.findAll('.user-settings__tab')[1].trigger('click')

        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(true)
        expect(wrapper.get('[data-testid="user-manager"]').attributes('style')).toContain('display: none')
        expect(wrapper.find('[data-testid="create-user"]').exists()).toBe(false)

        await wrapper.findAll('.user-settings__tab')[0].trigger('click')

        expect(wrapper.get('[data-testid="user-manager"]').attributes('style') ?? '').not.toContain('display: none')
        expect(wrapper.find('[data-testid="create-user"]').exists()).toBe(true)
    })

    it('asks the user list to open its create dialog from the toolbar', async () => {
        const wrapper = mountUserSettings()

        expect(wrapper.get('[data-testid="user-manager"]').attributes('data-ticket')).toBe('0')

        await wrapper.get('[data-testid="create-user"]').trigger('click')

        expect(wrapper.get('[data-testid="user-manager"]').attributes('data-ticket')).toBe('1')
    })
})
