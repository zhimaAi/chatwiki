<template>
  <div class="rate-limit-page" :class="{ embedded }">
    <div v-if="!embedded" class="page-title">{{ t('title') }}</div>

    <div class="page-body">
      <div v-if="pageLoading" class="page-loading">
        <a-spin />
      </div>

      <a-alert
        v-else-if="loadError"
        type="error"
        show-icon
        :message="loadError"
      ><template #action><a @click="loadConfig(robotKey)">{{ t('retry') }}</a></template></a-alert>

      <a-form
        v-else
        ref="formRef"
        class="rate-limit-card"
        layout="vertical"
        :model="formState"
        :rules="rules"
      >
        <div class="card-header">
          <span class="card-title">{{ t('config_title') }}</span>
          <a-switch
            :checked="formState.switch_status === 1"
            :disabled="submitLoading || !hasLoaded"
            @change="handleSwitchChange"
          />
        </div>

        <div class="card-content">
          <section class="limit-section">
            <div class="section-title">{{ t('five_minute.title') }}</div>
            <div class="section-desc">{{ t('five_minute.desc') }}</div>

            <a-form-item :label="t('five_minute.limit_label')" name="five_minute_limit">
              <a-input-number
                v-model:value="formState.five_minute_limit"
                class="full-width"
                :min="1"
              />
              <div class="field-help">{{ t('five_minute.limit_help') }}</div>
            </a-form-item>

            <a-form-item :label="t('trigger_action')" name="five_minute_reply_type">
              <a-radio-group v-model:value="formState.five_minute_reply_type">
                <a-radio :value="0">{{ t('action_no_reply') }}</a-radio>
                <a-radio :value="1">{{ t('action_reply_content') }}</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item
              v-if="formState.five_minute_reply_type === 1"
              name="five_minute_reply_content"
            >
              <a-textarea
                v-model:value="formState.five_minute_reply_content"
                :auto-size="{ minRows: 3, maxRows: 6 }"
                :placeholder="t('reply_placeholder')"
              />
            </a-form-item>
          </section>

          <section class="limit-section daily-section">
            <div class="section-title">{{ t('daily.title') }}</div>
            <div class="section-desc">{{ t('daily.desc') }}</div>

            <a-form-item :label="t('daily.limit_label')" name="daily_limit">
              <a-input-number
                v-model:value="formState.daily_limit"
                class="full-width"
                :min="1"
              />
              <div class="field-help">{{ t('daily.limit_help') }}</div>
            </a-form-item>

            <a-form-item :label="t('trigger_action')" name="daily_reply_type">
              <a-radio-group v-model:value="formState.daily_reply_type">
                <a-radio :value="0">{{ t('action_no_reply') }}</a-radio>
                <a-radio :value="1">{{ t('action_reply_content') }}</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item
              v-if="formState.daily_reply_type === 1"
              name="daily_reply_content"
            >
              <a-textarea
                v-model:value="formState.daily_reply_content"
                :auto-size="{ minRows: 3, maxRows: 6 }"
                :placeholder="t('reply_placeholder')"
              />
            </a-form-item>
          </section>
        </div>

        <div class="card-footer">
          <a-button
            type="primary"
            :loading="submitLoading"
            :disabled="!hasLoaded"
            @click="handleSave"
          >
            {{ t('save_config') }}
          </a-button>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useRobotStore } from '@/stores/modules/robot'
import { getRobotRateLimitConf, saveRobotRateLimitConf } from '@/api/robot/index'

const props = defineProps({
  robotKey: { type: String, default: '' },
  embedded: {
    type: Boolean,
    default: false
  }
})
const emit = defineEmits(['status-change'])

const { t } = useI18n('views.robot.robot-config.rate-limit.index')
const route = useRoute()
const clawbotStore = useClawbotStore()
const robotStore = useRobotStore()

const createDefaultConfig = (robotKey = '') => ({
  robot_key: robotKey,
  switch_status: 0,
  five_minute_limit: 10,
  five_minute_reply_type: 1,
  five_minute_reply_content: '',
  daily_limit: 100,
  daily_reply_type: 0,
  daily_reply_content: ''
})

