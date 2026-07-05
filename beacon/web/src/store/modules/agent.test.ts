import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAgentStore } from '@/store/modules/agent'

describe('useAgentStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('sets current agent from list', () => {
        const store = useAgentStore()
        store.setAgents([
            { agent_id: 'agent-a', hostname: 'host-a', status: 'online', last_seen: String(Date.now()) },
            { agent_id: 'agent-b', hostname: 'host-b', status: 'online', last_seen: String(Date.now()) },
        ])
        expect(store.currentAgentID).toBe('agent-a')
    })

    it('clears current agent when list is empty', () => {
        const store = useAgentStore()
        store.setCurrentAgent('agent-a')
        store.setAgents([])
        expect(store.currentAgentID).toBe('')
    })

    it('switches current agent explicitly', () => {
        const store = useAgentStore()
        store.setAgents([{ agent_id: 'agent-a', hostname: 'host-a', status: 'online', last_seen: String(Date.now()) }])
        store.setCurrentAgent('agent-b')
        expect(store.currentAgentID).toBe('agent-b')
    })
})
