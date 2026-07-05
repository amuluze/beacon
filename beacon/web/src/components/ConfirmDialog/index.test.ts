/**
 * Tests for the ConfirmDialog reusable component.
 *
 * Setup:
 *  - Mount via Vue Test Utils inside the jsdom environment defined by
 *    vitest.config.ts.
 *  - Provide minimal i18n stubs by mocking `vue-i18n` so titles and labels
 *    resolve to their key (sufficient for behavior coverage).
 */
import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

import ConfirmDialog from '@/components/ConfirmDialog/index.vue'

describe('ConfirmDialog', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('renders title and message via i18n', async () => {
        const wrapper = mount(ConfirmDialog, {
            props: {
                visible: true,
                title: 'container.deleteContainer',
                message: 'container.confirmDelete',
                i18nPrefix: 'container',
            },
            global: {
                stubs: {
                    'el-dialog': {
                        template: '<div class="el-dialog-stub"><header class="el-dialog__title">{{ title }}</header><slot /></div>',
                        props: ['title', 'modelValue', 'width', 'draggable'],
                    },
                    'el-button': {
                        template: '<button @click="$emit(\'click\')"><slot /></button>',
                    },
                },
                mocks: { $t: (k: string) => k },
            },
        })
        await nextTick()
        expect(wrapper.text()).toContain('container.confirmDelete')
        expect(wrapper.find('.el-dialog__title').text()).toBe('container.deleteContainer')
    })

    it('emits update:visible(false) when close() is invoked', async () => {
        const wrapper = mount(ConfirmDialog, {
            props: { visible: true },
            global: {
                stubs: {
                    'el-dialog': { template: '<div><slot /></div>', props: ['title', 'modelValue', 'width', 'draggable'] },
                },
                mocks: { $t: (k: string) => k },
            },
        })
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const inst = wrapper.vm as any
        inst.close()
        await nextTick()
        const events = wrapper.emitted('update:visible')
        expect(events).toBeTruthy()
        expect(events![0]).toEqual([false])
        expect(wrapper.emitted('close')).toBeTruthy()
    })

    it('emits confirm event when onConfirm is invoked from outside', async () => {
        // The ConfirmDialog exposes its confirm() handler through defineExpose.
        // Testing the button click path requires element-plus to be registered
        // globally; we cover the same behavior through the exposed method.
        const wrapper = mount(ConfirmDialog, {
            props: {
                visible: true,
                title: 'common.action',
                confirmation: true,
            },
            global: {
                stubs: {
                    'el-dialog': { template: '<div><slot /></div>', props: ['title', 'modelValue', 'width', 'draggable'] },
                    'el-button': { template: '<button></button>' },
                },
                mocks: { $t: (k: string) => k },
            },
        })
        await nextTick()
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const inst = wrapper.vm as any
        expect(typeof inst.confirm).toBe('function')
        inst.confirm()
        await nextTick()
        expect(wrapper.emitted('confirm')).toBeTruthy()
    })
})
