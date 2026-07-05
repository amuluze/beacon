<script setup lang="ts">
import type { SetDockerRegistryMirrorsArgs } from '@/interface/container.ts'
import { SetDockerRegistryMirrors } from '@/api/container'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  registryMirrors: string
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const drawerVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 设置docker镜像仓库
const textarea = ref('')

onMounted(() => {
  textarea.value = props.registryMirrors
})

async function confirmEditDockerRegistryMirrors() {
  const params: SetDockerRegistryMirrorsArgs = {
    registry_mirrors: textarea.value.split('\n').map(item => item.trim()),
  }
  console.log('set mirrors params: ', params)
  await SetDockerRegistryMirrors(params)
  drawerVisible.value = false
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="540" :title="t(props.title as string)">
        <el-form>
            <el-form-item>
                <el-input v-model="textarea" :rows="6" type="textarea" />
            </el-form-item>
        </el-form>
        <el-button type="primary" size="default" plain @click="drawerVisible = false">
            {{ t('setting.cancel') }}
        </el-button>
        <el-button type="primary" size="default" plain @click="confirmEditDockerRegistryMirrors">
            {{ t('setting.confirm') }}
        </el-button>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
