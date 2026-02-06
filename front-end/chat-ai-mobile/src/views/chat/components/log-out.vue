<style lang="less" scoped>
.log-out-box {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 9px 12px;
  border-radius: 12px 0 0 12px;
  background: #FFF;
  box-shadow: 0 4px 8px 0 #00000014;

  .log-out-icon {
    font-size: 16px;
  }

  .log-text {
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 1;
    overflow: hidden;
    color: #3a4559;
    text-overflow: ellipsis;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 16px;
    padding-left: 8px;
  }
}
</style>

<template>
  <div class="log-out-box">
    <svg-icon name="log-out" class="log-out-icon" />
    <span class="log-text">{{ t('btn_logout') }}</span>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.log-out')

const emit = defineEmits([ 'send' ])


defineExpose({
})
</script>
