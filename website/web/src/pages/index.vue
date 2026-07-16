<script setup lang="ts">
import { statisticQuery, statisticUpdate } from '~/api/statistics'
import { usePageSeo } from '~/composables/usePageSeo'
import { showToast } from '~/composables/useToast'

const officialBaseURL = 'https://help.beacon.amuluze.com'
const installCommand = `curl -fsSLO ${officialBaseURL}/release/latest/manager.sh && curl -fsSLO ${officialBaseURL}/release/latest/SHA256SUMS && grep ' manager.sh$' SHA256SUMS | sha256sum -c - && sudo sh manager.sh`
const statisticID = shallowRef<number | null>(null)
const statistic = shallowRef<number | null>(null)
const copying = shallowRef(false)

const heroTags = [
    { icon: 'mdi:package-variant-closed', label: 'Docker 管理' },
    { icon: 'mdi:chart-line-variant', label: '实时监控' },
    { icon: 'mdi:feather', label: '轻量部署' },
]

const statisticLabel = computed(() => typeof statistic.value === 'number' ? statistic.value.toLocaleString() : '--')

async function loadStatistic() {
    try {
        const reply = await statisticQuery()
        statisticID.value = reply.id
        statistic.value = reply.times
    }
    catch {
        statisticID.value = null
        statistic.value = null
    }
}

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

async function downloadInstallScript() {
    window.open('/release/latest/manager.sh', '_blank')
    if (statisticID.value === null)
        return
    try {
        await statisticUpdate({ id: statisticID.value })
        await loadStatistic()
    }
    catch {
        // 统计失败不影响安装脚本下载与页面主体。
    }
}

onMounted(() => {
    void loadStatistic()
})

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
                    <Icon name="mdi:sparkles" />
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
                        <Icon :name="copying ? 'mdi:loading' : 'mdi:content-copy'" />
                    </button>
                </div>
                <div class="hero__actions">
                    <button class="site-button site-button--primary" type="button" @click="copyInstallCommand">
                        <Icon name="mdi:console-line" />
                        <span>复制安装命令</span>
                    </button>
                    <button class="site-button site-button--ghost" type="button" @click="downloadInstallScript">
                        <Icon name="mdi:download-outline" />
                        <span>下载安装脚本</span>
                    </button>
                    <a class="site-button site-button--card" href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">
                        <Icon name="mdi:github" />
                        <span>GitHub 开源仓库</span>
                    </a>
                </div>
                <div class="hero__stats">
                    <Icon name="mdi:cloud-download-outline" />
                    <span>累计获取</span>
                    <strong>{{ statisticLabel }}</strong>
                    <span>次</span>
                </div>
            </div>
        </section>

        <div>
            <HomeFeatureSection
                overline="容器管理"
                title="全面的 Docker 容器管理"
                description="覆盖容器全生命周期，运行状态与资源占用一目了然"
                :points="['查看 Docker 版本与运行状态', '管理容器、镜像和网络', '按 Agent 安全执行远程控制']"
            >
                <HomeProductPreview type="containers" />
            </HomeFeatureSection>

            <HomeFeatureSection
                overline="主机管理"
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
                    <Icon name="mdi:flash-outline" />
                    <span>快速开始</span>
                </div>
                <h2>立即开始监控你的主机与容器</h2>
                <p>一行命令安装 Beacon Server，几分钟完成接入，开箱即用</p>
                <div class="hero__actions">
                    <NuxtLink class="site-button site-button--primary" to="/document">
                        <Icon name="mdi:rocket-launch-outline" />
                        <span>开始安装</span>
                    </NuxtLink>
                    <a class="site-button site-button--card" href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">
                        <Icon name="mdi:star-outline" />
                        <span>Star on GitHub</span>
                    </a>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped lang="scss">
.site-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  min-width: 112px;
  height: 40px;
  padding: 0 20px;
  font-size: var(--font-size-md);
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition:
    background 0.2s ease,
    border-color 0.2s ease;
}

.site-button--primary {
  color: var(--color-text-inverse);
  background: var(--primary);
  border: 1px solid var(--primary);
}

.site-button--primary:hover {
  background: var(--color-brand-hover);
  border-color: var(--color-brand-hover);
}

.site-button--ghost {
  color: var(--color-text-secondary);
  background: var(--color-surface-muted);
  border: 1px solid var(--border);
}

.site-button--ghost:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.site-button--card {
  color: var(--color-text-secondary);
  background: var(--card);
  border: 1px solid var(--border);
}

.site-button--card:hover {
  border-color: var(--primary);
  color: var(--primary);
}

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
  font-weight: 500;
}

.hero__badge :deep(svg) {
  color: var(--color-warning);
  font-size: 14px;
}

.hero h1 {
  margin: var(--space-6) 0 0;
  font-size: clamp(56px, 10vw, 80px);
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.04em;
}

.hero__subtitle {
  margin: var(--space-4) 0 0;
  color: var(--color-text-secondary);
  font-size: clamp(18px, 3vw, 22px);
  font-weight: 500;
  letter-spacing: 0.24em;
}

.hero__description {
  max-width: 640px;
  margin: var(--space-4) 0 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-lg);
  line-height: 1.6;
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
  font-weight: 500;
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
  font-weight: 700;
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

.hero__stats {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-8);
  padding: 10px 18px;
  color: var(--muted-foreground);
  background: var(--color-surface-muted);
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: var(--font-size-sm);
}

.hero__stats :deep(svg) {
  color: var(--primary);
  font-size: 16px;
}

.hero__stats strong {
  color: var(--foreground);
  font-family: var(--font-mono);
  font-size: var(--font-size-lg);
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
  font-size: clamp(28px, 5vw, 36px);
  line-height: 1.2;
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

  .hero h1 {
    font-size: 56px;
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

  .hero__stats {
    align-self: flex-start;
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
