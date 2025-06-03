<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import type { RouteRecordRaw } from 'vue-router';

defineProps<{
    items: RouteRecordRaw[]
}>()

const { t } = useI18n()

const router = useRouter()
function handleClickMenu(item: RouteRecordRaw) {
    router.replace(item.path)
}
</script>

<template>
    <template v-for="item in items" :key="item.name">
        <el-sub-menu v-if="item.children?.length" :index="item.path">
            <template #title>
                <svg-icon v-if="item.meta && item.meta.icon" :icon-class="item.meta.icon as string" size="1.2rem" />
                <span class="am-sle">{{ t(item.meta?.title as string) }}</span>
            </template>
            <menuitem :items="item.children" />
        </el-sub-menu>
        <el-menu-item v-else :index="item.path" @click="handleClickMenu(item)">
            <svg-icon v-if="item.meta && item.meta.icon" :icon-class="item.meta.icon as string" size="1.2rem" />
            <template #title>
                <span class="am-sle">{{ t(item.meta?.title as string) }}</span>
            </template>
        </el-menu-item>
    </template>
</template>

<style scoped lang="scss">
/* 文字单行省略号 */
@include b(sle) {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-left: 8px;
}
</style>
