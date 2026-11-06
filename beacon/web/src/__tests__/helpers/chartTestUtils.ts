/**
 * @Description: 监控页面测试共用工具：图表/选择器 stub、设计 token 注入与挂载回收。
 * 供 monitor 视图与其 Agent 上下文组件复用，避免每个用例重复声明。
 */

import { defineComponent } from 'vue'
import { vi } from 'vitest'

/** 与 DESIGN.md §5.1 chart-1..chart-6 对应的固定色值，用于断言序列色来源。 */
export const CHART_TOKEN_VALUES: Record<string, string> = {
    '--am-chart-1': '#0A7C79',
    '--am-chart-2': '#4FB8B4',
    '--am-chart-3': '#7C5CE0',
    '--am-chart-4': '#B45309',
    '--am-chart-5': '#047857',
    '--am-chart-6': '#E11D48',
    '--am-accent-warning': '#B45309',
    '--am-foreground-secondary': '#475569',
}

export const ElSelectStub = defineComponent({
    name: 'ElSelect',
    props: {
        modelValue: Number,
    },
    emits: ['update:modelValue'],
    template: '<div data-testid="density-select"><slot /></div>',
})

export const ElOptionStub = defineComponent({
    name: 'ElOption',
    props: {
        label: String,
        value: Number,
    },
    template: '<span />',
})

export const EchartsStub = defineComponent({
    name: 'Echarts',
    props: {
        option: {
            type: Object,
            required: true,
        },
    },
    template: '<div class="echarts-stub" />',
})

/** 用固定 token 取值替换 getComputedStyle，返回的 spy 需在用例收口处 mockRestore。 */
export function stubChartTokenValues(values: Record<string, string> = CHART_TOKEN_VALUES) {
    const original = window.getComputedStyle.bind(window)
    return vi.spyOn(window, 'getComputedStyle').mockImplementation((element) => {
        const style = original(element)
        return new Proxy(style, {
            get(target, property) {
                if (property === 'getPropertyValue')
                    return (name: string) => values[name] ?? target.getPropertyValue(name)
                const value = Reflect.get(target, property)
                return typeof value === 'function' ? value.bind(target) : value
            },
        })
    })
}

const trackedWrappers: Array<{ unmount: () => void }> = []

/** 登记挂载结果，配合 unmountTrackedWrappers 在 afterEach 中统一卸载。 */
export function trackWrapper<T extends { unmount: () => void }>(wrapper: T): T {
    trackedWrappers.push(wrapper)
    return wrapper
}

export function unmountTrackedWrappers(): void {
    trackedWrappers.splice(0).forEach(wrapper => wrapper.unmount())
}
