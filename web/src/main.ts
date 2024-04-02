import { setupStore } from '@/store/setup'
import '@/styles/index.scss'
import 'virtual:svg-icons-register'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)
setupStore(app)
app.use(router).mount('#app') //注册路由
