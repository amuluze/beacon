<script setup lang="ts">
import { getUserInfo, login } from '@/api/auth'
import useStore from '@/store'

import type { LoginForm } from '@/interface/auth'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const store = useStore()
const router = useRouter()

const loginFormRef = ref<FormInstance>()

const loginForm = reactive<LoginForm>({
  username: '',
  password: '',
})

// 登录只校验必填：凭据是否有效由服务端裁决，客户端不设长度门槛，避免误拦历史密码。
const loginFormRules = computed(() => ({
  username: [{ required: true, message: t('login.usernamePlaceholder'), trigger: 'blur' }],
  password: [{ required: true, message: t('login.passwordRequired'), trigger: 'blur' }],
}))

const submitting = shallowRef(false)
async function handleLogin() {
  // 防重入：回车既触发 @keyup.enter 又触发 form @submit.prevent，避免重复提交
  if (submitting.value)
    return
  if (!loginFormRef.value)
    return
  submitting.value = true
  try {
    await loginFormRef.value.validate()
    await login({ ...loginForm })
    const userInfo = await getUserInfo()
    store.user.setUserInfo(
      userInfo.data.username,
      userInfo.data.status,
      userInfo.data.is_admin,
      userInfo.data.must_change_password,
    )
    // 每次认证会话都重新确认可用 Agent，避免沿用上一个会话的目标节点。
    store.agent.clear()
    await router.replace(userInfo.data.must_change_password ? '/change-password' : '/monitor')
  }
  catch {
    // 认证失败与网络异常由全局响应拦截器统一提示，这里只负责结束提交态。
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
}
</script>

<template>
    <div class="am-login">
        <aside class="am-login__brand">
            <span class="am-login__orb am-login__orb--primary" aria-hidden="true" />
            <span class="am-login__orb am-login__orb--bright" aria-hidden="true" />

            <div class="am-login__logo">
                <img class="am-login__logo-mark" src="/beacon.svg" alt="" aria-hidden="true" />
                <span class="am-login__logo-text">Beacon</span>
            </div>

            <div class="am-login__pitch">
                <p class="am-login__slogan">
                    {{ $t('login.slogan') }}
                </p>
                <p class="am-login__subline">
                    {{ $t('login.subline') }}
                </p>
            </div>
        </aside>

        <main class="am-login__main">
            <section class="am-login__card">
                <div class="am-login__lang">
                    <el-dropdown trigger="click" @command="changeLanguage">
                        <button class="am-login__lang-chip" type="button" :aria-label="t('login.language')">
                            <span>{{ language === 'zh' ? '简体中文' : 'English' }}</span>
                            <svg-icon icon-class="down" size="12px" />
                        </button>
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

                <h1 class="am-login__title">
                    {{ $t('login.login') }}
                </h1>

                <el-form ref="loginFormRef" :model="loginForm" :rules="loginFormRules" class="am-login__form" @submit.prevent="handleLogin">
                    <el-form-item prop="username">
                        <div class="am-login__field">
                            <svg-icon class="am-login__field-icon" icon-class="user" size="16px" />
                            <el-input
                                v-model="loginForm.username"
                                :placeholder="t('login.usernamePlaceholder')"
                                :aria-label="t('login.usernamePlaceholder')"
                                autocomplete="username"
                                @keyup.enter="handleLogin"
                            />
                        </div>
                    </el-form-item>
                    <el-form-item prop="password">
                        <div class="am-login__field">
                            <svg-icon class="am-login__field-icon" icon-class="lock" size="16px" />
                            <el-input
                                v-model="loginForm.password"
                                type="password"
                                show-password
                                :placeholder="t('login.passwordPlaceholder')"
                                :aria-label="t('login.passwordPlaceholder')"
                                autocomplete="current-password"
                                @keyup.enter="handleLogin"
                            />
                        </div>
                    </el-form-item>
                    <el-button class="am-login__btn" type="primary" native-type="submit" :loading="submitting">
                        {{ $t('login.submit') }}
                    </el-button>
                </el-form>
            </section>
        </main>
    </div>
</template>

<style scoped lang="scss">
// 登录页按 .pen/beacon.ax 的 Login 页实现：左栏品牌渐变 + 装饰圆形，右栏居中卡片。
// 控制台最小支持宽度 1280（DESIGN.md §15），不提供重排稿；窄窗口横向滚动。
.am-login {
  display: flex;
  min-width: 1280px;
  min-height: 100vh;
  color: var(--am-foreground-primary);
  background: var(--am-surface-primary);
}

.am-login__brand {
  position: relative;
  flex: 0 0 50%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  gap: var(--am-spacing-lg);
  padding: 56px;
  overflow: hidden;
  // 品牌渐变深端固定在两档主题的深色范围内，保证白色品牌文案始终可读。
  background: linear-gradient(
    225deg,
    color-mix(in srgb, var(--am-accent-primary) 40%, var(--am-brand-deep)) 0%,
    var(--am-brand-deep) 100%
  );
}

// 稿面两处装饰圆形：上右 accent-primary、下左 accent-primary-bright，越界裁切。
.am-login__orb {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}

.am-login__orb--primary {
  top: -140px;
  right: -130px;
  width: 420px;
  height: 420px;
  background: var(--am-accent-primary);
  opacity: 0.55;
}

.am-login__orb--bright {
  bottom: -40px;
  left: -120px;
  width: 300px;
  height: 300px;
  background: var(--am-accent-primary-bright);
  opacity: 0.32;
}

.am-login__logo {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

.am-login__logo-mark {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.am-login__logo-text {
  font-size: var(--am-font-xl);
  font-weight: 700;
  color: #fff;
}

.am-login__pitch {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  max-width: 480px;
  text-align: center;
}

.am-login__slogan {
  margin: 0;
  font-size: var(--am-font-lg);
  font-weight: 500;
  color: #fff;
}

.am-login__subline {
  margin: 0;
  font-size: var(--am-font-md);
  color: rgb(255 255 255 / 72%);
}

.am-login__main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--am-spacing-xl);
  background: var(--am-surface-primary);
}

.am-login__card {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-lg);
  width: 400px;
  max-width: 100%;
  padding: var(--am-spacing-xl);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-lg);
}

