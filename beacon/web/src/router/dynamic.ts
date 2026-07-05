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
        redirect: '/monitor/host',
        meta: {
            title: 'menu.monitor',
            icon: 'menu-data',
            show: true,
        },
        children: [
            {
                path: '/monitor/host',
                name: 'hostMonitor',
                component: async () => import('@/views/monitor/host/index.vue'),
                meta: {
                    title: 'menu.hostMonitor',
                },
            },
            {
                path: '/monitor/container',
                name: 'containerMonitor',
                component: async () => import('@/views/monitor/container/index.vue'),
                meta: {
                    title: 'menu.containerMonitor',
                },
            },
        ],
    },
    {
        path: '/container',
        name: 'container',
        redirect: '/container/container',
        meta: {
            title: 'menu.container',
            icon: 'menu-multi',
            show: true,
        },
        children: [
            {
                path: '/container/container',
                name: 'containerManager',
                component: async () => import('@/views/container/container/index.vue'),
                meta: {
                    title: 'menu.containerManager',
                },
            },
            {
                path: '/container/image',
                name: 'imageManager',
                component: async () => import('@/views/container/image/index.vue'),
                meta: {
                    title: 'menu.imageManager',
                },
            },
            {
                path: '/container/network',
                name: 'networkManager',
                component: async () => import('@/views/container/network/index.vue'),
                meta: {
                    title: 'menu.networkManager',
                },
            },
        ],
    },
    {
        path: '/audit',
        name: 'audit',
        redirect: '/audit/operator',
        meta: {
            title: 'menu.audit',
            icon: 'menu-result',
            show: true,
        },
        children: [
            {
                path: '/audit/operator',
                name: 'operatorLog',
                component: async () => import('@/views/audit/operator/index.vue'),
                meta: {
                    title: 'menu.operatorLog',
                },
            },
            {
                path: '/audit/system',
                name: 'systemLog',
                component: async () => import('@/views/audit/system/index.vue'),
                meta: {
                    title: 'menu.systemLog',
                },
            },
        ],
    },
    {
        path: '/setting',
        name: 'setting',
        redirect: '/setting/alarm',
        meta: {
            title: 'menu.setting',
            icon: 'menu-system',
            show: true,
        },
        children: [
            {
                path: '/setting/alarm',
                name: 'alarm',
                component: async () => import('@/views/setting/alarm/index.vue'),
                meta: {
                    title: 'menu.alarmSetting',
                },
            },
            {
                path: '/setting/host',
                name: 'host',
                component: async () => import('@/views/setting/host/index.vue'),
                meta: {
                    title: 'menu.systemSetting',
                },
            },
            {
                path: '/setting/container',
                name: 'docker',
                component: async () => import('@/views/setting/docker/index.vue'),
                meta: {
                    title: 'menu.systemDocker',
                },
            },
        ],
    },
    {
        path: '/account',
        name: 'account',
        redirect: '/account/user',
        meta: {
            title: 'menu.account',
            icon: 'menu-user',
            show: true,
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
