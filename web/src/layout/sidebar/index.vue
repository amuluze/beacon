<template>
    <el-aside :width="store.app.isCollapse ? '54px' : '160px'">
        <div class="am-logo">
            <img class="am-logo__img" src="@/assets/images/amprobe.png" alt="" />
            <span v-show="!store.app.isCollapse" class="am-logo__text">Amprobe</span>
        </div>
        <div class="am-menu">
            <el-scrollbar>
                <el-menu
                    :default-active="onRoutes"
                    :collapse="store.app.isCollapse"
                    :unique-opened="true"
                    :collapse-transition="false"
                    background-color="#E9EFFD"
                    text-color="#000"
                    active-text-color="#105EEB"
                    mode="vertical"
                >
                    <menuitem v-for="(item, index) in dRoutes" :key="index" :item="item" />
                </el-menu>
            </el-scrollbar>
        </div>
        <div class="am-user">
            <avatar />
        </div>
        <div class="am-collapse">
            <collapse-icon />
        </div>
    </el-aside>
</template>
<script setup lang="ts">
import Avatar from '@/layout/sidebar/Avatar.vue'
import CollapseIcon from '@/layout/sidebar/CollapseIcon.vue'
import Menuitem from '@/layout/sidebar/Menuitem.vue'
import { dynamicRoutes } from '@/router/dynamic.ts'
import useStore from '@/store'

const currentRoute = useRoute()

const store = useStore()

const dRoutes = computed(() => {
    return dynamicRoutes.filter((item) => item.meta?.show)
})

const onRoutes = computed(() => {
    if (currentRoute.meta.activeMenu) return currentRoute.meta.activeMenu as string
    return currentRoute.path
})
</script>

<style scoped lang="scss">
@include b(logo) {
    @include bfc;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 64px;
    line-height: 64px;
    border-bottom: #cdcfd2 1px solid;

    @include e(img) {
        height: 28px;
        width: 28px;
        margin-right: 6px;
    }
    @include e(text) {
        font-size: 20px;
        font-weight: bold;
        color: #105eeb;
        white-space: nowrap;
    }
}

@include b(menu) {
    height: calc(100vh - 64px - 64px);
    .el-scrollbar {
        .is-active {
            color: #409eff !important;
        }
    }

    .el-menu {
        width: 100% !important;
        height: 100%;
        background-color: #e9effd !important;
        border: none;
    }
}

@include b(user) {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 64px;
    border-top: #cdcfd2 1px solid;
    cursor: pointer;
}

@include b(collapse) {
    width: 12px;
    height: 32px;
    background: #e9effd;
    border-radius: 0 5px 5px 0;
    display: flex;
    justify-content: center;
    align-items: center;
    position: absolute;
    bottom: 16px;
    z-index: 999;
    right: -12px;
    border-top: 1px solid;
    border-bottom: 1px solid;
    border-right: 1px solid;
    border-color: #e9e9ea;
}
</style>
