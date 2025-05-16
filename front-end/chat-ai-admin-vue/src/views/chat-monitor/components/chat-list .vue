<style lang="less" scoped>
.scroll-wrapper {
  position: relative;
  height: 100%;
  overflow: hidden;
}
.chat-list {
  padding: 0 16px 16px 16px;
  .loading-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px 0;
    .loading-text {
      margin-left: 8px;
      font-size: 14px;
      color: #8c8c8c;
    }
  }
  .finished-text {
    text-align: center;
    padding: 16px 0;
    font-size: 14px;
    color: #8c8c8c;
  }
  .chat-item {
    display: flex;
    padding: 12px 16px;
    border-radius: 6px;
    cursor: pointer;
    background-color: #fff;
    transition: background-color 0.2s ease;

    &:hover {
      background-color: rgba(242, 244, 247, 1);
    }

    &.active {
      background-color: rgba(230, 239, 255, 1);
    }

    .avatar {
      position: relative;
      width: 48px;
      height: 48px;
      margin-right: 12px;

      img {
        width: 100%;
        height: 100%;
        border-radius: 12px;
        object-fit: cover;
      }

      .dot {
        display: flex;
        align-items: center;
        justify-content: center;
        position: absolute;
        top: -8px;
        right: -8px;
        width: 16px;
        height: 16px;
        font-size: 12px;
        border-radius: 50%;
        font-weight: 400;
        color: #ffffff;
        background-color: #fb363f;
        &.plus {
          width: auto;
          padding: 0 4px;
          border-radius: 8px;
        }
      }
    }

    .chat-info {
      flex: 1;
      overflow: hidden;

      .nickname {
        line-height: 22px;
        font-size: 14px;
        color: rgb(38, 38, 38);
      }

      .last-message {
        margin-top: 2px;
        line-height: 22px;
        font-size: 14px;
        color: rgb(89, 89, 89);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .webapp-source {
        line-height: 20px;
        margin-top: 2px;
        font-size: 12px;
        color: rgb(140, 140, 140);
      }
    }
  }
}
</style>
<template>
  <div class="scroll-wrapper" ref="scrollRef">
    <list-empty size="140" text="暂无机器人接待中会话" v-if="receiverList.length == 0" />
    <div class="chat-list" v-else>
      <div
        v-for="item in receiverList"
        :key="item.id"
        class="chat-item"
        :class="{ active: item.id === selectedReceiverId }"
        @click="changeReceiver(item)"
      >
        <div class="avatar">
          <img :src="item.avatar" :alt="item.displayName" />
          <span class="dot" :class="{ plus: item.unread > 9 }" v-if="item.unread > 0">{{ item.unread }}</span>
        </div>
        <div class="chat-info">
          <div class="nickname">{{ item.displayName }}</div>
          <div class="last-message">{{ item.last_chat_message }}</div>
          <div class="webapp-source">来自：{{ item.come_from.app_name }}</div>
        </div>
      </div>
      <div v-if="loading" class="loading-wrapper">
        <a-spin size="small" />
        <span class="loading-text">加载中...</span>
      </div>
      <div v-if="finished" class="finished-text">没有更多了</div>
    </div>
  </div>
</template>
<script setup>
import BScroll from '@better-scroll/core'
import ScrollBar from '@better-scroll/scroll-bar'
import MouseWheel from '@better-scroll/mouse-wheel'
import ObserveDOM from '@better-scroll/observe-dom'
import Pullup from '@better-scroll/pull-up'
import { storeToRefs } from 'pinia'
import { useChatMonitorStore } from '@/stores/modules/chat-monitor.js'
import { ref, onMounted, nextTick, onUnmounted, watch } from 'vue'
import ListEmpty from './list-empty.vue'

BScroll.use(MouseWheel)
BScroll.use(ScrollBar)
BScroll.use(ObserveDOM)
BScroll.use(Pullup)

const emit = defineEmits(['switchChat'])

const chatMonitorStore = useChatMonitorStore()
const { getReceiverList } = chatMonitorStore
const { receiverList, selectedReceiverId } = storeToRefs(chatMonitorStore)

const scrollRef = ref(null)
let scrollController = null
const loading = ref(false)
const finished = ref(false)

const initScroll = () => {
  scrollController = new BScroll(scrollRef.value, {
    scrollY: true,
    click: true,
    observeDOM: true,
    mouseWheel: true,
    bounce: false,
    scrollbar: {
      fade: true,
      interactive: true
    },
    pullUpLoad: true
  })

  scrollController.on('pullingUp', () => {
    getData()
  })
}

const getData = async (params) => {
  if (loading.value || finished.value) {
    scrollController.finishPullUp()
    return
  }

  loading.value = true

  try {
    let res = await getReceiverList(params)

    if (res.data.list.length === 0 && !params.keyword) {
      finished.value = true
    }
  } catch (error) {
    console.error('加载更多数据失败:', error)
  } finally {
    loading.value = false
    scrollController.finishPullUp()
  }
}

const changeReceiver = (item) => {
  emit('switchChat', item)
}

watch(receiverList, (newValue) => {
  if (newValue.length > 0) {
    nextTick(() => {
      scrollController && scrollController.refresh()
    })
  }
})

onMounted(() => {
  nextTick(() => {
    initScroll()
  })
})

onUnmounted(() => {
  if (scrollController) {
    scrollController.destroy()
  }
})

defineExpose({
  getData
})
</script>
