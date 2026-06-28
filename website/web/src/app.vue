<template>
    <div class="app">
        <!-- <ParticleBackground /> -->
        <NuxtLayout>
            <NuxtPage />
        </NuxtLayout>
    </div>
</template>

<style lang="scss">
.app {
    position: relative;
    min-height: 100vh;
    @include bg-gradient('primary');
}
</style>
