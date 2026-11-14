<script setup lang="ts">
import { queryMail, testMail } from '@/api/mail'
import { queryDingTalk, testDingTalk, updateDingTalk } from '@/api/dingtalk'
import { queryFeishu, testFeishu, updateFeishu } from '@/api/feishu'
import { queryWeCom, testWeCom, updateWeCom } from '@/api/wecom'
import { info, warning } from '@/components/Message/message.ts'

import type { AlarmMessageSetting, AlarmThreshold, EmailSetting } from '@/interface/alarm.ts'
import type { DingTalkSetting } from '@/interface/dingtalk'
import type { FeishuSetting, FeishuUpdateArgs } from '@/interface/feishu'
import type { MailTestArgs } from '@/interface/mail.ts'
import type { WeComSetting, WeComUpdateArgs } from '@/interface/wecom'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import EditEmailSetting from '@/views/setting/alarm/components/EditEmailSetting.vue'
import EditCPUThreshold from '@/views/setting/alarm/components/EditCPUThreshold.vue'
import EditMemThreshold from '@/views/setting/alarm/components/EditMemThreshold.vue'
import EditDiskThreshold from '@/views/setting/alarm/components/EditDiskThreshold.vue'
import EditAlarmMessageTitle from '@/views/setting/alarm/components/EditAlarmMessageTitle.vue'
import { queryAlarmMessageSetting, queryAlarmThreshold } from '@/api/alarm'
import { useI18n } from 'vue-i18n'
import useStore from '@/store'
import IconBellRing from '~icons/lucide/bell-ring'
import IconGauge from '~icons/lucide/gauge'
import IconMailCheck from '~icons/lucide/mail-check'
import IconPencil from '~icons/lucide/pencil'
import IconWebhook from '~icons/lucide/webhook'

const { t } = useI18n()
const store = useStore()
const locale = computed(() => store.app.language)

const emailSetting = reactive<EmailSetting>({
  id: 0,
  server: '',
  port: 465,
  sender: '',
  password: '******',
  receiver: '',
})

const editEmailSetting = useCommandComponent(EditEmailSetting)

const dingTalkSetting = reactive<DingTalkSetting>({
  id: 0,
  enabled: false,
  webhook_masked: '',
  webhook_configured: false,
  secret_configured: false,
  at_all: false,
})
const dingTalkTesting = shallowRef(false)

// 群机器人卡片内联表单（参考通知管理形态）：webhook/secret 保持「留空即保留」
// 的服务端语义，已配置时以遮罩提示展示当前凭据。飞书/企业微信渠道服务端尚未
// 支持，先以禁用态呈现。
const dingTalkChannel = ref('dingtalk')
const dingTalkForm = reactive({
  enabled: false,
  webhook: '',
  secret: '',
  clear_secret: false,
  at_all: false,
})
const dingTalkSaving = shallowRef(false)

watch(
  () => ({
    enabled: dingTalkSetting.enabled,
    at_all: dingTalkSetting.at_all,
    webhook_masked: dingTalkSetting.webhook_masked,
    secret_configured: dingTalkSetting.secret_configured,
  }),
  () => {
    dingTalkForm.enabled = dingTalkSetting.enabled
    dingTalkForm.at_all = dingTalkSetting.at_all
    dingTalkForm.webhook = ''
    dingTalkForm.secret = ''
    dingTalkForm.clear_secret = false
  },
)

async function saveDingTalk() {
  if (dingTalkForm.enabled && !dingTalkSetting.webhook_configured && !dingTalkForm.webhook.trim()) {
    warning(t('setting.dingTalkWebhookRequired'))
    return
  }
  dingTalkSaving.value = true
  try {
    await updateDingTalk({
      enabled: dingTalkForm.enabled,
      webhook: dingTalkForm.webhook.trim(),
      secret: dingTalkForm.secret.trim(),
      clear_secret: dingTalkForm.clear_secret,
      at_all: dingTalkForm.at_all,
    })
    info(t('setting.dingTalkUpdated'))
    await loadDingTalkSetting()
  }
  catch {
    // 请求拦截器负责展示服务端错误；保留本地表单供用户重试。
  }
  finally {
    dingTalkSaving.value = false
  }
}

