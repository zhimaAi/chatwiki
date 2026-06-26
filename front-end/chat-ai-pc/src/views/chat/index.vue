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
          <div class="title">{{ t('label_form_info') }}</div>
          <div class="edit-block" @click="handleEditVariableForm">
            <img src="@/assets/icons/edit.svg" alt=""> {{ t('btn_edit') }}
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
        <ScrollToBottomBtn :visible="isShowBottomBtn" :loading="sendLoading" @click="onScrollBottom" />
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
          @stop="onStopMessage"
        />
        <div class="technical-support-text">
          <span class="ai-disclaimer-text">{{ t('label_ai_disclaimer') }}，</span>
          <span class="label_powered_by">{{ t('label_powered_by') }}</span>
        </div>
      </div>
    </div>
    <VariableModal ref="variableModalRef" />
  </div>
</template>

<script setup lang="ts">
import type { Message } from './types'
import type { Chat } from '@/stores/modules/chat'
import { getOpenid, getUuid } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'
import { useLocale } from '@/hooks/web/useLocale'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { checkChatRequestPermission } from '@/api/robot/index'
import { postInit } from '@/event/postMessage'
import { ref, onMounted, onUnmounted, toRaw, computed, watch } from 'vue'
import { showToast } from 'vant'
import { useRoute } from 'vue-router'
import { useWindowWidth } from './useWindowWidth'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { useChatStore } from '@/stores/modules/chat'
import ChatHeader from './components/chat-header.vue'
import MessageInput from './components/message-input.vue'
import ScrollToBottomBtn from './components/scroll-to-bottom-btn.vue'
import MessageList from './components/messages/message-list.vue'
import MessageItem from './components/messages/message-item.vue'
import FastComand from './components/fast-comand/index.vue'
import LeftSideBar from '@/views/chat/components/left-side-bar/index.vue'
import VariableModal from './components/variable-modal/index.vue'
import { getLang } from '@/utils/getLangConfig'

type MessageListComponent = {
  scrollToMessage: (id: number | string) => void
  scrollToBottom: () => void
}

// 初始化 i18n
const { changeLocale } = useLocale()
const localeStore = useLocaleStoreWithOut()
const { t } = useI18n('views.chat.index')

interface LeftSideBarRefState {
  handleShowH5Chat: any
}
const { windowWidth } = useWindowWidth()
const route = useRoute()

const getRouteChatData = (): Chat => {
  const query = route.query || {}

  return {
    isOpen: false,
    unreadNumber: 0,
    openid: String(query.openid || getOpenid()),
    robot_key: String(query.robot_key || ''),
    avatar: String(query.avatar || ''),
    name: String(query.name || ''),
    nickname: String(query.nickname || ''),
    dialogue_id: Number(query.dialogue_id) || 0
  }
}

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
  stopMessage,
  onGetChatMessage,
  $reset,
  robot,
  openChatWindow,
  closeChatWindow,
  getMyChatList,
  openChat,
  createChat 
} = chatStore

const { messageList, sendLock, dialogue_id, externalConfigPC } = storeToRefs(chatStore)


const isShowFromHeader = computed(()=>{
  return !chatStore.chat_variables.need_fill_variable && chatStore.chat_variables.fill_variables && chatStore.chat_variables.fill_variables.length
})

const isShortcut = computed(()=>{
  return robot.yunpc_fast_command_switch == '1' ? true : false
})

const sendLoading = computed(() => sendLock.value || checkChatRequestPermissionLoding.value)
let sendRequestSeq = 0

const showUpload = computed(() => {
  return robot.question_multiple_switch == 1
})

// 允许滚动到底部
let isAllowedScrollToBottom = true
const scrollEndDiff = 60
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
  const isAtBottom = Math.abs(event.scrollHeight - event.clientHeight - event.scrollTop) <= scrollEndDiff
  isShowBottomBtn.value = !isAtBottom
  // 只要用户离开底部，就关闭自动滚动；回到底部后再恢复。
  isAllowedScrollToBottom = isAtBottom
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
  isAllowedScrollToBottom = true

  if (messageListRef.value) {
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

  const chatData = getRouteChatData()

  await getMyChatList(chatData.robot_key, chatData.openid)

  const initialDialogueId = chatData.dialogue_id || Number(chatStore.myChatList?.[0]?.id) || 0

  await createChat({
    ...chatData,
    dialogue_id: initialDialogueId
  })

  let res = dialogue_id.value ? await onGetChatMessage() : null

  if (res) {
    handleMessageListScrollToBottom()
  }
  // 通知sdk 初始化完毕
  SDKInit({ robot: toRaw(robot), config: toRaw(externalConfigPC.value) })
}

const sendTextMessage = (val: string) => {
  if (!val) {
    return showToast(t('msg_input_required'))
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
  if (sendLoading.value) {
    return
  }

  let text = content.trim()

  if (!text && !fileList.value.length) {
    return showToast(t('msg_input_required'))
  }

  //检查是否含有敏感词
  checkChatRequestPermissionLoding.value = true
  const currentSendRequestSeq = ++sendRequestSeq

  let result
  try {
    result = await checkChatRequestPermission({
      from: 'yun_pc',
      robot_key: robot.robot_key,
      openid: robot.openid,
      question: text
    })
  } catch (error) {
    if (currentSendRequestSeq === sendRequestSeq) {
      checkChatRequestPermissionLoding.value = false
    }
    return showToast(t('msg_send_failed'))
  }

  if (currentSendRequestSeq !== sendRequestSeq) {
    return
  }

  checkChatRequestPermissionLoding.value = false

  if (result.data && result.data.words) {
    return showToast(t('msg_sensitive_word', { words: result.data.words.join(';') }))
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

const onStopMessage = () => {
  sendRequestSeq += 1
  checkChatRequestPermissionLoding.value = false
  stopMessage()
}

// 监听 updateAiMessage 触发消息列表滚动
const onUpdateAiMessage = (msg: any) => {
  if (msg?.prevent_auto_scroll) {
    return
  }

  if (messageListRef.value) {
    handleMessageListScrollToBottom()
  }
}

const onSocketMessage = () => {
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

// 监听 externalConfigH5.lang 变化，自动切换 i18n 语言
watch(
  () => chatStore.externalConfigPC.lang,
  async () => {
    let newLang = getLang()
    if(!['zh-CN', 'en-US'].includes(newLang)){
      newLang = 'en-US'
    }
    if (newLang && newLang !== localeStore.currentLocale.lang) {
      await changeLocale(newLang as LocaleType)

      // window.location.reload()
    }
  },
  { immediate: true }
)

onMounted(() => {
  init()

  // 监听 updateAiMessage 触发消息列表滚动
  emitter.on('updateAiMessage', onUpdateAiMessage)
  on('message', onSocketMessage)

  emitter.on('openWindow', onOpenWindow)

  emitter.on('closeWindow', onCloseWindow)
})

onUnmounted(() => {
  $reset()
})
</script>
