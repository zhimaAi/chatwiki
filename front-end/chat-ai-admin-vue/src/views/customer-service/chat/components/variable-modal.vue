<template>
  <a-modal
    wrapClassName="customer-service-variable-modal"
    :open="props.open"
    :title="null"
    :footer="null"
    :closable="!props.loading"
    :keyboard="!props.loading"
    :maskClosable="false"
    :width="560"
    :bodyStyle="{ padding: 0 }"
    @cancel="handleClose"
  >
    <div class="modal-header">
      <div class="banner-bg" aria-hidden="true"></div>
      <div class="title">{{ t('title_fill_info') }}</div>
    </div>
    <div class="form-body">
      <a-form :model="formState" ref="formRef" layout="vertical">
        <a-form-item
          v-for="item in fields"
          :key="item.variable_key"
          :label="item.variable_type === 'checkbox_switch' ? null : item.variable_name"
          :name="item.variable_key"
          :rules="[{ required: item.must_input == 1, message: t('msg_input_required', { name: item.variable_name }) }]"
        >
          <a-input
            v-if="item.variable_type === 'input_string'"
            v-model:value="formState[item.variable_key]"
            :disabled="props.loading"
            :maxLength="getLength(item)"
            :placeholder="t('ph_input')"
          />
          <a-input-number
            v-else-if="item.variable_type === 'input_number'"
            v-model:value="formState[item.variable_key]"
            :disabled="props.loading"
            style="width: 100%"
            stringMode
            :maxLength="50"
            :placeholder="t('ph_input')"
          />
          <a-select
            v-else-if="item.variable_type === 'select_one'"
            v-model:value="formState[item.variable_key]"
            :disabled="props.loading"
            style="width: 100%"
          >
            <a-select-option v-for="option in item.options" :key="option.label" :value="option.label">
              {{ option.label }}
            </a-select-option>
          </a-select>
          <a-checkbox
            v-else-if="item.variable_type === 'checkbox_switch'"
            v-model:checked="formState[item.variable_key]"
            :disabled="props.loading"
          >{{ item.variable_name }}</a-checkbox>
        </a-form-item>
      </a-form>
    </div>
    <div class="form-footer">
      <a-button block type="primary" :loading="props.loading || validating" @click="handleSubmit">
        {{ t('btn_submit') }}
      </a-button>
    </div>
  </a-modal>
</template>

<script setup>
import { computed, nextTick, onScopeDispose, reactive, ref, watch } from 'vue'
import { safeParseJson } from '@/utils/agent-process'
import { useI18n } from '@/hooks/web/useI18n'

// 只校验并提交变量副本，接口/首轮缓存及成功关闭归父页和客服聊天Store。
const { t } = useI18n('views.robot.robot-test.components.variable-modal.index')
const props = defineProps({
  open: { type: Boolean, default: false },
  variables: { type: Object, default: () => ({}) },
  loading: { type: Boolean, default: false }
})
const emit = defineEmits(['submit', 'close'])
const formRef = ref(null)
const formState = reactive({})
const validating = ref(false)
let formVersion = 0

const fields = computed(() => {
  const waiting = props.variables.wait_variables || []
  const definitions = waiting.length ? waiting : props.variables.fill_variables || []
  return definitions.map((item) => {
    const options = safeParseJson(item.options, [])
    return { ...item, options: Array.isArray(options) ? options.map((option) => ({ ...option })) : [] }
  })
})

const resetForm = () => {
  for (const key of Object.keys(formState)) delete formState[key]
  const filled = new Map((props.variables.fill_variables || []).map((item) => [item.variable_key, item]))
  for (const item of fields.value) {
    const value = filled.get(item.variable_key)?.value ?? item.default_value ?? ''
    formState[item.variable_key] = item.variable_type === 'checkbox_switch' ? value == 1 : value
  }
  nextTick(() => formRef.value?.clearValidate())
}

watch([() => props.open, () => props.variables], () => {
  // 重新打开或父页切换会话/机器人时换一份表单，旧异步校验不得提交到新会话。
  formVersion++
  validating.value = false
  resetForm()
}, { immediate: true, flush: 'sync' })

onScopeDispose(() => { formVersion++ })

const getLength = (item) => item.max_input_length > 0 ? Number(item.max_input_length) : 50

const handleClose = () => {
  if (props.loading) return
  formVersion++
  validating.value = false
  resetForm()
  emit('close')
}

const handleSubmit = async () => {
  if (!props.open || props.loading || validating.value) return
  const version = formVersion
  validating.value = true
  try {
    await formRef.value.validate()
    if (version !== formVersion || !props.open || props.loading) return
    const values = fields.value.map((item) => ({
      ...item,
      options: item.options.map((option) => ({ ...option })),
      value: item.variable_type === 'checkbox_switch'
        ? (formState[item.variable_key] ? 1 : 0)
        : formState[item.variable_key]
    }))
    // 不自行关闭/清空；保存失败时父页保持open，编辑内容仍可重试。
    emit('submit', values)
  } catch {
    // 必填错误由表单显示，保留当前输入。
  } finally {
    if (version === formVersion) validating.value = false
  }
}
</script>

<style lang="less">
// 弹层teleport到body，限定本组件专用类，沿原弹窗去除内容区外边距。
.customer-service-variable-modal .ant-modal .ant-modal-content {
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
    color: var(--app-text-primary);
    font-size: 24px;
    font-weight: 500;
  }
  .banner-bg {
    height: 96px;
    position: absolute;
    left: 0;
    right: 0;
    top: -20px;
    overflow: hidden;
    background: linear-gradient(180deg, var(--app-selected), var(--app-bg-elevated));
  }
}
.form-body {
  padding: 16px 24px 0;
  max-height: 500px;
  overflow-y: auto;
  :deep(.ant-form-item) {
    margin-bottom: 16px;
  }
}
.form-footer {
  padding: 8px 24px 24px;
}
</style>
