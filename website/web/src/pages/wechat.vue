<script setup lang="ts">
useHead({ title: '微信公众号 - Beacon' })
</script>

<template>
    <main class="wechat-page site-container">
        <section class="site-card wechat-card">
            <img src="/images/wechat.jpg" alt="微信公众号二维码">
            <div>
                <p class="site-overline">关注我们</p>
                <h1>扫码关注公众号</h1>
                <p>获取 Beacon 最新动态与版本更新</p>
            </div>
        </section>
    </main>
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
  gap: var(--site-space-xl);
  padding: var(--site-space-xl);
}

.wechat-card img {
  width: 220px;
  height: 220px;
  object-fit: cover;
  border-radius: var(--site-radius-md);
}

.wechat-card h1 {
  margin: 0;
  font-size: 28px;
}

.wechat-card div > p:last-child {
  margin-top: var(--site-space-sm);
  color: var(--site-foreground-secondary);
}

@media (max-width: 640px) {
  .wechat-card {
    align-items: flex-start;
    flex-direction: column;
    padding: var(--site-space-lg);
  }

  .wechat-card img {
    width: 100%;
    height: auto;
  }
}
</style>
