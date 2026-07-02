import { useAppStore } from '@/store/modules/app.ts'
import { useEChartsStore } from '@/store/modules/echarts.ts'
import { useAgentStore } from '@/store/modules/agent'
import { usePermissionStore } from '@/store/modules/permission'
import { useThemeStore } from '@/store/modules/theme'
import { useUserStore } from '@/store/modules/user'

// 注册子模块
function useStore() {
    return {
        user: useUserStore(),
        theme: useThemeStore(),
        app: useAppStore(),
        agent: useAgentStore(),
        echarts: useEChartsStore(),
        permissions: usePermissionStore(),
    }
}

export default useStore
