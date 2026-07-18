<script setup lang="ts">
import { usePageSeo } from '~/composables/usePageSeo'
import { showToast } from '~/composables/useToast'

const officialBaseURL = 'https://help.beacon.amuluze.com'
const installCommand = `bash -c "$(curl -fsSLk ${officialBaseURL}/release/latest/manager.sh)"`
const copying = shallowRef(false)

const heroTags = [
    { icon: 'lucide:package', label: 'Docker 管理' },
    { icon: 'lucide:activity', label: '实时监控' },
    { icon: 'lucide:feather', label: '轻量部署' },
]

async function copyInstallCommand() {
    if (copying.value)
        return
    copying.value = true
    try {
        await navigator.clipboard.writeText(installCommand)
        showToast('安装命令已复制', 'success')
    }
    catch {
        showToast('复制失败，请手动选择命令', 'error')
    }
    finally {
        copying.value = false
    }
}

usePageSeo({
    title: 'Beacon - 开源轻量级主机与容器监控工具',
    description: 'Beacon 是轻量级 Server-Agent 主机监控及 Docker 容器管理工具。',
    path: '/',
})
</script>

<template>
    <div class="landing">
        <section class="hero">
            <div class="hero__glow" aria-hidden="true" />
            <div class="site-container hero__inner">
                <div class="hero__badge">
                    <Icon name="lucide:sparkles" />
                    <span>开源 · MIT License · 持续维护</span>
                </div>
                <h1>Beacon</h1>
                <p class="hero__subtitle">
                    开源 · 轻量 · 现代化
                </p>
                <p class="hero__description">
                    轻量级主机及容器监控管理工具，实时掌控服务器与 Docker 资源
                </p>
                <div class="hero__tags">
                    <span v-for="tag in heroTags" :key="tag.label" class="hero__tag">
                        <Icon :name="tag.icon" />
                        {{ tag.label }}
                    </span>
                </div>
                <div class="hero__command site-code">
                    <span class="hero__prompt">$</span>
                    <code>{{ installCommand }}</code>
                    <button type="button" class="hero__copy" :aria-label="copying ? '正在复制' : '复制安装命令'" @click="copyInstallCommand">
                        <Icon :name="copying ? 'lucide:loader-circle' : 'lucide:copy'" />
                    </button>
                </div>
                <div class="hero__actions">
                    <button class="site-button site-button--primary" type="button" @click="copyInstallCommand">
                        <Icon name="lucide:terminal" />
                        <span>一键安装</span>
                    </button>
                </div>
            </div>
        </section>

        <div>
            <HomeFeatureSection
                overline="容器管理"
                overline-icon="lucide:package"
                title="全面的 Docker 容器管理"
                description="覆盖容器全生命周期，运行状态与资源占用一目了然"
                :points="['查看 Docker 版本与运行状态', '管理容器、镜像和网络', '按 Agent 安全执行远程控制']"
            >
                <HomeProductPreview type="containers" />
            </HomeFeatureSection>

            <HomeFeatureSection
                overline="主机管理"
                overline-icon="lucide:cpu"
                title="实时监控主机系统资源"
                description="采集主机核心指标，远程运维与管理触手可及"
                :points="['CPU、内存、磁盘与网络趋势', '系统时间与时区管理', '远程终端、重启与关机控制']"
                variant="muted"
                reverse
            >
                <HomeProductPreview type="host" />
            </HomeFeatureSection>

            <HomeFeatureSection
                overline="用户管理"
                overline-icon="lucide:shield-check"
                title="完善的权限与角色管理"
                description="精细化权限控制，操作全程可审计、可追溯"
                :points="['用户与角色权限管理', 'API 接口授权视图', '按 Agent 过滤的系统审计']"
            >
                <HomeProductPreview type="users" />
            </HomeFeatureSection>

            <HomeTechStack />
        </div>

        <section class="cta">
            <div class="cta__glow" aria-hidden="true" />
            <div class="site-container cta__inner">
                <div class="hero__badge">
                    <Icon name="lucide:zap" />
                    <span>快速开始</span>
                </div>
                <h2>立即开始监控你的主机与容器</h2>
                <p>一行命令安装 Beacon Server，几分钟完成接入，开箱即用</p>
                <div class="hero__actions">
                    <NuxtLink class="site-button site-button--primary" to="/document">
                        <Icon name="lucide:rocket" />
                        <span>开始安装</span>
                    </NuxtLink>
                    <a class="site-button site-button--card" href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">
                        <Icon name="lucide:github" />
                        <span>Star on GitHub</span>
                    </a>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped lang="scss">
