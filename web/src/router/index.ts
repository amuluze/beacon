import { createRouter, createWebHashHistory } from 'vue-router'
import routes from './routes'
import useStore from '@/store'

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

router.beforeEach(async (to, _, next) => {
    const store = useStore()
    if (store.user.token === '' && to.name !== 'login') {
        // 未登录，to 非登录页  --> 跳转登录页
        next({ name: 'login' })
    }
    if (store.user.token === '' && to.name === 'login') {
        // 未登录，to 登录页
        next()
    }
    if (store.user.token !== '' && to.name === 'login') {
        // 跳转首页
        next({ name: 'overview' })
    }
    next()
})

export default router
