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

const emailSetting = reactive<EmailSetting>({
  id: 0,
  server: '',
  port: 465,
  sender: '',
  password: '******',
  receiver: '',
})

const editEmailSetting = useCommandComponent(EditEmailSetting)

const testReceiver = ref('')
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
const locale = store.app.language
</script>

<template>
    <el-row class="am-email" :gutter="8">
        <el-col :span="12">
            <el-card shadow="never">
                <el-descriptions :title="t('setting.mailServerSetting')" :column="1">
                    <el-descriptions-item :label="t('setting.alarmEmail')">
                        {{ emailSetting.sender }}
                        <svg-icon icon-class="edit" style="cursor: pointer" @click="editEmailSetting({ title: 'setting.mailServerSetting', setting: emailSetting, onSaved: applyEmailSetting })" />
                    </el-descriptions-item>
                    <div class="am-alarm-mail-test">
                        <el-descriptions-item :label="t('setting.testSend')">
                            <el-input v-model="testReceiver" style="width: 240px" size="small" :placeholder="t('setting.receiverPlaceholder')" />
                            <el-button style="margin-left: 8px;" size="small" plain type="primary" @click="mailTest">
                                {{ t('setting.test') }}
                            </el-button>
                        </el-descriptions-item>
                    </div>
                </el-descriptions>
            </el-card>
        </el-col>
        <el-col :span="12">
            <el-card shadow="never">
                <el-descriptions :title="t('setting.alarmThresholdSetting')" :column="1">
                    <el-descriptions-item :label="t('setting.cpuAlarmThreshold')">
                        <span v-if="locale === 'zh'" style="margin-right: 8px">{{ t('setting.cpuUsage') }} {{ CPUThreshold.duration }} {{ t('setting.over') }} {{ CPUThreshold.threshold }}%</span>
                        <span v-else style="margin-right: 8px">{{ t('setting.cpuUsage') }} {{ CPUThreshold.threshold }}% for {{ CPUThreshold.duration }} {{ t('setting.over') }}</span>
                        <svg-icon icon-class="edit" style="cursor: pointer" @click="editCPUThreshold({ title: 'setting.cpuAlarmThreshold', threshold: CPUThreshold, update: loadAlarmThresholds })" />
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('setting.memAlarmThreshold')">
                        <span v-if="locale === 'zh'" style="margin-right: 8px">{{ t('setting.memUsage') }} {{ MemThreshold.duration }} {{ t('setting.over') }} {{ MemThreshold.threshold }}%</span>
                        <span v-else style="margin-right: 8px">{{ t('setting.memUsage') }} {{ MemThreshold.threshold }}% for {{ MemThreshold.duration }} {{ t('setting.over') }}</span>
                        <svg-icon icon-class="edit" style="cursor: pointer" @click="editMemThreshold({ title: 'setting.memAlarmThreshold', threshold: MemThreshold, update: loadAlarmThresholds })" />
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('setting.diskAlarmThreshold')">
                        <span style="margin-right: 8px">{{ t('setting.diskUsage') }} {{ DiskThreshold.threshold }}%</span>
                        <svg-icon icon-class="edit" style="cursor: pointer" @click="editDiskThreshold({ title: 'setting.diskAlarmThreshold', threshold: DiskThreshold, update: loadAlarmThresholds })" />
                    </el-descriptions-item>
                </el-descriptions>
            </el-card>
        </el-col>
    </el-row>
</template>

<style scoped lang="scss">
@include b(email) {
  .el-descriptions {
    :deep(.el-descriptions__title) {
      font-size: 14px;
    }
    :deep(.el-descriptions__label) {
      font-size: 14px;
    }
  }
}
</style>
