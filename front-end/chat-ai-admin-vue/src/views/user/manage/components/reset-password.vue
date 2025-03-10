<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="476px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="登录密码" v-bind="validateInfos.password">
            <a-input-password
              v-model:value="formState.password"
              placeholder="密码必须包含字母、数字或者字符中的两种，6-32位"
            />
          </a-form-item>
          <a-form-item label="确认密码" v-bind="validateInfos.check_password">
            <a-input-password
              v-model:value="formState.check_password"
              placeholder="请重新输入密码"
            />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { validatePassword } from '@/utils/validate.js'
import { ref, reactive, toRaw } from 'vue'
import { Form, message } from 'ant-design-vue'
import { resetPass } from '@/api/manage/index.js'
const emit = defineEmits(['ok'])

const useForm = Form.useForm
const show = ref(false)
const modalTitle = ref('重置登录密码')
const formState = reactive({
  password: '',
  check_password: '',
  id: ''
})

const formRules = reactive({
  password: [
    {
      message: '请输入登录密码',
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!validatePassword(value) && value) {
          return Promise.reject('密码必须包含字母、数字或者字符中的两种，6-32位')
        }
        return Promise.resolve()
      }
    }
  ],
  check_password: [
    {
      message: '请输入确认密码',
      required: true
    },
    {
      validator: async (rule, value) => {
        if (value != formState.password && value) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const open = (record) => {
  console.log(record.id,'==')
  show.value = true
  resetFields()
  formState.password = ''
  formState.check_password = ''
  formState.id = record.id
  
}

const handleOk = () => {
  validate().then(() => {
    resetPass({
      ...toRaw(formState)
    }).then((res) => {
      message.success(`修改成功`)
      show.value = false
      emit('ok')
    })
  })
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}
</style>
