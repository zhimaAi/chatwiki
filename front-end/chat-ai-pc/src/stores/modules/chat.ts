import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import { sendAiMessage, chatWelcome, getDialogueList, getChatMessage, getFastCommandList, addFeedback, delFeedback } from '@/api/chat'
import { editPrompt } from '@/api/robot/index'
import { getUuid, getOpenid } from '@/utils/index'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { DEFAULT_USER_AVATAR } from '@/constants/index'

export interface Message {
  dialogue_id: number
  openid: string
  msg_type: number
  is_customer: number
  loading: boolean
  isWelcome: boolean
  menu_json: any
  quote_file: any
  id: number
  uid: string
  avatar: string
  content: string
  debug: 0 | 1,
  feedback_type?: string,
  reasoning_content: string,
  reasoning_status: boolean,
  show_reasoning: boolean,
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
  yunpc_fast_command_switch: string
  id: number | null
  welcomes: Welcome,
  comand_list: any | string[]
  app_id: number,
  is_sending: boolean,
}

export interface PageStyle {
  headBackgroundColor: string
}

export interface ExternalConfigPc {
  headTitle: string
  headSubTitle: string
  headImage: string
  lang: string
  pageStyle: PageStyle
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
    yunpc_fast_command_switch: '',
    openid: '',
    welcomes: { content: '', question: [] },
    comand_list: [],
    app_id: -2, // webapp:-1,嵌入网站:-2
    is_sending: false, // 是否在发送中
  })
  // 样式配置
  const externalConfigPC = reactive<ExternalConfigPc>({
    headTitle: '',
    headSubTitle: 'Based on LLM, free and open-source.',
    headImage: '',
    lang: 'zh-CN',
    pageStyle: {
      headBackgroundColor: 'linear-gradient,to right,#2435E7,#01A0FB',
    }
  }) 

  // 创建对话
  const isNewChat = ref(false)
  const createChat = async (data: Chat) => {
    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }

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
    user.avatar = data.avatar || DEFAULT_USER_AVATAR
    user.name = data.name || ''
    user.nickname = data.nickname || ''

    const res = await chatWelcome({
      robot_key: robot.robot_key,
      openid: openid.value,
      nickname: user.nickname,
      name: user.name,
      // avatar: user.avatar || DEFAULT_USER_AVATAR
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
      robot.yunpc_fast_command_switch = robotInfo.yunpc_fast_command_switch
      robot.id = robotInfo.id

      if (robotInfo.welcomes) {
        robot.welcomes = JSON.parse(robotInfo.welcomes)
      }

      if (robotInfo.external_config_pc) {
        Object.assign(externalConfigPC, JSON.parse(robotInfo.external_config_pc))
      }else{
        externalConfigPC.headTitle = robotInfo.robot_name
        externalConfigPC.headImage = robotInfo.robot_avatar
      }

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

    if (msg && msg.dialogue_id == dialogue_id.value) {
      msg.uid = getUuid(32)
      msg.loading = false
      msg.isWelcome = true
      msg.avatar = robot.robot_avatar

      if (msg.menu_json && typeof msg.menu_json === 'string') {
        msg.menu_json = JSON.parse(msg.menu_json)
      }

      if (msg.quote_file && typeof msg.quote_file === 'string') {
        msg.quote_file = JSON.parse(msg.quote_file)
      }

      messageList.value.push(msg)
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
    msg.msg_type = 1
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

    if (type == 'reasoning_content') {
      const oldText = messageList.value[msgIndex].reasoning_content
      messageList.value[msgIndex].reasoning_content = oldText + content

      // 推理开始
      messageList.value[msgIndex].reasoning_status = true
      messageList.value[msgIndex].show_reasoning = true
    }

    if (type == 'sending') {
      // 推理结束
      messageList.value[msgIndex].reasoning_status = false

      const oldText = messageList.value[msgIndex].content
      messageList.value[msgIndex].content = oldText + content
    }

    if (type == 'quote_file') {
      messageList.value[msgIndex].quote_file = content.length > 0 ? content : []
    }

    if (type == 'ai_message') {
      if (content.menu_json && content.msg_type == 2) {
        const menu_json = JSON.parse(content.menu_json)
        messageList.value[msgIndex].content = menu_json.content
        messageList.value[msgIndex].menu_json = menu_json
      }
      messageList.value[msgIndex].id = content.id
      messageList.value[msgIndex].msg_type = content.msg_type // 更新真实的msg_type
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
      loading: true,
      id: '',
      content: '',
      reasoning_content: '',
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

    const params = {
      robot_key: robot.robot_key,
      openid: robot.openid,
      question: data.message,
      // prompt: robot.prompt,
      // library_ids: robot.library_ids,
      dialogue_id: dialogue_id.value
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
          isNewChat.value = false
        }
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

  // 获取对话记录
  const myChatListSize = 25
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
      robot_key: robot.robot_key
    })

    myChatListLoading.value = false

    const list = res.data || []

    if (list.length === 0) {
      myChatListLoadCompleted.value = true
      return false
    }

    myChatList.value = [...myChatList.value, ...list]

    return res
  }
  // 插入最新一条对话记录
  // const insertNewSession = () => {
  //   getDialogueList({
  //     min_id: 0,
  //     size: 1,
  //     robot_key: robot.robot_key
  //   }).then((res) => {
  //     const list = res.data || []

  //     if (list[0]) {
  //       myChatList.value.unshift(list[0])
  //     }
  //   })
  // }

  // 打开对话
  const openChat = async (data: Chat) => {
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

        if (item.is_customer == 1) {
          item.avatar = user.avatar
        } else {
          item.avatar = robot.robot_avatar
        }

        if (item.menu_json) {
          item.menu_json = JSON.parse(item.menu_json)
        }

        if (item.quote_file) {
          item.quote_file = []
        }
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
  }
})
