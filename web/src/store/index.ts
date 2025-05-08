import { useAppStore } from '@/store/modules/app.ts'
import { useEChartsStore } from '@/store/modules/echarts.ts'
import { usePermissionStore } from '@/store/modules/permission'
import { useThemeStore } from '@/store/modules/theme'
import { useUserStore } from '@/store/modules/user'

// 注册子模块
function useStore() {
    return {
        user: useUserStore(),
        theme: useThemeStore(),
        app: useAppStore(),
        echarts: useEChartsStore(),
        permissions: usePermissionStore(),
    }
}

export default useStore
