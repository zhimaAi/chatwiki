import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { createPinia } from 'pinia'
import ChatModelSelect from './chat-model-select.vue'
import AgentMessageInput from './message-input.vue'
import CustomerMessageInput from '@/views/customer-service/chat/components/message-input.vue'
import noticeZh from '@/locales/lang/zh-CN/components/model-select/model-select.json'

const api = vi.hoisted(() => ({ options: vi.fn(), status: vi.fn() }))
vi.mock('@/api/model/index', () => ({ getModelConfigOption: api.options, getBackupModelStatus: api.status }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => noticeZh[key] || key }) }))
vi.mock('@/stores/modules/permission', () => ({ usePermissionStore: () => ({ role_permission: ['ModelManage'] }) }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => ({ admin_domain: '', agent_domain: '' }) }))
vi.mock('@/stores/modules/user', () => ({ useUserStore: () => ({ getToken: '', userInfo: {} }) }))
vi.mock('@/hooks/web/useUpload.js', () => ({ useUpload: () => ({ openFileDialog: vi.fn() }) }))
vi.mock('./file-toolbar.vue', () => ({ default: { template: '<div />' } }))

const wrappers = []
const modelOptions = () => [{
  model_define: 'openai', model_name: 'Fixture Provider',
  config_info: { id: 7, config_name: 'Fixture Provider' },
  use_model_configs: [
    { id: 11, use_model_name: 'first', show_model_name: 'First model' },
    { id: 12, use_model_name: 'second', show_model_name: 'Second model' }
  ]
}]
async function render(component, path, props = {}) {
  const router = createRouter({ history: createMemoryHistory(), routes: [
    { path: '/customer-service/chat', component: { template: '<div />' } },
    { path: '/clawbot/chat', component: { template: '<div />' } },
    { path: '/user/model', component: { template: '<div />' } }
  ] })
  await router.push(path)
  const wrapper = mount(component, { props: { modeId: 7, modeName: 'first', ...props }, global: {
    plugins: [createPinia(), router], stubs: { ATextarea: { template: '<textarea />' } }
  } })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}
beforeEach(() => {
  vi.clearAllMocks()
  api.options.mockImplementation(() => Promise.resolve({ data: modelOptions() }))
  api.status.mockResolvedValue({ data: { enabled: false } })
})
afterEach(() => wrappers.splice(0).forEach(wrapper => wrapper.unmount()))

describe.each([
  ['AI customer service', CustomerMessageInput, '/customer-service/chat'],
  ['Agent', AgentMessageInput, '/clawbot/chat?system_mode=agent&layout=customer-service']
])('%s mode actual message input', (_mode, component, path) => {
  it('shows the shared setup notice only when opening the model selector and reloads it on reopening', async () => {
    const wrapper = await render(component, path)
    expect(api.status).not.toHaveBeenCalled()
    await wrapper.get('.chat-model-trigger').trigger('click')
    await flushPromises()
    expect(wrapper.get('.backup-model-tip').text()).toContain('备用模型暂未开启，')
    expect(wrapper.get('.backup-model-tip a').attributes()).toMatchObject({ href: '/user/model?activeTab=3', target: '_blank', rel: 'noopener noreferrer' })
    await wrapper.get('.chat-model-trigger').trigger('click')
    window.dispatchEvent(new Event('focus'))
    await flushPromises()
    expect(wrapper.find('.chat-model-panel').exists()).toBe(false)
    expect(api.status).toHaveBeenCalledTimes(1)
    api.status.mockResolvedValueOnce({ data: { enabled: true } })
    await wrapper.get('.chat-model-trigger').trigger('click')
    await flushPromises()
    expect(api.status).toHaveBeenCalledTimes(2)
    expect(wrapper.find('.backup-model-tip').exists()).toBe(false)
  })

  it('hides the notice when configured and does not report failed reads as unconfigured', async () => {
    api.status.mockResolvedValueOnce({ data: { enabled: true } })
    const wrapper = await render(component, path)
    await wrapper.get('.chat-model-trigger').trigger('click')
    await flushPromises()
    expect(wrapper.find('.backup-model-tip').exists()).toBe(false)
    await wrapper.get('.chat-model-trigger').trigger('click')
    api.status.mockRejectedValueOnce(new Error('isolated failure'))
    await wrapper.get('.chat-model-trigger').trigger('click')
    await flushPromises()
    expect(wrapper.find('.backup-model-tip').exists()).toBe(false)
  })

  it('preserves selecting a model and closes the notice together with the selector', async () => {
    const wrapper = await render(component, path)
    await wrapper.get('.chat-model-trigger').trigger('click')
    await flushPromises()
    const buttons = wrapper.findAll('.model-item')
    expect(buttons).toHaveLength(2)
    await buttons[1].trigger('click')
    expect(wrapper.emitted('changeModel')[0][0]).toMatchObject({ model_config_id: '7', name: 'second' })
    expect(wrapper.find('.chat-model-panel').exists()).toBe(false)
    expect(wrapper.find('.backup-model-tip').exists()).toBe(false)
  })
})

it('does not request or display a backup notice for non-LLM model selection', async () => {
  const wrapper = await render(ChatModelSelect, '/customer-service/chat', { modelType: 'RERANK' })
  await wrapper.get('.chat-model-trigger').trigger('click')
  await flushPromises()
  expect(api.status).not.toHaveBeenCalled()
  expect(wrapper.find('.panel-backup-notice').exists()).toBe(false)
})
