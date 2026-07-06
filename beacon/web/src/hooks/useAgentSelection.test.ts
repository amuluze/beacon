import type { AgentInfo } from '@/interface/agent'
import { queryAgentList } from '@/api/agent'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { describe, expect, it, vi } from 'vitest'

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
})
