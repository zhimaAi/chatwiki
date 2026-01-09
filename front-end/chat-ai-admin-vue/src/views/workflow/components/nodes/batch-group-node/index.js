import { h as flh } from '@logicflow/core'
import { dynamicGroup } from '@logicflow/extension'
import { createApp, h, nextTick } from 'vue'
import CustomGroupComponent from './index.vue'

function transformArray(arr, parent) {
  // 使用map处理数组并返回新的数组
  return arr.map((item) => {
    let node_id = parent ? parent.node_id : item.node_id
    let node_name = parent ? parent.node_name : item.node_name
    let node_type = parent ? parent.node_type : item.node_type
    let label = item.name || item.key
    let text = parent ? `${parent.text}.${label}` : label // 路径
    let children = item.children || item.subs || []

    let value = ''
    let original_value = ''

    if (node_type == 1) {
      original_value = `global.${text}`
    } else {
      if (parent) {
        original_value = `${parent.original_value}.${item.key}`
      } else {
        original_value = `${node_id}.${item.key}`
      }
    }

    value = `【${original_value}】`

    let data = {
      label: label,
      value: value,
      original_value: original_value,
      node_id: node_id,
      node_name: node_name,
      node_type: node_type,
      text: text,
      key: item.key,
      id: node_id,
      typ: item.typ,
      children: children || []
    }

    // 递归处理子节点
    if (data.children && data.children.length > 0) {
      data.children = transformArray(data.children, data)
    }

    return data
  })
}

// 自定义分组
class CustomGroup extends dynamicGroup.view {
  constructor(props) {
    super(props)
    this.div = null
    this.app = null
  }

  /**
   * getShape: 在 super.getShape() 的基础上，追加一个 foreignObject。
   */
  getShape() {
    // 1. 保留父类的基础 shape 结构
    const shape = super.getShape()
    const { model } = this.props
    const { width, height, x, y } = model
    this.div = flh('div', {
      xmlns: 'http://www.w3.org/1999/xhtml',
      style: 'width: 100%; height: 100%;'
    })

    // 3. 创建 foreignObject 并通过 push 添加到 shape 的子元素中
    shape.props.children.push(
      flh('foreignObject', {
        // foreignObject 的坐标是相对于分组左上角的，所以用 0, 0
        x: x - width / 2,
        y: y - height / 2,
        width,
        height,
        overflow: 'unset',
        children: [this.div]
      })
    )

    // 异步挂载Vue组件
    this.mountVueComponentOnNextTick(this.props)
    return shape
  }

  /**
   * mountVueComponentOnNextTick: 负责在 DOM 渲染后挂载和更新 Vue 实例。
   */
  async mountVueComponentOnNextTick() {
    // const { model, graphModel } = this.props
    const mountPoint = this.div.__e

    if (!mountPoint) {
      setTimeout(() => this.mountVueComponentOnNextTick(), 100)
      return
    }


    // 如果 Vue 实例不存在，则创建并挂载
    if (!this.app) {
      this.r = h(CustomGroupComponent, {
        isSelected: this.props.model.isSelected,
        isHovered: this.props.model.isHovered,
        properties: this.props.model.getProperties()
      })

      this.app = createApp({
        render: () => this.r,
        provide: () => ({
          getNode: () => this.props.model,
          getGraph: () => this.props.graphModel,
          addNode: (data) => {
            this.props.graphModel.eventCenter.emit('custom:addGroupNode', {
              data,
              group_id: this.props.model.id
            })
          },
          resetSize: () => {
            if (this.r && this.r.el) {
              // 获取高度
              let height = this.r.el.clientHeight
              let width = this.r.el.clientWidth
              this.props.model.properties = {
                ...this.props.model.properties,
                width: width,
                height: height
              }

              this.props.model._width = width
              this.props.model._height = height

              // 视图变化  边的线位置更新
              this.props.model.refreshBranch()
            }
          },
          setData: (data) => {
            data.dataRaw = data.node_params
            nextTick(() => {
              this.props.model.properties = {
                ...this.props.model.properties,
                ...data
              }
              // 获取高度
              let height = null

              if (this.r && this.r.el) {
                height = this.r.el.clientHeight
              }

              this.props.model._height = height || data.height || this.props.model.properties.height
              this.props.model.refreshBranch() // 视图变化  边的线位置更新
            })
          },
          setTitle: (title) => {
            this.props.model.properties.node_name = title
          }
        })
      })
      this.app.mount(mountPoint)
    } else if (this.r && this.r.component) {
      // 如果实例已存在，则只更新 props
      const newProperties = this.props.model.getProperties()
      const oldProperties = this.r.component.props.properties

      this.r.component.props.isSelected = this.props.model.isSelected
      this.r.component.props.isHovered = this.props.model.isHovered

      if (JSON.stringify(newProperties) !== JSON.stringify(oldProperties)) {
        this.r.component.props.properties = newProperties
      }
    }
  }

