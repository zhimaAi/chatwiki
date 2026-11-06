import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive } from 'vue'
import AiPromptDialog from './components/ai-prompt-dialog.vue'

const mocks = vi.hoisted(() => ({ generate: vi.fn() }))
vi.mock('@/api/robot', () => ({ createPromptByAi: mocks.generate }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => ({ chatVariables: [] }) }))
vi.mock('@/stores/modules/mini-card', () => ({ useMiniCardStore: () => ({ cardList: [] }) }))
const stubs = {
  'svg-icon': true, 'a-spin': true, 'a-alert': { props: ['message'], template: '<div role="alert">{{ message }}</div>' },
  'a-input': { props: ['value'], emits: ['update:value'], template: '<input :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-button': { props: ['disabled', 'loading'], template: '<button :disabled="disabled || loading"><slot /></button>' },
  'a-modal': { name: 'AModal', props: ['open'], emits: ['cancel'], template: '<section v-if="open"><slot /><footer><slot name="footer" /></footer></section>' }
}
const response = () => ({ res: 0, data: { promptStruct: JSON.stringify({ role: { subject: '角色', describe: '生成角色' }, task: { subject: '任务', describe: '生成任务' }, output: { subject: '输出格式', describe: '输出要求' } }) } })
const deferred = () => { let resolve; const promise = new Promise(done => { resolve = done }); return { promise, resolve } }
function render() { return mount(AiPromptDialog, { global: { provide: { customerRobotContext: { robotInfo: reactive({ id: 7 }) } }, stubs } }) }
async function show(wrapper) { wrapper.vm.show(); await flushPromises() }
async function generate(wrapper) { await wrapper.find('input').setValue('客服需求'); await wrapper.find('.demand-row button').trigger('click'); await flushPromises() }

describe('workspace AI prompt generation', () => {
  beforeEach(() => { vi.clearAllMocks(); mocks.generate.mockResolvedValue(response()) })

  it('starts empty and fills quick prompts without generating until requested', async () => {
    const wrapper = render()
    await show(wrapper)
    expect(wrapper.find('footer button:last-child').attributes('disabled')).toBeDefined()
    expect(wrapper.find('.empty-prompt').text()).toBe('prompt_ai_empty')
    await wrapper.find('.quick-prompts button').trigger('click')
    expect(wrapper.find('input').element.value).toBe('quick_ecommerce_after_sales')
    expect(mocks.generate).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('uses the active robot id, previews the result and applies only on explicit selection', async () => {
    const wrapper = render()
    await show(wrapper)
    await generate(wrapper)
    expect(mocks.generate).toHaveBeenCalledWith({ id: 7, demand: '客服需求' })
    expect(wrapper.find('.generated-prompt.ready').text()).toContain('生成角色')
    expect(wrapper.emitted('handleAiSave')).toBeUndefined()
    await wrapper.find('footer button:last-child').trigger('click')
    expect(wrapper.emitted('handleAiSave')[0][0]).toEqual(expect.objectContaining({ role: { subject: '角色', describe: '生成角色' }, skill: { subject: '', describe: '' }, custom: [] }))
    expect(wrapper.findComponent({ name: 'AModal' }).exists()).toBe(false)
    wrapper.unmount()
  })

  it('ignores late results after cancel/reopen and blocks duplicate generation', async () => {
    const wrapper = render()
    const pending = deferred()
    mocks.generate.mockReturnValueOnce(pending.promise)
    await show(wrapper)
    await generate(wrapper)
    await wrapper.find('.demand-row button').trigger('click')
    expect(mocks.generate).toHaveBeenCalledOnce()
    await wrapper.find('footer button:first-child').trigger('click')
    await show(wrapper)
    pending.resolve(response())
    await flushPromises()
    expect(wrapper.find('.generated-prompt.ready').exists()).toBe(false)
    expect(wrapper.find('footer button:last-child').attributes('disabled')).toBeDefined()
    wrapper.unmount()
  })

  it('keeps use disabled on failures and permits retrying a malformed result', async () => {
    const wrapper = render()
    mocks.generate.mockResolvedValueOnce({ data: { promptStruct: 'invalid json' } })
    await show(wrapper)
    await generate(wrapper)
    expect(wrapper.find('[role="alert"]').exists()).toBe(true)
    expect(wrapper.find('footer button:last-child').attributes('disabled')).toBeDefined()
    await generate(wrapper)
    expect(wrapper.find('.generated-prompt.ready').exists()).toBe(true)
    wrapper.unmount()
  })
})
