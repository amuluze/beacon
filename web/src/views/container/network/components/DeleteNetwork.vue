<script setup lang="ts">
import type { NetworkDeleteArgs } from '@/interface/container.ts'
import { useI18n } from 'vue-i18n'
import { removeNetwork } from '@/api/container'

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

// 删除网络
const networkDeleteLoading = ref(false)
async function confirmDelete() {
  networkDeleteLoading.value = true
  const params: NetworkDeleteArgs = {
    network_id: props.id,
  }
  removeNetwork(params).finally(() => {
    networkDeleteLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span> {{ t('network.confirmDelete') }}？</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('network.cancel') }}</el-button>
                <el-button type="primary" @click="confirmDelete">{{ t('network.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
