import ELK from 'elkjs'
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
    const originalGraphData = JSON.parse(JSON.stringify(graphData))

    const { nodes, edges } = originalGraphData

    // 1. 创建节点映射，方便通过ID快速查找节点
    const nodeMap = new Map(nodes.map((node) => [node.id, node]))

    const newNodes = nodes.map((item) => {
      let node = {
        id: item.id,
        parentId: item.properties?.loop_parent_key || null
      }

      // 分组节点特殊处理
      if (item.children && item.children.length > 0) {
        // 对于分组节点，设置初始宽高，ELK会根据子节点自动调整
        // node.width = item.properties.width || 800
        // node.height = item.properties.height || 400
        node.layoutOptions = {
          'elk.padding': '[top=112, right=32, bottom=32, left=32]',
          'elk.nodeSize.constraints': 'MINIMUM_SIZE',
          'elk.nodeSize.minimum': `[600, 420]`,
          'elk.spacing.nodeNode': 100, // 设置组内同层节点间距
          'elk.layered.spacing.nodeNodeBetweenLayers': 120, // 设置组内跨层间距
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

    // 5. 构建ELK图数据
    // 保留原始边数据的映射，以便后续恢复锚点信息等
    const edgeMap = new Map()
    edges.forEach(edge => {
      edgeMap.set(edge.id, { ...edge })
    })

    const elkNodes = arrayToTree(newNodes)
    const elkEdges = edges.map((edge) => {
      return {
        id: edge.id,
        sources: [edge.sourceNodeId],
        targets: [edge.targetNodeId]
      }
    })

    const elkGraph = {
      id: 'root',
      layoutOptions: {
        'elk.algorithm': 'layered', // 使用分层算法
        'elk.direction': 'RIGHT', // 布局方向：从左到右
        'elk.layered.nodePlacement.strategy': 'BRANDES_KOEPF',
        'elk.layered.nodePlacement.bk.fixedAlignment': 'TOP',
        'elk.spacing.nodeNode': 200, //节点间距
        'elk.layered.spacing.nodeNodeBetweenLayers': 200, // 层间间距
        // 嵌套节点相关配置
        'elk.layered.considerHierarchy': 'true',
        'elk.spacing.componentComponent': 80 // 组件间距
      },
      children: elkNodes,
      edges: elkEdges
    }

    // 执行ELK布局计算
    const newElkGraph = await this.elk.layout(elkGraph)
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
        if (elkNode.children && elkNode.children.length > 0) {
          updateNodePositions(elkNode.children, elkNode, actualX, actualY)
        }
      }
    }
    
    // 更新所有节点位置
    updateNodePositions(newElkGraph.children)
 
    // 9. 更新边路径
    edges.forEach((edge) => {
      const originalEdge = edgeMap.get(edge.id) // 获取原始边数据
      const newEdge = newElkGraph.edges.find((e) => e.id === edge.id)
      if (newEdge && newEdge.sections && newEdge.sections.length > 0) {
        const section = newEdge.sections[0]
        if (section) {
          // 保留原始的锚点信息
          if (originalEdge) {
            edge.sourceAnchorId = originalEdge.sourceAnchorId
            edge.targetAnchorId = originalEdge.targetAnchorId
          }

          // 设置起始点和结束点
          edge.startPoint = { ...section.startPoint }
          edge.endPoint = { ...section.endPoint }

          // 构建贝塞尔曲线的pointsList
          // LogicFlow Bezier边的pointsList格式为 [起点, 控制点1, 控制点2, 终点]
          const startPoint = section.startPoint
          const endPoint = section.endPoint

          // 计算控制点，使用偏移量来生成平滑的贝塞尔曲线
          // 偏移量可以根据线段长度动态计算
          const dx = endPoint.x - startPoint.x
          const dy = endPoint.y - startPoint.y
          const length = Math.sqrt(dx * dx + dy * dy)

          // 计算偏移量，用于创建平滑的贝塞尔曲线
          const baseOffset = 50
          const scale = 0.3
          let offset = baseOffset + length * scale
          const minOffset = 10
          const maxOffset = 100
          offset = Math.max(minOffset, Math.min(offset, maxOffset))

          // 创建控制点
          // 控制点1：从起点向右偏移（水平方向）
          const controlPoint1 = {
            x: startPoint.x + Math.abs(dx) * 0.3, // 水平方向偏移
            y: startPoint.y // 保持起点的y坐标或根据需要调整
          }

          // 控制点2：从终点向左偏移（水平方向）
          const controlPoint2 = {
            x: endPoint.x - Math.abs(dx) * 0.3, // 水平方向偏移
            y: endPoint.y // 保持终点的y坐标或根据需要调整
          }

          // 如果是垂直方向的线，调整控制点的计算方式
          if (Math.abs(dy) > Math.abs(dx)) {
            // 如果是垂直方向为主
            controlPoint1.x = startPoint.x
            controlPoint1.y = startPoint.y + Math.abs(dy) * 0.3
            controlPoint2.x = endPoint.x
            controlPoint2.y = endPoint.y - Math.abs(dy) * 0.3
          } else {
            // 如果是水平方向为主
            controlPoint1.x = startPoint.x + Math.abs(dx) * 0.3
            controlPoint1.y = startPoint.y
            controlPoint2.x = endPoint.x - Math.abs(dx) * 0.3
            controlPoint2.y = endPoint.y
          }

          // 更新边的属性
          edge.pointsList = [
            { ...startPoint },
            controlPoint1,
            controlPoint2,
            { ...endPoint }
          ]

          // 保存offset到properties中，以便后续使用
          if (!edge.properties) {
            edge.properties = {}
          }
          edge.properties.offset = Math.round(offset)
        }
      } else {
        // 如果ELK没有返回该边的信息（可能由于某些原因边被移除或未正确生成）
        // 保留原始边信息，但确保pointsList被初始化
        if (originalEdge) {
          edge.sourceAnchorId = originalEdge.sourceAnchorId
          edge.targetAnchorId = originalEdge.targetAnchorId
          edge.type = originalEdge.type || 'custom-bezier-edge'

          // 如果原始边有起始点和结束点，但没有pointsList，则初始化
          if (edge.startPoint && edge.endPoint && (!edge.pointsList || edge.pointsList.length === 0)) {
            // 创建默认的贝塞尔曲线控制点
            const dx = edge.endPoint.x - edge.startPoint.x
            const dy = edge.endPoint.y - edge.startPoint.y
            const length = Math.sqrt(dx * dx + dy * dy)

            const baseOffset = 50
            const scale = 0.3
            let offset = baseOffset + length * scale
            const minOffset = 10
            const maxOffset = 100
            offset = Math.max(minOffset, Math.min(offset, maxOffset))

            const controlPoint1 = {
              x: edge.startPoint.x + Math.abs(dx) * 0.3,
              y: edge.startPoint.y
            }

            const controlPoint2 = {
              x: edge.endPoint.x - Math.abs(dx) * 0.3,
              y: edge.endPoint.y
            }

            if (Math.abs(dy) > Math.abs(dx)) {
              controlPoint1.x = edge.startPoint.x
              controlPoint1.y = edge.startPoint.y + Math.abs(dy) * 0.3
              controlPoint2.x = edge.endPoint.x
              controlPoint2.y = edge.endPoint.y - Math.abs(dy) * 0.3
            } else {
              controlPoint1.x = edge.startPoint.x + Math.abs(dx) * 0.3
              controlPoint1.y = edge.startPoint.y
              controlPoint2.x = edge.endPoint.x - Math.abs(dx) * 0.3
              controlPoint2.y = edge.endPoint.y
            }

            edge.pointsList = [
              { ...edge.startPoint },
              controlPoint1,
              controlPoint2,
              { ...edge.endPoint }
            ]

            if (!edge.properties) {
              edge.properties = {}
            }
            edge.properties.offset = Math.round(offset)
          }
        }
      }
    })
    
    this.lf.clearSelectElements()
    this.lf.clearData()
    this.lf.graphModel.graphDataToModel(originalGraphData)

    this.fitView(originalGraphData.nodes)

    originalGraphData.nodes.forEach(node => {
      // 将组的zIndex设为-9999，确保它在所有节点的下方
      if(node.children && node.children.length){
        let groupNode = this.lf.graphModel.getElement(node.id)
        if(groupNode){
          groupNode.setZIndex(-9999)
        }
      }
    })
  }
  /**
 * 计算并应用最佳视图
 * @param {Array} nodes - 所有节点的数组
 */
  fitView(nodes) {
    if (!nodes || nodes.length === 0) {
      // 如果没有节点，可以重置视图到初始状态
      this.lf.resetZoom();
      this.lf.resetTranslate();
      return;
    }

    // 1. 获取画布尺寸
    const { width: viewWidth, height: viewHeight } = this.lf.graphModel;

    // 2. 计算所有节点的包围盒
    let minX = Infinity;
    let minY = Infinity;
    let maxX = -Infinity;
    let maxY = -Infinity;

    nodes.forEach(node => {
      // 节点的坐标是中心点，需要考虑节点的宽高
      const halfWidth = (node.width || node.properties.width) / 2;
      const halfHeight = (node.height || node.properties.height) / 2;
      minX = Math.min(minX, node.x - halfWidth);
      minY = Math.min(minY, node.y - halfHeight);
      maxX = Math.max(maxX, node.x + halfWidth);
      maxY = Math.max(maxY, node.y + halfHeight);
    });

    const graphWidth = maxX - minX;
    const graphHeight = maxY - minY;

    if (graphWidth === 0 || graphHeight === 0) return;

    // 3. 计算缩放比例，并留出一些边距 (e.g., 10%)
    const padding = 0.1;
    const scaleX = viewWidth / (graphWidth * (1 + padding));
    const scaleY = viewHeight / (graphHeight * (1 + padding));
    let newZoom = Math.min(scaleX, scaleY);

    // 4. 应用缩放
    // LogicFlow 的 zoom API 会以画布中心为基点进行缩放
    if(newZoom > 1){
      newZoom = 1
    }else if(newZoom < 0.3){
      newZoom = 0.3
    }

    this.lf.zoom(newZoom);
    this.lf.translateCenter();
  }

}

export { Elk }