.am-login__lang {
  display: flex;
  justify-content: flex-end;
}

.am-login__lang-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  color: var(--am-foreground-secondary);
  font-family: inherit;
  font-size: var(--am-font-sm);
  font-weight: 500;
  background: var(--am-surface-inset);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
  cursor: pointer;
  transition:
    color 0.2s ease,
    border-color 0.2s ease;

  &:hover {
    color: var(--am-foreground-primary);
    border-color: var(--am-accent-primary);
  }
}

.am-login__title {
  margin: 0;
  text-align: center;
  font-size: var(--am-font-xl);
  font-weight: 700;
  color: var(--am-foreground-primary);
}

.am-login__form {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);

  :deep(.el-form-item) {
    margin-bottom: 0;
  }

  // 校验收紧：错误态描边跟随字段容器，而不是外层 form-item。
  :deep(.el-form-item.is-error) .am-login__field {
    border-color: var(--am-accent-danger);
  }

  :deep(.el-form-item__error) {
    padding-top: 4px;
    font-size: var(--am-font-xs);
  }
}

.am-login__field {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  width: 100%;
  padding: 12px 16px;
  background: var(--am-surface-inset);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;

  &:focus-within {
    border-color: var(--am-accent-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--am-accent-primary) 16%, transparent);
  }

  .el-input {
    flex: 1;
  }

  :deep(.el-input__wrapper) {
    padding: 0;
    background: transparent;
    box-shadow: none !important;
  }

  :deep(.el-input__inner) {
    height: 22px;
    font-size: var(--am-font-md);
    color: var(--am-foreground-primary);
  }

  :deep(.el-input__inner::placeholder) {
    color: var(--am-foreground-placeholder);
  }
}

.am-login__field-icon {
  flex: 0 0 auto;
  color: var(--am-foreground-muted);
}

.am-login__btn {
  width: 100%;
  height: 48px;
  color: var(--am-foreground-on-accent);
  font-size: var(--am-font-md);
  font-weight: 600;
  border-radius: var(--am-radius-sm);
}
</style>
