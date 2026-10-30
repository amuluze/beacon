<!--
    BatchActionBar — 表格选中行后的批量操作入口。

    与设计稿一致：工具条右侧展示「已选 N 项」并给出一个批量操作下拉，
    未选中任何行时整体禁用。具体动作由宿主视图提供，组件本身只负责展示与派发。
-->
<script setup lang="ts">
export interface BatchAction {
  // 派发回宿主的动作标识
  key: string
  // i18n key
  label: string
  disabled?: boolean
}

const props = defineProps<{
  count: number
  actions: BatchAction[]
}>()

const emit = defineEmits<{
  select: [string]
}>()

const disabled = computed(() => props.count === 0)
</script>

<template>
    <div class="am-table-operator__batch" :class="{ 'is-empty': disabled }">
        <span class="am-table-operator__selection">
            {{ $t('container.selectedCount', { count: props.count }) }}
        </span>
        <el-dropdown trigger="click" :disabled="disabled" @command="emit('select', $event)">
            <el-button size="small" :disabled="disabled">
                {{ $t('container.batchActions') }}
                <svg-icon icon-class="more" />
            </el-button>
            <template #dropdown>
                <el-dropdown-menu>
                    <el-dropdown-item
                        v-for="action in props.actions"
                        :key="action.key"
                        :command="action.key"
                        :disabled="action.disabled"
                    >
                        {{ $t(action.label) }}
                    </el-dropdown-item>
                </el-dropdown-menu>
            </template>
        </el-dropdown>
    </div>
</template>

<style scoped lang="scss">
.am-table-operator__batch {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

.am-table-operator__selection {
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-sm);
  white-space: nowrap;
}

.am-table-operator__batch.is-empty .am-table-operator__selection {
  color: var(--am-foreground-muted);
}
</style>
