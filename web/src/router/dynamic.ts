import { RouteRecordRaw } from 'vue-router'

export const dynamicRoutes: RouteRecordRaw[] = [
    {
        path: '/overview',
        name: 'overview',
        component: () => import('@/views/overview/index.vue'),
        meta: {
            title: '总览',
            icon: 'homepage'
        }
    },
    {
        path: '/container/container',
        name: 'container',
        component: () => import('@/views/container/container/index.vue'), //路由懒加载
        meta: {
            title: '容器',
            icon: 'setting'
        }
    },
    {
        path: '/container/image',
        name: 'image',
        component: () => import('@/views/container/image/index.vue'), //路由懒加载
        meta: {
            title: '镜像',
            icon: 'lab'
        }
    },
    {
        path: '/host/monitor',
        name: 'monitor',
        component: () => import('@/views/host/monitor/index.vue'), //路由懒加载
        meta: {
            title: '监控',
            icon: 'monitor'
        }
    },
    {
        path: '/audit',
        name: 'audit',
        component: () => import('@/views/audit/index.vue'), //路由懒加载
        meta: {
            title: '审计',
            icon: 'audit'
        }
    }
]
