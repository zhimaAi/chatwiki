import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import {
  sendAiMessage,
  chatWelcome,
  getDialogueList,
  getChatMessage,
  questionGuide,
  getSessionRecordList,
  getSessionChannelList
} from '@/api/chat'
import { editPrompt } from '@/api/robot/index'
import { getUuid, getOpenid, devLog } from '@/utils/index'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { DEFAULT_USER_AVATAR } from '@/constants/index'

export const useChatStore = defineStore('chat', () => {
  const emitter = useEventBus()
  const messageList = ref([])
  const im = useIM()

  let mySSE = null

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
  // 机器人的信息
  const robot = reactive({
    id: '',
    library_ids: '',
    form_ids: '',
    prompt: '',
    robot_avatar: '',
    robot_intro: '',
    robot_key: '',
    robot_name: '',
    openid: '',
    welcomes: { reasoning_content: '', content: '', question: [] },
    enable_question_guide: false,
    enable_common_question: false,
    common_question_list: [],
    question_multiple_switch: 0,
  })

  // 存储聊天记录必备的数据
  const createMsg = async (data) => {
    messageList.value = []
    // 重置聊天记录是否加载完成的状态
    chatMessageLoadCompleted.value = false
    sendLock.value = false

    if (!data.dialogue_id) {
      dialogue_id.value = 0
    } else {
      dialogue_id.value = data.dialogue_id
    }

    rel_user_id.value = data.rel_user_id

    openid.value = data.openid || getOpenid(16)

    robot.robot_key = data.robot_key
    robot.openid = openid.value

    user.openid = openid.value
    user.avatar = data.avatar || DEFAULT_USER_AVATAR
    user.name = data.name || ''
    user.nickname = data.nickname || ''
  }

  // 创建对话
  const isNewChat = ref(false)
  const createChat = async (data) => {
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

    openid.value = data.openid || getOpenid(data.robot_key)

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
      // avatar: user.avatar || DEFAULT_USER_AVATAR, // lion要求把请求的avatar 干掉
      is_background: data.is_background || undefined
    })

    try {
      let userInfo = res.data.customer
      let robotInfo = res.data.robot
      user.admin_user_id = userInfo.admin_user_id
      user.avatar = userInfo.avatar
      user.id = userInfo.id
      user.name = userInfo.name
      user.nickname = userInfo.nickname
      robot.library_ids = robotInfo.library_ids
      robot.prompt = robotInfo.prompt
      robot.robot_avatar = robotInfo.robot_avatar
      robot.robot_intro = robotInfo.robot_intro
      robot.robot_key = robotInfo.robot_key
      robot.robot_name = robotInfo.robot_name
      robot.form_ids = robotInfo.form_ids
      robot.id = robotInfo.id
      robot.enable_question_guide = robotInfo.enable_question_guide == 'true'
      robot.enable_common_question = robotInfo.enable_common_question == 'true'
      robot.common_question_list = robotInfo.common_question_list
      robot.question_multiple_switch = Number(robotInfo.question_multiple_switch) || 0
      
      if (robotInfo.welcomes) {
        robot.welcomes = JSON.parse(robotInfo.welcomes)
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
  const onImMessage = (msg) => {
    if (!msg) {
      return
    }

    if (msg.msg_type == 'receiver_notify') {
      return
    }

    if (msg.dialogue_id != dialogue_id.value) {
      return
    }

    msg.uid = getUuid(32)
    msg.loading = false
    msg.isWelcome = true
    msg.name = msg.name || msg.nickname
    if (msg.is_customer == 1) {
      msg.name = msg.name || user.name
      msg.avatar = msg.avatar || user.avatar
    } else {
      msg.name = msg.name || robot.name
      msg.robot_avatar = msg.avatar || robot.robot_avatar
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

    emitter.emit('updateAiMessage', msg)
  }

  //  插入欢迎语
  const insertWelcomeMsg = (msg) => {
    if (msg) {
      msg.uid = getUuid(32)
      msg.loading = false
      msg.isWelcome = true
      msg.robot_avatar = robot.robot_avatar

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
  const pushUserMessage = (msg) => {
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
  const updateAiMessage = (type, content, uid) => {
    let msgIndex = messageList.value.findIndex((item) => item.uid == uid)

    if (type == 'reasoning_content') {
      let oldText = messageList.value[msgIndex].reasoning_content || ''
      messageList.value[msgIndex].reasoning_content = oldText + content

      // 推理开始
      messageList.value[msgIndex].reasoning_status = true
      messageList.value[msgIndex].show_reasoning = true
    }

    if (type == 'sending') {
      // 推理结束
      messageList.value[msgIndex].reasoning_status = false

      let oldText = messageList.value[msgIndex].content || ''
      // console.log('sending', content.length, content)
      messageList.value[msgIndex].content = oldText + content
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
      // 【ID1017779】【芝麻AI】机器人直连
      // 在实时聊天的时候，也需要把 ai_message 中的菜单内容显示出来
      if (content.menu_json && content.msg_type == 2) {
        let menu_json_obj = JSON.parse(content.menu_json)
        //messageList.value[msgIndex].content = menu_json_obj.content
        messageList.value[msgIndex].question = menu_json_obj.question
      }

      messageList.value[msgIndex].id = content.id

      messageList.value[msgIndex].content = content.content
      if (content.reply_content_list !== undefined && typeof content.reply_content_list === 'string') {
        messageList.value[msgIndex].reply_content_list = JSON.parse(content.reply_content_list)
      }
      if (content.quote_file && typeof content.quote_file === 'string') {
        messageList.value[msgIndex].quote_file = JSON.parse(content.quote_file)
      }
    }

    if (type == 'debug') {
      messageList.value[msgIndex].debug = content.length > 0 ? content : []
    }

    if (type == 'error') {
      messageList.value[msgIndex].error = content.length > 0 ? content : []
    }
    if (type == 'recall_time') {
      messageList.value[msgIndex].recall_time = content || 0
    }
    if (type == 'request_time') {
      messageList.value[msgIndex].request_time = content || 0
    }

    if (type == 'guess_you_want') {
      // 猜你想问 插入
      messageList.value = messageList.value.map((item) => {
        return {
          ...item,
          question_tabkey: -1,
          guess_you_want: []
        }
      })
      messageList.value[msgIndex].guess_you_want = content
      messageList.value[msgIndex].question_tabkey = content.length > 0 ? 1 : 2
    }
    if (type == 'set_question_tabkey') {
      messageList.value = messageList.value.map((item) => {
        return {
          ...item,
          question_tabkey: -1
        }
      })
      // 猜你想问 常见问题的tabkey  1为 猜你想问 2为 常见问题
      messageList.value[msgIndex].question_tabkey = content
    }

    emitter.emit('updateAiMessage', messageList.value[msgIndex])
  }

  // 关闭AI的消息加载状态
  const closeAiMessageLoading = () => {
    let msgIndex = messageList.value.length - 1

    messageList.value[msgIndex].loading = false
  }

  // 发送消息
  const sendLock = ref(false)

  const sendMessage = (data) => {
    if (sendLock.value) {
      return
    }

    let aiMsg = {
      loading: true,
      id: '',
      content: '',
      reasoning_content: '',
      uid: getUuid(32),
      robot_avatar: robot.robot_avatar,
      msg_type: 1,
      quote_file: [],
      debug: [],
      error: '',
      recall_time: '',
      request_time: '',
      show_reasoning: false,
      reasoning_status: false,
      quote_loading: false,
      show_quote_file: true,
    }
    let params = {
      robot_key: robot.robot_key,
      openid: robot.openid,
      question: data.message,
      // prompt: robot.prompt,
      // library_ids: robot.library_ids,
      form_ids: robot.form_ids,
      dialogue_id: dialogue_id.value,
      global: data.global
    }

    if (import.meta.env.DEV) {
      params.debug = 0
    }

    sendLock.value = true

    mySSE = sendAiMessage(params)

    mySSE.onMessage = (res) => {
      if (res.event !== 'sending') {
        devLog(res)
      }
      // 更新对话id
      if (res.event == 'dialogue_id') {
        dialogue_id.value = res.data
      }

      // 插入用户的问题到聊天记录
      if (res.event == 'c_message') {
        let data = JSON.parse(res.data)
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
        let data = JSON.parse(res.data)

        updateAiMessage('ai_message', data, aiMsg.uid)
      }

      if(res.event == 'start_quote_file'){
        updateAiMessage('start_quote_file', res.data, aiMsg.uid)
      }

      // 更新引用文件
      if (res.event == 'quote_file') {
        let data = JSON.parse(res.data)

        updateAiMessage('quote_file', data, aiMsg.uid)
      }

      // 更新prompt日志
      if (res.event == 'debug') {
        let data = JSON.parse(res.data)

        updateAiMessage('debug', data, aiMsg.uid)
      }

      // 更新prompt错误日志
      if (res.event == 'error') {
        let data = res.data

        updateAiMessage('error', data, aiMsg.uid)
      }
      // 更新prompt recall_time
      if (res.event == 'recall_time') {
        let data = res.data

        updateAiMessage('recall_time', data, aiMsg.uid)
      }
      // 更新prompt request_time
      if (res.event == 'request_time') {
        let data = res.data

        updateAiMessage('request_time', data, aiMsg.uid)
      }
      // 猜你想问
      if (res.event == 'finish') {
        if (robot.enable_question_guide) {
          // 相关问题开关开启了
          questionGuide({
            robot_key: robot.robot_key,
            openid: robot.openid,
            dialogue_id: dialogue_id.value
          }).then((res) => {
            updateAiMessage('guess_you_want', res.data || [], aiMsg.uid)
          })
        } else {
          updateAiMessage('set_question_tabkey', 2, aiMsg.uid)
        }
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
  const myChatList = ref([])
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
  const insertNewSession = () => {
    getDialogueList({
      min_id: 0,
      size: 1,
      robot_key: robot.robot_key
    }).then((res) => {
      const list = res.data || []

      if (list[0]) {
        myChatList.value.unshift(list[0])
      }
    })
  }

  // 打开对话
  const openChat = async (data) => {
    if (mySSE) {
      mySSE.abort()
      mySSE = null
    }
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
    let list = messageList.value.filter((item) => !item.isWelcome)

    if (list.length > 0) {
      min_id = list[0].id
    }

    const request_dialogue_id = dialogue_id.value;
    let params = {
      robot_key: robot.robot_key,
      openid: user.openid,
      min_id: min_id,
      size: chatMessagePageSize,
      dialogue_id: request_dialogue_id,
      rel_user_id: rel_user_id.value
    }

    try {
      const res = await getChatMessage(params)
      if (dialogue_id.value !== request_dialogue_id) {
        return Promise.resolve(null);
      }
      const list = res.data.list || []
      const _customer = res?.data?.customer || {}
      const newRobot = res.data.robot
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

        item.name = item.name || item.nickname
        if (item.is_customer == 1) {
          item.name = item.name || user.name || _customer.name
          item.avatar = item.avatar || user.avatar || _customer.avatar
        } else {
          item.name = item.name || robot.robot_name || newRobot.robot_name
          item.robot_avatar = item.avatar || robot.robot_avatar || newRobot.robot_avatar
          item.avatar = item.robot_avatar
        }

        if (item.menu_json) {
          item.menu_json = JSON.parse(item.menu_json)
        }

        if (item.quote_file) {
          item.quote_file = JSON.parse(item.quote_file)
        }

        if (item.reply_content_list && typeof item.reply_content_list === 'string') {
          try { item.reply_content_list = JSON.parse(item.reply_content_list) } catch (_) { item.reply_content_list = [] }
        }
      })

      messageList.value = [...list, ...messageList.value]
      return res
    } catch (err) {
      Promise.reject(err)
    }
  }

  const getRecordList = async (params) => {
    try {
      const res = await getSessionRecordList(params)
      if (!res) {
        return Promise.reject(res)
      }
      return res
    } catch (err) {
      Promise.reject(err)
    }
  }

  const getChannelList = async (params) => {
    try {
      const res = await getSessionChannelList(params)
      if (!res) {
        return Promise.reject(res)
      }
      return res
    } catch (err) {
      Promise.reject(err)
    }
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
    robot.id = ''
    robot.library_ids = ''
    robot.form_ids = ''
    robot.prompt = ''
    robot.robot_avatar = ''
    robot.robot_intro = ''
    robot.robot_key = ''
    robot.robot_name = ''
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
    createMsg,
    sendMessage,
    getMyChatList,
    myChatList,
    openChat,
    onGetChatMessage,
    changeRobotPrompt,
    saveRobotPrompt,
    getRecordList,
    getChannelList
  }
})
