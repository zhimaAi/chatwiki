<template>
  <div class="upload-box">
    <!-- <input type="text" v-model="inputVal" @paste="pasteUpload()"> -->
    <a-input class="hidden-input" ref="inputRef" @paste="pasteUpload"></a-input>
    <div @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
      <a-upload-dragger
        v-model:fileList="fileList"
        name="file"
        :multiple="false"
        list-type="picture"
        :show-upload-list="false"
        accept=".jpg,.png,.jpeg"
        :before-upload="beforeUpload"
        @drop="handleDrop"
      >
        <div class="img-list-box" @click.stop>
          <div class="img-item" v-for="(item, index) in imageUrl" :key="index">
            <div class="mask-box">
              <EyeOutlined @click="preview(item)" />
              <DeleteOutlined @click="del(index)" />
            </div>
            <img :src="item" alt="" />
          </div>
        </div>
        <p class="upload-text" :class="{ 'center-content': imageUrl.length == 0 }">
          支持点击空白处、拖拽、粘贴图片，上传图片不得超过2M，仅支持png、jpg、jpeg格式
        </p>
      </a-upload-dragger>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { InboxOutlined, EyeOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { uploadFile } from '@/api/app'
import { api as viewerApi } from 'v-viewer'
const emit = defineEmits(['update:value'])
const fileList = ref([])
const imageUrl = ref([])
const inputVal = ref('')
const maxUploadNum = 3
const props = defineProps({
  value: {
    type: [String, Array],
    default: ''
  }
})
watch(
  () => props.value,
  (val) => {
    imageUrl.value = val
  },
  {
    immediate: true
  }
)
const beforeUpload = (file) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'

  if (!isJpgOrPng) {
    message.error('只支持JPG、PNG格式的图片')
    return false
  }

  const isLt2M = file.size / 1024 < 1024 * 2
  if (imageUrl.value.length >= maxUploadNum) {
    return message.error('最多上传三张图片')
  }
  if (!isLt2M) {
    message.error('图片大小不能超过2M')
    return false
  }
  uploadFile({
    category: 'library_image',
    file
  }).then((res) => {
    imageUrl.value.push(res.data.link)
    emit('update:value', imageUrl.value)
  })
  return false
}
const preview = (img) => {
  viewerApi({
    images: [img]
  })
}
const del = (index) => {
  imageUrl.value.splice(index, 1)
}

let isUploading = false
const pasteUpload = async (e) => {
  if (imageUrl.value.length >= maxUploadNum) {
    return message.error('最多上传三张图片')
  }
  if (!(e.clipboardData && e.clipboardData.items)) {
    message.error('当前浏览器不支持粘贴上传操作！')
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
        if (items[i].type.indexOf('image') !== -1) {
          file = items[i].getAsFile()
          break
        }
      }
      inputVal.value = ''
      if (file) {
        e.preventDefault()
        isUploading = true
        uploadFile({
          category: 'library_image',
          file
        }).then((res) => {
          isUploading = false
          imageUrl.value.push(res.data.link)
          emit('update:value', imageUrl.value)
        })
      } else {
        error('请粘贴正确图片')
        isUploading = false
        inputVal.value = ''
      }
    }
  } catch (e) {
    isUploading = false
    message.error('上传失败')
  }
}

const inputRef = ref(null)
const handleMouseEnter = ()=>{
  // 移入鼠标 将input聚焦
  inputRef.value.focus();
}
const handleMouseLeave = ()=>{
  // 移出鼠标 将input 取消焦点
  inputRef.value.blur();
}
function handleDrop(e) {
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
    }
    img {
      width: 100%;
      height: 100%;
    }
  }
}
.hidden-input{
  position: absolute;
  left: 9999;
  opacity: 0;
}
</style>