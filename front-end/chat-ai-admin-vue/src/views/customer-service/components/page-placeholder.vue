<template>
  <section class="page-placeholder">
    <h1 class="page-title">{{ t(titleKey) }}</h1>
    <div class="page-content">
      <a-empty :image="Empty.PRESENTED_IMAGE_SIMPLE" :description="t('feature_developing')" />
    </div>
  </section>
</template>

<script setup>
import { Empty } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

defineProps({
  titleKey: { type: String, required: true }
})

const { t } = useI18n('views.customer-service.index')
</script>

<style lang="less" scoped>
.page-placeholder {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  min-height: 0;
}
.page-title {
  margin: 0;
  padding: 20px;
  border-bottom: 1px solid #e4e6eb;
  background: #fff;
  font-size: 24px;
  font-weight: 500;
}
.page-content {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  min-height: 0;
  padding: 24px;
  overflow: auto;
}
</style>
