import { reactive, ref, computed } from 'vue'
import { defineStore } from 'pinia'
import {
  sendAiMessage,
  chatWelcome,
  getDialogueList,
  getChatMessage,
  getSessionRecordList,
  getSessionChannelList,
  editVariables,
} from '@/api/chat'
import { editPrompt } from '@/api/robot/index'
import { useClawbotStore } from './clawbot'
import { getUuid, getOpenid, devLog, extractVoiceInfo, removeVoiceFormat } from '@/utils/index'
import { useEventBus } from '@/hooks/event/useEventBus'
import { DEFAULT_USER_AVATAR } from '@/constants/index'

const defaultWelcomes = () => ({ content: '', question: [] })

const normalizeWelcomes = (welcomes) => {
  if (!welcomes) return defaultWelcomes()

  if (typeof welcomes === 'string') {
    try {
      return JSON.parse(welcomes) || defaultWelcomes()
    } catch {
      return defaultWelcomes()
    }
  }

  return {
    ...defaultWelcomes(),
    ...welcomes,
    question: Array.isArray(welcomes.question) ? welcomes.question : []
  }
}

const safeParseJson = (value, fallback) => {
  if (typeof value !== 'string') {
    return value ?? fallback
  }

  try {
    return JSON.parse(value)
  } catch {
    return fallback
  }
}

const PROCESS_EVENT_KEYS = ['FileOperation', 'ExecuteCommand']

/**
 * @typedef {'finish' | 'stop' | 'close'} ProcessFinalizeReason
 * @typedef {Object} NormalizedChatEvent
 * @property {'thinking_delta' | 'answer_delta' | 'round_begin' | 'round_finish' | 'tool_start' | 'tool_finish' | 'operation' | 'process_finalize' | 'final_snapshot'} type
 * @property {string} [content]
 * @property {string} [source]
 * @property {'tool' | 'skill'} [stepType]
 * @property {string} [title]
 * @property {string} [paramsText]
 * @property {string} [toolCallId]
 * @property {string} [toolName]
 * @property {string} [result]
 * @property {string} [eventName]
 * @property {string} [rawData]
 * @property {ProcessFinalizeReason} [reason]
 * @property {Object} [message]
 */

/**
 * Clawbot 聊天 Store（页面级）
 *
 * 职责：管理当前对话的消息列表、SSE 连接、对话记录等
 * 生命周期：进入 chat 页时初始化，离开时通过 $reset() 清理
 * 依赖：需要从 useClawbotStore 获取 currentAssistant 的 robot_key
 */
