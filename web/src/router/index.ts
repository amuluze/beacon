import useStore from '@/store'
import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
const router = createRouter({
    history: createWebHistory(), //可传参数，配置base路径，例如'/app'
    routes
})

router.beforeEach(async (to, _, next) => {
    const store = useStore()
    // 未登录，to 登录页
    if (store.user.token === '' && to.name === 'login') {
        next()
    }

    // 未登录，to 非登录页
    if (store.user.token === '' && to.name !== 'login') {
        // 跳转登录页
        next({ name: 'login' })
    }

    // 登录
    if (store.user.token !== '' && to.name === 'login') {
        // to 登录页
        if (to.name === 'login') {
            // 跳转首页
            next({ name: 'home' })
        }
        next()
    }
    next()
})
export default router
