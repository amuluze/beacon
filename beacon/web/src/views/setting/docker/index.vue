<script setup lang="ts">
import { getDockerRegistryMirrors } from '@/api/container'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import SetRegistryMirrors from '@/views/setting/docker/components/SetRegistryMirrors.vue'
import { useI18n } from 'vue-i18n'
import IconBoxes from '~icons/lucide/boxes'
import IconSlidersHorizontal from '~icons/lucide/sliders-horizontal'

const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const initialized = shallowRef(false)
const textarea = shallowRef('')
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
    <div v-else class="settings-grid settings-grid--docker">
        <article class="settings-card settings-card--wide">
            <header class="settings-card__header">
                <span class="settings-card__icon"><IconBoxes /></span>
                <div class="settings-card__heading">
                    <h3 class="settings-card__title">
                        {{ t('setting.mirrorRegistry') }}
                    </h3>
                    <p class="settings-card__description">
                        {{ t('setting.mirrorRegistryDescription') }}
                    </p>
                </div>
            </header>

            <el-input
                v-model="textarea"
                class="settings-card__textarea"
                :rows="4"
                type="textarea"
                placeholder="https://docker.nju.edu.cn"
            />

            <p class="settings-card__hint">
                {{ t('setting.mirrorRegistryTips') }}
            </p>
            <div class="settings-card__actions">
                <el-button
                    type="primary"
                    :disabled="!selectedAgentID"
                    @click="editDockerRegistryMirrors({ title: 'setting.mirrorRegistry', registryMirrors: textarea })"
                >
                    <IconSlidersHorizontal />
                    {{ t('setting.save') }}
                </el-button>
            </div>
        </article>
    </div>
</template>

<style scoped lang="scss">
.settings-grid--docker {
  grid-template-columns: minmax(0, 1fr);
}
</style>
