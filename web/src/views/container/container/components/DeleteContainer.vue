<script setup lang="ts">
import { removeContainer } from '@/api/container'
import type { RemoveContainerArgs } from '@/interface/container.ts'
import { useI18n } from 'vue-i18n'

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
const containerDeleteLoading = ref(false)
async function confirmDelete() {
  containerDeleteLoading.value = true
  const params: RemoveContainerArgs = {
    container_id: props.id,
  }
  removeContainer(params).finally(() => {
    containerDeleteLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('container.confirmDelete') }}？</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('container.cancel') }}</el-button>
                <el-button type="primary" @click="confirmDelete">{{ t('container.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
