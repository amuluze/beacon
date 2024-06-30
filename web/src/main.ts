import { setupStore } from '@/store/setup'
import '@/styles/index.scss'
import 'virtual:svg-icons-register'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import dayjs from 'dayjs'

const app = createApp(App)
setupStore(app)
app.use(router)
// 全局使用 dayjs
app.config.globalProperties.$dayjs = dayjs
app.mount('#app') //注册路由
