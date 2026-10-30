import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import UpgradeAgent from './UpgradeAgent.vue'

const mocks = vi.hoisted(() => ({
    openUpgradeDialog: vi.fn(),
    loadAgents: vi.fn(),
    state: {
        isAdmin: 1,
        selectedAgentID: 'agent-1',
        agentList: [
            {
                agent_id: 'agent-1',
                hostname: 'beacon-server',
                version: 'v3.0.5',
                os: 'linux',
                arch: 'amd64',
                status: 'online',
            },
        ],
    },
}))

vi.mock('@/hooks/useCommandComponent.ts', () => ({
    default: () => mocks.openUpgradeDialog,
}))

vi.mock('@/hooks/useAgentSelection', async () => {
    const { computed } = await import('vue')
    return {
        useAgentSelection: () => ({
            agentList: computed(() => mocks.state.agentList),
            selectedAgentID: computed(() => mocks.state.selectedAgentID),
            loadAgents: mocks.loadAgents,
        }),
    }
})

vi.mock('@/store', () => ({
    default: () => ({
        user: {
            userInfo: { name: 'admin', isAdmin: mocks.state.isAdmin },
        },
    }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const stubs = {
    'svg-icon': true,
}

describe('navbar upgrade trigger', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.state.isAdmin = 1
        mocks.state.selectedAgentID = 'agent-1'
        mocks.state.agentList = [
            {
                agent_id: 'agent-1',
                hostname: 'beacon-server',
                version: 'v3.0.5',
                os: 'linux',
                arch: 'amd64',
                status: 'online',
            },
        ]
    })

    it('hides the entry for non administrators', () => {
        mocks.state.isAdmin = 0
        const wrapper = mount(UpgradeAgent, { global: { stubs } })

        expect(wrapper.find('button').exists()).toBe(false)
    })

    it('disables the entry until an agent is selected', () => {
        mocks.state.selectedAgentID = ''
        const wrapper = mount(UpgradeAgent, { global: { stubs } })

        expect(wrapper.find('button').attributes('disabled')).toBeDefined()
    })

    it('opens the upgrade dialog for the selected agent', async () => {
        const wrapper = mount(UpgradeAgent, { global: { stubs } })

        await wrapper.find('button').trigger('click')

        expect(mocks.openUpgradeDialog).toHaveBeenCalledWith(expect.objectContaining({
            title: 'agent.upgradeTitle',
            agentID: 'agent-1',
            hostname: 'beacon-server',
            version: 'v3.0.5',
        }))
    })
})
