<script setup lang="ts">
import API from './api/index.vue';
import Role from './role/index.vue';
import User from './user/index.vue';

const activeTab = ref('user');
function getActiveComponent() {
    switch (activeTab.value) {
        case 'user':
            return User;
        case 'role':
            return Role;
        case 'api':
            return API;
        default:
            return User;
    }
}
</script>

<template>
    <div class="am-account-container">
        <div class="am-container-header">
            <div class="am-button-group">
                <div class="am-button-group__tab">
                    <div class="tab-item" :class="{ active: activeTab === 'user' }" @click="activeTab = 'user'">
                        用户管理
                    </div>
                    <div class="tab-item" :class="{ active: activeTab === 'role' }" @click="activeTab = 'role'">
                        角色管理
                    </div>
                    <div class="tab-item" :class="{ active: activeTab === 'api' }" @click="activeTab = 'api'">
                        接口管理
                    </div>
                </div>
            </div>
        </div>

        <div class="am-container-content">
            <component :is="getActiveComponent()" />
        </div>
    </div>
</template>

<style scoped lang="scss">
@include b(account-container) {
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: var(--el-bg-color);
}

@include b(container-header) {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    gap: 16px;
    background-color: var(--el-bg-color);
}

@include b(button-group) {
    background: var(--el-bg-color-page);
    border-radius: 12px;
    padding: 4px;

    @include e(tab) {
        display: flex;
        background-color: var(--el-fill-color-light);
        border-radius: 8px;
        padding: 4px;
        height: 45px;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

        .tab-item {
            padding: 8px 16px;
            cursor: pointer;
            font-size: 14px;
            color: var(--el-text-color-regular);
            border-radius: 6px;
            transition: all 0.3s;
            position: relative;
            white-space: nowrap;

            &:hover:not(.active) {
                color: var(--el-color-primary);
                background-color: rgba(var(--el-color-primary-rgb), 0.05);
            }

            &.active {
                color: var(--el-color-primary);
                background-color: var(--el-menu-bg-color);
                box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
                font-weight: 500;
            }
        }
    }
}

@include b(container-content) {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;

    &::-webkit-scrollbar {
        width: 4px;
    }

    &::-webkit-scrollbar-thumb {
        background: var(--el-border-color-darker);
        border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
        background: var(--el-border-color-light);
        border-radius: 2px;
    }
}
</style>
