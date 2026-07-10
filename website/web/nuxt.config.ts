// https://nuxt.com/docs/api/configuration/nuxt-config
import { defineNuxtConfig } from 'nuxt/config'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'

import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'

const serverProxyTarget = process.env.NUXT_SERVER_PROXY_TARGET || 'http://127.0.0.1:8000'

export default defineNuxtConfig({
    srcDir: 'src/',
    // 打包部署时需要置为 false
    ssr: false,
    devtools: { enabled: true },
    compatibilityDate: '2024-09-23',

    typescript: { typeCheck: true },

    // 注入全局样式
    css: ['@/styles/index.scss'],
    unocss: {
        nuxtLayers: true,
    },

    app: {
        baseURL: '/',
        head: {
            title: 'Beacon - 轻量级主机及容器监控管理工具',
            meta: [
                { charset: 'utf-8' },
                { name: 'viewport', content: 'width=device-width, initial-scale=1' },
                { name: 'description', content: 'Beacon 是一款开源的轻量级主机监控及 Docker 容器管理工具，支持实时监控服务器资源使用情况，管理 Docker 容器、镜像和网络。' },
                { name: 'keywords', content: '主机监控,Docker管理,容器管理,服务器监控,开源监控工具' },
                { name: 'author', content: 'Beacon Team' },
                // Open Graph tags
                { property: 'og:title', content: 'Beacon - 轻量级主机及容器监控管理工具' },
                { property: 'og:description', content: 'Beacon 是一款开源的轻量级主机监控及 Docker 容器管理工具' },
                { property: 'og:type', content: 'website' },
                { property: 'og:url', content: 'https://official.beacon.amuluze.com' },
                { property: 'og:image', content: '/images/beacon.png' },
            ],
            link: [
                { rel: 'icon', type: 'image/x-icon', href: 'favicon.ico' },
                { rel: 'canonical', href: 'https://official.beacon.amuluze.com' },
            ],
        },
    },

    modules: [
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        '@unocss/nuxt',
        '@element-plus/nuxt',
        '@nuxt/icon',
        '@nuxt/image',
    ],

    build: {
        // 持久化插件配置，必须
        transpile: ['element-plus/nuxt', 'pinia-plugin-persistedstate/nuxt'],
    },

    vite: {
        css: {
            preprocessorOptions: {
                scss: {
                    additionalData: `@use "@/styles/bem.scss" as *; @use "@/styles/mixins.scss" as *;`,
                },
            },
        },
        plugins: [
            // 自动导入
            Components({
                resolvers: [
                    ElementPlusResolver(),
                    IconsResolver({
                        enabledCollections: ['mdi'],
                    }),
                ],
            }),
            Icons({
                // 自动安装图标库
                autoInstall: true,
                compiler: 'vue3',
            }),
        ],
    },

    plugins: [],
    icon: {
        serverBundle: {
            collections: ['uil', 'mdi'],
        },
    },
    // https://blog.csdn.net/qq_43231248/article/details/137127500
    runtimeConfig: {
        public: {
            baseUrl: process.env.NUXT_BASE_URL,
        },
    },

    nitro: {
        routeRules: {
            // 直接使用 /api/** 会导致图标加载失败
            '/api/v1/**': {
                proxy: `${serverProxyTarget}/api/v1/**`,
            },
            '/download/**': {
                proxy: `${serverProxyTarget}/download/**`,
            },
        },
    },
})
