<template>
  <div class="chat-page">
    <aside class="chat-left" :class="{ 'is-wide': isNonChineseLocale }">
      <div class="chat-toolbar">
        <div class="toolbar-title">
          <svg-icon name="chat-plus" size="16" />
          <span>{{ agentT('title_conversation') }}</span>
        </div>
        <!-- 历史会话加载时用 aria-disabled 保持按钮外观，点击仍由 openNewChat 的加载守卫拦截。 -->
        <a-button
          type="primary"
          :disabled="!robotReady || (chatLoading && Number(dialogue_id) === 0)"
          :aria-disabled="!robotReady || chatLoading"
          @click="openNewChat"
        >
          <PlusOutlined />{{ t('new_chat') }}
        </a-button>
      </div>
      <div class="chat-history-wrapper">
        <ChatList :key="identityVersion" :list="myChatList" :active="dialogue_id" @openChat="handleOpenChat" @onScrollEnd="loadMoreChats" />
        <div v-if="robotReady && !myChatList.length && !myChatListLoading && !myChatListError" class="empty-list">
          {{ t('empty_chat_list') }}
        </div>
      </div>
      <div v-if="myChatListLoading" class="list-loading"><a-spin size="small" /></div>
      <a-button v-else-if="robotReady && myChatListHasMore" type="text" @click="loadMoreChats">
        {{ t('load_more') }}
      </a-button>
    </aside>

    <section class="chat-main" :class="{ 'is-empty': showEmptyState }">
      <div v-if="!robotReady" class="chat-loading-state">
        <a-spin v-if="robotStatus === 'loading'" />
        <span>{{ t('waiting_robot') }}</span>
      </div>
      <template v-else>
        <div v-if="showVariableHeader" class="form-banner-top">
          <div class="title">{{ robotTestT('title_form_info') }}</div>
          <button class="edit-block" type="button" :disabled="sendLock || variablesLoading || chatLoading" @click="variableOpen = true">
            <EditOutlined />{{ robotTestT('btn_edit') }}
          </button>
        </div>
        <div v-if="chatLoading || (chatMessageLoading && !hasMessages)" class="chat-loading-state">
          <a-spin size="large" /><span>{{ t('loading_chat') }}</span>
        </div>
        <div v-else-if="chatStatus === 'error'" class="chat-loading-state">
          <span>{{ t('chat_unavailable') }}</span>
        </div>
        <template v-else-if="chatStatus === 'ready'">
          <div v-if="showEmptyState" class="chat-empty-state">
            <div class="chat-empty-content">
              <div class="empty-avatar-wrap">
                <div class="empty-avatar">
                  <img v-if="currentRobotAvatar" :src="currentRobotAvatar" alt="" />
                  <svg-icon v-else name="chat-plus" size="40" />
                </div>
              </div>
              <div class="empty-title">{{ robot.robot_name }}</div>
              <div v-if="robot.robot_intro" class="empty-desc">{{ robot.robot_intro }}</div>
            </div>
          </div>
          <template v-else>
            <ConversationBar
              v-if="activeChatTitle && primaryLibrary"
              :title="activeChatTitle"
              :primaryLibrary="primaryLibrary"
              :extraLibraryCount="Math.max(selectedLibraries.length - 1, 0)"
              :libraries="selectedLibraries"
              @openLibrary="openKnowledgeDetail"
            />
            <div class="chat-messages">
              <div v-if="chatMessageLoading" class="history-loading"><a-spin size="small" /></div>
              <MessageList
                :key="sessionVersion"
                ref="messageListRef"
                :messages="messageList"
                :robotInfo="robot"
                :reserveReplySpace="reserveReplySpace"
                @clickMsgMeun="onClickMsgMenu"
                @scroll="onMessageListScroll"
                @scrollStart="loadOlderMessages"
                @scrollEnd="showScrollBottomButton = false"
              />
              <ScrollBottomBtn :visible="showScrollBottomButton" :loading="sendLock" @click="scrollToBottom" />
            </div>
          </template>
          <div class="chat-input-area" :inert="chatMessageLoading || variablesLoading ? true : undefined">
            <MessageInput
              :key="sessionVersion"
              v-model:value="question"
              v-model:fileList="fileList"
              :loading="sendLock"
              :modelLoading="settingsPending"
              :modeId="robot.model_config_id"
              :modeName="robot.use_model"
              :showUpload="robot.question_multiple_switch == 1"
              :emptyMode="showEmptyState"
              @send="onSendMessage"
              @stop="chatStore.stopMessage"
              @changeModel="onChangeChatModel"
              @openPrompt="promptOpen = true"
              @openKnowledge="openCustomerKnowledge"
              @inputClick="openRequiredVariables"
            />
            <div class="chat-input-tip">{{ agentT('msg_ai_generated_tip') }}</div>
          </div>
        </template>
      </template>
    </section>
    <PromptDrawer
      :key="identityVersion"
      :open="promptOpen"
      :prompt="robot.prompt || ''"
      :robotKey="readyRobotKey"
      :loading="settingsPending"
      @close="promptOpen = false"
      @save="onSavePrompt"
      @variablesLoaded="onVariablesLoaded"
    />
    <VariableModal
      :key="sessionVersion"
      :open="variableOpen"
      :variables="chat_variables"
      :loading="variablesLoading"
      @submit="saveVariables"
      @close="variableOpen = false"
    />
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { App } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
import { useCustomerServiceChatStore } from '@/stores/modules/customer-service-chat'
import { useUserStore } from '@/stores/modules/user'
import { useLocaleStore } from '@/stores/modules/locale'
import { getLibraryList } from '@/api/library'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
import { generateRandomId } from '@/utils/index'
import { normalizeCustomerServiceQueryValue } from '../navigation'
import ChatList from '@/views/clawbot/chat/components/chat-list.vue'
import MessageList from '@/views/clawbot/chat/components/message-list.vue'
import ConversationBar from '@/views/clawbot/chat/components/conversation-bar.vue'
import ScrollBottomBtn from '@/views/clawbot/chat/components/scroll-bottom-btn.vue'
import MessageInput from './components/message-input.vue'
import PromptDrawer from './components/prompt-drawer.vue'
import VariableModal from './components/variable-modal.vue'

