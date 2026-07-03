<template>
  <div class="skills-page">
    <div class="page-header">
      <div class="page-title">
        <a-segmented :value="2" :options="titleOptios" />
      </div>
      <a-dropdown placement="bottomRight">
        <a-button type="primary" class="add-btn">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('btn_add_skill') }}
          <DownOutlined class="btn-arrow" />
        </a-button>
        <template #overlay>
          <a-menu @click="handleAddMenuClick">
            <!-- <a-menu-item key="skill">新增 Skill</a-menu-item> -->
            <a-menu-item key="tool">{{ t('menu_add_tool') }}</a-menu-item>
            <a-menu-item key="skill">{{ t('menu_add_skill') }}</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>

    <div class="type-list-box">
      <button
        v-for="item in skillSourceTabs"
        :key="item.value"
        type="button"
        class="type-tab"
        :class="{ active: activeSkillSource === item.value }"
        @click="handleSourceChange(item.value)"
      >
        {{ item.label }}
      </button>
    </div>
    <div class="skill-list">
      <!-- 固定技能：查询知识库 -->
      <div class="skill-card" v-if="activeSkillSource == 'all'">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ t('title_query_knowledge_base') }}</div>
            <div class="skill-tag tool">{{ t('tag_tool') }}</div>
          </div>
          <div class="skill-desc">{{ t('desc_query_knowledge_base') }}</div>
        </div>
        <div class="skill-actions">
          <a-switch
            class="skill-switch"
            :checked="knowledgeEnabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleKnowledge"
          />
        </div>
      </div>

      <!-- 固定技能：查询本地文档 -->
      <div class="skill-card" v-if="activeSkillSource == 'all'">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ t('title_query_local_docs') }}</div>
            <div class="skill-tag skill">{{ t('tag_skill') }}</div>
          </div>
          <div class="skill-desc">{{ t('desc_query_local_docs') }}</div>
        </div>
        <div class="skill-actions">
          <a-switch
            class="skill-switch"
            :checked="localDocsEnabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleLocalDocs"
          />
        </div>
      </div>

      <!-- 固定技能：写文件 -->
      <div class="skill-card" v-if="activeSkillSource == 'all'">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ t('title_agent_write_file') }}</div>
            <div class="skill-tag tool">{{ t('tag_tool') }}</div>
          </div>
          <div class="skill-desc">{{ t('desc_agent_write_file') }}</div>
        </div>
        <div class="skill-actions">
          <a-switch
            class="skill-switch"
            :checked="writeFileEnabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleWriteFile"
          />
        </div>
      </div>

      <!-- 固定技能：编辑文件 -->
      <div class="skill-card" v-if="activeSkillSource == 'all'">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ t('title_agent_edit_file') }}</div>
            <div class="skill-tag tool">{{ t('tag_tool') }}</div>
          </div>
          <div class="skill-desc">{{ t('desc_agent_edit_file') }}</div>
        </div>
        <div class="skill-actions">
          <a-switch
            class="skill-switch"
            :checked="editFileEnabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleEditFile"
          />
        </div>
      </div>

      <!-- 固定技能：执行命令 -->
      <div class="skill-card" v-if="activeSkillSource == 'all'">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ t('title_agent_execute') }}</div>
            <div class="skill-tag tool">{{ t('tag_tool') }}</div>
          </div>
          <div class="skill-desc">{{ t('desc_agent_execute') }}</div>
        </div>
        <div class="skill-actions">
          <a-switch
            class="skill-switch"
            :checked="executeEnabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleExecute"
          />
        </div>
      </div>

      <!-- 固定技能：从商品库推荐商品 -->
      <div class="skill-card" v-if="activeSkillSource == 'all'">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ t('title_goods_recommend') }}</div>
            <div class="skill-tag tool">{{ t('tag_tool') }}</div>
          </div>
          <div class="skill-desc">{{ t('desc_goods_recommend') }}</div>
        </div>
        <div class="skill-actions">
          <span class="recommend-scope-btn" @click="handleOpenScopeModal">{{ t('btn_recommend_scope') }}</span>
          <a-switch
            class="skill-switch"
            :checked="goodsRecommendEnabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            @change="toggleGoodsRecommend"
          />
        </div>
      </div>

      <div v-if="skillListLoading" class="skill-loading">
        <a-spin />
      </div>
      <template v-else>
        <!-- 动态技能列表 -->
        <div v-for="item in skills" :key="item.id" class="skill-card">
          <div class="skill-main">
            <div class="skill-title-row">
              <div class="skill-title">{{ item.title }}</div>
              <div v-if="item.sourceLabel" class="skill-source-tag" :class="item.sourceClass">
                {{ item.sourceLabel }}
              </div>
              <div class="skill-tag skill">{{ t('tag_skill') }}</div>
            </div>
            <div class="skill-desc">{{ item.desc }}</div>
          </div>
          <div class="skill-actions">
            <div v-if="item.editable" class="icon-action" @click="handleEditSkill(item)">
              <EditOutlined />
            </div>
            <div v-if="item.removable" class="icon-action delete-action" @click="handleRemoveSkill(item)">
              <DeleteOutlined />
            </div>
          </div>
        </div>
      </template>

      <!-- 已关联的 WorkFlow 工具列表 -->
      <template v-if="activeSkillSource == 'all'">
        <div v-for="item in workFlowSkills" :key="item.id" class="skill-card">
          <div class="skill-main">
            <div class="skill-title-row">
              <div class="skill-title">{{ item.name }}</div>
              <div class="skill-tag tool">{{ t('tag_tool') }}</div>
            </div>
            <div class="skill-desc">{{ item.desc || '—' }}</div>
          </div>
          <div class="skill-actions">
            <div class="delete-action" @click="handleRemoveWorkFlow(item.id)">
              <DeleteOutlined />
            </div>
          </div>
        </div>
      </template>

    </div>
    <SelectSkillModal
      v-model:visible="selectSkillModalVisible"
      :robotId="currentAssistant?.id"
      :refreshKey="selectSkillRefreshKey"
      @create="handleOpenUploadSkillModal"
      @confirm="handleSelectSkillConfirm"
    />

    <!-- 新增 Skill 弹窗 -->
    <AddSkillModal v-model:visible="skillModalVisible" @confirm="handleSkillConfirm" />

    <!-- 新增 Tool 弹窗 -->
    <AddToolModal
      v-model:visible="toolModalVisible"
      :robotId="currentAssistant?.id"
      :workFlowIds="robotInfo?.work_flow_ids"
      @confirm="handleToolConfirm"
    />

    <UploadSkillZipModal
      v-model:visible="uploadSkillZipVisible"
      :robotId="currentAssistant?.id"
      :skill-id="editingSkillId"
      @confirm="handleUploadSkillConfirm"
    />

    <!-- 推荐范围弹窗 -->
    <GoodsRecommendScopeModal v-model:visible="scopeModalVisible" @confirm="handleScopeConfirm" />
  </div>
