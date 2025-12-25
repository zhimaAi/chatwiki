<template>
  <node-common
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    style="width: 420px"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">模型</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key">
                <model-name-text
                  :useModel="formState.use_model"
                  :modelConfigId="formState.model_config_id"
                />
              </span>
            </div>
          </div>
        </div>
        <div class="field-item">
          <div class="field-item-label">提示词</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key">
                <at-text
                  :options="valueOptions"
                  :defaultSelectedList="formState.prompt_tags"
                  :defaultValue="formState.prompt"
                  ref="atInputRef"
                />
              </span>
            </div>
          </div>
        </div>
        <div class="field-item">
          <div class="field-item-label">输出</div>
          <div class="options-list">
            <div class="options-item" v-for="item in outputList" :key="item">
              <div class="option-label">{{ item }}</div>
              <div class="option-type">string</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { getUuid } from '@/utils/index'
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount, computed } from 'vue'
import { storeToRefs } from 'pinia'
import NodeCommon from '../base-node.vue'
import ModelNameText from '../model-name-text.vue'
import UserQuestionText from '../user-question-text.vue'
import AtText from '../../at-input/at-text.vue'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const getNode = inject('getNode')
const resetSize = inject('resetSize')

const valueOptions = ref([])

// --- State ---
const menus = ref([])
const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  size: void 0,
  image_num: '1',
  prompt: '',
  prompt_tags: [],
  input_images: [],
  image_watermark: '1',
  image_optimize_prompt: '1'
})

const outputList = computed(() => {
  let list = ['msg']
  if (formState.image_num > 0) {
    for (let i = 0; i < +formState.image_num; i++) {
      let letter = String.fromCharCode('a'.charCodeAt(0) + i)
      list.push(`picture_url_${letter}`)
    }
  }
  return list
})

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let image_generation = {}
  try {
    getValueOptions()
    image_generation = JSON.parse(dataRaw).image_generation || {}
    formState.model_config_id = image_generation.model_config_id
    formState.use_model = image_generation.use_model
    formState.size = image_generation.size
    formState.image_num = image_generation.image_num
    formState.prompt = image_generation.prompt
    formState.prompt_tags = image_generation.prompt_tags
    formState.input_images = image_generation.input_images
    formState.image_watermark = image_generation.image_watermark
    formState.image_optimize_prompt = image_generation.image_optimize_prompt
  } catch (e) {
    image_generation = {}
  }

  nextTick(() => {
    resetSize()
  })
}

function getValueOptions() {
  let options = getNode().getAllParentVariable()
  valueOptions.value = options || []
}

watch(
  () => props.properties,
  (newVal, oldVal) => {
    const newDataRaw = newVal.dataRaw || newVal.node_params || '{}'
    const oldDataRaw = oldVal.dataRaw || oldVal.node_params || '{}'

    if (newDataRaw != oldDataRaw) {
      reset()
    }
  },
  { deep: true }
)

onMounted(() => {
  reset()
  resetSize()
})

onBeforeUnmount(() => {})
</script>

<style lang="less" scoped>
.ai-dialogue-node {
  .field-list {
    .field-item {
      display: flex;
      margin-bottom: 8px;
      &:last-child {
        margin-bottom: 0;
      }
    }
    .field-item-label {
      width: 60px;
      line-height: 22px;
      margin-right: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #262626;
      text-align: right;
    }
    .field-item-content {
      flex: 1;
      overflow: hidden;
    }
    .category-list {
      width: 100%;
      overflow: hidden;
    }
    .field-value {
      width: 100%;
      line-height: 16px;
      padding: 3px 4px;
      height: 24px;
      margin-bottom: 8px;
      border-radius: 4px;
      font-size: 12px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      &:last-child {
        margin-bottom: 0;
      }
      .category-value {
        width: 100%;
      }

      .field-type {
        padding: 1px 8px;
        margin-left: 4px;
        border-radius: 4px;
        font-size: 12px;
        line-height: 16px;
        font-weight: 400;
        background: #e4e6eb;
      }
    }
  }

  .options-list {
    flex: 1;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .options-item {
    display: flex;
    align-items: center;
    height: 22px;
    padding: 2px 2px 2px 4px;
    border-radius: 4px;
    border: 1px solid #d9d9d9;

    &.is-required .option-label::before {
      vertical-align: middle;
      content: '*';
      color: #fb363f;
      margin-right: 2px;
    }

    .option-label {
      color: var(--wf-color-text-3);
      font-size: 12px;
      margin-right: 4px;
    }

    .option-type {
      height: 18px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 4px;
      font-size: 12px;
      background-color: #e4e6eb;
      color: var(--wf-color-text-3);
    }
  }
}
</style>
