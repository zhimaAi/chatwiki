<template>
  <div>
    <a-modal
      v-model:open="visible"
      :confirm-loading="loading"
      :maskClosable="false"
      title="设置更新频率"
      @ok="handleOk"
    >
      <a-form style="margin: 32px 0" layout="vertical" :model="formState">
        <a-form-item name="doc_auto_renew_frequency" label="更新频率" required>
          <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">不自动更新</a-select-option>
            <a-select-option :value="2">每天</a-select-option>
            <a-select-option :value="3">每3天</a-select-option>
            <a-select-option :value="4">每7天</a-select-option>
            <a-select-option :value="5">每30天</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { editLibFile } from '@/api/library'
import { message } from 'ant-design-vue'
const emit = defineEmits(['ok'])
const visible = ref(false)
const loading = ref(false)
const formState = reactive({
  doc_auto_renew_frequency: 1,
  id: ''
})
const open = (data) => {
  let { doc_auto_renew_frequency, id } = data
  formState.doc_auto_renew_frequency = +doc_auto_renew_frequency
  formState.id = id
  visible.value = true
}
const handleOk = () => {
  loading.value = true
  editLibFile({
    ...formState
  })
    .then((res) => {
      emit('ok', formState.doc_auto_renew_frequency)
      visible.value = false
      message.success('设置成功')
    })
    .finally(() => {
      loading.value = false
    })
}
defineExpose({ open })
</script>

<style lang="less" scoped>
</style>