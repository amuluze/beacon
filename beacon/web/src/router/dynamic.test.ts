import type { RouteRecordRaw } from 'vue-router'
import { describe, expect, it } from 'vitest'
import { dynamicRoutes } from './dynamic'
import appRoutes from './routes'

function flattenRoutes(routes: RouteRecordRaw[]): RouteRecordRaw[] {
    return routes.flatMap(route => [route, ...flattenRoutes(route.children ?? [])])
}

describe('dynamic workspace routes', () => {
    const allRoutes = flattenRoutes(dynamicRoutes)

    it('uses the complete monitor workspace as the authenticated landing page', () => {
        const layout = flattenRoutes(appRoutes).find(route => route.name === 'layout')
        expect(layout?.redirect).toBe('/monitor')
    })

    it.each([
        ['/monitor', 'monitor'],
        ['/container', 'container'],
        ['/setting', 'setting'],
    ])('renders %s as a complete workspace', (path, name) => {
        const route = allRoutes.find(item => item.path === path)

        expect(route?.name).toBe(name)
        expect(route?.component).toBeTypeOf('function')
        expect(route?.redirect).toBeUndefined()
        expect(route?.meta?.show).toBe(true)
    })

    it.each([
        ['/setting/alarm', 'alarm'],
        ['/setting/host', 'host'],
        ['/setting/container', 'docker'],
    ])('keeps the legacy deep link %s resolvable', (path, name) => {
        const route = allRoutes.find(item => item.path === path)

        expect(route?.name).toBe(name)
        expect(route?.component).toBeTypeOf('function')
    })

    it('does not expose the removed profile page', () => {
        expect(allRoutes.find(route => route.path === '/profile')).toBeUndefined()
        expect(allRoutes.find(route => route.name === 'profile')).toBeUndefined()
    })

    it('keeps account deep links without exposing account as primary navigation', () => {
        const accountRoute = allRoutes.find(route => route.name === 'account')

        expect(accountRoute?.meta?.show).toBe(false)
        expect(allRoutes.find(route => route.path === '/account/user')?.component).toBeTypeOf('function')
        expect(allRoutes.find(route => route.path === '/account/role')?.component).toBeTypeOf('function')
        expect(allRoutes.find(route => route.path === '/account/api')?.component).toBeTypeOf('function')
    })
})
