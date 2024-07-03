<style lang="less" scoped>
.layout-footer {
  padding: 16px 0;
  .copyright-text {
    line-height: 20px;
    font-size: 12px;
    color: #8c8c8c;
    text-align: center;
  }
  .user-agreement-box {
    display: flex;
    justify-content: center;
    line-height: 20px;
    margin-bottom: 4px;
    font-size: 12px;

    .link-item {
      margin: 0 8px;
      color: #595959;
    }
  }
}
</style>

<template>
  <div class="layout-footer">
    <!-- <div class="user-agreement-box">
      <a class="link-item">用户协议</a><a class="link-item">隐私政策</a>
    </div> -->
    <div class="copyright-text">{{ t('common.copyright') }}</div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()
</script>
