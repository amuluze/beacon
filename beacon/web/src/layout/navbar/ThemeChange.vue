<script setup lang="ts">
import useStore from '@/store'
import useTheme from '@/hooks/useTheme.ts'

const store = useStore()

const { switchDark } = useTheme()

function themeChange() {
  switchDark()
}
const dark = ref(store.theme.dark)
</script>

<template>
    <el-switch v-model="dark" class="mr-4" @change="themeChange">
        <template #active-action>
            <svg-icon icon-class="moon" />
        </template>
        <template #inactive-action>
            <svg-icon icon-class="sunny" />
        </template>
    </el-switch>
</template>

<style scoped lang="scss">
</style>