const feishuSetting = reactive<FeishuSetting>({
  id: 0,
  enabled: false,
  webhook_masked: '',
  webhook_configured: false,
  secret_configured: false,
})
const feishuForm = reactive<FeishuUpdateArgs>({
  enabled: false,
  webhook: '',
  secret: '',
  clear_secret: false,
})
const feishuSaving = shallowRef(false)
const feishuTestedAt = shallowRef('')

watch(() => ({
  enabled: feishuSetting.enabled,
  secret_configured: feishuSetting.secret_configured,
}), () => {
  feishuForm.enabled = feishuSetting.enabled
  feishuForm.webhook = ''
  feishuForm.secret = ''
  feishuForm.clear_secret = false
})

async function loadFeishuSetting() {
  const response = await queryFeishu()
  const setting = response.data
  feishuSetting.id = setting.id
  feishuSetting.enabled = setting.enabled
  feishuSetting.webhook_masked = setting.webhook_masked
  feishuSetting.webhook_configured = setting.webhook_configured
  feishuSetting.secret_configured = setting.secret_configured
}

async function saveFeishu() {
  if (feishuForm.enabled && !feishuSetting.webhook_configured && !feishuForm.webhook.trim()) {
    warning(t('setting.feishuWebhookRequired'))
    return
  }
  feishuSaving.value = true
  try {
    await updateFeishu({
      enabled: feishuForm.enabled,
      webhook: feishuForm.webhook.trim(),
      secret: feishuForm.secret.trim(),
      clear_secret: feishuForm.clear_secret,
    })
    info(t('setting.feishuUpdated'))
    await loadFeishuSetting()
  }
  catch {
    // 请求拦截器负责展示服务端错误；保留本地表单供用户重试。
  }
  finally {
    feishuSaving.value = false
  }
}

async function feishuTest() {
  feishuTesting.value = true
  try {
    await testFeishu()
    feishuTestedAt.value = now()
    info(t('setting.feishuTestSuccess'))
  }
  catch {
    // 请求拦截器负责展示服务端错误，失败时不显示成功提示。
  }
  finally {
    feishuTesting.value = false
  }
}
const feishuTesting = shallowRef(false)

const weComSetting = reactive<WeComSetting>({
  id: 0,
  enabled: false,
  webhook_masked: '',
  webhook_configured: false,
  at_all: false,
})
const weComForm = reactive<WeComUpdateArgs>({
  enabled: false,
  webhook: '',
  at_all: false,
})
const weComSaving = shallowRef(false)
const weComTestedAt = shallowRef('')
const weComTesting = shallowRef(false)

watch(() => ({
  enabled: weComSetting.enabled,
  webhook_configured: weComSetting.webhook_configured,
}), () => {
  weComForm.enabled = weComSetting.enabled
  weComForm.webhook = ''
  weComForm.at_all = weComSetting.at_all
})

async function loadWeComSetting() {
  const response = await queryWeCom()
  const setting = response.data
  weComSetting.id = setting.id
  weComSetting.enabled = setting.enabled
  weComSetting.webhook_masked = setting.webhook_masked
  weComSetting.webhook_configured = setting.webhook_configured
  weComSetting.at_all = setting.at_all
}

async function saveWeCom() {
  if (weComForm.enabled && !weComSetting.webhook_configured && !weComForm.webhook.trim()) {
    warning(t('setting.weComWebhookRequired'))
    return
  }
  weComSaving.value = true
  try {
    await updateWeCom({
      enabled: weComForm.enabled,
      webhook: weComForm.webhook.trim(),
      at_all: weComForm.at_all,
    })
    info(t('setting.weComUpdated'))
    await loadWeComSetting()
  }
  catch {
    // 请求拦截器负责展示服务端错误；保留本地表单供用户重试。
  }
  finally {
    weComSaving.value = false
  }
}

