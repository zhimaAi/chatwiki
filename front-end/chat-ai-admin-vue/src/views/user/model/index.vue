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
    overflow: auto;
    padding: 24px;
    font-size: 14px;
    line-height: 22px;
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
      margin-top: 24px;
      gap: 24px;
    }
    .list-box {
      margin-top: 8px;
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
</style>

<template>
  <div class="user-model-page">
    <PageTabs v-model:value="activeTab" @change="onChangeTab" />
    <div class="list-wrapper">
      <div class="alert-content-box" v-if="activeTab == 0">
        <a-alert show-icon>
          <template #message>
            <div class="title">使用说明</div>
            <div>
              1、模型主要分为嵌入模型（TEXT
              EMBEDDING）和大模型（LLM），嵌入模型用于将上传知识库的文本内容进行向量化便于后续检索，大模型用于将检索到的相关知识进行理解后返回答案。
            </div>
            <div>2、配置模型时尽量选择同时支持LLM和TEXT EMBEDDING的模型，否则要配置多个模型。</div>
            <div>
              3、国内模型厂商一般都需要认证后才可使用，所以在配置模型前，确保在模型服务厂商那里已认证，否则配置模型时会提示配置参数错误。
            </div>
          </template>
        </a-alert>
      </div>

      <div class="content-body-box" v-if="activeTab == 1">
        <div class="model-group-box">
          <div
            class="model-list-item"
            @click="handleChangeDefine(item)"
            :class="{ active: item.model_define == currentDefine }"
            v-for="item in addedModelList"
            :key="item.model_define"
          >
            <img class="avatar" :src="item.model_icon_url" alt="" />
            <div class="model-name">{{ item.model_name }}</div>
          </div>
        </div>
        <div class="model-content-box">
          <template v-if="currentDefine == 'chatwiki'">
            <div class="statics-block">
              <div class="statics-item">
                <div class="title">剩余积分</div>
                <div class="num">{{ formatPriceWithCommas(staticsData.all_surplus) }}</div>
              </div>
              <div class="statics-item">
                <div class="title">累计已使用</div>
                <div class="num">{{ formatPriceWithCommas(staticsData.all_use) }}</div>
              </div>
              <div class="statics-item">
                <div class="title">最近过期时间</div>
                <div class="num">{{ staticsData.min_expiretime }}</div>
              </div>
              <div class="btn-block">
                <a-button type="primary" @click="kefuModalRef.show()" block>购买资源包</a-button>
                <a-button @click="handleShowBuyRecord" block>购买记录</a-button>
              </div>
            </div>
            <div class="search-block">
              <div class="search-item">
                <span>模型供应商：</span>
                <a-select
                  allowClear
                  @change="getSelfList"
                  v-model:value="searchState.model_supplier"
                  placeholder="请选择模型供应商"
                  style="width: 172px"
                >
                  <a-select-option value="tongyi">通义千问</a-select-option>
                  <a-select-option value="deepseek">DeepSeek</a-select-option>
                  <a-select-option value="doubao">豆包</a-select-option>
                </a-select>
              </div>
              <div class="search-item">
                <span>模型类型：</span>
                <a-select
                  allowClear
                  @change="getSelfList"
                  v-model:value="searchState.model_type"
                  placeholder="请选择模型类型"
                  style="width: 172px"
                >
                  <a-select-option value="LLM">大语言模型</a-select-option>
                  <a-select-option value="TEXT EMBEDDING">嵌入模型</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="list-box">
              <a-table :data-source="selfLists">
                <a-table-column key="uni_model_name" title="模型" data-index="uni_model_name" />
                <a-table-column
                  key="model_supplier_desc"
                  title="模型服务商"
                  data-index="model_supplier_desc"
                >
                </a-table-column>
                <a-table-column key="model_type_desc" title="模型类型" data-index="model_type_desc">
                </a-table-column>
                <a-table-column key="price" title="价格" data-index="price"> 
                  <template #default="{ record }"> {{ record.price}}积分/千Token</template>
                </a-table-column>
              </a-table>
            </div>
          </template>
          <template v-else>
            <template v-if="['xinference', 'azure', 'doubao', 'openaiAgent', 'ollama'].includes(currentDefine)">
              <div class="opration-block">
                <a-button type="primary" @click="handleNewModel(currentModalItem)">
                  <PlusOutlined />添加模型</a-button
                >
              </div>
              <div class="list-box">
                <a-table :pagination="false" :data-source="currentModalItem.config_list">
                  <a-table-column key="deployment_name" data-index="deployment_name">
                    <template #title>
                      <template v-if="currentDefine == 'doubao'">接入点ID</template>
                      <template v-else>模型名称</template>
                    </template>
                  </a-table-column>
                  <a-table-column
                    key="show_model_name"
                    title="模型名称"
                    data-index="show_model_name"
                    v-if="currentDefine == 'doubao'"
                  >
                  </a-table-column>
                  <a-table-column
                    key="model_types"
                    title="模型类型"
                    data-index="model_types"
                  >
                    <template #default="{ record, text }"> {{ model_type_maps[text] }}</template>
                  </a-table-column>
                  <a-table-column
                    v-if="currentDefine == 'azure'"
                    key="api_version"
                    title="API版本"
                    data-index="api_version"
                  >
                  </a-table-column>
                  <a-table-column
                    v-if="['openaiAgent', 'ollama'].includes(currentDefine)"
                    key="api_endpoint"
                    title="API 域名"
                    data-index="api_endpoint"
                  >
                  </a-table-column>
                  <a-table-column :width="120" key="action" title="操作">
                    <template #default="{ record }">
                      <a-flex :gap="16">
                        <a @click="handleEditModel(currentModalItem, record)">编辑</a>
                        <a @click="handleRemoveModel(record)">移除</a>
                      </a-flex>
                    </template>
                  </a-table-column>
                </a-table>
              </div>
            </template>
            <template v-else>
              <div class="aleat-box">
                <a-alert>
                  <template #message>
                    <div class="alert-description">
                      <div
                        v-if="currentModalItem.config_list && currentModalItem.config_list.length"
                        class="text"
                      >
                        <GlobalOutlined />
                        API Key：{{ currentModalItem.config_list[0].api_key.slice(0, 10) + '...' }}
                      </div>
                      <div class="btn-box">
                        <a-button
                          @click="
                            handleEditModel(currentModalItem, currentModalItem.config_list[0])
                          "
                          size="small"
                          >修改配置</a-button
                        >
                        <a-button
                          size="small"
                          @click="handleRemoveModel(currentModalItem.config_list[0])"
                          >移除配置</a-button
                        >
                      </div>
                    </div>
                  </template>
                </a-alert>
              </div>
              <div class="message-total-tip" @click="getModalLists">
                支持的模型（{{ modalConfigs.rows.length }}）
              </div>
              <div class="list-box">
                <a-table :pagination="false" :data-source="modalConfigs.rows">
                  <a-table-column
                    v-for="item in modalConfigs.columns"
                    :key="item.key"
                    :title="item.title"
                    :data-index="item.dataIndex"
                  >
                    <template #default="{ record, text }">
                      <div v-if="item.key == 'model_type'">{{ model_type_maps[text] }}</div>
                      <div v-else>{{ text }}</div>
                    </template>
                  </a-table-column>
                </a-table>
              </div>
            </template>
          </template>
        </div>
      </div>

      <!-- <ModelList
        :list="addedModelList"
        :type="1"
        key="addedModelList"
        v-if="activeTab == 1"
        @edit="handleEditModel"
        @new="handleNewModel"
        @see="handleSeeModel"
        @remove="handleRemoveModel"
      /> -->
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
  <SetModelAlert ref="setModelAlertRef" @ok="saveModelSuccess" />
  <BuyRecord :list="staticsData.list" ref="buyRecordRef" />
  <kefuModal ref="kefuModalRef" />
</template>

<script setup>
import { ref, computed, createVNode, reactive, watch } from 'vue'
import { ExclamationCircleOutlined, PlusOutlined, GlobalOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import {
  getModelConfigList,
  delModelConfig,
  getSelfModelBuylog,
  getSelfModelConfigs
} from '@/api/model/index'
import PageTabs from './components/page-tabs.vue'
import ModelList from './components/model-list.vue'
import SetModelAlert from './components/set-model-alert.vue'
import SeeModelAlert from './components/see-model-alert.vue'
import { getModelTableConfig } from './model-config'
import BuyRecord from './components/buy-record.vue'
import kefuModal from './components/kefu-modal.vue'
import dayjs from 'dayjs'
import { formatPriceWithCommas } from '@/utils/index'

const { t } = useI18n()
const activeTab = ref(1)

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
  model_supplier: void 0,
  model_type: void 0
})

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

