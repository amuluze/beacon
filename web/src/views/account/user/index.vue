<script setup lang="ts">
import { queryUser } from '@/api/account'
import type { User } from '@/interface/account'
import type { TableInstance } from 'element-plus'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import AddUser from '@/views/account/user/components/AddUser.vue'
import DeleteUser from '@/views/account/user/components/DeleteUser.vue'
import EditUser from '@/views/account/user/components/EditUser.vue'

import { useTable } from '@/hooks/useTable.ts'
import { useI18n } from 'vue-i18n'

const tableRef = ref<TableInstance>()
const tableSelection = ref<User[]>([])

const selectable = (row: User) => !['1', '2'].includes(row.id)

function handleSelectionChange(val: User[]) {
    console.log('selection: ', val)
    tableSelection.value = val
}

const { tableData, pageable, loading, search, handleCurrentChange, handleSizeChange } = useTable(queryUser)

onMounted(async () => {
    await search()
})

// 用户禁止操作
function enableEdit(user: User) {
    console.log('......', user)
    return user.username === 'admin'
}

const addUser = useCommandComponent(AddUser)
const deleteUser = useCommandComponent(DeleteUser)
const editUser = useCommandComponent(EditUser)

const { t } = useI18n()
</script>

<template>
    <div class="am-container">
        <div class="am-table-operator">
            <el-button type="primary" plain size="small"
                @click="addUser({ title: 'account.newAccount', update: search })">
                <svg-icon icon-class="add" />
                {{ t('account.newAccount') }}
            </el-button>
        </div>
        <div class="am-table">
            <el-table ref="tableRef" v-loading="loading" :data="tableData as User[]"
                :header-cell-style="{ height: '45px', fontSize: '14px', color: '#000' }" border
                @selection-change="handleSelectionChange">
                <el-table-column type="selection" :selectable="selectable" width="55" />
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
                        <el-tag v-if="scope.row.status === 1" type="success">
                            正常
                        </el-tag>
                        <el-tag v-else type="danger">
                            禁用
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" min-width="160" align="center" sortable />
                <el-table-column label="操作" width="200" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text :disabled="enableEdit(scope.row)"
                            @click="editUser({ title: 'account.editAccount', entity: scope.row, update: search })">
                            <svg-icon icon-class="edit" />
                            编辑
                        </el-button>
                        <el-button type="danger" size="small" text :disabled="enableEdit(scope.row)"
                            @click="deleteUser({ title: 'account.deleteAccount', ids: [scope.row.id], update: search })">
                            <svg-icon icon-class="delete" />
                            删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="am-pagination">
            <el-pagination v-model:current-page="pageable.page" :page-size="pageable.size"
                layout="total, sizes, prev, pager, next, jumper" :page-sizes="pageable.options" :total="pageable.total"
                @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
    </div>
</template>

<style scoped lang="scss"></style>
