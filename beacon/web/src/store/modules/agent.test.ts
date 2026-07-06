import type { AgentInfo } from '@/interface/agent'
import { useAgentStore } from '@/store/modules/agent'
import { describe, expect, it } from 'vitest'

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

describe('useAgentStore', () => {
    it('starts with empty unloaded state', () => {
        const store = useAgentStore()

        expect(store.list).toEqual([])
        expect(store.selectedAgentID).toBe('')
        expect(store.loading).toBe(false)
        expect(store.loaded).toBe(false)
        expect(store.hasSelectedAgent).toBe(false)
    })

    it('selects the first agent when loading a non-empty list', () => {
        const store = useAgentStore()

        store.setAgents([agent('agent-a'), agent('agent-b')])

        expect(store.loaded).toBe(true)
        expect(store.list.map(item => item.agent_id)).toEqual(['agent-a', 'agent-b'])
        expect(store.selectedAgentID).toBe('agent-a')
        expect(store.hasSelectedAgent).toBe(true)
    })

    it('keeps the current selection when it still exists', () => {
        const store = useAgentStore()

        store.setAgents([agent('agent-a'), agent('agent-b')])
        store.setSelectedAgentID('agent-b')
        store.setAgents([agent('agent-b'), agent('agent-c')])

        expect(store.selectedAgentID).toBe('agent-b')
    })

    it('reselects the first agent when the selected agent disappears', () => {
        const store = useAgentStore()

        store.setAgents([agent('agent-a'), agent('agent-b')])
        store.setSelectedAgentID('agent-b')
        store.setAgents([agent('agent-c'), agent('agent-d')])

        expect(store.selectedAgentID).toBe('agent-c')
    })

    it('clears the selection when the loaded list is empty', () => {
        const store = useAgentStore()

        store.setAgents([agent('agent-a')])
        store.setAgents([])

        expect(store.loaded).toBe(true)
        expect(store.list).toEqual([])
        expect(store.selectedAgentID).toBe('')
        expect(store.hasSelectedAgent).toBe(false)
    })

    it('ignores unknown selections after agents are loaded', () => {
        const store = useAgentStore()

        store.setAgents([agent('agent-a')])
        store.setSelectedAgentID('missing-agent')

        expect(store.selectedAgentID).toBe('agent-a')
    })

    it('clears all state', () => {
        const store = useAgentStore()

        store.setLoading(true)
        store.setAgents([agent('agent-a')])
        store.clear()

        expect(store.list).toEqual([])
        expect(store.selectedAgentID).toBe('')
        expect(store.loading).toBe(false)
        expect(store.loaded).toBe(false)
    })
})
