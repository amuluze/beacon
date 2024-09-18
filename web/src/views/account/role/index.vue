<template>
    <div class="am-account-header">
        <span @click="$router.push('/account/user')">用户</span>
        <span @click="$router.push('/account/role')">角色</span>
        <span @click="$router.push('/account/api')">接口</span>
        <div></div>
    </div>
    <div class="am-role-operator">
        <el-card shadow="never">
            <el-button type="primary" plain @click="newRoleDraw = true">新增角色</el-button>
        </el-card>
    </div>
    <!-- 表格主体-->
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="data" height="100%" :key="roleKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="name" label="角色名" min-width="120" align="center" />
                <el-table-column prop="resource_ids" label="权限" min-width="200" align="center" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag v-for="(item, index) in scope.row.resources" :key="index">
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
                        <el-button type="primary" size="small" text @click="eidtRole(scope.row)"> 编辑 </el-button>
                        <el-button type="danger" size="small" text @click="deleteRole(scope.row.id)"> 删除 </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </el-card>

    <!-- 创建角色 -->
    <div class="am-user-create">
        <el-drawer v-model="newRoleDraw" title="创建用户" size="50%">
            <el-form ref="userCreateRef" :model="roleCreateMode" :rules="rules" label-width="120px" label-position="left">
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="roleCreateMode.name" placeholder="请输入用户名" />
                </el-form-item>
            </el-form>

            <div class="am-user-create__operator">
                <el-button type="default" size="default" plain @click="newRoleDraw = false">取消</el-button>
                <el-button type="primary" size="default" plain @click="confirmRoleCreate(roleCreateRef)" v-loading="roleCreateLoading"> 确定 </el-button>
            </div>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import { queryRole } from '@/api/account'
import { error } from '@/components/Message/message'
import { useTable } from '@/hooks/useTable'
import { FormInstance } from 'element-plus'

onMounted(() => {
    refresh()
})

// 角色列表
const roleKey = ref(0)
const params = {}
const { data, refresh, loading, pagination } = useTable(queryRole, {}, {})

// 角色创建
const newRoleDraw = ref(false)
const roleCreateRef = ref<FormInstance>()
const roleCreateLoading = ref(false)

interface RuleForm {
    name: string
    resource_ids: string[]
    status: number
}
const roleCreateMode = ref<RuleForm>({
    name: '',
    resource_ids: [],
    status: 1
})

const rules = reactive({
    name: [{ required: true, message: '请输入角色名', trigger: 'blur' }]
})
const confirmRoleCreate = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate((valid, fields) => {
        if (valid) {
            roleCreateLoading.value = true
            // TODO: 创建角色
            roleCreateLoading.value = false
            newRoleDraw.value = false
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

        &:nth-child(2) {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 84px;
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

@include b(role-operator) {
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
