import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive } from 'vue'
import SensitiveWords from '@/views/robot/robot-config/basic-config/components/sensitive-words.vue'
import ProblemOptimization from '@/views/robot/robot-config/basic-config/components/problem-optimization.vue'
import SuggestedIssues from '@/views/robot/robot-config/basic-config/components/suggested-issues.vue'
import DisplayCitations from '@/views/robot/robot-config/basic-config/components/display-aitations.vue'
import ShowLike from '@/views/robot/robot-config/basic-config/components/show-like.vue'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key, params) => key === 'msg_guide_description' ? `Before ${params.num} after` : key }) }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/reasoning-effort-select/index.vue', () => ({ default: { template: '<div />' } }))

const toggle = { name: 'ASwitch', props: ['checked', 'checkedValue', 'unCheckedValue'], emits: ['update:checked', 'change'], template: '<input type="checkbox" />' }
const modal = { name: 'AModal', props: ['open'], template: '<section v-if="open" role="dialog"><slot /></section>' }
function render(component, compact) {
  const updateRobotInfo = vi.fn()
  const robotInfo = reactive({ sensitive_words_switch: 0, enable_question_optimize: 'false', enable_question_guide: 'false',
    feedback_switch: 0, answer_source_switch: 'false', question_guide_num: 3, question_guide_mode: 1,
    optimize_question_model_config_id: 0 })
  const wrapper = mount(component, { props: { compact }, global: {
    provide: { robotInfo: { robotInfo, updateRobotInfo } },
    stubs: { 'svg-icon': true, 'a-switch': toggle, 'a-modal': modal,
      'a-tooltip': { template: '<span><slot /></span>' }, 'a-button': { template: '<button><slot /></button>' },
      'a-input-number': { name: 'AInputNumber', props: ['value', 'size'], emits: ['update:value', 'blur'], template: '<input type="number" />' },
      'a-radio-group': true, 'a-radio': true, 'a-textarea': true }
  } })
  return { wrapper, robotInfo, updateRobotInfo }
}

describe('compact workspace settings', () => {
  it.each([
    [SensitiveWords, 'sensitive_words_switch', 1], [ProblemOptimization, 'enable_question_optimize', 'true'],
    [SuggestedIssues, 'enable_question_guide', 'true'], [ShowLike, 'feedback_switch', 1],
    [DisplayCitations, 'answer_source_switch', true]
  ])('keeps the header switch connected to the existing save path', async (component, field, value) => {
    const { wrapper, robotInfo, updateRobotInfo } = render(component, true)
    expect(wrapper.classes()).toContain('compact-setting')
    expect(wrapper.find('.edit-box-header input[type="checkbox"]').exists()).toBe(true)
    expect(wrapper.find('.edit-box-body input[type="checkbox"]').exists()).toBe(false)
    const control = wrapper.findComponent(toggle)
    control.vm.$emit('update:checked', value)
    control.vm.$emit('change', value)
    await wrapper.vm.$nextTick()
    expect(updateRobotInfo).toHaveBeenCalledOnce()
    if (component === DisplayCitations) expect(updateRobotInfo).toHaveBeenCalledWith({ answer_source_switch: 'true' })
    else expect(robotInfo[field]).toBe(value)
    wrapper.unmount()
  })

  it.each([ProblemOptimization, SuggestedIssues])('opens settings from the header before the switch', async component => {
    const { wrapper } = render(component, true)
    const header = wrapper.find('.edit-box-header')
    const button = header.find('button')
    expect(button.element.compareDocumentPosition(header.find('input').element) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    await button.trigger('click')
    expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    expect(wrapper.find('.edit-box-body button').exists()).toBe(false)
    wrapper.unmount()
  })

  it('opens sensitive-word management from the header and preserves the legacy inline link', async () => {
    const open = vi.spyOn(window, 'open').mockImplementation(() => null)
    const workspace = render(SensitiveWords, true).wrapper
    await workspace.find('.edit-box-header button').trigger('click')
    expect(open).toHaveBeenCalledWith('/#/user/sensitive-words')
    expect(workspace.find('.manage-link').exists()).toBe(false)
    const legacy = render(SensitiveWords).wrapper
    expect(legacy.find('.manage-link').exists()).toBe(true)
    expect(legacy.find('.edit-box-body input[type="checkbox"]').exists()).toBe(true)
    workspace.unmount()
    legacy.unmount()
    open.mockRestore()
  })

  it.each([ProblemOptimization, SuggestedIssues, DisplayCitations, ShowLike])('keeps legacy layout by default', component => {
    const { wrapper } = render(component)
    expect(wrapper.classes()).not.toContain('compact-setting')
    if (component === ProblemOptimization) expect(wrapper.find('.dialog-bg-set').exists()).toBe(true)
    if (component === SuggestedIssues) expect(wrapper.find('.guide-action-row button').exists()).toBe(true)
    if (component === DisplayCitations) expect(wrapper.find('.switch-box-right input').exists()).toBe(true)
    if (component === ShowLike) expect(wrapper.find('.edit-box-body .switch-item').exists()).toBe(true)
    wrapper.unmount()
  })

  it('preserves the question-count input and current generation-mode summary', async () => {
    const { wrapper, robotInfo, updateRobotInfo } = render(SuggestedIssues, true)
    expect(wrapper.find('.edit-box-header').text()).toContain('mode_prompt')
    const count = wrapper.findComponent({ name: 'AInputNumber' })
    count.vm.$emit('update:value', 5)
    count.vm.$emit('blur')
    await wrapper.vm.$nextTick()
    expect(robotInfo.question_guide_num).toBe(5)
    expect(updateRobotInfo).toHaveBeenCalledOnce()
    wrapper.unmount()
  })
})
