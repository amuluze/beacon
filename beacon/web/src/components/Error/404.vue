<script setup lang="ts">
import IconArrowLeft from '~icons/lucide/arrow-left'
import IconFileQuestion from '~icons/lucide/file-question'

const router = useRouter()

function goHome() {
  // 草稿：点击返回首页回到监控工作区
  router.replace('/')
}
</script>

<template>
    <main class="am-error-page">
        <div class="am-error-page__body">
            <IconFileQuestion class="am-error-page__icon" aria-hidden="true" />
            <p class="am-error-page__code">
                404!
            </p>
            <p class="am-error-page__message">
                {{ $t('error.notFound') }}
            </p>
            <button class="am-error-page__action" type="button" @click="goHome">
                <IconArrowLeft aria-hidden="true" />
                <span>{{ $t('error.backHome') }}</span>
            </button>
        </div>
    </main>
</template>

<style scoped lang="scss">
// 稿面 `404`：独立错误页，不套用工作区导航栏与状态栏（DESIGN.md §10.4）。
.am-error-page {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 100vh;
  padding: var(--am-spacing-xl);
  color: var(--am-foreground-primary);
  background: var(--am-surface-primary);
}

.am-error-page__body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--am-spacing-sm);
  text-align: center;
}

.am-error-page__icon {
  width: 32px;
  height: 32px;
  color: var(--am-foreground-muted);
}

.am-error-page__code {
  margin: 0;
  font-size: var(--am-font-2xl);
  font-weight: 700;
  color: var(--am-foreground-primary);
}

.am-error-page__message {
  margin: 0;
  font-size: var(--am-font-md);
  color: var(--am-foreground-secondary);
}

.am-error-page__action {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: var(--am-spacing-lg);
  padding: 8px 16px;
  color: var(--am-foreground-on-accent);
  font-family: inherit;
  font-size: var(--am-font-sm);
  font-weight: 600;
  background: var(--am-accent-primary);
  border: 0;
  border-radius: var(--am-radius-sm);
  box-shadow: 0 1px 2px #0f172a1a;
  cursor: pointer;
  transition: background-color 0.2s ease;

  &:hover {
    background: var(--am-accent-primary-hover);
  }

  svg {
    width: 14px;
    height: 14px;
  }
}
</style>
