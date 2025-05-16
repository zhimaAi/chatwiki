<style lang="less" scoped>
.chart-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  background-color: #F2F4F7;
}

.chart {
  width: 100%;
  height: 100%;
  user-select: none;
  touch-action: none;
}

.chart-mini-map {
  position: absolute;
  right: 0px;
  bottom: 82px;
  width: 200px;
  height: 160px;
  border-radius: 6px;
  border: 1px solid #fff;
  background-color: #fff;
  filter: drop-shadow(0 4px 16px rgba(0, 0, 0, 0.1));
}

.zoom-toolbar-wrapper{
  position: absolute;
  right: 12px;
  bottom: 30px;
}
</style>

<template>
  <div class="chart-wrapper">
    <div class="chart" id="chartBox" ref="chartBox"></div>
    <div class="chart-mini-map" id="chatMap"></div>
    <div class="zoom-toolbar-wrapper">
      <ZoomToolbar :value="initialZoom" @change="onZoomChange" />
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, reactive } from 'vue'
import { NVL } from '@neo4j-nvl/base'
import {
  ZoomInteraction,
  PanInteraction,
  // DragNodeInteraction, 
  HoverInteraction,
  ClickInteraction
} from '@neo4j-nvl/interaction-handlers'
import ZoomToolbar from '@/views/library/knowledge-graph/components/zoom-toolbar/index.vue'

const emit = defineEmits(['nodeClick'])

const chartBox = ref(null)
let nvl = null

const state = reactive({
  nodes: [],
  edges: []
})

const initialZoom = ref(1)

const callbacks = {
  // onLayoutDone: () => console.log('Layout done')
}

const init = () => {
  if (nvl) {
    nvl.destroy()
  }

  const options = {
    layout: 'forceDirected',
    initialZoom: initialZoom.value,
    disableTelemetry: true,
    maxZoom: 10,
    minZoom: 0.1,
    minimapContainer: document.getElementById('chatMap'),
    styling: {
      minimapViewportBoxColor: '#4399FF',
      // dropShadowColor: 'rgba(255, 255, 255, 1)', // 用于悬停节点和关系的投影的颜色
      selectedBorderColor: 'rgba(255, 255, 255, 1)', // 用于节点和关系的选定边界的颜色
      // selectedInnerBorderColor: 'rgba(255, 255, 255, 1)', // 用于选定节点的边框的颜色
    }
  }

  const container = document.getElementById('chartBox')

  nvl = new NVL(container, state.nodes, state.edges, options, callbacks)

  // 开启缩放和拖拽
  const zoomInteraction = new ZoomInteraction(nvl)
  // eslint-disable-next-line no-unused-vars
  const pan = new PanInteraction(nvl)

  zoomInteraction.updateCallback('onZoom', (zoomLevel) => {
    initialZoom.value = Number(zoomLevel.toFixed(1))
  })

  // 用于拖拽节点的交互处理程序
  // const dragNodeInteraction = new DragNodeInteraction(nvl)
  
  // 用于悬停节点的交互处理程序
  new HoverInteraction(nvl, { drawShadowOnHover: true })

  // 用于点击节点的交互处理程序
  const clickInteraction = new ClickInteraction(nvl, {
    selectOnClick: true
  })
  clickInteraction.updateCallback('onNodeClick', (node) => {
    emit('nodeClick', node)
  })

 
}

const onZoomChange = (value) => {
  initialZoom.value = value
  if (nvl) {
    nvl.setZoom(value)
  }
}

const add = ({ nodes, edges }) => {
  state.nodes = [state.nodes, ...nodes]
  state.edges = [state.edges, ...edges]

  if (nvl) {
    nvl.addElementsToGraph(nodes, edges)
  }
}

const addAndUpdate = ({ nodes, edges }) => {
  state.nodes = [state.nodes, ...nodes]
  state.edges = [state.edges, ...edges]

  if (nvl) {
    nvl.addAndUpdateElementsInGraph(nodes, edges)
  }
}

const update = ({ nodes, edges }) => {
  state.nodes = nodes
  state.edges = edges

  if (nvl) {
    nvl.updateElementsInGraph(nodes, edges)
  }
}

const destroy = () => {
  if (nvl) {
    nvl.destroy()
  }
}

const render = ({ nodes, edges }) => {
  state.nodes = nodes
  state.edges = edges

  init()
}

onMounted(() => {
  
})

defineExpose({
  render,
  add,
  addAndUpdate,
  update,
  destroy
})
</script>
