<style lang="less" scoped>
.message-item {
  display: flex;
  margin: 24px 12px;

  .message-item-body {
    flex: 1;
    padding-left: 8px;
  }

  .avatar {
    display: block;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: #fff;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.08);
  }

  .msg-img {
    width: auto;
    height: auto;
    max-width: 100%;
    max-height: 100%;
  }

  .message-content {
    margin-right: 40px;
    &:hover {
      .hover-copy-tool-block {
        display: flex;
      }
    }
  }

  .text-message {
    position: relative;
    display: inline-block;
    padding: 12px;
    font-size: 14px;
    min-height: 44px;
    line-height: 20px;
    text-align: left;
    font-weight: 400;
    white-space: pre-wrap;
    word-break: break-all;
    max-width: 100%;
  }

  &.robot-message-item {
    .text-message {
      color: #1a1a1a;
      background-color: #edeff2;
      border-radius: 4px 16px 16px 16px;
    }
    .hover-copy-tool-block {
      right: -12px;
    }
  }

  &.user-message-item {
    flex-direction: row-reverse;

    .message-item-body {
      text-align: right;
      padding-left: 0;
      padding-right: 8px;
    }

    .message-content {
      margin-left: 40px;
      margin-right: 0;
    }

    .text-message {
      border-radius: 16px 4px 16px 16px;
      color: #f5f9ff;
      background-color: #2475fc;
    }
    .hover-copy-tool-block {
      left: -12px;
    }
  }

  &.welcome-message-item {
    .text-message {
      width: 100%;
      border-radius: 4px 16px 0 0;
    }
  }

  .question-list {
    border: 1px solid #edeff2;
    border-radius: 0 0 16px 16px;
    background-color: #fff;
    .question-item {
      line-height: 20px;
      padding: 6px 12px;
      font-size: 14px;
      font-weight: 400;
      border-bottom: 1px solid #edeff2;
      color: #2475fc;
      cursor: pointer;

      &:last-child {
        border-bottom: 0;
      }
    }

    .question-item:hover {
      color: #4d94ff;
    }
  }
  .copy-block {
    display: flex;
    align-items: center;
    color: #7a8699;
    font-size: 14px;
    height: 24px;
    margin-top: 12px;
    cursor: pointer;
    width: fit-content;
    transition: all 0.5s ease;
    padding: 0 8px;
    border-radius: 6px;
    span{
      display: flex;
      height: 100%;
      line-height: 25px;
    }
    .copy-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/copy.png);
      background-size: 16px;
      width: 16px;
      height: 16px;
      margin-right: 4px;
    }
    &:hover {
      background: #F2F4F7;
      color: #3a4559;
      .copy-icon {
        background-image: url(@/assets/img/copy-hover.png);
      }
    }
  }
  .hover-copy-tool-block {
    padding: 0;
    display: none;
    position: absolute;
    bottom: -12px;
    height: 24px;
    width: 24px;
    align-items: center;
    justify-content: center;
    background: #fff;
    border: 1px solid #d8dde6;
    border-radius: 4.5px;
    transition: all 0.5s ease;
    .copy-icon{
      margin: 0;
    }
  }
}
</style>

<template>
  <div class="message-item" :class="messageItemClasses" :id="'msg-' + msg.uid">
    <div class="message-item-left">
      <img class="avatar" :src="props.msg.avatar" />
    </div>
    <div class="message-item-body">
      <div class="message-content">
        <template v-if="props.msg.msg_type == 1">
          <div class="text-message" v-if="props.msg.content !== ''" v-viewer>
            <div v-if="props.msg.is_customer == 1" v-html="props.msg.content"></div>
            <cherry-markdown :content="props.msg.content" v-else />
            <div v-tooltip="'复制'" @click="handleCopy" class="hover-copy-tool-block copy-block" v-if="isShowHoverCopy">
              <div class="copy-icon"></div>
            </div>
            <div @click="handleCopy" class="copy-block" v-if="isShowCopy">
              <div class="copy-icon"></div>
              <span>复制</span>
            </div>
          </div>
          <div v-else class="text-message">{{ textMessage }}</div>
          <div
            class="question-list"
            v-if="props.msg.menu_json && props.msg.menu_json.question.length"
          >
            <div
              class="question-item"
              @click="sendTextMessage(item)"
              v-for="(item, index) in props.msg.menu_json.question"
              :key="index"
            >
              {{ item }}
            </div>
          </div>
        </template>

        <template v-else-if="props.msg.msg_type == 2">
          <div class="text-message" v-html="escapeHTML(props.msg.menu_json.content)"></div>
          <div
            class="question-list"
            v-if="props.msg.menu_json && props.msg.menu_json.question.length"
          >
            <div
              class="question-item"
              v-for="item in props.msg.menu_json.question"
              :key="item"
              @click="sendTextMessage(item)"
            >
              <span>{{ item }}</span>
            </div>
          </div>
        </template>

        <template v-else-if="props.msg.msg_type == 3">
          <img v-viewer class="msg-img" :src="props.msg.content" />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import type { Message } from '@/stores/modules/chat'
import { useChatStore } from '@/stores/modules/chat'
import { escapeHTML } from '@/utils/index'
import { showToast } from 'vant'
import useClipboard from 'vue-clipboard3'
const { toClipboard } = useClipboard()

const emit = defineEmits(['sendTextMessage'])
const chatStore = useChatStore()
const { robot } = chatStore

const textMessage = ref('.')
let interval: number

const props = defineProps({
  msg: {
    type: Object as () => Message,
    required: true
  },
  index: {
    type: [Number, String]
  },
  messageLength: {
    type: Number,
    default: 0
  }
})

const isShowCopy = computed(() => {
  // 最后一条消息 机器人的消息 消息类型为1 不是正在发送
  return (
    props.index === props.messageLength - 1 &&
    props.msg.msg_type == 1 &&
    !robot.is_sending &&
    !isCustomerMessage.value
  )
})

const isShowHoverCopy = computed(() => {
  return !isShowCopy.value && props.index !== props.messageLength - 1
})

const handleCopy = async () => {
  await toClipboard(props.msg.content)
  showToast('复制成功')
}

// 检查是否为用户消息
const isCustomerMessage = computed(() => props.msg.is_customer == 1)

// 计算消息项的类
const messageItemClasses = computed(() => ({
  'user-message-item': isCustomerMessage.value === true,
  'robot-message-item': isCustomerMessage.value === false,
  'welcome-message-item': props.msg.menu_json && props.msg.menu_json.question
}))

// 等待机器人回复增加动态...
const startLoadingAnimation = () => {
  const dots = ['.', '..', '...']
  let dotIndex = 0
  interval = window.setInterval(() => {
    dotIndex = (dotIndex + 1) % dots.length
    textMessage.value = dots[dotIndex]
  }, 500)
}

const sendTextMessage = (text: string) => {
  emit('sendTextMessage', text)
}

onMounted(() => {
  startLoadingAnimation()
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>
