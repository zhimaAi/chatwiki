<template>
  <div class="chat-test-panel" :class="{ 'has-result': hasRunTested }">
    <div class="chat-pane">
      <template v-if="!hasRunTested">
        <div class="chat-empty-state">
          <div class="workflow-icon">
            <img :src="workflowAvatar" alt="" />
          </div>
          <div class="workflow-title">{{ workflowTitle }}</div>
          <div class="workflow-desc" v-if="workflowDesc">{{ workflowDesc }}</div>
          <div class="chat-input-center">
            <MessageInput
              v-model:value="messageText"
              v-model:fileList="fileList"
              :loading="loading"
              :showUpload="questionMultipleSwitch"
              :showStop="false"
              @send="handleInputSend"
            />
          </div>
        </div>
      </template>
      <template v-else>
        <div class="chat-messages customize-scroll-style">
          <MessageList
            ref="messageListRef"
            :messages="messageList"
            :robotInfo="robotStore.robotInfo"
            :showAvatar="false"
            messageAlign="split"
            :bubbleStyle="messageBubbleStyle"
            @clickMsgMeun="handleMenuClick"
            @scroll="handleMessageListScroll"
          />
        </div>
        <div class="chat-input-fixed">
          <MessageInput
            v-model:value="messageText"
            v-model:fileList="fileList"
            :loading="loading"
            :showUpload="questionMultipleSwitch"
            :showStop="false"
            @send="handleInputSend"
          />
        </div>
      </template>
    </div>

    <NodeRunLogs
      v-if="hasRunTested"
      v-model:currentNodeKey="currentNodeKey"
      :resultList="resultList"
      :loading="loading"
      :hasRunTested="hasRunTested"
    />
  </div>
</template>

