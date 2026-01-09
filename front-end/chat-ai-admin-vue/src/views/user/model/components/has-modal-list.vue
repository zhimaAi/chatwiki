<template>
  <div class="modal-list-box">
    <div class="alert-box">
      <div class="left-block">
        <KeyOutlined />
        <template v-if="['baai', 'ollama', 'xinference'].includes(config_info.model_define)">
          {{ t('views.user.model.api_endpoint_label') }}<template v-if="config_info.api_endpoint"
            >{{ config_info.api_endpoint.slice(0, 30)
            }}{{ config_info.api_endpoint.length > 30 ? '...' : '' }}</template
          >
        </template>
        <template v-else>
          {{ t('views.user.model.api_key_label') }}<template v-if="config_info.api_key"
            >{{ config_info.api_key.slice(0, 30) }}
            {{ config_info.api_key.length > 30 ? '...' : '' }}</template
          >
        </template>
      </div>
      <div class="btn-box">
        <a-button @click="handleAddModel()">{{ t('views.user.model.add_model_btn') }}</a-button>
        <a-button @click="handleEditConfig">{{ t('views.user.model.edit_config_btn') }}</a-button>
        <a-button @click="handleDelconfig">{{ t('views.user.model.remove_config_btn') }}</a-button>
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
          :title="t('views.user.model.show_model_name')"
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
          v-if="model_type == 'IMAGE'"
          :width="140"
          key="image_sizes_desc"
          :title="t('views.user.model.image_sizes_support')"
          data-index="image_sizes_desc"
        >
        </a-table-column>
        <a-table-column
          v-if="model_type == 'IMAGE'"
          :width="160"
          key="image_max"
          :title="t('views.user.model.max_image_count')"
          data-index="image_max"
        >
        </a-table-column>
        <a-table-column
          v-if="model_type == 'IMAGE'"
          :width="140"
          key="image_watermark"
          :title="t('views.user.model.support_watermark')"
          data-index="image_watermark"
        >
          <template #default="{ record }">
            <span v-if="record.image_watermark == 1">{{ t('views.user.model.supported') }}</span>
            <span v-if="record.image_watermark == 0">{{ t('views.user.model.not_supported') }}</span>
          </template>
        </a-table-column>
        <a-table-column
          v-if="model_type == 'IMAGE'"
          :width="160"
          key="image_optimize_prompt"
          :title="t('views.user.model.support_optimize_prompt')"
          data-index="image_optimize_prompt"
        >
          <template #default="{ record }">
            <span v-if="record.image_optimize_prompt == 1">{{ t('views.user.model.supported') }}</span>
            <span v-if="record.image_optimize_prompt == 0">{{ t('views.user.model.not_supported') }}</span>
          </template>
        </a-table-column>
        <a-table-column
          v-if="model_type == 'LLM'"
          :width="140"
          key="thinking_type"
          :title="t('views.user.model.thinking_type')"
          data-index="thinking_type"
        >
          <template #default="{ record }">
            <span v-if="record.thinking_type == 0">{{ t('views.user.model.not_supported') }}</span>
            <span v-if="record.thinking_type == 1">{{ t('views.user.model.supported') }}</span>
            <span v-if="record.thinking_type == 2">{{ t('views.user.model.optional') }}</span>
          </template>
        </a-table-column>

        <a-table-column
          v-if="model_type == 'LLM'"
          :width="140"
          key="function_call"
          :title="t('views.user.model.tool_call')"
          data-index="function_call"
        >
          <template #default="{ record }">
            <span v-if="record.function_call == 0">{{ t('views.user.model.not_supported') }}</span>
            <span v-if="record.function_call == 1">{{ t('views.user.model.supported') }}</span>
          </template>
        </a-table-column>

        <a-table-column
          :width="120"
          v-if="model_type == 'TEXT EMBEDDING'"
          key="vector_dimension_list"
          :title="t('views.user.model.vector_dimension')"
          data-index="vector_dimension_list"
        />
        <a-table-column
          v-if="model_type != 'RERANK'"
          :width="140"
          key="input_desc"
          :title="t('views.user.model.input')"
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
          :title="t('views.user.model.output')"
          data-index="output_desc"
        >
          <template #default="{ record }">
            {{ record.output_desc }}
          </template>
        </a-table-column>
        <a-table-column key="action" :title="t('views.user.model.operate')" :width="120" fixed="right">
          <template #default="{ record }">
            <a @click="handleAddModel(record)">{{ t('views.user.model.edit') }}</a>
            <a-divider type="vertical" />
            <a @click="hadnleDelModel(record)">{{ t('views.user.model.delete') }}</a>
          </template>
        </a-table-column>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { KeyOutlined } from '@ant-design/icons-vue'
