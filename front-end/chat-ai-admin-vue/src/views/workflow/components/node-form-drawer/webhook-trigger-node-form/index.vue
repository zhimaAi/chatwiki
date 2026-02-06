<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        @changeTitle="handleTitleChange"
        @deleteNode="handleDeleteNode"
        :desc="t('desc_webhook_trigger')"
      >
      </NodeFormHeader>
    </template>

    <div class="variable-node">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <a-form-item :label="t('label_request_url')">
              <a-input-group compact>
                <a-select v-model:value="formState.method" style="width: 20%">
                  <a-select-option value="POST">POST</a-select-option>
                  <a-select-option value="GET">GET</a-select-option>
                </a-select>
                <a-input readonly v-model:value="formState.url" style="width: 65%" />
                <a-button @click="handleCopyText()" style="width: 15%" type="primary"
                  >{{ t('btn_copy') }}</a-button
                >
              </a-input-group>
            </a-form-item>

            <a-form-item :label="t('label_auth')">
              <a-radio-group v-model:value="formState.switch_verify">
                <a-radio value="1">{{ t('label_need_auth') }}</a-radio>
                <a-radio value="0">{{ t('label_no_auth') }}</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item :label="t('label_ip_whitelist')">
              <a-radio-group v-model:value="formState.switch_allow_ip">
                <a-radio value="0">{{ t('label_not_enable') }}</a-radio>
                <a-radio value="1">{{ t('label_enable') }}</a-radio>
              </a-radio-group>
              <div class="mt8">
                <a-textarea
                  v-model:value="formState.allow_ips"
                  style="min-height: 74px"
                  :placeholder="t('ph_input_ip')"
                />
              </div>
            </a-form-item>

            <div class="array-form-box">
              <div class="form-item-label">{{ t('label_params') }}</div>
              <div class="form-item-list" v-for="(item, index) in formState.params" :key="index">
                <a-form-item :label="null">
                  <div class="flex-block-item">
                    <a-form-item-rest>
                      <a-input
                        style="width: 120px"
                        v-model:value="item.key"
                        :placeholder="t('ph_input_key')"
                      ></a-input>
                    </a-form-item-rest>
                    <div class="at-input-flex1">
                      <a-select
                        style="width: 100%"
                        :placeholder="t('ph_select_variable')"
                        v-model:value="item.desc"
                        @dropdownVisibleChange="dropdownVisibleChange"
                      >
                        <a-select-option
                          :value="opt.value"
                          v-for="opt in strGlobalOptions"
                          :key="opt.key"
                        >
                          <span>{{ opt.label }}</span>
                        </a-select-option>
                      </a-select>
                    </div>

                    <div class="btn-hover-wrap" @click="onDelParams(index)">
                      <CloseCircleOutlined />
                    </div>
                  </div>
                </a-form-item>
              </div>
              <a-button @click="handleAddParams" :icon="h(PlusOutlined)" block type="dashed"
                >{{ t('btn_add_param') }}</a-button
              >
            </div>
            <template v-if="formState.method == 'POST'">
              <a-form-item :label="t('label_body')">
                <a-radio-group v-model:value="formState.request_content_type">
                  <a-radio value="none">none</a-radio>
                  <a-radio value="multipart/form-data">form-data</a-radio>
                  <a-radio value="application/x-www-form-urlencoded">x-www-form-urlencoded</a-radio>
                  <a-radio value="application/json">JSON</a-radio>
                </a-radio-group>
              </a-form-item>
              <div
                class="array-form-box"
                v-if="formState.request_content_type == 'multipart/form-data'"
              >
                <div class="array-form-box">
                  <div class="th-header-block">
                    <div class="td-title" style="width: 180px">{{ t('label_key') }}</div>
                    <div class="td-title" style="width: 90px">{{ t('label_value_type') }}</div>
                    <div class="td-title" style="width: 160px">{{ t('label_all_variables') }}</div>
                    <div class="td-title" style="flex: 1">{{ t('label_operation') }}</div>
                  </div>
                  <div class="form-item-list" v-for="(item, index) in formState.form" :key="index">
                    <a-form-item :label="null">
                      <div class="key-block-item">
                        <div class="key-item" style="width: 180px">{{ item.key }}</div>
                        <div class="key-item" style="width: 90px">{{ item.typ }}</div>
                        <div class="key-item" style="width: 160px">
                          <a-select
                            :placeholder="t('ph_select_variable')"
                            v-model:value="item.desc"
                            style="width: 90%"
                            @dropdownVisibleChange="dropdownVisibleChange"
                          >
                            <a-select-option
                              :value="opt.value"
                              v-for="opt in handleFilterOptions(item)"
                              :key="opt.key"
                            >
                              <span>{{ opt.label }}</span>
                            </a-select-option>
                          </a-select>
                        </div>
                        <div class="key-item btn-block" style="flex: 1">
                          <div
                            class="btn-hover-wrap"
                            v-if="item.typ == 'object'"
                            @click="onAddSubs(index)"
                          >
                            <PlusCircleOutlined />
                          </div>
                          <div class="btn-hover-wrap" @click="handleEditKey(item, index)">
                            <EditOutlined />
                          </div>
                          <div class="btn-hover-wrap" @click="onDelBody(index, 'form')">
                            <CloseCircleOutlined />
                          </div>
                        </div>
                      </div>
                      <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
                        <a-form-item-rest>
                          <SubKey
                            :data="item.subs"
                            :request_content_type="formState.request_content_type"
                            :level="2"
                            :globalOptions="globalOptions"
                          />
                        </a-form-item-rest>
                      </div>
                    </a-form-item>
                  </div>
                </div>
                <div class="mt8">
                  <a-button @click="handleAddKeyModal" :icon="h(PlusOutlined)" block type="dashed"
                    >{{ t('label_add_param') }}</a-button
                  >
                </div>
              </div>
              <div
                class="array-form-box"
                v-if="formState.request_content_type == 'application/x-www-form-urlencoded'"
              >
                <div class="array-form-box">
                  <div class="th-header-block">
                    <div class="td-title" style="width: 180px">{{ t('label_key') }}</div>
                    <div class="td-title" style="width: 90px">{{ t('label_value_type') }}</div>
                    <div class="td-title" style="width: 160px">{{ t('label_all_variables') }}</div>
                    <div class="td-title" style="flex: 1">{{ t('label_operation') }}</div>
                  </div>
                  <div
                    class="form-item-list"
                    v-for="(item, index) in formState.x_form"
                    :key="index"
                  >
                    <a-form-item :label="null">
                      <div class="key-block-item">
                        <div class="key-item" style="width: 180px">{{ item.key }}</div>
                        <div class="key-item" style="width: 90px">{{ item.typ }}</div>
                        <div class="key-item" style="width: 160px">
                          <a-select
                            :placeholder="t('ph_select_variable')"
                            v-model:value="item.desc"
                            style="width: 90%"
                            @dropdownVisibleChange="dropdownVisibleChange"
                          >
                            <a-select-option
                              :value="opt.value"
                              v-for="opt in handleFilterOptions(item)"
                              :key="opt.key"
                            >
                              <span>{{ opt.label }}</span>
                            </a-select-option>
                          </a-select>
                        </div>
                        <div class="key-item btn-block" style="flex: 1">
                          <div class="btn-hover-wrap" @click="handleEditKey(item, index)">
                            <EditOutlined />
                          </div>
                          <div class="btn-hover-wrap" @click="onDelBody(index, 'x_form')">
                            <CloseCircleOutlined />
                          </div>
                        </div>
                      </div>
                    </a-form-item>
                  </div>
                </div>
                <div class="mt8">
                  <a-button @click="handleAddKeyModal" :icon="h(PlusOutlined)" block type="dashed"
                    >{{ t('label_add_param') }}</a-button
                  >
                </div>
              </div>
              <div
                class="array-form-box"
                v-if="formState.request_content_type == 'application/json'"
              >
                <div class="array-form-box">
                  <div class="th-header-block">
                    <div class="td-title" style="width: 180px">{{ t('label_key') }}</div>
                    <div class="td-title" style="width: 90px">{{ t('label_value_type') }}</div>
                    <div class="td-title" style="width: 160px">{{ t('label_all_variables') }}</div>
                    <div class="td-title" style="flex: 1">{{ t('label_operation') }}</div>
                  </div>
                  <div class="form-item-list" v-for="(item, index) in formState.json" :key="index">
                    <a-form-item :label="null">
                      <div class="key-block-item">
                        <div class="key-item" style="width: 180px">{{ item.key }}</div>
                        <div class="key-item" style="width: 90px">{{ item.typ }}</div>
                        <div class="key-item" style="width: 160px">
                          <a-select
                            v-if="item.typ != 'object'"
                            :placeholder="t('ph_select_variable')"
                            v-model:value="item.desc"
                            style="width: 90%"
                            @dropdownVisibleChange="dropdownVisibleChange"
                          >
                            <a-select-option
                              :value="opt.value"
                              v-for="opt in handleFilterOptions(item)"
                              :key="opt.key"
                            >
                              <span>{{ opt.label }}</span>
                            </a-select-option>
                          </a-select>
                        </div>
                        <div class="key-item btn-block" style="flex: 1">
                          <div
                            class="btn-hover-wrap"
                            v-if="item.typ == 'object'"
                            @click="onAddSubs(index)"
                          >
                            <PlusCircleOutlined />
                          </div>
                          <div class="btn-hover-wrap" @click="handleEditKey(item, index)">
                            <EditOutlined />
                          </div>
                          <div class="btn-hover-wrap" @click="onDelBody(index, 'json')">
                            <CloseCircleOutlined />
                          </div>
                        </div>
                      </div>
                      <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
                        <a-form-item-rest>
                          <SubKey
                            :data="item.subs"
                            :request_content_type="formState.request_content_type"
                            :level="2"
                            :globalOptions="globalOptions"
                          />
                        </a-form-item-rest>
                      </div>
                    </a-form-item>
                  </div>
                </div>
                <div class="mt8">
                  <a-button @click="handleAddKeyModal" :icon="h(PlusOutlined)" block type="dashed"
                    >{{ t('label_add_param') }}</a-button
                  >
                </div>
              </div>
            </template>
          </div>
          <div class="gray-block mt16">
            <a-form-item :label="t('label_response')">
              <a-radio-group v-model:value="formState.response_type">
                <a-radio value="message_variable">{{ t('label_return_after_end') }}</a-radio>
                <a-radio value="now">{{ t('label_return_now') }}</a-radio>
              </a-radio-group>
              <div class="mt8" v-if="formState.response_type == 'now'">
                <a-textarea
                  v-model:value="formState.response_now"
                  style="min-height: 74px"
                  :placeholder="t('ph_input_response')"
                />
              </div>
            </a-form-item>
          </div>
        </a-form>
        <AddKeyModal
          :request_content_type="formState.request_content_type"
          @add="onKeyAdd"
          @edit="onKeyEdit"
          @addSub="onKeyAddSub"
          ref="addKeyModalRef"
        />
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { getUuid } from '@/utils/index'
import { ref, onMounted, inject, reactive, computed, watch, h } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  PlusCircleOutlined,
  EditOutlined
} from '@ant-design/icons-vue'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import AddKeyModal from './add-key-modal.vue'
import SubKey from './subs-key.vue'
import { copyText } from '@/utils/index'
import { message } from 'ant-design-vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.webhook-trigger-node-form.index')

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

