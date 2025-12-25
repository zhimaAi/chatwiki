import { getNodeTypes, getNodesMap } from './node-list.js'
export function getQuestionNodeAnchor(node) {
  if (node.categorys && node.categorys.length) {
    return node.categorys.map((item, index) => {
      return {
        id: node.nodeSortKey + '-anchor_' + index
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
    }
  ]
}

export const haveOutKeyNode = ['http-node', 'code-run-node']

const nodeTypeMaps = getNodeTypes()
const nodesMap = getNodesMap()

export function getNodeIconByType(type) {
  if (!type) {
    return ''
  }

  let node = nodesMap[type]

  if (!node) {
    return ''
  }

  return node.properties.node_icon
}

export function getImageUrl(node_type) {
  let type = nodeTypeMaps[node_type]

  if (!type) {
    return ''
  }

  return getNodeIconByType(type)
}

export function getSizeOptions() {
  return [
    {
      label: '自动适配比例(2k)',
      value: '2K'
    },
    {
      label: '自动适配比例(4k)',
      value: '4K'
    },
    {
      label: '1:1(2048x2048)',
      value: '2048x2048'
    },
    {
      label: '4:3(2304x1728)',
      value: '2304x1728'
    },
    {
      label: '3:4(2304x1728)',
      value: '1728x2304'
    },
    {
      label: '16:9(2560x1440)',
      value: '2560x1440'
    },
    {
      label: '9:16(1440x2560)',
      value: '1440x2560'
    },
    {
      label: '3:2(2496x1664)',
      value: '2496x1664'
    },
    {
      label: '2:3(1664x2496)',
      value: '1664x2496'
    },
    {
      label: '21:9(3024x1296)',
      value: '3024x1296'
    }
  ]
}
