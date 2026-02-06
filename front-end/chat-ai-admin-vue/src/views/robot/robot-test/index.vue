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

.form-banner-top{
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
    &:hover{
      background: var(--07, #E4E6EB);
      border-radius: 6px;
    }
  }
}

</style>

<template>
  <div class="chat-page-wrapper">
    <div class="chat-page-header">
      <div class="breadcrumb-box">
        <a-breadcrumb>
          <a-breadcrumb-item
            ><router-link to="/robot/list">{{ t('breadcrumb_robot_management') }}</router-link></a-breadcrumb-item
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
            ><PlusOutlined /> {{ t('btn_new_chat') }}</a-button
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
        <div class="form-banner-top" v-if="isShowFromHeader">
          <div class="title">{{ t('title_form_info') }}</div>
          <div class="edit-block" @click="handleEditVariableForm">
            <EditOutlined />{{ t('btn_edit') }}
          </div>
        </div>
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

        <div style="padding-top: 30px" class="chat-box-footer">
          <MessageInput
            v-model:value="message"
            v-model:fileList="fileList"
            :loading="sendLoading"
            :showUpload="robotInfo.question_multiple_switch == 1"
            @send="onSendMesage" />
        </div>
      </div>
      <!-- right -->
      <div>
        <RobotTestRight
          v-if="false"
          @promptChange="onPromptChange"
          @saveRobotPrompt="onSaveRobotPrompt"
          :robotInfo="robotInfo"
        />
      </div>
    </div>

    <PromptLogAlert ref="promptLogAlertRef" />

    <libraryInfoAlert ref="libraryInfoAlertRef" />
    <VariableModal ref="variableModalRef" />
  </div>
</template>

<script setup>
import { generateRandomId } from '@/utils/index'
import { ref, computed, onMounted, watch, onUnmounted, createVNode, toRaw } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useChatStore } from '@/stores/modules/chat'
import { useUserStore } from '@/stores/modules/user'
import { ExclamationCircleOutlined, PlusOutlined, EditOutlined } from '@ant-design/icons-vue'
import { Modal, message as antMessage } from 'ant-design-vue'
import { showErrorMsg, showSuccessMsg } from '@/utils/index'
import MessageList from './components/message-list.vue'
import ChatList from './components/chat-list.vue'
import RobotTestRight from './components/robot-test-right.vue'
import MessageInput from './components/message-input.vue'
import PromptLogAlert from './components/prompt-log-alert.vue'
import libraryInfoAlert from './components/library-info-alert.vue'
import VariableModal from './components/variable-modal/index.vue'
import { saveRobot, checkChatRequestPermission } from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-test.index')


const rotue = useRoute()
const query = rotue.query

const robotStore = useRobotStore()

const { getRobot } = robotStore

const robotInfo = computed(() => robotStore.robotInfo)
const emitter = useEventBus()
const chatStore = useChatStore()
const userStore = useUserStore()
// 是否允许自动滚动到底部
let isAllowedScrollToBottom = true
const message = ref('')
const fileList = ref([])
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

const isShowFromHeader = computed(()=>{
  return !chatStore.chat_variables.need_fill_variable && chatStore.chat_variables.fill_variables && chatStore.chat_variables.fill_variables.length
})

const route = useRoute()
const isRobotInfo = ref(false)
const checkChatRequestPermissionLoding = ref(false)
const sendLoading = computed(() => sendLock.value || checkChatRequestPermissionLoding.value)

const onSendMesage = async () => {
  if (sendLoading.value){
    return
  }

  if (!message.value && !fileList.value.length) {
    return showErrorMsg(t('msg_please_input_message'))
  }

  checkChatRequestPermissionLoding.value = true

  try {
    //检查是否含有敏感词
  let result = await checkChatRequestPermission({
    robot_key: robot.value.robot_key,
    openid: robot.value.openid,
    question: message.value,
    form_ids: robot.value.form_ids,
    dialogue_id: dialogue_id.value,
  })

  checkChatRequestPermissionLoding.value = false

  if(result.data && result.data.words){
    return antMessage.error(t('msg_sensitive_words', { words: result.data.words.join(';') }))
  }

  let msg_type = 1;
  let msg = message.value;

  if(robotInfo.value.question_multiple_switch == 1){
    let messageList = []

    msg_type = 99;

    if(message.value){
      messageList = [{
        type: "text",
        uid: generateRandomId(16),
        text: message.value
      }]
    }

    if(fileList.value.length){
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




  isAllowedScrollToBottom = true

  sendMessage({
    message: msg,
    global: JSON.stringify(query),
    msg_type: msg_type
  })

  message.value = ''
  fileList.value = []
  }catch(err){
    checkChatRequestPermissionLoding.value = false
  }

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

      showSuccessMsg(t('success_save'))

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
  let newFormState = JSON.parse(JSON.stringify(formState)) // 深拷贝，不能改变原对象
  Modal.confirm({
        title: t('confirm_save_prompt'),
        icon: createVNode(ExclamationCircleOutlined),
        content: createVNode('div', { style: 'color:red;' }, t('confirm_save_prompt_content')),
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

const variableModalRef = ref(null)
const handleEditVariableForm = () => {
  variableModalRef.value.handleEdit()
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
