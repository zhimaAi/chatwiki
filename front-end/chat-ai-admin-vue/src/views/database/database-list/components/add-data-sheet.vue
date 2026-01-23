<template>
  <div class="add-data-sheet">
    <a-modal v-model:open="open" :title="t(modalTitle)" @ok="handleOk" :width="540">
      <a-form layout="vertical">
        <a-form-item :label="t('data_table_name_label')" v-bind="validateInfos.name">
          <a-input :maxLength="64" v-model:value="formState.name" :placeholder="t('data_table_name_placeholder')" />
          <div class="form-tip">{{ t('data_table_name_tip') }}</div>
        </a-form-item>
        <a-form-item :label="t('data_table_description_label')" v-bind="validateInfos.description">
          <a-textarea
            :maxLength="500"
            style="height: 80px"
            v-model:value="formState.description"
            :placeholder="t('data_table_description_placeholder')"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { Form, message } from 'ant-design-vue'
import { addForm, editForm } from '@/api/database'
import { useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.database.database-list.components.add-data-sheet')
const router = useRouter()
const emit = defineEmits(['ok'])
const open = ref(false)
const modalTitle = ref('create_data_table')
const formState = reactive({
  name: '',
  description: '',
  id: ''
})
const show = (data = {}) => {
  open.value = true
  resetFields()
  formState.name = data.name || ''
  formState.description = data.description || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? 'edit_data_table' : 'create_data_table'
}
const formRules = reactive({
  name: [
    {
      required: true,
      validator: async (rule, value) => {
        if (!/^[a-z][a-z0-9_]*$/.test(value)) {
          return Promise.reject(t('name_validation_error'))
        }
        return Promise.resolve()
      }
    }
  ],
  description: [
    {
      required: true,
      message: t('description_required_error')
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then((res) => {
    if (formState.id) {
      editForm(formState).then(() => {
        postCallback(t('edit_success'))
      })
    } else {
      addForm(formState).then((res) => {
        postCallback(t('create_success'))
        router.push({
          path: '/database/details/field-manage',
          query: {
            form_id: res.data.id,
          }
        })
      })
    }
  })
}
function postCallback(tip) {
  message.success(tip)
  open.value = false
  emit('ok')
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
