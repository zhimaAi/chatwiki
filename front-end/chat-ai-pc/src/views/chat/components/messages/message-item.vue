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
        flex-direction: column;
        align-items: center;
        justify-content: center;

        .copy-block{
          width: 26px;
          height: 26px;
          padding: 0 4px;
        }
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
    cursor: pointer;
    width: fit-content;
    transition: all 0.5s ease;
    padding: 0 8px;
    height: 24px;
    border-radius: 6px;
    span{
      display: flex;
      height: 100%;
      line-height: 25px;
      margin-left: 2px;
    }
    .copy-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/copy.png);
      background-size: 16px;
      width: 16px;
      height: 16px;
    }
    &:hover {
      background: #F2F4F7;
      text-align: center;
      color: #3a4559;
      .copy-icon {
        background-image: url(@/assets/img/copy-hover.png);
      }
    }
  }

  .operation {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .operation-right {
      display: flex;
      gap: 5px;
      align-items: center;
      justify-content: space-between;
    }
  }

  .praise-block {
    text-align: center;
    width: 26px;
    cursor: pointer;
    border-radius: 6px;

    .praise-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/praise.png);
      background-size: 26px;
      width: 26px;
      height: 26px;
      margin-right: 2px;
    }
    &:hover {
      background: #F2F4F7;
    }
  }

  .praise-active {
    .praise-icon {
      background-image: url(@/assets/img/praise-hover.png);
    }
  }

  .trample-block {
    text-align: center;
    width: 26px;
    cursor: pointer;
    border-radius: 6px;

    .trample-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/trample.png);
      background-size: 26px;
      width: 26px;
      height: 26px;
      margin-right: 2px;
    }
    &:hover {
      background: #F2F4F7;
    }
  }

  .trample-active {
    .trample-icon {
      background-image: url(@/assets/img/trample-hover.png);
    }
  }

  .hover-copy-tool-block {
    padding: 0;
    display: none;
    position: absolute;
    bottom: -12px;
    height: 28px;
    max-width: 89px;
    align-items: center;
    justify-content: center;
    background: #fff;
    border: 1px solid #d8dde6;
    border-radius: 4.5px;
    transition: all 0.5s ease;

    .copy-block{
      display: flex;
      align-items: center;
      justify-content: center;
      width: 26px;
      height: 26px;
      padding: 0;
    }

    .operation-right {
      display: flex;
      gap: 5px;
      align-items: center;
      justify-content: space-between;
    }
  }
}

