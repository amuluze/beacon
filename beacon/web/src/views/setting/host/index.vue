<script setup lang="ts">
import { getSystemTime, getSystemTimezone, reboot, shutdown } from '@/api/system'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { error, success } from '@/components/Message/message.ts'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import SetSystemTime from '@/views/setting/host/components/SetSystemTime.vue'
import SetSystemTimezone from '@/views/setting/host/components/SetSystemTimezone.vue'
import { useI18n } from 'vue-i18n'
import IconClock3 from '~icons/lucide/clock-3'
import IconGlobe2 from '~icons/lucide/globe-2'
import IconPencil from '~icons/lucide/pencil'
import IconPower from '~icons/lucide/power'

const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const initialized = shallowRef(false)

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

const systemTime = shallowRef('')
const systemTimezone = shallowRef('')

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
    <div v-else class="settings-grid settings-grid--system">
        <article class="settings-card settings-card--wide settings-card--warning">
            <header class="settings-card__header">
                <span class="settings-card__icon settings-card__icon--warning"><IconPower /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.systemControl') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.systemControlTips') }}
                    </p>
                </div>
            </header>
            <div class="settings-card__actions settings-card__actions--end">
                <el-button type="warning" plain :disabled="!selectedAgentID" @click="rebootHost">
                    {{ t('setting.reboot') }}
                </el-button>
                <el-button type="danger" plain :disabled="!selectedAgentID" @click="shutdownHost">
                    {{ t('setting.shutdown') }}
                </el-button>
            </div>
        </article>

        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon"><IconGlobe2 /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.systemTimezone') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.systemTimezoneTips') }}
                    </p>
                </div>
            </header>
            <div class="settings-card__value-row">
                <code class="settings-card__value">{{ systemTimezone || '—' }}</code>
                <el-button
                    plain
                    circle
                    :aria-label="t('setting.edit')"
                    :title="t('setting.edit')"
                    :disabled="!selectedAgentID"
                    @click="editSystemTimezone({ title: 'setting.systemTimezone', systemTimezone, update: refreshSystemSettings })"
                >
                    <IconPencil />
                </el-button>
            </div>
        </article>

        <article class="settings-card">
            <header class="settings-card__header">
                <span class="settings-card__icon settings-card__icon--success"><IconClock3 /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.systemTime') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.systemTimeTips') }}
                    </p>
                </div>
            </header>
            <div class="settings-card__value-row">
                <code class="settings-card__value">{{ systemTime || '—' }}</code>
                <el-button
                    plain
                    circle
                    :aria-label="t('setting.edit')"
                    :title="t('setting.edit')"
                    :disabled="!selectedAgentID"
                    @click="editSystemTime({ title: 'setting.systemTime', systemTime, update: refreshSystemSettings })"
                >
                    <IconPencil />
                </el-button>
            </div>
        </article>
    </div>
</template>

<style scoped lang="scss">
.settings-grid--system {
  .settings-card__value {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
