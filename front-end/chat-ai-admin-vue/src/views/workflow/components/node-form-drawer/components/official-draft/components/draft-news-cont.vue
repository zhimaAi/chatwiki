<template>
  <a-modal
    :title="t('title_content')"
    width="880px"
    v-model:open="visible"
    @ok="save"
  >
    <a-alert type="info" class="zm-alert-info">
      <template #message>
        {{ t('msg_rich_text_content_intro') }}
      </template>
    </a-alert>

    <div class="editor-box">
      <Toolbar
        style="border-bottom: 1px solid #ccc;"
        :editor="editorRef"
        :defaultConfig="toolbarConfig"
        :mode="mode"
      />
      <Editor
        style="height: 480px; overflow-y: hidden;"
        v-model="valueHtml"
        :defaultConfig="editorConfig"
        :mode="mode"
        @onCreated="handleCreated"
        @onChange="handleChange"
      />
    </div>

    <div v-if="showVariable"
         :style="varDialogStyle"
         class="dropdown-list"
         ref="dropdownList"
    >
      <CascadePanel
        :options="variableOptions"
        @change="onSelectOption"
        @direction-change="handleDirectionChange">
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </CascadePanel>
    </div>
  </a-modal>
</template>

<script setup>
import {ref, shallowRef, onBeforeUnmount, toRaw, watch} from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import {Editor, Toolbar} from '@wangeditor/editor-for-vue'
import {Editor as SlateEditor, Range, Transforms, Text} from 'slate'
import CascadePanel from "@/views/workflow/components/at-input/cascade-panel.vue";
import {runPlugin} from "@/api/plugins/index.js";
import {getBase64, getTreeOptions} from "@/utils/index.js";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.official-draft.components.draft-news-cont')


const emit = defineEmits(['change'])
const props = defineProps({
  appData: {
    type: Object,
    default: () => ({})
  },
  variableOptions: {
    type: Array,
  },
  defaultSelectedList: {
    type: Array,
    default: () => [],
  },
})

const visible = ref(false)
const editorRef = shallowRef()
const activeEditor = shallowRef(null)
const mode = ref('default')
const valueHtml = ref('')
const toolbarConfig = ref({
  excludeKeys: [
    "insertTable",
    "codeBlock",
    "uploadVideo",
    "group-video",
  ]
})
const editorConfig = {
  placeholder: t('ph_input_content'),
  MENU_CONF: {
    uploadImage: {
      async customUpload(file, insertFn) {
        const base64 = await getBase64(file)
        const res = await runPlugin({
          name: 'official_draft',
          action: "default/exec",
          params: JSON.stringify({
            business: 'upload_image',
            arguments: {
              ...props.appData,
              media: base64
            }
          })
        })
        insertFn(res?.data?.url)
      },
    }
  }
}
const showVariable = ref(false)
const varDialogStyle = ref({})
const showOptions = ref([])
const tagAllData = ref([])
const selectedList = ref([])
const selectedIdSet = ref(new Set())

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

function show(value = '', tags = []) {
  initData(value, tags)
  visible.value = true
}

function initData(value, tags) {
  let treeOptions = getTreeOptions(props.variableOptions);
  selectedList.value = tags
  tagAllData.value = [...tags, ...treeOptions].filter(item => item && item.value && item.value != '' && item.typ !== 'node');
  // 记录已经替换过的的value
  const replacedValues = new Set();
  for (let opt of  tagAllData.value) {
    value = decodeVariable(value, opt, replacedValues)
  }
  valueHtml.value = value
}

function decodeVariable(html, opt, replacedValues) {
  const regex = new RegExp(`(${opt.value})`, 'g');
  if (!replacedValues.has(opt.value) && regex.test(html)) {
    console.log('open')
    let text = opt.node_name + '.' + opt.text
    text = text.replace(/\./g, '/')
    html = html.replace(regex, `<a href="zm-variable:${opt.value}" target="_blank"><strong>${text}</strong></a>`)
    replacedValues.add(opt.value)
  }
  return html
}

function handleCreated(editor) {
  editorRef.value = editor
}

function handleChange(editor) {
  const {selection} = editor
  if (!selection || !Range.isCollapsed(selection)) {
    showVariable.value = false
    return
  }

  const path = selection.focus.path
  const offset = selection.focus.offset
  const text = SlateEditor.string(editor, path)

  // 光标前一个字符
  const charBefore = text[offset - 1]
  if (charBefore === '/') {
    activeEditor.value = editor
    showVariable.value = true
    updateDialogPosition(editor)
  } else {
    showVariable.value = false
  }
}


// 获取当前光标定位
function getCaretRectSafe() {
  const selection = window.getSelection()
  if (!selection || selection.rangeCount === 0) return null

  const range = selection.getRangeAt(0)
  if (!range.collapsed) return null

  let rect = range.getClientRects()[0]
  if (rect) return rect

  // ====== 兜底方案 ======
  const span = document.createElement('span')
  span.textContent = '\u200b' // 零宽字符
  span.style.display = 'inline-block'

  range.insertNode(span)
  rect = span.getBoundingClientRect()
  span.remove()
  return rect
}

