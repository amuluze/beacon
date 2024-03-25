import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { setupStore } from '@/store/setup'
import 'virtual:svg-icons-register'
import '@/styles/index.scss'

const app = createApp(App)
setupStore(app)
app.use(router).mount('#app') //注册路由
