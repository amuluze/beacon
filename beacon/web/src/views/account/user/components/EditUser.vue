<script setup lang="ts">
import type { Role, User } from '@/interface/account.ts'
import type { FormInstance, FormRules } from 'element-plus'
import { queryRole, updateUser } from '@/api/account'
import DialogHeader from '@/components/DialogHeader/index.vue'
import { useI18n } from 'vue-i18n'
import IconUser from '~icons/lucide/user'

const props = defineProps<{
  visible: boolean
  entity: User
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

// * 用户编辑
const roleData = ref<Role[]>([])
async function roleQuery() {
  const { data } = await queryRole()
  roleData.value = data.data
}

const userUpdateRef = ref<FormInstance>()
const userUpdateMode = ref<{
  id: string
  username: string
  remark: string
  role_ids: string[]
  status: string
}>({
  id: '',
  username: '',
  remark: '',
  role_ids: [],
  status: '1',
})

interface RuleForm {
  username: string
  password: string
  remark: string
  role_ids: string[]
  status: string
}

const rules = reactive<FormRules<RuleForm>>({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role_ids: [{ required: true, message: '请选择角色', trigger: 'blur' }],
})

onMounted(async () => {
  await roleQuery()
  userUpdateMode.value = {
    id: props.entity.id,
    username: props.entity.username,
    remark: props.entity.remark,
    role_ids: props.entity.roles.map(item => item.id),
    status: props.entity.status.toString(),
  }
})

const userUpdateLoading = ref(false)
async function confirmUserUpdate(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  await formEl.validate((valid) => {
    if (valid) {
      userUpdateLoading.value = true
      const params = {
        id: userUpdateMode.value.id,
        username: userUpdateMode.value.username,
        remark: userUpdateMode.value.remark,
        role_ids: userUpdateMode.value.role_ids,
        status: Number(userUpdateMode.value.status),
      }
      updateUser(params).finally(() => {
        userUpdateLoading.value = false
        drawerVisible.value = false
        props.update && props.update()
      })
    }
  })
}

const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" :title="t('account.editAccount')" size="720px">
        <template #header>
            <DialogHeader :title="t('account.editAccount')" :icon="IconUser" />
        </template>
        <el-form
            ref="userUpdateRef"
            :model="userUpdateMode"
            :rules="rules"
            label-position="top"
        >
            <el-form-item :label="t('account.userName')" prop="username">
                <el-input v-model="userUpdateMode.username" :placeholder="t('account.inputUserName')" />
            </el-form-item>
            <el-form-item :label="t('account.roleName')" prop="role_ids">
                <el-select v-model="userUpdateMode.role_ids" multiple :placeholder="t('account.selectRole')">
                    <el-option
                        v-for="item in roleData"
                        :key="item.id"
                        :label="item.name"
                        :value="item.id"
                    />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('account.remark')" prop="remark">
                <el-input v-model="userUpdateMode.remark" :placeholder="t('account.inputRemark')" />
            </el-form-item>
            <el-form-item :label="t('account.status')" prop="status">
                <el-tooltip :content="t('account.statusHint')" placement="top">
                    <el-switch v-model="userUpdateMode.status" active-value="1" inactive-value="2" />
                </el-tooltip>
            </el-form-item>
        </el-form>

        <div class="am-user-create__operator">
            <el-button type="default" size="default" @click="drawerVisible = false">
                {{ t('account.cancel') }}
            </el-button>
            <el-button
                v-loading="userUpdateLoading"
                type="primary"
                size="default"
                @click="confirmUserUpdate(userUpdateRef)"
            >
                {{ t('account.confirm') }}
            </el-button>
        </div>
    </el-drawer>
</template>

<style scoped lang="scss">
// §11.1 抽屉底部：上边框 + 右对齐
.am-user-create__operator {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--am-spacing-sm);
  padding: 16px 24px;
  border-top: 1px solid var(--am-border-subtle);
}
</style>
