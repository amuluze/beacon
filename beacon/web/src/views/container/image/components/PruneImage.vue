<!--
    PruneImage — uses the shared ConfirmDialog via useConfirmCommand hook.
-->
<script setup lang="ts">
import { pruneImages } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'
import IconDelete from '~icons/lucide/trash-2'

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
  badge: IconDelete,
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
