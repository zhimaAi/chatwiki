<template>
  <div class="e2b-page">
    <a-alert
      v-if="loadError"
      class="status-alert"
      type="error"
      show-icon
      :message="loadError"
    />

    <section class="e2b-card intro-card">
      <div class="card-title">
        <span class="title-icon" aria-hidden="true">
          <svg-icon name="e2b-intro"></svg-icon>
        </span>
        <span>{{ t('e2b.intro_title') }}</span>
      </div>

      <div class="card-desc">{{ t('e2b.intro_desc') }}</div>

      <div class="intro-grid">
        <div class="intro-item">
          <div class="intro-item-title">{{ t('e2b.feature_isolation_title') }}</div>
          <div class="intro-item-desc">{{ t('e2b.feature_isolation_desc') }}</div>
        </div>
        <div class="intro-item">
          <div class="intro-item-title">{{ t('e2b.feature_template_title') }}</div>
          <div class="intro-item-desc">{{ t('e2b.feature_template_desc') }}</div>
        </div>
        <div class="intro-item">
          <div class="intro-item-title">{{ t('e2b.feature_command_title') }}</div>
          <div class="intro-item-desc">{{ t('e2b.feature_command_desc') }}</div>
        </div>
      </div>
    </section>

    <section class="e2b-card config-card">
      <div class="card-title-row">
        <div class="card-title">
          <span class="title-icon" aria-hidden="true">
            <svg-icon name="system-setting"></svg-icon>
          </span>
          <span>{{ t('e2b.config_title') }}</span>
        </div>
        <a
          class="operation-link"
          :href="E2B_DOC_URL"
          target="_blank"
          rel="noopener noreferrer"
        >
          {{ t('e2b.operation_guide') }}
        </a>
      </div>

      <div v-if="pageLoading" class="config-loading">
        <a-spin />
      </div>

      <div v-else class="config-content">
        <div class="switch-row">
          <div class="switch-copy">
            <div class="switch-title">{{ t('e2b.switch_title') }}</div>
            <div class="switch-desc">{{ t('e2b.switch_desc') }}</div>
          </div>

          <button
            type="button"
            class="e2b-switch"
            :class="{ checked: isEnabled, loading: switchLoading }"
            :disabled="switchLoading || submitLoading"
            @click="handleSwitchClick"
          >
            <span class="switch-knob" />
            <span class="switch-text">{{ isEnabled ? t('switch_on') : t('switch_off') }}</span>
          </button>
        </div>

        <div v-if="isEnabled" class="preview-card">
          <div class="preview-head">
            <button type="button" class="edit-setting-btn" @click="handleEdit">
              <span class="edit-setting-icon" aria-hidden="true">
                <svg-icon name="system-setting"></svg-icon>
              </span>
              <span>{{ t('e2b.btn_edit') }}</span>
            </button>
          </div>

          <div class="preview-row">
            <div class="preview-item">
              <span class="preview-key">APIKey:</span>
              <span class="preview-value apikey">{{ displayValue(config.api_key) }}</span>
            </div>

            <div v-if="config.template" class="preview-item">
                <span class="preview-key">Template:</span>
                <span class="preview-value">{{ displayValue(config.template) }}</span>
              </div>
          </div>

          <div class="preview-row">
            <div class="preview-item">
              <span class="preview-key">APIBaseURL:</span>
              <span class="preview-value">{{ displayValue(config.api_base_url) }}</span>
            </div>
          </div>

          <div class="preview-row">
            <div class="preview-item">
              <span class="preview-key">SandboxDomain:</span>
              <span class="preview-value">{{ displayValue(config.sandbox_domain) }}</span>
            </div>
          </div>

          <div class="preview-row">
            <div class="preview-item">
              <span class="preview-key">Timeout:</span>
              <span class="preview-value">{{ String(config.timeout || 0) }}</span>
            </div>
            <div class="preview-item">
              <span class="preview-key">COMMAND_TIMEOUT:</span>
              <span class="preview-value">{{ String(config.command_timeout || 0) }}</span>
            </div>
            <div class="preview-item">
              <span class="preview-key">COMMAND_USER:</span>
              <span class="preview-value">{{ displayValue(config.command_user) }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <E2bSettingModal
      v-model:visible="modalOpen"
      :config="config"
      :loading="submitLoading"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { storeToRefs } from 'pinia'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { getE2bConf, saveE2bConf } from '@/api/clawbot'
import { E2B_DOC_URL } from '@/constants'
import E2bSettingModal from './e2b-setting-modal.vue'

const { t } = useI18n('views.clawbot.settings.index')
const clawbotStore = useClawbotStore()
const { currentAssistant } = storeToRefs(clawbotStore)

const createDefaultConfig = (robotKey = '') => ({
  robot_key: robotKey,
  switch_status: 0,
  api_key: '',
  api_base_url: '',
  sandbox_domain: '',
  template: '',
  timeout: 0,
  command_timeout: 0,
  command_user: ''
})

const normalizeConfig = (raw = {}) => ({
  robot_key: raw.robot_key || '',
  switch_status: Number(raw.switch_status || 0),
  api_key: raw.api_key || '',
  api_base_url: raw.api_base_url || '',
  sandbox_domain: raw.sandbox_domain || '',
  template: raw.template || '',
  timeout: Number(raw.timeout || 0),
  command_timeout: Number(raw.command_timeout || 0),
  command_user: raw.command_user || ''
})

const pageLoading = ref(false)
const switchLoading = ref(false)
const submitLoading = ref(false)
const modalOpen = ref(false)
const loadError = ref('')
const config = ref(createDefaultConfig())
let requestVersion = 0

const robotKey = computed(() => currentAssistant.value?.robot_key || '')
const isEnabled = computed(() => Number(config.value.switch_status || 0) === 1)

const displayValue = (value) => {
  return value ? String(value) : '—'
}

const applyConfig = (nextConfig = {}) => {
  config.value = normalizeConfig(nextConfig)
}

const loadConfig = async (currentRobotKey) => {
  requestVersion += 1
  const currentVersion = requestVersion

  if (!currentRobotKey) {
    loadError.value = ''
    applyConfig(createDefaultConfig())
    return
  }

  pageLoading.value = true
  loadError.value = ''
  try {
    const res = await getE2bConf({ robot_key: currentRobotKey })
    if (currentVersion !== requestVersion) return

    if (res?.res !== 0) {
      loadError.value = res?.msg || t('e2b.msg_load_failed')
      applyConfig(createDefaultConfig(currentRobotKey))
      return
    }

    applyConfig({
      ...createDefaultConfig(currentRobotKey),
      ...(res?.data || {})
    })
  } catch {
    if (currentVersion !== requestVersion) return
    loadError.value = t('e2b.msg_load_failed')
    applyConfig(createDefaultConfig(currentRobotKey))
  } finally {
    if (currentVersion === requestVersion) {
      pageLoading.value = false
    }
  }
}

watch(
  robotKey,
  (currentRobotKey) => {
    modalOpen.value = false
    loadConfig(currentRobotKey)
  },
  { immediate: true }
)

const handleEdit = () => {
  modalOpen.value = true
}

const handleSwitchClick = async () => {
  if (!robotKey.value || switchLoading.value || submitLoading.value) {
    return
  }

  if (!isEnabled.value) {
    modalOpen.value = true
    return
  }

  switchLoading.value = true
  try {
    const res = await saveE2bConf({
      robot_key: robotKey.value,
      switch_status: 0
    })

    if (res?.res !== 0) {
      message.error(res?.msg || t('e2b.msg_save_failed'))
      return
    }

    applyConfig({
      ...config.value,
      ...(res?.data || {}),
      robot_key: robotKey.value
    })
    loadError.value = ''
    message.success(t('e2b.msg_close_success'))
  } catch (error) {
    console.error(error)
    // message.error(t('e2b.msg_save_failed'))
  } finally {
    switchLoading.value = false
  }
}

const handleSubmit = async (payload) => {
  if (!robotKey.value) {
    return
  }

  submitLoading.value = true
  try {
    const res = await saveE2bConf({
      robot_key: robotKey.value,
      switch_status: 1,
      ...payload
    })

    if (res?.res !== 0) {
      message.error(res?.msg || t('e2b.msg_save_failed'))
      return
    }

    applyConfig({
      ...config.value,
      ...(res?.data || {}),
      robot_key: robotKey.value
    })
    loadError.value = ''
    modalOpen.value = false
    message.success(t('e2b.msg_save_success'))
  } catch (error) {
    // message.error(t('e2b.msg_save_failed'))
    console.error(error)
  } finally {
    submitLoading.value = false
  }
}
</script>

<style lang="less" scoped>
.e2b-page {
  height: 100%;
  overflow-y: auto;
  padding: 24px 24px 32px;
  background: #fff;
}

.status-alert {
  margin-bottom: 16px;
}

.e2b-card {
  border-radius: 12px;
  background: #f2f4f7;
}

.intro-card,
.config-card {
  padding: 24px;
}

.config-card {
  margin-top: 16px;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
  font-size: 16px;
  line-height: 24px;
  font-weight: 600;
}

.title-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  color: #262626;

  :deep(.svg-action) {
    display: inline-flex;
  }

  :deep(.svg-icon) {
    width: 16px;
    height: 16px;
    display: block;
  }
}

