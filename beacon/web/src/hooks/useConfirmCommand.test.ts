/**
 * Smoke tests for useConfirmCommand — verifies the hook returns a callable
 * function that consumers can use to drive the ConfirmDialog imperatively.
 *
 * We deliberately do NOT call open() in jsdom: the imperative render() path
 * depends on a full Vue app context (vue-i18n provide chain). End-to-end
 * behaviour is exercised via the platform; the component-level behaviour
 * lives in ConfirmDialog's own unit tests.
 */
import { describe, expect, it, vi } from 'vitest'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

vi.mock('@/components/Message/message.ts', () => ({
    warning: vi.fn(),
}))

describe('useConfirmCommand', () => {
    it('returns a callable function', () => {
        const open = useConfirmCommand({
            title: 'x.delete',
            message: 'x.confirmDelete',
            i18nPrefix: 'x',
            action: async () => {},
        })
        expect(typeof open).toBe('function')
    })
})
