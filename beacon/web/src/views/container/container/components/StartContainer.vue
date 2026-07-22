<!--
    StartContainer — uses the shared ConfirmDialog via useConfirmCommand hook.
-->
<script setup lang="ts">
import { startContainer } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const props = defineProps<{
  visible: boolean
  id: string
  title?: string
  update?: () => void
}>()

const trigger = useConfirmCommand({
  title: props.title,
  message: 'container.confirmStart',
  i18nPrefix: 'container',
  action: id => startContainer({ container_id: id as string }),
  onResolved: () => props.update?.(),
})

watch(() => props.visible, (visible) => {
  if (visible)
    trigger(props.id)
}, { immediate: true })
</script>

<template>
    <span v-if="false" />
</template>
