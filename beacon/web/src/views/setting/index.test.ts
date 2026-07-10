import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import SettingWorkspace from './index.vue'

const stubStore = {
    user: {
        userInfo: { name: 'admin' },
    },
}

vi.mock('@/store', () => ({
    default: () => stubStore,
}))

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
                HostSettings: { template: '<div data-testid="host-settings" />' },
                UserSettings: { template: '<div data-testid="user-settings" />' },
            },
        },
    })
}

describe('setting workspace', () => {
    beforeEach(() => {
        stubStore.user.userInfo.name = 'admin'
    })

    it('places administrator user settings next to the audit log', () => {
        const wrapper = mountWorkspace()
        const panelContent = wrapper
            .findAll('.workspace__panel [data-testid]')
            .map(item => item.attributes('data-testid'))

        expect(panelContent).toEqual([
            'host-settings',
            'alarm-settings',
            'docker-settings',
            'user-settings',
            'audit-log',
        ])
    })

    it('keeps user settings hidden from non-administrators', () => {
        stubStore.user.userInfo.name = 'operator'

        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="user-settings"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="audit-log"]').exists()).toBe(true)
    })
})
