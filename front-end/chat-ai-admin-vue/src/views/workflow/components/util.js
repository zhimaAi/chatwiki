import { getNodeTypes, getNodesMap } from './node-list.js'
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
const nodesMap = getNodesMap()

export function getNodeIconByType(type){
  if (!type) {
    return ''
  }

  let node = nodesMap[type]

  if (!node) {
    return ''
  }

  return node.properties.node_icon;
}

export function getImageUrl(node_type) {

  let type = nodeTypeMaps[node_type]

  if (!type) {
    return ''
  }
  
  return getNodeIconByType(type);
}
