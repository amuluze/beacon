<script setup lang="ts">
import { dayjs } from 'element-plus'
import type { SetSystemTimeArgs } from '@/interface/system.ts'
import { setSystemTime } from '@/api/system'
import { success } from '@/components/Message/message.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  systemTime: string
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

// 设置系统时间

const systemTime = ref('')
onMounted(() => {
  systemTime.value = props.systemTime
})

function getClientTime() {
  // 获取客户端时间
  console.log(new Date().toLocaleString())
  systemTime.value = dayjs(new Date()).format('YYYY-MM-DD HH:mm:ss')
}

async function confirmSystemTimeEdit() {
  // 确定修改系统时间
  systemTime.value = dayjs(new Date()).format('YYYY-MM-DD HH:mm:ss')
  console.log(systemTime.value)
  const args: SetSystemTimeArgs = {
    system_time: systemTime.value,
  }
  const { data } = await setSystemTime(args)
  console.log(data)
  success('修改系统时间成功')
  drawerVisible.value = false
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="50%" :title="t(props.title as string)">
        <el-form>
            <el-form-item :label="t('setting.systemTime')">
                <el-input v-model="systemTime" style="width: 240px" />
                <el-button type="info" style="margin-left: 16px" @click="getClientTime">
                    {{ t('setting.clientTime') }}
                </el-button>
            </el-form-item>
        </el-form>
        <el-button type="primary" size="default" plain @click="drawerVisible = false">
            {{ t('setting.cancel') }}
        </el-button>
        <el-button type="primary" size="default" plain @click="confirmSystemTimeEdit">
            {{ t('setting.confirm') }}
        </el-button>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
