<template>
  <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="540">
    <a-form layout="vertical" style="margin-top: 16px">
      <a-form-item :label="t('token_limit')" v-bind="validateInfos.max_token">
        <a-input-number
          v-model:value="formState.max_token"
          style="width: 60%"
          :min="0.001"
          :precision="3"
          :placeholder="t('please_input')"
        >
          <template #addonAfter> k </template>
        </a-input-number>
      </a-form-item>
      <a-form-item :label="t('remark')" v-bind="validateInfos.description">
        <a-textarea
          :maxLength="500"
          style="height: 80px"
          v-model:value="formState.description"
          :placeholder="t('please_input')"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { Form, message } from 'ant-design-vue'
import { tokenLimitCreate } from '@/api/manage/index.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.usetoken.modal')

const emit = defineEmits(['ok'])
const open = ref(false)
const modalTitle = ref(t('token_limit'))
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
      message: t('please_input_token_limit'),
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
        message.success(t('modify_success'))
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
