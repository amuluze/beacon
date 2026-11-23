import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import SettingWorkspace from './index.vue'

const stubStore: { user: { userInfo: { name: string, isAdmin: number, readOnly?: boolean } } } = {
    user: {
        userInfo: { name: 'ops-admin', isAdmin: 1 },
    },
}

const routerReplace = vi.fn()
const routeState = vi.hoisted(() => ({ query: {} }))

vi.mock('@/store', () => ({
    default: () => stubStore,
}))

vi.mock('vue-router', async (importOriginal) => {
    const actual = await importOriginal<typeof import('vue-router')>()
    return {
        ...actual,
        useRoute: () => routeState,
        useRouter: () => ({ replace: routerReplace }),
    }
})

function mountWorkspace() {
    return mount(SettingWorkspace, {
        global: {
            mocks: {
                $t: (key: string) => key,
            },
            stubs: {
                'el-button': {
                    emits: ['click'],
                    template: '<button type="button" v-bind="$attrs" @click="$emit(\'click\')"><slot /></button>',
                },
                'AlarmSettings': { template: '<div data-testid="alarm-settings" />' },
                'AuditLog': { template: '<div data-testid="audit-log" />' },
                'DockerSettings': { template: '<div data-testid="docker-settings" />' },
                'ProbeManager': { template: '<div data-testid="probe-manager" />' },
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

function tabLabels(wrapper: ReturnType<typeof mountWorkspace>) {
    return wrapper.findAll('.workspace__tab').map(tab => tab.text())
}

describe('setting workspace', () => {
    beforeEach(() => {
        stubStore.user.userInfo = { name: 'ops-admin', isAdmin: 1 }
        routeState.query = {}
        routerReplace.mockReset()
    })

    it('renders the configuration tab by default', () => {
        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="alarm-settings"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="docker-settings"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="user-manager"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="audit-log"]').exists()).toBe(false)
    })

    it('exposes account, role and API management as first-level tabs', () => {
        expect(tabLabels(mountWorkspace())).toEqual([
            'menu.setting',
            'menu.probeManager',
            'menu.userManager',
            'menu.roleManager',
            'menu.apiManager',
            'audit.title',
        ])
    })

    it('switches each account tab and syncs the query', async () => {
        const wrapper = mountWorkspace()

        await wrapper.findAll('.workspace__tab')[3].trigger('click')
        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="alarm-settings"]').exists()).toBe(false)
        expect(routerReplace).toHaveBeenLastCalledWith({ query: { tab: 'roles' } })

        await wrapper.findAll('.workspace__tab')[4].trigger('click')
        expect(wrapper.find('[data-testid="api-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(false)
        expect(routerReplace).toHaveBeenLastCalledWith({ query: { tab: 'apis' } })

        await wrapper.findAll('.workspace__tab')[5].trigger('click')
        expect(wrapper.find('[data-testid="audit-log"]').exists()).toBe(true)
        expect(routerReplace).toHaveBeenLastCalledWith({ query: { tab: 'audit' } })
    })

    it('opens the user create dialog from the toolbar ticket', async () => {
        const wrapper = mountWorkspace()

        await wrapper.findAll('.workspace__tab')[2].trigger('click')

        expect(wrapper.get('[data-testid="user-manager"]').attributes('data-ticket')).toBe('0')

        await wrapper.get('[data-testid="create-user"]').trigger('click')

        expect(wrapper.get('[data-testid="user-manager"]').attributes('data-ticket')).toBe('1')
    })

    it('opens a role tab directly from the query string', () => {
        routeState.query = { tab: 'roles' }

        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="alarm-settings"]').exists()).toBe(false)
    })

    it('keeps permission tabs hidden from non-administrators', () => {
        stubStore.user.userInfo = { name: 'ops-admin', isAdmin: 0 }

        const wrapper = mountWorkspace()

        expect(tabLabels(wrapper)).toEqual(['menu.setting', 'menu.probeManager', 'audit.title'])
        expect(wrapper.find('[data-testid="user-manager"]').exists()).toBe(false)
    })

    it('does not grant permission tabs by literal username', () => {
        stubStore.user.userInfo = { name: 'admin', isAdmin: 0 }

        const wrapper = mountWorkspace()

        expect(tabLabels(wrapper)).toEqual(['menu.setting', 'menu.probeManager', 'audit.title'])
    })

    it('exposes permission tabs to anonymous read-only sessions', () => {
        stubStore.user.userInfo = { name: '', isAdmin: 0, readOnly: true }

        const wrapper = mountWorkspace()

        expect(tabLabels(wrapper)).toContain('menu.userManager')
        expect(tabLabels(wrapper)).toContain('menu.roleManager')
        expect(tabLabels(wrapper)).toContain('menu.apiManager')
    })

    it('falls back to configuration when a restricted tab is requested in the query', () => {
        stubStore.user.userInfo = { name: 'ops-admin', isAdmin: 0 }
        routeState.query = { tab: 'roles' }

        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="role-manager"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="alarm-settings"]').exists()).toBe(true)
    })
})
