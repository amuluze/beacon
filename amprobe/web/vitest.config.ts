import vue from '@vitejs/plugin-vue'
import { resolve } from 'node:path'
import { defineConfig } from 'vitest/config'

import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

import UnoCSS from 'unocss/vite'

export default defineConfig({
    plugins: [
        vue(),
        UnoCSS(),
        AutoImport({
            imports: ['vue', 'vue-router', 'pinia'],
            resolvers: [ElementPlusResolver(), IconsResolver({ enabledCollections: ['ep'] })],
            vueTemplate: true,
            dts: resolve(resolve(__dirname, 'types'), 'auto-imports.d.ts'),
        }),
        Components({
            resolvers: [ElementPlusResolver(), IconsResolver({ enabledCollections: ['ep'] })],
            dts: resolve(resolve(__dirname, 'types'), 'components.d.ts'),
        }),
        Icons({ autoInstall: true }),
        createSvgIconsPlugin({
            iconDirs: [resolve(__dirname, 'src/assets/icons')],
            symbolId: 'icon-[dir]-[name]',
        }),
    ],
    resolve: {
        alias: {
            '@': resolve(__dirname, 'src'),
        },
        extensions: ['.js', '.ts', '.jsx', '.tsx', '.json', '.vue', '.mjs'],
    },
    test: {
        environment: 'happy-dom',
        globals: true,
        setupFiles: ['src/__tests__/setup.ts'],
        coverage: {
            provider: 'v8',
            reporter: ['text', 'html'],
        },
    },
})
