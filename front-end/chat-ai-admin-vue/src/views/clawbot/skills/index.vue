<template>
  <div class="skills-page">
    <div class="page-header">
      <div class="page-title">{{ t('page_title') }}</div>
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
          </a-menu>
        </template>
      </a-dropdown>
    </div>

    <div class="skill-list">
      <!-- 固定技能：查询知识库 -->
      <div class="skill-card">
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
      <div class="skill-card">
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

      <!-- 已关联的 WorkFlow 工具列表 -->
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
            <svg-icon name="delete-line" />
          </div>
        </div>
      </div>

      <!-- 动态技能列表 -->
      <div v-for="item in skills" :key="item.id" class="skill-card">
        <div class="skill-main">
          <div class="skill-title-row">
            <div class="skill-title">{{ item.title }}</div>
            <div class="skill-tag" :class="item.tagType">{{ item.tag }}</div>
          </div>
          <div class="skill-desc">{{ item.desc }}</div>
        </div>
        <div class="skill-actions">
          <a-switch
            class="skill-switch"
            :checked="item.enabled"
            :checked-children="t('switch_on')"
            :un-checked-children="t('switch_off')"
            disabled
          />
          <div v-if="item.removable" class="delete-action">
            <svg-icon name="delete-line" />
          </div>
        </div>
      </div>
    </div>
    <!-- 新增 Skill 弹窗 -->
    <AddSkillModal
      v-model:visible="skillModalVisible"
      @confirm="handleSkillConfirm"
    />

    <!-- 新增 Tool 弹窗 -->
    <AddToolModal
      v-model:visible="toolModalVisible"
      :robotId="currentAssistant?.id"
      :workFlowIds="robotInfo?.work_flow_ids"
      @confirm="handleToolConfirm"
    />
  </div>
</template>

<script setup>
import { computed, ref, watch, createVNode } from 'vue'
import { DownOutlined, ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { storeToRefs } from 'pinia'
import { getRobotList, relationWorkFlow } from '@/api/robot/index'
import AddSkillModal from './components/AddSkillModal.vue'
import AddToolModal from './components/AddToolModal.vue'

const { t } = useI18n('views.clawbot.skills.index')
const clawbotStore = useClawbotStore()
const { robotInfo, currentAssistant } = storeToRefs(clawbotStore)
const { updateClawbotConf, fetchRobotInfo } = clawbotStore

// 查询知识库开关：search_knowledge_close=0 表示开启，=1 表示关闭
const knowledgeEnabled = computed(() => !Number(robotInfo.value?.search_knowledge_close || 0))
// 查询本地文档开关
const localDocsEnabled = computed(() => !Number(robotInfo.value?.query_local_docs_close || 0))

const toggleKnowledge = async (checked) => {
  if (!currentAssistant.value?.id) {
    return
  }

  await updateClawbotConf({
    id: currentAssistant.value.id,
    search_knowledge_close: checked ? 0 : 1,
    query_local_docs_close: localDocsEnabled.value ? 0 : 1
  })
}

const toggleLocalDocs = async (checked) => {
  if (!currentAssistant.value?.id) {
    return
  }

  await updateClawbotConf({
    id: currentAssistant.value.id,
    search_knowledge_close: knowledgeEnabled.value ? 0 : 1,
    query_local_docs_close: checked ? 0 : 1
  })
}

const workFlowSkills = ref([])

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
        .filter(item => savedIds.includes(item.id))
        .map(item => ({
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
watch(() => robotInfo.value?.work_flow_ids, () => {
  loadWorkFlowSkills()
}, { immediate: true })

const skillModalVisible = ref(false)
const toolModalVisible = ref(false)

const handleAddMenuClick = ({ key }) => {
  if (key === 'skill') {
    skillModalVisible.value = true
  } else if (key === 'tool') {
    toolModalVisible.value = true
  }
}

const handleSkillConfirm = (formData) => {
  // TODO: 对接新增 Skill 接口
  console.log('新增 Skill:', formData)
}

const handleToolConfirm = async () => {
  // 保存成功后刷新 robotInfo 再重新加载列表
  await fetchRobotInfo()
  loadWorkFlowSkills()
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
      const newIds = currentIds.filter(item => item !== String(id))
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

const skills = [
  // {
  //   id: 3,
  //   title: '订单查询',
  //   tag: 'SKILL',
  //   tagType: 'skill',
  //   desc: '查询用户订单信息，包括订单状态、物流信息等',
  //   enabled: true,
  //   removable: true,
  // },
  // {
  //   id: 4,
  //   title: '客户信息同步',
  //   tag: 'TOOL',
  //   tagType: 'tool',
  //   desc: '自动同步客户信息到CRM系统',
  //   enabled: true,
  //   removable: true,
  // }
]
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
    font-weight: 600;
    line-height: 28px;
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
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    background: #f5f5f5;
    color: #ff4d4f;
  }

  :deep(svg) {
    width: 16px;
    height: 16px;
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
</style>
