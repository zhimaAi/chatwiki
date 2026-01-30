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
  .lf-drag-able{
    cursor: grab !important;
  }
  .lf-dragging{
    cursor: grabbing !important;
  }
  .lf-node-content foreignObject {
    filter: drop-shadow(0 2px 3px rgba(0, 0, 0, 0.2)); /* 阴影效果 */
    /* 如需外边框，可补充stroke相关样式（但foreignObject本身已有stroke="#000"） */
    /* stroke: #000; */
    /* stroke-width: 2; */
  }

  .lf-edge-selected foreignObject{
    filter: drop-shadow(0 2px 3px rgba(0, 0, 0, 0.2)); /* 阴影效果 */
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
.dont-show-again-checkbox {
 position: absolute;
 bottom: 26px;
 left: 24px;
 font-size: 12px;
}
</style>

<template>
  <div class="logic-flow-wrapper">
    <div ref="containerRef" class="logic-flow-container"></div>
    <TeleportContainer :flow-id="flowId" />

    <CustomControl :lf="lf" v-if="lf" @runTest="handleRunTest" @addNode="onCustomAddNode" @zoomChange="handleZoomChange" @autoLayout="handleAutoLayout" />

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
import { onMounted, ref, onUnmounted, h, nextTick } from 'vue'
import { useStorage } from '@/hooks/web/useStorage'
import { Modal, Checkbox, message } from 'ant-design-vue'
import CustomControl from './custom-control/index.vue'
import NodeFormDrawer from './node-form-drawer/index.vue'
import { generateUniqueId } from '@/utils/index'
import LogicFlow from '@logicflow/core'
import { DndPanel, MiniMap, SelectionSelect, DynamicGroup } from '@logicflow/extension'
import { Elk } from './plugins/elk/elk.js'
import { GroupAutoResize } from './plugins/group-auto-resize'
import { register, getTeleport } from '@logicflow/vue-node-registry'
import { getNodesMap, createTriggerNode } from './node-list'
import customLineEdge from './edges/custom-line/index.js'
// 贝塞尔曲线
import customBezierEdge from './edges/custom-bezier/index.js'
import { ContextPad } from './plugins/context-pad/index.js'
import { CanvasHistory } from './plugins/canvas-history/index.js'
import { CustomKeyboard } from './plugins/custom-keyboard/index.js'
import registerCustomNodes from './nodes/index.js'

const emit = defineEmits(['selectNode', 'onDeleteNode', 'onDeleteEdge', 'runTest', 'blankClick'])

let lf = null
const lfRef = ref(null)
const { setStorage, getStorage } = useStorage();
const DONT_SHOW_DELETE_CONFIRM_KEY = 'dont_show_delete_node_confirm';
const dontShowDeleteConfirm = ref(getStorage(DONT_SHOW_DELETE_CONFIRM_KEY) || false);
const nodeFormDrawerRef = ref(null)
const nodeFormDrawerShow = ref(false)
const selectedNode = ref(null)
const containerRef = ref(null)
const TeleportContainer = getTeleport()
const flowId = ref('')
const selectedElements = ref([]) // 用于存储当前“选中元素”


const canvasData = {
  nodes: [
  ]
}

const miniMapOptions = {
  width: 160,
  height: 140,
  bottomPosition: 24,
  rightPosition: 16
}

const selectionSelectOptions = {
  exclusiveMode: false,
}

const handleZoomChange = (value) => {
  // zoom.value = value
  if (lf) {
    lf.zoom(value)
  }
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
      stopMoveGraph: false,
      stopZoomGraph: true,
      adjustEdge: false, // 允许调整边
      adjustEdgeStartAndEnd: false, // 是否允许拖动边的端点来调整连线
      adjustNodePosition: true, // 是否允许拖动节点
      edgeSelectedOutline: true, // 边被选中时是否显示边的外框
      nodeSelectedOutline: true, // 节点被选中时是否显示节点的外框
      hoverOutline: true, // 鼠标hover节点时是否显示节点的外框
      hideAnchors: false, // 是否隐藏锚点
      // overlapMode: 0, // 元素重合的堆叠模式, 默认模式：0，节点层级优先，该模式下连线在下、节点在上、选中元素临时置
      background: {
        backgroundColor: '#f0f2f5'
      },
      plugins: [
        DndPanel, 
        MiniMap, 
        ContextPad, 
        SelectionSelect, 
        CanvasHistory, 
        CustomKeyboard, 
        DynamicGroup,
        GroupAutoResize,
        Elk
      ],
      pluginsOptions: {
        miniMap: miniMapOptions,
        selectionSelect: selectionSelectOptions,
        canvasHistory: {
          maxHistorySize: 100
        }
      },
      // history: false, // 关闭历史记录功能会导致小地图无法更新
      keyboard: {
        enabled: false,
        shortcuts: [{
          keys: [ // 屏蔽自带的ctrl + z, ctrl + y, ctrl + c
            "cmd + z",
            "ctrl + z",
            "cmd + y",
            "ctrl + y",
            "cmd + c",
            "ctrl + c",
            'cmd + v',
            'ctrl + v'
          ],
          callback: () => {
            // 自定义逻辑
          },
        }],
      }
    })

    lfRef.value = lf

    registerCustomNodes(lf)
    register(customLineEdge, lf)
    register(customBezierEdge, lf)
    

    // lf.setDefaultEdgeType('custom-edge')
    lf.setDefaultEdgeType('custom-bezier-edge')

    lf.on('graph:rendered', ({ graphModel }) => {
      flowId.value = graphModel.flowId
    })

    // 设置托拽节点时的zIndex
    let dragNodeZIndexMap = {}
    let nodeChildren = []

    lf.on('node:dragstart', ({ data }) => {
      let node = lf.graphModel.getElement(data.id)
      dragNodeZIndexMap[node.id] = node.zIndex

      nodeChildren = []

      let nodeGroup = lf.graphModel.dynamicGroup.getGroupByNodeId(data.id)

      if(nodeGroup){
        nodeGroup.properties.width = nodeGroup._width
        nodeGroup.properties.height = nodeGroup._height
      }else{
        if(!node.properties.loop_parent_key){
          // 如果是非分组路面的节点  禁止拖到分组里面去
          node.properties.disabled_add_group = true
        }
      }

      // 判断是不是组
      if(data.children && data.children.length > 0){
        node.setZIndex(-9999)
        nodeChildren = data.children || []

        nodeChildren.forEach((item) => {
          let nodeModel = lf.getNodeModelById(item);

          dragNodeZIndexMap[nodeModel.id] = nodeModel.zIndex;
          nodeModel.setZIndex(999999)
        })
      }else{
        node.setZIndex(999999)
      }

      node.refreshBranch()
    })

    // 恢复托拽节点的zIndex
    lf.on('node:drop', ({ data }) => {
      let node = lf.graphModel.getElement(data.id)
      let nodeGroup = lf.graphModel.dynamicGroup.getGroupByNodeId(data.id)

      // 拖拽过程中组的子节点可能会丢失，需要重新添加
      if(data.properties.loop_parent_key && !nodeGroup){
        const groupModel = lf.getNodeModelById(data.properties.loop_parent_key);
        groupModel.addChild(data.id)
      }

      if(data.children && data.children.length){
        node.setZIndex(-9999)
        data.children.forEach(item => {
          let nodeModel = lf.getNodeModelById(item);
          nodeModel.setZIndex(dragNodeZIndexMap[nodeModel.id])
        })
      }else{
        node.setZIndex(dragNodeZIndexMap[node.id])
      }

      node.refreshBranch()
    })

    // 节点拖拽时动态调整边的offset
    lf.on('node:drag', ({ data }) => {
      const nodeId = data.id;
      // 获取与当前节点相连的所有边
      const relatedEdges = lf.graphModel.getNodeEdges(nodeId);
      const nodeGroup = lf.graphModel.dynamicGroup.getGroupByNodeId(data.id)

      // 拖拽过程中组的子节点可能会丢失，需要重新添加
      // if(data.properties.loop_parent_key && !nodeGroup){
      //   const groupModel = lf.getNodeModelById(data.properties.loop_parent_key);
      //   groupModel.addChild(data.id)
      // }

      setTimeout(() => {
        relatedEdges.forEach(edge => {
          // 只处理目标边类型（如自定义贝塞尔边）
          if (edge.type === 'custom-bezier-edge') {
            // 获取边的起点和终点坐标
            const { startPoint, endPoint } = edge;
            
            // 计算边的长度（勾股定理）
            const dx = endPoint.x - startPoint.x;
            const dy = endPoint.y - startPoint.y;
            const length = Math.sqrt(dx * dx + dy * dy);
            
            // 计算新的offset（同上）
            const baseOffset = 10;
            const scale = 0.3;
            let newOffset = baseOffset + length * scale;
            const minOffset = 10;
            const maxOffset = 1000;
            newOffset = Math.round(newOffset)       
            newOffset = Math.max(minOffset, Math.min(newOffset, maxOffset));          

            edge.setProperties({
              offset: newOffset,
            });
          }
        });

        if(nodeGroup){
          nodeGroup.refreshBranch()
        }
      }, 0)
    });

    lf.on('node:add', () => {
      // console.log('node:add', data)
    })

    // 添加自定义节点事件
    lf.on('custom:addNode', (options) => {
      onCustomAddNode(options)
    })

    lf.on('custom:addGroupNode', ({data, group_id}) => {
      onCustomAddGroupNode(data ,group_id)
    })

    lf.on('group:add-node', ({ childId }) => {
      const childModel = lf.getNodeModelById(childId);
      // const groupModel = lf.getNodeModelById(data.id);
      if(!childModel.properties.loop_parent_key){
        // 移除
        // console.log('移除')
        // groupModel.removeChild(childId);
      }
      // 可以在这里执行自定义逻辑，比如更新UI状态
      // updateGroupInfo(data.id, childId, 'add')
    })
    
    // 元素点击
    lf.on('element:click', (e) => {
      selectedElements.value = [e.data]
    })

    // 点击节点
    lf.on('node:click', ({ data }) => {
      handleSelectedNode(data)
    })

    // 画布单击
    lf.on('blank:click', () => {
      handleBlankClick()
    })

    // 选区框选后触发
    lf.on('selection:selected', (e) => {
      let items = []

      e.elements.forEach((element) => {
        let data = lf.getDataById(element.id);
        if(element.BaseType === 'node' || element.BaseType === 'edge'){
          items.push(data)
        }
      })

      selectedElements.value = items
    })

    // 节点删除
    lf.on('node:delete', ({ data }) => {
      let dataRaw = JSON.parse(JSON.stringify(data));
      onDeleteNode(dataRaw)
    })

    // 边删除
    lf.on('edge:delete', ({ data }) => {
      let dataRaw = JSON.parse(JSON.stringify(data));
      onDeleteEdge(dataRaw)
    })

    // 更新数据
    // lf.on('custom:setData', (node) => {
    //   console.log(node)
    // })

    // 更新NodeName
    // lf.on('custom:setNodeName', (data) => {
    //   console.log(data)
    // })

    // 自定义节点删除
    lf.on('custom:node:delete', (node) => {
      handleDeleteNode(node)
    })

    // 自定义边删除
    lf.on('custom:edge:delete', (edge) => {
      console.log(edge, '==edge')
      handleDeleteEdge(edge)
    })

    // 自定义键盘删除
    lf.on('custom:keyoard:delete', () => {
      handleKeyoardDelete()
    })

    // 监听批量粘贴事件
    lf.on('custom:paste', handleBatchPaste)
    // 监听历史记录状态变化
    // lf.extension.canvasHistory.onHistoryChange((state) => {
      // console.log('可撤销:', state.canUndo)
      // console.log('可重做:', state.canRedo)
      // console.log('历史记录数:', state.historySize)
      // console.log('重做记录数:', state.redoSize)
      // console.log('操作日志:', state)
      // 撤销或重做操作后，清空选中元素
      // selectedElements.value = []
    // })

    lf.render(canvasData)

    lf.setZoomMiniSize(0.01)

    lf.setZoomMaxSize(8)
  }
}


