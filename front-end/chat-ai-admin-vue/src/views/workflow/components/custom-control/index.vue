<style lang="less" scoped>
.custom-control-warpper {
  position: absolute;
  right: 28px;
  bottom: 24px;
  z-index: 10;
  &:hover {
    &::after {
      content: '';
      position: absolute;
      top: -10px;
      left: 0;
      width: 100%;
      height: 30px;
    }
  }
  .custom-control {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    padding: 4px 12px;
    border-radius: 8px;
    background-color: #fff;
    box-shadow: 0 4px 16px 0 #0000001a;
  }

  .action-btn {
    border-radius: 6px;
    transition: all 0.2s;
    &:hover {
      background-color: #e4e6eb;
      cursor: pointer;
    }
  }

  .zoom-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    font-size: 16px;
    color: #595959;
  }
}
</style>

<template>
  <div class="custom-control-warpper" @mouseenter="handleMouseenter" @mouseleave="handleMouseleave">
    <div class="custom-control">
      <div class="action-btn zoom-btn" @click="handleReduce">
        <svg-icon name="minus" size="16" />
      </div>
      <zoom-select :title="zoomSelectTitle" @change="zoomSelectChagne" />
      <div class="action-btn zoom-btn" @click="handleAmplify">
        <svg-icon name="plus" size="16" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, inject, onMounted, onBeforeUnmount } from 'vue'
import ZoomSelect from './zoom-select.vue'

const props = defineProps({
  lf: {
    type: Object,
    default: () => ({}),
    required: true
  }
})

const { eventCenter } = props.lf.graphModel

const zoom = ref(100)

const setZoom = () => {
  props.lf.zoom(zoom.value / 100)
}

const zoomSelectTitle = ref('100%')

const zoomSelectChagne = ({ label, value }) => {
  zoomSelectTitle.value = label
  zoom.value = value

  setZoom()
}

const handleReduce = () => {
  // 四舍五入到整数
  let value = (zoom.value = Math.floor(zoom.value / 2))

  if (value < 1) {
    value = 1
  }

  zoom.value = value
  zoomSelectTitle.value = `${value}%`
  setZoom()
}

const handleAmplify = () => {
  let value = zoom.value * 2

  if (value > 800) {
    value = 800
  }

  zoom.value = value
  zoomSelectTitle.value = `${value}%`
  setZoom()
}

const onGraphTransform = (args) => {
  let value = Math.floor(args.transform.SCALE_X * 100)
  zoom.value = value
  zoomSelectTitle.value = `${value}%`
}
let miniMap = null
const handleMouseenter = () => {
  props.lf.extension.miniMap.show()
  miniMap = document.querySelector('.lf-mini-map')
  starListenEvent()
}

const handleMouseleave = (e) => {
  if (e.y > 800) {
    props.lf.extension.miniMap.hide()
  }
}

function starListenEvent() {
  if (!miniMap) {
    return
  }
  miniMap.addEventListener('mouseenter', () => {
    props.lf.extension.miniMap.show()
  })
  miniMap.addEventListener('mouseleave', () => {
    props.lf.extension.miniMap.hide()
  })
}

onMounted(() => {
  eventCenter.on('graph:transform', onGraphTransform)
})

onBeforeUnmount(() => {
  eventCenter.off('graph:transform', onGraphTransform)
})
</script>
