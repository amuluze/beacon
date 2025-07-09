<script setup lang="ts">
import { restartContainer } from '@/api/container'
import { useI18n } from 'vue-i18n'
import type { RestartContainerArgs } from '@/interface/container.ts'

const props = defineProps<{
  visible: boolean
  id: string
  title?: string
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const dialogVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 重启容器
const containerRestartLoading = ref(false)
async function confirmRestart() {
  containerRestartLoading.value = true
  const params: RestartContainerArgs = {
    container_id: props.id,
  }
  restartContainer(params).finally(() => {
    containerRestartLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('container.confirmRestart') }}？</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('container.cancel') }}</el-button>
                <el-button type="primary" @click="confirmRestart">{{ t('container.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