</template>

<script setup>
import { computed, ref, watch, createVNode } from 'vue'
import { DeleteOutlined, DownOutlined, EditOutlined, ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { storeToRefs } from 'pinia'
import { getRobotList, relationWorkFlow } from '@/api/robot/index'
import { deleteClawbotSkill, getClawbotSkillList } from '@/api/clawbot'
import AddSkillModal from './components/AddSkillModal.vue'
import AddToolModal from './components/AddToolModal.vue'
import GoodsRecommendScopeModal from './components/GoodsRecommendScopeModal.vue'
import SelectSkillModal from './components/SelectSkillModal.vue'
import UploadSkillZipModal from './components/UploadSkillZipModal.vue'

const { t } = useI18n('views.clawbot.skills.index')
const clawbotStore = useClawbotStore()
const { robotInfo, currentAssistant } = storeToRefs(clawbotStore)
const { updateClawbotConf, fetchRobotInfo } = clawbotStore

const titleOptios = ref([
  // {
  //   label: 'skill市场',
  //   value: 1
  // },
  {
    label: '我的skill',
    value: 2
  }
])

const skillSourceTabs = [
  { label: t('tab_all'), value: 'all' },
  { label: t('tab_market'), value: 'market' },
  { label: t('tab_mine'), value: 'mine' }
]

// 查询知识库开关：search_knowledge_close=0 表示开启，=1 表示关闭
const knowledgeEnabled = computed(() => !Number(robotInfo.value?.search_knowledge_close || 0))
// 查询本地文档开关
const localDocsEnabled = computed(() => !Number(robotInfo.value?.query_local_docs_close || 0))
// Agent 工具开关：1 表示开启，0 表示关闭
const writeFileEnabled = computed(() => Number(robotInfo.value?.open_agent_write_file_tool || 0) === 1)
const executeEnabled = computed(() => Number(robotInfo.value?.open_agent_execute_tool || 0) === 1)
const editFileEnabled = computed(() => Number(robotInfo.value?.open_agent_edit_file_tool || 0) === 1)

function buildConfPayload(overrides = {}) {
  return {
    id: currentAssistant.value?.id,
    search_knowledge_close: knowledgeEnabled.value ? 0 : 1,
    query_local_docs_close: localDocsEnabled.value ? 0 : 1,
    open_agent_write_file_tool: writeFileEnabled.value ? 1 : 0,
    open_agent_execute_tool: executeEnabled.value ? 1 : 0,
    open_agent_edit_file_tool: editFileEnabled.value ? 1 : 0,
    goods_lib_recommend_switch: goodsRecommendEnabled.value ? 1 : 0,
    goods_lib_recommend_group_ids: robotInfo.value?.goods_lib_recommend_group_ids || '',
    ...overrides
  }
}
// 商品库推荐开关
const goodsRecommendEnabled = computed(() => Number(robotInfo.value?.goods_lib_recommend_switch || 0) === 1)

// 推荐范围弹窗状态
const scopeModalVisible = ref(false)

const toggleKnowledge = async (checked) => {
  if (!currentAssistant.value?.id) {
    return
  }

  await updateClawbotConf(buildConfPayload({ search_knowledge_close: checked ? 0 : 1 }))
}

const toggleLocalDocs = async (checked) => {
  if (!currentAssistant.value?.id) return
  await updateClawbotConf(buildConfPayload({ query_local_docs_close: checked ? 0 : 1 }))
}

const toggleWriteFile = async (checked) => {
  if (!currentAssistant.value?.id) return
  await updateClawbotConf(buildConfPayload({ open_agent_write_file_tool: checked ? 1 : 0 }))
}

const toggleExecute = async (checked) => {
  if (!currentAssistant.value?.id) return
  await updateClawbotConf(buildConfPayload({ open_agent_execute_tool: checked ? 1 : 0 }))
}

const toggleEditFile = async (checked) => {
  if (!currentAssistant.value?.id) return

  await updateClawbotConf(buildConfPayload({ open_agent_edit_file_tool: checked ? 1 : 0 }))
}

const toggleGoodsRecommend = async (checked) => {
  if (!currentAssistant.value?.id) {
    return
  }

  await updateClawbotConf(buildConfPayload({ goods_lib_recommend_switch: checked ? 1 : 0 }))
}

const handleOpenScopeModal = () => {
  scopeModalVisible.value = true
}

const handleScopeConfirm = () => {
  scopeModalVisible.value = false
}

const workFlowSkills = ref([])
const activeSkillSource = ref('all')
const skillListLoading = ref(false)
const uploadedSkills = ref([])
const selectSkillModalVisible = ref(false)
const selectSkillRefreshKey = ref(0)
let skillListRequestSeq = 0

// 加载已关联的 WorkFlow 工具列表
const loadWorkFlowSkills = async () => {
  const workFlowIds = robotInfo.value?.work_flow_ids
  if (!workFlowIds) {
    workFlowSkills.value = []
    return
  }
  try {
    const res = await getRobotList({ application_type: 1 })
    if (res && res.res === 0) {
      const allTools = res.data || []
      const savedIds = workFlowIds.split(',').filter(Boolean)
      workFlowSkills.value = allTools
        .filter((item) => savedIds.includes(item.id))
        .map((item) => ({
          id: item.id,
          name: item.robot_name,
          desc: item.robot_intro || ''
        }))
    }
  } catch (err) {
    console.error('加载已关联工具失败', err)
  }
}

// robotInfo 变化时自动加载
watch(
  () => robotInfo.value?.work_flow_ids,
  () => {
    loadWorkFlowSkills()
  },
  { immediate: true }
)

watch(
  [() => currentAssistant.value?.id, activeSkillSource],
  () => {
    loadSkillList()
  },
  { immediate: true }
)

const skillModalVisible = ref(false)
const toolModalVisible = ref(false)
const uploadSkillZipVisible = ref(false)
const editingSkillId = ref(0)

const handleAddMenuClick = ({ key }) => {
  if (key === 'skill') {
    handleOpenSelectSkillModal()
  } else if (key === 'tool') {
    toolModalVisible.value = true
  } else if (key === 'zip') {
    editingSkillId.value = 0
    uploadSkillZipVisible.value = true
  }
}

const handleSkillConfirm = (formData) => {
  // TODO: 对接新增 Skill 接口
  console.log('新增 Skill:', formData)
  loadSkillList()
  if (selectSkillModalVisible.value) {
    selectSkillRefreshKey.value += 1
  }
}

const handleToolConfirm = async () => {
  // 保存成功后刷新 robotInfo 再重新加载列表
  await fetchRobotInfo()
  loadWorkFlowSkills()
}

const handleUploadSkillConfirm = async () => {
  await fetchRobotInfo()
  await loadSkillList()
  if (selectSkillModalVisible.value) {
    selectSkillRefreshKey.value += 1
  }
}

const handleSourceChange = (source) => {
  if (activeSkillSource.value === source) {
    return
  }
  activeSkillSource.value = source
}

// 移除已关联的 WorkFlow 工具
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
        if (res && res.res === 0) {
          message.success(t('msg_remove_success'))
          await fetchRobotInfo()
          loadWorkFlowSkills()
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

async function loadSkillList() {
  if (!currentAssistant.value?.id) {
    uploadedSkills.value = []
    return
  }

  const assistantId = currentAssistant.value.id
  const source = activeSkillSource.value
  const requestSeq = ++skillListRequestSeq
  skillListLoading.value = true
  uploadedSkills.value = []
  try {
    const res = await getClawbotSkillList({
      id: assistantId,
      source
    })
    if (requestSeq !== skillListRequestSeq || currentAssistant.value?.id !== assistantId || activeSkillSource.value !== source) {
      return
    }
    if (res && res.res === 0) {
      uploadedSkills.value = (res.data || []).filter((item) => Number(item.is_selected) === 1).map((item, index) => {
        const isMine = Number(item.source_type) === 1 || Number(item.is_mine) === 1
        return {
          id: `${item.source_type || item.source || 'skill'}-${item.skill_id || 0}-${item.skill_name || index}`,
          skillId: item.skill_id,
          title: item.remark_name || item.skill_name || '—',
          desc: item.intro || item.description || '—',
          sourceLabel: isMine ? t('tab_mine') : t('tag_market_skill'),
          sourceClass: isMine ? 'mine' : 'market',
          editable: isMine && Number(item.skill_id) > 0,
          removable: isMine && Number(item.skill_id) > 0,
          raw: item
        }
      })
    } else {
      message.error(res?.msg || t('msg_fetch_skill_failed'))
    }
  } catch (err) {
    if (requestSeq !== skillListRequestSeq || currentAssistant.value?.id !== assistantId || activeSkillSource.value !== source) {
      return
    }
    console.error('获取 Skill 列表失败', err)
    message.error(err?.msg || t('msg_fetch_skill_failed'))
  } finally {
    if (requestSeq === skillListRequestSeq) {
      skillListLoading.value = false
    }
  }
}

const handleOpenSelectSkillModal = () => {
  selectSkillModalVisible.value = true
}

const handleOpenUploadSkillModal = () => {
  editingSkillId.value = 0
  uploadSkillZipVisible.value = true
}

const handleSelectSkillConfirm = async () => {
  await loadSkillList()
}

const handleEditSkill = (item) => {
  editingSkillId.value = item.skillId
  uploadSkillZipVisible.value = true
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
        if (res && res.res === 0) {
          message.success(t('msg_remove_success'))
          loadSkillList()
        } else {
          message.error(res?.msg || t('msg_remove_failed'))
        }
      } catch (err) {
        console.error('删除 Skill 失败', err)
        message.error(err?.msg || t('msg_remove_failed'))
      }
    }
  })
}