const formState = reactive({
  url: '',
  method: 'GET',
  switch_verify: '1',
  switch_allow_ip: '0',
  allow_ips: '',

  params: [
    // {
    //   key: '',
    //   typ: 'string',
    //   desc: void 0
    // }
  ],
  form: [
    // {
    //   key: '',
    //   typ: 'string',
    //   desc: void 0
    // }
  ],

  x_form: [
    // {
    //   key: '',
    //   typ: 'string',
    //   desc: void 0
    // }
  ],

  json: [
    // {
    //   key: '',
    //   typ: 'string',
    //   desc: void 0,
    //   subs: []
    // }
  ],

  request_content_type: 'none',

  response_type: 'now',
  response_now: ''
})

watch(
  () => formState,
  () => {
    update()
  },
  {
    deep: true
  }
)

const getNode = inject('getNode')
const getGraph = inject('getGraph')

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

  node_params.trigger.trigger_web_hook_config = {
    ...formState
  }

  let data = { ...props.node, node_params: JSON.stringify(node_params) }

  emit('update-node', data)

  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', data)
  }, 10)
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'

    dataRaw = JSON.parse(dataRaw)
    getOptions()
    let trigger_web_hook_config = dataRaw.trigger.trigger_web_hook_config
    formState.url = trigger_web_hook_config.url
    formState.method = trigger_web_hook_config.method || 'GET'
    formState.switch_verify = trigger_web_hook_config.switch_verify || '1'
    formState.switch_allow_ip = trigger_web_hook_config.switch_allow_ip || '0'
    formState.allow_ips = trigger_web_hook_config.allow_ips
      ? trigger_web_hook_config.allow_ips.split(',').join('\n')
      : ''
    formState.params = trigger_web_hook_config.params || []
    formState.form = trigger_web_hook_config.form || []
    formState.x_form = trigger_web_hook_config.x_form || []
    formState.json = trigger_web_hook_config.json || []
    formState.request_content_type = trigger_web_hook_config.request_content_type || 'none'
    formState.response_type = trigger_web_hook_config.response_type || 'now'
    formState.response_now = trigger_web_hook_config.response_now || ''
  } catch (error) {
    console.log(error)
  }
}