const createTriggerElements= (startNode) => {
  const nodesMap = getNodesMap()
  const node_params = JSON.parse(startNode.properties.node_params)
  const { trigger_list } = node_params.start;

  const triggerNodes = []
  const triggerEdges = []
  let offsetTop = 0;

  trigger_list.forEach((item) => {
    const type = `trigger_${item.trigger_type}`
    const nodeId = generateUniqueId(type)
    const nodeCongfig = nodesMap[type]
    const nodeSortKey = nodeId.substring(0, 8) + nodeId.substring(nodeId.length - 8)
    const edgeId = generateUniqueId(type)
    // const icon = item.trigger_icon ? item.trigger_icon : nodeCongfig.properties.node_icon
    const triggerNode = createTriggerNode(item)

    triggerNode.id = nodeId
    triggerNode.x = startNode.x - 600
    triggerNode.y = startNode.y + offsetTop
    triggerNode.properties.node_key = nodeId
    triggerNode.properties.nodeSortKey = nodeSortKey

    const edge = {
      id: edgeId,
      sourceNodeId: nodeId,
      targetNodeId: startNode.id,
      properties: {
        isTriggerEdge: true,
      }
    }

    offsetTop += nodeCongfig.height + 50

    triggerNodes.push(triggerNode)
    triggerEdges.push(edge)
  })
  
  return {
    nodes: triggerNodes,
    edges: triggerEdges,
  };
}

