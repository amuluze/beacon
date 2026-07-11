import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const setToken = vi.fn()
const setUserInfo = vi.fn()
const setLanguage = vi.fn()

const stubStore = {
  user: { setToken, setUserInfo },
  app: { language: 'zh', setLanguage },
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

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
    locale: { value: 'zh' },
  }),
}))

import LoginView from './index.vue'

function mountLogin() {
  return mount(LoginView, {
    global: {
      mocks: { $t: (key: string) => key },
      stubs: {
        'el-form': { template: '<form class="form-stub"><slot /></form>' },
        'el-form-item': { template: '<div class="form-item-stub"><slot /></div>' },
        'el-input': { template: '<input class="input-stub" />' },
        'el-button': {
          template: '<button class="btn-stub" @click="$emit(\'click\', $event)"><slot /></button>',
        },
        'el-dropdown': { template: '<div class="dropdown-stub"><slot /><slot name="dropdown" /></div>' },
        'el-dropdown-menu': { template: '<div><slot /></div>' },
        'el-dropdown-item': { template: '<div><slot /></div>' },
        'svg-icon': { template: '<i class="svg-stub" />' },
      },
    },
  })
}

describe('login view', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    loginMock.mockResolvedValue({ data: { access_token: 'a', refresh_token: 'r' } })
    getUserInfoMock.mockResolvedValue({ data: { username: 'admin', status: 1, is_admin: true } })
  })

  it('renders the branded split layout with slogan and subline', () => {
    const wrapper = mountLogin()

    expect(wrapper.find('.am-login__brand').exists()).toBe(true)
    expect(wrapper.find('.am-login__logo-text').text()).toBe('Beacon')
    expect(wrapper.find('.am-login__slogan').text()).toBe('login.slogan')
    expect(wrapper.find('.am-login__subline').text()).toBe('login.subline')
  })

  it('renders the login card with title and submit button on the right pane', () => {
    const wrapper = mountLogin()

    const card = wrapper.find('.am-login__card')
    expect(card.exists()).toBe(true)
    expect(card.find('.am-login__title').text()).toBe('login.login')
    expect(card.find('.am-login__btn').text()).toBe('login.login')
    expect(wrapper.findAll('.form-item-stub')).toHaveLength(2)
  })

  it('submits credentials and redirects to monitor on success', async () => {
    const wrapper = mountLogin()
    await wrapper.find('.am-login__btn').trigger('click')
    await Promise.resolve()
    await Promise.resolve()

    expect(loginMock).toHaveBeenCalled()
    expect(setToken).toHaveBeenCalledWith('a', 'r')
    expect(setUserInfo).toHaveBeenCalledWith('admin', 1, true)
    expect(replace).toHaveBeenCalledWith('/monitor/host')
  })
})
