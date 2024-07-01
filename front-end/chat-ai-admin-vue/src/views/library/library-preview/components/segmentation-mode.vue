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
      tooltipStr.value = `分段模式：手动分段\n文档类型：普通文档`
      return '自定义：普通文档'
    }
    tooltipStr.value = `分段模式：智能分段\n文档类型：普通文档`
    return '智能分段：普通文档'
  }
  if (is_table_file != 1 && is_qa_doc == 1) {
    if (doc_type == 3) {
      tooltipStr.value = `分段模式：手动分段\n文档类型：QA文档`
      return '自定义：QA文档'
    }
    tooltipStr.value = `分段模式：智能分段\n文档类型：QA文档\n问题开始标识符："${question_lable}"\n答案开始标识符："${answer_lable}"`
    return '智能分段：QA文档'
  }
  if (is_table_file != 1 && is_qa_doc != 1 && is_diy_split == 1) {
    tooltipStr.value = `分段模式：自定义分段\n文档标识符：${separators}\n分段最大长度：${chunk_size}字符\n文段重叠长度：${chunk_overlap}字符`
    return '自定义分段'
  }
  if (is_table_file == 1 && is_qa_doc != 1) {
    tooltipStr.value = `分段模式：智能分段\n文档类型：普通表格`
    return '智能分段：普通表格'
  }
  if (is_table_file == 1 && is_qa_doc == 1) {
    let qa_type_str = qa_index_type == 1 ? '问题与答案一起生成索引' : '仅对问题生成索引'
    tooltipStr.value = `分段模式：智能分段\n文档类型：QA文档\n问题所在列："${question_column}"\n答案所在列："${answer_column}"\n索引方式：${qa_type_str}`
    return '智能分段：QA文档'
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
