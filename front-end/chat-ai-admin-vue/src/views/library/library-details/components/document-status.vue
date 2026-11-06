<template><span class="workspace-document-status" :class="[kind, { 'on-card': card }]" :title="record.errmsg && record.errmsg !== 'success' ? record.errmsg : undefined"><img :src="kind === 'complete' ? icons.success : kind === 'error' ? icons.error : icons.clock" alt="" /><span>{{ label }}</span></span></template>
<script setup>
import { computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import icons from './workspace-document-icons'
const props = defineProps({ record: { type: Object, required: true }, card: Boolean })
const { t } = useI18n('views.library.library-details.knowledge-document')
const kind = computed(() => [3, 7, 8].includes(Number(props.record.status)) ? 'error' : Number(props.record.status) === 2 ? 'complete' : 'waiting')
const label = computed(() => {
  const status = Number(props.record.status)
  const key = { 0: 'status_converting', 1: 'status_learning', 2: 'status_learn_complete', 3: 'status_convert_failed', 4: 'status_wait_learn', 5: 'status_wait_fetch', 6: 'status_fetching', 7: 'status_fetch_failed', 8: 'status_convert_error', 9: 'status_cancel_parse', 10: 'status_segmenting' }[status]
  const total = Number(props.record.ocr_pdf_total)
  const current = Number(props.record.ocr_pdf_index)
  // 没有真实进度时只展示状态，不能用等待时长模拟学习百分比。
  const progress = props.card && [0, 1].includes(status) && total > 0 && Number.isFinite(current) ? `${Math.max(0, Math.min(100, Math.floor(current / total * 100)))}% ` : ''
  return progress + (key ? t(key) : '--')
})
</script>
<style scoped>
.workspace-document-status { display: inline-flex; align-items: center; gap: 4px; min-height: 24px; padding: 2px 6px; border-radius: 4px; font-size: 12px; line-height: 20px; white-space: nowrap; color: #5a5a5a; background: #f5f5f5; }
.workspace-document-status.complete { color: #007012; background: #f0fdf2; }
.workspace-document-status.error { color: #ee1d36; background: #fef1f3; }
.workspace-document-status.on-card { position: absolute; top: 8px; left: 8px; max-width: calc(100% - 16px); box-shadow: 0 0 0 1px #e8e8e8; }
.workspace-document-status.on-card.waiting { background: #fff; }
</style>
