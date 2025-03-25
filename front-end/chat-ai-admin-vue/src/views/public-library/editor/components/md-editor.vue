<template>
  <div class="md-editor-wrapper" :class="['is-' + state.type]">
    <div class="md-editor-body">
      <div :class="[state.type + '-box']" id="md-editor"></div>
    </div>
  </div>
</template>

<script setup>
import 'cherry-markdown/dist/cherry-markdown.css'
import '@/assets/md-editor/theme/ant-design.less'
import Cherry from 'cherry-markdown'
import basicConfig from './base-config'
import { reactive, onMounted, onBeforeUnmount } from 'vue'

const emit = defineEmits(['inited', 'input', 'blur'])

let MDEditor = null

const state = reactive({
  type: 'editor',
  content: ''
})

const initPreview = () => {
  MDEditor = new Cherry({
    ...basicConfig,
    id: 'md-editor',
    value: state.content,
    editor: {
      defaultModel: 'previewOnly',
      keepDocumentScrollAfterInit: true
    },
    previewer: {
      enablePreviewerBubble: false,
      floatWhenClosePreviewer: false
    },
    toolbars: {
      toolbar: false,
      // 配置目录
      toc: {
        updateLocationHash: true, // 要不要更新URL的hash
        defaultModel: 'full', // pure: 精简模式/缩略模式，只有一排小点； full: 完整模式，会展示所有标题
        position: 'fixed', // 悬浮目录的悬浮方式。当滚动条在cherry内部时，用absolute；当滚动条在cherry外部时，用fixed
        cssText: 'right: 20px;'
      }
    }
  })

  MDEditor.setTheme('ant-design')
}

const initEditor = () => {
  let config = {
    ...basicConfig,
    id: 'md-editor',
    value: state.content,
    event: {
      afterInit() {
        emit('inited')
      },
      afterChange(text) {
        onContentInput(text)
      },
      blur: () => {
        onContentBlur()
      }
    }
  }

  MDEditor = new Cherry(config)

  MDEditor.setTheme('ant-design')
}

const setContent = (content, clearStack = false) => {
  MDEditor.setValue(content, clearStack)
}

const getContent = () => {
  return MDEditor.getMarkdown()
}

const setDoc = (doc) => {
  setContent(doc.content)
}

const getDoc = () => {
  return {
    content: getContent()
  }
}

const clearDoc = () => {
  setContent('', true)
}

const onContentInput = (val) => {
  emit('input', val)
}

const onContentBlur = (val) => {
  emit('blur', val)
}

const init = (type, content) => {
  state.type = type
  state.content = content

  if (MDEditor) {
    MDEditor.destroy()
  }

  if (type === 'editor') {
    initEditor()
  } else if (type === 'preview') {
    initPreview()
  }
}

onMounted(() => {})

onBeforeUnmount(() => {
  if (MDEditor) {
    MDEditor.destroy()
  }
})

defineExpose({
  init,
  setContent,
  getContent,
  setDoc,
  getDoc,
  clearDoc
})
</script>

<style lang="less" scoped>
.md-editor-wrapper {
  position: relative;
  height: 100%;
  overflow: hidden;

  .md-editor-body {
    height: 100%;
  }
  .preview-box {
    max-width: 900px;
    margin: 0 auto;
  }

  &.is-editor {
    .md-editor-body {
      width: 100%;
      ::v-deep(.cherry-previewer) {
        border-left: 1px solid #f0f0f0 !important;
      }
    }
  }
  ::v-deep(.cherry) {
    box-shadow: none !important;
  }
  &.is-preview {
    .md-editor-body {
      overflow: hidden;
      overflow-y: auto;

      ::v-deep(.cherry-previewer) {
        padding: 0 !important;
        border-left: none !important;
      }
    }
  }
}
</style>