const getData = () => {
  const data = JSON.parse(JSON.stringify(lf.getGraphRawData()))
  const nodes = []
  const edges = []
  const triggers = []

  data.edges.forEach(edge => {
    if(!edge.properties.isTriggerEdge){
      edges.push(edge)
    }
  })

  data.nodes.forEach(node => {
    if(node.properties.isTriggerNode){
      triggers.push(node)
    }else{
      nodes.push(node)
    }
  })

  data.nodes = nodes;
  data.edges = edges;
  data.triggers = triggers;

  return data;
}

const setData = (data) => {
  const nodesMap = getNodesMap()
  let triggerElements = {
    nodes: [],
    edges: [],
  }

  data.nodes.forEach((node) => {
    // 设置开始节点的宽高
    let nodeCongfig = nodesMap[node.type]
    node.properties.node_icon_name = nodeCongfig.properties.node_icon_name
    if(!node.properties.node_header_bg_color) {
      node.properties.node_header_bg_color = nodeCongfig.properties.node_header_bg_color || '#fff'
    }

    if (!node.properties.width) {
      node.properties.width = nodeCongfig.width
    }

    if (!node.properties.height) {
      node.properties.height = nodeCongfig.height
    }
    // 处理开始节点的触发器
    if (node.type == 'start-node') {
      triggerElements = createTriggerElements(node)
    }
  })

  data.nodes.unshift(...triggerElements.nodes)

  lf.clearData()
  lf.graphModel.graphDataToModel(data)
  lf.graphModel.translateCenter()
  lf.extension.miniMap.show()

  // 添加触发器边, 调用addEdge来产生边这样就不用计算边的位置了
  triggerElements.edges.forEach(edge => {
    lf.graphModel.addEdge(edge)
  })

  const history = lf.extension.canvasHistory
  if (history) {
    history.setInitialState(data)
  }

  handleSetNodeToGroup(data.nodes) // 给节点增加绑定关系
}

