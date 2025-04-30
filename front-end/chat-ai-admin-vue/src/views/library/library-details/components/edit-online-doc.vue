<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :maskClosable="false"
      title="添加在线数据"
      width="746px"
      @ok="handleSaveUrl"
    >
      <a-form
        class="url-add-form"
        layout="vertical"
        ref="urlFormRef"
        :model="formState"
        :rules="rules"
      >
        <a-form-item name="file_name" label="文件名称">
          <a-input v-model:value="formState.file_name" placeholder="请输入" />
        </a-form-item>
        <a-form-item name="doc_auto_renew_frequency" label="更新频率">
          <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">不自动更新</a-select-option>
            <a-select-option :value="2">每天</a-select-option>
            <a-select-option :value="3">每3天</a-select-option>
            <a-select-option :value="4">每7天</a-select-option>
            <a-select-option :value="5">每30天</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          v-if="formState.doc_auto_renew_frequency > 1"
          name="doc_auto_renew_minute"
          label="更新时间"
        >
          <a-time-picker
            valueFormat="HH:mm"
            v-model:value="formState.doc_auto_renew_minute"
            format="HH:mm"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { editLibFile } from '@/api/library'
import { convertTime } from '@/utils/index'
const emit = defineEmits(['ok'])
const formState = reactive({
  id: '',
  file_name: '',
  doc_auto_renew_frequency: 1,
  doc_auto_renew_minute: ''
})
const rules = ref({
  file_name: [
    {
      message: '请文档名称',
      required: true
    }
  ],
  doc_auto_renew_frequency: [
    {
      message: '请选择更新频率',
      required: true
    }
  ]
})

const open = ref(false)
const show = (data) => {
  data = JSON.parse(JSON.stringify(data))
  formState.id = data.id || ''
  formState.file_name = data.file_name || ''
  formState.doc_auto_renew_frequency = +data.doc_auto_renew_frequency || 1
  if (data.doc_auto_renew_minute > 0) {
    formState.doc_auto_renew_minute = convertTime(data.doc_auto_renew_minute)
  }
  open.value = true
}
const confirmLoading = ref(false)

const urlFormRef = ref(null)
const handleSaveUrl = () => {
  urlFormRef.value
    .validate()
    .then(() => {
      confirmLoading.value = true
      editLibFile({
        ...formState,
        doc_auto_renew_minute: formState.doc_auto_renew_minute
          ? convertTime(formState.doc_auto_renew_minute)
          : 0
      }).then(() => {
        open.value = false
        confirmLoading.value = false
        emit('ok')
      })
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped></style>
