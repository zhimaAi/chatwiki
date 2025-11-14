<template>
  <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="540">
    <a-form layout="vertical" style="margin-top: 16px">
      <a-form-item label="token限额" v-bind="validateInfos.max_token">
        <a-input-number
          v-model:value="formState.max_token"
          style="width: 60%"
          :min="0.001"
          :precision="3"
          placeholder="请输入"
        >
          <template #addonAfter> k </template>
        </a-input-number>
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
import { tokenLimitCreate } from '@/api/manage/index.js'
const emit = defineEmits(['ok'])
const open = ref(false)
const modalTitle = ref('token限额')
const formState = reactive({
  max_token: '',
  description: '',
  robot_id: '',
  token_app_type: ''
})

function formatNum(num) {
  if (num <= 0 || !num) {
    return ''
  }
  return (num / 1000).toFixed(3)
}

let limit_use = ref(0)
const show = (data = {}) => {
  open.value = true
  resetFields()
  formState.max_token = formatNum(data.max_token)
  formState.description = data.description || ''
  formState.robot_id = data.robot_id || ''
  formState.token_app_type = data.token_app_type || ''

  limit_use.value = formatNum(data.max_token - data.use_token)
}
const formRules = reactive({
  max_token: [
    {
      required: true,
      message: '请输入token限额',
      trigger: 'blur'
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then((res) => {
    tokenLimitCreate({ ...formState, max_token: (formState.max_token * 1000).toFixed(3) }).then(
      (res) => {
        message.success('修改成功')
        open.value = false
        emit('ok')
      }
    )
  })
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
