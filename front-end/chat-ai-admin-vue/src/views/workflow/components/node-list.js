import {getTMcpProviders} from "@/api/robot/thirdMcp.js";
import {jsonDecode} from "@/utils/index.js";
import {getInstallPlugins} from "@/api/plugins/index.js";

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
    key: 'database-operation',
    name: '数据库操作',
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
    width: 420,
    height: 102,
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
    groupKey: 'start',
    type: 'group-start-node',
    x: 0,
    y: 0,
    width: 200,
    height: 42,
    hidden: true,
    properties: {
      ...getRowData(),
      node_type: 27,
      node_name: '循环开始',
      node_icon: getNodeIconUrl('start-node'),
      node_icon_name: 'start-node',
      node_params: JSON.stringify({})
    }
  },
  {
    id: '',
    groupKey: 'execute-action',
    type: 'specify-reply-node',
    width: 420,
    height: 94,
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
    width: 420,
    height: 162,
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
          question_value: 'global.question'
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'problem-optimization-node',
    width: 420,
    height: 154,
    properties: {
      ...getRowData(),
      node_type: 11,
      node_name: '问题优化',
      node_icon: getNodeIconUrl('problem-optimization-node'),
      node_icon_name: 'problem-optimization-node',
      node_params: JSON.stringify({
        question_optimize: {
          model_config_id: void 0,
          use_model: '',
          context_pair: 6,
          temperature: 0.5,
          max_token: 2000,
          prompt: '',
          prompt_tags: [],
          enable_thinking: false,
          question_value: 'global.question'
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'end',
    type: 'end-node',
    width: 420,
    height: 94,
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
    groupKey: 'end',
    type: 'terminate-node',
    width: 420,
    height: 94,
    properties: {
      ...getRowData(),
      node_type: 26,
      node_name: '终止循环',
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
    height: 160,
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
    width: 420,
    height: 94,
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
    groupKey: 'processing-logic',
    type: 'custom-group',
    width: 600,
    height: 300,
    properties: {
      ...getRowData(),
      node_type: 25,
      node_name: '循环',
      node_icon: getNodeIconUrl('custom-group-node'),
      node_icon_name: 'custom-group-node',
      node_params: JSON.stringify({
        loop:{
          loop_type: 'array',
          loop_arrays: [],
          loop_number: '',
          intermediate_params: [],
          output: [],
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'knowledge-retrieval',
    type: 'knowledge-base-node',
    width: 420,
    height: 160,
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
          question_value: 'global.question',
          libs_node_key: void 0,
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'question-node',
    width: 420,
    height: 184,
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
    width: 420,
    height: 152,
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
    height: 216,
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
  {
    id: '',
    groupKey: 'database-operation',
    type: 'add-data-node',
    width: 420,
    height: 124,
    properties: {
      ...getRowData(),
      node_type: 13,
      node_name: '新增数据',
      node_icon: getNodeIconUrl('data-node'),
      node_icon_name: 'data-node',
      node_params: JSON.stringify({
        form_insert: {
          form_id: null,
          datas: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'database-operation',
    type: 'delete-data-node',
    width: 420,
    height: 124,
    properties: {
      ...getRowData(),
      node_type: 14,
      node_name: '删除数据',
      node_icon: getNodeIconUrl('data-node'),
      node_icon_name: 'data-node',
      node_params: JSON.stringify({
        form_delete: {
          form_name: '',
          form_description: '',
          form_id: '',
          typ: 1,
          where: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'database-operation',
    type: 'update-data-node',
    width: 420,
    height: 154,
    properties: {
      ...getRowData(),
      node_type: 15,
      node_name: '更新数据',
      node_icon: getNodeIconUrl('data-node'),
      node_icon_name: 'data-node',
      node_params: JSON.stringify({
        form_update: {
          form_name: '',
          form_description: '',
          form_id: '',
          typ: 1,
          datas: [],
          where: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'database-operation',
    type: 'select-data-node',
    width: 420,
    height: 154,
    properties: {
      ...getRowData(),
      node_type: 16,
      node_name: '查询数据',
      node_icon: getNodeIconUrl('data-node'),
      node_icon_name: 'data-node',
      node_params: JSON.stringify({
        form_select: {
          form_name: '',
          form_description: '',
          form_id: '',
          typ: 1,
          fields: [],
          where: [],
          order: [],
          limit: 100
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'parameter-extraction-node',
    width: 420,
    height: 158,
    properties: {
      ...getRowData(),
      node_type: 12,
      node_name: '参数提取器',
      node_icon: getNodeIconUrl('parameter-extraction-node'),
      node_icon_name: 'parameter-extraction-node',
      node_params: JSON.stringify({
        params_extractor: {
          model_config_id: void 0,
          use_model: '',
          temperature: 0.5,
          max_token: 2000,
          context_pair: 2,
          prompt: '',
          prompt_tags: [],
          question_value: '',
          enable_thinking: false,
          output: [
            // {
            //   key: '',
            //   typ: 'string',
            //   required: false,
            //   default: '',
            //   enum: '',
            //   subs: []
            // }
          ]
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'external-service',
    type: 'code-run-node',
    width: 420,
    height: 160,
    properties: {
      ...getRowData(),
      node_type: 17,
      node_name: '代码运行',
      node_icon: getNodeIconUrl('code-run-node'),
      node_icon_name: 'code-run-node',
      node_params: JSON.stringify({
        code_run: {
          main_func: '',
          params: [
            {
              field: '',
              variable: ''
            }
          ],
          timeout: 30,
          output: [
            {
              key: '',
              typ: 'string'
            }
          ],
          exception: '',
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'mcp-tool',
    type: 'mcp-node',
    width: 320,
    height: 94,
    properties: {
      ...getRowData(),
      node_type: 20,
      node_icon: getNodeIconUrl('mcp-node'),
      node_icon_name: 'mcp-node',
      node_params: JSON.stringify({
        mcp: {
          provider_id: '',
          tool_name: '',
          arguments: {},
          tag_map: {},
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'plugins',
    type: 'zm-plugins-node',
    width: 320,
    height: 320,
    properties: {
      ...getRowData(),
      node_type: 21,
      node_name: '',
      node_icon: 'https://xkf-upload.oss-cn-hangzhou.aliyuncs.com/dev/p/chat_wiki_plugin/1/0/202511/56/364f850087bd429286e8c1bfa87f2a.jpg',
      node_icon_name: '',
      node_params: JSON.stringify({
        plugin: {
          name: "",
          type: "",
          params: {},
          tag_map: {}
        }
      })
    }
  }
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
    if (type == 'loop-node'){
      if(node.type == 'custom-group' || node.type == 'end-node') {
        return
      }
    } else {
      if (node.type == 'terminate-node') {
        return
      }
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

export const getAllPluginNodes = async () => {
  let {data} = await getInstallPlugins()
  data = Array.isArray(data) ? data : []
  let plugin
  data = data.map(item => {
    plugin = {
      ...item.remote,
      local: item.local,
    }
    return {
      id: '',
      groupKey: 'plugins',
      type: 'zm-plugins-node',
      width: 320,
      height: 320,
      properties: {
        ...getRowData(),
        node_type: 21,
        node_name: plugin.title,
        node_icon: plugin.icon,
        node_icon_name: '',
        node_desc: plugin.description,
        node_params: JSON.stringify({
          plugin: {
            name: plugin.name,
            type: plugin.type,
            params: {},
            tag_map: {},
          }
        }),
      },
    }
  })
  return data
}

export const getAllMcpNodes = async () => {
  const {data} = await getTMcpProviders({has_auth: 1})
  data.forEach(item => {
    item.expand = false
    item.tools = jsonDecode(item.tools, [])
  })
  return data
}

export const getMcpNode = (mcp, tool) => {
  return  {
    id: '',
    groupKey: 'mcp-tool',
    type: 'mcp-node',
    width: 320,
    height: 94,
    properties: {
      ...getRowData(),
      node_type: 20,
      node_name: mcp.name,
      node_icon: mcp.avatar,
      node_icon_name: '',
      node_params: JSON.stringify({
        mcp: {
          provider_id: Number(mcp.id),
          tool_name: tool.name,
          arguments: {},
          tag_map: {},
        }
      })
    }
  }
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
