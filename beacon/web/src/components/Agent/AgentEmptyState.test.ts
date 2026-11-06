import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { afterEach, describe, expect, it, vi } from 'vitest'
import AgentEmptyState from './AgentEmptyState.vue'
import { trackWrapper, unmountTrackedWrappers } from '@/__tests__/helpers/chartTestUtils'

const translator = vi.hoisted(() => vi.fn((key: string) => key))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: translator }),
}))

// 安装弹窗会拉取安装令牌，空态测试只关心打开时机
vi.mock('./InstallAgentDialog.vue', () => ({
    default: {
        name: 'InstallAgentDialog',
        props: {
            visible: { type: Boolean, default: false },
            title: String,
        },
        template: '<div data-testid="install-dialog" :data-visible="String(visible)" />',
    },
}))

const ElButtonStub = defineComponent({
    name: 'ElButton',
    props: {
        type: String,
    },
    emits: ['click'],
    template: '<button type="button" @click="$emit(\'click\')"><slot /></button>',
})

const SvgIconStub = defineComponent({
    name: 'SvgIcon',
    props: {
        iconClass: String,
    },
    template: '<i />',
})

interface EmptyStateProps {
    title?: string
    description?: string
    minHeight?: string
    showRefresh?: boolean
    showInstall?: boolean
}

function mountEmptyState(props: EmptyStateProps = {}) {
    const wrapper = mount(AgentEmptyState, {
        props,
        global: {
            stubs: {
                'el-button': ElButtonStub,
                'ElButton': ElButtonStub,
                'svg-icon': SvgIconStub,
                'SvgIcon': SvgIconStub,
            },
        },
    })
    return trackWrapper(wrapper)
}

afterEach(() => {
    unmountTrackedWrappers()
})

describe('agent empty state', () => {
    it('falls back to the shipped copy when no props are given', () => {
        const wrapper = mountEmptyState()

        expect(wrapper.find('.agent-empty-state__title').text()).toBe('agent.emptyTitle')
        expect(wrapper.find('.agent-empty-state__description').text()).toBe('agent.emptyDescription')
        expect(wrapper.attributes('style')).toContain('min-height: 320px')
    })

    it('lets the page override the copy and the minimum height', () => {
        const wrapper = mountEmptyState({
            title: '暂无采集节点',
            description: '先安装 collia Agent',
            minHeight: '480px',
        })

        expect(wrapper.find('.agent-empty-state__title').text()).toBe('暂无采集节点')
        expect(wrapper.find('.agent-empty-state__description').text()).toBe('先安装 collia Agent')
        expect(wrapper.attributes('style')).toContain('min-height: 480px')
    })

    it('offers install and refresh actions by default', () => {
        const wrapper = mountEmptyState()

        const buttons = wrapper.findAll('button')
        expect(buttons).toHaveLength(2)
        expect(buttons[0].text()).toContain('agent.installAction')
        expect(buttons[1].text()).toContain('agent.refresh')
    })

    it('hides the install action when the page owns installation', () => {
        const wrapper = mountEmptyState({ showInstall: false })

        const buttons = wrapper.findAll('button')
        expect(buttons).toHaveLength(1)
        expect(buttons[0].text()).toContain('agent.refresh')
    })

    it('hides the refresh action when the page refreshes elsewhere', () => {
        const wrapper = mountEmptyState({ showRefresh: false })

        const buttons = wrapper.findAll('button')
        expect(buttons).toHaveLength(1)
        expect(buttons[0].text()).toContain('agent.installAction')
    })

    it('emits refresh from the refresh action', async () => {
        const wrapper = mountEmptyState()

        await wrapper.findAll('button')[1].trigger('click')

        expect(wrapper.emitted('refresh')).toHaveLength(1)
    })

    it('opens the install dialog from the install action', async () => {
        const wrapper = mountEmptyState()
        expect(wrapper.find('[data-testid="install-dialog"]').attributes('data-visible')).toBe('false')

        await wrapper.findAll('button')[0].trigger('click')
        await nextTick()

        expect(wrapper.find('[data-testid="install-dialog"]').attributes('data-visible')).toBe('true')
    })
})
