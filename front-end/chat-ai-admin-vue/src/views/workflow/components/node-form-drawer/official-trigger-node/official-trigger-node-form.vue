<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        @changeTitle="handleTitleChange"
        @deleteNode="handleDeleteNode"
        desc="当对接的公众号触发相关事件时，自动触发流程"
      >
      </NodeFormHeader>
    </template>

    <div class="variable-node">
      <div class="node-form-content">
        <div class="gray-block">
          <div class="output-label">
            <img src="@/assets/svg/wx-app-icon.svg" alt="" class="output-label-icon" />
            <span class="output-label-text">公众号</span>
          </div>
          <div>
            <a-button @click="handleOpenAppModal" type="dashed"
              ><PlusOutlined />添加公众号</a-button
            >
            <div class="app-item-list">
              <div class="app-item" v-for="item in selectAppItems" :key="item.app_id">
                {{ item.app_name }}
                <span @click="handleDelOfficial(item)" class="close-icon"
                  ><CloseCircleFilled
                /></span>
              </div>
            </div>
          </div>
          <div class="form-item">
            <div class="form-item-label">触发事件</div>
            <div class="form-item-content">
              <a-select
                :options="trigger_options"
                v-model:value="formState.msg_type"
                style="width: 100%"
                disabled
              >
              </a-select>
            </div>
          </div>
        </div>
        <div class="gray-block">
          <div class="output-label">
            <img src="@/assets/svg/output.svg" alt="" class="output-label-icon" />
            <span class="output-label-text">输出</span>
            <a class="ml8" href="https://developers.weixin.qq.com/doc/subscription/guide/product/message/Receiving_standard_messages.html" target="_blank">查看官方接口文档</a>
          </div>
          <div class="field-items">
            <!-- <div class="field-items-label">文本消息</div> -->
            <div class="field-item" v-for="(item, index) in list" :key="index">
              <div class="field-name-box">
                <span class="field-name">{{ item.key }}</span>
              </div>
              <div class="field-value-box">
                <a-select
                  style="width: 200px"
                  placeholder="请输入选择变量"
                  v-model:value="item.variable"
                  allowClear
                  @dropdownVisibleChange="dropdownVisibleChange"
                  @change="update"
                >
                  <!--  :disabled="selectedValues.includes(opt.value)" -->
                  <a-select-option :value="opt.value" v-for="opt in options" :key="opt.key">
                    <span>{{ opt.label }}</span>
                  </a-select-option>
                </a-select>
              </div>
              <div class="field-desc">
                {{ item.desc }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <SelectWechatApp ref="wxAppRef" title="添加公众号" @ok="handleAppChange" />
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { ref, onMounted, inject, computed, reactive } from 'vue'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { PlusOutlined, CloseCircleFilled } from '@ant-design/icons-vue'
import SelectWechatApp from '@/components/common/select-wechat-app.vue'
import { useWorkflowStore } from '@/stores/modules/workflow'
const workflowStore = useWorkflowStore()

const officialList = computed(() => workflowStore.officialList)
const triggerOfficialList = computed(() => workflowStore.triggerOfficialList)

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

const trigger_options = ref([
  {
    label: '私信消息',
    value: 'message'
  },
  {
    label: '关注/取消关注事件',
    value: 'subscribe_unsubscribe'
  },
  {
    label: '扫描带参数二维码事件',
    value: 'qrcode_scan'
  },
  {
    label: '自定义菜单事件',
    value: 'menu_click'
  }
])

const formState = reactive({
  msg_type: '',
  app_ids: []
})

const selectAppItems = computed(() => {
  return formState.app_ids.map((item) => officialList.value.find((it) => it.app_id == item))
})

const list = ref([])
const options = ref([])

const selectedValues = computed(() => {
  return list.value.map((item) => item.variable)
})

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
    getGraph().eventCenter.emit('custom:trigger-change', { ...props.node })
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
  node_params.trigger.trigger_official_config = {
    msg_type: formState.msg_type,
    app_ids: formState.app_ids.join(',')
  }

  let data = { ...props.node, node_params: JSON.stringify(node_params) }

  emit('update-node', data)

  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', data)
  }, 10)
}

const init = () => {
  getOptions()

  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'

    dataRaw = JSON.parse(dataRaw)

    const trigger = dataRaw.trigger || {
      outputs: []
    }

    let trigger_official_config = trigger.trigger_official_config

    formState.msg_type = trigger_official_config.msg_type
    formState.app_ids = trigger_official_config.app_ids
      ? trigger_official_config.app_ids.split(',')
      : []

    let outputs = trigger.outputs
    if (!outputs) {
      outputs = triggerOfficialList.value.find(
        (item) => item.msg_type == trigger_official_config.msg_type
      ).fields
    }

    list.value = outputs.map((item) => {
      item.tags = item.tags || []

      return item
    })
  } catch (error) {
    console.log(error)
  }
}

const handleOpenAppModal = () => {
  wxAppRef.value.open({}, JSON.parse(JSON.stringify(formState.app_ids)))
}

const wxAppRef = ref(null)
const handleAppChange = (data) => {
  formState.app_ids = data || []
  wxAppRef.value.close()
  update()
}

const handleDelOfficial = (item) => {
  formState.app_ids = formState.app_ids.filter((id) => id != item.app_id)
  update()
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.variable-node {
  .gray-block {
    margin-top: 16px;
  }
  .ml8 {
    margin-left: 8px;
  }
  .app-item-list {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 8px;
    flex-wrap: wrap;
    .app-item {
      height: 32px;
      padding: 5px 16px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      border: 1px solid var(--06, #d9d9d9);
      background: var(--08, #f5f5f5);
      color: #595959;
      font-size: 14px;
      position: relative;
      &:hover {
        .close-icon {
          opacity: 1;
        }
      }
      .close-icon {
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        position: absolute;
        top: -8px;
        right: -12px;
        font-size: 20px;
        color: #fb363f;
        opacity: 0;
        transition: all 0.1s ease-in-out;
        cursor: pointer;
      }
    }
  }
  .form-item {
    margin-top: 12px;
    .form-item-label {
      color: #262626;
      font-size: 14px;
      line-height: 22px;
      margin-bottom: 4px;
    }
  }
  .field-items {
    .field-items-label {
      margin-bottom: 4px;
      color: #262626;
      line-height: 22px;
    }
    .field-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      &:last-child {
        margin-bottom: 0;
      }
    }

    .field-name-box {
      width: auto;
      margin-right: 8px;
    }

    .field-value-box {
      flex: 1;
      margin-right: 8px;

      .field-value {
        display: inline-flex;
        line-height: 20px;
        padding: 1px 8px;
        border-radius: 6px;
        overflow: hidden;
        background: #fff;
        border: 1px solid rgba(0, 0, 0, 0.15);
      }

      .value-arrow {
        font-size: 16px;
        padding: 1px 4px;
        margin-right: 4px;
        border-radius: 4px;
        background: #e4e6eb;
      }
    }

    .field-desc {
      line-height: 22px;
      font-size: 14px;
      color: #595959;
      text-align: left;
    }
  }
}
</style>
