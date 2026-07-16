<script setup lang="ts">
import { usePageSeo } from '~/composables/usePageSeo'

usePageSeo({
    title: '隐私与安装遥测 - Beacon',
    description: '了解 Beacon 官网在安装上报中收集的数据、用途与安全边界。',
    path: '/privacy',
})
</script>

<template>
    <div>
        <header class="site-page-header">
            <div class="site-container">
                <p class="site-overline">
                    隐私说明
                </p>
                <h1 class="site-page-title">
                    隐私与安装遥测
                </h1>
                <p class="site-page-description">
                    以透明、最小化的方式说明 Beacon 安装过程产生的遥测信息
                </p>
            </div>
        </header>

        <div class="site-container privacy-content">
            <section class="site-card">
                <h2>收集哪些信息</h2>
                <p>Beacon 启动时会以尽力而为的方式上报安装标识、镜像与版本、公开访问地址、安装目录、服务端口、容器名称、主机名、客户端 IP、User-Agent 和上报时间；稳定安装标识用于重复上报去重。</p>
            </section>
            <section class="site-card">
                <h2>为什么收集</h2>
                <p>这些信息仅用于去重统计安装规模、判断版本分布和改善安装兼容性，不用于身份认证、许可控制或收集业务数据。</p>
            </section>
            <section class="site-card">
                <h2>明确不收集</h2>
                <p>安装上报不包含登录密码、认证 Token、签名密钥、容器环境变量、业务日志或被监控主机的指标明细。</p>
            </section>
            <section class="site-card">
                <h2>联系我们</h2>
                <p>如需查询或删除与安装标识相关的数据，请通过 GitHub Issues 联系维护者，并仅提供必要的安装标识。</p>
                <a href="https://github.com/amuluze/beacon/issues" target="_blank" rel="noopener noreferrer">前往 GitHub Issues</a>
            </section>
        </div>
    </div>
</template>

<style scoped lang="scss">
.privacy-content {
  display: grid;
  gap: var(--space-4);
  max-width: 880px;
  padding-top: 64px;
  padding-bottom: 80px;
}

.privacy-content section {
  padding: var(--space-6);
}

.privacy-content h2,
.privacy-content p {
  margin: 0;
}

.privacy-content h2 {
  font-size: var(--font-size-lg);
}

.privacy-content p {
  margin-top: var(--space-2);
  color: var(--color-text-secondary);
}

.privacy-content a {
  display: inline-flex;
  margin-top: var(--space-4);
  color: var(--primary);
  font-weight: 600;
}

@media (max-width: 640px) {
  .privacy-content {
    padding-top: 32px;
    padding-bottom: 48px;
  }
}
</style>
