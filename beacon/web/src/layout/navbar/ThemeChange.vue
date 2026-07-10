<script setup lang="ts">
import useStore from '@/store'
import useTheme from '@/hooks/useTheme.ts'

const store = useStore()

const { switchDark } = useTheme()

function themeChange() {
  switchDark()
}
const dark = computed(() => store.theme.dark)
</script>

<template>
    <button class="am-topbar-action" type="button" :aria-label="dark ? '切换到浅色主题' : '切换到深色主题'" @click="themeChange">
        <svg-icon :icon-class="dark ? 'moon' : 'sunny'" size="14px" />
    </button>
</template>

<style scoped lang="scss">
</style>
