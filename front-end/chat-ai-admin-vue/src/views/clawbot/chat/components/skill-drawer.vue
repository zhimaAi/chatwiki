<template>
  <a-drawer
    :open="props.open"
    placement="right"
    :width="500"
    :closable="false"
    :bodyStyle="{ padding: '24px', background: '#fff' }"
    :headerStyle="{ padding: '10px 16px', borderBottom: '1px solid #f0f0f0' }"
    @close="handleClose"
  >
    <template #title>
      <div class="drawer-header">
        <span class="drawer-title">{{ t('title_skill') }}</span>
        <button class="drawer-close" type="button" @click="handleClose">
          <CloseOutlined />
        </button>
      </div>
    </template>

    <div class="skill-drawer">
      <div class="add-action-row">
        <a-button type="dashed" class="add-skill-btn" @click="handleOpenToolModal">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('menu_add_tool') }}
        </a-button>
        <a-button type="dashed" class="add-skill-btn" @click="handleOpenSelectSkillModal">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('btn_add_skill') }}
        </a-button>
      </div>

      <div class="skill-card">
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ t('title_query_knowledge_base') }}</span>
              <span class="skill-tag skill">{{ t('tag_skill') }}</span>
            </div>
            <div v-overflow-tooltip="t('desc_query_knowledge_base')" class="skill-desc">
              {{ t('desc_query_knowledge_base') }}
            </div>
          </div>
          <div class="skill-status-text-group" :class="{ inactive: !knowledgeEnabled }">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ knowledgeEnabled ? t('status_enabled') : t('status_disabled') }}</span>
          </div>
        </div>

        <div class="skill-content-block">
          <div v-if="knowledgeLoading" class="skill-loading">
            <a-spin />
          </div>
          <div v-else-if="selectedLibraryRows.length" class="content-list">
            <div
              v-for="item in selectedLibraryRows"
              :key="item.id"
              class="content-item"
              @click="toLibraryDetail(item)"
            >
              <div class="content-item-main">
                <div class="library-avatar">
                  <img :src="item.avatar" alt="" />
                </div>
                <div class="content-item-text">
                  <div class="content-item-title">{{ item.library_name }}</div>
                  <div class="content-item-desc">{{ item.library_intro || t('msg_no_intro') }}</div>
                </div>
              </div>
              

              <a-tooltip :title="t('tooltip_remove_knowledge_base')">
                <button class="remove-icon-btn remove-icon-btn-plain" type="button" @click.stop="handleRemoveCheckedLibrary(item)">
                  <MinusCircleFilled />
                </button>
              </a-tooltip>
            </div>
          </div>
          <div v-else class="empty-tip">
            {{ t('empty_no_knowledge_base') }}
          </div>
        </div>

        <div class="skill-footer">
          <div class="skill-footer-actions">
            <button class="skill-action-btn skill-action-btn-primary" type="button" @click="handleOpenSelectLibraryAlert">
              {{ t('btn_associate_knowledge_base') }}
            </button>
            <button class="skill-action-btn" type="button" @click="handleOpenRecallSettingsAlert">
              {{ t('btn_recall_settings') }}
            </button>
          </div>
          <a-switch
            class="skill-switch"
            :checked="knowledgeEnabled"
            :loading="knowledgeSwitchLoading"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleKnowledge"
          />
        </div>
      </div>

      <div class="skill-card">
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ t('title_query_local_docs') }}</span>
              <span class="skill-tag skill">{{ t('tag_skill') }}</span>
            </div>
            <div v-overflow-tooltip="t('desc_query_local_docs')" class="skill-desc">
              {{ t('desc_query_local_docs') }}
            </div>
          </div>
          <div class="skill-status-text-group" :class="{ inactive: !localDocsEnabled }">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ localDocsEnabled ? t('status_enabled') : t('status_disabled') }}</span>
          </div>
        </div>

        <div class="skill-content-block">
          <div v-if="localDocsLoading" class="skill-loading">
            <a-spin />
          </div>
          <div v-else-if="localDocList.length" class="content-list">
            <div v-for="item in localDocList" :key="item.name" class="content-item local-doc-item">
              <div class="content-item-main">
                <div class="file-ext-tag" :class="item.ext">{{ item.ext.toUpperCase() }}</div>
                <div class="content-item-text">
                  <div class="content-item-title">{{ item.name }}</div>
                  <div class="content-item-desc">
                    {{ formatSize(item.size) }}
                    <span v-if="formatDate(item.time)"> · {{ formatDate(item.time) }}</span>
                  </div>
                </div>
              </div>
              <a-tooltip :title="t('tooltip_remove_file')">
                <button
                  class="remove-icon-btn remove-icon-btn-box"
                  type="button"
                  :disabled="!!deletingMap[item.name]"
                  @click="handleRemoveLocalDoc(item)"
                >
                  <a-spin v-if="deletingMap[item.name]" size="small" />
                  <MinusCircleFilled v-else />
                </button>
              </a-tooltip>
            </div>
          </div>
          <div v-else class="empty-tip">
            {{ t('empty_no_local_docs') }}
          </div>
        </div>

        <div class="skill-footer">
          <div class="skill-footer-actions">
            <button class="skill-action-btn skill-action-btn-primary" type="button" @click="handleUploadClick" :disabled="uploading">
              <a-spin v-if="uploading" size="small" />
              <template v-else>{{ t('btn_upload_file') }}</template>
            </button>
          </div>
          <a-switch
            class="skill-switch"
            :checked="localDocsEnabled"
            :loading="localDocsSwitchLoading"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleLocalDocs"
          />
        </div>
      </div>

      <div class="skill-card">
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ t('title_agent_write_file') }}</span>
              <span class="skill-tag tool">{{ t('tag_tool') }}</span>
            </div>
            <div class="skill-desc">{{ t('desc_agent_write_file') }}</div>
          </div>
          <div class="skill-status-text-group" :class="{ inactive: !writeFileEnabled }">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ writeFileEnabled ? t('status_enabled') : t('status_disabled') }}</span>
          </div>
        </div>
        <div class="skill-footer">
          <div class="skill-footer-actions" />
          <a-switch
            class="skill-switch"
            :checked="writeFileEnabled"
            :loading="writeFileSwitchLoading"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleWriteFile"
          />
        </div>
      </div>

      <div class="skill-card">
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ t('title_agent_edit_file') }}</span>
              <span class="skill-tag tool">{{ t('tag_tool') }}</span>
            </div>
            <div class="skill-desc">{{ t('desc_agent_edit_file') }}</div>
          </div>
          <div class="skill-status-text-group" :class="{ inactive: !editFileEnabled }">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ editFileEnabled ? t('status_enabled') : t('status_disabled') }}</span>
          </div>
        </div>
        <div class="skill-footer">
          <div class="skill-footer-actions" />
          <a-switch
            class="skill-switch"
            :checked="editFileEnabled"
            :loading="editFileSwitchLoading"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleEditFile"
          />
        </div>
      </div>

      <div class="skill-card">
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ t('title_agent_execute') }}</span>
              <span class="skill-tag tool">{{ t('tag_tool') }}</span>
            </div>
            <div class="skill-desc">{{ t('desc_agent_execute') }}</div>
          </div>
          <div class="skill-status-text-group" :class="{ inactive: !executeEnabled }">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ executeEnabled ? t('status_enabled') : t('status_disabled') }}</span>
          </div>
        </div>
        <div class="skill-footer">
          <div class="skill-footer-actions" />
          <a-switch
            class="skill-switch"
            :checked="executeEnabled"
            :loading="executeSwitchLoading"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleExecute"
          />
        </div>
      </div>

      <!-- 从商品库推荐商品 -->
      <div class="skill-card">
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ t('title_goods_recommend') }}</span>
              <span class="skill-tag tool">{{ t('tag_tool') }}</span>
            </div>
            <div class="skill-desc">{{ t('desc_goods_recommend') }}</div>
          </div>
          <div class="skill-status-text-group" :class="{ inactive: !goodsRecommendEnabled }">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ goodsRecommendEnabled ? t('status_enabled') : t('status_disabled') }}</span>
          </div>
        </div>

        <div class="skill-content-block">
          <div class="scope-info-text">{{ scopeInfoText }}</div>
        </div>

        <div class="skill-footer">
          <div class="skill-footer-actions">
            <button class="skill-action-btn skill-action-btn-primary" type="button" @click="handleOpenScopeModal">
              {{ t('btn_recommend_scope') }}
            </button>
          </div>
          <a-switch
            class="skill-switch"
            :checked="goodsRecommendEnabled"
            :loading="goodsSwitchLoading"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleGoodsRecommend"
          />
        </div>
      </div>
      
      <div
        v-for="item in skills"
        :key="item.id"
        class="skill-card"
      >
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ item.name }}</span>
              <span class="skill-tag skill">{{ t('tag_skill') }}</span>
            </div>
            <div v-overflow-tooltip="item.desc || ''" class="skill-desc">
              {{ item.desc || '' }}
            </div>
          </div>
          <div class="skill-status-text-group">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ t('status_enabled') }}</span>
          </div>
        </div>

        <div class="skill-footer single-action">
          <div class="skill-footer-actions">
            <button class="skill-action-btn skill-action-btn-danger" type="button" @click="handleRemoveSkill(item)">
              {{ t('btn_remove') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="skillListLoading" class="tool-loading-card">
        <a-spin />
      </div>

      <div
        v-for="item in workFlowSkills"
        :key="item.id"
        class="skill-card"
      >
        <div class="skill-card-header">
          <div class="skill-card-header-left">
            <div class="skill-title-row">
              <span class="skill-title">{{ item.name }}</span>
              <span class="skill-tag tool">{{ t('tag_tool') }}</span>
            </div>
            <div v-overflow-tooltip="item.desc || '—'" class="skill-desc">
              {{ item.desc || '—' }}
            </div>
          </div>
          <div class="skill-status-text-group">
            <span class="skill-status-dot"></span>
            <span class="skill-status-text">{{ t('status_enabled') }}</span>
          </div>
        </div>

        <div class="skill-footer single-action">
          <div class="skill-footer-actions">
            <button class="skill-action-btn skill-action-btn-danger" type="button" @click="handleRemoveWorkFlow(item.id)">
              {{ t('btn_remove') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="workFlowLoading" class="tool-loading-card">
        <a-spin />
      </div>
    </div>

    <LibrarySelectAlert
      ref="librarySelectAlertRef"
      :showWxType="!!wxAppLibary"
      @change="onChangeLibrarySelected"
    />
    <RecallSettingsAlert ref="recallSettingsAlertRef" @change="onChangeRecallSettings" />
    <NoOpenGraphModal
      ref="noOpenGraphModalRef"
      :list="noOpenLibraryList"
      @refreshList="getLibraryOptions"
    />
    <AddToolModal
      v-model:visible="toolModalVisible"
      :robotId="currentAssistant?.id"
      :workFlowIds="robotInfo?.work_flow_ids || ''"
      @confirm="handleToolConfirm"
    />
    <SelectSkillModal
      v-model:visible="selectSkillModalVisible"
      :robotId="currentAssistant?.id"
      :refreshKey="selectSkillRefreshKey"
      @create="handleOpenUploadSkillModal"
      @confirm="handleSelectSkillConfirm"
    />
    <UploadSkillZipModal
      v-model:visible="uploadSkillZipVisible"
      :robotId="currentAssistant?.id"
      :skill-id="editingSkillId"
      @confirm="handleUploadSkillConfirm"
    />
    <LocalDocUploadModal
      v-model:open="uploadModalOpen"
      :loading="uploading"
      @confirm="handleUploadConfirm"
    />

    <!-- 推荐范围弹窗 -->
    <GoodsRecommendScopeModal v-model:visible="scopeModalVisible" />
  </a-drawer>
</template>

<script setup>
import { computed, createVNode, nextTick, reactive, ref, toRaw, watch, watchEffect } from 'vue'
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import {
  CloseOutlined,
  ExclamationCircleOutlined,
  MinusCircleFilled,
  PlusOutlined
} from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { getRobotList, relationWorkFlow } from '@/api/robot'
import { getLibraryList } from '@/api/library'
import { getSpecifyAbilityConfig } from '@/api/explore/index.js'
import {
  deleteClawbotSkill,
  uploadClawbotLocalDoc,
  getClawbotLocalDocList,
  deleteClawbotLocalDoc,
  getClawbotSkillList
} from '@/api/clawbot'
import AddToolModal from '@/views/clawbot/skills/components/AddToolModal.vue'
import SelectSkillModal from '@/views/clawbot/skills/components/SelectSkillModal.vue'
import UploadSkillZipModal from '@/views/clawbot/skills/components/UploadSkillZipModal.vue'
import GoodsRecommendScopeModal from '@/views/clawbot/skills/components/GoodsRecommendScopeModal.vue'
import LocalDocUploadModal from '@/views/clawbot/components/LocalDocUploadModal.vue'
import LibrarySelectAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/library-select-alert.vue'
import RecallSettingsAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import NoOpenGraphModal from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/no-open-graph-modal.vue'

const { t } = useI18n('views.clawbot.skills.index')
const props = defineProps({
  open: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const clawbotStore = useClawbotStore()
const { robotInfo, currentAssistant } = storeToRefs(clawbotStore)
const { updateClawbotConf, fetchRobotInfo } = clawbotStore

const toolModalVisible = ref(false)
const selectSkillModalVisible = ref(false)
const selectSkillRefreshKey = ref(0)
const uploadSkillZipVisible = ref(false)
const editingSkillId = ref(0)
const workFlowLoading = ref(false)
const workFlowSkills = ref([])
const skillListLoading = ref(false)
const uploadedSkills = ref([])
const knowledgeSwitchLoading = ref(false)
const localDocsSwitchLoading = ref(false)
const writeFileSwitchLoading = ref(false)
const editFileSwitchLoading = ref(false)
const executeSwitchLoading = ref(false)
const knowledgeLoading = ref(false)
const localDocsLoading = ref(false)
const uploading = ref(false)
const deletingMap = ref({})
const localDocList = ref([])
const uploadModalOpen = ref(false)

// 商品库推荐
const goodsSwitchLoading = ref(false)
const scopeModalVisible = ref(false)

const libraryList = ref([])
const librarySelectAlertRef = ref(null)
const recallSettingsAlertRef = ref(null)
const noOpenGraphModalRef = ref(null)
const wxAppLibary = ref(null)

const updateOverflowTooltip = (el, value) => {
  nextTick(() => {
    const text = value || ''
    const isOverflowing = el.scrollHeight > el.clientHeight + 1 || el.scrollWidth > el.clientWidth + 1

    if (text && isOverflowing) {
      el.setAttribute('title', text)
      return
    }

    el.removeAttribute('title')
  })
}

const vOverflowTooltip = {
  mounted(el, binding) {
    updateOverflowTooltip(el, binding.value)

    if (window.ResizeObserver) {
      el._overflowTooltipObserver = new ResizeObserver(() => {
        updateOverflowTooltip(el, binding.value)
      })
      el._overflowTooltipObserver.observe(el)
    }
  },
  updated(el, binding) {
    updateOverflowTooltip(el, binding.value)
  },
  unmounted(el) {
    el._overflowTooltipObserver?.disconnect()
  }
}

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
  meta_search_condition_list: '',
  rrf_weight: {},
  recall_neighbor_switch: false,
  recall_neighbor_top_k: 0,
  recall_neighbor_before_num: 1,
  recall_neighbor_after_num: 1,
  library_search_type: 'fullTextSearch'
})

const parseRrfWeight = (value) => {
  if (!value) {
    return { vector: 0, search: 0, graph: 0 }
  }

  try {
    return typeof value === 'string' ? JSON.parse(value) : value
  } catch {
    return { vector: 0, search: 0, graph: 0 }
  }
}

const knowledgeEnabled = computed(() => !Number(robotInfo.value?.search_knowledge_close || 0))
const localDocsEnabled = computed(() => !Number(robotInfo.value?.query_local_docs_close || 0))
const writeFileEnabled = computed(() => Number(robotInfo.value?.open_agent_write_file_tool || 0) === 1)
const editFileEnabled = computed(() => Number(robotInfo.value?.open_agent_edit_file_tool || 0) === 1)
const executeEnabled = computed(() => Number(robotInfo.value?.open_agent_execute_tool || 0) === 1)
const goodsRecommendEnabled = computed(() => Number(robotInfo.value?.goods_lib_recommend_switch || 0) === 1)

const scopeInfoText = computed(() => {
  const ids = (robotInfo.value?.goods_lib_recommend_group_ids || '').split(',').filter(Boolean)
  if (ids.length === 0) {
    return t('empty_scope_all')
  }
  return t('empty_scope_partial', { count: ids.length })
})

const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    if (!wxAppLibary.value) {
      return formState.library_ids.includes(item.id) && item.type != 3
    }
    return formState.library_ids.includes(item.id)
  })
})

