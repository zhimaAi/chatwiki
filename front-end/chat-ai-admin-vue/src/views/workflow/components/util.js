import { secondsToMinutes } from '@/utils/index'
export function getNodeIconName(key) {
  let icons = {
    'start-node': 'start-node'
  }

  return icons[key] || key
}

export function getNodeHtml(doms, node) {
  if (!doms || doms.length === 0) {
    return 0
  }

  const wrapperStyle = {
    width: getNodeWidth(node),
    paddingTop: 12,
    paddingBottom: 16,
    paddingLeft: 16,
    paddingRight: 16
  }

  let wrapper = document.createElement('div')
  wrapper.style.boxSizing = 'border-box'
  wrapper.style.whiteSpace = 'pre-wrap'
  wrapper.style.position = 'fixed'
  wrapper.style.left = '-9999px'
  wrapper.style.top = '-9999px'
  wrapper.style.background = '#fff'

  for (let key in wrapperStyle) {
    wrapper.style[key] = wrapperStyle[key] + 'px'
  }

  doms.forEach((item) => {
    let dom = document.createElement('div')

    dom.className = item.className || ''
    dom.id = item.id || ''

    if (item.type === 'text') {
      dom.innerHTML = item.content
    }

    for (let key in item.style) {
      dom.style[key] = item.style[key] + 'px'
    }

    wrapper.appendChild(dom)
  })

  return wrapper
}

export function getDomHeight(doms, node) {
  let wrapper = getNodeHtml(doms, node)

  document.body.appendChild(wrapper)

  let nodeHeight = wrapper.offsetHeight

  wrapper.remove()

  return nodeHeight
}

export function getMessageNodeDoms(node) {
  let doms = [
    {
      type: 'text',
      content: node.properties.node_name,
      style: {
        height: 24,
        marginBottom: 4
      }
    }
  ]

  return doms
}

export function getMessageNodeHeight(node) {
  let doms = getMessageNodeDoms(node)

  let height = getDomHeight(doms, node)

  return height
}

export function getMessageNodeAnchor(node) {
  let doms = getMessageNodeDoms(node)
  let wrapper = getNodeHtml(doms, node)
  const anchorList = []

  document.body.appendChild(wrapper)

  let anchorItems = wrapper.querySelectorAll('.msg-anchor-item') || []

  if (anchorItems.length > 0) {
    anchorItems.forEach((el) => {
      let anchor = {
        offsetTop: el.offsetTop,
        offsetHeight: el.offsetHeight,
        id: el.id
      }

      anchorList.push(anchor)
    })
  }

  wrapper.remove()

  return anchorList
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

export function getStartNodeHeight(node) {
  let doms = [
    {
      type: 'div',
      content: '',
      style: {
        height: 190 // 开始节点 固定高度 按设计稿获取 不包含上下padding
      }
    }
  ]

  return getDomHeight(doms, node)
}

export function getQuestionNodeHeight(node) {
  let doms = [
    {
      type: 'text',
      content: node.properties.node_name,
      style: {
        height: 24,
        lineHeight: 24,
        marginBottom: 4
      }
    }
  ]

  let content = ''

  if (node.properties.question_reference && node.properties.question_reference.length > 0) {
    content = node.properties.question_reference[0]
  }

  if (content) {
    doms.push({
      type: 'text',
      content: node.properties.delay_question,
      style: {
        lineHeight: 22,
        marginBottom: 4
      }
    })

    doms.push({
      type: 'text',
      content: content,
      style: {
        lineHeight: 22,
        minHeight: 38,
        paddingTop: 8,
        paddingBottom: 8,
        paddingLeft: 12,
        paddingRight: 12,
        marginTop: 12
      }
    })
  } else {
    doms.push({
      type: 'div',
      content: '',
      style: {
        height: 38
      }
    })
  }

  return getDomHeight(doms, node)
}

export function getKnowledgeBaseNodeHeight(node) {
  let libs = JSON.parse(node.properties.node_params).libs
  let linLen = libs.library_ids.split(',').length
  return 416 + linLen * 72
}

export function getActionNodeHeight(node) {
  let doms = [
    {
      type: 'text',
      content: node.properties.node_name,
      style: {
        height: 24,
        lineHeight: 24,
        marginBottom: 4
      }
    },
    {
      type: 'text',
      content: node.properties.content,
      style: {
        lineHeight: 22,
        minHeight: 22
      }
    }
  ]

  return getDomHeight(doms, node)
}

export function getQaNodeHeight(node) {
  let time = secondsToMinutes(node.properties.no_resp_jump_time)

  let doms = [
    {
      type: 'text',
      content: node.properties.node_name,
      style: {
        height: 24,
        lineHeight: 24,
        marginBottom: 4
      }
    },
    {
      type: 'text',
      content: `访客超过${time}分钟无响应自动跳出`,
      style: {
        lineHeight: 22,
        marginTop: 12,
        paddingTop: 8,
        paddingBottom: 8,
        paddingLeft: 12,
        paddingRight: 12
      }
    },
    ,
    {
      type: 'div',
      content: '',
      style: {
        height: 22,
        lineHeight: 22,
        marginTop: 12
      }
    }
  ]

  return getDomHeight(doms, node)
}

export function getNodeWidth(node) {
  let widths = {
    'start-node': 420,
    'judge-node': 568,
    'http-node': 568,
    'question-node': 568,
    'ai-dialogue-node': 568,
    'knowledge-base-node': 568,
    'end-node': 420
  }

  return widths[node.type] || 568
}

export function getNodeHeight(node) {
  if (node.type == 'start-node') {
    return getStartNodeHeight(node)
  }

  if (node.type == 'message-node') {
    return 800
  }

  if (node.type == 'knowledge-base-node') {
    return getKnowledgeBaseNodeHeight(node)
  }

  if (node.type == 'question-node') {
    return 800
  }

  if (node.type == 'action-node') {
    return getActionNodeHeight(node)
  }

  if (node.type == 'ai-dialogue-node') {
    return 670
  }

  if (node.type == 'end-node') {
    return 86
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
