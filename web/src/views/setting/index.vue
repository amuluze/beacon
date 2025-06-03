<script setup lang="ts">
// 导入所需的API和组件
import { queryAlarmThreshold } from '@/api/alarm'
import { getDockerRegistryMirrors } from '@/api/container'
import { queryMail, testMail } from '@/api/mail'
import { getSystemTime, getSystemTimezone, reboot, shutdown } from '@/api/system'
import { error, info, success } from '@/components/Message/message.ts'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import useStore from '@/store'
import { useI18n } from 'vue-i18n'

// 导入子组件
import EditCPUThreshold from './components/EditCPUThreshold.vue'
import EditDiskThreshold from './components/EditDiskThreshold.vue'
import EditEmailSetting from './components/EditEmailSetting.vue'
import EditMemThreshold from './components/EditMemThreshold.vue'
import SetRegistryMirrors from './components/SetRegistryMirrors.vue'
import SetSystemTime from './components/SetSystemTime.vue'
import SetSystemTimezone from './components/SetSystemTimezone.vue'

// 导入类型
import type { AlarmThreshold, EmailSetting } from '@/interface/alarm.ts'
import type { MailTestArgs } from '@/interface/mail.ts'

const { t } = useI18n();
const store = useStore()
const locale = store.app.language

// 主机设置相关
const systemTime = ref('')
const systemTimezone = ref('')

function rebootHost() {
    reboot()
        .then((res) => {
            console.log(res)
            success('重启成功')
        })
        .catch((err) => {
            error('重启失败')
            console.log(err)
        })
}

function shutdownHost() {
    shutdown()
        .then((res) => {
            console.log(res)
            success('关机成功')
        })
        .catch((err) => {
            console.log(err)
            error('关机失败')
        })
}

async function querySystemTime() {
    // 获取系统时间
    const { data } = await getSystemTime()
    systemTime.value = data.system_time
}

async function querySystemTimezone() {
    // 获取系统时区
    const { data } = await getSystemTimezone()
    systemTimezone.value = data.system_timezone
}

// Docker设置相关
const textarea = ref('')
async function queryDockerRegistryMirrors() {
    const { data } = await getDockerRegistryMirrors()
    textarea.value = data.registry_mirrors.join('\n')
}

// 告警设置相关
const emailSetting = reactive<EmailSetting>({
    id: 0,
    server: 'test@163.com',
    port: 465,
    sender: 'test@163.com',
    password: '******',
    receiver: 'test@163.com',
})

