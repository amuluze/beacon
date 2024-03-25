import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
const router = createRouter({
    history: createWebHistory(), //可传参数，配置base路径，例如'/app'
    routes
})
export default router
