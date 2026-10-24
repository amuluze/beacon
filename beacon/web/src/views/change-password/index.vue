<script setup lang="ts">
import UpdatePassword from '@/layout/navbar/components/UpdatePassword.vue'
import useStore from '@/store'

const store = useStore()
const visible = shallowRef(true)
const username = computed(() => store.user.userInfo.name)
</script>

<template>
    <main class="mandatory-password-change">
        <UpdatePassword
            v-model:visible="visible"
            :username="username"
            mandatory
        />
    </main>
</template>

<style scoped>
.mandatory-password-change {
  min-height: 100vh;
  background: var(--am-surface-primary);
}
</style>
