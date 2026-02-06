<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_generate_image')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/system-setting.svg" alt="" />{{ t('label_model_settings') }}
            </div>
            <div class="model-setting-form">
              <div class="model-label">{{ t('label_model') }}</div>
              <div class="model-content">
                <ModelSelect
                  modelType="IMAGE"
                  v-model:modeName="formState.use_model"
                  v-model:modeId="formState.model_config_id"
                  @change="handleChangeModel"
                />
              </div>
            </div>

            <div class="model-setting-form">
              <div class="model-label">{{ t('label_ratio') }}</div>
              <div class="model-content">
                <a-select :placeholder="t('ph_select')" v-model:value="formState.size" style="width: 100%">
                  <a-select-option
                    v-for="item in sizeOptions"
                    :value="item.value"
                    :key="item.value"
                    >{{ item.label }}</a-select-option
                  >
                </a-select>
              </div>
            </div>

            <div class="model-setting-form">
              <div class="model-label">
                {{ t('label_max_count') }}
                <a-tooltip>
                  <template #title>{{ t('tip_max_count') }}</template>
                  <QuestionCircleOutlined />
                </a-tooltip>
              </div>
              <div class="model-content">
                <a-input-number
                  :placeholder="t('ph_select')"
                  style="width: 100%"
                  v-model:value="formState.image_num"
                  :min="1"
                  :max="image_max"
                />
              </div>
            </div>

            <div class="model-setting-form" v-if="show_image_watermark">
              <div class="model-label">{{ t('label_watermark') }}</div>
              <div class="model-content">
                <a-segmented
                  class="customer-segmented"
                  v-model:value="formState.image_watermark"
                  :options="options1"
                />
              </div>
            </div>

            <div class="model-setting-form" v-if="show_image_optimize_prompt">
              <div class="model-label">{{ t('label_optimize_prompt') }}</div>
              <div class="model-content">
                <a-segmented
                  class="customer-segmented"
                  v-model:value="formState.image_optimize_prompt"
                  :options="options2"
                />
              </div>
            </div>
          </div>
          <div
            class="gray-block mt16"
            v-if="currentModelConfig.input_image == 1 || currentModelConfig.input_text == 1"
          >
            <div class="gray-block-title"><img src="@/assets/svg/input.svg" alt="" />{{ t('label_input') }}</div>
            <a-form-item name="use_model" v-if="currentModelConfig.input_text == 1">
              <template #label>
                <div style="width: 409px" class="flex-between-box">
                  <div>{{ t('label_prompt') }}</div>
                  <div class="btn-hover-wrap" @click="handleOpenFullAtModal">
                    <FullscreenOutlined />
                  </div>
                </div>
              </template>
              <at-input
                :options="valueOptions"
                :defaultSelectedList="formState.prompt_tags"
                :defaultValue="formState.prompt"
                ref="promptInputRef"
                :placeholder="t('ph_input_message')"
                input-style="height: 76px"
                type="textarea"
                @open="showAtList"
                @change="(text, selectedList) => changeValue(text, selectedList)"
              />
            </a-form-item>
            <a-form-item
              :label="t('label_reference_image')"
              name="use_model"
              v-if="currentModelConfig.input_image == 1"
            >
              <a-form-item-rest>
                <div class="array-form-box">
                  <div
                    class="form-item-list"
                    v-for="(item, index) in formState.input_images"
                    :key="item.key"
                  >
                    <div class="flex-block-item">
                      <at-input
                        inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                        :ref="(el) => setAtInputRef(el, 'input_images', index)"
                        :options="variableOptions"
                        :defaultSelectedList="item.tags"
                        :defaultValue="item.value"
                        @open="getVlaueVariableList"
                        @change="
                          (text, selectedList) => changeImgListValue(text, selectedList, item)
                        "
                        :placeholder="t('ph_input_value')"
                      >
                        <template #option="{ label, payload }">
                          <div class="field-list-item">
                            <div class="field-label">{{ label }}</div>
                            <div class="field-type">{{ payload.typ }}</div>
                          </div>
                        </template>
                      </at-input>

                      <div class="btn-hover-wrap" @click="onDelcategory(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </div>
                  <a-button @click="handleAddcategory" :icon="h(PlusOutlined)" block type="dashed"
                    >{{ t('btn_add_reference_image') }}</a-button
                  >
                </div>
              </a-form-item-rest>
            </a-form-item>
          </div>
          <div class="gray-block mt16">
            <div class="gray-block-title"><img src="@/assets/svg/output.svg" alt="" />{{ t('label_output') }}</div>
            <div class="output-item">
              <div class="key-label">msg</div>
              <div class="key-value">string</div>
            </div>
            <div class="output-item" v-for="item in outputImggeList" :key="item">
              <div class="key-label">{{ item }}</div>
              <div class="key-value">string</div>
            </div>
          </div>
        </a-form>
      </div>
    </div>
    <FullAtInput
      :options="valueOptions"
      :defaultSelectedList="formState.prompt_tags"
      :defaultValue="formState.prompt"
      :placeholder="t('ph_input_message')"
      input-style="height: 76px"
      type="textarea"
      @open="showAtList"
      @change="(text, selectedList) => changeValue(text, selectedList)"
      @ok="handleRefreshAtInput"
      ref="fullAtInputRef"
    />
  </NodeFormLayout>
</template>

