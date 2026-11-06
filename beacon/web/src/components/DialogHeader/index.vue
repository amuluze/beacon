<script setup lang="ts">
import type { Component } from 'vue'

// 稿面 `BeaconDialogHeader`（DESIGN.md §11.1）：28×28 soft 图标块 + 标题（14/600）。
// 关闭入口由 Element Plus 的 header 关闭按钮承担，不在此组件内重复渲染。
defineProps<{
  title: string
  icon: Component
}>()
</script>

<template>
    <div class="am-dialog-header">
        <span class="am-dialog-header__badge">
            <component :is="icon" />
        </span>
        <span class="am-dialog-header__title">{{ title }}</span>
    </div>
</template>

<style scoped lang="scss">
.am-dialog-header {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  min-width: 0;
}

.am-dialog-header__badge {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  color: var(--am-accent-primary);
  background: var(--am-accent-primary-soft);
  border-radius: var(--am-radius-sm);

  svg {
    width: 15px;
    height: 15px;
  }
}

.am-dialog-header__title {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-md);
  font-weight: 600;
}
</style>
