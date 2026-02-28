<template>
  <div class="pc-chat-layout" v-if="robot?.robot_name">
    <!-- 左侧边栏 -->
    <LeftSidebar
      @new-chat="handleNewChat"
      @delete="handleDelete"
      @select="handleSelectSession"
    />

    <!-- 右侧主内容区 -->
    <div class="main-content">
      <!-- 中央展示区域 -->
      <WelcomeDisplay
        v-if="showWelcomeDisplay"
        :avatar="robot.robot_avatar"
        :title="robot.robot_name"
        :description="robot.robot_intro"
        class="welcome-display"
      />

      <!-- 消息列表区域 -->
      <div class="messages-list-wrapper" v-else>
        <!-- 表单信息横幅 -->
        <div class="form-banner-top" v-if="isShowFromHeader">
          <div class="title">{{ tChat('label_form_info') }}</div>
          <div class="edit-block" @click="handleEditVariableForm">
            <img src="@/assets/icons/edit.svg">{{ tChat('btn_edit') }}
          </div>
        </div>
        <div class="messages-list-box">
          <MessageList
            ref="messageListRef"
            style="padding: 0 16px;"
            :messages="messageList"
            @scrollStart="onScrollStart"
            @scrollEnd="onScrollEnd"
            @scroll="onScroll"
          >
            <template v-for="(item, index) in messageList" :key="item.uid">
              <MessageItem
                class="message-item"
                :style="{maxWidth: listMaxWidthPx}"
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
          <transition name="slide-down">
            <div class="bottom-btn-box" @click="handleScrollToBottom" v-if="isShowBottomBtn">
              <svg-icon name="down-arrow" class="bottom-btn" />
            </div>
          </transition>
        </div>
      </div>

      <!-- 快捷命令 -->
      <div class="fast-command-wrapper" v-if="isShortcut">
        <FastCommand @send="handleSetMessageInputValue" :style="{maxWidth: listMaxWidthPx}" />
      </div>

      <!-- 底部输入框 -->
      <div class="chat-input-wrapper">
        <ChatInput
          class="chat-input-box"
          :style="{maxWidth: listMaxWidthPx}"
          v-model="inputMessage"
          v-model:fileList="fileList"
          :show-upload="showUpload"
          @send="handleSend"
        />
      </div>
    </div>
    <VariableModal ref="variableModalRef" />
  </div>
</template>

<script lang="ts" setup>
import type { Chat } from '@/stores/modules/chat'
import { checkChatRequestPermission } from '@/api/robot/index'
import { getUuid } from '@/utils/index.js'
import { storeToRefs } from 'pinia'
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { useChatStore } from '@/stores/modules/chat'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { useLocale } from '@/hooks/web/useLocale'
import { showToast } from 'vant'
import { useI18n } from '@/hooks/web/useI18n'
import LeftSidebar from './components/LeftSidebar.vue'
import WelcomeDisplay from './components/WelcomeDisplay.vue'
import ChatInput from './components/ChatInput.vue'
import MessageList from '@/views/chat/components/messages/message-list.vue'
import MessageItem from '@/views/chat/components/messages/message-item.vue'
import FastCommand from '@/views/chat/components/fast-comand/index.vue'

import VariableModal from '@/views/chat/components/variable-modal/index.vue'

const { t } = useI18n('views.pc.index')
const { t: tChat } = useI18n('views.chat.index')
const { changeLocale } = useLocale()
const localeStore = useLocaleStoreWithOut()


// 消息列表相关
interface MessageListComponent{
  scrollToMessage: (id: number | string) => void
  scrollToBottom: () => void
}

// 辅助函数：处理语言切换
const handleLanguageSwitch = () => {
  let newLang = 'zh-CN' // 默认语言

  // externalConfigH5 已经在 chat.ts 的 createChat 中解析完成
  if (externalConfigH5.value.lang) {
    newLang = externalConfigH5.value.lang
  }

  // 检查是否需要切换语言
  if (
    newLang &&
    newLang !== localeStore.currentLocale.lang &&
    ['zh-CN', 'en-US'].includes(newLang)
  ) {
    changeLocale(newLang as any)
  }
}

let isAllowedScrollToBottom = true

const route = useRoute()
// const userStore = useUserStore()
const chatStore = useChatStore()
const emitter = useEventBus()
const { on, close } = useIM()

const { sendMessage, getMyChatList, onGetChatMessage, $reset, openChat, createChat, getFastCommand } = chatStore
const { messageList, sendLock, externalConfigH5, robot, dialogue_id } = storeToRefs(chatStore)

