<script setup lang="ts">
import { dayjs } from 'element-plus'
import { useI18n } from 'vue-i18n'

interface DensityOption {
  value: number
  label: string
}

// §7.3 Agent 上下文条：主机分组与容器分组必须使用同一组件，保证呈现完全一致。
// 过期态不得继续显示「实时」（domain I002）：点与文案切到 accent-warning「数据过期」。
const props = withDefaults(defineProps<{
  density: number
  densityOptions: DensityOption[]
  lastUpdated?: number
  stale?: boolean
}>(), {
  lastUpdated: 0,
  stale: false,
})

const emit = defineEmits<{
  'update:density': [number]
}>()

const { t } = useI18n()
const updatedAt = computed(() => (props.lastUpdated ? dayjs(props.lastUpdated).format('HH:mm:ss') : '-'))
</script>

<template>
    <div class="am-monitor-context">
        <span class="am-monitor-context__live" :class="{ 'am-monitor-context__live--stale': stale }">
            <span class="am-monitor-context__dot" />
            {{ stale ? t('agent.stale') : t('monitor.realtime') }}
        </span>
        <span class="am-monitor-context__sep">·</span>
        <span class="am-monitor-context__time">{{ t('monitor.lastUpdated', { time: updatedAt }) }}</span>
        <span class="am-monitor-context__sep">·</span>
        <span class="am-monitor-context__density-label">{{ t('monitor.timeDensity') }}：</span>
        <el-select
            :model-value="density"
            size="small"
            class="am-monitor-context__density"
            @update:model-value="emit('update:density', $event)"
        >
            <el-option v-for="item in densityOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
    </div>
</template>

<style scoped lang="scss">
.am-monitor-context {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  flex-wrap: wrap;
}

.am-monitor-context__live {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--am-accent-success);
  font-size: var(--am-font-xs);
  font-weight: 500;
}

.am-monitor-context__dot {
  width: 6px;
  height: 6px;
  background: var(--am-accent-success);
  border-radius: 50%;
}

// 过期态：状态点与文案改用 warning 语义色，避免把停更数据表达成实时。
.am-monitor-context__live--stale {
  color: var(--am-accent-warning);

  .am-monitor-context__dot {
    background: var(--am-accent-warning);
  }
}

.am-monitor-context__sep {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}

.am-monitor-context__time {
  color: var(--am-foreground-muted);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
}

.am-monitor-context__density-label {
  color: var(--am-foreground-secondary);
  font-size: var(--am-font-sm);
  white-space: nowrap;
}

.am-monitor-context__density {
  width: 120px;

  :deep(.el-select__wrapper) {
    height: 32px;
    padding: 0 10px;
    background: var(--am-surface-inset);
    border-radius: var(--am-radius-sm);
    box-shadow: 0 0 0 1px var(--am-border-primary) inset;
  }
}
</style>
