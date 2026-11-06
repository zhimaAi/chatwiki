import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { computed, reactive } from 'vue'
import { message } from 'ant-design-vue'

const state = vi.hoisted(() => ({ route: null, robotStore: null, editorForm: null }))
const api = vi.hoisted(() => ({ getRate: vi.fn(), saveRate: vi.fn(), validateRate: vi.fn(), validateEditor: vi.fn(), edit: vi.fn(), create: vi.fn(), modelOptions: vi.fn(), groups: vi.fn(), getRobot: vi.fn(), push: vi.fn() }))
vi.mock('@/api/robot/index', () => ({ getRobotRateLimitConf: api.getRate, saveRobotRateLimitConf: api.saveRate, editBaseInfo: api.edit, saveRobot: api.create, getRobotGroupList: api.groups }))
vi.mock('@/api/model/index', () => ({ getModelConfigOption: api.modelOptions }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => state.robotStore }))
vi.mock('@/stores/modules/clawbot', () => ({ useClawbotStore: () => ({ currentAssistant: { robot_key: 'other-agent' } }) }))
vi.mock('pinia', () => ({ storeToRefs: store => ({ robotInfo: computed(() => store.robotInfo) }) }))
vi.mock('vue-router', () => ({ useRoute: () => state.route, useRouter: () => ({ push: api.push }) }))
vi.mock('@/utils/permission', () => ({ getRobotPermission: () => 2 }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('ant-design-vue', () => ({
  message: { success: vi.fn(), error: vi.fn() }, Modal: { confirm: vi.fn() },
  Form: { useForm: form => { state.editorForm = form; return { validate: api.validateEditor, validateInfos: {}, clearValidate: vi.fn() } } }
}))
vi.mock('@/views/robot/robot-list/components/avatar-input.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/components/left-menu.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/components/top-header.vue', () => ({ default: { template: '<div />' } }))

import RateLimit from '@/views/robot/robot-config/rate-limit/index.vue'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'
import RobotLayout from '@/views/robot/robot-config/index.vue'
import robotRoutes from '@/router/robot'

const container = { template: '<div><slot /><slot name="action" /></div>' }
const stubs = {
  'a-form': { ...container, name: 'AForm', props: ['model', 'rules'], setup: (_, { expose }) => { expose({ validate: api.validateRate }) } },
  'a-form-item': container,
  'a-modal': { ...container, name: 'AModal', props: ['open'], emits: ['ok', 'cancel'] },
  'a-button': { template: '<button><slot /></button>' },
  'a-input-number': { inheritAttrs: false, props: ['value'], template: '<input />' },
  'a-input': { inheritAttrs: false, props: ['value'], template: '<input />' },
  'a-select': { ...container, inheritAttrs: false, props: ['value'] },
  'a-alert': container
}
const mountRate = (props = {}) => shallowMount(RateLimit, { props, global: { stubs } })
const lastStatus = wrapper => wrapper.emitted('status-change').at(-1)[0]

describe('function center integration with shared pages', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    state.route = reactive({ name: 'customerServiceFunctions', query: { id: '7', robot_key: 'customer-seven' } })
    state.robotStore = { robotInfo: reactive({ id: 7, robot_key: 'customer-seven', robot_name: 'Customer', robot_avatar_url: 'customer.png', application_type: 0 }), getRobot: api.getRobot }
    api.getRate.mockReset().mockResolvedValue({ res: 0, data: { switch_status: 0, five_minute_limit: 8, five_minute_reply_type: 1, five_minute_reply_content: 'Slow down', daily_limit: 80, daily_reply_type: 0 } })
    api.saveRate.mockReset().mockResolvedValue({ res: 0 })
    api.validateRate.mockReset().mockResolvedValue()
    api.validateEditor.mockReset().mockResolvedValue()
    api.modelOptions.mockResolvedValue({ data: [{}] })
    api.groups.mockResolvedValue({ data: [] })
    api.edit.mockResolvedValue({ res: 0 })
    api.create.mockResolvedValue({ res: 0, data: { id: 900, robot_key: 'new-flow' } })
  })

  it('loads the explicit robot key and persists the card switch with the existing rate configuration', async () => {
    const wrapper = mountRate({ robotKey: 'explicit-customer', embedded: true })
    await flushPromises()
    expect(api.getRate).toHaveBeenCalledWith({ robot_key: 'explicit-customer' })
    expect(lastStatus(wrapper)).toMatchObject({ ready: true, enabled: false, loading: false })
    expect(await wrapper.vm.setEnabled(true)).toBe(true)
    await flushPromises()
    expect(api.saveRate).toHaveBeenCalledWith({ robot_key: 'explicit-customer', switch_status: 1, five_minute_limit: 8, five_minute_reply_type: 1, five_minute_reply_content: 'Slow down', daily_limit: 80, daily_reply_type: 0, daily_reply_content: '' })
    expect(lastStatus(wrapper).enabled).toBe(true)
    wrapper.unmount()
  })

  it('keeps the legacy route fallback and rolls the switch back after a failed save', async () => {
    api.saveRate.mockRejectedValueOnce(new Error('network'))
    const wrapper = mountRate()
    await flushPromises()
    expect(api.getRate).toHaveBeenCalledWith({ robot_key: 'customer-seven' })
    expect(await wrapper.vm.setEnabled(true)).toBe(false)
    await flushPromises()
    expect(lastStatus(wrapper)).toMatchObject({ enabled: false, loading: false })
    expect(message.error).toHaveBeenCalledWith('msg_save_failed')
    expect(await wrapper.vm.setEnabled(true)).toBe(true)
    wrapper.unmount()
  })

  it('validates before enabling, rejects duplicate toggles and does not save invalid settings', async () => {
    let finish
    api.validateRate.mockImplementationOnce(() => new Promise((_, reject) => { finish = reject }))
    const wrapper = mountRate()
    await flushPromises()
    const first = wrapper.vm.setEnabled(true)
    expect(await wrapper.vm.setEnabled(true)).toBe(false)
    const form = wrapper.findComponent({ name: 'AForm' })
    const rules = form.props('rules')
    await expect(rules.five_minute_limit[0].validator({}, 0)).rejects.toThrow('validator_positive_integer')
    await expect(rules.five_minute_reply_content[0].validator({}, '')).rejects.toThrow('validator_reply_required')
    finish(new Error('invalid'))
    expect(await first).toBe(false)
    expect(api.saveRate).not.toHaveBeenCalled()
    await flushPromises()
    expect(lastStatus(wrapper)).toMatchObject({ enabled: false, loading: false })
    wrapper.unmount()
  })

  it('does not save after leaving while validation is pending, or apply stale load responses', async () => {
    let finishValidation
    api.validateRate.mockImplementationOnce(() => new Promise(resolve => { finishValidation = resolve }))
    const wrapper = mountRate()
    await flushPromises()
    const pending = wrapper.vm.setEnabled(true)
    wrapper.unmount()
    finishValidation()
    await pending
    expect(api.saveRate).not.toHaveBeenCalled()
    let finishLoad
    api.getRate.mockImplementationOnce(() => new Promise(resolve => { finishLoad = resolve }))
    const status = vi.fn()
    const loading = mountRate({ onStatusChange: status })
    const events = status.mock.calls.length
    loading.unmount()
    finishLoad({ res: 0, data: { switch_status: 1 } })
    await flushPromises()
    expect(status).toHaveBeenCalledTimes(events)
  })

  it('edits the workflow card without overwriting the current customer or its identifier', async () => {
    const wrapper = shallowMount(AddRobotAlert, { global: { stubs } })
    await wrapper.vm.open(1, true, { id: 900, en_name: 'flow', robot_name: 'Workflow', robot_intro: 'Intro', robot_avatar: 'flow.png', group_id: '5' })
    expect(state.editorForm).toMatchObject({ robot_name: 'Workflow', en_name: 'flow', group_id: '5', robot_avatar_url: 'flow.png' })
    state.editorForm.robot_name = 'Edited workflow'
    wrapper.findComponent({ name: 'AModal' }).vm.$emit('ok')
    await flushPromises()
    expect(api.edit).toHaveBeenCalledWith(expect.objectContaining({ id: 900, robot_name: 'Edited workflow', group_id: '5' }), 1)
    expect(api.getRobot).not.toHaveBeenCalled()
    expect(state.robotStore.robotInfo.robot_name).toBe('Customer')
    expect(wrapper.emitted('editRobot')).toHaveLength(1)
    wrapper.unmount()
  })

  it('retains legacy robot editing and lets the workspace handle workflow creation redirects', async () => {
    const opened = vi.spyOn(window, 'open').mockImplementation(() => null)
    const wrapper = shallowMount(AddRobotAlert, { props: { createRedirect: 'none' }, global: { stubs } })
    await wrapper.vm.open(0, true)
    wrapper.findComponent({ name: 'AModal' }).vm.$emit('ok')
    await flushPromises()
    expect(api.edit).toHaveBeenCalledWith(expect.objectContaining({ id: '7', robot_name: 'Customer' }), 0)
    expect(api.getRobot).toHaveBeenCalledWith('7')
    await wrapper.vm.open(1)
    state.editorForm.en_name = 'new_flow'
    wrapper.findComponent({ name: 'AModal' }).vm.$emit('ok')
    await flushPromises()
    expect(wrapper.emitted('ok')[0][0]).toEqual({ id: 900, robot_key: 'new-flow' })
    expect(opened).not.toHaveBeenCalled()
    wrapper.unmount()
    opened.mockRestore()
  })

  it('uses the existing workflow route without a sidebar and preserves the legacy configuration sidebar', async () => {
    const workflowRoute = robotRoutes.children.find(item => item.path === '/robot/config').children.find(item => item.path === '/robot/config/workflow')
    state.route.name = workflowRoute.name
    const wrapper = shallowMount(RobotLayout, { global: { stubs: { 'router-view': container } } })
    expect(wrapper.find('.layout-left').exists()).toBe(false)
    state.route.name = 'basicConfig'
    await flushPromises()
    expect(wrapper.find('.layout-left').exists()).toBe(true)
    wrapper.unmount()
  })
})
