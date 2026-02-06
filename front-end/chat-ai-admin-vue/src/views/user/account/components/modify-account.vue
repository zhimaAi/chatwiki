<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="476px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item :label="t('label_old_password')" v-bind="validateInfos.old_password">
            <a-input-password
              v-model:value="formState.old_password"
              :placeholder="t('ph_input_old_password')"
            />
          </a-form-item>
          <a-form-item :label="t('label_password')" v-bind="validateInfos.password">
            <a-input-password
              v-model:value="formState.password"
              :placeholder="t('ph_password_rule')"
            />
          </a-form-item>
          <a-form-item :label="t('label_confirm_password')" v-bind="validateInfos.check_password">
            <a-input-password
              v-model:value="formState.check_password"
              :placeholder="t('ph_reinput_password')"
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
import { saveProfile } from '@/api/manage/index.js'
import { useI18n } from '@/hooks/web/useI18n'
const emit = defineEmits(['ok'])

const useForm = Form.useForm
const { t } = useI18n('views.user.account.components.modify-account')
const show = ref(false)
const modalTitle = ref(t('title_modify_password'))
const formState = reactive({
  password: '',
  check_password: '',
  old_password: '',
  id: ''
})

const formRules = reactive({
  old_password: [
    {
      message: t('msg_input_old_password'),
      required: true
    }
  ],
  password: [
    {
      message: t('msg_input_password'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!validatePassword(value) && value) {
          return Promise.reject(t('ph_password_rule'))
        }
        if(value == formState.old_password){
          return Promise.reject(t('msg_password_same'))
        }
        return Promise.resolve()
      }
    }
  ],
  check_password: [
    {
      message: t('msg_input_confirm_password'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (value != formState.password && value) {
          return Promise.reject(t('msg_password_mismatch'))
        }
        return Promise.resolve()
      }
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const open = (record) => {
  show.value = true
  resetFields()
  formState.password = ''
  formState.check_password = ''
  formState.old_password = ''
  formState.id = record.user_id;
}

const handleOk = () => {
  validate().then(() => {
    saveProfile({
      ...toRaw(formState)
    }).then((res) => {
      message.success(t('msg_modify_success'))
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
