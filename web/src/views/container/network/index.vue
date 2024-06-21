<template>
    <div class="am-network-title">
        <span @click="$router.push('/container/container')">容器</span>
        <span @click="$router.push('/container/image')">镜像</span>
        <span @click="$router.push('/container/network')">网络</span>
    </div>
    <div class="am-network-operator">
        <el-card>
            <el-button type="primary">清理网络</el-button>
        </el-card>
    </div>
    <div>Network</div>
</template>
<script setup lang="ts"></script>
<style scoped lang="scss">
@include b(network-title) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;
    span {
        display: flex;
        align-items: center;
        font-size: 16px;
        font-weight: 600;
        color: #424244;
        margin-left: 16px;
        margin-right: 16px;
        cursor: pointer;
        &:nth-child(3) {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 148px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
@include b(network-operator) {
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
