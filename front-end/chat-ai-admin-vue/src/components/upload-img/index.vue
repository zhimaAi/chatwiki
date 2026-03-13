<template>
  <div class="upload-box">
    <VideoPreviewModal ref="videoPreviewModalRef" />
    <a-input class="hidden-input" ref="inputRef" @paste="pasteUpload"></a-input>
    <div @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
      <a-upload-dragger
        v-model:fileList="fileList"
        name="file"
        :multiple="false"
        list-type="picture"
        :show-upload-list="false"
        :accept="accept"
        :before-upload="beforeUpload"
        @drop="handleDrop"
      >
        <div class="img-list-box" @click.stop>
          <div class="img-item" v-for="(item, index) in mediaList" :key="index">
            <div class="mask-box">
              <EyeOutlined @click.stop="preview(item)" />
              <DeleteOutlined @click.stop="del(index)" />
            </div>
            <video v-if="isVideoUrl(item)" :src="item" muted preload="metadata"></video>
            <img v-else :src="item" alt="" />
          </div>
        </div>
        <p class="upload-text" :class="{ 'center-content': mediaList.length == 0 }">
          {{ hintText }}
        </p>
      </a-upload-dragger>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { EyeOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { uploadFile } from '@/api/app'
import { api as viewerApi } from 'v-viewer'
import { useI18n } from '@/hooks/web/useI18n'
import VideoPreviewModal from './video-preview-modal.vue'

const { t } = useI18n('components.upload-img.index')
const emit = defineEmits(['update:value'])
const fileList = ref([])
const mediaList = ref([])
const videoPreviewModalRef = ref(null)
const props = defineProps({
  value: {
    type: [String, Array],
    default: ''
  },
  accept: {
    type: String,
    default: '.jpg,.png,.jpeg'
  },
  maxUploadNum: {
    type: Number,
    default: 3
  },
  allowVideo: {
    type: Boolean,
    default: false
  },
  imageMaxSizeMb: {
    type: Number,
    default: 2
  },
  videoMaxSizeMb: {
    type: Number,
    default: 20
  },
  uploadCategory: {
    type: String,
    default: 'library_image'
  },
  uploadHint: {
    type: String,
    default: ''
  }
})

const normalizeToArray = (val) => {
  if (Array.isArray(val)) {
    return val.filter(Boolean)
  }
  if (typeof val === 'string' && val) {
    return [val]
  }
  return []
}

watch(
  () => props.value,
  (val) => {
    mediaList.value = normalizeToArray(val)
  },
  {
    immediate: true
  }
)

const hintText = computed(() => {
  if (props.uploadHint) return props.uploadHint
  if (props.allowVideo) {
    return t('msg_upload_hint_with_video')
  }
  return t('msg_upload_hint')
})

const allowedExts = computed(() => {
  return props.accept
    .split(',')
    .map((ext) => ext.trim().toLowerCase().replace('.', ''))
    .filter(Boolean)
})
const extMimeMap = {
  jpg: 'image/jpeg',
  jpeg: 'image/jpeg',
  png: 'image/png',
  mp4: 'video/mp4'
}
const allowedMimes = computed(() => {
  return allowedExts.value.map((ext) => extMimeMap[ext]).filter(Boolean)
})

const getFileExt = (file) => {
  const name = file?.name || ''
  const parts = name.split('.')
  if (parts.length <= 1) return ''
  return parts.pop().toLowerCase()
}

const isVideoFile = (file) => {
  const ext = getFileExt(file)
  return (file?.type || '').startsWith('video/') || ext === 'mp4'
}

const validateFile = (file) => {
  const ext = getFileExt(file)
  const mimeType = (file?.type || '').toLowerCase()
  const isAccepted = allowedExts.value.includes(ext) || allowedMimes.value.includes(mimeType)
  if (!isAccepted) {
    message.error(props.allowVideo ? t('msg_invalid_media_format') : t('msg_invalid_format'))
    return false
  }

  if (mediaList.value.length >= props.maxUploadNum) {
    message.error(props.allowVideo ? t('msg_max_file_limit', { count: props.maxUploadNum }) : t('msg_max_limit'))
    return false
  }

  const maxSizeMb = isVideoFile(file) ? props.videoMaxSizeMb : props.imageMaxSizeMb
  const isValidSize = file.size / 1024 / 1024 <= maxSizeMb
  if (!isValidSize) {
    message.error(
      isVideoFile(file)
        ? t('msg_video_size_limit', { size: props.videoMaxSizeMb })
        : t('msg_image_size_limit', { size: props.imageMaxSizeMb })
    )
    return false
  }
  return true
}

const getUploadCategory = (file) => {
  if (isVideoFile(file)) {
    return 'library_video'
  }
  const mimeType = (file?.type || '').toLowerCase()
  if (mimeType.startsWith('image/')) {
    return 'library_image'
  }
  return props.uploadCategory
}

const uploadOne = (file) => {
  uploadFile({
    category: getUploadCategory(file),
    file
  }).then((res) => {
    mediaList.value.push(res.data.link)
    emit('update:value', [...mediaList.value])
  })
}

const beforeUpload = (file) => {
  if (!validateFile(file)) {
    return false
  }
  uploadOne(file)
  return false
}

const preview = (img) => {
  if (isVideoUrl(img)) {
    videoPreviewModalRef.value?.show(img)
    return
  }
  viewerApi({
    images: [img]
  })
}
const del = (index) => {
  mediaList.value.splice(index, 1)
  emit('update:value', [...mediaList.value])
}

let isUploading = false
const pasteUpload = async (e) => {
  if (mediaList.value.length >= props.maxUploadNum) {
    message.error(props.allowVideo ? t('msg_max_file_limit', { count: props.maxUploadNum }) : t('msg_max_limit'))
    return
  }
  if (!(e.clipboardData && e.clipboardData.items)) {
    message.error(t('msg_paste_not_supported'))
    return
  }
  if (isUploading) {
    return
  }
  try {
    if (e.clipboardData || e.originalEvent || window.clipboardData) {
      var clipboardData = e.clipboardData || e.originalEvent.clipboardData || window.clipboardData
      var items = clipboardData.items
      let file = null
      // 搜索剪切板items
      for (let i = 0; i < items.length; i++) {
        const itemType = items[i].type || ''
        const isImage = itemType.startsWith('image/')
        const isVideo = props.allowVideo && itemType === 'video/mp4'
        if (isImage || isVideo) {
          file = items[i].getAsFile()
          break
        }
      }
      if (file) {
        if (!validateFile(file)) {
          return
        }
        e.preventDefault()
        isUploading = true
        uploadFile({ category: getUploadCategory(file), file }).then((res) => {
          isUploading = false
          mediaList.value.push(res.data.link)
          emit('update:value', [...mediaList.value])
        })
      } else {
        message.error(props.allowVideo ? t('msg_paste_invalid_media') : t('msg_paste_invalid_image'))
        isUploading = false
      }
    }
  } catch (e) {
    isUploading = false
    message.error(t('msg_upload_failed'))
  }
}

const inputRef = ref(null)
const handleMouseEnter = ()=>{
  inputRef.value?.focus()
}
const handleMouseLeave = ()=>{
  inputRef.value?.blur()
}
function handleDrop() {
}

const isVideoUrl = (url) => {
  if (!url || typeof url !== 'string') return false
  return /\.mp4(\?|#|$)/i.test(url) || /library_video/i.test(url) || /\/video\//i.test(url)
}
</script>

<style lang="less" scoped>
.upload-box {
  &::v-deep(.ant-upload) {
    padding: 0;
  }
  &::v-deep(.ant-upload.ant-upload-btn) {
    padding: 16px;
    min-height: 86px;
  }
  &::v-deep(.ant-upload.ant-upload-drag) {
    border: none;
    background: #f2f4f7;
    &:hover {
      background: #e6efff;
    }
  }
  .upload-text {
    margin-bottom: 0;
    text-align: left;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    color: #3a4559;
    &.center-content {
      margin-top: 16px;
      text-align: center;
    }
  }
}
.img-list-box {
  display: flex;
  flex-wrap: wrap;
  margin-bottom: 8px;
  gap: 8px;
  width: fit-content;
  .img-item {
    width: 104px;
    height: 104px;
    border-radius: 2px;
    padding: 9px;
    border: 1px solid #d9d9d9;
    position: relative;
    &:hover {
      .mask-box {
        opacity: 0.75;
        span {
          cursor: pointer;
        }
      }
    }
    .mask-box {
      position: absolute;
      border-radius: 4px;
      left: 9px;
      right: 9px;
      top: 9px;
      bottom: 9px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #000;
      color: #fff;
      font-weight: 900;
      font-size: 16px;
      gap: 8px;
      opacity: 0;
      transition: all 0.3s;
      z-index: 2;
    }
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    video {
      width: 100%;
      height: 100%;
      object-fit: cover;
      background: #000;
    }
  }
}
.hidden-input{
  position: absolute;
  left: 9999;
  opacity: 0;
}
</style>
