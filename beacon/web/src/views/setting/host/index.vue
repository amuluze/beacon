<script setup lang="ts">
import { getSystemTime, getSystemTimezone, reboot, shutdown } from '@/api/system'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { error, success } from '@/components/Message/message.ts'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import SetSystemTime from '@/views/setting/host/components/SetSystemTime.vue'
import SetSystemTimezone from '@/views/setting/host/components/SetSystemTimezone.vue'
import { useI18n } from 'vue-i18n'

const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const initialized = ref(false)

function rebootHost() {
  if (!selectedAgentID.value)
    return
  reboot()
    .then(() => {
      success('重启成功')
    })
    .catch(() => {
      error('重启失败')
    })
}

function shutdownHost() {
  if (!selectedAgentID.value)
    return
  shutdown()
    .then(() => {
      success('关机成功')
    })
    .catch(() => {
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

async function refreshSystemSettings() {
  const agentID = await ensureSelectedAgent()
  if (!agentID)
    return
  await querySystemTime()
  await querySystemTimezone()
}
async function refreshAgents() {
  await loadAgents()
  await refreshSystemSettings()
}

onMounted(async () => {
  await refreshSystemSettings()
  initialized.value = true
})

watch(selectedAgentID, async () => {
  if (initialized.value && selectedAgentID.value) {
    await refreshSystemSettings()
  }
})

const editSystemTime = useCommandComponent(SetSystemTime)
const editSystemTimezone = useCommandComponent(SetSystemTimezone)

const { t } = useI18n()
</script>

<template>
    <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
    <template v-else>
        <div class="am-system">
            <el-card shadow="never">
                <el-button type="warning" plain size="small" :disabled="!selectedAgentID" @click="rebootHost">
                    {{ t('setting.reboot') }}
                </el-button>
                <el-button type="danger" plain size="small" :disabled="!selectedAgentID" @click="shutdownHost">
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
