import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'

const state = vi.hoisted(() => ({ form: null }))
const api = vi.hoisted(() => ({ info: vi.fn(), edit: vi.fn(), separators: vi.fn(), push: vi.fn() }))
vi.mock('@/api/library', () => ({ getLibraryInfo: api.info, editLibrary: api.edit, getSeparatorsList: api.separators }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: { id: '7' } }), useRouter: () => ({ push: api.push }) }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => ({ companyInfo: {} }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('ant-design-vue', () => ({
  Form: { useForm: form => { state.form = form; return { validateInfos: {} } } },
  message: { error: vi.fn(), success: vi.fn() }, Modal: { confirm: vi.fn() }
}))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { props: ['modeName', 'modeId'], emits: ['change', 'update:modeName', 'update:modeId'], template: '<div />' } }))
vi.mock('@/views/library/library-details/components/open-grapg-modal.vue', () => ({ default: { template: '<div />' } }))

import Settings from '@/views/library/library-details/knowledge-config.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import AvatarInput from '@/views/library/add-library/components/avatar-input.vue'

const container = { template: '<div><slot /><slot name="action" /><slot name="label" /></div>' }
const mountSettings = (props = { libraryId: '101', embedded: true }) => shallowMount(Settings, {
  props, global: { stubs: {
    'cu-scroll': container, 'a-form': container, 'a-form-item': container, 'a-flex': container,
    'a-input': { name: 'AInput', props: ['value'], emits: ['blur', 'update:value'], template: '<input />' },
    'a-input-number': { inheritAttrs: false, props: ['value'], template: '<input />' },
    'a-select': { ...container, inheritAttrs: false, props: ['value'] },
    'a-alert': { ...container, name: 'AAlert', props: ['message'] },
    'a-tooltip': { ...container, inheritAttrs: false, props: ['title'] },
    'a-button': { template: '<button><slot /></button>' }
  } }
})
const info = (type = 2) => ({ data: {
  id: type === 2 ? 101 : 202, type, library_name: 'Library', library_intro: 'Intro', avatar: 'cover.png',
  use_model: 'embedding-one', model_config_id: 11, is_offline: false, group_id: 0, graph_switch: '0',
  qa_index_type: 1, chunk_type: 1, normal_chunk_default_chunk_size: 512,
  normal_chunk_default_chunk_overlap: 50, normal_chunk_default_not_merged_text: 'false', is_permanent: 1
} })
const blurName = wrapper => wrapper.findComponent({ name: 'AInput' }).vm.$emit('blur')

