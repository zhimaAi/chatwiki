<template>
  <div
    class="final-answer-card"
    :class="{ 'is-indented': shouldIndent }"
    v-if="hasFinalAnswerContent(item)"
  >
    <div class="final-answer-header">
      <div class="final-answer-title">
        <svg-icon name="agent2"></svg-icon>
        <span>{{ t('label_final_answer') }}</span>
      </div>
      <div class="final-answer-source">
        {{ t('label_generated_by', { name: getRobotDisplayName(item) }) }}
      </div>
    </div>
    <div class="message-content final-answer-content" v-viewer>
      <CherryMarkdown
        class="markdown-body final-answer-markdown"
        :content="item.content"
      ></CherryMarkdown>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import '@/assets/style/markdown/markdown.less'

const { t } = useI18n('views.clawbot.chat.components.messages.final-answer-card')
const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  },
  robotInfo: {
    type: Object,
    default: () => ({})
  },
  shouldIndent: {
    type: Boolean,
    default: false
  }
})

const getRobotDisplayName = (item) => {
  return item?.name || props.robotInfo?.robot_name || t('label_agent_fallback')
}

const hasFinalAnswerContent = (item) => {
  return String(item?.content || '').trim() !== ''
}
</script>

<style lang="less" scoped>
.final-answer-card {
  position: relative;
  overflow: hidden;
  padding: 16px;
  border: 1px solid #bcd5ff;
  border-radius: 14px;
  background: linear-gradient(168.82deg, rgba(239, 246, 255, 0.5) 0%, #ffffff 100%);
  box-shadow:
    0 1px 3px 0 rgba(0, 0, 0, 0.1),
    0 1px 2px -1px rgba(0, 0, 0, 0.1);

  &::before {
    position: absolute;
    left: -1px;
    top: -1px;
    bottom: -1px;
    width: 4px;
    content: '';
    background: linear-gradient(180deg, #2b7fff 0%, #615fff 100%);
  }

  &.is-indented {
    margin-left: 32px;
  }
}

.final-answer-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.final-answer-title {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #1447e6;
  font-size: 14px;
  line-height: 22px;

  .svg-icon {
    font-size: 14px;
    color: #1447e6;
  }
}

.final-answer-source {
  margin-left: auto;
  color: #8c8c8c;
  font-size: 12px;
  line-height: 20px;
  text-align: right;
}

.message-content {
  line-height: 22px;
  font-size: 14px;
  font-weight: 400;
  color: #3a4559;
  white-space: pre-wrap;
}

.final-answer-content {
  white-space: normal;

  :deep(.final-answer-markdown.markdown-body) {
    max-width: none;
    margin: 0;
    padding: 0;
    font-size: 14px;
    line-height: 22px;
    background: transparent;
    white-space: normal;
  }
}
</style>
