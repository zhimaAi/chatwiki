<style lang="less" scoped>
.chat-between-box {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}
.left-menu-content {
  // width: 64px;
  height: 100%;
  overflow: hidden;
}
.chat-page {
  display: flex;
  flex-flow: column nowrap;
  height: 100vh;
  flex: 1;
  overflow: hidden;
  background: #fff;

  .chat-page-body {
    position: relative;
    margin: 0 auto;
    width: 100%;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-flow: column nowrap;

    .messages-list-wrap {
      flex: 1;
      overflow: hidden;
    }
    .open-chat-box{
      position: absolute;
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      right: 12px;
      top: 12px;
      background: #fff;
      border-radius: 8px;
      font-size: 24px;
    }
  }

  .fast-command-wrap {
    position: relative;
    padding-top: 5px;
    z-index: 2;
    background-color: #fff;
  }

  .technical-support-text {
    line-height: 20px;
    padding: 4px 0;
    font-size: 12px;
    color: #bfbfbf;
    text-align: center;
  }

  .bottom-btn-box {
    display: flex;
    width: 40px;
    height: 40px;
    padding: 12px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    border-radius: 40px;
    border: 1px solid #fff;
    background: #fff;
    box-shadow: 0 4px 16px 0 #0000001f;
    cursor: pointer;
    position: absolute;
    bottom: 20px;
    left: 50%;
    margin-left: -20px;
    .bottom-btn {
      font-size: 16px;
      color: #659efc;
    }
  }

  .bottom-btn-box:hover {
    border: 1px solid #659dfc;
  }

  /* 定义进入动画 */
  .slide-down-enter-active {
    animation: slide-down-in 0.3s ease-in;
    position: absolute;
    z-index: 1;
  }

  /* 定义进入完成后的状态 */
  .slide-down-enter-from {
    transform: translateY(150%);
  }

  /* 定义退出动画 */
  .slide-down-leave-active {
    animation: slide-down-out 0.3s ease-out;
    position: absolute;
    z-index: 1;
  }

  /* 定义退出完成后的状态 */
  .slide-down-leave-to {
    transform: translateY(150%);
  }

  @keyframes slide-down-in {
    from {
      transform: translateY(150%);
    }
    to {
      transform: translateY(0);
    }
  }

  @keyframes slide-down-out {
    from {
      transform: translateY(0);
    }
    to {
      transform: translateY(150%);
    }
  }
}


