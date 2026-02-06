<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_knowledge_base')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">{{ t('label_input') }}</div>

            <a-form-item :label="t('label_knowledge_base')">
              <div class="konwledge-list-item" v-for="item in selectedLibraryRows" :key="item.id">
                <div class="avatar-box">
                  <img :src="item.avatar" alt="" />
                </div>
                <div class="content-box">
                  <div class="list-name">{{ item.library_name }}</div>
                  <div class="list-intro">{{ item.library_intro }}</div>
                </div>
                <div class="btn-hover-wrap" @click="handleDelKonwledge(item)">
                  <CloseCircleOutlined />
                </div>
              </div>
              <div class="btn-block">
                <div>
                  <a-button
                    @click="handleOpenSelectLibraryAlert"
                    :icon="h(PlusOutlined)"
                    block
                    type="dashed"
                    >{{ t('btn_add_knowledge_base') }}</a-button
                  >
                </div>
                <div>
                  <a-button
                    @click="handleOpenRecallSettingsAlert"
                    :icon="h(SettingOutlined)"
                    block
                    type="dashed"
                    >{{ t('btn_recall_settings') }}</a-button
                  >
                </div>
              </div>
            </a-form-item>

            <div class="diy-form-item">
              <div class="form-label">{{ t('label_user_question') }}</div>
              <div class="form-content">
                <div class="form-content">
                  <a-cascader
                    v-model:value="formState.question_value"
                    @dropdownVisibleChange="onDropdownVisibleChange"
                    style="width: 220px"
                    :options="variableOptions"
                    :allowClear="false"
                    :displayRender="({ labels }) => labels.join('/')"
                    :field-names="{ children: 'children' }"
                    :placeholder="t('ph_select')"
                  />
                </div>
              </div>
            </div>
          </div>
          <div class="gray-block mt16">
            <div class="gray-block-title">{{ t('label_output') }}</div>
            <div class="options-item">
              <div class="option-label">{{ t('label_knowledge_quote') }}</div>
              <div class="option-type">string</div>
            </div>
          </div>
        </a-form>
        <LibrarySelectAlert
          ref="librarySelectAlertRef"
          :showWxType="!!wxAppLibary"
          @close="getList"
          @change="onChangeLibrarySelected"
        />
        <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import { ref, reactive, watch, computed, onMounted, h, toRaw } from 'vue'
import { CloseCircleOutlined, PlusOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { getLibraryList } from '@/api/library/index'
import LibrarySelectAlert from './components/library-select-alert.vue'
import RecallSettingsAlert from '../nodes/knowledge-base-node//recall-settings-alert.vue'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.knowledge-base-node-form')
const robotStore = useRobotStore()

const rrf_weight = computed(() => {
  if (robotStore.robotInfo.rrf_weight) {
    return JSON.parse(robotStore.robotInfo.rrf_weight)
  }
  return {
    vector: 0,
    search: 0,
    graph: 0
  }
})

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

const variableOptions = ref([])

function getOptions() {
  // const node = props.lf.getNodeDataById(props.nodeId)
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()

    variableOptions.value = handleOptions(list)
  }
}

// 递归处理Options
function handleOptions(options) {
  options.forEach((item) => {
    if (item.typ == 'node') {
      if (item.node_type == 1) {
        item.value = 'global'
      } else {
        item.value = item.node_id
      }
    } else {
      item.value = item.key
    }

    if (item.children && item.children.length > 0) {
      item.children = handleOptions(item.children)
    }
  })

  return options
}

function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return ['global', 'question']
}

const formRef = ref()
const wxAppLibary = ref(null)
const formState = reactive({
  library_ids: [],
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: void 0,
  top_k: 5,
  similarity: 0.5,
  search_type: 1,
  meta_search_switch: 0,
  meta_search_type: 1,
  meta_search_condition_list: "",
  question_value: '',
  rrf_weight: {},
  recall_neighbor_switch: false,
  recall_neighbor_before_num:1,
  recall_neighbor_after_num: 1,
})

