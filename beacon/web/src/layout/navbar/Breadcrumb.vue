<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const router = useRouter()

const breadcrumbList = computed(() => {
  return router.currentRoute.value.matched.map(item => item)
})

function toHomepage() {
  router.push('/overview')
}
const { t } = useI18n()
</script>

<template>
    <div class="am-breadcrumb">
        <svg-icon icon-class="menu-home" @click="toHomepage" />
        <el-breadcrumb separator="/">
            <el-breadcrumb-item v-for="(item, index) in breadcrumbList" :key="index" :to="{ path: item.path }">
                <span v-if="typeof item.meta?.title === 'string'">{{ t(item.meta?.title) }}</span>
                <span v-else>{{ item.meta?.title }}</span>
            </el-breadcrumb-item>
        </el-breadcrumb>
    </div>
</template>

<style scoped lang="scss">
@include b(breadcrumb) {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
}
</style>
