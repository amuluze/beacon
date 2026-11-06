import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import type { VueWrapper } from '@vue/test-utils'

import LoginView from './index.vue'

const setUserInfo = vi.fn()
const setLanguage = vi.fn()
const clearAgent = vi.fn()

const stubStore = {
    user: { setUserInfo },
    app: { language: 'zh', setLanguage },
    agent: { clear: clearAgent },
}

vi.mock('@/store', () => ({
    default: () => stubStore,
}))

const loginMock = vi.fn()
const getUserInfoMock = vi.fn()

vi.mock('@/api/auth', () => ({
    login: (...args: unknown[]) => loginMock(...args),
    getUserInfo: () => getUserInfoMock(),
}))

const replace = vi.fn()

vi.mock('vue-router', () => ({
    useRouter: () => ({ replace }),
}))

const locale = vi.hoisted(() => ({ value: 'zh' }))

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key, locale }),
}))

let validateImpl: () => Promise<boolean> = async () => true

function mountLogin() {
    return mount(LoginView, {
        global: {
            mocks: { $t: (key: string) => key },
            stubs: {
                'el-form': {
                    name: 'ElForm',
                    methods: { validate: async () => validateImpl() },
                    template: '<form class="form-stub"><slot /></form>',
                },
                'el-form-item': { template: '<div class="form-item-stub"><slot /></div>' },
                'el-input': { template: '<input class="input-stub" />' },
                'el-button': {
                    template: '<button class="btn-stub" @click="$emit(\'click\', $event)"><slot /></button>',
                },
                'el-dropdown': {
                    name: 'ElDropdown',
                    template: '<div class="dropdown-stub"><slot /><slot name="dropdown" /></div>',
                },
                'el-dropdown-menu': { template: '<div><slot /></div>' },
                'el-dropdown-item': { template: '<div><slot /></div>' },
                'svg-icon': { template: '<i class="svg-stub" />' },
            },
        },
    })
}

async function submitForm(wrapper: VueWrapper) {
    await wrapper.find('.form-stub').trigger('submit')
    await Promise.resolve()
    await Promise.resolve()
}

describe('login view', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        validateImpl = async () => true
        locale.value = 'zh'
        loginMock.mockResolvedValue({ data: { must_change_password: false } })
        getUserInfoMock.mockResolvedValue({ data: { username: 'admin', status: 1, is_admin: true, must_change_password: false } })
    })

    it('renders the brand panel and the login card with two credential fields', () => {
        const wrapper = mountLogin()
        const brand = wrapper.find('.am-login__brand')

        expect(brand.exists()).toBe(true)
        expect(brand.find('.am-login__logo-mark').attributes('src')).toBe('/beacon.svg')
        expect(brand.find('.am-login__slogan').exists()).toBe(true)
        expect(brand.find('.am-login__subline').exists()).toBe(true)

        const card = wrapper.find('.am-login__card')
        expect(card.exists()).toBe(true)
        expect(card.find('.am-login__title').exists()).toBe(true)
        expect(card.find('.am-login__btn').exists()).toBe(true)
        expect(card.findAll('.form-item-stub')).toHaveLength(2)
    })

    it('keeps the language switcher inside the card', () => {
        const wrapper = mountLogin()

        expect(wrapper.find('.am-login__card .am-login__lang-chip').exists()).toBe(true)
    })

    it('switches the interface language without navigating away', async () => {
        const wrapper = mountLogin()

        wrapper.findComponent({ name: 'ElDropdown' }).vm.$emit('command', 'en')
        await Promise.resolve()

        expect(locale.value).toBe('en')
        expect(setLanguage).toHaveBeenCalledWith('en')
        expect(replace).not.toHaveBeenCalled()
    })

    it('does not authenticate when form validation fails', async () => {
        validateImpl = async () => {
            throw new Error('invalid')
        }
        const wrapper = mountLogin()
        await submitForm(wrapper)

        expect(loginMock).not.toHaveBeenCalled()
        expect(replace).not.toHaveBeenCalled()
    })

    it('submits credentials and redirects to monitor on success', async () => {
        const wrapper = mountLogin()
        await submitForm(wrapper)

        expect(loginMock).toHaveBeenCalledWith({ username: '', password: '' })
        expect(setUserInfo).toHaveBeenCalledWith('admin', 1, true, false)
        expect(clearAgent).toHaveBeenCalledOnce()
        expect(replace).toHaveBeenCalledWith('/monitor')
    })

    it('redirects bootstrap accounts to mandatory password change', async () => {
        getUserInfoMock.mockResolvedValueOnce({
            data: { username: 'admin', status: 1, is_admin: true, must_change_password: true },
        })
        const wrapper = mountLogin()
        await submitForm(wrapper)

        expect(setUserInfo).toHaveBeenCalledWith('admin', 1, true, true)
        expect(replace).toHaveBeenCalledWith('/change-password')
    })

    it('stays on the login page when credentials are rejected', async () => {
        loginMock.mockRejectedValueOnce(new Error('unauthorized'))
        const wrapper = mountLogin()
        await submitForm(wrapper)

        expect(replace).not.toHaveBeenCalled()
        expect(setUserInfo).not.toHaveBeenCalled()
    })
})
