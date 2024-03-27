<template>
    <div class="am-host-container">
        <div class="am-host-container__operator">
            <el-card>
                <el-select v-model="timeDensity" placeholder="Select" size="default" style="width: 240px">
                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-card>
        </div>
        <el-row :gutter="4">
            <el-col :span="12">
                <el-card>
                    <echarts :option="cpuOption">
                        <div class="am-host-container__image-title">CPU 总使用率</div>
                        <div class="am-host-container__image-description">百分比： {{ cpuPercent }}</div>
                    </echarts>
                </el-card>
            </el-col>
            <el-col :span="12">
                <el-card>
                    <echarts :option="memOption">
                        <div class="am-host-container__image-title">内存使用率</div>
                        <div class="am-host-container__image-description">
                            总量：{{ memInfo.total }} 使用：{{ memInfo.used }} 百分比： {{ memInfo.percent }}
                        </div>
                    </echarts>
                </el-card>
            </el-col>
        </el-row>
        <el-row :gutter="4">
            <el-col :span="12">
                <el-card>
                    <echarts :option="diskOption">
                        <div class="am-host-container__image-title">磁盘使用率</div>
                        <div
                            v-for="(item, index) in diskInfo"
                            :key="index"
                            class="am-host-container__image-description"
                        >
                            {{ item.device }} 总量：{{ item.total }} 使用：{{ item.used }} 百分比： {{ item.percent }}
                        </div>
                    </echarts>
                </el-card>
            </el-col>
            <el-col :span="12">
                <el-card>
                    <echarts :option="netOption">
                        <div class="am-host-container__image-title">流量曲线图</div>
                        <div class="am-host-container__image-description">
                            {{ netInfo.ethernet }} 接收：{{ netInfo.read }} 发送：{{ netInfo.write }}
                        </div>
                    </echarts>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>
<script setup lang="ts">
import {
    queryCPUInfo,
    queryCPUUsage,
    queryDiskInfo,
    queryDiskUsage,
    queryMemInfo,
    queryMemUsage,
    queryNetworkUsage
} from '@/api/host'
import { EChartsOption } from '@/components/Echarts/echarts.ts'
import { cpuOptions, diskOptions, memOptions, netOptions } from '@/components/Echarts/line.ts'
import { CPUTrendingArgs, DiskIO, DiskTrendingArgs, MemTrendingArgs, NetIO, NetTrendingArgs } from '@/interface/host.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'

// 时间密度下拉框
const timeDensity = ref(600)
const options = [
    {
        value: 600,
        label: '10分钟'
    },
    {
        value: 1800,
        label: '30分钟'
    },
    {
        value: 3600,
        label: '1 小时'
    },
    {
        value: 43200,
        label: '12小时'
    },
    {
        value: 86400,
        label: '24小时'
    }
]

const cpuPercent = ref('0.0%')
const renderCPUPercent = async () => {
    const { data } = await queryCPUInfo()
    cpuPercent.value = data.percent.toFixed(2) + '%'
}

const cpuOption = reactive<EChartsOption>(cpuOptions)
const renderCPU = async () => {
    const param: CPUTrendingArgs = {
        start_time: dayjs().unix() - timeDensity.value,
        end_time: dayjs().unix()
    }
    console.log(param)
    const { data } = await queryCPUUsage(param)
    const cpuData = data.data
    console.log('cpu response:', cpuData)
    // set(cpuOption, 'title', { text: 'CPU使用率' });
    set(
        cpuOption,
        'xAxis.data',
        cpuData.map((item) => dayjs(item.timestamp * 1000).hour() + ':' + dayjs(item.timestamp * 1000).minute())
    )
    set(cpuOption, 'legend.data', ['CPU使用率'])
    set(cpuOption, 'series', [
        {
            name: 'CPU使用率',
            data: cpuData.map((item) => item.value),
            type: 'line',
            smooth: true,
            showSymbol: false
        }
    ])
    console.log('cpu: ', cpuOption)
}

const memInfo = ref({
    percent: '0%',
    total: '0',
    used: '0'
})

const renderMemInfo = async () => {
    const { data } = await queryMemInfo()
    memInfo.value.percent = data.percent.toFixed(2) + '%'
    memInfo.value.total = convertBytesToReadable(data.total)
    memInfo.value.used = convertBytesToReadable(data.used)
}

const memOption = reactive<EChartsOption>(memOptions) as EChartsOption
const renderMem = async () => {
    const param: MemTrendingArgs = {
        start_time: dayjs().unix() - timeDensity.value,
        end_time: dayjs().unix()
    }
    console.log(param)
    const { data } = await queryMemUsage(param)
    const memData = data.data
    console.log('mem response: ', memData)
    // set(memOption, 'title', { text: '内存使用率' });
    set(
        memOption,
        'xAxis.data',
        memData.map((item) => dayjs(item.timestamp * 1000).hour() + ':' + dayjs(item.timestamp * 1000).minute())
    )
    set(memOption, 'legend.data', ['内存使用率'])
    set(memOption, 'series', [
        {
            name: '内存使用率',
            data: memData.map((item) => item.value),
            type: 'line',
            smooth: true,
            showSymbol: false
        }
    ])
    console.log('mem: ', memOption)
}

