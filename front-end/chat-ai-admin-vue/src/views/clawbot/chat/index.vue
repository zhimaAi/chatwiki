<style lang="less" scoped>
.chat-page {
  display: flex;
  height: 100vh;
  background: #fff;
  color: #262626;
}

.chat-left {
  width: 224px;
  display: flex;
  flex-direction: column;
  background: #fff;
  flex-shrink: 0;
  position: relative;

  &.is-wide {
    width: 300px;
  }

  &::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    width: 1px;
    height: 100%;
    background: #f0f0f0;
  }
}

.chat-toolbar {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  flex-shrink: 0;
  background: #fff;
  min-width: 0;

  .toolbar-title {
    display: flex;
    align-items: center;
    flex: 1;
    min-width: 0;
    gap: 8px;
    font-size: 16px;
    line-height: 24px;
    font-weight: 600;
    color: #262626;

    span {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.chat-history-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: #fff;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
  background: #fff;

  &.is-empty {
    align-items: center;
    justify-content: center;
    
    .chat-empty-state {
      flex: none;
    }

    .chat-input-area {
      width: min(100%, 660px);
      margin: 48px auto 110px;
      padding: 0;
    }
  }
}

.chat-messages {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  flex: 1;
  position: relative;
  overflow: hidden;
  padding: 0 16px 16px;
}

.chat-loading-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 12px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.chat-input-area {
  padding: 0 32px 0;
  flex-shrink: 0;

  .chat-input-tip {
    margin-top: 10px;
    font-size: 12px;
    line-height: 20px;
    color: #8c8c8c;
    text-align: center;
  }
}

// Empty state (no messages)
.chat-empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-empty-content {
  text-align: center;

  .empty-avatar-wrap {
    position: relative;
    width: 80px;
    height: 80px;
    margin: 0 auto 34px;
  }

  .empty-avatar {
    width: 80px;
    height: 80px;
    border-radius: 16px;
    margin: 0 auto;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    font-size: 40px;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .empty-title {
    font-size: 36px;
    line-height: 44px;
    font-weight: 600;
    color: #0a0a0a;
  }

  .empty-desc {
    margin-top: 24px;
    font-size: 16px;
    line-height: 24px;
    color: #595959;
    white-space: pre-wrap;
  }
}

</style>

<template>
  <div class="chat-page">
    <!-- 左侧对话列表 -->
      <div :class="['chat-left', { 'is-wide': isNonChineseLocale }]">
        <div class="chat-toolbar">
          <div class="toolbar-title">
            <svg-icon name="chat-plus" size="16"></svg-icon>
            <span>{{ t('title_conversation') }}</span>
          </div>
          <a-button type="primary" @click="openNewChat">
            <PlusOutlined /> {{ t('btn_new_chat') }}
          </a-button>
        </div>

        <div class="chat-history-wrapper">
          <ChatList
            :list="myChatList"
            :active="dialogue_id"
            @openChat="handleOpenChat"
            @onScrollEnd="onChatListScrollEnd"
          />
        </div>
      </div>

      <!-- 中间对话区 -->
      <div class="chat-main" :class="{ 'is-empty': showEmptyState }">
        <!-- 无消息时展示 Hero -->
        <div class="chat-empty-state" v-if="showEmptyState">
          <div class="chat-empty-content">
            <div class="empty-avatar-wrap">
              <div class="empty-avatar">
                <img v-if="currentAssistantAvatar" :src="currentAssistantAvatar" alt="" />
                <span v-else>🔥</span>
              </div>
            </div>
            <div class="empty-title">{{ emptyTitle }}</div>
            <div class="empty-desc">{{ emptyDescription }}</div>
          </div>
        </div>

        <div class="chat-loading-state" v-else-if="isChatSwitching">
          <a-spin size="large" />
          <span>{{ t('msg_loading_conversation') }}</span>
        </div>

        <!-- 有消息时展示消息列表 -->
        <template v-else>
          <ConversationBar
            v-if="showConversationKnowledgeBar"
            :title="activeChatTitle"
            :primaryLibrary="primaryLibrary"
            :extraLibraryCount="extraLibraryCount"
            :libraries="selectedLibraryRows"
            @openLibrary="openKnowledgeDetail"
          />

          <div class="chat-messages">
            <MessageList
              ref="messageListRef"
              :messages="messageList"
              :robotInfo="robotInfo"
              :reserveReplySpace="reserveReplySpace"
              @clickMsgMeun="onClickMsgMenu"
              @scroll="onMessageListScroll"
              @scrollStart="onScrollStart"
              @scrollEnd="onScrollEnd"
            />
            <ScrollBottomBtn
              :visible="showScrollBottomButton"
              :loading="sendLoading"
              @click="scrollMessageListToBottom"
            />
          </div>
        </template>

        <!-- 输入区 -->
        <div class="chat-input-area" v-if="!isChatSwitching">
          <MessageInput
            v-model:value="message"
            v-model:fileList="fileList"
            :loading="sendLoading"
            :modelLoading="saveRobotModelLoading"
            :modeId="robotInfo.model_config_id"
            :modeName="robotInfo.use_model"
            :showUpload="robotInfo.question_multiple_switch == 1"
            :emptyMode="showEmptyState"
            @send="onSendMessage"
            @changeModel="onChangeChatModel"
            @openPrompt="onOpenPromptDrawer"
            @openSkill="onOpenSkillDrawer"
          />
          <div class="chat-input-tip">{{ t('msg_ai_generated_tip') }}</div>
        </div>
      </div>

      <SkillDrawer
        :open="skillDrawerOpen"
        @close="onCloseSkillDrawer"
      />
      <PromptDrawer
        :open="promptDrawerOpen"
        :prompt="robotInfo.prompt || ''"
        :loading="savePromptLoading"
        @close="onClosePromptDrawer"
        @save="onSavePrompt"
      />
  </div>
</template>

<script setup>
import { generateRandomId } from '@/utils/index'
import { ref, computed, onMounted, watch, onUnmounted, nextTick } from 'vue'
import { storeToRefs } from 'pinia'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotChatStore } from '@/stores/modules/clawbot-chat'
import { useUserStore } from '@/stores/modules/user'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useLocaleStore } from '@/stores/modules/locale'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message as antMessage } from 'ant-design-vue'
import { getLibraryList } from '@/api/library'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
import { showErrorMsg } from '@/utils/index'
import MessageList from './components/message-list.vue'
import ScrollBottomBtn from './components/scroll-bottom-btn.vue'
import ChatList from './components/chat-list.vue'
import ConversationBar from './components/conversation-bar.vue'
import MessageInput from './components/message-input.vue'
import PromptDrawer from './components/prompt-drawer.vue'
import SkillDrawer from './components/skill-drawer.vue'
import { checkChatRequestPermission } from '@/api/robot/index'

