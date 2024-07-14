import { dynamicRoutes } from '@/router/dynamic.ts'
import { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/app/login',
        name: 'login',
        component: () => import('@/views/login/index.vue'), //路由懒加载
        meta: {
            title: '登录'
        }
    },
    {
        // The main page
        path: '/',
        name: 'layout',
        component: () => import('@/layout/index.vue'),
        redirect: '/overview',
        children: [...dynamicRoutes]
    },
    // Resolve refresh page, route warnings
    {
        path: '/:pathMatch(.*)*',
        component: () => import('@/components/Error/404.vue')
    }
]

export default routes
