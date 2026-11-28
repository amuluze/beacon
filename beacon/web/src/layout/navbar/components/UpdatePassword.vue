<script setup lang="ts">
import { usePasswordChange } from '@/hooks/usePasswordChange'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  username: string
}>()
const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const { t } = useI18n()
const {
  formRef,
  loading,
  model,
  rules,
  submit,
} = usePasswordChange(props.username)

const drawerVisible = computed<boolean>({
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

async function confirmPasswordUpdate() {
  if (await submit())
    drawerVisible.value = false
}
</script>

<template>
    <el-drawer
        v-model="drawerVisible"
        :title="t('account.changePassword')"
        size="480px"
    >
        <p class="am-password-hint">
            {{ t('account.passwordPolicy') }}
        </p>
        <el-form
            ref="formRef"
            :model="model"
            :rules="rules"
            label-position="top"
            class="am-password-form"
        >
            <el-form-item :label="t('account.userName')" prop="username">
                <el-input v-model="model.username" disabled :placeholder="t('account.inputUserName')" />
            </el-form-item>
            <el-form-item :label="t('account.oldPassword')" prop="old_password">
                <el-input v-model="model.old_password" type="password" :placeholder="t('account.inputOldPassword')" />
            </el-form-item>
            <el-form-item :label="t('account.newPassword')" prop="new_password">
                <el-input v-model="model.new_password" type="password" :placeholder="t('account.inputNewPassword')" />
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button @click="drawerVisible = false">
                {{ t('account.cancel') }}
            </el-button>
            <el-button
                v-loading="loading"
                type="primary"
                @click="confirmPasswordUpdate"
            >
                {{ t('account.confirm') }}
            </el-button>
        </template>
    </el-drawer>
</template>

<style scoped lang="scss">
.am-password-hint {
  margin: 0 0 var(--am-spacing-lg);
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}

.am-password-form {
  :deep(.el-form-item) {
    margin-bottom: var(--am-spacing-md);
  }

  :deep(.el-form-item__label) {
    padding-bottom: 4px;
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-sm);
  }

  :deep(.el-input__wrapper) {
    min-height: 36px;
  }
}
</style>
