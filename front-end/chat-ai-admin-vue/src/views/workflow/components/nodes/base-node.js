import { HtmlNode, HtmlNodeModel, h as flh } from '@logicflow/core'
import { createApp, h } from 'vue'
import { generateUniqueId } from '@/utils/index'

function transformArray(arr, parent) {
  // 使用map处理数组并返回新的数组
  return arr.map((item) => {
    let node_id = parent ? parent.node_id : item.node_id
    let node_name = parent ? parent.node_name : item.node_name
    let node_type = parent? parent.node_type : item.node_type
    let label = item.name || item.key
    let text = parent ? `${parent.text}.${label}` : label // 路径
    let children = item.children || item.subs || []

    let value = ''
    let original_value = ''
    
    if(node_type == 1){
      original_value = `global.${text}`
    }else{
      if(parent){
        original_value = `${parent.original_value}.${item.key}`
      }else{
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
      children: children || [],
    }
    
    // 递归处理子节点
    if (data.children && data.children.length > 0) {
      data.children = transformArray(data.children, data)
    }

    return data
  })
}

export class BaseVueNodeView extends HtmlNode {
  constructor(props) {
    super(props)
  }

  renderVue(VueNode, props) {
    this.isMounted = false

    this.r = h(VueNode, {
      properties: props.model.getProperties(),
      isSelected: props.model.isSelected,
      isHovered: props.model.isHovered,
    })
    this.app = createApp({
      render: () => this.r,
      provide: () => ({
        getNode: () => props.model,
        getGraph: () => props.graphModel,
        setData: (data) => {
          props.model.properties = {
            ...props.model.properties,
            ...data,
          }

          props.model._height = data.height || props.model.properties.height

          let height = null
      
          if (this.r && this.r.el){
            height = this.r.el.clientHeight
          }
    
          props.model._height = height || data.height || props.model.properties.height

          props.model.refreshBranch()  // 视图变化  边的线位置更新

          props.graphModel.eventCenter.emit('custom:setData',  props.model)
        },
        setTitle: (title) => {
          props.model.properties.node_name = title;

          props.graphModel.eventCenter.emit('custom:setNodeName',  {
            node_name: title, 
            node_id: props.model.id,
            node_type: props.model.type
          })
        }
      }),
      mounted: () => {},
    })
  }

  shouldUpdate() {
    const data = {
      ...this.props.model.properties,
      isSelected: this.props.model.isSelected,
      isHovered: this.props.model.isHovered,
    }
    if (this.preProperties && this.preProperties === JSON.stringify(data)) return
    this.preProperties = JSON.stringify(data)
    return true
  }

  /**
   * 1.1.7版本后支持在view中重写锚点形状。
   * 重写锚点新增
   */
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
      'stroke-width': 2, // 你可以根据需要设置边框宽度
    })

    // 向右箭头的SVG路径元素
    const arrowPath = flh(
      'g',
      {
        className: `custom-anchor-arrow`,
        transform: `translate(${x - radius - 2}, ${y - radius - 2})`,
      },
      [
        flh('circle', {
          cx: 12,
          cy: 12,
          r: radius,
          fill: '#2475FC', // 你可以根据需要设置颜色
          stroke: '#2475FC', // 你可以根据需要设置边框颜色
          'stroke-width': 1, // 你可以根据需要设置边框宽度
        }),
        flh('path', {
          d: `M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z`,
          fill: '#fff', // 你可以根据需要设置箭头颜色
        }),
      ],
    )

    // 加号的SVG路径元素
    const plusPath = flh(
      'g',
      {
        className: `custom-anchor-plus ${anchorIsSelected ? 'anchor-hide' : ''}`,
        transform: `translate(${x - radius - 2}, ${y - radius - 2})`,
      },
      [
        flh('circle', {
          cx: 12,
          cy: 12,
          r: radius,
          fill: '#2475FC', // 你可以根据需要设置颜色
          stroke: '#2475FC', // 你可以根据需要设置边框颜色
          'stroke-width': 1, // 你可以根据需要设置边框宽度
        }),
        flh('line', {
          x1: 12,
          y1: 6,
          x2: 12,
          y2: 18,
          strokeWidth: 2.4,
          strokeLinecap: 'round',
          strokeLinejoin: 'round',
          stroke: '#fff', // 你可以根据需要设置箭头颜色
        }),
        flh('line', {
          x1: 6,
          y1: 12,
          x2: 18,
          y2: 12,
          strokeWidth: 2.4,
          strokeLinecap: 'round',
          strokeLinejoin: 'round',
          stroke: '#fff', // 你可以根据需要设置箭头颜色
        }),
      ],
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
        },
      },
      circle,
      arrowPath,
      plusPath,
    )

    return group
  }

  setHtml(rootEl) {
    if (!this.isMounted) {
      this.isMounted = true
      const node = document.createElement('div')
      rootEl.appendChild(node)
      this.app.mount(node)
    } else {
      this.r.component.props.isSelected = this.props.model.isSelected
      this.r.component.props.isHovered = this.props.model.isHovered
      this.r.component.props.properties = this.props.model.getProperties()
    }
  }
}

