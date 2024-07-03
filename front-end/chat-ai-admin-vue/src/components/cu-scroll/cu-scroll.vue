<style lang="less" scoped>
.scroll-wrapper {
  position: relative;
  height: 100%;
  overflow: hidden;

  .scroll-content {}
}
</style>
<template>
  <div class="scroll-wrapper" ref="scroller">
    <div class="scroll-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import BScroll from '@better-scroll/core'
import ScrollBar from '@better-scroll/scroll-bar'
import MouseWheel from '@better-scroll/mouse-wheel'
import ObserveDOM from '@better-scroll/observe-dom'
import Pullup from '@better-scroll/pull-up'

BScroll.use(MouseWheel)
BScroll.use(ScrollBar)
BScroll.use(ObserveDOM)
BScroll.use(Pullup)

const emit = defineEmits(['onScrollEnd'])

const props = defineProps({
  scrollX: {
    type: Boolean,
    default: false
  },
  scrollY: {
    type: Boolean,
    default: true
  },
  scrollbar: {
    type: [Boolean, Object],
    default: true
  },
  pullUpLoad: {
    // 布尔和对象
    type: [Boolean, Object],
    default: true
  }
})

const scroller = ref(null)
let scrollController = null

const refresh = () => {
  scrollController.refresh()
}

onMounted(() => {
  scrollController = new BScroll(scroller.value, {
    scrollX: props.scrollX,
    scrollY: props.scrollY,
    scrollbar: props.scrollbar,
    observeDOM: true,
    pullUpLoad: props.pullUpLoad,
    // 不派发 scroll 事件
    probeType: 0,
    mouseWheel: {
      speed: 20,
      invert: false,
      easeTime: 300
    },
    preventDefault: false,
    bounce: false,
    click: true
  })

  scrollController.on('pullingUp', () => {
    emit('onScrollEnd')

    nextTick(() => {
      scrollController.finishPullUp()
    })
  })
})

defineExpose({
  refresh
})
</script>
