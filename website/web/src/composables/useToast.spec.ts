import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { dismissToast, showToast, useToast } from './useToast'

describe('useToast', () => {
    beforeEach(() => {
        vi.useFakeTimers()
        dismissToast()
    })

    afterEach(() => vi.useRealTimers())

    it('展示消息并在超时后自动关闭', () => {
        const { toast } = useToast()

        showToast('安装命令已复制', 'success', 2000)
        expect(toast.value).toMatchObject({ message: '安装命令已复制', type: 'success' })

        vi.advanceTimersByTime(2000)
        expect(toast.value).toBeNull()
    })

    it('后续消息替换前一条消息', () => {
        const { toast } = useToast()

        showToast('第一条', 'info')
        showToast('第二条', 'error')

        expect(toast.value).toMatchObject({ message: '第二条', type: 'error' })
    })
})