import { getSizeOptions } from '@/views/workflow/components/util.js'

const { t } = useI18n()

const props = defineProps({
  currentModalItem: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

const emit = defineEmits(['addModel', 'delModel', 'editConfig', 'delConfig'])
const sizeOptions = getSizeOptions()
let model_type_maps = {
  LLM: t('views.user.model.llm_model'),
  'TEXT EMBEDDING': t('views.user.model.embedding_model'),
  RERANK: t('views.user.model.rerank_model'),
  IMAGE: t('views.user.model.image_generation_model')
}

let input_map = {
  input_text: t('views.user.model.text'),
  input_voice: t('views.user.model.voice'),
  input_image: t('views.user.model.image'),
  input_video: t('views.user.model.video'),
  input_document: t('views.user.model.document')
}

let output_map = {
  output_text: t('views.user.model.text'),
  output_voice: t('views.user.model.voice'),
  output_image: t('views.user.model.image'),
  output_video: t('views.user.model.video')
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
    return t('views.user.model.endpoint_id')
  }
  if (config_info.value.model_define == 'azure') {
    return t('views.user.model.deployment_name_label')
  }
  return 'model'
})

console.log(config_info.value, '==')

const use_model_configs = computed(() => {
  return props.currentModalItem.use_model_configs || []
})

function getSizeDesc(image_sizes) {
  console.log(image_sizes,'=')
  if (image_sizes) {
    let list = image_sizes.split(',')
    return list.map((item) => sizeOptions.find((it) => it.value == item).label).join(',')
  }
  return ''
}

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
    let image_generation_info = item.image_generation ? JSON.parse(item.image_generation) : {}
    let image_sizes_desc = ''
    if (item.image_generation) {
      image_sizes_desc = getSizeDesc(image_generation_info.image_sizes)
    }

    return {
      ...item,
      input_desc: input_desc.join(','),
      output_desc: output_desc.join(','),
      image_generation_info,
      image_sizes_desc,
      image_max: image_generation_info.image_max,
      image_optimize_prompt: image_generation_info.image_optimize_prompt,
      image_watermark: image_generation_info.image_watermark
    }
  })
})

// LLM  TEXT EMBEDDING RERANK
const typeOptions = computed(() => {
  return [
    {
      label: `${t('views.user.model.llm_model')}（${use_model_configs.value.filter((item) => item.model_type == 'LLM').length}）`,
      value: 'LLM'
    },
    {
      label: `${t('views.user.model.embedding_model')}（${use_model_configs.value.filter((item) => item.model_type == 'TEXT EMBEDDING').length}）`,
      value: 'TEXT EMBEDDING'
    },
    {
      label: `${t('views.user.model.rerank_model')}（${use_model_configs.value.filter((item) => item.model_type == 'RERANK').length}）`,
      value: 'RERANK'
    },
    {
      label: `${t('views.user.model.image_generation_model')}（${use_model_configs.value.filter((item) => item.model_type == 'IMAGE').length}）`,
      value: 'IMAGE'
    },
    {
      label: `${t('views.user.model.tts_model')}（${use_model_configs.value.filter((item) => item.model_type == 'TTS').length}）`,
      value: 'TTS'
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