const skills = computed(() => uploadedSkills.value)

const noOpenLibraryList = computed(() => {
  return selectedLibraryRows.value.filter((item) => item.graph_switch == 0)
})

watchEffect(() => {
  if (!robotInfo.value) {
    return
  }

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
  formState.rrf_weight = parseRrfWeight(info.rrf_weight)
  formState.recall_neighbor_switch = info.recall_neighbor_switch
  formState.recall_neighbor_top_k = info.recall_neighbor_top_k
  formState.recall_neighbor_before_num = info.recall_neighbor_before_num
  formState.recall_neighbor_after_num = info.recall_neighbor_after_num
  formState.library_search_type = info.library_search_type || 'fullTextSearch'
})

watch(
  () => props.open,
  (open) => {
    if (open) {
      hydrateDrawer()
      return
    }

    toolModalVisible.value = false
    uploadModalOpen.value = false
  }
)

watch(
  () => robotInfo.value?.work_flow_ids,
  () => {
    if (props.open) {
      loadWorkFlowSkills()
    }
  }
)

watch(
  () => currentAssistant.value?.id,
  () => {
    if (props.open) {
      loadSkillList()
    }
  }
)

const hydrateDrawer = async () => {
  await Promise.allSettled([
    loadSkillList(),
    loadWorkFlowSkills(),
    fetchLocalDocList(),
    getLibraryOptions(),
    loadWxLibraryStatus()
  ])
}

