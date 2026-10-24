import { dynamicRoutes } from '@/router/dynamic.ts'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        name: 'login',
        component: async () => import('@/views/login/index.vue'),
    },
    {
        path: '/change-password',
        name: 'change-password',
        component: async () => import('@/views/change-password/index.vue'),
    },
    {
        path: '/',
        name: 'layout',
        component: async () => import('@/layout/index.vue'),
        redirect: '/monitor',
        children: [...dynamicRoutes],
    },
    {
        path: '/:pathMatch(.*)*',
        component: async () => import('@/components/Error/404.vue'),
    },
]

export default routes
