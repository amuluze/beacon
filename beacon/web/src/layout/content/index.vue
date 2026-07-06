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

  width: 100%;
  height: calc(100vh - 72px);

  overflow: hidden;
  margin-top: 8px;
  position: relative;
  box-sizing: border-box;

  color: var(--el-header-text-color);
  background-color: var(--el-bg-color);
}
</style>
