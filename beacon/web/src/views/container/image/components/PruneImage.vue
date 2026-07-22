<!--
    PruneImage — uses the shared ConfirmDialog via useConfirmCommand hook.
-->
<script setup lang="ts">
import { pruneImages } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const props = defineProps<{
  visible: boolean
  id?: string
  title?: string
  update?: () => void
}>()

const trigger = useConfirmCommand({
  title: props.title,
  message: 'image.confirmPrune',
  i18nPrefix: 'image',
  action: () => pruneImages(),
  onResolved: () => props.update?.(),
})

watch(() => props.visible, (visible) => {
  if (visible)
    trigger()
}, { immediate: true })
</script>

<template>
    <span v-if="false" />
</template>
