import { RouteRecordRaw } from 'vue-router'

export const dynamicRoutes: RouteRecordRaw[] = [
    {
        path: '/overview',
        name: 'overview',
        component: () => import('@/views/overview/index.vue'),
        meta: {
            title: '总览',
            icon: 'homepage',
            show: true
        }
    },
    {
        path: '/host',
        name: 'host',
        redirect: '/host/monitor',
        meta: {
            title: '主机',
            icon: 'monitor',
            show: true
        }
    },
    {
        path: '/host/monitor',
        name: 'monitor',
        component: () => import('@/views/host/monitor/index.vue'), //路由懒加载
        meta: {
            title: '监控',
            icon: 'monitor',
            show: false,
            activeMenu: '/host'
        }
    },
    {
        path: '/host/file',
        name: 'file',
        component: () => import('@/views/host/file/index.vue'), //路由懒加载
        meta: {
            title: '文件',
            icon: 'setting',
            show: false,
            activeMenu: '/host'
        }
    },
    {
        path: '/host/terminal',
        name: 'terminal',
        component: () => import('@/views/host/terminal/index.vue'), //路由懒加载
        meta: {
            title: '终端',
            icon: 'setting',
            show: false,
            activeMenu: '/host'
        }
    },
    {
        path: '/host/settings',
        name: 'host-settings',
        component: () => import('@/views/host/settings/index.vue'), //路由懒加载
        meta: {
            title: '设置',
            icon: 'setting',
            show: false,
            activeMenu: '/host'
        }
    },
    {
        path: '/docker',
        name: 'docker',
        redirect: '/docker/container',
        meta: {
            title: '容器',
            icon: 'setting',
            show: true,
            activeMenu: '/container'
        }
    },
    {
        path: '/docker/container',
        name: 'container',
        component: () => import('@/views/container/container/index.vue'), //路由懒加载
        meta: {
            title: '容器',
            icon: 'setting',
            show: false,
            activeMenu: '/docker'
        }
    },
    {
        path: '/docker/image',
        name: 'image',
        component: () => import('@/views/container/image/index.vue'), //路由懒加载
        meta: {
            title: '镜像',
            icon: 'lab',
            show: false,
            activeMenu: '/docker'
        }
    },
    {
        path: '/docker/network',
        name: 'network',
        component: () => import('@/views/container/network/index.vue'), //路由懒加载
        meta: {
            title: '网络',
            icon: 'lab',
            show: false,
            activeMenu: '/docker'
        }
    },
    {
        path: '/docker/settings',
        name: 'docker-settings',
        component: () => import('@/views/container/settings/index.vue'), //路由懒加载
        meta: {
            title: '设置',
            icon: 'lab',
            show: false,
            activeMenu: '/docker'
        }
    },
    {
        path: '/audit',
        name: 'audit',
        component: () => import('@/views/audit/index.vue'), //路由懒加载
        meta: {
            title: '审计',
            icon: 'audit',
            show: true
        }
    }
    // {
    //     path: '/system',
    //     name: 'system',
    //     component: () => import('@/views/system/index.vue'), //路由懒加载
    //     meta: {
    //         title: '系统',
    //         icon: 'settings',
    //         show: true
    //     }
    // }
]
