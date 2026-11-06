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

function trendData(value: number, timestamp = 900) {
    return { data: { data: [{ timestamp, value }] } }
}

const AgentEmptyStateStub = defineComponent({
    name: 'AgentEmptyState',
    emits: ['refresh'],
    template: '<div data-testid="agent-empty" />',
})

const mountedWrappers: Array<{ unmount: () => void }> = []

function mountHostMonitor(options: { track?: boolean } = {}) {
    const wrapper = shallowMount(HostMonitor, {
        global: {
            stubs: {
                AgentEmptyState: AgentEmptyStateStub,
                DataStaleTag: { template: '<span />' },
                // 上下文条是真实组件，密度选择器需要能被测试直接驱动
                MonitorContextBar: false,
                Echarts: EchartsStub,
                echarts: EchartsStub,
                ElOption: { template: '<span />' },
                ElSelect: ElSelectStub,
            },
        },
    })
    if (options.track !== false)
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
        data: { usage: [{ device: 'disk0', data: [{ timestamp: 900, io_read: 1024, io_write: 2048 }] }] },
    })
    queryNetworkUsage.mockResolvedValue({
        data: { usage: [{ ethernet: 'eth0', data: [{ timestamp: 900, bytes_recv: 4096, bytes_sent: 8192 }] }] },
    })
})

afterEach(() => {
    mountedWrappers.splice(0).forEach(wrapper => wrapper.unmount())
})

describe('host monitor time density', () => {
    it('renders the four host charts inside the unified workspace grid', () => {
        const wrapper = mountHostMonitor()

        // 高度约束已统一到全局 .am-chart-area，host 视图复用全局 class，
        // 不再保留 host-monitor__* 局部覆盖，确保与容器监控图表高度一致。
        expect(wrapper.find('.host-monitor').exists()).toBe(true)
        expect(wrapper.findAll('.am-chart-grid')).toHaveLength(1)
        expect(wrapper.findAll('.am-chart-row')).toHaveLength(2)
        expect(wrapper.findAll('.am-chart-area')).toHaveLength(4)
        expect(wrapper.findAll('.am-chart-card')).toHaveLength(4)
        expect(wrapper.text()).not.toContain('agent.openTerminal')
    })

    it('preserves runtime axis formatters and timestamped numeric series values', async () => {
        const wrapper = mountHostMonitor()
        await flushPromises()

        const options = chartOptions(wrapper)
        expect(options).toHaveLength(4)
        expect(typeof options[2].yAxis[0].axisLabel.formatter).toBe('function')
        expect(typeof options[3].yAxis[0].axisLabel.formatter).toBe('function')
        options.forEach((option) => {
            option.series.forEach((series: { data: unknown[] }) => {
                expect(series.data.every((point) => {
                    return Array.isArray(point) && point[0] === 900_000 && typeof point[1] === 'number'
                })).toBe(true)
            })
        })
    })

    it('uses one shared five-minute window for every request and chart after density changes', async () => {
        const wrapper = mountHostMonitor()
        await flushPromises()
        clock.unixTime = 3_000

        wrapper.findComponent(ElSelectStub).vm.$emit('update:modelValue', 300)
        await flushPromises()

        const expectedRange = { agent_id: 'agent-a', start_time: 2_700, end_time: 3_000 }
        expect(queryCPUUsage).toHaveBeenLastCalledWith(expectedRange)
        expect(queryMemUsage).toHaveBeenLastCalledWith(expectedRange)
        expect(queryDiskUsage).toHaveBeenLastCalledWith(expectedRange)
        expect(queryNetworkUsage).toHaveBeenLastCalledWith(expectedRange)
        chartOptions(wrapper).forEach((option) => {
            expect(option.xAxis).toMatchObject({
                type: 'time',
                min: 2_700_000,
                max: 3_000_000,
            })
        })
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
        expect(chartOptions(wrapper)[0].series[0].data.map((point: [number, number]) => point[1])).toEqual([22])

        resolveOld(trendData(81))
        await flushPromises()
        expect(chartOptions(wrapper)[0].series[0].data.map((point: [number, number]) => point[1])).toEqual([22])
    })
})

