// import { BezierEdge, BezierEdgeModel, h } from '@logicflow/core'
// import { BaseEdgeModel, LineEdge, h } from '@logicflow/core'
import { PolylineEdge, PolylineEdgeModel, h } from '@logicflow/core'
import { createApp, h as vh } from 'vue'
import VueNode from './index.vue'
import { generateUniqueId } from '@/utils/index'

const DEFAULT_WIDTH = 70
const DEFAULT_HEIGHT = 36

class CustomEdge extends PolylineEdge {
  constructor(props) {
    super(props)
    this.isMounted = false
    this.r = vh(VueNode, {
      properties: props.model.getProperties(),
      isSelected: props.model.isSelected,
      isHovered: props.model.isHovered,
    })

    this.app = createApp({
      render: () => this.r,
      provide: () => ({
        getModel: () => props.model,
        getGraphModel: () => props.graphModel,
      }),
      mounted: () => {
        // console.log('mounted', this)
      },
    })
  }

  getText() {
    const { model } = this.props

    const { customWidth = DEFAULT_WIDTH, customHeight = DEFAULT_HEIGHT } = model.getProperties()
    const id = model.id
    // const edgeStyle = model.getEdgeStyle()
    const { startPoint, endPoint, arrowConfig, pointsList } = model

    // const lineData = {
    //   x1: startPoint.x,
    //   y1: startPoint.y,
    //   x2: endPoint.x,
    //   y2: endPoint.y,
    // }
    let positionData = {
      x: (startPoint.x + endPoint.x - customWidth) / 2,
      y: (startPoint.y + endPoint.y - customHeight) / 2,
      width: customWidth,
      height: customHeight,
    }

    if (pointsList.length > 2) {
      let index = Math.round(pointsList.length / 2)
      positionData.x = pointsList[index - 1].x + 5
      positionData.y = pointsList[index - 1].y
    }

    const wrapperStyle = {
      width: customWidth,
      height: customHeight,
    }

    setTimeout(() => {
      const dom = document.querySelector('#' + 'line_' + id).querySelector('.vue-app-wrapper')

      if (!this.isMounted) {
        this.isMounted = true
        this.app.mount(dom)
      } else {
        this.r.component.props.isSelected = this.props.model.isSelected
        this.r.component.props.isHovered = this.props.model.isHovered
        this.r.component.props.properties = this.props.model.getProperties()
      }
    }, 0)

    return h('foreignObject', { ...positionData, id: 'line_' + id }, [
      h('div', {
        id,
        style: wrapperStyle,
        className: 'vue-app-wrapper',
      }),
    ])
  }
}

class CustomEdgeModel extends PolylineEdgeModel {
  initEdgeData(data) {
    super.initEdgeData(data)
  }
  createId() {
    return generateUniqueId('customEdge');
  }

  getEdgeStyle() {
    const style = super.getEdgeStyle()
    // svg属性
    if (this.isSelected || this.isHovered) {
      style.stroke = '#2475FC'
    } else {
      style.stroke = '#A1A7B3'
    }

    style.strokeWidth = 2

    return style
  }

  getArrowStyle() {
    const style = super.getArrowStyle()
    style.refX = 0
    return style
  }

  getOutlineStyle() {
    const style = super.getOutlineStyle()
    style.stroke = 'none'
    style.hover.stroke = 'none'
    return style
  }

  /**
   * 重写此方法，使保存数据是能带上锚点数据。
   */
  getData() {
    const data = super.getData()
    data.sourceAnchorId = this.sourceAnchorId
    data.targetAnchorId = this.targetAnchorId
    return data
  }

  /**
   * 给边自定义方案，使其支持基于锚点的位置更新边的路径
   */
  updatePathByAnchor() {
    // TODO
    const sourceNodeModel = this.graphModel.getNodeModelById(this.sourceNodeId)
    const sourceAnchor = sourceNodeModel
      ?.getDefaultAnchor()
      .find((anchor) => anchor.id === this.sourceAnchorId)

    const targetNodeModel = this.graphModel.getNodeModelById(this.targetNodeId)
    const targetAnchor = targetNodeModel
      ?.getDefaultAnchor()
      .find((anchor) => anchor.id === this.targetAnchorId)

    if (sourceAnchor) {
      const startPoint = {
        x: sourceAnchor?.x,
        y: sourceAnchor?.y,
      }
      this.updateStartPoint(startPoint)
    }

    if (targetAnchor) {
      const endPoint = {
        x: targetAnchor?.x,
        y: targetAnchor?.y,
      }

      this.updateEndPoint(endPoint)
    }
    // 这里需要将原有的pointsList设置为空，才能触发bezier的自动计算control点。
    this.pointsList = []
    this.initPoints()
  }
}

export default {
  type: 'custom-edge',
  view: CustomEdge,
  model: CustomEdgeModel,
}
