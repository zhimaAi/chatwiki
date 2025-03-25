<style lang="less" scoped>
.doc-title-wrapper {
  padding: 0;
  display: flex;
  align-items: center;

  .file-icon {
    margin-right: 4px;
  }

  .doc-title-input {
    display: block;
    font-size: 16px;
    width: 100%;
    height: 24px;
    line-height: 24px;
    padding: 0;
    font-weight: 600;
    border: none;
    outline: none;
    resize: none;
    overflow-y: hidden;

    color: #262626;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
      'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
      'Noto Color Emoji';

    &::placeholder {
      color: #b0b1b8;
    }
  }
}
</style>

<template>
  <div class="doc-title-wrapper">
    <svg-icon class="file-icon" name="doc-file2" style="font-size: 16px; color: #262626"></svg-icon>
    <input
      class="doc-title-input"
      ref="autoHeightInput"
      v-model="text"
      resize="none"
      placeholder="无标题"
      type="text"
      :maxlength="100"
      @input="onInput"
      @blur="endEdit"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'

const emit = defineEmits(['update:title', 'input', 'blur'])

const props = defineProps({
  title: {
    type: String,
    default: ''
  }
})
// const defaultHeight = 24
const text = ref(props.title)

// const textareaHeight = ref(24)
const autoHeightInput = ref(null)

const endEdit = () => {
  emit('blur', text.value)
}

const adjustHeight = () => {
  // textareaHeight.value = defaultHeight
  // nextTick(() => {
  //   const newHeight = autoHeightInput.value.scrollHeight
  //   textareaHeight.value = newHeight > defaultHeight ? newHeight : defaultHeight
  // })
}

const onInput = () => {
  adjustHeight()

  emit('input', text.value)
  emit('update:title', text.value)
}

watch(
  () => props.title,
  () => {
    text.value = props.title
  }
)

onMounted(() => {
  adjustHeight()
})
</script>
