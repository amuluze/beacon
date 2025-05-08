<script setup lang="ts">
import { getSystemTime, getSystemTimezone, reboot, shutdown } from '@/api/system'
import { error, success } from '@/components/Message/message.ts'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import SetSystemTime from '@/views/setting/host/components/SetSystemTime.vue'
import SetSystemTimezone from '@/views/setting/host/components/SetSystemTimezone.vue'
import { useI18n } from 'vue-i18n'

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

const systemTime = ref('')
const systemTimezone = ref('')

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

onMounted(() => {
  querySystemTime()
  querySystemTimezone()
})

const editSystemTime = useCommandComponent(SetSystemTime)
const editSystemTimezone = useCommandComponent(SetSystemTimezone)

const { t } = useI18n()
</script>

<template>
    <div class="am-system">
        <el-card shadow="never">
            <el-button type="warning" plain size="small" @click="rebootHost">
                {{ t('setting.reboot') }}
            </el-button>
            <el-button type="danger" plain size="small" @click="shutdownHost">
                {{ t('setting.shutdown') }}
            </el-button>
        </el-card>
    </div>
    <el-row :gutter="4">
        <el-col :span="12">
            <el-card shadow="never">
                <h4>{{ t('setting.systemTimezone') }}</h4>
                <span>{{ t('setting.systemTimezone') }}：</span>
                <span style="margin-right: 4px">
                    <el-tag>{{ systemTimezone }}</el-tag>
                </span>
                <svg-icon icon-class="edit" style="cursor: pointer" @click="editSystemTimezone({ title: 'setting.systemTimezone', systemTimezone })" />
            </el-card>
        </el-col>
        <el-col :span="12">
            <el-card shadow="never">
                <h4>{{ t('setting.systemTime') }}</h4>
                <span>{{ t('setting.systemTime') }}：</span>
                <span style="margin-right: 4px">
                    <el-tag>{{ systemTime }}</el-tag>
                </span>
                <svg-icon icon-class="edit" style="cursor: pointer" @click="editSystemTime({ title: 'setting.systemTime', systemTime })" />
            </el-card>
        </el-col>
    </el-row>
</template>

<style scoped lang="scss">
@include b(system) {
  height: 48px;
  width: 100%;
  margin-bottom: 4px;
  .el-card {
    height: 100%;
    :deep(.el-card__body) {
      height: 100% !important;
      padding: 0 8px;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: flex-end;
    }
  }
}
</style>
