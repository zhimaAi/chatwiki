import { describe, expect, it, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/hooks/web/useUpload.js', () => ({ useUpload: () => ({ openFileDialog: vi.fn() }) }))
vi.mock('@/views/clawbot/chat/components/chat-model-select.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/clawbot/chat/components/file-toolbar.vue', () => ({ default: { template: '<div />' } }))
import MessageInput from './message-input.vue'

describe('chat input interaction', () => {
  it('opens required information only on a click, not on restored focus, and retains keyboard sending', async () => {
    const wrapper = shallowMount(MessageInput, {
      props: { value: 'draft' },
      global: { stubs: { ATextarea: {
        props: ['value'], emits: ['click', 'focus', 'keydown'],
        template: '<textarea @click="$emit(\'click\', $event)" @focus="$emit(\'focus\', $event)" @keydown="$emit(\'keydown\', $event)" />'
      } } }
    })
    const textarea = wrapper.find('textarea')
    await textarea.trigger('focus')
    expect(wrapper.emitted('inputClick')).toBeUndefined()
    await textarea.trigger('click')
    expect(wrapper.emitted('inputClick')).toHaveLength(1)
    await textarea.trigger('focus')
    expect(wrapper.emitted('inputClick')).toHaveLength(1)
    await textarea.trigger('keydown', { key: 'Enter' })
    expect(wrapper.emitted('send')).toEqual([['draft']])
    wrapper.unmount()
  })
})
