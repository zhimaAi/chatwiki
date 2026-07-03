<template>
  <a-modal
    :open="visible"
    :footer="null"
    :closable="false"
    :maskClosable="false"
    :width="472"
    :bodyStyle="{ padding: '0' }"
    centered
    @cancel="handleCancel"
  >
    <div class="e2b-modal">
      <div class="modal-header">
        <div class="modal-title">{{ t('e2b.modal_title') }}</div>
        <button type="button" class="close-btn" @click="handleCancel">
          <CloseOutlined />
        </button>
      </div>

      <div class="modal-body">
        <a-form class="e2b-form">
          <a-form-item v-bind="validateInfos.api_key">
            <div class="field-label"><span class="required">*</span><span>APIKey</span></div>
            <a-input
              v-model:value="formState.api_key"
              :maxlength="500"
              :placeholder="t('e2b.placeholder.api_key')"
            />
            <div class="field-tip">{{ t('e2b.help.api_key') }}</div>
          </a-form-item>

          <a-form-item v-bind="validateInfos.api_base_url">
            <div class="field-label"><span class="required">*</span><span>APIBaseURL</span></div>
            <a-input
              v-model:value="formState.api_base_url"
              :maxlength="1000"
              :placeholder="t('e2b.placeholder.api_base_url')"
            />
            <div class="field-tip">{{ t('e2b.help.api_base_url') }}</div>
          </a-form-item>

          <a-form-item v-bind="validateInfos.sandbox_domain">
            <div class="field-label"><span class="required">*</span><span>SandboxDomain</span></div>
            <a-input
              v-model:value="formState.sandbox_domain"
              :maxlength="1000"
              :placeholder="t('e2b.placeholder.sandbox_domain')"
            />
            <div class="field-tip">{{ t('e2b.help.sandbox_domain') }}</div>
          </a-form-item>

          <a-form-item v-bind="validateInfos.template">
            <div class="field-label"><span class="required">*</span><span>Template</span></div>
            <a-input
              v-model:value="formState.template"
              :maxlength="500"
              :placeholder="t('e2b.placeholder.template')"
            />
            <div class="field-tip">{{ t('e2b.help.template') }}</div>
          </a-form-item>

          <a-form-item v-bind="validateInfos.timeout">
            <div class="field-label"><span class="required">*</span><span>Timeout</span></div>
            <a-input-number
              v-model:value="formState.timeout"
              class="number-input"
              :min="1"
              :precision="0"
              :controls="false"
              :placeholder="t('e2b.placeholder.timeout')"
            />
            <div class="field-tip">{{ t('e2b.help.timeout') }}</div>
          </a-form-item>

          <a-form-item v-bind="validateInfos.command_timeout">
            <div class="field-label"><span class="required">*</span><span>COMMAND_TIMEOUT</span></div>
            <a-input-number
              v-model:value="formState.command_timeout"
              class="number-input"
              :min="1"
              :precision="0"
              :controls="false"
              :placeholder="t('e2b.placeholder.command_timeout')"
            />
            <div class="field-tip">{{ t('e2b.help.command_timeout') }}</div>
          </a-form-item>

          <a-form-item v-bind="validateInfos.command_user">
            <div class="field-label"><span class="required">*</span><span>COMMAND_USER</span></div>
            <a-input
              v-model:value="formState.command_user"
              :maxlength="100"
              :placeholder="t('e2b.placeholder.command_user')"
            />
            <div class="field-tip">{{ t('e2b.help.command_user') }}</div>
          </a-form-item>
        </a-form>
      </div>

      <div class="modal-footer">
        <a-button @click="handleCancel">{{ t('e2b.btn_cancel') }}</a-button>
        <a-button type="primary" :loading="loading" @click="handleConfirm">
          {{ t('e2b.btn_confirm') }}
        </a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, watch } from 'vue'
