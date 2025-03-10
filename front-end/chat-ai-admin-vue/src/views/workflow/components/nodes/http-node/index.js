import { h as flh } from '@logicflow/core'
import { BaseVueNodeView, BaseVueNodeModel } from '../base-node.js'
import VueNode from './index.vue'
import { getMessageNodeAnchor } from '../../util.js'

class VueHtmlNode extends BaseVueNodeView {
  constructor(props) {
    super(props)

    this.renderVue(VueNode, props)
  }
}

class VueHtmlNodeModel extends BaseVueNodeModel {
  setAttributes() {
    const sourceRule01 = {
      message: '只允许从右边的锚点连出',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.id == targetNode.id) {
          return false
        }

        const sourceAnchorIds = sourceNode.graphModel.edges.map(edge => edge.sourceAnchorId)
        const sourceAnchorId = sourceAnchor.id

        return sourceAnchor.type === 'right' && !sourceAnchorIds.includes(sourceAnchorId)
      },
    }
    // sourceRules - 当节点作为边的起始节点（source）时的校验规则
    this.sourceRules = [sourceRule01]

    const targetRule01 = {
      message: '只允许连接左边的锚点',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.id == targetNode.id) {
          return false
        }

        return targetAnchor.type === 'left'
      },
    }
    // targetRules - 当节点作为边的目标节点（target）时的校验规则
    this.targetRules = [targetRule01]
  }

  getAnchorPos() {
    let anchorListPos = getMessageNodeAnchor(this)

    return anchorListPos
  }

  getDefaultAnchor() {
    // 定义锚点位置
    const { id, x, y, width, height, isHovered, isSelected } = this
    const { nodeSortKey } = this.properties

    let defaultAnchor = [
      {
        x: x - width / 2,
        y: y - height / 2 + 24,
        id: nodeSortKey + '-anchor_left',
        type: 'left',
        nodeId: id,
      },
      {
        x: x + width / 2,
        y: y - height / 2 + 24,
        id: nodeSortKey + '-anchor_right',
        type: 'right',
        nodeId: id,
      },
    ]

    let anchorListPos = this.getAnchorPos()

    for (let i = 0; i < anchorListPos.length; i++) {
      let { offsetTop, offsetHeight } = anchorListPos[i]
      let y0 = y - height / 2

      defaultAnchor.push({
        id: anchorListPos[i].id,
        x: x + width / 2,
        y: y0 + offsetTop + offsetHeight / 2,
        type: 'right',
      })
    }

    return defaultAnchor
  }
}

export default {
  type: 'http-node',
  model: VueHtmlNodeModel,
  view: VueHtmlNode,
}
