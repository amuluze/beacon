/**
 * @Author     : Amu
 * @Date       : 2025/02/13 00:34:22
 * @Description:
 */

import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

export default defineNuxtPlugin((nuxtApp: any) => {
    nuxtApp.$pinia.use(piniaPluginPersistedstate)
})
