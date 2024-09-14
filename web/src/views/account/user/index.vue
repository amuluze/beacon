<template>
    <div class="am-account-header">
        <div class="am-account-header__menu">
            <span @click="$router.push('/account/user')">用户</span>
            <span @click="$router.push('/account/role')">角色</span>
            <span @click="$router.push('/account/api')">接口</span>
        </div>
        <div>
            <el-button type="primary" @click="newUserDraw = true">新增用户</el-button>
        </div>
    </div>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="data" :key="userKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="username" label="用户名" min-width="120" align="center" />
                <el-table-column prop="role" label="角色" min-width="160" align="center" />
                <el-table-column prop="status" label="状态" min-width="100" align="center">
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">正常</el-tag>
                        <el-tag v-else type="danger">禁用</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="role" label="创建时间" min-width="160" align="center" />
            </el-table>
        </div>

        <div class="am-pagination">
            <el-config-provider :locale="zhCn">
                <el-pagination
                    v-model:current-page="pagination.page"
                    :page-size="pagination.size"
                    :page-sizes="pagination.sizeOption"
                    :total="pagination.total"
                    layout="total, sizes, prev, pager, next, jumper"
                    @size-change="(size: number) => pagination.onPageChange(size, params)"
                    @current-change="(page: number) => pagination.onPageChange(page, params)"
                />
            </el-config-provider>
        </div>
    </el-card>

    <!-- 创建用户 -->
    <div class="am-user-create">
        <el-drawer v-model="newUserDraw" title="创建用户" size="50%">
            <el-form
                ref="userCreateRef"
                :model="userCreateMode"
                :rules="rules"
                label-width="120px"
                label-position="left"
            >
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="userCreateMode.username" placeholder="请输入用户名" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="userCreateMode.password" placeholder="请输入密码" />
                </el-form-item>
            </el-form>

            <div class="am-user-create__operator">
                <el-button type="default" size="default" plain @click="newUserDraw = false">取消</el-button>
                <el-button
                    type="primary"
                    size="default"
                    plain
                    @click="confirmUserCreate(userCreateRef)"
                    v-loading="userCreateLoading"
                >
                    确定
                </el-button>
            </div>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import { queryUser } from '@/api/account'
import { error } from '@/components/Message/message'
import { useTable } from '@/hooks/useTable'
import { FormInstance, FormRules } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

onMounted(() => {
    refresh()
})

// 用户列表
const userKey = ref(0)
const params = {}
const { data, refresh, loading, pagination } = useTable(queryUser, {}, {})

// 用户创建
const newUserDraw = ref(false)
const userCreateRef = ref<FormInstance>()
interface RuleForm {
    username: string
    password: string
    role_ids: string[]
    status: number
}

const userCreateMode = ref<RuleForm>({
    username: '',
    password: '',
    role_ids: [],
    status: 1
})

const rules = reactive<FormRules<RuleForm>>({
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

const userCreateLoading = ref(false)
const confirmUserCreate = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
        if (valid) {
            userCreateLoading.value = true
            setTimeout(() => {
                userCreateLoading.value = false
                newUserDraw.value = false
            }, 2000)
        } else {
            console.log('error submit!', fields)
            error('请检查表单')
            return
        }
    })
}
</script>

<style scoped lang="scss">
@include b(account-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    width: 100%;
    background-color: #fff;

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;

    @include e(menu) {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;

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
}
</style>
