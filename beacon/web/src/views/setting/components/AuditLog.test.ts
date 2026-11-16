import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import AuditLog from './AuditLog.vue'

const mocks = vi.hoisted(() => ({
    queryAudit: vi.fn(),
}))

vi.mock('@/api/audit', () => ({
    queryAudit: mocks.queryAudit,
}))

function mountAuditLog() {
    return mount(AuditLog, {
        global: {
            stubs: {
                'el-table': true,
                'el-table-column': true,
                'el-pagination': true,
            },
        },
    })
}

describe('audit log merged timeline', () => {
    beforeEach(() => {
        mocks.queryAudit.mockReset()
    })

    // 审计页必须同时展示系统与操作员记录：查询不得再固定 type=system 或按 Agent 过滤。
    it('queries both system and operator rows without source filters', async () => {
        mocks.queryAudit.mockResolvedValue({
            data: {
                data: [
                    { id: 1, username: 'system', agent_id: 'node-01', operate: 'Collia upgrade succeeded', created: '' },
                    { id: 2, username: 'admin', agent_id: '', operate: '创建网络', created: '' },
                ],
                total: 2,
            },
        })

        mountAuditLog()
        await flushPromises()

        expect(mocks.queryAudit).toHaveBeenCalledTimes(1)
        expect(mocks.queryAudit).toHaveBeenCalledWith({ page: 1, size: 10 })
    })
})
