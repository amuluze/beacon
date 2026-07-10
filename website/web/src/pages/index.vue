<script setup lang="ts">
import { statisticQuery, statisticUpdate } from '~/api/statistics'

const officialBaseURL = 'https://official.beacon.amuluze.com'
const installCommand = `curl -fsSL ${officialBaseURL}/download/install.sh | sh`
const statisticID = shallowRef<number | null>(null)
const statistic = shallowRef<number | null>(null)
const copying = shallowRef(false)

const statisticLabel = computed(() => statistic.value === null ? '--' : statistic.value.toLocaleString())

async function loadStatistic() {
    try {
        const reply = await statisticQuery()
        statisticID.value = reply.data.id
        statistic.value = reply.data.times
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
        ElMessage.success('安装命令已复制')
    }
    catch {
        ElMessage.error('复制失败，请手动选择命令')
    }
    finally {
        copying.value = false
    }
}

async function downloadInstallScript() {
    window.open('/download/install.sh', '_blank')
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

useHead({
    title: 'Beacon - 开源轻量级主机与容器监控工具',
    meta: [{
        name: 'description',
        content: 'Beacon 是轻量级 Server-Agent 主机监控及 Docker 容器管理工具。',
    }],
})
</script>

<template>
    <div class="landing">
        <section class="hero">
            <div class="site-container hero__inner">
                <div class="hero__badge">开源 · MIT License · 持续维护</div>
                <h1>Beacon</h1>
                <p class="hero__subtitle">开源 · 轻量 · 现代化</p>
                <p class="hero__description">轻量级主机及容器监控管理工具，实时掌控服务器与 Docker 资源</p>
                <div class="hero__command site-code">
                    <span>$</span>
                    <code>{{ installCommand }}</code>
                    <button type="button" :aria-label="copying ? '正在复制' : '复制安装命令'" @click="copyInstallCommand">
                        <Icon :name="copying ? 'mdi:loading' : 'mdi:content-copy'" />
                    </button>
                </div>
                <div class="hero__actions">
                    <button class="button button--primary" type="button" @click="downloadInstallScript">立即体验</button>
                    <a class="button" href="https://github.com/amuluze/amprobe" target="_blank" rel="noopener noreferrer">GitHub</a>
                </div>
                <div class="hero__stats">
                    <span>累计获取</span>
                    <strong>{{ statisticLabel }}</strong>
                    <span>次</span>
                </div>
            </div>
        </section>

        <div class="site-container">
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
            <div class="site-container cta__inner">
                <p class="site-overline">快速开始</p>
                <h2>立即开始监控你的主机与容器</h2>
                <p>一行命令安装 Agent，几分钟完成接入，开箱即用</p>
                <div class="hero__actions">
                    <NuxtLink class="button button--primary" to="/document">查看使用手册</NuxtLink>
                    <a class="button" href="https://github.com/amuluze/amprobe" target="_blank" rel="noopener noreferrer">访问 GitHub</a>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped lang="scss">
.hero {
  padding: 112px 0 80px;
  text-align: center;
  background: var(--site-surface-secondary);
  border-bottom: 1px solid var(--site-border-subtle);
}

.hero__inner {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.hero__badge {
  padding: 5px 10px;
  color: var(--site-accent);
  background: var(--site-accent-soft);
  border: 1px solid rgba(64, 158, 255, 0.28);
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}

.hero h1 {
  margin: 24px 0 0;
  font-size: clamp(56px, 10vw, 88px);
  line-height: 1;
  letter-spacing: -0.05em;
}

.hero__subtitle {
  margin: var(--site-space-md) 0 0;
  color: var(--site-foreground-secondary);
  font-size: clamp(18px, 3vw, 24px);
  letter-spacing: 0.08em;
}

.hero__description {
  max-width: 640px;
  margin: var(--site-space-md) 0 0;
  color: var(--site-foreground-muted);
  font-size: 16px;
}

.hero__command {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 12px;
  align-items: center;
  width: min(680px, 100%);
  margin-top: var(--site-space-xl);
  padding: 13px 16px;
  text-align: left;
}

.hero__command > span {
  color: var(--site-success);
}

.hero__command code {
  overflow: hidden;
  text-overflow: ellipsis;
}

.hero__command button {
  padding: 2px;
  color: var(--site-foreground-muted);
  background: transparent;
  border: 0;
  cursor: pointer;
}

.hero__actions {
  display: flex;
  justify-content: center;
  gap: var(--site-space-sm);
  margin-top: var(--site-space-lg);
}

.button {
  min-width: 112px;
  padding: 9px 16px;
  color: var(--site-foreground-secondary);
  background: var(--site-surface-card);
  border: 1px solid var(--site-border-primary);
  border-radius: var(--site-radius-sm);
  font-weight: 600;
  cursor: pointer;
}

.button--primary {
  color: var(--site-on-accent);
  background: var(--site-accent);
  border-color: var(--site-accent);
}

.hero__stats {
  display: flex;
  align-items: baseline;
  gap: var(--site-space-sm);
  margin-top: var(--site-space-xl);
  color: var(--site-foreground-muted);
  font-size: 12px;
}

.hero__stats strong {
  color: var(--site-foreground-primary);
  font-family: var(--site-font-mono);
  font-size: 24px;
}

.cta {
  padding: 80px 0;
  text-align: center;
  background: var(--site-surface-secondary);
  border-top: 1px solid var(--site-border-subtle);
}

.cta h2 {
  margin: 0;
  font-size: clamp(28px, 5vw, 40px);
}

.cta p:not(.site-overline) {
  color: var(--site-foreground-secondary);
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

  .hero__command {
    font-size: 11px;
  }

  .hero__actions {
    justify-content: flex-start;
  }

  .hero__stats {
    justify-content: flex-start;
  }

  .cta {
    padding: 56px 0;
    text-align: left;
  }

  .cta .hero__actions {
    flex-direction: column;
  }

  .cta .button {
    text-align: center;
  }
}
</style>
