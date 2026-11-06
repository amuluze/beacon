import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import StatusBar from './index.vue'

const mocks = vi.hoisted(() => ({
    queryHostInfo: vi.fn(),
    queryDockerInfo: vi.fn(),
    querySystemVersion: vi.fn(),
    ensureSelectedAgent: vi.fn(),
}))

vi.mock('@/api/host', () => ({ queryHostInfo: mocks.queryHostInfo }))
vi.mock('@/api/container', () => ({ queryDockerInfo: mocks.queryDockerInfo }))
vi.mock('@/api/system', () => ({ querySystemVersion: mocks.querySystemVersion }))
vi.mock('@/hooks/useAgentSelection', () => ({
    useAgentSelection: () => ({
        agentList: ref([{ agent_id: 'agent-1', hostname: 'beacon-server', os: 'linux', arch: 'amd64' }]),
        selectedAgentID: ref('agent-1'),
        ensureSelectedAgent: mocks.ensureSelectedAgent,
    }),
}))

function mountStatusBar() {
    return mount(StatusBar)
}

function items(wrapper: ReturnType<typeof mountStatusBar>) {
    return wrapper.findAll('.am-statusbar__item').map(node => ({
        label: node.find('.am-statusbar__label').text(),
        value: node.find('.am-statusbar__value').text(),
    }))
}

describe('statusbar beacon server version item', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.ensureSelectedAgent.mockResolvedValue('agent-1')
        mocks.queryHostInfo.mockResolvedValue({
            data: { hostname: 'beacon-server', platform_version: 'Ubuntu 22.04', kernel_version: '5.15.0-91-generic', os: 'linux', kernel_arch: 'amd64', uptime: '12d 5h 32m' },
        })
        mocks.queryDockerInfo.mockResolvedValue({
            data: { docker_version: '24.0.7', api_version: '1.43', min_api_version: '1.42', os: 'linux', arch: 'amd64' },
        })
    })

    it('renders the build version as the last item of the right group', async () => {
        mocks.querySystemVersion.mockResolvedValue({ data: { version: 'v3.0.15', update_available: false } })

        const wrapper = mountStatusBar()
        await flushPromises()

        const rendered = items(wrapper)
        expect(rendered.at(-1)).toEqual({ label: 'Beacon', value: 'v3.0.15' })
        expect(rendered.map(item => item.label)).toEqual(['主机', '系统', '内核', 'OS', 'Docker', 'API', 'OS', 'Beacon'])
    })

    it('falls back to a placeholder when the version query fails', async () => {
        mocks.querySystemVersion.mockRejectedValue(new Error('offline'))

        const wrapper = mountStatusBar()
        await flushPromises()

        expect(items(wrapper).at(-1)).toEqual({ label: 'Beacon', value: '-' })
    })
})