const globalOptions = ref([])

function getOptions() {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let globalVariable = nodeModel.getGlobalVariable()
    let diy_global = globalVariable.diy_global || []
    diy_global.forEach((item) => {
      item.label = item.key
      item.value = 'global.' + item.key
    })

    globalOptions.value = diy_global || []
  }
}

const strGlobalOptions = computed(()=>{
  return globalOptions.value.filter(item => item.typ == 'string')
})

const dropdownVisibleChange = (visible) => {
  if (visible) {
    getOptions()
  }
}

const handleFilterOptions = (item) => {
  let typ = item.typ
  if (typ == 'file') {
    typ = 'string'
  }
  if (typ) {
    return globalOptions.value.filter((item) => item.typ == typ)
  }
  return globalOptions.value
}

const handleAddParams = () => {
  formState.params.push({
    key: '',
    typ: 'string',
    desc: void 0
  })
}

const onDelParams = (index) => {
  formState.params.splice(index, 1)
}

const onDelBody = (index, key) => {
  formState[key].splice(index, 1)
}

const addKeyModalRef = ref(null)

const handleAddKeyModal = () => {
  addKeyModalRef.value.add()
}

const onAddSubs = (index) => {
  addKeyModalRef.value.addSub(index)
}

const handleEditKey = (data, index) => {
  addKeyModalRef.value.edit(data, index)
}

