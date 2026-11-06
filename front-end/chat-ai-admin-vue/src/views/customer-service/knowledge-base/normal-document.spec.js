import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { Table, Dropdown, Menu, MenuItem, Tooltip } from 'ant-design-vue'
import WorkspaceDocument from '@/views/library/library-details/components/workspace-document-view.vue'
import DocumentStatus from '@/views/library/library-details/components/document-status.vue'
import ColumnWidthHandle from '@/views/library/library-details/qa-knowledge-document/components/column-width-handle.vue'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key, values) => values?.total !== undefined ? `${key}:${values.total}` : key }) }))
const rows = [
  { id: 11, file_name: '产品手册.pdf', file_ext: 'pdf', file_size: 2000, file_size_str: '2 KB', paragraph_count: 4, status: '2', total_hits: 12, today_hits: 1, yesterday_hits: 2, update_time: '2026-10-28 12:00', group_id: 7, doc_type: 1, cover_url: '/live-cover.png' },
  { id: 12, file_name: '售后说明.docx', file_ext: 'docx', file_size: 1000, file_size_str: '1 KB', paragraph_count: 2, status: '3', total_hits: 3, today_hits: 0, yesterday_hits: 0, update_time: '2026-10-27 12:00', group_id: 0, doc_type: 1, errmsg: '真实转换错误' }
]
const groups = [{ id: -1, group_name: '全部分组', total: 2 }, { id: 0, group_name: '未分组', total: 1 }, { id: 7, group_name: '产品资料', total: 1 }]
const popover = { props: ['open'], emits: ['update:open'], template: '<div><slot /><slot name="content" /></div>' }
const wrappers = []
const resizeObservers = []
async function render(props = {}) {
  const wrapper = mount(WorkspaceDocument, {
    attachTo: document.body,
    props: { files: rows, groups, libraryId: 101, pagination: { page: 1, size: 100, total: 2, sort_field: '', sort_type: '' }, updateFrequency: () => '--', ...props },
    global: {
      components: { ATable: Table, ADropdown: Dropdown, AMenu: Menu, AMenuItem: MenuItem, ATooltip: Tooltip },
      stubs: {
        'svg-icon': true, 'a-popover': popover,
        'a-input': { name: 'AInput', props: ['value'], emits: ['update:value', 'change'], template: '<input :value="value" @input="$emit(\'update:value\', $event.target.value)" @change="$emit(\'change\')" />' },
        'a-select': { name: 'ASelect', props: ['value', 'options'], emits: ['update:value', 'change'], template: '<select />' },
        'a-button': { template: '<button><slot /></button>' }, 'a-spin': { template: '<div><slot /></div>' }, 'a-empty': { template: '<div>empty</div>' },
        'a-pagination': { name: 'APagination', props: ['total', 'current', 'pageSize'], emits: ['change'], template: '<nav />' },
        draggable: { name: 'Draggable', props: ['modelValue', 'disabled'], emits: ['update:modelValue'], template: '<div><template v-for="element in modelValue" :key="element.id"><slot name="item" :element="element" /></template></div>' }
      }
    }
  })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}
const button = (wrapper, text) => wrapper.findAll('button').find(item => item.text() === text)
const columns = wrapper => wrapper.findComponent(Table).props('columns')
const totalWidth = wrapper => columns(wrapper).reduce((sum, column) => sum + column.width, 0)
beforeEach(() => {
  localStorage.clear()
  resizeObservers.length = 0
  vi.stubGlobal('ResizeObserver', class { constructor(callback) { this.callback = callback; resizeObservers.push(this) } observe() {} disconnect() {} })
  vi.stubGlobal('matchMedia', () => ({ matches: false, addListener() {}, removeListener() {} }))
  vi.spyOn(Element.prototype, 'clientWidth', 'get').mockReturnValue(948)
  vi.spyOn(Element.prototype, 'clientHeight', 'get').mockReturnValue(600)
  const getStyle = window.getComputedStyle
  vi.spyOn(window, 'getComputedStyle').mockImplementation(element => getStyle(element))
})
afterEach(() => { wrappers.splice(0).forEach(wrapper => wrapper.unmount()); document.body.innerHTML = ''; vi.restoreAllMocks(); vi.unstubAllGlobals() })

