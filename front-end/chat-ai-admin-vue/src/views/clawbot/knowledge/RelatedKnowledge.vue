<template>
  <div>
    <div class="toolbar-row">
      <div class="section-desc">{{ t('section_desc') }}</div>

      <div class="toolbar">
        <a-button class="recall-btn" @click="handleOpenRecallSettingsAlert">{{ t('btn_recall_settings') }}</a-button>
        <a-button type="primary" class="upload-btn" @click="handleOpenSelectLibraryAlert">{{ t('btn_associate_knowledge_base') }}</a-button>
      </div>
    </div>

    <div class="file-list" v-if="selectedLibraryRows.length > 0">
      <div v-for="item in selectedLibraryRows" :key="item.id" class="file-item" @click="toLibraryDetail(item)">
        <div class="file-left">
          <div class="library-avatar">
            <img :src="item.avatar" alt="" />
          </div>
          <div>
            <div class="file-name">{{ item.library_name }}</div>
            <div class="file-size">{{ item.library_intro }}</div>
          </div>
        </div>
        <a-button danger ghost class="remove-btn" @click.stop="handleRemoveCheckedLibrary(item)">{{ t('btn_unlink') }}</a-button>
      </div>
    </div>

    <!-- <div class="setting-info-block" v-if="selectedLibraryRows.length > 0">
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
        <span v-if="formState.rerank_status == 0">关闭</span>
        <span v-else>{{ getModelName }}</span>
      </div>
    </div> -->

    <div v-if="!selectedLibraryRows.length" class="empty-tip">
      {{ t('empty_tip') }}
    </div>

    <LibrarySelectAlert
      ref="librarySelectAlertRef"
      @change="onChangeLibrarySelected"
      :showWxType="!!wxAppLibary"
    />
    <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
    <NoOpenGraphModal :list="noOpenLibraryList" @refreshList="getList" ref="noOpenGraphModalRef" />
  </div>
</template>

<script setup>
import { ref, reactive, watchEffect, computed, toRaw, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { getLibraryList } from '@/api/library/index'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
// import { getModelNameText } from '@/components/model-select/index.js'
import LibrarySelectAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/library-select-alert.vue'
import RecallSettingsAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import NoOpenGraphModal from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/no-open-graph-modal.vue'

const { t } = useI18n('views.clawbot.knowledge.RelatedKnowledge')
const clawbotStore = useClawbotStore()
const { robotInfo } = storeToRefs(clawbotStore)

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
  library_search_type: 'fullTextSearch'
})

// 知识库列表
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
  recallSettingsAlertRef.value.open(toRaw(formState), robotInfo.value)
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
  formState.recall_neighbor_top_k = data.recall_neighbor_top_k
  formState.recall_neighbor_before_num = data.recall_neighbor_before_num
  formState.recall_neighbor_after_num = data.recall_neighbor_after_num
  formState.library_search_type = data.library_search_type
  if (data.search_type == 1 || data.search_type == 4) {
    if (noOpenLibraryList.value.length > 0) {
      noOpenGraphModalRef.value.show()
    }
  }
  onSave()
}

// 保存：以完整 robotInfo 为基础，覆盖知识库相关字段
const onSave = async () => {
  let localState = toRaw(formState)
  let partialData = {}

  // 覆盖知识库和召回设置相关字段
  partialData.library_ids = localState.library_ids.join(',')
  partialData.rerank_status = localState.rerank_status
  partialData.rerank_use_model = localState.rerank_use_model || ''
  partialData.rerank_model_config_id = localState.rerank_model_config_id || 0
  partialData.top_k = localState.top_k
  partialData.similarity = localState.similarity
  partialData.search_type = localState.search_type
  partialData.meta_search_switch = localState.meta_search_switch
  partialData.meta_search_type = localState.meta_search_type
  partialData.meta_search_condition_list = localState.meta_search_condition_list
  partialData.rrf_weight = JSON.stringify(localState.rrf_weight || {})
  partialData.recall_neighbor_switch = localState.recall_neighbor_switch
  partialData.recall_neighbor_top_k = localState.recall_neighbor_top_k
  partialData.recall_neighbor_before_num = localState.recall_neighbor_before_num
  partialData.recall_neighbor_after_num = localState.recall_neighbor_after_num
  partialData.library_search_type = localState.library_search_type
  partialData.op_type_relation_library = 1

  try {
    await clawbotStore.saveAssistant(partialData, {
      optimistic: false,
      refreshAfterSave: true,
      successMessage: t('msg_saved')
    })
  } catch (err) {
    console.error('保存失败', err)
  }
}

// 获取知识库列表
const getList = async () => {
  const res = await getLibraryList({ type: '', show_open_docs: 1 })
  if (res) {
    libraryList.value = res.data || []
  }
}

// 公众号知识库状态
const loadWxLbStatus = () => {
  getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' }).then((res) => {
    let _data = res?.data || {}
    if (_data?.user_config?.switch_status == 1) {
      wxAppLibary.value = _data
    }
  })
}

// 跳转知识库详情
const toLibraryDetail = (item) => {
  window.open(`#/library/details/knowledge-document?id=${item.id}`)
}

// const getModelName = computed(() => {
//   return getModelNameText(formState.rerank_model_config_id, formState.rerank_use_model, 'RERANK')
// })

// 同步 robotInfo 到 formState
watchEffect(() => {
  if (!robotInfo.value) return
  const info = robotInfo.value
  formState.library_ids = (info.library_ids || '').split(',').filter(Boolean)
  formState.rerank_status = info.rerank_status || 0
  formState.rerank_use_model = info.rerank_use_model || undefined
  formState.rerank_model_config_id = info.rerank_model_config_id || ''
  formState.top_k = info.top_k
  formState.similarity = info.similarity
  formState.search_type = info.search_type
  formState.meta_search_switch = info.meta_search_switch
  formState.meta_search_type = info.meta_search_type
  formState.meta_search_condition_list = info.meta_search_condition_list
  formState.rrf_weight = info.rrf_weight != '' && info.rrf_weight ? JSON.parse(info.rrf_weight) : { vector: 0, search: 0, graph: 0 }
  formState.recall_neighbor_switch = info.recall_neighbor_switch
  formState.recall_neighbor_top_k = info.recall_neighbor_top_k
  formState.recall_neighbor_before_num = info.recall_neighbor_before_num
  formState.recall_neighbor_after_num = info.recall_neighbor_after_num
  formState.library_search_type = info.library_search_type
})

onMounted(() => {
  loadWxLbStatus()
  getList()
})
</script>

<style lang="less" scoped>
.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 16px 0 12px;
}

.section-desc {
  font-size: 13px;
  color: #999;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
}

.recall-btn {
  height: 28px;
  padding: 0 14px;
  border-radius: 6px;
  font-size: 13px;
}

.upload-btn {
  height: 28px;
  padding: 0 14px;
  border-radius: 6px;
  font-size: 13px;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.file-item {
  height: 60px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  padding: 0 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  cursor: pointer;

  &:hover {
    border-color: #d9d9d9;
  }
}

.file-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.library-avatar {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.file-name {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  margin-top: 2px;
  color: #999;
  font-size: 12px;
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.remove-btn {
  border-radius: 6px;
  height: 28px;
  font-size: 13px;
}

.setting-info-block {
  margin-top: 12px;
  padding: 12px 14px;
  background: #f9f9f9;
  border-radius: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  color: #595959;
  font-size: 13px;
  line-height: 22px;

  .set-item {
    display: flex;
    align-items: center;
  }
}

.empty-tip {
  padding: 40px 0;
  text-align: center;
  color: #999;
  font-size: 14px;
}
</style>
