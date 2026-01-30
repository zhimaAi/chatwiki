import {getTMcpProviders} from "@/api/robot/thirdMcp.js";
import {jsonDecode} from "@/utils/index.js";
import {getInstallPlugins, runPlugin} from "@/api/plugins/index.js";
import {getPluginActionDefaultArguments, pluginHasAction} from "@/constants/plugin.js";

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
    type: 'trigger_1',
    x: -600,
    y: 0,
    width: 420,
    height: 102,
    hidden: true,
    properties: {
      ...getRowData(),
      componentKey: 'session-trigger-node',
      isTriggerNode: true,
      node_type: 1,
      node_name: '',
      node_icon: getNodeIconUrl('session-trigger-node'),
      node_icon_name: 'session-trigger-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        trigger: {
          outputs: [],
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'start',
    type: 'trigger_2',
    x: -600,
    y: 0,
    width: 420,
    height: 102,
    hidden: true,
    properties: {
      ...getRowData(),
      componentKey: 'session-trigger-node',
      isTriggerNode: true,
      node_type: 2,
      node_name: '',
      node_icon: getNodeIconUrl('session-trigger-node'),
      node_icon_name: 'session-trigger-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        trigger: {
          outputs: [],
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'start',
    type: 'trigger_3',
    x: -600,
    y: 0,
    width: 420,
    height: 102,
    hidden: true,
    properties: {
      ...getRowData(),
      componentKey: 'timing-trigger-node',
      isTriggerNode: true,
      node_type: 2,
      node_name: '',
      node_icon: getNodeIconUrl('timing-trigger-node'),
      node_icon_name: 'timing-trigger-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        trigger: {
          outputs: [],
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'start',
    type: 'trigger_4',
    x: -600,
    y: 0,
    width: 420,
    height: 172,
    hidden: true,
    properties: {
      ...getRowData(),
      componentKey: 'official-trigger-node',
      isTriggerNode: true,
      node_type: 2,
      node_name: '',
      node_icon: getNodeIconUrl('official-trigger-node'),
      node_icon_name: 'official-trigger-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        trigger: {
          outputs: [],
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'start',
    type: 'trigger_5',
    x: -600,
    y: 0,
    width: 420,
    height: 102,
    hidden: true,
    properties: {
      ...getRowData(),
      componentKey: 'webhook-trigger-node',
      isTriggerNode: true,
      node_type: 5,
      node_name: '',
      node_icon: getNodeIconUrl('webhook-trigger-node'),
      node_icon_name: 'webhook-trigger-node',
      node_params: JSON.stringify({
        trigger: {
          outputs: [],
        }
      })
    }
  },

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
      node_header_bg_color: 'linear-gradient(180deg, #F0FFF8 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #F0FFF8 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({})
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
      node_type: 31,
      node_name: '批量执行开始',
      node_icon: getNodeIconUrl('start-node'),
      node_icon_name: 'start-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0FFF8 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({})
    }
  },
  {
    id: '',
    groupKey: 'execute-action',
    type: 'specify-reply-node',
    width: 420,
    height: 130,
    properties: {
      ...getRowData(),
      node_type: 9,
      node_name: '指定回复',
      node_icon: getNodeIconUrl('specify-reply-node'),
      node_icon_name: 'specify-reply-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
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
    groupKey: 'execute-action',
    type: 'qa-node',
    width: 420,
    height: 130,
    properties: {
      ...getRowData(),
      node_type: 43,
      node_name: '问答',
      node_icon: getNodeIconUrl('specify-reply-node'),
      node_icon_name: 'specify-reply-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        question: {
          answer_type: 'text',
          answer_text: '',
          reply_content_list: [
            {
              reply_type: 'smartMenu',
              smart_menu: {
                menu_description: '',
                menu_content: []
              }
            }
          ],
          outputs: [
            {
              key: 'question',
              typ: 'string',
              subs: []
            },
            {
              key: 'question_multiple',
              typ: 'array<object>',
              subs: []
            }
          ]
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'execute-action',
    type: 'immediately-reply-node',
    width: 420,
    height: 130,
    properties: {
      ...getRowData(),
      node_type: 42,
      node_name: '立即回复',
      node_icon: getNodeIconUrl('specify-reply-node'),
      node_icon_name: 'specify-reply-node',
      node_header_bg_color: 'linear-gradient(180deg, #FFF7F0 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        immediately_reply: {
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
      node_name: '大模型',
      node_icon: getNodeIconUrl('ai-dialogue-node'),
      node_icon_name: 'ai-dialogue-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #FFF0F7 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        finish:{
          out_type: 'message',
          messages: []
        }
      })
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
      node_header_bg_color: 'linear-gradient(180deg, #FFF0F7 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
    height: 420,
    properties: {
      ...getRowData(),
      node_type: 25,
      node_name: '循环',
      node_icon: getNodeIconUrl('custom-group-node'),
      node_icon_name: 'custom-group-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
    groupKey: 'processing-logic',
    type: 'batch-group',
    width: 600,
    height: 420,
    properties: {
      ...getRowData(),
      node_type: 30,
      node_name: '批量执行',
      node_icon: getNodeIconUrl('batch-group-node'),
      node_icon_name: 'batch-group-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        batch:{
          chan_number: 10,
          max_run_number: 500,
          batch_arrays: [],
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
      node_header_bg_color: 'linear-gradient(180deg, #F0FFFF 2%, rgba(229, 239, 255, 0) 100%)',
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
          recall_neighbor_switch: false,
          recall_neighbor_before_num: 1,
          recall_neighbor_after_num: 1,
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'knowledge-retrieval',
    type: 'import-library-node',
    width: 420,
    height: 160,
    properties: {
      ...getRowData(),
      node_type: 40,
      node_name: '导入知识库',
      node_icon: getNodeIconUrl('knowledge-base-node'),
      node_icon_name: 'knowledge-base-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0FFFF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        library_import: {
          library_group_id: '0',
          library_id: '',
          import_type: 'content',
          normal_url: '',
          normal_title: '',
          normal_content: '',
          normal_url_repeat_op: 'import',
          qa_question: '',
          qa_answer: '',
          qa_images_variable: '',
          qa_similar_question_variable: '',
          qa_repeat_op: 'import',
          outputs: [
            {
              key: 'msg',
              typ: 'string'
            }
          ]
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
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        cate: {
          model_config_id: void 0,
          use_model: '',
          context_pair: 2,
          temperature: 0.5,
          max_token: 2000,
          prompt: '',
          enable_thinking: false,
          question_value: 'global.question',
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
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
    groupKey: 'processing-logic',
    type: 'json-node',
    width: 420,
    height: 152,
    properties: {
      ...getRowData(),
      node_type: 36,
      node_name: 'JSON 序列化',
      node_icon: getNodeIconUrl('json-node'),
      node_icon_name: 'json-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        json_encode:{
          input_variable: '',
          output: [
            {
              key: '',
              typ: 'string',
              subs: []
            }
          ]
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'processing-logic',
    type: 'json-reverse-node',
    width: 420,
    height: 152,
    properties: {
      ...getRowData(),
      node_type: 37,
      node_name: 'JSON 反序列化',
      node_icon: getNodeIconUrl('json-reverse-node'),
      node_icon_name: 'json-reverse-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        json_decode:{
          input_variable: '',
          output: [
            {
              key: '',
              typ: 'object',
              subs: []
            }
          ]
        }
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
      node_header_bg_color: 'linear-gradient(180deg, #F7F0FF 2%, rgba(229, 239, 255, 0) 100%)',
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
          output: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'http-tool',
    type: 'http-tool-node',
    width: 568,
    height: 216,
    properties: {
      ...getRowData(),
      node_type: 45,
      node_name: 'http工具',
      node_icon: getNodeIconUrl('http-node'),
      node_icon_name: 'http-node',
      node_header_bg_color: 'linear-gradient(180deg, #F7F0FF 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #FFF0F7 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #FFF0F7 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #FFF0F7 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #FFF0F7 2%, rgba(229, 239, 255, 0) 100%)',
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
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
    groupKey: 'large-model-capability',
    type: 'image-generation-node',
    width: 420,
    height: 184,
    properties: {
      ...getRowData(),
      node_type: 33,
      node_name: '图片生成',
      node_icon: getNodeIconUrl('image-generation-node'),
      node_icon_name: 'image-generation-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        image_generation: {
          model_config_id: void 0,
          use_model: '',
          size: void 0,
          image_num: '1',
          prompt: '',
          prompt_tags: [],
          input_images: [],
          image_watermark: '1',
          image_optimize_prompt: '1'
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
      node_header_bg_color: 'linear-gradient(180deg, #F7F0FF 2%, rgba(229, 239, 255, 0) 100%)',
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
      width: 320,
      height: 94,
      node_icon: getNodeIconUrl('mcp-node'),
      node_icon_name: 'mcp-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
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
    height: 154,
    properties: {
      ...getRowData(),
      width: 320,
      height: 154,
      node_type: 21,
      node_name: '',
      node_icon: getNodeIconUrl('zm-plugins-node'),
      node_icon_name: 'zm-plugins-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0FFF8 2%, rgba(229, 239, 255, 0) 100%)',
      plugin_name: '',
      node_params: JSON.stringify({
        plugin: {
          name: "",
          type: "",
          params: {},
          tag_map: {}
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'voice-synthesis-node',
    width: 420,
    height: 158,
    properties: {
      ...getRowData(),
      node_type: 38,
      node_name: '语音合成',
      node_icon: getNodeIconUrl('voice-synthesis-node'),
      node_icon_name: 'voice-synthesis-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        text_to_audio: {
          model_config_id: void 0,
          voice_type: 'all',
          arguments: {},
          output: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'large-model-capability',
    type: 'voice-clone-node',
    width: 420,
    height: 158,
    properties: {
      ...getRowData(),
      node_type: 39,
      node_name: '声音复刻',
      node_icon: getNodeIconUrl('voice-clone-node'),
      node_icon_name: 'voice-clone-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0F5FF 2%, rgba(229, 239, 255, 0) 100%)',
      node_params: JSON.stringify({
        voice_clone: {
          model_config_id: void 0,
          arguments: {},
          output: []
        }
      })
    }
  },
  {
    id: '',
    groupKey: 'workflows',
    type: 'zm-workflow-node',
    width: 420,
    height: 154,
    properties: {
      ...getRowData(),
      width: 320,
      height: 154,
      node_type: 41,
      node_name: '',
      node_icon: getNodeIconUrl('zm-workflow-node'),
      node_icon_name: 'zm-workflow-node',
      node_header_bg_color: 'linear-gradient(180deg, #F0FFF8 2%, rgba(229, 239, 255, 0) 100%)',
      plugin_name: '',
      node_params: JSON.stringify({
        workflow: {
          name: "",
          robot_id: 0,
          params: {},
          output: []
        }
      })
    }
  },
]

// 获取分组和节点
export const getAllGroupNodes = ({excludedNodeTypes}) => {
  const nodesGroupMap = {}

  // 初始化所有组
  nodesGroup.forEach(group => {
    group.nodes = []

    nodesGroupMap[group.key] = group;
  })

  // 将节点按groupKey分组
  nodeList.forEach(node => {
    // 过滤掉excludedNodeTypes中的节点类型
    if (excludedNodeTypes.includes(node.type)) return
    // 当type不等于'node'时，过滤掉explain-node节点
    // if (node.type === 'explain-node' && type === 'node') {
    //   return
    // }
    // if (type == 'loop-node'){
    //   if(node.type == 'custom-group' || node.type == 'end-node') {
    //     return
    //   }
    // } else {
    //   if (node.type == 'terminate-node') {
    //     return
    //   }
    // }
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

export const PluginActionMap = {}

export const getPluginActions = (name) => {
  return PluginActionMap[name] || []
}

export const getAllPluginNodes = async () => {
  let {data} = await getInstallPlugins()
  const nodesMap = getNodesMap()
  const nodeCongfig = nodesMap['zm-plugins-node']

  data = Array.isArray(data) ? data : []
  data = data.filter(i => i?.local?.has_loaded)
  await loadPluginActions(data)
  let plugin
  data = data.map(item => {
    plugin = {
      ...item.remote,
      local: item.local,
    }

    const node_header_bg_color = item.node_header_bg_color ? item.node_header_bg_color : nodeCongfig.properties.node_header_bg_color

    return {
      id: '',
      groupKey: 'plugins',
      type: 'zm-plugins-node',
      width: 320,
      height: 154,
      properties: {
        ...getRowData(),
        node_type: 21,
        node_name: plugin.title,
        node_icon: plugin.icon,
        node_icon_name: '',
        node_desc: plugin.description,
        node_header_bg_color: node_header_bg_color,
        multiNode: plugin?.local?.multiNode || false,
        node_params: JSON.stringify({
          plugin: {
            name: plugin.name,
            type: plugin.type,
            multiNode: plugin?.local?.multiNode || false,
            params: {},
            tag_map: {},
          }
        }),
      },
      expand: false,
      plugin_name: plugin.name,
    }
  })
  return data
}

export const loadPluginActions = async (plugins) => {
  let name, actions
  for (let plugin of plugins) {
    name = plugin?.local?.name || ''
    const isMultiNode = plugin?.local?.multiNode || false
    if (pluginHasAction(name) || isMultiNode) {
      await runPlugin({
        name: name,
        action: "default/get-schema",
        params: {}
      }).then(res => {
        actions = res?.data || {}
        PluginActionMap[name] = []
        for (let key in actions) {
          if (actions[key].type == 'node') PluginActionMap[name].push({...actions[key], name: key})
        }
        PluginActionMap[name].sort((a, b) => (+(a.sort_num || 0)) - (+(b.sort_num || 0)))
      })
    }
  }
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
  const nodesMap = getNodesMap()
  const nodeCongfig = nodesMap['mcp-node']
  const node_header_bg_color = mcp.node_header_bg_color ? mcp.node_header_bg_color : nodeCongfig.properties.node_header_bg_color

  return  {
    id: '',
    groupKey: 'mcp-tool',
    type: 'mcp-node',
    width: 320,
    height: 94,
    properties: {
      ...getRowData(),
      node_type: 20,
      node_name: tool.name || mcp.name,
      node_icon: mcp.avatar,
      node_icon_name: '',
      node_header_bg_color: node_header_bg_color,
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

export const addHttpToolNode = (tool) => {
  const nodesMap = getNodesMap()
  const nodeConfig = nodesMap['http-tool-node']
  const node_header_bg_color = nodeConfig.properties.node_header_bg_color
  const width = nodeConfig.width
  const height = nodeConfig.height
  const method = String(tool.method || '').toUpperCase() || 'GET'
  const rawurl = String(tool.url || '').replace(/[`"]/g, '').trim()
  const headers = Array.isArray(tool.headers) ? tool.headers : []
  // 支持 query/params 两种字段名
  const params = Array.isArray(tool.params) ? tool.params : (Array.isArray(tool.query) ? tool.query : [])
  const body = Array.isArray(tool.body) ? tool.body : []
  const timeout = Number(tool.timeout || 60)
  const output = Array.isArray(tool.output) ? tool.output : [
    {
      key: '',
      typ: ''
    }
  ]
  const curl = {
    method,
    rawurl,
    headers,
    params,
    type: body.length ? 1 : 0,
    body,
    body_raw: '',
    timeout,
    output
  }
  if (Array.isArray(tool.http_auth) && tool.http_auth.length > 0) {
    curl.http_auth = tool.http_auth
  }
  if (tool.http_tool_info && Object.keys(tool.http_tool_info).length > 0) {
    curl.http_tool_info = tool.http_tool_info
  }
  return {
    id: '',
    groupKey: 'external-service',
    type: 'http-tool-node',
    width,
    height,
    properties: {
      ...getRowData(),
      node_type: 45,
      node_name: tool.name || nodeConfig.properties.node_name,
      node_icon: tool.avatar || nodeConfig.properties.node_icon,
      node_icon_name: nodeConfig.properties.node_icon_name,
      node_header_bg_color,
      node_params: JSON.stringify({
        curl
      })
    }
  }
}

export const getPluginActionNode = (node, action, actionName) => {
  node = JSON.parse(JSON.stringify(node))
  let params = JSON.parse(node.properties.node_params)
  let pluginName = params.plugin?.name
  let args = getPluginActionDefaultArguments(pluginName, actionName)
  if ((!args || !Object.keys(args)) && action?.params) {
    args = Object.fromEntries(
      Object.keys(action.params).map(key => [key, ''])
    )
  }
  params.plugin.params.arguments = args
  params.plugin.params.business = actionName
  node.properties.node_params = JSON.stringify(params)
  node.properties.node_name = action.title
  node.properties.plugin_name = `${pluginName}.${actionName}`
  node.width = 420
  return node
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

export const createTriggerNode = (item) => {
  const nodesMap = getNodesMap()
  const type = `trigger_${item.trigger_type}`
  const nodeCongfig = nodesMap[type]
  const icon = item.trigger_icon ? item.trigger_icon : nodeCongfig.properties.node_icon
  const node_header_bg_color = item.node_header_bg_color ? item.node_header_bg_color : nodeCongfig.properties.node_header_bg_color

  const node = {
    type: nodeCongfig.properties.componentKey,
    x: 0,
    y: 0,
    id: '',
    width: nodeCongfig.width,
    height: nodeCongfig.height,
    properties: {
      ...nodeCongfig.properties,
      node_icon: icon,
      isTriggerNode: true,
      width: nodeCongfig.width,
      height: nodeCongfig.height,
      node_key: '',
      nodeSortKey: '',
      node_name: item.trigger_name,
      node_header_bg_color: node_header_bg_color,
      node_params: JSON.stringify({
        trigger: {
          ...item
        }
      })
    }
  }

  return JSON.parse(JSON.stringify(node))
}


export const createWorkflowNode = (item) => {
  const nodesMap = getNodesMap()
  const nodeCongfig = nodesMap['zm-workflow-node']
  const icon = item.robot_avatar ? item.robot_avatar : nodeCongfig.properties.node_icon
  const node_header_bg_color = item.node_header_bg_color ? item.node_header_bg_color : nodeCongfig.properties.node_header_bg_color

  const node = {
    type: 'zm-workflow-node',
    x: 0,
    y: 0,
    id: '',
    width: nodeCongfig.width,
    height: nodeCongfig.height,
    properties: {
      ...nodeCongfig.properties,
      node_icon: icon,
      width: nodeCongfig.width,
      height: nodeCongfig.height,
      node_key: '',
      nodeSortKey: '',
      node_name: item.robot_name,
      node_header_bg_color: node_header_bg_color,
      node_params: JSON.stringify({
        workflow: {
          robot_info: item,
          robot_id: Number(item.id),
          output: [
            {
              "sys": false,
              "key": "data",
              "desc": "工作流回复内容",
              "typ": "string"
            }
          ]
        }
      })
    }
  }

  return JSON.parse(JSON.stringify(node))
}
