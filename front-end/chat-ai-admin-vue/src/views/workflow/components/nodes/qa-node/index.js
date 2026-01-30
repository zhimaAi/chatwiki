import { h as flh } from '@logicflow/core'
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

        const sourceAnchorIds = sourceNode.graphModel.edges.map((edge) => edge.sourceAnchorId)
        const sourceAnchorId = sourceAnchor.id

        return sourceAnchor.type === 'right' && !sourceAnchorIds.includes(sourceAnchorId)
      }
    }
    const sourceRule02 = {
      message: '分组里面的节点只能连接当前分组里面的节点,分组外面的节点只能连接分组外面的节点',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.properties.loop_parent_key) {
          return sourceNode.properties.loop_parent_key === targetNode.properties.loop_parent_key
        } else {
          if (targetNode.properties.loop_parent_key) {
            return false
          }
        }
        return true
      }
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

    // console.log(JSON.parse(this.properties.node_params), '==')
    let question = {}
    try {
      question = JSON.parse(this.properties.dataRaw).question
    } catch (error) {}

    if (question.answer_type == 'menu') {
      let customAnchor = []
      if (question.reply_content_list && question.reply_content_list.length > 0) {
        let menu_content = question.reply_content_list[0].smart_menu?.menu_content || []
        let y0 = y - height / 2
        menu_content.forEach((item, index) => {
          customAnchor.push({
            id: nodeSortKey + '-anchor_' + index,
            x: x + width / 2,
            y: y0 + 132 + index * 32,
            type: 'right',
            nodeId: id
          })
        })
      }
      return [
        {
          x: x - width / 2,
          y: y - height / 2 + 24,
          id: nodeSortKey + '-anchor_left',
          type: 'left',
          nodeId: id
        },
        ...customAnchor
      ]
    } else {
      return [
        {
          x: x - width / 2,
          y: y - height / 2 + 24,
          id: nodeSortKey + '-anchor_left',
          type: 'left',
          nodeId: id
        },
        {
          x: x + width / 2,
          y: y - height / 2 + 24,
          id: nodeSortKey + '-anchor_right',
          type: 'right',
          nodeId: id
        }
      ]
    }
  }
}

export default {
  type: 'qa-node',
  model: VueHtmlNodeModel,
  view: VueHtmlNode
}
