<template>
  <div class="process-timeline" v-if="showTimeline">
    <div class="process-timeline-header" @click="toggleTimeline">
      <LoadingOutlined class="timeline-loading" v-if="showHeaderLoading" />
      <span>{{ timelineTitle }}</span>
      <svg-icon
        name="arrow-down"
        class="timeline-arrow"
        :class="{ expanded: timelineExpanded }"
        v-if="displaySteps.length"
      ></svg-icon>
    </div>

    <div class="process-timeline-body" v-if="timelineExpanded && displaySteps.length">
      <div class="process-step" v-for="step in displaySteps" :key="step.id">
        <div class="process-step-header" @click="toggleStep(step)">
          <LoadingOutlined class="process-step-status is-running" v-if="step.status === 'running'" />
          <CheckCircleOutlined class="process-step-status is-done" v-else />
          <div class="process-step-title">
            <span
              class="process-step-thinking-summary"
              v-if="step.type === 'thinking' && !isStepExpanded(step)"
            >
              {{ getStepText(step) }}
            </span>
            <CherryMarkdown
              v-else-if="step.type === 'thinking'"
              class="process-step-thinking-content"
              :content="getStepText(step)"
              enable-image-preview
            ></CherryMarkdown>
            <template v-else>{{ getStepTitle(step) }}</template>
          </div>
          <svg-icon
            name="arrow-down"
            class="process-step-arrow"
            :class="{ expanded: isStepExpanded(step) }"
          ></svg-icon>
        </div>

        <div
          class="process-step-detail"
          v-if="isStepExpanded(step) && step.type !== 'thinking'"
        >
          <div class="process-step-section" v-if="step.type === 'tool' && step.paramsText">
            <span class="process-step-label">{{ t('label_input') }}</span>
            <div class="process-step-scroll">
              <code class="process-step-code">{{ step.paramsText }}</code>
            </div>
          </div>

          <div class="process-step-section" v-if="step.resultText">
            <span class="process-step-label">{{ t('label_output') }}</span>
            <div class="process-step-scroll chat-markdown-content">
              <CherryMarkdown
                class="markdown-body"
                :content="step.resultText"
                enable-image-preview
              ></CherryMarkdown>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { CheckCircleOutlined, LoadingOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import '@/assets/style/markdown/markdown.less'

const { t } = useI18n('views.clawbot.chat.components.messages.process-timeline')
const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  },
  runningLabel: {
    type: String,
    default: ''
  },
  quoteLoadingVisible: {
    type: Boolean,
    default: false
  }
})

const visibleSteps = computed(() => {
  if (!Array.isArray(props.item?.process_steps)) {
    return []
  }
  return props.item.process_steps.filter((step) => step?.hidden !== true)
})

const displaySteps = computed(() => {
  return visibleSteps.value.filter((step) => {
    if (step?.type === 'tool' || step?.type === 'skill') {
      return !!step.title || !!step.paramsText || !!step.resultText || step.status === 'running'
    }
    if (step?.type === 'thinking') {
      return !!step.contentText || !!step.resultText || step.status === 'running'
    }
    return !!step.resultText || !!step.contentText
  })
})

const hasRunningStep = computed(() => displaySteps.value.some((step) => step?.status === 'running'))
const isWaitingForAnswer = computed(() => !!props.item?.startLoading && !props.item?.is_stopped)
const showTimeline = computed(() => displaySteps.value.length > 0 || isWaitingForAnswer.value)
const showHeaderLoading = computed(() => {
  return isWaitingForAnswer.value && !hasRunningStep.value && !props.quoteLoadingVisible
})
const timelineExpanded = ref(props.item?.process_expanded !== false)
const stepExpandedOverrides = ref({})
const timelineTitle = computed(() => {
  if (props.item?.is_stopped) {
    return t('label_stopped')
  }
  if (isWaitingForAnswer.value || hasRunningStep.value) {
    return props.runningLabel || t('label_thinking')
  }
  return t('label_thinking_completed')
})

const toggleTimeline = () => {
  if (!displaySteps.value.length) {
    return
  }
  timelineExpanded.value = !timelineExpanded.value
}

