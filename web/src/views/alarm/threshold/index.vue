<template>
    <div class="am-alarm-title">
        <span @click="$router.push('/alarm/mail')">告警配置</span>
        <span @click="$router.push('/alarm/threshold')">告警阈值</span>
    </div>
    <el-card shadow="never">
        <el-descriptions title="主机运行状态" :column=1>
            <el-descriptions-item label="CPU 使用率过高">
                <span style="margin-right: 8px">CPU 使率连续 {{ CPUThreshold.duration }} 分钟超过 {{ CPUThreshold.threshold }}%</span>
                <svg-icon icon-class="edit" style="cursor: pointer" @click="editCPUThreshold" />
            </el-descriptions-item>
            <el-descriptions-item label="内存使用率过高">
                <span style="margin-right: 8px">内存使用率连续 {{ MemoryThreshold.duration }} 分钟超过 {{ MemoryThreshold.threshold }}%</span>
                <svg-icon icon-class="edit" style="cursor: pointer" @click="editMemoryThreshold" />
            </el-descriptions-item>
            <el-descriptions-item label="磁盘使用率过高">
                <span style="margin-right: 8px">磁盘使用率超过 {{ DiskThreshold.threshold }}%</span>
                <svg-icon icon-class="edit" style="cursor: pointer" @click="editDiskThreshold" />
            </el-descriptions-item>
        </el-descriptions>
    </el-card>
    <el-drawer v-model="CPUDrawer" size="560" title="编辑 CPU 使用率过高">
        <el-form ref="CPUSettingRef" :model="CPUThreshold" label-width="120px" label-position="left">
            <el-form-item prop="server" label="CPU 使用率过高">
            CPU 使用率连续
            <el-select v-model="CPUThreshold.duration" size="small" style="width: 100px;margin: 0 4px">
                <el-option v-for="item in optionOptions" :key="item.value" :value="item.value" :label="item.label" />
            </el-select>
            超过
            <el-input v-model="CPUThreshold.threshold" size="small" style="width: 140px;margin: 0 4px"><template #append>%</template></el-input>
            </el-form-item>
        </el-form>
        <el-button type="default" size="small" plain @click="cacelEditCPUThreshold">取消</el-button>
        <el-button type="primary" size="small" plain @click="confirmEditCPUThreshold">确定</el-button>
    </el-drawer>
    <el-drawer v-model="MemoryDrawer" size="560" title="编辑内存使用率过高">
        <el-form ref="MemorySettingRef" :model="MemoryThreshold" label-width="120px" label-position="left">
            <el-form-item prop="server" label="内存使用率过高">
            CPU 使用率连续
            <el-select v-model="MemoryThreshold.duration" size="small" style="width: 100px;margin: 0 4px">
                <el-option v-for="item in optionOptions" :key="item.value" :value="item.value" :label="item.label" />
            </el-select>
            超过
            <el-input v-model="MemoryThreshold.threshold" size="small" style="width: 140px;margin: 0 4px"><template #append>%</template></el-input>
            </el-form-item>
        </el-form>
        <el-button type="default" size="small" plain @click="cacelEditMemoryhreshold">取消</el-button>
        <el-button type="primary" size="small" plain @click="confirmEditMemoryThreshold">确定</el-button>
    </el-drawer>
    <el-drawer v-model="DiskDrawer" size="560" title="编辑磁盘使用率过高">
        <el-form ref="DiskSettingRef" :model="DiskThreshold" label-width="120px" label-position="left">
            <el-form-item prop="server" label="磁盘使用率过高">
            磁盘使用率超过
            <el-input v-model="DiskThreshold.threshold" size="small" style="width: 140px;margin: 0 4px"><template #append>%</template></el-input>
            </el-form-item>
        </el-form>
        <el-button type="default" size="small" plain @click="cacelEditDiskThreshold">取消</el-button>
        <el-button type="primary" size="small" plain @click="confirmEditDiskThreshold">确定</el-button>
    </el-drawer>
</template>
<script setup lang="ts">
import { queryAlarmThreshold, updateAlarmThreshold } from '@/api/alarm';
import { success } from '@/components/Message/message';
import { AlarmThresholdUpdateArgs } from '@/interface/alarm';


onMounted(async () => {
    const {data} = await queryAlarmThreshold()
    for (const el of data.data) {
        if (el.type === "cpu") {
            CPUThreshold.value.id = el.id
            CPUThreshold.value.duration = el.duration
            CPUThreshold.value.threshold = el.threshold
        } else if (el.type === 'memory') {
            MemoryThreshold.value.id = el.id
            MemoryThreshold.value.duration = el.duration
            MemoryThreshold.value.threshold = el.threshold
        } else if (el.type === "disk") {
            DiskThreshold.value.id = el.id
            DiskThreshold.value.duration = el.duration
            DiskThreshold.value.threshold = el.threshold
        }
    }
})


const CPUThreshold = ref({
    id: 1,
    duration: 2,
    threshold: 80,
})
const MemoryThreshold = ref({
    id: 2,
    duration: 2,
    threshold: 80,
})
const DiskThreshold = ref({
    id: 3,
    duration: 2,
    threshold: 80,
})
const CPUDrawer = ref(false)
const editCPUThreshold = () => CPUDrawer.value = true
const cacelEditCPUThreshold = () => CPUDrawer.value = false
const confirmEditCPUThreshold = () => {
    let parmas: AlarmThresholdUpdateArgs = {
        id: CPUThreshold.value.id,
        type: "cpu",
        duration: CPUThreshold.value.duration,
        threshold: Number(CPUThreshold.value.threshold)
    }
    console.log(parmas)
    updateAlarmThreshold(parmas).catch((err) => {
        console.log(err)
    }).finally(() => {
        success("修改成功")
        CPUDrawer.value = false
    })
}
const optionOptions = [
    {
        value: 2,
        label: '2分钟'
    },
    {
        value: 3,
        label: '3分钟'
    },
    {
        value: 4,
        label: '4分钟'
    },
    {
        value: 5,
        label: '5分钟'
    },
    {
        value: 6,
        label: '6分钟'
    },
    {
        value: 7,
        label: '7分钟'
    },
    {
        value: 8,
        label: '8分钟'
    },
    {
        value: 9,
        label: '9分钟'
    },
    {
        value: 10,
        label: '10分钟'
    },
]

const MemoryDrawer = ref(false)
const editMemoryThreshold = () => MemoryDrawer.value = true
const cacelEditMemoryhreshold = () => MemoryDrawer.value = false
const confirmEditMemoryThreshold = () => {
    let params: AlarmThresholdUpdateArgs = {
        id: MemoryThreshold.value.id,
        type: "memory",
        duration: MemoryThreshold.value.duration,
        threshold: Number(MemoryThreshold.value.threshold)
    }
    updateAlarmThreshold(params).finally(() => {
        success("修改成功")
        MemoryDrawer.value = false
    })

}

const DiskDrawer = ref(false)
const editDiskThreshold = () => DiskDrawer.value = true
const cacelEditDiskThreshold = () => DiskDrawer.value = false
const confirmEditDiskThreshold = () => {
    let params: AlarmThresholdUpdateArgs = {
        id: DiskThreshold.value.id,
        type: "disk",
        duration: DiskThreshold.value.duration,
        threshold: Number(DiskThreshold.value.threshold)
    }
    updateAlarmThreshold(params).finally(() => {
        success("修改成功")
        DiskDrawer.value = false
    })
}
</script>
<style scoped lang="scss">
@include b(alarm-title) {
    display: flex;
    align-items: center;
    justify-content: space-start;
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

        &:last-child {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 116px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
</style>
