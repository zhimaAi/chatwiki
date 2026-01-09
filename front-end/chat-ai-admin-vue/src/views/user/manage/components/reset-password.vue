<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="476px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item :label="t('login_password')" v-bind="validateInfos.password">
            <a-input-password
              v-model:value="formState.password"
              :placeholder="t('password_placeholder')"
            />
          </a-form-item>
          <a-form-item :label="t('confirm_password')" v-bind="validateInfos.check_password">
            <a-input-password
              v-model:value="formState.check_password"
              :placeholder="t('re_enter_password')"
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
import { useI18n } from '@/hooks/web/useI18n'
import { computed } from 'vue'

const emit = defineEmits(['ok'])

const useForm = Form.useForm
const { t } = useI18n('views.user.manage.components.reset-password')
const show = ref(false)
const modalTitle = computed(() => t('modal_title'))
const formState = reactive({
  password: '',
  check_password: '',
  id: ''
})

const formRules = reactive({
  password: [
    {
      message: t('enter_login_password'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!validatePassword(value) && value) {
          return Promise.reject(t('password_must_contain'))
        }
        return Promise.resolve()
      }
    }
  ],
  check_password: [
    {
      message: t('enter_confirm_password'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (value != formState.password && value) {
          return Promise.reject(t('passwords_not_match'))
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
      message.success(t('modification_successful'))
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
