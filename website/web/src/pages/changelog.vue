<script setup lang="ts">
import { usePageSeo } from '~/composables/usePageSeo'

const releases = [
    {
        version: 'v3.0.4',
        badge: '最新',
        date: '2026 · 07',
        changes: [
            '告警按 Agent 作用域评估，支持多节点独立阈值',
            'System Audit 页面支持按 Agent 过滤',
            'MailSender 接口与告警通知能力',
            'OpenAPI 类型自动生成（openapi-typescript）',
            'Health Probe 注入 DB / Tunnel 健康检查',
            '凭据强校验：空值、弱默认、长度不足拒绝启动',
            '并发安全与残留二进制清理',
        ],
    },
    {
        version: 'v3.0.0',
        badge: '里程碑',
        date: '2026 · 05',
        changes: [
            '重命名 amprobe → Beacon，品牌与命名统一',
            '统一 Agent 选择（X-Agent-ID / agent_id）',
            'Agent 版本上报、远程更新与自更新卸载',
            '反向 gRPC tunnel 控制通道重构',
            'CORS 白名单化、JWT 生产模式强密钥',
            '分层限流与 WS / report / tunnel 鉴权加固',
        ],
    },
    {
        version: 'v2.0.0',
        badge: '稳定版',
        date: '2025',
        changes: [
            '主机 CPU、内存、磁盘、网络实时监控',
            'Docker 容器全生命周期管理',
            '用户、角色与 API 接口权限管理',
            '操作审计与登录登出记录',
            'Vue 3 + TypeScript + Element Plus 前端架构',
        ],
    },
]

usePageSeo({
    title: '更新日志 - Beacon',
    description: '查看 Beacon 主机监控与 Docker 管理平台的版本演进和重要改进。',
    path: '/changelog',
})
</script>

<template>
    <div>
        <header class="site-page-header">
            <div class="site-container">
                <p class="site-overline">
                    更新日志
                </p>
                <h1 class="site-page-title">
                    Changelog
                </h1>
                <p class="site-page-description">
                    持续迭代，记录每一次进步
                </p>
            </div>
        </header>
        <div class="site-container timeline">
            <article v-for="release in releases" :key="release.version" class="timeline__item">
                <div class="timeline__marker" />
                <div class="site-card timeline__card">
                    <header>
                        <div><strong>{{ release.version }}</strong><span>{{ release.badge }}</span></div>
                        <time>{{ release.date }}</time>
                    </header>
                    <ul>
                        <li v-for="change in release.changes" :key="change">
                            <Icon name="mdi:check-circle" /><span>{{ change }}</span>
                        </li>
                    </ul>
                </div>
            </article>
        </div>
    </div>
</template>

<style scoped lang="scss">
.timeline {
  position: relative;
  max-width: 880px;
  padding-top: 64px;
  padding-bottom: 80px;
}

.timeline::before {
  content: '';
  position: absolute;
  top: 64px;
  bottom: 80px;
  left: 7px;
  width: 1px;
  background: var(--border);
}

.timeline__item {
  position: relative;
  padding: 0 0 var(--space-8) 40px;
}

.timeline__marker {
  position: absolute;
  top: 24px;
  left: 0;
  width: 15px;
  height: 15px;
  background: var(--primary);
  border: 4px solid var(--color-primary-soft);
  border-radius: 50%;
}

.timeline__card {
  padding: var(--space-6);
}

.timeline__card header,
.timeline__card header > div,
.timeline__card li {
  display: flex;
  align-items: center;
}

.timeline__card header {
  justify-content: space-between;
  gap: var(--space-4);
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--border);
}

.timeline__card header > div {
  gap: var(--space-2);
}

.timeline__card header strong {
  font-family: var(--font-mono);
  font-size: 20px;
}

.timeline__card header span {
  padding: 2px 8px;
  color: var(--primary);
  background: var(--color-primary-soft);
  border-radius: 999px;
  font-size: 11px;
}

.timeline__card time {
  color: var(--muted-foreground);
  font-family: var(--font-mono);
  font-size: 12px;
}

.timeline__card ul {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin: var(--space-4) 0 0;
  padding: 0;
  list-style: none;
}

.timeline__card li {
  align-items: flex-start;
  gap: var(--space-2);
  color: var(--color-text-secondary);
}

.timeline__card li :deep(svg) {
  flex: 0 0 auto;
  margin-top: 4px;
  color: var(--color-success);
}

@media (max-width: 640px) {
  .timeline {
    padding-top: 32px;
    padding-bottom: 48px;
  }

  .timeline::before {
    top: 32px;
    bottom: 48px;
  }

  .timeline__item {
    padding-left: 28px;
  }

  .timeline__card {
    padding: var(--space-4);
  }

  .timeline__card header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
