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
        ['/setting/alarm'],
        ['/setting/container'],
    ])('redirects the legacy deep link %s into the settings workspace', (path) => {
        const route = allRoutes.find(item => item.path === path)

        expect(route?.component).toBeUndefined()
        expect(route?.redirect).toEqual({ path: '/setting', query: { tab: 'config' } })
    })

    it('does not expose the removed system settings page', () => {
        expect(allRoutes.find(route => route.path === '/setting/host')).toBeUndefined()
        expect(allRoutes.find(route => route.name === 'host')).toBeUndefined()
    })

    it('does not expose the removed profile page', () => {
        expect(allRoutes.find(route => route.path === '/profile')).toBeUndefined()
        expect(allRoutes.find(route => route.name === 'profile')).toBeUndefined()
    })

    it('does not expose the removed terminal page', () => {
        expect(allRoutes.find(route => route.path === '/terminal')).toBeUndefined()
        expect(allRoutes.find(route => route.name === 'terminal')).toBeUndefined()
    })

    it('redirects account deep links into the users tab without exposing account as primary navigation', () => {
        const accountRoute = allRoutes.find(route => route.name === 'account')

        expect(accountRoute?.meta?.show).toBe(false)
        for (const path of ['/account/user', '/account/role', '/account/api']) {
            expect(allRoutes.find(route => route.path === path)?.redirect).toEqual({
                path: '/setting',
                query: { tab: 'users' },
            })
        }
    })
})
