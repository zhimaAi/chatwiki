import { BaseVueNodeView, BaseVueNodeModel } from '../base-node.js'
import VueNode from './index.vue'

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
      },
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
        y: y - height / 2 + 26,
        id: nodeSortKey + '-anchor_left',
        type: 'left',
        nodeId: id,
      },
      {
        x: x + width / 2,
        y: y - height / 2 + 26,
        id: nodeSortKey + '-anchor_right',
        type: 'right',
        nodeId: id,
      },
    ]

    return defaultAnchor
  }
}

export default {
  type: 'http-tool-node',
  model: VueHtmlNodeModel,
  view: VueHtmlNode,
}