const { t } = useI18n('views.clawbot.chat.index')
const chatStore = useClawbotChatStore()
const clawbotStore = useClawbotStore()
const localeStore = useLocaleStore()
const userStore = useUserStore()
const { currentLocale } = storeToRefs(localeStore)

const currentAssistant = computed(() => clawbotStore.currentAssistant)
const currentAssistantAvatar = computed(() => {
  return currentAssistant.value?.robot_avatar_url || currentAssistant.value?.robot_avatar || ''
})
const isNonChineseLocale = computed(() => {
  const lang = String(currentLocale.value?.lang || 'zh-CN')
  return !lang.startsWith('zh')
})

const {
  createChat,
  sendMessage,
  getMyChatList,
  openChat,
  onGetChatMessage,
  $reset
} = chatStore
const { messageList, sendLock, myChatList, robot, dialogue_id, lastPushedUserMessageUid } = storeToRefs(chatStore)

const robotInfo = computed(() => robot.value)
const emptyTitle = computed(() => {
  return robotInfo.value?.robot_name || currentAssistant.value?.robot_name || ''
})
const emptyDescription = computed(() => {
  return robotInfo.value?.robot_intro || currentAssistant.value?.robot_intro || ''
})

const isTruthyHidden = (value) => {
  if (value === true || value === 1) {
    return true
  }
  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase()
    return normalized === 'true' || normalized === '1'
  }
  return false
}

const isMessageHidden = (item) => {
  if (!item) {
    return false
  }
  if (item.visible === false) {
    return true
  }
  return isTruthyHidden(item.hide) || isTruthyHidden(item.hidden) || isTruthyHidden(item.is_hidden)
}

const hasMessages = computed(() => {
  return messageList.value.some((item) => !isMessageHidden(item))
})
const isChatSwitching = ref(false)
const showEmptyState = computed(() => !isChatSwitching.value && !hasMessages.value)