const update = () => {
  const data = JSON.stringify({
    libs: {
      ...formState,
      rerank_model_config_id: formState.rerank_model_config_id
        ? +formState.rerank_model_config_id
        : void 0,
      question_value: formState.question_value.join('.'),
      library_ids: formState.library_ids.join(','),
      rrf_weight: JSON.stringify(formState.rrf_weight),
      recall_neighbor_switch: formState.recall_neighbor_switch,
      recall_neighbor_before_num: formState.recall_neighbor_before_num,
      recall_neighbor_after_num: formState.recall_neighbor_before_num,
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let libs = JSON.parse(dataRaw).libs || {}

    getOptions()

    libs = JSON.parse(JSON.stringify(libs))

    for (let key in libs) {
      if (key == 'library_ids') {
        formState[key] = libs[key] ? libs[key].split(',') : []
      } else if (key == 'question_value') {
        formState.question_value = formatQuestionValue(libs['question_value'])
      } else if (key == 'rrf_weight') {
        formState.rrf_weight = libs[key] ? JSON.parse(libs[key]) : libs[key]
      } else {
        formState[key] = libs[key]
      }
    }

    if (!formState.rrf_weight || Object.keys(formState.rrf_weight).length == 0) {
      //  没有值 则去默认值
      formState.rrf_weight = rrf_weight.value
    }

    // 公众号知识库是否开启
    getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' }).then((res) => {
      let _data = res?.data || {}
      if (_data?.user_config?.switch_status == 1) {
        wxAppLibary.value = _data
      }
    })
  } catch (error) {
    console.log(error)
  }
}

const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    if (!wxAppLibary.value) {
      return formState.library_ids.includes(item.id) && item.type != 3
    } else {
      return formState.library_ids.includes(item.id)
    }
  })
})

// 移除知识库
const handleDelKonwledge = (item) => {
  let index = formState.library_ids.indexOf(item.id)
  formState.library_ids.splice(index, 1)
}

const onChangeLibrarySelected = (checkedList) => {
  getList()
  formState.library_ids = [...checkedList]
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.library_ids])
}

// 召回设置
const recallSettingsAlertRef = ref(null)

const handleOpenRecallSettingsAlert = () => {
  recallSettingsAlertRef.value.open(toRaw(formState), null, formState.library_ids)
}

const onChangeRecallSettings = (data) => {
  formState.rerank_status = data.rerank_status
  formState.rerank_use_model = data.rerank_use_model
  formState.rerank_model_config_id = data.rerank_model_config_id
  formState.top_k = data.top_k
  formState.similarity = data.similarity
  formState.search_type = data.search_type
  formState.meta_search_switch = Number(data.meta_search_switch)
  formState.meta_search_type = Number(data.meta_search_type)
  formState.meta_search_condition_list = data.meta_search_condition_list
  formState.rrf_weight = data.rrf_weight
  formState.recall_neighbor_switch = data.recall_neighbor_switch
  formState.recall_neighbor_before_num = data.recall_neighbor_before_num
  formState.recall_neighbor_after_num = data.recall_neighbor_after_num
}

// 获取知识库
const getList = async () => {
  const res = await getLibraryList({ type: '' })
  if (res) {
    libraryList.value = res.data || []
  }
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
  getList()
})
</script>

<style lang="less" scoped>
@import './form-block.less';
.node-form-content {
  ::v-deep(.ant-form-item-label) {
    padding: 0 0 4px;
  }
}
.konwledge-list-item {
  padding: 14px 16px;
  height: 72px;
  border: 1px solid #d8dde5;
  border-radius: 6px;
  display: flex;
  align-items: center;
  background-color: #fff;
  font-size: 14px;
  line-height: 22px;
  margin-bottom: 4px;
  .avatar-box {
    width: 40px;
    height: 40px;
    border: 1px solid var(--07, #f0f0f0);
    border-radius: 6px;
    margin-right: 8px;
    img {
      width: 100%;
      height: 100%;
    }
  }
  .list-name {
    color: #262626;
    font-weight: 600;
  }
  .list-intro {
    font-size: 12px;
    margin-top: 2px;
    line-height: 20px;
    color: #8c8c8c;
    max-width: 400px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  .content-box {
    flex: 1;
  }
}

.btn-block {
  display: flex;
  gap: 4px;
  & > div {
    flex: 1;
  }
}
</style>
