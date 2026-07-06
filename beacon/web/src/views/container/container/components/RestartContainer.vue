<!--
    RestartContainer — uses the shared ConfirmDialog via useConfirmCommand hook.
-->
<script setup lang="ts">
import { restartContainer } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const props = defineProps<{
    visible: boolean
    id: string
    title?: string
    update?: () => void
}>()

const trigger = useConfirmCommand({
    title: props.title,
    message: 'container.confirmRestart',
    i18nPrefix: 'container',
    action: (id) => restartContainer({ container_id: id as string }),
    onResolved: () => props.update?.(),
})

watch(() => props.visible, (visible) => {
    if (visible) trigger(props.id)
})
</script>

<template></template>
