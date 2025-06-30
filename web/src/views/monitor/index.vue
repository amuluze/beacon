<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import ContainerMonitor from './components/ContainerMonitor.vue';
import HostMonitor from './components/HostMonitor.vue';

const activeTab = ref('host')

// 时间密度下拉框
const timeDensity = ref(600)
const options = [
    {
        value: 120,
        label: '2分钟',
    },
    {
        value: 300,
        label: '5分钟',
    },
    {
        value: 600,
        label: '10分钟',
    },
    {
        value: 1800,
        label: '30分钟',
    },
    {
        value: 3600,
        label: '1 小时',
    },
    {
        value: 43200,
        label: '12小时',
    },
    {
        value: 86400,
        label: '24小时',
    },
]

const { t } = useI18n()
</script>

<template>
    <div class="am-monitor-page">
        <div class="am-header-container">
            <div class="am-button-group">
                <div class="am-button-group__tab">
                    <div class="tab-item" :class="{ active: activeTab === 'host' }" @click="activeTab = 'host'">
                        {{ t('menu.hostMonitor') }}
                    </div>
                    <div class="tab-item" :class="{ active: activeTab === 'container' }"
                        @click="activeTab = 'container'">
                        {{ t('menu.containerMonitor') }}
                    </div>
                </div>
            </div>

            <div class="am-density-selector" v-if="activeTab === 'host'">
                <div class="am-density-selector__tab">
                    <div class="density-item">
                        <el-select v-model="timeDensity" placeholder="Select" size="default" class="density-select"
                            :popper-append-to-body="false" :popper-style="{ minWidth: '200px' }">
                            <el-option v-for="item in options" :key="item.value" :label="item.label"
                                :value="item.value" />
                        </el-select>
                    </div>
                </div>
            </div>
        </div>

        <div class="am-monitor-content">
            <component :is="activeTab === 'host' ? HostMonitor : ContainerMonitor" :time-density="timeDensity" />
        </div>
    </div>
</template>

<style scoped lang="scss">
@include b(monitor-page) {
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: var(--el-bg-color);
}

@include b(header-container) {
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

@include b(density-selector) {
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

        .density-item {
            margin: auto 0;
            background-color: var(--el-bg-color);
            box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
            border-radius: 6px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .density-select {
            width: auto;
            min-width: 160px;

            :deep(.el-input__wrapper) {
                box-shadow: none !important;
                padding: 0 8px;
            }

            :deep(.el-input__inner) {
                min-width: 160px;
                height: 28px;
                line-height: 28px;
                color: var(--el-text-color-primary);
            }
        }
    }
}

@include b(monitor-content) {
    flex: 1;
    overflow: auto;
}
</style>
