<template>
    <div class="am-system-title">
        <span>系统设置</span>
    </div>
    <el-tabs type="border-card">
        <el-tab-pane label="DNS" class="am-settings-image"> DNS </el-tab-pane>
        <el-tab-pane label="时区" class="am-settings-image"> 系统时区 </el-tab-pane>
        <el-tab-pane label="NTP" class="am-settings-image"> NTP 服务器 </el-tab-pane>
    </el-tabs>
</template>
<script setup lang="ts"></script>
<style scoped lang="scss">
@include b(system-title) {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
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
            content: '';
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
</style>
