<template>
  <div class="process-timeline" v-if="displaySteps.length">
    <div class="process-timeline-header" @click="toggleTimeline">
      <span>{{ timelineTitle }}</span>
      <svg-icon name="arrow-down" class="timeline-arrow" :class="{ expanded: timelineExpanded }"></svg-icon>
    </div>

    <div class="process-timeline-body" v-if="timelineExpanded">
      <div class="process-step" v-for="(step, stepIndex) in displaySteps" :key="step.id">
        <div class="process-step-side">
          <span class="process-step-dot" :class="{ 'is-last': stepIndex === displaySteps.length - 1 }"></span>
          <span class="process-step-line" v-if="stepIndex < displaySteps.length - 1"></span>
        </div>
        <div class="process-step-content">
          <div class="process-step-title">{{ getStepTitle(step, stepIndex) }}</div>
          <div class="process-step-text chat-markdown-content" v-if="getStepText(step)">
            <CherryMarkdown
              class="markdown-body"
              :content="getStepText(step)"
              enable-image-preview
            ></CherryMarkdown>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import '@/assets/style/markdown/markdown.less'

const { t } = useI18n('views.clawbot.chat.components.messages.process-timeline')
const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
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
    if (step?.type === 'tool') {
      return !!step.title || !!step.resultText
    }
    if (step?.type === 'thinking') {
      return !!step.contentText || !!step.resultText || step.status === 'running'
    }
    return !!step.resultText || !!step.contentText
  })
})

const hasRunningStep = computed(() => displaySteps.value.some((step) => step?.status === 'running'))
const timelineExpanded = computed(() => {
  if (hasRunningStep.value) {
    return true
  }
  return displaySteps.value.some((step) => step?.expanded === true)
})
const timelineTitle = computed(() => hasRunningStep.value ? t('label_thinking') : t('label_thinking_completed'))

const setTimelineExpanded = (expanded) => {
  visibleSteps.value.forEach((step) => {
    step.expanded = expanded
  })
}

const toggleTimeline = () => {
  setTimelineExpanded(!timelineExpanded.value)
}

const getStepTitle = (step, stepIndex) => {
  if (step?.type === 'thinking') {
    const firstThinkingIndex = displaySteps.value.findIndex((item) => item?.type === 'thinking')
    return stepIndex === firstThinkingIndex ? t('label_thinking_process') : t('label_thinking_again')
  }
  return step?.title || step?.eventName || ''
}

const getStepText = (step) => {
  if (step?.type === 'tool') {
    return step.resultText
  }
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
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #8c8c8c;
  cursor: pointer;
  user-select: none;
}

.timeline-arrow {
  font-size: 14px;
  color: #8c8c8c;
  transform: rotate(-90deg);
  transition: transform 0.2s;

  &.expanded {
    transform: rotate(0deg);
  }
}

.process-timeline-body {
  margin-top: 12px;
}

.process-step {
  display: flex;
  align-items: stretch;
  gap: 12px;
}

.process-step-side {
  position: relative;
  flex: 0 0 12px;
  display: flex;
  justify-content: center;
}

.process-step-dot {
  position: relative;
  z-index: 1;
  width: 8px;
  height: 8px;
  margin-top: 7px;
  border-radius: 50%;
  background: #d9d9d9;

  &.is-last {
    width: 8px;
    height: 8px;
    border: 1px solid #d9d9d9;
    background: #fff;
  }
}

.process-step-line {
  position: absolute;
  top: 15px;
  bottom: -7px;
  left: 50%;
  width: 1px;
  background: #d9d9d9;
  transform: translateX(-50%);
}

.process-step-content {
  flex: 1;
  min-width: 0;
  padding-bottom: 12px;
  color: #8c8c8c;
  word-break: break-word;
}

.process-step-title {
  margin-bottom: 4px;
}

.process-step-text.chat-markdown-content {
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
</style>
