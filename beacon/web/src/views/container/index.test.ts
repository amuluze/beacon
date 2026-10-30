import { flushPromises, mount } from '@vue/test-utils'
import { computed, defineComponent } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ContainerWorkspace from './index.vue'

const routerState = vi.hoisted(() => ({
    replace: vi.fn(),
    query: {},
}))

vi.mock('vue-router', () => ({
    useRoute: () => ({ query: routerState.query }),
    useRouter: () => ({ replace: routerState.replace }),
}))

const agentState = vi.hoisted(() => ({ isEmpty: false, loadAgents: vi.fn() }))

vi.mock('@/hooks/useAgentSelection', async () => {
    const { shallowRef } = await import('vue')
    const selectedAgentID = shallowRef('agent-1')
    return {
        useAgentSelection: () => ({
            isAgentEmpty: agentState.isEmpty,
            loadAgents: agentState.loadAgents,
            selectedAgentID,
        }),
    }
})

vi.mock('@/api/container', () => ({
    queryContainerSummary: vi.fn(async () => ({
        data: { containers: 7, containers_running: 5, images: 6, image_bytes: 4718592, networks: 6 },
    })),
}))

// 最小 el-tabs 替身：按 modelValue 只渲染命中的 pane，用于断言当前激活的 Tab。
const ElTabsStub = defineComponent({
    name: 'ElTabsStub',
    props: { modelValue: String },
    emits: ['update:modelValue'],
    provide() {
        return { activeTab: computed(() => this.modelValue) }
    },
    template: '<div class="el-tabs-stub" :data-active="modelValue"><slot /></div>',
})

const ElTabPaneStub = defineComponent({
    name: 'ElTabPaneStub',
    props: { name: String, label: String },
    inject: ['activeTab'],
    template: `
    <section class="el-tab-pane-stub" :data-tab="name">
      <header class="el-tab-pane-stub__label"><slot name="label" /></header>
      <div v-if="activeTab === name" class="el-tab-pane-stub__content"><slot /></div>
    </section>
  `,
})

function mountWorkspace() {
    return mount(ContainerWorkspace, {
        global: {
            mocks: { $t: (key: string) => key },
            stubs: {
                'el-tabs': ElTabsStub,
                'el-tab-pane': ElTabPaneStub,
                'ContainerManager': { template: '<div data-testid="container-manager" />' },
                'ImageManager': { template: '<div data-testid="image-manager" />' },
                'NetworkManager': { template: '<div data-testid="network-manager" />' },
                'AgentEmptyState': { template: '<div data-testid="agent-empty" />' },
            },
        },
    })
}

describe('container workspace', () => {
    beforeEach(() => {
        routerState.query = {}
        routerState.replace.mockClear()
        agentState.loadAgents.mockClear()
        agentState.isEmpty = false
    })

    it('splits containers, images and networks into tabs and opens containers first', () => {
        const wrapper = mountWorkspace()

        expect(wrapper.findAll('.el-tab-pane-stub').map(pane => pane.attributes('data-tab'))).toEqual([
            'containers',
            'images',
            'networks',
        ])
        expect(wrapper.find('.el-tabs-stub').attributes('data-active')).toBe('containers')
        expect(wrapper.find('[data-testid="container-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="image-manager"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="network-manager"]').exists()).toBe(false)
        expect(wrapper.text()).toContain('container.tabContainers')
        expect(wrapper.text()).toContain('container.tabImages')
        expect(wrapper.text()).toContain('container.tabNetworks')
    })

    it('honours the tab deep link', () => {
        routerState.query = { tab: 'images' }

        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="image-manager"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="container-manager"]').exists()).toBe(false)
    })

    it('falls back to the containers tab for an unknown tab value', () => {
        routerState.query = { tab: 'volumes' }

        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="container-manager"]').exists()).toBe(true)
    })

    it('renders a single empty state instead of the three stacked tabs', () => {
        agentState.isEmpty = true

        const wrapper = mountWorkspace()

        expect(wrapper.find('[data-testid="agent-empty"]').exists()).toBe(true)
        expect(wrapper.find('.el-tabs-stub').exists()).toBe(false)
    })

    it('renders the server summary metrics in the workspace header', async () => {
        const wrapper = mountWorkspace()
        await flushPromises()

        expect(wrapper.findAll('.container-workspace__metric')).toHaveLength(3)
        expect(wrapper.text()).toContain('4.50 MB')
    })

    it('keeps the active tab in the route query', async () => {
        const wrapper = mountWorkspace()

        wrapper.findComponent(ElTabsStub).vm.$emit('update:modelValue', 'networks')
        await wrapper.vm.$nextTick()

        expect(routerState.replace).toHaveBeenCalledWith({ query: { tab: 'networks' } })
        expect(wrapper.find('[data-testid="network-manager"]').exists()).toBe(true)
    })
})
