<script setup lang="ts">
import type { StartContainerArgs } from '@/interface/container.ts'
import { useI18n } from 'vue-i18n'
import { startContainer } from '@/api/container'

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

// 删除容器
const containerStartLoading = ref(false)
async function confirmStart() {
  containerStartLoading.value = true
  const params: StartContainerArgs = {
    container_id: props.id,
  }
  startContainer(params).finally(() => {
    containerStartLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('container.confirmStart') }}？</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('container.cancel') }}</el-button>
                <el-button type="primary" @click="confirmStart">{{ t('container.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
