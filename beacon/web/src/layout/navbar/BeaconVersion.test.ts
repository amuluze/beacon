import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import BeaconVersion from './BeaconVersion.vue'

const mocks = vi.hoisted(() => ({
    querySystemVersion: vi.fn(),
}))

vi.mock('@/api/system', () => ({
    querySystemVersion: mocks.querySystemVersion,
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const stubs = {
    'svg-icon': true,
    'el-dialog': true,
}

describe('navbar beacon version trigger', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('keeps the upgrade icon disabled when no newer release is known', async () => {
        mocks.querySystemVersion.mockResolvedValue({ data: { version: 'v3.0.15', update_available: false } })

        const wrapper = mount(BeaconVersion, { global: { stubs } })
        await vi.waitFor(() => expect(wrapper.find('button').exists()).toBe(true))

        expect(wrapper.find('button').attributes('disabled')).toBeDefined()
        expect(wrapper.find('button').classes()).not.toContain('am-version-trigger--update')
    })

    it('highlights the upgrade icon and opens the detail dialog when a newer release exists', async () => {
        mocks.querySystemVersion.mockResolvedValue({ data: { version: 'v3.0.15', update_available: true, latest_version: 'v3.1.0' } })

        const wrapper = mount(BeaconVersion, { global: { stubs } })
        await vi.waitFor(() => expect(wrapper.find('button').classes()).toContain('am-version-trigger--update'))

        expect(wrapper.find('button').attributes('disabled')).toBeUndefined()

        await wrapper.find('button').trigger('click')

        expect(wrapper.findComponent({ name: 'ElDialog' }).props('modelValue')).toBe(true)
    })

    it('degrades to a disabled trigger when the version query fails', async () => {
        mocks.querySystemVersion.mockRejectedValue(new Error('offline'))

        const wrapper = mount(BeaconVersion, { global: { stubs } })
        await vi.waitFor(() => expect(wrapper.find('button').exists()).toBe(true))

        expect(wrapper.find('button').attributes('disabled')).toBeDefined()
    })
})
