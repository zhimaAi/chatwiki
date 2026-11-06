import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { reactive } from 'vue'
const state = vi.hoisted(() => ({ customer: null, chat: null }))
vi.mock('pinia', async () => {
  const { toRef } = await import('vue')
  return { storeToRefs: store => Object.fromEntries(Object.keys(store).map(key => [key, toRef(store, key)])) }
})
vi.mock('@/stores/modules/customer-service', () => ({ useCustomerServiceStore: () => state.customer }))
vi.mock('@/stores/modules/customer-service-chat', () => ({ useCustomerServiceChatStore: () => state.chat }))
vi.mock('@/stores/modules/user', () => ({ useUserStore: () => ({ user_id: 5 }) }))
vi.mock('@/stores/modules/locale', () => ({ useLocaleStore: () => ({ currentLocale: { lang: 'zh-CN' } }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('vue-router', () => ({ useRouter: () => ({ resolve: vi.fn() }) }))
vi.mock('ant-design-vue', () => ({ App: { useApp: () => ({ message: { error: vi.fn(), warning: vi.fn() } }) } }))
vi.mock('@/api/library', () => ({ getLibraryList: vi.fn().mockResolvedValue({ data: [] }) }))
vi.mock('@/api/explore/index.js', () => ({ getSpecifyAbilityConfig: vi.fn().mockResolvedValue({ data: {} }) }))
vi.mock('@/utils/index', () => ({ generateRandomId: () => 'draft' }))
vi.mock('@/views/clawbot/chat/components/chat-list.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/clawbot/chat/components/message-list.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/clawbot/chat/components/conversation-bar.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/clawbot/chat/components/scroll-bottom-btn.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/prompt-drawer.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/message-input.vue', () => ({ default: { name: 'MessageInput', emits: ['inputClick'], template: '<textarea @click="$emit(\'inputClick\')" />' } }))
vi.mock('./components/variable-modal.vue', () => ({ default: { name: 'VariableModal', props: ['open'], emits: ['close'], template: '<div />' } }))
import ChatPage from './index.vue'

describe('required information modal timing', () => {
  beforeEach(() => {
    state.customer = reactive({ robotStatus: 'ready', robot: { id: '7', robot_key: 'seven' } })
    state.chat = reactive({
      robot: { id: '7', robot_key: 'seven' }, messageList: [], myChatList: [], sensitiveWords: [], dialogue_id: 0, lastPushedUserMessageUid: '',
      chat_variables: { need_fill_variable: false, fill_variables: [] }, chatStatus: 'loading', chatLoading: false,
      sendLock: false, chatError: '', sendError: '', myChatListError: '', chatMessageError: '',
      myChatListLoading: false, myChatListHasMore: false, chatMessageLoading: false, variablesLoading: false,
      $reset: vi.fn(), initialize: vi.fn(async () => {
        await Promise.resolve()
        state.chat.chat_variables.need_fill_variable = true
        state.chat.chatStatus = 'ready'
      })
    })
  })
  it('does not open on entering chat; opens on input click and can be reopened after closing', async () => {
    const wrapper = shallowMount(ChatPage)
    await flushPromises()
    const modal = wrapper.findComponent({ name: 'VariableModal' })
    const input = wrapper.findComponent({ name: 'MessageInput' })
    expect(modal.props('open')).toBe(false)
    input.vm.$emit('inputClick')
    await flushPromises()
    expect(modal.props('open')).toBe(true)
    modal.vm.$emit('close')
    await flushPromises()
    expect(modal.props('open')).toBe(false)
    input.vm.$emit('inputClick')
    await flushPromises()
    expect(modal.props('open')).toBe(true)
    modal.vm.$emit('close')
    state.chat.chat_variables.need_fill_variable = false
    await flushPromises()
    input.vm.$emit('inputClick')
    await flushPromises()
    expect(modal.props('open')).toBe(false)
    wrapper.unmount()
  })
})
