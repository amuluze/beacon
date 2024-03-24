<template>
    <div class="am-menuitem-container" :class="{ 'am-hide-sidebar': hideSidebar }">
        <el-sub-menu v-if="item.children?.length" class="grid" :index="item.path">
            <template #title>
                <svg-icon v-if="item.meta && item.meta.icon" :icon-class="item.meta.icon" />
                <span v-show="!store.app.isCollapse">{{ item.meta?.title }}</span>
            </template>
            <menuitem v-for="i in item.children" :key="i.name" :item="i"></menuitem>
        </el-sub-menu>
        <el-menu-item v-else :index="item.path" @click="handleClickMenu(item)">
            <svg-icon v-if="item.meta && item.meta.icon" :icon-class="item.meta.icon" />
            <template #title>{{ item.meta?.title }}</template>
        </el-menu-item>
    </div>
</template>
<script setup lang="ts">
import useStore from '@/store';
import { RouteRecordRaw } from 'vue-router';

const store = useStore();
const hideSidebar = computed(() => store.app.isCollapse);

const router = useRouter();
defineProps<{
    item: RouteRecordRaw;
}>();
function handleClickMenu(item: RouteRecordRaw) {
    router.push(item.path);
}
</script>

<style scoped lang="scss">
@include b(menuitem-container) {
    .svg-icon {
        margin-right: 8px;
    }
}

@include b(hide-sidebar) {
    .el-sub-menu {
        overflow: hidden;
    }

    .svg-icon {
        margin-right: 0;
    }
}
</style>
