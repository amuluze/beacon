<template>
    <div class="am-navbar-container">
        <collapse-icon />
        <div>
            <a class="am-navbar-container__issue" href="https://support.qq.com/products/643578" target="_blank">
                问题反馈
            </a>
            <avatar />
        </div>
    </div>
</template>
<script setup lang="ts">
import Avatar from '@/layout/navbar/Avatar.vue'
import CollapseIcon from '@/layout/navbar/CollapseIcon.vue'
</script>

<style scoped lang="scss">
@include b(navbar-container) {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;

    @include e(issue) {
        font-size: 14px;
        font-weight: 400;
        margin-right: 16px;
    }
}
</style>
