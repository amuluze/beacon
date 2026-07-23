import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { EmailSetting } from '@/interface/alarm'
import type { DingTalkSetting } from '@/interface/dingtalk'

import EditCPUThreshold from './EditCPUThreshold.vue'
import EditEmailSetting from './EditEmailSetting.vue'
import EditDingTalkSetting from './EditDingTalkSetting.vue'

const mocks = vi.hoisted(() => ({
    updateAlarmThreshold: vi.fn(),
    updateMail: vi.fn(),
    createMail: vi.fn(),
    updateDingTalk: vi.fn(),
    success: vi.fn(),
    info: vi.fn(),
}))

vi.mock('@/api/alarm', () => ({
    updateAlarmThreshold: mocks.updateAlarmThreshold,
}))
vi.mock('@/api/mail', () => ({
    createMail: mocks.createMail,
    updateMail: mocks.updateMail,
}))
vi.mock('@/api/dingtalk', () => ({
    updateDingTalk: mocks.updateDingTalk,
}))
vi.mock('@/components/Message/message.ts', () => ({
    success: mocks.success,
    info: mocks.info,
}))
vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const stubs = {
    'el-drawer': { template: '<div><slot /><slot name="footer" /></div>', props: ['modelValue', 'size', 'title'] },
    'el-dialog': { template: '<div><slot /><slot name="footer" /></div>', props: ['modelValue', 'width', 'title'] },
    'el-form': {
        name: 'ElForm',
        methods: { async validate() { return true } },
        template: '<form><slot /></form>',
    },
    'el-form-item': { template: '<label><slot /></label>' },
    'el-input': {
        props: ['modelValue'],
        emits: ['update:modelValue'],
        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)">',
    },
    'el-select': {
        props: ['modelValue'],
        template: '<div><slot /></div>',
    },
    'el-option': { template: '<span />' },
    'el-button': { template: '<button type="button" @click="$emit(\'click\')"><slot /></button>' },
    'el-switch': {
        props: ['modelValue'],
        emits: ['update:modelValue'],
        template: '<input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)">',
    },
    'el-checkbox': {
        props: ['modelValue'],
        emits: ['update:modelValue'],
        template: '<label><input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)"><slot /></label>',
    },
}

describe('setting dialogs', () => {
    beforeEach(() => {
        Object.values(mocks).forEach(mock => mock.mockReset())
    })

    it('does not mutate email props when editing and cancelling', async () => {
        const setting: EmailSetting = {
            id: 7,
            server: 'smtp.example.com',
            port: 465,
            sender: 'sender@example.com',
            password: 'secret',
            receiver: 'ops@example.com',
        }
        const original = structuredClone(setting)
        const wrapper = mount(EditEmailSetting, {
            props: { visible: true, setting },
            global: { stubs },
        })

        await wrapper.findAll('input')[0].setValue('changed.example.com')
        await wrapper.findAll('button')[0].trigger('click')

        expect(setting).toEqual(original)
        expect(mocks.updateMail).not.toHaveBeenCalled()
    })

    it('keeps the threshold dialog open and does not report success when the request fails', async () => {
        mocks.updateAlarmThreshold.mockImplementation(() => {
            throw new Error('network down')
        })
        const wrapper = mount(EditCPUThreshold, {
            props: {
                visible: true,
                threshold: { id: 1, type: 'cpu', duration: 2, threshold: 80 },
            },
            global: { stubs },
        })

        await wrapper.findAll('button')[1].trigger('click')
        await flushPromises()

        expect(mocks.success).not.toHaveBeenCalled()
        expect(wrapper.emitted('update:visible')).toBeUndefined()
    })

    it('keeps stored DingTalk credentials when credential inputs are blank', async () => {
        const setting: DingTalkSetting = {
            id: 1,
            enabled: true,
            webhook_masked: 'https://oapi.dingtalk.com/robot/send?access_token=****oken',
            webhook_configured: true,
            secret_configured: true,
            at_all: false,
        }
        const wrapper = mount(EditDingTalkSetting, {
            props: { visible: true, setting },
            global: { stubs },
        })

        await wrapper.findAll('button').at(-1)?.trigger('click')
        await flushPromises()

        expect(mocks.updateDingTalk).toHaveBeenCalledWith({
            enabled: true,
            webhook: '',
            secret: '',
            clear_secret: false,
            at_all: false,
        })
        expect(mocks.info).toHaveBeenCalledWith('setting.dingTalkUpdated')
    })
})