describe('workspace library settings', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.info.mockReset().mockResolvedValue(info())
    api.edit.mockReset().mockResolvedValue({ res: 0 })
    api.separators.mockResolvedValue({ data: [{ name: 'Separator', no: 12 }] })
  })

  it('loads and saves the explicit QA library instead of the customer ID in the route', async () => {
    const wrapper = mountSettings()
    await flushPromises()
    expect(api.info).toHaveBeenCalledWith({ id: '101' })
    expect(wrapper.find('.indexing-methods-box').exists()).toBe(true)
    expect(wrapper.find('.select-card-box').exists()).toBe(false)
    state.form.library_name = 'Edited QA'
    blurName(wrapper)
    await wrapper.vm.waitForSave()
    expect(api.edit).toHaveBeenCalledWith(expect.objectContaining({ id: '101', library_name: 'Edited QA', use_model: 'embedding-one', qa_index_type: 1 }))
    expect(wrapper.emitted('saved')[0]).toEqual(['101'])
    wrapper.unmount()
  })

  it('uses normal-library segmentation, model and cover fields with the correct library ID', async () => {
    api.info.mockResolvedValue(info(0))
    const wrapper = mountSettings({ libraryId: '202', embedded: true })
    await flushPromises()
    expect(wrapper.find('.indexing-methods-box').exists()).toBe(false)
    expect(wrapper.find('.select-card-box').exists()).toBe(true)
    await wrapper.findAll('.select-card-item')[1].trigger('click')
    await flushPromises()
    const model = wrapper.findComponent(ModelSelect)
    model.vm.$emit('update:modeName', 'embedding-two')
    model.vm.$emit('update:modeId', 22)
    model.vm.$emit('change')
    const cover = new File(['image'], 'cover.png', { type: 'image/png' })
    wrapper.findComponent(AvatarInput).vm.$emit('change', { imageUrl: 'data:image/png;base64,aQ==', file: cover })
    await wrapper.vm.waitForSave()
    expect(api.edit).toHaveBeenLastCalledWith(expect.objectContaining({ id: '202', chunk_type: 2, use_model: 'embedding-two', model_config_id: 22, avatar: cover }))
    wrapper.unmount()
  })

  it('serializes independent configuration snapshots and lets closing wait for the last save', async () => {
    let finishFirst
    let finishSecond
    api.edit.mockImplementationOnce(() => new Promise(resolve => { finishFirst = resolve }))
      .mockImplementationOnce(() => new Promise(resolve => { finishSecond = resolve }))
    const wrapper = mountSettings()
    await flushPromises()
    state.form.library_name = 'First'
    blurName(wrapper)
    state.form.library_name = 'Second'
    blurName(wrapper)
    await flushPromises()
    expect(api.edit).toHaveBeenCalledTimes(1)
    expect(api.edit.mock.calls[0][0].library_name).toBe('First')
    const done = vi.fn()
    const closing = wrapper.vm.waitForSave().then(done)
    finishFirst({ res: 0 })
    await flushPromises()
    expect(api.edit.mock.calls[1][0].library_name).toBe('Second')
    expect(done).not.toHaveBeenCalled()
    finishSecond({ res: 0 })
    await closing
    expect(done).toHaveBeenCalledTimes(1)
    wrapper.unmount()
  })

  it('shows failed saves, keeps edited values and retries the same library', async () => {
    api.edit.mockRejectedValueOnce(new Error('save failed'))
    const wrapper = mountSettings()
    await flushPromises()
    state.form.library_intro = 'Keep this edit'
    blurName(wrapper)
    await flushPromises()
    await expect(wrapper.vm.waitForSave()).rejects.toThrow('save failed')
    expect(wrapper.findComponent({ name: 'AAlert' }).props('message')).toBe('save failed')
    await wrapper.find('.settings-error a').trigger('click')
    await wrapper.vm.waitForSave()
    await flushPromises()
    expect(api.edit).toHaveBeenLastCalledWith(expect.objectContaining({ id: '101', library_intro: 'Keep this edit' }))
    expect(wrapper.find('.settings-error').exists()).toBe(false)
    wrapper.unmount()
  })

  it('retries initial loading and ignores a response after the modal is destroyed', async () => {
    api.info.mockRejectedValueOnce(new Error('load failed'))
    const wrapper = mountSettings()
    await flushPromises()
    expect(wrapper.findComponent({ name: 'AAlert' }).props('message')).toBe('load failed')
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(wrapper.find('.indexing-methods-box').exists()).toBe(true)
    wrapper.unmount()
    let finish
    api.info.mockImplementationOnce(() => new Promise(resolve => { finish = resolve }))
    const pending = mountSettings({ libraryId: '202', embedded: true })
    const form = state.form
    pending.unmount()
    finish(info(0))
    await flushPromises()
    expect(form.library_name).toBe('')
    expect(api.edit).not.toHaveBeenCalled()
  })

  it('preserves the legacy route ID and immediate autosave when no props are passed', async () => {
    const wrapper = mountSettings({})
    await flushPromises()
    expect(api.info).toHaveBeenCalledWith({ id: '7' })
    state.form.library_name = 'Legacy'
    blurName(wrapper)
    expect(api.edit).toHaveBeenCalledWith(expect.objectContaining({ id: '7', library_name: 'Legacy' }))
    expect(wrapper.find('.settings-feedback').exists()).toBe(false)
    await flushPromises()
    wrapper.unmount()
  })
})
