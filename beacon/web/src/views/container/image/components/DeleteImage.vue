<!--
    DeleteImage — uses the shared ConfirmDialog via useConfirmCommand hook.
    The previous inline el-dialog implementation has been removed.
-->
<script setup lang="ts">
import { removeImage } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const props = defineProps<{
  visible: boolean
  id: string
  title?: string
  update?: () => void
}>()

const trigger = useConfirmCommand({
  title: props.title,
  message: 'image.confirmDelete',
  i18nPrefix: 'image',
  action: id => removeImage({ image_id: id as string }),
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
