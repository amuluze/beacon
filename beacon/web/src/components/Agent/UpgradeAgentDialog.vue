<script setup lang="ts">
import { upgradeAgent } from '@/api/agent'
import { success } from '@/components/Message/message'
import type { AgentUpgradeResult } from '@/interface/agent.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  agentID?: string
  hostname?: string
  version?: string
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
  (e: 'success', result: AgentUpgradeResult): void
}>()

const { t } = useI18n()

const dialogVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible)
      emits('close')
  },
})

const submitting = shallowRef(false)

const targetNode = computed(() => {
  const name = props.hostname || props.agentID || '-'
  return props.agentID ? `${name} · ${props.agentID}` : name
})

const tips = computed(() => [
  t('agent.upgradeTipConcurrency'),
  t('agent.upgradeTipInterrupt'),
])

async function submit() {
  if (!props.agentID || submitting.value)
    return

  submitting.value = true
  try {
    const { data } = await upgradeAgent(props.agentID)
    emits('success', data)
    success(t('agent.upgradeSuccess', {
      from: data?.current_version || '-',
      to: data?.target_version || '-',
    }))
    dialogVisible.value = false
  }
  catch {
    // 失败提示由全局响应拦截器统一给出，保留弹窗与目标节点便于重试。
  }
  finally {
    submitting.value = false
  }
}
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="560px">
        <div class="am-upgrade-agent">
            <p class="am-upgrade-agent__desc">
                {{ t('agent.upgradeDesc') }}
            </p>
            <dl class="am-upgrade-agent__facts">
                <div class="am-upgrade-agent__fact">
                    <dt class="am-upgrade-agent__label">
                        {{ t('agent.upgradeTargetLabel') }}
                    </dt>
                    <dd class="am-upgrade-agent__value">
                        {{ targetNode }}
                    </dd>
                </div>
                <div class="am-upgrade-agent__fact">
                    <dt class="am-upgrade-agent__label">
                        {{ t('agent.upgradeCurrentVersionLabel') }}
                    </dt>
                    <dd class="am-upgrade-agent__value">
                        {{ props.version || '-' }}
                    </dd>
                </div>
                <div class="am-upgrade-agent__fact">
                    <dt class="am-upgrade-agent__label">
                        {{ t('agent.upgradeTargetVersionLabel') }}
                    </dt>
                    <dd class="am-upgrade-agent__value">
                        {{ t('agent.upgradeTargetVersionValue') }}
                    </dd>
                </div>
            </dl>
            <ul class="am-upgrade-agent__tips">
                <li v-for="tip in tips" :key="tip">
                    {{ tip }}
                </li>
            </ul>
        </div>
        <template #footer>
            <el-button :disabled="submitting" @click="dialogVisible = false">
                {{ t('agent.upgradeCancel') }}
            </el-button>
            <el-button type="primary" :loading="submitting" @click="submit">
                {{ t('agent.upgradeConfirm') }}
            </el-button>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
@include b(upgrade-agent) {
  &__desc {
    margin: 0 0 16px;
    font-size: 13px;
    line-height: 1.6;
    color: var(--am-foreground-secondary);
  }

  &__facts {
    margin: 0;
    padding: 12px 16px;
    background: var(--am-surface-secondary);
    border-radius: 6px;
  }

  &__fact {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--am-spacing-md);

    & + & {
      margin-top: 8px;
    }
  }

  &__label {
    margin: 0;
    font-size: 12px;
    color: var(--am-foreground-muted);
    white-space: nowrap;
  }

  &__value {
    margin: 0;
    min-width: 0;
    overflow: hidden;
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: 12px;
    font-weight: 500;
    color: var(--am-foreground-primary);
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__tips {
    margin: 16px 0 0;
    padding: 12px 16px;
    list-style: none;
    background: var(--am-surface-secondary);
    border-radius: 6px;

    li {
      font-size: 12px;
      line-height: 1.8;
      color: var(--am-foreground-secondary);
    }
  }
}
</style>
