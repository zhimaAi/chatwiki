<template>
  <div class="message-item robot-message" :id="'msg-' + item.uid">
    <div class="itme-left">
      <img class="robot-avatar" :src="item.robot_avatar" />
    </div>

    <div class="itme-right">
      <ReplyContentList
        v-if="hasReplyList"
        :reply-content-list="item.reply_content_list"
        @clickMsgMeun="onClickMeun"
      />

      <div class="label-flex-block">
        <div
          class="thinking-label-wrapper"
          :class="{ reasoning_open: item.show_quote_file }"
          v-if="shouldShowQuoteFileProgress(item)"
        >
          <div class="thinking-label" @click="toggleQuoteFiel(item)">
            <template v-if="item.quote_loading">
              <LoadingOutlined class="loading" />
              <span class="label-text">{{ t('label_retrieving_knowledge_base') }}</span>
            </template>
            <template v-if="!item.quote_loading">
              <svg-icon class="think-icon" name="quote-file"></svg-icon>
              <span class="label-text">
                {{ t('label_retrieved_knowledge_base_docs', { count: getQuoteFileCount(item) }) }}
              </span>
            </template>
            <svg-icon
              name="arrow-down"
              class="arrow-down"
              v-if="getQuoteFileCount(item) > 0"
            ></svg-icon>
          </div>
        </div>

        <div class="stopped-label" v-if="showStoppedLabel">
          {{ t('label_stopped') }}
        </div>
      </div>

      <div class="item-body" :class="{ 'item-body-final': item.msg_type == 1 }" v-if="isShowBody">
        <QuoteFilePanel :item="item" />
        <ProcessTimeline
          v-if="shouldShowProcessTimeline(item)"
          :item="item"
          :running-label="processRunningLabel"
          :quote-loading-visible="shouldShowQuoteFileProgress(item) && item.quote_loading"
        />
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

      <div v-if="item.voice_content && item.voice_content.length">
        <VoiceMessage :voice_content="item.voice_content" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { LoadingOutlined } from '@ant-design/icons-vue'
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

const getVisibleProcessSteps = (item) => {
  if (!Array.isArray(item?.process_steps)) {
    return []
  }
  return item.process_steps.filter((step) => {
    if (step?.hidden === true) {
      return false
    }
    if (step?.type === 'tool' || step?.type === 'skill') {
      return !!step.title || !!step.paramsText || !!step.resultText || step.status === 'running'
    }
    return !!step?.contentText || !!step?.resultText || step?.status === 'running'
  })
}

const hasProcessSteps = (item) => {
  return getVisibleProcessSteps(item).length > 0
}

const processRunningLabel = computed(() => {
  if (!props.tipsBeforeAnswerSwitch) {
    return ''
  }
  return props.tipsBeforeAnswerContent?.trim() || ''
})

const shouldShowProcessTimeline = (item) => {
  return item?.msg_type == 1 && (
    hasProcessSteps(item) || (!!item.startLoading && !item.is_stopped)
  )
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

const showStoppedLabel = computed(() => {
  return !!props.item?.is_stopped &&
    !hasReplyList.value &&
    !hasProcessSteps(props.item) &&
    !hasFinalAnswerContent(props.item)
})

const toggleQuoteFiel = (item) => {
  item.show_quote_file = !item.show_quote_file
}

const getQuoteFileCount = (item) => {
  return Array.isArray(item?.quote_file) ? item.quote_file.length : 0
}

const shouldShowQuoteFileProgress = (item) => {
  return item?.msg_type == 1 &&
    props.isShowQuoteFileProgress &&
    (!item.startLoading || !props.tipsBeforeAnswerSwitch) &&
    (!item.is_stopped || item.quote_loading || getQuoteFileCount(item) > 0)
}

const isShowBody = computed(() => {
  if (shouldShowProcessTimeline(props.item)) {
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

.stopped-label {
  display: flex;
  align-items: center;
  height: 32px;
  padding: 0 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 400;
  color: #595959;
  background: #e4e6eb;
  margin-bottom: 8px;
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

  &.reasoning_open {
    .arrow-down {
      transform: rotate(180deg);
    }
  }
}

</style>
