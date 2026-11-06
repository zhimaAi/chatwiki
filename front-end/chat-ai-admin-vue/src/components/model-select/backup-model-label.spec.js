import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import BackupModelLabel from './backup-model-label.vue'
import zh from '@/locales/lang/zh-CN/components/model-select/model-select.json'

const mocks = vi.hoisted(() => ({
  getStatus: vi.fn(),
  permissions: { role_permission: ['ModelManage'], role_type: '2' },
  company: { admin_domain: '', agent_domain: '' },
  user: { getToken: '', token: {}, userInfo: {} }
}))
vi.mock('@/api/model/index', () => ({ getBackupModelStatus: mocks.getStatus }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => zh[key] }) }))
vi.mock('@/stores/modules/permission', () => ({ usePermissionStore: () => mocks.permissions }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => mocks.company }))
vi.mock('@/stores/modules/user', () => ({ useUserStore: () => mocks.user }))

const wrappers = []
const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/', component: { template: '<div />' } },
    { path: '/clawbot/settings', component: { template: '<div />' } },
    { path: '/clawbot/chat', component: { template: '<div />' } },
    { path: '/user/model', component: { template: '<div />' } }
  ]
})
function render(props = {}) {
  const wrapper = mount(BackupModelLabel, {
    props: { label: 'LLM模型', ...props },
    global: { plugins: [router] }
  })
  wrappers.push(wrapper)
  return wrapper
}
function deferred() {
  let resolve
  const promise = new Promise(done => { resolve = done })
  return { promise, resolve }
}
beforeEach(async () => {
  mocks.getStatus.mockReset().mockResolvedValue({ data: { enabled: false } })
  mocks.permissions.role_permission = ['ModelManage']
  mocks.permissions.role_type = '2'
  mocks.company.admin_domain = ''
  mocks.company.agent_domain = ''
  mocks.user.getToken = ''
  mocks.user.token = {}
  mocks.user.userInfo = {}
  await router.push('/')
  vi.spyOn(window, 'open').mockImplementation(() => null)
})
afterEach(() => {
  wrappers.splice(0).forEach(wrapper => wrapper.unmount())
  vi.restoreAllMocks()
})

