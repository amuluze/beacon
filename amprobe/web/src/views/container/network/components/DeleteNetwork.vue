<!--
    DeleteNetwork — uses the shared ConfirmDialog via useConfirmCommand hook.
-->
<script setup lang="ts">
import { removeNetwork } from '@/api/container'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const props = defineProps<{
    visible: boolean
    id: string
    title?: string
    update?: () => void
}>()

const trigger = useConfirmCommand({
    title: props.title,
    message: 'network.confirmDelete',
    i18nPrefix: 'network',
    action: (id) => removeNetwork({ network_id: id as string }),
    onResolved: () => props.update?.(),
})

watch(() => props.visible, (visible) => {
    if (visible) trigger(props.id)
})
</script>

<template></template>
