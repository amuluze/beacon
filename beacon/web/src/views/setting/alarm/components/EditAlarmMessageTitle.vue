<script setup lang="ts">
import { updateAlarmMessageSetting } from '@/api/alarm'
import { success } from '@/components/Message/message.ts'
import type { AlarmMessageSetting } from '@/interface/alarm'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  setting: AlarmMessageSetting
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

const { t } = useI18n()
const formRef = ref<FormInstance>()
const form = reactive<AlarmMessageSetting>({
  message_title: props.setting.message_title,
})
const loading = shallowRef(false)

function titleValidationMessage(value: string): string | undefined {
  if (/[\r\n]/.test(value))
    return t('setting.alarmMessageTitleNoLineBreak')

  const normalized = value.trim()
  if (!normalized)
    return t('setting.alarmMessageTitleRequired')
  if ([...normalized].length > 128)
    return t('setting.alarmMessageTitleTooLong')
}

const rules: FormRules = {
  message_title: [{
    validator: (_rule, value: string, callback) => {
      const message = titleValidationMessage(value)
      callback(message ? new Error(message) : undefined)
    },
    trigger: ['blur', 'change'],
  }],
}

async function submit() {
  if (titleValidationMessage(form.message_title)) {
    await formRef.value?.validate().catch(() => false)
    return
  }

  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid)
    return

  loading.value = true
  try {
    await updateAlarmMessageSetting({ message_title: form.message_title.trim() })
    success(t('setting.alarmMessageTitleUpdated'))
    emit('saved')
    dialogVisible.value = false
  }
  catch {
    // 请求拦截器负责展示服务端错误；保留弹窗和本地输入供用户重试。
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
    <el-dialog v-model="dialogVisible" width="480px" :title="t(props.title ?? 'setting.alarmMessageTitle')">
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-form-item prop="message_title" :label="t('setting.alarmMessageTitle')">
                <el-input
                    v-model="form.message_title"
                    :placeholder="t('setting.alarmMessageTitlePlaceholder')"
                    autocomplete="off"
                />
                <p class="message-title-hint">
                    {{ t('setting.alarmMessageTitleTips') }}
                </p>
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
.message-title-hint {
  margin: var(--am-spacing-xs) 0 0;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}
</style>
