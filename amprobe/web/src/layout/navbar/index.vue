<script setup lang="ts">
import Avatar from '@/layout/navbar/Avatar.vue';
import Language from '@/layout/navbar/Language.vue'; /* PartiallyEnd: #3632/scriptSetup.vue */
import ThemeChange from '@/layout/navbar/ThemeChange.vue';
import Breadcrumb from './Breadcrumb.vue';
</script>

<template>
    <div class="am-navbar">
        <Breadcrumb style="cursor: pointer" />
        <div class="am-navbar__right">
            <Language />
            <ThemeChange />
            <Avatar />
        </div>
    </div>
</template>

<style scoped lang="scss">
@include b(navbar) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  width: 100%;
  padding: 0 16px;
  color: var(--el-header-text-color);
  background-color: var(--el-menu-bg-color);
  border-radius: 4px;
  border-bottom: 1px solid var(--el-header-border-color);
  border-left: 1px solid var(--el-header-border-color);

  @include e(right) {
    display: flex;
    align-items: center;
    justify-content: flex-end;
  }
}
</style>
