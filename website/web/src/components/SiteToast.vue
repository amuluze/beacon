<script setup lang="ts">
import type { ToastType } from '~/composables/useToast'

import { useToast } from '~/composables/useToast'

const { toast, dismissToast } = useToast()

const iconNames: Record<ToastType, string> = {
    error: 'mdi:alert-circle-outline',
    info: 'mdi:information-outline',
    success: 'mdi:check-circle-outline',
}
</script>

<template>
    <Transition name="site-toast">
        <div
            v-if="toast"
            :key="toast.id"
            class="site-toast"
            :class="`site-toast--${toast.type}`"
            :role="toast.type === 'error' ? 'alert' : 'status'"
            :aria-live="toast.type === 'error' ? 'assertive' : 'polite'"
        >
            <Icon :name="iconNames[toast.type]" aria-hidden="true" />
            <span>{{ toast.message }}</span>
            <button type="button" aria-label="关闭提示" @click="dismissToast">
                <Icon name="mdi:close" aria-hidden="true" />
            </button>
        </div>
    </Transition>
</template>

<style scoped lang="scss">
.site-toast {
  position: fixed;
  z-index: 300;
  top: 80px;
  left: 50%;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: var(--space-2);
  align-items: center;
  width: min(420px, calc(100% - 32px));
  padding: 12px 14px;
  color: var(--foreground);
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--radius-m);
  box-shadow: var(--shadow-card);
  transform: translateX(-50%);
}

.site-toast--success > :deep(svg:first-child) {
  color: var(--color-success);
}

.site-toast--error > :deep(svg:first-child) {
  color: var(--color-error);
}

.site-toast--info > :deep(svg:first-child) {
  color: var(--color-link);
}

.site-toast button {
  display: grid;
  place-items: center;
  width: 28px;
  height: 28px;
  padding: 0;
  color: var(--color-text-secondary);
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.site-toast-enter-active,
.site-toast-leave-active {
  transition:
    opacity 0.18s ease,
    transform 0.18s ease;
}

.site-toast-enter-from,
.site-toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -8px);
}
</style>
