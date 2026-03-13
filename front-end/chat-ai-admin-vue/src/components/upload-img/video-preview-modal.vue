<template>
  <a-modal
    v-model:open="open"
    :title="t('title_video_preview')"
    :footer="null"
    :width="760"
    @cancel="close"
    destroyOnClose
  >
    <div class="video-preview-box">
      <video v-if="videoUrl" :src="videoUrl" controls autoplay></video>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('components.upload-img.index')
const open = ref(false)
const videoUrl = ref('')

const show = (url) => {
  videoUrl.value = url || ''
  open.value = true
}

const close = () => {
  open.value = false
  videoUrl.value = ''
}

defineExpose({ show, close })
</script>

<style lang="less" scoped>
.video-preview-box {
  width: 100%;
  background: #000;
  border-radius: 8px;
  overflow: hidden;

  video {
    width: 100%;
    max-height: 72vh;
    display: block;
    background: #000;
  }
}
</style>
