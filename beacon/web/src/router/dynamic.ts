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
    {
        path: '/setting/alarm',
        name: 'alarm',
        component: async () => import('@/views/setting/alarm/index.vue'),
        meta: { title: 'menu.alarmSetting', show: false },
    },
    {
        path: '/setting/host',
        name: 'host',
        component: async () => import('@/views/setting/host/index.vue'),
        meta: { title: 'menu.systemSetting', show: false },
    },
    {
        path: '/setting/container',
        name: 'docker',
        component: async () => import('@/views/setting/docker/index.vue'),
        meta: { title: 'menu.systemDocker', show: false },
    },
    {
        path: '/account',
        name: 'account',
        redirect: '/account/user',
        meta: {
            title: 'menu.account',
            icon: 'users',
            show: false,
        },
        children: [
            {
                path: '/account/user',
                name: 'userManager',
                component: async () => import('@/views/account/user/index.vue'),
                meta: {
                    title: 'menu.userManager',
                },
            },
            {
                path: '/account/role',
                name: 'roleManager',
                component: async () => import('@/views/account/role/index.vue'),
                meta: {
                    title: 'menu.roleManager',
                },
            },
            {
                path: '/account/api',
                name: 'apiManager',
                component: async () => import('@/views/account/api/index.vue'),
                meta: {
                    title: 'menu.apiManager',
                },
            },
        ],
    },
    {
        path: '/terminal',
        name: 'terminal',
        component: async () => import('@/views/terminal/index.vue'),
        meta: {
            title: 'menu.terminal',
            icon: 'terminal',
            show: true,
        },
    },
]
