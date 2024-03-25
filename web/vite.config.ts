import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

import Icons from 'unplugin-icons/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import IconsResolver from 'unplugin-icons/resolver'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        AutoImport({
            // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
            imports: ['vue', 'vue-router', 'pinia'],
            eslintrc: {
                enabled: false, // 是否自动生成 eslint 规则，建议生成之后设置 false，手动维护
                filepath: './.eslintrc.cjs-auto-import.json', // 指定自动导入函数 eslint 规则的文件路径
                globalsPropValue: true
            },
            resolvers: [
                // 自动导入 Element Plus 相关函数，如：ElMessage(带样式)
                ElementPlusResolver(),
                IconsResolver({
                    enabledCollections: ['ep']
                })
            ],
            vueTemplate: true, // 是否在 vue 模板中自动导入
            dts: resolve(resolve(__dirname, 'types'), 'auto-imports.d.ts') // 指定自动导入函数TS类型声明文件路径
        }),
        Components({
            resolvers: [
                // 自动导入 Element Plus 组件
                ElementPlusResolver(),
                // 自动导入图标组件
                IconsResolver({
                    // @iconify-json/ep 是 Element Plus 的图标库
                    enabledCollections: ['ep']
                })
            ],
            dts: resolve(resolve(__dirname, 'types'), 'components.d.ts') // 指定自动导入组件TS类型声明文件路径
        }),
        Icons({
            // 自动安装图标库
            autoInstall: true
        }),
        createSvgIconsPlugin({
            // 指定需要缓存的图标文件夹
            iconDirs: [resolve(process.cwd(), 'src/assets/icons'), resolve(process.cwd(), 'src/assets/error')],
            // iconDirs: [path.resolve(process.cwd(), 'src/assets/icons')],
            // 指定symbolId格式
            symbolId: 'icon-[dir]-[name]'
        })
    ],
    resolve: {
        alias: {
            '@': resolve(__dirname, 'src')
        },
        // 导入时，忽略后缀
        extensions: ['.js', '.ts', '.jsx', '.tsx', '.json', '.vue', '.mjs']
    },
    server: {
        host: '0.0.0.0',
        port: 3000,
        open: true,
        proxy: {
            '/api': {
                target: 'http://localhost:8000',
                changeOrigin: true,
                rewrite: (path: string) => path.replace(/^\/api/, '')
            }
        }
    },
    css: {
        preprocessorOptions: {
            scss: {
                additionalData: '@import "@/styles/bem.scss";'
            }
        }
    }
})
