<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="520px" @cancel="handleClose">
      <a-alert class="alert-box" message="您还未修改初始密码，为了保障您的数据安全，请您尽快重置密码。" type="info" show-icon />
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
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()
const useForm = Form.useForm
const show = ref(true)
const modalTitle = ref('重置登录密码')
const formState = reactive({
  password: '',
  check_password: '',
  id: userStore.userInfo.user_id,
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


const handleOk = () => {
  validate().then(() => {
    resetPass({
      ...toRaw(formState)
    }).then((res) => {
      message.success(`修改成功`)
      show.value = false
      userStore.reset(true)
    })
  })
}

const handleClose = () => {
  userStore.setResetPassModal()
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}
.alert-box{
  margin-top: 24px;
}
</style>
