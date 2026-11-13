<template>
  <div
    class="final-answer-card"
    :class="{ 'has-process': shouldIndent }"
    v-if="hasFinalAnswerContent(item)"
  >
    <div class="message-content chat-markdown-content">
      <CherryMarkdown
        class="markdown-body"
        :content="item.content"
        enable-image-preview
      ></CherryMarkdown>
    </div>
    <PromptOptimizeButton v-if="!item.loading && !item.startLoading" :robot="robotInfo" scene="conversation" :messageId="item.id" />
  </div>
</template>

<script setup>
import PromptOptimizeButton from '@/components/prompt-optimize/button.vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import '@/assets/style/markdown/markdown.less'

defineProps({
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

const hasFinalAnswerContent = (item) => {
  return String(item?.content || '').trim() !== ''
}
</script>

<style lang="less" scoped>
.final-answer-card {
  &.has-process {
    margin-top: 24px;
  }
}

.message-content {
  color: #262626;
}
</style>
