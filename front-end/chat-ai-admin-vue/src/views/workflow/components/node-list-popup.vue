<style lang="less" scoped>
.node-list-popup {
  .node-list-popup-content {
    width: 302px;
    padding: 2px;
    border-radius: 6px;
    background: #fff;
    box-shadow:
      0 6px 30px 5px #0000000d,
      0 16px 24px 2px #0000000a,
      0 8px 10px -5px #00000014;
  }

  .node-list {
    .node-type {
      height: 30px;
      padding-left: 16px;
      font-size: 12px;
      color: #8c8c8c;
      display: flex;
      align-items: center;
    }
    .node-flex-box {
      height: 32px;
      display: flex;
      align-items: center;
      .node-item {
        height: 100%;
        cursor: pointer;
        width: 50%;
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 0 16px;
        &:hover {
          background: #e4e6eb;
          border-radius: 6px;
        }
        img {
          width: 20px;
          height: 20px;
        }
      }
    }
  }
}
</style>

<template>
  <div class="node-list-popup">
    <div class="node-list-popup-content">
      <div class="node-list">
        <div class="node-type">大模型能力</div>
        <div class="node-flex-box">
          <div class="node-item" @click="handleAddNode('ai-dialogue-node')">
            <img src="@/assets/svg/ai-dialogue-node.svg" alt="" />
            <div class="node-name">AI对话</div>
          </div>
          <div class="node-item" @click="handleAddNode('question-node')">
            <img src="@/assets/svg/question-node.svg" alt="" />
            <div class="node-name">问题分类</div>
          </div>
        </div>
      </div>
      <div class="node-list">
        <div class="node-type">知识检索</div>
        <div class="node-flex-box">
          <div class="node-item" @click="handleAddNode('knowledge-base-node')">
            <img src="@/assets/svg/knowledge-base-node.svg" alt="" />
            <div class="node-name">检索知识库</div>
          </div>
        </div>
      </div>
      <div class="node-list">
        <div class="node-type">外部调用</div>
        <div class="node-flex-box">
          <div class="node-item" @click="handleAddNode('http-node')">
            <img src="@/assets/svg/http-node.svg" alt="" />
            <div class="node-name">Http请求</div>
          </div>
        </div>
      </div>
      <div class="node-list">
        <div class="node-type">处理逻辑</div>
        <div class="node-flex-box">
          <div class="node-item" @click="handleAddNode('judge-node')">
            <img src="@/assets/svg/judge-node.svg" alt="" />
            <div class="node-name">判断分支</div>
          </div>
        </div>
      </div>
      <div class="node-list">
        <div class="node-type">结束</div>
        <div class="node-flex-box">
          <div class="node-item" @click="handleAddNode('end-node')">
            <img src="@/assets/svg/end-node.svg" alt="" />
            <div class="node-name">结束流程</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { generateUniqueId } from '@/utils/index'
const emit = defineEmits(['addNode'])

const defaultRowData = {
  node_key: '',
  next_node_key: ''
}

const getRowData = () => {
  return JSON.parse(JSON.stringify(defaultRowData))
}

const data = {
  'ai-dialogue-node': {
    id: '',
    type: 'ai-dialogue-node',
    x: 600,
    y: 400,
    properties: {
      ...getRowData(),
      node_type: 6,
      node_name: 'AI对话',
      node_icon_name: 'ai-dialogue-node',
      node_params: JSON.stringify({
        llm: {
          model_config_id: void 0,
          use_model: '',
          context_pair: 6,
          temperature: 0.5,
          max_token: 2000,
          prompt: ''
        }
      })
    }
  },

  'end-node': {
    id: '',
    type: 'end-node',
    x: 600,
    y: 400,
    properties: {
      ...getRowData(),
      node_type: 7,
      node_name: '流程结束',
      node_icon_name: 'end-node',
      node_params: JSON.stringify({})
    }
  },
  'knowledge-base-node': {
    id: '',
    type: 'knowledge-base-node',
    x: 600,
    y: 400,
    properties: {
      ...getRowData(),
      node_type: 5,
      node_name: '知识库检索',
      node_icon_name: 'knowledge-base-node',
      node_params: JSON.stringify({
        libs: {
          library_ids: '',
          search_type: 1,
          top_k: 5,
          similarity: 0.5,
          rerank_status: 0,
          rerank_model_config_id: void 0,
          rerank_use_model: ''
        }
      })
    }
  },
  'question-node': {
    id: '',
    type: 'question-node',
    x: 600,
    y: 400,
    properties: {
      ...getRowData(),
      node_type: 3,
      node_name: '问题分类',
      node_icon_name: 'question-node',
      node_params: JSON.stringify({
        cate: {
          model_config_id: void 0,
          use_model: '',
          context_pair: 2,
          temperature: 0.5,
          max_token: 2000,
          prompt: '',
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
  'judge-node': {
    id: '',
    type: 'judge-node',
    x: 600,
    y: 400,
    properties: {
      ...getRowData(),
      node_type: 2,
      node_name: '判断分支',
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
  'http-node': {
    id: '',
    type: 'http-node',
    x: 600,
    y: 400,
    properties: {
      ...getRowData(),
      node_type: 4,
      node_name: 'http请求',
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
  }
}

const handleAddNode = (id) => {
  let node = data[id]

  node.id = generateUniqueId(node.type)
  emit('addNode', node)
}
</script>
