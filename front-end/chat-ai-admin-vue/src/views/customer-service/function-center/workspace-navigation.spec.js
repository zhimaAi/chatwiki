import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { computed, reactive } from 'vue'

const state = vi.hoisted(() => ({
  route: null, robot: null, push: vi.fn(), open: vi.fn(), setGroup: vi.fn(), setLayout: vi.fn()
}))
vi.mock('vue-router', () => ({ useRoute: () => state.route, useRouter: () => ({ push: state.push }) }))
vi.mock('pinia', () => ({ storeToRefs: store => ({ robotInfo: computed(() => store.robotInfo) }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => state.robot }))
vi.mock('@/stores/modules/permission', () => ({ usePermissionStore: () => ({ role_permission: ['RobotCreate'] }) }))
vi.mock('@/stores/modules/global', () => ({ useGlobalStore: () => ({ setHideLayoutTopAndBottom: state.setLayout }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/views/robot/robot-list/components/add-robot-alert.vue', () => ({ default: {
  setup: (_, { expose }) => { expose({ open: state.open, setGroupId: state.setGroup }) },
  template: '<div />'
} }))
vi.mock('@/views/robot/robot-config/components/robot-group.vue', () => ({ default: { template: '<div />' } }))

import TopHeader from '@/views/robot/robot-config/components/top-header.vue'
import PageSidebar from '@/views/workflow/components/page-sidebar.vue'

let wrappers
const container = { template: '<div><slot /></div>' }
const render = component => {
  const wrapper = shallowMount(component, { global: { stubs: {
    addRobotAlert: false,
    'a-popover': { ...container, name: 'APopover' },
    'a-dropdown': { ...container, name: 'ADropdown' }, 'a-tooltip': container,
    'cu-tooltip': container, 'svg-icon': true
  } } })
  wrappers.push(wrapper)
  return wrapper
}

beforeEach(() => {
  vi.useFakeTimers()
  vi.clearAllMocks()
  wrappers = []
  state.route = reactive({
    name: 'robotWorkflow', path: '/robot/config/workflow',
    query: { id: '900', robot_key: 'workflow-nine', source: 'customer-service' }
  })
  state.robot = {
    robotInfo: reactive({ id: 900, application_type: 1 }),
    robotList: reactive([{ id: 900, robot_name: 'Workflow nine', robot_avatar: 'flow.png', group_id: '5' }]),
    getRobotLists: vi.fn(), getGroupList: vi.fn()
  }
})
afterEach(() => {
  wrappers.forEach(wrapper => wrapper.unmount())
  vi.clearAllTimers()
  vi.useRealTimers()
})

describe('功能中心工作流导航', () => {
  it('显示工作流名称及编辑入口，移除主页、切换下拉和添加入口', async () => {
    const wrapper = render(TopHeader)
    expect(wrapper.find('.robot-name').text()).toBe('Workflow nine')
    expect(wrapper.find('.left-icon').exists()).toBe(false)
    expect(wrapper.find('.down-icon').exists()).toBe(false)
    expect(wrapper.findComponent({ name: 'APopover' }).exists()).toBe(false)
    expect(wrapper.findComponent({ name: 'ADropdown' }).exists()).toBe(false)
    await wrapper.find('.select-robot-block').trigger('click')
    expect(state.push).not.toHaveBeenCalled()
    await wrapper.find('.hover-block').trigger('click')
    expect(state.open).toHaveBeenCalledWith(1, true)
    expect(state.setGroup).toHaveBeenCalledWith('5')
  })

  it('旧版保留主页、切换下拉和添加入口', async () => {
    delete state.route.query.source
    const wrapper = render(TopHeader)
    expect(wrapper.find('.left-icon').exists()).toBe(true)
    expect(wrapper.find('.down-icon').exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'APopover' }).exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'ADropdown' }).exists()).toBe(true)
    await wrapper.find('.left-icon').trigger('click')
    expect(state.push).toHaveBeenCalledWith('/robot/list')
  })

  it('功能中心来源默认收起侧栏，菜单沿用工作流标识和来源', async () => {
    const wrapper = render(PageSidebar)
    expect(wrapper.find('.sidebar-wrapper').classes()).toContain('is-hide')
    expect(state.setLayout).not.toHaveBeenCalled()
    await wrapper.findAll('.sidebar-menu')[2].trigger('click')
    expect(state.push).toHaveBeenCalledWith({
      path: '/robot/config/basic-config',
      query: { id: '900', robot_key: 'workflow-nine', source: 'customer-service' }
    })
    await wrapper.find('.sidebar-handle').trigger('click')
    expect(wrapper.find('.sidebar-wrapper').classes()).not.toContain('is-hide')
    expect(state.setLayout).toHaveBeenLastCalledWith(false)
  })

  it('旧版侧栏默认折叠，展开行为保持原样', async () => {
    delete state.route.query.source
    const wrapper = render(PageSidebar)
    expect(wrapper.find('.sidebar-wrapper').classes()).toContain('is-hide')
    expect(state.setLayout).not.toHaveBeenCalled()
    await wrapper.find('.sidebar-handle').trigger('click')
    await flushPromises()
    expect(wrapper.find('.sidebar-wrapper').classes()).not.toContain('is-hide')
    expect(state.setLayout).toHaveBeenCalledWith(false)
  })
})
