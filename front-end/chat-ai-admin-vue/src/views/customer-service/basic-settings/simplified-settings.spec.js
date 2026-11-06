import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount, shallowMount } from '@vue/test-utils'
import { reactive, ref } from 'vue'

const mocks = vi.hoisted(() => ({
  leaf: name => ({ name, template: '<div />' }), robotStore: null, save: vi.fn(), legacySave: vi.fn(),
  libraries: vi.fn(), meta: vi.fn(), variables: vi.fn(), error: vi.fn(), graphShow: vi.fn(), verify: vi.fn()
}))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => mocks.robotStore }))
vi.mock('@/stores/modules/model', () => ({ useModelStore: () => ({ getAllmodelList: vi.fn() }) }))
vi.mock('@/stores/modules/mini-card', () => ({ useMiniCardStore: () => ({ cardList: [], fetchCardList: vi.fn() }) }))
vi.mock('@/api/robot', () => ({ saveRobot: mocks.save, getRobotList: vi.fn().mockResolvedValue({ data: [] }), getChatVariables: mocks.variables }))
vi.mock('@/api/library', () => ({ getLibraryList: mocks.libraries, getRobotMetaSchemaList: mocks.meta }))
vi.mock('ant-design-vue', async importOriginal => ({ Modal: (await importOriginal()).Modal, message: { error: mocks.error } }))
vi.mock('@/components/model-select', () => ({ getModelNameText: () => 'Rerank model' }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: mocks.leaf('ModelSelect') }))
vi.mock('@/components/weight-select/index.vue', () => ({ default: { name: 'WeightSelect', props: ['rrf_weight'], emits: ['update:rrf_weight'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/meta-filter-box.vue', () => ({ default: { name: 'MetaFilterBox', props: ['rule', 'type'], setup: (_, { expose }) => { expose({ verify: mocks.verify }) }, template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/associated-knowledge-base/no-open-graph-modal.vue', () => ({ default: { name: 'GraphWarning', props: ['list'], setup: (_, { expose }) => { expose({ show: mocks.graphShow }) }, template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/at-input/index.vue', () => ({ default: { name: 'AtInput', props: ['defaultValue'], emits: ['change'], template: '<textarea :value="defaultValue" @input="$emit(\'change\', $event.target.value)" />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/ai-create-prompt.vue', () => ({ default: mocks.leaf('AiCreatePrompt') }))
vi.mock('@/views/customer-service/basic-settings/components/ai-prompt-dialog.vue', () => ({ default: mocks.leaf('AiCreatePrompt') }))
vi.mock('@/components/import-prompt/index.vue', () => ({ default: mocks.leaf('ImportPrompt') }))
vi.mock('@/components/import-prompt/upload-prompt.vue', () => ({ default: mocks.leaf('UploadPrompt') }))
vi.mock('@/views/robot/robot-config/basic-config/components/slash-popup-panel.vue', () => ({ default: mocks.leaf('SlashPopupPanel') }))
vi.mock('@/views/robot/robot-config/basic-config/components/variable-setting/add-variable-modal.vue', () => ({ default: mocks.leaf('AddVariableModal') }))
vi.mock('@/components/mini-card/mini-card-form-modal.vue', () => ({ default: mocks.leaf('MiniCardFormModal') }))
vi.mock('@/views/robot/robot-config/basic-config/components/basic-config.vue', () => ({ default: mocks.leaf('BasicInfo') }))
vi.mock('@/views/robot/robot-config/basic-config/components/associated-knowledge-base/index.vue', () => ({ default: mocks.leaf('AssociatedKnowledgeBase') }))
vi.mock('@/views/robot/robot-config/basic-config/components/model-settings.vue', () => ({ default: mocks.leaf('ModelSettings') }))
vi.mock('@/views/robot/robot-config/basic-config/components/chat-mode/index.vue', () => ({ default: mocks.leaf('ChatMode') }))
vi.mock('@/views/robot/robot-config/basic-config/components/chat-cache.vue', () => ({ default: mocks.leaf('ChatCache') }))
vi.mock('@/views/robot/robot-config/basic-config/components/skill/index.vue', () => ({ default: mocks.leaf('Skill') }))
vi.mock('@/views/robot/robot-config/basic-config/components/data-base/index.vue', () => ({ default: mocks.leaf('DataBase') }))
vi.mock('@/views/robot/robot-config/basic-config/components/language-setting/index.vue', () => ({ default: mocks.leaf('LanguageSetting') }))
vi.mock('@/views/robot/robot-config/basic-config/components/language-setting/common-problem.vue', () => ({ default: mocks.leaf('CommonProblem') }))
vi.mock('@/views/robot/robot-config/basic-config/components/variable-setting/index.vue', () => ({ default: mocks.leaf('VariableSetting') }))
vi.mock('@/views/robot/robot-config/basic-config/components/sensitive-words.vue', () => ({ default: mocks.leaf('SensitiveWords') }))
vi.mock('@/views/robot/robot-config/basic-config/components/problem-optimization.vue', () => ({ default: mocks.leaf('ProblemOptimization') }))
vi.mock('@/views/robot/robot-config/basic-config/components/suggested-issues.vue', () => ({ default: mocks.leaf('SuggestedIssues') }))
vi.mock('@/views/robot/robot-config/basic-config/components/display-aitations.vue', () => ({ default: mocks.leaf('DisplayAitations') }))
vi.mock('@/views/robot/robot-config/basic-config/components/show-like.vue', () => ({ default: mocks.leaf('ShowLike') }))
vi.mock('@/views/robot/robot-config/basic-config/role-permission.vue', () => ({ default: mocks.leaf('RolePermission') }))

import BasicConfig from '@/views/robot/robot-config/basic-config/index.vue'
import SystemPromptWords from '@/views/robot/robot-config/basic-config/components/system-prompt-words.vue'
import RecallDialog from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import EditBox from '@/views/robot/robot-config/basic-config/components/edit-box.vue'
import PromptCard from './components/prompt-card.vue'
import RecallCard from './components/recall-settings-card.vue'
import ChatModeCard from './components/chat-mode-card.vue'

const group = name => ({ name, template: '<div><slot /></div>' })
const tabs = { ...group('ATabs'), props: ['activeKey'], emits: ['update:activeKey', 'change'] }
const modal = { name: 'AModal', props: ['open', 'confirmLoading', 'closable', 'cancelButtonProps'], emits: ['ok', 'cancel', 'update:open'], template: '<section v-if="open" role="dialog"><slot /></section>' }
const stubs = {
  'a-button': { template: '<button><slot /></button>' }, 'a-tabs': tabs, 'a-tab-pane': group('ATabPane'),
  'a-modal': modal, 'a-flex': group('AFlex'), 'a-tooltip': group('ATooltip'), 'svg-icon': true,
  'a-input-number': { name: 'AInputNumber', props: ['value'], emits: ['update:value'], template: '<input />' },
  'a-switch': { name: 'ASwitch', props: ['checked'], emits: ['update:checked', 'change'], template: '<input />' },
  'a-slider': true, 'a-select': group('ASelect'), 'a-select-option': true, 'a-radio-group': true, 'a-radio': true
}
const promptStruct = role => Object.fromEntries(['role', 'task', 'constraints', 'skill', 'output', 'tone'].map(key => [key, { subject: key, describe: key === 'role' ? role : '' }]).concat([['custom', []]]))
function setupContext(applicationType = 0) {
  const robotInfo = reactive({ id: 7, robot_key: 'seven', application_type: applicationType, prompt: 'original custom', prompt_type: 1,
    prompt_struct: promptStruct('original role'), prompt_struct_default: promptStruct('default'), library_ids: '101,202', default_library_id: 101, normal_library_id: 202,
    work_flow_ids: '900', form_ids: '12', search_type: 1, top_k: 5, similarity: 0.4, rerank_status: 0, meta_search_switch: 0, meta_search_type: 1, rrf_weight: '{"vector":0.5,"search":0.3,"graph":0.2}' })
  mocks.robotStore = reactive({ robotInfo, chatVariables: [], getRobot: vi.fn(), fetchChatVariables: vi.fn() })
  const context = { robotInfo, savedRevision: ref(0), updateRobotInfo: vi.fn(), saveAction: vi.fn(async task => {
    try {
      await task()
      const payload = mocks.save.mock.calls.at(-1)[0]
      Object.assign(robotInfo, payload, { prompt_struct: JSON.parse(payload.prompt_struct) })
      context.savedRevision.value++
      return true
    } catch { return false }
  }) }
  return context
}
const options = context => ({ global: { provide: { customerRobotContext: context, robotInfo: { robotInfo: context.robotInfo, updateRobotInfo: mocks.legacySave } }, stubs } })
const named = (wrapper, name) => wrapper.findComponent({ name })
const openPrompt = async wrapper => { await wrapper.find('.edit-box-header button').trigger('click'); await flushPromises() }
const sendAI = async (wrapper, text) => { named(wrapper, 'AiCreatePrompt').vm.$emit('handleAiSave', promptStruct(text)); await flushPromises() }
const sendImport = async (wrapper, text) => { named(wrapper, 'ImportPrompt').vm.$emit('ok', { prompt_type: 0, prompt: text }); await flushPromises() }

describe('workspace basic settings compatibility', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    mocks.save.mockResolvedValue({ res: 0 })
    mocks.legacySave.mockResolvedValue(true)
    mocks.libraries.mockResolvedValue({ data: [{ id: 101, graph_switch: 0 }] })
    mocks.meta.mockResolvedValue({ data: [] })
    mocks.variables.mockResolvedValue({ data: [] })
    mocks.verify.mockReturnValue(true)
  })

  it('ignores the legacy common-question tab cache in workspace mode without changing it', async () => {
    const context = setupContext()
    const key = '/robot/config/basic-config/activeKey'
    localStorage.setItem(key, '2')
    const reads = vi.spyOn(Storage.prototype, 'getItem')
    const wrapper = shallowMount(BasicConfig, { ...options(context), props: { variant: 'workspace' } })
    await flushPromises()
    expect(wrapper.find('basic-config-stub').exists()).toBe(true)
    expect(named(wrapper, 'ATabs').exists()).toBe(false)
    expect(wrapper.findComponent(RecallCard).exists()).toBe(true)
    expect(wrapper.findComponent(PromptCard).exists()).toBe(true)
    expect(wrapper.findComponent(ChatModeCard).exists()).toBe(true)
    expect(wrapper.find('.scroll-box').element.children.length).toBe(13)
    const sections = wrapper.findAll('.scroll-box > .setting-box')
    const expectedOrder = ['basic-config-stub', 'chat-mode-card-stub', 'model-settings-stub', 'prompt-card-stub', 'recall-settings-card-stub', 'language-setting-stub', 'variable-setting-stub', 'chat-cache-stub']
    expectedOrder.forEach((selector, index) => expect(sections[index].find(selector).exists()).toBe(true))
    for (const name of ['SensitiveWords', 'ProblemOptimization', 'SuggestedIssues', 'DisplayAitations', 'ShowLike']) expect(named(wrapper, name).attributes('compact')).toBe('true')
    expect(named(wrapper, 'ChatMode').exists()).toBe(false)
    for (const name of ['CommonProblem', 'Skill', 'DataBase', 'AssociatedKnowledgeBase']) expect(named(wrapper, name).exists()).toBe(false)
    expect(reads).not.toHaveBeenCalledWith(key)
    expect(localStorage.getItem(key)).toBe('2')
    reads.mockRestore()
    wrapper.unmount()
  })

  it('keeps the old tabs, associated libraries, skills, database and inline prompt editor by default', async () => {
    const context = setupContext()
    localStorage.setItem('/robot/config/basic-config/activeKey', '2')
    const wrapper = shallowMount(BasicConfig, options(context))
    expect(named(wrapper, 'CommonProblem').exists()).toBe(true)
    named(wrapper, 'ATabs').vm.$emit('update:activeKey', 1)
    named(wrapper, 'ATabs').vm.$emit('change', 1)
    await flushPromises()
    expect(wrapper.find('basic-config-stub').exists()).toBe(true)
    for (const name of ['AssociatedKnowledgeBase', 'Skill', 'DataBase']) expect(named(wrapper, name).exists()).toBe(true)
    expect(wrapper.findComponent(SystemPromptWords).exists()).toBe(true)
    expect(wrapper.findComponent(PromptCard).exists()).toBe(false)
    expect(named(wrapper, 'ChatMode').exists()).toBe(true)
    const sections = wrapper.findAll('.scroll-box > .setting-box')
    const expectedOrder = ['basic-config-stub', 'associated-knowledge-base-stub', 'model-settings-stub', 'chat-mode-stub', 'chat-cache-stub', 'system-prompt-words-stub', 'skill-stub', 'data-base-stub', 'language-setting-stub', 'variable-setting-stub']
    expectedOrder.forEach((selector, index) => expect(sections[index].find(selector).exists()).toBe(true))
    expect(wrapper.findComponent(ChatModeCard).exists()).toBe(false)
    expect(localStorage.getItem('/robot/config/basic-config/activeKey')).toBe('1')
    wrapper.unmount()
  })

  it('shows the workspace scrollbar only while scrolling and clears the idle timer on unmount', async () => {
    vi.useFakeTimers()
    let wrapper
    let legacy
    try {
      const context = setupContext()
      wrapper = shallowMount(BasicConfig, { ...options(context), props: { variant: 'workspace' } })
      const scroller = wrapper.find('.scroll-box')
      expect(scroller.classes()).not.toContain('is-scrolling')
      await scroller.trigger('scroll')
      expect(scroller.classes()).toContain('is-scrolling')
      await vi.advanceTimersByTimeAsync(600)
      await scroller.trigger('scroll')
      await vi.advanceTimersByTimeAsync(600)
      expect(scroller.classes()).toContain('is-scrolling')
      await vi.advanceTimersByTimeAsync(200)
      expect(scroller.classes()).not.toContain('is-scrolling')
      await scroller.trigger('scroll')
      wrapper.unmount()
      wrapper = null
      expect(vi.getTimerCount()).toBe(0)
      legacy = shallowMount(BasicConfig, options(context))
      await legacy.find('.scroll-box').trigger('scroll')
      expect(legacy.find('.scroll-box').classes()).not.toContain('is-scrolling')
      expect(vi.getTimerCount()).toBe(0)
    } finally {
      wrapper?.unmount()
      legacy?.unmount()
      vi.useRealTimers()
    }
  })

  it('discards manual and AI edits on cancel without changing the external preview', async () => {
    const context = setupContext()
    const wrapper = mount(PromptCard, options(context))
    await openPrompt(wrapper)
    named(wrapper, 'AtInput').vm.$emit('change', 'unsaved manual edit')
    await sendAI(wrapper, 'unsaved AI edit')
    expect(wrapper.find('.prompt-preview').text()).toContain('original role')
    expect(context.robotInfo.prompt_struct.role.describe).toBe('original role')
    expect(mocks.save).not.toHaveBeenCalled()
    named(wrapper, 'AModal').vm.$emit('cancel')
    await flushPromises()
    await openPrompt(wrapper)
    expect(named(wrapper, 'AtInput').props('defaultValue')).toBe('original role')
    wrapper.unmount()
  })

  it('uses the same card and small edit button as existing settings', () => {
    const context = setupContext()
    for (const component of [PromptCard, RecallCard]) {
      const wrapper = mount(component, options(context))
      expect(wrapper.findComponent(EditBox).exists()).toBe(true)
      expect(wrapper.classes()).toContain('edit-box')
      expect(wrapper.find('.edit-box-header button').attributes('size')).toBe('small')
      wrapper.unmount()
    }
  })

  it('expands and collapses a long read-only prompt', async () => {
    const height = vi.spyOn(HTMLElement.prototype, 'scrollHeight', 'get').mockReturnValue(240)
    const context = setupContext()
    context.robotInfo.prompt_type = 0
    context.robotInfo.prompt = 'A long prompt\n'.repeat(30)
    const wrapper = mount(PromptCard, options(context))
    await flushPromises()
    expect(wrapper.find('.prompt-preview').classes()).not.toContain('expanded')
    await wrapper.find('.expand button').trigger('click')
    expect(wrapper.find('.prompt-preview').classes()).toContain('expanded')
    await wrapper.find('.expand button').trigger('click')
    expect(wrapper.find('.prompt-preview').classes()).not.toContain('expanded')
    wrapper.unmount()
    height.mockRestore()
  })

  it('stages an import and saves only on confirmation, preserving application type and existing associations', async () => {
    const context = setupContext(2)
    const wrapper = mount(PromptCard, options(context))
    await openPrompt(wrapper)
    await sendImport(wrapper, 'Imported\n<img src="invalid">')
    expect(mocks.save).not.toHaveBeenCalled()
    expect(wrapper.find('.prompt-preview').text()).toContain('original role')
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(mocks.save).toHaveBeenCalledWith(expect.objectContaining({ id: 7, application_type: 2, library_ids: '101,202', work_flow_ids: '900', form_ids: '12', prompt_type: 0 }), 2)
    expect(wrapper.find('.prompt-preview').text()).toContain('Imported\n<img src="invalid">')
    expect(wrapper.find('.prompt-preview img').exists()).toBe(false)
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    wrapper.unmount()
  })

  it('retains a failed prompt draft for retry and updates the preview only after successful saving', async () => {
    const context = setupContext()
    mocks.save.mockRejectedValueOnce(new Error('network'))
    const wrapper = mount(PromptCard, options(context))
    await openPrompt(wrapper)
    await sendAI(wrapper, 'retry this prompt')
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(named(wrapper, 'AModal').exists()).toBe(true)
    expect(context.robotInfo.prompt_struct.role.describe).toBe('original role')
    expect(named(wrapper, 'AtInput').props('defaultValue')).toBe('retry this prompt')
    expect(wrapper.find('.prompt-preview').text()).not.toContain('retry this prompt')
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(wrapper.find('.prompt-preview').text()).toContain('retry this prompt')
    wrapper.unmount()
  })

  it('keeps immediate saving for AI generation and import in the legacy prompt editor', async () => {
    const context = setupContext()
    const wrapper = mount(SystemPromptWords, options(context))
    await sendAI(wrapper, 'legacy AI')
    expect(mocks.legacySave).toHaveBeenCalledWith(expect.objectContaining({ prompt_type: 1, prompt_struct: promptStruct('legacy AI') }))
    await sendImport(wrapper, 'legacy import')
    expect(mocks.legacySave).toHaveBeenLastCalledWith(expect.objectContaining({ prompt_type: 0, prompt: 'legacy import' }))
    wrapper.unmount()
  })

  it('keeps recall drafts private, validates Rerank and retries failures without changing library associations', async () => {
    const context = setupContext()
    mocks.save.mockRejectedValueOnce(new Error('network'))
    const wrapper = mount(RecallCard, options(context))
    await wrapper.find('.edit-box-header button').trigger('click')
    await flushPromises()
    const dialog = wrapper.findComponent(RecallDialog)
    named(dialog, 'AInputNumber').vm.$emit('update:value', 12)
    named(dialog, 'ASwitch').vm.$emit('update:checked', 1)
    named(dialog, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(mocks.error).toHaveBeenCalledWith('msg_please_select_rerank_model')
    expect(mocks.save).not.toHaveBeenCalled()
    named(dialog, 'ASwitch').vm.$emit('update:checked', 0)
    named(dialog, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(named(dialog, 'AModal').props('open')).toBe(true)
    expect(context.robotInfo.top_k).toBe(5)
    named(dialog, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(context.robotInfo.top_k).toBe(12)
    expect(mocks.save).toHaveBeenLastCalledWith(expect.objectContaining({ library_ids: '101,202', default_library_id: 101, normal_library_id: 202 }), 0)
    expect(mocks.graphShow).toHaveBeenCalledOnce()
    expect(named(dialog, 'AModal').props('open')).toBe(false)
    expect(dialog.emitted('change')).toBeUndefined()
    wrapper.unmount()
  })

  it('preserves legacy recall close/change behavior when no async submit is supplied', async () => {
    const context = setupContext()
    const wrapper = mount(RecallDialog, options(context))
    wrapper.vm.open({ ...context.robotInfo, rrf_weight: { vector: 0.5, search: 0.3, graph: 0.2 } }, context.robotInfo)
    await flushPromises()
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(wrapper.emitted('change')).toHaveLength(1)
    expect(named(wrapper, 'AModal').props('open')).toBe(false)
    expect(mocks.save).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('discards recall edits and reopens with saved values after cancellation', async () => {
    const context = setupContext()
    context.robotInfo.rrf_weight = ''
    const wrapper = mount(RecallCard, options(context))
    await wrapper.find('.edit-box-header button').trigger('click')
    await flushPromises()
    const dialog = wrapper.findComponent(RecallDialog)
    expect(named(dialog, 'WeightSelect').props('rrf_weight')).toEqual({ vector: 0, search: 0, graph: 0 })
    named(dialog, 'AInputNumber').vm.$emit('update:value', 29)
    named(dialog, 'AModal').vm.$emit('update:open', false)
    await flushPromises()
    expect(context.robotInfo.top_k).toBe(5)
    await wrapper.find('.edit-box-header button').trigger('click')
    await flushPromises()
    expect(named(dialog, 'AInputNumber').props('value')).toBe(5)
    expect(mocks.save).not.toHaveBeenCalled()
    wrapper.unmount()
  })
})
