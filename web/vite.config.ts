import vue from '@vitejs/plugin-vue'
import { resolve } from 'node:path'
import { defineConfig } from 'vite'

import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

import UnoCSS from 'unocss/vite'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        UnoCSS(),
        AutoImport({
            // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
            imports: ['vue', 'vue-router', 'pinia'],
            eslintrc: {
                enabled: false, // 是否自动生成 eslint 规则，建议生成之后设置 false，手动维护
                filepath: './.eslintrc.cjs-auto-import.json', // 指定自动导入函数 eslint 规则的文件路径
                globalsPropValue: true,
            },
            resolvers: [
                // 自动导入 Element Plus 相关函数，如：ElMessage(带样式)
                ElementPlusResolver(),
                IconsResolver({
                    enabledCollections: ['ep'],
                }),
            ],
            vueTemplate: true, // 是否在 vue 模板中自动导入
            dts: resolve(resolve(__dirname, 'types'), 'auto-imports.d.ts'), // 指定自动导入函数TS类型声明文件路径
        }),
        Components({
            resolvers: [
                // 自动导入 Element Plus 组件
                ElementPlusResolver(),
                // 自动导入图标组件
                IconsResolver({
                    // @iconify-json/ep 是 Element Plus 的图标库
                    enabledCollections: ['ep'],
                }),
            ],
            dts: resolve(resolve(__dirname, 'types'), 'components.d.ts'), // 指定自动导入组件TS类型声明文件路径
        }),
        Icons({
            // 自动安装图标库
            autoInstall: true,
        }),
        createSvgIconsPlugin({
            // 指定需要缓存的图标文件夹
            // iconDirs: [resolve(process.cwd(), 'src/assets/icons'), resolve(process.cwd(), 'src/assets/error')],
            iconDirs: [resolve(__dirname, 'src/assets/icons')],
            // 指定symbolId格式
            symbolId: 'icon-[dir]-[name]',
        }),
    ],
    resolve: {
        alias: {
            '@': resolve(__dirname, 'src'),
        },
        // 导入时，忽略后缀
        extensions: ['.js', '.ts', '.jsx', '.tsx', '.json', '.vue', '.mjs'],
    },
    server: {
        host: '0.0.0.0',
        port: 9000,
        open: true,
        proxy: {
            '/api/': {
                target: 'http://124.70.156.29:8000',
                changeOrigin: true,
            },
        },
    },
    css: {
        preprocessorOptions: {
            scss: {
                // https://blog.csdn.net/CssHero/article/details/142686148
                // Deprecation Warning [legacy-js-api]: The legacy JS API is deprecated and will be removed in Dart Sass 2.0.0.
                api: 'modern-compiler',
                additionalData: '@use "@/styles/bem.scss" as *;',
            },
        },
    },
    // 构建配置
    build: {
        chunkSizeWarningLimit: 2000, // 消除打包大小超过500kb警告
        minify: 'terser', // Vite 2.6.x 以上需要配置 minify: "terser", terserOptions 才能生效
        terserOptions: {
            compress: {
                keep_infinity: true, // 防止 Infinity 被压缩成 1/0，这可能会导致 Chrome 上的性能问题
                drop_console: true, // 生产环境去除 console
                drop_debugger: true, // 生产环境去除 debugger
            },
            format: {
                comments: false, // 删除注释
            },
        },
    },
    optimizeDeps: {
        include: [
            'vue',
            'vue-router',
            'pinia',
            'axios',
            'echarts',
            'path-browserify',
            'element-plus/es',
            'element-plus/es/components/base/style/css',
            'element-plus/es/components/overview/style/css',
            'element-plus/es/components/main/style/css',
            'element-plus/es/components/header/style/css',
            'element-plus/es/components/aside/style/css',
            'element-plus/es/components/avatar/style/css',
            'element-plus/es/components/form/style/css',
            'element-plus/es/components/form-item/style/css',
            'element-plus/es/components/button/style/css',
            'element-plus/es/components/input/style/css',
            'element-plus/es/components/input-number/style/css',
            'element-plus/es/components/switch/style/css',
            'element-plus/es/components/upload/style/css',
            'element-plus/es/components/menu/style/css',
            'element-plus/es/components/col/style/css',
            'element-plus/es/components/icon/style/css',
            'element-plus/es/components/row/style/css',
            'element-plus/es/components/tag/style/css',
            'element-plus/es/components/dialog/style/css',
            'element-plus/es/components/loading/style/css',
            'element-plus/es/components/radio/style/css',
            'element-plus/es/components/radio-group/style/css',
            'element-plus/es/components/popover/style/css',
            'element-plus/es/components/scrollbar/style/css',
            'element-plus/es/components/tooltip/style/css',
            'element-plus/es/components/dropdown/style/css',
            'element-plus/es/components/dropdown-menu/style/css',
            'element-plus/es/components/dropdown-item/style/css',
            'element-plus/es/components/sub-menu/style/css',
            'element-plus/es/components/menu-item/style/css',
            'element-plus/es/components/divider/style/css',
            'element-plus/es/components/card/style/css',
            'element-plus/es/components/link/style/css',
            'element-plus/es/components/breadcrumb/style/css',
            'element-plus/es/components/breadcrumb-item/style/css',
            'element-plus/es/components/table/style/css',
            'element-plus/es/components/tree-select/style/css',
            'element-plus/es/components/table-column/style/css',
            'element-plus/es/components/select/style/css',
            'element-plus/es/components/option/style/css',
            'element-plus/es/components/pagination/style/css',
            'element-plus/es/components/tree/style/css',
            'element-plus/es/components/alert/style/css',
            'element-plus/es/components/radio-button/style/css',
            'element-plus/es/components/checkbox-group/style/css',
            'element-plus/es/components/checkbox/style/css',
            'element-plus/es/components/tabs/style/css',
            'element-plus/es/components/tab-pane/style/css',
            'element-plus/es/components/rate/style/css',
            'element-plus/es/components/date-picker/style/css',
            'element-plus/es/components/notification/style/css',
            'element-plus/es/components/image/style/css',
            'element-plus/es/components/statistic/style/css',
            'element-plus/es/components/watermark/style/css',
            'element-plus/es/components/config-provider/style/css',
            'element-plus/es/components/text/style/css',
            'element-plus/es/components/drawer/style/css',
            'element-plus/es/components/color-picker/style/css',
        ],
    },
})