.form-banner-top{
  max-width: 736px;
  width: calc(100% - 24px);
  margin: 0 auto;
  margin-top: 12px;
  padding: 16px;
  border-radius: 12px;
  border: 2px solid var(--10, #FFF);
  height: 56px;
  background: linear-gradient(180deg, #EBF2FF 0%, #FFF 22.78%);
  box-shadow: 0 4px 24px 0 #00000014;
  display: flex;
  align-items: center;
  justify-content: space-between;
  .title{
    font-weight: 600;
    color: #000000;
    font-size: 16px;
  }
  .edit-block{
    display: flex;
    align-items: center;
    gap: 8px;
    color: #2475fc;
    font-size: 14px;
    cursor: pointer;
    padding: 1px 8px;
    img{
      width: 16px;
      height: 16px;
    }
    &:hover{
      background: var(--07, #E4E6EB);
      border-radius: 6px;
    }
  }
}

</style>

<template>
  <div class="chat-between-box">
    <div class="left-menu-content">
      <LeftSideBar
        @openChat="handleOpenChat"
        @openNewChat="openNewChat"
        :isMobileDevice="isMobileDevice"
        v-if="externalConfigPC.new_session_btn_show == 1"
        ref="leftSideBarRef"
      />
    </div>
    <div class="chat-page">
      <div class="chat-page-header">
        <ChatHeader />
        <div class="form-banner-top" v-if="isShowFromHeader">
          <div class="title">表单信息</div>
          <div class="edit-block" @click="handleEditVariableForm">
            <img src="@/assets/icons/edit.svg" alt="">编辑
          </div>
        </div>
      </div>
      <div class="chat-page-body">
        <div class="open-chat-box" @click="handleShowH5Chat" v-if="isMobileDevice && externalConfigPC.new_session_btn_show == 1">
          <svg-icon name="new-chat-btn"></svg-icon>
        </div>
        <div class="messages-list-wrap">
          <MessageList
            ref="messageListRef"
            :messages="messageList"
            @scrollStart="onScrollStart"
            @scrollEnd="onScrollEnd"
            @scroll="onScroll"
          >
            <template v-for="(item, index) in messageList" :key="item.uid">
              <MessageItem
                :index="index"
                :messageLength="messageList.length"
                :msg="item"
                :prevMsg="messageList[index - 1]"
                @sendTextMessage="sendTextMessage"
                @toggleReasonProcess="handleToggleReasonProcess"
                @toggleQuoteFiel="handleToggleQuoteFiel"
              />
            </template>
          </MessageList>
        </div>
        <transition name="slide-down">
          <div class="bottom-btn-box" @click="onScrollBottom" v-if="isShowBottomBtn">
            <svg-icon name="down-arrow" class="bottom-btn" />
          </div>
        </transition>
      </div>
      <div class="fast-command-wrap">
        <FastComand v-if="isShortcut" @send="handleSetMessageInputValue"></FastComand>
      </div>
      <div class="chat-page-footer">
        <MessageInput
          v-model:value="message"
          v-model:fileList="fileList"
          :loading="sendLoading"
          :showUpload="showUpload"
          @send="onSendMesage"
        />
        <div class="technical-support-text">{{ translate('由 ChatWiki 提供软件支持') }}</div>
      </div>
    </div>
    <VariableModal ref="variableModalRef" />
  </div>
</template>

<script setup lang="ts">
import type { Message } from './types'
import { getUuid } from '@/utils/index'
import { translate } from '@/utils/translate.js'
import { checkChatRequestPermission } from '@/api/robot/index'
import { postInit } from '@/event/postMessage'
import { ref, onMounted, onUnmounted, toRaw, computed } from 'vue'
import { showToast } from 'vant'
import { useWindowWidth } from './useWindowWidth'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { useChatStore } from '@/stores/modules/chat'
import ChatHeader from './components/chat-header.vue'
import MessageInput from './components/message-input.vue'
import MessageList from './components/messages/message-list.vue'
import MessageItem from './components/messages/message-item.vue'
import FastComand from './components/fast-comand/index.vue'
import LeftSideBar from '@/views/chat/components/left-side-bar/index.vue'
import VariableModal from './components/variable-modal/index.vue'

type MessageListComponent = {
  scrollToMessage: (id: number | string) => void
  scrollToBottom: () => void
}

interface LeftSideBarRefState {
  handleShowH5Chat: any
}
const { windowWidth } = useWindowWidth()
const isMobileDevice = computed(() => {
  return windowWidth.value <= 500
})

const isShowBottomBtn = ref(false)
const checkChatRequestPermissionLoding = ref(false)
const fileList = ref<any[]>([])
const message = ref('')

const emitter = useEventBus()
const { on } = useIM()
const chatStore = useChatStore()

const {
  sendMessage,
  onGetChatMessage,
  $reset,
  robot,
  externalConfigPC,
  openChatWindow,
  closeChatWindow,
  getMyChatList,
  openChat,
  createChat 
} = chatStore

const { messageList, sendLock, dialogue_id } = storeToRefs(chatStore)


const isShowFromHeader = computed(()=>{
  return !chatStore.chat_variables.need_fill_variable && chatStore.chat_variables.fill_variables && chatStore.chat_variables.fill_variables.length
})

const isShortcut = computed(()=>{
  return robot.yunpc_fast_command_switch == '1' ? true : false
})

const sendLoading = computed(() => sendLock.value || checkChatRequestPermissionLoding.value)

const showUpload = computed(() => {
  return robot.question_multiple_switch == 1
})

// 允许滚动到底部
let isAllowedScrollToBottom = true
let lastScrollTop = 0
const messageListRef = ref<null | MessageListComponent>(null)

const scrollToMessageById = (id: number | string) => {
  if (messageListRef.value) {
    messageListRef.value.scrollToMessage(id)
  }
}

const handleMessageListScrollToBottom = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
    isShowBottomBtn.value = false
  }
}

// 滚动
const onScroll = (event) => {
  if (event.scrollHeight - event.clientHeight > event.scrollTop) {
    // 不是在底部了，显示回到底部按钮
    isShowBottomBtn.value = true
  }

  if (lastScrollTop - event.scrollTop > 0) {
    isAllowedScrollToBottom = false
  }

  lastScrollTop = event.scrollTop
}

// 滚动到顶部
const onScrollStart = async () => {
  isAllowedScrollToBottom = true // 允许滚动到底部
  let msgId = messageList.value[0].uid

  let res = await onGetChatMessage()

  if (res) {
    scrollToMessageById(msgId)
  }
}

// 监听滚动到底部
const onScrollEnd = () => {
  isShowBottomBtn.value = false
  // console.log('滚动到底部')
}

// 回到底部
const onScrollBottom = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
    isShowBottomBtn.value = false
  }
}