const normalizeConfig = (raw = {}, robotKey = '') => ({
  robot_key: String(raw.robot_key || robotKey || ''),
  switch_status: Number(raw.switch_status ?? 0),
  five_minute_limit: Number(raw.five_minute_limit ?? 10),
  five_minute_reply_type: Number(raw.five_minute_reply_type ?? 1),
  five_minute_reply_content: String(raw.five_minute_reply_content ?? ''),
  daily_limit: Number(raw.daily_limit ?? 100),
  daily_reply_type: Number(raw.daily_reply_type ?? 0),
  daily_reply_content: String(raw.daily_reply_content ?? '')
})

const formRef = ref(null)
const pageLoading = ref(false)
const submitLoading = ref(false)
const hasLoaded = ref(false)
const loadError = ref('')
const formState = reactive(createDefaultConfig())
let requestVersion = 0
let alive = true

const robotKey = computed(() => {
  return String(
    props.robotKey ||
    route.query.robot_key ||
    clawbotStore.currentAssistant?.robot_key ||
    robotStore.robotInfo?.robot_key ||
    ''
  )
})

const validatePositiveInteger = (_rule, value) => {
  const numberValue = Number(value)
  if (value === '' || value === null || value === undefined || !Number.isInteger(numberValue) || numberValue <= 0) {
    return Promise.reject(new Error(t('validator_positive_integer')))
  }
  return Promise.resolve()
}

const validateReplyContent = (replyTypeKey) => (_rule, value) => {
  const shouldValidate = formState.switch_status === 1 && formState[replyTypeKey] === 1
  if (shouldValidate && !String(value || '').trim()) {
    return Promise.reject(new Error(t('validator_reply_required')))
  }
  return Promise.resolve()
}

const rules = {
  five_minute_limit: [{ validator: validatePositiveInteger, trigger: ['change', 'blur'] }],
  five_minute_reply_content: [
    { validator: validateReplyContent('five_minute_reply_type'), trigger: ['change', 'blur'] }
  ],
  daily_limit: [{ validator: validatePositiveInteger, trigger: ['change', 'blur'] }],
  daily_reply_content: [
    { validator: validateReplyContent('daily_reply_type'), trigger: ['change', 'blur'] }
  ]
}

const applyConfig = (raw = {}, currentRobotKey = robotKey.value) => {
  Object.assign(formState, normalizeConfig(raw, currentRobotKey))
}

const buildPayload = () => ({
  robot_key: robotKey.value,
  switch_status: Number(formState.switch_status),
  five_minute_limit: Number(formState.five_minute_limit),
  five_minute_reply_type: Number(formState.five_minute_reply_type),
  five_minute_reply_content: String(formState.five_minute_reply_content || '').trim(),
  daily_limit: Number(formState.daily_limit),
  daily_reply_type: Number(formState.daily_reply_type),
  daily_reply_content: String(formState.daily_reply_content || '').trim()
})

const loadConfig = async (currentRobotKey) => {
  const currentVersion = ++requestVersion
  pageLoading.value = false
  submitLoading.value = false
  hasLoaded.value = false
  loadError.value = ''
  applyConfig(createDefaultConfig(currentRobotKey), currentRobotKey)

  if (!currentRobotKey) {
    loadError.value = t('msg_missing_robot_key')
    return
  }

  pageLoading.value = true
  try {
    const res = await getRobotRateLimitConf({ robot_key: currentRobotKey })
    if (currentVersion !== requestVersion) return

    if (Number(res?.res) !== 0) {
      loadError.value = res?.msg || t('msg_load_failed')
      return
    }

    applyConfig({
      ...createDefaultConfig(currentRobotKey),
      ...(res?.data || {})
    }, currentRobotKey)
    hasLoaded.value = true
  } catch {
    if (currentVersion !== requestVersion) return
    loadError.value = t('msg_load_failed')
  } finally {
    if (currentVersion === requestVersion) {
      pageLoading.value = false
    }
  }
}

const validateForm = async () => {
  try {
    await formRef.value?.validate()
    return true
  } catch {
    return false
  }
}

