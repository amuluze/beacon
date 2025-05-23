<script setup lang="ts">
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
</script>

<template>
  <div class="monitor-page">
    <div class="am-header-container">
      <div class="am-button-group">
        <div class="tab-nav-container">
          <div
            class="tab-nav-item"
            :class="{ active: activeTab === 'host' }"
            @click="activeTab = 'host'"
          >
            主机监控
          </div>
          <div
            class="tab-nav-item"
            :class="{ active: activeTab === 'container' }"
            @click="activeTab = 'container'"
          >
            容器监控
          </div>
        </div>
      </div>

      <div class="am-density-selector" v-if="activeTab === 'host'">
        <div class="tab-nav-container">
          <div class="tab-nav-item density-item">
            <el-select
              v-model="timeDensity"
              placeholder="Select"
              size="default"
              class="density-select"
              :popper-append-to-body="false"
              :popper-style="{ minWidth: '200px' }"
              auto-width>
              <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </div>
        </div>
      </div>
    </div>

    <div class="monitor-content">
      <component :is="activeTab === 'host' ? HostMonitor : ContainerMonitor" :time-density="timeDensity" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.monitor-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--el-bg-color);
}

.am-header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 16px;
//   padding: 16px 24px;
  background-color: var(--el-bg-color);
//   border-bottom: 1px solid var(--el-border-color-light);

  .el-card {
    border-radius: 12px;
    background: var(--el-bg-color-page);

    :deep(.el-card__body) {
      padding: 8px 16px !important;
    }
  }
}

.am-button-group {
  background: var(--el-bg-color-page);
  border-radius: 12px;
  padding: 4px;

  .tab-nav-container {
    display: flex;
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 4px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .tab-nav-item {
    padding: 8px 16px;
    cursor: pointer;
    font-size: 14px;
    color: #606266;
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
      background-color: #fff;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
      font-weight: 500;
    }
  }
}

.am-density-selector {
  background: var(--el-bg-color-page);
  border-radius: 12px;
  padding: 4px;

  .tab-nav-container {
    display: flex;
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 4px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .tab-nav-item.density-item {
    padding: 4px 8px;
    background-color: #fff;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
    border-radius: 6px;
  }
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
    height: 32px;
    line-height: 32px;
  }
}

.monitor-content {
  flex: 1;
//   padding: 24px;
  overflow: auto;
}
</style>
