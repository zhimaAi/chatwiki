import { afterAll, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive } from 'vue'
import AtInput from '@/views/robot/robot-config/basic-config/components/at-input/index.vue'
import PromptVariableInput from './components/prompt-variable-input.vue'
import PromptEditor from './components/prompt-editor.vue'
import PromptTokenText from './components/prompt-token-text.vue'
import { getPromptOptions, getPromptSelection, readPromptText } from './prompt-tokens'

const mocks = vi.hoisted(() => ({ robot: null, cards: null }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => mocks.robot }))
vi.mock('@/stores/modules/mini-card', () => ({ useMiniCardStore: () => mocks.cards }))
vi.mock('@/components/mini-card/mini-card-preview.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/ai-prompt-dialog.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/import-prompt/index.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/import-prompt/upload-prompt.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/variable-setting/add-variable-modal.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/mini-card/mini-card-form-modal.vue', () => ({ default: { template: '<div />' } }))
const input = { props: ['value'], emits: ['update:value'], template: '<input :value="value" @input="$emit(\'update:value\', $event.target.value)" />' }
const container = { template: '<div><slot /></div>' }
const global = { stubs: { 'svg-icon': true, 'a-input': input, 'a-button': { template: '<button><slot /></button>' }, 'a-tabs': container, 'a-tab-pane': true, 'a-tooltip': container, 'a-flex': container, 'a-switch': true } }
const rectDescriptor = Object.getOwnPropertyDescriptor(Range.prototype, 'getBoundingClientRect')
Object.defineProperty(Range.prototype, 'getBoundingClientRect', { configurable: true, value: () => ({ top: 20, bottom: 40, left: 20, right: 20 }) })
afterAll(() => { if (rectDescriptor) Object.defineProperty(Range.prototype, 'getBoundingClientRect', rectDescriptor); else delete Range.prototype.getBoundingClientRect })
const choices = () => getPromptOptions(mocks.robot.chatVariables, mocks.cards.cardList)
const visibleText = element => element.textContent.replace(/\u200b/g, '')
function caret(editor, node = editor.lastChild, offset = node.textContent.length) {
  editor.setAttribute('tabindex', '0')
  editor.focus()
  const range = document.createRange()
  if (node.nodeType === 3) range.setStart(node, offset)
  else range.setStartAfter(node)
  range.collapse(true)
  const selection = window.getSelection()
  selection.removeAllRanges()
  selection.addRange(range)
}
function render(props = {}) {
  return mount(AtInput, { props: { type: 'textarea', semanticTags: true, options: choices().variables, cardOptions: choices().cards, ...props }, global, attachTo: document.body })
}

