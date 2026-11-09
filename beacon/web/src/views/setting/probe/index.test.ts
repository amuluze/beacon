import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ProbeManager from './index.vue'

const mocks = vi.hoisted(() => ({
    queryAgentList: vi.fn(),
    queryAgentUpdateSetting: vi.fn(),
    updateAgentUpdateSetting: vi.fn(),
    querySystemVersion: vi.fn(),
    openUpgradeDialog: vi.fn(),
    isAdmin: 1,
}))

vi.mock('@/api/agent', () => ({
    queryAgentList: mocks.queryAgentList,
    queryAgentUpdateSetting: mocks.queryAgentUpdateSetting,
    updateAgentUpdateSetting: mocks.updateAgentUpdateSetting,
}))

vi.mock('@/api/system', () => ({
    querySystemVersion: mocks.querySystemVersion,
}))

vi.mock('@/hooks/useCommandComponent.ts', () => ({
    default: () => mocks.openUpgradeDialog,
}))

vi.mock('@/store', () => ({
    default: () => ({
        user: {
            userInfo: { name: 'admin', isAdmin: mocks.isAdmin },
        },
    }),
}))

const agents = [
    {
        agent_id: 'agent-1',
        hostname: 'beacon-server',
        version: 'v3.0.15',
        os: 'linux',
        arch: 'amd64',
        status: 'online',
    },
    {
        agent_id: 'agent-2',
        hostname: 'worker-node',
        version: 'v3.0.14',
        os: 'linux',
        arch: 'arm64',
        status: 'offline',
    },
]

function mountProbe() {
    return mount(ProbeManager, {
        global: {
            mocks: {
                $t: (key: string, params?: Record<string, unknown>) => (params ? `${key}:${JSON.stringify(params)}` : key),
            },
            stubs: {
                'svg-icon': true,
            },
        },
    })
}

describe('probe management', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.isAdmin = 1
        mocks.queryAgentList.mockResolvedValue({ data: agents })
        mocks.queryAgentUpdateSetting.mockResolvedValue({ data: { auto_update: true, deployment_enabled: true } })
        mocks.updateAgentUpdateSetting.mockResolvedValue({ data: { auto_update: false, deployment_enabled: true } })
        mocks.querySystemVersion.mockResolvedValue({ data: { version: 'v3.0.15', update_available: false } })
    })

    it('lists installed probes with name, status and version', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        const table = wrapper.findComponent({ name: 'ElTable' })
        expect(table.props('data')).toEqual(agents)
        expect(wrapper.text()).toContain('beacon-server')
        expect(wrapper.text()).toContain('v3.0.15')
        expect(wrapper.text()).toContain('probe.online')
        expect(wrapper.text()).toContain('probe.offline')
    })

    it('opens the upgrade dialog for an online probe', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        const upgradeButtons = wrapper.findAll('button').filter(button => button.text().includes('probe.upgrade'))
        expect(upgradeButtons.length).toBeGreaterThan(0)

        await upgradeButtons[0].trigger('click')

        expect(mocks.openUpgradeDialog).toHaveBeenCalledWith(expect.objectContaining({
            agentID: 'agent-1',
            hostname: 'beacon-server',
            version: 'v3.0.15',
        }))
    })

    it('hides the upgrade action for non administrators', async () => {
        mocks.isAdmin = 0

        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        const upgradeButtons = wrapper.findAll('button').filter(button => button.text().includes('probe.upgrade'))
        expect(upgradeButtons.length).toBe(0)
    })

    it('filters probes by agent name keyword', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        wrapper.findComponent({ name: 'ElInput' }).vm.$emit('update:modelValue', 'worker')
        await nextTick()

        const table = wrapper.findComponent({ name: 'ElTable' })
        expect(table.props('data')).toHaveLength(1)
        expect(table.props('data')[0].hostname).toBe('worker-node')
    })

    it('filters probes by connection status', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        wrapper.findComponent({ name: 'ElRadioGroup' }).vm.$emit('update:modelValue', 'offline')
        await nextTick()

        const table = wrapper.findComponent({ name: 'ElTable' })
        expect(table.props('data')).toHaveLength(1)
        expect(table.props('data')[0].status).toBe('offline')
    })

    it('marks probes that lag behind the latest beacon version', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))
        await vi.waitFor(() => expect(wrapper.text()).toContain('probe.latest'))

        // beacon-server 与最新版本一致，worker-node 落后一个版本。
        expect(wrapper.text()).toContain('probe.upgradable')
        expect(wrapper.findAll('.am-probe__badge--upgradable')).toHaveLength(1)
    })

    it('renders the platform of each probe next to its name', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        expect(wrapper.text()).toContain('linux/amd64')
        expect(wrapper.text()).toContain('linux/arm64')
    })

    it('renders an enabled auto-update switch for administrators', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        const toggle = wrapper.findComponent({ name: 'ElSwitch' })
        expect(toggle.exists()).toBe(true)
        await vi.waitFor(() => expect(toggle.props('modelValue')).toBe(true))
        expect(toggle.props('disabled')).toBe(false)
    })

    it('persists the auto-update switch when toggled off', async () => {
        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        const toggle = wrapper.findComponent({ name: 'ElSwitch' })
        await vi.waitFor(() => expect(toggle.props('modelValue')).toBe(true))
        toggle.vm.$emit('change', false)

        await vi.waitFor(() => expect(mocks.updateAgentUpdateSetting).toHaveBeenCalledWith(false))
    })

    it('hides the auto-update switch for non administrators', async () => {
        mocks.isAdmin = 0

        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        expect(wrapper.findComponent({ name: 'ElSwitch' }).exists()).toBe(false)
        expect(mocks.queryAgentUpdateSetting).not.toHaveBeenCalled()
    })

    it('disables the switch when the deployment blocks auto updates', async () => {
        mocks.queryAgentUpdateSetting.mockResolvedValue({ data: { auto_update: true, deployment_enabled: false } })

        const wrapper = mountProbe()
        await vi.waitFor(() => expect(wrapper.findComponent({ name: 'ElTable' }).props('data')).toHaveLength(2))

        const toggle = wrapper.findComponent({ name: 'ElSwitch' })
        await vi.waitFor(() => expect(toggle.props('disabled')).toBe(true))
    })
})
