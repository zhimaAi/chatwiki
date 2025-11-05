<template>
  <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="540">
    <a-form layout="vertical" style="margin-top: 16px">
      <a-form-item label="token限额" v-bind="validateInfos.token">
        <a-input-number
          v-model:value="formState.token"
          style="width: 60%"
          :min="0"
          placeholder="请输入"
        >
          <template #addonAfter> k </template>
        </a-input-number>
        <div class="form-tip">剩余额度：xxx</div>
      </a-form-item>
      <a-form-item label="备注" v-bind="validateInfos.description">
        <a-textarea
          :maxLength="500"
          style="height: 80px"
          v-model:value="formState.description"
          placeholder="请输入"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { Form, message } from 'ant-design-vue'
const emit = defineEmits(['ok'])
const open = ref(false)
const modalTitle = ref('token限额')
const formState = reactive({
  token: '',
  description: ''
})
const show = (data = {}) => {
  open.value = true
  resetFields()
  formState.name = data.name || ''
  formState.description = data.description || ''
  formState.id = data.id || ''
}
const formRules = reactive({
  token: [
    {
      required: true,
      message: '请输入token限额'
    }
  ],
  description: [
    {
      required: true,
      message: '请输入数据表描述'
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then((res) => {})
}
defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-tip {
  color: #8c8c8c;
  margin-top: 8px;
}
</style>
