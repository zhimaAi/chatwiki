<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :title="modalTitle"
      width="476px"
      @ok="handleOk"
    >
      <a-form style="margin-top: 24px" layout="vertical" ref="formRef" :model="formState">
        <a-form-item
          name="group_name"
          label="分组名称"
          :rules="[{ required: true, message: '请输入分组名称' }]"
        >
          <a-input
            v-model:value="formState.group_name"
            placeholder="请输入，最多10个字"
            :maxLength="10"
          ></a-input>
        </a-form-item>
        <a-form-item name="group_desc" label="分组描述">
          <a-textarea
            style="height: 80px"
            :maxLength="100"
            v-model:value="formState.group_desc"
            placeholder="请输入"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { savePromptLibraryGroup } from '@/api/user/index.js'
import { message } from 'ant-design-vue'
const emit = defineEmits(['ok'])
const formState = reactive({
  id: '',
  group_name: '',
  group_desc: ''
})
const confirmLoading = ref(false)
const modalTitle = ref('新建分组')
const open = ref(false)
const formRef = ref(null)

const show = (data = {}) => {
  formRef.value && formRef.value.resetFields()
  formState.group_name = data.group_name || ''
  formState.group_desc = data.group_desc || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? '编辑分组' : '新建分组'
  open.value = true
}

const handleOk = () => {
  formRef.value.validate().then(() => {
    savePromptLibraryGroup({
      ...formState
    }).then(() => {
      emit('ok')
      message.success(`${modalTitle.value}成功`)
      open.value = false
    })
  })
}
defineExpose({
  show
})
</script>

<style lang="less" scoped></style>