const diskInfo = ref<
    {
        device: string
        total: string
        used: string
        percent: string
    }[]
>([])

const renderDiskInfo = async () => {
    const { data } = await queryDiskInfo()
    console.log(data.info)
    diskInfo.value = []
    data.info.map((item) => {
        diskInfo.value.push({
            device: item.device,
            total: convertBytesToReadable(item.total),
            used: convertBytesToReadable(item.used),
            percent: item.percent.toFixed(2) + '%'
        })
    })
}

const diskOption = reactive<EChartsOption>(diskOptions) as EChartsOption
const renderDisk = async () => {
    const param: DiskTrendingArgs = {
        start_time: dayjs().unix() - timeDensity.value,
        end_time: dayjs().unix()
    }
    console.log(param)
    const { data } = await queryDiskUsage(param)
    const diskData = data
    console.log('********<<<<>>>>>>*****', diskData, diskData.data, diskData.device)
    console.log('disk response: ', diskData.device, diskData.data)
    // set(diskOption, 'title', { text: '磁盘使用率' });
    set(
        diskOption,
        'xAxis.data',
        diskData.data.map(
            (item: DiskIO) => dayjs(item.timestamp * 1000).hour() + ':' + dayjs(item.timestamp * 1000).minute()
        )
    )
    set(diskOption, 'legend.data', ['Read', 'Write'])
    set(diskOption, 'series', [
        {
            name: 'Read',
            data: diskData.data.map((item: DiskIO) => item.io_read),
            type: 'line',
            smooth: true,
            showSymbol: false
        },
        {
            name: 'Write',
            data: diskData.data.map((item: DiskIO) => item.io_write),
            type: 'line',
            smooth: true,
            showSymbol: false
        }
    ])
    console.log('disk options: ', diskOption)
}

const netInfo = ref({
    ethernet: '',
    read: '',
    write: ''
})

const netOption = reactive<EChartsOption>(netOptions) as EChartsOption
const renderNet = async () => {
    const param: NetTrendingArgs = {
        start_time: dayjs().unix() - timeDensity.value,
        end_time: dayjs().unix()
    }
    console.log(param)
    const { data } = await queryNetworkUsage(param)
    const netData = data
    netInfo.value.ethernet = netData.ethernet
    netInfo.value.read = convertBytesToReadable(netData.data[netData.data.length - 1].bytes_recv)
    netInfo.value.write = convertBytesToReadable(netData.data[netData.data.length - 1].bytes_sent)
    console.log('net response: ', netData.ethernet, netData.data)
    // set(netOption, 'title', { text: '流量曲线图' });
    set(
        netOption,
        'xAxis.data',
        netData.data.map(
            (item: NetIO) => dayjs(item.timestamp * 1000).hour() + ':' + dayjs(item.timestamp * 1000).minute()
        )
    )
    set(netOption, 'legend.data', ['Receive', 'Send'])
    set(netOption, 'series', [
        {
            name: 'Receive',
            data: netData.data.map((item: NetIO) => item.bytes_recv),
            type: 'line',
            smooth: true,
            showSymbol: false
        },
        {
            name: 'Send',
            data: netData.data.map((item: NetIO) => item.bytes_sent),
            type: 'line',
            smooth: true,
            showSymbol: false
        }
    ])
    console.log('net options: ', netOption)
}
const timer = ref()
onMounted(() => {
    console.log('mounted')
    renderCPUPercent()
    renderCPU()
    renderMemInfo()
    renderMem()
    renderDiskInfo()
    renderDisk()
    renderNet()
    timer.value = setInterval(() => {
        console.log('start interval')
        renderCPUPercent()
        renderCPU()
        renderMemInfo()
        renderMem()
        renderDiskInfo()
        renderDisk()
        renderNet()
    }, 5000)
    console.log('timer: ', timer.value)
})

onUnmounted(() => {
    console.log('unmounted')
    clearInterval(timer.value)
})

watch(
    () => timeDensity.value,
    () => {
        renderCPUPercent()
        renderCPU()
        renderMemInfo()
        renderMem()
        renderDiskInfo()
        renderDisk()
        renderNet()
    }
)
</script>

<style scoped lang="scss">
@include b(host-container) {
    overflow: hidden;
    background-color: #ffffff;
    .el-row {
        margin-bottom: 4px;
        .el-col {
            height: 300px;
        }
    }

    .el-card {
        height: 100%;
        width: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            width: 100% !important;
        }
    }

    @include e(operator) {
        height: 48px;
        width: 100%;
        margin-bottom: 4px;
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

    @include e(image-title) {
        font-size: 14px;
        font-weight: bold;
    }

    @include e(image-description) {
        font-size: 12px;
        color: #73767a;
    }
}
</style>
