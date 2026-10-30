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

import ConfirmDialog from '@/components/ConfirmDialog/index.vue'

const { tMock } = vi.hoisted(() => ({
    tMock: vi.fn((key: string, params?: Record<string, unknown>) => {
        if (params && Object.keys(params).length > 0)
            return `${key}:${JSON.stringify(params)}`
        return key
    }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: tMock }),
}))

describe('confirmDialog', () => {
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

    it('interpolates message params for batch confirmations', async () => {
        tMock.mockClear()
        const wrapper = mount(ConfirmDialog, {
            props: {
                visible: true,
                title: 'container.batchDelete',
                message: 'container.confirmBatchDelete',
                params: { count: 3 },
                i18nPrefix: 'container',
            },
            global: {
                stubs: {
                    'el-dialog': {
                        template: '<div class="el-dialog-stub"><slot /></div>',
                        props: ['title', 'modelValue', 'width', 'draggable'],
                    },
                },
                mocks: { $t: (k: string) => k },
            },
        })
        await nextTick()

        expect(tMock).toHaveBeenCalledWith('container.confirmBatchDelete', { count: 3 })
        expect(wrapper.text()).toContain('container.confirmBatchDelete:{"count":3}')
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
        const inst = wrapper.vm
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
        const inst = wrapper.vm
        expect(typeof inst.confirm).toBe('function')
        inst.confirm()
        await nextTick()
        expect(wrapper.emitted('confirm')).toBeTruthy()
    })
})
