import request from '@/api'
import type { AgentInfo } from '@/interface/store'

export async function queryAgentList() {
    return request.get<AgentInfo[]>('/api/v1/agent/list', {})
}
