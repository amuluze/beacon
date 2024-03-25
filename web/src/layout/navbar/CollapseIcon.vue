<!-- @format -->
<template>
    <div class="collapse-ico" @click="changeCollapse">
        <span v-if="!store.app.isCollapse"><i-ep-fold /></span>
        <span v-if="store.app.isCollapse"><i-ep-expand /></span>
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
    font-size: 22px;
    color: #303133;
    cursor: pointer;
}
</style>
