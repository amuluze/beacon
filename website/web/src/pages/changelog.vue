<script setup lang="ts">
import { usePageSeo } from '~/composables/usePageSeo'

type ChangeType = '新功能' | '改进' | '安全' | '修复'

interface ReleaseChange {
    type: ChangeType
    text: string
}

interface Release {
    version: string
    badge: string
    date: string
    datetime: string
    changes: ReleaseChange[]
}

const releases: Release[] = [
    {
        version: 'v3.0.4',
        badge: '最新',
        date: '2026 · 07',
        datetime: '2026-07',
        changes: [
            { type: '新功能', text: '告警按 Agent 作用域评估，支持多节点独立阈值' },
            { type: '新功能', text: 'System Audit 页面支持按 Agent 过滤' },
            { type: '新功能', text: 'MailSender 接口与告警通知能力' },
            { type: '改进', text: 'OpenAPI 类型自动生成（openapi-typescript）' },
            { type: '改进', text: 'Health Probe 注入 DB / Tunnel 健康检查' },
            { type: '安全', text: '凭据强校验：空值、弱默认、长度不足拒绝启动' },
            { type: '修复', text: '并发安全与残留二进制清理' },
        ],
    },
    {
        version: 'v3.0.0',
        badge: '里程碑',
        date: '2026 · 05',
        datetime: '2026-05',
        changes: [
            { type: '新功能', text: '重命名 amprobe → Beacon，品牌与命名统一' },
            { type: '新功能', text: '统一 Agent 选择（X-Agent-ID / agent_id）' },
            { type: '新功能', text: 'Agent 版本上报、远程更新与自更新卸载' },
            { type: '改进', text: '反向 gRPC tunnel 控制通道重构' },
            { type: '安全', text: 'CORS 白名单化、JWT 生产模式强密钥' },
            { type: '安全', text: '分层限流与 WS / report / tunnel 鉴权加固' },
        ],
    },
    {
        version: 'v2.0.0',
        badge: '稳定版',
        date: '2025',
        datetime: '2025',
        changes: [
            { type: '新功能', text: '主机 CPU、内存、磁盘、网络实时监控' },
            { type: '新功能', text: 'Docker 容器全生命周期管理' },
            { type: '新功能', text: '用户、角色与 API 接口权限管理' },
            { type: '新功能', text: '操作审计与登录登出记录' },
            { type: '改进', text: 'Vue 3 + TypeScript + Element Plus 前端架构' },
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
    <div class="changelog">
        <header class="changelog__header">
            <div class="site-container changelog__header-inner">
                <h1>Changelog</h1>
                <p>持续迭代，记录每一次进步</p>
            </div>
        </header>
        <div class="site-container changelog__body">
            <article v-for="release in releases" :key="release.version" class="changelog__card">
                <header class="changelog__version">
                    <div class="changelog__version-info">
                        <strong>{{ release.version }}</strong>
                        <span class="changelog__badge">{{ release.badge }}</span>
                    </div>
                    <time :datetime="release.datetime">{{ release.date }}</time>
                </header>
                <div class="changelog__divider" />
                <ul class="changelog__entries">
                    <li v-for="change in release.changes" :key="change.text" class="release-change">
                        <span class="release-change__type" :data-type="change.type">{{ change.type }}</span>
                        <span class="release-change__text">{{ change.text }}</span>
                    </li>
                </ul>
            </article>
        </div>
    </div>
</template>

<style scoped lang="scss">
.changelog {
  background: var(--background);
}

.changelog__header {
  padding: 88px 0 40px;
  background: var(--background);
}

.changelog__header-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  text-align: center;
}

.changelog__header h1 {
  margin: 0;
  color: var(--foreground);
  font-size: 44px;
  font-weight: 800;
  letter-spacing: -1px;
  line-height: 1.1;
}

.changelog__header p {
  max-width: 480px;
  margin: 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-md);
  line-height: 1.6;
}

.changelog__body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-10);
  padding: 0 0 96px;
}

.changelog__card {
  width: 880px;
  max-width: 100%;
  padding: var(--space-8);
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
}

.changelog__version {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
}

.changelog__version-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.changelog__version-info strong {
  color: var(--foreground);
  font-family: var(--font-mono);
  font-size: 22px;
  font-weight: 700;
}

.changelog__badge {
  padding: 4px 9px;
  color: var(--foreground);
  background: var(--color-primary-soft);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: 600;
}

.changelog__version time {
  color: var(--muted-foreground);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
}

.changelog__divider {
  height: 1px;
  margin: var(--space-5) 0;
  background: var(--border);
}

.changelog__entries {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.release-change {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.release-change__type {
  flex: 0 0 auto;
  padding: 4px 9px;
  color: var(--foreground);
  background: var(--color-primary-soft);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-align: center;
}

.release-change__type[data-type='改进'] {
  color: var(--color-success);
  background: var(--color-improvement-soft);
}

.release-change__type[data-type='安全'] {
  color: var(--color-warning);
  background: var(--color-warning-soft);
}

.release-change__type[data-type='修复'] {
  color: var(--color-error);
  background: var(--color-error-soft);
}

.release-change__text {
  flex: 1;
  padding-top: 2px;
}

@media (max-width: 720px) {
  .changelog__header {
    padding: 56px 0 var(--space-8);
  }

  .changelog__header h1 {
    font-size: 32px;
  }

  .changelog__card {
    padding: var(--space-6);
  }

  .changelog__version {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-2);
  }
}
</style>
