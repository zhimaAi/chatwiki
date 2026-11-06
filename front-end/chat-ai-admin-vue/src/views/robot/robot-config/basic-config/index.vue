<template>
  <a-tabs v-if="!isWorkspace" class="tab-wrapper" @change="handleTabChange" v-model:activeKey="activeKey">
    <a-tab-pane :key="1" :tab="t('tab_basic_config')"></a-tab-pane>
    <a-tab-pane v-if="false" :key="3" :tab="t('tab_permission_management')"></a-tab-pane>
    <a-tab-pane :key="2" :tab="t('tab_common_questions')"></a-tab-pane>
  </a-tabs>
  <div class="robot-config-page" :class="{ 'workspace-config': isWorkspace }" v-if="isWorkspace || activeKey == 1">
    <div class="scroll-box" :class="{ 'is-scrolling': isWorkspace && isConfigScrolling }" ref="scrollBox" @scroll="handleConfigScroll">
      <div class="setting-box">
        <BasicConfig :inline-summary="isWorkspace" />
      </div>
      <template v-if="isWorkspace">
        <div class="setting-box" v-if="!isWorkflowRobot"><ChatModeCard /></div>
        <div class="setting-box" v-if="!isWorkflowRobot"><ModelSettings inline-summary /></div>
        <div class="setting-box" v-if="!isWorkflowRobot"><PromptCard :robot-list="workFlowRobotList" /></div>
        <div class="setting-box" v-if="!isWorkflowRobot"><RecallSettingsCard /></div>
        <div class="setting-box"><LanguageSetting variant="workspace" /></div>
        <div class="setting-box" v-if="!isWorkflowRobot"><VariableSetting compact /></div>
        <div class="setting-box" v-if="!isWorkflowRobot"><ChatCache compact /></div>
      </template>
      <template v-else>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <AssociatedKnowledgeBase />
        </div>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <ModelSettings />
        </div>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <ChatMode />
        </div>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <ChatCache />
        </div>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <SystemPromptWords :robotList="workFlowRobotList" />
        </div>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <Skill :robotList="workFlowRobotList" />
        </div>
        <div class="setting-box" v-if="!isWorkflowRobot">
          <DataBase />
        </div>
        <div class="setting-box">
          <LanguageSetting />
        </div>
        <!-- <div class="setting-box">
          <WelcomeWords />
        </div> -->

        <div class="setting-box" v-if="!isWorkflowRobot">
          <VariableSetting />
        </div>
      </template>

      <!-- <div class="setting-box">
        <UnknownProblemPrompt />
      </div> -->
      <div class="setting-box">
        <SensitiveWords :compact="isWorkspace" />
      </div>
      <!--
      <div class="setting-box">
        <MarkdownSetting />
      </div>
      -->
      <div class="setting-box" v-if="!isWorkflowRobot">
        <ProblemOptimization :compact="isWorkspace" />
      </div>
      <div class="setting-box" v-if="!isWorkflowRobot">
        <SuggestedIssues :robotList="workFlowRobotList" :compact="isWorkspace" />
      </div>
      <div class="setting-box" v-if="!isWorkflowRobot">
        <DisplayAitations :compact="isWorkspace" />
      </div>
      <!-- <div class="setting-box">
        <AnswerPrompt />
      </div> -->
      <div class="setting-box">
        <ShowLike :compact="isWorkspace" />
      </div>
    </div>
  </div>
  <div class="robot-config-page" v-if="!isWorkspace && activeKey == 2">
    <CommonProblem />
  </div>
  <div class="robot-config-page" v-if="!isWorkspace && activeKey == 3">
    <rolePermission />
  </div>
</template>

