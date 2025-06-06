<style lang="less" scoped>
.workflow-page {
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  overflow: hidden;
  position: relative;

  .page-body {
    flex: 1;
    display: flex;
    flex-flow: row nowrap;
    overflow: hidden;
    background: #f0f2f5;

    .page-left {
      height: 100%;
      padding-bottom: 8px;
    }

    .page-container {
      flex: 1;
      height: 100%;
    }
  }
}

.logic-flow-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>

<template>
  <div class="workflow-page">
    <PageHeader @edit="handleEdit" @save="handleSave" @release="handleRelease" :saveLoading="saveLoading" />
    <div class="page-body">
      <div class="page-left">
        <PageSidebar />
      </div>
      <div class="page-container">
        <WorkflowCanvas
          ref="workflowCanvasRef"
          @selectedNode="handleSelectedNode"
          @deleteNode="handleDeleteNode"
        />
      </div>
    </div>
    <AddRobotAlert ref="addRobotAlertRef" />
  </div>
</template>

<script setup>
import { getNodeList, saveNodes } from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'
import { generateUniqueId } from '@/utils/index'
import { getNodeIconName } from './components/util.js'
import { onMounted, ref, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import PageSidebar from './components/page-sidebar.vue'
import WorkflowCanvas from './components/workflow-canvas.vue'
import PageHeader from './components/page-header.vue'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'

import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { getModelConfigOption } from '@/api/model/index'

const route = useRoute()
const robot_key = ref(route.query.robot_key)

const addRobotAlertRef = ref(null)
const workflowCanvasRef = ref(null)

const robotStore = useRobotStore()

// const nodes = []
const nodeTypes = {
  '-1': 'explain-node',
  1: 'start-node',
  2: 'judge-node',
  3: 'question-node',
  4: 'http-node',
  5: 'knowledge-base-node',
  6: 'ai-dialogue-node',
  7: 'end-node',
  8: 'variable-assignment-node',
  9: 'specify-reply-node'
}

function getNode() {
  getNodeList({
    robot_key: robot_key.value,
    data_type: 1
  }).then((res) => {
    let list = res.data || []
    let nodes = []
    let edges = []

    list.forEach((item) => {
      // 边数据处理
      if (item.node_type == 0) {
        let edge = JSON.parse(item.node_info_json)

        edges.push(edge)
      } else {
        if(item.node_type == 1){
          let node_params = JSON.parse(item.node_params)
          
          if(!node_params.start){
            node_params.start = []
          }

          item.node_params = JSON.stringify(node_params)
        }
        // 节点数据处理
        let node = JSON.parse(item.node_info_json)

        if (item.node_type != 0) {
          node.type = nodeTypes[item.node_type]
          node.id = item.node_key || generateUniqueId(node.type)
          node.x = node.x || 0
          node.y = node.y || 0
        }

        if(item.node_type == -1){
          item.node_params = node.dataRaw
        }

        if (!node.nodeSortKey) {
          item.nodeSortKey = node.id.substring(0, 8) + node.id.substring(node.id.length - 8)
        } else {
          item.nodeSortKey = node.nodeSortKey
        }

        // 设置icon
        let iconKey = nodeTypes[item.node_type]

        item.iconKey = iconKey
        item.node_icon_name = getNodeIconName(iconKey)
        item.dataRaw = node.dataRaw || item.node_params

        // 删除不要的参数
        delete item.node_info_json

        // 设置 properties
        node.properties = item
        nodes.push(node)
      }
    })

    setWorkflowData({ nodes: nodes, edges: edges })
  })
}

const toAddRobot = (val) => {
  // router.push({ name: 'addRobot' })
  addRobotAlertRef.value.open(val)
}

const setWorkflowData = (data) => {
  workflowCanvasRef.value.setData(data)
}

const getCanvasData = () => {
  let data = workflowCanvasRef.value.getData()

  let list = []
  let edgeMap = {}
  // 先处理边数据
  data.edges.forEach((item) => {
    let obj = {
      node_key: item.id,
      node_name: 'edge',
      node_type: 0
    }

    let node_info_json = {
      ...item
    }

    obj.node_info_json = node_info_json

    list.push(obj)

    if (item.sourceAnchorId) {
      edgeMap[item.sourceAnchorId] = item.targetNodeId
    }

    if (item.targetAnchorId) {
      edgeMap[item.targetAnchorId] = item.sourceNodeId
    }
  })

  data.nodes.forEach((item) => {
    let obj = {
      ...item.properties,
      node_type: +item.properties.node_type,
      node_key: item.id
    }

    obj.node_info_json = {
      type: item.type,
      x: item.x,
      y: item.y,
      id: item.id,
      nodeSortKey: obj.nodeSortKey,
      dataRaw: item.properties.node_params,
    }

    // 关联next_node_key
    obj.next_node_key = edgeMap[obj.nodeSortKey + '-anchor_right'] || ''
    obj.prev_node_key = edgeMap[obj.nodeSortKey + '-anchor_left'] || ''

    let node_params = JSON.parse(obj.node_params)

    if (obj.node_type == 2) {
      // 判断分支
      if (node_params.term && node_params.term.length > 0) {
        node_params.term.forEach((msg, index) => {
          let key = obj.nodeSortKey + '-anchor_' + index
          msg.next_node_key = edgeMap[key] || ''
        })
      }
    }

    if (obj.node_type == 3) {
      // 判断分支
      if (node_params.cate && node_params.cate.categorys && node_params.cate.categorys.length > 0) {
        node_params.cate.categorys.forEach((msg, index) => {
          let key = obj.nodeSortKey + '-anchor_' + index
          msg.next_node_key = edgeMap[key] || ''
        })
      }
    }

    // for (let key in node_params) {
      // for (let key2 in node_params[key]) {
      //   node_params[key][key2] = obj[key2] || ''
      // }
    // }

    // obj.node_params = JSON.stringify(node_params)
    obj.node_params = node_params

    delete obj.dataRaw

    list.push(obj)
  })

  return list
}

const handleSave = (type) => {
  let list = getCanvasData()
  saveNodes({
    robot_key: robot_key.value,
    data_type: 1,
    node_list: JSON.stringify(list)
  }).then(() => {
    type == 'handle' && message.success('保存成功')
    robotStore.setDrafSaveTime({
      type,
      time: dayjs().format('MM/DD HH:mm:ss')
    })
  })
}

// 编辑
const handleEdit = () => {
  toAddRobot(1)
}

let timer = setInterval(
  () => {
    if (import.meta.env.PROD) {
      handleSave('automatic')
    }
  },
  1 * 60 * 1000
)

onUnmounted(() => {
  timer && clearInterval(timer)
})

const saveLoading = ref(false)
// 发布机器人
const handleRelease = async () => {
  let list = getCanvasData()

  let errorNodes = []
  for (let i = 0; i < list.length; i++) {
    let node = list[i]
    // 跳过边节点
    if (node.node_type == 0 || node.node_type == -1) {
      // 跳过
      continue
    }

    if (node.node_type == 1) {
      if (node.next_node_key == '') {
        errorNodes.push(node)
      }
    } else if (node.node_type == 7 || node.node_type == 3) {
      if (node.prev_node_key == '') {
        errorNodes.push(node)
      }
    } else {
      if (node.next_node_key == '' || node.prev_node_key == '') {
        errorNodes.push(node)
      }
    }
  }

  if (errorNodes.length > 0) {
    message.error('存在未关联的节点，请先关联')
    return
  }

  // 先保存草稿，再发布
  saveLoading.value = true;
  try {
    message.loading('保存中...')
    await saveNodes({
      robot_key: robot_key.value,
      data_type: 1,
      node_list: JSON.stringify(list)
    })

    const res = await saveNodes({
      robot_key: robot_key.value,
      data_type: 2,
      node_list: JSON.stringify(list)
    })
    setTimeout(()=>{
      message.destroy()
      saveLoading.value = false;
      if (res && res.res == 0) {
        message.success('发布成功')
      }
    },400)
  } catch (e) {
    saveLoading.value = false;
    // message.success('发布失败，请重试')
  }
}

// 选择节点
let selectedNode = ref(null)
const handleSelectedNode = (data) => {
  selectedNode.value = data
  // 结束节点不支持编辑
  if (data.properties.node_sub_type == 51) {
    return
  }

}

// 删除节点
const handleDeleteNode = () => {
  
}

const getModelList = async () => {
  let uniqueArr = (arr, arr1, key) => {
    const keyVals = new Set(arr.map((item) => item.model_define))
    arr1.filter((obj) => {
      let val = obj[key]
      if (keyVals.has(val)) {
        arr.filter((obj1) => {
          if (obj1.model_define == val) {
            obj1.children = removeRepeat(obj1.children, obj.children)
            return false
          }
        })
      }
    })
    return arr
  }
  return getModelConfigOption({
    model_type: 'LLM'
  }).then((res) => {
    let list = res.data || []
    let children = []
    let modelList = []
    modelList = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.llm_model_list.length; i++) {
        const ele = item.model_info.llm_model_list[i]
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define
        })
      }
      return {
        id: item.model_config.id,
        name: item.model_info.model_name,
        model_define: item.model_info.model_define,
        icon: item.model_info.model_icon_url,
        children: children,
        deployment_name: item.model_config.deployment_name
      }
    })

    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    modelList = uniqueArr(duplicateRemoval(modelList, 'model_define'), modelList, 'model_define')
    robotStore.setModelList(modelList)
  })
}

onMounted(async () => {
  await getModelList()
  getNode()
})
</script>