let isAllowedScrollToBottom = true
const message = ref('')
const fileList = ref([])
const libraryList = ref([])
const wxAppLibary = ref(null)
const messageListRef = ref(null)
const showScrollBottomButton = ref(false)
const reserveReplySpace = ref(false)
const checkChatRequestPermissionLoading = ref(false)
const saveRobotModelLoading = ref(false)
const savePromptLoading = ref(false)
const promptDrawerOpen = ref(false)
const skillDrawerOpen = ref(false)
const sendLoading = computed(() => sendLock.value || checkChatRequestPermissionLoading.value)
const relatedLibraryIds = computed(() => {
  return String(robotInfo.value?.library_ids || '')
    .split(',')
    .filter(Boolean)
    .map((item) => String(item))
})
const activeChatItem = computed(() => {
  return myChatList.value.find((item) => {
    const itemDialogueId = item.dialogue_id ?? item.id
    return String(itemDialogueId) === String(dialogue_id.value)
  }) || null
})
const activeChatTitle = computed(() => {
  return activeChatItem.value?.subject || activeChatItem.value?.last_chat_message || ''
})
const selectedLibraryRows = computed(() => {
  if (!relatedLibraryIds.value.length || !libraryList.value.length) {
    return []
  }

  const libraryMap = new Map()
  libraryList.value.forEach((item) => {
    const itemId = String(item.id)
    if (!wxAppLibary.value && item.type == 3) {
      return
    }
    libraryMap.set(itemId, item)
  })

  let rows = relatedLibraryIds.value
    .filter((itemId) => libraryMap.has(itemId))
    .map((itemId) => libraryMap.get(itemId))

  const defaultLibraryId = String(robotInfo.value?.default_library_id || '')
  if (defaultLibraryId) {
    const defaultIndex = rows.findIndex((item) => String(item.id) === defaultLibraryId)
    if (defaultIndex > 0) {
      const [defaultLibrary] = rows.splice(defaultIndex, 1)
      rows.unshift(defaultLibrary)
    }
  }

  return rows
})
const primaryLibrary = computed(() => selectedLibraryRows.value[0] || null)
const extraLibraryCount = computed(() => Math.max(selectedLibraryRows.value.length - 1, 0))
const showConversationKnowledgeBar = computed(() => {
  return !showEmptyState.value && !isChatSwitching.value && !!activeChatTitle.value && !!primaryLibrary.value
})

const loadKnowledgeOptions = async () => {
  try {
    const res = await getLibraryList({ type: '', show_open_docs: 1 })
    libraryList.value = res?.data || []
  } catch (err) {
    libraryList.value = []
    console.error('加载知识库列表失败', err)
  }
}

const loadWxLibraryStatus = async () => {
  try {
    const res = await getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' })
    const config = res?.data || {}
    wxAppLibary.value = config?.user_config?.switch_status == 1 ? config : null
  } catch (err) {
    wxAppLibary.value = null
    console.error('加载公众号知识库配置失败', err)
  }
}

const openKnowledgeDetail = (item) => {
  if (!item?.id) {
    return
  }

  window.open(`#/library/details/knowledge-document?id=${item.id}`, '_blank', 'noopener')
}

const onSendMessage = async () => {
  if (sendLoading.value) {
    return
  }

  if (!message.value && !fileList.value.length) {
    return showErrorMsg(t('msg_input_required'))
  }

  showScrollBottomButton.value = false
  checkChatRequestPermissionLoading.value = true

  try {
    let result = await checkChatRequestPermission({
      robot_key: robot.value.robot_key,
      openid: robot.value.openid,
      question: message.value,
      form_ids: robot.value.form_ids,
      dialogue_id: dialogue_id.value,
    })

    checkChatRequestPermissionLoading.value = false

    if (result.data && result.data.words) {
      return antMessage.error(t('msg_sensitive_words', { words: result.data.words.join(';') }))
    }

    let msg_type = 1
    let msg = message.value

    if (robotInfo.value.question_multiple_switch == 1) {
      let messageList = []
      msg_type = 99

      if (message.value) {
        messageList = [{
          type: 'text',
          uid: generateRandomId(16),
          text: message.value
        }]
      }

      if (fileList.value.length) {
        fileList.value.map((file) => {
          messageList.push({
            uid: file.uid,
            type: 'image_url',
            image_url: {
              url: file.url
            }
          })
        })
      }

      msg = JSON.stringify(messageList)
    }

    isAllowedScrollToBottom = false
    sendMessage({
      message: msg,
      global: JSON.stringify({}),
      msg_type: msg_type
    })

    message.value = ''
    fileList.value = []
  } catch (err) {
    checkChatRequestPermissionLoading.value = false
  }
}

const onChangeChatModel = async (model) => {
  if (!model || saveRobotModelLoading.value) {
    return
  }

  const sourceRobotInfo = clawbotStore.robotInfo || {}
  if (!sourceRobotInfo.id) {
    antMessage.error(t('msg_assistant_not_ready'))
    return
  }

  const nextModelConfigId = String(model.model_config_id || '')
  const nextUseModel = model.name || ''
  const previousModelConfigId = String(sourceRobotInfo.model_config_id || '')
  const previousUseModel = sourceRobotInfo.use_model || ''

  if (nextModelConfigId === previousModelConfigId && nextUseModel === previousUseModel) {
    return
  }

  saveRobotModelLoading.value = true

  try {
    await clawbotStore.saveAssistant({
      model_config_id: model.model_config_id,
      use_model: nextUseModel
    }, {
      optimistic: true,
      rollbackOnError: true,
      successMessage: t('msg_model_switched')
    })
  } catch (err) {
    // saveAssistant 已处理错误提示，这里只负责恢复按钮 loading。
  } finally {
    saveRobotModelLoading.value = false
  }
}

