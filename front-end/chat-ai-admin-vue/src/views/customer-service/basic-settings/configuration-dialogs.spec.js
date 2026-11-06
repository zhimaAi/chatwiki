import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive } from 'vue'
import VariableModal from '@/views/robot/robot-config/basic-config/components/variable-setting/add-variable-modal.vue'
import ProblemOptimization from '@/views/robot/robot-config/basic-config/components/problem-optimization.vue'
import SuggestedIssues from '@/views/robot/robot-config/basic-config/components/suggested-issues.vue'

const mocks = vi.hoisted(() => ({ createVariable: vi.fn(), saveRobot: vi.fn(), validate: vi.fn(), clearValidate: vi.fn(), error: vi.fn() }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: { robot_key: 'legacy-query' } }) }))
vi.mock('@/api/robot/index', () => ({ createChatVariable: mocks.createVariable, saveRobot: mocks.saveRobot }))
vi.mock('@/api/robot', () => ({ createChatVariable: mocks.createVariable, saveRobot: mocks.saveRobot }))
vi.mock('ant-design-vue', async importOriginal => ({ Modal: (await importOriginal()).Modal, message: { error: mocks.error, success: vi.fn() } }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { name: 'ModelSelect', props: ['modeName', 'modeId'], emits: ['update:modeName', 'update:modeId', 'loaded'], template: '<div />' } }))
vi.mock('@/components/reasoning-effort-select/index.vue', () => ({ default: { name: 'ReasoningEffortSelect', props: ['value'], emits: ['update:value', 'change'], template: '<input />' } }))

