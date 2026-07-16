import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import Header from './index.vue'

describe('siteHeader', () => {
    it('通过符合设计稿的菜单按钮展开移动导航', async () => {
        const wrapper = mount(Header, {
            global: {
                stubs: {
                    NuxtIcon: true,
                    NuxtLink: { template: '<a><slot /></a>' },
                },
            },
        })
        const toggle = wrapper.get('.site-header__menu-toggle')

        expect(toggle.attributes('aria-expanded')).toBe('false')
        expect(wrapper.find('.site-header__mobile-menu').exists()).toBe(false)

        await toggle.trigger('click')

        expect(toggle.attributes('aria-expanded')).toBe('true')
        expect(wrapper.get('.site-header__mobile-menu').text()).toContain('使用手册')
        expect(wrapper.get('.site-header__mobile-menu').text()).toContain('更新日志')
    })
})
