<script setup lang="ts">
interface Props {
    overline: string
    title: string
    description: string
    points: string[]
    reverse?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    reverse: false,
})
</script>

<template>
    <section class="feature" :class="{ 'feature--reverse': props.reverse }">
        <div class="feature__preview"><slot /></div>
        <div class="feature__content">
            <p class="site-overline">{{ props.overline }}</p>
            <h2>{{ props.title }}</h2>
            <p class="feature__description">{{ props.description }}</p>
            <ul>
                <li v-for="point in props.points" :key="point">
                    <Icon name="mdi:check-circle" />
                    <span>{{ point }}</span>
                </li>
            </ul>
        </div>
    </section>
</template>

<style scoped lang="scss">
.feature {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(320px, 0.9fr);
  gap: 80px;
  align-items: center;
  padding: 72px 0;
}

.feature--reverse .feature__preview {
  order: 2;
}

.feature__content h2 {
  margin: 0;
  font-size: clamp(26px, 4vw, 36px);
  line-height: 1.25;
}

.feature__description {
  margin: var(--site-space-md) 0 var(--site-space-lg);
  color: var(--site-foreground-secondary);
  font-size: 16px;
}

.feature__content ul {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.feature__content li {
  display: flex;
  align-items: flex-start;
  gap: var(--site-space-sm);
  color: var(--site-foreground-secondary);
}

.feature__content li :deep(svg) {
  flex: 0 0 auto;
  margin-top: 3px;
  color: var(--site-success);
}

@media (max-width: 800px) {
  .feature {
    grid-template-columns: 1fr;
    gap: var(--site-space-xl);
    padding: 48px 0;
  }

  .feature--reverse .feature__preview {
    order: 0;
  }
}
</style>
