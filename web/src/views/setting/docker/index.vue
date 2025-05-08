<script setup lang="ts">
import { getDockerRegistryMirrors } from '@/api/container'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import SetRegistryMirrors from '@/views/setting/docker/components/SetRegistryMirrors.vue'
import { useI18n } from 'vue-i18n'

const textarea = ref('')
async function queryDockerRegistryMirrors() {
  const { data } = await getDockerRegistryMirrors()
  textarea.value = data.registry_mirrors.join('\n')
}

onMounted(() => {
  queryDockerRegistryMirrors()
})

const editDockerRegistryMirrors = useCommandComponent(SetRegistryMirrors)
const { t } = useI18n()
</script>

<template>
    <el-row :gutter="8">
        <el-col :span="12">
            <el-card shadow="never">
                <h4>{{ t('setting.mirrorRegistry') }}</h4>
                <el-input v-model="textarea" :rows="6" type="textarea" placeholder="https://docker.nju.edu.cn" />
                <p>{{ t('setting.mirrorRegistryTips') }}</p>
                <el-button type="primary" plain size="small" @click="editDockerRegistryMirrors({ title: 'setting.mirrorRegistry', registryMirrors: textarea })">
                    <svg-icon icon-class="settings" />
                    {{ t('setting.setting') }}
                </el-button>
            </el-card>
        </el-col>
    </el-row>
</template>

<style scoped lang="scss">
</style>