const submitConfig = async ({ previousSwitch = null } = {}) => {
  submitLoading.value = true
  const payload = buildPayload()
  const submittedRobotKey = payload.robot_key
  const submittedVersion = requestVersion

  try {
    const res = await saveRobotRateLimitConf(payload)
    if (!alive || submittedVersion !== requestVersion || submittedRobotKey !== robotKey.value) return false

    if (Number(res?.res) !== 0) {
      if (previousSwitch !== null) {
        formState.switch_status = previousSwitch
      }
      message.error(res?.msg || t('msg_save_failed'))
      return false
    }

    applyConfig({
      ...payload,
      ...(res?.data || {})
    })
    message.success(t('msg_save_success'))
    return true
  } catch {
    if (!alive || submittedVersion !== requestVersion || submittedRobotKey !== robotKey.value) return false

    if (previousSwitch !== null) {
      formState.switch_status = previousSwitch
    }
    message.error(t('msg_save_failed'))
    return false
  } finally {
    if (alive && submittedVersion === requestVersion) submitLoading.value = false
  }
}

const handleSwitchChange = async (checked) => {
  if (!alive || !hasLoaded.value || submitLoading.value) return false

  const version = requestVersion
  const previousSwitch = formState.switch_status
  formState.switch_status = checked ? 1 : 0
  submitLoading.value = true

  const isValid = await validateForm()
  if (!alive || version !== requestVersion) return false
  if (!isValid) {
    formState.switch_status = previousSwitch
    submitLoading.value = false
    return false
  }

  return submitConfig({ previousSwitch })
}

const handleSave = async () => {
  if (!alive || !hasLoaded.value || submitLoading.value) return
  const version = requestVersion
  submitLoading.value = true
  const isValid = await validateForm()
  if (!alive || version !== requestVersion) return
  if (!isValid) { submitLoading.value = false; return }
  await submitConfig()
}

watch(
  robotKey,
  (currentRobotKey) => {
    loadConfig(currentRobotKey)
  },
  { immediate: true }
)
watch(
  [hasLoaded, pageLoading, submitLoading, loadError, () => formState.switch_status],
  () => emit('status-change', {
    ready: hasLoaded.value,
    loading: pageLoading.value || submitLoading.value,
    enabled: formState.switch_status === 1,
    error: loadError.value
  }),
  { immediate: true }
)
defineExpose({ setEnabled: handleSwitchChange })
onBeforeUnmount(() => { alive = false; requestVersion++ })
</script>

<style lang="less" scoped>
.rate-limit-page {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #fff;
}

.page-title {
  height: 56px;
  padding: 0 24px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
  border-bottom: 1px solid #edf1f6;
  color: #101828;
  font-size: 16px;
  font-weight: 600;
}

.page-body {
  flex: 1;
  min-height: 0;
  padding: 16px 24px;
  overflow-y: auto;
  background: #f5f7fa;
}

.page-loading {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.rate-limit-card {
  min-height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid #e6ebf3;
  border-radius: 10px;
  background: #fff;
}

.card-header {
  height: 54px;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  border-bottom: 1px solid #edf1f6;
}

.card-title,
.section-title {
  color: #101828;
  font-size: 15px;
  font-weight: 600;
}

.card-content {
  flex: 1;
  padding: 24px;
}

.limit-section {
  padding-bottom: 24px;
}

.daily-section {
  padding-top: 24px;
  border-top: 1px solid #edf1f6;
}

.section-desc,
.field-help {
  color: #98a2b3;
  font-size: 13px;
  line-height: 20px;
}

.section-desc {
  margin: 4px 0 16px;
}

.field-help {
  margin-top: 4px;
}

.full-width {
  width: 100%;
}

.card-content :deep(.ant-form-item) {
  margin-bottom: 16px;
}

.card-content :deep(.ant-form-item-label) {
  padding-bottom: 4px;
  font-weight: 500;
}

.card-footer {
  padding: 16px 24px;
  display: flex;
  justify-content: flex-end;
  flex-shrink: 0;
  border-top: 1px solid #edf1f6;
}

.embedded .page-body {
  height: 100%;
}

@media (max-width: 768px) {
  .page-body,
  .card-content {
    padding: 16px;
  }

  .card-header,
  .card-footer {
    padding-right: 16px;
    padding-left: 16px;
  }
}
</style>
