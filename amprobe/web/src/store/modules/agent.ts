import type { AgentInfo } from '@/interface/agent.ts'
import { defineStore } from 'pinia'

interface AgentState {
    list: AgentInfo[]
    selectedAgentID: string
    loading: boolean
    loaded: boolean
}

export const useAgentStore = defineStore('agent', {
    state: (): AgentState => ({
        list: [],
        selectedAgentID: '',
        loading: false,
        loaded: false,
    }),
    getters: {
        hasSelectedAgent: state => state.selectedAgentID !== '',
    },
    actions: {
        setLoading(loading: boolean) {
            this.loading = loading
        },
        setAgents(agents: AgentInfo[]) {
            this.list = agents
            this.loaded = true

            if (this.list.length === 0) {
                this.selectedAgentID = ''
                return
            }

            const selectedExists = this.list.some(agent => agent.agent_id === this.selectedAgentID)
            if (!selectedExists) {
                this.selectedAgentID = this.list[0].agent_id
            }
        },
        setSelectedAgentID(agentID: string) {
            if (agentID === '') {
                this.selectedAgentID = ''
                return
            }
            if (this.list.length === 0 || this.list.some(agent => agent.agent_id === agentID)) {
                this.selectedAgentID = agentID
            }
        },
        clear() {
            this.list = []
            this.selectedAgentID = ''
            this.loading = false
            this.loaded = false
        },
    },
})
