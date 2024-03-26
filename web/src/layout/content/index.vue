<template>
    <section class="am-content-container">
        <router-view />
    </section>
</template>
<script setup lang="ts"></script>

<style scoped lang="scss">
@include b(content-container) {
    position: relative;
    box-sizing: border-box;
    width: 100%;
    /* 50 = navbar 50 */
    min-height: calc(100vh - 48px);
    overflow: hidden;
}
</style>
