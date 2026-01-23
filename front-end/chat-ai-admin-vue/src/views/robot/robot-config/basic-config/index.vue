<template>
  <a-tabs class="tab-wrapper" @change="handleTabChange" v-model:activeKey="activeKey">
    <a-tab-pane :key="1" tab="基础配置"></a-tab-pane>
    <a-tab-pane v-if="false" :key="3" tab="权限管理"></a-tab-pane>
    <a-tab-pane :key="2" tab="常见问题"></a-tab-pane>
  </a-tabs>
  <div class="robot-config-page" v-if="activeKey == 1">
    <div class="scroll-box" ref="scrollBox">
      <div class="setting-box">
        <BasicConfig />
      </div>
      <div class="setting-box">
        <AssociatedKnowledgeBase />
      </div>
      <div class="setting-box">
        <ModelSettings />
      </div>
      <div class="setting-box">
        <ChatMode />
      </div>
      <div class="setting-box">
        <ChatCache />
      </div>
      <div class="setting-box">
        <SystemPromptWords :robotList="workFlowRobotList" />
      </div>
      <div class="setting-box">
        <Skill :robotList="workFlowRobotList" />
      </div>
      <div class="setting-box">
        <DataBase />
      </div>
      <div class="setting-box">
        <WelcomeWords />
      </div>

      <div class="setting-box">
        <UnknownProblemPrompt />
      </div>
      <div class="setting-box">
        <SensitiveWords />
      </div>
      <!--  
      <div class="setting-box">
        <MarkdownSetting />
      </div>
      -->
      <div class="setting-box">
        <ProblemOptimization />
      </div>
      <div class="setting-box">
        <SuggestedIssues />
      </div>
      <div class="setting-box">
        <DisplayAitations />
      </div>
      <div class="setting-box">
        <AnswerPrompt />
      </div>
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
import { reactive, ref, toRaw, provide, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { message } from 'ant-design-vue'
import BasicConfig from './components/basic-config.vue'
import SystemPromptWords from './components/system-prompt-words.vue'
import ModelSettings from './components/model-settings.vue'
import AssociatedKnowledgeBase from './components/associated-knowledge-base/index.vue'
import DataBase from './components/data-base/index.vue'
import WelcomeWords from './components/welcome-words.vue'
import UnknownProblemPrompt from './components/unknown-problem-prompt.vue'
import ProblemOptimization from './components/problem-optimization.vue'
import SuggestedIssues from './components/suggested-issues.vue'
import SensitiveWords from './components/sensitive-words.vue'
import ChatMode from './components/chat-mode/index.vue'
import MarkdownSetting from './components/markdown-setting.vue'
import CommonProblem from './components/common-problem.vue'
import DisplayAitations from './components/display-aitations.vue'
import { saveRobot, getRobotList } from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'
import { useModelStore } from '@/stores/modules/model'
import ShowLike from './components/show-like.vue'
import rolePermission from './role-permission.vue'
import Skill from './components/skill/index.vue'
import ChatCache from './components/chat-cache.vue'
import AnswerPrompt from './components/answer-prompt.vue'

const robotStore = useRobotStore()
const activeLocalKey = '/robot/config/basic-config/activeKey'

const scrollBox = ref(null)

const { robotInfo } = storeToRefs(robotStore)
const { getRobot } = robotStore

const modelStore = useModelStore()

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
  if (val.robot_avatar) {
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
})

const saveForm = () => {
  // 对机器人头像特殊处理
  let robot_avatar
  if (formState.robot_avatar) {
    robot_avatar = new File([formState.robot_avatar], formState.robot_avatar.name)
  }
  let formData = JSON.parse(JSON.stringify(toRaw(formState)))
  // 有机器人头像就赋值
  if (robot_avatar) {
    formData.robot_avatar = robot_avatar
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

  delete formData.robot_avatar_url

  saveLoading.value = true

  saveRobot(formData)
    .then((res) => {
      if (res.res != 0) {
        return message.error(res.msg)
      }

      saveLoading.value = false

      message.success('保存成功')

      getRobot(formState.id)
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
