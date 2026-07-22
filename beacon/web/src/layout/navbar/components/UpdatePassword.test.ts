import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import UpdatePassword from './UpdatePassword.vue'

const mocks = vi.hoisted(() => ({
    replace: vi.fn(),
    setToken: vi.fn(),
    success: vi.fn(),
    updatePassword: vi.fn(),
}))

vi.mock('@/api/auth', () => ({
    updatePassword: mocks.updatePassword,
}))

vi.mock('@/components/Message/message.ts', () => ({
    success: mocks.success,
}))

vi.mock('@/store', () => ({
    default: () => ({
        user: {
            setToken: mocks.setToken,
        },
    }),
}))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key }),
}))

vi.mock('vue-router', () => ({
    useRouter: () => ({ replace: mocks.replace }),
}))

const ElButton = defineComponent({
    inheritAttrs: false,
    emits: ['click'],
    setup(_, { emit, slots }) {
        return () => h('button', { onClick: () => emit('click') }, slots.default?.())
    },
})

const ElForm = defineComponent({
    setup(_, { expose, slots }) {
        expose({ validate: async () => true })
        return () => h('form', slots.default?.())
    },
})

const ElInput = defineComponent({
    props: {
        disabled: Boolean,
        modelValue: String,
        type: String,
    },
    emits: ['update:modelValue'],
    setup(props, { emit }) {
        return () => h('input', {
            disabled: props.disabled,
            type: props.type,
            value: props.modelValue,
            onInput: (event: Event) => emit('update:modelValue', (event.target as HTMLInputElement).value),
        })
    },
})

const global = {
    directives: {
        loading: {},
    },
    stubs: {
        'el-button': ElButton,
        'el-drawer': { template: '<section><slot /><slot name="footer" /></section>' },
        'el-form': ElForm,
        'el-form-item': { template: '<label><slot /></label>' },
        'el-input': ElInput,
    },
}

async function submitPasswordUpdate() {
    const wrapper = mount(UpdatePassword, {
        props: {
            username: 'admin',
            visible: true,
        },
        global,
    })
    const inputs = wrapper.findAll('input')
    await inputs[1].setValue('old-pass')
    await inputs[2].setValue('new-pass')
    await wrapper.findAll('button')[1].trigger('click')
    await flushPromises()
    return wrapper
}

describe('update password drawer', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mocks.replace.mockResolvedValue(undefined)
        mocks.updatePassword.mockResolvedValue({})
    })

    it('submits the current username and requires a new login after success', async () => {
        const wrapper = await submitPasswordUpdate()

        expect(wrapper.findAll('input')[0].attributes('disabled')).toBeDefined()
        expect(mocks.updatePassword).toHaveBeenCalledWith({
            username: 'admin',
            old_password: 'old-pass',
            new_password: 'new-pass',
        })
        expect(mocks.success).toHaveBeenCalledWith('更新成功')
        expect(mocks.setToken).toHaveBeenCalledWith('', '')
        expect(mocks.replace).toHaveBeenCalledWith('/login')
        expect(wrapper.emitted('update:visible')).toEqual([[false]])
    })

    it('keeps the drawer and login state when the update request fails', async () => {
        mocks.updatePassword.mockRejectedValueOnce(new Error('invalid password'))

        const wrapper = await submitPasswordUpdate()

        expect(mocks.success).not.toHaveBeenCalled()
        expect(mocks.setToken).not.toHaveBeenCalled()
        expect(mocks.replace).not.toHaveBeenCalled()
        expect(wrapper.emitted('update:visible')).toBeUndefined()
    })
})
