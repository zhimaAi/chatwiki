<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        @changeTitle="handleTitleChange"
        @deleteNode="handleDeleteNode"
        :desc="t('desc_session_trigger')"
      >
      </NodeFormHeader>
    </template>

    <div class="variable-node">
      <div class="node-form-content">
        <div class="gray-block multi-modal-box">
          <div class="output-label">
            <img src="@/assets/img/workflow/output.svg" alt="" class="output-label-icon" />
            <span class="output-label-text">{{ t('label_multi_modal_input') }}</span>
            &nbsp;&nbsp;
            <a-switch
              v-model:checked="formState.question_multiple_switch"
              @change="onChangeQuestionMultipleqSwitch"
            />
          </div>
          <div class="multi-modal-input">
            <div style="color: #8c8c8c;">{{ t('msg_multi_modal_desc') }}
              <a-tooltip placement="left">
                <template #title>
                  <div>{{ t('tip_question_multiple') }}</div>
              <pre><code>"question_multiple":[
  {
    "type":"text",
    "text":"{{ t('sample_text_content') }}"
  },{
    "type":"image_url",
    "image_url":{
      "url":"{{ t('sample_image_url') }}"
    }
  }
]
              </code></pre>
                </template>
                <QuestionCircleOutlined />
              </a-tooltip> 
            </div>
          </div>
        </div>
        <div class="gray-block">
          <div class="output-label">
            <img src="@/assets/img/workflow/output.svg" alt="" class="output-label-icon" />
            <span class="output-label-text">{{ t('label_output') }}</span>
            <span class="output-desc">{{ t('msg_output_auto_mapped') }}</span>
          </div>
          <div class="field-items">
            <div class="field-item" v-for="(item, index) in list" :key="index">
              <div class="field-name-box">
                <span class="field-name">{{item.key}}</span>
                <a-tooltip>
                  <template #title>{{ item.desc }}</template>
                  <QuestionCircleOutlined />
                </a-tooltip>
              </div>
              <div class="field-value-box">
                <a-select
                  style="width: 200px"
                  :placeholder="t('ph_select_variable')"
                  v-model:value="item.variable"
                  allowClear
                  @dropdownVisibleChange="dropdownVisibleChange"
                  @change="update"
                >
                  <a-select-option :disabled="selectedValues.includes(opt.value)" :value="opt.value" v-for="opt in options" :key="opt.key">
                    <span>{{ opt.label }}</span>
                  </a-select-option>
                </a-select>
              </div>
              <!-- <div class="field-desc">
                {{ item.desc }}
              </div> -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { ref, reactive, onMounted, inject, computed } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.session-trigger-form')

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

const getNode = inject('getNode')
const getGraph = inject('getGraph')

const list = ref([])
const options = ref([])
const formState = reactive({
  question_multiple_switch: false
})

const selectedValues = computed(() => { 
  return list.value.map(item => item.variable)
})

const onChangeQuestionMultipleqSwitch = (val) => {
  update()
}

function getOptions() {
  const nodeModel = getNode()

  if (nodeModel) {
    let globalVariable = nodeModel.getGlobalVariable()
    let diy_global = globalVariable.diy_global || []
    diy_global.forEach((item) => {
      item.label = item.key
      item.value = 'global.' + item.key
    })

    options.value = diy_global || []
  }
}

function dropdownVisibleChange(visible) {
  if (visible) {
    getOptions()
  }
}

const handleTitleChange = () => { 
  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change',  {...props.node})
  }, 10)
}

const handleDeleteNode = () => {
  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', null)
  }, 10)
}

const update = () => {
  let node_params = JSON.parse(props.node.node_params)

  node_params.trigger.outputs = [...list.value]
  node_params.trigger.chat_config = {...formState}

  let data = {...props.node, node_params: JSON.stringify(node_params)}

  emit('update-node', data)
  
  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', data)
  }, 10)
}

const init = () => {
  getOptions();

  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    
    dataRaw = JSON.parse(dataRaw)

    const trigger = dataRaw.trigger || {
      outputs: [],
      chat_config: {
        question_multiple_switch: false
      }
    }

    formState.question_multiple_switch = trigger.chat_config.question_multiple_switch

    list.value = trigger.outputs.map((item) => {
      item.tags = item.tags || []

      return item
    })
    
  } catch (error) {
    console.log(error)
  }
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import './form-block.less';
.variable-node {
  .multi-modal-box{
    margin-bottom: 8px;
  }
  .field-items {
    .field-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      justify-content: space-between;
      &:last-child {
        margin-bottom: 0;
      }
    }

    .field-name-box {
      width: auto;
      margin-right: 8px;
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .field-value-box {
      margin-right: 8px;
      
      .field-value{
        display: inline-flex;
        line-height: 20px;
        padding: 1px 8px;
        border-radius: 6px;
        overflow: hidden;
        background: #FFF;
        border: 1px solid rgba(0, 0, 0, 0.15);
      }

      .value-arrow{
        font-size: 16px;
        padding: 1px 4px;
        margin-right: 4px;
        border-radius: 4px;
        background: #E4E6EB;
      }
    }

    .field-desc{
      line-height: 22px;
      font-size: 14px;
      color: #595959;
      text-align: left;
    }
  }
}
</style>
