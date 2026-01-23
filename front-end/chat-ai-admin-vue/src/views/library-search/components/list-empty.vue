<template>
  <div class="no-data">
    <svg-icon name="empty-document" :style="{'font-size': `${props.size}px`, 'color': 'white'}" />
    <div class="no-data-text">
      <slot name="default">{{ displayText }}</slot>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n';
import { computed } from 'vue';

const { t } = useI18n('views.library-search.components.list-empty');

const props = defineProps({
  text: String,
  size: {
    type: [Number, String],
    default: 100 
  }
});

const displayText = computed(() => props.text || t('no_data'));
</script>

<style lang="less" scoped>
.no-data {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  .no-data-text {
    margin-top: 16px;
    font-size: 14px;
    color: #999;
    text-align: center;
  }
}
</style>
