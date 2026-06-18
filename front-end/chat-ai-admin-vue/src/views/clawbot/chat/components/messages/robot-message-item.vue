<template>
  <div class="message-item robot-message" :id="'msg-' + item.uid">
    <div class="itme-left">
      <a-spin size="small" :spinning="item.loading">
        <img class="robot-avatar" :src="item.robot_avatar" />
      </a-spin>
    </div>

    <div class="itme-right">
      <ReplyContentList
        v-if="hasReplyList"
        :reply-content-list="item.reply_content_list"
        @clickMsgMeun="onClickMeun"
      />

      <div class="label-flex-block">
        <div class="thinking-label-wrapper" v-if="tipsBeforeAnswerSwitch && item.startLoading">
          <div class="thinking-label">
            <LoadingOutlined class="loading" />
            <span class="label-text">{{ tipsBeforeAnswerContent }}</span>
          </div>
        </div>
        <div
          class="thinking-label-wrapper"
          :class="{ reasoning_open: item.show_quote_file }"
          v-if="item.msg_type == 1 && isShowQuoteFileProgress && (!item.startLoading || !tipsBeforeAnswerSwitch)"
        >
          <div class="thinking-label" @click="toggleQuoteFiel(item)">
            <template v-if="item.quote_loading">
              <LoadingOutlined class="loading" />
              <span class="label-text">{{ t('label_retrieving_knowledge_base') }}</span>
            </template>
            <template v-if="!item.quote_loading">
              <svg-icon class="think-icon" name="quote-file"></svg-icon>
              <span class="label-text">
                {{ t('label_retrieved_knowledge_base_docs', { count: item.quote_file.length }) }}
              </span>
            </template>
            <svg-icon
              name="arrow-down"
              class="arrow-down"
              v-if="item.quote_file.length"
            ></svg-icon>
          </div>
        </div>

        <div
          class="thinking-label-wrapper"
          :class="{ reasoning_open: isReasoningExpanded(item) }"
          v-if="hasReasoningContent(item) && !hasProcessSteps(item)"
        >
          <div class="thinking-label" @click="toggleReasonProcess(item)">
            <LoadingOutlined class="loading" v-if="item.reasoning_status" />
            <svg-icon class="think-icon" name="think" v-else></svg-icon>
            <span class="label-text">
              {{ item.reasoning_status ? t('label_thinking') : t('label_thinking_completed') }}
            </span>
            <svg-icon
              name="arrow-down"
              class="arrow-down"
              v-if="!item.reasoning_status"
            ></svg-icon>
          </div>
          <a-tooltip>
            <template #title>{{ t('msg_disable_reasoning_tooltip') }}</template>
            <InfoCircleOutlined class="tip" />
          </a-tooltip>
        </div>
      </div>

      <div class="item-body" :class="{ 'item-body-final': item.msg_type == 1 }" v-if="isShowBody">
        <QuoteFilePanel :item="item" />
        <ProcessTimeline v-if="hasProcessSteps(item)" :item="item" />
        <div class="thinking-content" v-if="isReasoningExpanded(item) && hasReasoningContent(item) && !hasProcessSteps(item)">
          <CherryMarkdown :content="item.reasoning_content"></CherryMarkdown>
        </div>

        <FinalAnswerCard
          v-if="item.msg_type == 1"
          :item="item"
          :robot-info="robotInfo"
          :should-indent="hasProcessSteps(item)"
          @clickMsgMeun="onClickMeun"
        />

        <template v-if="item.msg_type == 2">
          <CherryMarkdown
            class="markdown-content"
            :content="item.menu_json.content"
            v-if="item.isWelcome"
          ></CherryMarkdown>
          <div class="message-content" v-html="item.menu_json.content" v-else></div>
          <div class="message-menus">
            <div
              class="menu-item"
              @click="onClickMeun(question)"
              v-for="(question, index) in item.menu_json.question"
              :key="index"
            >
              {{ question }}
            </div>
          </div>
        </template>

        <template v-if="item.msg_type == 3">
          <div class="message-content">
            <img v-viewer class="msg-img" :src="item.content" alt="" />
          </div>
        </template>
      </div>

      <div class="message-loading" v-if="showMessageLoading">
        <span class="loading-text"></span>
        <span class="loading-dots">
          <span></span>
          <span></span>
          <span></span>
        </span>
      </div>

      <div v-if="item.voice_content && item.voice_content.length">
        <VoiceMessage :voice_content="item.voice_content" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { LoadingOutlined, InfoCircleOutlined } from '@ant-design/icons-vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import { useI18n } from '@/hooks/web/useI18n'
