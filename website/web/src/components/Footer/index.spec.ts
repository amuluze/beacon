import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import Footer from './index.vue'

describe('siteFooter', () => {
    it('以产品资源和相关作品组织链接，并移除团队故事', () => {
        const wrapper = mount(Footer, {
            global: {
                stubs: {
                    NuxtLink: {
                        props: ['to'],
                        template: '<a :href="to"><slot /></a>',
                    },
                },
            },
        })
        const text = wrapper.text()
        const logo = wrapper.get('.site-footer__brand-logo')

        expect(logo.element.tagName).toBe('IMG')
        expect(logo.attributes('src')).toBe('/beacon.svg')
        expect(text).toContain('产品资源')
        expect(text).toContain('相关作品')
        expect(text).toContain('更新日志')
        expect(text).not.toContain('关于我们')
        expect(text).not.toContain('友情链接')
        expect(text).not.toContain('团队故事')
        expect(wrapper.get('a[href="/changelog"]').attributes('href')).toBe('/changelog')
    })
})
