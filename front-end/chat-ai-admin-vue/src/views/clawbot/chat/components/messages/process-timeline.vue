<template>
  <div class="process-timeline" v-if="visibleSteps.length">
    <div class="process-step" v-for="(step, stepIndex) in visibleSteps" :key="step.id">
      <div class="process-step-side">
        <div class="process-step-icon">
          <svg-icon class="step-icon" v-if="step.type === 'thinking'" name="search-book"></svg-icon>
          <svg-icon class="step-icon" v-else-if="step.type === 'tool'" name="search-book"></svg-icon>
          <svg-icon class="step-icon" v-else name="search-book"></svg-icon>
        </div>
        <div class="process-step-line" v-if="stepIndex < visibleSteps.length - 1"></div>
      </div>
      <div class="process-step-card">
        <div class="process-step-header" @click="toggleProcessStep(step)">
          <div class="process-step-title" :class="{ expanded: step.expanded }">
            <span class="title-text">{{ step.title }}</span>
            <svg-icon name="arrow-down" class="arrow-down"></svg-icon>
          </div>
          <div
            class="process-step-status"
            :class="{ 'is-done': step.status === 'done' }"
            v-if="step.type === 'tool' && step.status"
          >
            <LoadingOutlined v-if="step.status === 'running'" class="loading" />
            <span v-else class="status-dot"></span>
            <span>{{ getProcessStepStatusText(step) }}</span>
          </div>
        </div>
        <div class="process-step-body" v-if="step.expanded">
          <template v-if="step.type === 'thinking'">
            <div class="process-step-running" v-if="step.status === 'running' && !step.contentText">
              <LoadingOutlined class="loading" />
              <span>{{ step.resultText || t('msg_thinking') }}</span>
            </div>
            <div class="process-step-text" v-else-if="!step.contentText">
              {{ step.resultText }}
            </div>
            <div class="process-step-text" v-else>
              <CherryMarkdown :content="step.contentText"></CherryMarkdown>
            </div>
          </template>
          <template v-else-if="step.type === 'tool'">
            <div class="process-step-block" v-if="step.paramsText">
              <div class="block-label">{{ t('label_params') }}</div>
              <pre class="process-step-code">{{ step.paramsText }}</pre>
            </div>
            <div class="process-step-block" v-if="step.resultText">
              <div class="block-label">{{ t('label_output') }}</div>
              <pre class="process-step-code">{{ step.resultText }}</pre>
            </div>
          </template>
          <template v-else>
            <div class="process-step-block" v-if="step.paramsText">
              <div class="block-label">{{ t('label_params') }}</div>
              <pre class="process-step-code">{{ step.paramsText }}</pre>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { LoadingOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'

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

const toggleProcessStep = (step) => {
  step.expanded = !step.expanded
}

const getProcessStepStatusText = (step) => {
  if (step.status === 'running') {
    return t('status_running')
  }
  if (step.status === 'done') {
    return t('status_done')
  }
  return ''
}
</script>

<style lang="less" scoped>
.process-timeline {
  margin-bottom: 12px;
}

.process-step {
  display: flex;
  gap: 8px;

  &:not(:last-child) {
    margin-bottom: 8px;
  }
}

.process-step-side {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  width: 24px;
  padding-top: 8px;
}

.process-step-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #7a8699;

  .step-icon {
    font-size: 20px;
    color: #7a8699;
  }
}

.process-step-line {
  width: 1px;
  flex: 1;
  margin-top: 6px;
  background: #d9d9d9;
}

.process-step-card {
  flex: 1;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #fcfdfe;
}

.process-step-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 14px;
  cursor: pointer;
}

.process-step-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  color: #262626;
  font-size: 14px;
  line-height: 22px;

  .title-text {
    min-width: 0;
    word-break: break-word;
  }

  .arrow-down {
    flex-shrink: 0;
    font-size: 14px;
    color: #8c8c8c;
    transition: transform 0.2s;
  }

  &.expanded {
    .arrow-down {
      transform: rotate(180deg);
    }
  }
}

.process-step-status {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  font-size: 12px;
  line-height: 20px;
  color: #8c8c8c;

  .status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #00bc7d;
  }

  &.is-done {
    color: #009966;
  }
}

.process-step-body {
  padding: 4px 14px 12px;
}

.process-step-text {
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
  white-space: pre-wrap;
  word-break: break-word;
}

.process-step-running {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;

  .loading {
    font-size: 14px;
    color: #8c8c8c;
  }
}

.process-step-block {
  &:not(:last-child) {
    margin-bottom: 12px;
  }

  .block-label {
    margin-bottom: 4px;
    color: #8c8c8c;
    font-size: 12px;
    line-height: 20px;
  }
}

.process-step-code {
  margin: 0;
  padding: 9px 13px;
  overflow-x: auto;
  border: 1px solid #f1f5f9;
  border-radius: 4px;
  background: #f2f4f7;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
}
</style>
