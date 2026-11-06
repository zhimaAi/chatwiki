import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
const api = vi.hoisted(() => ({ edit: vi.fn(), categories: vi.fn(), groups: vi.fn(), route: { name: 'customerServiceKnowledgeBase', query: { id: '236' } } }))
vi.mock('@/api/library', () => ({ editParagraph: api.edit, getCategoryList: api.categories, saveCategoryParagraph: api.edit, getLibraryGroup: api.groups, generateSimilarQuestions: vi.fn(), getSimilarQuestions: vi.fn().mockResolvedValue({ data: { count: 0 } }) }))
vi.mock('vue-router', () => ({ useRoute: () => api.route, useRouter: () => ({ resolve: vi.fn() }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('ant-design-vue', () => ({ message: { error: vi.fn(), success: vi.fn() } }))
vi.mock('@/components/upload-img/index.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/mini-card/mini-card-tab-content.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { template: '<div />' } }))
import Editor from './edit-subsection.vue'

function mountEditor(props = {}) {
  return shallowMount(Editor, {
    props: { detailsInfo: { is_qa_doc: '1', library_id: '501' }, ...props },
    global: { stubs: {
      AModal: { template: '<div><slot /><slot name="footer" /></div>' },
      AInput: { props: ['value', 'prefix'], template: '<input />' },
      AInputNumber: { props: ['value', 'prefix'], template: '<input />' },
      ASelect: { props: ['value', 'children'], template: '<div><slot /></div>' },
      AButton: { template: '<button><slot /></button>' }
    } }
  })
}

describe('paragraph resource context', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.edit.mockResolvedValue({ res: 0 })
    api.categories.mockResolvedValue({ data: [] })
    api.groups.mockResolvedValue({ data: [] })
  })
  it.each(['customerServiceKnowledgeBase', 'knowledgeDocument', 'libraryConfig'])('creates QA in its explicit library, not the route ID, on %s', async name => {
    api.route.name = name
    const wrapper = mountEditor()
    wrapper.vm.showModal({ question: '1+1=?', answer: '2' })
    await flushPromises()
    await wrapper.findAll('button').find(button => button.text() === 'btn_confirm').trigger('click')
    await flushPromises()
    expect(api.categories).toHaveBeenCalledWith({ library_id: '501' })
    const form = api.edit.mock.calls[0][0]
    expect(form.get('library_id')).toBe('501')
    expect(form.has('file_id')).toBe(false)
    wrapper.unmount()
  })
  it('preserves a real file context for the standalone document preview', async () => {
    api.route.name = 'libraryPreview'
    const wrapper = mountEditor({ fileId: '700' })
    wrapper.vm.showModal({ question: 'File question', answer: 'Answer' })
    await flushPromises()
    await wrapper.findAll('button').find(button => button.text() === 'btn_confirm').trigger('click')
    await flushPromises()
    expect(api.categories).toHaveBeenCalledWith({ file_id: '700' })
    expect(api.edit.mock.calls[0][0].get('file_id')).toBe('700')
    expect(api.edit.mock.calls[0][0].has('library_id')).toBe(false)
    wrapper.unmount()
  })
  it('preserves paragraph identity and group zero while editing in an embedded library', async () => {
    api.route.name = 'customerServiceKnowledgeBase'
    const wrapper = mountEditor()
    wrapper.vm.showModal({ id: '42', question: 'Existing', answer: 'Updated', group_id: 0 })
    await flushPromises()
    await wrapper.findAll('button').find(button => button.text() === 'btn_confirm').trigger('click')
    await flushPromises()
    const form = api.edit.mock.calls[0][0]
    expect(form.get('id')).toBe('42')
    expect(form.get('group_id')).toBe('0')
    expect(form.get('library_id')).toBe('501')
    expect(form.has('file_id')).toBe(false)
    wrapper.unmount()
  })
})
