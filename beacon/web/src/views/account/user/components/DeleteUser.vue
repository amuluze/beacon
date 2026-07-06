<script setup lang="ts">
import type { UserDeleteArgs } from '@/interface/account.ts'
import { deleteUser } from '@/api/account'
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
// * 用户删除
const userDeleteLoading = ref(false)
async function confirmDelete() {
  userDeleteLoading.value = true
  const params: UserDeleteArgs = {
    ids: [props.id],
  }
  deleteUser(params).finally(() => {
    userDeleteLoading.value = false
    dialogVisible.value = false
    props.update && props.update()
  })
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="500px" draggable>
        <span>{{ t('account.confirmDelete') }}？</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{ t('account.cancel') }}</el-button>
                <el-button type="primary" @click="confirmDelete">{{ t('account.confirm') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
</style>