  getAnchorShape(anchorData) {
    let edges = this.props.graphModel.getNodeEdges(this.props.model.id)
    let sourceAnchorIds = edges.map((edge) => edge.sourceAnchorId)
    let targetAnchorIds = edges.map((edge) => edge.targetAnchorId)
    let anchorIsSelected = false

    const { x, y, type } = anchorData
    const radius = 10 // 圆形的半径

    if (type === 'left') {
      anchorIsSelected = targetAnchorIds.includes(anchorData.id)
    } else {
      anchorIsSelected = sourceAnchorIds.includes(anchorData.id)
    }

    // 圆形的SVG元素
    const circle = flh('circle', {
      className: `custom-anchor-circle`,
      cx: x,
      cy: y,
      r: radius,
      fill: '#fff', // 你可以根据需要设置颜色
      stroke: '#2475FC', // 你可以根据需要设置边框颜色
      'stroke-width': 2 // 你可以根据需要设置边框宽度
    })

    // 向右箭头的SVG路径元素
    const arrowPath = flh(
      'g',
      {
        className: `custom-anchor-arrow`,
        transform: `translate(${x - radius - 2}, ${y - radius - 2})`
      },
      [
        flh('circle', {
          cx: 12,
          cy: 12,
          r: radius,
          fill: '#2475FC', // 你可以根据需要设置颜色
          stroke: '#2475FC', // 你可以根据需要设置边框颜色
          'stroke-width': 1 // 你可以根据需要设置边框宽度
        }),
        flh('path', {
          d: `M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z`,
          fill: '#fff' // 你可以根据需要设置箭头颜色
        })
      ]
    )

    // 加号的SVG路径元素
    const plusPath = flh(
      'g',
      {
        className: `custom-anchor-plus ${anchorIsSelected ? 'anchor-hide' : ''}`,
        transform: `translate(${x - radius - 2}, ${y - radius - 2})`
      },
      [
        flh('circle', {
          cx: 12,
          cy: 12,
          r: radius,
          fill: '#2475FC', // 你可以根据需要设置颜色
          stroke: '#2475FC', // 你可以根据需要设置边框颜色
          'stroke-width': 1 // 你可以根据需要设置边框宽度
        }),
        flh('line', {
          x1: 12,
          y1: 6,
          x2: 12,
          y2: 18,
          strokeWidth: 2.4,
          strokeLinecap: 'round',
          strokeLinejoin: 'round',
          stroke: '#fff' // 你可以根据需要设置箭头颜色
        }),
        flh('line', {
          x1: 6,
          y1: 12,
          x2: 18,
          y2: 12,
          strokeWidth: 2.4,
          strokeLinecap: 'round',
          strokeLinejoin: 'round',
          stroke: '#fff' // 你可以根据需要设置箭头颜色
        })
      ]
    )

    // 创建一个SVG组元素来包含圆形和箭头
    const group = flh(
      'g',
      {
        id: anchorData.id,
        className: `custom-anchor custom-anchor-${type} ${anchorIsSelected ? 'anchor-selected' : 'anchor-not-selected'}`,
        onClick: (e) => {
          e.stopPropagation()

          const { graphModel, model } = this.props

          if (type === 'right' && !anchorIsSelected) {
            graphModel.eventCenter.emit('custom:showPopupMenu', { anchorData, model })
          }
        }
      },
      circle,
      arrowPath,
      plusPath
    )

    return group
  }

  /**
   * 销毁时卸载 Vue 实例
   */
  destroy() {
    if (this.app) {
      this.app.unmount()
      this.app = null
    }
    super.destroy()
  }
}

