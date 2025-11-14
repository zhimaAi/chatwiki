<style lang="less">
.logic-flow-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}
/* 针对节点内容的foreignObject添加阴影 */


.logic-flow-container {
  width: 100%;
  height: 100%;
  overflow: hidden;

  .lf-node-content foreignObject {
    filter: drop-shadow(0 2px 3px rgba(0, 0, 0, 0.2)); /* 阴影效果 */
    /* 如需外边框，可补充stroke相关样式（但foreignObject本身已有stroke="#000"） */
    /* stroke: #000; */
    /* stroke-width: 2; */
  }

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
    position: fixed !important;
    padding: 6px;
    border: none;
    border-radius: 8px;
    background-color: #fff;
    filter: drop-shadow(0 4px 16px #0000001a);
    z-index: 100;
  }
}
</style>

<template>
  <div class="logic-flow-wrapper">
    <div ref="containerRef" class="logic-flow-container"></div>
    <TeleportContainer :flow-id="flowId" />

    <!-- <FloatAddBtn :lf="lf" v-if="lf" @addNode="onCustomAddNode" /> -->
    <CustomControl :lf="lf" v-if="lf" @runTest="handleRunTest" @addNode="onCustomAddNode" />

    <NodeFormDrawer
      ref="nodeFormDrawerRef"
      v-model:open="nodeFormDrawerShow"
      :destroyOnClose="true"
      :node-id="selectedNode.id"
      :node="selectedNode.properties"
      :node-type="selectedNode.type"
      :lf="lf"
      @update-node="handleNodeChange"
      @change-title="handleChangeNodeName"
      v-if="selectedNode"
    />
  </div>
</template>

<script setup>
import '@logicflow/core/lib/style/index.css'
import '@logicflow/extension/lib/style/index.css'
import { onMounted, ref } from 'vue'
// import FloatAddBtn from './float-add-btn.vue'
import CustomControl from './custom-control/index.vue'
import NodeFormDrawer from './node-form-drawer/index.vue'
import { generateUniqueId } from '@/utils/index'
import LogicFlow from '@logicflow/core'
import { DndPanel, MiniMap } from '@logicflow/extension'
import { register, getTeleport } from '@logicflow/vue-node-registry'
import { getNodesMap } from './node-list'
import customEdge from './edges/custom-line/index.js'
import startNode from './nodes/start-node/index.js'
import questionNode from './nodes/question-node/index.js'
import actionNode from './nodes/action-node/index.js'
import qaNode from './nodes/qa-node/index.js'
import aiDialogueNode from './nodes/ai-dialogue-node'
import httpNode from './nodes/http-node/index.js'
import knowledgeBaseNode from './nodes/knowledge-base-node/index.js'
import endNode from './nodes/end-node/index.js'
import explainNode from './nodes/explain-node/index.js'
import variableAssignmentNode from './nodes/variable-assignment-node/index.js'
import judgeNode from './nodes/judge-node/index.js'
import specifyReplyNode from './nodes/specify-reply-node/index.js'
import parameterExtractionNode from './nodes/parameter-extraction-node/index'
import problemOptimizationNode from './nodes/problem-optimization-node/index.js'
import addDataNode from './nodes/add-data-node/index.js'
import updateDataNode from './nodes/update-data-node/index.js'
import deleteDataNode from './nodes/delete-data-node/index.js'
import selectDataNode from './nodes/select-data-node/index.js'
import codeRunNode from './nodes/code-run-node/index.js'
import mcpNode from "./nodes/mcp-node/index.js";
import { ContextPad } from './plugins/context-pad/index.js'

const emit = defineEmits(['selectNode', 'deleteNode', 'runTest', 'blankClick'])

let lf = null
const nodeFormDrawerRef = ref(null)
const nodeFormDrawerShow = ref(false)
const selectedNode = ref(null)
const containerRef = ref(null)
const TeleportContainer = getTeleport()
const flowId = ref('')

const canvasData = {
  nodes: []
}

