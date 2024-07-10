<style lang="less">
.vue-markdown {
  white-space: initial;
  width: 100%;
  img{
    width: auto;
    height: auto;
    max-width: 100%;
    max-height: 100%;
    display: block;
    margin-top: 8px;
  }
}
</style>

<template>
  <div class="vue-markdown cherry-markdown" v-html="html"></div>
</template>

<script setup>
// cherry-markdow 配置详解 https://github.com/Tencent/cherry-markdown/wiki/%E9%85%8D%E7%BD%AE%E9%A1%B9%E5%85%A8%E8%A7%A3
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core'
import { computed } from 'vue'

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

const md = new CherryEngine({
  engine: {
    global: {
      classicBr: true,
      flowSessionContext: true
    },
    syntax: {
      codeBlock: {
        theme: 'dark', // 默认为深色主题
        wrap: true, // 超出长度是否换行，false则显示滚动条
        lineNumber: false, // 默认显示行号
        copyCode: true, // 是否显示“复制”按钮
        editCode: false, // 是否显示“编辑”按钮
        changeLang: true // 是否显示“切换语言”按钮
      },
      autoLink: {
        enableShortLink: false
      },
      header: {
        anchorStyle: 'none'
      }
    },
  }
})

const html = computed(() => {
  return md.makeHtml(props.content)
})
</script>
