import request from '@/api'
import type { AgentInfo, AgentUpgradeResult } from '@/interface/agent.ts'

export async function queryAgentList() {
    return request.get<AgentInfo[]>('/api/v1/agent/list', {})
}

// 手动触发 Collia 远程升级；后端按 agent_id 定位节点，使用 Server 镜像内置版本升级。
export async function upgradeAgent(agentID: string) {
    return request.post<AgentUpgradeResult>('/api/v1/agent/upgrade', { agent_id: agentID })
}
