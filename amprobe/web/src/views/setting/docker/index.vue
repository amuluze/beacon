<script setup lang="ts">
import { getDockerRegistryMirrors } from '@/api/container'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import SetRegistryMirrors from '@/views/setting/docker/components/SetRegistryMirrors.vue'
import { useI18n } from 'vue-i18n'

const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const initialized = ref(false)
const textarea = ref('')
async function queryDockerRegistryMirrors() {
  const { data } = await getDockerRegistryMirrors()
  textarea.value = data.registry_mirrors.join('\n')
}

async function refreshDockerSettings() {
  const agentID = await ensureSelectedAgent()
  if (!agentID)
    return
  await queryDockerRegistryMirrors()
}
async function refreshAgents() {
  await loadAgents()
  await refreshDockerSettings()
}

onMounted(async () => {
  await refreshDockerSettings()
  initialized.value = true
})

watch(selectedAgentID, async () => {
  if (initialized.value && selectedAgentID.value) {
    await refreshDockerSettings()
  }
})

const editDockerRegistryMirrors = useCommandComponent(SetRegistryMirrors)
const { t } = useI18n()
</script>

<template>
    <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
    <el-row v-else :gutter="8">
        <el-col :span="12">
            <el-card shadow="never">
                <h4>{{ t('setting.mirrorRegistry') }}</h4>
                <el-input v-model="textarea" :rows="6" type="textarea" placeholder="https://docker.nju.edu.cn" />
                <p>{{ t('setting.mirrorRegistryTips') }}</p>
                <el-button type="primary" plain size="small" :disabled="!selectedAgentID" @click="editDockerRegistryMirrors({ title: 'setting.mirrorRegistry', registryMirrors: textarea })">
                    <svg-icon icon-class="settings" />
                    {{ t('setting.setting') }}
                </el-button>
            </el-card>
        </el-col>
    </el-row>
</template>

<style scoped lang="scss">
</style>
