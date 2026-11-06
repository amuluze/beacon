import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'
import MonitorContextBar from './MonitorContextBar.vue'
import { ElOptionStub, ElSelectStub, trackWrapper, unmountTrackedWrappers } from '@/__tests__/helpers/chartTestUtils'

const translator = vi.hoisted(() => vi.fn((key: string, params?: Record<string, unknown>) => {
    return params ? `${key}(${Object.values(params).join(',')})` : key
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: translator }),
}))

const densityOptions = [
    { value: 300, label: '5分钟' },
    { value: 600, label: '10分钟' },
]

function mountContextBar(overrides: { density?: number, lastUpdated?: number, stale?: boolean } = {}) {
    const wrapper = mount(MonitorContextBar, {
        props: {
            density: 600,
            densityOptions,
            ...overrides,
        },
        global: {
            stubs: {
                'el-select': ElSelectStub,
                'ElSelect': ElSelectStub,
                'el-option': ElOptionStub,
                'ElOption': ElOptionStub,
            },
        },
    })
    return trackWrapper(wrapper)
}

afterEach(() => {
    unmountTrackedWrappers()
})

describe('monitor context bar', () => {
    it('renders the realtime marker and a dash before the first refresh', () => {
        const wrapper = mountContextBar()

        expect(wrapper.find('.am-monitor-context__live').text()).toBe('monitor.realtime')
        expect(wrapper.find('.am-monitor-context__time').text()).toBe('monitor.lastUpdated(-)')
        expect(wrapper.find('.am-monitor-context__density-label').text()).toBe('monitor.timeDensity：')
    })

    it('formats the last updated time as HH:mm:ss', () => {
        const lastUpdated = new Date(2024, 0, 2, 3, 4, 5).getTime()

        const wrapper = mountContextBar({ lastUpdated })

        expect(wrapper.find('.am-monitor-context__time').text()).toBe('monitor.lastUpdated(03:04:05)')
    })

    it('switches the live marker to the stale state instead of showing realtime', () => {
        const wrapper = mountContextBar({ stale: true })

        const marker = wrapper.find('.am-monitor-context__live')
        expect(marker.text()).toBe('agent.stale')
        expect(marker.classes()).toContain('am-monitor-context__live--stale')
    })

    it('renders every density option and forwards the selected density', () => {
        const wrapper = mountContextBar()

        const options = wrapper.findAllComponents(ElOptionStub)
        expect(options.map(option => option.props('value'))).toEqual([300, 600])
        expect(options.map(option => option.props('label'))).toEqual(['5分钟', '10分钟'])
        expect(wrapper.findComponent(ElSelectStub).props('modelValue')).toBe(600)

        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 300)

        expect(wrapper.emitted('update:density')).toEqual([[300]])
    })
})
