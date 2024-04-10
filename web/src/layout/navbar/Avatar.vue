<template>
    <el-dropdown>
        <el-avatar size="small">
            <i-ep-user />
        </el-avatar>
        <template #dropdown>
            <el-dropdown-menu>
                <el-dropdown-item @click.prevent="doLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
        </template>
    </el-dropdown>
</template>
<script setup lang="ts">
import { logout } from '@/api/auth'
import useStore from '@/store'

const store = useStore()
const router = useRouter()
const doLogout = async () => {
    // 1. 退出登录
    await logout()
    // 2.清除缓存
    store.user.setToken('', '')
    store.app.isCollapse = false
    // 3.重定向到登录页
    await router.replace('/login')
}
</script>

<style scoped lang="scss">
.el-dropdown {
    :deep(.el-avatar:hover) {
        border: none;
    }
}
</style>
