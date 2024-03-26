<!-- @format -->
<template>
    <div class="collapse-ico" @click="changeCollapse">
        <svg-icon class="hamburger" :class="{ 'is-active': !store.app.isCollapse }" icon-class="indent-decrease" />
    </div>
</template>
<script setup lang="ts">
import useStore from '@/store'

const store = useStore()
const changeCollapse = () => store.app.setCollapse(!store.app.isCollapse)
</script>
<style scoped lang="scss">
.collapse-ico {
    margin-right: 10px;
    font-size: 15px;
    color: #303133;
    cursor: pointer;
}

.hamburger {
    vertical-align: middle;
    cursor: pointer;
    transform: scaleX(-1);
}

.hamburger.is-active {
    transform: scaleX(1);
}
</style>