// 页面只编排UI/滚动/正式配置与知识导航；机器人身份归基础Store，会话协议归客服聊天Store。
const { t } = useI18n('views.customer-service.chat')
const { t: agentT } = useI18n('views.clawbot.chat.index')
const { t: robotTestT } = useI18n('views.robot.robot-test.index')
const { message: antMessage } = App.useApp()
const router = useRouter()
const customerServiceStore = useCustomerServiceStore()
const chatStore = useCustomerServiceChatStore()
const userStore = useUserStore()
const localeStore = useLocaleStore()
const { robotStatus } = storeToRefs(customerServiceStore)
const {
  robot, messageList, dialogue_id, lastPushedUserMessageUid, chat_variables,
  chatStatus, chatLoading, chatError, sendLock, sendError, sensitiveWords,
  myChatList, myChatListLoading, myChatListHasMore, myChatListError,
  chatMessageLoading, chatMessageError, variablesLoading, variablesError
} = storeToRefs(chatStore)

const readyRobotId = computed(() => robotStatus.value === 'ready'
  ? normalizeCustomerServiceQueryValue(customerServiceStore.robot?.id) : '')
const readyRobotKey = computed(() => robotStatus.value === 'ready' ? customerServiceStore.robot?.robot_key || '' : '')
const hasReadyRobot = () => customerServiceStore.robotStatus === 'ready' &&
  !!normalizeCustomerServiceQueryValue(customerServiceStore.robot?.id) && !!customerServiceStore.robot?.robot_key
