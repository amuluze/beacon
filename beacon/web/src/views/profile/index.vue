<script setup lang="ts">
import { getUserInfo } from '@/api/auth'
import { useI18n } from 'vue-i18n'

onMounted(() => {
  getProfile()
})

const profile = ref({
  username: '',
  status: '1',
  is_admin: 1,
})
async function getProfile() {
  const { data } = await getUserInfo()
  profile.value = {
    username: data.username,
    status: String(data.status),
    is_admin: data.is_admin,
  }
}

const { t } = useI18n()
</script>

<template>
    <el-row>
        <el-col :span="24">
            <el-card shadow="never">
                <el-descriptions :title="t('profile.currentAccount')" :column="1">
                    <el-descriptions-item :label="t('profile.username')">
                        <el-tag>{{ profile?.username }}</el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('profile.status')">
                        <el-tag>{{ profile?.status === '1' ? t('profile.enable') : t('profile.disable') }}</el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('profile.role')">
                        <el-tag>{{ profile?.is_admin === 1 ? t('profile.admin') : t('profile.normal') }}</el-tag>
                    </el-descriptions-item>
                </el-descriptions>
            </el-card>
        </el-col>
    </el-row>
</template>

<style scoped lang="scss">
@include b(profile-container) {
  font-size: 14px;

  @include e(title) {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    width: 100%;
    // background-color: var(--el-menu-bg-color);
    // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;
    span {
      display: flex;
      align-items: center;
      font-size: 16px;
      font-weight: 600;
      color: #105eeb;
      margin-left: 16px;
      &::before {
        content: ' ';
        position: absolute;
        left: 20px;
        width: 4px;
        height: 16px;
        text-align: center;
        background-color: #2f7bff;
        border-radius: 2px;
      }
    }
  }
}
@include b(profile-operator) {
  height: 48px;
  width: 100%;
  margin-bottom: 4px;
  .el-card {
    height: 100%;
    :deep(.el-card__body) {
      height: 100% !important;
      padding: 0 0 0 16px;
      display: flex;
      align-items: center;
      justify-content: flex-start;
    }
  }
}
</style>