async function weComTest() {
  weComTesting.value = true
  try {
    await testWeCom()
    weComTestedAt.value = now()
    info(t('setting.weComTestSuccess'))
  }
  catch {
    // 请求拦截器负责展示服务端错误，失败时不显示成功提示。
  }
  finally {
    weComTesting.value = false
  }
}

const alarmMessageSetting = reactive<AlarmMessageSetting>({
  message_title: 'Beacon 主机告警',
})
const editAlarmMessageTitle = useCommandComponent(EditAlarmMessageTitle)

// 卡片只展示本地观察到的最近一次投递结果，服务端未提供测试时间戳。
const mailTestedAt = shallowRef('')
const dingTalkTestedAt = shallowRef('')

function now(): string {
  return new Date().toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US', { hour12: false })
}

// 告警阈值
const CPUThreshold = ref<AlarmThreshold>({
  id: 1,
  type: 'cpu',
  duration: 2,
  threshold: 80,
})
const MemThreshold = ref<AlarmThreshold>({
  id: 2,
  type: 'memory',
  duration: 2,
  threshold: 80,
})
const DiskThreshold = ref<AlarmThreshold>({
  id: 3,
  type: 'disk',
  duration: 2,
  threshold: 80,
})

const mailConfigured = computed(() => Boolean(emailSetting.server && emailSetting.sender))
const dingTalkBadge = computed(() => (dingTalkSetting.enabled
  ? { text: t('setting.enabled'), tone: 'success' }
  : { text: t('setting.disabled'), tone: 'muted' }))

const cpuThresholdText = computed(() => (locale.value === 'zh'
  ? `${t('setting.cpuUsage')} ${CPUThreshold.value.duration} ${t('setting.over')} ${CPUThreshold.value.threshold}%`
  : `${t('setting.cpuUsage')} ${CPUThreshold.value.threshold}% for ${CPUThreshold.value.duration} ${t('setting.over')}`))

const memThresholdText = computed(() => (locale.value === 'zh'
  ? `${t('setting.memUsage')} ${MemThreshold.value.duration} ${t('setting.over')} ${MemThreshold.value.threshold}%`
  : `${t('setting.memUsage')} ${MemThreshold.value.threshold}% for ${MemThreshold.value.duration} ${t('setting.over')}`))

const diskThresholdText = computed(() => `${t('setting.diskUsage')} ${DiskThreshold.value.threshold}%`)

async function loadAlarmMessageSetting() {
  const response = await queryAlarmMessageSetting()
  Object.assign(alarmMessageSetting, response.data)
}

async function loadDingTalkSetting() {
  const response = await queryDingTalk()
  const setting = response.data
  dingTalkSetting.id = setting.id
  dingTalkSetting.enabled = setting.enabled
  dingTalkSetting.webhook_masked = setting.webhook_masked
  dingTalkSetting.webhook_configured = setting.webhook_configured
  dingTalkSetting.secret_configured = setting.secret_configured
  dingTalkSetting.at_all = setting.at_all
}

async function dingTalkTest() {
  dingTalkTesting.value = true
  try {
    await testDingTalk()
    dingTalkTestedAt.value = now()
    info(t('setting.dingTalkTestSuccess'))
  }
  catch {
    // 请求拦截器负责展示服务端错误，失败时不显示成功提示。
  }
  finally {
    dingTalkTesting.value = false
  }
}

const testReceiver = shallowRef('')
// 测试收件人默认跟随已保存的收信邮箱账号（receiver 只在编辑弹层维护，
// 面板内的测试输入框是同一收件人概念的第二个入口，避免两处值漂移）；
// 用户手动改过测试收件人后不再自动覆盖。
let lastSavedReceiver = ''
watch(() => emailSetting.receiver, (receiver) => {
  if (!testReceiver.value || testReceiver.value === lastSavedReceiver)
    testReceiver.value = receiver
  lastSavedReceiver = receiver
})
async function mailTest() {
  const Params: MailTestArgs = {
    receiver: testReceiver.value,
  }
  try {
    await testMail(Params)
    mailTestedAt.value = now()
    info('邮件发送成功')
  }
  catch {
    // 请求拦截器负责展示错误，失败时不显示成功提示。
  }
}

