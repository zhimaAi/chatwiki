const defaultRowData = {
  node_key: '',
  next_node_key: ''
}

const getRowData = () => {
  return JSON.parse(JSON.stringify(defaultRowData))
}

function getNodeIconUrl(name) {
  // 请注意，这不包括子目录中的文件
  return new URL(`../../../assets/svg/${name}.svg`, import.meta.url).href
}

export const nodesGroup = [
  {
    key: 'start',
    name: '开始',
    icon: '',
    hidden: true,
  },
  {
    key: 'large-model-capability',
    name: '大模型能力',
    icon: ''
  },
  {
    key: 'knowledge-retrieval',
    name: '知识检索',
    icon: ''
  },
  {
    key: 'external-service',
    name: '外部调用',
    icon: ''
  },
  {
    key: 'processing-logic',
    name: '处理逻辑',
    icon: ''
  },
  {
    key: 'execute-action',
    name: '执行动作',
    icon: ''
  },
  {
    key: 'end',
    name: '结束',
    icon: ''
  },
  {
    key: 'other',
    name: '其他',
    icon: ''
  }
]

export const nodeList = [
  {
    id: '',
    groupKey: 'start',
    type: 'start-node',
    x: 0,
    y: 0,
    width: 568,
    height: 322,
    hidden: true,
    properties: {
      ...getRowData(),
      node_type: 1,
      node_name: '流程开始',
      node_icon: getNodeIconUrl('start-node'),
      node_icon_name: 'start-node',
      node_params: JSON.stringify({
        start: {
          sys_global: [],
          diy_global: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'execute-action',
    type: 'specify-reply-node',
    width: 420,
    height: 312,
    properties: {
      ...getRowData(),
      node_type: 9,
      node_name: '指定回复',
      node_icon: getNodeIconUrl('specify-reply-node'),
      node_icon_name: 'specify-reply-node',
      node_params: JSON.stringify({
        reply: {
          content: '',
          content_tags: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'ai-dialogue-node',
    width: 568,
    height: 684,
    properties: {
      ...getRowData(),
      node_type: 6,
      node_name: 'AI对话',
      node_icon: getNodeIconUrl('ai-dialogue-node'),
      node_icon_name: 'ai-dialogue-node',
      node_params: JSON.stringify({
        llm: {
          model_config_id: void 0,
          use_model: '',
          context_pair: 6,
          temperature: 0.5,
          max_token: 2000,
          prompt: '',
          enable_thinking: false,
          question_value: 'global.question',
          libs_node_key: void 0,
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'end',
    type: 'end-node',
    width: 420,
    height: 82,
    properties: {
      ...getRowData(),
      node_type: 7,
      node_name: '结束流程',
      node_icon: getNodeIconUrl('end-node'),
      node_icon_name: 'end-node',
      node_params: JSON.stringify({})
    }
  },
  {
    id: '',
    groupKey: 'other',
    type: 'explain-node',
    width: 420,
    height: 152,
    properties: {
      ...getRowData(),
      node_type: -1,
      node_name: '注释卡片',
      node_icon: getNodeIconUrl('explain-node'),
      node_icon_name: 'explain-node',
      node_params: JSON.stringify({
        height: 88,
        content: ''
      })
    }
  },
  {
    id: '',
    groupKey: 'processing-logic',
    type: 'variable-assignment-node',
    width: 568,
    height: 170,
    properties: {
      ...getRowData(),
      node_type: 8,
      node_name: '变量赋值',
      node_icon: getNodeIconUrl('variable-assignment-node'),
      node_icon_name: 'variable-assignment-node',
      node_params: JSON.stringify({
        assign: [
          {
            variable: '',
            value: ''
          }
        ]
      })
    }
  },
  {
    id: '',
    groupKey: 'knowledge-retrieval',
    type: 'knowledge-base-node',
    width: 568,
    height: 386,
    properties: {
      ...getRowData(),
      node_type: 5,
      node_name: '检索知识库',
      node_icon: getNodeIconUrl('knowledge-base-node'),
      node_icon_name: 'knowledge-base-node',
      node_params: JSON.stringify({
        libs: {
          library_ids: '',
          search_type: 1,
          top_k: 5,
          similarity: 0.5,
          rerank_status: 0,
          rerank_model_config_id: void 0,
          rerank_use_model: '',
          question_value: 'global.question'
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'question-node',
    width: 568,
    height: 538,
    properties: {
      ...getRowData(),
      node_type: 3,
      node_name: '问题分类',
      node_icon: getNodeIconUrl('question-node'),
      node_icon_name: 'question-node',
      node_params: JSON.stringify({
        cate: {
          model_config_id: void 0,
          use_model: '',
          context_pair: 2,
          temperature: 0.5,
          max_token: 2000,
          prompt: '',
          enable_thinking: false,
          categorys: [
            {
              category: '',
              next_node_key: ''
            }
          ]
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'processing-logic',
    type: 'judge-node',
    width: 668,
    height: 364,
    properties: {
      ...getRowData(),
      node_type: 2,
      node_name: '判断分支',
      node_icon: getNodeIconUrl('judge-node'),
      node_icon_name: 'judge-node',
      node_params: JSON.stringify({
        term: [
          {
            is_or: false,
            terms: [
              {
                variable: '',
                is_mult: false,
                type: '',
                value: ''
              }
            ],
            next_node_key: ''
          }
        ]
      })
    }
  },
  {
    id: '',
    groupKey: 'external-service',
    type: 'http-node',
    width: 568,
    height: 820,
    properties: {
      ...getRowData(),
      node_type: 4,
      node_name: 'http请求',
      node_icon: getNodeIconUrl('http-node'),
      node_icon_name: 'http-node',
      node_params: JSON.stringify({
        curl: {
          method: 'POST',
          rawurl: '',
          headers: [
            {
              key: '',
              value: ''
            }
          ],
          params: [
            {
              key: '',
              value: ''
            }
          ],
          type: 1,
          body: [
            {
              key: '',
              value: ''
            }
          ],
          body_raw: '',
          timeout: 30,
          output: [
            {
              key: '',
              typ: ''
            }
          ]
        }
      })
    }
  },
]

// 获取分组和节点
export const getAllGroupNodes = (type) => {
  const nodesGroupMap = {}

  // 初始化所有组
  nodesGroup.forEach(group => {
    group.nodes = []

    nodesGroupMap[group.key] = group;
  })
  
  // 将节点按groupKey分组
  nodeList.forEach(node => {
    // 当type不等于'node'时，过滤掉explain-node节点
    if (node.type === 'explain-node' && type === 'node') {
      return
    }
    
    if (node.groupKey && nodesGroupMap[node.groupKey]) {
      nodesGroupMap[node.groupKey].nodes.push(node)
    }
  })

  // 过滤掉没有节点的组
  Object.keys(nodesGroupMap).forEach(key => {
    if (nodesGroupMap[key].nodes.length === 0) {
      delete nodesGroupMap[key]
    }
  })

  // 转换成数组
  let nodesGroupArr = Object.values(nodesGroupMap)
  return JSON.parse(JSON.stringify(nodesGroupArr))
}

export const getNodesMap = () => {
  const nodesMap = {}

  nodeList.forEach(node => {
    nodesMap[node.type] = node
  })

  return JSON.parse(JSON.stringify(nodesMap))
}

export const getNodeTypes = () => {
  let nodeTypes = {}

  nodeList.forEach(node => {
    nodeTypes[node.properties.node_type] = node.type
  })

  return nodeTypes  
}