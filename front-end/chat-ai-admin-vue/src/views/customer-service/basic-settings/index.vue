<template>
  <section class="basic-page">
    <CustomerPageHeader :title="t('basic_settings')"><a-tag v-if="context.saving.value" color="processing">{{ t('saving') }}</a-tag></CustomerPageHeader>
    <div class="basic-columns" :class="{ resizing: isResizing }">
      <div class="settings"><BasicConfig variant="workspace" /></div>
      <div id="basic-settings-preview" class="preview" :style="{ width: `${previewWidth}px` }">
        <div class="preview-divider" role="separator" tabindex="0" aria-orientation="vertical" aria-controls="basic-settings-preview" :aria-label="t('resize_preview')" :aria-valuemin="MIN_PREVIEW_WIDTH" :aria-valuemax="MAX_PREVIEW_WIDTH" :aria-valuenow="previewWidth" @pointerdown="startResize" @keydown="resizeWithKeyboard" @dblclick.prevent="previewWidth = DEFAULT_PREVIEW_WIDTH" />
        <ChatPreview />
      </div>
    </div>
  </section>
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import BasicConfig from '@/views/robot/robot-config/basic-config/index.vue'
import ChatPreview from './chat-preview.vue'
import CustomerPageHeader from '../components/page-header.vue'
const { t } = useI18n('views.customer-service.index')
const context = inject('customerRobotContext')
const MIN_PREVIEW_WIDTH = 360
const MAX_PREVIEW_WIDTH = 500
const DEFAULT_PREVIEW_WIDTH = 460
const PREVIEW_WIDTH_KEY = 'customer_service_basic_preview_width'
const clampWidth = width => Math.max(MIN_PREVIEW_WIDTH, Math.min(MAX_PREVIEW_WIDTH, width))
function readPreviewWidth() {
  try {
    const saved = localStorage.getItem(PREVIEW_WIDTH_KEY)
    const width = Number(saved)
    return saved !== null && saved.trim() !== '' && Number.isFinite(width) ? clampWidth(width) : DEFAULT_PREVIEW_WIDTH
  } catch { return DEFAULT_PREVIEW_WIDTH }
}
const previewWidth = ref(readPreviewWidth())
const isResizing = ref(false)
watch(previewWidth, width => {
  try { localStorage.setItem(PREVIEW_WIDTH_KEY, String(width)) }
  catch { /* 本地存储被禁用时仍允许调整布局。 */ }
}, { flush: 'sync' })
let stopResize
function startResize(event) {
  if (event.button !== 0 || event.isPrimary === false) return
  event.preventDefault()
  stopResize?.()
  const target = event.currentTarget
  const startX = event.clientX
  const startWidth = previewWidth.value
  const cursor = document.body.style.cursor
  const userSelect = document.body.style.userSelect
  const move = pointer => {
    if (pointer.pointerId === event.pointerId) previewWidth.value = clampWidth(startWidth + startX - pointer.clientX)
  }
  const stop = pointer => {
    if (pointer.pointerId === undefined || pointer.pointerId === event.pointerId) stopResize?.()
  }
  stopResize = () => {
    document.removeEventListener('pointermove', move)
    document.removeEventListener('pointerup', stop)
    document.removeEventListener('pointercancel', stop)
    window.removeEventListener('blur', stop)
    target.removeEventListener('lostpointercapture', stop)
    if (target.hasPointerCapture?.(event.pointerId)) target.releasePointerCapture(event.pointerId)
    document.body.style.cursor = cursor
    document.body.style.userSelect = userSelect
    isResizing.value = false
    stopResize = undefined
  }
  isResizing.value = true
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
  document.addEventListener('pointermove', move)
  document.addEventListener('pointerup', stop)
  document.addEventListener('pointercancel', stop)
  window.addEventListener('blur', stop)
  target.addEventListener('lostpointercapture', stop)
  target.setPointerCapture?.(event.pointerId)
}
function resizeWithKeyboard(event) {
  if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight') return
  event.preventDefault()
  previewWidth.value = clampWidth(previewWidth.value + (event.key === 'ArrowLeft' ? 8 : -8))
}
onBeforeUnmount(() => stopResize?.())
</script>
<style lang="less" scoped>
.basic-page { display: flex; flex: 1; min-width: 0; min-height: 0; margin: 8px 8px 8px 0; flex-direction: column; overflow: hidden; background: #fff; border-radius: 12px; }
.basic-columns { display: flex; flex: 1; min-height: 0; }
.settings { flex: 1; min-width: 0; min-height: 0; overflow: hidden; }
.preview { position: relative; flex-shrink: 0; min-width: 360px; max-width: 460px; min-height: 0; border-left: 1px solid #e8e8e8; }
.preview-divider { position: absolute; top: 0; bottom: 0; left: -4px; z-index: 1; width: 8px; cursor: col-resize; touch-action: none; outline: none; }
.preview-divider::after { position: absolute; top: 0; bottom: 0; left: 3px; width: 1px; content: ''; background: transparent; }
.preview-divider:hover::after, .preview-divider:focus-visible::after, .resizing .preview-divider::after { background: var(--app-primary, #146ef5); }
</style>
