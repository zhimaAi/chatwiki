import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive } from 'vue'
import { Modal } from 'ant-design-vue'
import VariableModal from '@/views/robot/robot-config/basic-config/components/variable-setting/add-variable-modal.vue'
import VariableSetting from '@/views/robot/robot-config/basic-config/components/variable-setting/index.vue'
import SystemPromptWords from '@/views/robot/robot-config/basic-config/components/system-prompt-words.vue'
import RecallDialog from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import SuggestedIssues from '@/views/robot/robot-config/basic-config/components/suggested-issues.vue'

const mocks = vi.hoisted(() => ({ robot: null, createVariable: vi.fn(), validate: vi.fn(), meta: vi.fn(), variables: vi.fn() }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: { robot_key: 'legacy-seven' } }) }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => ({ companyInfo: { neo4j_status: 'true' } }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => mocks.robot }))
vi.mock('@/stores/modules/mini-card', () => ({ useMiniCardStore: () => ({ cardList: [], fetchCardList: vi.fn() }) }))
vi.mock('@/api/library/index.js', () => ({ getRobotMetaSchemaList: mocks.meta }))
vi.mock('@/api/robot/index.js', () => ({ getChatVariables: mocks.variables, createChatVariable: mocks.createVariable, deleteChatVariable: vi.fn() }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/meta-filter-box.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/ai-create-prompt.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/customer-service/basic-settings/components/ai-prompt-dialog.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/at-input/index.vue', () => ({ default: { name: 'AtInput', props: ['defaultValue'], template: '<textarea :value="defaultValue" />' } }))
vi.mock('@/components/import-prompt/index.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/import-prompt/upload-prompt.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/mini-card/mini-card-form-modal.vue', () => ({ default: { template: '<div />' } }))
const container = { template: '<div><slot /></div>' }
const wrappers = []
const promptStruct = () => Object.fromEntries(['role', 'task', 'constraints', 'skill', 'output', 'tone'].map(key => [key, { subject: key, describe: '' }]).concat([['custom', []]]))
function render(component) {
  const robotInfo = reactive({ id: 7, robot_key: 'legacy-seven', prompt_type: 0, prompt: '原提示词', prompt_struct: promptStruct(), prompt_struct_default: promptStruct(), question_guide_mode: 0, question_guide_num: 3 })
  mocks.robot = reactive({ chatVariables: [{ id: 1, variable_key: 'member', variable_name: '会员', variable_type: 'input_string' }], fetchChatVariables: vi.fn() })
  const wrapper = mount(component, { attachTo: document.body, global: { provide: { robotInfo: { robotInfo, updateRobotInfo: vi.fn() } }, stubs: {
    'svg-icon': true, 'a-flex': container, 'a-tooltip': container, 'a-tabs': container, 'a-tab-pane': true,
    'a-button': { template: '<button><slot /></button>' },
    'a-form': { setup: (_, { expose }) => expose({ validate: mocks.validate }), template: '<form><slot /></form>' },
    'a-form-item': { props: ['label'], template: '<div><label>{{ label }}</label><slot /></div>' }, 'a-form-item-rest': container,
    'a-input': { template: '<input />' }, 'a-input-number': { props: ['size', 'value'], template: '<input />' }, 'a-textarea': { template: '<textarea />' },
    'a-select': container, 'a-select-option': true, 'a-checkbox': true, 'a-switch': true, 'a-slider': true, 'a-radio-group': container, 'a-radio': true,
    'a-dropdown': container, 'a-menu': container, 'a-menu-item': container
  } } })
  wrappers.push(wrapper)
  return wrapper
}
async function expectDialog(wrapper, text) {
  await flushPromises()
  await vi.waitFor(() => expect(document.querySelector('[role="dialog"]')?.textContent).toContain(text))
  expect(wrapper.find('a-modal').exists()).toBe(false)
  const dialog = document.querySelector('[role="dialog"]')
  expect(wrapper.element.contains(dialog)).toBe(false)
  dialog.querySelector('.ant-modal-close').click()
  await flushPromises()
  expect(wrapper.findComponent(Modal).props('open')).toBe(false)
}
beforeEach(() => {
  vi.clearAllMocks()
  mocks.validate.mockResolvedValue({})
  mocks.meta.mockResolvedValue({ data: [] })
  mocks.variables.mockResolvedValue({ data: [] })
  const computedStyle = window.getComputedStyle
  vi.spyOn(window, 'getComputedStyle').mockImplementation(element => computedStyle(element))
  vi.stubGlobal('matchMedia', () => ({ matches: false, addListener() {}, removeListener() {} }))
})
afterEach(() => {
  wrappers.splice(0).forEach(wrapper => wrapper.unmount())
  document.body.innerHTML = ''
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

describe('legacy basic settings dialogs without global modal registration', () => {
  it('keeps the variable form hidden until opening the actual legacy modal', async () => {
    const wrapper = render(VariableModal)
    expect(wrapper.find('form').exists()).toBe(false)
    wrapper.vm.show()
    await expectDialog(wrapper, 'label_field_key')
    expect(wrapper.find('form').exists()).toBe(false)
    expect(mocks.createVariable).not.toHaveBeenCalled()
  })

  it('keeps variable forms out of the prompt editor and opens them only from the add action', async () => {
    const wrapper = render(SystemPromptWords)
    expect(wrapper.find('form').exists()).toBe(false)
    expect(wrapper.text()).not.toContain('label_field_key')
    expect(wrapper.find('textarea').element.value).toBe('原提示词')
    wrapper.findComponent(VariableModal).vm.show()
    await expectDialog(wrapper, 'label_field_key')
    expect(wrapper.find('textarea').element.value).toBe('原提示词')
  })

  it('keeps the legacy variable list intact and opens add-variable as a modal', async () => {
    const wrapper = render(VariableSetting)
    expect(wrapper.find('form').exists()).toBe(false)
    expect(wrapper.findAll('.variable-item')).toHaveLength(1)
    await wrapper.findAll('button').find(button => button.text() === 'btn_add_variable').trigger('click')
    await expectDialog(wrapper, 'label_field_key')
    expect(wrapper.findAll('.variable-item')).toHaveLength(1)
  })

  it('opens legacy recall settings in a modal instead of displaying the form inline', async () => {
    const wrapper = render(RecallDialog)
    expect(wrapper.find('.recall-settings-box').exists()).toBe(false)
    wrapper.vm.open({ search_type: 3, top_k: 3, similarity: 0.4, rrf_weight: { vector: 36, search: 28, graph: 36 }, rerank_status: 0, meta_search_switch: 0, meta_search_type: 1 }, { id: 7, robot_key: 'legacy-seven' })
    await expectDialog(wrapper, 'label_retrieval_mode')
    expect(wrapper.find('.recall-settings-box').exists()).toBe(false)
  })

  it('keeps legacy suggested-question generation settings hidden until its modal is opened', async () => {
    const wrapper = render(SuggestedIssues)
    expect(wrapper.find('.modal-item-box').exists()).toBe(false)
    await wrapper.find('.guide-action-row button').trigger('click')
    await flushPromises()
    await vi.waitFor(() => expect(document.querySelector('[role="dialog"]')?.textContent).toContain('label_generation_mode'))
    expect(wrapper.find('a-modal').exists()).toBe(false)
    expect(wrapper.findComponent(Modal).props('width')).toBe(520)
    document.querySelector('.ant-modal-close').click()
    await flushPromises()
    expect(wrapper.findComponent(Modal).props('open')).toBe(false)
  })
})
