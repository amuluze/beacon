import { describe, expect, it } from 'vitest'
import { convertBytesToReadable } from '@/utils/convert'

describe('convertBytesToReadable', () => {
    it('returns bytes for small values', () => {
        expect(convertBytesToReadable(0)).toBe('0.00 B')
        expect(convertBytesToReadable(512)).toBe('512.00 B')
    })

    it('converts to KB', () => {
        expect(convertBytesToReadable(1024)).toBe('1.00 KB')
        expect(convertBytesToReadable(2048)).toBe('2.00 KB')
    })

    it('converts to MB', () => {
        expect(convertBytesToReadable(1048576)).toBe('1.00 MB')
    })

    it('converts to GB', () => {
        expect(convertBytesToReadable(1073741824)).toBe('1.00 GB')
    })

    it('converts to TB', () => {
        expect(convertBytesToReadable(1099511627776)).toBe('1.00 TB')
    })

    it('stops at TB for very large values', () => {
        const result = convertBytesToReadable(1099511627776 * 1024)
        expect(result).toContain('TB')
    })
})