.operation-link {
  color: #2475fc;
  font-size: 14px;
  line-height: 22px;
  text-decoration: none;
}

.card-desc {
  margin-top: 16px;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

.intro-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
}

.intro-item {
  min-height: 132px;
  padding: 16px;
  border: 1px solid #e5e7ec;
  border-radius: 6px;
  background: #fff;
}

.intro-item-title {
  color: #262626;
  font-size: 16px;
  line-height: 24px;
  font-weight: 600;
}

.intro-item-desc {
  margin-top: 10px;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

.config-loading {
  min-height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.config-content {
  margin-top: 16px;
}

.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.switch-copy {
  flex: 1;
  min-width: 0;
}

.switch-title {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
  font-weight: 600;
}

.switch-desc {
  margin-top: 4px;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

.e2b-switch {
  min-width: 45px;
  height: 22px;
  padding: 1px 2px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border: none;
  border-radius: 16px;
  background: #bfbfbf;
  color: #fff;
  cursor: pointer;
  transition: background-color 0.2s ease, opacity 0.2s ease;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.7;
  }

  &.checked {
    flex-direction: row-reverse;
    background: #2475fc;
  }
}

.switch-text {
  font-size: 12px;
  line-height: 20px;
  white-space: nowrap;
}

.switch-knob {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  border-radius: 77px;
  background: #fff;
  box-shadow: 0 2px 2px rgba(0, 35, 11, 0.2);
}

.preview-card {
  margin-top: 16px;
  padding: 16px 20px;
  border-radius: 10px;
  background: #fff;
}

.preview-head {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;

  .edit-setting-btn{
    position: absolute;
    right: 0;
    top: 0;
  }
}

.preview-row {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
  &:last-child {
    margin-bottom: 0;
  }
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

.preview-key {
  color: #262626;
  font-weight: 500;
  white-space: nowrap;
}

.preview-value {
  min-width: 0;
  color: #595959;
  word-break: break-all;
}

.preview-value.apikey{
  width: 300px;
}

.edit-setting-btn {
  padding: 0;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: none;
  background: transparent;
  color: #2475fc;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
}

.edit-setting-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;

  :deep(.svg-action) {
    display: inline-flex;
  }

  :deep(.svg-icon) {
    width: 16px;
    height: 16px;
    display: block;
  }
}

@media (max-width: 1200px) {
  .intro-grid {
    grid-template-columns: 1fr;
  }

  .preview-head {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .e2b-page {
    padding: 16px;
  }

  .intro-card,
  .config-card {
    padding: 16px;
  }

  .switch-row {
    align-items: flex-start;
    flex-direction: column;
  }

  .preview-head,
  .preview-row {
    align-items: flex-start;
    flex-direction: column;
    gap: 8px;
  }
}
</style>
