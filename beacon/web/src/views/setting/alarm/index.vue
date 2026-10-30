<script setup lang="ts">
import { queryMail, testMail } from '@/api/mail'
import { queryDingTalk, testDingTalk } from '@/api/dingtalk'
import { info } from '@/components/Message/message.ts'

import type { MailTestArgs } from '@/interface/mail.ts'
import type { AlarmMessageSetting, AlarmThreshold, EmailSetting } from '@/interface/alarm.ts'
import type { DingTalkSetting } from '@/interface/dingtalk'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import EditEmailSetting from '@/views/setting/alarm/components/EditEmailSetting.vue'
import EditDingTalkSetting from '@/views/setting/alarm/components/EditDingTalkSetting.vue'
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
const editDingTalkSetting = useCommandComponent(EditDingTalkSetting)
const dingTalkTesting = shallowRef(false)

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

            <section class="settings-panel">
                <header class="settings-panel__header">
                    <span class="settings-panel__title">
                        <span class="settings-panel__mark"><IconWebhook /></span>
                        {{ t('setting.dingTalkRobot') }}
                    </span>
                </header>
                <div class="settings-panel__cols">
                    <div class="settings-panel__col">
                        <div class="settings-row">
                            <span class="settings-row__label">{{ t('setting.dingTalkWebhook') }}</span>
                            <span
                                class="settings-row__value"
                                data-testid="dingtalk-webhook-mask"
                                :title="dingTalkSetting.webhook_masked || t('setting.notConfigured')"
                            >{{ dingTalkSetting.webhook_masked || t('setting.notConfigured') }}</span>
                        </div>
                        <div class="settings-row">
                            <span class="settings-row__label">{{ t('setting.dingTalkSecret') }}</span>
                            <span class="settings-row__value">{{ dingTalkSetting.secret_configured ? t('setting.configured') : t('setting.notConfigured') }}</span>
                        </div>
                    </div>
                    <div class="settings-panel__col">
                        <div class="settings-row">
                            <span class="settings-row__label">{{ t('setting.dingTalkMention') }}</span>
                            <span class="settings-row__value settings-row__value--text">
                                {{ dingTalkSetting.at_all ? t('setting.dingTalkAtAllOn') : t('setting.dingTalkAtUsers') }}
                            </span>
                        </div>
                        <div class="settings-row">
                            <span class="settings-row__label">{{ t('setting.lastTest') }}</span>
                            <span
                                v-if="dingTalkTestedAt"
                                class="settings-row__value settings-row__value--ok"
                            >{{ dingTalkTestedAt }} · {{ t('setting.testSuccess') }}</span>
                            <span v-else class="settings-row__value settings-row__value--text">{{ t('setting.notTested') }}</span>
                        </div>
                    </div>
                </div>
                <div class="settings-panel__footer">
                    <el-button
                        class="settings-chip settings-chip--sm"
                        data-testid="dingtalk-edit"
                        :aria-label="t('setting.edit')"
                        :title="t('setting.edit')"
                        @click="editDingTalkSetting({ title: 'setting.dingTalkSetting', setting: dingTalkSetting, onSaved: loadDingTalkSetting })"
                    >
                        {{ t('setting.edit') }}
                    </el-button>
                    <el-button
                        class="settings-chip settings-chip--sm"
                        :disabled="!dingTalkSetting.webhook_configured"
                        :loading="dingTalkTesting"
                        @click="dingTalkTest"
                    >
                        {{ t('setting.test') }}
                    </el-button>
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
