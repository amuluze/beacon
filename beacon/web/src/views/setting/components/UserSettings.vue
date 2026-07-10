<script setup lang="ts">
import ApiManager from '@/views/account/api/index.vue'
import RoleManager from '@/views/account/role/index.vue'
import UserManager from '@/views/account/user/index.vue'

const activeTab = shallowRef('users')
</script>

<template>
    <el-tabs v-model="activeTab" class="user-settings">
        <el-tab-pane :label="$t('menu.userManager')" name="users" lazy>
            <div class="user-settings__content">
                <UserManager />
            </div>
        </el-tab-pane>
        <el-tab-pane :label="$t('menu.roleManager')" name="roles" lazy>
            <div class="user-settings__content">
                <RoleManager />
            </div>
        </el-tab-pane>
        <el-tab-pane :label="$t('menu.apiManager')" name="apis" lazy>
            <div class="user-settings__content">
                <ApiManager />
            </div>
        </el-tab-pane>
    </el-tabs>
</template>

<style scoped lang="scss">
.user-settings__content {
  height: 560px;
  min-height: 0;

  :deep(.am-container) {
    padding: var(--am-spacing-md) 0 0;
    background: transparent;
  }
}

@media (max-width: 900px) {
  .user-settings__content {
    height: 520px;
  }
}
</style>
