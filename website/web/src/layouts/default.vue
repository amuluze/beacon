<template>
    <div class="site-layout">
        <Header />
        <main class="site-layout__main">
            <slot />
        </main>
        <Footer />
    </div>
</template>

<style scoped lang="scss">
.site-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.site-layout__main {
  flex: 1;
  padding-top: 64px;
}
</style>
