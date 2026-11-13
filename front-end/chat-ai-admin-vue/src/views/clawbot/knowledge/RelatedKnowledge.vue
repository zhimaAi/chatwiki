<template>
  <div>
    <div class="toolbar-row">
      <div class="section-desc">{{ t('section_desc') }}</div>

      <div class="toolbar">
        <a-button class="recall-btn" @click="handleOpenRecallSettingsAlert">{{ t('btn_recall_settings') }}</a-button>
        <a-button type="primary" class="upload-btn" @click="handleOpenSelectLibraryAlert">{{ t('btn_associate_knowledge_base') }}</a-button>
      </div>
    </div>

    <div class="list-box" v-if="selectedLibraryRows.length > 0">
      <div class="list-item-wrapper" v-for="item in selectedLibraryRows" :key="item.id">
        <div class="list-item" @click="toLibraryDetail(item)">
          <img
            class="default-icon"
            v-if="isDefaultLibrary(item)"
            src="@/assets/img/robot/default-allow.svg"
            alt=""
          />
          <div class="library-info">
            <img class="library-icon" :src="item.avatar" alt="" />
            <div class="library-info-content">
              <div class="library-title">{{ item.library_name }}</div>
              <div class="library-type">
                <span class="type-tag" v-if="item.type == 0">{{ t('label_normal_library') }}</span>
                <span class="type-tag" v-if="item.type == 1">{{ t('label_external_library') }}</span>
                <span class="type-tag" v-if="item.type == 2">{{ t('label_qa_library') }}</span>
                <span class="type-tag" v-if="item.type == 3">{{ t('label_official_account_library') }}</span>
              </div>
            </div>
          </div>

          <div class="item-body">
            <div class="library-desc">{{ item.library_intro }}</div>
          </div>

          <div class="item-footer">
            <div class="library-size">
              <span>{{ t('label_docs') }}：{{ item.file_total }}</span>
              <span>{{ t('label_size') }}：{{ item.file_size_str }}</span>
              <span>{{ t('label_related_apps') }}：{{ item.robot_nums || 0 }}</span>
              <LibraryExpireStatus :library="item" />
            </div>
            <div class="action-box" v-if="!isDefaultLibrary(item)" @click.stop>
              <a-dropdown>
                <div class="action-item" @click.stop>
                  <svg-icon class="action-icon" name="point-h"></svg-icon>
                </div>
                <template #overlay>
                  <a-menu>
                    <a-menu-item v-if="item.type != 1">
                      <a @click.stop="handleSetDefaultLibrary(item)">{{ t('btn_set_default') }}</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a class="delete-text-color" @click.stop="handleRemoveCheckedLibrary(item)">
                        {{ t('btn_unlink') }}
                      </a>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- <div class="setting-info-block" v-if="selectedLibraryRows.length > 0">
      <div class="set-item">
        检索模式：
        <span v-if="formState.search_type == 1">混合检索</span>
        <span v-if="formState.search_type == 2">向量检索</span>
        <span v-if="formState.search_type == 3">全文检索</span>
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
      :defaultLibraryId="formState.default_library_id"
    />
    <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
  </div>
</template>

<script setup>
import { ref, reactive, watchEffect, computed, toRaw, onMounted, createVNode } from 'vue'
import { storeToRefs } from 'pinia'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { getLibraryList } from '@/api/library/index'
import { relationLibrary } from '@/api/robot/index'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_QA_AVATAR } from '@/constants/index'
import { formatFileSize } from '@/utils/index'
import LibraryExpireStatus from '@/components/library-expire-status/index.vue'
// import { getModelNameText } from '@/components/model-select/index.js'
import LibrarySelectAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/library-select-alert.vue'
import RecallSettingsAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'

const { t } = useI18n('views.clawbot.knowledge.RelatedKnowledge')
const clawbotStore = useClawbotStore()
const { robotInfo } = storeToRefs(clawbotStore)

const formState = reactive({
  library_ids: [],
  default_library_id: '',
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
  const selectedIds = new Set(formState.library_ids.map((id) => String(id)))
  const defaultLibraryId = String(formState.default_library_id || '')
  if (defaultLibraryId) {
    selectedIds.add(defaultLibraryId)
  }

  const rows = libraryList.value.filter((item) => {
    return selectedIds.has(String(item.id)) && (wxAppLibary.value || item.type != 3)
  })

  if (!defaultLibraryId) return rows

  return rows.sort((a, b) => {
    if (String(a.id) === defaultLibraryId) return -1
    if (String(b.id) === defaultLibraryId) return 1
    return 0
  })
})

const isDefaultLibrary = (item) => {
  return String(item.id) === String(formState.default_library_id || '')
}

