<template>
  <div class="process-timeline" v-if="displaySteps.length">
    <div class="process-timeline-header-area" @click="toggleTimeline">
      <div class="process-timeline-header">
        <span>{{ timelineTitle }}</span>
        <svg-icon name="down-arrow" class="timeline-arrow" :class="{ expanded: timelineExpanded }"></svg-icon>
      </div>
    </div>

    <div class="process-timeline-body" v-if="timelineExpanded">
      <div class="process-step" v-for="(step, stepIndex) in displaySteps" :key="step.id">
        <div class="process-step-side">
          <span class="process-step-dot" :class="{ 'is-last': stepIndex === displaySteps.length - 1 }"></span>
          <span class="process-step-line" v-if="stepIndex < displaySteps.length - 1"></span>
        </div>
        <div class="process-step-content">
          <div class="process-step-title">{{ getStepTitle(step, stepIndex) }}</div>
          <template v-if="step.type === 'tool'">
            <div class="process-step-params" v-if="step.paramsText">
              <span class="process-step-label">{{ t('label_params') }}：</span>
              <code class="process-step-code">{{ step.paramsText }}</code>
            </div>
            <div class="process-step-result" v-if="step.resultText || step.status === 'running'">
              <span class="process-step-label">{{ t('label_output') }}：</span>
              <CherryMarkdown :content="normalizeText(step.resultText || (step.status === 'running' ? t('process_tool_running') : ''))" />
            </div>
          </template>
          <template v-else-if="getStepText(step)">
            <div class="process-step-text">
              <CherryMarkdown :content="getStepText(step)" />
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.messages.message-item')

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
      return !!step.title || !!step.resultText || step.status === 'running'
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
const timelineTitle = computed(() => hasRunningStep.value ? t('title_deep_thinking') : t('title_deep_thinking_completed'))

const setTimelineExpanded = (expanded: boolean) => {
  visibleSteps.value.forEach((step) => {
    step.expanded = expanded
  })
}

const toggleTimeline = () => {
  setTimelineExpanded(!timelineExpanded.value)
}

const getStepTitle = (step: any, stepIndex: number) => {
  if (step?.type === 'thinking') {
    const firstThinkingIndex = displaySteps.value.findIndex((item) => item?.type === 'thinking')
    return stepIndex === firstThinkingIndex ? t('process_thinking_step') : t('process_thinking_again')
  }
  return step?.title || t('process_tool_step')
}

const normalizeText = (value: any) => {
  if (value == null) {
    return ''
  }
  const text = typeof value === 'string' ? value : JSON.stringify(value)
  return text.length > 600 ? `${text.slice(0, 600)}...` : text
}

const getStepText = (step: any) => {
  return step?.contentText || step?.resultText || (step?.status === 'running' ? t('title_deep_thinking') : '')
}
</script>

<style lang="less" scoped>
.process-timeline {
  color: #666;
  font-size: 14px;
  line-height: 22px;
}

.process-timeline-header-area {
  cursor: pointer;
  user-select: none;
}

.process-timeline-header {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  min-height: 22px;
  padding: 0;
  color: #666;
}

.timeline-arrow {
  font-size: 14px;
  color: #666;
  transition: transform 0.2s;

  &.expanded {
    transform: rotate(180deg);
  }
}

.process-timeline-body {
  margin-top: 16px;
  padding-left: 0;
}

.process-step {
  display: flex;
  align-items: stretch;
  gap: 8px;
}

.process-step-side {
  position: relative;
  flex: 0 0 14px;
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
}

.process-step-line {
  position: absolute;
  top: 15px;
  bottom: -12px;
  left: 50%;
  width: 1px;
  background: #d9d9d9;
  transform: translateX(-50%);
}

.process-step-content {
  flex: 1;
  min-width: 0;
  padding-bottom: 12px;
  word-break: break-word;
}

.process-step-title {
  margin-bottom: 4px;
  color: #666;
}

.process-step-text {
  color: #666;

  :deep(.cherry-markdown) {
    color: #666;
    background: transparent;
    margin: 0;
    padding: 0;
    font-size: 14px;
    line-height: 22px;

    * {
      font-size: 14px;
      line-height: 22px;
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

.process-step-params {
  margin-bottom: 8px;
}

.process-step-label {
  display: block;
  margin-bottom: 4px;
  font-weight: 500;
  color: #666;
}

.process-step-code {
  display: block;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.04);
  border-radius: 4px;
  font-size: 12px;
  line-height: 20px;
  white-space: pre-wrap;
  word-break: break-all;
  color: #666;
}

.process-step-result {
  color: #666;

  :deep(.cherry-markdown) {
    color: #666;
    background: transparent;
    margin: 0;
    padding: 0;
    font-size: 14px;
    line-height: 22px;

    * {
      font-size: 14px;
      line-height: 22px;
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
