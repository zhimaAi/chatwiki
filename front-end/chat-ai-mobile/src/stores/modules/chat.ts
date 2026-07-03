import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import { sendAiMessage, chatWelcome, getDialogueList, getChatMessage, questionGuide, getFastCommandList, addFeedback, delFeedback, deleteDialogue, editDialogue, editVariables } from '@/api/chat'
import { editPrompt } from '@/api/robot/index'
import { getUuid, getOpenid, extractVoiceInfo, removeVoiceFormat } from '@/utils/index'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { useUserStore } from '@/stores/modules/user'
import { getCurrentConfig } from '@/utils/getLangConfig'
export interface Message {
  name: string
  nickname: string
  robot_avatar: string
  dialogue_id: number
  openid: string
  received_message_type?: string
  media_id_to_oss_url?: string
  msg_type: number | string
  is_customer: number
  loading: boolean
  isWelcome: boolean
  menu_json: any
  quote_file: any
  reply_content_list?: any
  id: number
  message_id: string
  uid: string
  avatar: string
  content: string
  debug: 0 | 1,
  guess_you_want: string[]
  question_tabkey: number
  feedback_type?: string,
  reasoning_content: string
  reasoning_status: boolean
  show_reasoning: boolean
  reasoning_expanded?: boolean
  process_steps?: ProcessStep[]
  current_round_index?: number
  active_thinking_step_id?: string
  quote_loading: boolean
  show_quote_file: boolean
  voice_content: any
  startLoading: boolean
  is_stopped?: boolean
  event?: string
  prevent_auto_scroll?: boolean
}

export interface ProcessStep {
  id: string
  type: 'thinking' | 'tool' | 'operation'
  title: string
  status: 'running' | 'done' | ''
  expanded: boolean
  hidden: boolean
  roundIndex: number
  contentText: string
  resultText: string
  eventName: string
  paramsText?: string
  tool_call_id?: string
}

export interface Chat {
  openid: string
  robot_key: string
  avatar: string
  name: string
  nickname: string
  dialogue_id: number
}

export interface Welcome {
  content: string
  question: any[]
}

export interface Robot {
  robot_key: string
  openid: string
  library_ids: string | string[]
  prompt: string
  robot_avatar: string
  robot_intro: string
  robot_name: string
  fast_command_switch: string
  id: number | null
  welcomes: Welcome
  enable_question_guide: boolean
  enable_common_question: boolean
  common_question_list: any | string[]
  comand_list: any | string[]
  app_id: number
  is_sending: boolean
  feedback_switch: boolean
  chat_type: any
  answer_source_switch: boolean
  application_type: string // 机器人类型 0普通 1工作流
  question_multiple_switch: number
  tips_before_answer_content: string
  tips_before_answer_switch: boolean
  multi_lang_configs: string
}

export interface PageStyle {
  navbarBackgroundColor: string
}

export interface ExternalConfigH5 {
  pageTitle: string
  logo: string
  lang: string
  navbarShow: number
  ai_generated_tip_show: number
  ai_generated_tip: string
  accessRestrictionsType: number
  pageStyle: PageStyle
  open_type: number
  window_width: number
  window_height: number
  new_session_btn_show: number // 显示新对话按钮
  avatarShow: number
}

