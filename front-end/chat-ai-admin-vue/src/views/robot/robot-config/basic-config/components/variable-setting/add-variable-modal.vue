<template>
  <div>
    <component
      :is="compact ? SettingsDialog : AntModal"
      :open="open"
      :title="t('title_field_detail')"
      @ok="handleOk"
      :width="compact ? 560 : 648"
      :wrapClassName="compact ? undefined : 'no-padding-modal'"
      :bodyStyle="compact ? undefined : { 'max-height': '70vh', 'overflow-y': 'auto', 'padding-right': '12px' }"
      :ok-text="compact ? commonT('confirm') : undefined"
      :cancel-text="compact ? commonT('cancel') : undefined"
      :confirm-loading="compact && saving"
      :closable="!compact || !saving"
      :mask-closable="!compact || !saving"
      :keyboard="!compact || !saving"
      :cancel-button-props="{ disabled: compact && saving }"
      @cancel="handleCancel"
    >
      <template v-if="compact" #closeIcon><img src="@/assets/img/customer-service/basic-settings/prompt/close.svg" alt="" /></template>
      <a-form
        :class="{ 'workspace-variable-form': compact }"
        :inert="compact && saving ? '' : null"
        :disabled="compact && saving"
        :model="form"
        ref="formRef"
        layout="vertical"
        :label-col="{ span: compact ? 24 : 4 }"
        :wrapper-col="{ span: 24 }"
      >
        <!-- 字段类型 -->
        <a-form-item :label="t('label_field_type')" name="variable_type">
          <a-select
            :aria-label="t('label_field_type')"
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
            :controls="!compact"
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
              :controls="!compact"
              style="width: 100%"
              stringMode
              v-model:value="form.default_value"
              :placeholder="t('ph_input_default_value')"
            />
          </div>
          <div v-if="form.variable_type == 'select_one'" :placeholder="t('ph_select_default_value')">
            <a-select v-model:value="form.default_value" style="width: 100%">
              <a-select-option v-for="item in filtersOptions" :key="item.label" :value="item.label">{{
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
        <a-form-item class="required-field" :label="null" name="must_input" v-if="form.variable_type != 'checkbox_switch'">
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
      <a-alert v-if="compact && saveFailed" type="error" :message="commonT('save_failed')" show-icon />
    </component>
  </div>
</template>

<script setup>
import { ref, reactive, computed, inject, nextTick } from 'vue'
import { createChatVariable } from '@/api/robot/index'
import { useRoute } from 'vue-router'
import draggable from 'vuedraggable'
import { message, Modal as AntModal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import SettingsDialog from '@/views/customer-service/basic-settings/components/settings-dialog.vue'
import { cloneConfig } from '@/views/customer-service/config-session'

const { t } = useI18n('views.robot.robot-config.basic-config.components.variable-setting.add-variable-modal')
const { t: commonT } = useI18n('views.customer-service.index')

const route = useRoute()
const customerContext = inject('customerRobotContext', null)
const props = defineProps({
  compact: { type: Boolean, default: false },
  robotKey: {
    type: String,
    default: ''
  }
})
const currentRobotKey = computed(() => props.robotKey || (props.compact ? customerContext?.robotInfo?.robot_key : '') || route.query.robot_key || '')
const open = ref(false)
const saving = ref(false)
const saveFailed = ref(false)

const emit = defineEmits(['ok'])

const form = reactive({
  robot_key: '',
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
  if (props.compact && saving.value) return
  form.robot_key = currentRobotKey.value
  form.id = data.id || ''
  form.variable_type = data.variable_type || 'input_string'
  form.variable_key = data.variable_key || ''
  form.variable_name = data.variable_name || ''
  form.max_input_length = data.max_input_length || 10
  form.default_value = data.default_value || ''
  form.must_input = data.must_input == '1'
  form.options = props.compact ? cloneConfig(data.options || [{ label: '' }]) : data.options || [{ label: '' }]
  saveFailed.value = false
  open.value = true
  if (props.compact) nextTick(() => formRef.value?.clearValidate())
}

const handleCancel = () => {
  if (!props.compact || !saving.value) open.value = false
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
const handleOk = async () => {
  if (props.compact && saving.value) return
  try {
    await formRef.value.validate()
  } catch { return }
  if (props.compact && (!open.value || saving.value)) return
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
  const params = {
    ...form,
    max_input_length,
    must_input: form.must_input ? 1 : 0,
    options: JSON.stringify(options)
  }
  if (props.compact) { saving.value = true; saveFailed.value = false }
  try {
    const res = await (customerContext ? customerContext.saveAction(() => createChatVariable(params)) : createChatVariable(params))
    if (props.compact && (res === false || (res?.res !== undefined && res.res !== 0))) { saveFailed.value = true; return }
    if (res === false) return
    message.success(form.id ? t('msg_update_success') : t('msg_create_success'))
    open.value = false
    emit('ok')
  } catch {
    if (props.compact) saveFailed.value = true
  } finally {
    saving.value = false
  }
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
.workspace-variable-form {
  padding: 8px 19px 8px 8px;
  color: #222;
  font-size: 14px;
  line-height: 22px;
  :deep(.ant-form-item) { margin-bottom: 24px; }
  :deep(.ant-form-item:last-child) { margin-bottom: 0; }
  :deep(.ant-form-item-label) { padding-bottom: 8px; }
  :deep(.ant-form-item-label > label) { height: 22px; color: #222; font-size: 14px; }
  :deep(.ant-form-item-required::before) { color: #ee1d36 !important; }
  :deep(.required-field .ant-form-item-control-input) { min-height: 24px; }
  :deep(.ant-input), :deep(.ant-input-number), :deep(.ant-select-selector) { min-height: 32px; border-radius: 6px; border-color: #d8d8d8; box-shadow: none; }
  :deep(.ant-input) { padding: 4px 12px; line-height: 22px; }
  :deep(.ant-input-number) { width: 100%; }
  :deep(.ant-input-number-input) { height: 30px; padding: 0 12px; }
  :deep(.ant-select-single .ant-select-selector) { height: 32px; padding: 0 12px; }
  :deep(.ant-checkbox-wrapper) { min-height: 24px; align-items: center; line-height: 22px; }
  :deep(.ant-checkbox) { top: 0; margin: 0 4px; }
  :deep(.ant-checkbox + span) { padding-inline-end: 0; }
  :deep(.ant-checkbox-inner) { border-radius: 4px; border-color: #d8d8d8; }
  :deep(.ant-checkbox-checked .ant-checkbox-inner) { border-color: #146ef5; background: #146ef5; }
}
</style>
