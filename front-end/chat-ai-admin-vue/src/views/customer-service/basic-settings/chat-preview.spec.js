import { beforeEach, describe, expect, it, vi } from 'vitest'
import { shallowMount, flushPromises } from '@vue/test-utils'
import { reactive, ref } from 'vue'
const mock = vi.hoisted(() => ({ session: null, permission: vi.fn(), error: vi.fn() }))
vi.mock('@/stores/modules/chat', () => ({ createRobotChatSession: () => mock.session }))
vi.mock('@/stores/modules/user', () => ({ useUserStore: () => ({ user_id: 5, user_name: 'Tester' }) }))
vi.mock('@/api/robot', () => ({ checkChatRequestPermission: mock.permission }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key, params) => params ? `${key}:${params.words}` : key }) }))
vi.mock('@/utils', () => ({ generateRandomId: () => 'image-text' }))
vi.mock('ant-design-vue', () => ({ message: { error: mock.error } }))
vi.mock('@/views/robot/robot-test/components/message-input.vue', () => ({ default: { name: 'MessageInput', props: ['value', 'fileList', 'loading'], emits: ['update:value', 'update:fileList', 'send', 'stop'], template: '<div />' } }))
vi.mock('@/views/robot/robot-test/components/message-list.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-test/components/variable-modal/index.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-test/components/prompt-log-alert.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-test/components/library-info-alert.vue', () => ({ default: { template: '<div />' } }))
import Preview from './chat-preview.vue'

