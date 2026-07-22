<!-- @format -->
<script setup lang="ts">
import useStore from '@/store'

const store = useStore()
const changeCollapse = () => store.app.setCollapse(!store.app.isCollapse)
</script>

<template>
    <div class="collapse-ico" @click="changeCollapse">
        <svg-icon v-if="store.app.isCollapse" icon-class="right" />
        <svg-icon v-else icon-class="left" />
    </div>
</template>

<style scoped lang="scss">
.collapse-ico {
  font-size: 8px;
  color: #424345;
  cursor: pointer;
}
</style>
