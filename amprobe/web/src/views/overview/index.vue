<script setup lang="ts">
import type { HostInfo } from '@/interface/host'
import type { DockerInfo } from '@/interface/container.ts'
import type { EChartsOption } from '@/components/Echarts/echarts.ts'

import { queryContainers, queryDockerInfo, queryImages } from '@/api/container'
import { queryCPUInfo, queryDiskInfo, queryHostInfo, queryMemInfo } from '@/api/host'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import DataStaleTag from '@/components/Agent/DataStaleTag.vue'
import { cpuOption, diskOption, memOption } from '@/config/echarts.ts'

import { useI18n } from 'vue-i18n'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { set } from 'lodash-es'
import type { Pagination } from '@/interface/pagination.ts'

const loading = ref(true)
const { t } = useI18n()
const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const initialized = ref(false)

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
const cpuStale = ref(false)
const cpuTimestamp = ref(0)
async function renderCPU() {
  const { data } = await queryCPUInfo()
  cpuStale.value = Boolean(data.stale)
  cpuTimestamp.value = data.timestamp
  set(cpuOption, 'title.text', 'CPU')
  set(cpuOption, 'series[0].data', [Math.round(data.percent) / 100])
}

const memOptionData: EChartsOption = reactive<EChartsOption>(memOption) as EChartsOption
const memStale = ref(false)
const memTimestamp = ref(0)
async function renderMem() {
  const { data } = await queryMemInfo()
  memStale.value = Boolean(data.stale)
  memTimestamp.value = data.timestamp
  set(memOption, 'title.text', 'Mem')
  set(memOption, 'series[0].data', [Math.round(data.percent) / 100])
}

const diskOptionData: EChartsOption = reactive<EChartsOption>(diskOption) as EChartsOption
const diskStale = ref(false)
const diskTimestamp = ref(0)
async function renderDisk() {
  const { data } = await queryDiskInfo()
  const staleDisk = data.info?.find(item => item.stale)
  diskStale.value = Boolean(staleDisk)
  diskTimestamp.value = staleDisk?.timestamp || data.info?.[0]?.timestamp || 0
  set(diskOption, 'title.text', 'Disk')
  set(diskOption, 'series[0].data', [data.info?.[0] ? Math.round(data.info[0].percent) / 100 : 0])
}

async function refreshAgents() {
  await loadAgents()
  if (selectedAgentID.value) {
    await refreshOverview()
  }
}

async function refreshOverview() {
  loading.value = true
  try {
    const agentID = await ensureSelectedAgent()
    if (!agentID)
      return
    await getHostInfo()
    await getDockerInfo()
    await statisticContainer()
    await statisticImage()
    await renderCPU()
    await renderMem()
    await renderDisk()
  }
  finally {
    loading.value = false
  }
}

onMounted(async () => {
  await refreshOverview()
  initialized.value = true
})

watch(selectedAgentID, () => {
  if (initialized.value) {
    refreshOverview()
  }
})
</script>

