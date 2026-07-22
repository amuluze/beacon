<script setup lang="ts">
import { getUserInfo, login } from '@/api/auth'
import useStore from '@/store'

import type { LoginForm } from '@/interface/auth'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const loginFormRef = ref<FormInstance>()

const loginForm = reactive<LoginForm>({
  username: '',
  password: '',
})

const loginFormRules = computed(() => ({
  username: [{ required: true, message: t('login.usernamePlaceholder'), trigger: 'blur' }],
  password: [{ required: true, trigger: 'blur', validator: passwordValidator }],
}))

function passwordValidator(_: any, value: string, callback: any) {
  if (value === '')
    callback(new Error(t('login.passwordRequired')))
  else if (value.length < 6)
    callback(new Error(t('login.passwordMinLength')))
  else
    callback()
}

const store = useStore()
const router = useRouter()
const submitting = shallowRef(false)
async function handleLogin() {
  // 防重入：回车既触发 @keyup.enter 又触发 form @submit.prevent，避免重复提交
  if (submitting.value)
    return
  if (!loginFormRef.value)
    return
  submitting.value = true
  try {
    // 提交前先做表单校验（用户名必填、密码长度），校验不通过则不发起登录请求
    await loginFormRef.value.validate()
    const { data } = await login({ ...loginForm })
    store.user.setToken(data.access_token, data.refresh_token)
    const userInfo = await getUserInfo()
    store.user.setUserInfo(userInfo.data.username, userInfo.data.status, userInfo.data.is_admin)
    // 每次认证会话都重新确认可用 Agent，避免沿用上一个会话的目标节点。
    store.agent.clear()
    await router.replace('/monitor')
  }
  catch (error) {
    if (error instanceof Error)
      ElMessage.error(error.message)
  }
  finally {
    submitting.value = false
  }
}

const languageList = [
  { label: '简体中文', value: 'zh' },
  { label: 'English', value: 'en' },
]

const language = computed(() => store.app.language)
function changeLanguage(lang: string) {
  locale.value = lang
  store.app.setLanguage(lang)
  router.replace('/login')
}
</script>

<template>
    <div class="am-login">
        <aside class="am-login__brand">
            <div class="am-login__logo">
                <img class="am-login__logo-mark" src="/beacon.svg" alt="" aria-hidden="true" />
                <span class="am-login__logo-text">Beacon</span>
            </div>
            <p class="am-login__slogan">
                {{ $t('login.slogan') }}
            </p>
            <p class="am-login__subline">
                {{ $t('login.subline') }}
            </p>
        </aside>

        <main class="am-login__main">
            <div class="am-login__lang">
                <el-dropdown trigger="click" @command="changeLanguage">
                    <span class="am-login__lang-text">
                        {{ language === 'zh' ? '简体中文' : 'English' }}
                        <svg-icon icon-class="down" size="12px" />
                    </span>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item
                                v-for="item in languageList"
                                :key="item.value"
                                :command="item.value"
                                :disabled="language === item.value"
                            >
                                {{ item.label }}
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>

            <section class="am-login__card">
                <h1 class="am-login__title">
                    {{ $t('login.login') }}
                </h1>

                <el-form ref="loginFormRef" :model="loginForm" :rules="loginFormRules" class="am-login__form" @submit.prevent="handleLogin">
                    <el-form-item prop="username">
                        <div class="am-login__field">
                            <svg-icon class="am-login__field-icon" icon-class="user" size="16px" />
                            <el-input v-model="loginForm.username" size="large" :placeholder="t('login.usernamePlaceholder')" @keyup.enter="handleLogin" />
                        </div>
                    </el-form-item>
                    <el-form-item prop="password">
                        <div class="am-login__field">
                            <svg-icon class="am-login__field-icon" icon-class="lock" size="16px" />
                            <el-input v-model="loginForm.password" size="large" type="password" :placeholder="t('login.passwordPlaceholder')" show-password @keyup.enter="handleLogin" />
                        </div>
                    </el-form-item>
                    <el-button class="am-login__btn" size="large" type="primary" native-type="submit" :loading="submitting">
                        {{ $t('login.login') }}
                    </el-button>
                </el-form>
            </section>
        </main>
    </div>
</template>

<style scoped lang="scss">
.am-login {
  display: flex;
  width: 100%;
  min-height: 100vh;
  background: var(--am-surface-primary);
}

.am-login__brand {
  flex: 0 0 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--am-spacing-lg);
  padding: var(--am-spacing-xl);
  color: #fff;
  background: var(--am-accent-primary);
}

.am-login__logo {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

.am-login__logo-mark {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.am-login__logo-text {
  font-size: var(--am-font-xl);
  font-weight: 700;
  color: #fff;
}

.am-login__slogan {
  max-width: 480px;
  margin: 0;
  text-align: center;
  font-size: var(--am-font-lg);
  font-weight: 500;
  color: #fff;
}

.am-login__subline {
  max-width: 480px;
  margin: 0;
  text-align: center;
  font-size: var(--am-font-md);
  color: #fff;
  opacity: 0.7;
}

.am-login__main {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--am-spacing-xl);
}

.am-login__lang {
  position: absolute;
  top: var(--am-spacing-lg);
  right: var(--am-spacing-lg);
}

.am-login__lang-text {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  cursor: pointer;
}

.am-login__card {
  width: 400px;
  max-width: 100%;
  padding: var(--am-spacing-xl);
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-lg);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-lg);
  box-shadow: var(--am-shadow-raised);
}

.am-login__title {
  margin: 0;
  text-align: center;
  font-size: var(--am-font-xl);
  font-weight: 700;
  color: var(--am-foreground-primary);
}

.am-login__form {
  .el-form-item {
    margin-bottom: 20px;
  }
}

.am-login__field {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  width: 100%;
  padding: 6px 16px;
  background: var(--am-surface-primary);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-sm);
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;

  &:focus-within {
    border-color: var(--am-accent-primary);
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--am-accent-primary) 16%, transparent);
  }

  .el-input {
    flex: 1;
  }
  :deep(.el-input__wrapper) {
    background: transparent;
    box-shadow: none !important;
  }
}

.am-login__field-icon {
  color: var(--am-foreground-muted);
  flex: 0 0 auto;
}

.am-login__btn {
  width: 100%;
  height: 46px;
  font-size: 15px;
  font-weight: 600;
  border-radius: var(--am-radius-sm);
}

@media (max-width: 960px) {
  .am-login__brand {
    display: none;
  }

  .am-login__main {
    flex: 1;
  }
}

@media (max-width: 520px) {
  .am-login__main {
    padding: var(--am-spacing-md);
  }

  .am-login__card {
    width: min(400px, 100%);
    padding: var(--am-spacing-lg);
  }

  .am-login__title {
    font-size: var(--am-font-lg);
  }
}
</style>
