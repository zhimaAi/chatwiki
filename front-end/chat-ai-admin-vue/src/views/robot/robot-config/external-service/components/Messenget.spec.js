import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'

const api = vi.hoisted(() => ({ list: vi.fn(), remove: vi.fn(), open: vi.fn() }))
vi.mock('@/api/robot', () => ({ getWechatAppList: api.list, deleteWechatApp: api.remove }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('./add-messenget-alert.vue', () => ({ default: {
  props: { compact: Boolean },
  emits: ['ok'],
  setup: (_, { expose }) => { expose({ open: api.open }) },
  template: '<div />'
} }))

import Messenger from './Messenget.vue'
import ChannelItem from './wechat-app-item.vue'
import AddChannel from './add-wechat-app.vue'
import ChannelModal from './add-messenget-alert.vue'

const sample = { id: 101, app_name: '拾光美妆旗舰店', app_id: 'style-test-101', app_type: 'messenger', app_avatar: '/upload/default/facebook-messenger.png' }
const mountPage = embedded => shallowMount(Messenger, {
  props: { embedded },
  global: {
    provide: { robotInfo: { robotInfo: { id: 7, wechat_ip: '127.0.0.1' } } },
    stubs: { AddMessagentAlert: false, 'a-button': { template: '<button><slot /></button>' } }
  }
})

beforeEach(() => {
  vi.clearAllMocks()
  api.list.mockResolvedValue({ data: [sample] })
})

describe('Messenger workspace presentation', () => {
  it('uses a creation toolbar and compact cards, reusing create and edit actions', async () => {
    const wrapper = mountPage(true)
    await flushPromises()
    expect(api.list).toHaveBeenCalledWith({ robot_id: 7, app_type: 'messenger', app_name: '' })
    expect(wrapper.findComponent(AddChannel).exists()).toBe(false)
    expect(wrapper.findComponent(ChannelItem).props('compact')).toBe(true)
    expect(wrapper.findComponent(ChannelModal).props('compact')).toBe(true)
    await wrapper.find('.messenger-toolbar button').trigger('click')
    expect(api.open).toHaveBeenCalledWith()
    wrapper.findComponent(ChannelItem).vm.$emit('edit', { ...sample })
    expect(api.open).toHaveBeenLastCalledWith({ ...sample, wechat_ip: '127.0.0.1' })
    api.list.mockResolvedValueOnce({ data: [] })
    wrapper.findComponent(ChannelModal).vm.$emit('ok')
    await flushPromises()
    expect(wrapper.findComponent(ChannelItem).exists()).toBe(false)
    expect(wrapper.find('.messenger-toolbar button').exists()).toBe(true)
    wrapper.unmount()
  })

  it('preserves the original creation tile and full cards outside the workspace', async () => {
    const wrapper = mountPage(false)
    await flushPromises()
    expect(wrapper.find('.messenger-toolbar').exists()).toBe(false)
    expect(wrapper.findComponent(AddChannel).exists()).toBe(true)
    expect(wrapper.findComponent(ChannelItem).props('compact')).toBe(false)
    expect(wrapper.findComponent(ChannelModal).props('compact')).toBe(false)
    wrapper.unmount()
  })
})