const onOpenPromptDrawer = () => {
  promptDrawerOpen.value = true
}

const onClosePromptDrawer = () => {
  promptDrawerOpen.value = false
}

const onOpenSkillDrawer = () => {
  skillDrawerOpen.value = true
}

const onCloseSkillDrawer = () => {
  skillDrawerOpen.value = false
}

const onSavePrompt = async (prompt) => {
  if (savePromptLoading.value) {
    return
  }

  savePromptLoading.value = true

  try {
    await clawbotStore.saveAssistant({
      prompt
    }, {
      optimistic: false,
      refreshAfterSave: true,
      successMessage: t('msg_saved')
    })
    promptDrawerOpen.value = false
  } catch (err) {
    // saveAssistant 已处理错误提示，这里只负责恢复保存状态。
  } finally {
    savePromptLoading.value = false
  }
}

const openNewChat = async () => {
  if (!currentAssistant.value) return

  isChatSwitching.value = false
  isAllowedScrollToBottom = true
  reserveReplySpace.value = false
  message.value = ''

  const assistant = currentAssistant.value
  const user_id = userStore.user_id

  let data = {
    robot_key: assistant.robot_key,
    avatar: userStore.avatar || '',
    name: userStore.user_name || '',
    nickname: userStore.user_name || '',
    is_background: 1,
    openid: user_id,
    dialogue_id: 0
  }

  resetScroll()

  await createChat(data)
}

const handleOpenChat = async (data) => {
  const currentDialogueId = data.dialogue_id || data.id

  if (dialogue_id.value == currentDialogueId) {
    return
  }

  isChatSwitching.value = true
  isAllowedScrollToBottom = true
  reserveReplySpace.value = false

  const assistant = currentAssistant.value

  let params = {
    robot_key: assistant.robot_key,
    avatar: data.avatar,
    name: data.name,
    nickname: data.nickname,
    is_background: data.is_background,
    openid: data.openid,
    dialogue_id: currentDialogueId
  }

  resetScroll()

  try {
    await openChat(params)

    let res = await onGetChatMessage()
    if (res) {
      isChatSwitching.value = false
      await nextTick()
      messageListScrollToBottom()
    }
  } finally {
    isChatSwitching.value = false
  }
}

const onClickMsgMenu = (text) => {
  isAllowedScrollToBottom = false

  sendMessage({
    message: text
  })
}

const onMessageListScroll = (scrollOption) => {
  showScrollBottomButton.value = !scrollOption.isAtBottom

  if (
    reserveReplySpace.value &&
    !sendLock.value &&
    scrollOption.scrollDirection === 'up' &&
    !scrollOption.isReplySpaceVisible
  ) {
    // 占位块离开视口后再清理，避免回答结束时直接回收高度造成页面抖动。
    reserveReplySpace.value = false
  }
}

const onScrollStart = async () => {
  isAllowedScrollToBottom = false

  let msgId = messageList.value[0].uid

  let res = await onGetChatMessage()

  if (res) {
    scrollToMessageById(msgId)
  }
}

const onScrollEnd = () => {
  showScrollBottomButton.value = false
}

const messageListScrollToBottom = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
    showScrollBottomButton.value = false
  }
}

const scrollMessageListToBottom = () => {
  if (messageListRef.value) {
    messageListRef.value.scrollToBottom()
    showScrollBottomButton.value = false
  }
}

const scrollToMessageById = (id, direction) => {
  if (messageListRef.value) {
    messageListRef.value.scrollToMessage(id, direction)
  }
}

const resetScroll = () => {
  if (messageListRef.value) {
    messageListRef.value.resetScroll()
  }
  reserveReplySpace.value = false
  showScrollBottomButton.value = false
}

const onChatListScrollEnd = () => {
  getMyChatList()
}

watch(
  lastPushedUserMessageUid,
  async (uid) => {
    if (!uid) {
      return
    }

    isAllowedScrollToBottom = false
    reserveReplySpace.value = true
    await nextTick()
    scrollToMessageById(uid, 'top')
  }
)

watch(
  () => currentAssistant.value?.id,
  async (assistantId) => {
    if (!assistantId) return

    await openNewChat()
    loadKnowledgeOptions()
    getMyChatList()
  },
  { immediate: true }
)

onMounted(() => {
  loadWxLibraryStatus()
})

onUnmounted(() => {
  $reset()
})
</script>