const handleRemoveCheckedLibrary = (item) => {
  if (isDefaultLibrary(item)) return

  formState.library_ids = formState.library_ids.filter((id) => String(id) !== String(item.id))
  onSave()
}

const onChangeLibrarySelected = (checkedList) => {
  const selectedIds = checkedList.map((id) => String(id))
  const defaultLibraryId = String(formState.default_library_id || '')
  if (defaultLibraryId && !selectedIds.includes(defaultLibraryId)) {
    selectedIds.unshift(defaultLibraryId)
  }
  formState.library_ids = selectedIds
  onSave()
}

const handleSetDefaultLibrary = (item) => {
  Modal.confirm({
    title: t('msg_confirm_set_default', { library_name: item.library_name }),
    icon: null,
    content: createVNode('div', { style: 'color: red;' }, t('msg_one_default_library')),
    async onOk() {
      const assistantId = robotInfo.value?.id
      const libraryId = String(item.id)
      const libraryIds = formState.library_ids.map((id) => String(id))
      if (!libraryIds.includes(libraryId)) {
        libraryIds.unshift(libraryId)
      }

      await relationLibrary({
        library_ids: libraryIds.join(','),
        default_library_id: libraryId,
        id: assistantId
      })
      await clawbotStore.fetchRobotInfo(String(assistantId))
      message.success(t('msg_saved'))
    }
  })
}

const handleOpenSelectLibraryAlert = () => {
  const selectedIds = formState.library_ids.map((id) => String(id))
  const defaultLibraryId = String(formState.default_library_id || '')
  if (defaultLibraryId && !selectedIds.includes(defaultLibraryId)) {
    selectedIds.unshift(defaultLibraryId)
  }
  librarySelectAlertRef.value.open(selectedIds)
}

// 召回设置
const recallSettingsAlertRef = ref(null)

const handleOpenRecallSettingsAlert = () => {
  recallSettingsAlertRef.value.open(toRaw(formState), robotInfo.value)
}

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
    const list = res.data || []
    list.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)
      if (!item.avatar) {
        item.avatar = item.type == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_QA_AVATAR
      }
    })
    libraryList.value = list
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
  formState.default_library_id = info.default_library_id || ''
  formState.rerank_status = info.rerank_status || 0
  formState.rerank_use_model = info.rerank_use_model || undefined
  formState.rerank_model_config_id = info.rerank_model_config_id || ''
  formState.top_k = info.top_k
  formState.similarity = info.similarity
  formState.search_type = info.search_type
  formState.meta_search_switch = info.meta_search_switch
  formState.meta_search_type = info.meta_search_type
  formState.meta_search_condition_list = info.meta_search_condition_list
  formState.rrf_weight = info.rrf_weight != '' && info.rrf_weight ? JSON.parse(info.rrf_weight) : { vector: 0, search: 0 }
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

.list-box {
  display: flex;
  flex-flow: row wrap;
  margin: 0 -8px;
}

.list-item-wrapper {
  width: 25%;
  padding: 8px;
}

.default-icon {
  position: absolute;
  top: 0;
  left: 0;
  width: 38px;
}

.list-item {
  position: relative;
  width: 100%;
  padding: 24px;
  border: 1px solid #e4e6eb;
  border-radius: 12px;
  background-color: #fff;
  cursor: pointer;
  transition: all 0.25s;

  &:hover {
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
  }
}

.library-info {
  position: relative;
  display: flex;
  align-items: center;
}

.library-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  overflow: hidden;
}

.library-info-content {
  flex: 1;
  padding-left: 12px;
  overflow: hidden;
}

.library-title {
  height: 24px;
  margin-bottom: 4px;
  color: #262626;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.library-type {
  display: flex;

  .type-tag {
    height: 22px;
    padding: 0 8px;
    border: 1px solid #cde0ff;
    border-radius: 6px;
    color: #2475fc;
    font-size: 12px;
    line-height: 20px;
  }

}

.item-body {
  margin-top: 12px;
}

.library-desc {
  display: -webkit-box;
  height: 44px;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
}

.item-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 24px;
  margin-top: 14px;
  color: #7a8699;
}

.library-size {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  color: #7a8699;
  font-size: 12px;
  line-height: 20px;
}

.action-box {
  display: flex;
  align-items: center;
  height: 24px;

  .action-item {
    display: flex;
    align-items: center;
    height: 100%;
    padding: 4px;
    border-radius: 6px;
    color: #595959;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #e4e6eb;
    }
  }

  .action-icon {
    font-size: 16px;
  }
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

@media screen and (max-width: 1680px) {
  .list-item-wrapper {
    width: 33.3%;
  }
}

@media screen and (max-width: 1024px) {
  .list-item-wrapper {
    width: 50%;
  }
}

@media screen and (max-width: 768px) {
  .list-item-wrapper {
    width: 100%;
  }
}
</style>