<script setup>
import { reactive, ref, toRaw, provide, onMounted, onBeforeUnmount, computed, inject } from 'vue'
import { storeToRefs } from 'pinia'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import BasicConfig from './components/basic-config.vue'
import SystemPromptWords from './components/system-prompt-words.vue'
import ModelSettings from './components/model-settings.vue'
import AssociatedKnowledgeBase from './components/associated-knowledge-base/index.vue'
import DataBase from './components/data-base/index.vue'
import ProblemOptimization from './components/problem-optimization.vue'
import SuggestedIssues from './components/suggested-issues.vue'
import SensitiveWords from './components/sensitive-words.vue'
import ChatMode from './components/chat-mode/index.vue'
import CommonProblem from './components/language-setting/common-problem.vue'
import DisplayAitations from './components/display-aitations.vue'
import { saveRobot, getRobotList } from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'
import { useModelStore } from '@/stores/modules/model'
import ShowLike from './components/show-like.vue'
import rolePermission from './role-permission.vue'
import Skill from './components/skill/index.vue'
import ChatCache from './components/chat-cache.vue'
import VariableSetting from './components/variable-setting/index.vue'
import LanguageSetting from './components/language-setting/index.vue'
import RecallSettingsCard from '@/views/customer-service/basic-settings/components/recall-settings-card.vue'
import PromptCard from '@/views/customer-service/basic-settings/components/prompt-card.vue'
import ChatModeCard from '@/views/customer-service/basic-settings/components/chat-mode-card.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.index')
const props = defineProps({ variant: { type: String, default: 'legacy' } })
const isWorkspace = computed(() => props.variant === 'workspace')

const emit = defineEmits(['tabChange'])
const customerContext = inject('customerRobotContext', null)
const robotStore = useRobotStore()
const activeLocalKey = '/robot/config/basic-config/activeKey'

const scrollBox = ref(null)
const isConfigScrolling = ref(false)
let scrollIdleTimer
const handleConfigScroll = () => {
  if (!isWorkspace.value) return
  isConfigScrolling.value = true
  clearTimeout(scrollIdleTimer)
  scrollIdleTimer = setTimeout(() => { isConfigScrolling.value = false }, 800)
}
onBeforeUnmount(() => { clearTimeout(scrollIdleTimer) })

const { robotInfo } = storeToRefs(robotStore)
const { getRobot, fetchChatVariables } = robotStore

const modelStore = useModelStore()
const isWorkflowRobot = computed(() => Number(robotInfo.value.application_type || 0) === 1)

const scrollBoxToBottom = () => {
  scrollBox.value.scrollTop = scrollBox.value.scrollHeight
}
const activeKey = ref(isWorkspace.value ? 1 : +localStorage.getItem(activeLocalKey) || 1)
const saveLoading = ref(false)

// 基本配置
const formState = reactive(robotInfo.value)

const updateRobotInfo = (val) => {
  if (customerContext) return customerContext.updateRobotInfo(val)
  let newState = JSON.parse(JSON.stringify(val))
  // 对机器人头像特殊处理
  if (val.robot_avatar && val.robot_avatar instanceof File) {
    newState.robot_avatar = new File([val.robot_avatar], val.robot_avatar.name)
  }
  if(!newState.op_type_relation_library && formState.op_type_relation_library){
    delete formState.op_type_relation_library
  }
  Object.assign(formState, newState)

  return saveForm()
}
const workFlowRobotList = ref([])
const getWorkFlowRobotList = ()=>{
  getRobotList({ application_type: 1 }).then(res=>{
    workFlowRobotList.value = res.data || []
  })
}

provide('robotInfo', {
  robotInfo: formState,
  updateRobotInfo,
  scrollBoxToBottom,
  getRobot: customerContext ? customerContext.refreshRobot : getRobot,
  getWorkFlowRobotList,
})

const saveForm = () => {
  let formData = JSON.parse(JSON.stringify(toRaw(formState)))
  // 有机器人头像就赋值
  if (formState.robot_avatar) {
    if (formState.robot_avatar instanceof File) {
      formData.robot_avatar = new File([formState.robot_avatar], formState.robot_avatar.name)
      delete formData.robot_avatar_url
    } else {
      formData.robot_avatar_url = formState.robot_avatar
    }
  }
  let welcomes = formData.welcomes

  welcomes.question = welcomes.question.map((item) => {
    return item.content
  })

  formData.welcomes = JSON.stringify(welcomes)
  formData.prompt_struct = JSON.stringify(formData.prompt_struct)

  let unknown_question_prompt = formData.unknown_question_prompt

  unknown_question_prompt.question = unknown_question_prompt.question.map((item) => {
    return item.content
  })

  formData.unknown_question_prompt = JSON.stringify(unknown_question_prompt)

  formData.cache_config = JSON.stringify(formData.cache_config)

  saveLoading.value = true

  return saveRobot(formData)
    .then((res) => {
      if (res.res != 0) {
        message.error(res.msg)
        return false
      }

      message.success(t('msg_save_success'))

      getRobot(formState.id)
      robotStore.getRobotLists()
      return true
    })
    .catch(() => false)
    .finally(() => {
      saveLoading.value = false
    })
}

