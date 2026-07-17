// https://nuxt.com/docs/api/configuration/nuxt-config
import { defineNuxtConfig } from 'nuxt/config'

const serverProxyTarget = process.env.NUXT_SERVER_PROXY_TARGET || 'http://127.0.0.1:8000'
const securityHeaders = {
    'Content-Security-Policy': [
        `default-src 'self'`,
        `script-src 'self' 'unsafe-inline'`,
        `script-src-attr 'none'`,
        `style-src 'self' 'unsafe-inline'`,
        `img-src 'self' data:`,
        `font-src 'self' data:`,
        `connect-src 'self'`,
        `object-src 'none'`,
        `base-uri 'self'`,
        `frame-ancestors 'none'`,
        `form-action 'self'`,
        'upgrade-insecure-requests',
    ].join('; '),
    'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
    'Referrer-Policy': 'strict-origin-when-cross-origin',
    'Strict-Transport-Security': 'max-age=31536000; includeSubDomains',
    'X-Content-Type-Options': 'nosniff',
    'X-Frame-Options': 'DENY',
}

export default defineNuxtConfig({
    srcDir: 'src/',
    // 开启 SSR 并预渲染静态页，让搜索引擎抓取到首页/文档/changelog 正文
    ssr: true,
    devtools: { enabled: process.env.NODE_ENV !== 'production' },
    compatibilityDate: '2024-09-23',

    typescript: { typeCheck: true },

    // 注入全局样式
    css: ['@/styles/index.scss'],
    app: {
        baseURL: '/',
        head: {
            htmlAttrs: { lang: 'zh-CN' },
            title: 'Beacon - 轻量级主机及容器监控管理工具',
            meta: [
                { charset: 'utf-8' },
                { name: 'viewport', content: 'width=device-width, initial-scale=1' },
                { name: 'description', content: 'Beacon 是一款开源的轻量级主机监控及 Docker 容器管理工具，支持实时监控服务器资源使用情况，管理 Docker 容器、镜像和网络。' },
                { name: 'keywords', content: '主机监控,Docker管理,容器管理,服务器监控,开源监控工具' },
                { name: 'author', content: 'Beacon Team' },
                { name: 'theme-color', content: '#237a62' },
            ],
            link: [
                { rel: 'icon', type: 'image/svg+xml', href: '/beacon.svg' },
            ],
            script: [
                { src: '/theme-bootstrap.js', tagPosition: 'head' },
            ],
        },
    },

    modules: [
        '@nuxt/icon',
    ],

    vite: {
        css: {
            preprocessorOptions: {
                scss: {
                    additionalData: `@use "@/styles/bem.scss" as *; @use "@/styles/mixins.scss" as *;`,
                },
            },
        },
    },

    icon: {
        serverBundle: {
            collections: ['mdi'],
        },
    },
    // https://blog.csdn.net/qq_43231248/article/details/137127500
    runtimeConfig: {
        public: {
            // 运行时通过 NUXT_PUBLIC_BASE_URL 覆盖；默认同源代理。
            baseUrl: '',
        },
    },

    nitro: {
        routeRules: {
            '/**': { headers: securityHeaders },
            // 直接使用 /api/** 会导致图标加载失败
            '/api/v1/**': {
                proxy: `${serverProxyTarget}/api/v1/**`,
            },
            // 发布物（manager.sh / compose.yaml / version.json 等）由官网 Go server 静态托管
            '/release/**': {
                proxy: `${serverProxyTarget}/release/**`,
            },
            '/healthz': {
                proxy: `${serverProxyTarget}/healthz`,
            },
            // 静态营销/文档页构建时预渲染为完整 HTML，利于 SEO 且运行时零渲染开销
            '/': { prerender: true },
            '/document': { prerender: true },
            '/changelog': { prerender: true },
            '/wechat': { prerender: true },
            '/privacy': { prerender: true },
        },
    },
})
