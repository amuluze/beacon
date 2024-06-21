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
        path: '/host/monitor',
        name: 'monitor',
        component: () => import('@/views/host/monitor/index.vue'), //路由懒加载
        meta: {
            title: '主机',
            icon: 'monitor',
            show: true
        }
    },
    {
        path: '/container/container',
        name: 'container',
        component: () => import('@/views/container/container/index.vue'), //路由懒加载
        meta: {
            title: '容器',
            icon: 'setting',
            show: true
        }
    },
    {
        path: '/container/image',
        name: 'image',
        component: () => import('@/views/container/image/index.vue'), //路由懒加载
        meta: {
            title: '镜像',
            icon: 'lab',
            show: false
        }
    },
    {
        path: '/container/network',
        name: 'network',
        component: () => import('@/views/container/network/index.vue'), //路由懒加载
        meta: {
            title: '网络',
            icon: 'lab',
            show: false
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
]
