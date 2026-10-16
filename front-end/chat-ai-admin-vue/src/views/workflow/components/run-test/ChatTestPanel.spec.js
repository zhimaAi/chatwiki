import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { callWorkFlowDialog } from '@/api/robot/index'
import ChatTestPanel from './ChatTestPanel.vue'

const mocks = vi.hoisted(() => ({ error: vi.fn(), warning: vi.fn(), id: 0 }))

vi.mock('ant-design-vue', () => ({ message: { error: mocks.error, warning: mocks.warning } }))
vi.mock('vue-router', () => ({ useRoute: () => ({ query: { robot_key: 'workflow-test' } }) }))
vi.mock('@/utils/index', () => ({ generateRandomId: () => `test-id-${++mocks.id}` }))
vi.mock('@/api/robot/index', () => ({
  callWorkFlowDialog: vi.fn((data) => ({ data, abort: vi.fn() }))
}))
vi.mock('../util', () => ({ getImageUrl: () => '/node.png' }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key) => key }) }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => ({ robotInfo: {} }) }))
vi.mock('@/constants/index', () => ({ DEFAULT_USER_AVATAR: '/user.png', DEFAULT_WORKFLOW_AVATAR: '/workflow.png' }))
vi.mock('@/views/robot/robot-test/components/message-input.vue', () => ({
  default: {
    name: 'MessageInput',
    props: ['value', 'fileList', 'loading'],
    emits: ['update:value', 'update:fileList', 'send'],
    render: () => null
  }
}))
vi.mock('@/views/robot/robot-test/components/message-list.vue', () => ({
  default: {
    name: 'MessageList',
    props: ['messages'],
    emits: ['clickMsgMeun', 'scroll'],
    methods: { scrollToBottom: vi.fn() },
    render: () => null
  }
}))
vi.mock('./NodeRunLogs.vue', () => ({
  default: {
    name: 'NodeRunLogs',
    props: ['resultList', 'currentNodeKey', 'loading'],
    render: () => null
  }
}))

const dialogueId = '90071992547409931'
const sessionId = '90071992547409932'
let wrapper

const input = () => wrapper.findComponent({ name: 'MessageInput' })
const messages = () => wrapper.findComponent({ name: 'MessageList' }).props('messages')
const logs = () => wrapper.findComponent({ name: 'NodeRunLogs' })
const context = (stream) => JSON.parse(stream.data.context_list)
const event = (stream, name, data) => stream.onMessage({
  event: name,
  data: typeof data === 'string' ? data : JSON.stringify(data)
})

const createPanel = (multimodal = false) => {
  wrapper = mount(ChatTestPanel, {
    props: {
      start_node_params: {
        trigger_list: [{ trigger_type: 1, chat_config: { question_multiple_switch: multimodal } }]
      }
    }
  })
  wrapper.vm.open()
}

const identify = (stream) => {
  event(stream, 'dialogue_id', dialogueId)
  event(stream, 'session_id', sessionId)
  return stream
}

const send = async (text, files = []) => {
  input().vm.$emit('update:value', text)
  input().vm.$emit('update:fileList', files)
  await nextTick()
  input().vm.$emit('send')
  await nextTick()
  return identify(callWorkFlowDialog.mock.results.at(-1).value)
}

const complete = async (stream, answer = '', nodeType = 7, finish = { is_paused: false }) => {
  event(stream, 'workflow_node_log', { node_key: 'node', node_type: nodeType, start_time: 1 })
  if (answer) event(stream, 'sending', answer)
  event(stream, 'finish', finish)
  stream.onClose()
  await nextTick()
}

beforeEach(() => {
  vi.clearAllMocks()
  mocks.id = 0
  vi.stubGlobal('requestAnimationFrame', vi.fn(() => 1))
  vi.stubGlobal('cancelAnimationFrame', vi.fn())
})

afterEach(() => {
  wrapper?.unmount()
  vi.unstubAllGlobals()
})