export class BaseVueNodeModel extends HtmlNodeModel {
  createId() {
    return generateUniqueId(this.type);
  }

  initNodeData(data) {
    this.width = data.properties.width
    this.height = data.properties.height

    super.initNodeData(data)
  }

  getAnchorLineStyle() {
    const style = super.getAnchorLineStyle()
    style.stroke = '#2475FC'
    return style
  }

  getOutlineStyle() {
    const style = super.getOutlineStyle()
    style.stroke = 'none'
    style.hover.stroke = 'none'

    return style
  }

  getGlobalVariable() {
    const { nodes } = this.graphModel

    let startNode = nodes.find((node) => node.type === 'start-node')

    if(!startNode){
      return {
        sys_global: [],
        diy_global: [],
        all_global: [],
      }
    }
    
    let node_params = JSON.parse(startNode.properties.node_params)

    let data = {
      sys_global: node_params.start.sys_global,
      diy_global: node_params.start.diy_global,
      all_global: [...node_params.start.sys_global, ...node_params.start.diy_global],
    }
  
    return data
  }

  getAllParentVariable(){
    const parentNodes = [];
    const visited = new Set();
    const edges = this.incoming.edges;
    const { nodes } = this.graphModel
    // 节点白名单
    const nodeWhiteList = ['start-node', 'http-node', 'knowledge-base-node', 'ai-dialogue-node', 'specify-reply-node']

    let startNode = nodes.find((node) => node.type === 'start-node')
    // 插入起始节点(起始节点必传)
    if(startNode){
      visited.add(startNode.id);
      parentNodes.push(startNode)
    }

    // 递归函数用于遍历父节点
    const traverseParents = (edges) => {
      for(let i=0;i<edges.length;i++){
        let edge = edges[i]
        let node = edge.sourceNode

        if (visited.has(node.id)) continue;

        visited.add(node.id);

        if (nodeWhiteList.includes(node.type)){
          parentNodes.push(node);
        };

        traverseParents(node.incoming.edges);
      }
    };
    
    // 获取所有父节点
    traverseParents(edges);

    // 取出输出的变量
    let variableArr = []

    for (const node of parentNodes) {
      // 如果节点类型既不是http-node也不是start-node，则跳过当前循环
      if (!nodeWhiteList.includes(node.type)) {
        continue;
      }

      let node_params = JSON.parse(node.properties.node_params)
  
      let obj = {
        label: node.properties.node_name,
        value: node.id,
        node_id: node.id,
        node_type: node.properties.node_type,
        typ: 'node',
        children: []
      }

      if(node.type === 'http-node'){
        obj.children = node_params.curl.output
      }

      if(node.type === 'start-node'){
        obj.children = [...node_params.start.diy_global, ...node_params.start.sys_global]
      }

      if(node.type === 'knowledge-base-node'){
        obj.children = [{
          key: 'special.lib_paragraph_list',
          typ: 'string',
          name: '知识库引用',
          label: '知识库引用'
        }]
      }

      if(node.type === 'ai-dialogue-node'){
        obj.children = [{
          key: 'special.llm_reply_content',
          typ: 'string',
          name: 'AI回复内容',
          label: 'AI回复内容',
        }]
      }

      if(node.type === 'specify-reply-node'){
        obj.children = [{
          key: 'special.llm_reply_content',
          typ: 'string',
          name: '消息内容',
          label: '消息内容',
        }]
      }
      
      obj.children.forEach(variable => {
        variable.node_id = node.id
        variable.node_name = node.properties.node_name
        variable.node_type = node.properties.node_type
      })

      obj.children = transformArray(obj.children, null)
      
      variableArr.push(obj)
    }

    return variableArr;
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
