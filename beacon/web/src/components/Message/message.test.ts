import { beforeEach, describe, expect, it, vi } from 'vitest'

import { error, info, success, warning } from './message'

const { ElMessage } = vi.hoisted(() => ({ ElMessage: vi.fn() }))

vi.mock('element-plus', () => ({ ElMessage }))

describe('message helpers', () => {
    beforeEach(() => {
        ElMessage.mockClear()
    })

    it.each([
        ['success', success],
        ['info', info],
        ['warning', warning],
        ['error', error],
    ])('%s renders the message as plain text', (_name, fn) => {
        fn('server <b>error</b>')
        expect(ElMessage).toHaveBeenCalledTimes(1)
        const options = ElMessage.mock.calls[0][0]
        expect(options.message).toBe('server <b>error</b>')
        // 服务端错误串可能包含用户可控内容，禁止以 HTML 注入 DOM。
        expect(options.dangerouslyUseHTMLString).toBeFalsy()
    })
})
