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
            icon: 'container',
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
            show: false,
            activeMenu: '/docker'
        }
    },
    {
        path: '/alarm',
        name: 'alarm',
        redirect: '/alarm/mail',
        meta: {
            title: '告警',
            icon: 'alarm',
            show: true,
            activeMenu: '/mail'
        }
    },
    {
        path: '/alarm/mail',
        name: 'mail',
        component: () => import('@/views/alarm/mail/index.vue'),
        meta: {
            title: '告警配置',
            show: false,
            activeMenu: '/alarm'
        }
    },
    {
        path: '/alarm/threshold',
        name: 'threshold',
        component: () => import('@/views/alarm/threshold/index.vue'),
        meta: {
            title: '告警阈值',
            show: false,
            activeMenu: '/alarm'
        }
    },
    {
        path: '/account',
        name: 'account',
        redirect: '/account/user',
        meta: {
            title: '用户',
            icon: 'user',
            show: true
        }
    },
    {
        path: '/account/user',
        name: 'user',
        component: () => import('@/views/account/user/index.vue'), //路由懒加载
        meta: {
            title: '用户',
            icon: 'user',
            show: false,
            activeMenu: '/account'
        }
    },
    {
        path: '/account/role',
        name: 'role',
        component: () => import('@/views/account/role/index.vue'), //路由懒加载
        meta: {
            title: '角色',
            icon: 'user',
            show: false,
            activeMenu: '/account'
        }
    },
    {
        path: '/account/api',
        name: 'api',
        component: () => import('@/views/account/api/index.vue'), //路由懒加载
        meta: {
            title: '接口',
            icon: 'user',
            show: false,
            activeMenu: '/account'
        }
    },
    {
        path: '/audit',
        name: 'audit',
        redirect: '/audit/operator',
        meta: {
            title: '日志',
            icon: 'audit',
            show: true
        }
    },
    {
        path: '/audit/operator',
        name: 'operator',
        component: () => import('@/views/audit/operator/index.vue'), //路由懒加载
        meta: {
            title: '操作日志',
            show: false
        }
    },
    {
        path: '/audit/system',
        name: 'system',
        component: () => import('@/views/audit/system/index.vue'), //路由懒加载
        meta: {
            title: '系统日志',
            show: false
        }
    },
    {
        path: '/profile',
        name: 'profile',
        component: () => import('@/views/profile/index.vue'), //路由懒加载
        meta: {
            title: '个人中心',
            icon: 'settings',
            show: false
        }
    }
]
