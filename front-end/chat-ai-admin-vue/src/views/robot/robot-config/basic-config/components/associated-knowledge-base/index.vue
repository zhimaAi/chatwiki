<style lang="less" scoped>
.setting-box {
  .actions-box {
    display: flex;
    align-items: center;
    line-height: 22px;
    font-size: 14px;
    color: #595959;

    .action-btn {
      cursor: pointer;
    }

    .save-btn {
      color: #2475fc;
    }
  }

  .library-list {
    display: flex;
    flex-flow: row wrap;
    gap: 16px;
    padding: 0 16px 0 16px;

    .library-item {
      width: 336px;
      padding: 14px 16px;
      border-radius: 6px;
      border: 1px solid #d9d9d9;
      background-color: #fff;
      cursor: pointer;
      display: flex;
      align-items: center;

      gap: 12px;
      .library-avatar {
        img {
          width: 40px;
          height: 40px;
          border-radius: 12px;
        }
      }
      .info-content-box {
        flex: 1;
        .library-title-block {
          display: flex;
          align-items: center;
          gap: 8px;
          .type-tag {
            height: 18px;
            line-height: 16px;
            padding: 0 4px;
            font-size: 12px;
            font-weight: 400;
            border-radius: 6px;
            color: #2475fc;
            border: 1px solid #99bffd;
            &.gray-tag {
              border: 1px solid var(--06, #d9d9d9);
              color: #8c8c8c;
            }
          }
        }
        .title-text {
          max-width: 150px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          font-size: 14px;
          color: #262626;
          font-weight: 600;
          line-height: 22px;
        }
        .desc-info-block {
          color: #8c8c8c;
          margin-top: 2px;
          font-size: 12px;
          line-height: 20px;
          width: 212px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
  .setting-info-block {
    padding: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    color: #595959;
    line-height: 22px;
    .set-item {
      display: flex;
      align-items: center;
    }
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    title="关联知识库"
    icon-name="guanlianzhishiku"
    v-model:isEdit="isEdit"
    :bodyStyle="{ padding: 0 }"
  >
    <template #extra>
      <div class="actions-box">
        <a-flex :gap="8">
          <a-button size="small" @click="handleOpenRecallSettingsAlert">召回设置</a-button>
          <a-button size="small" @click="handleOpenSelectLibraryAlert">关联知识库</a-button>
        </a-flex>
      </div>
    </template>
    <div class="library-list" v-if="selectedLibraryRows.length > 0">
      <div
        class="library-item"
        v-for="item in selectedLibraryRows"
        :key="item.id"
        @click="toLibraryDetail(item)"
      >
        <div class="library-avatar">
          <img :src="item.avatar" alt="" />
        </div>
        <div class="info-content-box">
          <div class="library-title-block">
            <div class="title-text">{{ item.library_name }}</div>
            <!-- <a-tooltip v-if="neo4j_status">
              <template #title>{{ item.graph_switch == 0 ? '未' : '已' }}开启知识图谱生成</template>
              <div class="graph-tag" :class="{ active: item.graph_switch == 1 }">Graph</div>
            </a-tooltip>
            <span class="type-tag" :class="{ 'gray-tag': item.graph_switch == 0 }">Graph</span> -->
          </div>
          <div class="desc-info-block">{{ item.library_intro }}</div>
        </div>
        <div class="close-btn" v-if="item.id != robotInfo.default_library_id" @click.stop="handleRemoveCheckedLibrary(item)">
          <CloseCircleOutlined />
        </div>
      </div>
    </div>
    <div class="setting-info-block">
      <div class="set-item">
        检索模式：
        <span v-if="formState.search_type == 1">混合检索</span>
        <span v-if="formState.search_type == 2">向量检索</span>
        <span v-if="formState.search_type == 3">全文检索</span>
        <span v-if="formState.search_type == 4">知识图谱检索</span>
      </div>
      <div class="set-item">
        Top K：
        <span>{{ formState.top_k }}</span>
      </div>
      <div class="set-item" v-if="formState.search_type <= 2">
        相似度阈值：
        <span>{{ formState.similarity }}</span>
      </div>
      <div class="set-item">
        Rerank模型：
        <span v-if="formState.rerank_status == 0">关</span>
        <span v-else>{{ getModelName }}</span>
      </div>
    </div>
    <LibrarySelectAlert
      ref="librarySelectAlertRef"
      @change="onChangeLibrarySelected"
      :showWxType="!!wxAppLibary"
    />
    <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
    <NoOpenGraphModal :list="noOpenLibraryList" @refreshList="getList" ref="noOpenGraphModalRef" />
  </edit-box>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { getLibraryList } from '@/api/library/index'
import { ref, reactive, inject, watchEffect, computed, toRaw, onMounted } from 'vue'
import { CloseCircleOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import EditBox from '../edit-box.vue'
import LibrarySelectAlert from './library-select-alert.vue'
import RecallSettingsAlert from './recall-settings-alert.vue'
import NoOpenGraphModal from './no-open-graph-modal.vue'
import { getModelNameText } from '@/components/model-select/index.js'

const { getStorage, removeStorage } = useStorage('localStorage')

import { useCompanyStore } from '@/stores/modules/company'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
const companyStore = useCompanyStore()
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})

const isEdit = ref(false)

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  library_ids: [],
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: '',
  top_k: 0,
  similarity: 0,
  search_type: 1,
  meta_search_switch: 0,
  meta_search_type: 1,
  meta_search_condition_list: "",
  rrf_weight: {},
  recall_neighbor_switch: false,
  recall_neighbor_before_num: 1,
  recall_neighbor_after_num: 1,
})

// 知识库
const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const wxAppLibary = ref(null)
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    if (!wxAppLibary.value) {
      return formState.library_ids.includes(item.id) && item.type != 3
    } else {
      return formState.library_ids.includes(item.id)
    }
  })
})

