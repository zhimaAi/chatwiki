<template>
  <div class="modal-list-box">
    <div class="alert-box">
      <div class="left-block">
        <KeyOutlined />
        <template v-if="['baai', 'ollama', 'xinference'].includes(config_info.model_define)">
          API endpoint：<template v-if="config_info.api_endpoint"
            >{{ config_info.api_endpoint.slice(0, 30)
            }}{{ config_info.api_endpoint.length > 30 ? '...' : '' }}</template
          >
        </template>
        <template v-else>
          API Key：<template v-if="config_info.api_key"
            >{{ config_info.api_key.slice(0, 30) }}
            {{ config_info.api_key.length > 30 ? '...' : '' }}</template
          >
        </template>
      </div>
      <div class="btn-box">
        <a-button @click="handleAddModel()">添加模型</a-button>
        <a-button @click="handleEditConfig">修改配置</a-button>
        <a-button @click="handleDelconfig">移除配置</a-button>
      </div>
    </div>
    <div class="tab-box">
      <a-segmented v-model:value="model_type" :options="typeOptions" />
    </div>
    <div class="table-box">
      <a-table sticky :data-source="tableData" :pagination="false" :scroll="{ x: 1000 }">
        <a-table-column
          :width="140"
          key="show_model_name"
          title="模型名称"
          data-index="show_model_name"
        />
        <a-table-column
          :width="140"
          key="use_model_name"
          :title="model_name"
          data-index="use_model_name"
        >
          <template #default="{ record }">
            {{ record.use_model_name }}
          </template>
        </a-table-column>
        <a-table-column
          v-if="model_type == 'LLM'"
          :width="140"
          key="thinking_type"
          title="深度思考"
          data-index="thinking_type"
        >
          <template #default="{ record }">
            <span v-if="record.thinking_type == 0">不支持</span>
            <span v-if="record.thinking_type == 1">支持</span>
            <span v-if="record.thinking_type == 2">可选</span>
          </template>
        </a-table-column>

        <a-table-column
          v-if="model_type == 'LLM'"
          :width="140"
          key="function_call"
          title="工具调用"
          data-index="function_call"
        >
          <template #default="{ record }">
            <span v-if="record.function_call == 0">不支持</span>
            <span v-if="record.function_call == 1">支持</span>
          </template>
        </a-table-column>

        <a-table-column
          :width="120"
          v-if="model_type == 'TEXT EMBEDDING'"
          key="vector_dimension_list"
          title="向量维度"
          data-index="vector_dimension_list"
        />
        <a-table-column
          v-if="model_type != 'RERANK'"
          :width="140"
          key="input_desc"
          title="输入"
          data-index="input_desc"
        >
          <template #default="{ record }">
            {{ record.input_desc }}
          </template>
        </a-table-column>
        <a-table-column
          v-if="model_type == 'LLM'"
          :width="140"
          key="output_desc"
          title="输出"
          data-index="output_desc"
        >
          <template #default="{ record }">
            {{ record.output_desc }}
          </template>
        </a-table-column>
        <a-table-column key="action" title="操作" :width="120" fixed="right">
          <template #default="{ record }">
            <a @click="handleAddModel(record)">编辑</a>
            <a-divider type="vertical" />
            <a @click="hadnleDelModel(record)">删除</a>
          </template>
        </a-table-column>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import { KeyOutlined } from '@ant-design/icons-vue'
const props = defineProps({
  currentModalItem: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

const emit = defineEmits(['addModel', 'delModel', 'editConfig', 'delConfig'])

let model_type_maps = {
  LLM: '大语言模型',
  'TEXT EMBEDDING': '嵌入模型',
  RERANK: '重排序模型'
}

let input_map = {
  input_text: '文本',
  input_voice: '语音',
  input_image: '图片',
  input_video: '视频',
  input_document: '文档'
}

let output_map = {
  output_text: '文本',
  output_voice: '语音',
  output_image: '图片',
  output_video: '视频'
}

const model_type = ref('LLM')

const config_info = computed(() => {
  return (
    props.currentModalItem.config_info || {
      api_key: ''
    }
  )
})

console.log(config_info.value, '==')

const model_name = computed(() => {
  if (config_info.value.model_define == 'doubao') {
    return '接入点id'
  }
  if (config_info.value.model_define == 'azure') {
    return '部署名称'
  }
  return 'model'
})

console.log(config_info.value, '==')

const use_model_configs = computed(() => {
  return props.currentModalItem.use_model_configs || []
})

const tableData = computed(() => {
  let list = use_model_configs.value.filter((item) => item.model_type == model_type.value)
  return list.map((item) => {
    let input_desc = []
    let output_desc = []
    for (let key in input_map) {
      if (item[key] == 1) {
        input_desc.push(input_map[key])
      }
    }
    for (let key in output_map) {
      if (item[key] == 1) {
        output_desc.push(output_map[key])
      }
    }
    return {
      ...item,
      input_desc: input_desc.join(','),
      output_desc: output_desc.join(',')
    }
  })
})

// LLM  TEXT EMBEDDING RERANK
const typeOptions = computed(() => {
  return [
    {
      label: `大语言模型（${use_model_configs.value.filter((item) => item.model_type == 'LLM').length}）`,
      value: 'LLM'
    },
    {
      label: `嵌入模型（${use_model_configs.value.filter((item) => item.model_type == 'TEXT EMBEDDING').length}）`,
      value: 'TEXT EMBEDDING'
    },
    {
      label: `重排序模型（${use_model_configs.value.filter((item) => item.model_type == 'RERANK').length}）`,
      value: 'RERANK'
    }
  ]
})

const handleAddModel = (record) => {
  emit('addModel', props.currentModalItem, record)
}

const handleEditConfig = () => {
  emit('editConfig', props.currentModalItem, config_info.value)
}
const handleDelconfig = () => {
  emit('delConfig', config_info.value)
}

const hadnleDelModel = (record) => {
  emit('delModel', record)
}
</script>

<style lang="less" scoped>
.modal-list-box {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  .alert-box {
    height: 40px;
    background: var(--01-, #e5efff);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    .left-block {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #7a8699;
      font-size: 14px;
    }
    .btn-box {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
  .tab-box {
    padding: 16px;
    padding-bottom: 0;
    &::v-deep(.ant-segmented) {
      color: #262626;
      .ant-segmented-item-selected {
        color: #2475fc;
      }
    }
  }
  .table-box {
    padding: 16px;
    padding-top: 0;
    margin-top: 16px;
    flex: 1;
    overflow-y: auto;
  }
}
</style>
