import { flushPromises, shallowMount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key) => key }) }))
vi.mock('@/api/app', () => ({ uploadFile: vi.fn() }))

import FloatIconSetting from './float-icon-setting.vue'

const form = {
  floatBtn: {
    displayType: 1,
    buttonText: '在线咨询',
    buttonIcon: '',
    bottomMargin: 50,
    rightMargin: 50,
    showUnreadCount: 0,
    showNewMessageTip: 0
  }
}

describe('float icon workspace', () => {
  it('applies a display change without a save button', async () => {
    const wrapper = shallowMount(FloatIconSetting, {
      props: { appearance: 'workspace', form },
      global: {
        stubs: {
          CardBox: { template: '<div><slot /></div>' },
          CuUpload: { template: '<div><slot /></div>' },
          'a-input': { template: '<input />' },
          'a-input-number': { template: '<input />' },
          'a-checkbox': { template: '<label><slot /></label>' }
        }
      }
    })
    await flushPromises()
    expect(wrapper.find('button').exists()).toBe(false)
    expect(wrapper.emitted('save')).toBeUndefined()

    await wrapper.findAll('.display-type')[1].trigger('click')
    const saved = wrapper.emitted('save')
    expect(saved.at(-1)[0].displayType).toBe(2)
  })
})
