<template>
  <div class="clawbot-loading">
    <div class="loading-content">
      <div class="loading-spinner">
        <div class="spinner-ring"></div>
        <MessageOutlined class="spinner-icon" />
      </div>
      <div class="loading-text">{{ t('msg_loading') }}</div>
    </div>
  </div>
</template>

<script setup>
import { MessageOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.components.clawbot-loading')
</script>

<style lang="less" scoped>
.clawbot-loading {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7f8fb;
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.loading-spinner {
  position: relative;
  width: 72px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner-ring {
  position: absolute;
  inset: 0;
  border: 3px solid #e6ebf3;
  border-top-color: #2f64eb;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.spinner-icon {
  font-size: 28px;
  color: #98a2b3;
}

.loading-text {
  font-size: 15px;
  color: #667085;
  letter-spacing: 0.5px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
