<style lang="less" scoped>
.message-list {
    height: 100%;
    overflow-y: auto;
    overflow-x: hidden;
}
</style>

<template>
    <div class="message-list" ref="scrollBoxRef" @scroll="onScroll">
        <slot></slot>
    </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'

interface ScrollOption {
  scrollTop: number;
  scrollHeight: number;
  clientHeight: number;
  scrollStartDiff: number;
  scrollEndDiff: number;
}

type ScrollDirection = 'top' | 'bottom';

const emit = defineEmits(['clickMsgMeun', 'scroll', 'scrollStart', 'scrollEnd'])

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  }
})

const scrollBoxRef = ref<HTMLDivElement | null>(null);
const scrollOption: ScrollOption  = {
  scrollTop: 0,
  scrollHeight: 0,
  clientHeight: 0,
  scrollStartDiff: 60,
  scrollEndDiff: 60
}

let scrollEventTimer: number | null = null
let onScrollEventLock = false

function handleScrollTop() {
  // 如果消息列表为空可能是断线重连等逻辑手动清空了消息列表造成的抖动，此时不触发事件
  if (props.messages.length == 0) {
    return
  }
  emit('scrollStart', {
    msg: props.messages[0]
  })
}

function handleScrollBottom() {
  // 如果消息列表为空可能是断线重连等逻辑手动清空了消息列表造成的抖动，此时不触发事件
  if (props.messages.length == 0) {
    return
  }

  emit('scrollEnd', {
    msg: props.messages[props.messages.length - 1]
  })
}

function onScroll(e: Event) {
  if (onScrollEventLock) {
    return
  }

  if (scrollEventTimer !== null) {
    clearTimeout(scrollEventTimer)
    scrollEventTimer = null
  }

  scrollEventTimer = window.setTimeout(() => {
    scrollOption.scrollTop = (e.target as HTMLDivElement).scrollTop
    scrollOption.scrollHeight = (e.target as HTMLDivElement).scrollHeight
    scrollOption.clientHeight =(e.target as HTMLDivElement).clientHeight

    emit('scroll', { ...scrollOption })

    let isAtTop = Math.abs(scrollOption.scrollTop) <= scrollOption.scrollStartDiff

    if (isAtTop) {
      handleScrollTop()
    }

    let isAtBottom =
      Math.abs(scrollOption.scrollHeight - scrollOption.scrollTop - scrollOption.clientHeight) <=
      scrollOption.scrollEndDiff

    if (isAtBottom) {
      handleScrollBottom()
    }
  }, 50)
}

const scrollToBottom = () => {
  if (!scrollBoxRef.value) {
    return
  }

  nextTick(() => {
    // 手动控制滚动到底部不触发触底事件
    onScrollEventLock = true
    if(!scrollBoxRef.value){
      return
    }

    scrollBoxRef.value.scrollTop = scrollBoxRef.value.scrollHeight + 1

    setTimeout(() => {
      onScrollEventLock = false
    }, 50)
  })
}
function scrollToMessage(id:string, direction: ScrollDirection) {
  nextTick(() => {
    // 手动控制滚动到底部不触发触底事件
    onScrollEventLock = true

    if (!direction) {
      direction = 'top'
    }

    let scroller = scrollBoxRef.value
    let element = document.querySelector<HTMLElement>('#msg-' + id)
    
    if(!scroller  || !element){
      return
    }

    if (direction == 'top') {
      scroller.scrollTop = element.offsetTop
    } else {
      scroller.scrollTop = element.offsetTop - scroller.clientHeight + element.clientHeight
    }

    setTimeout(() => {
      onScrollEventLock = false
    }, 50)
  })
}

defineExpose({
  scrollToBottom,
  scrollToMessage
})
</script>