const handleSetNodeToGroup = (nodes)=>{
  nextTick(()=>{
    nodes.forEach((node) => {
    if(node.loop_parent_key){
        const groupModel = lf.getNodeModelById(node.loop_parent_key);
        if(groupModel){
          groupModel.addChild(node.id)
          lf.graphModel.dynamicGroup.sendNodeToFront(groupModel);
        }
      }
    })
  })
}

// 自定义添加节点
let nodeNameMap = {}
// 生成trigger节点位置
const getTriggerNodePosition = () => { 
  // 获取开始节点的位置
  const { nodes } = lf.graphModel
  const startNode = nodes.find((node) => node.type === 'start-node')

  const x = startNode.x - 600
  const y = startNode.y

  return {
    x,
    y
  }
}

// 调整组大小以适应节点
const resizeGroup = (groupModel, node) => {
  // 计算节点的边界
  const nodeLeftBoundary = node.x - node.width / 2;
  const nodeRightBoundary = node.x + node.width / 2;
  const nodeTopBoundary = node.y - node.height / 2;
  const nodeBottomBoundary = node.y + node.height / 2;
  
  // 计算组的边界
  const groupLeftBoundary = groupModel.x - groupModel._width / 2;
  const groupRightBoundary = groupModel.x + groupModel._width / 2;
  const groupTopBoundary = groupModel.y - groupModel._height / 2;
  const groupBottomBoundary = groupModel.y + groupModel._height / 2;
  
  let needResize = false;
  let newWidth = groupModel._width;
  let newHeight = groupModel._height;
  let deltaX = 0;
  let deltaY = 0;
  
  // 检查右边界
  if(nodeRightBoundary > groupRightBoundary - 20){
    const extraWidth = nodeRightBoundary - groupRightBoundary + 50;
    newWidth += extraWidth;
    deltaX += extraWidth / 2;
    needResize = true;
  }
  
  // 检查左边界
  if(nodeLeftBoundary < groupLeftBoundary + 20){
    const extraWidth = groupLeftBoundary - nodeLeftBoundary + 50;
    newWidth += extraWidth;
    deltaX -= extraWidth / 2;
    needResize = true;
  }
  
  // 检查下边界
  if(nodeBottomBoundary > groupBottomBoundary - 20){
    const extraHeight = nodeBottomBoundary - groupBottomBoundary + 50;
    newHeight += extraHeight;
    deltaY += extraHeight / 2;
    needResize = true;
  }
  
  // 检查上边界
  if(nodeTopBoundary < groupTopBoundary + 20){
    const extraHeight = groupTopBoundary - nodeTopBoundary + 50;
    newHeight += extraHeight;
    deltaY -= extraHeight / 2;
    needResize = true;
  }
  
  // 如果需要调整大小，则更新组的尺寸和位置
  if(needResize){
    groupModel.properties.width = newWidth;
    groupModel.properties.height = newHeight;
    groupModel._width = newWidth;
    groupModel._height = newHeight;
    groupModel.width = newWidth;
    groupModel.height = newHeight;
    groupModel.x += deltaX;
    groupModel.y += deltaY;
  }
};


