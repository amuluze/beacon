<script setup lang="ts">
</script>

<template>
    <el-container class="layout-container">
        <Header />
        <el-main class="main-content">
            <NuxtPage :transition="{ name: 'page', mode: 'out-in' }" />
        </el-main>
        <Footer />
    </el-container>
</template>

<style scoped lang="scss">
.layout-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;

  .main-content {
    flex: 1;
    padding: 0 !important;
    transition: all 0.3s ease;
  }
}
</style>
