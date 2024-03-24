import { RouteRecordRaw } from 'vue-router';

export const dynamicRoutes: RouteRecordRaw[] = [
    {
        path: '/overview',
        name: 'overview',
        component: () => import('@/views/overview/index.vue'),
        meta: {
            title: '总览',
            icon: 'homepage',
        }
    },
    {
        path: '/container',
        name: 'container',
        redirect: '/container/container',
        meta: {
            title: '容器管理',
            icon: 'system',
        },
        children: [
            {
                path: '/container/container',
                name: 'container',
                component: () => import('@/views/container/container/index.vue'), //路由懒加载
                meta: {
                    title: '容器',
                    icon: 'system',
                }
            },
            {
                path: '/container/image',
                name: 'image',
                component: () => import('@/views/container/image/index.vue'), //路由懒加载
                meta: {
                    title: '镜像',
                    icon: 'system',
                }
            }
        ]
    },
    {
        path: '/host',
        name: 'host',
        // component: () => import('@/views/host/index.vue'), //路由懒加载
        redirect: '/host/monitor',
        meta: {
            title: '主机管理',
            icon: 'system',
        },
        children: [
            {
                path: '/host/monitor',
                name: 'monitor',
                component: () => import('@/views/host/monitor/index.vue'), //路由懒加载
                meta: {
                    title: '监控',
                    icon: 'system',
                }
            }
        ]
    },
];