// 通知sdk 初始化完毕
const SDKInit = (data: any) => {
  postInit(data)
}

const init = async () => {
  isAllowedScrollToBottom = true

  let res = await onGetChatMessage()

  if (res) {
    handleMessageListScrollToBottom()
  }
  // 通知sdk 初始化完毕
  SDKInit({ robot: toRaw(robot), config: toRaw(externalConfigPC) })
}

const sendTextMessage = (val: string) => {
  if (!val) {
    return showToast('请输入消息内容')
  }

  sendMessage({
    message: val
  })
}

const sendMultipleMessage = (messages: any[]) => {
  if (!messages.length) {
    return
  }

  sendMessage({
    message: JSON.stringify(messages)
  })
}

const onSendMesage = async (content: any) => {
  let text = content.trim()

  if (!text && !fileList.value.length) {
    return showToast('请输入消息内容')
  }

  //检查是否含有敏感词
  checkChatRequestPermissionLoding.value = true

  let result = await checkChatRequestPermission({
    from: 'yun_pc',
    robot_key: robot.robot_key,
    openid: robot.openid,
    question: text
  })

  checkChatRequestPermissionLoding.value = true

  if (result.data && result.data.words) {
    return showToast(`提交的内容包含敏感词：[${result.data.words.join(';')}] 请修改后再提交`)
  }

  isAllowedScrollToBottom = true

  if (showUpload.value) {
    let messages: Message[] = []

    if (text) {
      messages.push({
        type: 'text',
        uid: getUuid(16),
        text: text
      })
    }

    if (fileList.value.length) {
      fileList.value.map((file: any) => {
        messages.push({
          uid: file.uid,
          type: 'image_url',
          image_url: {
            url: file.url
          }
        })
      })
    }

    sendMultipleMessage(messages)
  } else {
    sendTextMessage(text)
  }

  message.value = ''
  fileList.value = []
}

// 监听 updateAiMessage 触发消息列表滚动
const onUpdateAiMessage = (msg) => {
  if (msg.event === 'reasoning_content') {
    return
  }

  if (messageListRef.value) {
    handleMessageListScrollToBottom()
  }
}

// 监听 打开窗口 触发消息列表滚动
const onOpenWindow = () => {
  openChatWindow()
  handleMessageListScrollToBottom()
}

const onCloseWindow = () => {
  closeChatWindow()
}

const handleSetMessageInputValue = (data: any) => {
  // 直接发出内容
  onSendMesage(data)
}

const handleToggleReasonProcess = (msgId: number) => {
  const msg = messageList.value.find((m) => {
    let id = m.message_id || m.id
    return id == msgId
  })

  if (msg) {
    msg.show_reasoning = !msg.show_reasoning
  }
}

const handleToggleQuoteFiel = (msgId: number) => {
  const msg = messageList.value.find((m) => {
    let id = m.message_id || m.id
    return id == msgId
  })
  if (msg) {
    msg.show_quote_file = !msg.show_quote_file
  }
}

const handleOpenChat = async (data : any) => {
  if (dialogue_id.value == data.id) {
    return
  }

  isAllowedScrollToBottom = true

  let params = {
    robot_key: robot.robot_key,
    openid: data.openid,
    dialogue_id: data.id
  }

  await openChat(params)

  let res = await onGetChatMessage()
  if (res) {
    onScrollBottom()
  }
}

const openNewChat = async () => {
  isAllowedScrollToBottom = true
  message.value = ''

  let data = {
    isOpen: false,
    unreadNumber: 0,
    openid: '',
    robot_key: robot.robot_key,
    avatar: '',
    name: '',
    nickname: '',
    dialogue_id: 0
  }

  onScrollBottom()

  await createChat(data)
}

const leftSideBarRef = ref<null | LeftSideBarRefState>(null)
const handleShowH5Chat = ()=>{
  leftSideBarRef.value && leftSideBarRef.value.handleShowH5Chat()
}

interface VariableModalRefState {
  handleEdit: (data?: any) => void
}
const variableModalRef = ref<null | VariableModalRefState>(null)
const handleEditVariableForm = () => {
  variableModalRef.value?.handleEdit()
}

onMounted(() => {
  init()

  // 获取对话记录
  getMyChatList()

  // 监听 updateAiMessage 触发消息列表滚动
  emitter.on('updateAiMessage', onUpdateAiMessage)

  // 监听im消息
  on('message', onUpdateAiMessage)

  emitter.on('openWindow', onOpenWindow)

  emitter.on('closeWindow', onCloseWindow)
})

onUnmounted(() => {
  $reset()
})
</script>
