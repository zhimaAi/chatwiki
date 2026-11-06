import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import RecallDialog from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import WeightSelect from '@/components/weight-select/index.vue'

const mocks = vi.hoisted(() => ({ error: vi.fn(), verify: vi.fn() }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('ant-design-vue', async importOriginal => ({ Modal: (await importOriginal()).Modal, message: { error: mocks.error } }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => ({ companyInfo: { neo4j_status: 'true' } }) }))
vi.mock('@/api/library', () => ({ getRobotMetaSchemaList: vi.fn().mockResolvedValue({ data: [] }) }))
vi.mock('@/api/robot', () => ({ getChatVariables: vi.fn().mockResolvedValue({ data: [] }) }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { name: 'ModelSelect', props: ['modeId', 'modeName', 'modelType'], emits: ['update:modeId', 'update:modeName'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/meta-filter-box.vue', () => ({ default: { name: 'MetaFilterBox', props: ['rule', 'type'], setup: (_, { expose }) => expose({ verify: mocks.verify }), template: '<div />' } }))
const group = { template: '<div><slot /></div>' }
const stubs = {
  'svg-icon': true, 'a-tooltip': group,
  'a-modal': { name: 'AModal', props: ['open', 'width', 'confirmLoading', 'closable'], emits: ['ok', 'cancel', 'update:open'], template: '<section v-if="open"><slot /></section>' },
  'a-slider': { name: 'ASlider', props: ['value', 'min', 'max', 'step', 'range'], emits: ['update:value', 'change'], template: '<div />' },
  'a-input-number': { name: 'AInputNumber', props: ['value', 'min', 'max', 'precision'], emits: ['update:value'], template: '<input />' },
  'a-switch': { name: 'ASwitch', props: ['checked'], emits: ['update:checked', 'change'], template: '<input />' },
  'a-radio-group': { name: 'ARadioGroup', props: ['value'], emits: ['update:value'], template: '<div><slot /></div>' },
  'a-radio': group, 'a-select': group, 'a-select-option': group
}
const config = () => ({ search_type: 1, top_k: 5, similarity: 0.4, rrf_weight: { vector: 36, search: 28, graph: 36 }, recall_neighbor_top_k: 5, recall_neighbor_before_num: 1, recall_neighbor_after_num: 1, rerank_status: 0, meta_search_switch: 0, meta_search_type: 1, library_search_type: 'fullTextSearch' })
async function render(submit = vi.fn().mockResolvedValue(true)) {
  const data = config()
  const wrapper = mount(RecallDialog, { props: { variant: 'workspace', submit }, global: { stubs }, attachTo: document.body })
  wrapper.vm.open(data, { id: 7, robot_key: 'seven' })
  await flushPromises()
  return { wrapper, data, submit }
}
const named = (wrapper, name) => wrapper.findComponent({ name })
const confirm = async wrapper => { named(wrapper, 'AModal').vm.$emit('ok'); await flushPromises() }
const deferred = () => { let resolve; const promise = new Promise(done => { resolve = done }); return { promise, resolve } }

describe('workspace recall settings', () => {
  beforeEach(() => { vi.clearAllMocks(); mocks.verify.mockReturnValue(true) })

  it('shows the corresponding fields for all four modes and preserves the existing weight control', async () => {
    const { wrapper } = await render()
    expect(named(wrapper, 'AModal').props('width')).toBe(880)
    const states = [[true, true, true], [false, true, false], [false, false, false], [false, false, true]]
    for (let index = 0; index < 4; index++) {
      await wrapper.findAll('.retrieval-mode-item')[index].trigger('click')
      expect(wrapper.findComponent(WeightSelect).exists()).toBe(states[index][0])
      expect(wrapper.findAll('.number-box')).toHaveLength(states[index][1] ? 2 : 1)
      expect(named(wrapper, 'ARadioGroup').exists()).toBe(states[index][2])
      expect(wrapper.findAll('.retrieval-mode-item')[index].attributes('aria-checked')).toBe('true')
    }
    await wrapper.findAll('.retrieval-mode-item')[0].trigger('click')
    const weights = wrapper.findComponent(WeightSelect)
    expect(weights.props('rrf_weight')).toEqual({ vector: 36, search: 28, graph: 36 })
    expect(weights.findAll('.list-item')).toHaveLength(3)
    expect(named(weights, 'ASlider').props('range')).toBe(true)
    expect(named(weights, 'ASlider').props('value')).toEqual([36, 64])
    wrapper.unmount()
  })

  it('keeps the three-way weight values and fulltext selection until confirmation', async () => {
    const { wrapper, submit, data } = await render()
    const slider = named(wrapper.findComponent(WeightSelect), 'ASlider')
    slider.vm.$emit('update:value', [40, 70])
    slider.vm.$emit('change')
    named(wrapper, 'ARadioGroup').vm.$emit('update:value', 'keywordSearch')
    await flushPromises()
    expect(submit).not.toHaveBeenCalled()
    expect(data.rrf_weight).toEqual({ vector: 36, search: 28, graph: 36 })
    await confirm(wrapper)
    expect(submit).toHaveBeenCalledWith(expect.objectContaining({ rrf_weight: { vector: 40, search: 30, graph: 30 }, library_search_type: 'keywordSearch' }))
    wrapper.unmount()
  })

  it('keeps Top K and similarity slider/number values synchronized and supports keyboard mode selection', async () => {
    const { wrapper } = await render()
    const topSlider = wrapper.find('.number-slider-box').findComponent({ name: 'ASlider' })
    topSlider.vm.$emit('update:value', 12)
    await flushPromises()
    expect(named(wrapper, 'AInputNumber').props('value')).toBe(12)
    named(wrapper, 'AInputNumber').vm.$emit('update:value', 27)
    await flushPromises()
    expect(topSlider.props('value')).toBe(27)
    await wrapper.findAll('.retrieval-mode-item')[0].trigger('keydown', { key: 'ArrowDown' })
    expect(wrapper.findAll('.retrieval-mode-item')[2].attributes('aria-checked')).toBe('true')
    expect(document.activeElement).toBe(wrapper.findAll('.retrieval-mode-item')[2].element)
    await wrapper.findAll('.retrieval-mode-item')[2].trigger('keydown', { key: 'Home' })
    expect(wrapper.findAll('.retrieval-mode-item')[0].attributes('tabindex')).toBe('0')
    wrapper.unmount()
  })

  it('validates typed ranges, the enabled Rerank model and metadata rules', async () => {
    const { wrapper, submit } = await render()
    named(wrapper, 'AInputNumber').vm.$emit('update:value', null)
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenLastCalledWith('recall_invalid_parameters')
    named(wrapper, 'AInputNumber').vm.$emit('update:value', 5)
    wrapper.findAllComponents({ name: 'AInputNumber' })[2].vm.$emit('update:value', 11)
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenLastCalledWith('recall_invalid_parameters')
    wrapper.findAllComponents({ name: 'AInputNumber' })[2].vm.$emit('update:value', 5)
    named(wrapper, 'ASwitch').vm.$emit('update:checked', 1)
    await confirm(wrapper)
    expect(mocks.error).toHaveBeenLastCalledWith('msg_please_select_rerank_model')
    named(wrapper, 'ModelSelect').vm.$emit('update:modeId', 9)
    named(wrapper, 'ModelSelect').vm.$emit('update:modeName', 'rerank')
    wrapper.findAllComponents({ name: 'ASwitch' })[1].vm.$emit('update:checked', 1)
    await flushPromises()
    mocks.verify.mockReturnValue(false)
    await confirm(wrapper)
    expect(mocks.verify).toHaveBeenCalledOnce()
    expect(submit).not.toHaveBeenCalled()
    mocks.verify.mockReturnValue(true)
    await confirm(wrapper)
    expect(submit).toHaveBeenCalledWith(expect.objectContaining({ rerank_status: 1, rerank_model_config_id: 9, meta_search_switch: 1 }))
    wrapper.unmount()
  })

  it('retains failed drafts, guards in-flight cancellation and duplicate submissions', async () => {
    const pending = deferred()
    const submit = vi.fn().mockResolvedValueOnce(false).mockReturnValueOnce(pending.promise)
    const { wrapper, data } = await render(submit)
    named(wrapper, 'AInputNumber').vm.$emit('update:value', 17)
    await confirm(wrapper)
    expect(named(wrapper, 'AModal').props('open')).toBe(true)
    expect(data.top_k).toBe(5)
    await confirm(wrapper)
    expect(named(wrapper, 'AModal').props('confirmLoading')).toBe(true)
    named(wrapper, 'AModal').vm.$emit('cancel')
    await confirm(wrapper)
    expect(named(wrapper, 'AModal').props('open')).toBe(true)
    expect(submit).toHaveBeenCalledTimes(2)
    pending.resolve(true)
    await flushPromises()
    expect(named(wrapper, 'AModal').props('open')).toBe(false)
    wrapper.unmount()
  })
})
