<script setup lang="ts">
import type { RemoveImageArgs } from '@/interface/container.ts'
import { removeImage } from '@/api/container'
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

// 删除镜像
const imageDeleteLoading = ref(false)
async function confirmDelete() {
  imageDeleteLoading.value = true
  const params: RemoveImageArgs = {
    image_id: props.id,
  }
  removeImage(params).finally(() => {
    imageDeleteLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('image.confirmDelete') }}</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('image.cancel') }}</el-button>
                <el-button type="primary" @click="confirmDelete">{{ t('image.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
