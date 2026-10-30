import { beforeEach, describe, expect, it, vi } from 'vitest'
import { queryAgentList, upgradeAgent } from './index'

const mocks = vi.hoisted(() => ({
    get: vi.fn(),
    post: vi.fn(),
}))

vi.mock('@/api', () => ({
    default: {
        get: mocks.get,
        post: mocks.post,
    },
}))

describe('agent API', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.get.mockResolvedValue({ data: [] })
        mocks.post.mockResolvedValue({
            data: {
                agent_id: 'agent-1',
                current_version: 'v3.0.5',
                target_version: 'v3.0.6',
                arch: 'amd64',
                triggered: true,
                message: 'collia upgrade triggered',
            },
        })
    })

    it('lists the agents reachable through the reverse tunnel', async () => {
        await queryAgentList()

        expect(mocks.get).toHaveBeenCalledWith('/api/v1/agent/list', {})
    })

    it('triggers a manual Collia upgrade for the given agent', async () => {
        await upgradeAgent('agent-1')

        expect(mocks.post).toHaveBeenCalledWith('/api/v1/agent/upgrade', {
            agent_id: 'agent-1',
        })
    })

    it('returns the server side upgrade result', async () => {
        const { data } = await upgradeAgent('agent-1')

        expect(data?.current_version).toBe('v3.0.5')
        expect(data?.target_version).toBe('v3.0.6')
        expect(data?.triggered).toBe(true)
    })
})
