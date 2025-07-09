/**
 * @Author     : Amu
 * @Date       : 2024/12/10 16:47
 * @Description:
 */
import { createI18n } from 'vue-i18n'

import { getBrowserLanguage } from '@/utils'
import en from './modules/en'
import zh from './modules/zh'

const i18n = createI18n({
    allowComposition: true,
    legacy: false,
    locale: getBrowserLanguage(),
    messages: {
        zh,
        en,
    },
})

export default i18n
