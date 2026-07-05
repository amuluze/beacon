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

// 编辑磁盘告警阈值
const diskThreshold = ref<AlarmThreshold>({
  id: 0,
  type: 'disk',
  duration: 0,
  threshold: 0,
})

onMounted(() => {
  diskThreshold.value.id = props.threshold.id
  diskThreshold.value.type = props.threshold.type
  diskThreshold.value.duration = props.threshold.duration
  diskThreshold.value.threshold = props.threshold.threshold
})

function confirmEditDiskThreshold() {
  const params: AlarmThresholdUpdateArgs = {
    id: diskThreshold.value.id,
    type: 'disk',
    duration: diskThreshold.value.duration,
    threshold: Number(diskThreshold.value.threshold),
  }
  updateAlarmThreshold(params).finally(() => {
    success('修改成功')
    drawerVisible.value = false
  })
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="50%" :title="t(props.title as string)">
        <el-form :model="diskThreshold" label-width="160px" label-position="left">
            <el-form-item prop="server" :label="t('setting.diskAlarmThreshold')">
                {{ t('setting.diskUsage') }}
                <el-input v-model="diskThreshold.threshold" size="small" style="width: 140px;margin: 0 4px">
                    <template #append>
                        %
                    </template>
                </el-input>
            </el-form-item>
        </el-form>
        <el-button type="primary" size="small" plain @click="drawerVisible = false">
            {{ t('setting.cancel') }}
        </el-button>
        <el-button type="primary" size="small" plain @click="confirmEditDiskThreshold">
            {{ t('setting.confirm') }}
        </el-button>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