class CustomGroupModel extends dynamicGroup.model {
  initNodeData(data) {
    super.initNodeData(data)
    this.width = data.width || 600
    this.height = data.height || 420
    this.collapsible = false
    this.isRestrict = false // 限制拖出去
    this.autoResize = false // 跟随节点 调整宽度
    this.radius = 8
    this.allowResize = false
    this.autoExpand = false
    this.transformWithContainer = true
    this.allowRotate = true
    this.stopMoveGraph = false
    this.autoToFront = false;
    this.properties = {
      ...this.properties,
      transformWithContainer: true, 
    }
  }

  // 自定义分组样式
  getNodeStyle() {
    const style = super.getNodeStyle()
    // 将默认背景设置为透明，因为背景由我们的 Vue 组件控制
    // 这可以防止父类矩形和我们的组件背景重叠
    style.stroke = 'none'
    style.strokeWidth = 2
    style.fill = '#F0F2F5'
    style.filter = 'drop-shadow(0 2px 3px rgba(0, 0, 0, 0.2))'
    // 边框也可以在Vue组件里实现，这里设为透明
    return style
  }

  getOutlineStyle() {
    const style = super.getOutlineStyle()
    style.stroke = 'none'
    style.hover.stroke = 'none'

    return style
  }
  getAnchorLineStyle() {
    const style = super.getAnchorLineStyle()
    style.stroke = '#2475FC'
    return style
  }

  getAddableOutlineStyle() {
    const style = super.getAddableOutlineStyle()
    style.stroke = '#52c41a'
    style.strokeDasharray = '4 4'
    style.strokeWidth = 2
    return style
  }