const robotReady = computed(hasReadyRobot)
const isNonChineseLocale = computed(() => !String(localeStore.currentLocale?.lang || 'zh-CN').startsWith('zh'))
const currentRobotAvatar = computed(() => robot.value.robot_avatar_url || robot.value.robot_avatar || '')
const question = ref('')
const fileList = ref([])
const promptOpen = ref(false)
const variableOpen = ref(false)
const identityVersion = ref(0)
const sessionVersion = ref(0)
const messageListRef = ref(null)
const reserveReplySpace = ref(false)
const showScrollBottomButton = ref(false)
const savingSettings = ref('')
const settingsPending = computed(() => !!savingSettings.value)
const libraryList = ref([])
const wxLibraryEnabled = ref(false)
let knowledgeRequestSequence = 0
let draftVersion = 0
let sendRequestSequence = 0
let alive = true

const isCurrentIdentity = (version) => alive && hasReadyRobot() && version === identityVersion.value
const isCurrentSession = (version) => alive && hasReadyRobot() && version === sessionVersion.value

watch([question, fileList], () => { draftVersion++ }, { deep: true, flush: 'sync' })

// Store在每次请求前清空错误；同步订阅保证连续两次相同错误也能提示。
watch([chatError, myChatListError, chatMessageError, sendError], (errors, previousErrors) => {
  if (!alive || !hasReadyRobot()) return
  errors.forEach((error, index) => {
    if (error && error !== previousErrors[index]) {
      antMessage.error({ key: `customer-service-${error}`, content: t(error, { words: sensitiveWords.value.join(';') }) })
    }
  })
}, { flush: 'sync' })

const isMessageHidden = (item) => {
  if (item.visible === false) return true
  return [item.hide, item.hidden, item.is_hidden].some((value) =>
    value === true || value === 1 || ['true', '1'].includes(String(value).trim().toLowerCase()))
}
const hasMessages = computed(() => messageList.value.some((item) => !isMessageHidden(item)))
const showEmptyState = computed(() => robotReady.value && chatStatus.value === 'ready' &&
  !hasMessages.value && !chatMessageLoading.value)
const needsVariables = computed(() => chat_variables.value.need_fill_variable === true || chat_variables.value.need_fill_variable == 1)
const showVariableHeader = computed(() => chatStatus.value === 'ready' && !needsVariables.value &&
  (chat_variables.value.fill_variables || []).length > 0)

const activeChatTitle = computed(() => {
  const row = myChatList.value.find((item) => String(item.dialogue_id ?? item.id) === String(dialogue_id.value))
  return row?.subject || row?.last_chat_message || ''
})
const selectedLibraries = computed(() => {
  const ids = String(robot.value.library_ids || '').split(',').filter(Boolean)
  const available = new Map(libraryList.value.filter((item) => item.type != 3 || wxLibraryEnabled.value)
    .map((item) => [String(item.id), item]))
  const rows = ids.filter((id) => available.has(id)).map((id) => available.get(id))
  const defaultIndex = rows.findIndex((item) => String(item.id) === String(robot.value.default_library_id || ''))
  if (defaultIndex > 0) rows.unshift(...rows.splice(defaultIndex, 1))
  return rows
})
const primaryLibrary = computed(() => selectedLibraries.value[0] || null)

const clearConversationUI = () => {
  // 同一机器人切会话也要换输入实例，旧上传/校验/发送等待不能消费新会话草稿。
  sessionVersion.value++
  sendRequestSequence++
  question.value = ''
  fileList.value = []
  promptOpen.value = false
  variableOpen.value = false
  reserveReplySpace.value = false
  showScrollBottomButton.value = false
  messageListRef.value?.resetScroll()
}

const newChatParams = () => ({
  robot_key: customerServiceStore.robot?.robot_key,
  openid: userStore.user_id ?? undefined,
  avatar: userStore.avatar || '',
  name: userStore.user_name || '',
  nickname: userStore.user_name || '',
  dialogue_id: 0
})

const openNewChat = async () => {
  if (!robotReady.value || chatLoading.value) return
  chatStore.stopMessage()
  clearConversationUI()
  await chatStore.createChat(newChatParams())
}

