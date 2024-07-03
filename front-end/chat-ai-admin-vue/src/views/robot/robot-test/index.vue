<style lang="less" scoped>
.chat-page-wrapper {
  flex: 1;
  display: flex;
  flex-flow: column nowrap;
  border-radius: 2px;
  overflow: hidden;
  background-color: #f5f9ff;
}
.chat-page-header {
  padding: 0 24px;
  .breadcrumb-box {
    display: flex;
    align-items: center;
    height: 56px;
  }
}
.chat-page {
  display: flex;
  justify-content: space-between;
  flex: 1;
  padding: 24px;
  overflow: hidden;

  .page-body {
    flex: 1;
    display: flex;
    flex-flow: column nowrap;
    max-width: 46.6%;
    height: 100%;
    overflow: hidden;

    .chat-box-body {
      flex: 1;
      overflow: hidden;
    }
  }
  .page-left,
  .page-right {
    width: 280px;
    height: 100%;
    border-radius: 8px;
    overflow-y: auto;
    background-color: #fafafa;
    box-shadow: 0 2px 16px 0 rgba(14, 27, 58, 0.16);
  }

  .page-left {
    display: flex;
    flex-flow: column nowrap;
    padding: 24px 0;

    .page-left-top,
    .page-left-footer {
      padding: 0 24px;
    }

    .page-left-top {
      margin-bottom: 15px;
    }

    .page-left-body {
      position: relative;
      flex: 1;
      overflow: hidden;
    }
  }
}

.page-right {
  display: flex;
  flex-flow: column nowrap;
  justify-content: space-between;
  padding: 24px;

  .prompt-tips {
    line-height: 22px;
    padding: 16px;
    font-size: 14px;
    border-radius: 4px;
    color: #3a4559;
    background: #f2f4f7;
  }
}
</style>

<template>
  <div class="chat-page-wrapper">
    <div class="chat-page-header">
      <div class="breadcrumb-box">
        <a-breadcrumb>
          <a-breadcrumb-item
            ><router-link to="/robot/list">机器人管理</router-link></a-breadcrumb-item
          >
          <a-breadcrumb-item>{{ robot.robot_name }}</a-breadcrumb-item>
        </a-breadcrumb>
      </div>
    </div>
    <div class="chat-page">
      <!-- left -->
      <div class="page-left">
        <div class="page-left-top">
          <a-button type="primary" block ghost @click="openNewChat"
            ><PlusOutlined /> 新建对话</a-button
          >
        </div>
        <div class="page-left-body">
          <ChatList
            :list="myChatList"
            :active="dialogue_id"
            @openChat="handleOpenChat"
            @onScrollEnd="onChatListScrollEnd"
          />
        </div>
        <div class="page-left-footer"></div>
      </div>
      <!-- body -->
      <div class="page-body">
        <div class="chat-box-body">
          <MessageList
            ref="messageListRef"
            :messages="messageList"
            :robotInfo="robotInfo"
            @clickMsgMeun="onClickMsgMenu"
            @scrollStart="onScrollStart"
            @scrollEnd="onScrollEnd"
            @openPromptLog="handleOpenPromptLog"
            @openLibrary="handleOpenLibraryInfo"
          />
        </div>

        <div style="margin-top: 30px" class="chat-box-footer">
          <MessageInput v-model:value="message" @send="onSendMesage" :loading="sendLock" />
        </div>
      </div>
      <!-- right -->
      <RobotTestRight
        v-if="isRobotInfo"
        @promptChange="onPromptChange"
        @saveRobotPrompt="onSaveRobotPrompt"
        :robotInfo="robotInfo"
      />
    </div>

    <PromptLogAlert ref="promptLogAlertRef" />

    <libraryInfoAlert ref="libraryInfoAlertRef" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onUnmounted, createVNode, toRaw, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useChatStore } from '@/stores/modules/chat'
import { useUserStore } from '@/stores/modules/user'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import { showErrorMsg, showSuccessMsg } from '@/utils/index'
import MessageList from './components/message-list.vue'
import ChatList from './components/chat-list.vue'
import RobotTestRight from './components/robot-test-right.vue'
import MessageInput from './components/message-input.vue'
import PromptLogAlert from './components/prompt-log-alert.vue'
import libraryInfoAlert from './components/library-info-alert.vue'
import { saveRobot } from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'

const rotue = useRoute()
const query = rotue.query

const robotStore = useRobotStore()

const { getRobot, robotInfo } = robotStore

const emitter = useEventBus()
const chatStore = useChatStore()
const userStore = useUserStore()

// 是否允许自动滚动到底部
let isAllowedScrollToBottom = true
const message = ref('')
const messageListRef = ref(null)

const {
  createChat,
  sendMessage,
  getMyChatList,
  openChat,
  onGetChatMessage,
  changeRobotPrompt,
  $reset
} = chatStore
const { messageList, sendLock, myChatList, robot, dialogue_id } = storeToRefs(chatStore)

const route = useRoute()
const isRobotInfo = ref(false)

