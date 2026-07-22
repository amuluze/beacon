<script setup lang="ts">
import { logout } from '@/api/auth'
import { useI18n } from 'vue-i18n'

import useStore from '@/store'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import UpdatePassword from '@/layout/navbar/components/UpdatePassword.vue'

const { t } = useI18n()
const store = useStore()
const router = useRouter()
const displayName = computed(() => store.user.userInfo.name || 'admin')
const avatarText = computed(() => displayName.value.slice(0, 1).toUpperCase())

async function doLogout() {
  // 1. 退出登录
  await logout()
  // 2.清除缓存
  store.user.setToken('', '')
  store.agent.clear()
  // 3.重定向到登录页
  await router.replace('/login')
}

const updatePasswordDraw = useCommandComponent(UpdatePassword)

function openPasswordDrawer() {
  updatePasswordDraw({
    title: '更新密码',
    username: displayName.value,
  })
}
</script>

<template>
    <el-dropdown trigger="click" placement="bottom">
        <div class="am-avatar-trigger">
            <span class="am-avatar-trigger__avatar">{{ avatarText }}</span>
            <span class="am-avatar-trigger__name">{{ displayName }}</span>
        </div>
        <template #dropdown>
            <el-dropdown-menu>
                <el-dropdown-item @click="openPasswordDrawer">
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
.am-avatar-trigger {
  display: inline-flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  color: var(--am-foreground-primary);
  cursor: pointer;

  &__avatar {
    width: 30px;
    height: 30px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: var(--am-surface-primary);
    background: var(--am-accent-primary);
    border-radius: 50%;
    font-size: var(--am-font-sm);
    font-weight: 700;
  }

  &__name {
    max-width: 120px;
    overflow: hidden;
    font-size: var(--am-font-sm);
    font-weight: 500;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.el-divider {
  margin: 4px;
  width: 90%;
}
</style>
