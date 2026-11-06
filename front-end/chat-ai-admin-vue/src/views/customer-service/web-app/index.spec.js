import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { reactive, ref } from 'vue'

const mocks = vi.hoisted(() => ({
  editExternalConfig: vi.fn(),
  getRobot: vi.fn(),
  warning: vi.fn(),
  error: vi.fn(),
  toDataURL: vi.fn()
}))

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/api/robot/index', () => ({ editExternalConfig: mocks.editExternalConfig }))
vi.mock('@/utils/index', () => ({ copyText: vi.fn() }))
vi.mock('qrcode', () => ({ default: { toDataURL: mocks.toDataURL } }))
vi.mock('ant-design-vue', () => ({ message: { warning: mocks.warning, error: mocks.error, success: vi.fn() } }))
vi.mock('@/stores/modules/customer-service', () => ({
  useCustomerServiceStore: () => ({ robot: { id: 7 } })
}))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => mocks.robotStore }))

import WebApp from './index.vue'

const container = { template: '<div><slot /><slot name="content" /></div>' }
const stubs = {
  CustomerPageHeader: false,
  'a-button': { name: 'AButton', template: '<button><slot /></button>' },
  'a-input': { name: 'AInput', props: ['value'], emits: ['update:value'], template: '<input />' },
  'a-popover': container,
  'a-switch': { name: 'ASwitch', props: ['checked'], emits: ['change'], template: '<input />' },
  'a-checkbox': { name: 'ACheckbox', props: ['checked'], emits: ['change'], template: '<input />' },
  'a-textarea': { name: 'ATextarea', props: ['value'], emits: ['update:value'], template: '<textarea />' },
  'a-radio-group': { name: 'ARadioGroup', props: ['value'], emits: ['update:value'], template: '<div><slot /></div>' },
  'a-radio': container,
  'a-select': { name: 'ASelect', props: ['value'], emits: ['update:value'], template: '<div><slot /></div>' },
  'a-select-option': container,
  'a-input-number': { name: 'AInputNumber', props: ['value'], emits: ['update:value'], template: '<input />' },
  'a-space': container,
  'a-empty': container,
  'a-spin': container
}

const mountPage = async () => {
  const wrapper = shallowMount(WebApp, { global: { stubs } })
  await flushPromises()
  return wrapper
}

describe('customer web app settings', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.clearAllMocks()
    mocks.robotStore = {
      robotInfo: ref({
        id: 7,
        robot_key: 'seven',
        robot_name: 'Robot',
        h5_domain: 'https://h5.example.com',
        h5_website: ''
      }),
      external_config_h5: reactive({
        pageTitle: 'Robot',
        ai_generated_tip_show: 1,
        ai_generated_tip: 'AI tip',
        navShow: 1,
        pageStyle: { navbarBackgroundColor: '#2475FC' }
      }),
      getRobot: mocks.getRobot
    }
    mocks.getRobot.mockResolvedValue(true)
    mocks.editExternalConfig.mockResolvedValue({})
    mocks.toDataURL.mockResolvedValue('data:image/png;base64,qr')
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('removes the manual save button and saves changes automatically after the debounce', async () => {
    const wrapper = await mountPage()
    expect(wrapper.find('.save-btn').exists()).toBe(false)

    const checkboxes = wrapper.findAllComponents({ name: 'ACheckbox' })
    checkboxes[0].vm.$emit('change', { target: { checked: false } })
    checkboxes[0].vm.$emit('change', { target: { checked: true } })
    await flushPromises()
    expect(mocks.editExternalConfig).not.toHaveBeenCalled()

    expect(wrapper.find('.auto-save-hint').exists()).toBe(false)
    await vi.advanceTimersByTimeAsync(500)
    await flushPromises()
    expect(mocks.editExternalConfig).toHaveBeenCalledTimes(1)
    expect(wrapper.find('.auto-save-hint').exists()).toBe(false)
    const payload = JSON.parse(mocks.editExternalConfig.mock.calls[0][0].external_config_h5)
    expect(payload.avatarShow).toBe(1)
    wrapper.unmount()
  })

  it('keeps the last local edit when saving without waiting for the debounce', async () => {
    const wrapper = await mountPage()
    const checkboxes = wrapper.findAllComponents({ name: 'ACheckbox' })
    checkboxes[0].vm.$emit('change', { target: { checked: false } })
    await flushPromises()

    wrapper.unmount()
    await flushPromises()
    expect(mocks.editExternalConfig).toHaveBeenCalledTimes(1)
    const payload = JSON.parse(mocks.editExternalConfig.mock.calls[0][0].external_config_h5)
    expect(payload.avatarShow).toBe(2)
  })

  it('does not save an empty AI tip while the tip is enabled', async () => {
    const wrapper = await mountPage()
    const tipInput = wrapper
      .findAllComponents({ name: 'AInput' })
      .find((input) => input.classes().includes('setting-control'))
    tipInput.vm.$emit('update:value', '')
    await flushPromises()

    await vi.advanceTimersByTimeAsync(500)
    await flushPromises()
    expect(mocks.editExternalConfig).not.toHaveBeenCalled()
    expect(mocks.warning).toHaveBeenCalledWith('tip_content_required')
    wrapper.unmount()
  })

  it('flushes pending web app settings when switching header tabs and restores the web app view', async () => {
    const wrapper = await mountPage()
    wrapper.findAllComponents({ name: 'ACheckbox' })[0].vm.$emit('change', { target: { checked: false } })
    await flushPromises()
    await wrapper.findAll('[role="tab"]')[1].trigger('click')
    await flushPromises()
    expect(mocks.editExternalConfig).toHaveBeenCalledTimes(1)
    expect(wrapper.find('.embed-body').exists()).toBe(true)
    expect(wrapper.findAll('[role="tab"]')[1].attributes('aria-selected')).toBe('true')
    await wrapper.findAll('[role="tab"]')[0].trigger('click')
    expect(wrapper.find('.web-app-body').exists()).toBe(true)
    wrapper.unmount()
  })
})
