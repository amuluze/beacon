import { afterEach, describe, expect, it, vi } from 'vitest'

import { useHttp } from '~/composables/useHttp'
import { statisticQuery, statisticUpdate } from './statistics'

// 用稳定的 mock 替换 useHttp，专注验证 api 层契约与错误转换，
// 不依赖 ofetch 网络层与 Nuxt 运行时上下文。
vi.mock('~/composables/useHttp', () => ({
    useHttp: {
        get: vi.fn(),
        post: vi.fn(),
    },
}))

afterEach(() => vi.clearAllMocks())

describe('api/statistics', () => {
    it('statisticQuery 解包后端 data 信封并静默返回业务体', async () => {
        (useHttp.get as any).mockResolvedValue({ data: { id: 9, times: 3 } })

        const r = await statisticQuery()

        expect(useHttp.get).toHaveBeenCalledWith('/api/v1/statistics/query', {}, { silent: true })
        expect(r).toEqual({ id: 9, times: 3 })
    })

    it('statisticQuery 拒绝不符合后端契约的响应', async () => {
        (useHttp.get as any).mockResolvedValue({ id: 9, times: 3 })

        await expect(statisticQuery()).rejects.toThrow('统计响应格式错误')
    })

    it('statisticUpdate 成功时透传参数', async () => {
        (useHttp.post as any).mockResolvedValue({})

        await expect(statisticUpdate({ id: 1 })).resolves.toEqual({})
        expect(useHttp.post).toHaveBeenCalledWith('/api/v1/statistics/update', { id: 1 })
    })

    it('statisticUpdate 失败时抛出友好错误并吞掉原始异常', async () => {
        (useHttp.post as any).mockRejectedValue(new Error('network down'))

        await expect(statisticUpdate({ id: 1 })).rejects.toThrow('数据更新失败')
    })
})
