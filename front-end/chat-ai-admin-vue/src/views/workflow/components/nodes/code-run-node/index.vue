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
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
    .field-list-item {
      display: flex;
      gap: 4px;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
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
        max-width: 150px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 12px;
      }
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
      .field-key{
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
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
}
</style>

<template>
  <node-common
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    style="width: 420px;"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">输入参数</div>
          <div class="field-item-content">
            <div class="field-list-item" v-for="item in formState.params" :key="item.cu_key">
              <div class="field-text">
                <UserQuestionText :value="item.variable" />
              </div>
              <div class="right-arrow"><ArrowRightOutlined /></div>
              <div class="field-text">{{ item.field }}</div>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">输出字段</div>
          <div class="field-item-content">
            <div class="field-value" v-for="item in formState.output" :key="item.cu_key">
              <span class="field-key"> {{ item.key }}</span>
              <span class="field-type">{{ item.typ }}</span>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">异常处理</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key">运行代码报错时执行该分支</span>
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
import UserQuestionText from '../user-question-text.vue'
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
const graphModel = inject('getGraph')
const getNode = inject('getNode')
const resetSize = inject('resetSize')

// --- State ---
const menus = ref([])
const formState = reactive({
  main_func: '',
  params: [
    {
      field: '',
      variable: []
    }
  ],
  timeout: 30,
  output: [
    {
      key: '',
      typ: '',
      subs: []
    }
  ],
  exception: ''
})

const defaultCode = `function main({data1, data2}){
	return {
		data1,
    	data2
	}
}`

const variableOptions = ref([])

let specialNodeList = [
  'special.lib_paragraph_list',
  'special.llm_reply_content',
  'specify-reply-node'
]

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
const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let code_run = {}
  try {
    code_run = JSON.parse(dataRaw).code_run || {}
  } catch (e) {
    code_run = {}
  }

  getVlaueVariableList()

  code_run = JSON.parse(JSON.stringify(code_run))

  for (let key in code_run) {
    if (key == 'params') {
      if (code_run[key] && code_run[key].length > 0) {
        formState[key] = code_run[key].map((item) => {
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
      formState['output'] = recursionData(code_run[key])
      continue
    }
    formState[key] = code_run[key]
  }
  formState.main_func = formState.main_func || defaultCode

  formState.params = formState.params.map((it) => {
    let specialKey = ''
    for (let i = 0; i < specialNodeList.length; i++) {
      if (it.variable.indexOf(specialNodeList[i]) > -1) {
        specialKey = specialNodeList[i]
        break
      }
    }
    if (specialKey != '') {
      let arr = it.variable.split('.')
      it.variable = [arr[0], specialKey]
    } else {
      it.variable = it.variable.split('.')
    }
    return {
      ...it,
      key: Math.random() * 10000
    }
  })

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = JSON.stringify({
    code_run: {
      ...formState,
      params: formatParams()
    }
  })

  setData({
    ...props.node,
    ...formState,
    node_params: data
  })
}

function formatParams() {
  let list = []
  formState.params.forEach((item) => {
    list.push({
      field: item.field,
      variable: item.variable && item.variable.length > 0 ? item.variable.join('.') : ''
    })
  })
  return list
}

const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()

  variableOptions.value = list
}

const onUpatateNodeName = (data) => {
  if (!haveOutKeyNode.includes(data.node_type)) {
    return
  }
  getVlaueVariableList()

  nextTick(() => {
    update()
  })
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
