import { defineVitestConfig } from '@nuxt/test-utils/config'

// 复用 Nuxt 的 vite 插件与 auto-import，composables 可直接测试
export default defineVitestConfig({
    test: {
        environment: 'happy-dom',
        include: ['src/**/*.{test,spec}.ts'],
    },
})
