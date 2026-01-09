<template>
  <div class="user-model-page">
    <PageTabs v-model:value="activeTab" @change="onChangeTab" />
    <div class="list-wrapper">
      <div class="alert-content-box" v-if="activeTab == 0">
        <a-alert show-icon>
          <template #message>
            <div class="title">{{ t('views.user.model.usage_title') }}</div>
            <div>
              {{ t('views.user.model.usage_point_1') }}
            </div>
            <div>{{ t('views.user.model.usage_point_2') }}</div>
            <div>
              {{ t('views.user.model.usage_point_3') }}
            </div>
          </template>
        </a-alert>
      </div>

      <div class="content-body-box" v-if="activeTab == 1">
        <div class="model-group-box">
          <div
            class="model-list-item"
            @click="handleChangeDefine(item)"
            :class="{ active: item.config_info_id == currentDefine }"
            v-for="item in hasAddModalList"
            :key="item.config_info_id"
          >
            <img class="avatar" :src="item.model_icon_url" alt="" />
            <div class="model-name">{{ item.config_info.config_name || item.model_name }}</div>
          </div>
        </div>
        <div class="model-content-box">
          <div v-if="currentDefine == 'chatwiki'" class="chatwiki-model-box">
            <div class="statics-block">
              <div class="statics-item">
                <div class="title">{{ t('views.user.model.remaining_points') }}</div>
                <div class="num">{{ formatPriceWithCommas(staticsData.all_surplus) }}</div>
              </div>
              <div class="statics-item">
                <div class="title">{{ t('views.user.model.total_used') }}</div>
                <div class="num">{{ formatPriceWithCommas(staticsData.all_use) }}</div>
              </div>
              <div class="statics-item">
                <div class="title">{{ t('views.user.model.nearest_expiry') }}</div>
                <div class="num">{{ staticsData.min_expiretime }}</div>
              </div>
              <div class="btn-block">
                <a-button type="primary" @click="openBuyPointsModal" block>{{ t('views.user.model.buy_resource_pack') }}</a-button>
                <a-button @click="handleShowBuyRecord" block>{{ t('views.user.model.buy_record') }}</a-button>
              </div>
            </div>
            <div class="search-block">
              <a-segmented v-model:value="model_type" :options="typeOptions" />
              <div class="search-item">
                <span>{{ t('views.user.model.model_supplier') }}：</span>
                <a-select
                  allowClear
                  @change="getSelfList"
                  v-model:value="searchState.model_supplier"
                  :placeholder="t('views.user.model.select_model_supplier')"
                  style="width: 172px"
                >
                  <a-select-option value="tongyi">{{ t('views.user.model.tongyi') }}</a-select-option>
                  <a-select-option value="deepseek">{{ t('views.user.model.deepseek') }}</a-select-option>
                  <a-select-option value="doubao">{{ t('views.user.model.doubao') }}</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="list-box">
              <a-table sticky style="margin-top: 8px" :pagination="false" :data-source="filterSelfLists">
                <a-table-column
                  key="uni_model_name"
                  :width="140"
                  :title="t('views.user.model.model')"
                  data-index="uni_model_name"
                />
                <a-table-column
                  :width="140"
                  key="model_supplier_desc"
                  :title="t('views.user.model.model_provider')"
                  data-index="model_supplier_desc"
                >
                </a-table-column>
                <a-table-column
                  v-if="model_type == 'LLM'"
                  :width="140"
                  key="thinking_type"
                  :title="t('views.user.model.deep_thinking')"
                  data-index="thinking_type"
                >
                  <template #default="{ record }">
                    <span v-if="record.thinking_type == 0">{{ t('views.user.model.not_supported') }}</span>
                    <span v-if="record.thinking_type == 1">{{ t('views.user.model.supported') }}</span>
                    <span v-if="record.thinking_type == 2">{{ t('views.user.model.optional') }}</span>
                  </template>
                </a-table-column>
                <a-table-column
                  v-if="model_type == 'LLM'"
                  :width="140"
                  key="function_call"
                  :title="t('views.user.model.tool_call')"
                  data-index="function_call"
                >
                  <template #default="{ record }">
                    <span v-if="record.function_call == 0">{{ t('views.user.model.not_supported') }}</span>
                    <span v-if="record.function_call == 1">{{ t('views.user.model.supported') }}</span>
                  </template>
                </a-table-column>
                <a-table-column
                  :width="120"
                  v-if="model_type == 'TEXT EMBEDDING'"
                  key="vector_dimension_list"
                  :title="t('views.user.model.vector_dimension')"
                  data-index="vector_dimension_list"
                />
                <a-table-column
                  v-if="model_type != 'RERANK'"
                  :width="100"
                  key="input_desc"
                  :title="t('views.user.model.input')"
                  data-index="input_desc"
                >
                  <template #default="{ record }">
                    {{ record.input_desc }}
                  </template>
                </a-table-column>
                <a-table-column
                  v-if="model_type == 'LLM'"
                  :width="100"
                  key="output_desc"
                  :title="t('views.user.model.output')"
                  data-index="output_desc"
                >
                  <template #default="{ record }">
                    {{ record.output_desc }}
                  </template>
                </a-table-column>
                <a-table-column :width="130" key="price" :title="t('views.user.model.price')" data-index="price">
                  <template #default="{ record }"> {{ record.price }}{{ t('views.user.model.points_per_token') }}</template>
                </a-table-column>
              </a-table>
            </div>
          </div>
          <template v-else>
            <HasModalList
              @addModel="handleOpenAddModelNew"
              @delModel="handleDelModelNew"
              @editConfig="handleEditModel"
              @delConfig="handleRemoveModel"
              :currentModalItem="currentModalItem"
            />
          </template>
        </div>
      </div>
      <!-- 新手指引 -->
       <div class="guide-box" v-if="activeTab == 2">
        <BeginnerGuide></BeginnerGuide>
       </div>
      <ModelList
        :list="canAddModelList"
        :type="2"
        key="canAddModelList"
        @add="handleAddModel"
        @see="handleSeeModel"
        v-if="activeTab == 0"
      />
    </div>
  </div>
  <!-- 查看模型 -->
  <SeeModelAlert
    ref="seeModelAlertRef"
    @new="handleNewModel"
    @edit="handleEditModel"
    @remove="handleRemoveModel"
  />
  <!-- 设置模型 -->
  <SetModelAlert ref="setModelAlertRef" @ok="configSaveOk" />
  <BuyRecord :list="staticsData.list" ref="buyRecordRef" />
  <KefuModal ref="kefuModalRef" />
  <!-- <BuyPointsModal ref="buyPointsModalRef" @ok="handleBuyPoints" /> -->
  <AddModelNew ref="addModelNewRef" @ok="saveModelSuccess" />