const container = { template: '<div><slot /></div>' }
const stubs = {
  'svg-icon': true, 'a-flex': container, 'a-tooltip': container, 'a-select-option': container,
  'a-button': { template: '<button><slot /></button>' },
  'a-modal': { name: 'AModal', props: ['open', 'width', 'confirmLoading', 'closable', 'getContainer'], emits: ['ok', 'cancel'], template: '<section v-if="open" role="dialog"><slot /></section>' },
  'a-form': { name: 'AForm', props: ['model'], setup(_props, { expose }) { expose({ validate: mocks.validate, clearValidate: mocks.clearValidate }) }, template: '<form><slot /></form>' },
  'a-form-item': { name: 'AFormItem', props: ['name', 'rules'], template: '<div><slot /></div>' },
  'a-form-item-rest': container,
  'a-input': { name: 'AInput', props: ['value'], emits: ['update:value'], template: '<input :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-input-number': { name: 'AInputNumber', props: ['value', 'size'], emits: ['update:value'], template: '<input type="number" />' },
  'a-textarea': { name: 'ATextarea', props: ['value'], emits: ['update:value'], template: '<textarea :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-select': { name: 'ASelect', props: ['value', 'options'], emits: ['update:value', 'change'], template: '<select><slot /></select>' },
  'a-checkbox': { name: 'ACheckbox', props: ['checked'], emits: ['update:checked'], template: '<input type="checkbox" />' },
  'a-switch': { name: 'ASwitch', props: ['checked'], emits: ['update:checked', 'change'], template: '<input type="checkbox" />' },
  'a-radio-group': { name: 'ARadioGroup', props: ['value'], emits: ['update:value', 'change'], template: '<div><slot /></div>' },
  'a-radio': { template: '<label><slot /></label>' },
  'a-slider': { name: 'ASlider', props: ['value'], emits: ['change'], template: '<input type="range" />' },
  'a-alert': { template: '<div role="alert" />' },
  draggable: { props: ['modelValue'], template: '<div><template v-for="(element, index) in modelValue" :key="index"><slot name="item" :element="element" :index="index" /></template></div>' }
}
function render(component, compact = true, props = {}) {
  const robotInfo = reactive({ id: 7, robot_key: 'workspace-seven', application_type: 0, use_model: 'main-model', model_config_id: 2,
    optimize_question_model_config_id: 0, optimize_question_use_model: '', optimize_question_enable_thinking: 1,
    optimize_question_reasoning_effort: 'medium', optimize_question_dialogue_background: '原对话背景',
    question_guide_mode: 0, question_guide_prompt: '原提示词', question_guide_workflow_key: '', question_guide_num: 3 })
  const context = { robotInfo, saveAction: vi.fn(async task => { await task(); return true }) }
  const updateRobotInfo = vi.fn()
  const wrapper = mount(component, { props: { compact, ...props }, global: {
    stubs, provide: { customerRobotContext: context, robotInfo: { robotInfo, updateRobotInfo } }
  } })
  return { wrapper, robotInfo, context, updateRobotInfo }
}
const named = (wrapper, name) => wrapper.findComponent({ name })
const confirm = async wrapper => { named(wrapper, 'AModal').vm.$emit('ok'); await flushPromises() }
const cancel = async wrapper => { named(wrapper, 'AModal').vm.$emit('cancel'); await flushPromises() }
const openSetting = wrapper => wrapper.find('.edit-box-header button').trigger('click')
const chooseRadio = async (wrapper, value) => {
  const radio = named(wrapper, 'ARadioGroup')
  radio.vm.$emit('update:value', value)
  radio.vm.$emit('change', { target: { value } })
  await flushPromises()
}
const field = (wrapper, name) => wrapper.findAllComponents({ name: 'AFormItem' }).find(item => item.props('name') === name)
const setField = async (wrapper, name, value, control = 'AInput') => {
  named(field(wrapper, name), control).vm.$emit('update:value', value)
  await flushPromises()
}
const chooseType = async (wrapper, value) => {
  const control = named(field(wrapper, 'variable_type'), 'ASelect')
  control.vm.$emit('update:value', value)
  control.vm.$emit('change', value)
  await flushPromises()
}
const savedPatch = () => mocks.saveRobot.mock.calls.at(-1)[0]
const thinkingMap = { '2#main-model': true, '9#other-model': true }
const loadModels = async wrapper => {
  named(wrapper, 'ModelSelect').vm.$emit('loaded', [{ model_define: 'custom', children: [{ name: 'other-model', model_config_id: 9 }] }], thinkingMap)
  await flushPromises()
}

describe('workspace configuration dialogs', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.validate.mockResolvedValue()
    mocks.createVariable.mockResolvedValue({ res: 0 })
    mocks.saveRobot.mockResolvedValue({ res: 0 })
  })

  it('uses the workspace robot key and discards variable option edits on cancel', async () => {
    const { wrapper } = render(VariableModal)
    const data = { id: 3, variable_type: 'select_one', variable_key: 'level', variable_name: '会员等级', options: [{ label: '普通会员' }] }
    wrapper.vm.show(data)
    await flushPromises()
    expect(named(wrapper, 'AModal').props('width')).toBe(560)
    await wrapper.find('.option-item input').setValue('高级会员')
    expect(data.options).toEqual([{ label: '普通会员' }])
    await cancel(wrapper)
    expect(mocks.createVariable).not.toHaveBeenCalled()
    wrapper.vm.show(data)
    await flushPromises()
    expect(wrapper.find('.option-item input').element.value).toBe('普通会员')
    await confirm(wrapper)
    expect(mocks.createVariable.mock.calls[0][0]).toMatchObject({ robot_key: 'workspace-seven', id: 3, options: '[{"label":"普通会员"}]', max_input_length: 1 })
    wrapper.unmount()
  })

  it('preserves key/form validation and serializes a numeric variable with its required state', async () => {
    const { wrapper } = render(VariableModal)
    wrapper.vm.show()
    await flushPromises()
    const validator = field(wrapper, 'variable_key').props('rules')[0].validator
    await expect(validator({}, '会员1')).rejects.toBe('msg_key_invalid')
    await expect(validator({}, 'member_level')).resolves.toBeUndefined()
    mocks.validate.mockRejectedValueOnce(new Error('required'))
    await confirm(wrapper)
    expect(mocks.createVariable).not.toHaveBeenCalled()
    await chooseType(wrapper, 'input_number')
    await setField(wrapper, 'variable_key', 'level')
    await setField(wrapper, 'variable_name', '等级')
    await setField(wrapper, 'max_input_length', 20, 'AInputNumber')
    await setField(wrapper, 'default_value', '123', 'AInputNumber')
    named(wrapper, 'ACheckbox').vm.$emit('update:checked', true)
    await confirm(wrapper)
    expect(mocks.createVariable.mock.calls[0][0]).toMatchObject({ variable_type: 'input_number', default_value: '123', max_input_length: 20, must_input: 1, options: '[]' })
    wrapper.unmount()
  })

  it('validates select options and switches to a checkbox default without the required control', async () => {
    const { wrapper } = render(VariableModal)
    wrapper.vm.show()
    await flushPromises()
    await chooseType(wrapper, 'select_one')
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenCalledWith('msg_fill_options')
    expect(mocks.createVariable).not.toHaveBeenCalled()
    await chooseType(wrapper, 'checkbox_switch')
    expect(field(wrapper, 'must_input')).toBeUndefined()
    await confirm(wrapper)
    expect(mocks.createVariable.mock.calls[0][0]).toMatchObject({ variable_type: 'checkbox_switch', default_value: '2', max_input_length: 1 })
    wrapper.unmount()
  })

  it('keeps variable drafts after failure and locks the dialog during a retry', async () => {
    const { wrapper, context } = render(VariableModal)
    wrapper.vm.show({ variable_key: 'member', variable_name: '会员' })
    await flushPromises()
    context.saveAction.mockResolvedValueOnce(false)
    await confirm(wrapper)
    expect(wrapper.find('[role="alert"]').exists()).toBe(true)
    let resolve
    context.saveAction.mockImplementationOnce(task => { task(); return new Promise(done => { resolve = done }) })
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    await confirm(wrapper)
    await cancel(wrapper)
    expect(context.saveAction).toHaveBeenCalledTimes(2)
    expect(named(wrapper, 'AModal').props('confirmLoading')).toBe(true)
    expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    resolve(true)
    await flushPromises()
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('does not submit a variable when canceled during form validation', async () => {
    let resolve
    mocks.validate.mockImplementationOnce(() => new Promise(done => { resolve = done }))
    const { wrapper } = render(VariableModal)
    wrapper.vm.show({ variable_key: 'member', variable_name: '会员' })
    await flushPromises()
    named(wrapper, 'AModal').vm.$emit('ok')
    await cancel(wrapper)
    resolve()
    await flushPromises()
    expect(mocks.createVariable).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('discards optimization background/model/effort drafts and reloads the latest saved values', async () => {
    const { wrapper, robotInfo } = render(ProblemOptimization)
    await openSetting(wrapper)
    await loadModels(wrapper)
    expect(named(wrapper, 'AModal').props('width')).toBe(640)
    expect(named(wrapper, 'ASlider').props('value')).toBe(3)
    await wrapper.find('textarea').setValue('取消的背景')
    await chooseRadio(wrapper, '1')
    named(wrapper, 'ASlider').vm.$emit('change', 4)
    await cancel(wrapper)
    expect(robotInfo.optimize_question_dialogue_background).toBe('原对话背景')
    expect(mocks.saveRobot).not.toHaveBeenCalled()
    robotInfo.optimize_question_dialogue_background = '最新背景'
    await openSetting(wrapper)
    expect(wrapper.find('textarea').element.value).toBe('最新背景')
    expect(named(wrapper, 'ARadioGroup').props('value')).toBe('0')
    expect(named(wrapper, 'ReasoningEffortSelect').props('value')).toBe('medium')
    wrapper.unmount()
  })

  it('synchronizes effort controls, validates the specified model and submits a private patch', async () => {
    const { wrapper, robotInfo, context } = render(ProblemOptimization)
    const original = JSON.stringify(robotInfo)
    await openSetting(wrapper)
    await chooseRadio(wrapper, '1')
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenCalledWith('ph_select_model')
    expect(mocks.saveRobot).not.toHaveBeenCalled()
    await loadModels(wrapper)
    named(wrapper, 'ASlider').vm.$emit('change', 4)
    await flushPromises()
    expect(named(wrapper, 'ReasoningEffortSelect').props('value')).toBe('high')
    const effort = named(wrapper, 'ReasoningEffortSelect')
    effort.vm.$emit('update:value', 'none')
    effort.vm.$emit('change', 'none')
    await flushPromises()
    expect(named(wrapper, 'ASlider').props('value')).toBe(0)
    await wrapper.find('textarea').setValue('保存的背景')
    await confirm(wrapper)
    expect(context.saveAction).toHaveBeenCalledTimes(1)
    expect(savedPatch()).toMatchObject({ optimize_question_model_config_id: 9, optimize_question_use_model: 'other-model', optimize_question_enable_thinking: 0, optimize_question_reasoning_effort: 'none', optimize_question_dialogue_background: '保存的背景' })
    expect(JSON.stringify(robotInfo)).toBe(original)
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('preserves custom reasoning levels and hides the thinking controls for unsupported models', async () => {
    const { wrapper, robotInfo } = render(ProblemOptimization)
    robotInfo.optimize_question_reasoning_effort = 'custom-effort'
    await openSetting(wrapper)
    named(wrapper, 'ModelSelect').vm.$emit('loaded', [], {})
    await flushPromises()
    expect(wrapper.find('.thinking-slider-row').exists()).toBe(false)
    await loadModels(wrapper)
    expect(named(wrapper, 'ReasoningEffortSelect').props('value')).toBe('custom-effort')
    await confirm(wrapper)
    expect(savedPatch()).toMatchObject({ optimize_question_model_config_id: 0, optimize_question_use_model: '', optimize_question_reasoning_effort: 'custom-effort' })
    wrapper.unmount()
  })

  it('keeps optimization edits after failure and prevents cancel or duplicate saves while pending', async () => {
    const { wrapper, context } = render(ProblemOptimization)
    await openSetting(wrapper)
    await wrapper.find('textarea').setValue('重试背景')
    context.saveAction.mockRejectedValueOnce(new Error('offline'))
    await confirm(wrapper)
    expect(wrapper.find('[role="alert"]').exists()).toBe(true)
    expect(wrapper.find('textarea').element.value).toBe('重试背景')
    let resolve
    context.saveAction.mockImplementationOnce(task => { task(); return new Promise(done => { resolve = done }) })
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    await confirm(wrapper)
    await cancel(wrapper)
    expect(context.saveAction).toHaveBeenCalledTimes(2)
    expect(named(wrapper, 'AModal').props('closable')).toBe(false)
    resolve(true)
    await flushPromises()
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('stages all generation modes, validates required fields and lists only published workflows', async () => {
    const { wrapper, robotInfo } = render(SuggestedIssues, true, { robotList: [
      { has_published: 1, robot_key: 'published', robot_name: '已发布流程' },
      { has_published: 0, robot_key: 'draft', robot_name: '草稿流程' }
    ] })
    await openSetting(wrapper)
    expect(named(wrapper, 'AModal').props('width')).toBe(640)
    expect(wrapper.find('textarea').exists()).toBe(false)
    await chooseRadio(wrapper, 1)
    await wrapper.find('textarea').setValue('  ')
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenCalledWith('msg_prompt_required')
    await wrapper.find('textarea').setValue('草稿提示词')
    await chooseRadio(wrapper, 2)
    expect(named(wrapper, 'ASelect').props('options')).toEqual([{ label: '已发布流程', value: 'published' }])
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenCalledWith('msg_workflow_required')
    named(wrapper, 'ASelect').vm.$emit('update:value', 'published')
    await chooseRadio(wrapper, 1)
    expect(wrapper.find('textarea').element.value).toBe('草稿提示词')
    await cancel(wrapper)
    expect(robotInfo.question_guide_prompt).toBe('原提示词')
    expect(mocks.saveRobot).not.toHaveBeenCalled()
    await openSetting(wrapper)
    expect(named(wrapper, 'ARadioGroup').props('value')).toBe(0)
    wrapper.unmount()
  })

  it.each([0, 1, 2])('saves generation mode %s through the existing configuration endpoint', async mode => {
    const { wrapper, context } = render(SuggestedIssues)
    await openSetting(wrapper)
    await chooseRadio(wrapper, mode)
    if (mode === 1) await wrapper.find('textarea').setValue('  自定义提示词  ')
    if (mode === 2) named(wrapper, 'ASelect').vm.$emit('update:value', 'published')
    await confirm(wrapper)
    expect(context.saveAction).toHaveBeenCalledTimes(1)
    expect(savedPatch()).toMatchObject({ question_guide_mode: mode, question_guide_prompt: mode === 1 ? '自定义提示词' : '', question_guide_workflow_key: mode === 2 ? 'published' : '' })
    wrapper.unmount()
  })

  it('keeps the generation dialog open after failure and retries the same draft', async () => {
    const { wrapper, context } = render(SuggestedIssues)
    await openSetting(wrapper)
    await chooseRadio(wrapper, 1)
    await wrapper.find('textarea').setValue('保留提示词')
    context.saveAction.mockResolvedValueOnce(false)
    await confirm(wrapper)
    expect(wrapper.find('[role="alert"]').exists()).toBe(true)
    expect(wrapper.find('textarea').element.value).toBe('保留提示词')
    await confirm(wrapper)
    expect(savedPatch().question_guide_prompt).toBe('保留提示词')
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('preserves legacy modal sizing, route keys and the inline save path', async () => {
    const variable = render(VariableModal, false).wrapper
    variable.vm.show({ variable_key: 'member', variable_name: '会员' })
    await flushPromises()
    expect(named(variable, 'AModal').props('width')).toBe(648)
    await confirm(variable)
    expect(mocks.createVariable.mock.calls[0][0].robot_key).toBe('legacy-query')
    variable.unmount()
    const { wrapper, updateRobotInfo } = render(SuggestedIssues, false)
    await wrapper.find('.guide-action-row button').trigger('click')
    expect(named(wrapper, 'AModal').props('width')).toBe(520)
    await confirm(wrapper)
    expect(updateRobotInfo).toHaveBeenCalledTimes(1)
    expect(mocks.saveRobot).not.toHaveBeenCalled()
    wrapper.unmount()
  })
})