.feedback-modal {
  padding: 16px;
  flex-shrink: 0;
  border-radius: 12px;
  background: var(--Conditional-pop-over, #FFF);
  box-shadow: 0 6px 30px 5px #0000000d, 0 16px 24px 2px #0000000a, 0 8px 10px -5px #00000014;

  .modal-title {
    color: #262626;
    text-align: left;
    font-family: "PingFang SC";
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
    margin-bottom: 4px;
  }

  .btn-box {
    margin-top: 12px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 8px;
  }
  
  .model-btn {
    width: 74px;
    height: 32px;
    display: flex;
    padding: 5px 16px;
    justify-content: center;
    align-items: center;
    border-radius: 6px;
  }
}

.text-input {
  width: 268px;
  height: 98px;
  display: flex;
  padding: 5px 12px;
  color: #666;
  flex-direction: column;
  align-items: flex-start;
  align-self: stretch;
  border-radius: 2px;
  border: 1px solid var(--Neutral-5, #D9D9D9);
  background: var(--Neutral-1, #FFF);
}

.answer-reference-box{
  border-top: 1px solid #EDEFF2;
  padding-top: 12px;
  margin-top: 16px;
  .title-block{
    color: #7a8699;
    font-size: 14px;
    line-height: 22px;
    font-weight: 400;
  }
  .list-item{
    cursor: pointer;
    margin-top: 8px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    color: #164799;
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
            <div class="operation" v-if="isShowCopy">
              <div @click="handleCopy" class="copy-block">
                <div class="copy-icon"></div>
                <span>复制</span>
              </div>
              <div ref="operationRef" class="operation-right" v-if="!isCustomerMessage && !isWelcomeMessage">
                <div @click="handlePraise(props.msg)" class="praise-block" :class="{'praise-active': !isTrampleClick && isPraiseActive}">
                  <div class="praise-icon"></div>
                </div>
                <van-popover placement="top-end" v-if="!isTrampleActive" @close="onCancel">
                  <div ref="feedbackModal" class="feedback-modal">
                    <div class="modal-title">结果反馈</div>
                    <div class="content">
                      <textarea
                        class="text-input"
                        v-model="feedbackContent"
                        placeholder="请反馈你觉得回答不满意的地方"
                      ></textarea>
                    </div>
                    <div class="btn-box">
                      <van-button class="model-btn cancel-btn" type="default" @click="onCancel">取消</van-button>
                      <van-button class="model-btn submit-btn" type="primary" @click="onSubmit">提交</van-button>
                    </div>
                  </div>
                  <template #reference>
                      <div @click="handleTrample(props.msg, 4)" class="trample-block" :class="{'trample-active': isTrampleClick}">
                        <div class="trample-icon"></div>
                      </div>
                  </template>
                </van-popover>
                <div v-else @click="handleTrample(props.msg, 3)" class="trample-block" :class="{'trample-active': isTrampleActive}">
                  <div class="trample-icon"></div>
                </div>
              </div>
            </div>
            <div class="answer-reference-box" v-if="props.msg.is_customer != 1 && props.msg.quote_file && props.msg.quote_file.length">
              <div class="title-block">回答参考</div>
              <div
                class="list-item"
                v-for="(item, index) in props.msg.quote_file"
                :key="index"
              >
                <svg-icon name="attachment"  />
                <span @click="handleToLink(item)">{{ item.file_name}}</span>
              </div>
            </div>
            <div class="hover-copy-tool-block" v-if="isShowHoverCopy">
              <div ref="operationRef" class="operation-right">
                <div v-tooltip="'复制'" @click="handleCopy" class="copy-block">
                  <div class="copy-icon"></div>
                </div>
                <template v-if="!isCustomerMessage && !isWelcomeMessage">
                  <div v-tooltip="'点赞'" @click="handlePraise(props.msg)" class="praise-block" :class="{'praise-active': !isTrampleClick && isPraiseActive}">
                    <div class="praise-icon"></div>
                  </div>
                  <van-popover placement="top-end" v-if="!isTrampleActive" @close="onCancel">
                    <div ref="feedbackModal" class="feedback-modal">
                      <div class="modal-title">结果反馈</div>
                      <div class="content">
                        <textarea
                          class="text-input"
                          v-model="feedbackContent"
                          placeholder="请反馈你觉得回答不满意的地方"
                        ></textarea>
                      </div>
                      <div class="btn-box">
                        <van-button class="model-btn cancel-btn" type="default" @click="onCancel">取消</van-button>
                        <van-button class="model-btn submit-btn" type="primary" @click="onSubmit">提交</van-button>
                      </div>
                    </div>
                    <template #reference>
                        <div v-tooltip="'点踩'" @click="handleTrample(props.msg, 2)" class="trample-block" :class="{'trample-active': isTrampleClick}">
                          <div class="trample-icon"></div>
                        </div>
                    </template>
                  </van-popover>
                  <div v-else v-tooltip="'点踩'" @click="handleTrample(props.msg, 1)" class="trample-block" :class="{'trample-active': isTrampleActive}">
                    <div class="trample-icon"></div>
                  </div>
                </template>
              </div>
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
    <quote-modal ref="quoteModalRef"></quote-modal>
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
import QuoteModal from '../quote-modal/index.vue'

interface praiseParams {
  ai_message_id: string;
  type: number;
  customer_message_id?: string;
}

const { toClipboard } = useClipboard()

const emit = defineEmits(['sendTextMessage'])
const chatStore = useChatStore()
const { robot, onAddFeedback, onDelFeedback } = chatStore
const textMessage = ref('.')
const feedbackContent = ref('')
const ai_message_id = ref('')
const customer_message_id = ref<any>('')
const operationRef = ref<HTMLElement | null>(null)
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
  },
  prevMsg: {
    type: Object as () => Message
  }
})

const isTrampleClick = ref(false)
const isPraiseActive = ref(props.msg.feedback_type == '1' ? true : false);
const isTrampleActive = ref(props.msg.feedback_type == '2' ? true : false);

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
  isTrampleClick.value = false
  await toClipboard(props.msg.content)
  showToast('复制成功')
}

const handlePraise = async (item) => {
  isTrampleClick.value = false
  if (props.prevMsg && props.prevMsg.is_customer == 1) {
    customer_message_id.value = props.prevMsg.id
  }
  ai_message_id.value = item.id
  const params:praiseParams  = {
    ai_message_id: ai_message_id.value,
    type: 1
  }
  if (!isPraiseActive.value) {
    params.customer_message_id = customer_message_id.value
    await onAddFeedback(params)
    isTrampleClick.value = false
    isTrampleActive.value = false
    isPraiseActive.value = true
  } else {
    await onDelFeedback(params)
    isPraiseActive.value = false
  }
  // showToast('操作成功')
}

const handleTrample = async (item, type) => {
  console.log(type)
  isTrampleClick.value = true
  if (props.prevMsg && props.prevMsg.is_customer == 1) {
    customer_message_id.value = props.prevMsg.id
  }
  ai_message_id.value = item.id
  if (isTrampleActive.value) {
    // 取消踩
    const params = {
      ai_message_id: ai_message_id.value
    }
    await onDelFeedback(params)
    isTrampleActive.value = false
    isTrampleClick.value = false
    showToast('操作成功')
  }
}

const onCancel = async () => {
  const params = {
    ai_message_id: ai_message_id.value,
    customer_message_id: customer_message_id.value,
    type: 2
  }
  await onAddFeedback(params)
  isTrampleActive.value = true
  // operationRef.value!.click()
  feedbackContent.value = ''
  isPraiseActive.value = false
  // showToast('操作成功')
}

const onSubmit = async () => {
  const params = {
    ai_message_id: ai_message_id.value,
    customer_message_id: customer_message_id.value,
    type: 2,
    content: feedbackContent.value
  }
  await onAddFeedback(params)
  isTrampleActive.value = true
  // operationRef.value!.click()
  feedbackContent.value = ''
  isPraiseActive.value = false
  showToast('感谢反馈')
}

// 检查是否为用户消息
const isCustomerMessage = computed(() => props.msg.is_customer == 1)

// 是否为欢迎语
const isWelcomeMessage = computed(() => props.msg.msg_type == 2)

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
const quoteModalRef = ref<any>(null)
const handleToLink = (item) => {
  quoteModalRef.value && quoteModalRef.value.showPopup({
    message_id: props.msg.id,
    file_id: item.id,
    robot_key: robot.robot_key,
    openid: robot.openid,
    file_name: item.file_name,
  })
}

onMounted(() => {
  startLoadingAnimation()
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>
