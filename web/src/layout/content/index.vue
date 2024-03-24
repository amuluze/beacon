<template>
    <section class="am-content-container">
        <router-view />
    </section>
</template>
<script setup lang="ts"></script>

<style scoped lang="scss">
@include b(content-container) {
    height: 100%;
    width: 100%;
}
</style>