const handleOpenChat = async (row) => {
  if (!robotReady.value || chatLoading.value) return
  const id = row.dialogue_id ?? row.id
  if (!normalizeCustomerServiceQueryValue(id)) return
  if (String(id) === String(dialogue_id.value) && chatStatus.value === 'ready') return
  chatStore.stopMessage()
  clearConversationUI()
  const version = sessionVersion.value
  const params = {
    robot_key: readyRobotKey.value, dialogue_id: id, openid: row.openid,
    avatar: row.avatar, name: row.name, nickname: row.nickname, rel_user_id: row.rel_user_id
  }
  // openChat已加载首批历史，不能再次请求同一页。
  const response = await chatStore.openChat(params)
  if (!response || !isCurrentSession(version)) return
  await nextTick()
  if (isCurrentSession(version)) messageListRef.value?.scrollToBottom()
}

const loadMoreChats = () => { void chatStore.getMyChatList() }

const loadOlderMessages = async (event) => {
  if (chatMessageLoading.value || sendLock.value) return
  const version = sessionVersion.value
  const anchor = event?.msg?.uid || messageList.value.find((item) => !isMessageHidden(item))?.uid
  const response = await chatStore.onGetChatMessage()
  if (!response || !isCurrentSession(version)) return
  await nextTick()
  if (!isCurrentSession(version)) return
  if (anchor) messageListRef.value?.scrollToMessage(anchor, 'top')
  else messageListRef.value?.scrollToBottom()
}

const onMessageListScroll = (scroll) => {
  showScrollBottomButton.value = !scroll.isAtBottom
  if (reserveReplySpace.value && !sendLock.value && scroll.scrollDirection === 'up' && !scroll.isReplySpaceVisible) {
    reserveReplySpace.value = false
  }
}
const scrollToBottom = () => {
  messageListRef.value?.scrollToBottom()
  showScrollBottomButton.value = false
}

const sendQuestion = async (text, files, consumeDraft) => {
  if (!robotReady.value || chatStatus.value !== 'ready' || sendLock.value || settingsPending.value ||
    chatMessageLoading.value || variablesLoading.value) return
  if (!text && !files.length) {
    antMessage.warning(t('empty_message'))
    return
  }
  if (files.some((file) => file.status === 'error')) {
    antMessage.warning(t('upload_failed'))
    return
  }
  if (files.some((file) => file.status !== 'done' || !file.url)) {
    antMessage.warning(t('upload_pending'))
    return
  }
  if (files.length && robot.value.question_multiple_switch != 1) {
    antMessage.warning(t('images_unavailable'))
    return
  }
  let content = text
  if (robot.value.question_multiple_switch == 1) {
    const parts = text ? [{ type: 'text', uid: generateRandomId(16), text }] : []
    files.forEach((file) => parts.push({ uid: file.uid, type: 'image_url', image_url: { url: file.url } }))
    content = JSON.stringify(parts)
  }
  const version = sessionVersion.value
  const submittedDraftVersion = draftVersion
  const requestId = ++sendRequestSequence
  try {
    const accepted = await chatStore.sendMessage({ message: content, checkQuestion: text, global: JSON.stringify({}) })
    if (!isCurrentSession(version) || requestId !== sendRequestSequence) return
    // 业务接纳才消费提交时草稿；用户在等待期间编辑/改附件则保留后来内容。
    if (accepted && consumeDraft && submittedDraftVersion === draftVersion) {
      question.value = ''
      fileList.value = []
    }
    if (!accepted && sendError.value === 'variables_required') variableOpen.value = true
  } catch {
    if (isCurrentSession(version) && requestId === sendRequestSequence) antMessage.error(t('error_send'))
  }
}

const onSendMessage = () => {
  void sendQuestion(question.value.trim(), fileList.value.map((item) => ({ ...item })), true)
}
const onClickMsgMenu = (text) => {
  // 菜单与输入共用发送流程，但菜单选择不覆盖/消费正在编辑的草稿和附件，失败可再次点击原选项。
  void sendQuestion(String(text || '').trim(), [], false)
}

