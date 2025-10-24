<template>
  <div>
    <Codemirror
      @wheel.stop=""
      originalStyle
      class="cm-component"
      v-model:value="code"
      :options="cmOptions"
      ref="cmRef"
      :height="height"
      :width="width"
      @change="onChange"
      @input="onInput"
      @ready="onReady"
    >
    </Codemirror>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import 'codemirror/mode/javascript/javascript.js'
import Codemirror from 'codemirror-editor-vue3'
import 'codemirror/theme/dracula.css'
const emit = defineEmits(['update:value', 'change'])

const code = ref('')

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  width: {
    type: Number,
    default: 512
  },
  height: {
    type: Number,
    default: 150
  }
})

watch(
  () => props.value,
  (val) => {
    code.value = val
  },
  {
    immediate: true
  }
)

const cmOptions = {
  mode: 'text/javascript',
  theme: 'dracula'
}

const cmRef = ref(null)

const onChange = (val) => {
  emit('update:value', val)
}

const handleRefresh = ()=>{
  cmRef.value.refresh()
}

const onReady = () => {}
const onInput = () => {}

defineExpose({
  handleRefresh
})
</script>

<style lang="less" scoped></style>