const handleClose = () => {
  emit('close')
}

function buildConfPayload(overrides = {}) {
  return {
    id: currentAssistant.value?.id,
    search_knowledge_close: knowledgeEnabled.value ? 0 : 1,
    query_local_docs_close: localDocsEnabled.value ? 0 : 1,
    open_agent_write_file_tool: writeFileEnabled.value ? 1 : 0,
    open_agent_edit_file_tool: editFileEnabled.value ? 1 : 0,
    open_agent_execute_tool: executeEnabled.value ? 1 : 0,
    goods_lib_recommend_switch: goodsRecommendEnabled.value ? 1 : 0,
    goods_lib_recommend_group_ids: robotInfo.value?.goods_lib_recommend_group_ids || '',
    ...overrides
  }
}

const toggleKnowledge = async (checked) => {
  if (!currentAssistant.value?.id || knowledgeSwitchLoading.value) {
    return
  }

  knowledgeSwitchLoading.value = true
  try {
    await updateClawbotConf(buildConfPayload({ search_knowledge_close: checked ? 0 : 1 }))
  } finally {
    knowledgeSwitchLoading.value = false
  }
}

const toggleLocalDocs = async (checked) => {
  if (!currentAssistant.value?.id || localDocsSwitchLoading.value) {
    return
  }

  localDocsSwitchLoading.value = true
  try {
    await updateClawbotConf(buildConfPayload({ query_local_docs_close: checked ? 0 : 1 }))
  } finally {
    localDocsSwitchLoading.value = false
  }
}

