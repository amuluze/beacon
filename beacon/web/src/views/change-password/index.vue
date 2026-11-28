<script setup lang="ts">
import { usePasswordChange } from '@/hooks/usePasswordChange'
import useStore from '@/store'
import { useI18n } from 'vue-i18n'

const store = useStore()
const { t } = useI18n()
const username = computed(() => store.user.userInfo.name)
const {
  formRef,
  loading,
  model,
  rules,
  submit,
} = usePasswordChange(username.value)
</script>

<template>
    <main class="mandatory-password-change">
        <section class="mandatory-password-change__card">
            <div class="mandatory-password-change__brand">
                <img src="/beacon.svg" alt="" aria-hidden="true">
                <span>Beacon</span>
            </div>
            <div class="mandatory-password-change__heading">
                <h1>{{ t('account.changePassword') }}</h1>
                <p>{{ t('account.passwordPolicyMandatory') }}</p>
            </div>
            <el-form
                ref="formRef"
                :model="model"
                :rules="rules"
                label-position="top"
                class="mandatory-password-change__form"
                @keyup.enter="submit"
            >
                <el-form-item :label="t('account.userName')" prop="username">
                    <el-input v-model="model.username" disabled :placeholder="t('account.inputUserName')" />
                </el-form-item>
                <el-form-item :label="t('account.oldPassword')" prop="old_password">
                    <el-input
                        v-model="model.old_password"
                        type="password"
                        show-password
                        :placeholder="t('account.inputOldPassword')"
                        autocomplete="current-password"
                    />
                </el-form-item>
                <el-form-item :label="t('account.newPassword')" prop="new_password">
                    <el-input
                        v-model="model.new_password"
                        type="password"
                        show-password
                        :placeholder="t('account.inputNewPassword')"
                        autocomplete="new-password"
                    />
                </el-form-item>
                <el-button
                    class="mandatory-password-change__submit"
                    type="primary"
                    :loading="loading"
                    @click="submit"
                >
                    {{ t('account.confirm') }}
                </el-button>
            </el-form>
        </section>
    </main>
</template>

<style scoped lang="scss">
.mandatory-password-change {
  display: grid;
  min-width: 1280px;
  min-height: 100vh;
  place-items: center;
  padding: var(--am-spacing-xl);
  color: var(--am-foreground-primary);
  background: var(--am-surface-primary);
}

.mandatory-password-change__card {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-lg);
  width: 480px;
  max-width: 100%;
  padding: var(--am-spacing-xl);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-lg);
}

.mandatory-password-change__brand {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  color: var(--am-foreground-primary);
  font-size: var(--am-font-lg);
  font-weight: 700;

  img {
    width: 28px;
    height: 28px;
    object-fit: contain;
  }
}

.mandatory-password-change__heading {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);

  h1 {
    margin: 0;
    font-size: var(--am-font-xl);
    font-weight: 700;
  }

  p {
    margin: 0;
    color: var(--am-foreground-muted);
    font-size: var(--am-font-sm);
    line-height: 1.55;
  }
}

.mandatory-password-change__form {
  display: flex;
  flex-direction: column;

  :deep(.el-form-item) {
    margin-bottom: var(--am-spacing-md);
  }

  :deep(.el-form-item__label) {
    padding-bottom: 4px;
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-sm);
  }
}

.mandatory-password-change__submit {
  width: 100%;
  height: 48px;
  margin-top: var(--am-spacing-sm);
  color: var(--am-foreground-on-accent);
  font-size: var(--am-font-md);
  font-weight: 600;
}
</style>
