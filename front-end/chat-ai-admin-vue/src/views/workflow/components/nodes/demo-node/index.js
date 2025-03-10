import { h as flh } from '@logicflow/core'
import { BaseVueNodeView, BaseVueNodeModel } from '../base-node.js'
import { createApp, h } from 'vue'
import VueNode from './index.vue'

class VueHtmlNode extends BaseVueNodeView {
  constructor(props) {
    super(props)

    this.renderVue(VueNode, props)
  }
}

class VueHtmlNodeModel extends BaseVueNodeModel {
  setAttributes() {
    this.width = 256
    this.height = this.properties.height

    const {
      properties: { fields },
    } = this
    // this.width = 256
    // this.height = 172

    const circleOnlyAsTarget = {
      message: '只允许从右边的锚点连出',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.id == targetNode.id) {
          return false
        }
        // console.log(sourceAnchor, targetAnchor)
        // if(sourceAnchor.id === 'q-left-1'){
        //   return targetAnchor.type === 'right';
        // }else{
        //   return targetAnchor.type === 'left';
        // }
        return sourceAnchor.type === 'right'
      },
    }
    // 当节点作为边的起始节点（source）时的校验规则
    this.sourceRules.push(circleOnlyAsTarget)

    // targetRules - 当节点作为边的目标节点（target）时的校验规则
    this.targetRules.push({
      message: '只允许连接左边的锚点',
      validate: (sourceNode, targetNode, sourceAnchor, targetAnchor) => {
        if (sourceNode.id == targetNode.id) {
          return false
        }
        // if(sourceAnchor.type === 'right'){
        //   return targetAnchor.type === 'left'
        // }

        return targetAnchor.type === 'left'
      },
    })
  }

  getDefaultAnchor() {
    // 定义锚点位置
    const { id, x, y, width, height, isHovered, isSelected } = this

    let anchors = [
      {
        x: x - width / 2 + 5,
        y: y - height / 2 + 28,
        id: 'q-left-1',
        type: 'left',
        nodeId: id,
      },
    ]

    if (this.properties.node_type != 5) {
      anchors.push({
        x: x + width / 2 - 5,
        y: y - height / 2 + 28,
        id: 'q-right-1',
        type: 'right',
        nodeId: id,
      })
    }

    return anchors
  }
}

export default {
  type: 'demo-node',
  model: VueHtmlNodeModel,
  view: VueHtmlNode,
}
