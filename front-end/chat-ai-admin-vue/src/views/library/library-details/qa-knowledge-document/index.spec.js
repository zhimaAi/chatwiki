import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
const api = vi.hoisted(() => ({ groups: vi.fn(), list: vi.fn(), remove: vi.fn(), move: vi.fn(), confirm: vi.fn(), selection: { selectedRowKeys: [] }, reset: vi.fn() }))
vi.mock('@/api/library', () => ({ getLibraryGroup: api.groups, getParagraphList: api.list, deleteParagraph: api.remove, setParagraphGroup: api.move, reconstructCategoryVector: vi.fn(), deleteLibraryGroup: vi.fn(), createExportLibFileTask: vi.fn() }))
vi.mock('ant-design-vue', () => ({ Modal: { confirm: api.confirm }, message: { success: vi.fn(), error: vi.fn() } }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: { id: '236' } }) }))
vi.mock('@/router', () => ({ default: { resolve: vi.fn() } }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/constants/index.js', () => ({ KNOWLEDGE_SOURCE_TYPE_MAP: () => ({}) }))
vi.mock('./components/subsection-box.vue', () => ({ default: { name: 'SubsectionBox', setup: (_, { expose }) => { expose({ state: api.selection, resetSelect: api.reset }) }, template: '<div />' } }))
vi.mock('@/components/cu-scroll/cu-scroll.vue', () => ({ default: { template: '<div><slot /></div>' } }))
vi.mock('@/views/library/library-preview/components/edit-subsection.vue', () => ({ default: { template: '<div />' } }))
vi.mock('../components/qa-upload-modal.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/add-group.vue', () => ({ default: { template: '<div />' } }))
vi.mock('../components/group-panel-toggle.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/library/library-details/components/metadata-manage-modal.vue', () => ({ default: { template: '<div />' } }))
import QaDocument from './index.vue'

function mountPage() {
  const wrapper = shallowMount({ ...QaDocument, render: () => null }, { props: { library_id: '501', compact: true } })
  wrapper.vm.$.setupState.subsectionBoxRef = { state: api.selection, resetSelect: api.reset }
  return wrapper
}

describe('QA batch count refresh', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.selection.selectedRowKeys = ['11', '12']
    api.groups.mockResolvedValue({ data: [{ id: 0, group_name: 'Ungrouped', total: 2 }] })
    api.list.mockResolvedValue({ data: { list: [{ id: 11 }], total: 2, exception_total: 0 } })
    api.remove.mockResolvedValue({ res: 0 })
    api.move.mockResolvedValue({ res: 0 })
  })
  it('refreshes the server group totals as well as the list after successful deletion', async () => {
    const wrapper = mountPage()
    await flushPromises()
    wrapper.vm.$.setupState.handleBathDel()
    await api.confirm.mock.calls[0][0].onOk()
    await flushPromises()
    expect(api.remove).toHaveBeenCalledWith({ id: '11,12' })
    expect(api.groups).toHaveBeenCalledTimes(2)
    expect(api.list).toHaveBeenCalledTimes(2)
    expect(api.reset).toHaveBeenCalledTimes(2)
    wrapper.unmount()
  })
  it('moves selected QA to the valid ungrouped ID zero and refreshes both sides', async () => {
    const wrapper = mountPage()
    await flushPromises()
    wrapper.vm.$.setupState.handleOpenMoveModal()
    wrapper.vm.$.setupState.moveState.group_id = 0
    await flushPromises()
    wrapper.vm.$.setupState.handleMove()
    await flushPromises()
    expect(api.move).toHaveBeenCalledWith({ group_id: 0, ids: '11,12' })
    expect(api.groups).toHaveBeenCalledTimes(2)
    expect(api.list).toHaveBeenCalledTimes(2)
    wrapper.unmount()
  })
  it('returns to the last available page when deletion empties a page', async () => {
    api.list.mockResolvedValue({ data: { list: [{ id: 11 }], total: 201, exception_total: 0 } })
    const wrapper = mountPage()
    await flushPromises()
    wrapper.vm.$.setupState.onShowSizeChange(3, 100)
    await flushPromises()
    api.list.mockResolvedValueOnce({ data: { list: [], total: 200, exception_total: 0 } })
    api.list.mockResolvedValue({ data: { list: [{ id: 10 }], total: 200, exception_total: 0 } })
    wrapper.vm.$.setupState.handleBathDel()
    await api.confirm.mock.calls[0][0].onOk()
    await flushPromises()
    expect(api.list).toHaveBeenLastCalledWith(expect.objectContaining({ library_id: '501', page: 2 }))
    wrapper.unmount()
  })
  it('does not refresh or decrement counts after a failed deletion', async () => {
    const wrapper = mountPage()
    await flushPromises()
    api.remove.mockRejectedValueOnce(new Error('Delete failed'))
    wrapper.vm.$.setupState.handleBathDel()
    await expect(api.confirm.mock.calls[0][0].onOk()).rejects.toThrow('Delete failed')
    expect(api.groups).toHaveBeenCalledTimes(1)
    expect(api.list).toHaveBeenCalledTimes(1)
    wrapper.unmount()
  })
})
