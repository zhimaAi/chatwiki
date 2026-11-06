import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { message, Modal } from 'ant-design-vue'

const api = vi.hoisted(() => ({ files: vi.fn(), groups: vi.fn(), remove: vi.fn(), add: vi.fn(), edit: vi.fn(), graph: vi.fn(), crawl: vi.fn(), cancel: vi.fn(), deleteGroup: vi.fn(), reorder: vi.fn(), relearn: vi.fn(), custom: vi.fn(), feishu: vi.fn() }))
vi.mock('@/api/library', () => ({ getLibraryFileList: api.files, getLibraryGroup: api.groups, delLibraryFile: api.remove, addLibraryFile: api.add, editLibrary: api.edit, createGraph: api.graph, manualCrawl: api.crawl, cancelOcrPdf: api.cancel, deleteLibraryGroup: api.deleteGroup, sortLibararyListGroup: api.reorder, restudyLibraryFile: api.relearn }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/stores/modules/user', () => ({ useUserStore: () => ({ getToken: 'test-token' }) }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => ({ companyInfo: { neo4j_status: 'false' } }) }))
vi.mock('@/stores/modules/library', () => ({ useLibraryStore: () => ({ changeGraphSwitch: vi.fn() }) }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: {} }), useRouter: () => ({ push: vi.fn(), replace: vi.fn(), resolve: value => ({ href: JSON.stringify(value) }) }) }))
vi.mock('@/router', () => ({ default: { push: vi.fn() } }))
vi.mock('@/components/model-select/model-select.vue', () => ({ default: { name: 'ModelSelect', props: ['modelType'], emits: ['loaded'], template: '<div />' } }))
vi.mock('@/views/library/library-details/components/add-custom-document.vue', () => ({ default: { name: 'AddCustomDocument', setup: (_, { expose }) => { expose({ add: api.custom }) }, template: '<div />' } }))
vi.mock('@/views/library/library-details/components/add-feishu-document.vue', () => ({ default: { name: 'AddFeishuDocument', setup: (_, { expose }) => { expose({ show: api.feishu }) }, template: '<div />' } }))
import KnowledgeDocument from '@/views/library/library-details/knowledge-document.vue'
import WorkspaceDocumentView from '@/views/library/library-details/components/workspace-document-view.vue'
import ModelSelect from '@/components/model-select/model-select.vue'

const wrappers = []
async function render(compact = true) {
  const wrapper = shallowMount(KnowledgeDocument, { props: { compact, library_id: '101' }, global: { directives: { copy: {} }, stubs: {
    ModelSelect: false, AddCustomDocument: false, AddFeishuDocument: false,
    'a-modal': { name: 'AModal', props: ['open', 'title'], template: '<div v-if="open"><slot /></div>' },
    'a-table': { template: '<div />' }, 'a-tooltip': { template: '<div><slot /></div>' }, 'a-dropdown': { template: '<div><slot /></div>' },
    'a-flex': { template: '<div><slot /></div>' }, 'a-select': { template: '<select />' }, 'a-select-option': { template: '<option />' }
  } } })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}
beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  api.files.mockImplementation(async () => ({ data: { info: { id: 101, type: 0, use_model: 'embedding', model_config_id: 1, graph_switch: '0' }, list: [{ id: 11, file_name: '真实文档.pdf', file_ext: 'pdf', status: '2', doc_type: 1 }], total: 1, count_data: {} } }))
  api.groups.mockResolvedValue({ data: [{ id: 0, group_name: '未分组', total: 0 }, { id: 7, group_name: '产品资料', total: 1 }] })
  api.reorder.mockResolvedValue({})
  api.remove.mockResolvedValue({})
  vi.spyOn(message, 'success').mockImplementation(() => {})
  vi.spyOn(Modal, 'confirm').mockImplementation(() => {})
  vi.stubGlobal('ResizeObserver', class { observe() {} disconnect() {} })
  vi.stubGlobal('IntersectionObserver', class { observe() {} disconnect() {} })
})
afterEach(() => { wrappers.splice(0).forEach(wrapper => wrapper.unmount()); vi.restoreAllMocks(); vi.unstubAllGlobals() })

