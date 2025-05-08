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
            <div class="am-login-left__title">
                <img class="am-login-left__img" src="@/assets/images/amprobe.png" alt="" />
                <h1>Amprobe</h1>
            </div>
            <h3>{{ t('about.info') }}</h3>
            <div class="am-login-left__item">
                <span>· {{ t('about.infoFirst') }}</span>
                <span>· {{ t('about.infoSecond') }}</span>
                <span>· {{ t('about.infoThird') }}</span>
            </div>
        </div>

        <div class="am-login-right">
            <div class="am-login-right__form">
                <div class="am-login-right__language">
                    <el-dropdown class="mr-4" trigger="click" @command="changeLanguage">
                        <svg-icon size="1.2rem" icon-class="translate" />
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
                <el-form :model="loginForm" :rules="loginFormRules">
                    <div class="title">
                        <span class="login">{{ t('login.login') }}</span>
                        <el-tag> <span class="version">v1.3.10</span> </el-tag>
                    </div>

                    <el-form-item prop="username">
                        <el-input v-model="loginForm.username" size="large" class="username-input" placeholder="请输入用户名" name="username" autocomplete="on" :prefix-icon="User" />
                    </el-form-item>
                    <el-form-item prop="password">
                        <el-input v-model="loginForm.password" size="large" type="password" class="password-input" placeholder="请输入密码" name="password" :prefix-icon="Lock" :show-password="true" />
                    </el-form-item>
                    <el-button class="btn" size="large" type="primary" @click.prevent="handleLogin">
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
  min-height: 100%;
  background-color: #eee;
  background-image: url('@/assets/images/login_bg.svg');
  background-size: auto;
}

@include b(login-left) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-image: url('@/assets/images/bg-bd839ea3.svg');
  width: 48%;
  height: 100%;
  min-height: 100%;
  background-color: #22325b;
  color: #fff;
  padding: 20px;

  @include e(title) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    margin-bottom: 30px;
  }

  @include e(img) {
    height: 36px;
    width: 36px;
    margin-right: 6px;
  }

  @include e(item) {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: center;
  }
}

@include b(login-right) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 60%;
  height: 100%;

  @include e(form) {
    width: 400px;
    padding: 50px 40px 45px;
    background-color: var(--el-bg-color);
    border-radius: 10px;
    box-shadow: rgb(0 0 0 / 10%) 0 2px 10px 2px;

    .title {
      display: flex;
      flex-direction: row;
      align-items: flex-start;
      justify-content: center;
      padding-bottom: 1rem;
      margin: 0 auto;
      text-align: center;

      .login {
        font-size: 24px;
        font-weight: bold;
        color: #34495e;
        white-space: nowrap;
      }
      .version {
        font-size: 14px;
        font-weight: bold;
        color: #2f7bff;
      }
    }

    .btn {
      width: 100%;
    }
  }
  @include e(language) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-end;
  }
}
</style>
