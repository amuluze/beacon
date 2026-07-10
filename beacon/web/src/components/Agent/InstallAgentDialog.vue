<script setup lang="ts">
import { getInstallToken } from '@/api/host'
import { success } from '@/components/Message/message'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  title?: string
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

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

const loading = shallowRef(false)
const token = ref('')
const node = ref('node-01')
const nodeRegex = /^[A-Za-z0-9][A-Za-z0-9._-]*$/
const nodeError = ref('')

const origin = window.location.origin

const curlCommand = computed(() => {
  const safeNode = nodeRegex.test(node.value) ? node.value : 'node-01'
  const tk = token.value || '<your-install-token>'
  return `curl -kfsSL '${origin}/api/v1/host/install?node=${safeNode}' | sudo bash -s -- --token=${tk}`
})

const tips = computed(() => [
  t('agent.installTipRoot'),
  t('agent.installTipArch'),
  t('agent.installTipService'),
])

function validateNode() {
  nodeError.value = node.value && !nodeRegex.test(node.value)
    ? t('agent.installNodeInvalid')
    : ''
}

async function loadToken() {
  loading.value = true
  try {
    const { data } = await getInstallToken()
    token.value = data?.token ?? ''
  }
  finally {
    loading.value = false
  }
}

async function copyCommand() {
  const text = curlCommand.value
  if (!text)
    return
  try {
    await navigator.clipboard.writeText(text)
    success(t('agent.installCopied'))
  }
  catch {
    // 剪贴板不可用时静默忽略
  }
}

onMounted(loadToken)
</script>

<template>
  <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="560px">
    <div v-loading="loading" class="am-install-agent">
      <p class="am-install-agent__desc">
        {{ t('agent.installDesc') }}
      </p>
      <el-form label-position="top">
        <el-form-item :label="t('agent.installServerLabel')">
          <el-input :model-value="origin" readonly />
        </el-form-item>
        <el-form-item :label="t('agent.installNodeLabel')" :error="nodeError">
          <el-input v-model="node" :placeholder="t('agent.installNodePlaceholder')" @input="validateNode" />
        </el-form-item>
        <el-form-item :label="t('agent.installTokenLabel')">
          <el-input v-model="token" type="password" show-password readonly />
        </el-form>
      </el-form>
      <div class="am-install-agent__code">
        <div class="am-install-agent__code-head">
          <span>{{ t('agent.installCommandLabel') }}</span>
          <el-button text size="small" @click="copyCommand">
            {{ t('agent.installCopy') }}
          </el-button>
        </div>
        <pre class="am-install-agent__code-body">{{ curlCommand }}</pre>
      </div>
      <ul class="am-install-agent__tips">
        <li v-for="tip in tips" :key="tip">{{ tip }}</li>
      </ul>
    </div>
  </el-dialog>
</template>

<style scoped lang="scss">
@include b(install-agent) {
  &__desc {
    margin: 0 0 16px;
    font-size: 13px;
    line-height: 1.6;
    color: var(--am-foreground-secondary);
  }

  &__code {
    margin-top: 16px;
    overflow: hidden;
    border: 1px solid var(--am-border-subtle);
    border-radius: 6px;
  }

  &__code-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: var(--am-surface-secondary);

    span {
      font-size: 12px;
      font-weight: 600;
      color: var(--am-foreground-secondary);
    }
  }

  &__code-body {
    margin: 0;
    padding: 12px 14px;
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: 13px;
    line-height: 1.6;
    color: #d4d4d4;
    white-space: pre-wrap;
    word-break: break-all;
    background: #1e1e1e;
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