describe('host monitor chart legends and series colors', () => {
    const tokenValues: Record<string, string> = {
        '--am-chart-1': '#0A7C79',
        '--am-chart-2': '#4FB8B4',
        '--am-chart-3': '#7C5CE0',
        '--am-chart-4': '#B45309',
        '--am-chart-5': '#047857',
        '--am-chart-6': '#E11D48',
        '--am-accent-warning': '#B45309',
        '--am-foreground-secondary': '#475569',
    }

    function stubTokenValues() {
        const original = window.getComputedStyle.bind(window)
        return vi.spyOn(window, 'getComputedStyle').mockImplementation((element) => {
            const style = original(element)
            return new Proxy(style, {
                get(target, property) {
                    if (property === 'getPropertyValue')
                        return (name: string) => tokenValues[name] ?? target.getPropertyValue(name)
                    const value = Reflect.get(target, property)
                    return typeof value === 'function' ? value.bind(target) : value
                },
            })
        })
    }

    it('names multi-series legends per design wording and colors them by chart tokens', async () => {
        const spy = stubTokenValues()
        const wrapper = mountHostMonitor()
        await flushPromises()

        const options = chartOptions(wrapper)
        const [cpu, mem, disk, net] = options

        // 单序列卡：CPU 取 chart-1，内存取 accent-warning
        expect(cpu.legend.data).toEqual(['CPU 使用率'])
        expect(cpu.series[0].lineStyle.color).toBe('#0A7C79')
        expect(mem.legend.data).toEqual(['内存使用率'])
        expect(mem.series[0].lineStyle.color).toBe('#B45309')

        // 多序列卡：图例文案按稿成文，序列色取 chart-1 / chart-5
        expect(disk.legend.data).toEqual(['读', '写'])
        expect(disk.series.map((series: { name: string }) => series.name)).toEqual(['读', '写'])
        expect(disk.series[0].lineStyle.color).toBe('#0A7C79')
        expect(disk.series[1].lineStyle.color).toBe('#047857')
        expect(disk.legend.textStyle.color).toBe('#475569')

        expect(net.legend.data).toEqual(['接收', '发送'])
        expect(net.series.map((series: { name: string }) => series.name)).toEqual(['接收', '发送'])
        expect(net.series[0].lineStyle.color).toBe('#0A7C79')
        expect(net.series[1].lineStyle.color).toBe('#047857')

        spy.mockRestore()
    })

    it('prefixes legend names with the device when multiple disks are reported', async () => {
        queryDiskUsage.mockResolvedValue({
            data: {
                usage: [
                    { device: 'disk0', data: [{ timestamp: 900, io_read: 1, io_write: 2 }] },
                    { device: 'disk1', data: [{ timestamp: 900, io_read: 3, io_write: 4 }] },
                ],
            },
        })
        const spy = stubTokenValues()
        const wrapper = mountHostMonitor()
        await flushPromises()

        const disk = chartOptions(wrapper)[2]
        expect(disk.legend.data).toEqual(['disk0 读', 'disk0 写', 'disk1 读', 'disk1 写'])
        expect(disk.series.map((series: { name: string }) => series.name)).toEqual(['disk0 读', 'disk0 写', 'disk1 读', 'disk1 写'])

        spy.mockRestore()
    })
})

describe('host monitor metric summaries', () => {
    it('summarises the latest sample under each chart title', async () => {
        const wrapper = mountHostMonitor()
        await flushPromises()

        const metrics = wrapper.findAll('.am-chart-card-percent').map(metric => metric.text())
        expect(metrics).toEqual([
            '当前 12.3%',
            '1.00 KB / 2.00 KB · 45.6%',
            '读 1.00 KB/s / 写 2.00 KB/s · 最高 -',
            '接收 4.00 KB/s · 发送 8.00 KB/s',
        ])
    })

    it('clears the disk and network series when the latest window has no samples', async () => {
        queryDiskUsage.mockResolvedValue({ data: { usage: [] } })
        queryNetworkUsage.mockResolvedValue({ data: { usage: [] } })

        const wrapper = mountHostMonitor()
        await flushPromises()

        const [, , disk, net] = chartOptions(wrapper)
        expect(disk.series).toEqual([])
        expect(disk.legend.data).toEqual([])
        expect(net.series).toEqual([])
        expect(net.legend.data).toEqual([])

        const metrics = wrapper.findAll('.am-chart-card-percent').map(metric => metric.text())
        expect(metrics[2]).toBe('读 - / 写 - · 最高 -')
        expect(metrics[3]).toBe('接收 - · 发送 -')
    })
})

describe('host monitor empty state and polling', () => {
    it('replaces the chart grid with the agent empty state when no agent is available', async () => {
        selectionState.isAgentEmpty = true

        const wrapper = mountHostMonitor()
        await flushPromises()

        expect(wrapper.find('[data-testid="agent-empty"]').exists()).toBe(true)
        expect(wrapper.findAll('.am-chart-grid')).toHaveLength(0)
    })

    it('reloads agents and refreshes the charts from the empty state action', async () => {
        selectionState.isAgentEmpty = true

        const wrapper = mountHostMonitor()
        await flushPromises()
        queryCPUUsage.mockClear()

        wrapper.findComponent(AgentEmptyStateStub).vm.$emit('refresh')
        await flushPromises()

        expect(selectionState.loadAgents).toHaveBeenCalled()
        expect(queryCPUUsage).toHaveBeenCalledTimes(1)
    })

    it('polls the agent every ten seconds and stops polling on unmount', async () => {
        const setIntervalSpy = vi.spyOn(globalThis, 'setInterval')
        const clearIntervalSpy = vi.spyOn(globalThis, 'clearInterval')

        const wrapper = mountHostMonitor({ track: false })
        await flushPromises()

        const [poll, interval] = setIntervalSpy.mock.calls.at(-1) as [() => void, number]
        expect(interval).toBe(10_000)

        queryCPUUsage.mockClear()
        poll()
        await flushPromises()
        expect(queryCPUUsage).toHaveBeenCalledTimes(1)

        wrapper.unmount()
        expect(clearIntervalSpy).toHaveBeenCalled()

        setIntervalSpy.mockRestore()
        clearIntervalSpy.mockRestore()
    })
})
