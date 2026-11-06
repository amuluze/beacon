import { flushPromises, shallowMount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ContainerMonitor from './index.vue'

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
    queryContainersUsage: vi.fn(),
}))

const { queryContainersUsage } = apiMocks

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

vi.mock('@/api/container', () => apiMocks)

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const ElSelectStub = defineComponent({
    name: 'ElSelect',
    props: {
        modelValue: Number,
    },
    emits: ['update:modelValue'],
    template: '<div data-testid="density-select"><slot /></div>',
})

const ElOptionStub = defineComponent({
    name: 'ElOption',
    props: {
        label: String,
        value: Number,
    },
    template: '<span />',
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

function usageData(name: string, cpuValue: number, memValue: number, timestamp = 900) {
    return {
        data: {
            names: [name],
            cpu_usage: { [name]: [{ timestamp, value: cpuValue }] },
            mem_usage: { [name]: [{ timestamp, value: memValue }] },
        },
    }
}

function emptyUsageData() {
    return {
        data: {
            names: [],
            cpu_usage: {},
            mem_usage: {},
        },
    }
}

function manyContainersUsage(names: string[]) {
    return {
        data: {
            names,
            cpu_usage: Object.fromEntries(names.map((name, index) => [name, [{ timestamp: 900, value: index + 1 }]])),
            mem_usage: Object.fromEntries(names.map((name, index) => [name, [{ timestamp: 900, value: (index + 1) * 1024 }]])),
        },
    }
}

const AgentEmptyStateStub = defineComponent({
    name: 'AgentEmptyState',
    emits: ['refresh'],
    template: '<div data-testid="agent-empty" />',
})

const mountedWrappers: Array<{ unmount: () => void }> = []

function mountContainerMonitor(options: { track?: boolean } = {}) {
    const wrapper = shallowMount(ContainerMonitor, {
        global: {
            stubs: {
                AgentEmptyState: AgentEmptyStateStub,
                // 上下文条是真实组件，密度选择器需要能被测试直接驱动
                MonitorContextBar: false,
                Echarts: EchartsStub,
                echarts: EchartsStub,
                ElOption: ElOptionStub,
                ElSelect: ElSelectStub,
            },
        },
    })
    if (options.track !== false)
        mountedWrappers.push(wrapper)
    return wrapper
}

function chartOptions(wrapper: ReturnType<typeof mountContainerMonitor>): Record<string, any>[] {
    const charts = wrapper.findAllComponents({ name: 'Index' })
    return charts.map(chart => chart.props('option') as Record<string, any>)
}

beforeEach(() => {
    vi.clearAllMocks()
    clock.unixTime = 1_000
    selectionState.selectedAgentID = 'agent-a'
    selectionState.isAgentEmpty = false
    queryContainersUsage.mockResolvedValue(usageData('container-a', 12.3, 2_048))
})

afterEach(() => {
    mountedWrappers.splice(0).forEach(wrapper => wrapper.unmount())
})

describe('container monitor time density', () => {
    it('offers and renders one shared five-minute window after density changes', async () => {
        const wrapper = mountContainerMonitor()
        await flushPromises()
        expect(wrapper.findAllComponents(ElOptionStub).map(option => option.props('value'))).toContain(300)
        clock.unixTime = 3_000

        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 300)
        await flushPromises()

        expect(queryContainersUsage).toHaveBeenLastCalledWith({
            agent_id: 'agent-a',
            start_time: 2_700,
            end_time: 3_000,
        })
        chartOptions(wrapper).forEach((option) => {
            expect(option.xAxis).toMatchObject({
                type: 'time',
                min: 2_700_000,
                max: 3_000_000,
            })
        })
    })

    it('ignores a stale response from the previous density', async () => {
        let resolveOld!: (value: ReturnType<typeof usageData>) => void
        let resolveCurrent!: (value: ReturnType<typeof usageData>) => void
        queryContainersUsage
            .mockReset()
            .mockImplementationOnce(async () => new Promise(resolve => resolveOld = resolve))
            .mockImplementationOnce(async () => new Promise(resolve => resolveCurrent = resolve))

        const wrapper = mountContainerMonitor()
        await flushPromises()
        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 1_800)
        await flushPromises()

        resolveCurrent(usageData('current-container', 22, 4_096))
        await flushPromises()
        expect(chartOptions(wrapper)[0].series[0].data).toEqual([[900_000, 22]])
        expect(chartOptions(wrapper)[1].series[0].data).toEqual([[900_000, 4_096]])

        resolveOld(usageData('old-container', 81, 8_192))
        await flushPromises()
        expect(chartOptions(wrapper)[0].series[0].data).toEqual([[900_000, 22]])
        expect(chartOptions(wrapper)[1].series[0].data).toEqual([[900_000, 4_096]])
        expect(wrapper.findAll('.am-legend-label').map(label => label.text())).toEqual(['current-container', 'current-container'])
    })

    it('clears old chart state when the latest window has no data', async () => {
        const wrapper = mountContainerMonitor()
        await flushPromises()
        expect(wrapper.findAll('.am-legend-label')).toHaveLength(2)

        queryContainersUsage.mockResolvedValueOnce(emptyUsageData())
        clock.unixTime = 3_000
        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 1_800)
        await flushPromises()

        expect(wrapper.findAll('.am-legend-label')).toHaveLength(0)
        chartOptions(wrapper).forEach((option) => {
            expect(option.legend.data).toEqual([])
            expect(option.xAxis).toMatchObject({
                type: 'time',
                min: 1_200_000,
                max: 3_000_000,
            })
            expect(option.series).toEqual([])
        })
    })

    it('keeps high-cardinality legends scrollable without duplicating them inside charts', async () => {
        const names = Array.from({ length: 30 }, (_, index) => index === 29
            ? 'container-with-a-very-long-production-name-that-needs-truncation'
            : `container-${index + 1}`)
        queryContainersUsage.mockResolvedValueOnce(manyContainersUsage(names))

        const wrapper = mountContainerMonitor()
        await flushPromises()

        const legends = wrapper.findAll('.am-chart-card-legend--scrollable')
        expect(legends).toHaveLength(2)
        expect(legends.map(legend => legend.attributes('role'))).toEqual(['list', 'list'])
        expect(legends.map(legend => legend.findAll('.am-legend-item').length)).toEqual([30, 30])

        const items = wrapper.findAll('.am-legend-item')
        expect(items).toHaveLength(60)
        expect(items[29].attributes('title')).toBe(names[29])
        expect(items[59].attributes('title')).toBe(names[29])

        chartOptions(wrapper).forEach((option) => {
            expect(option.legend).toMatchObject({ show: false, data: [] })
            expect(option.series).toHaveLength(30)
            expect(option.series.map((series: { name: string }) => series.name)).toEqual(names)
        })
    })
})