export const useChatStore = defineStore('chat', () => {
  const emitter = useEventBus()
  const im = useIM()
  const messageList = ref<Message[]>([])
  let mySSE: any = null

  // 对话id
  const dialogue_id = ref(0)

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
  // 机器人的信息
  const robot = reactive<Robot>({
    id: null,
    library_ids: '',
    prompt: '',
    robot_avatar: '',
    robot_intro: '',
    robot_key: '',
    robot_name: '',
    fast_command_switch: '',
    openid: '',
    welcomes: { content: '', question: [] },
    enable_question_guide: false,
    common_question_list: [],
    enable_common_question: false,
    comand_list: [],
    app_id: -1, // webapp:-1,嵌入网站:-2
    is_sending: false, // 是否在发送中
    feedback_switch: false,
    chat_type: '',
    answer_source_switch: false,
    application_type: '0',
    question_multiple_switch: 0,
    tips_before_answer_content: '思考中、请稍候',
    tips_before_answer_switch: true,
    multi_lang_configs: '',
  })

  // 样式配置
  const externalConfigH5 = reactive<ExternalConfigH5>({
    pageTitle: 'ZHIMA CHATAI',
    lang: 'zh-CN',
    logo: '',
    navbarShow: 2,
    ai_generated_tip_show: 1,
    ai_generated_tip: '',
    accessRestrictionsType: 1,
    pageStyle: {
      navbarBackgroundColor: '#2475FC',
    },
    open_type: 1,
    window_width: 1200,
    window_height: 650,
    new_session_btn_show: 2,
    avatarShow: 1,
  })

  const getDefaultAvatarShow = (applicationType: any) => {
    return Number(applicationType) === 2 ? 2 : 1
  }

  const ensureAvatarShow = (config: ExternalConfigH5, applicationType: any, sourceConfig?: Partial<ExternalConfigH5>) => {
    const sourceAvatarShow = Number(sourceConfig?.avatarShow)
    if (sourceAvatarShow === 1 || sourceAvatarShow === 2) {
      config.avatarShow = sourceAvatarShow
    } else {
      config.avatarShow = getDefaultAvatarShow(applicationType)
    }
  }

  const getDefaultChatVariables = () => ({
    need_fill_variable: false,
    fill_variables: [],
    wait_variables: [],
    session_id: 0,
    dialogue_id: 0,
  })
  const chat_variables = ref<any>(getDefaultChatVariables())

  // 创建对话
  const isNewChat = ref(false)

  const setH5Config = async (data: Chat) => {

    openid.value = data.openid || getOpenid()

    robot.robot_key = data.robot_key
    robot.openid = openid.value

    user.openid = openid.value
    user.avatar = data.avatar || ''
    user.name = data.name || ''
    user.nickname = data.nickname || ''

    const res = await chatWelcome({
      robot_key: robot.robot_key,
      openid: openid.value,
      nickname: user.nickname,
      name: user.name,
      avatar: user.avatar,
      dialogue_id: 0,
    })

    try {
      const robotInfo = res.data.robot
      // 设置网页标题
      if(robotInfo.external_config_h5){
        const h5Config = JSON.parse(robotInfo.external_config_h5)
        Object.assign(externalConfigH5, h5Config)
        ensureAvatarShow(externalConfigH5, robotInfo.application_type, h5Config)
      }else{
        externalConfigH5.pageTitle = robotInfo.robot_name
        externalConfigH5.logo = robotInfo.robot_avatar
        ensureAvatarShow(externalConfigH5, robotInfo.application_type)
      }
      return res
    } catch (e) {
      Promise.reject(e)
    }
  }

  const createChat = async (data: Chat, autoInsertWelcomeMsg = true, _isForceNewChat = false) => {
    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }

    
    messageList.value = []
    // 重置聊天记录是否加载完成的状态
    chatMessageLoadCompleted.value = false
    chatMessageLoading.value = false
    sendLock.value = false
    chat_variables.value = getDefaultChatVariables()

    if (!data.dialogue_id) {
      isNewChat.value = true
      dialogue_id.value = 0
    } else {
      isNewChat.value = false
      dialogue_id.value = data.dialogue_id
    }

    openid.value = data.openid || getOpenid()

    robot.robot_key = data.robot_key
    robot.openid = openid.value

    user.openid = openid.value
    user.avatar = data.avatar || ''
    user.name = data.name || ''
    user.nickname = data.nickname || ''

    const res = await chatWelcome({
      robot_key: robot.robot_key,
      openid: openid.value,
      nickname: user.nickname,
      name: user.name,
      avatar: user.avatar,
      dialogue_id: dialogue_id.value,
      use_new_dialogue: isNewChat.value ? 1 : 0,
    })

    try {
      const userInfo = res.data.customer
      const robotInfo = res.data.robot

      user.admin_user_id = userInfo.admin_user_id
      user.avatar = userInfo.avatar
      user.id = userInfo.id
      user.name = userInfo.name
      user.nickname = userInfo.nickname
      robot.prompt = robotInfo.prompt
      robot.robot_avatar = robotInfo.robot_avatar
      robot.robot_intro = robotInfo.robot_intro
      robot.robot_key = robotInfo.robot_key
      robot.robot_name = robotInfo.robot_name
      robot.library_ids = robotInfo.library_ids
      robot.fast_command_switch = robotInfo.fast_command_switch
      robot.id = robotInfo.id
      robot.enable_question_guide = robotInfo.enable_question_guide == 'true';
      robot.enable_common_question = robotInfo.enable_common_question == 'true';
      robot.feedback_switch = robotInfo.feedback_switch == '1';
      robot.chat_type = robotInfo.chat_type;
      robot.answer_source_switch = robotInfo.answer_source_switch == 'true';
      robot.application_type = robotInfo.application_type
      robot.multi_lang_configs = robotInfo.multi_lang_configs

      robot.question_multiple_switch = Number(robotInfo.question_multiple_switch) || 0
      if (robotInfo.common_question_list) {
        robot.common_question_list = JSON.parse(robotInfo.common_question_list)
      }
      if (robotInfo.welcomes) {
        robot.welcomes = JSON.parse(robotInfo.welcomes)
      }

      // 插入欢迎语
      if(autoInsertWelcomeMsg){
        insertWelcomeMsg(res.data.message)
      }
      
      // 连接im
      im.connect(openid.value)
      im.on('message', onImMessage)

      // 设置网页标题
      if(robotInfo.external_config_h5){
        const h5Config = JSON.parse(robotInfo.external_config_h5)
        Object.assign(externalConfigH5, h5Config)
        ensureAvatarShow(externalConfigH5, robotInfo.application_type, h5Config)
      }else{
        externalConfigH5.pageTitle = robotInfo.robot_name
        externalConfigH5.logo = robotInfo.robot_avatar
        ensureAvatarShow(externalConfigH5, robotInfo.application_type)
      }

      document.title = externalConfigH5.pageTitle

      const faviconLink = document.querySelector('link[rel="icon"]');

      if(faviconLink && externalConfigH5.logo){
        faviconLink.setAttribute('href', externalConfigH5.logo);
      }

      let currentConfig = getCurrentConfig(robotInfo.multi_lang_configs)
      robot.tips_before_answer_content = currentConfig?.tips_before_answer_content || ''
      robot.tips_before_answer_switch = currentConfig?.tips_before_answer_switch == 'true';

      setTimeout(() => {
        const chatVariable = res.data.chat_variable || {}
        chat_variables.value = {
          ...getDefaultChatVariables(),
          ...chatVariable,
          session_id: Number(chatVariable.session_id || res.data.session_id || 0),
          dialogue_id: Number(chatVariable.dialogue_id || res.data.dialog_id || dialogue_id.value || 0),
          fill_variables: chatVariable.fill_variables || [],
          wait_variables: chatVariable.wait_variables || [],
        }
      })

      return res
    } catch (e) {
      Promise.reject(e)
    }
  }

  // 插入来自im的聊天记录
  const onImMessage = (msg: Message) => {
    if (import.meta.env.DEV) {
      msg.dialogue_id = dialogue_id.value
    }

    if (msg.msg_type == 'receiver_notify') {
      return
    }

    if (msg && msg.dialogue_id == dialogue_id.value) {
      msg.uid = getUuid(32)
      msg.loading = false
      msg.isWelcome = true
      msg.name = msg.name || msg.nickname
      if (msg.is_customer == 1) {
        msg.name = msg.name || user.name
        msg.avatar = msg.avatar || user.avatar
      } else {
        msg.name = msg.name || robot.robot_name
        msg.avatar = msg.avatar || robot.robot_avatar
      }

      if (msg.menu_json && typeof msg.menu_json === 'string') {
        msg.menu_json = JSON.parse(msg.menu_json)
      }

      if (msg.quote_file && typeof msg.quote_file === 'string') {
        msg.quote_file = JSON.parse(msg.quote_file)
      }
      if (msg.reply_content_list && typeof msg.reply_content_list === 'string') {
        try { msg.reply_content_list = JSON.parse(msg.reply_content_list) } catch (_) { msg.reply_content_list = [] }
      }
      messageList.value.push(msg)
    }
  }

  function checkIsPushWeclome(msg: Message){
    let menu_json = msg.menu_json
    let quote_file = msg.quote_file
    if(menu_json){
      menu_json = JSON.parse(menu_json)
    }
    if(quote_file){
      quote_file = JSON.parse(quote_file)
    }
    if(menu_json === '' && quote_file.length == 0){
      return true
    }
    if(quote_file.length == 0 && !menu_json.content && !menu_json?.question?.length){
      return true
    }
    return false
  }
  //  插入欢迎语
  const insertWelcomeMsg = (msg: Message) => {
    if (msg) {
      msg.uid = getUuid(32)
      msg.loading = false
      msg.isWelcome = true
      msg.avatar = robot.robot_avatar

      if(checkIsPushWeclome(msg)){
        return
      }
      if (msg.menu_json) {
        msg.menu_json = JSON.parse(msg.menu_json)
      }

      if (msg.quote_file) {
        msg.quote_file = JSON.parse(msg.quote_file)
      }

      messageList.value.push(msg)
    }
  }
  // 推送用户的消息到列表
  const pushUserMessage = (msg: Message) => {
    msg.uid = getUuid(32)
    msg.loading = false
    msg.avatar = user.avatar
    msg.openid = user.openid
    msg.msg_type = msg.msg_type || 1
    msg.is_customer = 1
    messageList.value.push(msg)
  }

  const pushAiMessage = (msg: any) => {
    messageList.value.push(msg)

    emitter.emit('updateAiMessage', msg)
  }

  const safeParseJson = (value: any, fallback: any) => {
    if (typeof value !== 'string') {
      return value ?? fallback
    }

    try {
      return JSON.parse(value)
    } catch {
      return fallback
    }
  }

  const markAutoScroll = (msg: Message, prevent: boolean) => {
    msg.prevent_auto_scroll = prevent
  }

  const ensureProcessStepState = (msg: Message) => {
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

  const createProcessStep = (step: Partial<ProcessStep> = {}): ProcessStep => ({
    id: step.id || getUuid(32),
    type: step.type || 'thinking',
    title: step.title || '',
    status: step.status || '',
    expanded: step.expanded === true,
    hidden: step.hidden === true,
    roundIndex: step.roundIndex || 0,
    contentText: step.contentText || '',
    resultText: step.resultText || '',
    eventName: step.eventName || '',
    paramsText: step.paramsText || '',
    tool_call_id: step.tool_call_id || '',
  })

  const appendProcessStep = (msg: Message, step: Partial<ProcessStep> = {}) => {
    ensureProcessStepState(msg)
    const nextStep = createProcessStep(step)
    msg.process_steps!.push(nextStep)
    return nextStep
  }

  const getProcessStepById = (msg: Message, stepId?: string) => {
    ensureProcessStepState(msg)
    return msg.process_steps!.find((item) => item.id === stepId)
  }

  const setRunningProcessStepsDone = (msg: Message) => {
    ensureProcessStepState(msg)
    msg.process_steps = msg.process_steps!
      .filter((step) => !(step.type === 'thinking' && step.status === 'running' && !step.contentText))
      .map((step) => {
        if (step.status !== 'running') {
          return step
        }
        return {
          ...step,
          status: 'done' as const,
          resultText: step.type === 'thinking' && !step.contentText ? '' : step.resultText,
        }
      })
    msg.active_thinking_step_id = ''
    msg.reasoning_status = false
  }

  const collapseCompletedProcessBlocks = (msg: Message) => {
    ensureProcessStepState(msg)
    msg.process_steps = msg.process_steps!.map((step) => ({
      ...step,
      expanded: false,
    }))
    msg.reasoning_expanded = false
  }

  // 更新AI的消息到列表
  const updateAiMessage = (type: string, content: any, uid: string) => {
    const msgIndex = messageList.value.findIndex((item) => item.uid == uid)
    if (msgIndex === -1) {
      return
    }
    const currentMessage = messageList.value[msgIndex]
    markAutoScroll(currentMessage, false)

    if (type == 'reply_content_list') {
      if (content !== undefined && typeof content === 'string') {
        currentMessage.reply_content_list = JSON.parse(content)
      }
    }

    if (type == 'reasoning_content') {
      const oldText = currentMessage.reasoning_content
      currentMessage.reasoning_content = oldText + content

      // 推理开始
      currentMessage.reasoning_status = true
      currentMessage.show_reasoning = true
      markAutoScroll(currentMessage, true)
    }

    if (type == 'sending') {

      // 开始生成中答案
      currentMessage.startLoading = false
      // 推理结束
      currentMessage.reasoning_status = false

      const oldText = currentMessage.content
      currentMessage.content = oldText + content

      currentMessage.voice_content = extractVoiceInfo(currentMessage.content)
      currentMessage.content = removeVoiceFormat(currentMessage.content)
      markAutoScroll(currentMessage, false)
    }

    if(type == 'start_quote_file'){
      currentMessage.quote_loading = true
    }

    if (type == 'quote_file') {
      currentMessage.quote_file = content.length > 0 ? content : []
      currentMessage.show_quote_file = true
      currentMessage.quote_loading = false
    }

    if (type == 'ai_message') {
      const hadContent = Boolean(currentMessage.content)
      const hasProcessSteps = Array.isArray(currentMessage.process_steps) && currentMessage.process_steps.length > 0
      if (content.menu_json && content.msg_type == 2) {
        const menu_json = JSON.parse(content.menu_json)
        // messageList.value[msgIndex].content = menu_json.content
        currentMessage.menu_json = menu_json
      }
      currentMessage.startLoading = false
      currentMessage.id = content.id
      currentMessage.msg_type = content.msg_type // 更新真实的msg_type

      currentMessage.content = content.content
      // 提取语音消息
      currentMessage.voice_content = extractVoiceInfo(currentMessage.content)
      currentMessage.content = removeVoiceFormat(currentMessage.content)

      if (content.reply_content_list !== undefined) {
        currentMessage.reply_content_list = content.reply_content_list
      }
      if (content.quote_file && typeof content.quote_file === 'string') {
        currentMessage.quote_file = JSON.parse(content.quote_file)
      }
      // 流式链路里 ai_message 是收尾同步，正文或过程已渲染过时不再回拉到底部。
      markAutoScroll(currentMessage, hadContent || hasProcessSteps)
    }

    if (type == 'debug') {
      currentMessage.debug = content.length > 0 ? content : []
      markAutoScroll(currentMessage, true)
    }

    if (type == 'llm_rounds') {
      ensureProcessStepState(currentMessage)
      if (content === 'begin') {
        currentMessage.current_round_index! += 1
      }

      if (content === 'finish') {
        const activeStep = getProcessStepById(currentMessage, currentMessage.active_thinking_step_id)
        if (activeStep) {
          if (!activeStep.contentText) {
            currentMessage.process_steps = currentMessage.process_steps!.filter((item) => item.id !== activeStep.id)
          } else {
            activeStep.status = 'done'
            activeStep.resultText = ''
          }
        }
        currentMessage.active_thinking_step_id = ''
        currentMessage.reasoning_status = false
      }
      markAutoScroll(currentMessage, true)
    }

    if (type == 'stream_message') {
      const nextData = safeParseJson(content, {})
      const reasoningText = nextData?.reasoning_content || ''
      const messageText = nextData?.content || ''

      if (reasoningText) {
        let activeStep = getProcessStepById(currentMessage, currentMessage.active_thinking_step_id)
        if (!activeStep) {
          ensureProcessStepState(currentMessage)
          currentMessage.current_round_index! += 1
          activeStep = appendProcessStep(currentMessage, {
            type: 'thinking',
            title: '思考过程',
            status: 'running',
            expanded: true,
            roundIndex: currentMessage.current_round_index || 0,
            resultText: '思考中...',
            eventName: 'stream_message',
          })
          currentMessage.active_thinking_step_id = activeStep.id
        }

        activeStep.contentText = `${activeStep.contentText || ''}${reasoningText}`
        activeStep.resultText = ''
        const oldText = currentMessage.reasoning_content || ''
        currentMessage.reasoning_content = oldText + reasoningText
        currentMessage.show_reasoning = true
        currentMessage.reasoning_status = true
      }

      if (messageText) {
        currentMessage.startLoading = false
        currentMessage.reasoning_status = false
        const oldText = currentMessage.content || ''
        currentMessage.content = oldText + messageText
        currentMessage.voice_content = extractVoiceInfo(currentMessage.content)
        currentMessage.content = removeVoiceFormat(currentMessage.content)
      }
      // Agent 新链路下，思考和正文都交给用户手动查看，不跟随 stream_message 自动滚动。
      markAutoScroll(currentMessage, true)
    }

    if (type == 'tool_call') {
      const nextData = safeParseJson(content, {})
      appendProcessStep(currentMessage, {
        type: 'tool',
        title: nextData?.name || nextData?.tool_name || 'tool',
        status: 'running',
        expanded: true,
        roundIndex: currentMessage.current_round_index || 0,
        eventName: 'tool_call',
      })
      currentMessage.show_reasoning = true
      markAutoScroll(currentMessage, true)
    }

    if (type == 'tool_call_full') {
      const functionInfo = content?.function || {}
      appendProcessStep(currentMessage, {
        type: 'tool',
        title: functionInfo?.name || 'tool',
        status: 'running',
        expanded: true,
        roundIndex: currentMessage.current_round_index || 0,
        paramsText: functionInfo?.arguments || '',
        tool_call_id: content?.id || '',
        eventName: 'tool_call_full',
      })
      currentMessage.show_reasoning = true
      markAutoScroll(currentMessage, true)
    }

    if (type == 'tool_result') {
      const nextData = safeParseJson(content, {})
      ensureProcessStepState(currentMessage)
      const toolCallId = nextData?.tool_call_id || ''
      const matchedStep = toolCallId
        ? currentMessage.process_steps!.find((step) => {
            return step.type === 'tool' && step.status === 'running' && step.tool_call_id === toolCallId
          })
        : [...currentMessage.process_steps!].reverse().find((step) => {
            return step.type === 'tool' && step.status === 'running' && (!nextData?.tool_name || step.title === nextData.tool_name)
          })

      if (matchedStep) {
        matchedStep.status = 'done'
        matchedStep.resultText = nextData?.content || ''
      } else {
        appendProcessStep(currentMessage, {
          type: 'tool',
          title: nextData?.tool_name || 'tool',
          status: 'done',
          expanded: true,
          roundIndex: currentMessage.current_round_index || 0,
          resultText: nextData?.content || '',
          tool_call_id: toolCallId,
          eventName: 'tool_result',
        })
      }
      markAutoScroll(currentMessage, true)
    }

    if (type == 'finalize_process_steps') {
      setRunningProcessStepsDone(currentMessage)
      collapseCompletedProcessBlocks(currentMessage)
      markAutoScroll(currentMessage, true)
    }
    if (type == 'guess_you_want') {
      // 猜你想问 插入
      messageList.value = messageList.value.map((item) => {
        return {
          ...item,
          question_tabkey: -1,
          guess_you_want: [],
        }
      })
      messageList.value[msgIndex].guess_you_want = content;
      messageList.value[msgIndex].question_tabkey = content.length > 0 ? 1 : 2;
      markAutoScroll(messageList.value[msgIndex], true)
    }
    if (type == 'set_question_tabkey') {
      messageList.value = messageList.value.map((item) => {
        return {
          ...item,
          question_tabkey: -1,
        }
      })
      // 猜你想问 常见问题的tabkey  1为 猜你想问 2为 常见问题
      messageList.value[msgIndex].question_tabkey = content;
      markAutoScroll(messageList.value[msgIndex], true)
    }
    emitter.emit('updateAiMessage', messageList.value[msgIndex])
  }

  // 关闭AI的消息加载状态
  const closeAiMessageLoading = () => {
    const msgIndex = messageList.value.length - 1

    if (!messageList.value[msgIndex]) {
      return
    }

    messageList.value[msgIndex].loading = false
  }

  // 发送消息
  const sendLock = ref(false)

  const getRunningAiMessageIndex = () => {
    for (let i = messageList.value.length - 1; i >= 0; i--) {
      const item = messageList.value[i]
      if (
        item.is_customer != 1 &&
        !item.is_stopped &&
        (item.loading || item.startLoading || item.quote_loading || item.reasoning_status)
      ) {
        return i
      }
    }

    return -1
  }

  const stopMessage = () => {
    const msgIndex = getRunningAiMessageIndex()
    if (msgIndex > -1) {
      messageList.value[msgIndex].is_stopped = true
      messageList.value[msgIndex].loading = false
      messageList.value[msgIndex].startLoading = false
      messageList.value[msgIndex].quote_loading = false
      messageList.value[msgIndex].reasoning_status = false
    }

    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }

    robot.is_sending = false
    sendLock.value = false
    closeAiMessageLoading()
  }

  const sendMessage = (data: any) => {
    if (sendLock.value) {
      return
    }

    const aiMsg = {
      startLoading: true, // // 对话开始状态
      loading: true,
      id: '',
      content: '',
      reasoning_content: '',
      reply_content_list : [],
      uid: getUuid(32),
      avatar: robot.robot_avatar,
      msg_type: 1,
      quote_file: [],
      is_customer: 0,
      debug: [],
      reasoning_status: false,
      show_reasoning: false,
      quote_loading: false,
      show_quote_file: true,
      voice_content: [],
      is_stopped: false,
      reasoning_expanded: false,
      process_steps: [],
      current_round_index: 0,
      active_thinking_step_id: '',
      event: 'robot',
      prevent_auto_scroll: false,
    }
    const userStore = useUserStore()
    const params: any = {
      robot_key: robot.robot_key,
      openid: robot.openid,
      question: data.message,
      // prompt: robot.prompt,
      // library_ids: robot.library_ids,
      dialogue_id: dialogue_id.value,
      global: data.global,
      rel_user_id: userStore.userInfo ? userStore.userInfo.user_id : '',
      use_new_dialogue: dialogue_id.value ? 0 : 1,
    }

    let variables_key = `chat_prompt_variables_${robot.robot_key}`

    const localVariables = localStorage.getItem(variables_key)
    const isNewDialogue = Number(dialogue_id.value || 0) === 0

    if (isNewDialogue && localVariables) {
      params.chat_prompt_variables = localVariables
      localStorage.removeItem(variables_key)
    }

    sendLock.value = true

    mySSE = sendAiMessage(params)

    mySSE.onMessage = (res) => {
      if (import.meta.env.MODE !== 'production') {
        console.log(res)
      }

      aiMsg.event = res.event;
      robot.is_sending = true;
      // 更新对话id
      if (res.event == 'dialogue_id') {
        dialogue_id.value = res.data
      }

      // 插入用户的问题到聊天记录
      if (res.event == 'c_message') {
        const data = JSON.parse(res.data)
        // 插入用户的问题到聊天记录
        pushUserMessage(data)
      }

      // 插入AI的回答到聊天记录
      if (res.event == 'robot') {
        // 插入AI的回答到聊天记录
        pushAiMessage(aiMsg)

        if (isNewChat.value) {
          // 插入新的对话记录
          insertNewSession()
          isNewChat.value = false
        }
      }
      
       // 更新功能中心回复
      if (res.event == 'reply_content_list') {
        updateAiMessage('reply_content_list', res.data, aiMsg.uid)
      }

      // 更新机器人深度思考的内容
      if (res.event == 'reasoning_content') {
        updateAiMessage('reasoning_content', res.data, aiMsg.uid)
      }

      if (res.event == 'llm_rounds') {
        updateAiMessage('llm_rounds', res.data, aiMsg.uid)
      }

      if (res.event == 'stream_message') {
        updateAiMessage('stream_message', res.data, aiMsg.uid)
      }

      // 更新机器人的消息
      if (res.event == 'sending') {
        updateAiMessage('sending', res.data, aiMsg.uid)
      }

      // tool_call 已废弃，使用 tool_call_full 替代
      // if (res.event == 'tool_call') {
      //   updateAiMessage('tool_call', res.data, aiMsg.uid)
      // }
      if (res.event == 'tool_call_full') {
        const data = safeParseJson(res.data, {})
        updateAiMessage('tool_call_full', data, aiMsg.uid)
      }

      if (res.event == 'tool_result') {
        updateAiMessage('tool_result', res.data, aiMsg.uid)
      }

      // 更新机器人消息的消息id时间等
      if (res.event == 'ai_message') {
        const data = JSON.parse(res.data)

        updateAiMessage('ai_message', data, aiMsg.uid)
      }

      if(res.event == 'start_quote_file'){
        updateAiMessage('start_quote_file', res.data, aiMsg.uid)
      }

      // 更新引用文件
      if (res.event == 'quote_file') {
        const data = JSON.parse(res.data)

        updateAiMessage('quote_file', data, aiMsg.uid)
      }

      // 更新prompt日志
      if (res.event == 'debug') {
        const data = JSON.parse(res.data)

        updateAiMessage('debug', data, aiMsg.uid)
      }

      if (res.event == 'chat_prompt_variables') {
        let data = res.data
        if(data){
          data = JSON.parse(data)
          chat_variables.value.need_fill_variable = data.need_fill_variable
          chat_variables.value.fill_variables = data.fill_variables || []
          chat_variables.value.session_id = data.session_id
          chat_variables.value.dialogue_id = data.dialogue_id
        }
      }

      if (res.event == 'finish') {
        robot.is_sending = false;
        if (robot.enable_question_guide) {
          // 相关问题开关开启了
          questionGuide({
            robot_key: robot.robot_key,
            openid: robot.openid,
            dialogue_id: dialogue_id.value,
          }).then(res => {
            updateAiMessage('guess_you_want', res.data || [], aiMsg.uid)
          })
        } else {
          updateAiMessage('set_question_tabkey', 2, aiMsg.uid)
        }
      }
    }

    mySSE.onClose = () => {
      if (aiMsg.process_steps.length) {
        updateAiMessage('finalize_process_steps', null, aiMsg.uid)
      }
      closeAiMessageLoading()
      sendLock.value = false

      mySSE = null
    }
  }

  const handleEditVariables = (data : any) => {
    return editVariables({
      robot_key: robot.robot_key,
      openid: robot.openid,
      dialogue_id: dialogue_id.value,
      chat_prompt_variables: JSON.stringify(data.chat_prompt_variables),
      session_id: chat_variables.value.session_id
    }).then((res)=>{
      chat_variables.value.fill_variables = data.chat_prompt_variables
      let variables_key = `chat_prompt_variables_${robot.robot_key}`
      localStorage.removeItem(variables_key)
      return res
    })
  }

  // 获取对话记录
  const myChatListSize = 35
  const myChatList = ref<any[]>([])
  const myChatListLoading = ref(false)
  const myChatListLoadCompleted = ref(false)

  const getMyChatList = async (robot_key?: string, openid?: string) => {
    if (myChatListLoadCompleted.value || myChatListLoading.value) {
      return false
    }

    let min_id = 0
    if (myChatList.value.length > 0) {
      min_id = myChatList.value[myChatList.value.length - 1].id
    }

    myChatListLoading.value = true

    const res = await getDialogueList({
      min_id: min_id,
      size: myChatListSize,
      robot_key: robot.robot_key || robot_key,
      openid: robot.openid || openid,
    })

    myChatListLoading.value = false

    const list = res.data || []

    if (list.length === 0) {
      myChatListLoadCompleted.value = true
      return false
    }

    myChatList.value = [...myChatList.value, ...list]

    // 根据 id 去重
    myChatList.value = myChatList.value.filter(
      (item, index, self) =>
        index === self.findIndex((t) => t.id === item.id)
    )

    return res
  }

  const delDialogue = async (data: any) => {
    const res = await deleteDialogue({
      robot_key: robot.robot_key,
      openid: robot.openid,
      ids: data.id,
    })
    if(data.id == -1){
      // 删除的是全部
      myChatList.value = []
    }else{
      myChatList.value = myChatList.value.filter(item => item.id != data.id)
    }
    if(data.id == -1 || data.id == dialogue_id.value){
      dialogue_id.value = 0;
      createChat({
        openid: '',
        robot_key: robot.robot_key,
        avatar: '',
        name: '',
        nickname: '',
        dialogue_id: 0
      })
    }
    return res
  }

  const editDialogueChat = async (data: any)=>{
    const res = await editDialogue({
      robot_key: robot.robot_key,
      openid: robot.openid,
      id: data.id,
      subject: data.subject
    })
    if(res.res == 0){
      // 将myChatList.value 中id为data.id的项的subject改为data.subject
      myChatList.value = myChatList.value.map(item => {
        if(item.id == data.id){
          item.subject = data.subject
        }
        return item
      })
    }
    return res
  }
  // 插入最新一条对话记录
  const insertNewSession = () => {
    getDialogueList({
      min_id: 0,
      size: 1,
      robot_key: robot.robot_key,
      openid: robot.openid,
    }).then((res) => {
      const list = res.data || []

      if (list[0]) {
        if(myChatList.value.filter(item => item.id === list[0].id).length === 0){
          myChatList.value.unshift(list[0])
        }
      }
    })
  }

  // 打开对话
  const openChat = async (data: any) => {
    const res = await createChat(data)

    return res
  }

  // 获取聊天记录
  const chatMessagePageSize = 20
  const chatMessageLoadCompleted = ref(false)
  const chatMessageLoading = ref(false)

  const onGetChatMessage = async () => {
    if (chatMessageLoadCompleted.value || chatMessageLoading.value) {
      return
    }

    let min_id = 0
    const list = messageList.value.filter((item) => !item.isWelcome)

    if (list.length > 0) {
      min_id = list[0].id
    }

    const params = {
      robot_key: robot.robot_key,
      openid: user.openid,
      min_id: min_id,
      size: chatMessagePageSize,
      dialogue_id: dialogue_id.value
    }

    chatMessageLoading.value = true

    try {
      const res = await getChatMessage(params)
      const list = res.data.list || []
      const _customer = res?.data?.customer || {}
      const _robot = res?.data?.robot || {}
      // 消息加载完了
      if (list.length === 0) {
        chatMessageLoadCompleted.value = true
        return
      }
      // 把消息倒过来
      list.sort((a, b) => {
        return a.id - b.id
      })
      list.forEach((item) => {
        item.loading = false
        item.uid = getUuid(32)
        item.show_reasoning = false;
        item.reasoning_status = false;
        item.reasoning_expanded = false;
        item.process_steps = Array.isArray(item.process_steps) ? item.process_steps : safeParseJson(item.process_steps, [])
        item.current_round_index = Number(item.current_round_index || 0)
        item.active_thinking_step_id = item.active_thinking_step_id || ''

        item.name = item.name || item.nickname
        if (item.is_customer == 1) {
          item.name = item.name || _customer.name
          item.avatar = item.avatar || user.avatar || _customer.avatar
        } else {
          item.name = item.name || robot.robot_name || _robot.robot_name
          item.avatar = item.avatar || robot.robot_avatar || _robot.robot_avatar
        }

        if (item.menu_json) {
          item.menu_json = JSON.parse(item.menu_json)
        }

        if (item.quote_file) {
          item.quote_file = JSON.parse(item.quote_file) || []
        }
        if (item.reply_content_list) {
          try { item.reply_content_list = JSON.parse(item.reply_content_list) } catch (_) { item.reply_content_list = [] }
        }

        item.voice_content = extractVoiceInfo(item.content)
        item.content = removeVoiceFormat(item.content)

      })

      messageList.value = [...list, ...messageList.value]
      return res
    } catch (err) {
      Promise.reject(err)
    } finally {
      chatMessageLoading.value = false
    }
  }

  // 点赞/点踩
  const onAddFeedback = async (data) => {

    const params = {
      robot_key: robot.robot_key,
      openid: user.openid,
      ai_message_id: data.ai_message_id,
      customer_message_id: data.customer_message_id,
      type: data.type,
      content: data.content
    }

    const res = await addFeedback(params)

    return res
  }

  // 取消点赞/点踩
  const onDelFeedback = async (data) => {

    const params = {
      robot_key: robot.robot_key,
      openid: user.openid,
      ai_message_id  : data.ai_message_id,
      customer_message_id: data.customer_message_id
    }

    const res = await delFeedback(params)

    return res
  }

  // 提示词
  const changeRobotPrompt = (text) => {
    robot.prompt = text
  }

  const saveRobotPrompt = () => {
    return editPrompt({
      id: robot.id,
      prompt: robot.prompt
    })
  }

  const getFastCommand = () => {
    if (!robot.robot_key || !robot.openid) {
      return Promise.resolve(false)
    }

    return getFastCommandList({
      robot_key: robot.robot_key,
      openid: robot.openid,
      app_id: robot.app_id
    }).then(res => {
      robot.comand_list = res.data;
      return res
    })
  }

  // 更新预览 ui
  const upDataUiStyle = (data) => {
    Object.assign(externalConfigH5, data)
    ensureAvatarShow(externalConfigH5, robot.application_type, data)
    let currentConfig = getCurrentConfig(robot.multi_lang_configs)
    robot.tips_before_answer_content = currentConfig?.tips_before_answer_content || ''
    robot.tips_before_answer_switch = currentConfig?.tips_before_answer_switch == 'true';
  }

  const updataQuickComand = (data) => {
    if (!robot.robot_key) {
      return
    }

   // 更新快捷指令
    robot.comand_list = data.comand_list || [];
    robot.fast_command_switch = data.fast_command_switch
  }


  function $reset() {
    // 关闭im
    im.close()

    dialogue_id.value = 0
    messageList.value = []
    openid.value = ''
    // 用户信息
    user.admin_user_id = ''
    user.avatar = ''
    user.id = ''
    user.name = ''
    user.nickname = ''
    user.openid = ''

    // 机器人的信息
    robot.id = null
    robot.library_ids = ''
    robot.prompt = ''
    robot.robot_avatar = ''
    robot.robot_intro = ''
    robot.robot_key = ''
    robot.robot_key = ''
    robot.robot_name = ''
    robot.fast_command_switch = ''
    robot.openid = ''
    robot.welcomes = { content: '', question: [] }
    // 是否是新的对话
    isNewChat.value = false
    // 消息加载完了
    chatMessageLoadCompleted.value = false
    chatMessageLoading.value = false
    // 是否正在发送消息
    sendLock.value = false
    // 对话记录
    myChatListLoading.value = false
    myChatListLoadCompleted.value = false
    myChatList.value = []
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
    createChat,
    sendMessage,
    stopMessage,
    getMyChatList,
    myChatList,
    openChat,
    onGetChatMessage,
    changeRobotPrompt,
    onAddFeedback,
    onDelFeedback,
    saveRobotPrompt,
    externalConfigH5,
    upDataUiStyle,
    getFastCommand,
    updataQuickComand,
    delDialogue,
    editDialogueChat,
    chat_variables,
    handleEditVariables,
    setH5Config,
  }
})
