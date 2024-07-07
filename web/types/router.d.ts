/**
 * 扩展 vue-router RouteRecordRaw RouteMeta
 */

import 'vue-router'

declare module 'vue-router' {
    interface RouteMeta {
        title?: string
        show?: boolean
    }
}
