<template>
  <div class="chat-message-scroll" ref="scrollContainer" @scroll="onScroll">
    <slot></slot>
  </div>
</template>

<script setup>
/**
 * 聊天消息滚动容器组件
 * 用于处理消息列表的滚动事件，包括触顶和触底检测
 * 支持自定义触发距离，可用于实现上拉加载更多、下拉刷新等功能
 */

import { ref, watch } from 'vue'

// 组件属性定义
const props = defineProps({
  // 触发顶部滚动事件的距离，默认40px
  topTriggerDistance: {
    type: Number,
    default: 10
  },
  // 触发底部滚动事件的距离，默认40px
  bottomTriggerDistance: {
    type: Number,
    default: 40
  },
  // 是否正在加载数据，加载数据时不触发触顶和触底事件
  isLoading: {
    type: Boolean,
    default: false
  },
  // 控制滚动位置
  scrollTop: {
    type: Number,
    default: 0
  }
})

// 定义组件事件
const emit = defineEmits([
  'scroll-top', // 滚动到顶部时触发
  'scroll-bottom', // 滚动到底部时触发
  'scroll' // 滚动时触发
])

// 滚动容器引用
const scrollContainer = ref(null)

// 滚动状态配置
const scrollOption = ref({
  scrollTop: 0, // 当前滚动位置
  scrollHeight: 0, // 滚动内容总高度
  clientHeight: 0, // 可视区域高度
  scrollDirection: 'down', // 滚动方向：up-向上滚动，down-向下滚动
  scrollStartDiff: props.topTriggerDistance, // 触发顶部事件的距离
  scrollEndDiff: props.bottomTriggerDistance // 触发底部事件的距离
})

// 滚动事件锁，防止重复触发
let isScrollLocked = false

// 使用Vue的ref来存储防抖定时器
const scrollDebounceTimer = ref(null)

/**
 * 更新滚动状态
 * @param {HTMLElement} target - 滚动容器元素
 */
const updateScrollState = (target) => {
  const { scrollTop: currentScrollTop } = target
  const scrollDiff = scrollOption.value.scrollTop - currentScrollTop

  // 更新滚动方向
  scrollOption.value.scrollDirection = scrollDiff > 0 ? 'up' : 'down'

  // 更新滚动状态
  Object.assign(scrollOption.value, {
    scrollTop: currentScrollTop,
    scrollHeight: target.scrollHeight,
    clientHeight: target.clientHeight
  })
}

/**
 * 检查是否触发滚动边界事件
 */
const checkScrollBoundaries = () => {
  const { scrollTop, scrollHeight, clientHeight, scrollStartDiff, scrollEndDiff, scrollDirection } =
    scrollOption.value

  // 触发通用滚动事件
  emit('scroll', { ...scrollOption.value })

  if (!props.isLoading) {
    // 检测是否触顶
    const isAtTop = Math.abs(scrollTop) <= scrollStartDiff
    if (isAtTop && scrollDirection === 'up') {
      emit('scroll-top', { ...scrollOption.value })
    }

    // 检测是否触底
    const isAtBottom = Math.abs(scrollHeight - scrollTop - clientHeight) <= scrollEndDiff
    if (isAtBottom && scrollDirection === 'down') {
      emit('scroll-bottom', { ...scrollOption.value })
    }
  }
}

/**
 * 监听滚动事件
 * 处理滚动方向判断、触顶和触底检测
 * @param {Event} e - 滚动事件对象
 */
const onScroll = (e) => {
  updateScrollState(e.target)

  if (isScrollLocked) return

  if (scrollDebounceTimer.value) {
    clearTimeout(scrollDebounceTimer.value)
    scrollDebounceTimer.value = null
  }

  scrollDebounceTimer.value = setTimeout(() => {
    checkScrollBoundaries()
  }, 100)
}

/**
 * 滚动到顶部
 */
function scrollToTop() {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = 0
  }
}

/**
 * 滚动到底部
 */
function scrollToBottom() {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight
  }
}

function getState() {
  if (scrollContainer.value) {
    Object.assign(scrollOption.value, {
      scrollTop: scrollContainer.value.scrollTop,
      scrollHeight: scrollContainer.value.scrollHeight,
      clientHeight: scrollContainer.value.clientHeight
    })
  }

  return JSON.parse(JSON.stringify(scrollOption.value))
}

// 监听scrollTop属性变化
watch(
  () => props.scrollTop,
  (newVal) => {
    if (scrollContainer.value) {
      scrollContainer.value.scrollTop = newVal
    }
  }
)

// 暴露方法给父组件
defineExpose({
  scrollToTop,
  scrollToBottom,
  getState
})
</script>

<style scoped>
.chat-message-scroll {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