describe('backup model notice', () => {
  it('shows a new-tab setup link only after the read-only status confirms no backup', async () => {
    const pending = deferred()
    mocks.getStatus.mockReturnValue(pending.promise)
    const wrapper = render()
    expect(wrapper.text()).toBe('LLM模型')
    expect(wrapper.find('a').exists()).toBe(false)
    pending.resolve({ data: { enabled: false } })
    await flushPromises()
    expect(wrapper.text()).toContain('备用模型暂未开启，')
    const link = wrapper.get('a')
    expect(link.text()).toBe('点此设置')
    expect(link.attributes()).toMatchObject({
      href: '/user/model?activeTab=3', target: '_blank', rel: 'noopener noreferrer'
    })
    expect(window.open).not.toHaveBeenCalled()
  })

  it('does not display the notice when a backup model is configured', async () => {
    mocks.getStatus.mockResolvedValue({ data: { enabled: true } })
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toBe('LLM模型')
    expect(wrapper.find('a').exists()).toBe(false)
  })

  it('hides stale status during refresh after returning from settings and reflects clearing the backup', async () => {
    const wrapper = render()
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(true)
    const pending = deferred()
    mocks.getStatus.mockReturnValueOnce(pending.promise)
    window.dispatchEvent(new Event('focus'))
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(false)
    pending.resolve({ data: { enabled: true } })
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(false)
    window.dispatchEvent(new Event('focus'))
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(true)
  })

  it('loads only while the robot model dialog is open and reloads on reopening', async () => {
    const wrapper = render({ active: false })
    window.dispatchEvent(new Event('focus'))
    expect(mocks.getStatus).not.toHaveBeenCalled()
    await wrapper.setProps({ active: true })
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(true)
    await wrapper.setProps({ active: false })
    window.dispatchEvent(new Event('focus'))
    expect(mocks.getStatus).toHaveBeenCalledTimes(1)
    expect(wrapper.find('a').exists()).toBe(false)
    await wrapper.setProps({ active: true })
    await flushPromises()
    expect(mocks.getStatus).toHaveBeenCalledTimes(2)
  })

  it('does not treat failed or incomplete status responses as unconfigured', async () => {
    mocks.getStatus.mockRejectedValueOnce(new Error('network'))
    const wrapper = render()
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(false)
    mocks.getStatus.mockResolvedValueOnce({ data: {} })
    window.dispatchEvent(new Event('focus'))
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(false)
  })

  it('does not let an earlier request override the latest status', async () => {
    const pending = deferred()
    mocks.getStatus.mockReturnValueOnce(pending.promise)
    const wrapper = render()
    mocks.getStatus.mockResolvedValueOnce({ data: { enabled: true } })
    window.dispatchEvent(new Event('focus'))
    await flushPromises()
    pending.resolve({ data: { enabled: false } })
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(false)
  })

  it('ignores a status response after the dialog closes', async () => {
    const pending = deferred()
    mocks.getStatus.mockReturnValueOnce(pending.promise)
    const wrapper = render()
    await wrapper.setProps({ active: false })
    pending.resolve({ data: { enabled: false } })
    await flushPromises()
    expect(wrapper.find('a').exists()).toBe(false)
  })

  it('removes the focus listener on unmount', async () => {
    const wrapper = render()
    await flushPromises()
    wrapper.unmount()
    window.dispatchEvent(new Event('focus'))
    expect(mocks.getStatus).toHaveBeenCalledTimes(1)
  })

  it.each(['1', '2'])('matches the main model route when role type %s has no model management permission', async roleType => {
    mocks.permissions.role_type = roleType
    mocks.permissions.role_permission = []
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('请联系管理员设置')
    expect(wrapper.find('a').exists()).toBe(false)
    expect(window.open).not.toHaveBeenCalled()
  })

  it.each([
    '/clawbot/chat?system_mode=agent&layout=customer-service',
    '/clawbot/settings?layout=customer-service'
  ])('keeps the customer-service Agent layout in the current domain for %s', async path => {
    await router.push(path)
    mocks.company.admin_domain = 'https://admin.example.com'
    mocks.company.agent_domain = 'https://agent.example.com'
    mocks.user.getToken = 'fixture-token'
    mocks.user.userInfo = { user_id: 7, user_name: 'fixture' }
    const wrapper = render()
    await flushPromises()
    expect(wrapper.get('a').attributes('href')).toBe('/user/model?activeTab=3')
    expect(wrapper.find('button').exists()).toBe(false)
    expect(wrapper.html()).not.toContain('fixture-token')
    expect(window.open).not.toHaveBeenCalled()
  })

  it('uses the existing token handoff when opening admin settings from a separate Clawbot domain', async () => {
    await router.push('/clawbot/settings')
    mocks.company.admin_domain = 'https://admin.example.com'
    mocks.company.agent_domain = 'https://agent.example.com'
    mocks.user.getToken = 'fixture-token'
    mocks.user.token = { exp: 123, ttl: 60 }
    mocks.user.userInfo = { user_id: 7, user_name: 'fixture' }
    const wrapper = render()
    await flushPromises()
    const button = wrapper.get('button')
    expect(button.attributes('type')).toBe('button')
    expect(wrapper.find('a').exists()).toBe(false)
    expect(wrapper.html()).not.toContain('fixture-token')
    await button.trigger('click')
    const url = new URL(window.open.mock.calls[0][0])
    expect(url.origin).toBe('https://admin.example.com')
    expect(url.hash.split('?')[0]).toBe('#/set_token')
    const params = new URLSearchParams(url.hash.split('?')[1])
    expect(params.get('redirect_url')).toBe('/user/model?activeTab=3')
    expect(params.get('token')).toBe('fixture-token')
    expect(router.currentRoute.value.path).toBe('/clawbot/settings')
  })
})
