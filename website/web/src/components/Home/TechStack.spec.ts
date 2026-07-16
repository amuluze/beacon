import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import TechStack from './TechStack.vue'

describe('techStack', () => {
    it('渲染 4 个技术栈卡片（shallow 隔离 Icon 等 Nuxt 子组件）', () => {
        const wrapper = mount(TechStack, { shallow: true })
        expect(wrapper.findAll('.tech-stack__card')).toHaveLength(4)
    })

    it('包含关键栈项', () => {
        const wrapper = mount(TechStack, { shallow: true })
        const text = wrapper.text()
        expect(text).toContain('Golang')
        expect(text).toContain('Vue 3')
        expect(text).toContain('SQLite')
    })
})
