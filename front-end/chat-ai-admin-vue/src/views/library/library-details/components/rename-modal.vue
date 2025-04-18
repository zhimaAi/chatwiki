<template>
  <div>
    <a-modal v-model:open="open" title="重命名" @ok="handleOk">
      <div class="rename-input">
        <a-input v-model:value="file_name" placeholder="请输入文档名称" />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { editLibFile } from '@/api/library/index.js'
import { message } from 'ant-design-vue'
const emit = defineEmits(['ok'])
const open = ref(false)
const file_name = ref('')
let id = ''
const show = (data = {}) => {
  data = JSON.parse(JSON.stringify(data))
  file_name.value = data.file_name
  id = data.id
  open.value = true
}
const handleOk = () => {
  if (!file_name.value) {
    message.error('请输入文档名称')
  }
  editLibFile({
    id,
    file_name: file_name.value
  }).then((res) => {
    message.success('重命名成功')
    emit('ok', file_name.value)
    open.value = false
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.rename-input {
  margin: 32px 0;
}
</style>