const toggleWriteFile = async (checked) => {
  if (!currentAssistant.value?.id || writeFileSwitchLoading.value) return
  writeFileSwitchLoading.value = true
  try {
    await updateClawbotConf(buildConfPayload({ open_agent_write_file_tool: checked ? 1 : 0 }))
  } finally {
    writeFileSwitchLoading.value = false
  }
}

const toggleEditFile = async (checked) => {
  if (!currentAssistant.value?.id || editFileSwitchLoading.value) return
  editFileSwitchLoading.value = true
  try {
    await updateClawbotConf(buildConfPayload({ open_agent_edit_file_tool: checked ? 1 : 0 }))
  } finally {
    editFileSwitchLoading.value = false
  }
}

const toggleExecute = async (checked) => {
  if (!currentAssistant.value?.id || executeSwitchLoading.value) return
  executeSwitchLoading.value = true
  try {
    await updateClawbotConf(buildConfPayload({ open_agent_execute_tool: checked ? 1 : 0 }))
  } finally {
    executeSwitchLoading.value = false
  }
}

const toggleGoodsRecommend = async (checked) => {
  if (!currentAssistant.value?.id || goodsSwitchLoading.value) {
    return
  }

  goodsSwitchLoading.value = true
  try {
    await updateClawbotConf(buildConfPayload({ goods_lib_recommend_switch: checked ? 1 : 0 }))
  } finally {
    goodsSwitchLoading.value = false
  }
}

