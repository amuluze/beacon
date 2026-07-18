<script setup lang="ts">
interface Props {
    type: 'containers' | 'host' | 'users'
}

const props = defineProps<Props>()

const containers = [
    { name: 'nginx-proxy', cpu: 12, memory: 34 },
    { name: 'redis-cache', cpu: 8, memory: 21 },
    { name: 'postgres-db', cpu: 45, memory: 62 },
    { name: 'worker-cron', cpu: 3, memory: 10 },
    { name: 'metrics-scan', cpu: 0, memory: 0 },
]
const users = [
    { initial: 'A', name: '阿慕', role: 'admin', time: '2 分钟前' },
    { initial: 'Z', name: '张伟', role: '运维', time: '15 分钟前' },
    { initial: 'L', name: '李娜', role: '访客', time: '1 小时前' },
    { initial: 'W', name: '王强', role: '运维', time: '刚刚' },
]
const hostMetrics = [
    { label: 'CPU', value: '34%', sub: '8 核 · 3.2 GHz', width: 34 },
    { label: '内存', value: '62%', sub: '19.8 / 32 GB', width: 62 },
    { label: '磁盘', value: '71%', sub: '355 / 500 GB', width: 71 },
    { label: '网络', value: '↓ 1.2 MB/s', sub: '↑ 240 KB/s', width: 45 },
]
</script>

<template>
    <div class="preview site-card">
        <template v-if="props.type === 'containers'">
            <header class="preview__header">
                <strong>Containers</strong><span class="preview__badge">6</span>
            </header>
            <div class="preview__list">
                <div v-for="item in containers" :key="item.name" class="preview__row">
                    <span class="preview__status" />
                    <strong>{{ item.name }}</strong>
                    <span class="preview__metric">CPU {{ item.cpu }}%</span>
                    <span class="preview__metric">MEM {{ item.memory }}%</span>
                </div>
            </div>
        </template>
        <template v-else-if="props.type === 'host'">
            <header class="preview__header">
                <strong>Host · prod-node-01</strong><span class="preview__tag">uptime 42d</span>
            </header>
            <div class="preview__metrics">
                <article v-for="metric in hostMetrics" :key="metric.label" class="preview__tile">
                    <span>{{ metric.label }}</span>
                    <strong>{{ metric.value }}</strong>
                    <div class="preview__track">
                        <span :style="{ width: `${metric.width}%` }" />
                    </div>
                    <small>{{ metric.sub }}</small>
                </article>
            </div>
        </template>
        <template v-else>
            <header class="preview__header">
                <strong>用户管理</strong><span class="preview__badge">12</span>
            </header>
            <div class="preview__list">
                <div v-for="user in users" :key="user.name" class="preview__row preview__row--user">
                    <span class="preview__avatar">{{ user.initial }}</span>
                    <strong>{{ user.name }}</strong>
                    <span class="preview__tag">{{ user.role }}</span>
                    <small>{{ user.time }}</small>
                </div>
            </div>
        </template>
    </div>
</template>

<style scoped lang="scss">
.preview {
  min-height: 360px;
  padding: var(--space-6);
}

.preview__header,
.preview__row {
  display: flex;
  align-items: center;
}

.preview__header {
  justify-content: space-between;
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--border);
}

.preview__badge,
.preview__tag {
  padding: 2px 8px;
  color: var(--primary);
  background: var(--color-primary-soft);
  border-radius: var(--radius-pill);
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
}

.preview__list {
  display: flex;
  flex-direction: column;
}

.preview__row {
  min-height: 56px;
  gap: var(--space-2);
  border-bottom: 1px solid var(--border);
}

.preview__row strong {
  flex: 1;
  font-size: var(--font-size-sm);
}

.preview__status {
  width: 8px;
  height: 8px;
  background: var(--color-success);
  border-radius: 50%;
}

.preview__metric,
.preview__row small {
  color: var(--muted-foreground);
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
}

.preview__avatar {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  color: var(--color-text-inverse);
  background: var(--primary);
  border-radius: 50%;
  font-size: var(--font-size-xs);
}

.preview__metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--space-4);
  padding-top: var(--space-6);
}

.preview__tile {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  padding: var(--space-4);
  background: var(--background);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}

.preview__tile span,
.preview__tile small {
  color: var(--muted-foreground);
  font-size: var(--font-size-xs);
}

.preview__tile strong {
  font-family: var(--font-mono);
  font-size: var(--font-metric-md);
  font-weight: var(--font-weight-bold);
}

.preview__track {
  height: 6px;
  overflow: hidden;
  background: var(--color-surface-muted);
  border-radius: var(--radius-pill);
}

.preview__track span {
  display: block;
  height: 100%;
  background: var(--primary);
}

@media (max-width: 640px) {
  .preview {
    min-height: auto;
    padding: var(--space-4);
  }

  .preview__row {
    min-height: 52px;
  }

  .preview__metric:last-child,
  .preview__row small {
    display: none;
  }

  .preview__metrics {
    grid-template-columns: 1fr;
  }
}
</style>
