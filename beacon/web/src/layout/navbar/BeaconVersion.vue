<script setup lang="ts">
import { querySystemVersion } from '@/api/system'
import type { SystemVersion } from '@/interface/system.ts'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const version = ref<SystemVersion>()
const loading = ref(false)
const detailVisible = ref(false)

// 升级图标常驻顶栏：无新版本时置灰并显示当前版本，有新版本时高亮并可查看升级指引。
// Beacon 自身不执行升级，弹层只做只读说明，是否升级由管理员手动执行。
const hasUpdate = computed(() => version.value?.update_available === true)

const triggerTitle = computed(() => {
  if (loading.value)
    return t('system.versionLoading')
  if (hasUpdate.value)
    return t('system.versionUpdateAvailable', { version: version.value?.latest_version || '-' })
  return t('system.versionCurrent', { version: version.value?.version || '-' })
})

async function loadVersion() {
  loading.value = true
  try {
    const { data } = await querySystemVersion()
    version.value = data
  }
  catch {
    version.value = undefined
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  void loadVersion()
})

function openDetail() {
  if (!hasUpdate.value)
    return
  detailVisible.value = true
}
</script>

<template>
    <button
        class="am-version-trigger"
        :class="{ 'am-version-trigger--update': hasUpdate }"
        type="button"
        :disabled="!hasUpdate"
        :title="triggerTitle"
        @click="openDetail"
    >
        <svg-icon icon-class="update" size="14px" />
        <span v-if="hasUpdate" class="am-version-trigger__label">{{ t('system.newVersion') }}</span>
    </button>

    <el-dialog v-model="detailVisible" :title="t('system.versionDialogTitle')" width="480px">
        <div class="am-version-detail">
            <div class="am-version-detail__row">
                <span class="am-version-detail__label">{{ t('system.currentVersion') }}</span>
                <span class="am-version-detail__value">{{ version?.version || '-' }}</span>
            </div>
            <div class="am-version-detail__row">
                <span class="am-version-detail__label">{{ t('system.latestVersion') }}</span>
                <span class="am-version-detail__value">{{ version?.latest_version || '-' }}</span>
            </div>
            <p class="am-version-detail__tip">
                {{ t('system.versionUpgradeTip') }}
            </p>
        </div>
    </el-dialog>
</template>

<style scoped lang="scss">
.am-version-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  gap: var(--am-spacing-xs);
  height: 32px;
  padding: 0 12px;
  color: var(--am-foreground-secondary);
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--am-radius-sm);
  cursor: pointer;

  &:hover:not(:disabled) {
    color: var(--am-foreground-primary);
    background: var(--am-surface-inset);
  }

  &:disabled {
    cursor: default;
    opacity: 0.6;
  }

  &--update {
    color: var(--am-accent-primary);
    border-color: var(--am-accent-primary-soft);
    background: var(--am-accent-primary-soft);
  }
}

.am-version-detail {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);

  &__row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--am-spacing-md);
  }

  &__label {
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
  }

  &__value {
    color: var(--am-foreground-primary);
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: var(--am-font-xs);
    font-weight: 500;
  }

  &__tip {
    margin: 0;
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
    line-height: 1.6;
  }
}
</style>
