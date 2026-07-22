<script setup lang="ts">
import { queryMail, testMail } from '@/api/mail'
import { info } from '@/components/Message/message.ts'

import type { MailTestArgs } from '@/interface/mail.ts'
import type { AlarmThreshold, EmailSetting } from '@/interface/alarm.ts'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import EditEmailSetting from '@/views/setting/alarm/components/EditEmailSetting.vue'
import EditCPUThreshold from '@/views/setting/alarm/components/EditCPUThreshold.vue'
import EditMemThreshold from '@/views/setting/alarm/components/EditMemThreshold.vue'
import EditDiskThreshold from '@/views/setting/alarm/components/EditDiskThreshold.vue'
import { queryAlarmThreshold } from '@/api/alarm'
import { useI18n } from 'vue-i18n'
import useStore from '@/store'
import IconBellRing from '~icons/lucide/bell-ring'
import IconGauge from '~icons/lucide/gauge'
import IconHardDrive from '~icons/lucide/hard-drive'
import IconMailCheck from '~icons/lucide/mail-check'
import IconMemoryStick from '~icons/lucide/memory-stick'
import IconPencil from '~icons/lucide/pencil'

const emailSetting = reactive<EmailSetting>({
  id: 0,
  server: '',
  port: 465,
  sender: '',
  password: '******',
  receiver: '',
})

const editEmailSetting = useCommandComponent(EditEmailSetting)

const testReceiver = shallowRef('')
async function mailTest() {
  const Params: MailTestArgs = {
    receiver: testReceiver.value,
  }
  try {
    await testMail(Params)
    info('邮件发送成功')
  }
  catch {
    // 请求拦截器负责展示错误，失败时不显示成功提示。
  }
}

function applyEmailSetting(setting: EmailSetting) {
  Object.assign(emailSetting, setting)
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
  type: 'mem',
  duration: 2,
  threshold: 80,
})
const DiskThreshold = ref<AlarmThreshold>({
  id: 3,
  type: 'disk',
  duration: 2,
  threshold: 80,
})

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
  await Promise.all([loadMailSetting(), loadAlarmThresholds()])
})

const editCPUThreshold = useCommandComponent(EditCPUThreshold)
const editMemThreshold = useCommandComponent(EditMemThreshold)
const editDiskThreshold = useCommandComponent(EditDiskThreshold)

const { t } = useI18n()
const store = useStore()
const locale = computed(() => store.app.language)
</script>

<template>
    <div class="settings-grid settings-grid--alarm">
        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon"><IconMailCheck /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.mailServerSetting') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.mailServerSettingTips') }}
                    </p>
                </div>
            </header>

            <div class="settings-card__value-row">
                <div class="settings-card__meta">
                    <span class="settings-card__label">{{ t('setting.alarmEmail') }}</span>
                    <span class="settings-card__value settings-card__value--plain">{{ emailSetting.sender || '—' }}</span>
                </div>
                <el-button
                    plain
                    circle
                    :aria-label="t('setting.edit')"
                    :title="t('setting.edit')"
                    @click="editEmailSetting({ title: 'setting.mailServerSetting', setting: emailSetting, onSaved: applyEmailSetting })"
                >
                    <IconPencil />
                </el-button>
            </div>

            <div class="settings-card__form">
                <label class="settings-card__label" for="alarm-test-receiver">{{ t('setting.testSend') }}</label>
                <div class="settings-card__controls">
                    <el-input
                        id="alarm-test-receiver"
                        v-model="testReceiver"
                        :placeholder="t('setting.receiverPlaceholder')"
                    />
                    <el-button plain type="primary" @click="mailTest">
                        {{ t('setting.test') }}
                    </el-button>
                </div>
            </div>
        </article>

        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon settings-card__icon--warning"><IconBellRing /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.alarmThresholdSetting') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.alarmThresholdSettingTips') }}
                    </p>
                </div>
            </header>

            <div class="settings-thresholds">
                <div class="settings-threshold">
                    <span class="settings-threshold__icon"><IconGauge /></span>
                    <div class="settings-threshold__main">
                        <strong class="settings-threshold__title">{{ t('setting.cpuAlarmThreshold') }}</strong>
                        <span v-if="locale === 'zh'" class="settings-threshold__value">{{ t('setting.cpuUsage') }} {{ CPUThreshold.duration }} {{ t('setting.over') }} {{ CPUThreshold.threshold }}%</span>
                        <span v-else class="settings-threshold__value">{{ t('setting.cpuUsage') }} {{ CPUThreshold.threshold }}% for {{ CPUThreshold.duration }} {{ t('setting.over') }}</span>
                    </div>
                    <el-button
                        circle
                        text
                        :aria-label="t('setting.edit')"
                        :title="t('setting.edit')"
                        @click="editCPUThreshold({ title: 'setting.cpuAlarmThreshold', threshold: CPUThreshold, update: loadAlarmThresholds })"
                    >
                        <IconPencil />
                    </el-button>
                </div>

                <div class="settings-threshold">
                    <span class="settings-threshold__icon settings-threshold__icon--success"><IconMemoryStick /></span>
                    <div class="settings-threshold__main">
                        <strong class="settings-threshold__title">{{ t('setting.memAlarmThreshold') }}</strong>
                        <span v-if="locale === 'zh'" class="settings-threshold__value">{{ t('setting.memUsage') }} {{ MemThreshold.duration }} {{ t('setting.over') }} {{ MemThreshold.threshold }}%</span>
                        <span v-else class="settings-threshold__value">{{ t('setting.memUsage') }} {{ MemThreshold.threshold }}% for {{ MemThreshold.duration }} {{ t('setting.over') }}</span>
                    </div>
                    <el-button
                        circle
                        text
                        :aria-label="t('setting.edit')"
                        :title="t('setting.edit')"
                        @click="editMemThreshold({ title: 'setting.memAlarmThreshold', threshold: MemThreshold, update: loadAlarmThresholds })"
                    >
                        <IconPencil />
                    </el-button>
                </div>

                <div class="settings-threshold">
                    <span class="settings-threshold__icon settings-threshold__icon--danger"><IconHardDrive /></span>
                    <div class="settings-threshold__main">
                        <strong class="settings-threshold__title">{{ t('setting.diskAlarmThreshold') }}</strong>
                        <span class="settings-threshold__value">{{ t('setting.diskUsage') }} {{ DiskThreshold.threshold }}%</span>
                    </div>
                    <el-button
                        circle
                        text
                        :aria-label="t('setting.edit')"
                        :title="t('setting.edit')"
                        @click="editDiskThreshold({ title: 'setting.diskAlarmThreshold', threshold: DiskThreshold, update: loadAlarmThresholds })"
                    >
                        <IconPencil />
                    </el-button>
                </div>
            </div>
        </article>
    </div>
</template>

<style scoped lang="scss">
.settings-grid--alarm {
  align-items: stretch;
}
</style>