const miniMapOptions = {
  width: 160,
  height: 140,
  bottomPosition: 24,
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
      hideAnchors: false, // 是否隐藏锚点
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
    register(questionNode, lf)
    register(actionNode, lf)
    register(httpNode, lf)
    register(qaNode, lf)
    register(aiDialogueNode, lf)
    register(knowledgeBaseNode, lf)
    register(endNode, lf)
    register(explainNode, lf)
    register(judgeNode, lf)
    register(variableAssignmentNode, lf)
    register(specifyReplyNode, lf)
    register(parameterExtractionNode, lf)
    register(problemOptimizationNode, lf)
    register(addDataNode, lf)
    register(updateDataNode, lf)
    register(deleteDataNode, lf)
    register(selectDataNode, lf)
    register(codeRunNode, lf)
    register(mcpNode, lf)

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
      handleSelectedNode(data)
    })

    // 画布单击
    lf.on('blank:click', () => {
      handleBlankClick()
    })

    // 节点删除
    lf.on('node:delete', ({ data }) => {
      let dataRaw = JSON.parse(JSON.stringify(data));
      handleDeleteNode(dataRaw)
    })

    // 更新数据
    // lf.on('custom:setData', (node) => {
    //   console.log(node)
    // })

    // 更新NodeName
    // lf.on('custom:setNodeName', (data) => {
    //   console.log(data)
    // })

    lf.render(canvasData)

    lf.setZoomMiniSize(0.01)

    lf.setZoomMaxSize(8)

    // 监听 wheel 事件
    lf.container.addEventListener(
      'wheel',
      (e) => {
        // 检测是否按下 Shift 键
        if (e.shiftKey) {
          // 阻止默认行为（避免页面垂直滚动）
          e.preventDefault()

          // 获取当前画布的 transform 状态
          const transform = lf.getTransform()
          let SCALE_X = transform.SCALE_X
          let SCALE_Y = transform.SCALE_Y
          // 根据滚轮方向调整 x 坐标（deltaY 用于判断滚轮方向）
          // console.log(e.deltaY)
          // 设置新的 transform
          if (e.deltaY > 0) {
            lf.translate(-100 * SCALE_X, 100 * SCALE_Y)
          } else {
            lf.translate(100 * SCALE_X, -100 * SCALE_Y)
          }
        }
      },
      { passive: false }
    ) // passive: false 允许调用 preventDefault()

    lf.extension.miniMap.show()
  }
}

const getData = () => {
  let data = lf.getGraphRawData()
  return data
}

const setData = (data) => {
  let nodesMap = getNodesMap()

  data.nodes.forEach((node) => {
    // 设置开始节点的宽高
    let nodeCongfig = nodesMap[node.type]

    node.properties.node_icon_name = nodeCongfig.properties.node_icon_name

    if (!node.properties.width) {
      node.properties.width = nodeCongfig.width
    }

    if (!node.properties.height) {
      node.properties.height = nodeCongfig.height
    }
  })
  lf.clearData()

  lf.graphModel.graphDataToModel(data)
  lf.graphModel.translateCenter()
}

// 自定义添加节点
const onCustomAddNode = (data, model, anchorData) => {
  data.id = generateUniqueId(data.type)
  data.nodeSortKey = data.id.substring(0, 8) + data.id.substring(data.id.length - 8)

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
  if (data.height) {
    data.properties.height = data.height
  }

  let node = lf.getNodeModelById(data.id)

  node.height = data.properties.height || data.height || node.height
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

const noShowDrawerNode = ['explain-node', 'end-node']
// 选择节点
const handleSelectedNode = (data) => {
  console.log('----handleSelectedNode', data)
  const node = JSON.parse(JSON.stringify(data))
  node.properties.dataRaw =  node.properties.dataRaw || node.properties.node_params
  
  emit('selectNode', node)

  
  // 结束节点不支持编辑
  if(noShowDrawerNode.includes(data.type)){
    return
  }

  selectedNode.value = node
  nodeFormDrawerShow.value = true
}

const handleNodeChange = (data) => {
  selectedNode.value.properties = data

  console.log('---handleNodeChange', selectedNode.value)
  // 更新节点
  updateNode(JSON.parse(JSON.stringify(selectedNode.value)))
}

const handleChangeNodeName = (node_name) => {
  selectedNode.value.properties.node_name = node_name;
  console.log('---handleChangeNodeName', selectedNode.value)
  // 先更新数据
  updateNode(JSON.parse(JSON.stringify(selectedNode.value)))
  // 在发送事件之前，确保数据已经更新
  lf.graphModel.eventCenter.emit('custom:setNodeName',  {
    node_name: node_name, 
    node_id: selectedNode.value.id,
    node_type: selectedNode.value.type
  })
}

const handleDeleteNode = (data) => {
  emit('deleteNode', data)

  if(selectedNode.value && data.id === selectedNode.value.id) {
    nodeFormDrawerShow.value = false
    setTimeout(() => {
      selectedNode.value = null
    }, 350)
  }
}

const handleRunTest = () => {
  emit('runTest')
}

const handleBlankClick = () => {
  emit('blankClick')
  if(nodeFormDrawerShow.value) {
    nodeFormDrawerShow.value = false
  }
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