const handleOpenScopeModal = () => {
  scopeModalVisible.value = true
}

const loadWorkFlowSkills = async () => {
  const workFlowIds = robotInfo.value?.work_flow_ids
  if (!workFlowIds) {
    workFlowSkills.value = []
    return
  }

  workFlowLoading.value = true
  try {
    const res = await getRobotList({ application_type: 1 })
    if (res?.res === 0) {
      const allTools = res.data || []
      const savedIds = workFlowIds.split(',').filter(Boolean)
      workFlowSkills.value = allTools
        .filter((item) => savedIds.includes(String(item.id)))
        .map((item) => ({
          id: String(item.id),
          name: item.robot_name,
          desc: item.robot_intro || ''
        }))
    }
  } catch (err) {
    console.error('加载已关联工具失败', err)
  } finally {
    workFlowLoading.value = false
  }
}

const handleOpenToolModal = () => {
  toolModalVisible.value = true
}

const handleOpenSelectSkillModal = () => {
  selectSkillModalVisible.value = true
}

const handleCloseSelectSkillModal = () => {
  selectSkillModalVisible.value = false
}

const handleOpenUploadSkillModal = () => {
  editingSkillId.value = 0
  uploadSkillZipVisible.value = true
}

const handleUploadSkillConfirm = async () => {
  await fetchRobotInfo()
  await loadSkillList()
  if (selectSkillModalVisible.value) {
    selectSkillRefreshKey.value += 1
  }
}