describe('workflow chat multi-turn history', () => {
  it('retains three rounds and string IDs while resetting logs for each new run', async () => {
    createPanel()
    const first = await send('第一问')
    expect(first.data).toMatchObject({ robot_key: 'workflow-test', dialogue_id: '', session_id: '', context_list: '[]' })
    await complete(first, '第一答')

    const second = await send('第二问')
    expect(second.data).toMatchObject({ dialogue_id: dialogueId, session_id: sessionId, openid: first.data.openid })
    expect(context(second)).toEqual([{ question: '第一问', answer: '第一答' }])
    expect(logs().props('resultList')).toEqual([])
    expect(logs().props('currentNodeKey')).toBe('')
    await complete(second, '第二答')

    const third = await send('第三问')
    expect(context(third)).toEqual([
      { question: '第一问', answer: '第一答' },
      { question: '第二问', answer: '第二答' }
    ])
    expect(third.data.openid).toBe(first.data.openid)
    expect(messages()).toHaveLength(6)
    expect(messages()[4].content).toBe('第三问')
    await complete(third, '第三答')
    expect(messages()[5]).toMatchObject({ content: '第三答', loading: false })
  })

  it('resumes a paused menu with plain text and keeps its logs and question in history', async () => {
    createPanel(true)
    const first = await send('开始')
    event(first, 'workflow_node_log', {
      node_key: 'menu', node_type: 43,
      output: { special: { reply_content_list: [{ smart_menu: { menu_description: '请选择服务' } }] } }
    })
    event(first, 'finish', { is_paused: true })
    first.onClose()
    await nextTick()
    wrapper.findComponent({ name: 'MessageList' }).vm.$emit('clickMsgMeun', '售后服务')
    await nextTick()
    const second = callWorkFlowDialog.mock.results.at(-1).value
    expect(second.data.question).toBe('售后服务')
    expect(second.data.dialogue_id).toBe(dialogueId)
    expect(context(second)).toEqual([{ question: '开始', answer: '请选择服务' }])
    expect(logs().props('resultList')).toHaveLength(1)
    expect(logs().props('currentNodeKey')).toBe('menu')
    expect(messages()).toHaveLength(4)
  })

  it.each([
    [true, 7, true],
    [false, 43, false],
    [undefined, 43, true],
    [undefined, 7, false]
  ])('uses pause flag %s with last node %s (preserve logs: %s)', async (paused, nodeType, preserve) => {
    createPanel()
    const first = await send('问题')
    await complete(first, '回答', nodeType, paused === undefined ? '[DONE]' : { is_paused: paused })
    await send('继续')
    expect(logs().props('resultList')).toHaveLength(preserve ? 1 : 0)
  })

  it('extracts multimodal text and image-only questions', async () => {
    createPanel(true)
    const image = { uid: 'image', status: 'done', url: '/image.png' }
    const first = await send('图片问题', [image])
    expect(JSON.parse(first.data.question)).toHaveLength(2)
    await complete(first, '文本回答')
    const second = await send('', [image])
    expect(context(second)).toEqual([{ question: '图片问题', answer: '文本回答' }])
    await complete(second, '图片回答')
    const third = await send('继续')
    expect(context(third)[1]).toEqual({ question: '[图片]', answer: '图片回答' })
  })

  it.each([
    ['[{"type":"input_audio","input_audio":{}}]', '[语音]'],
    ['[{"type":"video_url","video_url":{}}]', '[视频]'],
    ['[{"type":"text","text":"首段"},{"type":"text","text":"次段"}]', '首段'],
    ['{"question":"普通 JSON"}', '{"question":"普通 JSON"}'],
    ['[1,2,3]', '[1,2,3]'],
    ['[{"type":"text","text":123}]', '[{"type":"text","text":123}]'],
    ['[不完整 JSON', '[不完整 JSON']
  ])('extracts question %s as %s', async (question, expected) => {
    createPanel()
    const first = await send(question)
    await complete(first, '回答')
    const second = await send('继续')
    expect(context(second)).toEqual([{ question: expected, answer: '回答' }])
  })

  it.each([
    ['流式回答', [{ description: '描述' }], '流式回答'],
    ['', [{ description: '描述', title: '标题' }], '描述'],
    ['', [{ title: '标题', smart_menu: { menu_description: '菜单' } }], '标题'],
    ['', [{}, { smart_menu: { menu_description: '菜单' } }], '菜单']
  ])('extracts answer content %s and structured replies %j', async (text, replies, expected) => {
    createPanel()
    const first = await send('问题')
    event(first, 'reply_content_list', replies)
    await complete(first, text)
    const second = await send('继续')
    expect(context(second)).toEqual([{ question: '问题', answer: expected }])
  })

  it('omits empty questions, empty answers and unfinished replies without mismatching pairs', async () => {
    createPanel()
    const first = await send('   ')
    await complete(first, '空问题的回答')
    const second = await send('无回答的问题')
    await complete(second)
    const third = await send('未完成的问题')
    await complete(third, '部分回答')
    messages().at(-1).loading = true
    const fourth = await send('有效问题')
    expect(context(fourth)).toEqual([])
    await complete(fourth, '有效回答')
    const fifth = await send('继续')
    expect(context(fifth)).toEqual([{ question: '有效问题', answer: '有效回答' }])
  })

  it('sends the most recent 50 pairs while keeping every displayed message', async () => {
    createPanel()
    for (let i = 0; i < 52; i++) {
      const stream = await send(`问题${i}`)
      await complete(stream, `回答${i}`)
    }
    const next = await send('下一问')
    expect(context(next)).toHaveLength(50)
    expect(context(next)[0]).toEqual({ question: '问题2', answer: '回答2' })
    expect(context(next).at(-1)).toEqual({ question: '问题51', answer: '回答51' })
    expect(messages()).toHaveLength(106)
  })

  it.each(['transport', 'business'])('starts a new session after %s errors and ignores late callbacks', async (kind) => {
    createPanel()
    const first = await send('第一问')
    await complete(first, '第一答')
    const failed = await send('失败问题')
    event(failed, 'sending', '失败前的内容')
    event(failed, 'workflow_node_log', { node_key: 'question', node_type: 43 })
    if (kind === 'transport') failed.onError(new Error('connection lost'))
    else event(failed, 'error', '业务失败')
    event(failed, 'finish', { is_paused: true })
    failed.onClose()
    await nextTick()
    expect(messages()).toHaveLength(4)
    expect(messages().at(-1)).toMatchObject({ content: '失败前的内容', loading: false })
    expect(input().props('loading')).toBe(false)
    expect(mocks.error).toHaveBeenCalledTimes(1)
    expect(failed.abort).toHaveBeenCalled()

    const retry = await send('重试')
    expect(retry.data).toMatchObject({ dialogue_id: '', session_id: '', context_list: '[]' })
    expect(retry.data.openid).not.toBe(first.data.openid)
    expect(messages()).toHaveLength(2)
    expect(logs().props('resultList')).toEqual([])
    event(failed, 'sending', '迟到内容')
    event(failed, 'dialogue_id', '旧 ID')
    failed.onClose()
    failed.onError()
    expect(messages().at(-1).content).toBe('')
    expect(input().props('loading')).toBe(true)
    await complete(retry, '重试成功')
    const next = await send('继续')
    expect(next.data.dialogue_id).toBe(dialogueId)
    expect(context(next)).toEqual([{ question: '重试', answer: '重试成功' }])
  })

  it('resets on reopen and isolates callbacks from the previous stream', async () => {
    createPanel()
    const first = await send('第一问')
    await complete(first, '第一答')
    const active = await send('进行中')
    wrapper.vm.open()
    await nextTick()
    expect(active.abort).toHaveBeenCalled()
    expect(wrapper.findComponent({ name: 'MessageList' }).exists()).toBe(false)
    expect(logs().exists()).toBe(false)
    const fresh = await send('新问题')
    event(active, 'sending', '旧内容')
    active.onClose()
    expect(fresh.data).toMatchObject({ dialogue_id: '', session_id: '', context_list: '[]' })
    expect(fresh.data.openid).not.toBe(first.data.openid)
    expect(messages()).toHaveLength(2)
    expect(messages().at(-1).content).toBe('')
  })

  it('prevents duplicate sends during streaming', async () => {
    createPanel()
    await send('问题')
    input().vm.$emit('send')
    await nextTick()
    expect(callWorkFlowDialog).toHaveBeenCalledTimes(1)
    expect(messages()).toHaveLength(2)
  })
})
