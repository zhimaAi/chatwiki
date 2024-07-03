<style lang="less" scoped>
.ignore-message-item {
  max-width: 760px;
  display: flex;
  padding: 24px 12px 0px;
  margin: 0 auto;

  .message-item-body {
    flex: 1;
    padding-left: 10px;
    overflow: hidden;
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
    position: relative;
    display: inline-block;
    padding: 12px;
    margin-right: 50px;
    border-radius: 4px 16px 16px 16px;
    background-color: #fff;
    white-space: pre-wrap;
    word-break: break-all;
    max-width: 100%;

    .triangle {
      position: absolute;
      top: 20px;
      width: 0;
      height: 0;
      border-top: 6px solid transparent;
      border-bottom: 6px solid transparent;
      border-right: 6px solid white;
    }
  }

  .text-message {
    min-height: 20px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    text-align: left;
    color: rgb(26, 26, 26);
  }

  .question-list {
    padding-bottom: 4px;

    .question-item {
      cursor: pointer;
      line-height: 20px;
      padding: 6px 12px;
      margin-top: 8px;
      font-size: 14px;
      font-weight: 400;
      border-radius: 4px;
      color: #164799;
      background-color: #e6efff;
    }

    .question-item:hover {
      background: #dee8fb;
    }

    .question-item:active {
      background-color: #e6efff;
    }
  }
  .question-list.guess-you-want {
    border-top: 1px solid #edeff2;
    margin-top: 16px;
    .question-list-title {
      margin-top: 9px;
      color: #7a8699;
      font-size: 14px;
      line-height: 22px;
    }
  }

  &.robot-message-item {
    .triangle {
      left: -5px;
    }
  }

  &.user-message-item {
    max-width: 760px;
    flex-direction: row-reverse;

    .message-item-body {
      text-align: right;
      padding-left: 0;
      padding-right: 10px;
    }

    .message-content {
      border-radius: 16px 4px 16px 16px;
      margin-left: 50px;
      margin-right: 0;
      background-color: #2475fc;
      color: white;
    }

    .triangle {
      right: -6px;
      transform: rotate(180deg);
      border-right: 6px solid #2475fc;
    }

    .text-message {
      color: #f5f9ff;
    }
  }
}
.ignore-message-item:last-child {
  padding-bottom: 24px;
}
</style>

<template>
  <div class="ignore-message-item" :class="messageItemClasses" :id="'msg-' + msg.uid">
    <div class="message-item-left">
      <img class="avatar" :src="props.msg.avatar" />
    </div>
    <div class="message-item-body">
      <div class="message-content">
        <!-- <span class="triangle"></span> -->
        <template v-if="props.msg.msg_type == 1">
          <div v-if="props.msg.content !== ''" class="text-message">
            <cherry-markdown :content="props.msg.content" />
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
          <img class="msg-img" :src="props.msg.content" />
        </template>
      </div>

      <div v-if="props.msg.msg_type == 1">
        <GuessYouWant
          v-if="
            ((props.msg.guess_you_want && props.msg.guess_you_want.length) ||
              (common_question_list.length && robot.enable_common_question)) &&
            props.msg.question_tabkey > 0
          "
          :msg="props.msg"
          :common_question_list="common_question_list"
          :enable_common_question="robot.enable_common_question"
          @sendTextMessage="sendTextMessage"
        />
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
import GuessYouWant from './guess-you-want.vue'

const emit = defineEmits(['sendTextMessage'])
const chatStore = useChatStore()
const { robot } = chatStore
const textMessage = ref('.')
let interval: number

const props = defineProps({
  msg: {
    type: Object as () => Message,
    required: true
  }
})
const common_question_list = computed(() => robot.common_question_list)
// 检查是否为用户消息
const isCustomerMessage = computed(() => props.msg.is_customer == 1)

// 计算消息项的类
const messageItemClasses = computed(() => ({
  'user-message-item': isCustomerMessage.value === true,
  'robot-message-item': isCustomerMessage.value === false,
  'welcome-message-item': props.msg.menu_json && props.msg.menu_json.question
}))

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
