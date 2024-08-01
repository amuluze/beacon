<template>
    <div class="am-layout-container">
        <el-container>
            <div class="am-layout-sidebar">
                <sidebar />
            </div>
            <div class="am-layout-content">
                <content />
            </div>
        </el-container>
    </div>
</template>

<script setup lang="ts">
import Content from '@/layout/content/index.vue'
import Sidebar from '@/layout/sidebar/index.vue'
</script>

<style scoped lang="scss">
@include b(layout-container) {
    @include bfc;
    display: flex;
    height: 100%;
    width: 100%;
    background-color: #f4f4f4;
    overflow: hidden;
}
@include b(layout-sidebar) {
    display: block;
    position: relative;
    height: 100%;
    transition: 0.5s;
    background-color: #e9effd;
}

@include b(layout-content) {
    display: block;
    position: relative;
    margin: 16px 16px;
    width: 100%;
    height: calc(100vh - 32px);
    overflow-y: scroll;
}
</style>
