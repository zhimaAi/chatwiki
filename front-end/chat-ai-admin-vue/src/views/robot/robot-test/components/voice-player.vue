<template>
  <div class="voice-play-container">
    <div class="voice-controls">
      <button
        class="play-btn"
        :class="{ playing: isPlaying }"
        @click="togglePlay"
        :disabled="!url || loading"
      >
        <div class="play-content" v-if="isPlaying">
          <img
            class="icon1"
            src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/voice/play-active.gif"
            alt=""
          />
          <span class="text1">{{ loading ? t('label_loading') : t('label_playing') }}...</span>
          <span class="text2">{{ t('label_click_to_pause') }} </span>
          <img class="icon2" src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/voice/open.svg" alt="">
        </div>
        <div class="play-content" v-else>
          <img
            class="icon1"
            src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/voice/play-static.svg"
            alt=""
          />
          <span class="text1">{{ t('label_voice_message') }}</span>
          <span class="text2">{{ t('label_click_to_play') }}</span>
          <img class="icon2" src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/voice/stop.svg" alt="">
        </div>
      </button>
      <div v-if="showProgress" class="progress-container">
        <input
          type="range"
          class="progress-bar"
          v-model="currentTime"
          :max="duration"
          @input="onProgressChange"
        />
        <div class="time-info">{{ formatTime(currentTime) }} / {{ formatTime(duration) }}</div>
      </div>
    </div>
    <audio
      ref="audioRef"
      :src="url"
      @loadedmetadata="onLoadedMetadata"
      @timeupdate="onTimeUpdate"
      @ended="onEnded"
      @error="onError"
    ></audio>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-test.components.voice-player')

const props = defineProps({
  url: { type: String, default: '' },
  showProgress: { type: Boolean, default: false }
})

const audioRef = ref(null)
const isPlaying = ref(false)
const duration = ref(0)
const currentTime = ref(0)
const loading = ref(false)

// 播放/暂停切换
const togglePlay = () => {
  if (!props.url) return

  const audio = audioRef.value
  if (isPlaying.value) {
    audio.pause()
    isPlaying.value = false
  } else {
    loading.value = true
    audio
      .play()
      .then(() => {
        isPlaying.value = true
        loading.value = false
      })
      .catch((error) => {
        console.error(t('msg_play_failed'), error)
        loading.value = false
      })
  }
}

// 音频元数据加载完成
const onLoadedMetadata = () => {
  duration.value = audioRef.value.duration
}

// 时间更新
const onTimeUpdate = () => {
  currentTime.value = audioRef.value.currentTime
}

// 播放结束
const onEnded = () => {
  isPlaying.value = false
  currentTime.value = 0
}

// 音频错误处理
const onError = () => {
  console.error(t('msg_audio_error'))
  isPlaying.value = false
  loading.value = false
}

// 进度条改变
const onProgressChange = () => {
  if (audioRef.value) {
    audioRef.value.currentTime = currentTime.value
  }
}

// 格式化时间显示
const formatTime = (time) => {
  if (isNaN(time)) return '00:00'
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60)
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}

// 监听URL变化，重新设置音频源
watch(
  () => props.url,
  (newUrl) => {
    if (audioRef.value) {
      // 重置状态
      isPlaying.value = false
      currentTime.value = 0
      audioRef.value.load() // 重新加载音频
    }
  }
)
</script>

<style lang="less" scoped>
.voice-play-container {
  display: flex;
  align-items: center;
  gap: 10px;

  .voice-controls {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .play-btn {
      padding: 12px;
      border: 1px solid #d9d9d9;
      border-radius: 16px;
      background: #fff;
      cursor: pointer;
      transition: all 0.3s;
      height: 46px;

      &:hover:not(:disabled) {
        border-color: #2475fc;
        background: #e6f7ff;
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
    }

    .play-content {
      display: flex;
      align-items: center;
      gap: 8px;
      .icon1 {
        width: 16px;
      }
      .text1 {
        color: #262626;
        font-size: 14px;
      }
      .text2 {
        color: #8c8c8c;
        font-size: 14px;
        margin-right: 4px;
      }
      .icon2{
        width: 16px;
      }
    }
    .progress-container {
      display: flex;
      align-items: center;
      gap: 10px;
      width: 200px;

      .progress-bar {
        flex: 1;
        height: 4px;
        -webkit-appearance: none;
        background: #f0f0f0;
        border-radius: 2px;

        &::-webkit-slider-thumb {
          -webkit-appearance: none;
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background: #1890ff;
          cursor: pointer;
        }
      }

      .time-info {
        font-size: 12px;
        color: #666;
        min-width: 50px;
        text-align: right;
      }
    }
  }
}
</style>