const handleTabChange = () => {
  if (isWorkspace.value) return
  emit('tabChange', activeKey.value)
  localStorage.setItem(activeLocalKey, activeKey.value)
}


onMounted(()=>{
  getWorkFlowRobotList()
  modelStore.getAllmodelList()
  fetchChatVariables()
})

</script>

<style lang="less" scoped>
.robot-config-page {
  height: calc(100% - 46px);
  width: 100%;
  padding: 8px 10px 24px 24px;
  padding-bottom: 0;
  overflow: hidden;
  background-color: #fff;

  .scroll-box {
    height: 100%;
    overflow-y: auto;
  }

  .setting-box {
    margin-bottom: 16px;
  }
}
.robot-config-page.workspace-config {
  height: 100%;
  padding: 24px;

  .scroll-box {
    display: flex;
    flex-direction: column;
    gap: 48px;
    padding-right: 16px;
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
    &::-webkit-scrollbar { width: 6px; height: 6px; }
    &::-webkit-scrollbar-track { background: transparent; }
    &::-webkit-scrollbar-thumb { background: transparent; border-radius: 4px; }
    &::-webkit-scrollbar-button { display: none; width: 0; height: 0; }
    &.is-scrolling {
      scrollbar-color: #d1d5db transparent;
      &::-webkit-scrollbar-thumb { background: #d1d5db; }
    }
    @supports selector(::-webkit-scrollbar) {
      scrollbar-width: auto;
      scrollbar-color: auto;
      &.is-scrolling { scrollbar-color: auto; }
    }
  }
}
.workspace-config {
  .scroll-box > .setting-box { flex-shrink: 0; margin: 0; }
  :deep(.setting-box) { margin-bottom: 0; }
  :deep(.edit-box:not(.profile-summary)),
  :deep(.workspace-language) { height: auto; background: #fff; border-radius: 0; }

  :deep(.edit-box:not(.profile-summary) > .edit-box-header),
  :deep(.workspace-language > .edit-box-header) {
    min-height: 40px;
    height: auto;
    gap: 12px;
    padding: 4px 20px;
    border-radius: 8px;
    background: #f5f5f5;
    .left-box { min-width: 0; }
    .right-box, .extra-box { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
    .edit-box-icon, .left-box > .svg-icon { display: none; }
    .edit-box-title { color: #222; font-size: 16px; font-weight: 600; line-height: 24px; }
  }
  :deep(.edit-box-header .ant-btn) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
    height: 32px;
    min-width: 56px;
    padding: 0 12px;
    font-size: 14px;
    border-radius: 6px;
    box-shadow: none;
    &:not(.ant-btn-primary) { color: #222; border-color: #d8d8d8; background: #fff; }
  }
  :deep(.edit-box:not(.profile-summary) > .edit-box-body) { padding: 24px 20px 0; }
  :deep(.title-tip) {
    display: inline-flex;
    align-items: center;
    .anticon-question-circle {
      width: 16px;
      height: 16px;
      background: url('@/assets/img/customer-service/basic-settings/help.svg') center no-repeat;
      svg { visibility: hidden; }
    }
  }
  :deep(.workspace-summary) {
    display: grid;
    grid-template-columns: repeat(var(--summary-columns, 4), minmax(0, 1fr));
    gap: 24px;
    font-size: 14px;
    line-height: 22px;
    .summary-item, .set-item { display: flex; flex-direction: column; align-items: flex-start; gap: 8px; min-width: 0; }
    .summary-label { color: #5a5a5a; }
    .summary-value { font-weight: 500; color: #080808; overflow-wrap: anywhere; }
    .summary-status {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      color: #8c8c8c;
      &::before { content: ''; width: 8px; height: 8px; flex-shrink: 0; border-radius: 4px; background: #bfbfbf; }
      &.is-enabled { color: #007012; &::before { background: #00d722; } }
    }
  }
  :deep(.setting-info-block.inline-summary) { padding: 0; }
  :deep(.prompt-preview) { padding: 16px; border: 1px solid #e8e8e8; border-radius: 8px; color: #5a5a5a; }
  :deep(.ant-switch) {
    min-width: 36px;
    width: 36px;
    height: 20px;
    line-height: 20px;
    flex-shrink: 0;
    .ant-switch-handle { width: 16px; height: 16px; top: 2px; }
    .ant-switch-inner { display: none; }
    &.ant-switch-checked .ant-switch-handle { inset-inline-start: calc(100% - 18px); }
  }
  :deep(.compact-setting .edit-box-body .robot-prompt) { white-space: normal; overflow: visible; text-overflow: clip; line-height: 24px; color: #7a7a7a; }
  :deep(.current-mode) { font-size: 12px; color: #7a7a7a; }
  :deep(.variable-setting .form-box) { padding: 0; }
  :deep(.variable-setting .variable-item) {
    gap: 0;
    border: 1px solid #d8d8d8;
    border-radius: 6px;
    .variable-info-block { border: 0; }
    .btn-box { border-radius: 0 6px 6px 0; }
  }
  :deep(.workspace-cache .robot-info-box) { display: block; color: #7a7a7a; line-height: 24px; }
  :deep(.workspace-cache .set-block) { display: flex; align-items: center; flex-wrap: wrap; gap: 8px; margin-top: 24px; }
  :deep(.workspace-cache .switch-item) { top: 4px; right: 20px; }
  :deep(.workspace-cache .ant-input-group) { display: inline-flex; align-items: center; gap: 8px; width: auto; }
  :deep(.workspace-language) {
    > .edit-box-header .left-box { gap: 8px; font-size: 16px; line-height: 24px; color: #222; font-weight: 600; }
    > .edit-box-header .left-box > svg { display: none; }
    > .edit-box-body { padding: 24px 20px 0; display: flex; flex-direction: column; gap: 24px; }
    .tabs-box { padding: 0 20px; }
    .edit-box { display: grid; grid-template-columns: 176px minmax(0, 1fr); gap: 20px; position: relative; overflow: visible; }
    .edit-box > .edit-box-header { min-height: 0; padding: 0; margin: 0; align-items: flex-start; background: transparent; }
    .edit-box > .edit-box-header .edit-box-title { font-weight: 400; color: #7a7a7a; line-height: 24px; }
    .edit-box > .edit-box-body { min-width: 0; padding: 0; color: #5a5a5a; line-height: 24px; }
    .edit-box:last-child > .edit-box-body { padding-top: 32px; }
    .switch-item { position: absolute; left: 196px; top: 2px; }
    .question-item { background: #f5f5f5; border: 0; border-radius: 6px; padding: 8px 12px; }
    .ant-input { max-width: 100%; }
  }
}
@media (max-width: 1200px) {
  .robot-config-page.workspace-config {
    padding: 24px 16px;
    .scroll-box { padding-right: 8px; }
    :deep(.workspace-summary) { grid-template-columns: repeat(2, minmax(0, 1fr)); }
    :deep(.edit-box:not(.profile-summary) > .edit-box-header) { flex-wrap: wrap; }
    :deep(.workspace-language .edit-box) { grid-template-columns: 128px minmax(0, 1fr); gap: 12px; }
    :deep(.workspace-language .switch-item) { left: 140px; }
    :deep(.basic-config .robot-form-box) { flex-direction: column; gap: 24px; }
  }
}
::v-deep(.ant-tabs-nav-wrap) {
  padding-left: 24px;
}
.tab-wrapper ::v-deep(.ant-tabs-nav){
  margin-bottom: 0;
}
</style>
