<script setup lang="ts">
import { getUserInfo, login } from '@/api/auth'
import useStore from '@/store'
import { Lock, User } from '@element-plus/icons-vue'

import type { LoginForm } from '@/interface/auth'
import { useI18n } from 'vue-i18n'

const loginForm = reactive<LoginForm>({
  username: 'amprobe',
  password: '123456',
})

const loading = ref(false)

const loginFormRules = {
  username: [{ required: true, trigger: 'blur' }],
  password: [{ required: true, trigger: 'blur', validator: passwordValidator }],
}

/**
 * 密码校验器
 */
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

/**
 * 登录
 */
const store = useStore()
const router = useRouter()
async function handleLogin() {
  try {
    loading.value = true
    const { data } = await login({ ...loginForm })
    store.user.setToken(data.access_token, data.refresh_token)
    const userInfo = await getUserInfo()
    store.user.setUserInfo(userInfo.data.username, userInfo.data.status, userInfo.data.is_admin)

    // 跳转首页
    await router.replace('/')
  }
  catch (error) {
    console.log(error)
    if (error instanceof Error)
      ElMessage.error(error.message)
  }
  finally {
    loading.value = false
  }
}

const languageList = [
  { label: '简体中文', value: 'zh' },
  { label: 'English', value: 'en' },
]

const i18n = useI18n()
const language = computed(() => store.app.language)
const t = i18n.t
function changeLanguage(lang: string) {
  i18n.locale.value = lang
  store.app.setLanguage(lang)
  router.replace('/login')
}
</script>

<template>
    <div class="am-login">
        <div class="am-login-left">
            <div class="am-login-left__content">
                <div class="am-login-left__title">
                    <img class="am-login-left__img" src="@/assets/images/amprobe.png" alt="" />
                    <h1>Amprobe</h1>
                </div>
                <h3 class="am-login-left__subtitle">{{ t('about.info') }}</h3>
                <div class="am-login-left__item">
                    <span class="feature-item">· {{ t('about.infoFirst') }}</span>
                    <span class="feature-item">· {{ t('about.infoSecond') }}</span>
                    <span class="feature-item">· {{ t('about.infoThird') }}</span>
                </div>
            </div>
        </div>

        <div class="am-login-right">
            <div class="am-login-right__form">
                <div class="am-login-right__language">
                    <el-dropdown class="mr-4" trigger="click" @command="changeLanguage">
                        <div class="language-switch">
                            <svg-icon size="1.2rem" icon-class="translate" />
                            <span>{{ language === 'zh' ? '简体中文' : 'English' }}</span>
                        </div>
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
                <el-form :model="loginForm" :rules="loginFormRules" class="login-form">
                    <div class="title">
                        <span class="login">{{ t('login.login') }}</span>
                        <el-tag size="small" effect="plain" class="version-tag">
                            <span class="version">v2.0.0</span>
                        </el-tag>
                    </div>

                    <el-form-item prop="username">
                        <el-input
                            v-model="loginForm.username"
                            size="large"
                            class="username-input"
                            placeholder="请输入用户名"
                            name="username"
                            autocomplete="on"
                            :prefix-icon="User"
                        />
                    </el-form-item>
                    <el-form-item prop="password">
                        <el-input
                            v-model="loginForm.password"
                            size="large"
                            type="password"
                            class="password-input"
                            placeholder="请输入密码"
                            name="password"
                            :prefix-icon="Lock"
                            :show-password="true"
                        />
                    </el-form-item>
                    <el-button
                        class="login-btn"
                        size="large"
                        type="primary"
                        :loading="loading"
                        @click.prevent="handleLogin"
                    >
                        {{ t('login.login') }}
                    </el-button>
                </el-form>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
@include b(login) {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(64, 158, 255, 0.1) 0%, transparent 60%);
    animation: rotate 30s linear infinite;
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@include b(login-left) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 50%;
  height: 100vh;
  background: linear-gradient(135deg, #22325b 0%, #1a2540 100%);
  color: #fff;
  position: relative;
  overflow: hidden;

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('@/assets/images/login_bg.svg') no-repeat center;
    background-size: cover;
    opacity: 0.1;
    animation: float 6s ease-in-out infinite;
  }

  @include e(content) {
    position: relative;
    z-index: 1;
    padding: 40px;
    text-align: center;
  }

  @include e(title) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    margin-bottom: 30px;
    animation: fadeInDown 1s ease-out;
  }

  @include e(img) {
    height: 48px;
    width: 48px;
    margin-right: 12px;
    filter: drop-shadow(0 0 8px rgba(255, 255, 255, 0.3));
  }

  @include e(subtitle) {
    font-size: 1.5rem;
    margin-bottom: 2rem;
    opacity: 0.9;
    animation: fadeInUp 1s ease-out 0.3s both;
  }

  @include e(item) {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: center;
    gap: 1rem;
    animation: fadeInUp 1s ease-out 0.6s both;
  }
}

.feature-item {
  font-size: 1.1rem;
  opacity: 0.8;
  transition: all 0.3s ease;

  &:hover {
    opacity: 1;
    transform: translateX(10px);
  }
}

@include b(login-right) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 50%;
  height: 100vh;
  position: relative;

  @include e(form) {
    width: 420px;
    padding: 50px 40px;
    background-color: rgba(255, 255, 255, 0.95);
    border-radius: 16px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(10px);
    animation: fadeIn 1s ease-out;
    position: relative;
    z-index: 1;
  }
}

.language-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background-color: var(--el-color-primary-light-9);
  }

  span {
    font-size: 14px;
    color: var(--el-text-color-regular);
  }
}

.login-form {
  .title {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 12px;
    margin-bottom: 2rem;
    text-align: center;

    .login {
      font-size: 28px;
      font-weight: 600;
      background: linear-gradient(45deg, var(--el-color-primary), var(--el-color-primary-light-3));
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }

  .version-tag {
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 12px;
  }
}

:deep(.el-input) {
  --el-input-hover-border-color: var(--el-color-primary);

  .el-input__wrapper {
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 2px 12px 0 rgba(64, 158, 255, 0.1);
    }

    &.is-focus {
      box-shadow: 0 0 0 1px var(--el-color-primary) inset;
    }
  }
}

.login-btn {
  width: 100%;
  margin-top: 24px;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  background: linear-gradient(45deg, var(--el-color-primary), var(--el-color-primary-light-3));
  border: none;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  }

  &:active {
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
}
</style>
