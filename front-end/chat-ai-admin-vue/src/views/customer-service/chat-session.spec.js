import { beforeEach, describe, expect, it, vi } from 'vitest'
const mocks = vi.hoisted(() => ({ welcome: vi.fn(), streams: [], connections: [] }))
vi.mock('@/api/chat', () => ({
  chatWelcome: mocks.welcome,
  sendAiMessage: () => { const stream = { abort: vi.fn() }; mocks.streams.push(stream); return stream },
  getDialogueList: vi.fn().mockResolvedValue({ data: [] }),
  getChatMessage: vi.fn(), questionGuide: vi.fn(), getSessionRecordList: vi.fn(), getSessionChannelList: vi.fn(), editVariables: vi.fn()
}))
vi.mock('@/api/robot/index', () => ({ editPrompt: vi.fn() }))
vi.mock('@/hooks/event/useEventBus', () => ({ useEventBus: () => ({ emit: vi.fn() }) }))
vi.mock('@/hooks/event/useIM', () => ({ useIM: isolated => { const connection = { isolated, connect: vi.fn().mockResolvedValue(), close: vi.fn(), on: vi.fn() }; mocks.connections.push(connection); return connection } }))
vi.mock('@/utils/index', () => ({ getUuid: () => Math.random().toString(), getOpenid: key => key, devLog: vi.fn(), extractVoiceInfo: vi.fn(), removeVoiceFormat: value => value }))
import { createRobotChatSession } from '@/stores/modules/chat'

const response = key => ({ data: { customer: {}, robot: { robot_key: key, welcomes: '{"content":"","question":[]}' }, message: [] } })
describe('isolated robot test sessions', () => {
  beforeEach(() => { mocks.welcome.mockReset(); mocks.streams.length = 0; mocks.connections.length = 0 })
  it.each([
    [{ content: '你好，欢迎咨询！', question: [] }, 1],
    [{ content: '', question: ['菜单问题'] }, 1],
    [{ content: '', question: [] }, 0]
  ])('preserves a valid welcome and respects an intentionally empty configuration: %j', async (menu, expectedCount) => {
    const session = createRobotChatSession({ isolatedIM: true })
    mocks.welcome.mockResolvedValue({ data: { ...response('same').data, message: { menu_json: JSON.stringify(menu) } } })
    await session.createChat({ robot_key: 'same', openid: 'user' })
    expect(session.messageList.value).toHaveLength(expectedCount)
    if (expectedCount) expect(session.messageList.value[0].menu_json).toEqual(menu)
    session.$reset()
  })
  it('ignores a welcome response after the session is reset', async () => {
    let resolve
    mocks.welcome.mockImplementation(() => new Promise(done => { resolve = done }))
    const session = createRobotChatSession({ isolatedIM: true })
    const pending = session.createChat({ robot_key: 'old', openid: 'user' })
    session.$reset()
    resolve(response('old'))
    await pending
    expect(session.robot.robot_key).toBe('')
    expect(mocks.connections[0].connect).not.toHaveBeenCalled()
  })
  it('stops old streams from replacing a new dialogue or unlocking another send', async () => {
    const session = createRobotChatSession({ isolatedIM: true })
    session.sendMessage({ message: 'first' })
    const old = mocks.streams[0]
    session.stopMessage()
    session.sendMessage({ message: 'second' })
    old.onMessage({ event: 'dialogue_id', data: 999 })
    old.onClose()
    expect(session.dialogue_id.value).toBe(0)
    expect(session.sendLock.value).toBe(true)
    expect(old.abort).toHaveBeenCalledOnce()
    session.$reset()
    expect(mocks.streams[1].abort).toHaveBeenCalledOnce()
  })
  it('keeps preview and normal chat state and IM lifetimes separate', () => {
    const first = createRobotChatSession({ isolatedIM: true })
    const second = createRobotChatSession({ isolatedIM: true })
    first.dialogue_id.value = 12
    first.$reset()
    expect(second.dialogue_id.value).toBe(0)
    expect(mocks.connections.every(item => item.isolated)).toBe(true)
    expect(mocks.connections[1].close).not.toHaveBeenCalled()
  })
  it('keeps preview variables out of the normal chat local storage', () => {
    localStorage.setItem('chat_prompt_variables_same', '["normal"]')
    const preview = createRobotChatSession({ isolatedIM: true })
    preview.robot.robot_key = 'same'
    preview.setPendingVariables(['preview'])
    preview.sendMessage({ message: 'question' })
    expect(localStorage.getItem('chat_prompt_variables_same')).toBe('["normal"]')
    preview.$reset()
    localStorage.removeItem('chat_prompt_variables_same')
  })
  it('refreshes saved configuration and variable rules without clearing the current conversation', async () => {
    const session = createRobotChatSession({ isolatedIM: true })
    session.dialogue_id.value = 99
    session.messageList.value = [{ content: 'existing message' }]
    mocks.welcome.mockResolvedValue({ data: { ...response('same').data, chat_variable: { need_fill_variable: true, wait_variables: [{ variable_key: 'new_variable' }] } } })
    await session.createChat({ robot_key: 'same', dialogue_id: 99, refreshOnly: true })
    expect(session.messageList.value).toEqual([{ content: 'existing message' }])
    expect(session.dialogue_id.value).toBe(99)
    expect(session.chat_variables.value.need_fill_variable).toBe(true)
    expect(mocks.welcome).toHaveBeenCalledWith(expect.objectContaining({ dialogue_id: 99 }), { cancelToken: false })
    expect(mocks.connections[0].connect).not.toHaveBeenCalled()
    session.$reset()
  })
})
