<script setup lang="ts">
import IconTriangleAlert from '~icons/lucide/triangle-alert'
import { formatMonitorAge, formatSampleClock, formatSampleTime } from '@/utils/monitorStale'

// 稿面 `Monitor · 数据过期`（DESIGN.md §7.5）：警示条 + KPI 大数字 + 采样明细。
// 过期是独立态而非错误态：保留最后采样值，只标注过期，并明示恢复后自动续传。
export interface StaleKpi {
  key: string
  label: string
  value: string
  note: string
  lastSampleAt: number
}

export interface StaleDetailRow {
  key: string
  label: string
  value: string
  lastSampleAt: number
}

const props = defineProps<{
  lastSampleAt: number
  kpis: StaleKpi[]
  rows: StaleDetailRow[]
}>()

const interrupted = computed(() => formatMonitorAge(Math.floor(Date.now() / 1000) - props.lastSampleAt))
const lastSampleText = computed(() => formatSampleTime(props.lastSampleAt))

function rowAge(lastSampleAt: number): string {
  return formatMonitorAge(Math.floor(Date.now() / 1000) - lastSampleAt)
}
</script>

<template>
    <div class="am-stale">
        <div class="am-stale__banner" role="status">
            <IconTriangleAlert class="am-stale__banner-icon" aria-hidden="true" />
            <span class="am-stale__banner-title">Agent 上报已中断 {{ interrupted }}</span>
            <span class="am-stale__banner-sep" aria-hidden="true">·</span>
            <span class="am-stale__banner-time">最后一次采样 {{ lastSampleText }}</span>
            <span class="am-stale__banner-hint">数据恢复后自动续传，无需手动刷新</span>
        </div>

        <div class="am-stale__kpis">
            <section v-for="item in kpis" :key="item.key" class="am-stale__kpi">
                <div class="am-stale__kpi-head">
                    <span class="am-stale__kpi-label">{{ item.label }}</span>
                    <span class="am-stale__kpi-badge">数据过期</span>
                    <span class="am-stale__kpi-spacer" />
                    <span class="am-stale__kpi-sampled">最后采样 {{ formatSampleClock(item.lastSampleAt) }}</span>
                </div>
                <p class="am-stale__kpi-value">
                    {{ item.value }}
                </p>
                <p class="am-stale__kpi-note">
                    {{ item.note }}
                </p>
            </section>
        </div>

        <section class="am-stale__detail">
            <header class="am-stale__detail-head">
                <span class="am-stale__detail-title">采样明细</span>
                <span class="am-stale__detail-summary">共 {{ rows.length }} 项指标 · 全部过期</span>
            </header>
            <div class="am-stale__table">
                <div class="am-stale__row am-stale__row--header">
                    <span>指标</span>
                    <span class="am-stale__cell-value">最后采样值</span>
                    <span class="am-stale__cell-age">过期时长</span>
                    <span class="am-stale__cell-status">状态</span>
                </div>
                <div v-for="row in rows" :key="row.key" class="am-stale__row">
                    <span class="am-stale__cell-label">{{ row.label }}</span>
                    <span class="am-stale__cell-value">{{ row.value }}</span>
                    <span class="am-stale__cell-age">{{ rowAge(row.lastSampleAt) }}</span>
                    <span class="am-stale__cell-status">等待 Agent 续传</span>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped lang="scss">
.am-stale {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
}

.am-stale__banner {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  padding: 12px 16px;
  color: var(--am-accent-warning);
  font-size: var(--am-font-sm);
  background: var(--am-accent-warning-soft);
  border-radius: var(--am-radius-sm);
}

.am-stale__banner-icon {
  flex: 0 0 auto;
  width: 16px;
  height: 16px;
}

.am-stale__banner-title {
  font-weight: 600;
}

.am-stale__banner-sep {
  color: var(--am-foreground-secondary);
}

.am-stale__banner-time {
  color: var(--am-foreground-secondary);
  font-size: var(--am-font-xs);
}

.am-stale__banner-hint {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}

.am-stale__kpis {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--am-spacing-md);
}

.am-stale__kpi {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
  padding: var(--am-spacing-lg);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
}

.am-stale__kpi-head {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

.am-stale__kpi-label {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-sm);
  font-weight: 600;
}

.am-stale__kpi-badge {
  padding: 2px 8px;
  color: var(--am-accent-warning);
  font-size: var(--am-font-xs);
  font-weight: 600;
  background: var(--am-accent-warning-soft);
  border-radius: var(--am-radius-xs);
}

.am-stale__kpi-spacer {
  flex: 1;
}

.am-stale__kpi-sampled {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}

.am-stale__kpi-value {
  margin: 0;
  color: var(--am-foreground-primary);
  font-size: var(--am-font-2xl);
  font-weight: 700;
  line-height: 1.2;
}

.am-stale__kpi-note {
  margin: 0;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  line-height: 1.5;
}

.am-stale__detail {
  display: flex;
  flex-direction: column;
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
}

.am-stale__detail-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
}

.am-stale__detail-title {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-md);
  font-weight: 600;
}

.am-stale__detail-summary {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}

.am-stale__table {
  display: flex;
  flex-direction: column;
}

.am-stale__row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 150px 120px 180px;
  gap: var(--am-spacing-md);
  align-items: center;
  padding: 10px 16px;
  background: var(--am-surface-primary);
  border-bottom: 1px solid var(--am-border-subtle);
  font-size: var(--am-font-sm);

  &:last-child {
    border-bottom: 0;
  }
}

.am-stale__row--header {
  color: var(--am-foreground-secondary);
  font-size: var(--am-font-xs);
  font-weight: 500;
  background: var(--am-surface-secondary);
}

.am-stale__cell-label {
  color: var(--am-foreground-primary);
}

.am-stale__cell-value {
  color: var(--am-foreground-secondary);
  font-size: var(--am-font-xs);
}

.am-stale__cell-age {
  color: var(--am-accent-warning);
  font-size: var(--am-font-xs);
  font-weight: 600;
}

.am-stale__cell-status {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}
</style>