const handleSelectSkillConfirm = async () => {
  await loadSkillList()
}

const handleToolConfirm = async () => {
  await fetchRobotInfo()
  await loadWorkFlowSkills()
}

const handleRemoveWorkFlow = (id) => {
  Modal.confirm({
    title: t('title_remove_tool'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_remove_tool'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    okType: 'danger',
    onOk: async () => {
      const currentIds = robotInfo.value?.work_flow_ids?.split(',').filter(Boolean) || []
      const newIds = currentIds.filter((item) => item !== String(id))
      try {
        const res = await relationWorkFlow({
          id: currentAssistant.value?.id,
          work_flow_ids: newIds.join(',')
        })
        if (res?.res === 0) {
          message.success(t('msg_remove_success'))
          await fetchRobotInfo()
          await loadWorkFlowSkills()
        } else {
          message.error(res?.msg || t('msg_remove_failed'))
        }
      } catch (err) {
        console.error('移除工具失败', err)
        message.error(t('msg_remove_failed'))
      }
    }
  })
}

const loadSkillList = async () => {
  if (!currentAssistant.value?.id) {
    uploadedSkills.value = []
    return
  }

  skillListLoading.value = true
  try {
    const res = await getClawbotSkillList({
      id: currentAssistant.value.id,
      source: 'all'
    })
    if (res?.res === 0) {
      uploadedSkills.value = (res.data || []).filter((item) => Number(item.is_selected) === 1).map((item, index) => ({
        id: `${item.source_type || item.source || 'skill'}-${item.skill_id || 0}-${item.skill_name || index}`,
        skillId: item.skill_id,
        name: item.remark_name || item.skill_name || '',
        desc: item.intro || item.description || '',
        raw: item
      }))
    } else {
      message.error(res?.msg || t('msg_fetch_skill_failed'))
    }
  } catch (err) {
    console.error('鑾峰彇 Skill 鍒楄〃澶辫触', err)
    message.error(err?.msg || t('msg_fetch_skill_failed'))
  } finally {
    skillListLoading.value = false
  }
}

const handleRemoveSkill = (item) => {
  Modal.confirm({
    title: t('title_remove_skill'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_remove_skill'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    okType: 'danger',
    onOk: async () => {
      try {
        const res = await deleteClawbotSkill({
          id: currentAssistant.value?.id,
          skill_id: item.skillId
        })
        if (res?.res === 0) {
          message.success(t('msg_remove_success'))
          await loadSkillList()
        } else {
          message.error(res?.msg || t('msg_remove_failed'))
        }
      } catch (err) {
        console.error('鍒犻櫎 Skill 澶辫触', err)
        message.error(err?.msg || t('msg_remove_failed'))
      }
    }
  })
}

const getLibraryOptions = async () => {
  knowledgeLoading.value = true
  try {
    const res = await getLibraryList({ type: '', show_open_docs: 1 })
    libraryList.value = res?.data || []
  } finally {
    knowledgeLoading.value = false
  }
}

const loadWxLibraryStatus = async () => {
  const res = await getSpecifyAbilityConfig({ ability_type: 'library_ability_official_account' })
  const config = res?.data || {}
  wxAppLibary.value = config?.user_config?.switch_status == 1 ? config : null
}

const handleRemoveCheckedLibrary = (item) => {
  const index = formState.library_ids.indexOf(item.id)
  if (index > -1) {
    formState.library_ids.splice(index, 1)
    saveKnowledgeConfig()
  }
}

const onChangeLibrarySelected = (checkedList) => {
  formState.library_ids = [...checkedList]
  saveKnowledgeConfig()
}

const handleOpenSelectLibraryAlert = () => {
  librarySelectAlertRef.value?.open([...formState.library_ids])
}

const handleOpenRecallSettingsAlert = () => {
  recallSettingsAlertRef.value?.open(toRaw(formState), robotInfo.value)
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

  if ((data.search_type == 1 || data.search_type == 4) && noOpenLibraryList.value.length > 0) {
    noOpenGraphModalRef.value?.show()
  }

  saveKnowledgeConfig()
}

const saveKnowledgeConfig = async () => {
  const localState = toRaw(formState)
  const partialData = {
    library_ids: localState.library_ids.join(','),
    rerank_status: localState.rerank_status,
    rerank_use_model: localState.rerank_use_model || '',
    rerank_model_config_id: localState.rerank_model_config_id || 0,
    top_k: localState.top_k,
    similarity: localState.similarity,
    search_type: localState.search_type,
    meta_search_switch: localState.meta_search_switch,
    meta_search_type: localState.meta_search_type,
    meta_search_condition_list: localState.meta_search_condition_list,
    rrf_weight: JSON.stringify(localState.rrf_weight || {}),
    recall_neighbor_switch: localState.recall_neighbor_switch,
    recall_neighbor_top_k: localState.recall_neighbor_top_k,
    recall_neighbor_before_num: localState.recall_neighbor_before_num,
    recall_neighbor_after_num: localState.recall_neighbor_after_num,
    library_search_type: localState.library_search_type,
    op_type_relation_library: 1
  }

  try {
    await clawbotStore.saveAssistant(partialData, {
      optimistic: false,
      refreshAfterSave: true,
      successMessage: t('msg_saved')
    })
  } catch (err) {
    console.error('保存知识库设置失败', err)
  }
}

const toLibraryDetail = (item) => {
  window.open(`#/library/details/knowledge-document?id=${item.id}`)
}

const fetchLocalDocList = async () => {
  const id = currentAssistant.value?.id
  if (!id) {
    localDocList.value = []
    return
  }

  localDocsLoading.value = true
  try {
    const res = await getClawbotLocalDocList({ id })
    if (res?.res === 0) {
      localDocList.value = (res.data || []).map((item) => ({
        name: item.name,
        size: item.size,
        time: item.time,
        ext: String(item.ext || '').toLowerCase()
      }))
    } else {
      message.error(res?.msg || t('msg_fetch_docs_failed'))
    }
  } catch (err) {
    console.error('获取文档列表失败', err)
    message.error(t('msg_fetch_docs_failed'))
  } finally {
    localDocsLoading.value = false
  }
}

const handleUploadClick = () => {
  uploadModalOpen.value = true
}

const handleUploadConfirm = async ({ file, description, keywords }) => {
  const pureName = file.name.replace(/[\\/]/g, '')
  if (!pureName || pureName === '.' || pureName === '..') {
    message.error(t('msg_invalid_name'))
    return
  }

  const id = currentAssistant.value?.id
  if (!id) {
    message.error(t('msg_no_agent_selected'))
    return
  }

  uploading.value = true
  try {
    const res = await uploadClawbotLocalDoc({
      id,
      file,
      description,
      'keywords[]': keywords
    })
    if (res?.res === 0) {
      message.success(t('msg_upload_success'))
      uploadModalOpen.value = false
      await fetchLocalDocList()
    } else {
      message.error(res?.msg || t('msg_upload_failed'))
    }
  } catch (err) {
    console.error('上传失败', err)
    message.error(t('msg_upload_failed'))
  } finally {
    uploading.value = false
  }
}

const handleRemoveLocalDoc = (item) => {
  Modal.confirm({
    title: t('title_delete_document'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete_document', { name: item.name }),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    okType: 'danger',
    onOk: async () => {
      const id = currentAssistant.value?.id
      if (!id) {
        return
      }

      deletingMap.value[item.name] = true
      try {
        const res = await deleteClawbotLocalDoc({ id, name: item.name })
        if (res?.res === 0) {
          message.success(t('msg_delete_success'))
          await fetchLocalDocList()
        } else {
          message.error(res?.msg || t('msg_delete_failed'))
        }
      } catch (err) {
        console.error('删除失败', err)
        message.error(t('msg_delete_failed'))
      } finally {
        deletingMap.value[item.name] = false
      }
    }
  })
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const index = Math.floor(Math.log(bytes) / Math.log(1024))
  return `${parseFloat((bytes / Math.pow(1024, index)).toFixed(2))} ${units[index]}`
}

const formatDate = (time) => {
  if (!time) return ''
  const value = dayjs(time)
  return value.isValid() ? value.format('YYYY/MM/DD') : ''
}
</script>

<style lang="less" scoped>
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.drawer-title {
  font-size: 16px;
  line-height: 24px;
  font-weight: 600;
  color: #262626;
}

.drawer-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: #595959;
  cursor: pointer;

  &:hover {
    background: #f5f5f5;
  }
}

.skill-drawer {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.add-action-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.add-skill-btn {
  height: 40px;
  border-radius: 8px;
  border-color: #99bffd;
  color: #2475fc;
  font-size: 14px;
  line-height: 22px;
  font-weight: 500;
}

.skill-card {
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
  background: #fff;
  overflow: hidden;
}

.skill-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px;
}

.skill-card-header-left {
  min-width: 0;
  flex: 1;
}

.skill-title-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.skill-title {
  font-size: 16px;
  line-height: 24px;
  font-weight: 400;
  color: #262626;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.skill-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 18px;
  padding: 1px 6px;
  border-radius: 4px;
  color: #2475fc;
  font-size: 12px;
  line-height: 16px;
  background: #e5efff;
  flex-shrink: 0;

  &.skill {
    color: #2475fc;
    background: #e5efff;
  }
}

.skill-desc {
  margin-top: 4px;
  font-size: 14px;
  line-height: 22px;
  color: #8c8c8c;
  display: -webkit-box;
  overflow: hidden;
  word-break: break-word;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
}

.skill-status-text-group {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;

  &.inactive {
    .skill-status-dot {
      background: #d9d9d9;
    }

    .skill-status-text {
      color: #8c8c8c;
    }
  }
}

.skill-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00bc7d;
}

.skill-status-text {
  font-size: 12px;
  line-height: 20px;
  color: #009966;
}

.skill-content-block {
  margin: 0 16px;
}

.skill-loading,
.tool-loading-card {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 72px;
}

.content-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 232px;
  overflow-y: auto;
}