import VoiceMessage from '../voice-message.vue'
import ReplyContentList from './reply-content-list.vue'
import ProcessTimeline from './process-timeline.vue'
import FinalAnswerCard from './final-answer-card.vue'
import QuoteFilePanel from './quote-file-panel.vue'

const { t } = useI18n('views.clawbot.chat.components.message-list')

const emit = defineEmits(['clickMsgMeun'])

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  },
  robotInfo: {
    type: Object,
    default: () => ({})
  },
  tipsBeforeAnswerContent: {
    type: String,
    default: ''
  },
  tipsBeforeAnswerSwitch: {
    type: Boolean,
    default: false
  },
  isShowQuoteFileProgress: {
    type: Boolean,
    default: false
  }
})

const parseReplyList = (val) => {
  try {
    if (!val) return []
    if (Array.isArray(val)) return val
    if (typeof val === 'string') return JSON.parse(val || '[]')
    return []
  } catch (_e) {
    return []
  }
}

const hasReplyList = computed(() => parseReplyList(props.item.reply_content_list).length > 0)

const hasReasoningContent = (item) => {
  return typeof item?.reasoning_content === 'string'
    ? item.reasoning_content.trim().length > 0
    : !!item?.reasoning_content
}

const getVisibleProcessSteps = (item) => {
  if (!Array.isArray(item?.process_steps)) {
    return []
  }
  return item.process_steps.filter((step) => {
    if (step?.hidden === true) {
      return false
    }
    if (step?.type === 'tool') {
      return true
    }
    return hasReasoningContent(item)
  })
}

const hasProcessSteps = (item) => {
  return getVisibleProcessSteps(item).length > 0
}

const hasFinalAnswerContent = (item) => {
  if (!item) {
    return false
  }
  if (item.msg_type == 2) {
    return !!item.menu_json?.content
  }
  if (item.msg_type == 3) {
    return !!item.content
  }
  if (typeof item.content === 'string') {
    return item.content.trim().length > 0
  }
  return !!item.content
}

const showMessageLoading = computed(() => {
  return !!props.item?.loading && !hasFinalAnswerContent(props.item)
})

const isReasoningExpanded = (item) => {
  if (!item) {
    return false
  }
  if (typeof item.reasoning_expanded === 'boolean') {
    return item.reasoning_expanded
  }
  return !!item.show_reasoning
}

const toggleReasonProcess = (item) => {
  item.reasoning_expanded = !isReasoningExpanded(item)
}

const toggleQuoteFiel = (item) => {
  item.show_quote_file = !item.show_quote_file
}

const isShowBody = computed(() => {
  if (hasProcessSteps(props.item)) {
    return true
  }
  if (hasReplyList.value && props.item.content == '') {
    return false
  }
  return true
})

const onClickMeun = (item) => {
  emit('clickMsgMeun', item)
}
</script>

<style lang="less" scoped>
.markdown-content {
  white-space: normal;
}

.robot-message {
  .item-body {
    padding: 0;
    width: auto;
    min-height: 32px;
    max-width: 100%;
    overflow: hidden;
    background: transparent;
  }

  .item-body.item-body-final {
    padding: 0;
    min-height: auto;
    background: transparent;
  }
}

.message-content {
  line-height: 22px;
  font-size: 14px;
  font-weight: 400;
  color: #262626;
  white-space: pre-wrap;
  word-break: break-word;
}

.message-loading {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding-left: 0;
  color: #8c8c8c;
  font-size: 13px;
  line-height: 20px;

  .loading-dots {
    display: inline-flex;
    align-items: center;
    gap: 4px;

    span {
      width: 4px;
      height: 4px;
      border-radius: 50%;
      background: currentColor;
      opacity: 0.3;
      animation: loading-dot-blink 1.2s infinite ease-in-out;

      &:nth-child(2) {
        animation-delay: 0.2s;
      }

      &:nth-child(3) {
        animation-delay: 0.4s;
      }
    }
  }
}

@keyframes loading-dot-blink {
  0%, 80%, 100% {
    opacity: 0.25;
    transform: translateY(0);
  }

  40% {
    opacity: 1;
    transform: translateY(-1px);
  }
}

.msg-img {
  width: 100%;
  height: 100%;
  max-width: 300px;
  max-height: 300px;
}

.message-menus {
  .menu-item {
    line-height: 22px;
    padding: 8px 16px;
    margin-top: 8px;
    font-size: 14px;
    border-radius: 4px;
    color: rgb(22, 71, 153);
    background: #f2f4f7;
    cursor: pointer;
  }
}

.label-flex-block {
  display: flex;
  align-items: center;
  gap: 8px;
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
    .arrow-down {
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
