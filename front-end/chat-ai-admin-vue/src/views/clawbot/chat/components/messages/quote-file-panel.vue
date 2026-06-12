<template>
  <div class="file-items" v-if="item.show_quote_file && item.quote_file && item.quote_file.length > 0">
    <div class="file-item" v-for="file in item.quote_file" :key="file.id">
      <a class="file-name">
        <svg-icon class="think-icon" name="quote-file"></svg-icon>
        <span v-if="file.file_name">{{ file.file_name }}</span>
        <span v-else>{{ file.library_name }}-{{ t('library_featured') }}</span>
      </a>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.chat.components.message-list')

defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
})
</script>

<style lang="less" scoped>
.file-items {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 24px;
  line-height: 22px;
  font-weight: 400;
  color: #164799;
  font-size: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid #edeff2;
  margin-bottom: 12px;

  .file-name {
    cursor: pointer;
    display: flex;
    align-items: baseline;
    gap: 4px;
    color: #164799;
    position: relative;

    .svg-action {
      position: relative;
      top: 2px;
    }
  }
}
</style>
