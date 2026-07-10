<script setup lang="ts">
import { updateAlarmThreshold } from '@/api/alarm'
import { success } from '@/components/Message/message.ts'
import type { AlarmThreshold } from '@/interface/alarm.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  threshold: AlarmThreshold
  update?: () => void | Promise<void>
}>()

const emit = defineEmits<{
  'update:visible': [visible: boolean]
  'close': []
}>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
    if (!visible)
      emit('close')
  },
})
const form = reactive<AlarmThreshold>({ ...props.threshold })
const loading = shallowRef(false)
const durationOptions = Array.from({ length: 9 }, (_, index) => index + 2)
const { t } = useI18n()

async function submit() {
  loading.value = true
  try {
    await updateAlarmThreshold({ ...form, threshold: Number(form.threshold) })
    success('修改成功')
    dialogVisible.value = false
    await props.update?.()
  }
  catch {
    // 请求拦截器负责展示服务端错误；保留弹窗和本地输入供用户重试。
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
    <el-dialog v-model="dialogVisible" width="480px" :title="t(props.title ?? 'setting.memAlarmThreshold')">
        <el-form :model="form" label-position="top">
            <el-form-item :label="t('setting.memAlarmThreshold')">
                <div class="threshold-row">
                    <span>{{ t('setting.memUsage') }}</span>
                    <el-select v-model="form.duration">
                        <el-option v-for="duration in durationOptions" :key="duration" :value="duration" :label="`${duration} 分钟`" />
                    </el-select>
                    <span>{{ t('setting.over') }}</span>
                    <el-input v-model="form.threshold">
                        <template #append>
                            %
                        </template>
                    </el-input>
                </div>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="dialogVisible = false">
                {{ t('setting.cancel') }}
            </el-button>
            <el-button :loading="loading" type="primary" @click="submit">
                {{ t('setting.confirm') }}
            </el-button>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
.threshold-row {
  display: grid;
  grid-template-columns: auto 110px auto minmax(120px, 1fr);
  gap: var(--am-spacing-sm);
  align-items: center;
  width: 100%;
}

@media (max-width: 520px) {
  .threshold-row {
    grid-template-columns: 1fr;
  }
}
</style>