export const useClawbotChatStore = defineStore('clawbotChat', () => {
  const emitter = useEventBus()
  const clawbotStore = useClawbotStore()
  const messageList = ref([])
  const lastPushedUserMessageUid = ref('')

  let mySSE = null
  let sseRequestSeq = 0
  let chatCreateSeq = 0

  const abortCurrentSSE = () => {
    sseRequestSeq += 1
    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }
  }

  // 对话id
  const dialogue_id = ref(0)
  const rel_user_id = ref()

  const openid = ref('')
  // 用户信息
  const user = reactive({
    admin_user_id: '',
    avatar: '',
    id: '',
    name: '',
    nickname: '',
    openid: ''
  })
  // 机器人信息：从 clawbotStore.robotInfo 派生，附加 chat 专属的 openid
  const robot = computed(() => {
    const info = clawbotStore.robotInfo || {}
    return {
      ...info,
      openid: openid.value,
      enable_common_question: info.enable_common_question == 'true' || info.enable_common_question === true,
      common_question_list: info.common_question_list || [],
      question_multiple_switch: Number(info.question_multiple_switch) || 0,
      welcomes: normalizeWelcomes(info.welcomes),
    }
  })

  const getDefaultChatVariables = () => ({
    need_fill_variable: false,
    fill_variables: [],
    wait_variables: [],
    session_id: 0,
    dialogue_id: 0,
  })
  const chat_variables = ref(getDefaultChatVariables())

  // 创建对话
  const isNewChat = ref(false)
  const sendLock = ref(false)

  // 获取对话记录
  const myChatListSize = 25
  const myChatList = ref([])
  const myChatListLoading = ref(false)
  const myChatListLoadCompleted = ref(false)
  const chatListRobotKey = ref('')

  const resetMyChatList = (robotKey = '') => {
    chatListRobotKey.value = robotKey
    myChatList.value = []
    myChatListLoading.value = false
    myChatListLoadCompleted.value = false
  }

  const createChat = async (data) => {
    abortCurrentSSE()
    const requestCreateSeq = ++chatCreateSeq

    const nextRobotKey = data.robot_key || ''
    if (chatListRobotKey.value !== nextRobotKey) {
      resetMyChatList(nextRobotKey)
    } else if (!chatListRobotKey.value) {
      chatListRobotKey.value = nextRobotKey
    }

    messageList.value = []
    lastPushedUserMessageUid.value = ''
    chatMessageLoadCompleted.value = false
    sendLock.value = false
    chat_variables.value = getDefaultChatVariables()

    if (!data.dialogue_id) {
      isNewChat.value = true
      dialogue_id.value = 0
    } else {
      isNewChat.value = false
      dialogue_id.value = data.dialogue_id
    }

    openid.value = data.openid || getOpenid(data.robot_key)

    user.openid = openid.value
    user.avatar = data.avatar || DEFAULT_USER_AVATAR
    user.name = data.name || ''
    user.nickname = data.nickname || ''

    const res = await chatWelcome({
      robot_key: data.robot_key,
      openid: openid.value,
      nickname: user.nickname,
      name: user.name,
      is_background: data.is_background || undefined,
      dialogue_id: dialogue_id.value,
    })

    if (requestCreateSeq !== chatCreateSeq) {
      return null
    }

    const resData = res?.data || {}
    const userInfo = resData.customer || {}
    const robotData = resData.robot || {}
    const nextRobotData = { ...robotData }
    user.admin_user_id = userInfo.admin_user_id
    user.avatar = userInfo.avatar
    user.id = userInfo.id
    user.name = userInfo.name
    user.nickname = userInfo.nickname
    const normalizedWelcomes = normalizeWelcomes(nextRobotData.welcomes)
    // /chat/welcome 可能返回空 welcomes，保留 manage/getRobotInfo 已拉取到的欢迎语配置
    if (!normalizedWelcomes.content && !normalizedWelcomes.question.length) {
      delete nextRobotData.welcomes
    }
    // 将 chatWelcome 返回的 robot 数据合并回 clawbotStore，避免部分字段覆盖完整详情
    clawbotStore.updateRobotInfo(nextRobotData)

    setTimeout(() => {
      if (requestCreateSeq !== chatCreateSeq) {
        return
      }
      const chatVariable = resData.chat_variable || {}
      chat_variables.value = {
        ...getDefaultChatVariables(),
        ...chatVariable,
        session_id: Number(chatVariable.session_id || resData.session_id || 0),
        dialogue_id: Number(chatVariable.dialogue_id || resData.dialog_id || dialogue_id.value || 0),
        fill_variables: chatVariable.fill_variables || [],
        wait_variables: chatVariable.wait_variables || [],
      }
    })

    // clawbot 对话页需要保持空态 Hero，不在初始化时注入欢迎语消息
    return res
  }

  // 推送用户的消息到列表
  const pushUserMessage = (msg) => {
    msg.uid = getUuid(32)
    msg.loading = false
    msg.avatar = user.avatar
    msg.openid = user.openid
    msg.msg_type = msg.msg_type || 1
    msg.is_customer = 1
    messageList.value.push(msg)
    lastPushedUserMessageUid.value = msg.uid
  }

  const pushAiMessage = (msg) => {
    messageList.value.push(msg)
    emitter.emit('updateAiMessage', msg)
  }

  const ensureProcessStepState = (msg) => {
    if (!Array.isArray(msg.process_steps)) {
      msg.process_steps = []
    }
    if (typeof msg.current_round_index !== 'number') {
      msg.current_round_index = 0
    }
    if (typeof msg.active_thinking_step_id !== 'string') {
      msg.active_thinking_step_id = ''
    }
  }

  const createProcessStep = (step = {}) => ({
    id: step.id || getUuid(32),
    type: step.type || 'thinking',
    title: step.title || '',
    status: step.status || '',
    expanded: step.expanded === true,
    hidden: step.hidden === true,
    roundIndex: step.roundIndex || 0,
    contentText: step.contentText || '',
    paramsText: step.paramsText || '',
    resultText: step.resultText || '',
    eventName: step.eventName || '',
    tool_call_id: step.tool_call_id || '',
  })

  const appendProcessStep = (msg, step = {}) => {
    ensureProcessStepState(msg)
    const nextStep = createProcessStep(step)
    msg.process_steps.push(nextStep)
    return nextStep
  }

  const getProcessStepById = (msg, stepId) => {
    ensureProcessStepState(msg)
    return msg.process_steps.find((item) => item.id === stepId)
  }

  const finishActiveThinkingStep = (msg) => {
    const activeStep = getProcessStepById(msg, msg.active_thinking_step_id)
    if (activeStep) {
      if (!activeStep.contentText) {
        msg.process_steps = msg.process_steps.filter((item) => item.id !== activeStep.id)
      } else {
        activeStep.status = 'done'
        activeStep.resultText = ''
      }
    }
    msg.active_thinking_step_id = ''
  }

  const finalizeRunningProcessSteps = (msg, _reason) => {
    ensureProcessStepState(msg)
    msg.process_steps = msg.process_steps
      .filter((step) => !(step.type === 'thinking' && step.status === 'running' && !step.contentText))
      .map((step) => {
        if (step.status !== 'running') {
          return step
        }
        return {
          ...step,
          status: 'done',
          resultText: step.type === 'thinking' && !step.contentText ? '' : step.resultText,
        }
      })
    msg.active_thinking_step_id = ''
    msg.loading = false
    msg.startLoading = false
    msg.quote_loading = false
  }

  /**
   * 将后端原始 SSE 事件转换为页面无关的过程事件。
   *
   * 后端根据 application_type 选择旧链路或新链路，正常情况下两套正文事件互斥：
   * - reasoning_content -> thinking_delta
   * - stream_message.reasoning_content -> thinking_delta
   * - sending -> answer_delta
   * - stream_message.content -> answer_delta
   * - llm_rounds:begin / finish -> round_begin / round_finish
   * - tool_call_full / tool_result -> tool_start / tool_finish
   * - FileOperation、ExecuteCommand -> operation
   * - finish、用户停止、SSE 关闭 -> process_finalize
   * - ai_message -> final_snapshot
   *
   * stream_message 可能同时包含思考和正文，因此返回事件数组并保证先思考、后正文。
   * 首段非空正文表示当前思考阶段结束；空正文、空思考和元数据帧不改变过程状态。
   * ai_message 转为 final_snapshot，继续作为服务端持久化后的最终权威快照。
   * error、debug、耗时、引用和会话变量属于管理端附属状态，不进入过程状态机。
   *
   * @param {string} type 原始事件名
   * @param {*} content 原始事件数据
   * @returns {NormalizedChatEvent[]}
   */
  const normalizeSseEvent = (type, content) => {
    if (type === 'reasoning_content') {
      return [{ type: 'thinking_delta', content: String(content || ''), source: 'reasoning_content' }]
    }

    if (type === 'sending') {
      return [{ type: 'answer_delta', content: String(content || ''), source: 'sending' }]
    }

    if (type === 'stream_message') {
      const data = safeParseJson(content, {})
      const events = []
      if (data?.reasoning_content) {
        events.push({
          type: 'thinking_delta',
          content: String(data.reasoning_content),
          source: 'stream_message',
        })
      }
      if (data?.content) {
        events.push({
          type: 'answer_delta',
          content: String(data.content),
          source: 'stream_message',
        })
      }
      return events
    }

    if (type === 'llm_rounds') {
      if (content === 'begin') {
        return [{ type: 'round_begin' }]
      }
      if (content === 'finish') {
        return [{ type: 'round_finish' }]
      }
      return []
    }

    if (type === 'tool_call_full') {
      const functionInfo = content?.function || {}
      const functionArgs = safeParseJson(functionInfo.arguments, {})
      const isSkill = functionInfo.name === 'skill'
      return [{
        type: 'tool_start',
        stepType: isSkill ? 'skill' : 'tool',
        title: isSkill ? functionArgs?.skill || 'skill' : functionInfo.name || 'tool',
        paramsText: isSkill ? '' : functionInfo.arguments || '',
        toolCallId: content?.id || '',
      }]
    }

    if (type === 'tool_result') {
      const data = safeParseJson(content, {})
      return [{
        type: 'tool_finish',
        toolCallId: data?.tool_call_id || '',
        toolName: data?.tool_name || '',
        result: data?.content || '',
      }]
    }

    if (type === 'process_event') {
      return [{
        type: 'operation',
        eventName: content?.eventName || '',
        rawData: content?.rawData || '',
      }]
    }

    if (type === 'finalize_process_steps') {
      return [{ type: 'process_finalize', reason: content?.reason || 'finish' }]
    }

    if (type === 'ai_message') {
      return [{ type: 'final_snapshot', message: content }]
    }

    return []
  }

  /**
   * @param {Object} msg 当前机器人消息
   * @param {NormalizedChatEvent} event 统一过程事件
   */
  const applyNormalizedChatEvent = (msg, event) => {
    if (event.type === 'thinking_delta') {
      if (!event.content) {
        return
      }
      ensureProcessStepState(msg)
      let activeStep = getProcessStepById(msg, msg.active_thinking_step_id)
      if (!activeStep) {
        // 旧链路没有 llm_rounds:begin，首个思考增量至少归入第 1 轮。
        msg.current_round_index = Math.max(1, msg.current_round_index || 0)
        activeStep = appendProcessStep(msg, {
          type: 'thinking',
          status: 'running',
          expanded: false,
          roundIndex: msg.current_round_index,
          eventName: event.source,
        })
        msg.active_thinking_step_id = activeStep.id
      }
      activeStep.contentText = `${activeStep.contentText || ''}${event.content}`
      activeStep.resultText = ''
      // 原字段只用于接口和历史兼容，页面统一从 process_steps 渲染。
      msg.reasoning_content = `${msg.reasoning_content || ''}${event.content}`
      return
    }

    if (event.type === 'answer_delta') {
      // 空正文及元数据帧不能结束思考，首个非空正文才代表当前思考阶段完成。
      if (!event.content) {
        return
      }
      msg.startLoading = false
      finishActiveThinkingStep(msg)
      msg.content = `${msg.content || ''}${event.content}`
      msg.voice_content = extractVoiceInfo(msg.content)
      msg.content = removeVoiceFormat(msg.content)
      return
    }

    if (event.type === 'round_begin') {
      ensureProcessStepState(msg)
      msg.current_round_index += 1
      return
    }

    if (event.type === 'round_finish') {
      finishActiveThinkingStep(msg)
      return
    }

    if (event.type === 'tool_start') {
      appendProcessStep(msg, {
        type: event.stepType,
        title: event.title,
        status: 'running',
        expanded: false,
        roundIndex: msg.current_round_index || 0,
        paramsText: event.paramsText,
        tool_call_id: event.toolCallId,
        eventName: 'tool_call_full',
      })
      return
    }

    if (event.type === 'tool_finish') {
      ensureProcessStepState(msg)
      const matchedStep = event.toolCallId
        ? msg.process_steps.find((step) => {
            return ['tool', 'skill'].includes(step.type) &&
              step.status === 'running' &&
              step.tool_call_id === event.toolCallId
          })
        : [...msg.process_steps].reverse().find((step) => {
            return ['tool', 'skill'].includes(step.type) &&
              step.status === 'running' &&
              (!event.toolName || step.title === event.toolName)
          })

      // 结果无法关联时不补造步骤，避免错误完成其他并行工具。
      if (matchedStep) {
        matchedStep.status = 'done'
        matchedStep.resultText = event.result
      }
      return
    }

    if (event.type === 'operation') {
      // 文件和命令事件只保留审计数据，前端不展示也不会执行其中内容。
      appendProcessStep(msg, {
        type: 'operation',
        title: event.eventName,
        status: '',
        hidden: true,
        expanded: false,
        roundIndex: msg.current_round_index || 0,
        paramsText: event.rawData,
        eventName: event.eventName,
      })
      return
    }

    if (event.type === 'process_finalize') {
      // 正常完成、用户停止和断流统一清理 running，避免页面残留 spinner。
      finalizeRunningProcessSteps(msg, event.reason)
      return
    }

    if (event.type === 'final_snapshot') {
      const snapshot = event.message || {}
      finalizeRunningProcessSteps(msg, 'finish')
      if (snapshot.menu_json && snapshot.msg_type == 2) {
        const menuJson = safeParseJson(snapshot.menu_json, {})
        msg.question = menuJson.question || []
      }
      msg.id = snapshot.id
      msg.message_id = snapshot.message_id || snapshot.id || msg.message_id
      if (snapshot.msg_type !== undefined) {
        msg.msg_type = snapshot.msg_type
      }
      msg.content = snapshot.content || ''
      msg.voice_content = extractVoiceInfo(msg.content)
      msg.content = removeVoiceFormat(msg.content)
      if (snapshot.reply_content_list !== undefined) {
        msg.reply_content_list = safeParseJson(snapshot.reply_content_list, snapshot.reply_content_list || [])
      }
      if (snapshot.quote_file !== undefined) {
        msg.quote_file = safeParseJson(snapshot.quote_file, snapshot.quote_file || [])
      }
    }
  }

  /**
   * 旧历史机器人文本消息没有 process_steps 时，将 reasoning_content 转成确定性的完成步骤。
   * 已有步骤、用户消息、菜单、图片和空思考均不转换，避免分页或重进时重复生成。
   */
  const normalizeHistoricalProcessSteps = (msg) => {
    const parsedSteps = Array.isArray(msg.process_steps)
      ? msg.process_steps
      : safeParseJson(msg.process_steps, [])
    if (Array.isArray(parsedSteps) && parsedSteps.length > 0) {
      return parsedSteps
    }

    const reasoningContent = typeof msg.reasoning_content === 'string' ? msg.reasoning_content : ''
    if (msg.is_customer == 1 || msg.msg_type != 1 || !reasoningContent.trim()) {
      return []
    }

    return [createProcessStep({
      id: `legacy-thinking-${msg.message_id || msg.id || msg.uid}`,
      type: 'thinking',
      status: 'done',
      expanded: false,
      hidden: false,
      roundIndex: 1,
      contentText: reasoningContent,
      resultText: '',
      eventName: 'reasoning_content',
    })]
  }

  // 更新AI的消息到列表
  const updateAiMessage = (type, content, uid) => {
    let msgIndex = messageList.value.findIndex((item) => item.uid == uid)
    if (msgIndex === -1) {
      return
    }

    const currentMessage = messageList.value[msgIndex]
    ensureProcessStepState(currentMessage)

    const normalizedEvents = normalizeSseEvent(type, content)
    if (normalizedEvents.length > 0) {
      normalizedEvents.forEach((event) => applyNormalizedChatEvent(currentMessage, event))
      emitter.emit('updateAiMessage', currentMessage)
      return
    }

    if (type == 'reply_content_list') {
      if (content !== undefined && typeof content === 'string') {
        currentMessage.reply_content_list = safeParseJson(content, [])
      }
    }

    if (type == 'start_quote_file') {
      currentMessage.quote_loading = true
    }

    if (type == 'quote_file') {
      currentMessage.quote_file = content?.length > 0 ? content : []
      currentMessage.show_quote_file = true
      currentMessage.quote_loading = false
    }

    if (type == 'debug') {
      currentMessage.debug = content?.length > 0 ? content : []
    }

    if (type == 'error') {
      currentMessage.error = content?.length > 0 ? content : []
    }
    if (type == 'recall_time') {
      currentMessage.recall_time = content || 0
    }
    if (type == 'request_time') {
      currentMessage.request_time = content || 0
    }

    emitter.emit('updateAiMessage', messageList.value[msgIndex])
  }

  // 关闭AI的消息加载状态
  const closeAiMessageLoading = () => {
    let msgIndex = -1
    for (let i = messageList.value.length - 1; i >= 0; i -= 1) {
      if (messageList.value[i]?.loading) {
        msgIndex = i
        break
      }
    }
    if (msgIndex === -1) {
      return
    }
    messageList.value[msgIndex].loading = false
  }

  // 查找当前正在生成中的 AI 消息索引
  const getRunningAiMessageIndex = () => {
    for (let i = messageList.value.length - 1; i >= 0; i--) {
      const item = messageList.value[i]
      const hasRunningProcess = Array.isArray(item.process_steps) && item.process_steps.some((step) => {
        return step?.status === 'running'
      })
      if (
        item.is_customer != 1 &&
        !item.is_stopped &&
        (item.loading || item.startLoading || item.quote_loading || hasRunningProcess)
      ) {
        return i
      }
    }
    return -1
  }

  // 手动终止当前会话
  const stopMessage = () => {
    const msgIndex = getRunningAiMessageIndex()
    if (msgIndex > -1) {
      const msg = messageList.value[msgIndex]
      msg.is_stopped = true
      msg.loading = false
      msg.startLoading = false
      msg.quote_loading = false
      updateAiMessage('finalize_process_steps', { reason: 'stop' }, msg.uid)
    }

    abortCurrentSSE()
    sendLock.value = false
    closeAiMessageLoading()
  }

  // 发送消息
  const sendMessage = (data) => {
    if (sendLock.value) {
      return
    }

    let aiMsg = {
      startLoading: true,
      loading: true,
      id: '',
      content: '',
      reasoning_content: '',
      uid: getUuid(32),
      robot_avatar: robot.value.robot_avatar,
      msg_type: 1,
      quote_file: [],
      debug: [],
      error: '',
      recall_time: '',
      request_time: '',
      process_expanded: true,
      process_steps: [],
      current_round_index: 0,
      active_thinking_step_id: '',
      quote_loading: false,
      show_quote_file: true,
      is_stopped: false,
      voice_content: [],
    }
    let params = {
      robot_key: robot.value.robot_key,
      openid: robot.value.openid,
      question: data.message,
      form_ids: robot.value.form_ids,
      dialogue_id: dialogue_id.value,
      global: data.global
    }

    let variables_key = `chat_prompt_variables_${robot.value.robot_key}`
    const localVariables = localStorage.getItem(variables_key)
    const isNewDialogue = Number(dialogue_id.value || 0) === 0

    if (isNewDialogue && localVariables) {
      params.chat_prompt_variables = localVariables
      localStorage.removeItem(variables_key)
    }

    if (import.meta.env.DEV) {
      params.debug = 0
    }

    sendLock.value = true
    const requestSeq = ++sseRequestSeq
    mySSE = sendAiMessage(params)

    mySSE.onMessage = (res) => {
      if (requestSeq !== sseRequestSeq) {
        return
      }
      if (res.event !== 'sending') {
        devLog(res)
      }
      if (res.event == 'dialogue_id') {
        dialogue_id.value = res.data
      }
      if (res.event == 'c_message') {
        let data = safeParseJson(res.data, {})
        pushUserMessage(data)
      }
      if (res.event == 'robot') {
        pushAiMessage(aiMsg)
        if (isNewChat.value) {
          insertNewSession()
          isNewChat.value = false
        }
      }
      if (res.event == 'reply_content_list') {
        updateAiMessage('reply_content_list', res.data, aiMsg.uid)
      }
      if (res.event == 'reasoning_content') {
        updateAiMessage('reasoning_content', res.data, aiMsg.uid)
      }
      if (res.event == 'llm_rounds') {
        updateAiMessage('llm_rounds', res.data, aiMsg.uid)
      }
      if (res.event == 'stream_message') {
        updateAiMessage('stream_message', res.data, aiMsg.uid)
      }
      if (res.event == 'sending') {
        updateAiMessage('sending', res.data, aiMsg.uid)
      }
      if (res.event == 'tool_call_full') {
        const data = safeParseJson(res.data, {})
        updateAiMessage('tool_call_full', data, aiMsg.uid)
      }
      if (res.event == 'tool_result') {
        updateAiMessage('tool_result', res.data, aiMsg.uid)
      }
      if (PROCESS_EVENT_KEYS.includes(res.event)) {
        updateAiMessage('process_event', {
          eventName: res.event,
          rawData: typeof res.data === 'string' ? res.data : JSON.stringify(res.data || {})
        }, aiMsg.uid)
      }
      if (res.event == 'ai_message') {
        let data = safeParseJson(res.data, {})
        updateAiMessage('ai_message', data, aiMsg.uid)
      }
      if (res.event == 'finish') {
        updateAiMessage('finalize_process_steps', { reason: 'finish' }, aiMsg.uid)
      }
      if (res.event == 'start_quote_file') {
        updateAiMessage('start_quote_file', res.data, aiMsg.uid)
      }
      if (res.event == 'quote_file') {
        let data = safeParseJson(res.data, [])
        updateAiMessage('quote_file', data, aiMsg.uid)
      }
      if (res.event == 'debug') {
        let data = safeParseJson(res.data, [])
        updateAiMessage('debug', data, aiMsg.uid)
      }
      if (res.event == 'error') {
        let data = res.data
        updateAiMessage('error', data, aiMsg.uid)
      }
      if (res.event == 'recall_time') {
        let data = res.data
        updateAiMessage('recall_time', data, aiMsg.uid)
      }
      if (res.event == 'request_time') {
        let data = res.data
        updateAiMessage('request_time', data, aiMsg.uid)
      }
      if (res.event == 'chat_prompt_variables') {
        let data = res.data
        if (data) {
          data = safeParseJson(data, {})
          chat_variables.value.need_fill_variable = data.need_fill_variable
          chat_variables.value.fill_variables = data.fill_variables || []
          chat_variables.value.session_id = data.session_id
          chat_variables.value.dialogue_id = data.dialogue_id || dialogue_id.value
        }
      }
    }

    mySSE.onClose = () => {
      if (requestSeq !== sseRequestSeq) {
        return
      }
      updateAiMessage('finalize_process_steps', { reason: 'close' }, aiMsg.uid)
      closeAiMessageLoading()
      sendLock.value = false
      mySSE = null
    }
  }

  const handleEditVariables = (data) => {
    return editVariables({
      robot_key: robot.value.robot_key,
      openid: robot.value.openid,
      dialogue_id: dialogue_id.value,
      chat_prompt_variables: JSON.stringify(data.chat_prompt_variables),
      session_id: chat_variables.value.session_id
    }).then(res => {
      chat_variables.value.fill_variables = data.chat_prompt_variables
      let variables_key = `chat_prompt_variables_${robot.value.robot_key}`
      localStorage.removeItem(variables_key)
      return res
    })
  }

  const getMyChatList = async () => {
    if (myChatListLoadCompleted.value || myChatListLoading.value) {
      return false
    }

    const requestRobotKey = robot.value.robot_key || ''
    if (!requestRobotKey) {
      return false
    }

    let min_id = 0
    if (myChatList.value.length > 0) {
      min_id = myChatList.value[myChatList.value.length - 1].id
    }

    myChatListLoading.value = true

    try {
      const res = await getDialogueList({
        min_id: min_id,
        size: myChatListSize,
        robot_key: requestRobotKey
      })

      if ((robot.value.robot_key || '') !== requestRobotKey || chatListRobotKey.value !== requestRobotKey) {
        return false
      }

      const list = res.data || []
      if (list.length === 0) {
        myChatListLoadCompleted.value = true
        return false
      }

      myChatList.value = [...myChatList.value, ...list]
      return res
    } finally {
      if ((robot.value.robot_key || '') === requestRobotKey && chatListRobotKey.value === requestRobotKey) {
        myChatListLoading.value = false
      }
    }
  }

  // 插入最新一条对话记录
  const insertNewSession = () => {
    const requestRobotKey = robot.value.robot_key || ''
    if (!requestRobotKey) {
      return
    }

    getDialogueList({
      min_id: 0,
      size: 1,
      robot_key: requestRobotKey
    }).then((res) => {
      if ((robot.value.robot_key || '') !== requestRobotKey || chatListRobotKey.value !== requestRobotKey) {
        return
      }
      const list = res.data || []
      if (list[0]) {
        myChatList.value.unshift(list[0])
      }
    })
  }

  // 打开对话
  const openChat = async (data) => {
    abortCurrentSSE()
    let res = await createChat(data)
    return res
  }

  // 获取聊天记录
  const chatMessagePageSize = 20
  const chatMessageLoadCompleted = ref(false)

  const onGetChatMessage = async () => {
    if (chatMessageLoadCompleted.value) {
      return
    }

    let min_id = 0
    const currentMessages = messageList.value.filter((item) => !item.isWelcome)

    if (currentMessages.length > 0) {
      min_id = currentMessages[0].id
    }

    const request_dialogue_id = dialogue_id.value
    const request_robot_key = robot.value.robot_key || ''
    let params = {
      robot_key: request_robot_key,
      openid: user.openid,
      min_id: min_id,
      size: chatMessagePageSize,
      dialogue_id: request_dialogue_id,
      rel_user_id: rel_user_id.value
    }

    const res = await getChatMessage(params)
    if (dialogue_id.value !== request_dialogue_id || (robot.value.robot_key || '') !== request_robot_key) {
      return null
    }
    const historyList = res?.data?.list || []
    const _customer = res?.data?.customer || {}
    const newRobot = res?.data?.robot || {}
    if (historyList.length === 0) {
      chatMessageLoadCompleted.value = true
      return
    }
    historyList.sort((a, b) => a.id - b.id)
    historyList.forEach((item) => {
      item.loading = false
      item.uid = getUuid(32)
      item.name = item.name || item.nickname
      if (item.is_customer == 1) {
        item.name = item.name || user.name || _customer.name
        item.avatar = item.avatar || user.avatar || _customer.avatar
      } else {
        item.name = item.name || robot.value.robot_name || newRobot.robot_name
        item.robot_avatar = item.avatar || robot.value.robot_avatar || newRobot.robot_avatar
        item.avatar = item.robot_avatar
      }
      if (item.menu_json) {
        item.menu_json = safeParseJson(item.menu_json, {})
      }
      item.process_steps = normalizeHistoricalProcessSteps(item)
      item.process_expanded = typeof item.process_expanded === 'boolean'
        ? item.process_expanded
        : true
      item.current_round_index = Number(item.current_round_index || 0)
      item.active_thinking_step_id = item.active_thinking_step_id || ''
      if (item.quote_file) {
        item.quote_file = safeParseJson(item.quote_file, [])
      }
      if (item.reply_content_list && typeof item.reply_content_list === 'string') {
        item.reply_content_list = safeParseJson(item.reply_content_list, [])
      }
      item.voice_content = extractVoiceInfo(item.content)
      item.content = removeVoiceFormat(item.content)
    })

    messageList.value = [...historyList, ...messageList.value]
    return res
  }

  const getRecordList = async (params) => {
    const res = await getSessionRecordList(params)
    if (!res) {
      throw new Error('Failed to get session record list')
    }
    return res
  }

  const getChannelList = async (params) => {
    const res = await getSessionChannelList(params)
    if (!res) {
      throw new Error('Failed to get session channel list')
    }
    return res
  }

  // 提示词
  const changeRobotPrompt = (text) => {
    const data = { prompt: text }
    if (robot.value.id) {
      data.id = robot.value.id
    }
    clawbotStore.updateRobotInfo(data)
  }

  const saveRobotPrompt = () => {
    return editPrompt({
      id: robot.value.id,
      prompt: robot.value.prompt
    }).then(() => {
      // 保存后刷新权威源确保一致
      clawbotStore.fetchRobotInfo()
    })
  }

  function $reset() {
    abortCurrentSSE()
    chatCreateSeq += 1
    dialogue_id.value = 0
    messageList.value = []
    lastPushedUserMessageUid.value = ''
    openid.value = ''
    user.admin_user_id = ''
    user.avatar = ''
    user.id = ''
    user.name = ''
    user.nickname = ''
    user.openid = ''
    // robot 配置由 clawbotStore 管理，此处不清除
    isNewChat.value = false
    chatMessageLoadCompleted.value = false
    sendLock.value = false
    resetMyChatList()
    chat_variables.value = getDefaultChatVariables()
  }

  return {
    $reset,
    user,
    robot,
    dialogue_id,
    openid,
    sendLock,
    messageList,
    lastPushedUserMessageUid,
    createChat,
    sendMessage,
    stopMessage,
    getMyChatList,
    myChatList,
    openChat,
    onGetChatMessage,
    changeRobotPrompt,
    saveRobotPrompt,
    getRecordList,
    getChannelList,
    chat_variables,
    handleEditVariables
  }
})