describe('normal document integration', () => {
  it('renders only the compact workspace and keeps model loading available for all add sources', async () => {
    const wrapper = await render()
    expect(wrapper.find('.group-content-box').exists()).toBe(false)
    expect(wrapper.findAllComponents(WorkspaceDocumentView)).toHaveLength(1)
    const model = wrapper.findComponent(ModelSelect)
    expect(model.props('modelType')).toBe('TEXT EMBEDDING')
    model.vm.$emit('loaded', [{ children: [{ name: 'embedding', model_config_id: 1 }] }])
    await flushPromises()
    const view = wrapper.findComponent(WorkspaceDocumentView)
    view.vm.$emit('add', { key: 1 })
    await flushPromises()
    expect(wrapper.findAllComponents({ name: 'AModal' }).find(item => item.props('title') === 'title_upload_doc').props('open')).toBe(true)
    view.vm.$emit('add', { key: 2 })
    await flushPromises()
    expect(wrapper.findAllComponents({ name: 'AModal' }).some(item => item.props('open') && item.props('title') !== 'title_upload_doc')).toBe(true)
    view.vm.$emit('add', { key: 3 }); view.vm.$emit('add', { key: 4 })
    expect(api.custom).toHaveBeenCalledTimes(1)
    expect(api.feishu).toHaveBeenCalledTimes(1)
    expect(Modal.confirm).not.toHaveBeenCalled()
    expect(api.edit).not.toHaveBeenCalled()
  })

  it('persists compact navigation separately and preserves group/search/status when switching views', async () => {
    const wrapper = await render()
    const view = wrapper.findComponent(WorkspaceDocumentView)
    view.vm.$emit('select-group', { id: 7 })
    view.vm.$emit('update:fileName', '手册'); view.vm.$emit('update:status', '2'); view.vm.$emit('search')
    await flushPromises()
    view.vm.$emit('view-change', 'folder')
    await flushPromises()
    expect(view.props()).toMatchObject({ groupId: 7, fileName: '手册', status: '2', viewMode: 'folder' })
    expect(localStorage.getItem('customer_service_document_view_mode_101')).toBe('folder')
    expect(localStorage.getItem('qa_document_view_mode_101')).toBeNull()
    view.vm.$emit('toggle-groups')
    await flushPromises()
    expect(view.props('collapsed')).toBe(true)
    expect(localStorage.getItem('customer_service_document_group_hide_key')).toBe('1')
    view.vm.$emit('page-change', 2, 20)
    await flushPromises()
    expect(api.files.mock.calls.at(-1)[0]).toMatchObject({ library_id: '101', group_id: 7, file_name: '手册', status: '2', page: 2, size: 20 })
  })

  it('uses existing delete/reorder APIs and preserves the legacy entry without compact UI', async () => {
    const wrapper = await render()
    const view = wrapper.findComponent(WorkspaceDocumentView)
    view.vm.$emit('action', { key: 'delete', record: { id: 11 } })
    expect(api.remove).not.toHaveBeenCalled()
    Modal.confirm.mock.calls[0][0].onOk()
    await flushPromises()
    expect(api.remove).toHaveBeenCalledWith({ id: 11 })
    view.vm.$emit('reorder-groups', [...view.props('groups')])
    await flushPromises()
    expect(api.reorder.mock.calls[0][0]).toMatchObject({ library_id: '101', sense: 1 })
    const legacy = await render(false)
    expect(legacy.findComponent(WorkspaceDocumentView).exists()).toBe(false)
    expect(legacy.find('.group-content-box').exists()).toBe(true)
    expect(api.files.mock.calls.at(-1)[0].size).toBe(10)
    expect(localStorage.getItem('qa_document_view_mode_101')).toBe('list')
  })
})
