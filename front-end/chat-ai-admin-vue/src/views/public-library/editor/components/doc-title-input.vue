<style lang="less" scoped>
.doc-title-input-wrapper {
  width: 100%;
  .doc-title-input {
    display: block;
    width: 100%;
    height: 32px;
    line-height: 32px;
    padding: 0;
    font-size: 24px;
    font-weight: 600;
    outline: none;
    resize: none;
    box-shadow: none;
    border: 0;
    background: none;
    overflow-y: hidden;
    color: #262626;

    &::placeholder {
      color: #b0b1b8;
    }
    &:focus {
      border: none;
      outline: none;
      box-shadow: none;
    }
  }
}
</style>

<template>
  <div class="doc-title-input-wrapper">
    <a-input class="doc-title-input" v-model:value.trim="text" type="text" placeholder="" :maxlength="100" @input="onInput"
      @blur="endEdit" ref="inputRef" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'

const emit = defineEmits(['update:title', 'input', 'blur', 'titleClick'])

const inputRef = ref(null)

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  isHighlight: {
    type: Boolean,
    default: false
  }
})

const text = ref(props.title)

const endEdit = () => {
  emit('blur', text.value)
}

const onInput = () => {
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

})
</script>