const createNodeInfo = (options) => {
  const data = options.node || options;
  const anchorData = options.anchorData; // 来自右键菜单添加
  const dropEvent = options.event;


  data.id = data.id || generateUniqueId(data.type)
  data.nodeSortKey = data.id.substring(0, 8) + data.id.substring(data.id.length - 8)
  
  // 处理节点的group绑定关系
  if(anchorData){
    let nodeGroup = lf.graphModel.dynamicGroup.getGroupByNodeId(anchorData.nodeId)
    if(!nodeGroup){
      data.properties.disabled_add_group = true
    }else{
      data.properties.loop_parent_key = nodeGroup.id
    }
  }

  // 如果不是复制粘贴的节点，则处理节点名称递增
  if (!options.isCopy) {
    // 同一类型的节点多次添加时，从第二次添加开始，默认名称后面加上序号
    if(nodeNameMap[data.type]){
      nodeNameMap[data.type] = nodeNameMap[data.type] + 1
      data.properties.node_name = data.properties.node_name + nodeNameMap[data.type]
    }else{
      nodeNameMap[data.type] = 1
    }
  }

  data.properties.width = data.width
  data.properties.height = data.height
  data.properties.nodeSortKey = data.nodeSortKey

  // 核心改动：判断坐标来源
  if (dropEvent) {
    // === 拖拽添加 ===
    const { transformModel } = lf.graphModel;
    // 1. 将浏览器视口坐标(clientX, clientY)转换为相对于画布容器的坐标
    const { left, top } = containerRef.value.getBoundingClientRect();
    const [canvasX, canvasY] = transformModel.HtmlPointToCanvasPoint([
      dropEvent.clientX - left,
      dropEvent.clientY - top
    ]);
    // 2. LogicFlow默认(x, y)为中心点，需减去宽高一半将其校正为左上角
    data.x = canvasX + data.width / 2;
    data.y = canvasY + data.height / 2;
  } else if(!data.x && !data.y) { // 只有在没有预设坐标时才计算位置
    if (anchorData) {
    // === 从锚点添加 ===
      data.x = anchorData.x + data.width + 100;
      data.y = anchorData.y + data.height / 2 - 24;
    } else {
      // 触发器放到开始节点左边
      if(options.isTriggerNode){
        let pos = getTriggerNodePosition()
        data.x = pos.x;
        data.y = pos.y;
        data.properties.disabled_add_group = true
      }else{
        // === 默认（点击）添加，放在画布中心 ===
        const { transformModel } = lf.graphModel;
        const point = transformModel.HtmlPointToCanvasPoint([
          lf.graphModel.width / 2,
          lf.graphModel.height / 2
        ]);
        data.x = point[0];
        data.y = point[1];
        data.properties.disabled_add_group = true // 从底部按钮添加的 禁止添加进入分组
      }
    }
  }
  return data;
}

