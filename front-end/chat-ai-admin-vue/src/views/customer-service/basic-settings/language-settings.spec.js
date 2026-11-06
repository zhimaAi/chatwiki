import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive } from 'vue'
import LanguageSetting from '@/views/robot/robot-config/basic-config/components/language-setting/index.vue'

const mocks = vi.hoisted(() => ({ store: null, save: vi.fn(), leaf: name => ({ name, props: ['currentMultiLangConfig', 'isEdit'], template: '<div />' }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => mocks.store }))
vi.mock('@/api/robot/index', () => ({ saveRobotLangConfigs: mocks.save }))
vi.mock('@/views/robot/robot-config/basic-config/components/language-setting/welcome-words.vue', () => ({ default: mocks.leaf('WelcomeWords') }))
vi.mock('@/views/robot/robot-config/basic-config/components/language-setting/unknown-problem-prompt.vue', () => ({ default: mocks.leaf('UnknownProblemPrompt') }))
vi.mock('@/views/robot/robot-config/basic-config/components/language-setting/answer-prompt.vue', () => ({ default: mocks.leaf('AnswerPrompt') }))

const legacyKey = 'robot_config_more_language_visible_keys'
const workspaceKey = 'customer_service_more_language_visible_keys'
const stubs = {
  'svg-icon': true,
  'a-popover': { template: '<div><slot /><slot name="content" /></div>' },
  'a-button': { template: '<button><slot /></button>' },
  'a-modal': { name: 'AModal', props: ['open', 'width', 'confirmLoading', 'closable'], emits: ['ok', 'cancel'], template: '<section v-if="open" role="dialog"><header><slot name="title" /></header><slot /></section>' },
  'a-select': { name: 'ASelect', props: ['value', 'options'], emits: ['update:value'], template: '<select />' },
  'a-input': { name: 'AInput', props: ['value'], emits: ['update:value'], template: '<input :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-textarea': { name: 'ATextarea', props: ['value'], emits: ['update:value'], template: '<textarea :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-alert': { name: 'AAlert', props: ['message'], template: '<div role="alert">{{ message }}</div>' },
  draggable: { name: 'Draggable', props: ['modelValue'], emits: ['update:modelValue'], template: '<div><template v-for="(element, index) in modelValue" :key="element.id"><slot name="item" :element="element" :index="index" /></template></div>' },
  'a-tabs': { name: 'ATabs', props: ['activeKey'], emits: ['update:activeKey'], template: '<div><slot /></div>' },
  'a-tab-pane': { name: 'ATabPane', props: ['tab'], template: '<span>{{ tab }}</span>' },
  'a-switch': { name: 'ASwitch', props: ['checked', 'size'], emits: ['change', 'update:checked'], template: '<input type="checkbox" />' }
}
const render = (variant, context) => mount(LanguageSetting, { props: { variant }, global: { stubs, provide: context ? { customerRobotContext: context } : {} } })
const named = (wrapper, name) => wrapper.findComponent({ name })
const englishSwitch = wrapper => wrapper.findAll('.more-language-item').find(item => item.find('.language-name').text() === '英文').findComponent({ name: 'ASwitch' })
const open = wrapper => wrapper.find('.actions-box button').trigger('click')
const confirm = async wrapper => { named(wrapper, 'AModal').vm.$emit('ok'); await flushPromises() }
const cancel = async wrapper => { named(wrapper, 'AModal').vm.$emit('cancel'); await flushPromises() }
const changeLanguage = async (wrapper, key) => { named(wrapper, 'ASelect').vm.$emit('update:value', key); await flushPromises() }

describe('workspace system-prompt language visibility', () => {
  beforeEach(() => {
    localStorage.clear()
    mocks.save.mockReset().mockResolvedValue({ res: 0 })
    mocks.store = reactive({ getRobot: vi.fn().mockResolvedValue(), robotInfo: { id: 7, application_type: 0,
      multi_lang_configs: ['zh-CN', 'en-US', 'ja'].map(lang_key => ({ lang_key, content: lang_key,
        welcomes: JSON.stringify({ content: `欢迎${lang_key}`, question: ['问题一', '问题二'] }),
        unknown_question_prompt: JSON.stringify({ content: `未知${lang_key}`, question: ['帮助一', '帮助二'] }),
        tips_before_answer_switch: 'true', tips_before_answer_content: '思考中', common_question_list: '["保留"]' })) } })
  })

  it('defaults to Chinese without tabs and offers English as a disabled option', async () => {
    localStorage.setItem(legacyKey, JSON.stringify(['ja']))
    const wrapper = render('workspace')
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    await open(wrapper)
    expect(wrapper.find('.tabs-box').exists()).toBe(false)
    expect(englishSwitch(wrapper).props('checked')).toBe(false)
    expect(named(wrapper, 'WelcomeWords').props('currentMultiLangConfig').lang_key).toBe('zh-CN')
    expect(localStorage.getItem(legacyKey)).toBe('["ja"]')
    expect(mocks.save).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('shows enabled language tabs, remembers the choice and falls back to Chinese when English is hidden', async () => {
    let wrapper = render('workspace')
    await open(wrapper)
    englishSwitch(wrapper).vm.$emit('change', true)
    await flushPromises()
    expect(wrapper.findAllComponents({ name: 'ATabPane' }).map(item => item.props('tab'))).toEqual(['中文', '英文'])
    expect(localStorage.getItem(workspaceKey)).toBe('["en-US"]')
    wrapper.unmount()
    wrapper = render('workspace')
    await open(wrapper)
    expect(englishSwitch(wrapper).props('checked')).toBe(true)
    named(wrapper, 'ATabs').vm.$emit('update:activeKey', 'en-US')
    await flushPromises()
    expect(named(wrapper, 'WelcomeWords').props('currentMultiLangConfig').lang_key).toBe('en-US')
    englishSwitch(wrapper).vm.$emit('change', false)
    await flushPromises()
    expect(wrapper.find('.tabs-box').exists()).toBe(false)
    expect(named(wrapper, 'WelcomeWords').props('currentMultiLangConfig').lang_key).toBe('zh-CN')
    expect(mocks.store.robotInfo.multi_lang_configs).toHaveLength(3)
    expect(mocks.save).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('preserves legacy Chinese/English tabs and isolates its saved visibility preferences', async () => {
    localStorage.setItem(legacyKey, JSON.stringify(['ja']))
    const workspace = render('workspace')
    await open(workspace)
    englishSwitch(workspace).vm.$emit('change', true)
    await flushPromises()
    workspace.unmount()
    const legacy = render()
    expect(legacy.findAllComponents({ name: 'ATabPane' }).map(item => item.props('tab'))).toEqual(['中文', '英文', '日语'])
    expect(legacy.findAll('.language-name').map(item => item.text())).not.toContain('英文')
    expect(localStorage.getItem(legacyKey)).toBe('["ja"]')
    legacy.unmount()
  })

  it('opens an 800px dialog and discards edits in every language on cancel', async () => {
    localStorage.setItem(workspaceKey, '["en-US"]')
    const wrapper = render('workspace')
    const original = JSON.stringify(mocks.store.robotInfo.multi_lang_configs)
    await open(wrapper)
    expect(named(wrapper, 'AModal').props('width')).toBe(800)
    expect(named(wrapper, 'WelcomeWords').props('isEdit')).toBe(false)
    await wrapper.find('[data-section="welcomes"] textarea').setValue('新的中文')
    await changeLanguage(wrapper, 'en-US')
    await wrapper.find('[data-section="welcomes"] textarea').setValue('New welcome')
    await changeLanguage(wrapper, 'zh-CN')
    expect(wrapper.find('textarea').element.value).toBe('新的中文')
    expect(JSON.stringify(mocks.store.robotInfo.multi_lang_configs)).toBe(original)
    await cancel(wrapper)
    expect(mocks.save).not.toHaveBeenCalled()
    await open(wrapper)
    expect(wrapper.find('textarea').element.value).toBe('欢迎zh-CN')
    await changeLanguage(wrapper, 'en-US')
    expect(wrapper.find('textarea').element.value).toBe('欢迎en-US')
    wrapper.unmount()
  })

  it('saves edited languages, guide order/add/delete and the answer switch through the existing endpoint', async () => {
    localStorage.setItem(workspaceKey, '["en-US"]')
    const context = { saveAction: vi.fn(async task => { await task(); return true }) }
    const wrapper = render('workspace', context)
    const original = JSON.stringify(mocks.store.robotInfo.multi_lang_configs)
    await open(wrapper)
    let welcome = wrapper.find('[data-section="welcomes"]')
    await welcome.find('textarea').setValue('新欢迎语')
    await welcome.find('.question-drag').trigger('keydown', { key: 'ArrowDown' })
    expect(welcome.findAll('.guide-question input').map(item => item.element.value)).toEqual(['问题二', '问题一'])
    await welcome.find('.guide-question .icon-button:last-child').trigger('click')
    await welcome.find('.add-question').trigger('click')
    await welcome.findAll('.guide-question input')[1].setValue('新增问题')
    const unknown = wrapper.find('[data-section="unknown_question_prompt"]')
    const drag = named(unknown, 'Draggable')
    drag.vm.$emit('update:modelValue', [...drag.props('modelValue')].reverse())
    await flushPromises()
    named(wrapper.find('.answer-section'), 'ASwitch').vm.$emit('update:checked', 'false')
    await changeLanguage(wrapper, 'en-US')
    await wrapper.find('textarea').setValue('Hello')
    expect(mocks.save).not.toHaveBeenCalled()
    await confirm(wrapper)
    expect(context.saveAction).toHaveBeenCalledTimes(1)
    const payload = mocks.save.mock.calls[0][0]
    expect(payload.id).toBe(7)
    const configs = JSON.parse(payload.multi_lang_configs)
    expect(JSON.parse(configs[0].welcomes)).toEqual({ content: '新欢迎语', question: ['问题一', '新增问题'] })
    expect(JSON.parse(configs[0].unknown_question_prompt).question).toEqual(['帮助二', '帮助一'])
    expect(configs[0].tips_before_answer_switch).toBe('false')
    expect(JSON.parse(configs[1].welcomes).content).toBe('Hello')
    expect(configs[2]).toEqual(mocks.store.robotInfo.multi_lang_configs[2])
    expect(configs[0].common_question_list).toBe('["保留"]')
    expect(JSON.stringify(mocks.store.robotInfo.multi_lang_configs)).toBe(original)
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('creates a missing language without copying Chinese and keeps its draft when temporarily hidden', async () => {
    const wrapper = render('workspace')
    await open(wrapper)
    const vietnamese = wrapper.findAll('.more-language-item').find(item => item.find('.language-name').text() === '越南语')
    const toggle = named(vietnamese, 'ASwitch')
    toggle.vm.$emit('change', true)
    await flushPromises()
    await changeLanguage(wrapper, 'vi')
    expect(wrapper.find('textarea').element.value).toBe('')
    await wrapper.find('textarea').setValue('Xin chào')
    toggle.vm.$emit('change', false)
    await flushPromises()
    expect(named(wrapper, 'ASelect').props('value')).toBe('zh-CN')
    toggle.vm.$emit('change', true)
    await flushPromises()
    await changeLanguage(wrapper, 'vi')
    expect(wrapper.find('textarea').element.value).toBe('Xin chào')
    await confirm(wrapper)
    const config = JSON.parse(mocks.save.mock.calls[0][0].multi_lang_configs).find(item => item.lang_key === 'vi')
    expect(JSON.parse(config.welcomes)).toEqual({ content: 'Xin chào', question: [] })
    wrapper.unmount()
  })

  it('retains the draft after failure, retries and blocks duplicate saves or cancel while pending', async () => {
    let resolve
    const pending = new Promise(done => { resolve = done })
    const context = { saveAction: vi.fn().mockResolvedValueOnce(false).mockImplementationOnce(task => { task(); return pending }) }
    const wrapper = render('workspace', context)
    await open(wrapper)
    await wrapper.find('textarea').setValue('失败后保留')
    await confirm(wrapper)
    expect(wrapper.find('[role="alert"]').exists()).toBe(true)
    expect(wrapper.find('textarea').element.value).toBe('失败后保留')
    named(wrapper, 'AModal').vm.$emit('ok')
    await flushPromises()
    expect(named(wrapper, 'AModal').props('confirmLoading')).toBe(true)
    expect(named(wrapper, 'AModal').props('closable')).toBe(false)
    await confirm(wrapper)
    await cancel(wrapper)
    expect(context.saveAction).toHaveBeenCalledTimes(2)
    expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    resolve(true)
    await flushPromises()
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('keeps workflow unknown prompts unchanged and preserves the legacy inline editor', async () => {
    mocks.store.robotInfo.application_type = 1
    const originalUnknown = mocks.store.robotInfo.multi_lang_configs[0].unknown_question_prompt
    const wrapper = render('workspace')
    await open(wrapper)
    expect(wrapper.find('[data-section="unknown_question_prompt"]').exists()).toBe(false)
    await confirm(wrapper)
    expect(JSON.parse(mocks.save.mock.calls[0][0].multi_lang_configs)[0].unknown_question_prompt).toBe(originalUnknown)
    wrapper.unmount()
    const legacy = render()
    await open(legacy)
    expect(named(legacy, 'WelcomeWords').props('isEdit')).toBe(true)
    expect(legacy.find('[role="dialog"]').exists()).toBe(false)
    legacy.unmount()
  })
})
