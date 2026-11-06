import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { afterEach, describe, expect, it, vi } from 'vitest'
import DataStaleTag from './DataStaleTag.vue'
import { trackWrapper, unmountTrackedWrappers } from '@/__tests__/helpers/chartTestUtils'

const translator = vi.hoisted(() => vi.fn((key: string, params?: Record<string, unknown>) => {
    return params ? `${key}(${Object.values(params).join(',')})` : key
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: translator }),
}))

const ElTooltipStub = defineComponent({
    name: 'ElTooltip',
    props: {
        content: String,
        placement: String,
    },
    template: '<div data-testid="stale-tooltip" :data-content="content"><slot /></div>',
})

const ElTagStub = defineComponent({
    name: 'ElTag',
    props: {
        type: String,
        effect: String,
        size: String,
    },
    template: '<span data-testid="stale-tag"><slot /></span>',
})

function mountStaleTag(props: { stale?: boolean, timestamp?: number } = {}) {
    const wrapper = mount(DataStaleTag, {
        props,
        global: {
            stubs: {
                'el-tooltip': ElTooltipStub,
                'ElTooltip': ElTooltipStub,
                'el-tag': ElTagStub,
                'ElTag': ElTagStub,
            },
        },
    })
    return trackWrapper(wrapper)
}

afterEach(() => {
    unmountTrackedWrappers()
})

describe('data stale tag', () => {
    it('renders nothing while the sample is fresh', () => {
        const wrapper = mountStaleTag({ stale: false, timestamp: 100 })

        expect(wrapper.find('[data-testid="stale-tag"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="stale-tooltip"]').exists()).toBe(false)
    })

    it('marks stale data with a warning tag', () => {
        const wrapper = mountStaleTag({ stale: true, timestamp: 100 })

        const tag = wrapper.findComponent(ElTagStub)
        expect(tag.exists()).toBe(true)
        expect(tag.text()).toBe('agent.stale')
        expect(tag.props()).toMatchObject({ type: 'warning', effect: 'plain', size: 'small' })
    })

    it('explains the stale sample with the default tooltip when no timestamp is reported', () => {
        const wrapper = mountStaleTag({ stale: true, timestamp: 0 })

        expect(wrapper.findComponent(ElTooltipStub).props('content')).toBe('agent.staleTooltipDefault')
    })

    it('formats the last sample time inside the stale tooltip', () => {
        const timestamp = new Date(2024, 0, 2, 3, 4, 5).getTime() / 1000

        const wrapper = mountStaleTag({ stale: true, timestamp })

        expect(wrapper.findComponent(ElTooltipStub).props('content')).toBe('agent.staleTooltipWithTime(2024-01-02 03:04:05)')
    })
})
