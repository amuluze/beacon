import { flushPromises, shallowMount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import HostMonitor from './index.vue'

const clock = vi.hoisted(() => ({ unixTime: 1_000 }))

const selectionState = vi.hoisted(() => {
    const state = {
        selectedAgentID: 'agent-a',
        isAgentEmpty: false,
        ensureSelectedAgent: vi.fn<() => Promise<string>>(),
        loadAgents: vi.fn<() => Promise<string>>(),
    }
    state.ensureSelectedAgent.mockImplementation(async () => state.selectedAgentID)
    state.loadAgents.mockImplementation(async () => state.selectedAgentID)
    return state
})

const apiMocks = vi.hoisted(() => ({
    queryCPUInfo: vi.fn(),
    queryCPUUsage: vi.fn(),
    queryDiskInfo: vi.fn(),
    queryDiskUsage: vi.fn(),
    queryMemInfo: vi.fn(),
    queryMemUsage: vi.fn(),
    queryNetworkUsage: vi.fn(),
}))

const {
    queryCPUInfo,
    queryCPUUsage,
    queryDiskInfo,
    queryDiskUsage,
    queryMemInfo,
    queryMemUsage,
    queryNetworkUsage,
} = apiMocks

vi.mock('element-plus', async (importOriginal) => {
    const actual = await importOriginal<typeof import('element-plus')>()
    return {
        ...actual,
        dayjs: (value?: number) => ({
            format: () => String(value ?? ''),
            unix: () => clock.unixTime++,
        }),
    }
})

vi.mock('@/hooks/useAgentSelection', async () => {
    const { computed } = await import('vue')
    return {
        useAgentSelection: () => ({
            selectedAgentID: computed({
                get: () => selectionState.selectedAgentID,
                set: value => selectionState.selectedAgentID = value,
            }),
            isAgentEmpty: computed(() => selectionState.isAgentEmpty),
            agentParams: computed(() => selectionState.selectedAgentID ? { agent_id: selectionState.selectedAgentID } : {}),
            ensureSelectedAgent: selectionState.ensureSelectedAgent,
            loadAgents: selectionState.loadAgents,
        }),
    }
})

vi.mock('@/api/host', () => apiMocks)

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

vi.mock('vue-router', () => ({
    useRouter: () => ({ push: vi.fn() }),
}))

const ElSelectStub = defineComponent({
    name: 'ElSelect',
    props: {
        modelValue: Number,
    },
    emits: ['update:modelValue'],
    template: '<div data-testid="density-select"><slot /></div>',
})

const EchartsStub = defineComponent({
    name: 'Echarts',
    props: {
        option: {
            type: Object,
            required: true,
        },
    },
    template: '<div class="echarts-stub" />',
})

function trendData(value: number) {
    return { data: { data: [{ timestamp: 100, value }] } }
}

const mountedWrappers: Array<{ unmount: () => void }> = []

function mountHostMonitor() {
    const wrapper = shallowMount(HostMonitor, {
        global: {
            stubs: {
                AgentEmptyState: { template: '<div data-testid="agent-empty" />' },
                DataStaleTag: { template: '<span />' },
                Echarts: EchartsStub,
                echarts: EchartsStub,
                ElButton: { template: '<button><slot /></button>' },
                ElOption: { template: '<span />' },
                ElSelect: ElSelectStub,
            },
        },
    })
    mountedWrappers.push(wrapper)
    return wrapper
}

function chartOptions(wrapper: ReturnType<typeof mountHostMonitor>): Record<string, any>[] {
    const charts = wrapper.findAllComponents({ name: 'Index' })
    return charts.map(chart => chart.props('option') as Record<string, any>)
}

beforeEach(() => {
    vi.clearAllMocks()
    clock.unixTime = 1_000
    selectionState.selectedAgentID = 'agent-a'
    selectionState.isAgentEmpty = false

    queryCPUInfo.mockResolvedValue({ data: { percent: 12.3, stale: false, timestamp: 100 } })
    queryCPUUsage.mockResolvedValue(trendData(12.3))
    queryMemInfo.mockResolvedValue({ data: { percent: 45.6, total: 2048, used: 1024, stale: false, timestamp: 100 } })
    queryMemUsage.mockResolvedValue(trendData(45.6))
    queryDiskInfo.mockResolvedValue({ data: { info: [] } })
    queryDiskUsage.mockResolvedValue({
        data: { usage: [{ device: 'disk0', data: [{ timestamp: 100, io_read: 1024, io_write: 2048 }] }] },
    })
    queryNetworkUsage.mockResolvedValue({
        data: { usage: [{ ethernet: 'eth0', data: [{ timestamp: 100, bytes_recv: 4096, bytes_sent: 8192 }] }] },
    })
})

afterEach(() => {
    mountedWrappers.splice(0).forEach(wrapper => wrapper.unmount())
})

describe('host monitor time density', () => {
    it('uses a dedicated taller plot area for every host chart', () => {
        const wrapper = mountHostMonitor()

        expect(wrapper.find('.host-monitor').exists()).toBe(true)
        expect(wrapper.find('.host-monitor__chart-grid').exists()).toBe(true)
        expect(wrapper.findAll('.host-monitor__chart-row')).toHaveLength(2)
        expect(wrapper.findAll('.host-monitor__chart-area')).toHaveLength(4)
    })

    it('preserves runtime axis formatters and numeric series values', async () => {
        const wrapper = mountHostMonitor()
        await flushPromises()

        const options = chartOptions(wrapper)
        expect(options).toHaveLength(4)
        expect(typeof options[2].yAxis[0].axisLabel.formatter).toBe('function')
        expect(typeof options[3].yAxis[0].axisLabel.formatter).toBe('function')
        options.forEach((option) => {
            option.series.forEach((series: { data: unknown[] }) => {
                expect(series.data.every(value => typeof value === 'number')).toBe(true)
            })
        })
    })

    it('uses one shared time window for every trend request after density changes', async () => {
        const wrapper = mountHostMonitor()
        await flushPromises()
        clock.unixTime = 3_000

        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 120)
        await flushPromises()

        const expectedRange = { agent_id: 'agent-a', start_time: 2_880, end_time: 3_000 }
        expect(queryCPUUsage).toHaveBeenLastCalledWith(expectedRange)
        expect(queryMemUsage).toHaveBeenLastCalledWith(expectedRange)
        expect(queryDiskUsage).toHaveBeenLastCalledWith(expectedRange)
        expect(queryNetworkUsage).toHaveBeenLastCalledWith(expectedRange)
    })

    it('ignores a stale response from the previous density', async () => {
        let resolveOld!: (value: ReturnType<typeof trendData>) => void
        let resolveCurrent!: (value: ReturnType<typeof trendData>) => void
        queryCPUUsage
            .mockReset()
            .mockImplementationOnce(async () => new Promise(resolve => resolveOld = resolve))
            .mockImplementationOnce(async () => new Promise(resolve => resolveCurrent = resolve))

        const wrapper = mountHostMonitor()
        await flushPromises()
        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 120)
        await flushPromises()

        resolveCurrent(trendData(22))
        await flushPromises()
        expect(chartOptions(wrapper)[0].series[0].data.map(Number)).toEqual([22])

        resolveOld(trendData(81))
        await flushPromises()
        expect(chartOptions(wrapper)[0].series[0].data.map(Number)).toEqual([22])
    })
})
