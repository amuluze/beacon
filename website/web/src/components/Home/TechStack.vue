<script setup lang="ts">
const stacks = [
    { title: '前端', items: ['Vue 3', 'TypeScript', 'Element Plus', 'Pinia'] },
    { title: '后端', items: ['Golang', 'Fiber', 'Wire DI'] },
    { title: '存储', items: ['SQLite', 'GORM'] },
    { title: '传输', items: ['gRPC Tunnel', 'WebSocket', 'HTTP/2'] },
]
</script>

<template>
    <section class="tech-stack">
        <header>
            <p class="site-overline">技术栈</p>
            <h2>现代化技术栈驱动</h2>
            <p>Server-Agent 架构，前后端分离，简洁高效</p>
        </header>
        <div class="tech-stack__grid">
            <article v-for="stack in stacks" :key="stack.title" class="site-card tech-stack__card">
                <h3>{{ stack.title }}</h3>
                <ul>
                    <li v-for="item in stack.items" :key="item">{{ item }}</li>
                </ul>
            </article>
        </div>
    </section>
</template>

<style scoped lang="scss">
.tech-stack {
  padding: 72px 0;
  text-align: center;
}

.tech-stack h2 {
  margin: 0;
  font-size: clamp(26px, 4vw, 36px);
}

.tech-stack header > p:last-child {
  color: var(--site-foreground-secondary);
}

.tech-stack__grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--site-space-md);
  margin-top: var(--site-space-xl);
  text-align: left;
}

.tech-stack__card {
  padding: var(--site-space-lg);
}

.tech-stack__card h3 {
  margin: 0 0 var(--site-space-md);
  padding-bottom: var(--site-space-sm);
  border-bottom: 1px solid var(--site-border-subtle);
  font-size: 14px;
}

.tech-stack__card ul {
  display: flex;
  flex-direction: column;
  gap: var(--site-space-sm);
  margin: 0;
  padding: 0;
  color: var(--site-foreground-secondary);
  font-family: var(--site-font-mono);
  font-size: 12px;
  list-style: none;
}

@media (max-width: 800px) {
  .tech-stack__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 480px) {
  .tech-stack {
    padding: 48px 0;
  }

  .tech-stack__grid {
    grid-template-columns: 1fr;
  }
}
</style>
