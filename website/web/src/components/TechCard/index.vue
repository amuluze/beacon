<template>
  <div
    :class="[
      'tech-card',
      `tech-card--${variant}`,
      { 'tech-card--hoverable': hoverable }
    ]"
  >
    <div class="tech-card__glow"></div>
    <div class="tech-card__content">
      <div v-if="$slots.header" class="tech-card__header">
        <slot name="header" />
      </div>
      <div class="tech-card__body">
        <slot />
      </div>
      <div v-if="$slots.footer" class="tech-card__footer">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'default' | 'primary' | 'secondary'
  hoverable?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'default',
  hoverable: true
})
</script>

<style scoped lang="scss">
.tech-card {
  position: relative;
  border-radius: 12px;
  @include glass-effect(0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;

  &__glow {
    position: absolute;
    inset: -1px;
    border-radius: inherit;
    padding: 1px;
    background: linear-gradient(135deg, transparent, rgba(0, 102, 255, 0.3), transparent);
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask-composite: exclude;
    -webkit-mask-composite: xor;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  &__content {
    position: relative;
    z-index: 2;
    padding: 24px;
  }

  &__header {
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid rgba(0, 102, 255, 0.2);
  }

  &__body {
    flex: 1;
  }

  &__footer {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid rgba(0, 102, 255, 0.2);
  }

  &--hoverable {
    &:hover {
      transform: translateY(-4px);
      @include glass-effect(0.1);

      .tech-card__glow {
        opacity: 1;
      }
    }
  }

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    @include bg-gradient('glow');
  }

  &--primary {
    @include bg-gradient('btn-primary');
  }

  &--secondary {
    @include bg-gradient('btn-accent');
  }
}
</style>
