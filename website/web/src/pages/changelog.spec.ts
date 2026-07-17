import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

import ChangelogPage from './changelog.vue'

vi.mock('~/composables/usePageSeo', () => ({
    usePageSeo: vi.fn(),
}))

describe('changelogPage', () => {
    it('按版本展示带类型标签的更新条目', () => {
        const wrapper = mount(ChangelogPage, {
            global: {
                stubs: {
                    NuxtIcon: true,
                },
            },
        })
        const types = wrapper.findAll('.release-change__type').map(node => node.text())

        expect(wrapper.findAll('.release-change')).toHaveLength(18)
        expect(types).toContain('新功能')
        expect(types).toContain('改进')
        expect(types).toContain('安全')
        expect(types).toContain('修复')
    })
})
