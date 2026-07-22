/**
 * Smoke tests for useConfirmCommand — verifies the hook returns a callable
 * function that consumers can use to drive the ConfirmDialog imperatively.
 *
 * We deliberately do NOT call open() in jsdom: the imperative render() path
 * depends on a full Vue app context (vue-i18n provide chain). End-to-end
 * behaviour is exercised via the platform; the component-level behaviour
 * lives in ConfirmDialog's own unit tests.
 */
import { flushPromises, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { defineComponent, h, onMounted } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

const { confirmDialogTranslation } = vi.hoisted(() => ({
    confirmDialogTranslation: vi.fn(),
}))

vi.mock('@/components/ConfirmDialog/index.vue', async () => {
    const { defineComponent, h } = await import('vue')
    const { useI18n } = await import('vue-i18n')
    return {
        default: defineComponent({
            emits: ['update:visible'],
            setup(_, { emit }) {
                const { t } = useI18n()
                confirmDialogTranslation(t('common.confirm'))
                queueMicrotask(() => emit('update:visible', false))
                return () => h('div')
            },
        }),
    }
})

vi.mock('@/components/Message/message.ts', () => ({
    warning: vi.fn(),
}))

describe('useConfirmCommand', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('returns a callable function', () => {
        const open = useConfirmCommand({
            title: 'x.delete',
            message: 'x.confirmDelete',
            i18nPrefix: 'x',
            action: async () => {},
        })
        expect(typeof open).toBe('function')
    })

    it('passes the calling component app context to the imperative dialog', async () => {
        const Host = defineComponent({
            setup() {
                const open = useConfirmCommand({
                    action: async () => {},
                })
                onMounted(() => open('container-1'))
                return () => h('div')
            },
        })
        const i18n = createI18n({
            legacy: false,
            locale: 'zh',
            messages: { zh: { common: { confirm: '确认' } } },
        })

        mount(Host, { global: { plugins: [i18n] } })
        await flushPromises()

        expect(confirmDialogTranslation).toHaveBeenCalledWith('确认')
    })
})