.hero {
  position: relative;
  overflow: hidden;
  padding: 112px 0 80px;
  text-align: center;
  background: var(--background);
  border-bottom: 1px solid var(--border);
}

.hero__glow,
.cta__glow {
  position: absolute;
  top: 0;
  left: 50%;
  width: 680px;
  height: 560px;
  background: var(--primary);
  filter: blur(130px);
  opacity: 0.13;
  transform: translateX(-50%);
  pointer-events: none;
}

.cta__glow {
  width: 500px;
  height: 340px;
  opacity: 0.11;
}

.hero__inner {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.hero__badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  color: var(--color-text-secondary);
  background: var(--color-surface-muted);
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.hero__badge :deep(svg) {
  color: var(--color-warning);
  font-size: 14px;
}

.hero h1 {
  margin: var(--space-6) 0 0;
  font-size: var(--font-display-xl);
  font-weight: var(--font-weight-extrabold);
  line-height: var(--line-height-none);
  letter-spacing: var(--letter-spacing-tighter);
}

.hero__subtitle {
  margin: var(--space-4) 0 0;
  color: var(--color-text-secondary);
  font-size: var(--font-display-2xs);
  font-weight: var(--font-weight-medium);
  letter-spacing: var(--letter-spacing-wider);
}

.hero__description {
  max-width: 640px;
  margin: var(--space-4) 0 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-lg);
  line-height: var(--line-height-relaxed);
}

.hero__tags {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: var(--space-3);
  margin-top: var(--space-5);
}

.hero__tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  color: var(--color-text-secondary);
  background: var(--color-surface-muted);
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.hero__tag :deep(svg) {
  color: var(--primary);
  font-size: 16px;
}

.hero__command {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: var(--space-3);
  align-items: center;
  width: min(680px, 100%);
  margin-top: var(--space-8);
  padding: 14px 16px;
  text-align: left;
}

.hero__prompt {
  color: var(--color-success);
  font-weight: var(--font-weight-bold);
}

.hero__command code {
  overflow: hidden;
  text-overflow: ellipsis;
}

.hero__copy {
  display: inline-grid;
  place-items: center;
  width: 28px;
  height: 28px;
  padding: 0;
  color: var(--color-text-secondary);
  background: var(--border);
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: color 0.2s ease;
}

.hero__copy:hover {
  color: var(--primary);
}

.hero__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: var(--space-3);
  margin-top: var(--space-6);
}

.cta {
  position: relative;
  overflow: hidden;
  padding: 88px 0;
  text-align: center;
  background: var(--color-surface-muted);
  border-top: 1px solid var(--border);
}

.cta__inner {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.cta h2 {
  margin: var(--space-5) 0 0;
  font-size: var(--font-display-md);
  font-weight: var(--font-weight-bold);
  line-height: var(--line-height-snug);
}

.cta p {
  max-width: 520px;
  margin: var(--space-3) 0 0;
  color: var(--muted-foreground);
}

@media (max-width: 640px) {
  .hero {
    padding: 72px 0 48px;
    text-align: left;
  }

  .hero__inner {
    align-items: stretch;
  }

  .hero__badge {
    align-self: flex-start;
  }

  .hero__tags {
    justify-content: flex-start;
  }

  .hero__command {
    font-size: var(--font-size-xs);
  }

  .hero__actions {
    flex-direction: column;
    align-items: stretch;
  }

  .cta {
    padding: 56px 0;
    text-align: left;
  }

  .cta__inner {
    align-items: stretch;
  }

  .cta .hero__actions {
    flex-direction: column;
  }
}
</style>