describe('prompt variables keep semantic labels and encoded values', () => {
  beforeEach(() => {
    mocks.robot = reactive({ chatVariables: [{ id: 8, variable_key: 'member_level', variable_name: '会员等级' }], fetchChatVariables: vi.fn() })
    mocks.cards = reactive({ cardList: [{ id: 12, title: '商品说明' }], fetchCardList: vi.fn(), removeCard: vi.fn() })
  })

  it('renders system/custom variables and cards as labels while serializing original codes', () => {
    const value = '当前：【system_var:current_time】\n等级：【chat_variable:member_level】\n推荐：【mini_card:12】'
    const wrapper = render({ defaultValue: value })
    expect(wrapper.find('.j-mention').text()).toContain('当前时间-系统变量')
    expect(wrapper.find('.j-mention').text()).toContain('会员等级')
    expect(wrapper.find('.mini-card-tag').text()).toBe('小程序卡片：商品说明')
    expect(wrapper.find('.j-mention').text()).not.toContain('【')
    wrapper.vm.updateValue()
    expect(wrapper.emitted('change').at(-1)[0]).toBe(value)
    wrapper.unmount()
  })

  it('refreshes asynchronously loaded names without losing unsaved text or the caret', async () => {
    const wrapper = render({ defaultValue: '开头【chat_variable:member_level】', options: [] })
    expect(wrapper.find('.j-mention').text()).toContain('未找到的变量')
    const editor = wrapper.find('.j-mention').element
    editor.append(document.createTextNode('未保存'))
    caret(editor)
    const offsets = getPromptSelection(editor)
    await wrapper.setProps({ options: choices().variables })
    expect(visibleText(editor)).toBe('开头会员等级未保存')
    expect(getPromptSelection(editor)).toEqual(offsets)
    wrapper.vm.updateValue()
    expect(wrapper.emitted('change').at(-1)[0]).toBe('开头【chat_variable:member_level】未保存')
    wrapper.unmount()
  })

  it('inserts at the slash position after focusing and searching the popup', async () => {
    const wrapper = mount(PromptVariableInput, { props: { type: 'textarea', defaultValue: '前/后', options: choices().variables, cardOptions: choices().cards }, global, attachTo: document.body })
    const mention = wrapper.findComponent(AtInput)
    const editor = wrapper.find('.j-mention').element
    caret(editor, editor.firstChild, 2)
    await mention.vm.handleDropdownLogic({ inputType: 'insertText' })
    await flushPromises()
    const search = document.querySelector('.workspace-slash-panel .variable-search input')
    search.focus()
    search.value = 'member_level'
    search.dispatchEvent(new Event('input', { bubbles: true }))
    await flushPromises()
    const options = [...document.querySelectorAll('.workspace-slash-panel .tag-item')]
    expect(options).toHaveLength(1)
    expect(options[0].textContent).toContain('会员等级')
    options[0].click()
    await flushPromises()
    expect(visibleText(editor)).toBe('前会员等级后')
    expect(mention.emitted('change').at(-1)[0]).toBe('前【chat_variable:member_level】后')
    wrapper.unmount()
  })

  it('inserts card chips, preserves adjacent text and removes deleted tokens from saved content', () => {
    const wrapper = render({ defaultValue: '/' })
    const editor = wrapper.find('.j-mention').element
    caret(editor)
    wrapper.vm.atText = ''
    wrapper.vm.selectCard(choices().cards[0])
    expect(wrapper.find('.mini-card-tag').text()).toBe('小程序卡片：商品说明')
    expect(wrapper.emitted('change').at(-1)[0]).toBe('【mini_card:12】')
    wrapper.find('.mini-card-tag').element.remove()
    wrapper.vm.updateValue()
    expect(wrapper.emitted('change').at(-1)[0]).toBe('')
    wrapper.unmount()
  })

  it('keeps labels and pasted markup as safe text and preserves nested line breaks', async () => {
    const options = [{ value: '【chat_variable:x.y】', label: '<img src=x onerror=alert(1)>' }]
    const wrapper = render({ defaultValue: '【chat_variable:x.y】', options })
    expect(wrapper.find('.j-mention img').exists()).toBe(false)
    expect(visibleText(wrapper.find('.j-mention').element)).toBe('<img src=x onerror=alert(1)>')
    const editor = wrapper.find('.j-mention').element
    caret(editor)
    await wrapper.find('.j-mention').trigger('paste', { clipboardData: { getData: () => '<script>safe text</script>' } })
    expect(editor.querySelector('script')).toBeNull()
    const block = document.createElement('div')
    block.innerHTML = '<div>第一行</div><div><span class="j-mention-at" data-value="【system_var:current_date】">当前日期</span><br>第二行</div>'
    expect(readPromptText(block)).toBe('第一行\n【system_var:current_date】\n第二行')
    wrapper.unmount()
  })

  it('uses the same labels in the read-only preview and retains missing tokens for saving', () => {
    const options = [...choices().variables, ...choices().cards]
    const wrapper = mount(PromptTokenText, { props: { value: '【system_var:current_date】【mini_card:12】【chat_variable:deleted】', options }, global })
    expect(wrapper.text()).toContain('当前日期-系统变量')
    expect(wrapper.text()).toContain('小程序卡片：商品说明')
    expect(wrapper.text()).toContain('未找到的变量')
    expect(wrapper.text()).not.toContain('【')
    wrapper.unmount()
  })

  it('flushes the latest editor DOM into a private encoded draft before saving', async () => {
    const sections = Object.fromEntries(['role', 'task', 'constraints', 'skill', 'output', 'tone'].map(key => [key, { subject: key, describe: '' }]))
    const robotInfo = { prompt: '', prompt_type: 0, prompt_struct: { ...sections, custom: [] }, prompt_struct_default: sections, work_flow_ids: '' }
    const wrapper = mount(PromptEditor, { props: { robotInfo }, global, attachTo: document.body })
    const editor = wrapper.find('.diy-prompt-input .j-mention').element
    editor.textContent = '刚输入的【system_var:current_time】'
    expect(wrapper.vm.getDraft().prompt).toBe('刚输入的【system_var:current_time】')
    expect(robotInfo.prompt).toBe('')
    await wrapper.find('.diy-prompt-input .j-mention').trigger('blur')
    expect(visibleText(editor)).toBe('刚输入的当前时间-系统变量')
    wrapper.unmount()
  })

  it('places a real editable caret after a loaded tag and saves text without caret placeholders', async () => {
    const code = '【system_var:current_time】'
    const wrapper = render({ defaultValue: code })
    const token = wrapper.find('.j-mention-at').element
    caret(wrapper.find('.j-mention').element)
    Object.defineProperty(token, 'getBoundingClientRect', { value: () => ({ left: 20, right: 120, width: 100 }) })
    await wrapper.find('.j-mention-at').trigger('click', { clientX: 119 })
    const selection = window.getSelection()
    expect(selection.anchorNode).toBe(token.nextSibling.firstChild)
    expect(selection.anchorOffset).toBe(1)
    const range = selection.getRangeAt(0)
    const text = document.createTextNode('继续提问')
    range.insertNode(text)
    range.setStartAfter(text)
    range.collapse(true)
    selection.removeAllRanges()
    selection.addRange(range)
    wrapper.vm.updateValue()
    expect(wrapper.emitted('change').at(-1)[0]).toBe(code + '继续提问')
    await wrapper.setProps({ options: [...choices().variables] })
    expect(visibleText(wrapper.find('.j-mention').element)).toBe('当前时间-系统变量继续提问')
    expect(getPromptSelection(wrapper.find('.j-mention').element).start).toBe((code + '继续提问').length)
    wrapper.unmount()
  })

  it('inserts consecutive variables after tags without nesting chips inside caret slots', () => {
    const wrapper = render({ defaultValue: '/' })
    const editor = wrapper.find('.j-mention').element
    caret(editor)
    wrapper.vm.selectOption(choices().variables[0])
    const slot = editor.querySelector('.j-mention-at').nextSibling
    slot.firstChild.appendData('/')
    caret(editor, slot.firstChild, 2)
    wrapper.vm.selectOption(choices().variables[1])
    expect(wrapper.emitted('change').at(-1)[0]).toBe('【system_var:current_time】【system_var:current_date】')
    expect(editor.querySelector('[data-prompt-caret] .j-mention-at')).toBeNull()
    expect(editor.querySelectorAll('.j-mention-at')).toHaveLength(2)
    wrapper.unmount()
  })

  it('deletes a variable with one Backspace after its editable caret slot', async () => {
    const wrapper = render({ defaultValue: '【system_var:current_time】' })
    const editor = wrapper.find('.j-mention').element
    const slot = editor.querySelector('.j-mention-at').nextSibling
    caret(editor, slot.firstChild, 1)
    await wrapper.find('.j-mention').trigger('keydown', { key: 'Backspace' })
    expect(editor.querySelector('.j-mention-at')).toBeNull()
    expect(wrapper.emitted('change').at(-1)[0]).toBe('')
    wrapper.unmount()
  })

  it('mounts the custom prompt popup on body even when the editor is inside a clipped container', async () => {
    const sections = Object.fromEntries(['role', 'task', 'constraints', 'skill', 'output', 'tone'].map(key => [key, { subject: key, describe: '' }]))
    const wrapper = mount(PromptEditor, { props: { robotInfo: { prompt: '/', prompt_type: 0, prompt_struct: { ...sections, custom: [] }, prompt_struct_default: sections, work_flow_ids: '' } }, global, attachTo: document.body })
    wrapper.element.style.overflow = 'hidden'
    const editor = wrapper.find('.j-mention').element
    caret(editor)
    await wrapper.findComponent(AtInput).vm.handleDropdownLogic({ inputType: 'insertText' })
    expect(document.querySelector('.dropdown-list-semantic').parentElement).toBe(document.body)
    wrapper.unmount()
  })

  it('keeps the popup open on editor/popup scroll and repositions it inside the viewport on resize', async () => {
    const wrapper = render({ defaultValue: '/' })
    const editor = wrapper.find('.j-mention').element
    caret(editor)
    await wrapper.vm.handleDropdownLogic({ inputType: 'insertText' })
    const popup = document.querySelector('.dropdown-list-semantic')
    Object.defineProperties(popup, { offsetWidth: { value: 640 }, offsetHeight: { value: 240 }, scrollHeight: { configurable: true, value: 240 } })
    let anchor = { top: 100, bottom: 120, left: 900 }
    Object.defineProperty(wrapper.vm.dropdownAnchorRange, 'getBoundingClientRect', { value: () => anchor })
    wrapper.vm.positionSemanticDropdown()
    await flushPromises()
    expect(popup.style.top).toBe('124px')
    expect(Number.parseFloat(popup.style.left)).toBeLessThanOrEqual(window.innerWidth - 640 - 8)
    anchor = { top: 300, bottom: 320, left: 40 }
    editor.dispatchEvent(new Event('scroll'))
    await flushPromises()
    expect(wrapper.vm.showDropdown).toBe(true)
    expect(document.querySelector('.dropdown-list-semantic')).toBe(popup)
    expect(popup.style.top).toBe('324px')
    popup.dispatchEvent(new Event('scroll'))
    expect(wrapper.vm.showDropdown).toBe(true)
    anchor = { top: window.innerHeight - 68, bottom: window.innerHeight - 48, left: 40 }
    window.dispatchEvent(new Event('resize'))
    await flushPromises()
    expect(Number.parseFloat(popup.style.top) + 240).toBeLessThanOrEqual(anchor.top)
    Object.defineProperty(popup, 'scrollHeight', { value: 1000 })
    window.dispatchEvent(new Event('resize'))
    await flushPromises()
    expect(Number.parseFloat(popup.style.top)).toBeGreaterThanOrEqual(8)
    expect(Number.parseFloat(popup.style.maxHeight)).toBeLessThanOrEqual(window.innerHeight - 16)
    document.body.click()
    await flushPromises()
    expect(wrapper.vm.showDropdown).toBe(false)
    wrapper.unmount()
  })

  it('preserves the legacy behavior of closing its popup on outside scroll', async () => {
    const wrapper = render({ semanticTags: false, defaultValue: '/' })
    wrapper.vm.showDropdown = true
    await flushPromises()
    wrapper.vm.onScroll({ target: document.body })
    expect(wrapper.vm.showDropdown).toBe(false)
    wrapper.unmount()
  })
})
