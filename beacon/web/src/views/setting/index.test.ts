import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import SettingWorkspace from './index.vue'

const stubStore = {
    user: {
        userInfo: { name: 'admin' },
    },
}

const routerReplace = vi.fn()

vi.mock('@/store', () => ({
    default: () => stubStore,
}))

vi.mock('vue-router', async (importOriginal) => {
    const actual = await importOriginal<typeof import('vue-router')>()
    return {
        ...actual,
        useRoute: () => ({ query: {} }),
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
                AlarmSettings: { template: '<div data-testid="alarm-settings" />' },
                AuditLog: { template: '<div data-testid="audit-log" />' },
                DockerSettings: { template: '<div data-testid="docker-settings" />' },
                UserSettings: { template: '<div data-testid="user-settings" />' },
            },
        },
    })
}

describe('setting workspace', () => {
    beforeEach(() => {
        stubStore.user.userInfo.name = 'admin'
        routerReplace.mockReset()
    })

    it('renders the configuration tab by default', () => {
        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="alarm-settings"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="docker-settings"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="user-settings"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="audit-log"]').exists()).toBe(false)
    })

    it('switches between the workspace tabs and syncs the query', async () => {
        const wrapper = mountWorkspace()
        const tabs = wrapper.findAll('.workspace__tab')

        expect(tabs.map(tab => tab.text())).toEqual([
            'menu.setting',
            'menu.userManager',
            'audit.title',
        ])

        await tabs[1].trigger('click')
        expect(wrapper.find('[data-testid="user-settings"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="alarm-settings"]').exists()).toBe(false)
        expect(routerReplace).toHaveBeenLastCalledWith({ query: { tab: 'users' } })

        await wrapper.findAll('.workspace__tab')[2].trigger('click')
        expect(wrapper.find('[data-testid="audit-log"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="user-settings"]').exists()).toBe(false)
        expect(routerReplace).toHaveBeenLastCalledWith({ query: { tab: 'audit' } })
    })

    it('keeps user management hidden from non-administrators', () => {
        stubStore.user.userInfo.name = 'operator'

        const wrapper = mountWorkspace()

        expect(wrapper.findAll('.workspace__tab').map(tab => tab.text())).toEqual([
            'menu.setting',
            'audit.title',
        ])
        expect(wrapper.find('[data-testid="user-settings"]').exists()).toBe(false)
    })
})
