import { RouteRecordRaw } from 'vue-router'
import { dynamicRoutes } from '@/router/dynamic.ts'

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        component: () => import('@/views/login/index.vue') //路由懒加载
    },
    {
        // The main page
        path: '/',
        name: 'layout',
        component: () => import('@/layout/index.vue'),
        redirect: '/container',
        children: [...dynamicRoutes]
    },
    // Resolve refresh page, route warnings
    {
        path: '/:pathMatch(.*)*',
        component: () => import('@/components/Error/404.vue')
    }
]

export default routes
