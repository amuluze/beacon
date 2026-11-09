import request from '@/api'
import type { AgentInfo, AgentUpdateSetting, AgentUpgradeResult } from '@/interface/agent.ts'

export async function queryAgentList() {
    return request.get<AgentInfo[]>('/api/v1/agent/list', {})
}

// 手动触发 Collia 远程升级；后端按 agent_id 定位节点，使用 Server 镜像内置版本升级。
export async function upgradeAgent(agentID: string) {
    return request.post<AgentUpgradeResult>('/api/v1/agent/upgrade', { agent_id: agentID })
}

// Agent 自动更新开关：记录缺失时后端按默认值（开启）返回。
export async function queryAgentUpdateSetting() {
    return request.get<AgentUpdateSetting>('/api/v1/agent/update_setting', {})
}

// 管理员开关仅影响注册后的自动调度；手动升级仍走 /api/v1/agent/upgrade。
export async function updateAgentUpdateSetting(autoUpdate: boolean) {
    return request.post<AgentUpdateSetting>('/api/v1/agent/update_setting', { auto_update: autoUpdate })
}
