import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const api = vi.hoisted(() => ({ getMode: vi.fn(), saveMode: vi.fn(), replace: vi.fn() }))
vi.mock('@/api/user', () => ({ getSystemMode: api.getMode, saveSystemMode: api.saveMode }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => ({ id: 1 }) }))
vi.mock('vue-router', () => ({ useRouter: () => ({ replace: api.replace, resolve: () => ({ href: '/#/' }) }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))

import Setting from '@/views/user/enterprise/components/system-mode-setting.vue'
import { SYSTEM_MODES, useSystemModeStore } from '@/stores/modules/system-mode'

const mountSetting = (onboarding = true) => shallowMount(Setting, {
  props: { onboarding },
  global: { stubs: {
    'a-button': { props: ['disabled', 'loading'], template: '<button :disabled="disabled || loading"><slot /></button>' },
    'a-modal': { name: 'AModal', props: ['open'], emits: ['ok', 'cancel'], template: '<div />' }
  } }
})

describe('first enterprise mode selection', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
    api.getMode.mockReset().mockResolvedValue({ data: { system_mode: SYSTEM_MODES.AI_SERVICE, system_mode_selected: 0 } })
    api.saveMode.mockReset().mockResolvedValue({ res: 0 })
    api.replace.mockResolvedValue()
  })

  it('defaults to AI service without saving and starts even when the selected mode is unchanged', async () => {
    const wrapper = mountSetting()
    await flushPromises()
    expect(wrapper.find('input:checked').element.value).toBe(SYSTEM_MODES.AI_SERVICE)
    expect(api.saveMode).not.toHaveBeenCalled()
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(api.saveMode).toHaveBeenCalledWith({ system_mode: SYSTEM_MODES.AI_SERVICE })
    expect(useSystemModeStore().needsModeSelection).toBe(false)
    expect(api.replace).toHaveBeenCalledWith('/')
    wrapper.unmount()
  })

  it('stages another choice, preserves it on failure, and only navigates after successful retry', async () => {
    api.saveMode.mockRejectedValueOnce(new Error('network'))
    const wrapper = mountSetting()
    await flushPromises()
    await wrapper.find('input[value="agent"]').setValue(true)
    expect(api.saveMode).not.toHaveBeenCalled()
    expect(wrapper.findComponent({ name: 'AModal' }).props('open')).toBe(false)
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(api.replace).not.toHaveBeenCalled()
    expect(wrapper.find('input:checked').element.value).toBe(SYSTEM_MODES.AGENT)
    expect(useSystemModeStore().needsModeSelection).toBe(true)
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(api.saveMode).toHaveBeenLastCalledWith({ system_mode: SYSTEM_MODES.AGENT })
    expect(api.replace).toHaveBeenCalledWith('/')
    wrapper.unmount()
  })

  it('does not repeat onboarding after it has been completed on another login', async () => {
    api.getMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.AGENT, system_mode_selected: 1 } })
    const wrapper = mountSetting()
    await flushPromises()
    expect(api.replace).toHaveBeenCalledWith('/')
    expect(api.saveMode).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('blocks starting until loading succeeds and offers retry', async () => {
    api.getMode.mockRejectedValueOnce(new Error('network'))
    const wrapper = mountSetting()
    await flushPromises()
    const start = () => wrapper.findAll('button').find(button => button.text() === 'start_now')
    expect(start().attributes('disabled')).toBeDefined()
    await wrapper.findAll('button').find(button => button.text() === 'retry').trigger('click')
    await flushPromises()
    expect(start().attributes('disabled')).toBeUndefined()
    expect(api.saveMode).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('keeps legacy settings as a confirmation flow without a get-started button', async () => {
    api.getMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.ADVANCED, system_mode_selected: 1 } })
    const wrapper = mountSetting(false)
    await flushPromises()
    expect(wrapper.find('.start-action').exists()).toBe(false)
    await wrapper.find('input[value="agent"]').setValue(true)
    expect(wrapper.findComponent({ name: 'AModal' }).props('open')).toBe(true)
    expect(api.saveMode).not.toHaveBeenCalled()
    wrapper.findComponent({ name: 'AModal' }).vm.$emit('cancel')
    await flushPromises()
    expect(wrapper.find('input:checked').element.value).toBe(SYSTEM_MODES.ADVANCED)
    wrapper.unmount()
  })
})
