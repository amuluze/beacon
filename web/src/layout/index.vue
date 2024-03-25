<template>
    <div class="am-layout-container">
        <el-container>
            <div class="am-layout-sidebar">
                <sidebar />
            </div>
            <div class="am-layout-content">
                <el-header>
                    <navbar />
                </el-header>
                <el-main>
                    <content />
                </el-main>
            </div>
        </el-container>
    </div>
</template>

<script setup lang="ts">
import Navbar from '@/layout/navbar/index.vue'
import Sidebar from '@/layout/sidebar/index.vue'
import Content from '@/layout/content/index.vue'
</script>

<style scoped lang="scss">
@include b(layout-container) {
    height: 100%;
    width: 100%;
    background-color: #f4f4f4;
    overflow: hidden;

    .el-container {
        height: 100%;
        width: 100%;
    }
}

@include b(layout-sidebar) {
    //width: 210px;
    height: 100%;
    background-color: #e9effd;
}

@include b(layout-content) {
    margin: 8px 8px 0 8px;
    width: 100%;
    height: 100%;

    .el-header {
        height: 48px;
    }

    .el-main {
        height: calc(100vh - 48px);
    }
}
</style>
