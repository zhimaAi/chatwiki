<template>
  <a-modal
    wrapClassName="variable-modal-class"
    v-model:open="open"
    :title="null"
    :footer="null"
    :closable="isEdit"
    :maskClosable="false"
    :width="560"
  >
    <div class="modal-header">
      <div class="banner-bg">
        <img src="@/assets/img/form-modal-header.png" alt="" />
      </div>
      <div class="title">{{ t('title_fill_info') }}</div>
    </div>
    <div class="form-body">
      <a-form :model="formState" ref="formRef" layout="vertical">
        <a-form-item
          v-for="item in waitVariables"
          :key="item.variable_key"
          :label="item.variable_type == 'checkbox_switch' ? null : item.variable_name"
          :name="item.variable_key"
          :rules="[{ required: item.must_input == 1, message: t('msg_input_required', { name: item.variable_name }) }]"
        >
          <template v-if="item.variable_type == 'input_string'">
            <a-input
              v-model:value="formState[item.variable_key]"
              :maxLength="getLength(item)"
              :placeholder="t('ph_input')"
            />
          </template>
          <template v-if="item.variable_type == 'input_number'">
            <a-input-number
              style="width: 100%"
              v-model:value="formState[item.variable_key]"
              stringMode
              :maxLength="50"
              :placeholder="t('ph_input')"
            />
          </template>
          <template v-if="item.variable_type == 'select_one'">
            <a-select v-model:value="formState[item.variable_key]" style="width: 100%">
              <a-select-option
                v-for="option in item.options"
                :key="option.label"
                :value="option.label"
              >{{ option.label }}</a-select-option>
            </a-select>
          </template>
          <template v-if="item.variable_type == 'checkbox_switch'">
            <a-checkbox v-model:checked="formState[item.variable_key]">{{ item.variable_name }}</a-checkbox>
          </template>
        </a-form-item>
      </a-form>
    </div>
    <div class="from-footer">
      <div class="btn-box">
        <a-button @click="handleSave" block type="primary">{{ t('btn_submit') }}</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useClawbotChatStore } from '@/stores/modules/clawbot-chat'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-test.components.variable-modal.index')
const props = defineProps({
  robotKey: {
    type: String,
    default: ''
  }
})
const chatStore = useClawbotChatStore()
const chatVariables = computed(() => chatStore.chat_variables || {})
const waitVariables = computed(() => (chatVariables.value.wait_variables || []).map((item) => ({
  ...item,
  options: item.options ? JSON.parse(item.options) : []
})))
const fillVariables = computed(() => chatVariables.value.fill_variables || [])
const needFillVariable = computed(() => chatVariables.value.need_fill_variable || false)
const isEdit = ref(false)
const open = ref(false)
const formRef = ref(null)
const formState = reactive({})

const clearForm = () => {
  for (const key of Object.keys(formState)) delete formState[key]
}

const fillForm = (list, valueKey) => {
  clearForm()
  list.forEach((item) => {
    const value = item[valueKey]
    if (item.variable_type === 'checkbox_switch') {
      formState[item.variable_key] = value == 1
    } else {
      formState[item.variable_key] = value || ''
    }
  })
}

const handleEdit = () => {
  isEdit.value = true
  fillForm(fillVariables.value, 'value')
  open.value = true
}

const getLength = (item) => item.max_input_length > 0 ? +item.max_input_length : 50

watch(() => props.robotKey, () => {
  open.value = false
  isEdit.value = false
  clearForm()
})

watch(waitVariables, (list) => fillForm(list, 'default_value'), { immediate: true })
watch(needFillVariable, (needed) => {
  if (needed) {
    isEdit.value = false
    open.value = true
  }
}, { immediate: true })

const handleSave = async () => {
  try {
    await formRef.value.validate()
    const robotKey = props.robotKey
    const dialogueId = Number(chatStore.dialogue_id || 0)
    const sessionId = chatStore.chat_variables.session_id
    if (!robotKey || chatStore.robot.robot_key !== robotKey) return
    const chatPromptVariables = waitVariables.value.map((item) => ({
      ...item,
      value: item.variable_type === 'checkbox_switch'
        ? (formState[item.variable_key] ? 1 : 0)
        : formState[item.variable_key]
    }))
    const variablesKey = `chat_prompt_variables_${robotKey}`
    if (dialogueId > 0) {
      await chatStore.handleEditVariables({ chat_prompt_variables: chatPromptVariables })
      if (props.robotKey !== robotKey || chatStore.robot.robot_key !== robotKey ||
        Number(chatStore.dialogue_id || 0) !== dialogueId || chatStore.chat_variables.session_id !== sessionId) return
      localStorage.removeItem(variablesKey)
    } else {
      localStorage.setItem(variablesKey, JSON.stringify(chatPromptVariables))
      chatStore.chat_variables.fill_variables = chatPromptVariables
    }
    chatStore.chat_variables.need_fill_variable = false
    open.value = false
  } catch {
    // Keep the form open when validation or saving fails.
  }
}

defineExpose({ handleEdit })
</script>

<style lang="less">
.variable-modal-class .ant-modal .ant-modal-content {
  padding: 0;
}
</style>

<style lang="less" scoped>
.modal-header {
  height: 76px;
  display: flex;
  align-items: center;
  position: relative;
  .title {
    position: relative;
    z-index: 9;
    padding-left: 24px;
    color: #000000;
    font-size: 24px;
    font-weight: 600;
  }
  .banner-bg {
    height: 96px;
    position: absolute;
    left: 0;
    right: 0;
    top: -20px;
    overflow: hidden;
    img {
      width: 100%;
      height: 100%;
    }
  }
}
.form-body {
  padding: 16px 24px 0 24px;
  max-height: 500px;
  overflow-y: auto;
  &::v-deep(.ant-form-item) {
    margin-bottom: 16px;
  }
}
.from-footer {
  padding: 8px 24px 24px 24px;
}
</style>
