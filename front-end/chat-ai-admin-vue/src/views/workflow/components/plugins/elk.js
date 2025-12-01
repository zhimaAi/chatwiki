import ELK from 'elkjs'

const elk = new ELK()

// 定义一个简单的图形结构
const graph = {
  id: 'root',
  children: [
    {
      id: 'parent',
      children: [
        { id: 'c1', width: 40, height: 40 },
        { id: 'c2', width: 40, height: 40 }
      ],
      width: 200,
      height: 150
    },
    { id: 'n3', width: 50, height: 50 }
  ],
  edges: [{ id: 'e1', sources: ['c1'], targets: ['n3'] }]
}

// 计算布局
elk
  .layout(graph)
  .then((layoutedGraph) => {
    console.log(layoutedGraph)
  })
  .catch(console.error)

/**
 * 将扁平化数组转换为树形结构
 * @param {Array} flatData - 包含 id 和 parentId 的扁平数组
 * @param {string} [idKey='id'] - 唯一标识符的键名
 * @param {string} [parentIdKey='parentId'] - 父节点标识符的键名
 * @param {string} [childrenKey='children'] - 子节点数组的键名
 * @returns {Array} - 树形结构数组
 */
function arrayToTree(flatData, idKey = 'id', parentIdKey = 'parentId', childrenKey = 'children') {
  const tree = []
  const map = new Map()

  // 1. 第一次遍历：将所有节点以 id 为键存入 Map，方便快速查找
  for (const item of flatData) {
    map.set(item[idKey], { ...item, [childrenKey]: [] })
  }

  // 2. 第二次遍历：构建树形结构
  for (const item of flatData) {
    const mappedItem = map.get(item[idKey])
    const parentId = item[parentIdKey]

    // 检查是否存在父节点
    if (parentId !== null && parentId !== undefined && map.has(parentId)) {
      const parentItem = map.get(parentId)
      // 将当前节点添加到其父节点的 children 数组中
      parentItem[childrenKey].push(mappedItem)
    } else {
      // 如果没有父节点，说明是根节点
      tree.push(mappedItem)
    }
  }

  return tree
}

class Elk {
  static pluginName = 'elk'
  constructor({ lf }) {
    this.lf = lf
    this.elk = new ELK()
  }
  /**
   * @param {object} graphData
   * @return {Promise<object>}
   */
  async layout(graphData) {
    console.log(graphData)
    const originalGraphData = JSON.parse(JSON.stringify(graphData))

    const { nodes, edges } = originalGraphData

    // 1. 创建节点映射，方便通过ID快速查找节点
    const nodeMap = new Map(nodes.map((node) => [node.id, node]))

    const newNodes = nodes.map((item) => {
      let node = {
        id: item.id,
        parentId: item.properties?.loop_parent_key || null
      }

      console.log('item', item.type)
      // 分组节点特殊处理
      if (item.type === 'custom-group') {
        // 对于分组节点，设置初始宽高，ELK会根据子节点自动调整
        // node.width = item.properties.width || 800
        // node.height = item.properties.height || 400
        node.layoutOptions = {
          'elk.padding': '[top=126, right=32, bottom=32, left=32]',
          'elk.nodeSize.constraints': 'MINIMUM_SIZE',
          // 'elk.nodeSize.minimum': `[${item.properties.width || 800}, ${item.properties.height || 400}]`,
          'elk.spacing.nodeNode': 100, // 设置组内同层节点间距
          'elk.layered.spacing.nodeNodeBetweenLayers': 120, // 设置组内跨层间距
        }
        console.log('custom-group', node.width, node.height)
      } else if (item.type === 'group-start-node' || item.type === 'group-end-node') {
        // 嵌套节点设置固定大小
        node.width = item.properties.width || 200
        node.height = item.properties.height || 80
        node.layoutOptions = {
          'elk.nodeSize.constraints': 'FIXED'
        }
      } else {
        // 普通节点
        node.width = item.properties.width || 420
        node.height = item.properties.height || 100
        node.layoutOptions = {
          'elk.nodeSize.constraints': 'FIXED'
        }
      }

      return node
    })

    console.log(newNodes)

    // 5. 构建ELK图数据
    const elkNodes = arrayToTree(newNodes)
    const elkEdges = edges.map((edge) => {
      return {
        id: edge.id,
        sources: [edge.sourceNodeId],
        targets: [edge.targetNodeId]
      }
    })

    console.log('elkNodes', elkNodes)

    const elkGraph = {
      id: 'root',
      layoutOptions: {
        'elk.algorithm': 'layered', // 使用分层算法
        'elk.direction': 'RIGHT', // 布局方向：从左到右
        'elk.layered.nodePlacement.bk.fixedAlignment': 'TOP',
        'elk.layered.nodePlacement.strategy': 'BRANDES_KOEPF',
        'elk.spacing.nodeNode': 200, //节点间距
        'elk.layered.spacing.nodeNodeBetweenLayers': 200, // 层间间距
        // 嵌套节点相关配置
        'elk.layered.considerHierarchy': 'true',
        'elk.spacing.componentComponent': 80 // 组件间距
      },
      children: elkNodes,
      edges: elkEdges
    }

    console.log(elkGraph)
    // 执行ELK布局计算
    const newElkGraph = await this.elk.layout(elkGraph)
    console.log(newElkGraph)
    // 更新节点位置的辅助函数（递归处理层次结构）
    function updateNodePositions(elkNodes, parentElkNode = null, parentX = 0, parentY = 0) {
      for (const elkNode of elkNodes) {
        const originalNode = nodeMap.get(elkNode.id)
        if (!originalNode) continue

        // 计算实际坐标
        let actualX, actualY

        if (parentElkNode) {
          // 嵌套子节点：坐标相对于父节点
          actualX = parentX + elkNode.x
          actualY = parentY + elkNode.y
        } else {
          // 顶层节点：直接使用绝对坐标
          actualX = elkNode.x
          actualY = elkNode.y
        }

        // 更新节点坐标
        originalNode.x = actualX + elkNode.width / 2;
        originalNode.y = actualY + elkNode.height / 2;

        // 更新宽高（可选）
        if (!originalNode.properties) originalNode.properties = {}

        originalNode.width = elkNode.width
        originalNode.height = elkNode.height
        originalNode.properties.width = elkNode.width
        originalNode.properties.height = elkNode.height

        // 递归子节点，传入当前 elkNode 作为 parent
        if (elkNode.children) {
          updateNodePositions(elkNode.children, elkNode, actualX, actualY)
        }
      }
    }

    // 更新所有节点位置
    updateNodePositions(newElkGraph.children)

    // 9. 更新边路径
    edges.forEach((edge) => {
      const newEdge = newElkGraph.edges.find((e) => e.id === edge.id)
      if (newEdge) {
        let point = newEdge.sections[0]
        if (point) {
          edge.startPoint = ''
          edge.endPoint = ''
          edge.pointsList = []
        }
      }
    })

    console.log(originalGraphData)

    this.lf.clearSelectElements()
    this.lf.clearData()

    this.lf.graphModel.graphDataToModel(originalGraphData)

    return originalGraphData
  }
}

export { Elk }
