<style lang="less" scoped>
.agent-node {
  .field-list {
    .field-item {
      display: flex;
      margin-bottom: 8px;

      &:last-child {
        margin-bottom: 0;
      }
    }

    .field-item-label {
      width: auto;
      line-height: 22px;
      margin-right: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #262626;
      text-align: right;
    }

    .field-item-content {
      flex: 1;
      display: flex;
      flex-wrap: wrap;
    }

    .field-value {
      width: 100%;
      display: flex;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      font-size: 12px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;

      &.is-required .field-key::before {
        content: '*';
        color: #FB363F;
        display: inline-block;
        margin-right: 2px;
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

  .text-input {
    width: 100%;
    height: 58px;
    line-height: 18px;
    font-size: 13px;
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    white-space: pre-wrap;
  }
}
</style>

<template>
  <node-common
    :properties="properties"
    :title="properties.node_name"
    :menus="menus"
    :icon-url="robotAvatar"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    style="width: 420px;"
    @handleMenu="handleMenu"
  >
    <div class="agent-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">提问内容</div>
          <div class="field-item-content">
            <div class="field-value">
              <AtText
                :options="valueOptions"
                :defaultSelectedList="formState.question_tags"
                :defaultValue="formState.question_value"
                ref="atInputRef"
                class="text-input"
                @resize="onInputResize"
                v-if="formState.question_value.length > 0"
              />
            </div>
          </div>
        </div>

        <div class="field-item" v-if="false">
          <div class="field-item-label">输出</div>
          <div class="field-item-content">
            <div
              v-for="item in formState.output"
              :key="item.key"
              class="field-value is-required"
            >
              <span class="field-key">{{ item.desc || item.key }}</span>
              <span class="field-type">{{ item.typ }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { reactive, ref, watch, inject, nextTick } from 'vue'
import { jsonDecode } from '@/utils/index'
import NodeCommon from '../base-node.vue'
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
const getGraph = inject('getGraph')
const resetSize = inject('resetSize')

const menus = ref([])
const robotAvatar = ref('')
const atInputRef = ref(null)
const valueOptions = ref([])
const formState = reactive({
  question_value: '',
  question_tags: [],
  output: []
})

function getValueOptions() {
  valueOptions.value = getNode().getAllParentVariable() || []
}

function reset() {
  getValueOptions()

  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  const agent = jsonDecode(dataRaw)?.agent || {}
  robotAvatar.value = agent?.robot_info?.robot_avatar || props.properties.node_icon || ''
  formState.question_value = agent.question_value || ''
  formState.question_tags = agent.question_tags || []
  formState.output = Array.isArray(agent.output) ? agent.output : []

  nextTick(() => {
    resetSize()
  })
}

function onInputResize() {
  nextTick(() => {
    resetSize()
  })
}

function handleMenu(item) {
  if (item.key === 'delete') {
    let node = getNode()

    getGraph().deleteNode(node.id)
  }
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

reset()
</script>