// getLogs()
const kefuModalRef = ref(null)
const buyRecordRef = ref(null)
const handleShowBuyRecord = () => {
  getLogs()
  buyRecordRef.value.show()
}
let model_supplier_maps = {
  tongyi: '通义千问',
  deepseek: 'DeepSeek',
  doubao: '豆包'
}
let model_type_maps = {
  LLM: '大语言模型',
  'TEXT EMBEDDING': '嵌入模型',
  RERANK: '重排序模型'
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
      return {
        ...item,
        model_supplier_desc: model_supplier_maps[item.model_supplier],
        model_type_desc: model_type_maps[item.model_type]
      }
    })
  })
}

// getSelfList()

const currentModalItem = computed(() => {
  return modelList.value.filter((item) => item.model_define == currentDefine.value)[0] || {}
})

const modalConfigs = reactive({
  columns: [],
  rows: []
})
const handleChangeDefine = (item) => {
  currentDefine.value = item.model_define
  setTimeout(() => {
    getModalLists()
  }, 50)
}

const getModalLists = () => {
  let tableConfig = getModelTableConfig(currentModalItem.value)
  for (let key in tableConfig) {
    modalConfigs[key] = tableConfig[key]
  }
}

const addedModelList = computed(() => {
  return modelList.value.filter((item) => item.listType == 1)
})

watch(
  () => addedModelList,
  () => {
    if (
      addedModelList.value.filter((item) => item.model_define == currentDefine.value).length == 0
    ) {
      currentDefine.value = addedModelList.value[0].model_define
      handleChangeDefine(addedModelList.value[0])
    }
  },
  {
    deep: true
  }
)

const canAddModelList = computed(() => {
  return modelList.value
})

const getModelList = () => {
  getModelConfigList().then((res) => {
    let list = res.data || []

    list.forEach((item) => {
      item.listType = item.config_list.length > 0 ? 1 : 0
    })

    // 添加成功后移动到底部
    list.sort((a, b) => a.listType - b.listType)
    modelList.value = list
  })
}

getModelList()

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
      })
    }
  })
}
</script>