const noOpenLibraryList = computed(() => {
  return selectedLibraryRows.value.filter((item) => item.graph_switch == 0)
})

// 移除知识库
const handleRemoveCheckedLibrary = (item) => {
  let index = formState.library_ids.indexOf(item.id)

  formState.library_ids.splice(index, 1)

  onSave()
}

const onChangeLibrarySelected = (checkedList) => {
  formState.library_ids = [...checkedList]

  onSave()
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value.open([...formState.library_ids])
}

// 召回设置
const recallSettingsAlertRef = ref(null)

const handleOpenRecallSettingsAlert = () => {
  recallSettingsAlertRef.value.open(toRaw(formState), robotInfo)
}

const noOpenGraphModalRef = ref(null)

const onChangeRecallSettings = (data) => {
  formState.rerank_status = data.rerank_status
  formState.rerank_use_model = data.rerank_use_model
  formState.rerank_model_config_id = data.rerank_model_config_id
  formState.top_k = data.top_k
  formState.similarity = data.similarity
  formState.search_type = data.search_type
  formState.meta_search_switch = data.meta_search_switch
  formState.meta_search_type = data.meta_search_type
  formState.meta_search_condition_list = data.meta_search_condition_list
  formState.rrf_weight = data.rrf_weight
  formState.recall_neighbor_switch = data.recall_neighbor_switch
  formState.recall_neighbor_before_num = data.recall_neighbor_before_num
  formState.recall_neighbor_after_num = data.recall_neighbor_after_num

  if (data.search_type == 1 || data.search_type == 4) {
    if (noOpenLibraryList.value.length > 0) {
      noOpenGraphModalRef.value.show()
    }
  }
  onSave()
}

const onSave = () => {
  let formData = { ...toRaw(formState) }

  formData.library_ids = formData.library_ids.join(',')
  formData.op_type_relation_library = 1
  formData.rrf_weight = JSON.stringify(formData.rrf_weight)

  updateRobotInfo({ ...formData })
}

// 获取知识库
const getList = async () => {
  const res = await getLibraryList({ type: '', show_open_docs: 1 })
  if (res) {
    libraryList.value = res.data || []
  }
}

const loadWxLbStatus = () => {
  // 公众号知识库是否开启
  getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' }).then((res) => {
    let _data = res?.data || {}
    if (_data?.user_config?.switch_status == 1) {
      wxAppLibary.value = _data
    }
  })
}

// 显示未关联知识库提示
const handleShowNoLibraryTip = () => {
  Modal.confirm({
    title: '还未关联知识库',
    content: '关联知识库后，机器人会根据知识库作答。',
    okText: '立即关联',
    cancelText: '暂不关联',
    onOk() {
      handleOpenSelectLibraryAlert()
    }
  })
}

formState.rerank_use_model = robotInfo.rerank_use_model || undefined

watchEffect(() => {
  formState.library_ids = robotInfo.library_ids.split(',')
  formState.rerank_status = robotInfo.rerank_status || 0
  formState.rerank_use_model = robotInfo.rerank_use_model || undefined
  formState.rerank_model_config_id = robotInfo.rerank_model_config_id || ''
  formState.top_k = robotInfo.top_k
  formState.similarity = robotInfo.similarity
  formState.search_type = robotInfo.search_type
  formState.meta_search_switch = robotInfo.meta_search_switch
  formState.meta_search_type = robotInfo.meta_search_type
  formState.meta_search_condition_list = robotInfo.meta_search_condition_list
  formState.rrf_weight = robotInfo.rrf_weight != '' ? JSON.parse(robotInfo.rrf_weight) : {vector: 0, search: 0, graph: 0}
  formState.recall_neighbor_switch = robotInfo.recall_neighbor_switch
  formState.recall_neighbor_before_num = robotInfo.recall_neighbor_before_num
  formState.recall_neighbor_after_num = robotInfo.recall_neighbor_after_num
})

const toLibraryDetail = (item) => {
  window.open(`#/library/details/knowledge-document?id=${item.id}`)
}

const getModelName = computed(() => {
  return getModelNameText(formState.rerank_model_config_id, formState.rerank_use_model, 'RERANK')
})

onMounted(() => {
  loadWxLbStatus()
  getList()
})
</script>
