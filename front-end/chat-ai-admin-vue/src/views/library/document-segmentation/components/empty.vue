<template>
  <div class="empty-box">
    <img src="../../../../assets/img/library/preview/empty.png" alt="" />
    <div class="title">{{ t('title') }}</div>
    <div class="info">{{ t('info') }}</div>
  </div>
</template>
<script setup>
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.document-segmentation.components.empty')
</script>
<style lang="less" scoped>
.empty-box {
  padding-top: 70px;
  width: 100%;
  height: calc(100% - 20px);
  background: #f2f4f7;
  text-align: center;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #595959;
    margin-bottom: 8px;
  }

  .info {
    color: rgba(0, 0, 0, 0.45);
  }
}
</style>