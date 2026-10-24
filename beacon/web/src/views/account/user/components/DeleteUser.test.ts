import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import DeleteUser from './DeleteUser.vue'

const mocks = vi.hoisted(() => ({
    deleteUser: vi.fn(),
}))

vi.mock('@/api/account', () => ({
    deleteUser: mocks.deleteUser,
}))
vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

const stubs = {
    'el-dialog': { template: '<div><slot /><slot name="footer" /></div>', props: ['modelValue', 'title', 'width'] },
    'el-button': { template: '<button type="button" @click="$emit(\'click\')"><slot /></button>', emits: ['click'] },
}

describe('delete user dialog', () => {
    beforeEach(() => {
        mocks.deleteUser.mockReset()
    })

    it('submits the selected ids and refreshes only after deletion succeeds', async () => {
        const update = vi.fn()
        mocks.deleteUser.mockResolvedValue(undefined)
        const wrapper = mount(DeleteUser, {
            props: {
                visible: true,
                ids: ['019f8dbc-8781-7af3-bf49-e7fafa912e87'],
                update,
            },
            global: { stubs },
        })

        await wrapper.findAll('button')[1].trigger('click')
        await flushPromises()

        expect(mocks.deleteUser).toHaveBeenCalledWith({
            ids: ['019f8dbc-8781-7af3-bf49-e7fafa912e87'],
        })
        expect(wrapper.emitted('update:visible')).toEqual([[false]])
        expect(update).toHaveBeenCalledOnce()
    })

    it('keeps the dialog and list state when deletion fails', async () => {
        const update = vi.fn()
        mocks.deleteUser.mockRejectedValue(new Error('delete failed'))
        const wrapper = mount(DeleteUser, {
            props: {
                visible: true,
                ids: ['019f8dbc-8781-7af3-bf49-e7fafa912e87'],
                update,
            },
            global: { stubs },
        })

        await wrapper.findAll('button')[1].trigger('click')
        await flushPromises()

        expect(wrapper.emitted('update:visible')).toBeUndefined()
        expect(update).not.toHaveBeenCalled()
    })
})