  isAllowAppendIn(nodeData) {
    return !nodeData.properties.disabled_add_group
  }
  setAttributes() {
    const sourceRule01 = {
      message: '只允许从右边的锚点连出',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.id == targetNode.id) {
          return false
        }

        const sourceAnchorIds = sourceNode.graphModel.edges.map((edge) => edge.sourceAnchorId)
        const sourceAnchorId = sourceAnchor.id

        return sourceAnchor.type === 'right' && !sourceAnchorIds.includes(sourceAnchorId)
      }
    }
    const sourceRule02 = {
      message: '分组里面的节点只能连接当前分组里面的节点,分组外面的节点只能连接分组外面的节点',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if(sourceNode.properties.loop_parent_key){
          return sourceNode.properties.loop_parent_key === targetNode.properties.loop_parent_key
        }else{
          if(targetNode.properties.loop_parent_key){
            return false
          }
        }
        return true
      },
    }
    // sourceRules - 当节点作为边的起始节点（source）时的校验规则
    this.sourceRules = [sourceRule01, sourceRule02]

    const targetRule01 = {
      message: '只允许连接左边的锚点',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.id == targetNode.id) {
          return false
        }

        return targetAnchor.type === 'left'
      }
    }
    // targetRules - 当节点作为边的目标节点（target）时的校验规则
    this.targetRules = [targetRule01]
  }
  getDefaultAnchor() {
    // 定义锚点位置
    const { id, x, y, width, height, isHovered, isSelected } = this
    const { nodeSortKey } = this.properties

    let defaultAnchor = [
      {
        x: x - width / 2,
        y: y - height / 2 + 28,
        id: nodeSortKey + '-anchor_left',
        type: 'left',
        nodeId: id
      },
      {
        x: x + width / 2,
        y: y - height / 2 + 28,
        id: nodeSortKey + '-anchor_right',
        type: 'right',
        nodeId: id
      }
    ]

    return defaultAnchor
  }

  getAllParentVariable() {
    const parentNodes = []
    const visited = new Set()
    const edges = this.incoming.edges
    const { nodes } = this.graphModel
    // 节点白名单
    const nodeWhiteList = [
      'start-node',
      'http-node',
      'parameter-extraction-node',
      'knowledge-base-node',
      'ai-dialogue-node',
      'specify-reply-node',
      'problem-optimization-node',
      'select-data-node',
      'code-run-node',
      'mcp-node',
      'custom-group',
      'image-generation-node',
      'zm-plugins-node'
    ]

    let startNode = nodes.find((node) => node.type === 'start-node')
    // 插入起始节点(起始节点必传)
    if (startNode) {
      visited.add(startNode.id)
      parentNodes.push(startNode)
    }

    // 递归函数用于遍历父节点
    const traverseParents = (edges) => {
      for (let i = 0; i < edges.length; i++) {
        let edge = edges[i]
        let node = edge.sourceNode

        if (visited.has(node.id)) continue

        visited.add(node.id)

        if (nodeWhiteList.includes(node.type)) {
          parentNodes.push(node)
        }

        traverseParents(node.incoming.edges)
      }
    }

    // 获取所有父节点
    traverseParents(edges)


    let childernNodes = nodes.filter(item => item.properties.loop_parent_key == this.id)
    childernNodes.forEach(item => {
      parentNodes.push(item)
    })


    // 取出输出的变量
    let variableArr = []

    for (const node of parentNodes) {
      // 如果节点类型既不是http-node也不是start-node，则跳过当前循环
      if (!nodeWhiteList.includes(node.type)) {
        continue
      }

      let node_params = JSON.parse(node.properties.node_params)

      let obj = {
        label: node.properties.node_name,
        value: node.id,
        node_id: node.id,
        node_type: node.properties.node_type,
        typ: 'node',
        children: [],
        loop_parent_key: node.properties.loop_parent_key,
        group_node_key: this.id
      }

      if (node.type === 'http-node') {
        obj.children = node_params.curl.output
      }

      if (node.type === 'code-run-node') {
        obj.children = node_params.code_run.output
      }

      if (node.type === 'parameter-extraction-node') {
        obj.children = node_params.params_extractor.output
      }

      if (node.type === 'start-node') {
        obj.children = [...node_params.start.diy_global, ...node_params.start.sys_global]
      }

      if (node.type === 'custom-group') {
        obj.children = node_params.loop.output
      }

      if (node.type === 'select-data-node') {
        obj.children = [
          {
            key: 'output_list',
            typ: 'array<object>',
            name: 'output_list',
            label: 'output_list'
          },
          {
            key: 'row_num',
            typ: 'integer',
            name: 'row_num',
            label: 'row_num'
          }
        ]
      }

      if (node.type === 'knowledge-base-node') {
        obj.children = [
          {
            key: 'special.lib_paragraph_list',
            typ: 'string',
            name: '知识库引用',
            label: '知识库引用'
          }
        ]
      }

      if (node.type === 'ai-dialogue-node') {
        obj.children = [
          {
            key: 'special.llm_reply_content',
            typ: 'string',
            name: 'AI回复内容',
            label: 'AI回复内容'
          }
        ]
      }

      if (node.type === 'specify-reply-node') {
        obj.children = [
          {
            key: 'special.llm_reply_content',
            typ: 'string',
            name: '消息内容',
            label: '消息内容'
          }
        ]
      }

      if (node.type === 'problem-optimization-node') {
        obj.children = [
          {
            key: 'special.question_optimize_reply_content',
            typ: 'string',
            name: '问题优化结果',
            label: '问题优化结果'
          }
        ]
      }

      if (node.type === 'mcp-node') {
        obj.children = [
          {
            key: 'special.mcp_reply_content',
            typ: 'string',
            name: 'text',
            label: '工具生成的内容'
          }
        ]
      }
      if(node.type === 'image-generation-node'){
        let image_num = node_params.image_generation.image_num 
        if(image_num > 0){
          let list = []
          for (let i = 0; i < +image_num; i++) {
            let letter = String.fromCharCode('a'.charCodeAt(0) + i)
            list.push({
              key: `picture_url_${letter}`,
              typ: 'string',
              name: `picture_url_${letter}`,
              label: `picture_url_${letter}`,
            })
          }
          obj.children = list
        }
      }

      if (node.type === 'zm-plugins-node') {
        let output = node_params.plugin.output_obj ? JSON.parse(JSON.stringify(node_params.plugin.output_obj)) : []
        
        const loop = (arr) => {
          arr.forEach((item) => {
            item.name = item.key || ''
            if (item.subs && item.subs.length > 0) {
              loop(item.subs)
            }
          })
        }

        loop(output)

        obj.children = output || []
      }

      obj.children.forEach((variable) => {
        variable.node_id = node.id
        variable.node_name = node.properties.node_name
        variable.node_type = node.properties.node_type
      })

      obj.children = transformArray(obj.children, null)

      variableArr.push(obj)
    }

    return variableArr
  }
  
  refreshBranch() {
    // 更新节点连接边的path
    this.incoming.edges.forEach((edge) => {
      // 调用自定义的更新方案
      edge.updatePathByAnchor()
    })
    this.outgoing.edges.forEach((edge) => {
      edge.updatePathByAnchor()
    })
  }
}

export default {
  type: 'batch-group',
  view: CustomGroup,
  model: CustomGroupModel
}