describe('container monitor metric chips', () => {
    it('shows the latest cpu and memory percentage per container', async () => {
        queryContainersUsage.mockResolvedValueOnce({
            data: {
                names: ['app'],
                cpu_usage: { app: [{ timestamp: 900, value: 12.3 }] },
                mem_usage: { app: [{ timestamp: 900, value: 2_048 }] },
                mem_percent: { app: [{ timestamp: 900, value: 10.04 }, { timestamp: 960, value: 42.36 }] },
            },
        })

        const wrapper = mountContainerMonitor()
        await flushPromises()

        const values = wrapper.findAll('.am-legend-value').map(value => value.text())
        expect(values).toEqual(['12.3%', '42.4%'])
    })
})

describe('container monitor empty state and polling', () => {
    it('replaces the chart card with the agent empty state when no agent is available', async () => {
        selectionState.isAgentEmpty = true

        const wrapper = mountContainerMonitor()
        await flushPromises()

        expect(wrapper.find('[data-testid="agent-empty"]').exists()).toBe(true)
        expect(wrapper.findAll('.am-chart-grid')).toHaveLength(0)
    })

    it('reloads agents and re-renders from the empty state action', async () => {
        selectionState.isAgentEmpty = true

        const wrapper = mountContainerMonitor()
        await flushPromises()
        queryContainersUsage.mockClear()

        wrapper.findComponent(AgentEmptyStateStub).vm.$emit('refresh')
        await flushPromises()

        expect(selectionState.loadAgents).toHaveBeenCalled()
        expect(queryContainersUsage).toHaveBeenCalledTimes(1)
    })

    it('polls every ten seconds and stops polling on unmount', async () => {
        const setIntervalSpy = vi.spyOn(globalThis, 'setInterval')
        const clearIntervalSpy = vi.spyOn(globalThis, 'clearInterval')

        const wrapper = mountContainerMonitor({ track: false })
        await flushPromises()

        const [poll, interval] = setIntervalSpy.mock.calls.at(-1) as [() => void, number]
        expect(interval).toBe(10_000)

        queryContainersUsage.mockClear()
        poll()
        await flushPromises()
        expect(queryContainersUsage).toHaveBeenCalledTimes(1)

        wrapper.unmount()
        expect(clearIntervalSpy).toHaveBeenCalled()

        setIntervalSpy.mockRestore()
        clearIntervalSpy.mockRestore()
    })
})
