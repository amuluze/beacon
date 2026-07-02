// Agent 列表项契约，与后端 /api/v1/agent/list 返回结构对齐。
export interface AgentInfo {
    agent_id: string
    hostname: string
    version: string
    os: string
    arch: string
    status: string
}
