<template>
    <div class="am-host-header">
        <span @click="$router.push('/host/monitor')">监控</span>
        <span @click="$router.push('/host/file')">文件</span>
        <span @click="$router.push('/host/terminal')">终端</span>
        <span @click="$router.push('/host/settings')">设置</span>
    </div>
    <div class="am-host-operator">
        <el-card shadow="never">
            <el-button type="warning" plain @click="rebootHost">重启</el-button>
            <el-button type="danger" plain @click="shutdownHost">关机</el-button>
        </el-card>
    </div>
    <el-row :gutter="4" class="am-host-settings">
        <!-- <el-col :span="8">
            <el-card shadow="never">
                <h4>DNS 设置</h4>
                <el-divider />
                <div class="am-system-settings__content">
                    <span>DNS 设置：</span>
                    <span></span>
                    <svg-icon icon-class="edit" />
                </div>
            </el-card>
        </el-col> -->
        <el-col :span="12">
            <el-card shadow="never">
                <h4>系统时区设置</h4>
                <el-divider />
                <div class="am-system-settings__content">
                    <span>系统时区设置：</span>
                    <span style="margin-right: 4px">
                        <el-tag>{{ systemTimezone }}</el-tag>
                    </span>
                    <svg-icon icon-class="edit" style="cursor: pointer" @click="editSystemTimezone" />
                </div>
            </el-card>
        </el-col>
        <el-col :span="12">
            <el-card shadow="never">
                <h4>系统时间设置</h4>
                <el-divider />
                <div class="am-system-settings__content">
                    <span>系统时间设置：</span>
                    <span style="margin-right: 4px">
                        <el-tag>{{ systemTime }}</el-tag>
                    </span>
                    <svg-icon icon-class="edit" style="cursor: pointer" @click="editSystemTime" />
                </div>
            </el-card>
        </el-col>
    </el-row>
    <div class="am-host-edit-systemtime">
        <el-drawer v-model="systemTimeDrawer" size="540" title="系统时间设置">
            <div class="am-host-edit-systemtime__content">
                <el-input v-model="systemTime" style="width: 240px" />
                <el-button type="info" style="margin-left: 16px" @click="getClientTime">获取客户端时间</el-button>
            </div>
            <el-button type="default" size="default" plain @click="cancelSystemTimeEdit">取消</el-button>
            <el-button type="primary" size="default" plain @click="confirmSystemTimeEdit">确定</el-button>
        </el-drawer>
    </div>
    <div class="am-host-edit-systemtimezone">
        <el-drawer v-model="systemTimezoneDrawer" size="540" title="系统时区设置">
            <div class="am-host-edit-systemtimezone__content">
                <el-select v-model="systemTimezone" style="width: 240px" filterable>
                    <el-option
                        v-for="item in timezoneOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                    />
                </el-select>
            </div>
            <div class="am-host-edit-systemtimezone__operator">
                <el-button type="default" size="default" plain @click="cancelSystemTimezoneEdit">取消</el-button>
                <el-button type="primary" size="default" plain @click="confirmSystemTimezoneEdit">确定</el-button>
            </div>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import {
    getSystemTime,
    getSystemTimezone,
    getSystemTimezoneList,
    reboot,
    setSystemTime,
    setSystemTimezone,
    shutdown
} from '@/api/system'
import { error, success } from '@/components/Message/message'
import { SetSystemTimeArgs, SetSystemTimezoneArgs } from '@/interface/system'
import { dayjs } from 'element-plus'

const rebootHost = () => {
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

const shutdownHost = () => {
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

const systemTime = ref('')
const systemTimezone = ref('')

const querySystemTime = async () => {
    // 获取系统时间
    const { data } = await getSystemTime()
    systemTime.value = data.system_time
}

const systemTimeDrawer = ref(false)
const systemTimezoneDrawer = ref(false)
const timezoneOptions = ref<
    {
        label: string
        value: string
    }[]
>([])

const editSystemTime = () => {
    // 修改系统时间
    systemTimeDrawer.value = true
}

const cancelSystemTimeEdit = () => {
    // 取消修改系统时间
    systemTimeDrawer.value = false
}

const getClientTime = () => {
    // 获取客户端时间
    console.log(new Date().toLocaleString())
    systemTime.value = dayjs(new Date()).format('YYYY-MM-DD HH:mm:ss')
}

const confirmSystemTimeEdit = async () => {
    // 确定修改系统时间
    systemTimeDrawer.value = false
    systemTime.value = dayjs(new Date()).format('YYYY-MM-DD HH:mm:ss')
    console.log(systemTime.value)
    let args: SetSystemTimeArgs = {
        system_time: systemTime.value
    }
    const { data } = await setSystemTime(args)
    console.log(data)
    success('修改系统时间成功')
}

const editSystemTimezone = () => {
    // 修改系统时区
    systemTimezoneDrawer.value = true
}

const cancelSystemTimezoneEdit = () => {
    // 取消修改系统时区
    systemTimezoneDrawer.value = false
}

const confirmSystemTimezoneEdit = async () => {
    // 确定修改系统时区
    systemTimezoneDrawer.value = false
    console.log(systemTimezone.value)
    let args: SetSystemTimezoneArgs = {
        system_timezone: systemTimezone.value
    }
    const { data } = await setSystemTimezone(args)
    console.log(data)
    success('修改系统时区成功')
}

const querySystemTimezone = async () => {
    // 获取系统时区
    const { data } = await getSystemTimezone()
    systemTimezone.value = data.system_timezone
}

const querySystemTimezoneOptions = async () => {
    // 获取系统时区选项
    const { data } = await getSystemTimezoneList()
    console.log(data)
    console.log('zone list: ', data.system_timezone_list)
    data.system_timezone_list.forEach((item) => {
        timezoneOptions.value.push({
            label: item,
            value: item
        })
    })
    console.log(timezoneOptions.value)
}

onMounted(() => {
    querySystemTime()
    querySystemTimezone()
    querySystemTimezoneOptions()
})
</script>
<style scoped lang="scss">
@include b(host-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
    // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;
    span {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        font-size: 16px;
        font-weight: 600;
        margin-left: 16px;
        margin-right: 16px;
        color: #424244;
        cursor: pointer;
        &:first-child {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 212px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 8px 0 0;
            display: flex;
            align-items: center;
            justify-content: flex-end;
        }
    }
}

@include b(host-operator) {
    height: 48px;
    width: 100%;

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 16px 0 16px;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-end;
        }
    }

    border-radius: 4px;
    margin-bottom: 4px;
}

@include b(host-edit-systemtime) {
    :deep(.el-drawer__header) {
        margin-bottom: 0;
    }

    @include e(content) {
        margin-top: 16px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        margin-bottom: 32px;
    }

    @include e(operator) {
        margin-top: 16px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
    }
}

@include b(host-edit-systemtimezone) {
    :deep(.el-drawer__header) {
        margin-bottom: 0;
    }

    @include e(content) {
        margin-top: 16px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        margin-bottom: 32px;
    }

    @include e(operator) {
        margin-top: 16px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
    }
}
</style>
