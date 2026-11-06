<script setup lang="ts">
import type { FormInstance } from 'element-plus'
import { updatePassword } from '@/api/auth'
import { success } from '@/components/Message/message.ts'
import useStore from '@/store'
import { useI18n } from 'vue-i18n'

const props = withDefaults(defineProps<{
  visible: boolean
  username: string
  mandatory?: boolean
}>(), {
  mandatory: false,
})
const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const { t } = useI18n()

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

// 更新密码
const passwordUpdateRef = ref<FormInstance>()
interface RuleForm {
  username: string
  old_password: string
  new_password: string
}
const passwordUpdateMode = ref<RuleForm>({
  username: props.username,
  old_password: '',
  new_password: '',
})
const passwordUpdateLoading = shallowRef(false)

// 与 beacon/pkg/auth/password.go 的 ValidatePassword 同源：≥12 位且至少三类字符。
const MIN_PASSWORD_LENGTH = 12
const PASSWORD_CLASS_PATTERNS = [/[a-z]/, /[A-Z]/, /\d/, /[^a-z0-9]/i]

function validateNewPassword(_: unknown, value: string, callback: (error?: Error) => void) {
  if (!value) {
    callback(new Error(t('account.inputNewPassword')))
    return
  }
  if ([...value].length < MIN_PASSWORD_LENGTH) {
    callback(new Error(t('account.passwordTooShort')))
    return
  }
  const classes = PASSWORD_CLASS_PATTERNS.filter(pattern => pattern.test(value)).length
  if (classes < 3) {
    callback(new Error(t('account.passwordTooWeak')))
    return
  }
  callback()
}

const rules = computed(() => ({
  username: [
    { required: true, message: t('account.usernameRequired'), trigger: 'blur' },
  ],
  old_password: [
    { required: true, message: t('account.oldPasswordRequired'), trigger: 'blur' },
  ],
  new_password: [
    { required: true, validator: validateNewPassword, trigger: 'blur' },
  ],
}))

const store = useStore()
const router = useRouter()
async function confirmPasswordUpdate(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  const valid = await formEl.validate().catch(() => false)
  if (!valid)
    return

  passwordUpdateLoading.value = true
  try {
    await updatePassword(passwordUpdateMode.value)
    drawerVisible.value = false
    success(t('account.passwordUpdated'))
    store.user.clearSession()
    await router.replace('/login')
  }
  catch {
  }
  finally {
    passwordUpdateLoading.value = false
  }
}
</script>

<template>
    <el-drawer
        v-model="drawerVisible"
        :title="t('account.changePassword')"
        size="480px"
        :close-on-click-modal="!props.mandatory"
        :close-on-press-escape="!props.mandatory"
        :show-close="!props.mandatory"
    >
        <p class="am-password-hint">
            {{ props.mandatory ? t('account.passwordPolicyMandatory') : t('account.passwordPolicy') }}
        </p>
        <el-form
            ref="passwordUpdateRef"
            :model="passwordUpdateMode"
            :rules="rules"
            label-position="top"
            class="am-password-form"
        >
            <el-form-item :label="t('account.userName')" prop="username">
                <el-input v-model="passwordUpdateMode.username" disabled :placeholder="t('account.inputUserName')" />
            </el-form-item>
            <el-form-item :label="t('account.oldPassword')" prop="old_password">
                <el-input v-model="passwordUpdateMode.old_password" type="password" :placeholder="t('account.inputOldPassword')" />
            </el-form-item>
            <el-form-item :label="t('account.newPassword')" prop="new_password">
                <el-input v-model="passwordUpdateMode.new_password" type="password" :placeholder="t('account.inputNewPassword')" />
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button v-if="!props.mandatory" @click="drawerVisible = false">
                {{ t('account.cancel') }}
            </el-button>
            <el-button
                v-loading="passwordUpdateLoading"
                type="primary"
                @click="confirmPasswordUpdate(passwordUpdateRef)"
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