<template>
    <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
    <el-row v-else class="am-overview" :gutter="8" justify="space-between">
        <el-col :span="16">
            <el-row :gutter="8">
                <el-col :xl="12" :lg="12" :md="12" :sm="12" :xs="24">
                    <el-card shadow="never" class="am-overview__card">
                        <el-skeleton :loading="loading" animated :rows="2">
                            <template #default>
                                <div class="am-panel">
                                    <div class="am-description">
                                        <div class="am-description__text">
                                            {{ t('content.containerNumber') }}
                                        </div>
                                        <div class="am-description__number">
                                            {{ containerCount }}
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </el-skeleton>
                    </el-card>
                </el-col>
                <el-col :xl="12" :lg="12" :md="12" :sm="12" :xs="24">
                    <el-card shadow="never" class="am-overview__card">
                        <el-skeleton :loading="loading" animated :rows="2">
                            <template #default>
                                <div class="am-panel">
                                    <div class="am-description">
                                        <div class="am-description__text">
                                            {{ t('content.imageNumber') }}
                                        </div>
                                        <div class="am-description__number">
                                            {{ imageCount }}
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </el-skeleton>
                    </el-card>
                </el-col>
            </el-row>
            <el-row :gutter="8">
                <el-col :xl="8" :lg="8" :md="8" :sm="8" :xs="24">
                    <div class="am-chart">
                        <el-card shadow="never">
                            <el-skeleton :loading="loading" animated>
                                <DataStaleTag :stale="cpuStale" :timestamp="cpuTimestamp" />
                                <Echarts :option="cpuOptionData" />
                            </el-skeleton>
                        </el-card>
                    </div>
                </el-col>
                <el-col :xl="8" :lg="8" :md="8" :sm="8" :xs="24">
                    <div class="am-chart">
                        <el-card shadow="never">
                            <el-skeleton :loading="loading" animated>
                                <DataStaleTag :stale="memStale" :timestamp="memTimestamp" />
                                <Echarts :option="memOptionData" />
                            </el-skeleton>
                        </el-card>
                    </div>
                </el-col>
                <el-col :xl="8" :lg="8" :md="8" :sm="8" :xs="24">
                    <div class="am-chart">
                        <el-card shadow="never">
                            <el-skeleton :loading="loading" animated>
                                <DataStaleTag :stale="diskStale" :timestamp="diskTimestamp" />
                                <Echarts :option="diskOptionData" />
                            </el-skeleton>
                        </el-card>
                    </div>
                </el-col>
            </el-row>
        </el-col>
        <el-col :span="8">
            <content-wrap :style="{ 'margin-bottom': '8px', 'height': '250px' }" :title="t('content.hostInfo')" :message="t('content.hostInfo')">
                <el-skeleton :loading="loading" animated>
                    <DataStaleTag :stale="hostInfo?.stale" :timestamp="hostInfo?.timestamp" />
                    <p>
                        {{ t('content.hostName') }}：<el-tag>{{ hostInfo?.hostname }}</el-tag>
                    </p>
                    <p>
                        {{ t('content.upTime') }}：<el-tag>{{ hostInfo?.uptime }}</el-tag>
                    </p>
                    <p>
                        {{ t('content.releaseVersion') }}：<el-tag>{{ hostInfo?.platform }}-{{ hostInfo?.platform_version }}</el-tag>
                    </p>
                    <p>
                        {{ t('content.kernelVersion') }}：<el-tag>{{ hostInfo?.kernel_version }}</el-tag>
                    </p>
                    <p>
                        {{ t('content.os') }}：<el-tag>{{ hostInfo?.os }}/{{ hostInfo?.kernel_arch }}</el-tag>
                    </p>
                </el-skeleton>
            </content-wrap>
            <content-wrap :style="{ 'margin-bottom': '8px', 'height': '200px' }" :title="t('content.dockerInfo')" :message="t('content.dockerInfo')">
                <el-skeleton :loading="loading" animated>
                    <p>
                        {{ t('content.dockerVersion') }}：<el-tag>{{ dockerInfo?.docker_version }}</el-tag>
                    </p>
                    <p>
                        {{ t('content.apiVersion') }}： <el-tag>{{ dockerInfo?.min_api_version }}-{{ dockerInfo?.api_version }}</el-tag>
                    </p>
                    <p>
                        {{ t('content.os') }}： <el-tag>{{ dockerInfo?.os }}/{{ dockerInfo?.arch }}</el-tag>
                    </p>
                </el-skeleton>
            </content-wrap>
        </el-col>
    </el-row>
</template>

<style scoped lang="scss">
@include b(overview) {
  overflow: auto;
  @include e(card) {
    margin-bottom: 8px;
    height: 140px;
    //background-color: var(--el-menu-active-bg-color) !important;
  }
}

@include b(panel) {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

@include b(chart) {
  height: 310px;
  margin-bottom: 8px;
  .el-card {
    height: 100%;
    :deep(.el-card__body) {
      height: 100%;
    }
  }
}

@include b(description) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;

  @include e(text) {
    font-size: 14px;
    font-weight: normal;
  }

  @include e(number) {
    font-size: 28px;
  }
}

h4 {
  margin: 16px 0;
}
p {
  font-size: 14px;
  margin: 8px 0;
}
</style>
