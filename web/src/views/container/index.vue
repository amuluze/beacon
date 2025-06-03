<script setup lang="ts">
import Container from './container/index.vue';
import Image from './image/index.vue';
import Network from './network/index.vue';

const activeTab = ref('container');
function getActiveComponent() {
    switch (activeTab.value) {
        case 'container':
            return Container;
        case 'image':
            return Image;
        case 'network':
            return Network;
        default:
            return Container;
    }
}
</script>
<template>
    <div class="am-container-container">
        <div class="am-container-header">
            <div class="am-button-group">
                <div class="am-button-group__tab">
                    <div class="tab-item" :class="{ active: activeTab === 'container' }"
                        @click="activeTab = 'container'">
                        容器管理
                    </div>
                    <div class="tab-item" :class="{ active: activeTab === 'image' }" @click="activeTab = 'image'">
                        镜像管理
                    </div>
                    <div class="tab-item" :class="{ active: activeTab === 'network' }" @click="activeTab = 'network'">
                        网络管理
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
@include b(container-container) {
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
    overflow: auto;
    display: flex;
    flex-direction: column;
}
</style>