function mountPreview(robotInfo = {}) {
  const context = { robotInfo: reactive({ id: 7, robot_key: 'seven', question_multiple_switch: 1, ...robotInfo }), savedRevision: ref(0), waitForSave: vi.fn().mockResolvedValue() }
  const wrapper = shallowMount(Preview, { global: { provide: { customerRobotContext: context } } })
  return { wrapper, context }
}
function input(wrapper) { return wrapper.findComponent({ name: 'MessageInput' }) }
describe('configuration chat preview', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mock.session = { robot: { robot_key: 'seven', openid: 5 }, chat_variables: {}, messageList: [], sendLock: false, dialogue_id: 0, createChat: vi.fn().mockResolvedValue(), sendMessage: vi.fn(), stopMessage: vi.fn(), $reset: vi.fn() }
    mock.permission.mockResolvedValue({ data: {} })
  })
  it('waits for saves before checking permissions and sending text with an attachment', async () => {
    const { wrapper, context } = mountPreview()
    await flushPromises()
    let finish
    context.waitForSave.mockImplementation(() => new Promise(resolve => { finish = resolve }))
    input(wrapper).vm.$emit('update:value', 'Hello')
    input(wrapper).vm.$emit('update:fileList', [{ uid: 'image', status: 'done', url: 'https://example.com/image.png' }])
    input(wrapper).vm.$emit('send')
    await flushPromises()
    expect(mock.permission).not.toHaveBeenCalled()
    finish()
    await flushPromises()
    const sent = mock.session.sendMessage.mock.calls[0][0]
    expect(JSON.parse(sent.message)).toEqual([{ type: 'text', uid: 'image-text', text: 'Hello' }, { type: 'image_url', uid: 'image', image_url: { url: 'https://example.com/image.png' } }])
    wrapper.unmount()
    expect(mock.session.$reset).toHaveBeenCalledTimes(2)
  })
  it('blocks sends after a failed save and drops pending sends on unmount', async () => {
    const { wrapper, context } = mountPreview()
    await flushPromises()
    context.waitForSave.mockRejectedValueOnce(new Error('save failed'))
    input(wrapper).vm.$emit('update:value', 'Hello')
    input(wrapper).vm.$emit('send')
    await flushPromises()
    expect(mock.error).toHaveBeenCalledWith('save failed')
    expect(mock.permission).not.toHaveBeenCalled()
    let finish
    context.waitForSave.mockImplementation(() => new Promise(resolve => { finish = resolve }))
    input(wrapper).vm.$emit('send')
    await flushPromises()
    wrapper.unmount()
    finish()
    await flushPromises()
    expect(mock.session.sendMessage).not.toHaveBeenCalled()
  })
  it('maps the avatar URL without overwriting session identity with form fields', async () => {
    const { wrapper, context } = mountPreview()
    await flushPromises()
    Object.assign(context.robotInfo, { robot_avatar: new File(['avatar'], 'avatar.png'), robot_avatar_url: '/saved-avatar.png', openid: 'wrong-user' })
    input(wrapper).vm.$emit('update:value', 'Hello')
    input(wrapper).vm.$emit('send')
    await flushPromises()
    expect(mock.session.robot.robot_avatar).toBe('/saved-avatar.png')
    expect(mock.session.robot.openid).toBe(5)
    context.robotInfo.robot_avatar = undefined
    context.robotInfo.robot_avatar_url = ''
    input(wrapper).vm.$emit('update:value', 'Again')
    input(wrapper).vm.$emit('send')
    await flushPromises()
    expect(mock.session.robot.robot_avatar).toBe('/saved-avatar.png')
    wrapper.unmount()
  })
  it('reuses the existing management welcome for legacy robots without changing the welcome API', async () => {
    mock.session.createChat.mockImplementation(async () => {
      mock.session.messageList.push({ isWelcome: true, menu_json: '' })
      return { data: { robot: { welcomes: '', lang_key: 'zh-CN' } } }
    })
    const { wrapper } = mountPreview({ welcomes: { content: 'Legacy welcome', question: [{ content: 'Legacy question' }] } })
    await flushPromises()
    expect(mock.session.messageList[0].menu_json).toEqual({ content: 'Legacy welcome', question: ['Legacy question'] })
    expect(mock.session.robot.welcomes).toEqual({ content: 'Legacy welcome', question: ['Legacy question'] })
    expect(mock.session.createChat).toHaveBeenCalledWith(expect.not.objectContaining({ welcomes: expect.anything() }))
    wrapper.unmount()
  })
  it.each([
    ['en-US', '{"content":"English welcome","question":[]}', 'English welcome'],
    ['ja', '', undefined],
    ['zh-CN', '{"content":"","question":[]}', ''],
  ])('preserves the existing %s language welcome, including intentionally empty values', async (language, welcomes, expected) => {
    mock.session.createChat.mockImplementation(async () => {
      mock.session.messageList.push({ isWelcome: true, menu_json: '' })
      return { data: { robot: { welcomes: '', lang_key: language } } }
    })
    const { wrapper } = mountPreview({ welcomes: { content: 'Chinese legacy welcome', question: [] }, multi_lang_configs: [{ lang_key: language, welcomes }] })
    await flushPromises()
    expect(mock.session.messageList[0].menu_json?.content).toBe(expected)
    wrapper.unmount()
  })
  it('keeps the welcome returned by the chat API instead of overwriting it with form data', async () => {
    mock.session.createChat.mockImplementation(async () => {
      mock.session.messageList.push({ isWelcome: true, menu_json: { content: 'Server welcome', question: [] } })
      return { data: { robot: { welcomes: '{"content":"Server welcome","question":[]}' } } }
    })
    const { wrapper } = mountPreview({ welcomes: { content: 'Form welcome', question: [] } })
    await flushPromises()
    expect(mock.session.messageList[0].menu_json.content).toBe('Server welcome')
    wrapper.unmount()
  })
  it('shows the full sensitive-word hint and preserves the draft and attachments without sending', async () => {
    const { wrapper } = mountPreview()
    await flushPromises()
    mock.permission.mockResolvedValue({ data: { words: ['99', 'blocked'] } })
    input(wrapper).vm.$emit('update:value', '99')
    input(wrapper).vm.$emit('update:fileList', [{ uid: 'image', status: 'done', url: '/image.png' }])
    input(wrapper).vm.$emit('send')
    await flushPromises()
    expect(mock.error).toHaveBeenCalledWith('msg_sensitive_words:99;blocked')
    expect(mock.session.sendMessage).not.toHaveBeenCalled()
    expect(input(wrapper).props('value')).toBe('99')
    expect(input(wrapper).props('fileList')).toHaveLength(1)
    wrapper.unmount()
  })
  it('refreshes variable rules after a save and waits for required variables before sending', async () => {
    const { wrapper, context } = mountPreview()
    await flushPromises()
    context.savedRevision.value++
    mock.session.createChat.mockImplementationOnce(async () => { mock.session.chat_variables.need_fill_variable = true })
    input(wrapper).vm.$emit('update:value', 'Hello')
    input(wrapper).vm.$emit('send')
    await flushPromises()
    expect(mock.session.createChat).toHaveBeenLastCalledWith(expect.objectContaining({ robot_key: 'seven', refreshOnly: true }))
    expect(mock.session.sendMessage).not.toHaveBeenCalled()
    expect(input(wrapper).props('value')).toBe('Hello')
    wrapper.unmount()
  })
})
