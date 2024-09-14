<template>
    <div class="am-profile-container">
        <div class="am-profile-container__title">
            <span>个人中心</span>
        </div>
    </div>
</template>
<script setup lang="ts"></script>
<style scoped lang="scss">
@include b(profile-container) {
    font-size: 14px;

    @include e(title) {
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 48px;
        width: 100%;
        background-color: #ffffff;
        // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

        border-radius: 4px;
        margin-bottom: 8px;
        padding: 0 16px;
        span {
            display: flex;
            align-items: center;
            font-size: 16px;
            font-weight: 600;
            color: #105eeb;
            margin-left: 16px;
            &::before {
                content: ' ';
                position: absolute;
                left: 20px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
</style>
