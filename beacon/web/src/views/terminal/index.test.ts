import { mount } from '@vue/test-utils'
import { defineComponent, h, ref } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import type { AgentInfo } from '@/interface/agent'
import TerminalView from './index.vue'

function agent(agent_id: string, hostname = agent_id): AgentInfo {
    return { agent_id, hostname, version: 'v1.0.0', os: 'linux', arch: 'amd64', status: 'online' }
}

const stubClear = vi.fn()
const stubNewSession = vi.fn()

const TerminalStub = defineComponent({
    name: 'Terminal',
    emits: ['statusChange', 'resize'],
    setup(_props, { expose }) {
        expose({
            clear: stubClear,
            newSession: stubNewSession,
        })
        return () => h('div', { 'data-testid': 'terminal-stub' })
    },
})

const agentList = ref<AgentInfo[]>([])
const selectedAgentID = ref('')
const loading = ref(false)
const isAgentEmpty = ref(false)
const loadAgents = vi.fn(async () => selectedAgentID.value)
const ensureSelectedAgent = vi.fn(async () => selectedAgentID.value)

vi.mock('@/hooks/useAgentSelection', () => ({
    useAgentSelection: () => ({
        agentList,
        selectedAgentID,
        loading,
        isAgentEmpty,
        loadAgents,
        ensureSelectedAgent,
    }),
}))

vi.mock('vue-router', () => ({
    useRoute: () => ({ query: {} }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

async function mountView() {
    return mount(TerminalView, {
        global: {
            mocks: { $t: (key: string) => key },
            stubs: {
                'Terminal': TerminalStub,
                'el-select': { template: '<div class="select-stub"><slot /></div>' },
                'el-option': { template: '<div class="option-stub" />' },
                'svg-icon': { template: '<i class="svg-stub" />' },
                'AgentEmptyState': { template: '<div data-testid="agent-empty" />' },
                'i-lucide-server': { template: '<i />' },
                'i-lucide-eraser': { template: '<i />' },
                'i-lucide-plus': { template: '<i />' },
            },
        },
    })
}

describe('terminal workspace', () => {
    beforeEach(async () => {
        vi.clearAllMocks()
        agentList.value = [agent('agent-a', 'host-a')]
        selectedAgentID.value = 'agent-a'
        loading.value = false
        isAgentEmpty.value = false
    })

    afterEach(() => {
        agentList.value = []
        selectedAgentID.value = ''
        isAgentEmpty.value = false
    })

    it('renders workspace header and toolbar controls', async () => {
        const wrapper = await mountView()
        await wrapper.vm.$nextTick()
        expect(wrapper.find('.am-terminal-page__eyebrow').text()).toBe('agent.terminalEyebrow')
        expect(wrapper.find('.am-terminal-page__title').text()).toBe('agent.terminalTitle')
        expect(wrapper.find('.am-terminal-page__hint').text()).toBe('agent.terminalHint')
        expect(wrapper.find('[data-testid="terminal-stub"]').exists()).toBe(true)
        expect(wrapper.find('.am-terminal-page__btn--ghost').text()).toContain('agent.terminalClear')
        expect(wrapper.find('.am-terminal-page__btn--primary').text()).toContain('agent.terminalNewSession')
    })

    it('reflects connection status in the status chip', async () => {
        const wrapper = await mountView()
        await wrapper.vm.$nextTick()
        const terminal = wrapper.findComponent(TerminalStub)

        terminal.vm.$emit('statusChange', 'idle')
        await wrapper.vm.$nextTick()
        expect(wrapper.find('.am-terminal-page__status').text()).toContain('agent.terminalStatusDisconnected')

        terminal.vm.$emit('statusChange', 'connecting')
        await wrapper.vm.$nextTick()
        expect(wrapper.find('.am-terminal-page__status').text()).toContain('agent.terminalStatusConnecting')

        terminal.vm.$emit('statusChange', 'connected')
        await wrapper.vm.$nextTick()
        expect(wrapper.find('.am-terminal-page__status').text()).toContain('agent.terminalStatusConnected')

        terminal.vm.$emit('statusChange', 'error')
        await wrapper.vm.$nextTick()
        expect(wrapper.find('.am-terminal-page__status').text()).toContain('agent.terminalStatusError')
    })

    it('derives the panel title from hostname and terminal dimensions', async () => {
        const wrapper = await mountView()
        await wrapper.vm.$nextTick()
        wrapper.findComponent(TerminalStub).vm.$emit('resize', { rows: 40, cols: 120 })
        await wrapper.vm.$nextTick()
        expect(wrapper.find('.am-terminal-page__panel-title').text()).toBe('host-a — bash — 120x40')
    })

    it('invokes clear when the clear button is clicked', async () => {
        const wrapper = await mountView()
        await wrapper.vm.$nextTick()
        await wrapper.find('.am-terminal-page__btn--ghost').trigger('click')
        expect(stubClear).toHaveBeenCalled()
    })

    it('invokes newSession when the new session button is clicked', async () => {
        const wrapper = await mountView()
        await wrapper.vm.$nextTick()
        await wrapper.find('.am-terminal-page__btn--primary').trigger('click')
        expect(stubNewSession).toHaveBeenCalled()
    })

    it('disables actions and shows the empty state when no agent is available', async () => {
        agentList.value = []
        selectedAgentID.value = ''
        isAgentEmpty.value = true
        const wrapper = await mountView()
        await wrapper.vm.$nextTick()

        expect(wrapper.find('[data-testid="agent-empty"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="terminal-stub"]').exists()).toBe(false)
        expect(wrapper.find('.am-terminal-page__btn--ghost').attributes('disabled')).toBeDefined()
        expect(wrapper.find('.am-terminal-page__btn--primary').attributes('disabled')).toBeDefined()
    })
})