describe('normal knowledge workspace', () => {
  it('renders live document values and statistics with fixed table scrolling and a true total', async () => {
    const wrapper = await render()
    expect(wrapper.find('.ant-table-tbody').text()).toContain('产品手册.pdf')
    expect(wrapper.find('.ant-table-tbody').text()).toContain('2 KB')
    expect(wrapper.find('.ant-table-tbody').text()).toContain('document_manual_update')
    expect(wrapper.find('.workspace-document-pagination').text()).toContain('document_total:2')
    expect(wrapper.findComponent(Table).props('scroll').x).toBeUndefined()
    expect(totalWidth(wrapper)).toBe(948)
    expect(wrapper.findComponent(Table).props('scroll').y).toBeGreaterThan(0)
    expect(wrapper.findAllComponents(DocumentStatus)[1].attributes('title')).toBe('真实转换错误')
  })

  it('has no selection column, row checkbox or bulk selection controls', async () => {
    const wrapper = await render()
    expect(wrapper.findComponent(Table).props('rowSelection')).toBeUndefined()
    expect(wrapper.find('input[type="checkbox"]').exists()).toBe(false)
    expect(wrapper.find('.ant-table-selection-column').exists()).toBe(false)
    expect(wrapper.find('.document-selection-bar').exists()).toBe(false)
    expect(columns(wrapper)[0].key).toBe('file_name')
    await wrapper.find('.ant-table-tbody tr').trigger('mouseenter')
    expect(wrapper.find('input[type="checkbox"]').exists()).toBe(false)
  })

  it('resizes and resets columns independently of sorting and supports dynamic metadata', async () => {
    const wrapper = await render()
    const handle = wrapper.findAllComponents(ColumnWidthHandle)[0]
    await handle.trigger('keydown', { key: 'ArrowRight' })
    expect(totalWidth(wrapper)).toBe(948)
    expect(columns(wrapper)[0].width).toBe(188)
    expect(wrapper.emitted('sort')).toBeUndefined()
    expect(JSON.parse(localStorage.getItem('customer_service_document_columns_101')).file_name).toBe(188)
    await handle.trigger('dblclick')
    expect(columns(wrapper)[0].width).toBe(180)
    await handle.trigger('pointerdown', { button: 0, clientX: 100 })
    document.dispatchEvent(new MouseEvent('pointermove', { clientX: 124 }))
    document.dispatchEvent(new MouseEvent('pointerup'))
    await flushPromises()
    expect(totalWidth(wrapper)).toBe(948)
    expect(columns(wrapper)[0].width).toBe(204)
    expect(document.body.style.cursor).not.toBe('col-resize')
    expect(wrapper.emitted('sort')).toBeUndefined()
    await handle.trigger('dblclick')
    await wrapper.setProps({ extraColumns: [{ key: 'meta_level', dataIndex: 'meta_level', title: '会员等级' }], files: rows.map(row => ({ ...row, meta_level: '黄金' })) })
    expect(totalWidth(wrapper)).toBe(948)
    expect(wrapper.find('.ant-table-thead').text()).toContain('会员等级')
    expect(wrapper.find('.ant-table-tbody').text()).toContain('黄金')
  })

  it('routes supported sorting to the API and sorts size locally without changing the source rows', async () => {
    const wrapper = await render()
    await button(wrapper, 'col_total').trigger('click')
    expect(wrapper.emitted('sort')).toEqual([[{ field: 'total_hits', order: 'ascend' }]])
    await button(wrapper, 'col_size_segment').trigger('click')
    expect(wrapper.findAll('.document-file-name').map(item => item.text())).toEqual(['售后说明.docx', '产品手册.pdf'])
    expect(wrapper.emitted('sort')).toHaveLength(1)
    expect(rows[0].id).toBe(11)
  })

  it('filters groups, selects from a collapsed popover and preserves real names/counts', async () => {
    const wrapper = await render()
    await wrapper.find('.group-search').setValue('产品')
    expect(wrapper.findAll('.document-group-row')).toHaveLength(1)
    expect(wrapper.findComponent({ name: 'Draggable' }).props('disabled')).toBe(true)
    await wrapper.setProps({ collapsed: true })
    expect(wrapper.find('.document-groups').exists()).toBe(false)
    await wrapper.findAll('.collapsed-group-list button')[2].trigger('click')
    expect(wrapper.emitted('select-group').at(-1)).toEqual([groups[2]])
    await wrapper.find('.expand-groups').trigger('click')
    expect(wrapper.emitted('toggle-groups')).toHaveLength(1)
  })

  it('forwards the four add sources, metadata, search/status and pagination to existing handlers', async () => {
    const wrapper = await render()
    for (const item of wrapper.findAll('.document-add-grid button')) await item.trigger('click')
    expect(wrapper.emitted('add')).toEqual([[{ key: 1 }], [{ key: 2 }], [{ key: 3 }], [{ key: 4 }]])
    await wrapper.find('.metadata-button').trigger('click')
    expect(wrapper.emitted('metadata')).toHaveLength(1)
    await wrapper.find('.document-search').setValue('售后')
    expect(wrapper.emitted('update:fileName').at(-1)).toEqual(['售后'])
    const status = wrapper.findComponent({ name: 'ASelect' })
    status.vm.$emit('update:value', '3'); status.vm.$emit('change')
    expect(wrapper.emitted('update:status').at(-1)).toEqual(['3'])
    wrapper.findComponent({ name: 'APagination' }).vm.$emit('change', 2, 20)
    expect(wrapper.emitted('page-change').at(-1)).toEqual([2, 20])
  })

  it('uses real thumbnails, preserves group navigation and exposes hover menu actions in card mode', async () => {
    const wrapper = await render({ viewMode: 'folder' })
    expect(wrapper.findAll('.document-card')).toHaveLength(3)
    expect(wrapper.find('.file-cover').attributes('src')).toBe('/live-cover.png')
    expect(wrapper.findAll('.document-card-name').map(item => item.text())).toEqual(['产品资料', '产品手册.pdf', '售后说明.docx'])
    await wrapper.find('.group-card .document-card-open').trigger('click')
    expect(wrapper.emitted('select-group').at(-1)[0].id).toBe(7)
    await wrapper.find('.document-card:not(.group-card) .document-more').trigger('click')
    await flushPromises()
    const entries = [...document.querySelectorAll('.workspace-document-actions .ant-dropdown-menu-item')]
    expect(entries.map(item => item.textContent)).toContain('btn_base_data')
    entries.find(item => item.textContent.includes('btn_rename')).dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flushPromises()
    expect(wrapper.emitted('action').at(-1)).toEqual([{ key: 'rename', record: rows[0] }])
  })

  it('shows only real progress, handles an empty response and isolates column widths by library', async () => {
    const wrapper = await render({ viewMode: 'folder', files: [{ ...rows[0], status: '1', cover_url: '' }] })
    expect(wrapper.find('.on-card').text()).toBe('status_learning')
    await wrapper.setProps({ files: [{ ...rows[0], status: '1', ocr_pdf_index: 2, ocr_pdf_total: 5 }] })
    expect(wrapper.find('.on-card').text()).toBe('40% status_learning')
    await wrapper.setProps({ files: [], groups: [], viewMode: 'list' })
    expect(wrapper.findAll('.document-file-name')).toHaveLength(0)
    await wrapper.findAllComponents(ColumnWidthHandle)[0].trigger('keydown', { key: 'ArrowRight' })
    await wrapper.setProps({ libraryId: 202 })
    expect(wrapper.findAllComponents(ColumnWidthHandle)[0].props('width')).toBe(180)
  })

  it('fits the actual body width after viewport changes, metadata additions and oversized saved widths', async () => {
    let bodyWidth = 942
    vi.spyOn(Element.prototype, 'clientWidth', 'get').mockImplementation(function () { return this.classList.contains('ant-table-body') ? bodyWidth : bodyWidth + 6 })
    localStorage.setItem('customer_service_document_columns_101', JSON.stringify({ file_name: 6000 }))
    const wrapper = await render()
    expect(totalWidth(wrapper)).toBe(942)
    expect(columns(wrapper).at(-1).width).toBe(96)
    bodyWidth = 694
    resizeObservers.forEach(observer => observer.callback([]))
    await flushPromises()
    await wrapper.setProps({ extraColumns: Array.from({ length: 4 }, (_, index) => ({ key: `meta_${index}`, dataIndex: `meta_${index}`, title: `元数据${index}` })) })
    expect(totalWidth(wrapper)).toBe(694)
    expect(columns(wrapper)).toHaveLength(12)
    expect(columns(wrapper).every(column => column.width > 0)).toBe(true)
    expect(columns(wrapper).at(-1).width).toBe(96)
    expect(wrapper.findComponent(Table).props('scroll').x).toBeUndefined()
    expect(wrapper.findComponent(Table).props('scroll').y).toBeGreaterThan(0)
  })

  it('preserves graph states and offers generation only for waiting or failed documents', async () => {
    const wrapper = await render({ files: Array.from({ length: 6 }, (_, graph_status) => ({ ...rows[0], id: 20 + graph_status, graph_status, graph_err_count: 3 })), extraColumns: [{ key: 'graph_status', title: '知识图谱' }], loading: true })
    expect(wrapper.findComponent(Table).props('loading')).toBe(true)
    const generate = wrapper.findAll('button').filter(item => item.text() === 'btn_generate')
    expect(generate).toHaveLength(2)
    await generate[1].trigger('click')
    expect(wrapper.emitted('action').at(-1)[0]).toMatchObject({ key: 'generate-graph', record: { id: 23 } })
    await button(wrapper, 'btn_details').trigger('click')
    expect(wrapper.emitted('action').at(-1)[0]).toMatchObject({ key: 'graph-details', record: { id: 25 } })
  })
})
