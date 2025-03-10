import { HtmlNode, HtmlNodeModel, h as flh } from '@logicflow/core'
import { createApp, h } from 'vue'
import { generateUniqueId } from '@/utils/index'

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
          props.model.refreshBranch()  // 视图变化  边的线位置更新
        },
        setTitle: (title) => {
          props.model.properties.node_name = title;
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

    const { x, y, type, id } = anchorData
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

  getAnchorLineStyle(anchorInfo) {
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
