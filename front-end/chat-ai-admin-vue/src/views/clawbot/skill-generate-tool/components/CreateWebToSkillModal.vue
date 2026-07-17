<template>
  <a-modal
    class="create-web-skill-modal"
    :open="visible"
    :title="t('modal_add_web_skill')"
    :width="472"
    :maskClosable="false"
    :destroyOnClose="false"
    @cancel="handleCancel"
  >
    <a-form class="create-form" layout="vertical">
      <a-form-item required :label="t('label_web_urls')">
        <a-textarea
          v-model:value="formState.urlsText"
          :auto-size="{ minRows: 4, maxRows: 8 }"
          :placeholder="t('placeholder_web_urls')"
          :disabled="submitLoading"
        />
      </a-form-item>

      <a-form-item :label="t('label_custom_prompt')" v-if="false">
        <a-textarea
          v-model:value="formState.custom_prompt"
          :auto-size="{ minRows: 4, maxRows: 8 }"
          :placeholder="t('placeholder_custom_prompt')"
          :disabled="submitLoading"
        />
      </a-form-item>

      <a-form-item required :label="t('label_generate_model')">
        <ModelSelect
          modelType="LLM"
          v-model:modeName="formState.use_model"
          v-model:modeId="formState.model_config_id"
          :placeholder="t('placeholder_select_model')"
        />
      </a-form-item>

      <a-form-item required :label="t('label_max_token')">
        <div class="max-token-control">
          <a-slider
            v-model:value="formState.max_token"
            :min="1"
            :max="102400"
            :disabled="submitLoading"
          />
          <a-input-number
            v-model:value="formState.max_token"
            :min="1"
            :max="102400"
            :precision="0"
            :disabled="submitLoading"
          />
        </div>
      </a-form-item>
    </a-form>

    <template #footer>
      <a-button :disabled="submitLoading" @click="handleCancel">{{ t('btn_cancel') }}</a-button>
      <a-button type="primary" :loading="submitLoading" @click="handleConfirm">{{ t('btn_confirm') }}</a-button>
    </template>
  </a-modal>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import ModelSelect from '@/components/model-select/model-select.vue'
import { createWebToSkillTask } from '@/api/clawbot'

const { t } = useI18n('views.clawbot.skill-generate-tool.index')

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const submitLoading = ref(false)
const formState = reactive({
  urlsText: '',
  custom_prompt: '',
  model_config_id: '',
  use_model: '',
  max_token: 32768
})

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      resetForm()
    }
  }
)

const resetForm = () => {
  formState.urlsText = ''
  formState.custom_prompt = ''
  formState.model_config_id = ''
  formState.use_model = ''
  formState.max_token = 32768
  submitLoading.value = false
}

const parseUrls = () => {
  const urls = formState.urlsText
    .split(/\r?\n/)
    .map((item) => item.trim())
    .filter(Boolean)

  return [...new Set(urls)]
}

const validateUrls = (urls) => {
  if (!urls.length) {
    message.error(t('msg_enter_web_url'))
    return false
  }

  for (const item of urls) {
    try {
      const url = new URL(item)
      if (url.protocol !== 'http:' && url.protocol !== 'https:') {
        message.error(t('msg_invalid_web_url'))
        return false
      }
    } catch {
      message.error(t('msg_invalid_web_url'))
      return false
    }
  }
  return true
}

const validateForm = () => {
  const urls = parseUrls()
  if (!validateUrls(urls)) {
    return null
  }
  if (!formState.model_config_id || !formState.use_model) {
    message.error(t('msg_select_model'))
    return null
  }
  if (
    !Number.isInteger(Number(formState.max_token)) ||
    Number(formState.max_token) < 1 ||
    Number(formState.max_token) > 102400
  ) {
    message.error(t('msg_invalid_max_token'))
    return null
  }
  return urls
}

const handleConfirm = async () => {
  const urls = validateForm()
  if (!urls) {
    return
  }

  submitLoading.value = true
  try {
    const res = await createWebToSkillTask({
      urls,
      custom_prompt: formState.custom_prompt,
      model_config_id: Number(formState.model_config_id),
      use_model: formState.use_model,
      max_token: Number(formState.max_token)
    })
    if (res && (res.res === 0 || res.code === 0)) {
      message.success(t('msg_task_created'))
      emit('confirm', res.data)
      emit('update:visible', false)
    } else {
      message.error(res?.msg || t('msg_task_create_failed'))
    }
  } catch (error) {
    console.error('创建Web转Skill任务失败', error)
  } finally {
    submitLoading.value = false
  }
}

const handleCancel = () => {
  if (!submitLoading.value) {
    emit('update:visible', false)
  }
}
</script>

<style lang="less" scoped>
.create-web-skill-modal {
  :deep(.ant-modal-content) {
    border-radius: 16px;
  }

  :deep(.ant-modal-header) {
    margin-bottom: 0;
    padding: 16px 24px;
  }

  :deep(.ant-modal-title) {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }

  :deep(.ant-modal-body) {
    padding: 24px 32px;
  }

  :deep(.ant-modal-footer) {
    margin-top: 0;
    padding: 10px 24px;
  }
}

.create-form {
  :deep(.ant-form-item) {
    margin-bottom: 16px;
  }

  :deep(.ant-form-item:last-child) {
    margin-bottom: 0;
  }

  :deep(.ant-form-item-label) {
    padding-bottom: 4px;
  }

  :deep(.ant-form-item-label > label) {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
  }

  :deep(.ant-input),
  :deep(.ant-input-number),
  :deep(.ant-select-selector) {
    border-radius: 6px;
  }

  :deep(.ant-input-number) {
    width: 100%;
  }
}

.max-token-control {
  display: flex;
  align-items: center;

  :deep(.ant-slider) {
    flex: 1;
    min-width: 0;
  }

  :deep(.ant-input-number) {
    width: 120px;
    margin-left: 16px;
  }
}
</style>
