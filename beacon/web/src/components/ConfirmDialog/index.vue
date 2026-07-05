<!--
    Generic confirmation dialog used to ask the user before invoking a
    single-pass destructive or operational action.

    Pattern extracted from:
      - container/container/components/{Start,Stop,Restart,Delete}Container.vue
      - container/network/components/DeleteNetwork.vue
      - container/image/components/DeleteImage.vue
      - container/image/components/PruneImage.vue

    Each of those previously reimplemented the same v-model, loading flag,
    and footer button row.
-->
<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = withDefaults(defineProps<{
    visible: boolean
    title?: string
    // The text shown in the dialog body. May be a key into vue-i18n.
    message?: string
    // i18n locale key prefix for cancel/confirm labels, e.g. 'container' or 'network'.
    // Falls back to 'common' if unset.
    i18nPrefix?: string
    // If true, only renders the footer confirm button as a primary action
    // (used for risky confirmations on Start / Stop / Restart).
    confirmation?: boolean
    // Width override for the dialog (defaults to 500px to match existing popups).
    width?: string
}>(), {
    confirmation: false,
    width: '500px',
    i18nPrefix: 'common',
})

const emits = defineEmits<{
    (e: 'update:visible', visible: boolean): void
    (e: 'close'): void
    (e: 'confirm'): void
}>()

const dialogVisible = computed<boolean>({
    get() {
        return props.visible
    },
    set(visible) {
        emits('update:visible', visible)
        if (!visible) {
            emits('close')
        }
    },
})

const loading = ref(false)
function onConfirm() {
    emits('confirm')
}

const { t } = useI18n()

const cancelLabel = computed(() => t(`${props.i18nPrefix}.cancel`))
const confirmLabel = computed(() => t(`${props.i18nPrefix}.confirm`))

// Expose a way for callers to dismiss the dialog after their async action resolves.
defineExpose({
    close() {
        dialogVisible.value = false
    },
    confirm() {
        onConfirm()
    },
})
</script>

<template>
    <el-dialog
        v-model="dialogVisible"
        :title="title ? t(title) : ''"
        :width="width"
        draggable
    >
        <span>{{ message ? t(message) : '' }}</span>
        <template #footer>
            <span class="dialog-footer">
                <el-button :disabled="loading" @click="dialogVisible = false">
                    {{ cancelLabel }}
                </el-button>
                <el-button v-if="confirmation" type="primary" :loading="loading" @click="onConfirm">
                    {{ confirmLabel }}
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
</style>
