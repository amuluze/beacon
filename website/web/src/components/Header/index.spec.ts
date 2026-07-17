import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import Header from './index.vue'

describe('siteHeader', () => {
    it('通过符合设计稿的菜单按钮展开移动导航', async () => {
        const wrapper = mount(Header, {
            global: {
                stubs: {
                    NuxtIcon: true,
                    NuxtLink: {
                        props: ['to'],
                        template: '<a :href="to"><slot /></a>',
                    },
                },
            },
        })
        const toggle = wrapper.get('.site-header__menu-toggle')
        const desktopNav = wrapper.get('.site-header__nav--desktop')
        const logo = wrapper.get('.site-header__mark')

        expect(logo.element.tagName).toBe('IMG')
        expect(logo.attributes('src')).toBe('/beacon.svg')
        expect(toggle.attributes('aria-expanded')).toBe('false')
        expect(wrapper.find('.site-header__mobile-menu').exists()).toBe(false)
        expect(desktopNav.text()).toContain('更新日志')
        expect(desktopNav.text()).not.toContain('团队故事')
        expect(desktopNav.get('a[href="/changelog"]').attributes('href')).toBe('/changelog')

        await toggle.trigger('click')

        expect(toggle.attributes('aria-expanded')).toBe('true')
        expect(wrapper.get('.site-header__mobile-menu').text()).toContain('使用手册')
        expect(wrapper.get('.site-header__mobile-menu').text()).toContain('更新日志')
        expect(wrapper.get('.site-header__mobile-menu').text()).not.toContain('团队故事')
    })
})