import { Form } from 'ant-design-vue'
import { CloseOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.settings.index')

const buildDefaultForm = (config = {}) => ({
  api_key: config.api_key || '',
  api_base_url: config.api_base_url || '',
  sandbox_domain: config.sandbox_domain || '',
  template: config.template || '',
  timeout: Number(config.timeout || 0) || undefined,
  command_timeout: Number(config.command_timeout || 0) || undefined,
  command_user: config.command_user || ''
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:visible', 'submit'])

const formState = reactive(buildDefaultForm())

const requireMessage = (label) => t('e2b.validator_required', { label })

const createStringRule = (label, max) => ([
  {
    validator: async (_, value) => {
      const normalized = String(value || '').trim()
      if (!normalized) {
        return Promise.reject(requireMessage(label))
      }
      if (normalized.length > max) {
        return Promise.reject(t('e2b.validator_max', { label, max }))
      }
      return Promise.resolve()
    },
    trigger: 'change'
  }
])

const rules = reactive({
  api_key: createStringRule('APIKey', 500),
  api_base_url: [
    {
      validator: async (_, value) => {
        const normalized = String(value || '').trim()
        if (!normalized) {
          return Promise.reject(requireMessage('APIBaseURL'))
        }
        if (normalized.length > 1000) {
          return Promise.reject(t('e2b.validator_max', { label: 'APIBaseURL', max: 1000 }))
        }

        try {
          const parsed = new URL(normalized)
          if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
            throw new Error('invalid protocol')
          }
          return Promise.resolve()
        } catch {
          return Promise.reject(t('e2b.validator_url'))
        }
      },
      trigger: 'change'
    }
  ],
  sandbox_domain: createStringRule('SandboxDomain', 1000),
  template: createStringRule('Template', 500),
  timeout: [
    {
      required: true,
      message: requireMessage('Timeout'),
      trigger: 'change'
    },
    {
      validator: async (_, value) => {
        if (value === undefined || value === null || value === '') {
          return Promise.resolve()
        }
        if (!Number.isInteger(Number(value)) || Number(value) <= 0) {
          return Promise.reject(t('e2b.validator_positive_integer', { label: 'Timeout' }))
        }
        return Promise.resolve()
      },
      trigger: 'change'
    }
  ],
  command_timeout: [
    {
      required: true,
      message: requireMessage('COMMAND_TIMEOUT'),
      trigger: 'change'
    },
    {
      validator: async (_, value) => {
        if (value === undefined || value === null || value === '') {
          return Promise.resolve()
        }
        if (!Number.isInteger(Number(value)) || Number(value) <= 0) {
          return Promise.reject(t('e2b.validator_positive_integer', { label: 'COMMAND_TIMEOUT' }))
        }
        return Promise.resolve()
      },
      trigger: 'change'
    }
  ],
  command_user: createStringRule('COMMAND_USER', 100)
})

const { validate, validateInfos, clearValidate } = Form.useForm(formState, rules)

const syncFormState = (config = {}) => {
  const nextState = buildDefaultForm(config)
  Object.keys(nextState).forEach((key) => {
    formState[key] = nextState[key]
  })
}

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      syncFormState(props.config)
      clearValidate()
    }
  }
)

const handleCancel = () => {
  clearValidate()
  emit('update:visible', false)
}

const handleConfirm = () => {
  validate().then(() => {
    emit('submit', {
      api_key: formState.api_key.trim(),
      api_base_url: formState.api_base_url.trim(),
      sandbox_domain: formState.sandbox_domain.trim(),
      template: formState.template.trim(),
      timeout: Number(formState.timeout),
      command_timeout: Number(formState.command_timeout),
      command_user: formState.command_user.trim()
    })
  })
}
</script>

<style lang="less" scoped>
.e2b-modal {
  border-radius: 16px;
  background: #fff;
}

.modal-header {
  padding: 0 0 16px 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title {
  color: #262626;
  font-size: 16px;
  line-height: 24px;
  font-weight: 500;
}

.close-btn {
  width: 16px;
  height: 16px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: #8c8c8c;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
}

.modal-footer {
  padding: 10px 24px 0;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  gap: 8px;
}

.e2b-form {
  :deep(.ant-form-item) {
    margin-bottom: 16px;
  }

  :deep(.ant-form-item-control-input) {
    min-height: auto;
  }

  :deep(.ant-form-item-explain-error) {
    margin-top: 4px;
    font-size: 12px;
    line-height: 20px;
  }

  :deep(.ant-input),
  :deep(.ant-input-affix-wrapper),
  :deep(.ant-input-number) {
    border-color: #d9d9d9;
    border-radius: 6px;
    box-shadow: none;
  }

  :deep(.ant-input),
  :deep(.ant-input-affix-wrapper) {
    padding: 4px 12px;
    font-size: 14px;
    line-height: 22px;
  }

  :deep(.ant-input-number) {
    width: 100%;
    padding: 4px 12px;
  }

  :deep(.ant-input-number-input) {
    height: 22px;
    padding: 0;
    font-size: 14px;
    line-height: 22px;
  }

  :deep(.ant-input::placeholder),
  :deep(.ant-input-number-input::placeholder) {
    color: #bfbfbf;
  }
}

.field-label {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-bottom: 4px;
  color: #262626;
  font-size: 14px;
  line-height: 22px;
}

.required {
  color: #fb363f;
}

.field-tip {
  margin-top: 2px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.number-input {
  width: 100%;
}
</style>
