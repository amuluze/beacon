import type { FormInstance } from 'element-plus'
import { updatePassword } from '@/api/auth'
import { success } from '@/components/Message/message.ts'
import useStore from '@/store'
import { useI18n } from 'vue-i18n'

interface PasswordForm {
  username: string
  old_password: string
  new_password: string
}

const MIN_PASSWORD_LENGTH = 12
const PASSWORD_CLASS_PATTERNS = [/[a-z]/, /[A-Z]/, /\d/, /[^a-z0-9]/i]

export function usePasswordChange(username: string) {
  const { t } = useI18n()
  const store = useStore()
  const router = useRouter()
  const formRef = ref<FormInstance>()
  const loading = shallowRef(false)
  const model = ref<PasswordForm>({
    username,
    old_password: '',
    new_password: '',
  })

  // 与 beacon/pkg/auth/password.go 的 ValidatePassword 同源：≥12 位且至少三类字符。
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

  async function submit(): Promise<boolean> {
    if (loading.value || !formRef.value)
      return false
    const valid = await formRef.value.validate().catch(() => false)
    if (!valid)
      return false

    loading.value = true
    try {
      await updatePassword(model.value)
      success(t('account.passwordUpdated'))
      store.user.clearSession()
      await router.replace('/login')
      return true
    }
    catch {
      return false
    }
    finally {
      loading.value = false
    }
  }

  return {
    formRef,
    loading,
    model,
    rules,
    submit,
  }
}
