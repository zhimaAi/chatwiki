<template>
  <div class="mode-box">
    <a-tooltip>
      <template #title>
        <div class="tooltip-text-box" v-html="tooltipStr"></div>
      </template>
      <div class="tag-item">{{ name }}</div>
    </a-tooltip>
  </div>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-preview.components.segmentation-mode')

const props = defineProps({
  detailsInfo: {
    type: Object
  }
})
const tooltipStr = ref('')
const name = computed(() => {
  const {
    is_table_file,
    is_qa_doc,
    is_diy_split,
    question_lable,
    answer_lable,
    separators,
    chunk_size,
    chunk_overlap,
    question_column,
    answer_column,
    qa_index_type,
    doc_type
  } = props.detailsInfo
  if (is_table_file != 1 && is_qa_doc != 1 && is_diy_split != 1) {
    if (doc_type == 3) {
      tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_manual_segmentation')}\n${t('msg_doc_type')}：${t('msg_normal_doc')}`
      return `${t('msg_custom')}：${t('msg_normal_doc')}`
    }
    tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_smart_segmentation')}\n${t('msg_doc_type')}：${t('msg_normal_doc')}`
    return `${t('msg_smart_segmentation')}：${t('msg_normal_doc')}`
  }
  if (is_table_file != 1 && is_qa_doc == 1) {
    if (doc_type == 3) {
      tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_manual_segmentation')}\n${t('msg_doc_type')}：${t('msg_qa_doc')}`
      return `${t('msg_custom')}：${t('msg_qa_doc')}`
    }
    tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_smart_segmentation')}\n${t('msg_doc_type')}：${t('msg_qa_doc')}\n${t('msg_question_start_marker')}"${question_lable}"\n${t('msg_answer_start_marker')}"${answer_lable}"`
    return `${t('msg_smart_segmentation')}：${t('msg_qa_doc')}`
  }
  if (is_table_file != 1 && is_qa_doc != 1 && is_diy_split == 1) {
    tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_custom_segmentation')}\n${t('msg_doc_marker')}：${separators}\n${t('msg_chunk_max_length')}：${chunk_size}${t('msg_characters')}\n${t('msg_chunk_overlap_length')}：${chunk_overlap}${t('msg_characters')}`
    return t('msg_custom_segmentation')
  }
  if (is_table_file == 1 && is_qa_doc != 1) {
    tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_smart_segmentation')}\n${t('msg_doc_type')}：${t('msg_normal_table')}`
    return `${t('msg_smart_segmentation')}：${t('msg_normal_table')}`
  }
  if (is_table_file == 1 && is_qa_doc == 1) {
    let qa_type_str = qa_index_type == 1 ? t('msg_qa_generate_index_together') : t('msg_question_only_index')
    tooltipStr.value = `${t('msg_segmentation_mode')}：${t('msg_smart_segmentation')}\n${t('msg_doc_type')}：${t('msg_qa_doc')}\n${t('msg_question_column')}"${question_column}"\n${t('msg_answer_column')}"${answer_column}"\n${t('msg_index_method')}：${qa_type_str}`
    return `${t('msg_smart_segmentation')}：${t('msg_qa_doc')}`
  }
})
</script>

<style lang="less" scoped>
.mode-box {
  margin-left: 8px;
  cursor: context-menu;
}
.tooltip-text-box {
  white-space: pre-wrap;
}
.tag-item {
  border-radius: 2px;
  border: 1px solid #99bffd;
  background: #e9f1fe;
  padding: 0 8px;
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  color: #2475fc;
}
</style>
