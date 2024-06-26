import useStore from '@/store'
import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
const router = createRouter({
    history: createWebHistory('/app'), //可传参数，配置base路径，例如'/app'
    routes
})

router.beforeEach(async (to, _, next) => {
    const store = useStore()
    console.log('to', to)
    if (store.user.token === '' && to.name !== 'login') {
        // 未登录，to 非登录页  --> 跳转登录页
        next({ name: 'login' })
    } else if (store.user.token === '' && to.name === 'login') {
        // 未登录，to 登录页
        next()
    } else if (store.user.token !== '' && to.name === 'login') {
        // 跳转首页
        next({ name: 'overview' })
    } else {
        next()
    }
    // next()
})
export default router
