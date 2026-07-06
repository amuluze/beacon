/**
 * @Author     : Amu
 * @Date       : 2025/2/12 01:29
 * @Description:
 */
import SvgIcon from '@/components/SvgIcon/index.vue'
import 'virtual:svg-icons-register'

// 注册脚本，不进行nuxt.config.ts里面的配置，会报错
export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.component('svg-icon', SvgIcon)
})
