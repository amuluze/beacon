import { beforeEach, describe, expect, it } from 'vitest'
import { initTheme, useTheme } from './useTheme'

const STORAGE_KEY = 'beacon-theme'

describe('useTheme', () => {
    beforeEach(() => {
        window.localStorage.clear()
        document.documentElement.classList.remove('dark', 'light')
    })

    it('initTheme 读取 localStorage 中的持久化值', () => {
        window.localStorage.setItem(STORAGE_KEY, 'dark')
        initTheme()
        const { isDark } = useTheme()
        expect(isDark.value).toBe(true)
        expect(document.documentElement.classList.contains('dark')).toBe(true)
    })

    it('initTheme 在无存储时回退到 prefers-color-scheme', () => {
        // matchMedia 默认不匹配 dark → 浅色
        initTheme()
        const { isDark } = useTheme()
        expect(isDark.value).toBe(false)
    })

    it('setTheme 写入 localStorage 并切换 html class', () => {
        initTheme()
        const { setTheme, isDark } = useTheme()

        setTheme('dark')
        expect(isDark.value).toBe(true)
        expect(window.localStorage.getItem(STORAGE_KEY)).toBe('dark')
        expect(document.documentElement.classList.contains('dark')).toBe(true)

        setTheme('light')
        expect(isDark.value).toBe(false)
        expect(window.localStorage.getItem(STORAGE_KEY)).toBe('light')
        expect(document.documentElement.classList.contains('dark')).toBe(false)
    })

    it('toggleTheme 在两种模式间来回切换', () => {
        initTheme()
        const { toggleTheme, isDark } = useTheme()

        expect(isDark.value).toBe(false)
        toggleTheme()
        expect(isDark.value).toBe(true)
        toggleTheme()
        expect(isDark.value).toBe(false)
    })
})
