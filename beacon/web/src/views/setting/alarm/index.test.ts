import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import AlarmSettings from './index.vue'

const mocks = vi.hoisted(() => ({
    command: vi.fn(),
    queryAlarmMessageSetting: vi.fn(),
    queryAlarmThreshold: vi.fn(),
    queryDingTalk: vi.fn(),
    queryMail: vi.fn(),
    testDingTalk: vi.fn(),
    testMail: vi.fn(),
}))

vi.mock('@/api/alarm', () => ({
    queryAlarmMessageSetting: mocks.queryAlarmMessageSetting,
    queryAlarmThreshold: mocks.queryAlarmThreshold,
}))
vi.mock('@/api/dingtalk', () => ({
    queryDingTalk: mocks.queryDingTalk,
    testDingTalk: mocks.testDingTalk,
}))
vi.mock('@/api/mail', () => ({
    queryMail: mocks.queryMail,
    testMail: mocks.testMail,
}))
vi.mock('@/components/Message/message.ts', () => ({
    info: vi.fn(),
}))
vi.mock('@/hooks/useCommandComponent.ts', () => ({
    default: () => mocks.command,
}))
vi.mock('@/store', () => ({
    default: () => ({ app: { language: 'zh' } }),
}))
vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const stubs = {
    'el-button': {
        emits: ['click'],
        template: '<button type="button" @click="$emit(\'click\')"><slot /></button>',
    },
    'el-input': {
        props: ['modelValue'],
        template: '<input :value="modelValue">',
    },
    'el-tag': { template: '<span><slot /></span>' },
}

function dingTalkSetting(mask: string, plainWebhook: string) {
    return {
        id: 1,
        enabled: true,
        webhook_masked: mask,
        webhook_configured: true,
        secret_configured: true,
        at_all: false,
        // 防御性回归：即使运行期响应意外携带额外明文字段，页面也不得消费或渲染它。
        webhook: plainWebhook,
    }
}

function mountAlarmSettings() {
    return mount(AlarmSettings, {
        global: { stubs },
    })
}

describe('alarm settings DingTalk credential display', () => {
    beforeEach(() => {
        Object.values(mocks).forEach(mock => mock.mockReset())
        mocks.queryMail.mockResolvedValue({
            data: { id: 1, server: 'smtp.example.com', port: 465, sender: 'ops@example.com', receiver: 'team@example.com' },
        })
        mocks.queryAlarmMessageSetting.mockResolvedValue({ data: { message_title: 'Beacon 主机告警' } })
        mocks.queryAlarmThreshold.mockResolvedValue({ data: { data: [] } })
    })

    it('renders only the masked Webhook and refreshes the mask after saving', async () => {
        const firstPlainWebhook = 'https://oapi.dingtalk.com/robot/send?access_token=first-plain-token'
        const nextPlainWebhook = 'https://oapi.dingtalk.com/robot/send?access_token=next-plain-token'
        const firstMask = 'https://oapi.dingtalk.com/robot/send?access_token=****oken'
        const nextMask = 'https://oapi.dingtalk.com/robot/send?access_token=****next'
        mocks.queryDingTalk
            .mockResolvedValueOnce({ data: dingTalkSetting(firstMask, firstPlainWebhook) })
            .mockResolvedValueOnce({ data: dingTalkSetting(nextMask, nextPlainWebhook) })

        const wrapper = mountAlarmSettings()
        await flushPromises()

        const mask = wrapper.get('[data-testid="dingtalk-webhook-mask"]')
        expect(mask.text()).toBe(firstMask)
        expect(wrapper.html()).not.toContain(firstPlainWebhook)

        await wrapper.findAll('.settings-card')[1].find('button').trigger('click')
        const editOptions = mocks.command.mock.calls.at(-1)?.[0]
        expect(editOptions.setting).not.toHaveProperty('webhook')

        await editOptions.onSaved()
        await flushPromises()

        expect(mask.text()).toBe(nextMask)
        expect(wrapper.html()).not.toContain(nextPlainWebhook)
        expect(mocks.queryDingTalk).toHaveBeenCalledTimes(2)
    })
})
