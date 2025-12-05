<style lang="less">
.vue-markdown {
  white-space: pre-wrap;
  width: 100%;

  ul, ol{
    list-style-type: circle;
    margin-block-start: 0px;
    margin-block-end: 0px;
  }
  ol {
    display: block;
    list-style-type: decimal;
    margin-block-start: 1em;
    margin-block-end: 1em;
    margin-inline-start: 0px;
    margin-inline-end: 0px;
    padding-inline-start: 40px;
    unicode-bidi: isolate;
  }
  img{
    width: auto;
    height: auto;
    max-width: 100%;
    max-height: 100%;
    display: block;
    margin-top: 8px;
  }
  p:last-child {
    margin-bottom: 0 !important;
  }
  div:last-child {
    margin-bottom: 0 !important;
  }
}
</style>

<template>
  <div ref="containerRef" class="vue-markdown cherry-markdown" v-html="html"></div>
</template>

<script setup>
// cherry-markdow 配置详解 https://github.com/Tencent/cherry-markdown/wiki/%E9%85%8D%E7%BD%AE%E9%A1%B9%E5%85%A8%E8%A7%A3
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core'
import { computed, onMounted, ref, watch, nextTick } from 'vue'
import textParseProcessing from '@/utils/textParseProcessing'
import { useChatStore } from '@/stores/modules/chat'
const chatStore = useChatStore()
const externalConfigPC = computed(()=> chatStore.externalConfigPC)

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

const containerRef = ref(null);

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
    }
  }
})

const html = computed(() => {
  let str = textParseProcessing(props.content)

  return md.makeHtml(str)
})


watch(html, () => {
  nextTick(() => {
    bindLinkEvents();
  });
}, { immediate: true });


function isMobileDevice() {
  const userAgent = navigator.userAgent || navigator.vendor || window.opera;
  const mobileRegex = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i;
  return mobileRegex.test(userAgent);
}
function bindLinkEvents() {
  const links = containerRef.value.querySelectorAll('a');
  links.forEach(link => {
    link.addEventListener('click', function(e) {
      e.preventDefault(); // 阻止默认跳转行为
      const href = this.getAttribute('href');
      if(isMobileDevice() || externalConfigPC.value.open_type == 1){
        // 移动设备上点击链接时，使用自定义的跳转逻辑
        window.open(href, '_blank')
        return
      }
      // 处理点击事件
      const width = +externalConfigPC.value.window_width || 1200;
      const height = +externalConfigPC.value.window_height || 650;
      const left = (screen.width - width) / 2;  // 居中
      const top = (screen.height - height) / 2; // 居中
      window.open(href, '_blank', `width=${width},height=${height},left=${left},top=${top},resizable=yes`)
    });
  });
}


</script>