<script setup>
import { computed, nextTick, onUnmounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { generateRandomId } from '@/utils/index'
import { callWorkFlowDialog } from '@/api/robot/index'
import { getImageUrl } from '../util'
import { useI18n } from '@/hooks/web/useI18n'
import { useRobotStore } from '@/stores/modules/robot'
import { DEFAULT_USER_AVATAR, DEFAULT_WORKFLOW_AVATAR } from '@/constants/index'
import MessageInput from '@/views/robot/robot-test/components/message-input.vue'
import MessageList from '@/views/robot/robot-test/components/message-list.vue'
import NodeRunLogs from './NodeRunLogs.vue'

const { t } = useI18n('views.workflow.components.run-test.index')
const robotStore = useRobotStore()

const props = defineProps({
  start_node_params: {
    default: () => {},
    type: Object
  }
})

const emit = defineEmits(['stateChange'])
const query = useRoute().query

const messageText = ref('')
const fileList = ref([])
const messageList = ref([])
const messageListRef = ref(null)
const resultList = ref([])
const currentNodeKey = ref('')
const loading = ref(false)
const hasRunTested = ref(false)
const use_token = ref(0)
const use_mills = ref(0)
const messageBubbleStyle = {
  user: {
    backgroundColor: '#e6efff',
    color: '#262626'
  },
  robot: {
    backgroundColor: '#fff',
    color: '#262626'
  }
}
// 会话 ID 可能超过安全整数范围，始终按字符串保存
const dialogue_id = ref('')
const session_id = ref('')
const testUserId = ref('')
const isAwaitingFollowUp = ref(false)
let mySSE = null
let currentAiMessage = null
let hasFinalizedCurrentRun = false
let shouldAutoScroll = true
let scrollFrameId = 0
let forceNextScroll = false
const AUTO_SCROLL_THRESHOLD = 60

const handleMessageListScroll = ({ scrollTop, scrollHeight, clientHeight }) => {
  shouldAutoScroll = scrollHeight - scrollTop - clientHeight <= AUTO_SCROLL_THRESHOLD
}

const scheduleMessageScroll = (force = false) => {
  if (force) {
    forceNextScroll = true
  }
  if ((!shouldAutoScroll && !forceNextScroll) || scrollFrameId) {
    return
  }
  scrollFrameId = requestAnimationFrame(async () => {
    scrollFrameId = 0
    const shouldForce = forceNextScroll
    forceNextScroll = false
    if (!shouldForce && !shouldAutoScroll) {
      return
    }
    await nextTick()
    messageListRef.value?.scrollToBottom()
  })
}

const resetMessageScroll = () => {
  shouldAutoScroll = true
  forceNextScroll = false
  if (scrollFrameId) {
    cancelAnimationFrame(scrollFrameId)
    scrollFrameId = 0
  }
}

const workflowTitle = computed(() => {
  return robotStore.robotInfo.robot_name || t('title_chat_test_default_name')
})

const workflowDesc = computed(() => {
  return robotStore.robotInfo.robot_intro || ''
})

const workflowAvatar = computed(() => {
  return robotStore.robotInfo.robot_avatar_url || robotStore.robotInfo.robot_avatar || DEFAULT_WORKFLOW_AVATAR
})

const questionMultipleSwitch = computed(() => {
  let trigger_list = props.start_node_params.trigger_list || []
  let result = false
  trigger_list.forEach((item) => {
    if (item.trigger_type == 1) {
      result = item.chat_config?.question_multiple_switch
    }
  })
  return result
})

const syncState = () => {
  emit('stateChange', {
    hasRunTested: hasRunTested.value,
    hasRunResult: resultList.value.length > 0,
    use_token: use_token.value,
    use_mills: use_mills.value
  })
}

const open = () => {
  abortSSE()
  resetMessageScroll()
  messageText.value = ''
  fileList.value = []
  messageList.value = []
  resultList.value = []
  currentNodeKey.value = ''
  loading.value = false
  hasRunTested.value = false
  use_token.value = 0
  use_mills.value = 0
  dialogue_id.value = ''
  session_id.value = ''
  testUserId.value = ''
  isAwaitingFollowUp.value = false
  currentAiMessage = null
  hasFinalizedCurrentRun = false
  syncState()
}

const handleSend = (menuText = '', isMultimodal = questionMultipleSwitch.value) => {
  if (loading.value) {
    return
  }
  const content = menuText || messageText.value
  const selectedFiles = menuText ? [] : fileList.value
  if (!content && !selectedFiles.length) {
    return message.warning(t('msg_input_chat_question'))
  }
  const uploading = selectedFiles.some((file) => file.status === 'uploading')
  if (uploading) {
    return message.warning(t('msg_wait_image_upload'))
  }
  const failed = selectedFiles.some((file) => file.status === 'error')
  if (failed) {
    return message.warning(t('msg_remove_failed_image'))
  }

  abortSSE()
  const shouldContinue = isAwaitingFollowUp.value && dialogue_id.value && session_id.value && testUserId.value
  use_token.value = 0
  use_mills.value = 0
  if (!shouldContinue) {
    resultList.value = []
    currentNodeKey.value = ''
    dialogue_id.value = ''
    session_id.value = ''
    testUserId.value = generateRandomId(16)
    messageList.value = []
  }
  isAwaitingFollowUp.value = false
  hasFinalizedCurrentRun = false
  hasRunTested.value = true
  syncState()

  const questionItems = buildQuestionItems(content, selectedFiles)
  const question = isMultimodal ? JSON.stringify(questionItems) : content
  const userMessage = {
    uid: generateRandomId(16),
    is_customer: 1,
    msg_type: isMultimodal ? 99 : 1,
    content: question,
    avatar: DEFAULT_USER_AVATAR
  }
  const aiMessage = {
    uid: generateRandomId(16),
    is_customer: 0,
    msg_type: 1,
    content: '',
    reply_content_list: [],
    robot_avatar: workflowAvatar.value,
    loading: true,
    quote_file: [],
    debug: [],
    error: '',
    show_reasoning: false,
    reasoning_status: false,
    quote_loading: false,
    show_quote_file: true,
    is_stopped: false,
    voice_content: []
  }
  messageList.value.push(userMessage, aiMessage)
  shouldAutoScroll = true
  scheduleMessageScroll(true)
  currentAiMessage = messageList.value[messageList.value.length - 1]
  messageText.value = ''
  fileList.value = []
  loading.value = true

  const sse = callWorkFlowDialog({
    robot_key: query.robot_key,
    openid: testUserId.value,
    question,
    dialogue_id: dialogue_id.value,
    session_id: session_id.value
  })
  mySSE = sse

  sse.onMessage = (res) => {
    // 忽略旧请求的迟到回调，避免污染当前测试结果
    if (mySSE === sse) {
      handleSSEMessage(res)
    }
  }
  sse.onClose = () => {
    if (mySSE !== sse) {
      return
    }
    finalizeCurrentRun()
    mySSE = null
  }
  sse.onError = () => {
    if (mySSE !== sse) {
      return
    }
    loading.value = false
    mySSE = null
    dialogue_id.value = ''
    session_id.value = ''
    testUserId.value = ''
    isAwaitingFollowUp.value = false
    if (currentAiMessage) {
      currentAiMessage.loading = false
    }
    message.error(t('msg_chat_test_failed'))
    syncState()
  }
}

const handleInputSend = () => {
  handleSend()
}

const handleMenuClick = (text) => {
  // 智能菜单仅支持纯文本，不能继承输入框的多模态消息格式。
  handleSend(text, false)
}

const buildQuestionItems = (content, selectedFiles) => {
  const list = []
  if (content) {
    list.push({
      uid: generateRandomId(16),
      type: 'text',
      text: content
    })
  }
  selectedFiles.forEach((file) => {
    list.push({
      uid: file.uid,
      type: 'image_url',
      image_url: {
        url: file.url
      }
    })
  })
  return list
}

const handleSSEMessage = (res) => {
  if (res.event === 'ping') {
    return
  }
  // 流式文本必须保持原始字符串，避免数字或 JSON 格式文本被错误转换
  if (res.event === 'sending') {
    appendEventLog(res.event, res.data)
    appendAiContent(res.data)
    return
  }

  const data = parseEventData(res.data)

  appendEventLog(res.event, data)

  if (res.event === 'dialogue_id') {
    dialogue_id.value = res.data ? String(res.data).trim() : ''
    return
  }
  if (res.event === 'session_id') {
    session_id.value = res.data ? String(res.data).trim() : ''
    return
  }
  if (res.event === 'use_token') {
    use_token.value = Number(data || 0)
    syncState()
    return
  }
  if (res.event === 'use_mills') {
    use_mills.value = Number(data || 0)
    syncState()
    return
  }
  if (res.event === 'robot' || res.event === 'reply_content') {
    appendAiContent(typeof data === 'string' ? data : data?.content || data?.answer || '')
    return
  }
  if (res.event === 'reply_content_list') {
    appendReplyContentList(data)
    return
  }
  if (res.event === 'ai_message' || res.event === 'message') {
    appendAiContent(typeof data === 'string' ? data : data?.content || data?.answer || data?.message || '')
  }

  const nodeLogs = pickNodeLogs(data)
  if (nodeLogs.length) {
    mergeNodeLogs(nodeLogs)
  }

  const token = data?.use_token ?? data?.token_usage
  const mills = data?.use_mills
  if (token != null) {
    use_token.value = Number(token || 0)
  }
  if (mills != null) {
    use_mills.value = Number(mills || 0)
  }
  if (res.event === 'stat') {
    syncState()
  }

  if (res.event === 'finish') {
    if (data?.use_token != null) {
      use_token.value = Number(data.use_token || 0)
    }
    if (data?.use_mills != null) {
      use_mills.value = Number(data.use_mills || 0)
    }
    // finish 仅表示当前 SSE 结束，问答节点仍需保留会话标识供下一次回复续跑。
    finalizeCurrentRun()
  }
}

const appendAiContent = (content) => {
  if (content == null || !currentAiMessage) {
    return
  }
  currentAiMessage.content += String(content)
  scheduleMessageScroll()
}

function appendReplyContentList(value) {
  const replyContentList = parseReplyContentList(value)
  if (!replyContentList.length || !currentAiMessage) {
    return false
  }
  currentAiMessage.reply_content_list = replyContentList
  scheduleMessageScroll()
  return true
}

function parseReplyContentList(value) {
  if (Array.isArray(value)) {
    return value
  }
  if (typeof value !== 'string' || !value) {
    return []
  }
  try {
    const parsed = JSON.parse(value)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

const appendEventLog = (event, data) => {
  if (!resultList.value.length) {
    return
  }
  const current = resultList.value.find((item) => getNodeLogKey(item) === currentNodeKey.value) || resultList.value[0]
  const logs = Array.isArray(current.output?.events) ? current.output.events : []
  current.output = {
    ...(current.output || {}),
    events: [
      ...logs,
      {
        event,
        data
      }
    ]
  }
}

function parseEventData(data) {
  if (data == null || data === '') {
    return ''
  }
  if (typeof data !== 'string') {
    return data
  }
  const text = data.trim()
  if (!text.startsWith('{') && !text.startsWith('[')) {
    return data
  }
  try {
    return JSON.parse(text)
  } catch {
    return data
  }
}

function pickNodeLogs(data) {
  if (Array.isArray(data)) {
    return data
  }
  const candidates = [
    data?.node_logs,
    data?.nodeLogs,
    data?.workflow_node_logs,
    data?.data?.node_logs
  ]
  for (const item of candidates) {
    if (Array.isArray(item)) {
      return item
    }
  }
  if (data && (data.node_key || data.node_name || data.node_type)) {
    return [data]
  }
  return []
}

function mergeNodeLogs(logs) {
  logs.forEach((item) => {
    const formatted = formatNodeLog(item)
    const index = resultList.value.findIndex((log) => getNodeLogKey(log) === formatted.log_key)
    if (index === -1) {
      resultList.value.push(formatted)
    } else {
      resultList.value.splice(index, 1, {
        ...resultList.value[index],
        ...formatted
      })
    }
    // 兼容智能菜单仅包含在问答节点日志、未单独发送事件的情况
    if (formatted.node_type === 43) {
      appendReplyContentList(
        formatted.output?.special?.reply_content_list || formatted.node_output?.special?.reply_content_list
      )
    }
  })
  if (!currentNodeKey.value) {
    currentNodeKey.value = getNodeLogKey(resultList.value[0])
  }
  updateRunSummaryFromNodeLogs()
  syncState()
}

function updateRunSummaryFromNodeLogs() {
  use_mills.value = resultList.value.reduce((total, item) => total + toNumber(item.use_time ?? item.use_mills), 0)

  const tokenNodeList = resultList.value.filter((item) => getNodeUseToken(item) > 0 && item.node_type !== 7)
  const summaryNodeList = tokenNodeList.length ? tokenNodeList : resultList.value
  use_token.value = summaryNodeList.reduce((total, item) => total + getNodeUseToken(item), 0)
}

function getNodeUseToken(item) {
  const llmResult = item.output?.llm_result || item.node_output?.llm_result
  if (!llmResult) {
    return 0
  }
  const totalToken = toNumber(llmResult.total_token ?? llmResult.total_tokens)
  if (totalToken > 0) {
    return totalToken
  }
  return toNumber(llmResult.prompt_token) + toNumber(llmResult.completion_token)
}

function toNumber(value) {
  const numberValue = Number(value)
  return Number.isFinite(numberValue) ? numberValue : 0
}

function formatNodeLog(item) {
  let nodeIcon = getImageUrl(item.node_type)
  if (item.node_type == 45) {
    nodeIcon = item.node_icon || getImageUrl(item.node_type)
  }
  const nodeKey = item.node_key || item.key || generateRandomId(12)
  return {
    input: {},
    output: {},
    node_output: {},
    use_time: item.use_time || item.use_mills || 0,
    ...item,
    node_key: nodeKey,
    // 同一节点在一个会话中可能执行多次，使用开始时间区分每次执行日志
    log_key: item.log_key || (item.start_time != null ? `${nodeKey}_${item.start_time}` : nodeKey),
    node_name: item.node_name || item.name || t('label_unknown_node'),
    is_success: item.error_msg ? item.error_msg === '<nil>' : item.is_success !== false,
    node_icon: item.node_icon || nodeIcon
  }
}

function getNodeLogKey(item) {
  return item?.log_key || item?.node_key || ''
}

function finalizeCurrentRun() {
  if (hasFinalizedCurrentRun) {
    return
  }
  hasFinalizedCurrentRun = true
  loading.value = false
  const lastNode = resultList.value[resultList.value.length - 1]
  // 问答节点会暂停流程，保留会话标识供下一次回复继续执行
  isAwaitingFollowUp.value = lastNode?.node_type === 43
  if (!isAwaitingFollowUp.value) {
    dialogue_id.value = ''
    session_id.value = ''
    testUserId.value = ''
  }
  if (currentAiMessage) {
    currentAiMessage.loading = false
  }
  syncState()
}

function abortSSE() {
  if (mySSE) {
    mySSE.abort()
    mySSE = null
  }
}

onUnmounted(() => {
  abortSSE()
  resetMessageScroll()
})

defineExpose({
  open,
  abort: abortSSE
})
</script>

<style lang="less" scoped>
.chat-test-panel {
  display: flex;
  height: 70vh;
  overflow: hidden;
  align-items: stretch;
  background: #fff;

  &.has-result {
    .chat-pane {
      flex: 0 0 var(--chat-test-chat-width, 382px);
      width: var(--chat-test-chat-width, 382px);
      min-width: var(--chat-test-chat-width, 382px);
      box-sizing: border-box;
      border-right: 1px solid #f0f0f0;
    }
  }
}

.chat-pane {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.chat-empty-state {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0 72px;
}

.workflow-icon {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.workflow-title {
  margin-top: 28px;
  font-size: 32px;
  line-height: 40px;
  font-weight: 700;
  color: #141414;
}

.workflow-desc {
  margin-top: 18px;
  font-size: 14px;
  color: #6b7890;
}

.chat-input-center {
  width: 700px;
  max-width: 100%;
  margin-top: 48px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.chat-input-fixed {
  padding: 14px 16px 16px;
  background: #fff;
}
</style>