const saveSettings = async (patch, kind) => {
  if (!robotReady.value || settingsPending.value) return
  const version = identityVersion.value
  savingSettings.value = kind
  try {
    const result = await customerServiceStore.saveRobotSettings(patch)
    if (!isCurrentIdentity(version)) return
    if (kind === 'prompt') promptOpen.value = false
    if (result.refreshed) antMessage.success(t('settings_saved'))
    else antMessage.warning(t('settings_saved_refresh_failed'))
  } catch {
    if (isCurrentIdentity(version)) {
      antMessage.error(t('error_settings_save'))
    }
  } finally {
    if (isCurrentIdentity(version)) savingSettings.value = ''
  }
}
const onChangeChatModel = (model) => {
  if (!model) return
  if (String(model.model_config_id) === String(robot.value.model_config_id) && model.name === robot.value.use_model) return
  void saveSettings({ model_config_id: model.model_config_id, use_model: model.name || '' }, 'model')
}
const onSavePrompt = (prompt) => { void saveSettings({ prompt }, 'prompt') }

const onVariablesLoaded = (list) => {
  if (!robotReady.value || !promptOpen.value || chatStatus.value !== 'ready') return
  const variables = chat_variables.value
  chatStore.chat_variables = {
    ...variables, wait_variables: list,
    need_fill_variable: list.length > 0 && (variables.fill_variables || []).length === 0
  }
}
const saveVariables = async (values) => {
  if (sendLock.value) {
    antMessage.warning(t('wait_response_for_variables'))
    return
  }
  const version = sessionVersion.value
  const response = await chatStore.handleEditVariables({ chat_prompt_variables: values })
  if (!isCurrentSession(version)) return
  if (response !== null) variableOpen.value = false
  else if (variablesError.value) antMessage.error(t(variablesError.value))
}

const loadKnowledgeOptions = async () => {
  if (!hasReadyRobot()) return
  const version = identityVersion.value
  const requestId = ++knowledgeRequestSequence
  const results = await Promise.allSettled([
    getLibraryList({ type: '', show_open_docs: 1 }),
    getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' })
  ])
  if (!isCurrentIdentity(version) || requestId !== knowledgeRequestSequence) return
  const [libraries, ability] = results
  const librariesLoaded = libraries.status === 'fulfilled' && Array.isArray(libraries.value?.data)
  if (librariesLoaded) libraryList.value = libraries.value.data
  if (ability.status === 'fulfilled') wxLibraryEnabled.value = ability.value?.data?.user_config?.switch_status == 1
  if (!librariesLoaded || ability.status === 'rejected') antMessage.warning(t('error_knowledge'))
}
const openKnowledgeDetail = (item) => {
  if (!item?.id) return
  const target = router.resolve({ path: '/library/details/knowledge-document', query: { id: item.id } })
  window.open(target.href, '_blank', 'noopener,noreferrer')
}
const openCustomerKnowledge = async () => {
  if (!robotReady.value) return
  const version = identityVersion.value
  try {
    await router.push({ name: 'customerServiceKnowledgeBase', query: { id: readyRobotId.value, robot_key: readyRobotKey.value } })
  } catch {
    if (isCurrentIdentity(version)) antMessage.error(t('error_open_knowledge'))
  }
}

const openRequiredVariables = () => {
  if (chatStatus.value === 'ready' && needsVariables.value) variableOpen.value = true
}
watch(lastPushedUserMessageUid, async (uid) => {
  if (!uid) return
  const version = sessionVersion.value
  reserveReplySpace.value = true
  await nextTick()
  if (isCurrentSession(version)) messageListRef.value?.scrollToMessage(uid, 'top')
})

// 只观察已就绪身份，不观察robot整个配置对象；同机器人保存刷新不会清会话。
watch([
  () => customerServiceStore.robotStatus === 'ready' ? normalizeCustomerServiceQueryValue(customerServiceStore.robot?.id) : '',
  () => customerServiceStore.robotStatus === 'ready' ? customerServiceStore.robot?.robot_key || '' : ''
], ([id, key]) => {
  identityVersion.value++
  knowledgeRequestSequence++
  clearConversationUI()
  savingSettings.value = ''
  libraryList.value = []
  wxLibraryEnabled.value = false
  chatStore.$reset()
  if (id && key) {
    void chatStore.initialize(newChatParams())
    void loadKnowledgeOptions()
  }
}, { immediate: true, flush: 'sync' })

