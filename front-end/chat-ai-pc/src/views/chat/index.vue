<style lang="less" scoped>
.chat-page {
  display: flex;
  flex-flow: column nowrap;
  height: 100vh;
  width: 100vw;
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
  }

  .fast-command-wrap {
    position: relative;
    padding-top: 5px;
    z-index: 2;
    background-color: #FFF;
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
    border: 1px solid #FFF;
    background: #FFF;
    box-shadow: 0 4px 16px 0 #0000001f;
    cursor: pointer;
    position: absolute;
    bottom: 20px;
    left: 50%;
    margin-left: -20px;
    .bottom-btn {
      font-size: 16px;
      color: #659EFC;
    }
  }

  .bottom-btn-box:hover {
    border: 1px solid #659DFC;
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
</style>

<template>
  <div class="chat-page">
    <div class="chat-page-header">
      <ChatHeader />
    </div>
    <div class="chat-page-body">
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
              :prevMsg="messageList[index-1]"
              @sendTextMessage="sendTextMessage"
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
      <MessageInput v-model:value="message" @send="onSendMesage" :loading="sendLock" />
      <div class="technical-support-text">由 ChatWiki 提供软件支持</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { postInit } from '@/event/postMessage'
import { ref, onMounted, onUnmounted, toRaw, computed } from 'vue'
import { showToast } from 'vant'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useIM } from '@/hooks/event/useIM'
import { useChatStore } from '@/stores/modules/chat'
import ChatHeader from './components/chat-header.vue'
import MessageInput from './components/message-input.vue'
import MessageList from './components/messages/message-list.vue'
import MessageItem from './components/messages/message-item.vue'
import FastComand from './components/fast-comand/index.vue'

import { checkChatRequestPermission } from '@/api/robot/index'

type MessageListComponent = {  
  scrollToMessage: (id: number | string) => void;
  scrollToBottom: () => void;   
};  

const isShowBottomBtn = ref(false)

const emitter = useEventBus()
const { on } = useIM()
const chatStore = useChatStore()

const { sendMessage, onGetChatMessage, $reset, robot, externalConfigPC, openChatWindow, closeChatWindow } = chatStore

const { messageList, sendLock} = storeToRefs(chatStore)

// const isShortcut = ref(robot.yunpc_fast_command_switch == '1' ? true : false)

const isShortcut = computed(()=>{
  return robot.yunpc_fast_command_switch == '1' ? true : false
})

// 允许滚动到底部
let isAllowedScrollToBottom = true
const messageListRef = ref<null | MessageListComponent>(null);  

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
  SDKInit({robot: toRaw(robot), config: toRaw(externalConfigPC)})
}

// 发送消息
const message = ref('')

const sendTextMessage = (val: string) => {
  if (!val) {
    return showToast('请输入消息内容')
  }

  sendMessage({
    message: val
  })
}

const onSendMesage = async (content) => {
  if (!content) {
    return showToast('请输入消息内容')
  }

  //检查是否含有敏感词
  let result = await checkChatRequestPermission({
    from: 'yun_pc',
    robot_key: robot.robot_key,
    openid: robot.openid,
    question: content
  })
  if(result.data && result.data.words){
    return showToast(`提交的内容包含敏感词：[${result.data.words.join(';')}] 请修改后再提交`)
  }

  isAllowedScrollToBottom = true

  sendTextMessage(content)

  message.value = ''
}

// 监听 updateAiMessage 触发消息列表滚动
const onUpdateAiMessage = (msg) => {
  if(msg.event === 'reasoning_content'){
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

// const messageInputRef = ref<InstanceType<typeof MessageInput> | null>(null);  
const handleSetMessageInputValue = (data: any) => {
  // if(messageInputRef.value){
    // 直接发出内容
    onSendMesage(data)
    // messageInputRef.value.handleSetValue(data)
  // }
 
}

onMounted(() => {
  init()

  // 获取对话记录
  // getMyChatList()

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