const isStepExpanded = (step) => {
  const stepId = step?.id
  if (stepId && Object.prototype.hasOwnProperty.call(stepExpandedOverrides.value, stepId)) {
    return stepExpandedOverrides.value[stepId]
  }
  // 未手动切换的思考步骤在完成后默认收起，显式的用户选择优先。
  if (step?.type === 'thinking' && step?.status === 'done') {
    return false
  }
  return step?.expanded === true
}

const toggleStep = (step) => {
  if (!step?.id) {
    return
  }
  stepExpandedOverrides.value = {
    ...stepExpandedOverrides.value,
    [step.id]: !isStepExpanded(step)
  }
}

const getStepTitle = (step) => {
  if (step?.type === 'thinking') {
    return getStepText(step)
  }
  if (step?.type === 'skill') {
    return t('label_execute_skill', { name: step?.title || '' })
  }
  return t('label_execute_tool', { name: step?.title || step?.eventName || '' })
}

const getStepText = (step) => {
  return step?.contentText || step?.resultText || (step?.status === 'running' ? t('msg_thinking') : '')
}
</script>

<style lang="less" scoped>
.process-timeline {
  margin-bottom: 24px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.process-timeline-header {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
  padding: 8px 0;
  color: #595959;
  border-bottom: 1px solid #d9d9d9;
  cursor: pointer;
  user-select: none;
}

.timeline-loading {
  color: #2475fc;
}

.timeline-arrow {
  font-size: 14px;
  color: #595959;
  transform: rotate(-90deg);
  transition: transform 0.2s;

  &.expanded {
    transform: rotate(0deg);
  }
}

.process-timeline-body {
  padding-top: 8px;
}

.process-step {
  width: 100%;
}

.process-step-header {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  min-height: 38px;
  cursor: pointer;
  user-select: none;
}

.process-step-status {
  flex: 0 0 auto;
  width: 16px;
  height: 16px;
  margin-top: 11px;
  font-size: 16px;

  &.is-running {
    color: #2475fc;
  }

  &.is-done {
    color: #8c8c8c;
  }
}

.process-step-title {
  min-width: 0;
  padding: 8px 0;
  color: #8c8c8c;
  word-break: break-word;
}

.process-step-thinking-summary {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.process-step-arrow {
  flex: 0 0 auto;
  margin-top: 12px;
  font-size: 14px;
  color: #8c8c8c;
  transform: rotate(-90deg);
  transition: transform 0.2s;

  &.expanded {
    transform: rotate(0deg);
  }
}

.process-step-detail {
  padding: 0 0 12px 22px;
}

.process-step-thinking-content {
  :deep(.cherry-markdown) {
    color: #8c8c8c;
    background: transparent;
    margin: 0;
    padding: 0;

    *{
      line-height: 22px;
      font-size: 14px;
    }
    p,
    li,
    blockquote {
      color: #8c8c8c;
    }

    h1, h2, h3, h4, h5, h6 {
      padding: 0;
      margin: 0;
    }
    p,
    ul,
    ol,
    blockquote,
    pre,
    table {
      margin-top: 0;
      margin-bottom: 4px;
    }

    > :last-child {
      margin-bottom: 0 !important;
    }
  }
}

.process-step-section {
  & + & {
    margin-top: 12px;
  }
}

.process-step-label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  color: #595959;
}

.process-step-scroll {
  max-height: 200px;
  padding: 8px 12px;
  overflow-y: auto;
  background: #f5f6f8;
  border-radius: 6px;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: #d9d9d9;
    border-radius: 3px;
  }
}

.process-step-code {
  display: block;
  font-size: 12px;
  line-height: 20px;
  white-space: pre-wrap;
  word-break: break-all;
  color: #595959;
}

.process-step-scroll.chat-markdown-content {
  :deep(.cherry-markdown) {
    color: #595959;
    background: transparent;
    margin: 0;
    padding: 0;

    * {
      line-height: 20px;
      font-size: 13px;
    }

    p,
    li,
    blockquote {
      color: #595959;
    }

    h1, h2, h3, h4, h5, h6 {
      padding: 0;
      margin: 0;
    }

    p,
    ul,
    ol,
    blockquote,
    pre,
    table {
      margin-top: 0;
      margin-bottom: 4px;
    }

    > :last-child {
      margin-bottom: 0 !important;
    }
    
    pre {
      border: none;
    }
  }
}
</style>