const onAddLoopGroup = (options) => {
  setTimeout(() => {
    onCustomAddGroupNode(
      {
        node: {
          type: 'group-start-node',
          with: 200,
          properties: {
            node_type: 27,
            node_name: '循环开始',
            node_icon_name: 'start-node',
            node_params: JSON.stringify({})
          }
        }
      },
      options.id
    )

  }, 100)
}

const onAddBatchGroup = (options) => {
  setTimeout(() => {
    onCustomAddGroupNode(
      {
        node: {
          type: 'group-start-node',
          with: 200,
          properties: {
            node_type: 31,
            node_name: '批量执行开始',
            node_icon_name: 'start-node',
            node_params: JSON.stringify({})
          }
        }
      },
      options.id
    )

  }, 100)
}

const onCustomAddNode = (options) => {
  const nodeData = createNodeInfo(options);
  const model = options.model; // 来自右键菜单添加
  const anchorData = options.anchorData; // 来自右键菜单添加

  // 情况选中状态
  let hasSessionNode = false
  // let zIndex = 0

  lf.graphModel.nodes.forEach((node) => {
    // if (node.zIndex > zIndex) {
    //   zIndex = node.zIndex
    // }

    if(!hasSessionNode && node.type == 'session-trigger-node'){
      hasSessionNode = true
    }
  })

  // zIndex = zIndex + 1

  if(hasSessionNode && nodeData.type == 'session-trigger-node'){
    return message.error('请勿添加多个会话触发器')
  }

  // nodeData.zIndex = zIndex
  // nodeData.groupZindex = zIndex

  let node = lf.addNode(nodeData)

  if (node.type == 'custom-group') {
    onAddLoopGroup(node)
  }

  if(node.type == 'batch-group'){
    onAddBatchGroup(node)
  }

  if(options.isTriggerNode){
    const { nodes } = lf.graphModel
    const startNode = nodes.find((node) => node.type === 'start-node')

    let edge = {
      type: 'custom-bezier-edge',
      sourceNodeId: nodeData.id,
      targetNodeId: startNode.id,
      properties: {
        isTriggerEdge: true,
      }
    }

    lf.graphModel.addEdge(edge)
    lf.graphModel.eventCenter.emit("custom:trigger-add", nodeData);
  }

  if (anchorData) {
    let targetAnchor = node.anchors.find((anchor) => anchor.type === 'left')

    lf.graphModel.addEdge({
      type: 'custom-bezier-edge',
      sourceNodeId: model.id,
      targetNodeId: node.id,
      sourceAnchorId: anchorData.id,
      targetAnchorId: targetAnchor.id,
      startPoint: {
        x: anchorData.x,
        y: anchorData.y
      },
      endPoint: {
        x: targetAnchor.x,
        y: targetAnchor.y
      }
    })

    // 获取锚点节点的分组
    const nodeGroup = lf.graphModel.dynamicGroup.getGroupByNodeId(anchorData.nodeId)
    // 如果有分组，则将节点添加到分组中并且调整分组大小和位置
    if(nodeGroup && nodeGroup.id){
      const groupModel = lf.getNodeModelById(nodeGroup.id);

      groupModel.addChild(node.id);
      
      resizeGroup(groupModel, node);
      
      lf.graphModel.dynamicGroup.sendNodeToFront(groupModel);
    }
  }

  // node.setZIndex(zIndex)

  lf.graphModel.clearSelectElements()
  node.setSelected(true)
}

const onCustomAddGroupNode = (options, group_id) => {
  const data = options.node || options;
  data.id = generateUniqueId(data.type)
  data.nodeSortKey = data.id.substring(0, 8) + data.id.substring(data.id.length - 8)

  data.properties.width = data.width
  data.properties.height = data.height
  data.properties.nodeSortKey = data.nodeSortKey
  data.properties.loop_parent_key = group_id  // 父节点的id
  const groupModel = lf.getNodeModelById(group_id);

  data.x = groupModel.x 
  data.y = groupModel.y 
  // 情况选中状态
  let zIndex = 0
  lf.graphModel.nodes.forEach((node) => {
    if (node.zIndex > zIndex) {
      zIndex = node.zIndex
    }
  })

  zIndex = zIndex + 99

  let node = lf.addNode(data)

  groupModel.addChild(node.id);
  node.setZIndex(zIndex)
  lf.graphModel.clearSelectElements()
  lf.graphModel.dynamicGroup.sendNodeToFront(groupModel);
  // node.setSelected(true)
}

