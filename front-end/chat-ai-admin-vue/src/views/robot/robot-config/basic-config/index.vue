<template>
  <a-tabs class="tab-wrapper" @change="handleTabChange" v-model:activeKey="activeKey">
    <a-tab-pane :key="1" :tab="t('tab_basic_config')"></a-tab-pane>
    <a-tab-pane v-if="false" :key="3" :tab="t('tab_permission_management')"></a-tab-pane>
    <a-tab-pane :key="2" :tab="t('tab_common_questions')"></a-tab-pane>
  </a-tabs>
  <div class="robot-config-page" v-if="activeKey == 1">
    <div class="scroll-box" ref="scrollBox">
      <div class="setting-box">
        <BasicConfig />
      </div>
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

      <!-- <div class="setting-box">
        <UnknownProblemPrompt />
      </div> -->
      <div class="setting-box">
        <SensitiveWords />
      </div>
      <!--
      <div class="setting-box">
        <MarkdownSetting />
      </div>
      -->
      <div class="setting-box" v-if="!isWorkflowRobot">
        <ProblemOptimization />
      </div>
      <div class="setting-box" v-if="!isWorkflowRobot">
        <SuggestedIssues :robotList="workFlowRobotList" />
      </div>
      <div class="setting-box" v-if="!isWorkflowRobot">
        <DisplayAitations />
      </div>
      <!-- <div class="setting-box">
        <AnswerPrompt />
      </div> -->
      <div class="setting-box">
        <ShowLike />
      </div>
    </div>
  </div>
  <div class="robot-config-page" v-if="activeKey == 2">
    <CommonProblem />
  </div>
  <div class="robot-config-page" v-if="activeKey == 3">
    <rolePermission />
  </div>
</template>

<script setup>
import { reactive, ref, toRaw, provide, onMounted, computed } from 'vue'
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

const { t } = useI18n('views.robot.robot-config.basic-config.index')

const robotStore = useRobotStore()
const activeLocalKey = '/robot/config/basic-config/activeKey'

const scrollBox = ref(null)

const { robotInfo } = storeToRefs(robotStore)
const { getRobot, fetchChatVariables } = robotStore

const modelStore = useModelStore()
const isWorkflowRobot = computed(() => Number(robotInfo.value.application_type || 0) === 1)

const scrollBoxToBottom = () => {
  scrollBox.value.scrollTop = scrollBox.value.scrollHeight
}
const activeKey = ref(+localStorage.getItem(activeLocalKey) || 1)
const saveLoading = ref(false)

// 基本配置
const formState = reactive(robotInfo.value)

const updateRobotInfo = (val) => {
  let newState = JSON.parse(JSON.stringify(val))
  // 对机器人头像特殊处理
  if (val.robot_avatar && val.robot_avatar instanceof File) {
    newState.robot_avatar = new File([val.robot_avatar], val.robot_avatar.name)
  }
  if(!newState.op_type_relation_library && formState.op_type_relation_library){
    delete formState.op_type_relation_library
  }
  Object.assign(formState, newState)

  saveForm()
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
  getRobot,
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

  saveRobot(formData)
    .then((res) => {
      if (res.res != 0) {
        return message.error(res.msg)
      }

      saveLoading.value = false

      message.success(t('msg_save_success'))

      getRobot(formState.id)
      robotStore.getRobotLists()
    })
    .catch(() => {
      saveLoading.value = false
    })
}

const handleTabChange = () => {
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
::v-deep(.ant-tabs-nav-wrap) {
  padding-left: 24px;
}
.tab-wrapper ::v-deep(.ant-tabs-nav){
  margin-bottom: 0;
}
</style>
