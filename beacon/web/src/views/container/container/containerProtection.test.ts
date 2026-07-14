import { describe, expect, it } from 'vitest'

import { isProtectedContainerName } from './containerProtection'

describe('isProtectedContainerName', () => {
    it.each(['beacon', 'amprobe'])('protects the server container name %s', (containerName) => {
        expect(isProtectedContainerName(containerName)).toBe(true)
    })

    it('allows operations on unrelated containers', () => {
        expect(isProtectedContainerName('nginx')).toBe(false)
    })
})
