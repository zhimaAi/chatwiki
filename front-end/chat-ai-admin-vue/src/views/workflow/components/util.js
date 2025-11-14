import { getNodeTypes } from './node-list.js'
export function getQuestionNodeAnchor(node) {
  if (node.categorys && node.categorys.length) {
    return node.categorys.map((item, index) => {
      return {
        id: node.nodeSortKey + '-anchor_' + index,
      }
    })
  }

  return []
}

export function getKnowledgeBaseNodeHeight(node) {
  let libs = JSON.parse(node.properties.node_params).libs
  let linLen = libs.library_ids.split(',').length
  return 416 + linLen * 72
}


export function getNodeWidth(node) {
  let widths = {
    'start-node': 420,
    'judge-node': 668,
    'http-node': 568,
    'question-node': 568,
    'ai-dialogue-node': 568,
    'knowledge-base-node': 568,
    'end-node': 420,
    'specify-reply-node': 420,
    'code-run-node': 568,
  }

  return widths[node.type] || 568
}

export function getNodeHeight(node) {
  if (node.type == 'start-node') {
    return 274
  }

  if (node.type == 'message-node') {
    return 800
  }

  if (node.type == 'knowledge-base-node') {
    return getKnowledgeBaseNodeHeight(node)
  }

  if (node.type == 'question-node') {
    return 502
  }

  if (node.type == 'ai-dialogue-node') {
    return 684
  }

  if (node.type == 'problem-optimization-node') {
    return 628
  }

  if (node.type == 'parameter-extraction-node') {
    return 645
  }

  if (node.type == 'end-node') {
    return 86
  }

  if (node.type == 'variable-assignment-node') {
    return 130
  }

  if (node.type == 'specify-reply-node') {
    return 312
  }
  if (node.type == 'explain-node') {
    return 152
  }
  return 800
}

export function getTargetUserOptions() {
  let options = [
    { value: 1, label: '历史' },
    { value: 2, label: '近30天' },
    { value: 3, label: '近60天' },
    { value: 4, label: '近90天' }
  ]

  return options
}

export function getSystemVariable() {
  return [
    {
      label: '用户消息',
      value: '【global.question】',
      payload: { typ: 'string' }
    },
    {
      label: 'open_id',
      value: '【global.openid】',
      payload: { typ: 'string' }
    },
  ]
}

export const haveOutKeyNode = ['http-node', 'code-run-node']


const nodeTypeMaps = getNodeTypes()

export function getImageUrl(node_type) {
  let name = nodeTypeMaps[node_type]
  let url = new URL(`../../../assets/svg/${name}.svg`, import.meta.url)
  return url.href
}