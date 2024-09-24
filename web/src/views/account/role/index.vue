<template>
    <div class="am-account-header">
        <span @click="$router.push('/account/user')">用户</span>
        <span @click="$router.push('/account/role')">角色</span>
        <span @click="$router.push('/account/api')">接口</span>
        <div></div>
    </div>
    <!-- 表格主体-->
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="roleData" height="100%" :key="roleKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="name" label="角色名" min-width="120" align="center" />
                <el-table-column prop="resource_ids" label="权限" min-width="200" align="center" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tree :data="generageTree(scope.row.resources)" :props="defaultProps" />
                        <!-- <el-tag v-for="(item, index) in scope.row.resources" :key="index">
                            {{ item.name }}
                        </el-tag> -->
                    </template>
                </el-table-column>
                <el-table-column prop="status" label="状态" min-width="100" align="center" sortable>
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">正常</el-tag>
                        <el-tag v-else type="danger">禁用</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" min-width="160" align="center" sortable />
            </el-table>
        </div>
    </el-card>
</template>
<script setup lang="ts">
import { queryRole } from '@/api/account';
import { Resource, Role } from '@/interface/account';

onMounted(() => {
    roleQuery()
})

// 角色列表
const roleKey = ref(0)
const roleData = ref<Role[]>([])
const loading = ref<boolean>(false)
const defaultProps = {
    children: 'children',
    label: 'label',
    id: 'id'
}

interface Tree {
    id: string
    label: string
    children?: Tree[]
}
const generageTree = (data: Resource[]) => {
    let child: Tree[] = []
    data.forEach((item) => {
        const treeItem: Tree = {
            id: item.id,
            label: item.name,
            children: []
        }
        child.push(treeItem)
    })
    const tree: Tree[] = [{
        id: '',
        label: '权限列表',
        children: child
    }]
    console.log(">>", tree)
    return tree
}
const roleQuery = async () => {
    loading.value = true
    const { data } = await queryRole()
    roleData.value = data.data
    roleKey.value += 1
    loading.value = false
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

@include b(table) {
    width: 100%;
    height: calc(100vh - 136px);
    overflow-y: auto;
}
</style>
