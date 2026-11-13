<template>
  <div class="reply-text" v-html="html"></div>
</template>

<script>
import MarkdownIt from 'markdown-it'

const markdown = new MarkdownIt({ html: false, breaks: true }).disable('image')
</script>

<script setup>
import { computed } from 'vue'

const props = defineProps({ content: { type: String, required: true } })
const html = computed(() => markdown.render(props.content))
</script>

<style scoped lang="less">
.reply-text {
  line-height: 1.7;
  overflow-wrap: anywhere;

  :deep(> :first-child) { margin-top: 0; }
  :deep(> :last-child) { margin-bottom: 0; }
  :deep(pre) {
    max-width: 100%;
    overflow: auto;
    padding: 8px 12px;
    border-radius: 6px;
    background: #fff;
  }
}
</style>
