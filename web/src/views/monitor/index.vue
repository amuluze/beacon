<script setup lang="ts">
import { useRoute } from 'vue-router'
import ContainerMonitor from './components/ContainerMonitor.vue'
import HostMonitor from './components/HostMonitor.vue'

const route = useRoute()

const activeTab = computed(() => {
  return route.path.startsWith('/monitor/host') ? 'host' : 'container'
})
</script>

<template>
  <div class="monitor-page">
    <el-tabs v-model="activeTab" class="monitor-tabs">
      <el-tab-pane label="主机监控" name="host">
        <HostMonitor />
      </el-tab-pane>
      <el-tab-pane label="容器监控" name="container">
        <ContainerMonitor />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped lang="scss">
.monitor-page {
  height: 100%;

  .monitor-tabs {
    height: 100%;

    :deep(.el-tabs__content) {
      height: calc(100% - 40px);
      overflow: auto;
    }
  }
}
</style>