const skills = computed(() => uploadedSkills.value)
</script>

<style lang="less" scoped>
.skills-page {
  --skills-primary: #2475fc;
  --skills-border: #d9d9d9;
  --skills-title: #262626;
  --skills-text: #595959;
  --skills-text-light: #8c8c8c;
  --skills-card-shadow: 0 4px 4px rgba(0, 0, 0, 0.08);

  position: relative;
  min-height: 100vh;
  padding: 22px 24px 32px;
  background: #fff;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -88px;
    right: -46px;
    width: 284px;
    height: 274px;
    border-radius: 999px;
    background: rgba(198, 210, 255, 0.3);
    filter: blur(64px);
    pointer-events: none;
  }
}

.page-header {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;

  .page-title {
    color: var(--skills-title);
    font-size: 20px;
    line-height: 28px;
    .ant-segmented {
      background: #edeff2;
    }
    &::v-deep(.ant-segmented-item-selected) {
      color: #2475fc;
    }
  }

  .add-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    height: 32px;
    padding: 0 16px;
    border: none;
    border-radius: 6px;
    background: var(--skills-primary);
    box-shadow: none;
    font-size: 14px;
    line-height: 22px;

    &:hover,
    &:focus {
      background: #4a8dff;
    }
  }

  .btn-arrow {
    font-size: 12px;
  }
}