.content-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 64px;
  padding: 9px;
  border-radius: 8px;
  background: #f2f4f7;
  border: 1px solid #e4e6eb;
}

.content-item-main {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.content-item-text {
  min-width: 0;
  flex: 1;
}

.content-item-title {
  font-size: 14px;
  line-height: 22px;
  color: #262626;
  font-weight: 400;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.content-item-desc {
  margin-top: 2px;
  font-size: 12px;
  line-height: 20px;
  color: #7a8699;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.library-avatar {
  width: 32px;
  height: 32px;
  border-radius: 16px;
  overflow: hidden;
  flex-shrink: 0;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.file-ext-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #fb4d4f;
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  flex-shrink: 0;
}

.remove-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: 0;
  padding: 0;
  border-radius: 8px;
  background: transparent;
  color: #fb363f;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.2s;

  &:hover {
    background: #e4e6eb;
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
}

.empty-tip {
  font-size: 13px;
  line-height: 22px;
  color: #8c8c8c;
  min-height: 64px;
  padding: 20px 12px;
  border-radius: 8px;
  border: 1px solid #e4e6eb;
  background: #f2f4f7;
}

.skill-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin: 12px 16px 0;
  padding-top: 11px;
  padding-bottom: 16px;
  border-top: 1px dashed #f1f5f9;
}

.single-action {
  margin-top: 12px;
}

.skill-footer-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.skill-action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 32px;
  padding: 5px 16px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }

  &:hover {
    color: #fb363f;
  }
}

.skill-action-btn-primary {
  background: #e5efff;
  color: #2475fc;

  &:hover {
    color: #2475fc;
    background: #d8e8ff;
  }
}

.skill-action-btn-danger {
  &:hover {
    color: #fb363f;
  }
}


.tool-loading-card {
  padding: 12px 0 4px;
}
</style>
