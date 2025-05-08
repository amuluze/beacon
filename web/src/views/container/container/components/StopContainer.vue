<script setup lang="ts">
import type { StopContainerArgs } from '@/interface/container.ts'
import { useI18n } from 'vue-i18n'
import { stopContainer } from '@/api/container'

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

// 停止容器
const containerStopLoading = ref(false)
async function confirmStop() {
  containerStopLoading.value = true
  const params: StopContainerArgs = {
    container_id: props.id,
  }
  stopContainer(params).finally(() => {
    containerStopLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('container.confirmStop') }}？</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('container.cancel') }}</el-button>
                <el-button type="primary" @click="confirmStop">{{ t('container.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
