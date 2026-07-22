import type { AgentInfo } from '@/interface/agent'
import { queryAgentList } from '@/api/agent'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import useStore from '@/store'
import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/api/agent', () => ({
    queryAgentList: vi.fn(),
}))

const queryAgentListMock = vi.mocked(queryAgentList)

function agent(agent_id: string, hostname = agent_id): AgentInfo {
    return {
        agent_id,
        hostname,
        version: 'v1.0.0',
        os: 'linux',
        arch: 'amd64',
        status: 'online',
    }
}

describe('useAgentSelection', () => {
    beforeEach(() => {
        queryAgentListMock.mockReset()
    })

    it('exposes empty derived state before agents are loaded', () => {
        const selection = useAgentSelection({ immediate: false })

        expect(selection.agentList.value).toEqual([])
        expect(selection.selectedAgentID.value).toBe('')
        expect(selection.loading.value).toBe(false)
        expect(selection.hasAgents.value).toBe(false)
        expect(selection.hasSelectedAgent.value).toBe(false)
        expect(selection.isAgentEmpty.value).toBe(false)
        expect(selection.agentParams.value).toEqual({})
    })

    it('loads agents and derives request params from the selected agent', async () => {
        queryAgentListMock.mockResolvedValueOnce({ data: [agent('agent-a'), agent('agent-b')] })
        const selection = useAgentSelection({ immediate: false })

        await expect(selection.loadAgents()).resolves.toBe('agent-a')

        expect(selection.loading.value).toBe(false)
        expect(selection.hasAgents.value).toBe(true)
        expect(selection.isAgentEmpty.value).toBe(false)
        expect(selection.selectedAgentID.value).toBe('agent-a')
        expect(selection.agentParams.value).toEqual({ agent_id: 'agent-a' })

        selection.selectedAgentID.value = 'agent-b'
        expect(selection.agentParams.value).toEqual({ agent_id: 'agent-b' })
    })

    it('marks the selection as empty after loading an empty agent list', async () => {
        queryAgentListMock.mockResolvedValueOnce({ data: [] })
        const selection = useAgentSelection({ immediate: false })

        await expect(selection.loadAgents()).resolves.toBe('')

        expect(selection.loading.value).toBe(false)
        expect(selection.hasAgents.value).toBe(false)
        expect(selection.isAgentEmpty.value).toBe(true)
        expect(selection.agentParams.value).toEqual({})
    })

    it('reuses an existing loaded selection without querying again', async () => {
        queryAgentListMock.mockResolvedValueOnce({ data: [agent('agent-a')] })
        const selection = useAgentSelection({ immediate: false })

        await selection.loadAgents()
        queryAgentListMock.mockClear()

        await expect(selection.ensureSelectedAgent()).resolves.toBe('agent-a')
        expect(queryAgentListMock).not.toHaveBeenCalled()
    })

    it('loads agents when no selected agent exists', async () => {
        queryAgentListMock.mockResolvedValueOnce({ data: [agent('agent-a')] })
        const selection = useAgentSelection({ immediate: false })

        await expect(selection.ensureSelectedAgent()).resolves.toBe('agent-a')

        expect(queryAgentListMock).toHaveBeenCalledTimes(1)
        expect(selection.selectedAgentID.value).toBe('agent-a')
    })

    it('deduplicates concurrent agent list requests across workspace sections', async () => {
        let resolveQuery!: (value: { data: AgentInfo[] }) => void
        queryAgentListMock.mockReturnValueOnce(new Promise(resolve => resolveQuery = resolve))
        const first = useAgentSelection({ immediate: false })
        const second = useAgentSelection({ immediate: false })

        const firstLoad = first.loadAgents()
        const secondLoad = second.loadAgents()

        expect(queryAgentListMock).toHaveBeenCalledTimes(1)
        resolveQuery({ data: [agent('agent-a')] })
        await expect(Promise.all([firstLoad, secondLoad])).resolves.toEqual(['agent-a', 'agent-a'])
    })

    it('clears the shared load promise after a rejected request', async () => {
        const requestError = new Error('session expired')
        queryAgentListMock
            .mockRejectedValueOnce(requestError)
            .mockResolvedValueOnce({ data: [agent('agent-after-login')] })
        const selection = useAgentSelection({ immediate: false })

        await expect(selection.loadAgents()).rejects.toBe(requestError)
        expect(selection.loading.value).toBe(false)

        await expect(selection.loadAgents()).resolves.toBe('agent-after-login')
        expect(queryAgentListMock).toHaveBeenCalledTimes(2)
        expect(selection.loading.value).toBe(false)
    })

    it('starts a fresh load when the Agent store is cleared during an old request', async () => {
        let rejectStale!: (reason: unknown) => void
        queryAgentListMock
            .mockReturnValueOnce(new Promise((_resolve, reject) => rejectStale = reject))
            .mockResolvedValueOnce({ data: [agent('agent-after-login')] })
        const selection = useAgentSelection({ immediate: false })
        const staleRequest = selection.loadAgents()
        expect(selection.loading.value).toBe(true)

        useStore().agent.clear()
        const freshRequest = selection.loadAgents()

        await expect(freshRequest).resolves.toBe('agent-after-login')
        expect(queryAgentListMock).toHaveBeenCalledTimes(2)
        rejectStale(new Error('old session expired'))
        await expect(staleRequest).rejects.toThrow('old session expired')
        expect(selection.selectedAgentID.value).toBe('agent-after-login')
        expect(selection.loading.value).toBe(false)
    })
})
