import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { reactive, toRef } from 'vue'

const api = vi.hoisted(() => ({ list: vi.fn(), keys: vi.fn(), robot: vi.fn() }))
let customerStore
let robotStore
vi.mock('pinia', () => ({ storeToRefs: store => ({ robotInfo: toRef(store, 'robotInfo') }) }))
vi.mock('@/stores/modules/customer-service', () => ({ useCustomerServiceStore: () => customerStore }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => robotStore }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key, values) => values ? `${key}:${values.count}` : key }) }))
vi.mock('@/api/robot/index', () => ({ getWechatAppList: api.list, listRobotApikey: api.keys }))
vi.mock('@/views/robot/api-key-manage/index.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/wechat-official-account.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/wechat-mini-program.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/wechat-customer-service.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/wecom-robot.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/feishu-robot.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/dingding-robot.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/Messenget.vue', () => ({ default: { props: { embedded: Boolean }, template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/telegram-robot.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/Line.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-config/external-service/components/whats-app.vue', () => ({ default: { template: '<div />' } }))

import Channels from './index.vue'
import Messenger from '@/views/robot/robot-config/external-service/components/Messenget.vue'
import wechatCustomerServiceIcon from '@/assets/img/channels/wechat-customer-service.png'
import wecomIcon from '@/assets/img/channels/wecom.svg'

const mountPage = () => shallowMount(Channels, {
  global: { stubs: {
    CustomerPageHeader: false,
    'a-button': { template: '<button><slot /></button>' },
    'a-empty': { template: '<div><slot /></div>' }
  } }
})
const names = wrapper => wrapper.findAll('.channel-name').map(item => item.text())
const messengerCard = wrapper => wrapper.findAll('.channel-card').find(item => item.text().includes('channel_messenger'))

beforeEach(() => {
  vi.clearAllMocks()
  customerStore = reactive({ robot: { id: 7, robot_key: 'seven' } })
  robotStore = reactive({ robotInfo: { id: 7 }, getRobot: api.robot })
  api.list.mockResolvedValue({ data: [] })
  api.keys.mockResolvedValue({ data: { list: [] } })
})

describe('other channel categories and status', () => {
  it('uses the updated independent WeChat service and WeCom robot icons', async () => {
    const wrapper = mountPage()
    await flushPromises()
    const cards = wrapper.findAll('.channel-card')
    const serviceIcon = cards.find(card => card.text().includes('channel_wechat_customer_service')).find('.channel-icon img')
    const robotIcon = cards.find(card => card.text().includes('channel_wecom_robot')).find('.channel-icon img')
    expect(serviceIcon.attributes('src')).toBe(wechatCustomerServiceIcon)
    expect(robotIcon.attributes('src')).toBe(wecomIcon)
    expect(serviceIcon.classes()).toContain('small')
    expect(robotIcon.classes()).toContain('small')
    expect(wrapper.findAll('.channel-desc')).toHaveLength(11)
    wrapper.unmount()
  })

  it('filters all eleven channels without refetching and restores the complete list', async () => {
    const wrapper = mountPage()
    await flushPromises()
    expect(names(wrapper)).toHaveLength(11)
    const filters = wrapper.findAll('.channel-filters button')
    await filters[1].trigger('click')
    expect(names(wrapper)).toEqual(['channel_official_account', 'channel_mini_program', 'channel_wechat_customer_service', 'channel_wecom_robot'])
    expect(filters[1].attributes('aria-pressed')).toBe('true')
    await filters[2].trigger('click')
    expect(names(wrapper)).toEqual(['channel_messenger', 'channel_telegram', 'channel_line', 'channel_whatsapp'])
    await filters[3].trigger('click')
    expect(names(wrapper)).toEqual(['channel_feishu_robot', 'channel_dingtalk_robot', 'channel_api_key'])
    await filters[0].trigger('click')
    expect(names(wrapper)).toHaveLength(11)
    expect(api.list).toHaveBeenCalledTimes(10)
    wrapper.unmount()
  })

  it('shows real connection counts, enabled API keys and unavailable channels', async () => {
    api.list.mockImplementation(({ app_type }) => app_type === 'line'
      ? Promise.reject(new Error('unavailable'))
      : Promise.resolve({ data: app_type === 'messenger' ? [{ id: 1 }, { id: 2 }] : [] }))
    api.keys.mockResolvedValue({ data: { list: [{ status: '1' }] } })
    const wrapper = mountPage()
    await flushPromises()
    expect(messengerCard(wrapper).find('.channel-badge').text()).toBe('status_connected:2')
    expect(messengerCard(wrapper).find('.channel-badge').classes()).toContain('green')
    expect(wrapper.findAll('.channel-badge').map(badge => badge.text())).toContain('status_unavailable')
    expect(wrapper.findAll('.channel-badge.green').map(badge => badge.text())).toContain('status_enabled')
    wrapper.unmount()
  })

  it('opens Messenger in the workspace and refreshes counts when returning to the selected category', async () => {
    const wrapper = mountPage()
    await flushPromises()
    await wrapper.findAll('.channel-filters button')[2].trigger('click')
    await messengerCard(wrapper).trigger('click')
    expect(wrapper.findComponent(Messenger).props('embedded')).toBe(true)
    expect(wrapper.find('.secondary-header h1').text()).toBe('channel_messenger')
    api.list.mockResolvedValue({ data: [{ id: 1 }] })
    await wrapper.find('.back-button').trigger('click')
    await flushPromises()
    expect(names(wrapper)).toHaveLength(4)
    expect(messengerCard(wrapper).find('.channel-badge').text()).toBe('status_connected:1')
    wrapper.unmount()
  })

  it('retries a failed detail load before showing the embedded channel', async () => {
    robotStore.robotInfo = null
    api.robot.mockRejectedValueOnce(new Error('failed')).mockImplementationOnce(async () => {
      robotStore.robotInfo = { id: 7 }
      return robotStore.robotInfo
    })
    const wrapper = mountPage()
    await flushPromises()
    await messengerCard(wrapper).trigger('click')
    await flushPromises()
    expect(wrapper.findComponent(Messenger).exists()).toBe(false)
    await wrapper.find('.secondary-body button').trigger('click')
    await flushPromises()
    expect(wrapper.findComponent(Messenger).props('embedded')).toBe(true)
    wrapper.unmount()
  })

  it('ignores old counts and failures after switching applications', async () => {
    let finishOld
    let failOld
    api.list.mockImplementation(({ robot_id, app_type }) => {
      if (robot_id === 7 && app_type === 'messenger') return new Promise(resolve => { finishOld = resolve })
      if (robot_id === 7 && app_type === 'line') return new Promise((resolve, reject) => { failOld = reject })
      return Promise.resolve({ data: [] })
    })
    const wrapper = mountPage()
    await flushPromises()
    customerStore.robot = { id: 8, robot_key: 'eight' }
    await flushPromises()
    finishOld({ data: [{ id: 1 }] })
    failOld(new Error('old failure'))
    await flushPromises()
    expect(messengerCard(wrapper).find('.channel-badge').text()).toBe('status_not_connected')
    expect(wrapper.findAll('.channel-badge').some(badge => badge.text() === 'status_unavailable')).toBe(false)
    expect(api.keys).toHaveBeenCalledTimes(1)
    expect(api.keys).toHaveBeenCalledWith({ robot_key: 'eight' })
    wrapper.unmount()
  })
})
