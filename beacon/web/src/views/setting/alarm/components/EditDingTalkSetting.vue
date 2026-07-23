<script setup lang="ts">
import { updateDingTalk } from '@/api/dingtalk'
import { info, warning } from '@/components/Message/message.ts'
import type { DingTalkSetting, DingTalkUpdateArgs } from '@/interface/dingtalk'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  setting: DingTalkSetting
}>()

const emit = defineEmits<{
  'update:visible': [visible: boolean]
  'close': []
  'saved': []
}>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
    if (!visible)
      emit('close')
  },
})

const formRef = ref<FormInstance>()
const form = reactive<DingTalkUpdateArgs>({
  enabled: props.setting.enabled,
  webhook: '',
  secret: '',
  clear_secret: false,
  at_all: props.setting.at_all,
})
const loading = shallowRef(false)
const { t } = useI18n()

async function submit() {
  if (form.enabled && !props.setting.webhook_configured && !form.webhook.trim()) {
    warning(t('setting.dingTalkWebhookRequired'))
    return
  }
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid)
    return

  loading.value = true
  try {
    await updateDingTalk({
      ...form,
      webhook: form.webhook.trim(),
      secret: form.secret.trim(),
    })
    info(t('setting.dingTalkUpdated'))
    emit('saved')
    dialogVisible.value = false
  }
  catch {
    // 请求拦截器负责展示服务端错误；保留本地表单供用户重试。
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
    <el-dialog v-model="dialogVisible" width="520px" :title="t(props.title ?? 'setting.dingTalkSetting')">
        <el-form ref="formRef" :model="form" label-position="top">
            <el-form-item :label="t('setting.dingTalkEnabled')">
                <el-switch v-model="form.enabled" />
            </el-form-item>
            <el-form-item prop="webhook" :label="t('setting.dingTalkWebhook')">
                <el-input
                    v-model="form.webhook"
                    type="password"
                    show-password
                    autocomplete="off"
                    :placeholder="setting.webhook_configured ? t('setting.credentialKeepPlaceholder') : t('setting.dingTalkWebhookPlaceholder')"
                />
                <p v-if="setting.webhook_configured" class="credential-hint">
                    {{ t('setting.currentCredential') }}：{{ setting.webhook_masked }}
                </p>
            </el-form-item>
            <el-form-item prop="secret" :label="t('setting.dingTalkSecret')">
                <el-input
                    v-model="form.secret"
                    type="password"
                    show-password
                    autocomplete="off"
                    :disabled="form.clear_secret"
                    :placeholder="setting.secret_configured ? t('setting.credentialKeepPlaceholder') : t('setting.dingTalkSecretPlaceholder')"
                />
                <el-checkbox v-if="setting.secret_configured" v-model="form.clear_secret">
                    {{ t('setting.clearDingTalkSecret') }}
                </el-checkbox>
            </el-form-item>
            <el-form-item :label="t('setting.dingTalkMention')">
                <el-checkbox v-model="form.at_all">
                    {{ t('setting.dingTalkAtAll') }}
                </el-checkbox>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="dialogVisible = false">
                {{ t('setting.cancel') }}
            </el-button>
            <el-button :loading="loading" type="primary" @click="submit">
                {{ t('setting.confirm') }}
            </el-button>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
.credential-hint {
  margin: var(--am-spacing-xs) 0 0;
  overflow-wrap: anywhere;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}
</style>
