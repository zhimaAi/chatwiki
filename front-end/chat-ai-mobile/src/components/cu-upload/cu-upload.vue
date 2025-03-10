<style lang="less" scoped>
.cu-upload {
  .file-input {
    display: none;
  }
  .upload-btn {
    display: inline-block;
    cursor: pointer;
  }
}
</style>

<template>
  <div class="cu-upload">
    <input
      ref="fileInput"
      class="file-input"
      type="file"
      :name="name"
      :accept="props.accept"
      @change="onChange"
    />
    <span class="upload-btn" @click="handleUpload"><slot></slot></span>
  </div>
</template>
<script setup>
import { ref } from 'vue'

const emit = defineEmits(['change'])
const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  accept: {
    type: String,
    default: '*'
  }
})

const fileInput = ref(null)

const onChange = (e) => {
  emit('change', e.target.files)
}

const handleUpload = () => {
  fileInput.value.value = ''
  fileInput.value.click()
}
</script>