function updateDialogPosition(editor) {
  const rect = getCaretRectSafe()
  if (!rect) return
  const left = rect.left + window.scrollX
  const top = rect.bottom + window.scrollY + 6
  varDialogStyle.value = `left:${left}px; top:${top}px;`;
}

function onSelectOption(value, selectedValuePath, selectedPath) {
  let item = selectedPath[selectedPath.length - 1];
  delete item['children'];
  selectOption(item)
}

function selectOption(opt, isInit = false) {
  if (!opt) return;
  // 处理用户选择逻辑
  showVariable.value = false;
  let text = opt.node_name + '.' + opt.text;
  text = text.replace(/\./g, '/')
  const dataSet = {
    id: opt.id,
    label: opt.label,
    value: opt.value,
    text: opt.text,
    index: selectedList.value.length,
  };
  selectedList.value.push(opt);
  const set = new Set(selectedIdSet.value)
  set.add(opt.value)
  selectedIdSet.value = set

  if (!isInit) {
    insertAtCaret(text, dataSet);
  }
}

function insertAtCaret(text, dataSet) {
  const editor = activeEditor.value
  const {selection} = editor
  if (!selection) return

  // 如果有选中内容，先删
  if (!Range.isCollapsed(selection)) SlateEditor.deleteFragment(editor)

  // 删除触发字符 `/`
  SlateEditor.deleteBackward(editor, 'character')

  SlateEditor.insertNode(editor, {
    type: 'link',
    url: `zm-variable:${dataSet.value}`,
    children: [
      {text, bold: true}
    ]
  })
  // 光标落点
  SlateEditor.insertText(editor, ' ')
}

function initShowOptionList() {
  showOptions.value = props.variableOptions.filter((opt) => {
    if (this.atText) {
      return opt.label.startsWith(this.atText);
    }
    if (this.canRepeat) return true;
    return !this.selectedIdSet.has(opt.id + "");
  });
}

function handleDirectionChange(direction, width) {
  const leftMatch = varDialogStyle.value.match(/left:\s*(-?\d+\.?\d*)px/);
  const topMatch = varDialogStyle.value.match(/top:\s*(-?\d+\.?\d*)px/);

  if (!leftMatch || !topMatch) return;

  const currentLeft = parseFloat(leftMatch[1]);
  const currentTop = parseFloat(topMatch[1]);
  let newLeft;

  if (direction === 'left') {
    newLeft = currentLeft - width;
  } else {
    newLeft = currentLeft + width;
  }

  varDialogStyle.value = `left:${newLeft}px; top:${currentTop}px;`;
}

function unwrapRoot(html) {
  if (typeof html !== 'string') return ''
  const match = html.match(/^<div[^>]*class=["'][^"']*\bzm-editor-root\b[^"']*["'][^>]*>([\s\S]*)<\/div>$/)
  return match ? match[1] : html
}

function save() {
  let html = toRaw(valueHtml.value)
  let selectedListMap = {};
  selectedList.value.forEach(item => {
    if(!selectedListMap[item.value] && html.indexOf(item.value) > -1){
      selectedListMap[item.value] = item;
    }
  });
  let mySelectedList = Object.values(selectedListMap);
  html = unwrapRoot(html)
  html = html.replace(/<a[^>]*href="zm-variable:([^"]+)"[^>]*>[\s\S]*?<\/a>/g, '$1')
  emit('change', html, toRaw(mySelectedList))
  visible.value = false
}

defineExpose({
  show,
})
</script>

<style scoped lang="less">
.editor-box {
  border: 1px solid #ccc;
  border-radius: 4px;
  margin-top: 16px;

  &.w-e-full-screen-container {
    z-index: 9999;
  }
}

.dropdown-list {
  position: fixed;
  padding: 0;
  margin: 0;
  background-color: #fff;
  z-index: 999999;
  min-width: 180px;
  overflow: hidden;
  user-select: none;
  border-radius: 8px;
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);

  .dropdown-list-nodata {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    font-size: 13px;
    padding: 8px;
  }

  .dropdown-list-item {
    padding: 8px;
    cursor: pointer;
    background: #fff;

    &:hover {
      background-color: rgba(0, 0, 0, 0.04) !important;
    }

    &.active {
      background-color: rgba(0, 0, 0, 0.09) !important;
    }
  }


  &::-webkit-scrollbar {
    width: 0px;
    height: 0px;
  }

  /*定义滚动条轨道 内阴影+圆角*/

  &::-webkit-scrollbar-track {
    border-radius: 0px;
    background-color: #fafafa;
  }

  /*定义滑块 内阴影+圆角*/

  &::-webkit-scrollbar-thumb {
    border-radius: 0px;
    background: rgb(191, 191, 191);
  }
}
</style>
