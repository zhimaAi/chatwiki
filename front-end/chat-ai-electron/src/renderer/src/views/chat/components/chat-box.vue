<style lang="less" scoped>
.chat-box {
  display: flex;
  flex-flow: column nowrap;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: #f0f2f5;
  .robot-name {
    line-height: 22px;
    padding: 8px 0;
    font-size: 14px;
    font-weight: 600;
    color: #242933;
    text-align: center;
  }
  .iframe-box {
    flex: 1;
    width: 100%;
    .chat-iframe {
      width: 100%;
      height: 100%;
      border: none;
    }
  }
}
</style>

<template>
  <div class="chat-box">
    <div class="robot-name">{{ props.name }}</div>
    <div class="iframe-box">
      <iframe
        v-show="show"
        class="chat-iframe"
        :src="props.src"
        frameborder="0"
        @load="onLoad"
      ></iframe>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'

const props = defineProps({
  name: {
    type: String,
    defaul: ''
  },
  src: {
    type: String,
    defaul: ''
  }
})

const show = ref(false)

watch(
  () => props.src,
  () => {
    show.value = false
  }
)

const onLoad = () => {
  setTimeout(() => {
    show.value = true
  }, 300)
}

onMounted(() => {})
</script>
