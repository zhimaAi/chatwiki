<template>
  <div class="voice-list-box">
    <div class="voice-item" v-for="(item, index) in props.voice_content" :key="index">
      <VoicePlayer :url="item.voice" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import VoicePlayer from './voice-player.vue'
const props = defineProps({
  voice_content: {
    type: Array,
    default: () => []
  }
})
</script>

<style lang="less" scoped>

.voice-list-box{
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

</style>
