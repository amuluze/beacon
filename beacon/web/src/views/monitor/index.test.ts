import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import MonitorWorkspace from './index.vue'

// 子视图各自持有轮询与 API 依赖，工作区测试只关心分组结构
vi.mock('./host/index.vue', () => ({
    default: { name: 'HostMonitor', template: '<div data-testid="host-monitor" />' },
}))

vi.mock('./container/index.vue', () => ({
    default: { name: 'ContainerMonitor', template: '<div data-testid="container-monitor" />' },
}))

function mountWorkspace() {
    return mount(MonitorWorkspace, {
        global: {
            mocks: {
                $t: (key: string) => key,
            },
        },
    })
}

describe('monitor workspace', () => {
    it('renders the observability header and workspace description', () => {
        const wrapper = mountWorkspace()

        expect(wrapper.find('.workspace__eyebrow').text()).toBe('OBSERVABILITY')
        expect(wrapper.find('.workspace__title').text()).toBe('menu.monitor')
        expect(wrapper.find('.workspace__description').text()).toBe('monitor.workspaceDescription')
        expect(wrapper.find('.workspace__hint').text()).toBe('Host & Container')

        wrapper.unmount()
    })

    it('stacks the host monitor panel above the container monitor panel', () => {
        const wrapper = mountWorkspace()

        const panels = wrapper.findAll('.workspace__panel--monitor')
        expect(panels).toHaveLength(2)
        expect(panels[0].find('[data-testid="host-monitor"]').exists()).toBe(true)
        expect(panels[0].find('[data-testid="container-monitor"]').exists()).toBe(false)
        expect(panels[1].find('[data-testid="container-monitor"]').exists()).toBe(true)
        expect(panels[1].find('[data-testid="host-monitor"]').exists()).toBe(false)

        wrapper.unmount()
    })
})