</template>

<script setup>
import { ref, computed, createVNode, reactive, watch, onMounted } from 'vue'
import { ExclamationCircleOutlined, PlusOutlined, GlobalOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import {
  getModelConfigList,
  delModelConfig,
  getSelfModelBuylog,
  getSelfModelConfigs,
  showModelConfigList,
  delUseModelConfig
} from '@/api/model/index'
import PageTabs from './components/page-tabs.vue'
import ModelList from './components/model-list.vue'
import SetModelAlert from './components/set-model-alert.vue'
import SeeModelAlert from './components/see-model-alert.vue'
import { getModelTableConfig } from './model-config'
import BuyRecord from './components/buy-record.vue'
import KefuModal from './components/kefu-modal.vue'
import dayjs from 'dayjs'
import { formatPriceWithCommas } from '@/utils/index'
import { useUserStore } from '@/stores/modules/user'
import HasModalList from './components/has-modal-list.vue'
import AddModelNew from './components/add-model-new.vue'
import BeginnerGuide from '@/components/beginner-guide/index.vue'
import { useRoute } from 'vue-router'

const query = useRoute().query
const { t } = useI18n()
const activeTab = ref(+query.activeTab >= 0 ? +query.activeTab : 1)

const onChangeTab = () => {}
// 获取模型列表
const modelList = ref([])
const selfLists = ref([])

const currentDefine = ref('')
const staticsData = reactive({
  all_surplus: 0,
  all_use: 0,
  min_expiretime: '',
  list: []
})

const searchState = reactive({
  model_supplier: void 0
})

const model_type = ref('LLM')

const filterSelfLists = computed(() => {
  return selfLists.value.filter((item) => item.model_type == model_type.value)
})

// LLM  TEXT EMBEDDING RERANK
const typeOptions = computed(() => {
  return [
    {
      label: `${t('views.user.model.llm_model')}（${selfLists.value.filter((item) => item.model_type == 'LLM').length}）`,
      value: 'LLM'
    },
    {
      label: `${t('views.user.model.embedding_model')}（${selfLists.value.filter((item) => item.model_type == 'TEXT EMBEDDING').length}）`,
      value: 'TEXT EMBEDDING'
    },
    {
      label: `${t('views.user.model.rerank_model')}（${selfLists.value.filter((item) => item.model_type == 'RERANK').length}）`,
      value: 'RERANK'
    },
    {
      label: `${t('views.user.model.image_generation_model')}（${selfLists.value.filter((item) => item.model_type == 'IMAGE').length}）`,
      value: 'IMAGE'
    },
    {
      label: `${t('views.user.model.tts_model')}（${selfLists.value.filter((item) => item.model_type == 'TTS').length}）`,
      value: 'TTS'
    }
  ]
})

const hasAddModalList = ref([])

const getLogs = () => {
  getSelfModelBuylog().then((res) => {
    staticsData.all_surplus = res.data.all_surplus
    staticsData.all_use = res.data.all_use
    staticsData.list = res.data.list || []
    staticsData.min_expiretime =
      res.data.min_expiretime > 0
        ? dayjs(res.data.min_expiretime * 1000).format('YYYY-MM-DD')
        : '--'
  })
}

const kefuModalRef = ref(null)
const buyRecordRef = ref(null)
const handleShowBuyRecord = () => {
  getLogs()
  buyRecordRef.value.show()
}
const model_supplier_maps = {
  tongyi: t('views.user.model.tongyi'),
  deepseek: t('views.user.model.deepseek'),
  doubao: t('views.user.model.doubao')
}
const model_type_maps = {
  LLM: t('views.user.model.llm_model'),
  'TEXT EMBEDDING': t('views.user.model.embedding_model'),
  RERANK: t('views.user.model.rerank_model'),
  IMAGE: t('views.user.model.image_generation_model')
}

const input_map = {
  input_text: t('views.user.model.text'),
  input_voice: t('views.user.model.voice'),
  input_image: t('views.user.model.image'),
  input_video: t('views.user.model.video'),
  input_document: t('views.user.model.document')
}

const output_map = {
  output_text: t('views.user.model.text'),
  output_voice: t('views.user.model.voice'),
  output_image: t('views.user.model.image'),
  output_video: t('views.user.model.video')
}

const getSelfList = () => {
  getSelfModelConfigs({
    ...searchState
  }).then((res) => {
    if (!res.data) {
      selfLists.value = []
      return
    }

    selfLists.value = res.data.map((item) => {
      let input_desc = []
      let output_desc = []
      for (let key in input_map) {
        if (item[key] == 1) {
          input_desc.push(input_map[key])
        }
      }
      for (let key in output_map) {
        if (item[key] == 1) {
          output_desc.push(output_map[key])
        }
      }
      return {
        ...item,
        model_supplier_desc: model_supplier_maps[item.model_supplier],
        model_type_desc: model_type_maps[item.model_type],
        input_desc: input_desc.join(','),
        output_desc: output_desc.join(',')
      }
    })
  })
}

const currentModalItem = computed(() => {
  return hasAddModalList.value.filter((item) => item.config_info_id == currentDefine.value)[0] || {}
})

const handleChangeDefine = (item) => {
  currentDefine.value = item.config_info_id
}

watch(
  () => hasAddModalList,
  () => {
    if (
      hasAddModalList.value.filter((item) => item.config_info_id == currentDefine.value).length == 0
    ) {
      currentDefine.value = hasAddModalList.value[0].config_info_id
    }
  },
  {
    deep: true
  }
)

const canAddModelList = ref([])

const getModelList = () => {
  getModelConfigList().then((res) => {
    let list = res.data || []

    list.forEach((item) => {
      item.config_info_id = item.model_define == 'chatwiki' ? 'chatwiki' : item.config_info.id
    })
    hasAddModalList.value = list
  })
}

const addModelNewRef = ref(null)
const handleOpenAddModelNew = (data, record) => {
  addModelNewRef.value.show(data, record)
}

const handleDelModelNew = (record) => {
  Modal.confirm({
    title: t('views.user.model.delete_model_confirm'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('views.user.model.confirm_delete_model', { name: record.show_model_name }),
    onOk() {
      delUseModelConfig({ id: record.id }).then(() => {
        message.success(t('views.user.model.delete_success'))
        getModelList()
        getCanAddModelList()
      })
    }
  })
}

const getCanAddModelList = () => {
  showModelConfigList().then((res) => {
    canAddModelList.value = res.data || []
  })
}

onMounted(() => {
  getModelList()
  // getSelfList()
  getCanAddModelList()
})

// 查看模型
const seeModelAlertRef = ref(null)
const handleSeeModel = (model) => {
  seeModelAlertRef.value.open(model)
}

// 新增模型
const handleNewModel = (model) => {
  setModelAlertRef.value.open(model)
}

// 添加模型
const handleAddModel = (model) => {
  setModelAlertRef.value.open(model)
}

// 修改模型配置
const setModelAlertRef = ref(null)
const handleEditModel = (model, record) => {
  setModelAlertRef.value.open(model, record)
}

const saveModelSuccess = () => {
  getModelList()
  getCanAddModelList()
}

const configSaveOk = (id) => {
  if (activeTab.value == 0) {
    activeTab.value = 1
    currentDefine.value = id
  }
  saveModelSuccess()
}

// 删除模型
const handleRemoveModel = (data) => {
  Modal.confirm({
    title: t('views.user.model.deleteModel'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('views.user.model.deleteModelContent'),
    onOk() {
      delModelConfig({ id: data.id }).then(() => {
        message.success(t('common.deleteSuccessful'))
        getModelList()
        getCanAddModelList()
      })
    }
  })
}

const buyPointsModalRef = ref(null)
const openBuyPointsModal = () => {
  buyPointsModalRef.value.open()
}
const handleBuyPoints = ({ points }) => {
  const userStore = useUserStore()
  let token = userStore.getToken
  let url = `/manage/createOrder?token=${token}&buy_unit=${points}`

  window.location.href = url
}
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;

  .list-wrapper {
    height: calc(100% - 47px);
    overflow-x: hidden;
    overflow-y: auto;
  }
}
.alert-content-box {
  padding: 16px 24px 0 24px;
  .ant-alert {
    align-items: baseline;
  }
  .title {
    font-size: 15px;
    font-weight: 600;
    color: #333;
  }
}

.content-body-box {
  height: 100%;
  background: #fff;
  overflow: hidden;
  display: flex;
  border-top: 1px solid #f0f0f0;
  .model-group-box {
    width: 236px;
    overflow-x: hidden;
    overflow-y: auto;
    padding: 8px;
    scrollbar-width: none;
    &::-webkit-scrollbar {
      display: none;
    }
    .model-list-item {
      height: 42px;
      display: flex;
      align-items: center;
      padding: 0 8px;
      border-radius: 6px;
      gap: 8px;
      cursor: pointer;
      line-height: 22px;
      font-size: 14px;
      .model-name {
        flex: 1;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        color: #595959;
      }
      .avatar {
        width: 22px;
        // height: 24px;
      }
      &:hover {
        background-color: #e4e6eb;
      }
      &.active {
        background-color: #e6efff;
      }
    }
  }
  .model-content-box {
    border-left: 1px solid #f0f0f0;
    flex: 1;
    height: 100%;
    overflow: hidden;
    font-size: 14px;
    line-height: 22px;
    .chatwiki-model-box {
      padding: 24px;
      height: 100%;
      overflow: hidden;
      padding-bottom: 0;
      display: flex;
      flex-direction: column;
      padding-right: 0;
    }
    .statics-block {
      display: flex;
      gap: 24px;
      .statics-item {
        width: 208px;
        height: 90px;
        border-radius: 6px;
        background: #f2f4f7;
        padding: 16px 24px;
        .title {
          color: #7a8699;
        }
        .num {
          line-height: 32px;
          font-weight: 600;
          font-size: 24px;
          color: #242933;
          margin-top: 4px;
        }
      }
    }
    .btn-block {
      width: 104px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 8px;
    }

    .search-block {
      color: #262626;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding-right: 24px;
      margin-top: 24px;
      gap: 24px;
    }
    .list-box {
      margin-top: 8px;
      flex: 1;
      overflow-y: auto;
      &::v-deep(.ant-segmented) {
        color: #262626;
        .ant-segmented-item-selected {
          color: #2475fc;
        }
      }
    }

    .aleat-box {
      .ant-alert-info {
        background: var(--01-, #e5efff);
        border: 1px solid var(--01-, #659dfc);
      }
      .text {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #3a4559;
        line-height: 22px;
        font-size: 14px;
      }
      .alert-description {
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
      .btn-box {
        display: flex;
        align-items: center;
        gap: 8px;
      }
    }
    .message-total-tip {
      margin-top: 16px;
      color: #262626;
      line-height: 22px;
    }
  }
}
.guide-box{
  // width: 100%;
  // height: 100%;
  // overflow-y: auto;
}
</style>