function applyEmailSetting(setting: EmailSetting) {
  Object.assign(emailSetting, setting)
}

async function loadMailSetting() {
  const mailSetting = await queryMail()
  emailSetting.id = mailSetting.data.id
  emailSetting.server = mailSetting.data.server
  emailSetting.port = mailSetting.data.port
  emailSetting.sender = mailSetting.data.sender
  emailSetting.receiver = mailSetting.data.receiver
}

async function loadAlarmThresholds() {
  const { data } = await queryAlarmThreshold()
  for (const el of data.data) {
    if (el.type === 'cpu') {
      CPUThreshold.value.id = el.id
      CPUThreshold.value.type = el.type
      CPUThreshold.value.duration = el.duration
      CPUThreshold.value.threshold = el.threshold
    }
    else if (el.type === 'memory') {
      MemThreshold.value.id = el.id
      MemThreshold.value.type = el.type
      MemThreshold.value.duration = el.duration
      MemThreshold.value.threshold = el.threshold
    }
    else if (el.type === 'disk') {
      DiskThreshold.value.id = el.id
      DiskThreshold.value.type = el.type
      DiskThreshold.value.duration = el.duration
      DiskThreshold.value.threshold = el.threshold
    }
  }
}

onMounted(async () => {
  await Promise.all([
    loadMailSetting(),
    loadDingTalkSetting(),
    loadFeishuSetting(),
    loadWeComSetting(),
    loadAlarmMessageSetting(),
    loadAlarmThresholds(),
  ])
})

const editCPUThreshold = useCommandComponent(EditCPUThreshold)
const editMemThreshold = useCommandComponent(EditMemThreshold)
const editDiskThreshold = useCommandComponent(EditDiskThreshold)
</script>

