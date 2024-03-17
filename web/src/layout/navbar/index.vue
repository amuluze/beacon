<template>
    <div class="am-navbar-container">
        <div class="am-navbar-container__logo">
            <img src="@/assets/amprobe.png" alt="" />
            <span class="am-logo__text">Amprobe</span>
        </div>
        <div class="am-navbar-container__tab">
            <el-menu :default-active="activeIndex" :ellipsis="false" mode="horizontal" @select="handleSelect">
                <el-menu-item index="1">容器</el-menu-item>
                <el-menu-item index="2">主机</el-menu-item>
            </el-menu>
            <div>
                <avatar></avatar>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import Avatar from '@/layout/navbar/Avatar.vue';

const router = useRouter();
const activeIndex = computed(() => (router.currentRoute.value.name === 'container' ? '1' : '2'));

const handleSelect = (key: string) => {
    console.log(key);
    if (key === '1') {
        router.replace({ name: 'container' });
    }
    if (key === '2') {
        router.replace({ name: 'host' });
    }
};
</script>

<style scoped lang="scss">
@include b(navbar-container) {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;

    @include e(logo) {
        color: #409eff;
        font-size: 24px;
        font-weight: bold;
        width: 150px;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        img {
            height: 48px;
            width: 48px;
            margin-right: 6px;
        }
    }

    @include e(tab) {
        display: flex;
        margin-left: 32px;
        justify-content: space-between;
        align-items: center;
        height: 100%;
        width: 100%;
    }
}
</style>