// 强制使用新会话
const isForceNewChat = computed(() => route.query.isForceNewChat == '1')
// 是否显示退出登录按钮
// const isShowLogOut = computed(() => userStore.getLoginStatus)
// 退出登录按钮是否收起
const isScrolled = ref(false)
// 是否显示表单信息横幅
const isShowFromHeader = computed(() => {
  return !chatStore.chat_variables.need_fill_variable && chatStore.chat_variables.fill_variables && chatStore.chat_variables.fill_variables.length
}) 
// 
const listMaxWidth = computed(() => messageList.value.length ? 900 : 800)
// 
const listMaxWidthPx = computed(() => `${listMaxWidth.value}px`)

const messageListRef = ref<null | MessageListComponent>(null)
const isShowBottomBtn = ref(false)

// 滚动到指定消息
const scrollToMessageById = (id: number | string) => {
  if (messageListRef.value) {
    messageListRef.value.scrollToMessage(id)
  }
}

// 输入框
const inputMessage = ref('')
const fileList = ref([])
const checkChatRequestPermissionLoading = ref(false)
const sendLoading = computed(() => sendLock.value || checkChatRequestPermissionLoading.value)
const showUpload = computed(() => robot.value.question_multiple_switch == 1)
const showWelcomeDisplay = computed(() => messageList.value.length == 0)
const isShortcut = computed(() => robot.value.fast_command_switch == '1')

// 事件处理
const handleNewChat = async () => {
  isAllowedScrollToBottom = true
  inputMessage.value = ''

  const data = {
    openid: '',
    robot_key: robot.value.robot_key,
    avatar: '',
    name: '',
    nickname: '',
    dialogue_id: 0
  }

  scrollToBottom ()

  await createChat(data, false, isForceNewChat.value)
}

const handleDelete = () => {
  // console.log('清空所有记录')
}

const handleSelectSession = async (data: any) => {
  if (dialogue_id.value == data.id) {
    return
  }

  isAllowedScrollToBottom = true

  const query = route.query || {}

  const params = {
    robot_key: query.robot_key,
    openid: data.openid,
    dialogue_id: data.id
  }

  await openChat(params)

  const res = await onGetChatMessage()

  if (res) {
    scrollToBottom ()
  }
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

// 滚动到底部
const scrollToBottom  = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
    isShowBottomBtn.value = false
  }
}

// 点击按钮滚动到底部
const handleScrollToBottom = () => {
  isAllowedScrollToBottom = true
  scrollToBottom()
}



// 滚动事件
let lastScrollTop = 0
const onScroll = (event: any) => {
  if (event.scrollHeight - event.clientHeight > event.scrollTop) {
    // 不是在底部了，显示回到底部按钮
    isShowBottomBtn.value = true
  }

  if (lastScrollTop && lastScrollTop - event.scrollTop > 0) {
    isAllowedScrollToBottom = false
  }

  lastScrollTop = event.scrollTop

  // 滚动页面就收起
  isScrolled.value = true
}

// 滚动到顶部
const onScrollStart = async () => {
  // dialogue_id == 0加载消息列表没有意义
  if(dialogue_id.value == 0){
    return
  }

  isAllowedScrollToBottom = true
  const msgId = messageList.value[0]?.uid
  const res = await onGetChatMessage()
  if (res && msgId) {
    scrollToMessageById(msgId)
  }
}

// 滚动到底部
const onScrollEnd = () => {
  isShowBottomBtn.value = false
}

const sendTextMessage = (val: string) => {
  if (!val) {
    return
  }

  let query = route.query || {}

  sendMessage({
    message: val,
    global: JSON.stringify(query)
  })
}

const sendMultipleMessage = (messages: any[]) => {
  if (!messages.length) {
    return
  }

  const query = route.query || {}

  sendMessage({
    message: JSON.stringify(messages),
    global: JSON.stringify(query)
  })
}

