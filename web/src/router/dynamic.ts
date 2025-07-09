import type { RouteRecordRaw } from 'vue-router'

export const dynamicRoutes: RouteRecordRaw[] = [
    {
        path: '/overview',
        name: 'overview',
        component: async () => import('@/views/overview/index.vue'),
        meta: {
            title: 'menu.overview',
            icon: 'menu-home',
            show: true,
        },
    },
    {
        path: '/monitor',
        name: 'monitor',
        component: async () => import('@/views/monitor/index.vue'),
        meta: {
            title: 'menu.monitor',
            icon: 'menu-data',
            show: true,
        },
    },
    {
        path: '/container',
        name: 'container',
        component: async () => import('@/views/container/index.vue'),
        meta: {
            title: 'menu.container',
            icon: 'menu-multi',
            show: true,
        },
    },
    {
        path: '/setting',
        name: 'setting',
        component: async () => import('@/views/setting/index.vue'),
        meta: {
            title: 'menu.setting',
            icon: 'menu-system',
            show: true,
        },
    },
    {
        path: '/account',
        name: 'account',
        component: async () => import('@/views/account/index.vue'),
        meta: {
            title: 'menu.account',
            icon: 'menu-user',
            show: true,
        },
    },
    {
        path: '/about',
        name: 'about',
        component: async () => import('@/views/about/index.vue'),
        meta: {
            title: 'menu.about',
            icon: 'menu-about',
            show: true,
        },
    },
    {
        path: '/profile',
        name: 'profile',
        component: async () => import('@/views/profile/index.vue'),
        meta: {
            title: 'menu.profile',
            show: false,
        },
    },
]
