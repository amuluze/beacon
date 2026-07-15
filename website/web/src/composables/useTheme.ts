export type ThemeMode = 'light' | 'dark'

const STORAGE_KEY = 'beacon-theme'
const state = ref<ThemeMode>('light')

function readStored(): ThemeMode {
    if (typeof window === 'undefined')
        return 'light'
    const stored = window.localStorage.getItem(STORAGE_KEY) as ThemeMode | null
    if (stored === 'light' || stored === 'dark')
        return stored
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function apply(mode: ThemeMode) {
    if (typeof document === 'undefined')
        return
    document.documentElement.classList.toggle('dark', mode === 'dark')
}

export function initTheme() {
    const mode = readStored()
    state.value = mode
    apply(mode)
}

export function useTheme() {
    const mode = computed(() => state.value)
    const isDark = computed(() => state.value === 'dark')

    const setTheme = (next: ThemeMode) => {
        state.value = next
        apply(next)
        if (typeof window !== 'undefined')
            window.localStorage.setItem(STORAGE_KEY, next)
    }

    const toggleTheme = () => setTheme(isDark.value ? 'light' : 'dark')

    return { mode, isDark, setTheme, toggleTheme }
}
