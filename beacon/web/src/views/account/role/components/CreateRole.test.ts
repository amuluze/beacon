import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import CreateRole from './CreateRole.vue'

const mocks = vi.hoisted(() => ({
    createRole: vi.fn(),
    queryResource: vi.fn(),
    success: vi.fn(),
}))

vi.mock('@/api/account', () => ({
    createRole: mocks.createRole,
    queryResource: mocks.queryResource,
}))

vi.mock('@/components/Message/message', () => ({
    success: mocks.success,
}))

// 弹层内部用 useI18n 读取文案，测试里替换成与 $t 一致的格式。
vi.mock('vue-i18n', () => ({
    useI18n: () => ({
        t: (key: string, params?: Record<string, unknown>) => (params ? `${key}:${JSON.stringify(params)}` : key),
    }),
}))

const resources = [
    { id: 'r1', name: '主机查询', path: '/api/v1/host/query', method: 'GET', status: 1 },
    { id: 'r2', name: '主机上报', path: '/api/v1/host/report', method: 'POST', status: 1 },
    { id: 'r3', name: '用户查询', path: '/api/v1/user/user_query', method: 'GET', status: 1 },
]

function mountDialog() {
    return mount(CreateRole, {
        props: { visible: true },
        global: {
            mocks: {
                $t: (key: string, params?: Record<string, unknown>) => (params ? `${key}:${JSON.stringify(params)}` : key),
            },
            stubs: {
                'svg-icon': true,
            },
        },
    })
}

function actionButton(wrapper: ReturnType<typeof mountDialog>, label: string) {
    return wrapper.findAll('button').find(button => button.text().includes(label))
}

describe('create role dialog', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.queryResource.mockResolvedValue({ data: { data: resources } })
        mocks.createRole.mockResolvedValue({})
    })

    it('groups the callable apis by their path segment', async () => {
        const wrapper = mountDialog()
        await vi.waitFor(() => expect(wrapper.text()).toContain('/api/v1/host/query'))

        // 分组名取自 /api/v1/<group>/<action> 的首段，接口名与路径一并展示。
        expect(wrapper.text()).toContain('host')
        expect(wrapper.text()).toContain('user')
        expect(wrapper.text()).toContain('主机上报')
        expect(wrapper.text()).toContain('POST')
    })

    it('selects every api through the select-all action', async () => {
        const wrapper = mountDialog()
        await vi.waitFor(() => expect(wrapper.text()).toContain('/api/v1/host/query'))

        await actionButton(wrapper, 'account.rolePermissionSelectAll')!.trigger('click')

        expect(wrapper.text()).toContain('account.rolePermissionSelected:{"selected":3,"total":3}')
    })

    it('clears the selection through the clear action', async () => {
        const wrapper = mountDialog()
        await vi.waitFor(() => expect(wrapper.text()).toContain('/api/v1/host/query'))

        await actionButton(wrapper, 'account.rolePermissionSelectAll')!.trigger('click')
        await actionButton(wrapper, 'account.rolePermissionClear')!.trigger('click')

        expect(wrapper.text()).toContain('account.rolePermissionSelected:{"selected":0,"total":3}')
    })
})
