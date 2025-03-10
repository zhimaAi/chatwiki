import { HtmlNode, HtmlNodeModel, h as flh } from '@logicflow/core'
import { createApp, h } from 'vue'
import VueNode from './index.vue'
import { HtmlResize } from '@logicflow/extension'

class VueHtmlNode extends HtmlNode {
  isMounted
  r
  app

  constructor(props) {
    super(props)
    this.isMounted = false

    this.r = h(VueNode, {
      properties: props.model.getProperties(),
    })

    this.app = createApp({
      render: () => this.r,
    })
  }

  /**
   * 1.1.7版本后支持在view中重写锚点形状。
   * 重写锚点新增
   */
  getAnchorShape(anchorData) {
    const { x, y, type } = anchorData
    return flh('circle', {
      cx: x - 5,
      cy: y - 5,
      r: 10,
      className: `custom-anchor ${type === 'left' ? 'incomming-anchor' : 'outgoing-anchor'}`,
    })
  }

  setHtml(rootEl) {
    if (!this.isMounted) {
      this.isMounted = true
      const node = document.createElement('div')
      rootEl.appendChild(node)
      this.app.mount(node)
    } else {
      this.r.component.props.properties = this.props.model.getProperties()
    }
  }
}

class VueHtmlNodeModel extends HtmlNodeModel {
  /**
   * 给model自定义添加字段方法
   */
  addField(item) {
    this.properties.fields.unshift(item)
    this.setAttributes()
    // 为了保持节点顶部位置不变，在节点变化后，对节点进行一个位移,位移距离为添加高度的一半。
    this.move(0, 24 / 2)
    // 更新节点连接边的path
    this.incoming.edges.forEach((egde) => {
      // 调用自定义的更新方案
      egde.updatePathByAnchor()
    })
    this.outgoing.edges.forEach((edge) => {
      // 调用自定义的更新方案
      edge.updatePathByAnchor()
    })
  }

  getOutlineStyle() {
    const style = super.getOutlineStyle()
    style.stroke = 'none'
    style.hover.stroke = 'none'
    return style
  }

  // 如果不用修改锚地形状，可以重写颜色相关样式
  getAnchorStyle(anchorInfo) {
    const style = super.getAnchorStyle()
    if (anchorInfo.type === 'left') {
      style.fill = 'red'
      style.hover.fill = 'transparent'
      style.hover.stroke = 'transpanrent'
      style.className = 'lf-hide-default'
    } else {
      style.fill = 'green'
    }
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

  setAttributes() {
    this.width = 256
    const {
      properties: { fields },
    } = this
    this.height = 204
    // const circleOnlyAsTarget = {
    //   message: '只允许从右边的锚点连出',
    //   validate: (sourceNode, targetNode, sourceAnchor) => {
    //     return sourceAnchor.type === 'right'
    //   },
    // }
    // this.sourceRules.push(circleOnlyAsTarget)
    // this.targetRules.push({
    //   message: '只允许连接左边的锚点',
    //   validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
    //     return targetAnchor.type === 'left'
    //   },
    // })
  }

  getDefaultAnchor() {
    // 定义锚点位置
    const { id, x, y, width, height, isHovered, isSelected } = this

    return [
      {
        x: x - width / 2 + 5,
        y: y - height / 2 + 28,
        id: 'a-left-1',
        type: 'let',
      },
      {
        x: x + width / 2 + 5,
        y: y - height / 2 + 125 + 0 * 50,
        id: 'a1',
        type: 'right',
      },
      {
        x: x + width / 2 + 5,
        y: y - height / 2 + 125 + 1 * 50,
        id: 'a2',
        type: 'right',
      },
    ]
  }
}

export default {
  type: 'custom-vue-node',
  model: VueHtmlNodeModel,
  view: VueHtmlNode,
}
