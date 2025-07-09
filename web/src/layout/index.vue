<script setup lang="ts">
import Content from '@/layout/content/index.vue';
import Navbar from '@/layout/navbar/index.vue';
import Sidebar from '@/layout/sidebar/index.vue';
</script>

<template>
    <el-container class="am-layout">
        <div class="am-layout-sidebar">
            <Sidebar />
        </div>
        <div class="am-layout-main">
            <Navbar />
            <Content />
        </div>
    </el-container>
</template>

<style scoped lang="scss">
@include b(layout) {
    width: 100%;
    height: 100%;
    background-color: var(--el-bg-color);
    overflow: hidden;
}

@include b(layout-sidebar) {
    display: block;
    position: relative;
    height: 100%;
    transition: 0.5s;
    border-radius: 16px;
    border-right: var(--el-aside-border-color) 1px solid;
    background-color: var(--el-menu-bg-color);
}

@include b(layout-main) {
    display: block;
    position: relative;
    width: 100%;
    height: 100%;
    margin-left: 12px;
    border-radius: 4px;
    background-color: var(--el-bg-color);
    overflow-y: scroll;
}
</style>