<template>
    <div class="settings-grid">
        <article class="settings-card settings-card--warning">
            <header class="settings-card__header">
                <span class="settings-card__icon settings-card__icon--warning"><IconGauge /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.alarmThresholdSetting') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.alarmThresholdDescription') }}
                    </p>
                </div>
            </header>

            <div class="settings-thresholds">
                <div class="settings-threshold">
                    <span class="settings-threshold__label">{{ t('setting.cpuAlarmThreshold') }}</span>
                    <span class="settings-threshold__value">{{ cpuThresholdText }}</span>
                    <el-button
                        class="settings-chip"
                        :aria-label="t('setting.edit')"
                        @click="editCPUThreshold({ title: 'setting.cpuAlarmThreshold', threshold: CPUThreshold, update: loadAlarmThresholds })"
                    >
                        {{ t('setting.edit') }}
                    </el-button>
                </div>
                <div class="settings-threshold">
                    <span class="settings-threshold__label">{{ t('setting.memAlarmThreshold') }}</span>
                    <span class="settings-threshold__value">{{ memThresholdText }}</span>
                    <el-button
                        class="settings-chip"
                        :aria-label="t('setting.edit')"
                        @click="editMemThreshold({ title: 'setting.memAlarmThreshold', threshold: MemThreshold, update: loadAlarmThresholds })"
                    >
                        {{ t('setting.edit') }}
                    </el-button>
                </div>
                <div class="settings-threshold">
                    <span class="settings-threshold__label">{{ t('setting.diskAlarmThreshold') }}</span>
                    <span class="settings-threshold__value">{{ diskThresholdText }}</span>
                    <el-button
                        class="settings-chip"
                        :aria-label="t('setting.edit')"
                        @click="editDiskThreshold({ title: 'setting.diskAlarmThreshold', threshold: DiskThreshold, update: loadAlarmThresholds })"
                    >
                        {{ t('setting.edit') }}
                    </el-button>
                </div>
            </div>
        </article>

        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon"><IconMailCheck /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.mailServerSetting') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.mailServerCardTips') }}
                    </p>
                </div>
                <span
                    class="settings-card__badge"
                    :class="mailConfigured ? 'settings-card__badge--success' : 'settings-card__badge--muted'"
                >
                    {{ mailConfigured ? t('setting.enabled') : t('setting.notConfigured') }}
                </span>
            </header>

            <div class="settings-panels">
                <section class="settings-panel">
                    <header class="settings-panel__header">
                        <span class="settings-panel__label">{{ t('setting.mailSenderPanel') }}</span>
                        <el-button
                            class="settings-chip settings-chip--sm"
                            :aria-label="t('setting.edit')"
                            :title="t('setting.edit')"
                            @click="editEmailSetting({ title: 'setting.mailServerSetting', setting: emailSetting, onSaved: applyEmailSetting })"
                        >
                            <IconPencil />{{ t('setting.edit') }}
                        </el-button>
                    </header>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.mailServerHost') }}</span>
                        <span class="settings-row__value">{{ emailSetting.server || '—' }}</span>
                    </div>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.mailServerPort') }}</span>
                        <span class="settings-row__value">{{ emailSetting.port }}</span>
                    </div>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.mailServerAccount') }}</span>
                        <span class="settings-row__value">{{ emailSetting.sender || '—' }}</span>
                    </div>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.mailServerPassword') }}</span>
                        <span class="settings-row__value settings-row__value--text">•••••••••••• {{ t('setting.saved') }}</span>
                    </div>
                </section>

                <section class="settings-panel">
                    <header class="settings-panel__header">
                        <span class="settings-panel__label">{{ t('setting.mailVerifyPanel') }}</span>
                    </header>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.mailReceiver') }}</span>
                        <span class="settings-row__value">{{ emailSetting.receiver || '—' }}</span>
                    </div>
                    <p class="settings-row__hint">
                        {{ t('setting.mailReceiverHint') }}
                    </p>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.testSend') }}</span>
                        <el-input
                            v-model="testReceiver"
                            class="settings-row__input"
                            :placeholder="t('setting.receiverPlaceholder')"
                        />
                        <el-button class="settings-chip settings-chip--sm" @click="mailTest">
                            {{ t('setting.test') }}
                        </el-button>
                    </div>
                    <div class="settings-row">
                        <span class="settings-row__label">{{ t('setting.lastTest') }}</span>
                        <span
                            v-if="mailTestedAt"
                            class="settings-row__value settings-row__value--ok"
                        >{{ mailTestedAt }} · {{ t('setting.testSuccess') }}</span>
                        <span v-else class="settings-row__value settings-row__value--text">{{ t('setting.notTested') }}</span>
                    </div>
                </section>
            </div>

            <p class="settings-card__hint">
                {{ t('setting.mailChannelHint') }}
            </p>
        </article>

        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon"><IconWebhook /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.dingTalkGroupSetting') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.dingTalkGroupSettingTips') }}
                    </p>
                </div>
                <span class="settings-card__badge" :class="`settings-card__badge--${dingTalkBadge.tone}`">
                    {{ dingTalkBadge.text }}
                </span>
            </header>

            <section v-if="dingTalkChannel === 'dingtalk'" class="settings-panel">
                <header class="settings-panel__header">
                    <span class="settings-panel__label">{{ t('setting.dingTalkRobot') }}</span>
                </header>
                <div class="settings-field" role="radiogroup" :aria-label="t('setting.dingTalkRobot')">
                    <el-radio-group v-model="dingTalkChannel">
                        <el-radio value="dingtalk">{{ t('setting.channelDingTalk') }}</el-radio>
                        <el-radio value="feishu">{{ t('setting.channelFeishu') }}</el-radio>
                        <el-radio value="wecom">{{ t('setting.channelWeCom') }}</el-radio>
                    </el-radio-group>
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.dingTalkEnabled') }}</span>
                    <el-switch v-model="dingTalkForm.enabled" />
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.dingTalkWebhook') }}</span>
                    <el-input
                        v-model="dingTalkForm.webhook"
                        data-testid="dingtalk-webhook-input"
                        type="password"
                        show-password
                        autocomplete="off"
                        :placeholder="dingTalkSetting.webhook_configured ? t('setting.credentialKeepPlaceholder') : t('setting.dingTalkWebhookPlaceholder')"
                    />
                    <p v-if="dingTalkSetting.webhook_configured" class="settings-field__hint">
                        {{ t('setting.currentCredential') }}：<code
                            class="settings-field__mask"
                            data-testid="dingtalk-webhook-mask"
                            :title="dingTalkSetting.webhook_masked"
                        >{{ dingTalkSetting.webhook_masked }}</code>
                    </p>
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.dingTalkSecret') }}</span>
                    <el-input
                        v-model="dingTalkForm.secret"
                        type="password"
                        show-password
                        autocomplete="off"
                        :disabled="dingTalkForm.clear_secret"
                        :placeholder="dingTalkSetting.secret_configured ? t('setting.credentialKeepPlaceholder') : t('setting.dingTalkSecretPlaceholder')"
                    />
                    <el-checkbox v-if="dingTalkSetting.secret_configured" v-model="dingTalkForm.clear_secret">
                        {{ t('setting.clearDingTalkSecret') }}
                    </el-checkbox>
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.dingTalkMention') }}</span>
                    <el-checkbox v-model="dingTalkForm.at_all">
                        {{ t('setting.dingTalkAtAll') }}
                    </el-checkbox>
                </div>
                <div class="settings-panel__footer">
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="dingtalk-save"
                        :loading="dingTalkSaving"
                        @click="saveDingTalk"
                    >
                        {{ t('setting.save') }}
                    </el-button>
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="dingtalk-test"
                        :disabled="!dingTalkSetting.webhook_configured"
                        :loading="dingTalkTesting"
                        @click="dingTalkTest"
                    >
                        {{ t('setting.test') }}
                    </el-button>
                    <span class="settings-row__label">{{ t('setting.lastTest') }}</span>
                    <span
                        v-if="dingTalkTestedAt"
                        class="settings-row__value settings-row__value--ok"
                    >{{ dingTalkTestedAt }} · {{ t('setting.testSuccess') }}</span>
                    <span v-else class="settings-row__value settings-row__value--text">{{ t('setting.notTested') }}</span>
                </div>
            </section>

            <section v-if="dingTalkChannel === 'feishu'" class="settings-panel">
                <header class="settings-panel__header">
                    <span class="settings-panel__label">{{ t('setting.feishuRobot') }}</span>
                </header>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.feishuEnabled') }}</span>
                    <el-switch v-model="feishuForm.enabled" />
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.feishuWebhook') }}</span>
                    <el-input
                        v-model="feishuForm.webhook"
                        data-testid="feishu-webhook-input"
                        type="password"
                        show-password
                        autocomplete="off"
                        :placeholder="feishuSetting.webhook_configured ? t('setting.credentialKeepPlaceholder') : t('setting.feishuWebhookPlaceholder')"
                    />
                    <p v-if="feishuSetting.webhook_configured" class="settings-field__hint">
                        {{ t('setting.currentCredential') }}：<code
                            class="settings-field__mask"
                            data-testid="feishu-webhook-mask"
                            :title="feishuSetting.webhook_masked"
                        >{{ feishuSetting.webhook_masked }}</code>
                    </p>
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.feishuSecret') }}</span>
                    <el-input
                        v-model="feishuForm.secret"
                        type="password"
                        show-password
                        autocomplete="off"
                        :disabled="feishuForm.clear_secret"
                        :placeholder="feishuSetting.secret_configured ? t('setting.credentialKeepPlaceholder') : t('setting.feishuSecretPlaceholder')"
                    />
                    <el-checkbox v-if="feishuSetting.secret_configured" v-model="feishuForm.clear_secret">
                        {{ t('setting.clearFeishuSecret') }}
                    </el-checkbox>
                </div>
                <div class="settings-panel__footer">
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="feishu-save"
                        :loading="feishuSaving"
                        @click="saveFeishu"
                    >
                        {{ t('setting.save') }}
                    </el-button>
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="feishu-test"
                        :disabled="!feishuSetting.webhook_configured"
                        :loading="feishuTesting"
                        @click="feishuTest"
                    >
                        {{ t('setting.test') }}
                    </el-button>
                    <span class="settings-row__label">{{ t('setting.lastTest') }}</span>
                    <span
                        v-if="feishuTestedAt"
                        class="settings-row__value settings-row__value--ok"
                    >{{ feishuTestedAt }} · {{ t('setting.testSuccess') }}</span>
                    <span v-else class="settings-row__value settings-row__value--text">{{ t('setting.notTested') }}</span>
                </div>
            </section>

            <section v-if="dingTalkChannel === 'wecom'" class="settings-panel">
                <header class="settings-panel__header">
                    <span class="settings-panel__label">{{ t('setting.weComRobot') }}</span>
                </header>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.weComEnabled') }}</span>
                    <el-switch v-model="weComForm.enabled" />
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.weComWebhook') }}</span>
                    <el-input
                        v-model="weComForm.webhook"
                        data-testid="wecom-webhook-input"
                        type="password"
                        show-password
                        autocomplete="off"
                        :placeholder="weComSetting.webhook_configured ? t('setting.credentialKeepPlaceholder') : t('setting.weComWebhookPlaceholder')"
                    />
                    <p v-if="weComSetting.webhook_configured" class="settings-field__hint">
                        {{ t('setting.currentCredential') }}：<code
                            class="settings-field__mask"
                            data-testid="wecom-webhook-mask"
                            :title="weComSetting.webhook_masked"
                        >{{ weComSetting.webhook_masked }}</code>
                    </p>
                </div>
                <div class="settings-field">
                    <span class="settings-row__label">{{ t('setting.dingTalkMention') }}</span>
                    <el-checkbox v-model="weComForm.at_all">
                        {{ t('setting.dingTalkAtAll') }}
                    </el-checkbox>
                </div>
                <div class="settings-panel__footer">
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="wecom-save"
                        :loading="weComSaving"
                        @click="saveWeCom"
                    >
                        {{ t('setting.save') }}
                    </el-button>
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="wecom-test"
                        :disabled="!weComSetting.webhook_configured"
                        :loading="weComTesting"
                        @click="weComTest"
                    >
                        {{ t('setting.test') }}
                    </el-button>
                    <span class="settings-row__label">{{ t('setting.lastTest') }}</span>
                    <span
                        v-if="weComTestedAt"
                        class="settings-row__value settings-row__value--ok"
                    >{{ weComTestedAt }} · {{ t('setting.testSuccess') }}</span>
                    <span v-else class="settings-row__value settings-row__value--text">{{ t('setting.notTested') }}</span>
                </div>
            </section>
        </article>

        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon"><IconBellRing /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.alarmMessageTitle') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.alarmMessageTitleCardTips') }}
                    </p>
                </div>
            </header>

            <div class="settings-threshold">
                <span class="settings-threshold__label">{{ t('setting.alarmMessageTitle') }}</span>
                <span class="settings-threshold__value">{{ alarmMessageSetting.message_title }}</span>
                <el-button
                    class="settings-chip"
                    :aria-label="t('setting.editAlarmMessageTitle')"
                    @click="editAlarmMessageTitle({ title: 'setting.editAlarmMessageTitle', setting: alarmMessageSetting, onSaved: loadAlarmMessageSetting })"
                >
                    {{ t('setting.edit') }}
                </el-button>
            </div>

            <p class="settings-card__hint">
                {{ t('setting.alarmMessageTitleTips') }}
            </p>
        </article>
    </div>
</template>

<style scoped lang="scss">
.settings-row__input {
  flex: 1;
  min-width: 0;
}

.settings-row__input :deep(.el-input__wrapper) {
  min-height: 28px;
  background: var(--am-surface-primary);
}
</style>