const handleAutoLayout = async () => {
  const graphData = lf.getGraphRawData()
  const history = lf.extension.canvasHistory

  // 开始事务，将布局前的状态记录下来
  if (history) {
    history.beginTransaction()
  }

  lf.extension.elk
    .layout(graphData)
    .then(() => {
      // if(graphData.nodes.length > 10){
      //   lf.fitView(32)
      // }

      // lf.translateCenter();

      // 布局完成后提交事务
      if (history) {
        history.commitTransaction('auto-layout')
      }
    })
    .catch((e) => {
      message.error('布局失败，请检查' + e)
      // 如果布局失败，回滚事务
      if (history) {
        history.rollbackTransaction()
      }
    })
}

const handleBatchPaste = ({ originalNodes, basePoint, pasteBasePoint }) => {
  const history = lf.extension.canvasHistory;

  if (history) {
    history.beginTransaction();
  };

  try {
    const nodesToCreate = originalNodes.map(nodeData => {
      const deltaX = nodeData.x - basePoint.x;
      const deltaY = nodeData.y - basePoint.y;

      // 关键：这里要创建一个新的对象，而不是修改剪贴板里的原始数据
      const newNodeData = JSON.parse(JSON.stringify(nodeData));
      newNodeData.x = pasteBasePoint.x + deltaX;
      newNodeData.y = pasteBasePoint.y + deltaY;

      // 调用核心函数生成最终节点信息
      return createNodeInfo({ node: newNodeData, isCopy: true });
    });

    const { nodes: newNodes } = lf.addElements({ nodes: nodesToCreate });
    lf.clearSelectElements();
    newNodes.forEach(node => {
      node.setSelected(true);
    });

    // selectedElements.value = newNodes

  } finally {
    if (history) {
      history.commitTransaction('paste');
    }
  }
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

function isNode(data){
  let node = lf.getNodeModelById(data.id)
  return node
}

function isEdge(data){
  let edge = lf.getEdgeModelById(data.id)
  return edge
}

function getTriggerNodes(){
  const { nodes } = lf.getGraphData()

  let triggerNodes = nodes.filter((node) => node.properties.isTriggerNode)

  return triggerNodes;
}

const handleDeleteNode = (node) => {
  // 将右键点击的节点设置为当前唯一选中的元素
  if (node && node.id) {
    selectedElements.value = [node]
  }

  // 统一调用删除处理函数
  deleteSelectedElements()
}

const handleDeleteEdge = (data) => {
  // 将右键点击的边设置为当前唯一选中的元素
  if (data && data.id) {
    selectedElements.value = [data]
  }

  // 统一调用删除处理函数
  deleteSelectedElements()
}

const handleKeyoardDelete = () => {
  let selectElements = lf.getSelectElements(true)

  selectedElements.value = [...selectElements.nodes, ...selectElements.edges]

  // 统一调用删除处理函数
  deleteSelectedElements()
}

// 处理删除选中元素
const deleteSelectedElements = () => {
  if (selectedElements.value.length === 0) {
    return
  }

  // 过滤掉不能删除的开始节点
  const elementsToDelete = selectedElements.value.filter((el) => {
    if (isNode(el)) {
      return !['start-node', 'group-start-node'].includes(el.type)
    }
    return true // 边可以删除
  })

  if (elementsToDelete.length === 0) {
    selectedElements.value = [] // 清空无效选择
    return
  }

  // 限制删除触发器
  const deletetTriggerNodes = selectedElements.value.filter((el) => el.properties.isTriggerNode === true)

  if (deletetTriggerNodes.length > 0) {
    // 必须保留一个触发器
    const triggerNodes = getTriggerNodes();
    if (triggerNodes.length - deletetTriggerNodes.length < 1) {
      message.error('请至少保留一个触发器！');
      return;
    }
  }

  // 定义核心删除逻辑，用于复用
  const performDelete = () => {
    const history = lf.extension.canvasHistory

    // 1. 操作前：生成一个 "即将删除" 的快照，并用它替换掉最后一条历史记录
    if (history) {
      history.replaceLastState('replace-before-delete')
    }

    // 2. 执行所有删除操作
    let isTriggerNode = false;

    elementsToDelete.forEach((el) => {
      if (isNode(el)) {
        if(el.properties.isTriggerNode){
          isTriggerNode = true;
        }

        lf.deleteNode(el.id)
      } else if (isEdge(el)) {
        lf.deleteEdge(el.id)
      }
    })

    // 更新一下触发器的数据
    if (isTriggerNode) {
      lf.graphModel.eventCenter.emit("custom:trigger-change", null);
    }

    // 3. 操作后：保存一次最终状态
    if (history) {
      history.saveCurrentState('delete')
    }

    // 4. 清空选中状态
    selectedElements.value = []
    lf.clearSelectElements()
  }

  // 判断是否需要弹窗
  // 条件：当选中元素大于1个，或者选中元素为1个且是节点时，需要弹窗确认
  if (dontShowDeleteConfirm.value) {
    performDelete();
    return;
  }

  if (elementsToDelete.length > 1 || (elementsToDelete.length === 1 && isNode(elementsToDelete[0]))) {
    const title = elementsToDelete.length > 1 ? '批量删除' : '删除节点'
    const content = `确定要删除选中的 ${elementsToDelete.length} 个元素吗？`
    let checkboxChecked = false

    Modal.confirm({
      title: title,
      content: h('div', {}, [
        h('p', content),
        h(Checkbox, {
          class: 'dont-show-again-checkbox',
          onChange: (e) => {
            checkboxChecked = e.target.checked
          }
        },
        () => '不再提示')
      ]),
      onOk: () => {
        if(checkboxChecked){
          dontShowDeleteConfirm.value = true;
          setStorage(DONT_SHOW_DELETE_CONFIRM_KEY, true);
        }
        performDelete()
      },
      onCancel: () => {
        // 如果取消删除，不清空选择，以便用户继续操作
      }
    })
  } else {
    // 如果只选中1个边，则直接删除，不弹窗
    performDelete()
  }
}

const noShowDrawerNode = ['group-start-node', 'terminate-node']
// 选择节点
const handleSelectedNode = (data) => {
  const node = JSON.parse(JSON.stringify(data))

  node.properties.dataRaw =  node.properties.dataRaw || node.properties.node_params

  emit('selectNode', node)

  // 结束节点不支持编辑
  if (noShowDrawerNode.includes(data.type)) {
    return
  }

  selectedNode.value = node
  nodeFormDrawerShow.value = true
}

const handleNodeChange = (data) => {
  selectedNode.value.properties = data

  // 更新节点
  updateNode(JSON.parse(JSON.stringify(selectedNode.value)))
}

const handleChangeNodeName = (node_name) => {
  selectedNode.value.properties.node_name = node_name;
  // 先更新数据
  updateNode(JSON.parse(JSON.stringify(selectedNode.value)))
  // 在发送事件之前，确保数据已经更新
  lf.graphModel.eventCenter.emit('custom:setNodeName',  {
    node_name: node_name,
    node_id: selectedNode.value.id,
    node_type: selectedNode.value.type
  })
}

const onDeleteNode = (data) => {
  emit('onDeleteNode', data)

  if (selectedNode.value && data.id === selectedNode.value.id) {
    nodeFormDrawerShow.value = false
    setTimeout(() => {
      selectedNode.value = null
    }, 350)
  }
}

const onDeleteEdge = (data) => {
  emit('onDeleteEdge', data)
}

const handleRunTest = () => {
  emit('runTest')
}

const handleBlankClick = () => {
  emit('blankClick')
  if (nodeFormDrawerShow.value) {
    nodeFormDrawerShow.value = false
  }
  // 清空“选中元素”
  if (selectedElements.value.length > 0) {
    selectedElements.value = []
    lf.clearSelectElements()
  }
}

onMounted(() => {
  initLogicFlow()
})

onUnmounted(() => {
  if (lf && lf.destroy) {
    lf.destroy()
  }
})

defineExpose({
  lfRef,
  setData,
  getData,
  updateNode
})
</script>
