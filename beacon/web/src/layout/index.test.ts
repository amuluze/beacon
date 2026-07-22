import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import AppLayout from './index.vue'

const ensureSelectedAgent = vi.fn()

vi.mock('@/hooks/useAgentSelection', () => ({
    useAgentSelection: () => ({ ensureSelectedAgent }),
}))

function deferred<T>() {
    let resolve!: (value: T) => void
    const promise = new Promise<T>((promiseResolve) => {
        resolve = promiseResolve
    })
    return { promise, resolve }
}

function mountLayout() {
    return mount(AppLayout, {
        global: {
            mocks: { $t: (key: string) => key },
            stubs: {
                Content: { template: '<main data-testid="content" />' },
                Navbar: { template: '<nav data-testid="navbar" />' },
                StatusBar: { template: '<footer data-testid="status-bar" />' },
            },
        },
    })
}

describe('application layout Agent readiness gate', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('mounts workspace content only after Agent initialization settles', async () => {
        const agentReady = deferred<string>()
        ensureSelectedAgent.mockReturnValue(agentReady.promise)

        const wrapper = mountLayout()

        expect(wrapper.find('[data-testid="navbar"]').exists()).toBe(true)
        expect(wrapper.find('[role="status"]').text()).toContain('agent.loadingWorkspace')
        expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="status-bar"]').exists()).toBe(false)

        agentReady.resolve('agent-1')
        await flushPromises()

        expect(wrapper.find('[role="status"]').exists()).toBe(false)
        expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
        expect(wrapper.find('[data-testid="status-bar"]').exists()).toBe(true)
    })
})
