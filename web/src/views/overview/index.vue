<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { DockerInfo } from '@/interface/container.ts'
import type { HostInfo } from '@/interface/host'

import { queryContainers, queryDockerInfo, queryImages } from '@/api/container'
import { queryCPUInfo, queryDiskInfo, queryHostInfo, queryMemInfo } from '@/api/host'
import { cpuOption, diskOption, memOption } from '@/config/echarts.ts'

import type { Pagination } from '@/interface/pagination.ts'
import { set } from 'lodash-es'
import { useI18n } from 'vue-i18n'

const loading = ref(true)
const { t } = useI18n()

const hostInfo = ref<HostInfo>()
async function getHostInfo() {
  const { data } = await queryHostInfo()
  hostInfo.value = data
}

const dockerInfo = ref<DockerInfo>()
async function getDockerInfo() {
  const { data } = await queryDockerInfo()
  dockerInfo.value = data
}

const containerCount = ref(0)
const imageCount = ref(0)

async function statisticContainer() {
  const params: Pagination = {
    page: 1,
    size: 10,
  }
  const { data } = await queryContainers(params)
  containerCount.value = data.total
}

async function statisticImage() {
  const params: Pagination = {
    page: 1,
    size: 10,
  }
  const { data } = await queryImages(params)
  imageCount.value = data.total
}

const cpuOptionData: EChartsOption = reactive<EChartsOption>(cpuOption) as EChartsOption
async function renderCPU() {
  const { data } = await queryCPUInfo()
  set(cpuOption, 'title.text', 'CPU')
  set(cpuOption, 'series[0].data', [Math.round(data.percent) / 100])
  console.log('cpu option: ', cpuOptionData)
}

const memOptionData: EChartsOption = reactive<EChartsOption>(memOption) as EChartsOption
async function renderMem() {
  const { data } = await queryMemInfo()
  set(memOption, 'title.text', 'Mem')
  set(memOption, 'series[0].data', [Math.round(data.percent) / 100])
  console.log('mem option: ', memOptionData)
}

const diskOptionData: EChartsOption = reactive<EChartsOption>(diskOption) as EChartsOption
async function renderDisk() {
  const { data } = await queryDiskInfo()
  set(diskOption, 'title.text', 'Disk')
  set(diskOption, 'series[0].data', [Math.round(data.info[0].percent) / 100])
  console.log('disk option: ', diskOptionData)
}

onMounted(async () => {
  await getHostInfo()
  await getDockerInfo()
  await statisticContainer()
  await statisticImage()
  await renderCPU()
  await renderMem()
  await renderDisk()
  loading.value = false
})
</script>

<template>
    <el-card shadow="hover" class="am-overview">
        <el-skeleton :loading="loading" animated>
            <div class="overview-container">
                <div class="overview-stats">
                    <div class="stat-item">
                        <div class="stat-content">
                            <div class="stat-label">{{ t('content.containerNumber') }}</div>
                            <div class="stat-value">{{ containerCount }}</div>
                        </div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-content">
                            <div class="stat-label">{{ t('content.imageNumber') }}</div>
                            <div class="stat-value">{{ imageCount }}</div>
                        </div>
                    </div>
                </div>
                <div class="overview-charts">
                    <div class="chart-item">
                        <Echarts :option="cpuOptionData" />
                    </div>
                    <div class="chart-item">
                        <Echarts :option="memOptionData" />
                    </div>
                    <div class="chart-item">
                        <Echarts :option="diskOptionData" />
                    </div>
                </div>
                <div class="overview-info">
                    <div class="info-section">
                        <div class="info-title">{{ t('content.hostInfo') }}</div>
                        <div class="info-content">
                            <div class="info-item">
                                <span class="info-label">{{ t('content.hostName') }}：</span>
                                <el-tag>{{ hostInfo?.hostname }}</el-tag>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('content.upTime') }}：</span>
                                <el-tag>{{ hostInfo?.uptime }}</el-tag>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('content.releaseVersion') }}：</span>
                                <el-tag>{{ hostInfo?.platform }}-{{ hostInfo?.platform_version }}</el-tag>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('content.kernelVersion') }}：</span>
                                <el-tag>{{ hostInfo?.kernel_version }}</el-tag>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('content.os') }}：</span>
                                <el-tag>{{ hostInfo?.os }}/{{ hostInfo?.kernel_arch }}</el-tag>
                            </div>
                        </div>
                    </div>
                    <div class="info-section">
                        <div class="info-title">{{ t('content.dockerInfo') }}</div>
                        <div class="info-content">
                            <div class="info-item">
                                <span class="info-label">{{ t('content.dockerVersion') }}：</span>
                                <el-tag>{{ dockerInfo?.docker_version }}</el-tag>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('content.apiVersion') }}：</span>
                                <el-tag>{{ dockerInfo?.min_api_version }}-{{ dockerInfo?.api_version }}</el-tag>
                            </div>
                            <div class="info-item">
                                <span class="info-label">{{ t('content.os') }}：</span>
                                <el-tag>{{ dockerInfo?.os }}/{{ dockerInfo?.arch }}</el-tag>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </el-skeleton>
    </el-card>
</template>

<style lang="scss" scoped>
.am-overview {
    border-radius: 16px;
    transition: all 0.3s;
    height: calc(100vh - 80px);
    overflow: hidden;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    border: 1px solid var(--el-border-color-light);

    &:hover {
        box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
    }

    :deep(.el-card__body) {
        padding: 16px;
        height: 100%;
        overflow: hidden;
    }
}

.overview-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: 100%;
    overflow-y: auto;
    padding-right: 4px;

    &::-webkit-scrollbar {
        width: 4px;
    }

    &::-webkit-scrollbar-thumb {
        background: var(--el-border-color-darker);
        border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
        background: var(--el-border-color-light);
        border-radius: 2px;
    }
}

.overview-stats {
    display: flex;
    gap: 16px;
    flex-shrink: 0;
}

.stat-item {
    flex: 1;
    background: var(--el-bg-color-page);
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.05);
    border: 1px solid var(--el-border-color-light);
    transition: all 0.3s;

    &:hover {
        box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
    }
}

.stat-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.stat-label {
    font-size: 14px;
    color: var(--el-text-color-secondary);
}

.stat-value {
    font-size: 24px;
    font-weight: 500;
    color: var(--el-text-color-primary);
}

.overview-charts {
    display: flex;
    gap: 16px;
    flex-shrink: 0;
}

.chart-item {
    flex: 1;
    background: var(--el-bg-color-page);
    border-radius: 8px;
    padding: 12px;
    box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.05);
    border: 1px solid var(--el-border-color-light);
    height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;

    &:hover {
        box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
    }

    :deep(.echarts) {
        width: 100%;
        height: 100%;
    }
}

.overview-info {
    display: flex;
    gap: 16px;
    flex-shrink: 0;
}

.info-section {
    flex: 1;
    background: var(--el-bg-color-page);
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.05);
    border: 1px solid var(--el-border-color-light);
    transition: all 0.3s;

    &:hover {
        box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
    }
}

.info-title {
    font-size: 16px;
    font-weight: 500;
    color: var(--el-text-color-primary);
    margin-bottom: 12px;
}

.info-content {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.info-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.info-label {
    color: var(--el-text-color-secondary);
    font-size: 14px;
    min-width: 110px;
}

:deep(.el-tag) {
    border-radius: 4px;
    font-size: 14px;
    padding: 0 8px;
    height: 24px;
    line-height: 24px;
}
</style>
