import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import { sendAiMessage, chatWelcome, getDialogueList, getChatMessage, getFastCommandList, addFeedback, delFeedback, deleteDialogue, editDialogue, editVariables } from '@/api/chat'
import { editPrompt } from '@/api/robot/index'
import { getUuid, getOpenid, extractVoiceInfo, removeVoiceFormat } from '@/utils/index'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { postDot, postNewMessage } from '@/event/postMessage'
import { DEFAULT_SDK_FLOAT_AVATAR, DEFAULT_SDK_FLOAT_AVATAR2 } from '@/constants/index'

export interface Message {
  robot_name: any
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
  name: string
  nickname: string
  content: string
  debug: 0 | 1
  feedback_type?: string
  reasoning_content: string
  reasoning_status: boolean
  show_reasoning: boolean
  quote_loading: boolean
  show_quote_file: boolean
  voice_content: any
  startLoading: boolean
}

export interface Chat {
  isOpen: boolean
  unreadNumber: number
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
  yunpc_fast_command_switch: string
  id: number | null
  welcomes: Welcome
  comand_list: any | string[]
  app_id: number,
  is_sending: boolean,
  feedback_switch: boolean,
  question_multiple_switch: number,
  chat_type: any
  answer_source_switch: boolean
  application_type: string
  tips_before_answer_content: string
  tips_before_answer_switch: boolean
}

export interface PageStyle {
  headBackgroundColor: string
}
export interface FloatBtn {
  displayType: number
  buttonText: string
  buttonIcon: string
  bottomMargin: number
  rightMargin: number
  showUnreadCount: number
  showNewMessageTip: number
}

export interface ExternalConfigPc {
  headTitle: string
  headSubTitle: string
  headImage: string
  lang: string
  pageStyle: PageStyle
  floatBtn: FloatBtn
  open_type: number
  window_width: number
  window_height: number
  new_session_btn_show: number
}

const SDK_STATIC_HOST = window.location.origin

