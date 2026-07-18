<script setup lang="ts">
interface Props {
    overline: string
    overlineIcon?: string
    title: string
    description: string
    points: string[]
    reverse?: boolean
    variant?: 'default' | 'muted'
}

const props = withDefaults(defineProps<Props>(), {
    overlineIcon: 'lucide:package',
    reverse: false,
    variant: 'default',
})
</script>

<template>
    <section class="feature" :class="[`feature--${props.variant}`, { 'feature--reverse': props.reverse }]">
        <div class="site-container feature__inner">
            <div class="feature__preview">
                <slot />
            </div>
            <div class="feature__content">
                <p class="feature__overline">
                    <Icon :name="props.overlineIcon" />
                    <span>{{ props.overline }}</span>
                </p>
                <h2>{{ props.title }}</h2>
                <p class="feature__description">
                    {{ props.description }}
                </p>
                <ul>
                    <li v-for="point in props.points" :key="point">
                        <span class="feature__check"><Icon name="lucide:check" /></span>
                        <span>{{ point }}</span>
                    </li>
                </ul>
            </div>
        </div>
    </section>
</template>

<style scoped lang="scss">
.feature {
  border-top: 1px solid var(--border);
}

.feature--muted {
  background: var(--color-surface-muted);
}

.feature__inner {
  display: grid;
  grid-template-columns: minmax(0, 580px) minmax(320px, 460px);
  gap: 80px;
  align-items: center;
  padding: 96px 0;
}

.feature--reverse .feature__preview {
  order: 2;
}

.feature__overline {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  margin: 0 0 var(--space-2);
  color: var(--primary);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  letter-spacing: var(--letter-spacing-wide);
}

.feature__overline :deep(svg) {
  width: 16px;
  height: 16px;
}

.feature__content h2 {
  margin: 0;
  font-size: var(--font-display-sm);
  font-weight: var(--font-weight-bold);
  line-height: var(--line-height-snug);
}

.feature__description {
  margin: var(--space-4) 0 var(--space-5);
  color: var(--muted-foreground);
  font-size: var(--font-size-md);
  line-height: var(--line-height-relaxed);
}

.feature__content ul {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.feature__content li {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
}

.feature__check {
  display: inline-grid;
  place-items: center;
  flex: 0 0 auto;
  width: 20px;
  height: 20px;
  margin-top: 1px;
  color: var(--color-text-inverse);
  background: var(--color-success);
  border-radius: var(--radius-pill);
  font-size: var(--font-size-sm);
}

@media (max-width: 800px) {
  .feature__inner {
    grid-template-columns: 1fr;
    gap: var(--space-8);
    padding: 48px 0;
  }

  .feature--reverse .feature__preview {
    order: 0;
  }
}
</style>
