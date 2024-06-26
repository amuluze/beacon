<template>
    <div class="am-system-title">
        <span>系统设置</span>
    </div>
    <el-row :gutter="4" class="am-system-settings">
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

@include b(system-settings) {
    @include e(content) {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;

        .svg-icon {
            color: #2f7bff;
        }
    }
}
</style>