const handleSend = async (message: string) => {
  if (sendLoading.value) {
    return
  }

  let text = message.trim()

  if (!text && !fileList.value.length) {
    return showToast(t('msg_input_required'))
  }

  checkChatRequestPermissionLoading.value = true

  try {
    //检查是否含有敏感词
    let result = await checkChatRequestPermission({
      robot_key: robot.value.robot_key,
      openid: robot.value.openid,
      question: text
    })

    checkChatRequestPermissionLoading.value = false

    if (result.data && result.data.words) {
      return showToast(t('msg_sensitive_content', { words: result.data.words.join(';') }))
    }

    isAllowedScrollToBottom = true

    if (showUpload.value && fileList.value.length > 0) {
      let messages: any[] = []

      if (text) {
        messages.push({
          type: 'text',
          uid: getUuid(16),
          text: text
        })
      }

      if (fileList.value.length) {
        fileList.value.forEach((file: any) => {
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

    inputMessage.value = ''
    fileList.value = []
  } catch (err) {
    showToast(t('msg_send_failed'))
    checkChatRequestPermissionLoading.value = false
  }
}

interface VariableModalRefState {
  handleEdit: (data?: any) => void
}
const variableModalRef = ref<null | VariableModalRefState>(null)
const handleEditVariableForm = () => {
  variableModalRef.value?.handleEdit()
}

// 快捷命令发送消息
const handleSetMessageInputValue = (data: string) => {
  if (!data) return
  isAllowedScrollToBottom = true
  sendTextMessage(data)
}

// 监听 updateAiMessage 触发消息列表滚动
const onUpdateAiMessage = () => {
  // AI消息更新时滚动到底部
  scrollToBottom ()
}

// 绑定事件监听
const bindEventListeners = () => {
  emitter.on('updateAiMessage', onUpdateAiMessage)
  on('message', onUpdateAiMessage)
}

const init = async () => {
  try {
    isAllowedScrollToBottom = true
    const query = route.query || {}
    const chatData: Chat = {
      openid: String(query.openid || ''),
      robot_key: String(query.robot_key || ''),
      avatar: String(query.avatar || ''),
      name: String(query.name || ''),
      nickname: String(query.nickname || ''),
      dialogue_id: Number(query.dialogue_id) || 0,
    }

    await createChat(chatData, false, isForceNewChat.value)

    // 在 createChat 完成后处理语言切换
    handleLanguageSwitch()

    // 如果不是强制新建对话，则获取对话记录
    if(!isForceNewChat.value){
      let res = await onGetChatMessage()
      if (res) {
        scrollToBottom ()
      }
    }

    // 获取对话记录
    getMyChatList(chatData.robot_key, chatData.openid)

    // 获取快捷命令列表
    getFastCommand()
  } catch (error) {
    console.error('初始化聊天失败:', error)
    showToast(t('msg_init_failed'))
  }
}

onMounted(async () => {

  await init()

  // 绑定事件监听
  bindEventListeners()
})

onUnmounted(() => {
  $reset()
  emitter.off('updateAiMessage', onUpdateAiMessage)
  close()
})
</script>

<style lang="less" scoped>
.pc-chat-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background-color: #fff;
}

/* 右侧主内容区 */
.main-content {
  flex: 1;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  background: #fff;
}

.welcome-display{
  margin-bottom: 16px;
}

/* 消息列表区域 */
.messages-list-wrapper {
  flex: 1;
  width: 100%;
  margin-bottom: 16px;
  display: flex;
  flex-flow: column nowrap;
  overflow: hidden;
  position: relative;
  .messages-list-box{
    flex: 1;
    overflow: hidden;
  }
  .message-item{
    max-width: 800px;
    padding: 48px 0 0 0;
  }
}

/* 回到底部按钮 */
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
  z-index: 10;

  .bottom-btn {
    font-size: 16px;
    color: #659efc;
  }
}

.bottom-btn-box:hover {
  border: 1px solid #659dfc;
}

.chat-input-wrapper{
  width: 100%;
  padding: 0 16px;
  .chat-input-box{
    margin: 0 auto;
  }
}

/* 快捷命令区域 */
.fast-command-wrapper {
  width: 100%;
  padding: 0 16px;
  margin-bottom: 12px;
  .fast-comand-container{
    max-width: 900px;
    margin: 0 auto;
    height: auto;
  }
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

.log-out {
  position: fixed;
  top: 68px;
  right: 0;
  transform: translateY(-50%);
  width: 104px;
  height: 40px;
  transition: all 0.3s ease;
  z-index: 100;
}

.log-out.scrolled {
  transform: translateY(-50%) translateX(84px); /* 露出20px (104px - 20px) */
}

/* 表单信息横幅 */
.form-banner-top {
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

  .title {
    font-weight: 600;
    color: #000000;
    font-size: 16px;
  }

  .edit-block {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #2475fc;
    font-size: 14px;
    cursor: pointer;
    padding: 1px 8px;

    img {
      width: 16px;
      height: 16px;
    }

    &:hover {
      background: var(--07, #E4E6EB);
      border-radius: 6px;
    }
  }
}
</style>
