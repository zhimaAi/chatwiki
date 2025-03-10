<style lang="less">
.logic-flow-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}
.logic-flow-container {
  width: 100%;
  height: 100%;
  overflow: hidden;

  /* 自定义锚点样式 */
  .custom-anchor {
    cursor: pointer;

    .custom-anchor-arrow,
    .custom-anchor-plus {
      display: none;
    }
  }
  .custom-anchor.anchor-selected {
    .custom-anchor-arrow {
      display: block;
    }
  }
  .custom-anchor.anchor-not-selected.custom-anchor-right:hover {
    .custom-anchor-plus {
      display: block !important;
    }
  }

  .custom-anchor .anchor-hide {
    display: none !important;
  }

  .lf-mini-map {
    padding: 6px;
    border: none;
    border-radius: 8px;
    background-color: #fff;
    filter: drop-shadow(0 4px 16px #0000001a);
  }
}
</style>

<template>
  <div class="logic-flow-wrapper">
    <div ref="containerRef" class="logic-flow-container"></div>
    <TeleportContainer :flow-id="flowId" />

    <FloatAddBtn :lf="lf" v-if="lf" @addNode="onCustomAddNode" />
    <CustomControl :lf="lf" v-if="lf" />
  </div>
</template>

<script setup>
import { getNodeHeight, getNodeWidth } from './util'
import '@logicflow/core/lib/style/index.css'
import '@logicflow/extension/lib/style/index.css'
import { onMounted, ref } from 'vue'
import FloatAddBtn from './float-add-btn.vue'
import CustomControl from './custom-control/index.vue'
import { generateUniqueId } from '@/utils/index'
import LogicFlow from '@logicflow/core'
import { DndPanel, MiniMap } from '@logicflow/extension'
import { register, getTeleport } from '@logicflow/vue-node-registry'
import customEdge from './edges/custom-line/index.js'
import startNode from './nodes/start-node/index.js'
import questionNode from './nodes/question-node/index.js'
import messageNode from './nodes/message-node/index.js'
import actionNode from './nodes/action-node/index.js'
import qaNode from './nodes/qa-node/index.js'
import aiDialogueNode from './nodes/ai-dialogue-node'
import httpNode from './nodes/http-node/index.js'
import knowledgeBaseNode from './nodes/knowledge-base-node/index.js'
import endNode from './nodes/end-node/index.js'
import judgeNode from './nodes/judge-node'
import { ContextPad } from './plugins/context-pad/index.js'

const emit = defineEmits(['selectedNode', 'deleteNode'])

let lf = null
const containerRef = ref(null)
const TeleportContainer = getTeleport()
const flowId = ref('')

const canvasData = {
  nodes: []
}

const miniMapOptions = {
  width: 200,
  height: 160,
  bottomPosition: 76,
  rightPosition: 16
}

function initLogicFlow() {
  if (containerRef.value) {
    lf = new LogicFlow({
      container: containerRef.value,
      // width: containerRef.value.offsetWidth,
      // height: containerRef.value.offsetHeight,
      nodeTextEdit: false,
      edgeTextEdit: false,
      textEdit: false,
      grid: false,
      adjustEdge: false, // 允许调整边
      adjustEdgeStartAndEnd: false, // 是否允许拖动边的端点来调整连线
      adjustNodePosition: true, // 是否允许拖动节点
      edgeSelectedOutline: true, // 边被选中时是否显示边的外框
      nodeSelectedOutline: true, // 节点被选中时是否显示节点的外框
      hoverOutline: true, // 鼠标hover节点时是否显示节点的外框
      background: {
        backgroundColor: '#f0f2f5'
      },
      plugins: [DndPanel, MiniMap, ContextPad],
      pluginsOptions: {
        miniMap: miniMapOptions
      }
    })

    register(customEdge, lf)
    register(startNode, lf)
    register(messageNode, lf)
    register(questionNode, lf)
    register(actionNode, lf)
    register(httpNode, lf)
    register(qaNode, lf)
    register(aiDialogueNode, lf)
    register(knowledgeBaseNode, lf)
    register(endNode, lf)
    register(judgeNode, lf)
    

    lf.setDefaultEdgeType('custom-edge')

    lf.on('graph:rendered', ({ graphModel }) => {
      flowId.value = graphModel.flowId
    })

    // 设置托拽节点时的zIndex
    let dragNodeZIndex = 1

    lf.on('node:dragstart', ({ data }) => {
      let node = lf.graphModel.getElement(data.id)

      dragNodeZIndex = node.zIndex

      node.setZIndex(999999)
    })

    // 恢复托拽节点的zIndex
    lf.on('node:drop', ({ data }) => {
      let node = lf.graphModel.getElement(data.id)

      node.setZIndex(dragNodeZIndex)
    })

    // 添加自定义节点事件
    lf.on('custom:addNode', ({ data, model, anchorData }) => {
      onCustomAddNode(data, model, anchorData)
    })
    // 点击节点
    lf.on('node:click', ({ data }) => {
      emit('selectedNode', JSON.parse(JSON.stringify(data)))
    })

    // 节点删除
    lf.on('node:delete', ({ data }) => {
      emit('deleteNode', JSON.parse(JSON.stringify(data)))
    })

    lf.render(canvasData)

    lf.setZoomMiniSize(0.01)

    lf.setZoomMaxSize(8)

    lf.extension.miniMap.show()
  }
}

const getData = () => {
  let data = lf.getGraphRawData()
  return data
}

const setData = (data) => {
  data.nodes.forEach((node) => {
    // 目前所有的节点width
    node.properties.width = node.width = getNodeWidth(node)
    node.properties.height = node.height = getNodeHeight(node)
  })

  lf.graphModel.graphDataToModel(data)
  lf.graphModel.translateCenter()
}

// 自定义添加节点
const onCustomAddNode = (data, model, anchorData) => {
  console.log(data,model,anchorData,'===')
  data.id = generateUniqueId(data.type)
  data.nodeSortKey = data.id.substring(0, 8) + data.id.substring(data.id.length - 8)
  data.width = getNodeWidth(data)
  data.height = getNodeHeight(data)
  data.properties.width = data.width
  data.properties.height = data.height
  data.properties.nodeSortKey = data.nodeSortKey
  
  if (anchorData) {
    data.x = anchorData.x + data.width + 100
    data.y = anchorData.y + data.height / 2 - 24
  } else {
    const { transformModel } = lf.graphModel
    const point = transformModel.HtmlPointToCanvasPoint([
      lf.graphModel.width / 2,
      lf.graphModel.height / 2
    ])

    data.x = point[0]
    data.y = point[1]
  }
  // 情况选中状态
  let zIndex = 0
  lf.graphModel.nodes.forEach((node) => {
    if (node.zIndex > zIndex) {
      zIndex = node.zIndex
    }
  })

  zIndex = zIndex + 1
  
  let node = lf.addNode(data)

  if (anchorData) {
    lf.graphModel.addEdge({
      type: 'custom-edge',
      sourceNodeId: model.id,
      targetNodeId: node.id,
      sourceAnchorId: anchorData.id,
      startPoint: {
        x: anchorData.x,
        y: anchorData.y
      }
    })
  }
  node.setZIndex(zIndex)
  lf.graphModel.clearSelectElements()
  node.setSelected(true)
}

const updateNode = (data) => {
  data.properties.height = data.height = getNodeHeight(data)

  let node = lf.getNodeModelById(data.id)

  node.height = data.height
  node.properties = data.properties

  if (data.properties.node_type == 2 && data.properties.node_sub_type == 21) {
    // 删除已经没有了的边
    const edgeModels = node.graphModel.getNodeEdges(data.id)

    edgeModels.forEach((edge) => {
      // 锁定当前节点上的锚点，如果锚点不存在了，则删除连线
      if (edge.sourceNodeId == node.id) {
        let anchor = node.getAnchorInfo(edge.sourceAnchorId)
        // 如果锚点不存在了，则删除连线
        if (!anchor) {
          node.graphModel.deleteEdgeById(edge.id)
        }
      }
    })
  }

  node.refreshBranch()
}

onMounted(() => {
  initLogicFlow()
})

defineExpose({
  setData,
  getData,
  updateNode
})
</script>
