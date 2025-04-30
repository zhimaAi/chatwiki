export function getNodeIconName(key) {
  let icons = {
    'start-node': 'start-node'
  }

  return icons[key] || key
}

export function getJudgeNodeAnchor(node) {
  let baseTop = 84 // 根据设计稿计算得到
  let anchors = []
  let itemLen = 0
  let termLen = 0
  if (node.term && node.term.length) {
    termLen = node.term.length
    node.term.forEach((item, index) => {
      itemLen += item.terms.length
      let needDis = 0
      if (index > 0) {
        let lens = 2
        for (let i = 0; i < index; i++) {
          lens += node.term[i].terms.length
        }
        needDis = 46 * lens + 106 * (index - 1)
        if(index > 3){
          needDis = needDis + 10
        }
      }
      anchors.push({
        id: node.nodeSortKey + '-anchor_' + index,
        offsetHeight: 38,
        offsetTop: baseTop + needDis
      })
    })
  }

  anchors.push({
    id: node.nodeSortKey + '-anchor_right',
    offsetHeight: 38,
    offsetTop: 225 + 42 * itemLen +  116 * (termLen - 1)
  })

  return anchors
}

export function getQuestionNodeAnchor(node) {
  let baseTop = 388 // 根据设计稿计算得到
  if (node.showMoreBtn) {
    baseTop = baseTop + 144
  }
  if (node.categorys && node.categorys.length) {
    return node.categorys.map((item, index) => {
      return {
        id: node.nodeSortKey + '-anchor_' + index,
        offsetHeight: 38,
        offsetTop: baseTop + 36 * index
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
    'specify-reply-node': 420
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
    return 670
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