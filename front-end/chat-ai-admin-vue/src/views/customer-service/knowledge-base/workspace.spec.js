import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { Table, Dropdown, Menu, MenuItem, Checkbox, Tooltip, Popover } from 'ant-design-vue'
import { Modal } from 'ant-design-vue'
import SubsectionBox from '@/views/library/library-details/qa-knowledge-document/components/subsection-box.vue'
import ColumnWidthHandle from '@/views/library/library-details/qa-knowledge-document/components/column-width-handle.vue'
import QaDocument from '@/views/library/library-details/qa-knowledge-document/index.vue'

const api = vi.hoisted(() => ({ paragraphs: vi.fn(), groups: vi.fn(), remove: vi.fn(), preview: vi.fn() }))
vi.mock('@/api/library', () => ({
  getParagraphList: api.paragraphs, getLibraryGroup: api.groups, deleteParagraph: api.remove,
  reconstructCategoryVector: vi.fn(), deleteLibraryGroup: vi.fn(), setParagraphGroup: vi.fn(),
  createExportLibFileTask: vi.fn(), updateParagraphCategory: vi.fn()
}))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: { id: '7' } }) }))
vi.mock('@/router', () => ({ default: { push: vi.fn() } }))
vi.mock('@/components/cu-scroll/cu-scroll.vue', () => ({ default: { name: 'CuScroll', template: '<div><slot /></div>' } }))
vi.mock('@/components/upload-img/video-preview-modal.vue', () => ({ default: { setup: (_, { expose }) => expose({ show: api.preview }), template: '<div />' } }))
vi.mock('@/views/library/library-preview/components/edit-subsection.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/library/library-details/components/qa-upload-modal.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/library/library-details/components/metadata-manage-modal.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/library/library-details/qa-knowledge-document/components/add-group.vue', () => ({ default: { template: '<div />' } }))

