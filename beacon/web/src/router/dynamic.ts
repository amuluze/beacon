import type { RouteRecordRaw } from 'vue-router'

export const dynamicRoutes: RouteRecordRaw[] = [
    {
        path: '/monitor',
        name: 'monitor',
        component: async () => import('@/views/monitor/index.vue'),
        meta: {
            title: 'menu.monitor',
            icon: 'activity',
            show: true,
        },
    },
    {
        path: '/container',
        name: 'container',
        component: async () => import('@/views/container/index.vue'),
        meta: {
            title: 'menu.container',
            icon: 'package',
            show: true,
        },
    },
    {
        path: '/setting',
        name: 'setting',
        component: async () => import('@/views/setting/index.vue'),
        meta: {
            title: 'menu.setting',
            icon: 'settings',
            show: true,
        },
    },
    // 设置工作区改为一级 Tab 后，旧深链仍指向同一工作区的对应 Tab。
    {
        path: '/setting/alarm',
        redirect: { path: '/setting', query: { tab: 'config' } },
    },
    {
        path: '/setting/container',
        redirect: { path: '/setting', query: { tab: 'config' } },
    },
    {
        path: '/account',
        name: 'account',
        redirect: { path: '/setting', query: { tab: 'users' } },
        meta: {
            title: 'menu.account',
            icon: 'users',
            show: false,
        },
        children: [
            {
                path: '/account/user',
                name: 'userManager',
                redirect: { path: '/setting', query: { tab: 'users' } },
                meta: {
                    title: 'menu.userManager',
                },
            },
            {
                path: '/account/role',
                name: 'roleManager',
                redirect: { path: '/setting', query: { tab: 'users' } },
                meta: {
                    title: 'menu.roleManager',
                },
            },
            {
                path: '/account/api',
                name: 'apiManager',
                redirect: { path: '/setting', query: { tab: 'users' } },
                meta: {
                    title: 'menu.apiManager',
                },
            },
        ],
    },
]
