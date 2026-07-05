<script setup lang="ts">
import { getUserInfo, login } from '@/api/auth'
import useStore from '@/store'

import type { LoginForm } from '@/interface/auth'
import { useI18n } from 'vue-i18n'

const loginForm = reactive<LoginForm>({
  username: 'beacon',
  password: '123456',
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
    await router.replace('/')
  }
  catch (error) {
    console.log(error)
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
    <!-- Left Brand Panel -->
    <div class="am-login-left">
      <div class="am-login-left__inner">
        <div class="am-login-left__logo-row">
          <div class="am-login-left__logo-icon" />
          <span class="am-login-left__app-name">Beacon</span>
        </div>
        <p class="am-login-left__desc">轻量级 Docker 容器管理平台</p>
        <div class="am-login-left__features">
          <span>· 实时监控主机与容器资源</span>
          <span>· 可视化容器生命周期管理</span>
          <span>· 镜像仓库与网络配置</span>
          <span>· 多维度告警与审计日志</span>
        </div>
      </div>
    </div>

    <!-- Right Login Card -->
    <div class="am-login-right">
      <div class="am-login-right__card">
        <div class="am-login-right__lang">
          <el-dropdown trigger="click" @command="changeLanguage">
            <span class="am-login-right__lang-text">{{ language === 'zh' ? '简体中文' : 'English' }} ▼</span>
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

        <div class="am-login-right__title-row">
          <span class="am-login-right__title">登录</span>
          <span class="am-login-right__version">v3.0.0</span>
        </div>

        <el-form :model="loginForm" :rules="loginFormRules" class="am-login-right__form">
          <el-form-item prop="username">
            <div class="am-login-right__field">
              <span class="am-login-right__field-icon">👤</span>
              <el-input v-model="loginForm.username" size="large" placeholder="请输入用户名" />
            </div>
          </el-form-item>
          <el-form-item prop="password">
            <div class="am-login-right__field">
              <span class="am-login-right__field-icon">🔒</span>
              <el-input v-model="loginForm.password" size="large" type="password" placeholder="请输入密码" show-password />
            </div>
          </el-form-item>
          <el-button class="am-login-right__btn" size="large" type="primary" @click.prevent="handleLogin">
            登 录
          </el-button>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.am-login {
  display: flex;
  height: 100vh;
  background: #f5f6fa;
}

// ── Left brand panel ──
.am-login-left {
  width: 580px;
  min-width: 580px;
  background: #22325b;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;

  &__inner {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  &__logo-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  &__logo-icon {
    width: 40px;
    height: 40px;
    border-radius: 8px;
    background: #4f7cff;
  }

  &__app-name {
    font-size: 28px;
    font-weight: 700;
    color: #fff;
    font-family: Inter, sans-serif;
  }

  &__desc {
    font-size: 15px;
    color: rgba(255, 255, 255, 0.7);
    margin: 0;
  }

  &__features {
    display: flex;
    flex-direction: column;
    gap: 14px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 13px;
  }
}

// ── Right card ──
.am-login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;

  &__card {
    width: 400px;
    background: #fff;
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    padding: 32px;
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  &__lang {
    display: flex;
    justify-content: flex-end;
  }
  &__lang-text {
    font-size: 13px;
    color: #999;
    cursor: pointer;
  }

  &__title-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
  }
  &__title {
    font-size: 22px;
    font-weight: 700;
    color: #1a1a2e;
  }
  &__version {
    font-size: 11px;
    font-weight: 600;
    color: #fff;
    background: #4f7cff;
    padding: 2px 8px;
    border-radius: 4px;
  }

  &__form {
    .el-form-item {
      margin-bottom: 20px;
    }
  }

  &__field {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    background: #f5f6fa;
    border-radius: 6px;
    padding: 4px 14px;

    .el-input {
      flex: 1;
    }
    .el-input__wrapper {
      background: transparent;
      box-shadow: none !important;
    }
  }
  &__field-icon {
    font-size: 15px;
  }

  &__btn {
    width: 100%;
    height: 46px;
    font-size: 15px;
    font-weight: 600;
    border-radius: 6px;
  }
}
</style>
