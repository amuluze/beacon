<template>
    <div class="am-host-header">
        <span @click="$router.push('/host/monitor')">监控</span>
        <span @click="$router.push('/host/file')">文件</span>
        <span @click="$router.push('/host/terminal')">终端</span>
        <span @click="$router.push('/host/settings')">设置</span>
    </div>
    <div class="am-host-operator">
        <el-card shadow="never">
            <el-button type="warning" @click="reboot">重启</el-button>
            <el-button type="danger" @click="shutdown">关机</el-button>
        </el-card>
    </div>
    <el-row :gutter="4" class="am-host-settings">
        <!-- <el-col :span="8">
            <el-card shadow="never">
                <h4>DNS 设置</h4>
                <el-divider />
                <div class="am-system-settings__content">
                    <span>DNS 设置：</span>
                    <span></span>
                    <svg-icon icon-class="edit" />
                </div>
            </el-card>
        </el-col> -->
        <el-col :span="12">
            <el-card shadow="never">
                <h4>系统时区设置</h4>
                <el-divider />
                <div class="am-system-settings__content">
                    <span>系统时区设置：</span>
                    <span></span>
                    <svg-icon icon-class="edit" />
                </div>
            </el-card>
        </el-col>
        <el-col :span="12">
            <el-card shadow="never">
                <h4>系统时间设置</h4>
                <el-divider />
                <div class="am-system-settings__content">
                    <span>系统时间设置：</span>
                    <span></span>
                    <svg-icon icon-class="edit" />
                </div>
            </el-card>
        </el-col>
    </el-row>
</template>
<script setup lang="ts">
import { reboot, shutdown } from '@/api/system'
</script>
<style scoped lang="scss">
@include b(host-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
    // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;
    span {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        font-size: 16px;
        font-weight: 600;
        margin-left: 16px;
        margin-right: 16px;
        color: #424244;
        cursor: pointer;
        &:first-child {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 212px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 8px 0 0;
            display: flex;
            align-items: center;
            justify-content: flex-end;
        }
    }
}

@include b(host-operator) {
    height: 48px;
    width: 100%;

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 16px 0 16px;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-end;
        }
    }

    border-radius: 4px;
    margin-bottom: 8px;
}

@include b(host-settings) {
    @include e(content) {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
    }
}
</style>
