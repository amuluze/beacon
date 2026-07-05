<script setup lang="ts">
import type { AlarmThreshold, AlarmThresholdUpdateArgs } from '@/interface/alarm.ts'
import { updateAlarmThreshold } from '@/api/alarm'
import { success } from '@/components/Message/message.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  threshold: AlarmThreshold
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const drawerVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 编辑 CPU 告警阈值
const cpuThreshold = ref<AlarmThreshold>({
  id: 0,
  type: 'cpu',
  duration: 0,
  threshold: 0,
})

const optionOptions = [
  {
    value: 2,
    label: '2分钟',
  },
  {
    value: 3,
    label: '3分钟',
  },
  {
    value: 4,
    label: '4分钟',
  },
  {
    value: 5,
    label: '5分钟',
  },
  {
    value: 6,
    label: '6分钟',
  },
  {
    value: 7,
    label: '7分钟',
  },
  {
    value: 8,
    label: '8分钟',
  },
  {
    value: 9,
    label: '9分钟',
  },
  {
    value: 10,
    label: '10分钟',
  },
]

onMounted(() => {
  cpuThreshold.value.id = props.threshold.id
  cpuThreshold.value.type = props.threshold.type
  cpuThreshold.value.duration = props.threshold.duration
  cpuThreshold.value.threshold = props.threshold.threshold
})

function confirmEditCPUThreshold() {
  const params: AlarmThresholdUpdateArgs = {
    id: cpuThreshold.value.id,
    type: cpuThreshold.value.type,
    duration: cpuThreshold.value.duration,
    threshold: Number(cpuThreshold.value.threshold),
  }
  console.log(params)
  updateAlarmThreshold(params).catch((err) => {
    console.log(err)
  }).finally(() => {
    success('修改成功')
    drawerVisible.value = false
  })
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="50%" :title="t(props.title as string)">
        <el-form :model="cpuThreshold" label-width="160px" label-position="left">
            <el-form-item prop="server" :label="t('setting.cpuAlarmThreshold')">
                {{ t('setting.cpuUsage') }}
                <el-select v-model="cpuThreshold.duration" size="small" style="width: 100px;margin: 0 4px">
                    <el-option v-for="item in optionOptions" :key="item.value" :value="item.value" :label="item.label" />
                </el-select>
                {{ t('setting.over') }}
                <el-input v-model="cpuThreshold.threshold" size="small" style="width: 140px;margin: 0 4px">
                    <template #append>
                        %
                    </template>
                </el-input>
            </el-form-item>
        </el-form>
        <el-button type="primary" size="small" plain @click="drawerVisible = false">
            {{ t('setting.cancel') }}
        </el-button>
        <el-button type="primary" size="small" plain @click="confirmEditCPUThreshold">
            {{ t('setting.confirm') }}
        </el-button>
    </el-drawer>
</template>

<style scoped lang="scss">
</style>