.skill-list {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.type-list-box {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.type-tab {
  height: 32px;
  padding: 0 16px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--skills-text);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;

  &:hover {
    color: var(--skills-primary);
  }

  &.active {
    color: var(--skills-primary);
    background: #dce9ff;
  }
}

.skill-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 120px;
}

.skill-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  min-height: 82px;
  padding: 16px 20px;
  border-radius: 8px;
  border: 1px solid var(--skills-border);
  background: #fff;
  box-shadow: var(--skills-card-shadow);
}

.skill-main {
  flex: 1;
  min-width: 0;
}

.skill-title-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.skill-title {
  color: var(--skills-title);
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

.skill-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 18px;
  padding: 0 6px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;

  &.tool {
    color: #10ae8a;
    background: #cbfaf0;
  }

  &.skill {
    color: var(--skills-primary);
    background: #e5efff;
  }
}

.skill-source-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 18px;
  padding: 0 6px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;

  &.market {
    color: #fa8c16;
    background: #fff1dc;
  }

  &.mine {
    color: #fa8c16;
    background: #fff1dc;
  }
}

.skill-desc {
  margin-top: 4px;
  color: var(--skills-text-light);
  font-size: 14px;
  line-height: 22px;
}

.skill-actions {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 10px;
}

.delete-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  color: var(--skills-text);
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease;

  &:hover {
    background: #f5f5f5;
    color: #ff4d4f;
  }

  :deep(svg) {
    width: 16px;
    height: 16px;
  }
}

.icon-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  color: var(--skills-text);
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    background: #f5f5f5;
    color: var(--skills-primary);
  }

  &.delete-action:hover {
    color: #ff4d4f;
  }
}

@media (max-width: 768px) {
  .skills-page {
    padding: 20px 16px 24px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .page-header .add-btn {
    justify-content: center;
  }

  .skill-card {
    align-items: flex-start;
    flex-direction: column;
  }

  .skill-actions {
    width: 100%;
    justify-content: flex-end;
  }
}

.recommend-scope-btn {
  color: var(--skills-primary);
  font-size: 14px;
  cursor: pointer;
  white-space: nowrap;
  user-select: none;
}
</style>
