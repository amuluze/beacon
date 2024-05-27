import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { defineConfig } from 'vite'

import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

// https://vitejs.dev/config/
export default defineConfig({
    base: './', // 设置打包路径
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
    },
    // 构建配置
    build: {
        chunkSizeWarningLimit: 2000, // 消除打包大小超过500kb警告
        minify: 'terser', // Vite 2.6.x 以上需要配置 minify: "terser", terserOptions 才能生效
        terserOptions: {
            compress: {
                keep_infinity: true, // 防止 Infinity 被压缩成 1/0，这可能会导致 Chrome 上的性能问题
                drop_console: true, // 生产环境去除 console
                drop_debugger: true // 生产环境去除 debugger
            },
            format: {
                comments: false // 删除注释
            }
        },
        rollupOptions: {
            output: {
                // 用于从入口点创建的块的打包输出格式[name]表示文件名,[hash]表示该文件内容hash值
                entryFileNames: 'js/[name].[hash].js',
                // 用于命名代码拆分时创建的共享块的输出命名
                chunkFileNames: 'js/[name].[hash].js',
                // 用于输出静态资源的命名，[ext]表示文件扩展名
                assetFileNames: (assetInfo: any) => {
                    const info = assetInfo.name.split('.')
                    let extType = info[info.length - 1]
                    // console.log('文件信息', assetInfo.name)
                    if (/\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/i.test(assetInfo.name)) {
                        extType = 'media'
                    } else if (/\.(png|jpe?g|gif|svg)(\?.*)?$/.test(assetInfo.name)) {
                        extType = 'img'
                    } else if (/\.(woff2?|eot|ttf|otf)(\?.*)?$/i.test(assetInfo.name)) {
                        extType = 'fonts'
                    }
                    return `${extType}/[name].[hash].[ext]`
                }
            }
        }
    }
})
