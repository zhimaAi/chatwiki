<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="t('title_field_detail')"
      @ok="handleOk"
      :width="648"
      wrapClassName="no-padding-modal"
      :bodyStyle="{ 'max-height': '70vh', 'overflow-y': 'auto', 'padding-right': '12px' }"
    >
      <a-form
        :model="form"
        ref="formRef"
        layout="vertical"
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 24 }"
      >
        <!-- 字段类型 -->
        <a-form-item :label="t('label_field_type')" name="variable_type">
          <a-select
            @change="handleTypChange"
            v-model:value="form.variable_type"
            :placeholder="t('ph_select_field_type')"
          >
            <a-select-option value="input_string">{{ t('option_text') }}</a-select-option>
            <a-select-option value="input_number">{{ t('option_number') }}</a-select-option>
            <a-select-option value="select_one">{{ t('option_select_one') }}</a-select-option>
            <a-select-option value="checkbox_switch">{{ t('option_checkbox_switch') }}</a-select-option>
          </a-select>
        </a-form-item>

        <!-- 字段key -->
        <a-form-item
          :label="t('label_field_key')"
          name="variable_key"
          :rules="[
            {
              required: true,
              validator: (rule, value) => checkKey(rule, value)
            }
          ]"
        >
          <a-input v-model:value="form.variable_key" :maxLength="10" :placeholder="t('ph_input_field_key')" />
        </a-form-item>

        <!-- 字段名 -->
        <a-form-item :label="t('label_field_name')" name="variable_name" :rules="[{ required: true }]">
          <a-input v-model:value="form.variable_name" :maxLength="10" :placeholder="t('ph_input_field_name')" />
        </a-form-item>

        <!-- 最大长度 -->
        <a-form-item
          :label="t('label_max_length')"
          name="max_input_length"
          v-if="['input_string', 'input_number'].includes(form.variable_type)"
        >
          <a-input-number
            v-model:value="form.max_input_length"
            :placeholder="t('ph_input')"
            :precision="0"
            :min="1"
            :max="50"
          />
        </a-form-item>

        <!-- 默认值 -->
        <a-form-item :label="t('label_default_value')" name="default_value">
          <div v-if="form.variable_type == 'input_string'">
            <a-input v-model:value="form.default_value" :placeholder="t('ph_input_default_value')" />
          </div>
          <div v-if="form.variable_type == 'input_number'">
            <a-input-number
              :maxLength="50"
              style="width: 100%"
              stringMode
              v-model:value="form.default_value"
              :placeholder="t('ph_input_default_value')"
            />
          </div>
          <div v-if="form.variable_type == 'select_one'" :placeholder="t('ph_select_default_value')">
            <a-select v-model:value="form.default_value" style="width: 100%">
              <a-select-option v-for="item in filtersOptions" :value="item.label">{{
                item.label
              }}</a-select-option>
            </a-select>
          </div>
          <div v-if="form.variable_type == 'checkbox_switch'" :placeholder="t('ph_select_default_value')">
            <a-select v-model:value="form.default_value" style="width: 100%">
              <a-select-option value="1">{{ t('option_selected') }}</a-select-option>
              <a-select-option value="2">{{ t('option_not_selected') }}</a-select-option>
            </a-select>
          </div>
        </a-form-item>

        <!-- 是否必填 -->
        <a-form-item :label="null" name="must_input" v-if="form.variable_type != 'checkbox_switch'">
          <a-checkbox v-model:checked="form.must_input">{{ t('label_required') }}</a-checkbox>
        </a-form-item>

        <!-- 下拉选项（仅当类型为 select_one 时显示） -->
        <a-form-item v-if="form.variable_type === 'select_one'" :label="t('label_options')" name="options">
          <a-button type="dashed" @click="addOption">{{ t('btn_add_option') }}</a-button>
          <a-form-item-rest>
            <draggable
              v-model="form.options"
              item-key="index"
              group="table-rows"
              handle=".drag-btn"
            >
              <template #item="{ element, index }">
                <div :key="index" class="option-item">
                  <span class="drag-btn"><svg-icon name="drag" /></span>
                  <a-input
                    v-model:value="element.label"
                    :maxLength="50"
                    :placeholder="t('ph_input_option_label')"
                  />
                  <a-button type="link" danger @click="removeOption(index)">{{ t('btn_delete') }}</a-button>
                </div>
              </template>
            </draggable>
          </a-form-item-rest>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { createChatVariable } from '@/api/robot/index'
import { useRoute } from 'vue-router'
import draggable from 'vuedraggable'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.variable-setting.add-variable-modal')

const query = useRoute().query
const open = ref(false)

const emit = defineEmits(['ok'])

const form = reactive({
  robot_key: query.robot_key,
  id: '',
  variable_type: 'input_string',
  variable_key: '',
  variable_name: '',
  max_input_length: 10,
  default_value: '',
  must_input: false,
  options: [{ label: '' }] // 初始化一个空选项
})

const filtersOptions = computed(() => {
  return form.options.filter((item) => item.label)
})

const show = (data = {}) => {
  form.id = data.id || ''
  form.variable_type = data.variable_type || 'input_string'
  form.variable_key = data.variable_key || ''
  form.variable_name = data.variable_name || ''
  form.max_input_length = data.max_input_length || 10
  form.default_value = data.default_value || ''
  form.must_input = data.must_input == '1'
  form.options = data.options || [{ label: '' }]
  open.value = true
}

const checkKey = (rule, value) => {
  if (!value) {
    return Promise.reject(t('msg_input_key_required'))
  }
  // 校验是否只包含英文字母和下划线
  const regex = /^[a-zA-Z_]+$/
  if (!regex.test(value)) {
    return Promise.reject(t('msg_key_invalid'))
  }
  return Promise.resolve()
}

const formRef = ref()
const handleOk = () => {
  formRef.value.validate().then(() => {
    // 验证成功，可以提交表单
    let options = []
    if (form.variable_type == 'select_one') {
      options = form.options.filter((item) => item.label != '')
      if (options.length == 0) {
        return message.error(t('msg_fill_options'))
      }
    }
    let max_input_length = form.max_input_length
    if (form.variable_type == 'select_one' || form.variable_type == 'checkbox_switch') {
      max_input_length = 1
    }
    createChatVariable({
      ...form,
      max_input_length,
      must_input: form.must_input ? 1 : 0,
      options: JSON.stringify(options)
    }).then((res) => {
      message.success(form.id ? t('msg_update_success') : t('msg_create_success'))
      open.value = false
      emit('ok')
    })
  })
}

const addOption = () => {
  form.options.push({ label: '' })
}

const removeOption = (index) => {
  if (form.options.length > 1) {
    form.options.splice(index, 1)
  }
}

const handleTypChange = () => {
  form.options = [{ label: '' }]
  if (form.variable_type == 'checkbox_switch') {
    form.default_value = '2'
  } else {
    form.default_value = void 0
  }
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.option-item {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  .drag-btn {
    cursor: grab;
  }
}
</style>
