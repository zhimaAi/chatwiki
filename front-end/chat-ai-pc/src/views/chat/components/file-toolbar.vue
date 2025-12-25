<template>
  <div class="file-toolbar">
    <div class="toolbar-left">
      <div class="file-list">
        <div
          v-for="(file, index) in fileList"
          :key="index"
          class="file-item"
          :class="{
            'uploading': file.status === 'uploading',
            'error': file.status === 'error'
          }"
        >
          <span class="delete-btn action-btn" v-if="file.status == 'done' || file.status == 'error'">
            <van-icon class="delete-icon" @click="handleDelete(index)" name="cross" />
          </span>
          <span class="preview-btn action-btn" v-if="file.status == 'done'">
            <svg-icon class="preview-icon" name="eye-open" @click="handlePreview(index)"></svg-icon>
          </span>

          <img :src="file.url" alt="avatar" />
          <div class="mask"></div>
          <div v-if="file.status === 'uploading'" class="progress-bar">
            <div class="progress" :style="{ width: file.percent + '%' }"></div>
          </div>
          <div v-if="file.status === 'error'" class="error-text">上传失败</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { api as viewerApi } from "v-viewer"
import { computed } from 'vue'

const props = defineProps({
  fileList: {
    type: Array,
    default: () => []
  }
})

// 事件处理
const emit = defineEmits(['add', 'delete'])

const images = computed(() => props.fileList.map(i => i.url))

const handlePreview = (index) => {
  const viewer = viewerApi({
    options: {
      title: false,
      toolbar: true,
      initialViewIndex: index
    },
    index: index,
    images: images.value,
  })
  viewer.show()
}

const handleDelete = (index) => {
  emit('delete', index)
}
</script>

<style lang="less" scoped>
.file-toolbar {
  display: flex;
  align-items: center;
  background-color: #fff;
  margin-bottom: 10px;

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.file-list {
  display: flex;
  align-items: center;
  flex-flow: row wrap;
  gap: 8px;
}

.file-item {
  position: relative;
  flex-shrink: 0;
  width: 62px;
  height: 62px;
  border-radius: 5px;
  cursor: pointer;
  border: 1px dashed #d9d9d9;

  .action-btn{
    display: none;
  }

  &:hover .action-btn{
    display: block;
  }
  .delete-btn{
    position: absolute;
    top: -5px;
    right: -5px;
    z-index: 10;
  }
  .delete-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    color: white;
    border-radius: 50%;
    font-size: 14px;
    background: #1E2226;
  }

  .preview-btn{
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 10;
    font-size: 16px;
    color: #fff;

    .preview-icon {
      font-size: 20px;
    }
  }

  img {
    width: 100%;
    height: 100%;
    border-radius: 5px;
    object-fit: cover;
  }

  .mask {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 6px;
    opacity: 0;
    transition: opacity 0.3s;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  &:hover .mask {
    opacity: 1;
  }

  &.uploading {
    img {
      filter: blur(2px);
    }
    .progress-bar {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 10;
      width: 90%;
      height: 4px;
      background-color: rgba(0, 0, 0, 0.4);
      border-radius: 2px;
      .progress {
        height: 100%;
        background-color: #2475FC;
        border-radius: 2px;
        transition: width 0.3s;
      }
    }
  }

  &.error {
    .mask{
      display: block;
      opacity: 1;
    }
    .error-text {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      color: white;
      font-size: 12px;
      text-align: center;
      white-space: nowrap;
    }
  }

  &.add-item {
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px dashed #d9d9d9;
    background-color: #fafafa;
    font-size: 18px;
    color: #666;
    &:hover {
      border-color: #1890ff;
      color: #1890ff;
    }
  }
}
</style>