<template>
  <div>
    <!-- <a-button @click="handleOpenTestModal" style="background-color: #00ad3a" type="primary"
      ><CaretRightOutlined />运行测试</a-button
    > -->
    <a-modal
      v-model:open="show"
      :footer="null"
      :width="820"
      wrapClassName="no-padding-modal"
    >
      <template #title>
        <div class="modal-title-block">运行测试
          <div class="run-detail" v-if="resultList.length">
            <span>总耗时：{{ formatTime(use_mills) }}</span>
            <span>token消耗：{{ use_token }} Tokens</span>
          </div>
        </div>
      </template>
      <div class="flex-content-box">
        <div class="test-model-box">
          <div class="top-title">开始节点参数</div>
          <a-form
            :model="formState"
            ref="formRef"
            layout="vertical"
            :wrapper-col="{ span: 24 }"
            autocomplete="off"
          >
            <a-form-item
              v-if="is_need_question"
              name="question"
              :rules="[{ required: true, message: '请输入question!' }]"
            >
              <template #label>
                <a-flex :gap="4">question <a-tag style="margin: 0">string</a-tag> </a-flex>
              </template>
              <a-input placeholder="请输入" v-model:value="formState.question" />
            </a-form-item>
            <a-form-item name="openid" :rules="[{ required: true, message: '请输入openid!' }]" v-if="is_need_openid ">
              <template #label>
                <a-flex :gap="4">openid <a-tag style="margin: 0">string</a-tag> </a-flex>
              </template>
              <a-input placeholder="请输入" v-model:value="formState.openid" />
            </a-form-item>
            <a-form-item
              v-for="item in batch_test_params"
              :key="item.node_key"
              :rules="[{ required: item.field.required, message: `请输入${item.node_name}` }]"
            >
              <template #label>
                <a-flex :gap="4"
                  >{{ item.node_name }} {{ formatStr(item)
                  }}<a-tag style="margin: 0">{{ item.field.typ }}</a-tag>
                </a-flex>
              </template>
              <template v-if="item.field.typ == 'string'">
                <a-input placeholder="请输入" v-model:value="item.field.Vals" />
              </template>
              <template v-if="item.field.typ == 'number'">
                <a-input-number
                  style="width: 100%"
                  placeholder="请输入"
                  v-model:value="item.field.Vals"
                />
              </template>
              <template v-if="item.field.typ.includes('array')">
                <div class="input-list-box">
                  <div class="input-list-item" v-for="(input, i) in item.field.Vals" :key="i">
                    <a-form-item-rest
                      ><a-input placeholder="请输入" v-model:value="input.value"
                    /></a-form-item-rest>

                    <CloseCircleOutlined
                      v-if="item.field.Vals.length > 1"
                      @click="handleDelItem(item.field.Vals, i)"
                    />
                  </div>
                  <div class="add-btn-box">
                    <a-button @click="handleAddItem(item.field.Vals)" block type="dashed"
                      >添加</a-button
                    >
                  </div>
                </div>
              </template>
            </a-form-item>
          </a-form>

          <div class="result-list-box loading-box" v-if="loading">
            <a-spin v-if="loading" tip="测试结果生成中..." />
          </div>

          <div class="result-list-box" v-if="resultList.length > 0">
            <div
              class="list-item-block"
              :class="{ active: currentNodeKey == item.node_key }"
              v-for="(item, index) in resultList"
              @click="handleChangeNodeKey(item)"
              :key="index"
            >
              <div class="status-block">
                <CheckCircleFilled v-if="item.is_success" style="color: #138b1b" />
                <CloseCircleFilled v-else style="color: #d81e06" />
              </div>
              <div class="icon-name-box">
                <img :src="item.node_icon" alt="" />
                <div class="node-name">{{ item.node_name }}</div>
              </div>
              <div class="time-tag" v-if="item.is_success">{{ item.use_time }}ms</div>
              <div class="right-active-icon"><RightCircleOutlined /></div>
              <!-- <div class="out-put-box" v-if="item.is_success">
                <a-tooltip>
                  <template #title>{{ item.output }}</template>
                  <div class="out-text-box">{{ item.output }}</div>
                </a-tooltip>
              </div> -->
            </div>
          </div>

          <div class="save-btn-box">
            <a-button
              :loading="loading"
              @click="handleSubmit"
              style="background-color: #00ad3a"
              type="primary"
              ><CaretRightOutlined />运行测试</a-button
            >
          </div>
        </div>
        <div class="preview-box">
          <template v-if="cuttentItem">
            <div class="preview-title">
              <div class="title-text">日志详情</div>
              <div class="icon-name-box">
                <img :src="cuttentItem.node_icon" alt="" />
                <div class="node-name">{{ cuttentItem.node_name }}</div>
              </div>
              <div class="time-tag" v-if="cuttentItem.is_success">{{ cuttentItem.use_time }}ms</div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">运行日志<CopyOutlined @click="handleCopy" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="cuttentItem.output" />
              </div>
            </div>
          </template>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import {
  CaretRightOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  CloseCircleOutlined,
  RightCircleOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { reactive, ref, computed, nextTick } from 'vue'
import { useRobotStore } from '@/stores/modules/robot'
import { callBatchWorkFlow, callBatchWorkFlowParams } from '@/api/robot/index'
import { getImageUrl, formatTime } from '../../../util'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'
const robotStore = useRobotStore()

const isLockedByOther = computed(() => {
  return robotStore.robotInfo.isLockedByOther
})

const props = defineProps({
  batch_node_key: {
    type: String,
    default: ''
  }
})
const diy_global = computed(() => {
  return []
})

const emit = defineEmits(['save'])

const show = ref(false)
const currentNodeKey = ref('')
const resultList = ref([])

const cuttentItem = computed(() => {
  if (!currentNodeKey.value) {
    return null
  }
  return resultList.value.filter((item) => item.node_key == currentNodeKey.value)[0]
})

const loading = ref(false)

const formState = reactive({
  is_draft: true,
  robot_key: robotStore.robotInfo.robot_key,
  question: '',
  openid: ''
})

const handleOpenTestModal = async () => {
  if (isLockedByOther.value) {
    message.warning('当前已有其他用户在编辑中，无法运行测试')
    return
  }
  robotStore.robotInfo.loop_save_canvas_status++ // 触发保存草稿操作

  await nextTick()

  let localData = localStorage.getItem('workflow_batch_run_test_data') || '{}'

  localData = JSON.parse(localData)

  formState.question = localData.question || ''
  formState.openid = localData.openid || ''
  formState.global = localData.global || ''
  
  resultList.value = []
  currentNodeKey.value = ''
  setTimeout(() => {
    getFieldLists()
  }, 1000)
  show.value = true
}

const batch_test_params = ref([])
const is_need_question = ref(false)
const is_need_openid = ref(false)

const getFieldLists = () => {
  callBatchWorkFlowParams({
    robot_key: formState.robot_key,
    batch_node_key: props.batch_node_key,
  }).then((res) => {
    is_need_question.value = res.data.is_need_question
    is_need_openid.value = res.data.is_need_openid
    
    if (!is_need_question.value) {
      formState.question = ''
    }
    batch_test_params.value = res.data.loop_test_params.map((item) => {
      return {
        ...item,
        field: {
          ...item.field,
          Vals: setGlobalDefaultVal(item.field)
        }
      }
    })
  })
}

function setGlobalDefaultVal(item) {
  if (item.typ == 'string' || item.typ == 'number') {
    return ''
  }
  return [
    {
      value: '',
      key: Math.random() * 10000
    }
  ]
}

const handleDelItem = (item, index) => {
  item.splice(index, 1)
}
const handleAddItem = (item) => {
  item.push({
    value: '',
    key: Math.random() * 10000
  })
}

const formRef = ref(null)

const use_token = ref(0)
const use_mills = ref(0)

const handleSubmit = () => {
  formRef.value.validate().then(() => {
    let postData = { ...formState, batch_node_key: props.batch_node_key }
    let batch_test_params_data = JSON.parse(JSON.stringify(batch_test_params.value))

    batch_test_params_data.forEach((item) => {
      if (item.field.typ.includes('array')) {
        item.field.Vals = item.field.Vals.map((it) => it.value)
      }
    })

    postData.batch_test_params = JSON.stringify(batch_test_params_data)

    loading.value = true
    resultList.value = []
    localStorage.setItem('workflow_batch_run_test_data', JSON.stringify({ ...postData }))
    callBatchWorkFlow({
      ...postData
    })
      .then((res) => {
        message.success('测试结果生成完成')
        let node_logs = res.data.node_logs || []
        use_token.value = res.data.use_token
        use_mills.value = res.data.use_mills
        formateData(node_logs)
      })
      .catch((res) => {
        resultList.value = []
        let node_logs = res.data.node_logs || []
        if (node_logs && node_logs.length) {
          formateData(node_logs)
        }
      })
      .finally(() => {
        loading.value = false
      })
  })
}

const formateData = (data) => {
  resultList.value = data.map((item) => {
    return {
      ...item,
      is_success: item.error_msg === '<nil>',
      node_icon: getImageUrl(item.node_type)
    }
  })
  currentNodeKey.value = resultList.value[0]?.node_key
}

const handleChangeNodeKey = (item) => {
  currentNodeKey.value = item.node_key
}

const handleCopy = () => {
  copyText(JSON.stringify(cuttentItem.value.output))
  message.success('复制成功')
}

function formatStr(item) {
  return item.field.key.replace(item.node_key + '.', '')
}

const open = () => {
  handleOpenTestModal()
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.flex-content-box {
  display: flex;
  max-height: 600px;
  overflow: hidden;
}
.test-model-box {
  flex: 1;
  margin: 24px 0 0 0;
  padding-right: 16px;
  overflow-y: auto;
  .top-title {
    font-weight: 600;
    margin-bottom: 16px;
  }
  .save-btn-box {
    margin: 32px 0;
    margin-top: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
.tooltip-content {
  white-space: pre-wrap;
}
.loading-box {
  height: 100px;
  justify-content: center;
}
.result-list-box {
  margin: 24px 0;
  width: 100%;
  border: 1px solid #ebebeb;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  padding: 8px;
  .list-item-block {
    display: flex;
    align-items: center;
    overflow: hidden;
    gap: 8px;
    padding: 8px;
    color: #333;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    .right-active-icon {
      margin-left: auto;
      color: #2475fc;
      opacity: 0;
    }
    &:hover {
      background: #f2f4f7;
      .right-active-icon {
        opacity: 1;
      }
    }
    &.active {
      color: #2475fc;
      background: #e6efff;
      .right-active-icon {
        opacity: 0;
      }
    }
    .status-block {
      font-size: 20px;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-weight: 600;
      img {
        width: 24px;
        height: 24px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
    .out-put-box {
      flex: 1;
      margin-left: 24px;
      overflow: hidden;
      .out-text-box {
        background: #f2f2f2;
        border-radius: 6px;
        padding: 8px;
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.input-list-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
  .input-list-item {
    display: flex;
    gap: 8px;
  }
}

.preview-box {
  flex: 1;
  border-left: 1px solid #d9d9d9;
  padding: 16px;
  overflow-y: auto;
  .preview-title {
    display: flex;
    align-items: center;
    gap: 8px;
    .title-text {
      font-size: 15px;
      font-weight: 600;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      margin-left: 12px;
      img {
        width: 16px;
        height: 16px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
  }
  .preview-content-block {
    margin-top: 16px;
    .title-block {
      font-size: 15px;
      color: #262626;
      display: flex;
      align-items: center;
      gap: 4px;
      .anticon-copy {
        cursor: pointer;
        &:hover {
          color: #2475fc;
        }
      }
    }
    .preview-code-box {
      width: fit-content;
      margin-top: 16px;
      padding: 8px;
      border-radius: 8px;
      border: 1px solid #d9d9d9;

      &::v-deep(.vjs-tree) {
        width: fit-content;
      }

      &::v-deep(.vjs-tree-node) {
        width: calc(100% + 16px);
        padding-right: 16px;
      }
    }
  }
}

.modal-title-block{
  display: flex;
  align-items: center;
  gap: 12px;
  .run-detail{
    display: flex;
    align-items: center;
    gap: 16px;
    background: #BFFBD7;
    padding: 4px 16px;
    font-size: 14px;
    color: #595959;
    border-radius: 8px;
  }
}

</style>
