import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive, ref } from 'vue'
import ChatModeCard from './components/chat-mode-card.vue'
import LegacyChatMode from '@/views/robot/robot-config/basic-config/components/chat-mode/index.vue'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('./configuration', () => ({ saveConfigurationPatch: (state, patch) => state.savePatch(patch) }))
const container = { template: '<div><slot /></div>' }
const stubs = {
  'svg-icon': true, 'a-tooltip': container,
  'a-button': { template: '<button><slot /></button>' },
  'a-radio-group': { name: 'ARadioGroup', props: ['value'], emits: ['change'], template: '<div><slot /></div>' },
  'a-radio': { props: ['value'], template: '<label><slot /></label>' },
  'a-checkbox': { name: 'ACheckbox', props: ['checked', 'disabled'], emits: ['change'], template: '<input type="checkbox" :disabled="disabled" />' },
  'a-input-number': { name: 'AInputNumber', props: ['value', 'disabled', 'min', 'max', 'precision'], emits: ['change', 'blur'], template: '<input />' },
  'a-alert': { name: 'AAlert', template: '<div><slot name="action" /></div>' },
  'a-modal': { name: 'AModal', props: ['open', 'confirmLoading', 'closable', 'cancelButtonProps'], emits: ['ok', 'cancel'], template: '<section v-if="open"><slot /></section>' },
  'a-switch': true
}
function context() {
  const state = {
    robotInfo: reactive({ id: 7, robot_key: 'seven', application_type: 2, chat_type: '3', library_qa_direct_reply_switch: 'false', library_qa_direct_reply_score: '0.850', mixture_qa_direct_reply_switch: 'true', mixture_qa_direct_reply_score: '0.900' }),
    savedRevision: ref(0), updateRobotInfo: vi.fn(), savePatch: vi.fn()
  }
  state.updateRobotInfo.mockImplementation(async patch => { Object.assign(state.robotInfo, patch); state.savedRevision.value++; return true })
  state.savePatch.mockImplementation(async patch => { Object.assign(state.robotInfo, patch); state.savedRevision.value++; return true })
  return state
}
const options = state => ({ global: { provide: { customerRobotContext: state, robotInfo: { robotInfo: state.robotInfo, updateRobotInfo: state.updateRobotInfo } }, stubs } })
const named = (wrapper, name) => wrapper.findComponent({ name })
const deferred = () => { let resolve; const promise = new Promise(done => { resolve = done }); return { promise, resolve } }
const openDialog = async wrapper => { await wrapper.find('.edit-box-header button').trigger('click') }
const confirmDialog = async wrapper => { named(wrapper, 'AModal').vm.$emit('ok'); await flushPromises() }

