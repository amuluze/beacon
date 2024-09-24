<template>
    <div class="am-profile-container">
        <div class="am-profile-container__title">
            <span>个人中心</span>
        </div>
    </div>
    <div class="am-profile-operator">
        <el-card shadow="never">
            <el-button type="primary" plain @click="passwordUpdateDraw = true">更新密码</el-button>
        </el-card>
    </div>
    <el-card shadow="never">
        <p>
            用户名：<el-tag>{{ profile?.username }}</el-tag>
        </p>
        <p>
            状态：<el-tag>{{ profile?.status == '1' ? '启用' : '禁用'}}</el-tag>
        </p>
        <p>
            角色：<el-tag>{{ profile?.is_admin == 1 ? '管理员' : '普通用户'}}</el-tag>
        </p>
    </el-card>
    <el-drawer v-model="passwordUpdateDraw" title="更新密码" size="30%">
        <el-form
            ref="passwordUpdateRef"
            :model="passwordUpdateMode"
            :rules="rules"
            label-width="120px"
            label-position="left"
        >
            <el-form-item label="用户名" prop="username">
                <el-input v-model="passwordUpdateMode.username" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="旧密码" prop="old_password">
                <el-input v-model="passwordUpdateMode.old_password" type="password" placeholder="请输入旧密码" />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
                <el-input v-model="passwordUpdateMode.new_password" type="password" placeholder="请输入新密码" />
            </el-form-item>
        </el-form>

        <div class="am-user-create__operator">
            <el-button type="default" size="default" plain @click="passwordUpdateDraw = false">取消</el-button>
            <el-button
                type="primary"
                size="default"
                plain
                @click="confirmPasswordUpdate(passwordUpdateRef)"
                v-loading="passwordUpdateLoading"
            >
                确定
            </el-button>
        </div>
    </el-drawer>
</template>
<script setup lang="ts">
import { getUserInfo, updatePassword } from '@/api/auth';
import { success } from '@/components/Message/message';
import useStore from '@/store';
import { FormInstance } from 'element-plus';

onMounted(() => {
    getProfile()
})

const profile = ref({
    username: '',
    status: '1',
    is_admin: 1
})
const getProfile = async () => {
    const { data } = await getUserInfo()
    profile.value = {
        username: data.username,
        status: data.status,
        is_admin: data.is_admin
    }
}
// * 更新密码
const passwordUpdateDraw = ref(false)
const passwordUpdateRef = ref<FormInstance>()
interface RuleForm {
    username: string;
    old_password: string;
    new_password: string;
}
const passwordUpdateMode = ref<RuleForm>({
    username: profile.value.username,
    old_password: '',
    new_password: ''
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
const confirmPasswordUpdate = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid) => {
        if (valid) {
            passwordUpdateLoading.value = true
            updatePassword(passwordUpdateMode.value).finally(() => {
                passwordUpdateLoading.value = false
                passwordUpdateDraw.value = false
                success('更新成功')
                store.user.setToken('', '')
                router.replace('/app')
            })
        }
    })
}
</script>
<style scoped lang="scss">
@include b(profile-container) {
    font-size: 14px;

    @include e(title) {
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 48px;
        width: 100%;
        background-color: #ffffff;
        // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

        border-radius: 4px;
        margin-bottom: 8px;
        padding: 0 16px;
        span {
            display: flex;
            align-items: center;
            font-size: 16px;
            font-weight: 600;
            color: #105eeb;
            margin-left: 16px;
            &::before {
                content: ' ';
                position: absolute;
                left: 20px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
@include b(profile-operator) {
    height: 48px;
    width: 100%;
    margin-bottom: 4px;
    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 0 0 16px;
            display: flex;
            align-items: center;
            justify-content: flex-start;
        }
    }
}
</style>
