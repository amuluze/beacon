<script setup lang="ts">
import { dayjs } from 'element-plus'
import { useI18n } from 'vue-i18n'

interface Props {
  stale?: boolean
  timestamp?: number
}

const props = withDefaults(defineProps<Props>(), {
  stale: false,
  timestamp: 0,
})

const { t } = useI18n()
const tooltip = computed(() => {
  if (!props.timestamp)
    return t('agent.staleTooltipDefault')
  return t('agent.staleTooltipWithTime', {
    time: dayjs(props.timestamp * 1000).format('YYYY-MM-DD HH:mm:ss'),
  })
})
</script>

<template>
    <el-tooltip v-if="stale" :content="tooltip" placement="top">
        <el-tag type="warning" effect="plain" size="small">
            {{ t('agent.stale') }}
        </el-tag>
    </el-tooltip>
</template>
