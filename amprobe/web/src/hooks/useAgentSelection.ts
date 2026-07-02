import { queryAgentList } from '@/api/agent'
import useStore from '@/store'
import { computed, onMounted } from 'vue'

interface UseAgentSelectionOptions {
    immediate?: boolean
}

export function useAgentSelection(options: UseAgentSelectionOptions = {}) {
    const store = useStore()

    const agentList = computed(() => store.agent.list)
    const selectedAgentID = computed({
        get: () => store.agent.selectedAgentID,
        set: (agentID: string) => store.agent.setSelectedAgentID(agentID),
    })
    const loading = computed(() => store.agent.loading)
    const hasAgents = computed(() => store.agent.list.length > 0)
    const hasSelectedAgent = computed(() => store.agent.hasSelectedAgent)
    const isAgentEmpty = computed(() => store.agent.loaded && !store.agent.loading && !hasAgents.value)
    const agentParams = computed<Record<string, string>>((): Record<string, string> => {
        if (!selectedAgentID.value)
            return {}
        return { agent_id: selectedAgentID.value }
    })

    async function loadAgents() {
        store.agent.setLoading(true)
        try {
            const { data } = await queryAgentList()
            store.agent.setAgents(data || [])
        }
        finally {
            store.agent.setLoading(false)
        }
        return store.agent.selectedAgentID
    }

    async function ensureSelectedAgent() {
        if (!store.agent.loaded || !store.agent.selectedAgentID) {
            return loadAgents()
        }
        return store.agent.selectedAgentID
    }

    if (options.immediate !== false) {
        onMounted(() => {
            void loadAgents()
        })
    }

    return {
        agentList,
        selectedAgentID,
        loading,
        hasAgents,
        hasSelectedAgent,
        isAgentEmpty,
        agentParams,
        loadAgents,
        ensureSelectedAgent,
    }
}
