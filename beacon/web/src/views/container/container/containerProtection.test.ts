import { describe, expect, it } from 'vitest'

import { isProtectedContainerName } from './containerProtection'

describe('isProtectedContainerName', () => {
    it('protects the server container name', () => {
        expect(isProtectedContainerName('beacon')).toBe(true)
    })

    it('normalizes case and surrounding whitespace like the server', () => {
        expect(isProtectedContainerName('Beacon')).toBe(true)
        expect(isProtectedContainerName(' BEACON ')).toBe(true)
    })

    it('allows operations on unrelated containers', () => {
        expect(isProtectedContainerName('nginx')).toBe(false)
    })
})
