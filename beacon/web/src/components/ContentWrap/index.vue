<script setup lang="ts">
defineProps<{
  title: string
  message: string
}>()
</script>

<template>
    <el-card shadow="never" class="am-card">
        <template v-if="title" #header>
            <div>
                <span class="am-card__title">
                    {{ title }}
                </span>
                <el-tooltip v-if="message" :content="message" placement="right" trigger="hover">
                    <template #content>
                        <div>{{ message }}</div>
                    </template>
                    <svg-icon icon-class="about" />
                </el-tooltip>
            </div>
        </template>
        <div>
            <slot />
        </div>
    </el-card>
</template>

<style scoped lang="scss">
@include b(card) {
  @include e(title) {
    font-size: 14px;
    font-weight: 500;
    margin-right: 4px;
    color: var(--el-text-color-primary);
  }
}
</style>
