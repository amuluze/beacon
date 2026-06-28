<template>
  <button
    :class="[
      'tech-button',
      `tech-button--${variant}`,
      { 'tech-button--loading': loading }
    ]"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <div class="tech-button__content">
      <Icon v-if="icon && !loading" :name="icon" class="tech-button__icon" />
      <Icon v-if="loading" name="mdi:loading" class="tech-button__icon tech-button__icon--loading" />
      <span class="tech-button__text">
        <slot />
      </span>
    </div>
    <div class="tech-button__shine"></div>
    <div class="tech-button__glow"></div>
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'ghost'
  icon?: string
  loading?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  loading: false,
  disabled: false
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped lang="scss">
.tech-button {
  position: relative;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  @include glass-effect(0.1);

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  &__content {
    position: relative;
    z-index: 2;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  &__icon {
    font-size: 18px;
    transition: transform 0.3s ease;

    &--loading {
      animation: spin 1s linear infinite;
    }
  }

  &__text {
    white-space: nowrap;
  }

  &__shine {
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    transition: left 0.6s ease;
    z-index: 1;
  }

  &__glow {
    position: absolute;
    inset: -2px;
    border-radius: inherit;
    padding: 2px;
    background: linear-gradient(135deg, transparent, rgba(0, 102, 255, 0.5), transparent);
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask-composite: exclude;
    -webkit-mask-composite: xor;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  &:hover:not(:disabled) {
    transform: translateY(-2px);

    .tech-button__shine {
      left: 100%;
    }

    .tech-button__glow {
      opacity: 1;
    }

    .tech-button__icon {
      transform: scale(1.1);
    }
  }

  &--primary {
    @include bg-gradient('btn-primary');
    color: white;
    @include glow-effect(#0066ff, 0.3);

    &:hover:not(:disabled) {
      @include bg-gradient('btn-secondary');
    }
  }

  &--secondary {
    @include bg-gradient('btn-accent');
    color: white;
    @include glow-effect(#7c3aed, 0.3);

    &:hover:not(:disabled) {
      background: linear-gradient(135deg, rgba(124, 58, 237, 0.9), rgba(0, 212, 255, 0.9));
    }
  }

  &--ghost {
    background: transparent;
    color: #e2e8f0;
    border: 1px solid rgba(0, 102, 255, 0.3);

    &:hover:not(:disabled) {
      background: rgba(0, 102, 255, 0.1);
      border-color: rgba(0, 102, 255, 0.5);
    }
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
