import { computed, reactive, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { chatWelcome, getDialogueList, getChatMessage, sendAiMessage, editVariables } from '@/api/chat'
import { checkChatRequestPermission } from '@/api/robot'
import { getUuid, getOpenid, extractVoiceInfo, removeVoiceFormat } from '@/utils/index'
import {
  safeParseJson,
  ensureProcessStepState,
  normalizeSseEvent,
  applyNormalizedChatEvent,
  normalizeHistoricalProcessSteps
} from '@/utils/agent-process'
import { DEFAULT_USER_AVATAR } from '@/constants'
import { normalizeCustomerServiceQueryValue } from '@/views/customer-service/navigation'
import { useCustomerServiceStore } from './customer-service'

const defaultVariables = () => ({
  need_fill_variable: false,
  fill_variables: [],
  wait_variables: [],
  session_id: 0,
  dialogue_id: 0
})

function normalizeWelcomes(value) {
  const config = safeParseJson(value, {}) || {}
  return { ...config, content: config.content || '', question: Array.isArray(config.question) ? config.question : [] }
}

// 客服页面私有会话状态：只读基础Store身份，不选择/创建机器人，也不使用Agent业务Store。
export const useCustomerServiceChatStore = defineStore('customer-service-chat', () => {
  const customerServiceStore = useCustomerServiceStore()
  const messageList = ref([])
  const lastPushedUserMessageUid = ref('')
  const dialogue_id = ref(0)
  const rel_user_id = ref()
  const openid = ref('')
  const user = reactive({ admin_user_id: '', avatar: '', id: '', name: '', nickname: '', openid: '' })
  const chat_variables = ref(defaultVariables())
  const isNewChat = ref(false)
  const chatStatus = ref('idle')
  const chatError = ref('')
  const chatLoading = computed(() => chatStatus.value === 'loading')
  const sendStatus = ref('idle')
  const sendLock = computed(() => sendStatus.value !== 'idle')
  const sendError = ref('')
  const sensitiveWords = ref([])
  const myChatListPageSize = 25
  const myChatList = ref([])
  const myChatListLoading = ref(false)
  const myChatListLoadCompleted = ref(false)
  const myChatListHasMore = computed(() => !myChatListLoadCompleted.value && myChatList.value.length >= myChatListPageSize)
  const myChatListError = ref('')
  const chatMessageLoading = ref(false)
  const chatMessageLoadCompleted = ref(false)
  const chatMessageError = ref('')
  const variablesLoading = ref(false)
  const variablesError = ref('')

  // 身份值相同也可能已经A→B→A；上下文/会话代次不能重置归零。
  let contextGeneration = 0
  let sessionGeneration = 0
  let contextId = ''
  let contextRobotKey = ''
  let listRequestSequence = 0
  let listRefreshSequence = 0
  let historyRequestSequence = 0
  let variableRequestSequence = 0
  let activeSend = null

  const robot = computed(() => {
    const info = customerServiceStore.robot || {}
    return {
      ...info,
      openid: openid.value,
      enable_common_question: info.enable_common_question === true || info.enable_common_question === 'true',
      common_question_list: info.common_question_list || [],
      question_multiple_switch: Number(info.question_multiple_switch) || 0,
      welcomes: normalizeWelcomes(info.welcomes)
    }
  })

  const isContextCurrent = (generation) => {
    return generation === contextGeneration && !!contextId && !!contextRobotKey &&
      customerServiceStore.robotStatus === 'ready' &&
      normalizeCustomerServiceQueryValue(customerServiceStore.robot?.id) === contextId &&
      customerServiceStore.robot?.robot_key === contextRobotKey
  }
  const isSessionCurrent = (context, session) => {
    return isContextCurrent(context) && session === sessionGeneration
  }

  const updateAiMessage = (message, type, content) => {
    ensureProcessStepState(message)
    const events = normalizeSseEvent(type, content)
    if (events.length) {
      events.forEach((event) => applyNormalizedChatEvent(message, event))
      return
    }
    if (type === 'reply_content_list') message.reply_content_list = safeParseJson(content, [])
    if (type === 'start_quote_file') message.quote_loading = true
    if (type === 'quote_file') {
      message.quote_file = Array.isArray(content) ? content : []
      message.show_quote_file = true
      message.quote_loading = false
    }
    if (type === 'debug') message.debug = content || []
    if (type === 'error') {
      message.error = content?.length > 0 ? content : []
      applyNormalizedChatEvent(message, { type: 'process_finalize', reason: 'error' })
    }
    if (type === 'recall_time') message.recall_time = content || 0
    if (type === 'request_time') message.request_time = content || 0
  }

  // 接纳和结束可能先后触发；每次发送只能结算一次，不能让迟到的错误撤销已接纳结果。
  const settleSendResult = (attempt, accepted) => {
    if (attempt.resultSettled) return
    attempt.resultSettled = true
    attempt.resolveResult(accepted)
  }

  // 所有终止入口只收尾自己的连接；旧onClose/onError不能解锁后来建立的发送。
  const finishSend = (attempt, reason, error = '') => {
    if (!attempt || attempt.ended) return
    attempt.ended = true
    settleSendResult(attempt, false)
    if (activeSend === attempt) {
      activeSend = null
      if (isSessionCurrent(attempt.context, attempt.session)) {
        // 接纳后仍未收到创建id，不能伪造列表项；保留显式刷新恢复入口。
        if (attempt.newSessionAccepted && !attempt.sessionInserted) myChatListError.value = 'error_chat_list'
        if (attempt.aiMessage) {
          updateAiMessage(attempt.aiMessage, 'finalize_process_steps', { reason })
          if (reason === 'stop') attempt.aiMessage.is_stopped = true
        }
        sendStatus.value = 'idle'
        if (error) sendError.value = error
      }
    }
    attempt.connection?.abort()
  }

  const stopMessage = () => finishSend(activeSend, 'stop')

  const clearConversation = () => {
    stopMessage()
    sessionGeneration++
    historyRequestSequence++
    variableRequestSequence++
    dialogue_id.value = 0
    rel_user_id.value = undefined
    openid.value = ''
    Object.assign(user, { admin_user_id: '', avatar: '', id: '', name: '', nickname: '', openid: '' })
    messageList.value = []
    lastPushedUserMessageUid.value = ''
    chat_variables.value = defaultVariables()
    isNewChat.value = false
    chatStatus.value = 'idle'
    chatError.value = ''
    sendStatus.value = 'idle'
    sendError.value = ''
    sensitiveWords.value = []
    chatMessageLoading.value = false
    chatMessageLoadCompleted.value = false
    chatMessageError.value = ''
    variablesLoading.value = false
    variablesError.value = ''
  }

  const $reset = () => {
    clearConversation()
    contextGeneration++
    listRequestSequence++
    contextId = ''
    contextRobotKey = ''
    myChatList.value = []
    myChatListLoading.value = false
    myChatListLoadCompleted.value = false
    myChatListError.value = ''
  }

  // 基础详情切换时立即使旧回调失效；同身份配置对象替换不会重建/清空会话。
  watch(
    [() => customerServiceStore.robotStatus,
      () => normalizeCustomerServiceQueryValue(customerServiceStore.robot?.id),
      () => customerServiceStore.robot?.robot_key || ''],
    ([status, id, key]) => {
      if (contextId && (status !== 'ready' || id !== contextId || key !== contextRobotKey)) $reset()
    },
    { flush: 'sync' }
  )

  /** 新建只清会话，不清同机器人的列表；data沿Agent传openid/avatar/name/nickname/dialogue_id。 */
  const createChat = async (data = {}) => {
    const source = customerServiceStore.robot
    if (customerServiceStore.robotStatus !== 'ready' || !source?.robot_key || chatLoading.value) return null
    const id = normalizeCustomerServiceQueryValue(source.id)
    if (!id || (data.robot_key && data.robot_key !== source.robot_key)) return null
    if (contextId !== id || contextRobotKey !== source.robot_key) {
      $reset()
      contextId = id
      contextRobotKey = source.robot_key
    }
    clearConversation()
    const context = contextGeneration
    const session = sessionGeneration
    // 每次welcome捕获当前原robot对象，只供配置回填校验，不能替代会话代次。
    const robotSnapshot = customerServiceStore.robot
    dialogue_id.value = data.dialogue_id || 0
    rel_user_id.value = data.rel_user_id
    isNewChat.value = Number(data.dialogue_id || 0) === 0
    openid.value = data.openid || getOpenid()
    Object.assign(user, {
      openid: openid.value,
      avatar: data.avatar || DEFAULT_USER_AVATAR,
      name: data.name || '',
      nickname: data.nickname || ''
    })
    chatStatus.value = 'loading'
    try {
      const response = await chatWelcome({
        robot_key: contextRobotKey,
        openid: openid.value,
        nickname: user.nickname,
        name: user.name,
        is_background: 1,
        dialogue_id: dialogue_id.value
      })
      if (!isSessionCurrent(context, session)) return null
      const responseData = response?.data
      if (!responseData || typeof responseData !== 'object') throw new Error('Invalid welcome response')
      const customer = responseData.customer || {}
      for (const field of ['admin_user_id', 'avatar', 'id', 'name', 'nickname']) {
        if (customer[field] !== undefined) user[field] = customer[field]
      }
      const robotInfo = { ...(responseData.robot || {}) }
      const welcomes = normalizeWelcomes(robotInfo.welcomes)
      if (!welcomes.content && !welcomes.question.length) delete robotInfo.welcomes
      // 快照过期只拒绝配置，不丢弃同会话的customer/变量，也不再次新建会话。
      customerServiceStore.mergeCurrentRobotInfo(robotSnapshot, robotInfo)
      const variables = responseData.chat_variable || {}
      chat_variables.value = {
        ...defaultVariables(),
        ...variables,
        session_id: Number(variables.session_id || responseData.session_id || 0),
        dialogue_id: Number(variables.dialogue_id || responseData.dialog_id || dialogue_id.value || 0),
        fill_variables: Array.isArray(variables.fill_variables) ? variables.fill_variables : [],
        wait_variables: Array.isArray(variables.wait_variables) ? variables.wait_variables : []
      }
      chatStatus.value = 'ready'
      // 新建保持Hero，不插入欢迎消息。
      return response
    } catch {
      if (isSessionCurrent(context, session)) {
        chatStatus.value = 'error'
        chatError.value = 'error_create_chat'
      }
      return null
    }
  }

  /** 首次进入：welcome和列表互不阻塞；返回welcome结果，列表错误单独展示/重试。 */
  const initialize = async (data = {}) => {
    if (chatLoading.value) return null
    $reset()
    const welcome = createChat({ ...data, dialogue_id: 0 })
    void getMyChatList()
    return welcome
  }

  /** refresh用于列表错误重试，成功后才替换旧数据；普通分页失败不推进游标。 */
  const getMyChatList = async ({ refresh = false } = {}) => {
    if (!isContextCurrent(contextGeneration) || myChatListLoading.value ||
      (!refresh && myChatListLoadCompleted.value)) return null
    const context = contextGeneration
    const requestId = ++listRequestSequence
    if (refresh) listRefreshSequence++
    const previousIds = new Set(myChatList.value.map((item) => String(item.id)))
    const minId = refresh ? 0 : myChatList.value[myChatList.value.length - 1]?.id || 0
    myChatListLoading.value = true
    myChatListError.value = ''
    try {
      const response = await getDialogueList({ robot_key: contextRobotKey, min_id: minId, size: myChatListPageSize })
      if (!isContextCurrent(context) || requestId !== listRequestSequence) return null
      if (!Array.isArray(response?.data)) throw new Error('Invalid dialogue list response')
      const rows = response.data
      // 刷新等待期间可能补入新的真实会话；保留这些新增项，避免旧刷新响应将其抹掉。
      const existing = refresh ? myChatList.value.filter((item) => !previousIds.has(String(item.id))) : myChatList.value
      const ids = new Set(existing.map((item) => String(item.id)))
      myChatList.value = [...existing, ...rows.filter((item) => !ids.has(String(item.id)))]
      myChatListLoadCompleted.value = rows.length < myChatListPageSize
      return response
    } catch {
      if (isContextCurrent(context) && requestId === listRequestSequence) myChatListError.value = 'error_chat_list'
      return null
    } finally {
      if (isContextCurrent(context) && requestId === listRequestSequence) myChatListLoading.value = false
    }
  }

  const insertNewSession = async (context, robotKey, createdDialogueId) => {
    const refreshSequence = listRefreshSequence
    const isCurrentList = () => isContextCurrent(context) && contextRobotKey === robotKey &&
      refreshSequence === listRefreshSequence
    try {
      const response = await getDialogueList({ robot_key: robotKey, min_id: 0, size: 1 })
      if (!isCurrentList()) return
      if (!Array.isArray(response?.data)) throw new Error('Invalid new dialogue response')
      if (myChatList.value.some((row) => String(row.dialogue_id ?? row.id) === createdDialogueId)) return
      const item = response.data[0]
      // size:1只返回当时最新项；若已被另一新会话抢先，不冒充本次记录，交列表刷新恢复。
      if (!item || String(item.dialogue_id ?? item.id) !== createdDialogueId) throw new Error('New dialogue not found')
      // 列表归机器人所有，不随切会话失效；补入之间不递增请求序号，避免互相吞掉结果。
      if (!myChatList.value.some((row) => String(row.id) === String(item.id))) myChatList.value.unshift(item)
    } catch {
      if (isCurrentList()) myChatListError.value = 'error_chat_list'
    }
  }

  const onGetChatMessage = async () => {
    if (!isContextCurrent(contextGeneration) || chatStatus.value !== 'ready' || Number(dialogue_id.value || 0) === 0 ||
      chatMessageLoading.value || chatMessageLoadCompleted.value || sendLock.value) return null
    const context = contextGeneration
    const session = sessionGeneration
    const requestId = ++historyRequestSequence
    const minId = messageList.value[0]?.id || 0
    chatMessageLoading.value = true
    chatMessageError.value = ''
    try {
      const response = await getChatMessage({
        robot_key: contextRobotKey,
        openid: openid.value,
        dialogue_id: dialogue_id.value,
        min_id: minId,
        size: 20,
        rel_user_id: rel_user_id.value
      })
      if (!isSessionCurrent(context, session) || requestId !== historyRequestSequence) return null
      if (!Array.isArray(response?.data?.list)) throw new Error('Invalid message list response')
      const customer = response.data.customer || {}
      const historyRobot = response.data.robot || {}
      const rows = response.data.list.slice().sort((a, b) => a.id - b.id).map((row) => {
        const item = { ...row, loading: false, uid: getUuid(32), name: row.name || row.nickname }
        if (item.is_customer == 1) {
          item.name = item.name || user.name || customer.name
          item.avatar = item.avatar || user.avatar || customer.avatar
        } else {
          item.name = item.name || robot.value.robot_name || historyRobot.robot_name
          item.robot_avatar = item.avatar || robot.value.robot_avatar || historyRobot.robot_avatar
          item.avatar = item.robot_avatar
        }
        if (item.menu_json) item.menu_json = safeParseJson(item.menu_json, {})
        item.process_steps = normalizeHistoricalProcessSteps(item)
        item.process_expanded = false
        item.current_round_index = Number(item.current_round_index || 0)
        item.active_thinking_step_id = item.active_thinking_step_id || ''
        if (item.quote_file) item.quote_file = safeParseJson(item.quote_file, [])
        if (item.reply_content_list) item.reply_content_list = safeParseJson(item.reply_content_list, [])
        item.voice_content = extractVoiceInfo(item.content)
        item.content = removeVoiceFormat(item.content)
        return item
      })
      const ids = new Set(messageList.value.map((item) => String(item.id)))
      messageList.value = [...rows.filter((item) => !ids.has(String(item.id))), ...messageList.value]
      chatMessageLoadCompleted.value = rows.length === 0
      return response
    } catch {
      if (isSessionCurrent(context, session) && requestId === historyRequestSequence) chatMessageError.value = 'error_chat_messages'
      return null
    } finally {
      if (isSessionCurrent(context, session) && requestId === historyRequestSequence) chatMessageLoading.value = false
    }
  }

  /** 打开历史包含首批消息加载；分页仍由onGetChatMessage处理。 */
  const openChat = async (data) => {
    const response = await createChat(data)
    if (!response) return null
    return onGetChatMessage()
  }

  const pushUserMessage = (data) => {
    const message = {
      ...data,
      uid: getUuid(32),
      loading: false,
      avatar: user.avatar,
      openid: user.openid,
      msg_type: data.msg_type || 1,
      is_customer: 1
    }
    messageList.value.push(message)
    lastPushedUserMessageUid.value = message.uid
  }

  /**
   * message为协议question（可为图片JSON），checkQuestion可传原始文本；global默认JSON空对象。
   * true表示当前发送收到c_message/robot业务接纳；接纳前失败/停止/切换/关闭返回false，保留草稿。
   * 调用方仍须核对提交时的页面和草稿版本，不能用旧结果清掉后来编辑的内容。
   */
  const sendMessage = async (data = {}) => {
    if (sendLock.value || chatStatus.value !== 'ready' || chatMessageLoading.value || variablesLoading.value ||
      !isContextCurrent(contextGeneration)) return false
    sendError.value = ''
    sensitiveWords.value = []
    if (typeof data.message !== 'string' || !data.message) {
      sendError.value = 'empty_message'
      return false
    }
    if (chat_variables.value.need_fill_variable === true || chat_variables.value.need_fill_variable == 1) {
      sendError.value = 'variables_required'
      return false
    }
    const attempt = {
      context: contextGeneration,
      session: sessionGeneration,
      connection: null,
      ended: false,
      aiMessage: null,
      userPushed: false,
      robotKey: contextRobotKey,
      createdDialogueId: '',
      newSessionAccepted: false,
      sessionInserted: false,
      resultSettled: false,
      resolveResult: null
    }
    const result = new Promise((resolve) => { attempt.resolveResult = resolve })
    activeSend = attempt
    sendStatus.value = 'checking'
    const isCurrentSend = () => activeSend === attempt && !attempt.ended && isSessionCurrent(attempt.context, attempt.session)
    try {
      // 停止/切换即返回false，不必等旧检查请求结束；请求本身仍沿用现有封装。
      const permission = await Promise.race([checkChatRequestPermission({
        robot_key: contextRobotKey,
        openid: openid.value,
        question: data.checkQuestion ?? data.message,
        form_ids: robot.value.form_ids,
        dialogue_id: dialogue_id.value
      }), result])
      if (!isCurrentSend()) return false
      const words = permission?.data?.words
      if (words && words.length) {
        sensitiveWords.value = Array.isArray(words) ? words : [words]
        finishSend(attempt, 'stop', 'sensitive_words')
        return false
      }
      const params = {
        robot_key: contextRobotKey,
        openid: openid.value,
        question: data.message,
        form_ids: robot.value.form_ids,
        dialogue_id: dialogue_id.value,
        global: data.global ?? JSON.stringify({})
      }
      const variablesKey = `chat_prompt_variables_${contextRobotKey}`
      const cachedVariables = Number(dialogue_id.value || 0) === 0 ? localStorage.getItem(variablesKey) : null
      if (cachedVariables) params.chat_prompt_variables = cachedVariables
      if (import.meta.env.DEV) params.debug = 0
      const aiMessage = reactive({
        uid: getUuid(32), id: '', content: '', reasoning_content: '', msg_type: 1,
        robot_avatar: robot.value.robot_avatar || robot.value.robot_avatar_url,
        startLoading: true, loading: true, quote_file: [], debug: [], error: '',
        recall_time: '', request_time: '', process_expanded: true, process_steps: [],
        current_round_index: 0, active_thinking_step_id: '', quote_loading: false,
        show_quote_file: true, is_stopped: false, voice_content: []
      })
      attempt.aiMessage = aiMessage
      let aiPushed = false
      const ensureAiMessage = () => {
        if (!aiPushed) {
          messageList.value.push(aiMessage)
          aiPushed = true
        }
      }
      const insertAcceptedSession = () => {
        if (!attempt.newSessionAccepted || !attempt.createdDialogueId || attempt.sessionInserted) return
        attempt.sessionInserted = true
        void insertNewSession(attempt.context, attempt.robotKey, attempt.createdDialogueId)
      }
      const acceptMessage = () => {
        settleSendResult(attempt, true)
        if (isNewChat.value) {
          isNewChat.value = false
          attempt.newSessionAccepted = true
        }
        insertAcceptedSession()
      }
      attempt.connection = sendAiMessage(params)
      sendStatus.value = 'streaming'
      attempt.connection.onMessage = (response) => {
        if (!isCurrentSend()) return
        try {
          const { event, data: content } = response
          if (event === 'dialogue_id') {
            dialogue_id.value = content
            attempt.createdDialogueId = normalizeCustomerServiceQueryValue(content)
            if (attempt.createdDialogueId === '0') attempt.createdDialogueId = ''
            insertAcceptedSession()
            return
          }
          if (event === 'c_message') {
            if (!attempt.userPushed) {
              const message = safeParseJson(content, null)
              if (!message || typeof message !== 'object' || Array.isArray(message)) throw new Error('Invalid customer message')
              pushUserMessage(message)
              attempt.userPushed = true
            }
            acceptMessage()
            if (cachedVariables) localStorage.removeItem(variablesKey)
            return
          }
          if (event === 'robot') {
            ensureAiMessage()
            acceptMessage()
            if (cachedVariables) localStorage.removeItem(variablesKey)
            return
          }
          if (event === 'chat_prompt_variables') {
            const variables = safeParseJson(content, {}) || {}
            chat_variables.value = {
              ...chat_variables.value,
              ...variables,
              fill_variables: variables.fill_variables || [],
              dialogue_id: variables.dialogue_id || dialogue_id.value
            }
            return
          }
          if (event === 'finish') {
            // finish只结束过程步骤，后面仍可能有最终内容/引用/耗时/变量，等onClose才释放连接。
            updateAiMessage(aiMessage, 'finalize_process_steps', { reason: 'finish' })
            return
          }
          if (event === 'error') {
            ensureAiMessage()
            // 业务error也是消息，沿Clawbot记录错误与过程状态，继续接收后续正文。
            updateAiMessage(aiMessage, 'error', content)
            return
          }
          const processEvents = ['reasoning_content', 'llm_rounds', 'stream_message', 'sending',
            'tool_call_full', 'tool_result', 'ai_message', 'FileOperation', 'ExecuteCommand']
          const extraEvents = ['reply_content_list', 'start_quote_file', 'quote_file', 'debug', 'recall_time', 'request_time']
          if (!processEvents.includes(event) && !extraEvents.includes(event)) return
          ensureAiMessage()
          if (event === 'FileOperation' || event === 'ExecuteCommand') {
            updateAiMessage(aiMessage, 'process_event', {
              eventName: event,
              rawData: typeof content === 'string' ? content : JSON.stringify(content || {})
            })
          } else if (event === 'tool_call_full' || event === 'ai_message') {
            updateAiMessage(aiMessage, event, safeParseJson(content, {}) || {})
          } else if (event === 'quote_file' || event === 'debug') {
            updateAiMessage(aiMessage, event, safeParseJson(content, []))
          } else {
            updateAiMessage(aiMessage, event, content)
          }
        } catch {
          finishSend(attempt, 'error', 'error_send')
        }
      }
      attempt.connection.onError = () => finishSend(attempt, 'error', 'error_send')
      attempt.connection.onClose = () => finishSend(attempt, 'close')
      return result
    } catch {
      if (isCurrentSend()) finishSend(attempt, 'error', 'error_send')
      return false
    }
  }

  /** 首轮按机器人key缓存；历史会话提交接口，失败保留原变量及调用方表单。 */
  const handleEditVariables = async ({ chat_prompt_variables: values } = {}) => {
    if (!isContextCurrent(contextGeneration) || chatStatus.value !== 'ready' || variablesLoading.value ||
      sendLock.value || !Array.isArray(values)) return null
    const context = contextGeneration
    const session = sessionGeneration
    const requestId = ++variableRequestSequence
    const requestDialogueId = dialogue_id.value
    const requestSessionId = chat_variables.value.session_id
    const variablesKey = `chat_prompt_variables_${contextRobotKey}`
    variablesLoading.value = true
    variablesError.value = ''
    try {
      // 捕获提交值，避免等待接口期间表单对象变化被误当作已保存内容。
      const encodedValues = JSON.stringify(values)
      const savedValues = JSON.parse(encodedValues)
      let response = true
      if (Number(requestDialogueId || 0) > 0) {
        response = await editVariables({
          robot_key: contextRobotKey,
          openid: openid.value,
          dialogue_id: requestDialogueId,
          session_id: requestSessionId,
          chat_prompt_variables: encodedValues
        })
      }
      if (!isSessionCurrent(context, session) || requestId !== variableRequestSequence ||
        chat_variables.value.session_id !== requestSessionId) return null
      if (Number(requestDialogueId || 0) > 0) localStorage.removeItem(variablesKey)
      else localStorage.setItem(variablesKey, encodedValues)
      chat_variables.value = { ...chat_variables.value, fill_variables: savedValues, need_fill_variable: false }
      return response
    } catch {
      if (isSessionCurrent(context, session) && requestId === variableRequestSequence) variablesError.value = 'error_variables'
      return null
    } finally {
      if (isSessionCurrent(context, session) && requestId === variableRequestSequence) variablesLoading.value = false
    }
  }

  return {
    user, robot, dialogue_id, rel_user_id, openid, messageList, lastPushedUserMessageUid,
    isNewChat, chat_variables, chatStatus, chatLoading, chatError,
    sendStatus, sendLock, sendError, sensitiveWords,
    myChatList, myChatListLoading, myChatListLoadCompleted, myChatListHasMore, myChatListError,
    chatMessageLoading, chatMessageLoadCompleted, chatMessageError, variablesLoading, variablesError,
    initialize, createChat, openChat, getMyChatList, onGetChatMessage,
    sendMessage, stopMessage, handleEditVariables, $reset
  }
})
