<style lang="less">
.vue-markdown {
  white-space: normal;
  word-break: break-all;
  width: 100%;
  
  img {
    width: auto;
    height: auto;
    max-width: 100%;
    max-height: 100%;
    display: block;
    margin-top: 8px;
    min-height: 100px;
  }
  video {
    width: 100%;
    height: 240px;
    max-width: 480px;
    max-height: 240px;
    display: block;
    margin-top: 8px;
    border-radius: 8px;
    background: #000;
    object-fit: contain;
  }
  p:last-child {
    margin-bottom: 0;
  }
}
</style>

<template>
  <div
    ref="markdownRef"
    class="vue-markdown cherry-markdown"
    v-html="html"
    @click="handleImagePreview"
  ></div>
</template>

<script setup>
// cherry-markdow 配置详解 https://github.com/Tencent/cherry-markdown/wiki/%E9%85%8D%E7%BD%AE%E9%A1%B9%E5%85%A8%E8%A7%A3
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core'
import { computed, ref } from 'vue'
import { api as viewerApi } from 'v-viewer'
import textParseProcessing from '@/utils/textParseProcessing'

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  enableImagePreview: {
    type: Boolean,
    default: false
  }
})

const markdownRef = ref(null)

const md = new CherryEngine({
  themeSettings: {
    // 目前应用的主题
    mainTheme: 'orange',
  },
  engine: {
    global: {
      classicBr: true,
      flowSessionContext: true,
      htmlWhiteList: 'video',
      htmlAttrWhiteList: 'src|controls|preload|playsinline|webkit-playsinline|poster|class|muted|loop|autoplay'
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
      link: {
        target: '_blank',
      },
      autoLink: {
        target: '_blank',
        enableShortLink: false
      },
      header: {
        anchorStyle: 'none'
      },
      mathBlock: {
        plugins: true, 
        engine: 'MathJax', // katex或MathJax
        src: './libs/MathJax/es5/tex-svg.js',
        // src: 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js',
      },
      inlineMath: {
        engine: 'MathJax', // katex或MathJax
      },
    }
  }
})

const html = computed(() => {
  let str = textParseProcessing(props.content)
  return md.makeHtml(str)
})

const handleImagePreview = (event) => {
  if (!props.enableImagePreview) {
    return
  }

  const target = event.target instanceof Element ? event.target : null
  const currentImg = target?.closest('img')
  if (!currentImg || !markdownRef.value?.contains(currentImg)) {
    return
  }

  const imageEntries = Array.from(markdownRef.value.querySelectorAll('img'))
    .map((img) => ({
      el: img,
      src: img.currentSrc || img.getAttribute('src') || ''
    }))
    .filter((item) => item.src)

  const initialViewIndex = imageEntries.findIndex((item) => item.el === currentImg)
  if (initialViewIndex < 0) {
    return
  }

  event.preventDefault()

  const viewer = viewerApi({
    options: {
      title: false,
      navbar: false,
      keyboard: false,
      loop: false,
      initialViewIndex
    },
    index: initialViewIndex,
    images: imageEntries.map((item) => item.src)
  })

  viewer.show()
}
</script>
