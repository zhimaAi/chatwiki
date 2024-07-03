<style lang="less">
.vue-markdown {
  line-height: 1.2;
  code {
    width: 100%;
    overflow: auto;
    white-space: pre-wrap;
  }
}
</style>

<template>
  <div>
    <div class="vue-markdown" v-html="html"></div>
  </div>
</template>

<script setup>
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/docco.css'
import { computed } from 'vue'

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

const md = new MarkdownIt({
  html: true, // 在源码中启用 HTML 标签
  breaks: true, // 转换段落里的 '\n' 到 <br>。
  linkify: true,
  typographer: true,
  xhtmlOut: true,
  langPrefix: 'language-', // 给围栏代码块的 CSS 语言前缀。对于额外的高亮代码非常有用。
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return '<pre class="hljs"><code>' + hljs.highlight(lang, str, true).value + '</code></pre>'
      } catch (__) {}
    }

    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>'
  }
})

const html = computed(() => {
  return md.render(props.content)
})
</script>
