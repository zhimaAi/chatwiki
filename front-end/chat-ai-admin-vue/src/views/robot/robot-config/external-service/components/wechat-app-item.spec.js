import { describe, expect, it, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))

import ChannelItem from './wechat-app-item.vue'

const item = { id: 1, app_type: 'messenger', app_name: '样式测试', app_id: 'style-test' }
const mountItem = compact => shallowMount(ChannelItem, {
  props: { item, compact },
  global: {
    stubs: {
      'a-dropdown': { template: '<div><slot /><slot name="overlay" /></div>' },
      'a-menu': { template: '<div><slot /></div>' },
      'a-menu-item': { template: '<div class="menu-item"><slot name="icon" /><slot /></div>' }
    }
  }
})

describe('channel operation menu', () => {
  it('keeps compact delete text and icon on the same inherited menu color', async () => {
    const wrapper = mountItem(true)
    const actions = wrapper.findAll('.menu-item')
    expect(actions.map(action => action.text())).toEqual(['edit', 'delete'])
    expect(actions[1].find('a').element.style.color).toBe('inherit')
    expect(actions[1].find('delete-outlined-stub').exists()).toBe(true)
    await actions[0].trigger('click')
    await actions[1].trigger('click')
    expect(wrapper.emitted('edit')).toEqual([[item]])
    expect(wrapper.emitted('delete')).toEqual([[item]])
    wrapper.unmount()
  })

  it('preserves the original red deletion text outside compact mode', () => {
    const wrapper = mountItem(false)
    expect(wrapper.findAll('.menu-item')[1].find('a').element.style.color).toBe('rgb(251, 54, 63)')
    expect(wrapper.find('delete-outlined-stub').exists()).toBe(false)
    wrapper.unmount()
  })
})