export const useChatStore = defineStore('chat', () => {
  const emitter = useEventBus()

  const isOpen = ref(false)

  const openChatWindow = () => {
    isOpen.value = true
    unreadNumber.value = 0;
    newMessageList.value = []
    postDot(unreadNumber.value)
    postNewMessage(newMessageList.value)
  }

  const closeChatWindow = () => {
    isOpen.value = false
  }

  const unreadNumber = ref(0)

  const im = useIM()

  const newMessageList = ref<Message[]>([])
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
    yunpc_fast_command_switch: '',
    openid: '',
    welcomes: { content: '', question: [] },
    comand_list: [],
    app_id: -2, // webapp:-1,嵌入网站:-2
    is_sending: false, // 是否在发送中
    feedback_switch: false,
    question_multiple_switch: 0,
    chat_type: '',
    answer_source_switch: false,
    application_type: '0',
    tips_before_answer_content: '思考中、请稍侯...',
    tips_before_answer_switch: true,
  })
  // 样式配置
  const externalConfigPC = reactive<ExternalConfigPc>({
    headTitle: '',
    headSubTitle: 'Based on LLM, free and open-source.',
    headImage: '',
    lang: 'zh-CN',
    pageStyle: {
      headBackgroundColor: 'linear-gradient,to right,#2435E7,#01A0FB',
    },
    floatBtn: {
      displayType: 1,
      buttonText: '快来聊聊吧~',
      buttonIcon: DEFAULT_SDK_FLOAT_AVATAR,
      bottomMargin: 32,
      rightMargin: 32,
      showUnreadCount: 1,
      showNewMessageTip: 1
    },
    open_type: 1,
    window_width: 1200,
    window_height: 650,
    new_session_btn_show: 2,
  })

  const chat_variables = ref<any>({
    need_fill_variable: false,
    fill_variables: [],
    wait_variables: [],
    session_id: 0,
    dialogue_id: 0,
  })

  let isFirstLoad = true
  // 创建对话
  const isNewChat = ref(false)
  const createChat = async (data: Chat) => {
    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }
    newMessageList.value = []
    messageList.value = []
    // 重置聊天记录是否加载完成的状态
    chatMessageLoadCompleted.value = false
    sendLock.value = false

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
      avatar: user.avatar
    })

    try {
      const userInfo = res.data.customer
      const robotInfo = res.data.robot

      if(isFirstLoad) {
        dialogue_id.value = res.data.dialog_id || 0
        isFirstLoad = false
      }

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
      robot.yunpc_fast_command_switch = robotInfo.yunpc_fast_command_switch
      robot.feedback_switch = robotInfo.feedback_switch == '1';
      robot.question_multiple_switch = robotInfo.question_multiple_switch;

      robot.chat_type = robotInfo.chat_type;
      robot.answer_source_switch = robotInfo.answer_source_switch == 'true';
      robot.application_type = robotInfo.application_type
      robot.tips_before_answer_content = robotInfo.tips_before_answer_content
      robot.tips_before_answer_switch = robotInfo.tips_before_answer_switch == 'true';

      robot.id = robotInfo.id

      if (robotInfo.welcomes) {
        robot.welcomes = JSON.parse(robotInfo.welcomes)
      }

      if (robotInfo.external_config_pc) {
        Object.assign(externalConfigPC, JSON.parse(robotInfo.external_config_pc))

        if(externalConfigPC.floatBtn.displayType == 1) {
          externalConfigPC.floatBtn.buttonIcon = DEFAULT_SDK_FLOAT_AVATAR
        }else if(externalConfigPC.floatBtn.displayType == 2) {
          externalConfigPC.floatBtn.buttonIcon = DEFAULT_SDK_FLOAT_AVATAR2
        }else if(externalConfigPC.floatBtn.displayType == 3) {
          externalConfigPC.floatBtn.buttonIcon = SDK_STATIC_HOST + externalConfigPC.floatBtn.buttonIcon
        }
      }else{
        externalConfigPC.headTitle = robotInfo.robot_name
        externalConfigPC.headImage = robotInfo.robot_avatar
      }

      chat_variables.value = {}
      
      setTimeout(()=>{
        chat_variables.value = res.data.chat_variable || {}
      })

      // 插入欢迎语
      insertWelcomeMsg(res.data.message)

      // 连接im
      im.connect(openid.value)
      im.on('message', onImMessage)

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

      if(!isOpen.value){
        const { showUnreadCount, showNewMessageTip } = externalConfigPC.floatBtn;

        if(!msg.robot_name){
          if (msg.avatar.indexOf('http') < 0) {
            msg.avatar = SDK_STATIC_HOST + msg.avatar
          }
          msg.robot_name = robot.robot_name
        }

        if(showUnreadCount == 1){
          unreadNumber.value ++;
         postDot(unreadNumber.value)
        }

        if(showNewMessageTip == 1){
          newMessageList.value.push(msg)
          postNewMessage(newMessageList.value)
        }
      }
    }
  }
  //  插入欢迎语
  const insertWelcomeMsg = (msg: Message) => {
    if (msg) {
      msg.uid = getUuid(32)
      msg.loading = false
      msg.isWelcome = true
      msg.avatar = robot.robot_avatar

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

  const pushAiMessage = (msg) => {
    messageList.value.push(msg)

    emitter.emit('updateAiMessage', msg)
  }

  // 更新AI的消息到列表
  const updateAiMessage = (type: string, content: any, uid: string) => {
    const msgIndex = messageList.value.findIndex((item) => item.uid == uid)

    if (type == 'reply_content_list') {
      if (content !== undefined && typeof content === 'string') {
        messageList.value[msgIndex].reply_content_list = JSON.parse(content)
      }
    }

    if (type == 'reasoning_content') {
      const oldText = messageList.value[msgIndex].reasoning_content
      messageList.value[msgIndex].reasoning_content = oldText + content

      // 推理开始
      messageList.value[msgIndex].reasoning_status = true
      messageList.value[msgIndex].show_reasoning = true
    }

    if (type == 'sending') {
      // 开始生成中答案
      messageList.value[msgIndex].startLoading = false
      // 推理结束
      messageList.value[msgIndex].reasoning_status = false

      const oldText = messageList.value[msgIndex].content
      messageList.value[msgIndex].content = oldText + content

      messageList.value[msgIndex].voice_content = extractVoiceInfo(messageList.value[msgIndex].content)
      messageList.value[msgIndex].content = removeVoiceFormat(messageList.value[msgIndex].content)
    }

    if(type == 'start_quote_file'){
      messageList.value[msgIndex].quote_loading = true
    }
    if (type == 'quote_file') {
      messageList.value[msgIndex].quote_file = content.length > 0 ? content : []
      messageList.value[msgIndex].show_quote_file = true
      messageList.value[msgIndex].quote_loading = false
    }

    if (type == 'ai_message') {
      if (content.menu_json && content.msg_type == 2) {
        const menu_json = JSON.parse(content.menu_json)
        // messageList.value[msgIndex].content = menu_json.content
        messageList.value[msgIndex].menu_json = menu_json
      }
      messageList.value[msgIndex].startLoading = false
      messageList.value[msgIndex].id = content.id
      messageList.value[msgIndex].msg_type = content.msg_type // 更新真实的msg_type

      messageList.value[msgIndex].content = content.content
      // 提取语音消息
      messageList.value[msgIndex].voice_content = extractVoiceInfo(messageList.value[msgIndex].content)
      messageList.value[msgIndex].content = removeVoiceFormat(messageList.value[msgIndex].content)

      if (content.reply_content_list !== undefined) {
        messageList.value[msgIndex].reply_content_list = content.reply_content_list
      }
      if (content.quote_file && typeof content.quote_file === 'string') {
        messageList.value[msgIndex].quote_file = JSON.parse(content.quote_file)
      }
    }

    if (type == 'debug') {
      messageList.value[msgIndex].debug = content.length > 0 ? content : []
    }

    emitter.emit('updateAiMessage', messageList.value[msgIndex])
  }

  // 关闭AI的消息加载状态
  const closeAiMessageLoading = () => {
    const msgIndex = messageList.value.length - 1

    messageList.value[msgIndex].loading = false
  }

  // 发送消息
  const sendLock = ref(false)

  const sendMessage = (data) => {
    if (sendLock.value) {
      return
    }

    const aiMsg = {
      startLoading: true, // // 对话开始状态
      loading: true,
      id: '',
      content: '',
      reasoning_content: '',
      reply_content_list: [],
      uid: getUuid(32),
      avatar: robot.robot_avatar,
      msg_type: 1,
      quote_file: [],
      is_customer: 0,
      debug: [],
      reasoning_status: false,
      show_reasoning: false,
      event: 'robot',
    }

    const params: any = {
      robot_key: robot.robot_key,
      openid: robot.openid,
      question: data.message,
      // prompt: robot.prompt,
      // library_ids: robot.library_ids,
      dialogue_id: dialogue_id.value,
      use_new_dialogue: externalConfigPC.new_session_btn_show == 1 ? 1 : 0,
    }

    let variables_key = `chat_prompt_variables_${robot.robot_key}`

    if(localStorage.getItem(variables_key)){
      params.chat_prompt_variables = localStorage.getItem(variables_key)
      localStorage.removeItem(variables_key)
    }

    sendLock.value = true

    mySSE = sendAiMessage(params)

    mySSE.onMessage = (res) => {
      if (import.meta.env.MODE !== 'production') {
        console.log(res)
      }
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

      // 更新机器人的消息
      if (res.event == 'sending') {
        updateAiMessage('sending', res.data, aiMsg.uid)
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
      }
    }

    mySSE.onClose = () => {
      closeAiMessageLoading()
      sendLock.value = false

      mySSE = null
    }
  }

  const handleEditVariables = (data : any) => {
    editVariables({
      robot_key: robot.robot_key,
      openid: robot.openid,
      dialogue_id: chat_variables.value.dialogue_id,
      chat_prompt_variables: JSON.stringify(data.chat_prompt_variables),
      session_id: chat_variables.value.session_id
    }).then(()=>{
      chat_variables.value.fill_variables = data.chat_prompt_variables
    })
  }

  // 获取对话记录
  const myChatListSize = 35
  const myChatList = ref<any[]>([])
  const myChatListLoading = ref(false)
  const myChatListLoadCompleted = ref(false)

  const getMyChatList = async () => {
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
      robot_key: robot.robot_key,
      openid: robot.openid,
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
        isOpen: false,
        unreadNumber: 0,
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

  const onGetChatMessage = async () => {
    if (chatMessageLoadCompleted.value) {
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
    getFastCommandList({
      robot_key: robot.robot_key,
      openid: robot.openid,
      app_id: robot.app_id
    }).then(res => {
      robot.comand_list = res.data;
    })
  }

  // 更新预览 ui
  const upDataUiStyle = (data) => {
    Object.assign(externalConfigPC, data)
  }

  const updataQuickComand = (data) => {
    // 更新快捷指令
     robot.comand_list = data.comand_list || [];
     robot.yunpc_fast_command_switch = data.fast_command_switch
   }

  function $reset() {
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
    robot.yunpc_fast_command_switch = ''
    robot.openid = ''
    robot.welcomes = { content: '', question: [] }
    // 是否是新的对话
    isNewChat.value = false
    // 消息加载完了
    chatMessageLoadCompleted.value = false
    // 是否正在发送消息
    sendLock.value = false
    // 对话记录
    myChatListLoading.value = false
    myChatListLoadCompleted.value = false
    myChatList.value = []
  }

  return {
    $reset,
    isOpen,
    openChatWindow,
    closeChatWindow,
    user,
    robot,
    dialogue_id,
    openid,
    sendLock,
    messageList,
    createChat,
    sendMessage,
    getMyChatList,
    myChatList,
    openChat,
    onGetChatMessage,
    changeRobotPrompt,
    onAddFeedback,
    onDelFeedback,
    saveRobotPrompt,
    externalConfigPC,
    upDataUiStyle,
    getFastCommand,
    updataQuickComand,
    delDialogue,
    editDialogueChat,
    chat_variables,
    handleEditVariables,
  }
})
