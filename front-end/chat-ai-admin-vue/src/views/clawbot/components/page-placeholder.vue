<template>
  <div class="page-placeholder">
    <div class="page-header">
      <div class="page-title">{{ title }}</div>
      <div class="page-subtitle">{{ displayDescription }}</div>
    </div>
    <div class="page-body">
      <a-empty :description="displayEmptyText" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.components.page-placeholder')

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  emptyText: {
    type: String,
    default: '',
  }
})

const displayDescription = computed(() => props.description || t('desc_default'))
const displayEmptyText = computed(() => props.emptyText || t('empty_default'))
</script>

<style lang="less" scoped>
.page-placeholder {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 28px 32px;
  background: #f7f8fb;
}

.page-header {
  margin-bottom: 24px;

  .page-title {
    font-size: 34px;
    font-weight: 700;
    color: #101828;
  }

  .page-subtitle {
    margin-top: 10px;
    font-size: 15px;
    color: #667085;
  }
}

.page-body {
  flex: 1;
  border-radius: 20px;
  border: 1px solid #e4e7ec;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
