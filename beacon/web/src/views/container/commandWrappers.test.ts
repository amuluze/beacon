import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import en from '@/languages/modules/en'
import zh from '@/languages/modules/zh'

import DeleteContainer from './container/components/DeleteContainer.vue'
import RestartContainer from './container/components/RestartContainer.vue'
import StartContainer from './container/components/StartContainer.vue'
import StopContainer from './container/components/StopContainer.vue'
import DeleteImage from './image/components/DeleteImage.vue'
import PruneImage from './image/components/PruneImage.vue'
import DeleteNetwork from './network/components/DeleteNetwork.vue'

const { triggerConfirm, useConfirmCommand } = vi.hoisted(() => ({
    triggerConfirm: vi.fn(),
    useConfirmCommand: vi.fn(),
}))

vi.mock('@/api/container', () => ({
    pruneImages: vi.fn(),
    removeContainer: vi.fn(),
    removeImage: vi.fn(),
    removeNetwork: vi.fn(),
    restartContainer: vi.fn(),
    startContainer: vi.fn(),
    stopContainer: vi.fn(),
}))

vi.mock('@/hooks/useConfirmCommand', () => ({
    useConfirmCommand,
}))

const commandCases = [
    ['start container', StartContainer, { visible: true, id: 'container-1' }, ['container-1']],
    ['stop container', StopContainer, { visible: true, id: 'container-1' }, ['container-1']],
    ['restart container', RestartContainer, { visible: true, id: 'container-1' }, ['container-1']],
    ['delete container', DeleteContainer, { visible: true, id: 'container-1' }, ['container-1']],
    ['delete image', DeleteImage, { visible: true, id: 'image-1' }, ['image-1']],
    ['prune images', PruneImage, { visible: true }, []],
    ['delete network', DeleteNetwork, { visible: true, id: 'network-1' }, ['network-1']],
] as const

describe('container command wrappers', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        useConfirmCommand.mockReturnValue(triggerConfirm)
    })

    it.each(commandCases)('opens confirmation immediately for %s', (_name, component, props, expectedArgs) => {
        mount(component, { props })

        expect(triggerConfirm.mock.calls).toEqual([expectedArgs])
    })

    it('provides localized network deletion copy', () => {
        expect(zh.network.deleteNetwork).toBe('删除网络')
        expect(zh.network.confirmDelete).toBe('确认要删除网络吗？')
        expect(en.network.deleteNetwork).toBe('Delete Network')
        expect(en.network.confirmDelete).toBe('Confirm Delete?')
    })
})
