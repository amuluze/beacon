<!--
    DeleteContainer dialog wrapper.

    Historically this file reimplemented the el-dialog + loading state +
    footer button row inline. We keep it as the smallest possible wrapper to
    preserve existing call sites (`useCommandComponent(DeleteContainer)`)
    while delegating the actual dialog UI to the shared ConfirmDialog.
-->
<script setup lang="ts">
import { removeContainer } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const props = defineProps<{
  visible: boolean
  id: string
  title?: string
  update?: () => void
}>()

// The shared command hook renders ConfirmDialog imperatively. We expose a
// function-style interface by binding visibility to the close lifecycle.
const trigger = useConfirmCommand({
  title: props.title,
  message: 'container.confirmDelete',
  i18nPrefix: 'container',
  action: id => removeContainer({ container_id: id as string }),
  onResolved: () => props.update?.(),
})

// The parent component calls `componentApi({ id, ... })` via
// useCommandComponent — we trigger the dialog when `visible` flips to true.
watch(
  () => props.visible,
  (visible) => {
    if (visible)
      trigger(props.id)
  },
  { immediate: true },
)
</script>

<template>
    <span v-if="false" />
</template>
