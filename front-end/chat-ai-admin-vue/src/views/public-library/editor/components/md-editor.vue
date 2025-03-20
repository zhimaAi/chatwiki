<template>
  <div class="md-editor-wrapper" :class="['is-' + state.type]">
    <div class="md-editor-body">
      <div :class="[state.type + '-box']" id="vditor"></div>
    </div>
  </div>
</template>

<script setup>
import { uploadFile } from '@/api/app'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { reactive, onMounted, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'

const emit = defineEmits(['inited', 'input', 'blur'])

// 允许上传的图片类型
const allowUploadImageType = ['image/png', 'image/jpeg', 'image/jpg', 'image/gif']

let vditor = null

const uploadImage = (files) => {
  /* 上传逻辑 */
  const file = files[0]

  uploadFile({
    file: file,
    category: 'library_doc_image'
  })
    .then((res) => {
      if (res.res != 0) {
        return message.error(res.msg)
      }

      let url = res.data.link
      if (file.type.indexOf('image') > -1) {
        let text = `![${file.name}](${url})`

        vditor.insertValue(text)
      }
    })
    .catch((err) => {
      message.error('图片上传失败')
      console.log(err)
    })
}

const getMinHeight = () => {
  return document.documentElement.clientHeight - 52 - 300
}

const toolbar = [
  'emoji',
  'headings',
  'bold',
  'italic',
  'strike',
  'link',
  '|',
  'list',
  'ordered-list',
  'check',
  'outdent',
  'indent',
  '|',
  'quote',
  'line',
  'code',
  'inline-code',
  'insert-before',
  'insert-after',
  '|',
  'upload',
  'table',
  '|',
  'undo',
  'redo',
  '|',
  'fullscreen',
  'edit-mode',
  {
    name: 'more',
    toolbar: ['both', 'code-theme', 'content-theme', 'export', 'outline', 'preview']
  }
]

const state = reactive({
  type: 'editor',
  content: ''
})

const initPreview = () => {
  Vditor.preview(document.getElementById('vditor'), state.content, {
    // 对选中后的内容进行阅读
    speech: {
      enable: false
    },
    // 为标题添加锚点 0：不渲染；1：渲染于标题前；2：渲染于标题后，默认 0
    anchor: 0,
    theme: {
      current: 'ant-design'
    }
  })
}

const initEditor = () => {
  vditor = new Vditor('vditor', {
    // width: 800,
    // height: opt.minHeight,
    minHeight: getMinHeight(),
    mode: 'wysiwyg',
    theme: 'light',
    value: state.content,
    toolbar: toolbar,
    toolbarConfig: {
      hide: false
    },
    cache: {
      enable: false
    },
    outline: {
      enable: false,
      position: 'right'
    },
    preview: {
      maxWidth: 900,
      markdown: {
        toc: true,
        codeBlockPreview: true, // wysiwyg 和 ir 模式下是否对代码块进行渲染
        mathBlockPreview: true // wysiwyg 和 ir 模式下是否对数学公式进行渲染
      },
      theme: {
        current: 'ant-design'
      }
    },
    upload: {
      multiple: false,
      accept: allowUploadImageType.join(','),
      handler: (files) => {
        uploadImage(files)
      }
    },
    input(value) {
      onContentInput(value)
    },
    blur(value) {
      onContentBlur(value)
    },
    after() {
      emit('inited', vditor)

      // const demo = document.querySelector('#vditor .vditor-toolbar')
      // const newContainer = document.getElementById('mdEditorToolbarWarrper')

      // // 直接移动元素（事件保留）
      // newContainer.appendChild(demo)
    }
  })
}

const setContent = (content, clearStack = false) => {
  vditor.setValue(content, clearStack)
}

const getContent = () => {
  return vditor.getValue()
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

  if (vditor) {
    vditor.destroy()
  }

  if (type === 'editor') {
    initEditor()
  } else if (type === 'preview') {
    initPreview()
  }
}

onMounted(() => {})

onBeforeUnmount(() => {
  if (vditor) {
    vditor.destroy()
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
  &.is-editor {
    padding-top: 37px;
  }
  &.is-preview {
    .md-editor-body {
      overflow: hidden;
      overflow-y: auto;
    }
  }
  .md-editor-body {
    height: 100%;
  }
  .preview-box {
    max-width: 900px;
    margin: 0 auto;
  }

  .vditor {
    width: 100%;
    border: none;
    &::v-deep(.vditor-toolbar) {
      // display: none;
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      border-bottom: 0px solid #f0f0f0 !important;
      background-color: #fff;
      z-index: 2;

      .vditor-toolbar__item {
        border-top: 1px solid #f0f0f0;
        border-bottom: 1px solid #f0f0f0;
        padding: 0 5px;
        .vditor-tooltipped {
          padding: 10px 5px;
        }
        .vditor-tooltipped::before {
          display: none;
        }
        &:first-child {
          border-left: 1px solid #f0f0f0;
          border-radius: 6px 0 0 6px;
        }
        &:last-child {
          border-right: 1px solid #f0f0f0;
          border-radius: 0 6px 6px 0;
        }
      }

      .vditor-toolbar__divider {
        height: 37px;
        border-top: 1px solid #f0f0f0;
        border-bottom: 1px solid #f0f0f0;
        border-left: 0;
        padding: 8px 5px;
        margin: 0;
        &::after {
          display: block;
          content: '';
          height: 100%;
          border-left: 1px solid #f0f0f0;
        }
      }
    }
    &::v-deep(.vditor-content) {
      .vditor-reset {
        padding-top: 26px !important;
        background-color: #fff !important;
        & > h1:first-child,
        & > h2:first-child,
        & > h3:first-child,
        & > h4:first-child,
        & > h5:first-child,
        & > h6:first-child {
          margin-top: 0;
        }
      }
    }
  }
}
</style>
