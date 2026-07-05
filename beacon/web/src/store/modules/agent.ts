import type { AgentInfo, AgentState } from '@/interface/store'
import { defineStore } from 'pinia'

export const useAgentStore = defineStore('agent', {
    state: (): AgentState => <AgentState>({
        currentAgentID: '',
        agents: [],
    }),
    actions: {
        setAgents(agents: AgentInfo[]) {
            this.agents = agents
            if (agents.length === 0) {
                this.currentAgentID = ''
                return
            }
            if (!this.currentAgentID || !agents.some(item => item.agent_id === this.currentAgentID)) {
                this.currentAgentID = agents[0].agent_id
            }
        },
        setCurrentAgent(agentID: string) {
            this.currentAgentID = agentID
        },
    },
    persist: true,
})
