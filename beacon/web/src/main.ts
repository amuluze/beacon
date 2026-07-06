import { createApp } from 'vue'
import App from './App.vue'

import { setupStore } from '@/store/setup'

import '@/styles/index.scss'
import 'virtual:svg-icons-register'

import 'virtual:uno.css'

// element+ 默认主题
import 'element-plus/dist/index.css'
// element+ 内置黑暗主题
import 'element-plus/theme-chalk/dark/css-vars.css'

import router from './router'

import registerDirective from '@/directive'
import i18n from '@/languages/index'

const app = createApp(App)
// 注册状态管理
setupStore(app)
// 注册指令
app.use(registerDirective)
// 注册路由
app.use(router)
// 注册i18n
app.use(i18n)
app.mount('#app')
