<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_variable_assignment')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>

    <div class="variable-node">
      <div class="node-form-content">
        <div class="gray-block">
          <div class="field-items">
            <div class="field-item" v-for="(item, index) in list" :key="index">
              <div class="field-name-box">
                <a-select
                  style="width: 100%"
                  :placeholder="t('ph_select_variable')"
                  v-model:value="item.variable"
                  @dropdownVisibleChange="dropdownVisibleChange"
                  @change="update"
                >
                  <a-select-option :value="opt.value" v-for="opt in options" :key="opt.key">
                    <span>{{ opt.label }}</span>
                  </a-select-option>
                </a-select>
              </div>
              <div class="field-value-box">
                <at-input
                  :ref="'atInput' + index"
                  inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                  :options="valueOptions"
                  :defaultSelectedList="item.tags"
                  :defaultValue="item.value"
                  @open="showAtList"
                  @change="(text, selectedList) => changeValue(text, selectedList, item, index)"
                  :placeholder="t('ph_input_variable_value')"
                >
                  <template #option="{ label, payload }">
                    <div class="field-value-option">
                      <div class="option-label">{{ label }}</div>
                      <div class="option-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </at-input>
              </div>
              <div class="field-actions">
                <CloseCircleOutlined class="action-btn" @click="handleDel(index)" />
              </div>
            </div>
          </div>

          <div>
            <a-button style="width: 100%" class="add-btn" type="dashed" @click="handleAdd">
              <template #icon>
                <PlusOutlined />
              </template>
              {{ t('btn_add_assign_variable') }}
            </a-button>
          </div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import { ref, onMounted } from 'vue'
import { CloseCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import AtInput from '../at-input/at-input.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.variable-assignment-node-form')

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

const valueOptions = ref([])
const options = ref([])

function getOptions() {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let loop_parent_key = nodeModel.properties.loop_parent_key
    if (loop_parent_key) {
      const gropModel = props.lf.getNodeModelById(loop_parent_key)
      
      if (gropModel) {
        let dataRaw = gropModel.properties.dataRaw || gropModel.properties.node_params || '{}'
        let loop = JSON.parse(dataRaw).loop || {}

        options.value =loop.intermediate_params.map((item) => {
          return {
            label: item.key,
            value: loop_parent_key + '.' + item.key
          }
        })
        
        return
      }
    }
    let globalVariable = nodeModel.getGlobalVariable()
    let diy_global = globalVariable.diy_global || []
    diy_global.forEach((item) => {
      item.label = item.key
      item.value = 'global.' + item.key
    })

    options.value = diy_global || []
  }
}
function getValueOptions() {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()
    valueOptions.value = list
  }
}

const list = ref([])

const update = () => {
  let node_params = JSON.parse(props.node.node_params)
  node_params.assign = [...list.value]
  emit('update-node', {
    ...props.node,
    node_params: JSON.stringify(node_params)
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let assign = JSON.parse(dataRaw).assign || []

    assign = JSON.parse(JSON.stringify(assign))

    assign.forEach((item) => {
      item.tags = item.tags || []
    })

    list.value = assign

    getOptions()
  } catch (error) {
    console.log(error)
  }
}

const showAtList = (val) => {
  if (val) {
    getValueOptions()
  }
}

const changeValue = (text, selectedList, item) => {
  item.tags = selectedList
  item.value = text
  update()
}

const handleAdd = () => {
  list.value.push({
    variable: void 0,
    value: ''
  })
  update()
}

const handleDel = (index) => {
  list.value.splice(index, 1)
  update()
}

const dropdownVisibleChange = (visible) => {
  if (visible) {
    getOptions()
  }
}

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import './form-block.less';
.variable-node {
  .field-items {
    .field-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
    }

    .field-name-box {
      width: 130px;
      margin-right: 8px;
    }
    .field-value-box {
      flex: 1;
      overflow: hidden;
    }
    .field-actions {
      margin-left: 12px;
      .action-btn {
        cursor: pointer;
        font-size: 16px;
        color: #595959;
      }
    }
  }
}

.field-value-option {
  display: flex;
  align-items: center;
  color: #595959;
  font-size: 14px;
  font-weight: 400;
  .option-label {
    font-weight: 400;
    margin-right: 6px;
  }
  .option-type {
    width: fit-content;
    padding: 1px 8px;
    border-radius: 6px;
    border: 1px solid #00000026;
    background: var(--10, #fff);
    display: flex;
    align-items: center;
    font-size: 12px;
  }
}
</style>
