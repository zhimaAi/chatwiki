<style scoped>
.auto-size-textarea {
  display: block;
  padding: 0;
  height: 22px;
  max-height: 66px;
  line-height: 22px;
  resize: none;
  overflow-y: auto; /* 允许内容滚动 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* Internet Explorer and Edge */

  &::-webkit-scrollbar {
    display: none; /* Chrome, Safari, and Opera */
  }
}
</style>

<template>
  <textarea
    class="auto-size-textarea"
    rows="1"
    ref="textareaRef"
    :value="props.value"
    @input="updateValue"
    @focus="onFocus"
    @blur="onBlur"
    @keydown.enter.prevent.exact="onEnter"
    @keydown.shift.enter.prevent="onShiftEnter"
    placeholder="在此输入，Shift+Enter换行"
  ></textarea>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:value', 'change', 'focus', 'blur', 'enter', 'shiftEnter'])

const inputLineHeight = 22
const inputMaxRow = 3
const inputMaxHeight = inputLineHeight * inputMaxRow

const textareaRef = ref<HTMLTextAreaElement | null>(null)

function onFocus() {
  emit('focus')
}

function onBlur() {
  emit('blur')
}

function onEnter(event: KeyboardEvent) {
  if (!event.shiftKey) {
    emit('enter')
  }
}

function isCursorAtEnd(inputElement) {  
    var val = inputElement.value;  
    var pos = inputElement.selectionStart;  
    return pos === val.length; 
}

function insertNewLine() {
  if (textareaRef.value) {
    const cursorPosition = textareaRef.value.selectionStart
    const newValue =
    textareaRef.value.value.slice(0, cursorPosition) + '\n' + textareaRef.value.value.slice(cursorPosition)
    textareaRef.value.value = newValue

    if(isCursorAtEnd(textareaRef.value)) {
      textareaRef.value.scrollTop = textareaRef.value.scrollHeight
    }
    
    // 将光标移至新行开头
    textareaRef.value.setSelectionRange(cursorPosition + 1, cursorPosition + 1)

    updateHeight()
  }
}

function onShiftEnter(event: KeyboardEvent) {
  event.preventDefault()
  insertNewLine()
  emit('shiftEnter')
}

// 监听 textarea 的输入事件，并更新其高度和 value
function updateValue(event: Event) {
  const target = event.target as HTMLTextAreaElement

  emit('update:value', target.value)
  emit('change', target.value)
  updateHeight()
}

// 更新 textarea 的高度以匹配内容
function updateHeight() {
  if (textareaRef.value) {
    textareaRef.value.style.height = `auto`

    if (textareaRef.value.scrollHeight > inputMaxHeight) {
      textareaRef.value.style.height = `${inputMaxHeight}px`
    } else if(textareaRef.value.value.length == 0) {
      textareaRef.value.style.height = `${inputLineHeight}px`
    } else {
      textareaRef.value.style.height = `${textareaRef.value.scrollHeight}px`
    }
  }
}

watch(
  () => props.value,
  () => {
    if (!props.value && textareaRef.value) {
      textareaRef.value.style.height = `${inputLineHeight}px`
    }
  }
)

// 组件挂载后，确保 textarea 初始高度正确
onMounted(() => {
  if (textareaRef.value) {
    updateHeight()
  }
})
</script>