describe('customer chat mode summary and dialog', () => {
  beforeEach(() => vi.clearAllMocks())

  it('shows a summary and stages mode changes until the dialog is confirmed', async () => {
    const state = context()
    const wrapper = mount(ChatModeCard, options(state))
    expect(wrapper.findAll('.mode-option')).toHaveLength(0)
    expect(wrapper.find('.summary').text()).toContain('chat_mode_mixed')
    expect(wrapper.find('.summary').text()).toContain('0.9')
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    await openDialog(wrapper)
    expect(wrapper.findAll('.mode-option')).toHaveLength(3)
    await wrapper.findAll('.mode-option')[0].trigger('click')
    await flushPromises()
    expect(state.savePatch).not.toHaveBeenCalled()
    expect(wrapper.find('.summary').text()).toContain('chat_mode_mixed')
    await confirmDialog(wrapper)
    expect(state.savePatch).toHaveBeenCalledWith({ chat_type: '1', library_qa_direct_reply_switch: 'false', library_qa_direct_reply_score: '0.85', mixture_qa_direct_reply_switch: 'true', mixture_qa_direct_reply_score: '0.9' })
    expect(wrapper.find('.summary').text()).toContain('chat_mode_knowledge_base_only')
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    await openDialog(wrapper)
    await wrapper.findAll('.mode-option')[2].trigger('click')
    await flushPromises()
    expect(wrapper.find('.qa-settings').classes()).toContain('is-disabled')
    expect(named(wrapper, 'ACheckbox').props('disabled')).toBe(true)
    expect(named(wrapper, 'AInputNumber').props('disabled')).toBe(true)
    await confirmDialog(wrapper)
    expect(wrapper.find('.summary').text()).toContain('chat_mode_direct')
    expect(wrapper.find('.summary').text()).not.toContain('label_direct_reply_qa')
    expect(state.robotInfo.application_type).toBe(2)
    wrapper.unmount()
  })

  it('keeps QA switches and thresholds separate for library and mixed modes', async () => {
    const state = context()
    const wrapper = mount(ChatModeCard, options(state))
    await openDialog(wrapper)
    named(wrapper, 'AInputNumber').vm.$emit('change', 0.95)
    await flushPromises()
    expect(state.robotInfo.mixture_qa_direct_reply_score).toBe('0.900')
    await wrapper.findAll('.mode-option')[0].trigger('click')
    await flushPromises()
    expect(named(wrapper, 'AInputNumber').props('disabled')).toBe(true)
    named(wrapper, 'ACheckbox').vm.$emit('change', { target: { checked: true } })
    await flushPromises()
    expect(named(wrapper, 'AInputNumber').props('disabled')).toBe(false)
    named(wrapper, 'AInputNumber').vm.$emit('change', 0.8)
    await flushPromises()
    expect(state.savePatch).not.toHaveBeenCalled()
    await wrapper.findAll('.mode-option')[1].trigger('click')
    await flushPromises()
    expect(named(wrapper, 'AInputNumber').props('value')).toBe(0.95)
    await confirmDialog(wrapper)
    expect(state.robotInfo.library_qa_direct_reply_score).toBe('0.8')
    expect(state.robotInfo.mixture_qa_direct_reply_score).toBe('0.95')
    wrapper.unmount()
  })

  it('does not submit empty or out-of-range thresholds and restores the saved value on blur', async () => {
    const state = context()
    const wrapper = mount(ChatModeCard, options(state))
    await openDialog(wrapper)
    for (const value of [null, -0.1, 1.1]) {
      named(wrapper, 'AInputNumber').vm.$emit('change', value)
      await flushPromises()
      await confirmDialog(wrapper)
      expect(state.savePatch).not.toHaveBeenCalled()
      named(wrapper, 'AInputNumber').vm.$emit('blur')
      await flushPromises()
      expect(named(wrapper, 'AInputNumber').props('value')).toBe(0.9)
    }
    wrapper.unmount()
  })

  it('discards edits when cancelled and reopens with the saved configuration', async () => {
    const state = context()
    const wrapper = mount(ChatModeCard, options(state))
    await openDialog(wrapper)
    named(wrapper, 'AInputNumber').vm.$emit('change', 0.95)
    await wrapper.findAll('.mode-option')[0].trigger('click')
    named(wrapper, 'AModal').vm.$emit('cancel')
    await flushPromises()
    expect(state.savePatch).not.toHaveBeenCalled()
    expect(wrapper.find('.summary').text()).toContain('chat_mode_mixed')
    await openDialog(wrapper)
    expect(wrapper.findAll('.mode-option')[1].classes()).toContain('active')
    expect(named(wrapper, 'AInputNumber').props('value')).toBe(0.9)
    wrapper.unmount()
  })

  it('retains failed edits for retry and keeps the last saved summary', async () => {
    const state = context()
    state.savePatch.mockResolvedValueOnce(false)
    const wrapper = mount(ChatModeCard, options(state))
    await openDialog(wrapper)
    await wrapper.findAll('.mode-option')[2].trigger('click')
    await confirmDialog(wrapper)
    expect(wrapper.findAll('.mode-option')[2].classes()).toContain('active')
    expect(wrapper.find('.summary').text()).toContain('chat_mode_mixed')
    expect(named(wrapper, 'AAlert').exists()).toBe(true)
    await named(wrapper, 'AAlert').find('button').trigger('click')
    await flushPromises()
    expect(state.robotInfo.chat_type).toBe('2')
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    expect(wrapper.find('.summary').text()).toContain('chat_mode_direct')
    wrapper.unmount()
  })

  it('prevents duplicate submissions while waiting for saving and refresh', async () => {
    const state = context()
    const pending = deferred()
    state.savePatch.mockReturnValueOnce(pending.promise)
    const wrapper = mount(ChatModeCard, options(state))
    await openDialog(wrapper)
    await confirmDialog(wrapper)
    expect(named(wrapper, 'AModal').props('confirmLoading')).toBe(true)
    expect(named(wrapper, 'AModal').props('closable')).toBe(false)
    await confirmDialog(wrapper)
    expect(state.savePatch).toHaveBeenCalledOnce()
    pending.resolve(true)
    await flushPromises()
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    wrapper.unmount()
  })

  it('keeps the legacy summary and edit modal until the user confirms changes', async () => {
    const state = context()
    const wrapper = mount(LegacyChatMode, options(state))
    expect(wrapper.find('.setting-info-block').exists()).toBe(true)
    await wrapper.find('.edit-box-header button').trigger('click')
    expect(state.updateRobotInfo).not.toHaveBeenCalled()
    expect(named(wrapper, 'AModal').props('open')).toBe(true)
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(state.updateRobotInfo).toHaveBeenCalledOnce()
    expect(named(wrapper, 'AModal').props('open')).toBe(false)
    wrapper.unmount()
  })

  it('keeps direct-mode controls inert and restores each QA draft when switching back', async () => {
    const state = context()
    const wrapper = mount(ChatModeCard, options(state))
    await openDialog(wrapper)
    named(wrapper, 'AInputNumber').vm.$emit('change', 0.94)
    await wrapper.findAll('.mode-option')[2].trigger('click')
    named(wrapper, 'AInputNumber').vm.$emit('change', 0.1)
    named(wrapper, 'ACheckbox').vm.$emit('change', { target: { checked: false } })
    await wrapper.findAll('.mode-option')[1].trigger('click')
    expect(named(wrapper, 'AInputNumber').props('value')).toBe(0.94)
    expect(named(wrapper, 'ACheckbox').props('checked')).toBe(true)
    await wrapper.findAll('.mode-option')[0].trigger('click')
    expect(named(wrapper, 'AInputNumber').props('value')).toBe(0.85)
    expect(named(wrapper, 'ACheckbox').props('checked')).toBe(false)
    expect(state.savePatch).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('supports arrow keys and Home/End without submitting until confirmation', async () => {
    const state = context()
    const wrapper = mount(ChatModeCard, { ...options(state), attachTo: document.body })
    await openDialog(wrapper)
    await wrapper.findAll('.mode-option')[1].trigger('keydown', { key: 'ArrowRight' })
    expect(wrapper.findAll('.mode-option')[2].attributes('aria-checked')).toBe('true')
    expect(document.activeElement).toBe(wrapper.findAll('.mode-option')[2].element)
    await wrapper.findAll('.mode-option')[2].trigger('keydown', { key: 'Home' })
    expect(wrapper.findAll('.mode-option')[0].attributes('tabindex')).toBe('0')
    await wrapper.findAll('.mode-option')[0].trigger('keydown', { key: 'End' })
    expect(wrapper.findAll('.mode-option')[2].attributes('aria-checked')).toBe('true')
    expect(state.savePatch).not.toHaveBeenCalled()
    wrapper.unmount()
  })
})
