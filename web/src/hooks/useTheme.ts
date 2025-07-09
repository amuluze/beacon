/**
 * @Author     : Amu
 * @Date       : 2024/11/14 00:00
 * @Description:
 */
import { themeConfig } from '@/config/theme.ts'
import useStore from '@/store'

function useTheme() {
    const store = useStore()

    // 修改主题颜色
    const changePrimary = () => {
        let type: Theme.ThemeType = 'light'
        if (store.theme.dark)
            type = 'dark'
        const theme = themeConfig[type]

        for (const [key, value] of Object.entries(theme)) {
            console.log(key, value)
            document.documentElement.style.setProperty(key, value)
        }
    }

    // 切换暗黑模式
    const switchDark = () => {
        const html = document.documentElement
        console.log('======', !store.theme.dark)
        store.theme.setDark(!store.theme.dark)
        if (store.theme.dark)
            html.setAttribute('class', 'dark')
        else html.setAttribute('class', '')
        console.log('is dark: ', store.theme.dark)
        changePrimary()
    }

    // 初始化主题#e9effd
    const initTheme = () => {
        const html = document.documentElement
        console.log('is dark: ', store.theme.dark)
        if (store.theme.dark)
            html.setAttribute('class', 'dark')
        else html.setAttribute('class', 'light')
        changePrimary()
    }

    return { initTheme, switchDark, changePrimary }
}

export default useTheme
