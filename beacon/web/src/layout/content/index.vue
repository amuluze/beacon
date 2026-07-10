<script setup lang="ts"></script>

<template>
    <section class="am-content">
        <router-view />
    </section>
</template>

<style scoped lang="scss">
@include b(content) {
  display: flex;
  flex-direction: column;
  flex: 1;
  width: 100%;
  min-height: 0;
  overflow: auto;
  position: relative;
  box-sizing: border-box;

  color: var(--am-foreground-primary);
  background-color: var(--am-surface-primary);
}
</style>
