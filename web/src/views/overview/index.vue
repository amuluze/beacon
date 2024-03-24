<template>
    <div class="am-overview-container">
        <el-row :gutter="8">
            <el-col :span="18">
                <el-card class="am-overview-overview">
                    <h4>概览</h4>
                </el-card>
                <el-card>
                    <h4>状态</h4>
                </el-card>
            </el-col>
            <el-col :span="6">
                <el-card class="am-overview-host">
                    <h4>主机信息</h4>
                    <p>主机名称：{{ HostInfo?.hostname }}</p>
                    <p>启动时间： {{ HostInfo?.uptime }}</p>
                    <p>os：{{ HostInfo?.os }}</p>
                    <p>platform: {{ HostInfo?.platform }}</p>
                    <p>platform version: {{HostInfo?.platform_version }}</p>
                    <p>kernel version: {{HostInfo?.kernel_version }}</p>
                    <p>kernel arch: {{HostInfo?.kernel_arch }}</p>
                </el-card>
                <el-card>
                    <h4>容器信息</h4>
                    <p>主机名称：{{ HostInfo?.hostname }}</p>
                    <p>启动时间： {{ HostInfo?.uptime }}</p>
                    <p>os：{{ HostInfo?.os }}</p>
                    <p>platform: {{ HostInfo?.platform }}</p>
                    <p>platform version: {{HostInfo?.platform_version }}</p>
                    <p>kernel version: {{HostInfo?.kernel_version }}</p>
                    <p>kernel arch: {{HostInfo?.kernel_arch }}</p>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { queryHostInfo } from '@/api/host';

onMounted(() => {
    getHostInfo();
});

type HostInfoType = {
    hostname: string;
    uptime: string;
    os: string;
    platform: string;
    platform_version: string;
    kernel_version: string;
    kernel_arch: string;
};

const HostInfo = ref<HostInfoType>();

const getHostInfo = async () => {
    const { data } = await queryHostInfo();
    HostInfo.value = {
        hostname: data.hostname,
        uptime: data.uptime,
        os: data.os,
        platform: data.platform,
        platform_version: data.platform_version,
        kernel_version: data.kernel_version,
        kernel_arch: data.kernel_arch,
    };
};
</script>

<style scoped lang="scss">
@include b(overview-container) {
    font-size: 14px;
    .el-row {
        margin-bottom: 8px;
    }
}
@include b(overview-overview) {
    margin-bottom: 8px;
}

@include b(overview-host) {
    margin-bottom: 8px;
}
</style>
