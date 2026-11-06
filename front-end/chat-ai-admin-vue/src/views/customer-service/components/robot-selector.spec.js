import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { h, nextTick, reactive } from 'vue'
import { Modal, Popover } from 'ant-design-vue'
import RobotSelector from './robot-selector.vue'
import CustomerSidebar from './customer-sidebar.vue'

const state = vi.hoisted(() => ({ roles: [], permissions: {}, route: { meta: {} } }))
const api = vi.hoisted(() => ({ sort: vi.fn(), copy: vi.fn(), remove: vi.fn(), createRobot: vi.fn(), createAgent: vi.fn(), setGroup: vi.fn(), confirm: null, confirmation: null }))
vi.mock('@/api/robot', () => ({ moveRobotSort: api.sort, robotCopy: api.copy, deleteRobot: api.remove }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/utils/permission', () => ({ checkRole: roles => roles.some(role => state.roles.includes(role)), getRobotPermission: id => state.permissions[id] ?? 4 }))
vi.mock('vue-router', () => ({ useRoute: () => state.route }))
vi.mock('./customer-user-menu.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-list/components/add-robot-alert.vue', () => ({ default: { name: 'AddRobotAlert', setup: (_, { expose }) => expose({ open: api.createRobot, setGroupId: api.setGroup }), template: '<div />' } }))
vi.mock('@/views/clawbot/components/add-agent-alert.vue', () => ({ default: { name: 'AddAgentAlert', setup: (_, { expose }) => expose({ open: api.createAgent }), template: '<div />' } }))
vi.mock('vuedraggable', () => ({ default: {
  name: 'Draggable', props: ['modelValue', 'itemKey', 'handle', 'disabled'], emits: ['update:modelValue', 'start', 'end'],
  template: '<div><slot v-for="(option, index) in modelValue" name="item" :element="option" :index="index" /></div>'
} }))

const popover = {
  name: 'APopover', props: ['open', 'trigger', 'getPopupContainer'], emits: ['openChange'],
  methods: { toggle() { if (this.trigger?.includes('click')) this.$emit('openChange', !this.open) } },
  template: '<div><span class="popover-trigger" @click="toggle"><slot /></span><div v-if="open" class="popover-content"><slot name="content" /></div></div>'
}
const container = { template: '<div><slot /></div>' }
const globalOptions = { stubs: { APopover: popover, ASpin: container, AAlert: container, AEmpty: container, AButton: { template: '<button><slot /></button>' }, SvgIcon: true } }
const options = [
  { id: '1', name: 'Beauty', intro: 'Beauty description', robotKey: 'private-beauty-key', applicationType: 0, avatar: 'beauty.png' },
  { id: '2', name: 'Course', intro: 'Course description', robotKey: 'private-course-key', applicationType: 2, avatar: '' }
]
let wrappers
beforeEach(() => {
  wrappers = []
  vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout', 'setInterval', 'clearInterval'] })
  vi.clearAllMocks()
  state.roles = ['RobotCreate', 'ClawBotCreate']
  state.permissions = {}
  state.route = { meta: { customerServiceMenuKey: 'function-center', customerServiceGroupKey: 'ai-agent' } }
  api.sort.mockResolvedValue({ res: 0 })
  api.copy.mockResolvedValue({ res: 0 })
  api.remove.mockResolvedValue({ res: 0 })
  api.createRobot.mockResolvedValue()
  api.createAgent.mockResolvedValue()
  api.confirmation = { update: vi.fn(), destroy: vi.fn() }
  vi.spyOn(Modal, 'confirm').mockImplementation(config => { api.confirm = config; return api.confirmation })
})
afterEach(() => {
  wrappers.forEach(wrapper => wrapper.unmount())
  document.body.innerHTML = ''
  vi.clearAllTimers()
  vi.useRealTimers()
  vi.restoreAllMocks()
})

async function mountPicker(props = {}, realPopover = false) {
  const values = reactive({ open: false, options, value: '1', status: 'ready', ...props })
  const events = { change: vi.fn(), retry: vi.fn(), removed: vi.fn(), created: vi.fn() }
  const wrapper = mount({ render: () => h(RobotSelector, {
    ...values, 'onUpdate:open': value => { values.open = value }, onChange: events.change, onRetry: events.retry, onRemoved: events.removed, onCreated: events.created
  }, () => h('button', { class: 'test-trigger' }, 'Select app')) }, {
    attachTo: document.body,
    global: realPopover ? { ...globalOptions, components: { APopover: Popover }, stubs: { ...globalOptions.stubs, APopover: false } } : globalOptions
  })
  wrappers.push(wrapper)
  await wrapper.find('.test-trigger').trigger('click')
  await flushPromises()
  return { wrapper, picker: wrapper.findComponent(RobotSelector), values, events }
}

describe('顶部应用选择浮层', () => {
  it('真实 Popover 将列表挂载到浮层，详情操作区属于同一浮层容器', async () => {
    const { wrapper, values, events } = await mountPicker({}, true)
    await vi.advanceTimersByTimeAsync(250)
    expect(values.open).toBe(true)
    const panel = document.querySelector('#customer-service-app-picker')
    expect(panel).not.toBeNull()
    panel.querySelectorAll('.robot-row')[1].dispatchEvent(new MouseEvent('mouseenter'))
    await flushPromises()
    await vi.advanceTimersByTimeAsync(250)
    expect(panel.querySelector('.robot-preview h3').textContent).toBe('Course')
    expect(events.change).not.toHaveBeenCalled()
    await wrapper.find('.test-trigger').trigger('click')
    expect(values.open).toBe(false)
  })
  it('顶部头像和名称保留，按钮可重复开关浮层，选择使用实际应用 ID', async () => {
    const wrapper = mount(CustomerSidebar, { props: { currentRobot: { id: '1', name: 'Beauty', avatar: 'beauty.png' }, robotOptions: options, robotListStatus: 'ready' }, global: globalOptions })
    wrappers.push(wrapper)
    expect(wrapper.find('.brand-logo').attributes('src')).toBe('beauty.png')
    await wrapper.find('.app-brand').trigger('click')
    expect(wrapper.find('.app-brand').attributes('aria-expanded')).toBe('true')
    expect(wrapper.findAll('.robot-row')).toHaveLength(2)
    await wrapper.find('.app-brand').trigger('click')
    expect(wrapper.find('.robot-selector').exists()).toBe(false)
    await wrapper.find('.app-brand').trigger('click')
    await wrapper.findAll('.robot-row-select')[1].trigger('click')
    expect(wrapper.emitted('robot-change')[0]).toEqual(['2'])
    expect(wrapper.find('.app-brand').attributes('aria-expanded')).toBe('false')
    expect(wrapper.find('.brand-app-name').text()).toBe('Beauty')
  })

  it('悬停详情不切换应用，类型、简介和头像按当前行更新，密钥不进入 DOM', async () => {
    const { wrapper, events } = await mountPicker()
    await wrapper.findAll('.robot-row')[1].trigger('mouseenter')
    expect(wrapper.find('.robot-preview h3').text()).toBe('Course')
    expect(wrapper.find('.robot-type').text()).toBe('type_agent')
    expect(wrapper.find('.robot-description').text()).toBe('Course description')
    expect(wrapper.find('.preview-avatar').attributes('src')).toContain('agent.svg')
    expect(events.change).not.toHaveBeenCalled()
    expect(wrapper.html()).not.toContain('private-course-key')
    await wrapper.findAll('.robot-row')[0].trigger('focusin')
    expect(wrapper.find('.robot-type').text()).toBe('type_chat_robot')
    expect(wrapper.find('.preview-avatar').attributes('src')).toBe('beauty.png')
  })

  it('从列表移入详情保留卡片，移出后关闭；Esc 关闭浮层并恢复按钮焦点', async () => {
    const { wrapper, values } = await mountPicker()
    await wrapper.find('.robot-row').trigger('mouseenter')
    await wrapper.find('.robot-selector').trigger('mouseleave')
    await wrapper.find('.robot-preview').trigger('mouseenter')
    await vi.advanceTimersByTimeAsync(200)
    expect(wrapper.find('.robot-preview').exists()).toBe(true)
    await wrapper.find('.robot-preview').trigger('mouseleave')
    await vi.advanceTimersByTimeAsync(200)
    expect(wrapper.find('.robot-preview').exists()).toBe(false)
    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    await nextTick()
    expect(values.open).toBe(false)
    expect(document.activeElement).toBe(wrapper.find('.test-trigger').element)
  })

  it('外部关闭会清空预览，再次打开时不残留详情', async () => {
    const { wrapper } = await mountPicker()
    await wrapper.find('.robot-row').trigger('mouseenter')
    wrapper.findAllComponents(popover)[0].vm.$emit('openChange', false)
    await nextTick()
    expect(wrapper.find('.robot-selector').exists()).toBe(false)
    await wrapper.find('.test-trigger').trigger('click')
    expect(wrapper.find('.robot-preview').exists()).toBe(false)
  })

  it('无效应用禁止选择，加载、失败和空状态继续提供反馈与重试', async () => {
    const { wrapper, values, events } = await mountPicker({ options: [{ ...options[0], robotKey: '' }] })
    expect(wrapper.find('.robot-row-select').element.disabled).toBe(true)
    await wrapper.find('.robot-row-select').trigger('click')
    expect(events.change).not.toHaveBeenCalled()
    values.status = 'loading'
    await nextTick()
    expect(wrapper.find('[role="status"]').text()).toContain('loading_robots')
    values.status = 'error'
    await nextTick()
    await wrapper.find('.picker-feedback button').trigger('click')
    expect(events.retry).toHaveBeenCalledOnce()
    values.status = 'ready'
    values.options = []
    await nextTick()
    expect(wrapper.find('.picker-feedback').exists()).toBe(true)
    expect(wrapper.find('.robot-list').exists()).toBe(false)
  })

  it('新增入口按权限展示，创建表单打开前关闭选择浮层', async () => {
    state.roles = ['RobotCreate']
    const { wrapper, values } = await mountPicker()
    const submenu = wrapper.findAllComponents(popover).at(-1)
    submenu.vm.$emit('openChange', true)
    await nextTick()
    expect(wrapper.findAll('.picker-create-menu button').map(button => button.text())).toEqual(['add_robot'])
    await wrapper.find('.picker-create-menu button').trigger('click')
    await flushPromises()
    expect(values.open).toBe(false)
    expect(api.setGroup).toHaveBeenCalledWith('0')
    expect(api.createRobot).toHaveBeenCalledWith(0)
    expect(api.createAgent).not.toHaveBeenCalled()
  })

  it('智能体新增仍调用已有表单，成功事件继续交由上层切换应用', async () => {
    const { wrapper, values, events } = await mountPicker()
    wrapper.findAllComponents(popover).at(-1).vm.$emit('openChange', true)
    await nextTick()
    await wrapper.findAll('.picker-create-menu button')[1].trigger('click')
    await flushPromises()
    expect(api.createAgent).toHaveBeenCalledOnce()
    expect(values.open).toBe(false)
    wrapper.findComponent({ name: 'AddAgentAlert' }).vm.$emit('ok', { id: '3' })
    expect(events.created).toHaveBeenCalledWith({ id: '3' })
  })

  it('拖动只保存列表顺序，失败回滚且不改变当前应用', async () => {
    api.sort.mockRejectedValueOnce(new Error('network'))
    const { wrapper, events } = await mountPicker()
    const drag = wrapper.findComponent({ name: 'Draggable' })
    drag.vm.$emit('start')
    await nextTick()
    expect(drag.props('disabled')).toBe(false)
    drag.vm.$emit('update:modelValue', [options[1], options[0]])
    drag.vm.$emit('end', { oldIndex: 1, newIndex: 0 })
    await flushPromises()
    expect(api.sort).toHaveBeenCalledWith({ move_robot_id: '2', to_robot_id: 0 })
    expect(wrapper.findAll('.robot-name').map(name => name.text())).toEqual(['Beauty', 'Course'])
    expect(events.change).not.toHaveBeenCalled()
    drag.vm.$emit('start')
    drag.vm.$emit('update:modelValue', [options[1], options[0]])
    drag.vm.$emit('end', { oldIndex: 1, newIndex: 0 })
    await flushPromises()
    expect(events.retry).toHaveBeenCalledOnce()
  })

  it('复制和删除保留确认倒计时，并操作悬停应用而非当前应用', async () => {
    const { wrapper, events } = await mountPicker()
    await wrapper.findAll('.robot-row')[1].trigger('mouseenter')
    await wrapper.find('.robot-actions button').trigger('click')
    expect(api.copy).not.toHaveBeenCalled()
    expect(api.confirm.okButtonProps.disabled).toBe(true)
    await vi.advanceTimersByTimeAsync(3000)
    expect(api.confirmation.update).toHaveBeenLastCalledWith(expect.objectContaining({ okButtonProps: { disabled: false } }))
    await api.confirm.onOk()
    expect(api.copy).toHaveBeenCalledWith({ id: '2' })
    expect(events.retry).toHaveBeenCalledOnce()
    api.confirm.afterClose()
    await wrapper.findAll('.robot-row')[1].trigger('mouseenter')
    await wrapper.findAll('.robot-actions button')[1].trigger('click')
    expect(api.remove).not.toHaveBeenCalled()
    await api.confirm.onOk()
    expect(api.remove).toHaveBeenCalledWith({ id: '2' })
    expect(events.removed).toHaveBeenCalledWith('2')
    expect(events.change).not.toHaveBeenCalled()
  })

  it('没有创建或删除权限时隐藏对应入口，复制仍可使用', async () => {
    state.roles = []
    state.permissions = { '2': 1 }
    const { wrapper } = await mountPicker()
    expect(wrapper.find('.picker-add').exists()).toBe(false)
    await wrapper.findAll('.robot-row')[1].trigger('mouseenter')
    expect(wrapper.findAll('.robot-actions button').map(button => button.text())).toEqual(['copy'])
  })
})
