<template>
    <div class="am-account-header">
        <span @click="$router.push('/account/user')">用户</span>
        <span @click="$router.push('/account/role')">角色</span>
        <span @click="$router.push('/account/api')">接口</span>
    </div>
    <div class="am-user-operator">
        <el-card shadow="never">
            <el-button type="primary" plain @click="newUserDraw = true">新增用户</el-button>
        </el-card>
    </div>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="userData" height="100%" :key="userKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="username" label="用户名" min-width="120" align="center" />
                <el-table-column prop="role" label="角色" min-width="160" align="center">
                    <template #default="scope">
                        <el-tag v-for="(item, index) in scope.row.roles" :key="index">
                            {{ item.name }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="status" label="状态" min-width="100" align="center" sortable>
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">正常</el-tag>
                        <el-tag v-else type="danger">禁用</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" min-width="160" align="center" sortable />
                <el-table-column label="操作" width="200" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text @click="userEdit(scope.row)" :disabled="enableEdit(scope.row)"> 编辑 </el-button>
                        <el-button type="danger" size="small" text @click="userDelete(scope.row.id)" v-loading="userDeleteLoading" :disabled="enableEdit(scope.row)"> 删除 </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </el-card>

    <!-- 创建用户 -->
    <div class="am-user-create">
        <el-drawer v-model="newUserDraw" title="创建用户" size="50%">
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

            <div class="am-user-create__operator">
                <el-button type="default" size="default" plain @click="newUserDraw = false">取消</el-button>
                <el-button type="primary" size="default" plain @click="confirmUserCreate(userCreateRef)" v-loading="userCreateLoading"> 确定 </el-button>
            </div>
        </el-drawer>
    </div>

    <!-- 编辑用户 -->
    <div class="am-user-create">
        <el-drawer v-model="userEditDraw" title="编辑用户" size="50%">
            <el-form
                ref="userUpdateRef"
                :model="userUpdateMode"
                :rules="rules"
                label-width="120px"
                label-position="left"
            >
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="userUpdateMode.username" placeholder="请输入用户名" />
                </el-form-item>
                <el-form-item label="角色" prop="role_ids">
                    <el-select v-model="userUpdateMode.role_ids" multiple placeholder="请选择角色">
                        <el-option
                            v-for="item in roleData"
                            :key="item.id"
                            :label="item.name"
                            :value="item.id"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item label="备注" prop="remark">
                    <el-input v-model="userUpdateMode.remark" placeholder="请输入备注" />
                </el-form-item>
                <el-form-item label="状态" prop="status">
                    <el-tooltip content="用户状态，1为正常，2为禁用" placement="top">
                        <el-switch v-model="userUpdateMode.status" active-value="1" inactive-value="2" />
                    </el-tooltip>
                </el-form-item>
            </el-form>

            <div class="am-user-create__operator">
                <el-button type="default" size="default" plain @click="userEditDraw = false">取消</el-button>
                <el-button
                    type="primary"
                    size="default"
                    plain
                    @click="confirmUserUpdate(userUpdateRef)"
                    v-loading="userUpdateLoading"
                >
                    确定
                </el-button>
            </div>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import { createUser, deleteUser, queryRole, queryUser, updateUser } from '@/api/account';
import { Role, User, UserDeleteArgs } from '@/interface/account';
import { FormInstance, FormRules } from 'element-plus';

onMounted(() => {
    userQuery()
    roleQuery()
})

// 用户列表
const userKey = ref(0)
const userData = ref<User[]>([])
const loading = ref(false)
const userQuery = async () => {
    loading.value = true
    const { data } = await queryUser()
    userData.value = data.data
    userKey.value += 1
    loading.value = false
}

// 用户禁止操作
const enableEdit = (user: User) => {
    console.log('......', user)
    return user.username === 'admin'
}

// 用户创建
const roleData = ref<Role[]>([])
const roleQuery = async () => {
    const { data } = await queryRole()
    roleData.value = data.data
}
const newUserDraw = ref(false)
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
    status: '1'
})

const rules = reactive<FormRules<RuleForm>>({
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    role_ids: [{ required: true, message: '请选择角色', trigger: 'blur' }],
})

const userCreateLoading = ref(false)
const confirmUserCreate = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid) => {
        if (valid) {
            userCreateLoading.value = true
            console.log('submit!', userCreateMode.value)
            const params = {
                username: userCreateMode.value.username,
                password: userCreateMode.value.password,
                remark: userCreateMode.value.remark,
                role_ids: userCreateMode.value.role_ids,
                status: Number(userCreateMode.value.status)
            }
            createUser(params).finally(() => {
                userCreateLoading.value = false
                newUserDraw.value = false
                userQuery()
            })

        }
    })
}

// * 用户删除
const userDeleteLoading = ref(false)
const userDelete = async (id: string) => {
    userDeleteLoading.value = true
    const params: UserDeleteArgs = {
        ids: [id]
    }
    deleteUser(params).finally(() => {
        userDeleteLoading.value = false
        userQuery()
    })
}

// * 用户编辑
const userUpdateRef = ref<FormInstance>()
const userEditDraw = ref(false)
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
    status: '1'
})
const userEdit = (user: User) => {
    userEditDraw.value = true
    userUpdateMode.value = {
        id: user.id,
        username: user.username,
        remark: user.remark,
        role_ids: user.roles.map((item) => item.id),
        status: user.status.toString()
    }
}
const userUpdateLoading = ref(false)
const confirmUserUpdate = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid) => {
        if (valid) {
            userUpdateLoading.value = true
            const params = {
                id: userUpdateMode.value.id,
                username: userUpdateMode.value.username,
                remark: userUpdateMode.value.remark,
                role_ids: userUpdateMode.value.role_ids,
                status: Number(userUpdateMode.value.status)
            }
            updateUser(params).finally(() => {
                userUpdateLoading.value = false
                userEditDraw.value = false
                userQuery()
            })
        }
    })
}
</script>

<style scoped lang="scss">
@include b(account-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    height: 48px;
    width: 100%;
    background-color: #fff;

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;

    span {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        font-size: 16px;
        font-weight: 600;
        margin-left: 16px;
        margin-right: 16px;
        color: #424244;
        cursor: pointer;

        &:first-child {
            color: #2f7bff;
            &::before {
                content: '';
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

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 8px 0 0;
            display: flex;
            align-items: center;
            justify-content: flex-end;
        }
    }
}

@include b(user-operator) {
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

@include b(table) {
    width: 100%;
    height: calc(100vh - 188px);
    overflow-y: auto;
}
</style>
