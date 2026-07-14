<script setup lang="ts">
const contacts = [
    { icon: 'mdi:wechat', label: '微信', value: 'Amuluze' },
    { icon: 'mdi:github', label: 'GitHub', value: 'amuluze/beacon', href: 'https://github.com/amuluze/beacon' },
    { icon: 'mdi:qrcode', label: '公众号', value: '阿慕的订阅', href: '/wechat' },
]

useHead({ title: '团队故事 - Beacon' })
</script>

<template>
    <div class="about-page">
        <header class="site-page-header">
            <div class="site-container">
                <p class="site-overline">关于 Beacon</p>
                <h1 class="site-page-title">团队故事</h1>
                <p class="site-page-description">一名见习独立开发者的开源练手之作</p>
            </div>
        </header>

        <main class="site-container about-page__content">
            <article class="site-card letter">
                <p>你好，我是阿慕 —— 一个 35 岁的临“退”程序员，一名见习独立开发者。</p>
                <p>程序员生涯的进度条已快到头，下一步打算成为一名独立开发者，做自己喜欢的事，也能养家糊口。如果你也面临职业生涯的窘境，想探索新的可能，欢迎加我微信 Amuluze 交流。</p>
                <p>Beacon 是我作为见习独立开发者的第一款开源产品，既满足自己监控云服务器的需求，也用来锻炼独立开发的能力。</p>
                <p>它是一款开源的云服务器监控与 Docker 容器管理工具，帮助实时掌控 CPU、内存、磁盘等资源占用，也能便捷管理部署在云服务上的容器状态。</p>
                <p>Beacon 仍在持续开发与完善中，如果你有需求或建议，欢迎通过 GitHub Issues 或微信联系我，也可以关注公众号了解最新进展。</p>
                <footer class="letter__signature">
                    <strong>阿慕</strong>
                    <span>2024.05.26 夜</span>
                </footer>
            </article>

            <section class="contact">
                <p class="site-overline">与我联系</p>
                <div class="contact__grid">
                    <component
                        :is="item.href ? (item.href.startsWith('/') ? 'NuxtLink' : 'a') : 'div'"
                        v-for="item in contacts"
                        :key="item.label"
                        class="site-card contact__card"
                        :to="item.href?.startsWith('/') ? item.href : undefined"
                        :href="item.href && !item.href.startsWith('/') ? item.href : undefined"
                        :target="item.href && !item.href.startsWith('/') ? '_blank' : undefined"
                        rel="noopener noreferrer"
                    >
                        <span class="contact__icon"><Icon :name="item.icon" /></span>
                        <span>{{ item.label }}</span>
                        <strong>{{ item.value }}</strong>
                    </component>
                </div>
            </section>

            <section class="site-card qrcode-card">
                <img src="/images/wechat.jpg" alt="Beacon 公众号二维码">
                <div>
                    <strong>扫码关注公众号</strong>
                    <p>获取 Beacon 最新动态与版本更新</p>
                </div>
            </section>
        </main>
    </div>
</template>

<style scoped lang="scss">
.about-page__content {
  display: flex;
  flex-direction: column;
  gap: 64px;
  padding-top: 64px;
  padding-bottom: 80px;
}

.letter {
  max-width: 800px;
  margin: 0 auto;
  padding: 48px;
}

.letter p {
  margin: 0 0 var(--site-space-lg);
  color: var(--site-foreground-secondary);
  font-size: 16px;
  line-height: 1.9;
}

.letter__signature {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--site-space-xs);
  color: var(--site-foreground-muted);
}

.letter__signature strong {
  color: var(--site-foreground-primary);
}

.contact {
  text-align: center;
}

.contact__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--site-space-md);
  margin-top: var(--site-space-lg);
}

.contact__card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--site-space-xs);
  padding: var(--site-space-lg);
}

.contact__icon {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  margin-bottom: var(--site-space-sm);
  color: var(--site-accent);
  background: var(--site-accent-soft);
  border-radius: var(--site-radius-md);
  font-size: 22px;
}

.contact__card > span:not(.contact__icon) {
  color: var(--site-foreground-muted);
  font-size: 12px;
}

.contact__card strong {
  font-size: 13px;
}

.qrcode-card {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--site-space-lg);
  max-width: 520px;
  margin: 0 auto;
  padding: var(--site-space-lg);
}

.qrcode-card img {
  width: 84px;
  height: 84px;
  object-fit: cover;
}

.qrcode-card p {
  margin: var(--site-space-xs) 0 0;
  color: var(--site-foreground-muted);
}

@media (max-width: 640px) {
  .about-page__content {
    gap: 40px;
    padding-top: 32px;
    padding-bottom: 48px;
  }

  .letter {
    padding: var(--site-space-lg);
  }

  .letter p {
    font-size: 14px;
  }

  .contact__grid {
    grid-template-columns: 1fr;
  }
}
</style>
