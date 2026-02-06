import { getNodeTypes, getNodesMap } from './node-list.js'
import { useI18n } from '@/hooks/web/useI18n'

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
  const { t } = useI18n('views.workflow.components.util')
  let options = [
    { value: 1, label: t('opt_history') },
    { value: 2, label: t('opt_last_30_days') },
    { value: 3, label: t('opt_last_60_days') },
    { value: 4, label: t('opt_last_90_days') }
  ]

  return options
}

export function getSystemVariable() {
  const { t } = useI18n('views.workflow.components.util')
  return [
    {
      label: t('var_user_message'),
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

export const haveOutKeyNode = ['http-node', 'http-tool-node', 'code-run-node']

export const allParentVariableNodeWhiteList = [
  'start-node',
  'http-node',
  'http-tool-node',
  'parameter-extraction-node',
  'knowledge-base-node',
  'ai-dialogue-node',
  'specify-reply-node',
  'immediately-reply-node',
  'problem-optimization-node',
  'select-data-node',
  'code-run-node',
  'mcp-node',
  'custom-group',
  'zm-plugins-node',
  'zm-workflow-node',
  'batch-group',
  'voice-synthesis-node',
  'voice-clone-node',
  'image-generation-node',
  'import-library-node',
  'json-node',
  'json-reverse-node',
  'qa-node',
]

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
  const { t } = useI18n('views.workflow.components.util')
  return [
    {
      label: t('size_auto_2k'),
      value: '2K'
    },
    {
      label: t('size_auto_4k'),
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
      label: '3:4(1728x2304)',
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

export const specialNodeList = [
  'special.lib_paragraph_list',
  'special.llm_reply_content',
  'specify-reply-node',
  'immediately-reply-node'
]

export const formatSpacialKey = (value) => {
  if(typeof value != 'string'){
    return []
  }

  let result = []
  let specialKey = ''
  for (let i = 0; i < specialNodeList.length; i++) {
    if (value.indexOf(specialNodeList[i]) > -1) {
      specialKey = specialNodeList[i]
      break
    }
  }

  if (specialKey != '') {
    let arr = value.split('.')
    result = [arr[0], specialKey]
  } else {
    result= value.split('.')
  }
  return result
}

export function formatTime(milliseconds) {
  if(!milliseconds){
    return 0
  }
  if (milliseconds < 1000) {
    // 小于1秒，返回毫秒
    return `${milliseconds}ms`;
  } else if (milliseconds < 60000) {
    // 大于等于1秒但小于1分钟，返回秒和毫秒
    const seconds = Math.floor(milliseconds / 1000);
    const remainingMs = milliseconds % 1000;
    if (remainingMs > 0) {
      return `${seconds}s${remainingMs}ms`;
    } else {
      return `${seconds}s`;
    }
  } else {
    // 大于等于1分钟，返回分、秒和毫秒
    const minutes = Math.floor(milliseconds / 60000);
    const remainingMsAfterMinutes = milliseconds % 60000;
    const seconds = Math.floor(remainingMsAfterMinutes / 1000);
    const remainingMs = remainingMsAfterMinutes % 1000;
    
    let result = `${minutes}m`;
    if (seconds > 0 || remainingMs > 0) {
      if (seconds > 0) {
        result += `${seconds}s`;
      }
      if (remainingMs > 0) {
        result += `${remainingMs}ms`;
      }
    }
    
    return result;
  }
}
