<style lang="less" scoped>
.multiple-message {
  .message-items-wrapper{
    margin-bottom: 8px;
    &:last-child {
      margin-bottom: 0;
    }
  }
  .text-message {
    font-size: 14px;
    line-height: 22px;
    color: #3a4559;
  }
  .file-message-wrapper {
    display: flex;
    flex-flow: row wrap;
    gap: 8px;
  }
  .image-message {
    width: 100px;
    height: 100px;
    border-radius: 4.5px;
    border: 0.75px dashed #D9D9D9;
    overflow: hidden;
    cursor: pointer;
    .img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
}
</style>

<template>
  <div class="multiple-message">
    <div class="message-items-wrapper text-message-wrapper" v-if="textMessages.length > 0">
      <template  v-for="(message, index) in textMessages" :key="index">
        <div class="text-message" v-if="message.type === 'text'">{{ message.text }}</div>
      </template>
    </div>
    <div class="message-items-wrapper file-message-wrapper" v-if="imageMessages.length > 0">
      <div class="image-message" v-for="(message, index) in imageMessages" :key="index" @click="handlePreview(index)">
        <img class="img" :src="message.image_url.url" alt="" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { api as viewerApi } from "v-viewer"

const props = defineProps({
  message: {
    type: String,
    default: ''
  }
})

const messages = computed(() => JSON.parse(props.message) || [])
const textMessages = computed(() => messages.value.filter(i => i.type === 'text'))
const imageMessages = computed(() => messages.value.filter(i => i.type === 'image_url'))
const images = computed(() => imageMessages.value.map(i => i.image_url.url))

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
</script>
