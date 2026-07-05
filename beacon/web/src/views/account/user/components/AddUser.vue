<script setup lang="ts">
import { createUser, queryRole } from '@/api/account'
import type { FormInstance, FormRules } from 'element-plus'
import type { Role } from '@/interface/account.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  update?: () => void
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

// 用户创建
const roleData = ref<Role[]>([])
async function roleQuery() {
  const { data } = await queryRole()
  roleData.value = data.data
}
onMounted(() => {
  roleQuery()
})

const userCreateRef = ref<FormInstance>()
interface RuleForm {
  username: string
  password: string
  remark: string
  role_ids: string[]
  status: string
}

const userCreateMode = ref<RuleForm>({
  username: '',
  password: '',
  remark: '',
  role_ids: [],
  status: '1',
})

const rules = reactive<FormRules<RuleForm>>({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role_ids: [{ required: true, message: '请选择角色', trigger: 'blur' }],
})

const userCreateLoading = ref(false)
async function confirmUserCreate(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  await formEl.validate((valid) => {
    if (valid) {
      userCreateLoading.value = true
      console.log('submit!', userCreateMode.value)
      const params = {
        username: userCreateMode.value.username,
        password: userCreateMode.value.password,
        remark: userCreateMode.value.remark,
        role_ids: userCreateMode.value.role_ids,
        status: Number(userCreateMode.value.status),
      }
      createUser(params).finally(() => {
        userCreateLoading.value = false
        drawerVisible.value = false
        props.update && props.update()
      })
    }
  })
}

const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" title="创建用户" size="50%">
        <el-form ref="userCreateRef" :model="userCreateMode" :rules="rules" label-width="120px" label-position="left">
            <el-form-item label="用户名" prop="username">
                <el-input v-model="userCreateMode.username" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
                <el-input v-model="userCreateMode.password" placeholder="请输入密码" />
            </el-form-item>
            <el-form-item label="角色" prop="role_ids">
                <el-select v-model="userCreateMode.role_ids" multiple placeholder="请选择角色">
                    <el-option
                        v-for="item in roleData"
                        :key="item.id"
                        :label="item.name"
                        :value="item.id"
                    />
                </el-select>
            </el-form-item>
            <el-form-item label="备注" prop="remark">
                <el-input v-model="userCreateMode.remark" placeholder="请输入备注" />
            </el-form-item>
            <el-form-item label="状态" prop="status">
                <el-tooltip content="用户状态，1为正常，2为禁用" placement="top">
                    <el-switch v-model="userCreateMode.status" active-value="1" inactive-value="2" />
                </el-tooltip>
            </el-form-item>
        </el-form>

        <template #footer>
            <span>
                <el-button @click="drawerVisible = false">{{ t('account.cancel') }}</el-button>
                <el-button type="primary" @click="confirmUserCreate(userCreateRef)"> {{ t('account.confirm') }}</el-button>
            </span>
        </template>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
