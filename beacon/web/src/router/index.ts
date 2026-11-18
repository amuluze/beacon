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
            store.user.setUserInfo(
                data.username,
                data.status,
                data.is_admin,
                data.must_change_password,
                data.read_only,
                data.no_auth,
            )
        }
        catch {
            store.user.clearSession()
        }
        finally {
            store.user.markSessionChecked()
        }
    }

    const isAuthenticated = store.user.authenticated
    const isAnonymous = store.user.userInfo.noAuth === true
    const mustChangePassword = store.user.userInfo.mustChangePassword === true

    if (!isAuthenticated && to.name !== 'login') {
        // 未登录，to 非登录页  --> 跳转登录页
        return next({ name: 'login' })
    }
    if (!isAuthenticated && to.name === 'login') {
        // 未登录，to 登录页
        return next()
    }
    // 匿名只读会话可以主动进入登录页；登录成功后由登录页切换到正常会话。
    if (isAuthenticated && isAnonymous && to.name === 'login') {
        return next()
    }
    if (isAuthenticated && isAnonymous && to.name === 'change-password') {
        return next({ name: 'monitor' })
    }
    if (isAuthenticated && mustChangePassword && to.name !== 'change-password') {
        return next({ name: 'change-password' })
    }
    if (isAuthenticated && !mustChangePassword && to.name === 'change-password') {
        return next({ name: 'monitor' })
    }
    if (isAuthenticated && !isAnonymous && to.name === 'login') {
        // 已登录用户访问登录页时回首页
        return next({ name: 'monitor' })
    }
    return next()
})

export default router
