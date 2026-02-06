<template>
  <div>
    <a-modal v-model:open="open" :title="t('title_rename')" @ok="handleOk">
      <div class="rename-input">
        <a-input v-model:value="file_name" :placeholder="t('ph_input_document_name')" />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { editLibFile } from '@/api/library/index.js'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.components.rename-modal')
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
    message.error(t('msg_input_document_name'))
  }
  editLibFile({
    id,
    file_name: file_name.value
  }).then((res) => {
    message.success(t('msg_rename_success'))
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
