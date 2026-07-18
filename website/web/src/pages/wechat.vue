<script setup lang="ts">
import { usePageSeo } from '~/composables/usePageSeo'

usePageSeo({
    title: '微信公众号 - Beacon',
    description: '关注 Beacon 公众号，获取版本更新和开源项目动态。',
    path: '/wechat',
})
</script>

<template>
    <div class="wechat-page site-container">
        <section class="site-card wechat-card">
            <img src="/images/wechat.jpg" alt="微信公众号二维码">
            <div>
                <p class="site-overline">
                    关注我们
                </p>
                <h1>扫码关注公众号</h1>
                <p>获取 Beacon 最新动态与版本更新</p>
            </div>
        </section>
    </div>
</template>

<style scoped lang="scss">
.wechat-page {
  display: grid;
  place-items: center;
  min-height: 70vh;
  padding-top: 64px;
  padding-bottom: 64px;
}

.wechat-card {
  display: flex;
  align-items: center;
  gap: var(--space-8);
  padding: var(--space-8);
}

.wechat-card img {
  width: 220px;
  height: 220px;
  object-fit: cover;
  border-radius: var(--radius-m);
}

.wechat-card h1 {
  margin: 0;
  font-size: var(--font-display-md);
  font-weight: var(--font-weight-extrabold);
  line-height: var(--line-height-tight);
}

.wechat-card div > p:last-child {
  margin-top: var(--space-2);
  color: var(--color-text-secondary);
}

@media (max-width: 640px) {
  .wechat-card {
    align-items: flex-start;
    flex-direction: column;
    padding: var(--space-6);
  }

  .wechat-card img {
    width: 100%;
    height: auto;
  }
}
</style>
