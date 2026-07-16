<template>
    <div class="site-layout">
        <a class="site-skip-link" href="#main-content">跳转到主要内容</a>
        <Header />
        <main id="main-content" class="site-layout__main" tabindex="-1">
            <slot />
        </main>
        <Footer />
        <SiteToast />
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
