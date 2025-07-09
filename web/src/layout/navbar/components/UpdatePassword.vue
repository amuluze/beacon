<script setup lang="ts">
import type { FormInstance } from 'element-plus'
import { updatePassword } from '@/api/auth'
import { success } from '@/components/Message/message.ts'
import useStore from '@/store'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  username: string
  title?: string
}>()
const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

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
const passwordUpdateLoading = ref(false)
const rules = ref({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
  ],

  old_password: [
    { required: true, message: '请输入旧密码', trigger: 'blur' },
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
  ],
})

const store = useStore()
const router = useRouter()
async function confirmPasswordUpdate(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  await formEl.validate(async (valid) => {
    if (valid) {
      passwordUpdateLoading.value = true
      updatePassword(passwordUpdateMode.value).finally(() => {
        passwordUpdateLoading.value = false
        drawerVisible.value = false
        success('更新成功')
        store.user.setToken('', '')
        router.replace('/login')
      })
    }
  })
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" :title="t('account.updatePassword')" size="30%">
        <el-form
            ref="passwordUpdateRef"
            :model="passwordUpdateMode"
            :rules="rules"
            label-width="120px"
            label-position="left"
        >
            <el-form-item :label="t('account.userName')" prop="username">
                <el-input v-model="passwordUpdateMode.username" :placeholder="t('account.inputUserName')" />
            </el-form-item>
            <el-form-item :label="t('account.oldPassword')" prop="old_password">
                <el-input v-model="passwordUpdateMode.old_password" type="password" :placeholder="t('account.inputOldPassword')" />
            </el-form-item>
            <el-form-item :label="t('account.newPassword')" prop="new_password">
                <el-input v-model="passwordUpdateMode.new_password" type="password" :placeholder="t('account.inputNewPassword')" />
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button type="primary" plain @click="drawerVisible = false">
                {{ t('account.cancel') }}
            </el-button>
            <el-button
                v-loading="passwordUpdateLoading"
                type="primary"
                plain
                @click="confirmPasswordUpdate(passwordUpdateRef)"
            >
                {{ t('account.confirm') }}
            </el-button>
        </template>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
