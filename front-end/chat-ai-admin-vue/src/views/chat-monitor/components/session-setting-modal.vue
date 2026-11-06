<template>
  <a-modal
    v-model:open="open"
    :title="t('title_session_setting')"
    :confirm-loading="saving"
    :ok-button-props="{ disabled: !loaded || loading }"
    @ok="handleOk"
  >
    <a-spin :spinning="loading">
      <div class="session-setting-row">
        <span>{{ t('label_auto_end') }}</span>
        <a-input-number
          v-model:value="hours"
          :min="0.1"
          :max="48"
          :step="0.1"
          :precision="1"
          :disabled="!loaded"
        />
        <span>{{ t('label_hour') }}</span>
        <span>{{ t('desc_auto_end') }}</span>
      </div>
    </a-spin>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { showGlobalErrorMsg } from '@/utils/http/axios/config'
import { getAdminConfig } from '@/api/robot'
import { saveSessionExpireTime } from '@/api/chat-monitor'
import { DEFAULT_SESSION_EXPIRE_HOURS, isValidSessionExpireHours, sessionExpireHoursFromSeconds } from '../session-expire'

const { t } = useI18n('views.chat-monitor')
const { t: tCommon } = useI18n()
const open = ref(false)
const loading = ref(false)
const loaded = ref(false)
const saving = ref(false)
const hours = ref(DEFAULT_SESSION_EXPIRE_HOURS)

const show = async () => {
  open.value = true
  loading.value = true
  loaded.value = false
  hours.value = DEFAULT_SESSION_EXPIRE_HOURS
  try {
    const res = await getAdminConfig()
    hours.value = sessionExpireHoursFromSeconds(res?.data?.session_expire_seconds)
    loaded.value = true
  } catch {
    // 请求封装已经弹出错误，这里只保持未加载，避免再提示一次。
  } finally {
    loading.value = false
  }
}

const handleOk = async () => {
  if (!loaded.value || loading.value) {
    return Promise.reject(new Error('session setting is not loaded'))
  }
  if (!isValidSessionExpireHours(hours.value)) {
    showGlobalErrorMsg(t('msg_invalid_session_hours'))
    return Promise.reject(new Error('invalid session hours'))
  }
  saving.value = true
  try {
    await saveSessionExpireTime({ session_expire_hours: hours.value })
    message.success(tCommon('common.saveSuccess'))
    open.value = false
  } catch (error) {
    return Promise.reject(error)
  } finally {
    saving.value = false
  }
}

defineExpose({ show })
</script>

<style lang="less" scoped>
.session-setting-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-height: 72px;
  color: #262626;
}
</style>
