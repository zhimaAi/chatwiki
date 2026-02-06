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
      width: 90px;
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
    .field-list-item {
      display: flex;
      align-items: center;
      gap: 4px;
      width: 100%;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      overflow: hidden;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      overflow: hidden;
      margin-bottom: 4px;
      &:last-child {
        margin-bottom: 0;
      }
      .right-arrow {
        width: 24px;
        height: 100%;
        border-radius: 4px;
        background: #e4e6eb;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 13px;
      }

      .field-text{
        max-width: 120px;
        font-size: 12px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
      .field-text2{
        max-width: 120px;
        font-size: 12px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }
    .body_raw{
      width: 100%;
      display: block;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .field-value {
      display: flex;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      font-size: 12px; 
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      overflow: hidden;

      .field-type {
        flex: 1;
        padding: 1px 8px;
        margin-left: 4px;
        border-radius: 4px;
        font-size: 12px;
        line-height: 16px;
        font-weight: 400;
        background: #e4e6eb;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }
  }
}
</style>

<template>
  <node-common
    :properties="properties"
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    @click="handleConsole"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">{{ t('label_request_url') }}</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key"> {{ formState.method }}</span>
              <span class="field-type" v-if="formState.rawurl">{{ formState.rawurl }}</span>
            </div>
          </div>
        </div>
        <div class="field-item">
          <div class="field-item-label">{{ t('label_headers') }}</div>
          <div class="field-item-content">
            <div class="field-list-item" v-for="item in formState.headers" :key="item.cu_key">
              <div class="field-text">{{ item.key }}</div>
              <div class="right-arrow"><ArrowRightOutlined /></div>
              <div class="field-text2">
                <AtText :options="variableOptions" :default-value="item.value" :defaultSelectedList="item.tags" />
              </div>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">{{ t('label_params') }}</div>
          <div class="field-item-content">
            <div class="field-list-item" v-for="item in formState.params" :key="item.cu_key">
              <div class="field-text">{{ item.key }}</div>
              <div class="right-arrow"><ArrowRightOutlined /></div>
              <div class="field-text2">
                <AtText :options="variableOptions" :default-value="item.value" :defaultSelectedList="item.tags" />
              </div>
            </div>
          </div>
        </div>

        <div class="field-item" v-if="formState.type == 1">
          <div class="field-item-label">{{ t('label_body') }}</div>
          <div class="field-item-content">
            <div class="field-list-item" v-for="item in formState.body" :key="item.cu_key">
              <div class="field-text">{{ item.key }}</div>
              <div class="right-arrow"><ArrowRightOutlined /></div>
              <div class="field-text2">
                <AtText :options="variableOptions" :default-value="item.value" :defaultSelectedList="item.tags" />
              </div>
            </div>
          </div>
        </div>
        <div class="field-item" v-if="formState.type == 2">
          <div class="field-item-label">{{ t('label_body') }}</div>
          <div class="field-item-content">
            <div class="field-list-item body_raw">
              {{ formState.body_raw }}
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">{{ t('label_output_fields') }}</div>
          <div class="field-item-content">
            <div class="field-value" v-for="item in formState.output" :key="item.cu_key">
              <span class="field-key"> {{ item.key }}</span>
              <span class="field-type">{{ item.typ }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount } from 'vue'
import NodeCommon from '../base-node.vue'
import { ArrowRightOutlined } from '@ant-design/icons-vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
import AtText from '../../at-input/at-text.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.nodes.http-node.index')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const setData = inject('setData')
const graphModel = inject('getGraph')
const getNode = inject('getNode')
const resetSize = inject('resetSize')

// --- State ---
const menus = ref([])
const formState = reactive({
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
  body: [],
  body_raw: '',
  body_raw_tags: [],
  timeout: 30,
  output: []
})

function recursionData(data) {
  data.forEach((item) => {
    item.cu_key = Math.random() * 10000
    if (item.subs && item.subs.length) {
      recursionData(item.subs)
    } else {
      item.subs = []
    }
  })
  return data
}

const variableOptions = ref([])

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let curl = {}
  try {
    curl = JSON.parse(dataRaw).curl || {}
  } catch (e) {
    curl = {}
  }

  getVlaueVariableList()

  curl = JSON.parse(JSON.stringify(curl))

  for (let key in curl) {
    if (key == 'headers' || key == 'params' || key == 'body') {
      if (curl[key] && curl[key].length > 0) {
        formState[key] = curl[key].map((item) => {
          return {
            ...item,
            cu_key: Math.random() * 10000
          }
        })
      } else {
        formState[key] = []
      }
      continue
    }
    if (key == 'output') {
      formState['output'] = recursionData(curl[key])
      continue
    }
    formState[key] = curl[key]
  }

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = JSON.stringify({
    curl: {
      ...formState
    }
  })

  setData({
    ...props.node,
    ...formState,
    node_params: data
  })
}

const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()

  list.forEach((item) => {
    item.tags = item.tags || []
  })
  variableOptions.value = list
}

const onUpatateNodeName = (data) => {
  if (!haveOutKeyNode.includes(data.node_type)) {
    return
  }

  getVlaueVariableList()

  nextTick(() => {
    if (formState.body_raw_tags && formState.body_raw_tags.length > 0) {
      formState.body_raw_tags.forEach((tag) => {
        if (tag.node_id == data.node_id) {
          let arr = tag.label.split('/')
          arr[0] = data.node_name
          tag.label = arr.join('/')
          tag.node_name = data.node_name
        }
      })
    }

    let keys = ['headers', 'params', 'body']

    keys.forEach((key) => {
      let items = formState[key]

      items.forEach((item) => {
        if (item.tags && item.tags.length > 0) {
          item.tags.forEach((tag) => {
            if (tag.node_id == data.node_id) {
              let arr = tag.label.split('/')
              arr[0] = data.node_name
              tag.label = arr.join('/')
              tag.node_name = data.node_name
            }
          })
        }
      })
    })

    update()
  })
}

const handleConsole = () => {

}

// --- Watchers and Lifecycle Hooks ---
watch(() => props.properties, reset, { deep: true })

onMounted(() => {
  reset()
  resetSize()
  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
