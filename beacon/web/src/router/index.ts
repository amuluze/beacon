import { createRouter, createWebHashHistory } from 'vue-router'
import routes from './routes'
import useStore from '@/store'
import { getUserInfo } from '@/api/auth'

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

router.beforeEach(async (to, _, next) => {
    const store = useStore()
    if (!store.user.sessionChecked) {
        try {
            const { data } = await getUserInfo()
            store.user.setUserInfo(data.username, data.status, data.is_admin, data.must_change_password)
        }
        catch {
            store.user.clearSession()
        }
        finally {
            store.user.markSessionChecked()
        }
    }

    const isAuthenticated = store.user.authenticated
    const mustChangePassword = store.user.userInfo.mustChangePassword === true

    if (!isAuthenticated && to.name !== 'login') {
        // 未登录，to 非登录页  --> 跳转登录页
        return next({ name: 'login' })
    }
    if (!isAuthenticated && to.name === 'login') {
        // 未登录，to 登录页
        return next()
    }
    if (isAuthenticated && mustChangePassword && to.name !== 'change-password') {
        return next({ name: 'change-password' })
    }
    if (isAuthenticated && !mustChangePassword && to.name === 'change-password') {
        return next({ name: 'monitor' })
    }
    if (isAuthenticated && to.name === 'login') {
        // 跳转首页
        return next({ name: 'monitor' })
    }
    return next()
})

export default router
