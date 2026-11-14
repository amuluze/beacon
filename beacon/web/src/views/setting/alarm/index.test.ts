import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import AlarmSettings from './index.vue'

const mocks = vi.hoisted(() => ({
    command: vi.fn(),
    queryAlarmMessageSetting: vi.fn(),
    queryAlarmThreshold: vi.fn(),
    queryDingTalk: vi.fn(),
    updateDingTalk: vi.fn(),
    queryFeishu: vi.fn(),
    updateFeishu: vi.fn(),
    testFeishu: vi.fn(),
    queryWeCom: vi.fn(),
    updateWeCom: vi.fn(),
    testWeCom: vi.fn(),
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
    updateDingTalk: mocks.updateDingTalk,
    testDingTalk: mocks.testDingTalk,
}))
vi.mock('@/api/feishu', () => ({
    queryFeishu: mocks.queryFeishu,
    updateFeishu: mocks.updateFeishu,
    testFeishu: mocks.testFeishu,
}))
vi.mock('@/api/wecom', () => ({
    queryWeCom: mocks.queryWeCom,
    updateWeCom: mocks.updateWeCom,
    testWeCom: mocks.testWeCom,
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
        emits: ['update:modelValue'],
        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)">',
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
        mocks.queryFeishu.mockResolvedValue({
            data: { id: 0, enabled: false, webhook_masked: '', webhook_configured: false, secret_configured: false },
        })
        mocks.queryWeCom.mockResolvedValue({
            data: { id: 0, enabled: false, webhook_masked: '', webhook_configured: false, at_all: false },
        })
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

        // 内联表单：输入新 Webhook 并保存（留空即保留语义的对面路径），
        // 保存成功后重新拉取配置刷新遮罩。
        await wrapper.get('[data-testid="dingtalk-webhook-input"]').setValue(nextPlainWebhook)
        await wrapper.get('[data-testid="dingtalk-save"]').trigger('click')
        await flushPromises()

        expect(mocks.updateDingTalk).toHaveBeenCalledTimes(1)
        expect(mocks.updateDingTalk).toHaveBeenCalledWith(expect.objectContaining({ webhook: nextPlainWebhook }))
        expect(mask.text()).toBe(nextMask)
        expect(wrapper.html()).not.toContain(nextPlainWebhook)
        expect(mocks.queryDingTalk).toHaveBeenCalledTimes(2)
    })

    it('switches to the Feishu channel and saves its inline form', async () => {
        mocks.queryFeishu.mockResolvedValue({
            data: { id: 1, enabled: false, webhook_masked: 'https://open.feishu.cn/open-apis/bot/v2/hook/****ook', webhook_configured: true, secret_configured: false },
        })
        const wrapper = mountAlarmSettings()
        await flushPromises()

        // 渠道单选切到飞书：钉钉表单卸载，飞书表单出现。
        const radios = wrapper.findAll('input[type="radio"]')
        await radios[1].setValue(true)
        await flushPromises()
        expect(wrapper.find('[data-testid="dingtalk-webhook-input"]').exists()).toBe(false)

        const feishuInput = wrapper.get('[data-testid="feishu-webhook-input"]')
        await feishuInput.setValue('https://open.feishu.cn/open-apis/bot/v2/hook/new-hook-id')
        await wrapper.get('[data-testid="feishu-save"]').trigger('click')
        await flushPromises()

        expect(mocks.updateFeishu).toHaveBeenCalledTimes(1)
        expect(mocks.updateFeishu).toHaveBeenCalledWith(expect.objectContaining({
            enabled: false,
            webhook: 'https://open.feishu.cn/open-apis/bot/v2/hook/new-hook-id',
        }))
    })

    it('blocks enabling the WeCom channel without a webhook', async () => {
        const wrapper = mountAlarmSettings()
        await flushPromises()

        const radios = wrapper.findAll('input[type="radio"]')
        await radios[2].setValue(true)
        await flushPromises()

        // 未配置 Webhook 时启用开关 + 直接保存：前端拦截且不发起请求。
        await wrapper.get('.el-switch').trigger('click')
        await wrapper.get('[data-testid="wecom-save"]').trigger('click')
        await flushPromises()

        expect(mocks.updateWeCom).not.toHaveBeenCalled()
    })
})
