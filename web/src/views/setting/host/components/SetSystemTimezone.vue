<script setup lang="ts">
import type { SetSystemTimezoneArgs } from '@/interface/system.ts'
import { getSystemTimezoneList, setSystemTimezone } from '@/api/system'
import { success } from '@/components/Message/message.ts'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  systemTimezone: string
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

// 设置系统时区
const systemTimezone = ref('')
onMounted(async () => {
  systemTimezone.value = props.systemTimezone
  await querySystemTimezoneOptions()
})

const timezoneOptions = ref<
  {
    label: string
    value: string
  }[]
>([])

async function querySystemTimezoneOptions() {
  // 获取系统时区选项
  const { data } = await getSystemTimezoneList()
  console.log(data)
  console.log('zone list: ', data.system_timezone_list)
  data.system_timezone_list.forEach((item) => {
    timezoneOptions.value.push({
      label: item,
      value: item,
    })
  })
  console.log(timezoneOptions.value)
}

async function confirmSystemTimezoneEdit() {
  // 确定修改系统时区
  console.log(systemTimezone.value)
  const args: SetSystemTimezoneArgs = {
    system_timezone: systemTimezone.value,
  }
  const { data } = await setSystemTimezone(args)
  console.log(data)
  success('修改系统时区成功')
  drawerVisible.value = false
}

const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="540" :title="t(props.title as string)">
        <el-form>
            <el-form-item :label="t('setting.systemTimezone')">
                <el-select v-model="systemTimezone" style="width: 240px" filterable>
                    <el-option v-for="item in timezoneOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
        </el-form>
        <el-button type="primary" size="default" plain @click="drawerVisible = false">
            {{ t('setting.cancel') }}
        </el-button>
        <el-button type="primary" size="default" plain @click="confirmSystemTimezoneEdit">
            {{ t('setting.confirm') }}
        </el-button>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
