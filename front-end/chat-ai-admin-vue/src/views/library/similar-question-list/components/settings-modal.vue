<template>
  <a-modal
    :open="visible"
    :title="t('modal_title')"
    @ok="handleConfirm"
    @cancel="handleCancel"
    :width="550"
  >
    <div class="settings-modal-content">
      <a-alert
        class="settings-alert"
        :message="t('alert_message')"
        type="info"
        show-icon
      />
      <div class="threshold-setting">
        <div class="threshold-label">{{ t('threshold_label') }}</div>
        <div class="threshold-control">
          <a-slider
            v-model:value="thresholdValue"
            :min="0"
            :max="1"
            :step="0.01"
            class="threshold-slider"
          />
          <a-input-number
            v-model:value="thresholdValue"
            :min="0"
            :max="1"
            :step="0.01"
            :precision="2"
            class="threshold-input"
          />
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getUserConfig, saveUserConfig } from '@/api/library'

const { t } = useI18n('views.library.similar-question-list.components.settings-modal')

const emit = defineEmits(['confirm', 'cancel'])

const visible = ref(false)
const thresholdValue = ref(0.8)

const open = ({ defaultValue = 0.8 }) => {
  thresholdValue.value = defaultValue
  visible.value = true
  // 打开弹窗时获取用户配置
  fetchUserConfig()
}

// 获取用户配置
const fetchUserConfig = () => {
  getUserConfig().then((res) => {
    if (res.data && res.data.qa_merge_similarity !== undefined) {
      thresholdValue.value = res.data.qa_merge_similarity
    }
  }).catch((err) => {
    console.error('getUserConfig error:', err)
  })
}

const handleConfirm = async () => {
  try {
    await saveUserConfig({
      qa_merge_similarity: thresholdValue.value
    })
    emit('confirm', thresholdValue.value)
    visible.value = false
  } catch (err) {
    console.error('saveUserConfig error:', err)
  }
}

const handleCancel = () => {
  emit('cancel')
  visible.value = false
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.settings-modal-content {
  padding: 8px 0;

  .settings-alert {
    margin-bottom: 24px;
  }

  .threshold-setting {
    .threshold-label {
      font-size: 14px;
      color: #262626;
      margin-bottom: 12px;
    }

    .threshold-control {
      display: flex;
      align-items: center;
      gap: 16px;

      .threshold-slider {
        flex: 1;
      }

      .threshold-input {
        width: 80px;
      }
    }

    .threshold-default {
      margin-top: 8px;
      font-size: 12px;
      color: #8c8c8c;
    }
  }
}
</style>
