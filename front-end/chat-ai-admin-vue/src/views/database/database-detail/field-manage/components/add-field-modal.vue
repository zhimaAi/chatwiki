<template>
  <div class="add-data-sheet">
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="540">
      <a-form layout="vertical">
        <a-form-item v-bind="validateInfos.name">
          <template #label
            >{{ t('label_name') }}&nbsp;
            <a-tooltip>
              <template #title>{{ t('label_name_tooltip') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-input :maxLength="64" v-model:value="formState.name" :placeholder="t('ph_name')" />
        </a-form-item>
        <a-form-item v-bind="validateInfos.description">
          <template #label
            >{{ t('label_description') }}&nbsp;
            <a-tooltip>
              <template #title>{{ t('label_description_tooltip') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-input :maxLength="64" v-model:value="formState.description" :placeholder="t('ph_description')" />
        </a-form-item>
        <a-form-item v-bind="validateInfos.type">
          <template #label
            >{{ t('label_type') }}&nbsp;
            <a-tooltip>
              <template #title>{{ t('label_type_tooltip') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-select v-model:value="formState.type" style="width: 100%" :placeholder="t('ph_type')">
            <a-select-option v-for="item in typeOption" :value="item.value" :key="item.value">{{
              item.label
            }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <template #label
            >{{ t('label_required') }}&nbsp;
            <a-tooltip>
              <template #title>
                <div>{{ t('label_required_tooltip_required') }}</div>
                <div>{{ t('label_required_tooltip_optional') }}</div>
              </template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>
          <a-switch
            v-model:checked="formState.required"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            checkedValue="true"
            unCheckedValue="false"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { PlusOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Form, message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { addFormField } from '@/api/database'

const { t } = useI18n('views.database.database-detail.field-manage.components.add-field-modal')

const modalTitle = ref(t('modal_title_add'))
const rotue = useRoute()
const query = rotue.query
const open = ref(false)
const emit = defineEmits(['ok'])
const typeOption = [
  {
    label: 'string',
    value: 'string'
  },
  {
    label: 'integer',
    value: 'integer'
  },
  {
    label: 'number',
    value: 'number'
  },
  {
    label: 'boolean',
    value: 'boolean'
  }
]
const formState = reactive({
  name: '',
  description: '',
  type: void 0,
  required: 'false',
  form_id: query.form_id,
  id: ''
})
const show = (data) => {
  open.value = true
  resetFields()
  formState.name = data.name || ''
  formState.description = data.description || ''
  formState.type = data.type
  formState.required = data.required
  formState.id = data.id || ''
  modalTitle.value = data.id ? t('modal_title_edit') : t('modal_title_add')
}
const formRules = reactive({
  name: [
    {
      required: true,
      validator: async (rule, value) => {
        if (!/^[a-z][a-z0-9_]*$/.test(value)) {
          return Promise.reject(t('validator_name_pattern'))
        }
        return Promise.resolve()
      }
    }
  ],
  description: [
    {
      required: true,
      message: t('validator_description_required')
    }
  ],
  type: [
    {
      required: true,
      message: t('validator_type_required')
    }
  ]
})
const useForm = Form.useForm
const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const handleOk = () => {
  validate().then(() => {
    addFormField(formState).then((res) => {
      let tip = formState.id ? t('msg_edit_success') : t('msg_add_success')
      message.success(tip)
      open.value = false
      emit('ok')
    })
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
