import { describe, expect, it, vi, beforeEach } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { reactive } from 'vue'
import dayjs from 'dayjs'

const api = vi.hoisted(() => ({
  getRecordList: vi.fn(),
  getChannelList: vi.fn(),
  createMsg: vi.fn(),
  onGetChatMessage: vi.fn()
}))

vi.mock('pinia', async () => {
  const { ref } = await import('vue')
  return { storeToRefs: () => ({ messageList: ref([]) }) }
})
vi.mock('@/stores/modules/chat', () => ({
  useChatStore: () => reactive({
    messageList: [],
    getRecordList: api.getRecordList,
    getChannelList: api.getChannelList,
    createMsg: api.createMsg,
    onGetChatMessage: api.onGetChatMessage
  })
}))
vi.mock('@/hooks/web/useI18n', () => ({
  useI18n: () => ({ t: key => key })
}))
vi.mock('./design/components/message-list.vue', () => ({
  default: { name: 'DesignMessageList', template: '<div />' }
}))
vi.mock('@/views/robot/robot-test/components/prompt-log-alert.vue', () => ({
  default: { name: 'PromptLogAlert', template: '<div />' }
}))

import SessionRecordDesign from './design/index.vue'

const DateRangePanelStub = { name: 'DateRangePanel', template: '<div />' }
const PopoverStub = {
  name: 'APopover',
  template: '<div><slot /><slot name="content" /></div>'
}

const mountPage = () => shallowMount(SessionRecordDesign, {
  props: {
    robot: { id: 7, robot_key: 'customer-seven' }
  },
  global: {
    stubs: {
      APopover: PopoverStub,
      DateRangePanel: DateRangePanelStub,
      ATooltip: { template: '<div><slot /></div>' },
      AInput: { template: '<input />' },
      AButton: { template: '<button><slot /></button>' },
      AAvatar: { template: '<span><slot /></span>' },
      AEmpty: { template: '<div />' },
      ASpin: { template: '<div />' }
    }
  }
})

describe('session record list loading', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.getChannelList.mockResolvedValue({ data: [] })
    api.createMsg.mockResolvedValue()
    api.onGetChatMessage.mockResolvedValue(null)
  })

  it('lets a date filter reset supersede an in-flight list request', async () => {
    let resolveFirst
    api.getRecordList
      .mockImplementationOnce(() => new Promise((resolve) => { resolveFirst = resolve }))
      .mockResolvedValue({ data: { list: [], has_more: false } })

    const wrapper = mountPage()
    await flushPromises()
    expect(api.getRecordList).toHaveBeenCalledTimes(1)

    wrapper.findComponent(DateRangePanelStub).vm.$emit('confirm', [
      dayjs('2026-10-01'),
      dayjs('2026-10-02')
    ])
    await flushPromises()

    expect(api.getRecordList).toHaveBeenCalledTimes(2)
    const params = api.getRecordList.mock.calls[1][0]
    expect(params.start_time).toBe(dayjs('2026-10-01').startOf('day').unix())
    expect(params.end_time).toBe(dayjs('2026-10-02').endOf('day').unix())

    resolveFirst({ data: { list: [], has_more: false } })
    await flushPromises()
    wrapper.unmount()
  })
  it('keeps historical and zero-count channels visible and uses full server counts, not the loaded page', async () => {
    api.getChannelList.mockResolvedValue({ data: { total: 28, list: [
      { app_type: 'yun_h5', app_name: 'WebAPP', count: 0 },
      { app_type: 'feishu_robot', app_id: 'unbound', app_name: 'Historical Feishu', count: 25 },
      { app_type: 'feishu_robot', app_id: 'current', app_name: 'Current Feishu', count: 3 }
    ] } })
    api.getRecordList.mockResolvedValue({ data: { list: [{ session_id: 1, app_type: 'feishu_robot', app_id: 'unbound', is_transfer: '0' }], has_more: true } })
    const wrapper = mountPage()
    await flushPromises()
    const options = wrapper.findAll('.channel-option')
    expect(options.map(option => option.text())).toEqual(['all_channels28', 'WebAPP0', 'Historical Feishu25', 'Current Feishu3'])
    expect(wrapper.find('.chip-text').text()).toBe('Historical Feishu')
    await options[2].trigger('click')
    await flushPromises()
    expect(api.getRecordList).toHaveBeenLastCalledWith(expect.objectContaining({ app_type: 'feishu_robot', app_id: 'unbound' }))
    expect(wrapper.findAll('.channel-option')[2].text()).toBe('Historical Feishu25')
    wrapper.unmount()
  })
  it.each([0, '0', false, 'false', 1, '1', true])('never claims a handoff for an unsupported capability, including state %s', async state => {
    api.getRecordList.mockResolvedValue({ data: { list: [{ session_id: 1, app_type: 'yun_h5', is_transfer: state }], has_more: false } })
    const wrapper = mountPage()
    await flushPromises()
    expect(wrapper.find('.chip.transfer').exists()).toBe(false)
    expect(wrapper.find('.transfer-tag').exists()).toBe(false)
    wrapper.unmount()
  })
})
