<script setup lang="ts">
import { getUserInfo, login } from '@/api/auth'
import useStore from '@/store'

import type { LoginForm } from '@/interface/auth'
import { useI18n } from 'vue-i18n'

const loginForm = reactive<LoginForm>({
  username: '',
  password: '',
})

const loginFormRules = {
  username: [{ required: true, trigger: 'blur' }],
  password: [{ required: true, trigger: 'blur', validator: passwordValidator }],
}

function passwordValidator(_: any, value: string, callback: any) {
  if (value === '') {
    callback(new Error('password is required'))
  }
  else if (value.length < 6) {
    callback(new Error('The password can not be less than 6 digits'))
  }
  else {
    callback()
  }
}

const store = useStore()
const router = useRouter()
async function handleLogin() {
  try {
    const { data } = await login({ ...loginForm })
    store.user.setToken(data.access_token, data.refresh_token)
    const userInfo = await getUserInfo()
    store.user.setUserInfo(userInfo.data.username, userInfo.data.status, userInfo.data.is_admin)
    await router.replace('/monitor/host')
  }
  catch (error) {
    if (error instanceof Error)
      ElMessage.error(error.message)
  }
}

const languageList = [
  { label: '简体中文', value: 'zh' },
  { label: 'English', value: 'en' },
]

const i18n = useI18n()
const language = computed(() => store.app.language)
function changeLanguage(lang: string) {
  i18n.locale.value = lang
  store.app.setLanguage(lang)
  router.replace('/login')
}
</script>

<template>
    <div class="am-login">
        <main class="am-login__main">
            <section class="am-login__card">
                <div class="am-login__lang">
                    <el-dropdown trigger="click" @command="changeLanguage">
                        <span class="am-login__lang-text">{{ language === 'zh' ? '简体中文' : 'English' }} ▼</span>
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

                <div class="am-login__title-row">
                    <span class="am-login__title">登录</span>
                    <span class="am-login__version">v3.0.0</span>
                </div>

                <el-form :model="loginForm" :rules="loginFormRules" class="am-login__form">
                    <el-form-item prop="username">
                        <div class="am-login__field">
                            <svg-icon class="am-login__field-icon" icon-class="user" size="16px" />
                            <el-input v-model="loginForm.username" size="large" placeholder="请输入用户名" />
                        </div>
                    </el-form-item>
                    <el-form-item prop="password">
                        <div class="am-login__field">
                            <svg-icon class="am-login__field-icon" icon-class="lock" size="16px" />
                            <el-input v-model="loginForm.password" size="large" type="password" placeholder="请输入密码" show-password />
                        </div>
                    </el-form-item>
                    <el-button class="am-login__btn" size="large" type="primary" @click.prevent="handleLogin">
                        登 录
                    </el-button>
                </el-form>
            </section>
        </main>
    </div>
</template>

<style scoped lang="scss">
.am-login {
  width: 100%;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--am-surface-primary);

  &__main {
    width: 100%;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--am-spacing-xl);
  }

  &__card {
    width: 400px;
    max-width: 100%;
    padding: var(--am-spacing-xl);
    display: flex;
    flex-direction: column;
    gap: var(--am-spacing-lg);
    background: var(--am-surface-card);
    border: 1px solid var(--am-border-subtle);
    border-radius: var(--am-radius-md);
    box-shadow: var(--am-shadow-raised);
  }

  &__lang {
    display: flex;
    justify-content: flex-end;
  }
  &__lang-text {
    color: var(--am-foreground-muted);
    font-size: var(--am-font-sm);
    cursor: pointer;
  }

  &__title-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--am-spacing-sm);
  }
  &__title {
    color: var(--am-foreground-primary);
    font-size: var(--am-font-xl);
    font-weight: 700;
  }
  &__version {
    padding: 2px 8px;
    color: var(--am-foreground-on-accent);
    background: var(--am-accent-primary);
    border-radius: 2px;
    font-size: var(--am-font-xs);
    font-weight: 600;
  }

  &__form {
    .el-form-item {
      margin-bottom: 20px;
    }
  }

  &__field {
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
  &__field-icon {
    color: var(--am-foreground-muted);
    flex: 0 0 auto;
  }

  &__btn {
    width: 100%;
    height: 46px;
    font-size: 15px;
    font-weight: 600;
    border-radius: var(--am-radius-sm);
  }
}

@media (max-width: 520px) {
  .am-login {
    &__main {
      padding: var(--am-spacing-md);
    }

    &__card {
      width: min(400px, 100%);
      padding: var(--am-spacing-lg);
    }

    &__title {
      font-size: var(--am-font-lg);
    }
  }
}
</style>
