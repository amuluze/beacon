<script setup lang="ts">
const stacks = [
    { title: '前端', icon: 'lucide:monitor', items: ['Vue 3', 'TypeScript', 'Element Plus', 'Pinia'] },
    { title: '后端', icon: 'lucide:server', items: ['Golang', 'Fiber', 'Wire DI'] },
    { title: '存储', icon: 'lucide:database', items: ['SQLite', 'GORM'] },
    { title: '传输', icon: 'lucide:share-2', items: ['gRPC Tunnel', 'WebSocket', 'HTTP/2'] },
]
</script>

<template>
    <section class="tech-stack">
        <div class="site-container tech-stack__inner">
            <header class="tech-stack__header">
                <h2>现代化技术栈驱动</h2>
                <p>Server-Agent 架构，前后端分离，简洁高效</p>
            </header>
            <div class="tech-stack__grid">
                <article v-for="stack in stacks" :key="stack.title" class="tech-stack__card">
                    <div class="tech-stack__head">
                        <Icon :name="stack.icon" />
                        <h3>{{ stack.title }}</h3>
                    </div>
                    <div class="tech-stack__divider" />
                    <ul class="tech-stack__list">
                        <li v-for="item in stack.items" :key="item">
                            {{ item }}
                        </li>
                    </ul>
                </article>
            </div>
        </div>
    </section>
</template>

<style scoped lang="scss">
.tech-stack {
  padding: 96px 0;
  background: var(--color-surface-muted);
  border-top: 1px solid var(--border);
}

.tech-stack__header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  text-align: center;
}

.tech-stack h2 {
  margin: 0;
  font-size: 30px;
  font-weight: 700;
  line-height: 1.2;
}

.tech-stack__header p {
  max-width: 480px;
  margin: 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-md);
  line-height: 1.6;
}

.tech-stack__grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--space-4);
  margin-top: var(--space-10);
}

.tech-stack__card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  height: 199px;
  padding: var(--space-6);
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
}

.tech-stack__head {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--primary);
}

.tech-stack__head :deep(svg) {
  width: 18px;
  height: 18px;
}

.tech-stack__head h3 {
  margin: 0;
  color: var(--foreground);
  font-size: var(--font-size-md);
  font-weight: 600;
}

.tech-stack__divider {
  height: 1px;
  background: var(--border);
}

.tech-stack__list {
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
  margin: 0;
  padding: 0;
  color: var(--color-text-secondary);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  list-style: none;
}

@media (max-width: 800px) {
  .tech-stack {
    padding: 48px 0;
  }

  .tech-stack__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 480px) {
  .tech-stack__grid {
    grid-template-columns: 1fr;
  }
}
</style>