<script setup>
import { getUuid } from '@/utils/index'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, computed, onMounted, h, inject } from 'vue'
import {
  CloseCircleOutlined,
  QuestionCircleOutlined,
  PlusOutlined,
  FullscreenOutlined
} from '@ant-design/icons-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import AtInput from '../../at-input/at-input.vue'
import { getCurrentModelConfig } from '@/components/model-select/index.js'
import { getSizeOptions } from '@/views/workflow/components/util.js'
import { message } from 'ant-design-vue'
import FullAtInput from '../../at-input/full-at-input.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.image-generation-node-form.index')

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const options1 = computed(() => [
  {
    label: t('label_show'),
    value: '1'
  },
  {
    label: t('label_hide'),
    value: '0'
  }
])

const options2 = computed(() => [
  {
    label: t('label_enable'),
    value: '1'
  },
  {
    label: t('label_disable'),
    value: '0'
  }
])

const variableOptions = ref([])

const getNode = inject('getNode')

const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  size: void 0,
  image_num: '1',

  prompt: '',
  prompt_tags: [],
  input_images: [
    {
      value: '',
      key: getUuid(16)
    }
  ],
  image_watermark: '1',
  image_optimize_prompt: '1'
})

const promptInputRef = ref(null)
const fullAtInputRef = ref(null)

const handleRefreshAtInput = () => {
  promptInputRef.value.refresh()
}

const handleOpenFullAtModal = () => {
  fullAtInputRef.value.show()
}

const outputImggeList = computed(() => {
  let list = []
  if (formState.image_num > 0) {
    for (let i = 0; i < +formState.image_num; i++) {
      let letter = String.fromCharCode('a'.charCodeAt(0) + i)
      list.push(`picture_url_${letter}`)
    }
  }
  return list
})

const handleChangeModel = () => {
  formState.size = void 0
  formState.image_num = '1'
  formState.input_images = [
    {
      value: '',
      key: getUuid(16)
    }
  ]
}

const currentModelConfig = computed(() => {
  return getCurrentModelConfig(formState.model_config_id, formState.use_model, 'IMAGE') || {}
})

const image_generation_config = computed(() => {
  if (currentModelConfig.value.image_generation) {
    return JSON.parse(currentModelConfig.value.image_generation)
  }
  return {}
})

const allSizeOptions = getSizeOptions()
const sizeOptions = computed(() => {
  if (image_generation_config.value.image_sizes) {
    return image_generation_config.value.image_sizes
      .split(',')
      .map((item) => allSizeOptions.find((it) => it.value == item))
  }
  return []
})

const image_max = computed(() => {
  return +image_generation_config.value.image_max || 1
})

const show_image_optimize_prompt = computed(() => {
  return image_generation_config.value.image_optimize_prompt == '1'
})

const show_image_watermark = computed(() => {
  return image_generation_config.value.image_watermark == '1'
})

const image_inputs_image_max = computed(() => {
  return +image_generation_config.value.image_inputs_image_max || 1
})

const update = () => {
  let output = [
    {
      key: 'msg',
      typ: 'string'
    }
  ]
  outputImggeList.value.forEach((item) => {
    output.push({
      key: item,
      typ: 'string'
    })
  })
  const data = JSON.stringify({
    image_generation: {
      ...formState,
      image_num: formState.image_num > 0 ? formState.image_num + '' : '1',
      input_images: formState.input_images.map((item) => item.value).filter(Boolean),
      output
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let image_generation = JSON.parse(dataRaw).image_generation || {}
    getValueOptions()
    getVlaueVariableList()
    image_generation = JSON.parse(JSON.stringify(image_generation))
    formState.model_config_id = image_generation.model_config_id
    formState.use_model = image_generation.use_model
    formState.size = image_generation.size
    formState.image_num = image_generation.image_num
    formState.prompt = image_generation.prompt
    formState.prompt_tags = image_generation.prompt_tags
    formState.input_images = image_generation.input_images.map((item) => {
      return {
        value: item,
        key: getUuid(16)
      }
    })
    formState.image_watermark = image_generation.image_watermark
    formState.image_optimize_prompt = image_generation.image_optimize_prompt
  } catch (error) {
    console.log(error)
  }
}

const valueOptions = ref([])
const showAtList = (val) => {
  if (val) {
    getValueOptions()
  }
}

const getValueOptions = () => {
  let options = getNode().getAllParentVariable()
  valueOptions.value = options || []
}

const changeValue = (text, selectedList) => {
  formState.prompt_tags = selectedList
  formState.prompt = text
}

const getVlaueVariableList = () => {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()
    list.forEach((item) => {
      item.tags = item.tags || []
    })

    variableOptions.value = list
  }
}

const atInputRefs = reactive({})
const setAtInputRef = (el, name, index) => {
  if (el) {
    let key = `at_input_${name}_${index}`
    atInputRefs[key] = el
  }
}

const changeImgListValue = (text, selectedList, item) => {
  item.tags = selectedList
  item.value = text
}

const handleAddcategory = () => {
  if (formState.input_images.length >= image_inputs_image_max.value) {
    return message.error(t('msg_max_images_limit', { val: image_inputs_image_max.value }))
  }
  formState.input_images.push({
    value: '',
    key: getUuid(16)
  })
}
const onDelcategory = (index) => {
  formState.input_images.splice(index, 1)
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.flex-between-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.btn-hover-wrap {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease-in;
  &:hover {
    background: #e4e6eb;
  }
}
.model-setting-form {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  .model-label {
    width: 114px;
    text-align: left;
    color: #262626;
  }
  .model-content {
    flex: 1;
  }
}
.form-item-list {
  margin-bottom: 8px;
}
.model-icon {
  height: 18px;
}

.output-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  .key-label {
    color: #262626;
    font-size: 14px;
  }
  .key-value {
    padding: 1px 8px;
    display: flex;
    align-items: center;
    border: 1px solid #00000026;
    background: var(--10, #fff);
    border-radius: 6px;
    font-size: 12px;
    color: #595959;
  }
}
</style>