const onSendMesage = async () => {
  if (!message.value) {
    return showErrorMsg('请输入消息内容')
  }

  isAllowedScrollToBottom = true

  sendMessage({
    message: message.value
  })

  message.value = ''
}

const openNewChat = async () => {
  isAllowedScrollToBottom = true
  message.value = ''

  let query = route.query || {}
  let user_id = userStore.user_id

  let data = {
    robot_key: query.robot_key,
    avatar: query.avatar,
    name: query.name,
    nickname: query.nickname,
    is_background: 1,
    openid: user_id,
    dialogue_id: 0
  }

  resetScroll()

  await createChat(data)
}

const handleOpenChat = async (data) => {
  if (dialogue_id.value == data.dialogue_id) {
    return
  }

  isAllowedScrollToBottom = true

  let query = route.query || {}

  let params = {
    robot_key: query.robot_key,
    avatar: data.avatar,
    name: data.name,
    nickname: data.nickname,
    is_background: data.is_background,
    openid: data.openid,
    dialogue_id: data.id
  }

  resetScroll()

  await openChat(params)

  let res = await onGetChatMessage()
  if (res) {
    messageListScrollToBottom()
  }
}

// 点击菜单
const onClickMsgMenu = (text) => {
  isAllowedScrollToBottom = true

  sendMessage({
    message: text
  })
}

const saveLoading = ref(false)

const saveForm = (formState) => {
  let formData = JSON.parse(JSON.stringify(toRaw(formState)))
  saveLoading.value = true
  let welcomes = formData.welcomes

  welcomes.question = welcomes.question.map((item) => {
    return item.content
  })

  let unknown_question_prompt = formData.unknown_question_prompt

  unknown_question_prompt.question = unknown_question_prompt.question.map((item) => {
    return item.content
  })

  formData.unknown_question_prompt = JSON.stringify(unknown_question_prompt)

  formData.welcomes = JSON.stringify(welcomes)

  saveLoading.value = true
  saveRobot(formData)
    .then((res) => {
      if (res.res != 0) {
        return showErrorMsg(res.msg)
      }

      saveLoading.value = false

      showSuccessMsg('保存成功')

      getRobot(formState.id)
    })
    .catch(() => {
      saveLoading.value = false
    })
}

// 提示词
const onPromptChange = (e) => {
  changeRobotPrompt(e.currentTarget.value)
}

const onSaveRobotPrompt = async (formState, isDefault) => {
  let newFormState = JSON.parse(JSON.stringify(formState)); // 深拷贝，不能改变原对象
  let content =
    '如果您已对外提供本机器人，保存后修改的提示问也会立刻生效。如果您只是想要测试优化提示词的效果，直接修改后测试即可，无需保存。'
  Modal.confirm({
    title: '确定保存提示词吗?',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode('div', { style: 'color:red;' }, content),
    onOk: async () => {
      // updateRobotInfo({ ...toRaw(formState) })
      if (isDefault) {
        // 传给后端的是默认，渲染的是真实名称
        newFormState.use_model = '默认'
      }
      saveForm(newFormState)
    },
    onCancel() {}
  })
}

// 监听 updateAiMessage 触发消息列表滚动
const onUpdateAiMessage = () => {
  messageListScrollToBottom()
}

// 滚动到顶部
const onScrollStart = async () => {
  isAllowedScrollToBottom = false

  let msgId = messageList.value[0].uid

  let res = await onGetChatMessage()

  if (res) {
    scrollToMessageById(msgId)
  }
}

// 监听滚动到底部
const onScrollEnd = () => {
  // console.log('滚动到底部')
}

const messageListScrollToBottom = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
  }
}

const scrollToMessageById = (id) => {
  if (messageListRef.value) {
    messageListRef.value.scrollToMessage(id)
  }
}

const resetScroll = () => {
  if (messageListRef.value) {
    messageListRef.value.resetScroll()
  }
}

// 监听 messageList 改变
const onMessageListChange = () => {
  if (!isAllowedScrollToBottom) {
    return
  }
}

// 加载对话记录(对话记录加载下一页)
const onChatListScrollEnd = () => {
  // 获取对话记录
  getMyChatList()
}

// 打开Prompt日志
const promptLogAlertRef = ref(null)

const handleOpenPromptLog = (item) => {
  promptLogAlertRef.value.open(item)
}

// 打开文件知识库
const libraryInfoAlertRef = ref(null)

const handleOpenLibraryInfo = (files, file) => {
  libraryInfoAlertRef.value.open(files, file)
}

watch(
  () => messageList.value,
  () => {
    onMessageListChange()
  }
)

onMounted(async () => {
  // 创建对话
  openNewChat()
  // 获取对话记录
  getMyChatList()

  await getRobot(query.id)
  isRobotInfo.value = true
  // 监听 updateAiMessage 触发消息列表滚动
  emitter.on('updateAiMessage', onUpdateAiMessage)
})

onUnmounted(() => {
  $reset()
  emitter.off('updateAiMessage', onUpdateAiMessage)
})
</script>
