<template>
    <div class="am-overview-container">
        <div class="am-overview-container__title">
            <span>总览</span>
        </div>
        <el-row :gutter="8">
            <el-col :span="16">
                <el-card class="am-overview-overview">
                    <h4>概览</h4>
                    <el-row :gutter="4">
                        <el-col :span="12">
                            <el-statistic title="容器数量" :value="containerCount" />
                        </el-col>
                        <el-col :span="12">
                            <el-statistic title="镜像数量" :value="imageCount" />
                        </el-col>
                    </el-row>
                </el-card>
                <el-card class="am-overview-status">
                    <h4>状态</h4>
                    <el-row :gutter="4">
                        <el-col :span="12"> <echarts :option="cpuGaugeOption" height="200px"></echarts> </el-col>
                        <el-col :span="12"> <echarts :option="memGaugeOption" height="200px"></echarts> </el-col>
                    </el-row>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card class="am-overview-host">
                    <h4>主机信息</h4>
                    <p>
                        主机名称：<el-tag>{{ HostInfo?.hostname }}</el-tag>
                    </p>
                    <p>
                        启动时间：<el-tag>{{ HostInfo?.uptime }}</el-tag>
                    </p>
                    <p>
                        发行版本：<el-tag>{{ HostInfo?.platform }}-{{ HostInfo?.platform_version }}</el-tag>
                    </p>
                    <p>
                        内核版本：<el-tag>{{ HostInfo?.kernel_version }}</el-tag>
                    </p>
                    <p>
                        系统类型：<el-tag>{{ HostInfo?.os }}/{{ HostInfo?.kernel_arch }}</el-tag>
                    </p>
                    <h4>Docker信息</h4>
                    <p>
                        Docker 版本：<el-tag>{{ DockerInfo?.docker_version }}</el-tag>
                    </p>
                    <p>
                        API 版本： <el-tag>{{ DockerInfo?.min_api_version }}-{{ DockerInfo?.api_version }}</el-tag>
                    </p>
                    <p>
                        系统类型： <el-tag>{{ DockerInfo?.os }}/{{ DockerInfo?.arch }}</el-tag>
                    </p>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { queryContainers, queryDockerInfo, queryImages } from '@/api/container'
import { queryCPUInfo, queryHostInfo, queryMemInfo } from '@/api/host'
import { EChartsOption } from '@/components/Echarts/echarts'
import { cpuGaugeOptions, memGaugeOptions } from '@/components/Echarts/gauge'
import { set } from 'lodash-es'

onMounted(() => {
    getHostInfo(), getDockerInfo(), renderCPUGauge(), renderMemGauge(), statisticContainer(), statisticImage()
})

const containerCount = ref(0)
const imageCount = ref(0)

const statisticContainer = async () => {
    const { data } = await queryContainers({ page: 1, size: 10 })
    containerCount.value = data.total
}

const statisticImage = async () => {
    const { data } = await queryImages({ page: 1, size: 10 })
    imageCount.value = data.total
}

const HostInfo = ref()
const getHostInfo = async () => {
    const { data } = await queryHostInfo()
    HostInfo.value = {
        hostname: data.hostname,
        uptime: data.uptime,
        os: data.os,
        platform: data.platform,
        platform_version: data.platform_version,
        kernel_version: data.kernel_version,
        kernel_arch: data.kernel_arch
    }
}

const DockerInfo = ref()
const getDockerInfo = async () => {
    const { data } = await queryDockerInfo()
    console.log(data)
    DockerInfo.value = {
        docker_version: data.docker_version,
        api_version: data.api_version,
        min_api_version: data.min_api_version,
        git_commit: data.git_commit,
        go_version: data.go_version,
        os: data.os,
        arch: data.arch
    }
    console.log(DockerInfo.value)
}

const cpuGaugeData = [
    {
        value: 20,
        name: 'CPU',
        title: {
            offsetCenter: ['0%', '-15%']
        },
        detail: {
            valueAnimation: true,
            offsetCenter: ['0%', '15%']
        }
    }
]

const cpuGaugeOption = reactive<EChartsOption>(cpuGaugeOptions)
const renderCPUGauge = async () => {
    const { data } = await queryCPUInfo()
    // 保留小数点后两位
    cpuGaugeData[0].value = Math.round(data.percent * 100) / 100
    set(cpuGaugeOption, 'series[0].data', cpuGaugeData)
    console.log('gauge', cpuGaugeOption)
}

const memGaugeData = [
    {
        value: 20,
        name: 'Memory',
        title: {
            offsetCenter: ['0%', '-15%']
        },
        detail: {
            valueAnimation: true,
            offsetCenter: ['0%', '15%']
        }
    }
]

const memGaugeOption = reactive<EChartsOption>(memGaugeOptions)
const renderMemGauge = async () => {
    const { data } = await queryMemInfo()
    // 保留小数点后两位
    memGaugeData[0].value = Math.round(data.percent * 100) / 100
    set(memGaugeOption, 'series[0].data', memGaugeData)
    console.log('gauge', memGaugeOption)
}
</script>

<style scoped lang="scss">
@include b(overview-container) {
    font-size: 14px;

    @include e(title) {
        display: flex;
        align-items: center;
        justify-content: space-between;
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
            color: #105eeb;
            margin-left: 16px;
            &::before {
                content: ' ';
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
}

@include b(overview-overview) {
    margin-bottom: 8px;
    height: 160px;
    .el-row {
        text-align: center;
    }
}

@include b(overview-status) {
    height: 300px;
}

@include b(overview-host) {
    height: 468px;
    overflow-y: auto;
}
</style>
