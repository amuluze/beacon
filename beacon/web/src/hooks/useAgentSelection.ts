import { queryAgentList } from '@/api/agent'
import useStore from '@/store'
import { computed, onMounted } from 'vue'

interface UseAgentSelectionOptions {
    immediate?: boolean
}

let agentLoadPromise: Promise<string> | null = null
let agentLoadGeneration = 0

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

    async function loadAgents(): Promise<string> {
        if (agentLoadPromise && store.agent.loading)
            return agentLoadPromise

        // clear() 会把 loading 复位；此时仍存在的 Promise 属于上一个认证会话，
        // 不能阻止新会话重新请求 Agent 列表。
        agentLoadPromise = null
        store.agent.setLoading(true)
        const loadGeneration = ++agentLoadGeneration
        const currentLoad = (async () => {
            try {
                const { data } = await queryAgentList()
                if (loadGeneration !== agentLoadGeneration || !store.agent.loading)
                    return store.agent.selectedAgentID
                store.agent.setAgents(data || [])
                return store.agent.selectedAgentID
            }
            finally {
                if (loadGeneration === agentLoadGeneration) {
                    store.agent.setLoading(false)
                    agentLoadPromise = null
                }
            }
        })()
        agentLoadPromise = currentLoad
        return currentLoad
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
