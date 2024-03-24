<template>
    <div class="am-navbar-container">
        <collapse-icon />
        <avatar></avatar>
    </div>
</template>
<script setup lang="ts">
import Avatar from '@/layout/navbar/Avatar.vue';
import CollapseIcon from '@/layout/navbar/CollapseIcon.vue';

</script>

<style scoped lang="scss">
@include b(navbar-container) {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    background-color: #ffffff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;
}
</style>