onBeforeUnmount(() => {
  alive = false
  identityVersion.value++
  sessionVersion.value++
  knowledgeRequestSequence++
  chatStore.$reset()
})
</script>

<style lang="less" scoped>
.chat-page {
  display: flex;
  height: 100%;
  min-height: 0;
  flex: 1;
  background: var(--app-bg-container);
  color: var(--app-text-primary);
}

.chat-left {
  width: 224px;
  display: flex;
  flex-direction: column;
  background: var(--app-bg-container);
  color: var(--app-text-primary);
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
    background: var(--app-border-secondary);
  }
}

.chat-toolbar {
  height: 64px;
  padding: 0 16px;
  border-bottom: 1px solid var(--app-border-secondary);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: nowrap;
  gap: 8px;
  flex-shrink: 0;
  box-sizing: border-box;
  background: var(--app-bg-container);
  min-width: 0;

  .toolbar-title {
    display: flex;
    align-items: center;
    flex: 1;
    min-width: 0;
    gap: 8px;
    font-size: 16px;
    line-height: 24px;
    font-weight: 500;
    color: var(--app-text-primary);

    span {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.chat-history-wrapper {
  flex: 1;
  min-height: 0;
  position: relative;
  overflow: hidden;
  background: var(--app-bg-container);
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  position: relative;
  overflow: hidden;
  background: var(--app-bg-container);

  // 会话标题栏与左侧「对话 / 新对话」同高，底边才是一条线。
  :deep(.chat-conversation-bar) {
    height: 64px;
    min-height: 64px;
    padding: 0 24px;
    box-sizing: border-box;
  }

  &.is-empty {
    align-items: center;
    justify-content: center;
    
    .chat-empty-state {
      flex: none;
    }

    .form-banner-top {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
    }

    .chat-input-area {
      width: min(100%, 660px);
      margin: 48px auto 110px;
      padding: 0;
    }
  }
}

.form-banner-top {
  align-self: stretch;
  flex: none;
  height: 56px;
  margin: 12px 16px;
  padding: 16px;
  border-radius: var(--app-radius-lg);
  border: 2px solid var(--app-bg-container);
  background: linear-gradient(180deg, var(--app-selected) 0%, var(--app-bg-container) 22.78%);
  box-shadow: var(--app-shadow-card);
  display: flex;
  align-items: center;
  justify-content: space-between;
  .title {
    font-weight: 500;
    color: var(--app-text-primary);
    font-size: 16px;
  }
  .edit-block {
    border: 0;
    background: transparent;
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--app-primary);
    font-size: 14px;
    cursor: pointer;
    padding: 1px 8px;
    border-radius: var(--app-radius-control);
    &:hover:not(:disabled) {
      color: var(--app-primary-hover);
      background: var(--app-bg-hover);
    }
    &:active:not(:disabled) {
      color: var(--app-primary-active);
    }
    &:focus-visible {
      outline: 2px solid var(--app-primary);
      outline-offset: 2px;
    }
  }
}

.chat-messages {
  min-height: 0;
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  flex: 1;
  position: relative;
  overflow: hidden;
  padding: 0 16px 16px;

  :deep(.cherry-markdown) {
    background-color: transparent;
  }
}

.chat-loading-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 12px;
  color: var(--app-text-tertiary);
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
    color: var(--app-text-tertiary);
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
    border-radius: var(--app-radius-lg);
    margin: 0 auto;
    box-shadow: var(--app-shadow-card);
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
    font-weight: 500;
    color: var(--app-text-primary);
  }

  .empty-desc {
    margin-top: 24px;
    font-size: 16px;
    line-height: 24px;
    color: var(--app-text-secondary);
    white-space: pre-wrap;
  }
}

.list-loading {
  padding: 12px;
  text-align: center;
}
.empty-list {
  position: relative;
  padding: 24px 16px;
  text-align: center;
  color: var(--app-text-tertiary);
}
.history-loading {
  position: absolute;
  top: 8px;
  left: 50%;
  z-index: 1;
}
.chat-input-area {
  padding-bottom: 16px;
}
.edit-block:disabled {
  cursor: not-allowed;
  color: var(--app-text-disabled);
}

</style>
