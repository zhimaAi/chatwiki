import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'
import { flushPromises, mount, shallowMount } from '@vue/test-utils'
import { useTableScroll } from '@/views/library/library-details/components/use-table-scroll'
import SubsectionBox from '@/views/library/library-details/qa-knowledge-document/components/subsection-box.vue'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key) => key }) }))
vi.mock('@/api/library', () => ({
  deleteParagraph: vi.fn(),
  editParagraph: vi.fn(),
  getCategoryList: vi.fn(),
  updateParagraphCategory: vi.fn()
}))
vi.mock('@/components/upload-img/video-preview-modal.vue', () => ({ default: { template: '<div />' } }))

let availableHeight
let headerHeight
let observers
let wrappers

beforeEach(() => {
  availableHeight = 600
  headerHeight = 56
  observers = []
  wrappers = []
  vi.spyOn(Element.prototype, 'clientHeight', 'get').mockImplementation(() => availableHeight)
  vi.spyOn(Element.prototype, 'getBoundingClientRect').mockImplementation(function () {
    return { height: this.classList.contains('ant-table-header') ? headerHeight : 0 }
  })
  vi.stubGlobal('ResizeObserver', class {
    constructor(callback) {
      this.callback = callback
      this.observe = vi.fn()
      this.disconnect = vi.fn()
      observers.push(this)
    }
  })
})

afterEach(() => {
  wrappers.forEach((wrapper) => wrapper.unmount())
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

async function mountScroll(enabled = ref(true)) {
  let scrollY
  const wrapper = mount({
    setup() {
      const container = ref(null)
      scrollY = useTableScroll(container, () => enabled.value).scrollY
      return { container }
    },
    template: '<section ref="container"><div class="ant-table-header" /></section>'
  })
  wrappers.push(wrapper)
  await flushPromises()
  return { wrapper, scrollY, enabled }
}

describe('知识库表格内容区高度', () => {
  it('扣除表头高度，并响应可用空间和表头换行变化', async () => {
    const { wrapper, scrollY } = await mountScroll()
    expect(scrollY.value).toBe(544)
    expect(observers[0].observe).toHaveBeenCalledWith(wrapper.element)
    expect(observers[0].observe).toHaveBeenCalledWith(wrapper.find('.ant-table-header').element)

    availableHeight = 420
    headerHeight = 80
    observers[0].callback()
    expect(scrollY.value).toBe(340)
  })

  it('极小窗口保持内容区为正高度，空间恢复后重新计算', async () => {
    availableHeight = 30
    const { scrollY } = await mountScroll()
    expect(scrollY.value).toBe(1)
    availableHeight = 500
    observers[0].callback()
    expect(scrollY.value).toBe(444)
  })

  it('旧版不监听高度，切换模式释放监听且忽略旧回调', async () => {
    const { enabled, scrollY } = await mountScroll(ref(false))
    expect(observers).toHaveLength(0)
    enabled.value = true
    await flushPromises()
    expect(scrollY.value).toBe(544)
    enabled.value = false
    await flushPromises()
    expect(observers[0].disconnect).toHaveBeenCalledOnce()
    availableHeight = 200
    observers[0].callback()
    expect(scrollY.value).toBe(544)
  })

  it('卸载后释放监听，迟到的通知不再更新高度', async () => {
    const { wrapper, scrollY } = await mountScroll()
    wrapper.unmount()
    wrappers = []
    expect(observers[0].disconnect).toHaveBeenCalledOnce()
    availableHeight = 200
    observers[0].callback()
    expect(scrollY.value).toBe(544)
  })

  it('没有 ResizeObserver 时使用窗口事件并在卸载时清理', async () => {
    vi.stubGlobal('ResizeObserver', undefined)
    const removeListener = vi.spyOn(window, 'removeEventListener')
    const { wrapper, scrollY } = await mountScroll()
    availableHeight = 400
    window.dispatchEvent(new Event('resize'))
    expect(scrollY.value).toBe(344)
    wrapper.unmount()
    wrappers = []
    expect(removeListener).toHaveBeenCalledWith('resize', expect.any(Function))
    availableHeight = 200
    window.dispatchEvent(new Event('resize'))
    expect(scrollY.value).toBe(344)
  })
})

const tableStub = {
  name: 'ATable',
  props: ['dataSource', 'rowSelection', 'pagination', 'scroll'],
  emits: ['change'],
  template: '<div><div v-if="scroll.y" class="ant-table-header" /><slot /></div>'
}

async function mountTable(props = {}) {
  const wrapper = shallowMount(SubsectionBox, {
    props,
    global: { stubs: { ATable: tableStub, ATableColumn: true }, directives: { viewer: {} } }
  })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}

describe('问答表格工作区兼容', () => {
  it('旧版保留横向滚动、选择及排序事件', async () => {
    const wrapper = await mountTable()
    const table = wrapper.findComponent(tableStub)
    expect(table.props('scroll')).toEqual({ x: 1200 })
    expect(observers).toHaveLength(0)
    table.props('rowSelection').onChange([7])
    await flushPromises()
    expect(table.props('rowSelection').selectedRowKeys).toEqual([7])
    table.vm.$emit('change', {}, {}, { field: 'total_hits', order: 'descend' })
    expect(wrapper.emitted('handleSort')[0]).toEqual([{ sort_field: 'total_hits', sort_type: 'desc' }])
  })

  it('工作区启用表体滚动，恢复旧版时移除纵向限制', async () => {
    const wrapper = await mountTable({ fixedHeight: true })
    const table = wrapper.findComponent(tableStub)
    expect(wrapper.find('.fixed-table-host').exists()).toBe(true)
    expect(table.props('scroll')).toEqual({ x: 1200, y: 544 })
    availableHeight = 400
    observers[0].callback()
    await flushPromises()
    expect(table.props('scroll').y).toBe(344)
    await wrapper.setProps({ fixedHeight: false })
    expect(table.props('scroll')).toEqual({ x: 1200 })
    expect(observers[0].disconnect).toHaveBeenCalledOnce()
  })
})
