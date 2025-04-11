<style lang="less" scoped>
.ignore-message-item {
  max-width: 760px;
  display: flex;
  padding: 48px 12px 0px;
  margin: 0 auto;

  .message-item-body {
    flex: 1;
    padding-left: 10px;
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
    &:hover{
      .hover-copy-tool-block{
        opacity: 1;
      }
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
    margin-top: 0px;
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
    .hover-copy-tool-block {
      bottom: -35px;
      left: 8px;
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
    .hover-copy-tool-block {
      left: -30px;
      width: 36px;
      bottom: 0px;
      height: 38px;
      display: flex;
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
  }

  .operation {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 12px;
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
    opacity: 0;
    position: absolute;
    max-width: 90px;
    align-items: center;
    justify-content: center;
    transition: all 0.5s ease;

    .copy-block{
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      width: 26px;
    }

    .operation-right {
      display: flex;
      gap: 5px;
      align-items: center;
      justify-content: space-between;
    }
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
}
.ignore-message-item:last-child {
  padding-bottom: 24px;
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
.thinking-label-wrapper {
  display: flex;
  align-items: center;
  margin-bottom: 8px;

  .thinking-label {
    display: flex;
    align-items: center;
    height: 32px;
    padding: 0 16px;
    border-radius: 8px;
    background: #e4e6eb;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #d8dde6;
    }
    .think-icon,
    .loading {
      margin-right: 8px;
      font-size: 16px;
      color: #262626;
    }
    .label-text {
      font-size: 14px;
      font-weight: 400;
      color: #262626;
    }
    .arrow-down {
      margin-left: 8px;
      font-size: 16px;
      color: #262626;
      cursor: pointer;
    }
  }
  .tip {
    margin-left: 8px;
    font-size: 16px;
    color: #8c8c8c;
    cursor: pointer;
  }

  &.reasoning_open {
    .down-arrow {
      transform: rotate(180deg);
    }
  }
}
.thinking-content {
  position: relative;
  line-height: 22px;
  padding-bottom: 0;
  padding-left: 16px;
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 400;
  color: #8c8c8c;
  border-bottom: 1px solid #edeff2;
  // 4px竖线
  &::before {
    display: block;
    position: absolute;
    content: '';
    left: 0;
    top: 4px;
    bottom: 20px;
    width: 4px;
    background-color: #d9d9d9;
  }
}
</style>

<template>
  <div class="ignore-message-item" :class="messageItemClasses" :id="'msg-' + msg.uid">
    <div class="message-item-left">
      <img class="avatar" :src="props.msg.avatar" />
    </div>
    <div class="message-item-body">
      <!-- 思考过程label -->
      <div
        class="thinking-label-wrapper"
        :class="{ reasoning_open: props.msg.show_reasoning }"
        v-if="props.msg.reasoning_content && props.msg.reasoning_content.length > 0"
      >
        <div class="thinking-label" @click="toggleReasonProcess()">
          <van-loading class="loading" color="#262626" size="16px" type="spinner" v-if="props.msg.reasoning_status" />
          <svg-icon class="think-icon" name="think" v-else></svg-icon>
          <span class="label-text">{{
            props.msg.reasoning_status ? '深度思考中...' : '已完成深度思考'
          }}</span>

          <svg-icon
            name="down-arrow"
            class="down-arrow"
            v-if="!props.msg.reasoning_status"
          ></svg-icon>
        </div>
      </div>
      <div class="message-content">
        <!-- <span class="triangle"></span> -->
        <div class="thinking-content" v-if="props.msg.show_reasoning">
            <cherry-markdown :content="props.msg.reasoning_content" />
        </div>
        <template v-if="props.msg.msg_type == 1">
          <div class="text-message" v-if="props.msg.content !== ''" v-viewer>
            <div v-if="props.msg.is_customer == 1" v-html="props.msg.content"></div>
            <cherry-markdown :content="props.msg.content" v-else />
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
          <div class="operation" v-if="isShowCopy">
            <div @click="handleCopy" class="copy-block">
              <div class="copy-icon"></div>
              <span>复制</span>
            </div>
            <div ref="operationRef" class="operation-right" v-if="!isCustomerMessage && !isWelcomeMessage && robot.feedback_switch">
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
                    <div v-tooltip="'点踩'" @click="handleTrample(props.msg, 4)" class="trample-block" :class="{'trample-active': isTrampleClick}">
                      <div class="trample-icon"></div>
                    </div>
                </template>
              </van-popover>
              <div v-tooltip="'点踩'" v-else @click="handleTrample(props.msg, 3)" class="trample-block" :class="{'trample-active': isTrampleActive}">
                <div class="trample-icon"></div>
              </div>
            </div>
          </div>
           

          <div class="answer-reference-box" v-if="props.msg.is_customer != 1 && props.msg.quote_file && props.msg.quote_file.length">
            <div class="title-block">回答参考
              <van-icon v-tooltip="'如果不想显示“回答参考”，可以在机器人配置界面“基础配置 - 显示引文”，关闭开关。'" name="question-o" />
            </div>
            <div
              class="list-item"
              v-for="(item, index) in props.msg.quote_file"
              :key="index"
            >
              <svg-icon name="attachment" />
              <span @click="handleToLink(item)">{{ item.file_name}}</span>
            </div>
          </div>
          <div class="hover-copy-tool-block" v-if="isShowHoverCopy">
            <div ref="operationRef" class="operation-right">
              <div v-tooltip="'复制'" @click="handleCopy" class="copy-block">
                <div class="copy-icon"></div>
              </div>
                <template v-if="!isCustomerMessage && !isWelcomeMessage && robot.feedback_switch">
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
    <quote-modal ref="quoteModalRef"></quote-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import type { Message } from '@/stores/modules/chat'
import { useChatStore } from '@/stores/modules/chat'
import { escapeHTML } from '@/utils/index'
import GuessYouWant from './guess-you-want.vue'
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
    isTrampleActive.value = false
    isTrampleClick.value = false
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

const common_question_list = computed(() => robot.common_question_list)
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

const toggleReasonProcess = () => {
  props.msg.show_reasoning = ! props.msg.show_reasoning
}

onMounted(() => {
  startLoadingAnimation()
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>
