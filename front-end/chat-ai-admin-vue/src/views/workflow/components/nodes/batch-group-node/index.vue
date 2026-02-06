<template>
  <div class="custom-group-content">
    <BorderLine :isHovered="isHovered" :isSelected="isSelected" />
    <div class="header-box" :style="{ background: properties.node_header_bg_color }">
      <div class="top-block">
        <div class="title-block">
          <img src="@/assets/svg/batch-group-node.svg" alt="" />
          <div class="title-text">{{ props.properties.node_name }}</div>
        </div>
        <div class="btn-block">
          <a-popover :title="null" trigger="click" v-model:open="isShowMenu" placement="right">
            <template #content>
              <NodeListPopup style="width: 500px;height: 550px;" :excludedNodeTypes="excludedNodeTypes" @addNode="handleAddNode" type="loop-node" />
            </template>
            <a-tooltip :title="t('tooltip_create_inner_node')">
              <div class="btn-item" @click.stop="handleClick">
                <PlusCircleOutlined />
              </div>
            </a-tooltip>
          </a-popover>

          <a-tooltip :title="t('tooltip_run_test')">
            <div class="btn-item" @click.stop="handleOpenTestModal">
              <CaretRightOutlined />
            </div>
          </a-tooltip>

          <a-dropdown trigger="click">
            <div class="btn-item" @click.stop="">
              <img class="btn-icon" src="@/assets/img/workflow/node-menu-btn.svg" alt="" />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="handleDelete">
                  <div style="color: #fb363f">{{ t('btn_delete') }}</div>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
      <div class="info-block">
        <div class="info-item">
          <div class="info-label">{{ t('label_execution_array') }}</div>
          <div class="info-content">
            <div class="label-text">input</div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="node-body">
      <div class="node-content"></div>
    </div>
    <RunTest ref="runTestRef" :batch_node_key="batch_node_key" />
  </div>
</template>

<script setup>
import { generateRandomId } from '@/utils/index'
import { ref, reactive, watch, onMounted, inject, nextTick, computed } from 'vue'
import { message } from 'ant-design-vue'
import { CaretRightOutlined, PlusCircleOutlined } from '@ant-design/icons-vue'
import NodeListPopup from '../../node-list-popup/index.vue'
import RunTest from './components/run-test.vue'
import BorderLine from './components/border-line.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.nodes.batch-group-node.index')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  model: {
    type: Object,
    default: () => ({})
  },
  isSelected: {
    type: Boolean,
    default: false
  },
  isHovered: {
    type: Boolean,
    default: false
  }
})
const getGraph = inject('getGraph')
const getNode = inject('getNode')
const addNode = inject('addNode')
const setData = inject('setData')
const resetSize = inject('resetSize')
const excludedNodeTypes = ref(['custom-group', 'batch-group', 'end-node', 'terminate-node', 'qa-node'])

const isShowMenu = ref(false)

const batch_node_key = computed(() => {
  return getNode().id
})

const handleAddNode = (node) => {
  addNode(node)
  isShowMenu.value = false
}
const handleClick = () => {
  isShowMenu.value = true
}

const runTestRef = ref(null)

const handleOpenTestModal = () => {
  let data = formState.batch_arrays[0]

  if(data && data.key && data.typ && data.value){
    runTestRef.value.open()
  }else{
    message.error(t('msg_fill_execution_array'))
  }
}

const formState = reactive({
  chan_number: 10,
  max_run_number: 500,
  batch_arrays: [
    {
      key: '',
      typ: '',
      value: '',
      _id: generateRandomId(16)
    }
  ],
  output: [
    {
      key: '',
      typ: '',
      value: '',
      _id: generateRandomId(16)
    }
  ]
})

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let batch = {}

  try {
    batch = JSON.parse(dataRaw).batch || {}
  } catch (e) {
    batch = {}
  }

  formState.chan_number = batch.chan_number
  formState.batch_arrays = batch.batch_arrays
  formState.output = batch.output

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = JSON.stringify({
    batch: {
      ...formState
    }
  })

  setData({
    ...props.node,
    node_params: data
  })
}

watch(
  () => props.properties,
  (val, oldVal) => {
    if(JSON.stringify(val) !== JSON.stringify(oldVal)){
      reset()
    }
    
  },
  { deep: true }
)

const handleDelete = () => {
  let node = getNode()
  getGraph().deleteNode(node.id)
}

onMounted(() => {
  reset()
  resetSize()
})
</script>

<style lang="less" scoped>
.custom-group-content {
  width: 100%;
  height: 100%;
  position: relative;
  background: #fff;
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  // border: 1px solid #fff;
  &.isHovered {
    // border: 1px solid #2475fc;
  }
  &.isSelected {
    // border: 2px solid #2475fc;
    // &::before {
    //   content: '';
    //   position: absolute;
    //   top: 0;
    //   left: 0;
    //   width: 100%;
    //   height: 100%;
    //   border-radius: 8px;
    //   border: 2px solid #2475fc;
    // }
  }
  .header-box {
    overflow: hidden;
    padding: 16px;

    .top-block {
      display: flex;
      align-items: center;
      justify-content: space-between;
      .title-block {
        display: flex;
        align-items: center;
        flex: 1;
        gap: 8px;
        color: #262626;
        font-size: 16px;
        font-weight: 600;
        img {
          width: 20px;
          height: 20px;
        }
      }
      .btn-block {
        display: flex;
        align-items: center;
        gap: 8px;
        .btn-item {
          width: 24px;
          height: 24px;
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          border-radius: 6px;
          transition: all 0.2s ease-in-out;
          font-size: 16px;
          &:hover {
            background: #e4e6eb;
          }
        }
      }
    }
  }

  .info-block {
    margin-top: 16px;
    padding-left: 20px;
    .info-item {
      display: flex;
      align-items: center;
      gap: 8px;
      .info-label {
        color: #262626;
        font-size: 14px;
      }
      .info-content {
        display: flex;
        width: fit-content;
        border: 1px solid #d9d9d9;
        background: #fff;
        gap: 4px;
        border-radius: 4px;
        color: #595959;
        height: 22px;
        .label-text {
          padding: 2px 2px 2px 4px;
          font-size: 12px;
          display: flex;
          align-items: center;
        }
        .label-value {
          border-radius: 4px;
          background: #e4e6eb;
          color: #595959;
          font-size: 12px;
          padding: 1px 4px;
        }
      }
    }
  }

  .node-body{
    flex: 1;
    border-radius: 6px;
    padding: 0 16px 16px;
  }
  .node-content {
    width: 100%;
    height: 100%;
    background: #f2f4f7;
  }
}
</style>
