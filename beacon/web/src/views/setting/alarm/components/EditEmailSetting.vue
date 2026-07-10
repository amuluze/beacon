<script setup lang="ts">
import { createMail, updateMail } from '@/api/mail'
import { info } from '@/components/Message/message.ts'
import type { EmailSetting } from '@/interface/alarm.ts'
import type { MailCreateArgs, MailUpdateArgs } from '@/interface/mail.ts'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  setting: EmailSetting
}>()

const emit = defineEmits<{
  'update:visible': [visible: boolean]
  'close': []
  'saved': [setting: EmailSetting]
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
const form = reactive<EmailSetting>(structuredClone(toRaw(props.setting)))
const loading = shallowRef(false)
const { t } = useI18n()

const rules: FormRules = {
  server: [{ required: true, message: '请输入邮箱服务器地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入邮箱服务器端口', trigger: 'blur' }],
  sender: [{ required: true, message: '请输入发信邮箱账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入发信邮箱密码', trigger: 'blur' }],
}

async function submit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid)
    return

  loading.value = true
  try {
    if (form.id === 0) {
      const params: MailCreateArgs = {
        server: form.server,
        port: Number(form.port),
        sender: form.sender,
        password: form.password,
        receiver: form.receiver,
      }
      await createMail(params)
      info('邮件设置创建成功')
    }
    else {
      const params: MailUpdateArgs = {
        id: form.id,
        server: form.server,
        port: Number(form.port),
        sender: form.sender,
        password: form.password,
        receiver: form.receiver,
      }
      await updateMail(params)
      info('邮件设置更新成功')
    }
    emit('saved', structuredClone(toRaw(form)))
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
    <el-dialog v-model="dialogVisible" width="480px" :title="t(props.title ?? 'setting.mailServerSetting')">
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-form-item prop="server" :label="t('setting.mailServerHost')">
                <el-input v-model="form.server" :placeholder="t('setting.hostPlaceholder')" />
            </el-form-item>
            <el-form-item prop="port" :label="t('setting.mailServerPort')">
                <el-input v-model="form.port" :placeholder="t('setting.portPlaceholder')" />
            </el-form-item>
            <el-form-item prop="sender" :label="t('setting.mailServerAccount')">
                <el-input v-model="form.sender" :placeholder="t('setting.senderPlaceholder')" />
            </el-form-item>
            <el-form-item prop="password" :label="t('setting.mailServerPassword')">
                <el-input v-model="form.password" type="password" show-password :placeholder="t('setting.passwordPlaceholder')" />
            </el-form-item>
            <el-form-item prop="receiver" :label="t('setting.mailReceiver')">
                <el-input v-model="form.receiver" :placeholder="t('setting.receiverPlaceholder')" />
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

<style scoped lang="scss"></style>
