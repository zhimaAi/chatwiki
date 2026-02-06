<style lang="less" scoped>
.chat-between-box{
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}
.left-menu-content{
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
  background: #f0f2f5;

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
    background-color: #f0f2f5;
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
    <div class="left-menu-content" v-if="isCookieAccepted">
      <LeftSideBar 
        @openChat="handleOpenChat" 
        @openNewChat="openNewChat" 
        :isMobileDevice="isMobileDevice"
        v-if="externalConfigH5.new_session_btn_show == 1"
        ref="leftSideBarRef"
      />
    </div>
    <div class="chat-page" id="chatPage" v-if="isCookieAccepted">
      <div class="chat-page-header">
        <CuNavbar
          :title="externalConfigH5.pageTitle"
          :background-color="externalConfigH5.pageStyle.navbarBackgroundColor"
          v-if="externalConfigH5.navbarShow == 1"
        />
        <div class="form-banner-top" v-if="isShowFromHeader">
          <div class="title">{{ t('label_form_info') }}</div>
          <div class="edit-block" @click="handleEditVariableForm">
            <img src="@/assets/icons/edit.svg">{{ t('btn_edit') }}
          </div>
        </div>
      </div>
      <div class="chat-page-body">
        <div class="open-chat-box" @click="handleShowH5Chat" v-if="isMobileDevice && externalConfigH5.new_session_btn_show == 1">
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
          v-if="isMobileDevice"
          ref="messageInputRef"
          :show-upload="showUpload"
          v-model:value="message"
          v-model:fileList="fileList"
          @showLogin="onShowLogin"
          @send="onSendMesage"
          :loading="sendLoading"
        />
        <MessageInputPc
          v-else
          ref="messageInputRef"
          v-model:value="message"
          v-model:fileList="fileList"
          :show-upload="showUpload"
          @showLogin="onShowLogin"
          @send="onSendMesage"
          :loading="sendLoading"
        />

        <div class="technical-support-text">{{ t('tech_support') }}</div>
      </div>

      <LogOut
        v-if="isShowLogOut && externalConfigH5.accessRestrictionsType > 1"
        class="log-out"
        :class="{ scrolled: isScrolled }"
        @click="onTrigger"
      />
      <LoginModal ref="loginModalRef" />
      <VariableModal ref="variableModalRef" />
    </div>
    <CookieModal ref="cookieModalRef" :isMobileDevice="isMobileDevice" @onAccept="handleCookieDecision" @onDecline="handleCookieDecision" />
  </div>
</template>

<script setup lang="ts">
import type { Message } from './types'
import { getUuid } from '@/utils/index.js'
import { checkChatRequestPermission } from '@/api/robot/index'
import { useRoute } from 'vue-router'
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { showToast, showConfirmDialog } from 'vant'
import { storeToRefs } from 'pinia'
import { useWindowWidth } from './useWindowWidth'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { useChatStore } from '@/stores/modules/chat'
import { useUserStore } from '@/stores/modules/user'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { useI18n } from '@/hooks/web/useI18n'
import { useLocale } from '@/hooks/web/useLocale'
import CuNavbar from '@/components/cu-navbar/index.vue'
import MessageInput from './components/message-input.vue'
import MessageList from './components/messages/message-list.vue'
import MessageItem from './components/messages/message-item.vue'
import FastComand from './components/fast-comand/index.vue'
import LogOut from './components/log-out.vue'
import LoginModal from './components/login-modal.vue'
import MessageInputPc from './components/message-input-pc.vue'
import LeftSideBar from '@/views/chat/components/left-side-bar/index.vue'
import VariableModal from './components/variable-modal/index.vue'
import CookieModal from './components/cookie-modal.vue'

const { t } = useI18n('views.chat.index')
const { changeLocale } = useLocale()
const localeStore = useLocaleStoreWithOut()
import type { LocaleType } from '@/locales/config'

const { windowWidth } = useWindowWidth()
const route = useRoute()

type MessageListComponent = {
  scrollToMessage: (id: number | string) => void
  scrollToBottom: () => void
}

interface LoginModalRefState {
  show: any
}

interface LeftSideBarRefState {
  handleShowH5Chat: any
}

const isMobileDevice = computed(() => {
  return windowWidth.value <= 500
})

const userStore = useUserStore()
const isShowLogOut = computed(() => userStore.getLoginStatus)
const loginModalRef = ref<null | LoginModalRefState>(null)
const isScrolled = ref(false) // 退出登录按钮是否收起
const isShowBottomBtn = ref(false)

const emitter = useEventBus()
const { on } = useIM()
const chatStore = useChatStore()

const { sendMessage, getMyChatList, onGetChatMessage, $reset, robot, openChat, createChat } = chatStore

const { messageList, sendLock, externalConfigH5, dialogue_id } = storeToRefs(chatStore)

const isShowFromHeader = computed(()=>{
  return !chatStore.chat_variables.need_fill_variable && chatStore.chat_variables.fill_variables && chatStore.chat_variables.fill_variables.length
})

const fileList = ref([])
const message = ref('')
const checkChatRequestPermissionLoding = ref(false)
const sendLoading = computed(() => sendLock.value || checkChatRequestPermissionLoding.value)
const showUpload = computed(() => robot.question_multiple_switch == 1)

const isShortcut = computed(() => {
  return robot.fast_command_switch == '1' ? true : false
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

// 回到底部
const onScrollBottom = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
    isShowBottomBtn.value = false
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

  if (lastScrollTop && lastScrollTop - event.scrollTop > 0) {
    isAllowedScrollToBottom = false
  }

  lastScrollTop = event.scrollTop

  // 滚动页面就收起
  isScrolled.value = true
}

// 点击了退出按钮
const onTrigger = () => {
  if (isScrolled.value) {
    // 显示出退出按钮
    isScrolled.value = false
  } else {
    // 触发退出操作
    showConfirmDialog({
      title: t('title_warning'),
      message: t('msg_logout_confirm')
    })
      .then(() => {
        // on close
        userStore.reset()
      })
      .catch(() => {
        // on cancel
      })
  }
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
}

const init = async () => {
  isAllowedScrollToBottom = true

  let res = await onGetChatMessage()

  if (res) {
    handleMessageListScrollToBottom()
  }
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

  let query = route.query || {}

  sendMessage({
    message: JSON.stringify(messages),
    global: JSON.stringify(query)
  })
}

const onSendMesage = async () => {
  if (sendLoading.value) {
    return
  }

  let text = message.value.trim()

  if (!text && !fileList.value.length) {
    return showToast(t('msg_input_message'))
  }

  checkChatRequestPermissionLoding.value = true

  try {
    //检查是否含有敏感词
    let result = await checkChatRequestPermission({
      robot_key: robot.robot_key,
      openid: robot.openid,
      question: text
    })

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
  } catch (err) {
    checkChatRequestPermissionLoding.value = false
  }
}

// 监听 updateAiMessage 触发消息列表滚动
const onUpdateAiMessage = (msg) => {
  // if(msg.event === 'reasoning_content'){
  //   return
  // }

  if (messageListRef.value) {
    handleMessageListScrollToBottom()
  }
}

function setChatPageHeight() {
  // 适配移动端 高度为浏览器可视区域高度
  setTimeout(() => {
    document.getElementById('chatPage')!.style.height =
      document.documentElement.clientHeight - 1 + 'px'
  }, 20)
}

const handleSetMessageInputValue = (data: any) => {
  if (!data) {
    return
  }

  isAllowedScrollToBottom = true

  sendTextMessage(data)
}

const onShowLogin = () => {
  if (loginModalRef.value) {
    loginModalRef.value.show()
  }
}

const checkLogin = async () => {
  //检查是否含有敏感词
  let result = await checkChatRequestPermission({
    robot_key: robot.robot_key,
    openid: robot.openid,
    question: ''
  })

  checkChatRequestPermissionLoding.value = false

  // 未登录
  if (result.data?.code == 10002) {
    // 弹出登录
    showToast(t('msg_please_login'))
    onShowLogin()
    userStore.setLoginStatus(false)
    return false
  }

  // 有权限的账号登录后才可访问
  if (result.data?.code == 10003) {
    // 弹出登录
    showToast(t('msg_no_permission'))
    onShowLogin()
    userStore.setLoginStatus(false)
    return false
  }

  userStore.setLoginStatus(true)

  if (result.data && result.data.words) {
    return showToast(t('msg_sensitive_word', { words: result.data.words.join(';') }))
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

const handleOpenChat = async (data : any) => {
  if (dialogue_id.value == data.id) {
    return
  }

  isAllowedScrollToBottom = true

  let query = route.query || {}

  let params = {
    robot_key: query.robot_key,
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
  () => chatStore.externalConfigH5.lang,
  async (newLang) => {
    if (newLang && newLang !== localeStore.currentLocale.lang && ['zh-CN', 'en-US'].includes(newLang)) {
      await changeLocale(newLang as LocaleType)

      // window.location.reload()
    }
  },
  { immediate: true }
)

interface CookieModalRefState {
  show: () => void
}
const cookieModalRef = ref<null | CookieModalRefState>(null)
const isCookieAccepted = ref(false)

const handleCookieDecision = () => {
  isCookieAccepted.value = true
  initChatPage()
}

const initChatPage = async () => {
  if (!isCookieAccepted.value) {
    return
  }

  await checkLogin()
  init()
  getMyChatList()

  emitter.on('updateAiMessage', onUpdateAiMessage)
  on('message', onUpdateAiMessage)

  setChatPageHeight()
  window.addEventListener('resize', setChatPageHeight)
}


onMounted(async () => {
  cookieModalRef.value?.show()
})

onUnmounted(() => {
  $reset()
  emitter.off('updateAiMessage', onUpdateAiMessage)
  window.removeEventListener('resize', setChatPageHeight)
})
</script>
