import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import type { AgentUpgradeResult } from '@/interface/agent.ts'
import UpgradeAgentDialog from './UpgradeAgentDialog.vue'

const mocks = vi.hoisted(() => ({
    upgradeAgent: vi.fn(),
    success: vi.fn(),
}))

vi.mock('@/api/agent', () => ({
    upgradeAgent: mocks.upgradeAgent,
}))

vi.mock('@/components/Message/message', () => ({
    success: mocks.success,
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const stubs = {
    'el-dialog': { template: '<div><slot /><slot name="footer" /></div>', props: ['modelValue', 'title', 'width'] },
    'el-button': { template: '<button type="button" @click="$emit(\'click\')"><slot /></button>', emits: ['click'] },
}

function mountDialog() {
    return mount(UpgradeAgentDialog, {
        props: {
            visible: true,
            title: 'agent.upgradeTitle',
            agentID: 'agent-1',
            hostname: 'beacon-server',
            version: 'v3.0.5',
        },
        global: { stubs },
    })
}

describe('upgrade agent dialog', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.upgradeAgent.mockResolvedValue({
            data: {
                agent_id: 'agent-1',
                current_version: 'v3.0.5',
                target_version: 'v3.0.6',
                arch: 'amd64',
                triggered: true,
                message: 'collia upgrade triggered',
            },
        })
    })

    it('shows the target node, the current version and the upgrade constraints', () => {
        const wrapper = mountDialog()

        expect(wrapper.text()).toContain('beacon-server · agent-1')
        expect(wrapper.text()).toContain('v3.0.5')
        expect(wrapper.text()).toContain('agent.upgradeTargetVersionValue')
        expect(wrapper.text()).toContain('agent.upgradeTipConcurrency')
        expect(wrapper.text()).toContain('agent.upgradeTipInterrupt')
    })

    it('triggers the upgrade for the selected agent and closes on success', async () => {
        const wrapper = mountDialog()

        await wrapper.findAll('button')[1].trigger('click')
        await flushPromises()

        expect(mocks.upgradeAgent).toHaveBeenCalledWith('agent-1')
        expect(mocks.success).toHaveBeenCalledOnce()
        expect(wrapper.emitted('update:visible')?.at(-1)).toEqual([false])
        const emitted = wrapper.emitted('success')?.[0]?.[0] as AgentUpgradeResult | undefined
        expect(emitted?.target_version).toBe('v3.0.6')
    })

    it('keeps the dialog open when the upgrade request fails', async () => {
        mocks.upgradeAgent.mockRejectedValue(new Error('upgrade rejected'))
        const wrapper = mountDialog()

        await wrapper.findAll('button')[1].trigger('click')
        await flushPromises()

        expect(mocks.success).not.toHaveBeenCalled()
        expect(wrapper.emitted('update:visible')).toBeFalsy()
    })
})
