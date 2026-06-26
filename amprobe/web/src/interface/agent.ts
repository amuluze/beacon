// Agent 列表项契约，与后端 /api/v1/agent/list 返回结构对齐。
export interface AgentInfo {
  agent_id: string
  hostname: string
  version: string
  os: string
  arch: string
  status: string
}

// 构造一个填充默认值的 Agent，用于前端兜底场景（如 agent 列表查询失败）。
export function createDefaultAgent(agent_id = 'default', hostname = 'default'): AgentInfo {
  return {
    agent_id,
    hostname,
    version: '',
    os: '',
    arch: '',
    status: 'offline',
  }
}