let wrappers
const rows = [
  { id: 1, question: 'First question', answer: 'A long answer', images: ['one.png', 'two.png'], total_hits: 26, yesterday_hits: 0, today_hits: 0 },
  { id: 2, question: 'Second question', answer: 'Another answer', images: [], total_hits: 22, yesterday_hits: 1, today_hits: 0 }
]
beforeEach(() => {
  wrappers = []
  localStorage.clear()
  vi.clearAllMocks()
  vi.stubGlobal('ResizeObserver', class { observe() {} disconnect() {} })
  vi.stubGlobal('matchMedia', () => ({ matches: false, addListener() {}, removeListener() {} }))
  vi.spyOn(Element.prototype, 'clientWidth', 'get').mockReturnValue(932)
  vi.spyOn(Element.prototype, 'clientHeight', 'get').mockReturnValue(600)
  const getStyle = window.getComputedStyle
  vi.spyOn(window, 'getComputedStyle').mockImplementation(element => getStyle(element))
  api.paragraphs.mockResolvedValue({ data: { list: rows.map(row => ({ ...row })), total: 2, exception_total: 0 } })
  api.groups.mockResolvedValue({ data: [{ id: 0, group_name: 'Ungrouped', total: 1 }, { id: 8, group_name: 'Product', total: 1 }] })
  api.remove.mockResolvedValue({ res: 0 })
})
afterEach(() => {
  wrappers.forEach(wrapper => wrapper.unmount())
  document.body.innerHTML = ''
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

async function mountTable(props = {}) {
  const wrapper = mount(SubsectionBox, {
    attachTo: document.body,
    props: { compact: true, fixedHeight: true, paragraphLists: rows, total: 2, ...props },
    global: { components: { ATable: Table, ATableColumn: Table.Column, ADropdown: Dropdown, AMenu: Menu, AMenuItem: MenuItem, ACheckbox: Checkbox, ATooltip: Tooltip, APopover: Popover }, directives: { viewer: {} } }
  })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}

describe('新版问答表格交互', () => {
  it('选择、全选和清空使用现有选中 ID，并且不触发编辑', async () => {
    const wrapper = await mountTable()
    expect(wrapper.find('.subsection-box').classes()).not.toContain('has-selection')
    await wrapper.find('.ant-table-tbody input[type="checkbox"]').setValue(true)
    expect(wrapper.vm.state.selectedRowKeys).toEqual([1])
    expect(wrapper.find('.subsection-box').classes()).toContain('has-selection')
    await wrapper.find('.ant-table-thead input[type="checkbox"]').setValue(true)
    expect(wrapper.vm.state.selectedRowKeys).toEqual([1, 2])
    wrapper.vm.resetSelect()
    await flushPromises()
    expect(wrapper.find('.subsection-box').classes()).not.toContain('has-selection')
    expect(wrapper.findAll('.ant-table-tbody input:checked')).toHaveLength(0)
    expect(wrapper.emitted('openEditSubscription')).toBeUndefined()
  })

  it('列宽默认适应容器，调整统计列不排序，双击恢复默认值', async () => {
    const wrapper = await mountTable()
    const table = wrapper.findComponent(Table)
    expect(table.props('scroll').x).toBe(932)
    const handle = wrapper.findAllComponents(ColumnWidthHandle).find(item => item.props('width') === 72)
    await handle.trigger('keydown', { key: 'ArrowRight' })
    expect(handle.props('width')).toBe(80)
    expect(table.props('scroll').x).toBe(932)
    expect(wrapper.emitted('handleSort')).toBeUndefined()
    await handle.trigger('dblclick')
    expect(handle.props('width')).toBe(72)
    const questionHandle = wrapper.findAllComponents(ColumnWidthHandle).find(item => item.props('minWidth') === 240)
    await questionHandle.trigger('keydown', { key: 'ArrowRight' })
    expect(table.props('scroll').x).toBe(940)
    await questionHandle.trigger('dblclick')
    expect(table.props('scroll').x).toBe(932)
    await wrapper.find('th.qa-hit-cell').trigger('click')
    expect(wrapper.emitted('handleSort')[0]).toEqual([{ sort_field: 'total_hits', sort_type: 'asc' }])
  })

  it('动态元数据允许横向滚动，切换旧版恢复旧宽度和操作入口', async () => {
    const data = rows.map(row => ({ ...row, meta_list: [{ key: 'source', name: 'Source' }], meta_source: 'Uploaded' }))
    const wrapper = await mountTable({ paragraphLists: data })
    expect(wrapper.findComponent(Table).props('scroll').x).toBe(1092)
    expect(wrapper.findAll('.qa-more-button')).toHaveLength(2)
    await wrapper.setProps({ compact: false })
    expect(wrapper.findComponent(Table).props('scroll').x).toBe(1200)
    expect(wrapper.find('.qa-more-button').exists()).toBe(false)
    expect(wrapper.findAll('.list-label')).toHaveLength(4)
  })

  it('更多菜单编辑和删除操作使用当前记录，删除仍需确认', async () => {
    const wrapper = await mountTable()
    const confirm = vi.spyOn(Modal, 'confirm').mockImplementation(() => ({}))
    await wrapper.find('.qa-more-button').trigger('click')
    await vi.waitFor(() => expect(document.querySelector('.ant-dropdown-menu-item')).not.toBeNull())
    document.querySelector('.ant-dropdown-menu-item').click()
    expect(wrapper.emitted('openEditSubscription')[0]).toEqual([rows[0]])
    await wrapper.find('.qa-more-button').trigger('click')
    await vi.waitFor(() => expect(document.querySelectorAll('.ant-dropdown-menu-item').length).toBe(2))
    document.querySelectorAll('.ant-dropdown-menu-item')[1].click()
    expect(api.remove).not.toHaveBeenCalled()
    await confirm.mock.calls[0][0].onOk()
    expect(api.remove).toHaveBeenCalledWith({ id: 1 })
    expect(wrapper.emitted('handleDelParagraph')[0]).toEqual([1])
  })

  it('搜索高亮、异常提示和媒体展示保留真实内容与编辑入口', async () => {
    const record = { ...rows[0], status: 2, errmsg: 'Embedding failed', similar_questions: ['Related question'], images: ['photo.png', 'https://example.test/demo.mp4'] }
    const wrapper = await mountTable({ paragraphLists: [record], search: 'question' })
    expect(wrapper.find('.question-text .highlight').text()).toBe('question')
    expect(wrapper.find('.status-error').exists()).toBe(true)
    expect(wrapper.find('.image-item img').attributes('src')).toBe('photo.png')
    await wrapper.find('.video-item').trigger('click')
    expect(api.preview).toHaveBeenCalledWith('https://example.test/demo.mp4')
    await wrapper.find('.question-text').trigger('keydown', { key: 'Enter' })
    expect(wrapper.emitted('openEditSubscription')[0]).toEqual([record])
  })

  it('空列表正常渲染，加载状态仍由现有属性驱动', async () => {
    const wrapper = await mountTable({ paragraphLists: [], total: 0, isLoading: true })
    expect(wrapper.findComponent(Table).props('loading')).toBe(true)
    await wrapper.setProps({ isLoading: false })
    expect(wrapper.find('.ant-table-placeholder').exists()).toBe(true)
    expect(wrapper.find('.qa-more-button').exists()).toBe(false)
  })
})

describe('新版问答分组', () => {
  const popover = { name: 'APopover', props: ['open'], template: '<div><slot /><slot v-if="open" name="content" /></div>' }
  async function mountGroups(compact = true) {
    const wrapper = mount(QaDocument, {
      props: { compact, library_id: '101' },
      global: { stubs: {
        APopover: popover, ATooltip: { template: '<div><slot /></div>' }, CuScroll: false,
        AInput: { inheritAttrs: false, props: ['value', 'prefix', 'suffix'], template: '<input />' },
        SubsectionBox: { setup: (_, { expose }) => expose({ resetSelect() {} }), template: '<div />' }
      } }
    })
    wrappers.push(wrapper)
    await flushPromises()
    return wrapper
  }
  it('标题按钮折叠并记忆状态，浮层选择真实分组后保持折叠', async () => {
    const wrapper = await mountGroups()
    expect(wrapper.find('.group-content-box').attributes('style')).toContain('width: 200px')
    await wrapper.find('.collapse-group').trigger('click')
    expect(localStorage.getItem('customer_service_qa_document_group_hide_key')).toBe('1')
    const popup = wrapper.findComponent(popover)
    popup.vm.$emit('update:open', true)
    await flushPromises()
    expect(wrapper.findAll('.workspace-group-option').map(item => item.text())).toEqual(['all_groups2', 'Ungrouped1', 'Product1'])
    await wrapper.findAll('.workspace-group-option')[2].trigger('click')
    await flushPromises()
    expect(api.paragraphs).toHaveBeenLastCalledWith(expect.objectContaining({ library_id: '101', group_id: 8, page: 1 }))
    expect(wrapper.find('.current-group-name').text()).toBe('Product')
    expect(wrapper.find('.group-content-box').classes()).toContain('collapsed')
    expect(wrapper.find('.workspace-group-option').exists()).toBe(false)
    await wrapper.find('.expand-group').trigger('click')
    expect(wrapper.find('.group-content-box').classes()).not.toContain('collapsed')
  })
  it('旧版不出现新版工具栏，仍读取独立的折叠状态', async () => {
    localStorage.setItem('customer_service_qa_document_group_hide_key', '1')
    const wrapper = await mountGroups(false)
    expect(wrapper.find('.group-content-box').classes()).not.toContain('collapsed')
    expect(wrapper.find('.workspace-tools').exists()).toBe(false)
  })
  it('重新进入新版恢复分组宽度和折叠状态，折叠内容不参与键盘导航', async () => {
    localStorage.setItem('customer_service_qa_document_group_hide_key', '1')
    localStorage.setItem('customer_service_qa_document_group_box_width', '240')
    const wrapper = await mountGroups()
    expect(wrapper.find('.group-content-box').attributes('inert')).toBe('')
    await wrapper.find('.expand-group').trigger('click')
    expect(wrapper.find('.group-content-box').attributes('inert')).toBeUndefined()
    expect(wrapper.find('.group-content-box').attributes('style')).toContain('width: 240px')
  })
})

describe('列宽拖动生命周期', () => {
  it('拖动遵守最小宽度，指针取消和卸载恢复文档状态', async () => {
    const resize = vi.fn()
    const wrapper = mount(ColumnWidthHandle, { props: { width: 72, minWidth: 64, label: 'Resize', onResize: resize } })
    wrappers.push(wrapper)
    await wrapper.trigger('pointerdown', { button: 0, clientX: 100 })
    document.dispatchEvent(new MouseEvent('pointermove', { clientX: 50 }))
    expect(wrapper.emitted('resize')[0]).toEqual([64])
    document.dispatchEvent(new MouseEvent('pointercancel'))
    expect(document.body.style.cursor).toBe('')
    await wrapper.trigger('pointerdown', { button: 0, clientX: 100 })
    wrapper.unmount()
    document.dispatchEvent(new MouseEvent('pointermove', { clientX: 180 }))
    expect(resize).toHaveBeenCalledTimes(1)
    expect(document.body.style.userSelect).toBe('')
    wrappers.pop()
  })
})
