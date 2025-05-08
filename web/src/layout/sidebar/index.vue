<script setup lang="ts">
import CollapseIcon from '@/layout/sidebar/CollapseIcon.vue'
import Menuitem from '@/layout/sidebar/Menuitem.vue'
import { dynamicRoutes } from '@/router/dynamic.ts'
import useStore from '@/store'

const currentRoute = useRoute()
const store = useStore()

const menus = computed(() => {
  return dynamicRoutes.filter((item) => {
    if (item.meta?.show) {
      // 当 store.user.userInfo.name 不为 'admin' 时，不返回 item.name === 'user-manage' 的项
      if (store.user.userInfo.name !== 'admin' && item.name === 'account') {
        return false
      }
      // 其他情况下，返回所有 item.meta?.show 为 true 的项
      return true
    }
    return false
  })
})
</script>

<template>
    <el-aside :width="store.app.isCollapse ? '64px' : '210px'">
        <div class="am-logo">
            <img class="am-logo__img" src="@/assets/images/amprobe.png" alt="template" />
            <span v-show="!store.app.isCollapse" class="am-logo__text">Amprobe</span>
        </div>
        <div class="am-menu">
            <el-scrollbar>
                <el-menu
                    :default-active="currentRoute.path"
                    :collapse="store.app.isCollapse"
                    :unique-opened="true"
                    :collapse-transition="false"
                    mode="vertical"
                >
                    <Menuitem :items="menus" />
                </el-menu>
            </el-scrollbar>
        </div>
        <div class="am-collapse">
            <CollapseIcon />
        </div>
    </el-aside>
</template>

<style scoped lang="scss">
.el-aside {
  background-color: var(--el-menu-bg-color);
  //border-right: 1px solid var(--el-aside-border-color);
}

@include b(logo) {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  line-height: 64px;
  border-bottom: var(--el-aside-border-color) 1px solid;

  @include e(img) {
    height: 28px;
    object-fit: contain;
  }

  @include e(text) {
    font-size: 20px;
    font-weight: bold;
    margin-left: 8px;
    color: var(--el-aside-logo-text-color);
    white-space: nowrap;
  }
}

@include b(menu) {
  height: calc(100vh - 64px);
  :deep(.el-menu) {
    border: none !important;
  }
  :deep(.el-menu-item) {
    margin: 0 !important;
  }
  :deep(.el-sub-menu) {
    margin: 0 !important;
  }
}

@include b(collapse) {
  width: 12px;
  height: 32px;
  background: var(--el-menu-bg-color);
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
  border-color: var(--el-aside-border-color);
}
</style>
