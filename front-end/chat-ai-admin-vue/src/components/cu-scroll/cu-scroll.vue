<style lang="less" scoped>
.scroll-wrapper {
  position: relative;
  height: 100%;
  overflow: hidden;

  .scroll-content {
    min-height: 100%;
  }
}
.cu-scrolbar-box {
  /deep/ .bscroll-vertical-scrollbar {
    transition: opacity 0.4s;
    opacity: 0;
  }
  &:hover {
    /deep/ .bscroll-vertical-scrollbar {
      opacity: 1;
    }
  }
}
</style>
<template>
  <div
    class="scroll-wrapper"
    :class="{ 'cu-scrolbar-box': props.scrollbar.interactive }"
    ref="scroller"
  >
    <div class="scroll-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import BScroll from '@better-scroll/core'
import ScrollBar from '@better-scroll/scroll-bar'
import MouseWheel from '@better-scroll/mouse-wheel'
import ObserveDOM from '@better-scroll/observe-dom'
import Pullup from '@better-scroll/pull-up'

BScroll.use(MouseWheel)
BScroll.use(ScrollBar)
BScroll.use(ObserveDOM)
BScroll.use(Pullup)

const emit = defineEmits(['onScrollEnd', 'scroll'])

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
    default: {
      minSize: 0
    }
  },
  pullUpLoad: {
    // 布尔和对象
    type: [Boolean, Object],
    default: true
  },
  disableMouse: {
    // 禁止数据按下拖动
    type: Boolean,
    default: false
  }
})

let mouseDownHandler = null // 保存 mousedown 监听器的引用

const scroller = ref(null)
let scrollController = null

const refresh = () => {
  scrollController.refresh()
}

const scrollToElement = (option) => {
  let { el, time, offsetX, offsetY, easing } = option
  /**
   * scrollToElement(el, time, offsetX, offsetY, easing)
   * {DOM | string} el 滚动到的目标元素, 如果是字符串，则内部会尝试调用 querySelector 转换成 DOM 对象。
   * {number} time 滚动动画执行的时长（单位 ms）
   * {number | boolean} offsetX 相对于目标元素的横轴偏移量，如果设置为 true，则滚到目标元素的中心位置
   * {number | boolean} offsetY 相对于目标元素的纵轴偏移量，如果设置为 true，则滚到目标元素的中心位置
   * {Object} easing 缓动函数，一般不建议修改，如果想修改，参考源码中的 packages/shared-utils/src/ease.ts 里的写法
   */
  scrollController.scrollToElement(el, time, offsetX, offsetY, easing)
}

const enable = () => {
  scrollController.enable()
}

const disable = () => {
  scrollController.disable()
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

  scrollController.on('scroll', (position)=>{
    emit('scroll', position)
  })

  // 定义 mousedown 事件处理函数
  mouseDownHandler = (e) => {
    const startY = e.clientY

    const mouseMoveHandler = (e) => {
      const deltaY = e.clientY - startY
      if (Math.abs(deltaY) > 5) {
        e.preventDefault()
        e.stopPropagation()
      }
    }

    const removeListeners = () => {
      document.removeEventListener('mousemove', mouseMoveHandler)
      document.removeEventListener('mouseup', removeListeners)
    }

    document.addEventListener('mousemove', mouseMoveHandler)
    document.addEventListener('mouseup', removeListeners, { once: true })
  }

  // 添加 mousedown 监听
  if (props.disableMouse) {
    scroller.value?.addEventListener('mousedown', mouseDownHandler)
  }
})

onUnmounted(() => {
  // 移除 mousedown 监听（确保 scroller.value 存在）
  if (scroller.value && mouseDownHandler) {
    scroller.value.removeEventListener('mousedown', mouseDownHandler)
  }
})

defineExpose({
  refresh,
  scrollToElement,
  enable,
  disable,
})
</script>