const onKeyAdd = (data) => {
  if (formState.request_content_type == 'multipart/form-data') {
    formState.form.push(data)
  }
  if (formState.request_content_type == 'application/x-www-form-urlencoded') {
    formState.x_form.push(data)
  }
  if (formState.request_content_type == 'application/json') {
    formState.json.push(data)
  }
}

const onKeyEdit = (data, index) => {
  if (formState.request_content_type == 'multipart/form-data') {
    formState.form.splice(index, 1, data)
  }
  if (formState.request_content_type == 'application/x-www-form-urlencoded') {
    formState.x_form.splice(index, 1, data)
  }
  if (formState.request_content_type == 'application/json') {
    formState.json.splice(index, 1, data)
  }
}

const onKeyAddSub = (data, index) => {
  if (formState.request_content_type == 'multipart/form-data') {
    if(!formState.form[index].subs){
      formState.form[index].subs = []
    }
    formState.form[index].subs.push(data)
  }
  if (formState.request_content_type == 'application/x-www-form-urlencoded') {
    if(!formState.x_form[index].subs){
      formState.x_form[index].subs = []
    }
    formState.x_form[index].subs.push(data)
  }
  if (formState.request_content_type == 'application/json') {
    if(!formState.json[index].subs){
      formState.json[index].subs = []
    }
    formState.json[index].subs.push(data)
  }
}

const handleCopyText = () => {
  copyText(formState.url)
  message.success(t('msg_copy_success'))
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.at-input-flex1 {
  flex: 1;
  overflow: hidden;
}

.th-header-block {
  display: flex;
  align-items: center;
  .td-title {
    font-size: 14px;
    color: #262626;
  }
}

.key-block-item {
  display: flex;
  align-items: center;
  margin-top: 12px;
  .key-item {
    font-size: 14px;
    color: #262626;
  }
  .btn-block {
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .btn-hover-wrap {
    width: 24px;
    height: 24px;
  }
}
</style>
