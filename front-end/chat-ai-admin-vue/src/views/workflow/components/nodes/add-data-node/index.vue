<template>
  <node-common
    :properties="properties"
    :title="props.properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    style="width: 420px"
  >
    <div class="add-data-node">
      <div class="static-field-list">
        <div class="static-field-item">
          <div class="static-field-item-label">数据表</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              {{ state.formData.form_name || '--' }}
            </div>
          </div>
        </div>

        <div class="static-field-item">
          <div class="static-field-item-label">插入数据</div>
          <div class="static-field-item-content">
            <div class="static-field-value" v-if="state.formData.datas.length == 0">--</div>
            <div class="static-field-value" v-for="item in state.formData.datas" :key="item.id">
              <span class="field-name">
                <AtText
                  :options="atInputOptions"
                  :default-value="item.value"
                  :defaultSelectedList="item.atTags"
                  ref="atInputRef"
                />
              </span>

              <span class="field-arrow"><img src="@/assets/img/workflow/arrow-right.svg" alt=""/></span>
              
              <span class="field-name">{{ item.name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, inject, onMounted, reactive, nextTick, watch, toRaw, onUnmounted } from 'vue'
import NodeCommon from '../base-node.vue'
import AtText from '../../at-input/at-text.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const setData = inject('setData')
const resetSize = inject('resetSize')
const getNode = inject('getNode')
const getGraph = inject('getGraph')
const atInputRef = ref(null)

let node_params = {}

const state = reactive({
  tableList: [],
  formData: {
    form_name: '',
    form_description: '',
    form_id: '',
    datas: []
  }
})

const atInputOptions = ref([])

const getAtInputOptions = () => {
  let options = getNode().getAllParentVariable()

  atInputOptions.value = options || []
}

const init = () => {
  let dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'

  node_params = JSON.parse(dataRaw)

  getAtInputOptions()

  state.formData = node_params.form_insert

  update()

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  node_params.form_insert = toRaw(state.formData)

  setData({
    ...props.node,
    node_params: JSON.stringify(node_params)
  })
}

watch(
  () => props.properties.dataRaw,
  (newVal, oldVal) => {
    if (newVal !== oldVal) {
      init()
    }
  }
)

onMounted(() => {
  const graphModel = getGraph()
  graphModel.eventCenter.on('custom:setNodeName', onUpatateNodeName)
  init()
})

onUnmounted(() => {
  const graphModel = getGraph()
  graphModel.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})

const onUpatateNodeName = (data) => {
  if (!state.formData.datas || state.formData.datas.length == 0) {
    return
  }
  if (!haveOutKeyNode.includes(data.node_type)) {
    return
  }

  getAtInputOptions()

  state.formData.datas.forEach((item) => {
    if (item.atTags && item.atTags.length > 0) {
      item.atTags.forEach((tag) => {
        if (tag.node_id == data.node_id) {
          let arr = tag.label.split('/')
          arr[0] = data.node_name
          tag.label = arr.join('/')
          tag.node_name = data.node_name
        }
      })
    }
  })

  update()
  nextTick(() => {
    atInputRef.value.map(item => item.refresh())
    resetSize()
  })
}
</script>

<style lang="less" scoped>
@import '../form-block.less';
.add-data-node {
  position: relative;
  .static-field-item-content {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .static-field-value {
    font-size: 12px;

    .field-arrow {
      display: inline-flex;
      align-items: center;
      margin: 0 4px;
      padding: 0 4px;
      height: 18px;
      font-size: 12px;
      border-radius: 4px;
      background: #e4e6eb;
      img {
        width: 16px;
        height: 16px;
      }
    }
  }
  &:deep(.j-mention-at) {
    padding: 0 8px;
    border-radius: 4px;
    font-size: 12px;
    line-height: 16px;
    font-weight: 400;
    background: #f2f4f5;
  }
}
</style>