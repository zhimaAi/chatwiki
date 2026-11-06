import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick, ref } from 'vue'
import BasicSettings from './index.vue'
import BasicConfig from '@/views/robot/robot-config/basic-config/index.vue'
import ChatPreview from './chat-preview.vue'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/views/robot/robot-config/basic-config/index.vue', () => ({ default: { props: ['variant'], template: '<div />' } }))
vi.mock('./chat-preview.vue', () => ({ default: { template: '<div />' } }))
vi.mock('../components/page-header.vue', () => ({ default: { props: ['title'], template: '<header><slot /></header>' } }))
const wrappers = []
const STORAGE_KEY = 'customer_service_basic_preview_width'
function render() {
  const wrapper = mount(BasicSettings, { global: { provide: { customerRobotContext: { saving: ref(false) } } } })
  wrappers.push(wrapper)
  return wrapper
}
async function pointer(target, type, clientX = 900, pointerId = 1, button = 0) {
  const event = new MouseEvent(type, { bubbles: true, cancelable: true, clientX, button })
  Object.defineProperty(event, 'pointerId', { value: pointerId })
  target.dispatchEvent(event)
  await nextTick()
}
const width = wrapper => wrapper.find('.preview').element.style.width
beforeEach(() => { localStorage.removeItem(STORAGE_KEY) })
afterEach(() => {
  wrappers.splice(0).forEach(wrapper => wrapper.unmount())
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
  vi.restoreAllMocks()
})

describe('basic settings preview resizing', () => {
  it('defaults to 400px and resizes the existing panes without remounting their contents', async () => {
    const wrapper = render()
    const form = wrapper.findComponent(BasicConfig).vm
    const preview = wrapper.findComponent(ChatPreview).vm
    expect(width(wrapper)).toBe('400px')
    expect(wrapper.findComponent(BasicConfig).props('variant')).toBe('workspace')
    const separator = wrapper.find('[role="separator"]')
    expect(separator.attributes()).toMatchObject({ 'aria-valuemin': '360', 'aria-valuemax': '460', 'aria-valuenow': '400', 'aria-controls': 'basic-settings-preview' })
    await pointer(separator.element, 'pointerdown')
    await pointer(document, 'pointermove', 870)
    expect(width(wrapper)).toBe('430px')
    expect(wrapper.findComponent(BasicConfig).vm).toBe(form)
    expect(wrapper.findComponent(ChatPreview).vm).toBe(preview)
    await pointer(document, 'pointerup')
  })

  it('clamps dragging to 360–460px, ignores other pointers and restores cursor/selection after release', async () => {
    const wrapper = render()
    document.body.style.cursor = 'crosshair'
    document.body.style.userSelect = 'text'
    await pointer(wrapper.find('[role="separator"]').element, 'pointerdown')
    await pointer(document, 'pointermove', 850, 2)
    expect(width(wrapper)).toBe('400px')
    await pointer(document, 'pointerup', 850, 2)
    await pointer(document, 'pointermove', 1000)
    expect(width(wrapper)).toBe('360px')
    await pointer(document, 'pointermove', 700)
    expect(width(wrapper)).toBe('460px')
    expect(document.body.style.cursor).toBe('col-resize')
    expect(document.body.style.userSelect).toBe('none')
    await pointer(document, 'pointerup')
    expect(document.body.style.cursor).toBe('crosshair')
    expect(document.body.style.userSelect).toBe('text')
    await pointer(document, 'pointermove', 850)
    expect(width(wrapper)).toBe('460px')
  })

  it('supports bounded keyboard resizing and double-click reset', async () => {
    const wrapper = render()
    const separator = wrapper.find('[role="separator"]')
    await separator.trigger('keydown', { key: 'ArrowLeft' })
    expect(width(wrapper)).toBe('408px')
    for (let i = 0; i < 15; i++) await separator.trigger('keydown', { key: 'ArrowRight' })
    expect(width(wrapper)).toBe('360px')
    for (let i = 0; i < 15; i++) await separator.trigger('keydown', { key: 'ArrowLeft' })
    expect(width(wrapper)).toBe('460px')
    await separator.trigger('dblclick')
    expect(width(wrapper)).toBe('400px')
  })

  it('cleans up a cancelled drag, lost window focus or unmount and ignores right-clicks', async () => {
    const wrapper = render()
    const separator = wrapper.find('[role="separator"]').element
    await pointer(separator, 'pointerdown', 900, 1, 2)
    expect(document.body.style.cursor).toBe('')
    await pointer(separator, 'pointerdown')
    await pointer(document, 'pointercancel')
    await pointer(document, 'pointermove', 850)
    expect(width(wrapper)).toBe('400px')
    expect(document.body.style.cursor).toBe('')
    await pointer(separator, 'pointerdown')
    window.dispatchEvent(new Event('blur'))
    await pointer(document, 'pointermove', 850)
    expect(width(wrapper)).toBe('400px')
    expect(document.body.style.userSelect).toBe('')
    await pointer(separator, 'pointerdown')
    wrapper.unmount()
    wrappers.splice(wrappers.indexOf(wrapper), 1)
    expect(document.body.style.cursor).toBe('')
    expect(document.body.style.userSelect).toBe('')
    await pointer(document, 'pointermove', 700)
    expect(document.body.style.cursor).toBe('')
  })

  it('restores the last dragged width on re-entry and persists keyboard changes and reset', async () => {
    const first = render()
    await pointer(first.find('[role="separator"]').element, 'pointerdown')
    await pointer(document, 'pointermove', 870)
    await pointer(document, 'pointerup')
    expect(localStorage.getItem(STORAGE_KEY)).toBe('430')
    first.unmount()
    wrappers.splice(wrappers.indexOf(first), 1)
    const restored = render()
    expect(width(restored)).toBe('430px')
    await restored.find('[role="separator"]').trigger('keydown', { key: 'ArrowLeft' })
    expect(localStorage.getItem(STORAGE_KEY)).toBe('438')
    await restored.find('[role="separator"]').trigger('dblclick')
    expect(localStorage.getItem(STORAGE_KEY)).toBe('400')
  })

  it.each([['420', '420px'], ['200', '360px'], ['900', '460px'], ['', '400px'], ['invalid', '400px'], ['Infinity', '400px']])('safely restores cached width %s as %s', (cached, expected) => {
    localStorage.setItem(STORAGE_KEY, cached)
    expect(width(render())).toBe(expected)
  })

  it('still permits resizing when browser storage is unavailable', async () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('storage disabled') })
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage disabled') })
    const wrapper = render()
    expect(width(wrapper)).toBe('400px')
    await pointer(wrapper.find('[role="separator"]').element, 'pointerdown')
    await pointer(document, 'pointermove', 870)
    await pointer(document, 'pointerup')
    expect(width(wrapper)).toBe('430px')
  })
})
