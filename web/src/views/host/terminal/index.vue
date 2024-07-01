<template>
    <div class="am-host-header">
        <span @click="$router.push('/host/monitor')">监控</span>
        <span @click="$router.push('/host/file')">文件</span>
        <span @click="$router.push('/host/terminal')">终端</span>
        <span @click="$router.push('/host/settings')">设置</span>
    </div>
    <div class="am-host-operator">
        <el-card shadow="never">
            <el-row :gutter="1">
                <el-col :span="4">
                    <span>主机: </span>
                    <el-input type="text" v-model="host" style="width: 140px" placeholder="192.168.1.1" />
                </el-col>
                <el-col :span="4">
                    <span>端口: </span>
                    <el-input type="number" v-model="port" style="width: 140px" placeholder="22" />
                </el-col>
                <el-col :span="4">
                    <span>用户: </span>
                    <el-input type="text" v-model="username" style="width: 140px" placeholder="root" />
                </el-col>
                <el-col :span="4">
                    <span>密码: </span>
                    <el-input type="password" v-model="password" style="width: 140px" placeholder="123456" />
                </el-col>
                <el-col :span="4">
                    <el-button type="primary" plain @click="createConnection">连接</el-button>
                </el-col>
            </el-row>
        </el-card>
    </div>
    <el-card shadow="never" class="am-host-terminal">
        <el-empty description="暂无终端连接" />
    </el-card>
</template>
<script setup lang="ts">
import { warning } from '@/components/Message/message'

const host = ref('')
const port = ref('')
const username = ref('')
const password = ref('')

const createConnection = () => {
    warning('该功能尚未实现')
}
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
    margin-bottom: 4px;
    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 16px 0 16px;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-start;
            .el-row {
                width: 100%;
                .el-col {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: flex-start;

                    font-size: 14px;
                    color: #575758;
                    span {
                        margin-right: 4px;
                    }
                }
            }
        }
    }
}

@include b(host-terminal) {
    width: 100%;
    height: calc(100vh - 146px);
    overflow-y: auto;
}
</style>
