<script setup lang="ts">
import type { FormRules } from 'element-plus'
import type { MailCreateArgs, MailUpdateArgs } from '@/interface/mail.ts'
import { createMail, updateMail } from '@/api/mail'
import { info } from '@/components/Message/message.ts'
import type { EmailSetting } from '@/interface/alarm.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  setting: EmailSetting
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const drawerVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 编辑发信服务配置
let updateEmailSetting = reactive<EmailSetting>({
  id: 0,
  server: '',
  port: 0,
  sender: '',
  password: '',
  receiver: '',
})

onMounted(() => {
  updateEmailSetting = props.setting
})

const loading = ref(false)

interface RuleForm {
  id: number
  server: string
  port: number
  sender: string
  password: string
  receiver: string
}

const rules = reactive<FormRules<RuleForm>>({
  server: [
    {
      required: true,
      message: '请输入邮箱服务器地址',
      trigger: 'blur',
    },
  ],
  port: [
    {
      required: true,
      message: '请输入邮箱服务器端口',
      trigger: 'blur',
    },
  ],
  sender: [
    {
      required: true,
      message: '请输入发信邮箱账号',
      trigger: 'blur',
    },
  ],
  password: [
    {
      required: true,
      message: '请输入发信邮箱密码',
      trigger: 'blur',
    },
  ],
})

function confirmMailEdit() {
  if (updateEmailSetting.id === 0) {
    loading.value = true
    const params: MailCreateArgs = {
      server: updateEmailSetting.server,
      port: Number(updateEmailSetting.port),
      sender: updateEmailSetting.sender,
      password: updateEmailSetting.password,
      receiver: updateEmailSetting.receiver,
    }
    createMail(params).finally(() => {
      info('邮件设置创建成功')
      drawerVisible.value = false
      loading.value = false
    })
  }
  else {
    loading.value = true
    const params: MailUpdateArgs = {
      id: updateEmailSetting.id,
      server: updateEmailSetting.server,
      port: Number(updateEmailSetting.port),
      sender: updateEmailSetting.sender,
      password: updateEmailSetting.password,
      receiver: updateEmailSetting.receiver,
    }
    updateMail(params).finally(() => {
      info('邮件设置更新成功')
      drawerVisible.value = false
      loading.value = false
    })
  }
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="50%" :title="t(props.title as string)">
        <el-form :model="updateEmailSetting" :rules="rules" label-width="160px" label-position="left">
            <el-form-item prop="server" :label="t('setting.mailServerHost')">
                <el-input v-model="updateEmailSetting.server" :placeholder="t('setting.hostPlaceholder')" />
            </el-form-item>
            <el-form-item prop="port" :label="t('setting.mailServerPort')">
                <el-input v-model="updateEmailSetting.port" :placeholder="t('setting.portPlaceholder')" />
            </el-form-item>
            <el-form-item prop="sender" :label="t('setting.mailServerAccount')">
                <el-input v-model="updateEmailSetting.sender" :placeholder="t('setting.senderPlaceholder')" />
            </el-form-item>
            <el-form-item prop="password" :label="t('setting.mailServerPassword')">
                <el-input v-model="updateEmailSetting.password" type="password" :placeholder="t('setting.passwordPlaceholder')" show-password />
            </el-form-item>
            <el-form-item prop="receiver" :label="t('setting.mailReceiver')">
                <el-input v-model="updateEmailSetting.receiver" :placeholder="t('setting.receiverPlaceholder')" />
            </el-form-item>
        </el-form>
        <el-button type="primary" size="default" plain @click="drawerVisible = false">
            {{ t('setting.cancel') }}
        </el-button>
        <el-button type="primary" size="default" plain @click="confirmMailEdit">
            {{ t('setting.confirm') }}
        </el-button>
    </el-drawer>
</template>

<style scoped lang="scss">
</style>
