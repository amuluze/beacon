<script setup lang="ts">
import { logout } from '@/api/auth'
import { useI18n } from 'vue-i18n'

import useStore from '@/store'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import UpdatePassword from '@/layout/navbar/components/UpdatePassword.vue'

const { t } = useI18n()
const store = useStore()
const router = useRouter()
async function doLogout() {
  // 1. 退出登录
  await logout()
  // 2.清除缓存
  store.user.setToken('', '')
  // 3.重定向到登录页
  await router.replace('/login')
}

function toProfile() {
  router.push('/profile')
}

const updatePasswordDraw = useCommandComponent(UpdatePassword)
</script>

<template>
    <el-dropdown trigger="click" placement="bottom">
        <el-avatar size="small">
            <i-ep-user />
        </el-avatar>
        <template #dropdown>
            <el-dropdown-menu>
                <el-dropdown-item @click.prevent="toProfile">
                    <svg-icon icon-class="people" style="margin-right: 4px" />
                    {{ t('avatar.profile') }}
                </el-dropdown-item>
                <el-dropdown-item @click="updatePasswordDraw({ title: '更新密码' })">
                    <svg-icon icon-class="edit" style="margin-right: 4px" />
                    {{ t('avatar.updatePassword') }}
                </el-dropdown-item>
                <el-divider />
                <el-dropdown-item @click.prevent="doLogout">
                    <svg-icon icon-class="power" style="margin-right: 4px" />
                    {{ t('avatar.logout') }}
                </el-dropdown-item>
            </el-dropdown-menu>
        </template>
    </el-dropdown>
</template>

<style scoped lang="scss">
.el-divider {
  margin: 4px;
  width: 90%;
}
</style>
