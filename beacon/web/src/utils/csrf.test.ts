import { beforeEach, describe, expect, it } from 'vitest'
import { getCSRFHeaders, getCSRFToken } from './csrf'

describe('csrf cookie helper', () => {
    beforeEach(() => {
        document.cookie = 'beacon_csrf_token=; Max-Age=0; path=/'
    })

    it('reads and decodes the CSRF cookie', () => {
        document.cookie = 'beacon_csrf_token=csrf%2Dvalue; path=/'
        expect(getCSRFToken()).toBe('csrf-value')
        expect(getCSRFHeaders()).toEqual({ 'X-CSRF-Token': 'csrf-value' })
    })

    it('returns no header when the CSRF cookie is missing', () => {
        expect(getCSRFToken()).toBe('')
        expect(getCSRFHeaders()).toEqual({})
    })
})
