<script setup lang="ts">
import { pruneImages } from '@/api/container'
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

// 清理虚悬镜像
const imagePruneLoading = ref(false)
async function confirmPrune() {
  imagePruneLoading.value = true
  pruneImages().finally(() => {
    imagePruneLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}
const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('image.confirmPrune') }}</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('image.cancel') }}</el-button>
                <el-button type="primary" @click="confirmPrune">{{ t('image.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
