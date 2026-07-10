<script setup lang="ts">
import Content from '@/layout/content/index.vue'
import Navbar from '@/layout/navbar/index.vue'
import StatusBar from '@/layout/statusbar/index.vue'
</script>

<template>
    <section class="am-layout">
        <Navbar />
        <Content />
        <StatusBar />
    </section>
</template>

<style scoped lang="scss">
@include b(layout) {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  color: var(--am-foreground-primary);
  background-color: var(--am-surface-primary);
  overflow: hidden;
}
</style>