const testReceiver = ref('test@163.com')
function mailTest() {
    const Params: MailTestArgs = {
        receiver: testReceiver.value,
    }
    testMail(Params).finally(() => {
        info('邮件发送成功')
    })
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

async function initAlarmSettings() {
    const mailSetting = await queryMail()
    emailSetting.id = mailSetting.data.id
    emailSetting.server = mailSetting.data.server
    emailSetting.port = mailSetting.data.port
    emailSetting.sender = mailSetting.data.sender
    emailSetting.receiver = mailSetting.data.receiver

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

// 初始化数据
onMounted(() => {
    querySystemTime()
    querySystemTimezone()
    queryDockerRegistryMirrors()
    initAlarmSettings()
})

// 组件调用
const editSystemTime = useCommandComponent(SetSystemTime)
const editSystemTimezone = useCommandComponent(SetSystemTimezone)
const editDockerRegistryMirrors = useCommandComponent(SetRegistryMirrors)
const editEmailSetting = useCommandComponent(EditEmailSetting)
const editCPUThreshold = useCommandComponent(EditCPUThreshold)
const editMemThreshold = useCommandComponent(EditMemThreshold)
const editDiskThreshold = useCommandComponent(EditDiskThreshold)
</script>

<template>
    <div class="am-setting-container">
        <div class="am-setting-content">
            <div class="am-system">
                <el-card shadow="hover">
                    <el-button type="warning" plain size="small" @click="rebootHost">
                        {{ t('setting.reboot') }}
                    </el-button>
                    <el-button type="danger" plain size="small" @click="shutdownHost">
                        {{ t('setting.shutdown') }}
                    </el-button>
                </el-card>
            </div>
            <el-row :gutter="8">
                <el-col :span="12">
                    <el-card shadow="hover">
                        <h4>{{ t('setting.systemTimezone') }}</h4>
                        <span>{{ t('setting.systemTimezone') }}：</span>
                        <span style="margin-right: 4px">
                            <el-tag>{{ systemTimezone }}</el-tag>
                        </span>
                        <svg-icon icon-class="edit" style="cursor: pointer"
                            @click="editSystemTimezone({ title: 'setting.systemTimezone', systemTimezone })" />
                    </el-card>
                </el-col>
                <el-col :span="12">
                    <el-card shadow="never">
                        <h4>{{ t('setting.systemTime') }}</h4>
                        <span>{{ t('setting.systemTime') }}：</span>
                        <span style="margin-right: 4px">
                            <el-tag>{{ systemTime }}</el-tag>
                        </span>
                        <svg-icon icon-class="edit" style="cursor: pointer"
                            @click="editSystemTime({ title: 'setting.systemTime', systemTime })" />
                    </el-card>
                </el-col>
                <el-col :span="12">
                    <el-card shadow="never">
                        <el-descriptions :title="t('setting.mailServerSetting')" :column="1">
                            <el-descriptions-item :label="t('setting.alarmEmail')">
                                <span>{{ emailSetting.sender }}</span>
                                <svg-icon icon-class="edit" style="cursor: pointer"
                                    @click="editEmailSetting({ title: 'setting.mailServerSetting', setting: emailSetting })" />
                            </el-descriptions-item>
                            <div class="am-alarm-mail-test">
                                <el-descriptions-item :label="t('setting.testSend')">
                                    <el-input v-model="testReceiver" style="width: 240px" size="small"
                                        :placeholder="t('setting.receiverPlaceholder')" />
                                    <el-button style="margin-left: 8px;" size="small" plain type="primary"
                                        @click="mailTest">
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
                                <span v-if="locale === 'zh'" style="margin-right: 8px">{{ t('setting.cpuUsage') }}
                                    {{ CPUThreshold.duration }} {{ t('setting.over') }} {{ CPUThreshold.threshold
                                    }}%</span>
                                <span v-else style="margin-right: 8px">{{ t('setting.cpuUsage') }} {{
                                    CPUThreshold.threshold }}% for {{ CPUThreshold.duration }} {{ t('setting.over')
                                    }}</span>
                                <svg-icon icon-class="edit" style="cursor: pointer"
                                    @click="editCPUThreshold({ title: 'setting.cpuAlarmThreshold', threshold: CPUThreshold })" />
                            </el-descriptions-item>
                            <el-descriptions-item :label="t('setting.memAlarmThreshold')">
                                <span v-if="locale === 'zh'" style="margin-right: 8px">{{ t('setting.memUsage') }}
                                    {{ MemThreshold.duration }} {{ t('setting.over') }} {{ MemThreshold.threshold
                                    }}%</span>
                                <span v-else style="margin-right: 8px">{{ t('setting.memUsage') }} {{
                                    MemThreshold.threshold }}% for {{ MemThreshold.duration }} {{ t('setting.over')
                                    }}</span>
                                <svg-icon icon-class="edit" style="cursor: pointer"
                                    @click="editMemThreshold({ title: 'setting.memAlarmThreshold', threshold: MemThreshold })" />
                            </el-descriptions-item>
                            <el-descriptions-item :label="t('setting.diskAlarmThreshold')">
                                <span style="margin-right: 8px">{{ t('setting.diskUsage') }} {{
                                    DiskThreshold.threshold }}%</span>
                                <svg-icon icon-class="edit" style="cursor: pointer"
                                    @click="editDiskThreshold({ title: 'setting.diskAlarmThreshold', threshold: DiskThreshold })" />
                            </el-descriptions-item>
                        </el-descriptions>
                    </el-card>
                </el-col>
                <el-col :span="12">
                    <el-card shadow="never">
                        <h4>{{ t('setting.mirrorRegistry') }}</h4>
                        <el-input v-model="textarea" :rows="6" type="textarea"
                            placeholder="https://docker.nju.edu.cn" />
                        <p>{{ t('setting.mirrorRegistryTips') }}</p>
                        <el-button type="primary" plain size="small"
                            @click="editDockerRegistryMirrors({ title: 'setting.mirrorRegistry', registryMirrors: textarea })">
                            <svg-icon icon-class="settings" />
                            {{ t('setting.setting') }}
                        </el-button>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </div>
</template>

<style scoped lang="scss">
@include b(setting-container) {
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: var(--el-bg-color);
}

@include b(setting-content) {

    &::-webkit-scrollbar {
        width: 4px;
    }

    &::-webkit-scrollbar-thumb {
        background: var(--el-border-color-darker);
        border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
        background: var(--el-border-color-light);
        border-radius: 2px;
    }

    .el-row {
        margin-bottom: 8px;
    }

    .el-col {
        margin-bottom: 8px;
    }

    .el-card {
        border-radius: 8px;
        transition: all 0.3s;
        box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.05);
        border: 1px solid var(--el-border-color-light);
        height: 100%;

        &:hover {
            box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
        }

        :deep(.el-card__body) {
            padding: 16px;
        }

        h4 {
            font-size: 16px;
            font-weight: 500;
            color: var(--el-text-color-primary);
            margin-top: 0;
            margin-bottom: 12px;
        }

        .el-tag {
            border-radius: 4px;
            font-size: 14px;
            padding: 0 8px;
            height: 24px;
            line-height: 24px;
        }

        .el-button {
            transition: all 0.3s;
        }

        p {
            color: var(--el-text-color-secondary);
            font-size: 14px;
            margin: 8px 0;
        }

        span {
            color: var(--el-text-color-secondary);
            font-size: 14px;
        }
    }
}

@include b(system) {
    height: 48px;
    width: 100%;
    margin-bottom: 8px;

    .el-card {
        height: 100%;
        border-radius: 8px;
        transition: all 0.3s;
        box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.05);
        border: 1px solid var(--el-border-color-light);

        &:hover {
            box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
        }

        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 16px;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-end;
            gap: 8px;
        }
    }
}
</style>
