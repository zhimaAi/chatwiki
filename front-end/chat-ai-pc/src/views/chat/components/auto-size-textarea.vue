<style scoped>
.text-input {
  line-height: 22px;
  height: 22px;
  width: 100%;
  padding: 0;
  font-size: 14px;
  font-weight: 400;
  color: rgb(26, 26, 26);
  overflow: hidden;
  white-space: pre-wrap; /* 保持内容的换行，并允许自动换行 */
  resize: none;
  border: none;
  background: none;
  transition: height 0.1s ease-in-out;

  &::placeholder {
    font-size: 14px;
    font-weight: 400;
    color: rgb(191, 191, 191);
  }
}

/* 滚动条样式 */
.text-input::-webkit-scrollbar {
    width: 4px; /*  设置纵轴（y轴）轴滚动条 */
    height: 4px; /*  设置横轴（x轴）轴滚动条 */
}
/* 滚动条滑块（里面小方块） */
.text-input::-webkit-scrollbar-thumb {
    cursor: pointer;
    border-radius: 5px;
    background: transparent;
}
/* 滚动条轨道 */
.text-input::-webkit-scrollbar-track {
    border-radius: 5px;
    background: transparent;
}

/* hover时显色 */
.text-input:hover::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.2);
}
.text-input:hover::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.1);
}
</style>

<template>
  <textarea
    :style="{height: inputHeight}"
    class="text-input"
    ref="textareaRef"
    :focus="focus"
    :value="props.value"
    @input="updateValue"
    @focus="onFocus"
    @blur="onBlur"
    @keydown.enter.prevent.exact="onEnter"
    @keydown.shift.enter.prevent="onTextEnter"
    @keydown.ctrl.enter.prevent="onTextEnter"
    @keydown.alt.enter.prevent="onTextEnter"
    :placeholder="t('ph_input')"
  ></textarea>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import calcTextareaHeight from '@/utils/calcTextareaHeight'
import { useI18n } from '@/hooks/web/useI18n'

// 初始化 i18n
const { t } = useI18n('views.chat.components.auto-size-textarea')

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  focus: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:value', 'change', 'focus', 'blur', 'enter', 'shiftEnter'])

const inputHeight = ref("22px")

const textareaRef = ref<HTMLTextAreaElement | null>(null)

function onFocus(event) {
  emit('focus', event)
}

function onBlur(event) {
  emit('blur', event)
}

// function isCursorAtEnd(inputElement) {  
//     var val = inputElement.value;  
//     var pos = inputElement.selectionStart;  
//     return pos === val.length; 
// }

function insertNewLine() {
  if (textareaRef.value) {
    const cursorPosition = textareaRef.value.selectionStart
    const newValue =
    textareaRef.value.value.slice(0, cursorPosition) + '\n' + textareaRef.value.value.slice(cursorPosition)
    textareaRef.value.value = newValue
    
    // 将光标移至新行开头
    textareaRef.value.setSelectionRange(cursorPosition + 1, cursorPosition + 1)

    textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px';  // 调整高度
  }
}

function onEnter(event: KeyboardEvent) {
  // 上面阻止了按键，进这里的都是直接按enter
  emit('enter')
}

function onTextEnter(event: KeyboardEvent) {
  event.preventDefault()
  insertNewLine()
  // emit('shiftEnter')
}

// 监听 textarea 的输入事件，并更新其高度和 value
function updateValue(event: Event) {
  const target = event.target as HTMLTextAreaElement

  emit('update:value', target.value)
  emit('change', target.value)
  
  inputHeight.value = calcTextareaHeight(textareaRef.value).height // 调整高度

  if (target) {
    if (parseInt(inputHeight.value) >= 160) {
      target.style.overflow = 'auto'
    } else {
      target.style.overflow = 'hidden'
    }
  }
}

watch(
  () => props.value,
  () => {
    if (!props.value && textareaRef.value) {
      // 消息清空后输入框回到最初的高度
      inputHeight.value = '22px'
      // 回车后输入框失去焦点
      // textareaRef.value.blur()
    }
  }
)

// 组件挂载后，确保 textarea 初始高度正确
onMounted(() => {
 
})
</script>
