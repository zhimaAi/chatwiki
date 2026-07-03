<template>
  <div class="clawbot-settings-page">
    <aside class="settings-sidebar">
      <div class="settings-sidebar-title">{{ t('sidebar_title') }}</div>
      <div class="settings-menu">
        <div
          v-for="item in menuItems"
          :key="item.key"
          class="settings-menu-item"
          :class="{ active: activeMenuKey === item.key }"
          @click="handleMenuClick(item)"
        >
          <svg-icon v-if="item.iconName" class="menu-icon" :name="item.iconName"></svg-icon>
          <component v-else :is="item.icon" class="menu-icon" />
          <span class="menu-label">{{ item.label }}</span>
        </div>
      </div>
    </aside>

    <section class="settings-content">
      <div class="settings-content-title">{{ activeMenu?.title }}</div>
      <div class="settings-content-body">
        <div v-if="!isContextReady" class="settings-loading">
          <a-spin />
        </div>

        <template v-else>
          <div v-if="activeMenuKey === 'basic'" ref="scrollBox" class="basic-config-scroll">
            <div class="setting-box">
              <BasicConfig />
            </div>
            <div class="setting-box">
              <ModelSettings />
            </div>
            <div class="setting-box">
              <LanguageSetting />
            </div>
            <div class="setting-box">
              <DisplayAitations />
            </div>
            <div class="setting-box">
              <ShowLike />
            </div>
            <div class="setting-box">
              <VariableSetting />
            </div>
          </div>

          <component
            v-else-if="activePageComponent"
            :is="activePageComponent"
            :key="`${activeMenuKey}-${robotId}`"
          />
        </template>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, markRaw, nextTick, onMounted, provide, ref, toRaw, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { message } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  DownloadOutlined,
  ExclamationCircleOutlined,
  KeyOutlined,
  ProfileOutlined,
  QuestionCircleOutlined,
  SettingOutlined
} from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { saveRobot, getRobotList } from '@/api/robot/index'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useRobotStore } from '@/stores/modules/robot'
import { useModelStore } from '@/stores/modules/model'
import BasicConfig from '@/views/robot/robot-config/basic-config/components/basic-config.vue'
import ModelSettings from '@/views/robot/robot-config/basic-config/components/model-settings.vue'
import DisplayAitations from '@/views/robot/robot-config/basic-config/components/display-aitations.vue'
import ShowLike from '@/views/robot/robot-config/basic-config/components/show-like.vue'
import VariableSetting from '@/views/robot/robot-config/basic-config/components/variable-setting/index.vue'
import LanguageSetting from './components/language-setting.vue'
import E2bSettingPage from './components/e2b-setting.vue'
import QaFeedbackPage from '@/views/robot/robot-config/qa-feedback/index.vue'
import SessionRecordPage from '@/views/robot/robot-config/session-record/index.vue'
import ApiKeyManagePage from '@/views/robot/api-key-manage/index.vue'
import UnknownIssuePage from '@/views/robot/robot-config/unknown_issue/unknow-index.vue'
import ExportRecordPage from '@/views/robot/robot-config/export-record/index.vue'

const { t } = useI18n('views.clawbot.settings.index')

const route = useRoute()
const router = useRouter()
const clawbotStore = useClawbotStore()
const robotStore = useRobotStore()
const modelStore = useModelStore()
const { robotInfo } = storeToRefs(robotStore)

const scrollBox = ref(null)
const isRobotReady = ref(false)
const workFlowRobotList = ref([])

const menuItems = computed(() => ([
  { key: 'basic', label: t('menu_basic_config'), title: t('menu_basic_config'), icon: SettingOutlined },
  { key: 'qa-feedbacks', label: t('menu_qa_feedback'), title: t('menu_qa_feedback'), icon: ExclamationCircleOutlined, component: markRaw(QaFeedbackPage) },
  { key: 'session-record', label: t('menu_session_record'), title: t('menu_session_record'), icon: ClockCircleOutlined, component: markRaw(SessionRecordPage) },
  { key: 'api-key-manage', label: t('menu_api_key_manage'), title: t('menu_api_key_manage'), icon: KeyOutlined, component: markRaw(ApiKeyManagePage) },
  { key: 'unknown_issue', label: t('menu_unknown_issue'), title: t('menu_unknown_issue'), icon: QuestionCircleOutlined, component: markRaw(UnknownIssuePage) },
  { key: 'export-record', label: t('menu_export_record'), title: t('menu_export_record'), icon: DownloadOutlined, component: markRaw(ExportRecordPage) },
  { key: 'model-management', label: t('menu_model_management'), title: t('menu_model_management'), icon: ProfileOutlined },
  { key: 'e2b-settings', label: t('menu_e2b_settings'), title: t('menu_e2b_settings'), iconName: 'e2b-settings', component: markRaw(E2bSettingPage) }
]))

const robotId = computed(() => String(clawbotStore.currentAssistant?.id || ''))
const robotKey = computed(() => clawbotStore.currentAssistant?.robot_key || '')
const validMenuKeys = computed(() => menuItems.value.map((item) => item.key))
const lastExplicitMenuKey = ref('')

const activeMenuKey = computed(() => {
  const key = String(route.query.menu || '')
  if (validMenuKeys.value.includes(key) && key !== 'model-management') {
    return key
  }
  return lastExplicitMenuKey.value || 'basic'
})

const activeMenu = computed(() => {
  return menuItems.value.find((item) => item.key === activeMenuKey.value) || menuItems.value[0]
})

const activePageComponent = computed(() => activeMenu.value?.component)

const isContextReady = computed(() => {
  return (
    isRobotReady.value &&
    !!robotId.value &&
    String(route.query.id || '') === robotId.value &&
    String(route.query.robot_key || '') === String(robotKey.value)
  )
})

const isSameQuery = (nextQuery) => {
  const keys = new Set([...Object.keys(route.query), ...Object.keys(nextQuery)])
  for (const key of keys) {
    if (String(route.query[key] ?? '') !== String(nextQuery[key] ?? '')) {
      return false
    }
  }
  return true
}

const buildSettingsQuery = (query = {}) => {
  const nextQuery = {
    ...route.query,
    ...query,
    id: robotId.value,
    robot_key: robotKey.value
  }

  return nextQuery
}

const syncRouteQuery = () => {
  if (route.path !== '/clawbot/settings') return
  if (!robotId.value) return

  const menu = activeMenuKey.value
  const nextQuery = buildSettingsQuery({ menu })

  if (!isSameQuery(nextQuery)) {
    router.replace({
      path: '/clawbot/settings',
      query: nextQuery
    })
  }
}

const handleMenuClick = (item) => {
  if (item.key === 'model-management') {
    router.push('/user/model')
    return
  }

  router.push({
    path: '/clawbot/settings',
    query: buildSettingsQuery({
      menu: item.key,
    })
  })
}

const scrollBoxToBottom = () => {
  nextTick(() => {
    if (scrollBox.value) {
      scrollBox.value.scrollTop = scrollBox.value.scrollHeight
    }
  })
}

const getWorkFlowRobotList = () => {
  return getRobotList({ application_type: 1 }).then((res) => {
    workFlowRobotList.value = res.data || []
  })
}

const refreshRobot = async (id = robotId.value) => {
  if (!id) return false
  const res = await robotStore.getRobot(id)
  robotStore.fetchChatVariables()
  return res
}

const normalizeQuestion = (item) => {
  if (typeof item === 'string') return item
  return item?.content || ''
}

const stringifyQuestionConfig = (data = {}) => {
  return JSON.stringify({
    ...data,
    question: (data.question || []).map(normalizeQuestion)
  })
}

const updateRobotInfo = (val = {}) => {
  const newState = JSON.parse(JSON.stringify(val))
  if (val.robot_avatar && val.robot_avatar instanceof File) {
    newState.robot_avatar = new File([val.robot_avatar], val.robot_avatar.name)
  }

  if (!newState.op_type_relation_library && robotInfo.value.op_type_relation_library) {
    delete robotInfo.value.op_type_relation_library
  }

  Object.assign(robotInfo.value, newState)
  saveForm()
}

const refreshSavedRobot = async (id) => {
  await refreshRobot(id)
  await clawbotStore.fetchRobotInfo(id)
  await robotStore.getRobotLists()
}

const saveForm = () => {
  const rawState = toRaw(robotInfo.value)
  const formData = JSON.parse(JSON.stringify(rawState))

  if (rawState.robot_avatar) {
    if (rawState.robot_avatar instanceof File) {
      formData.robot_avatar = new File([rawState.robot_avatar], rawState.robot_avatar.name)
      delete formData.robot_avatar_url
    } else {
      formData.robot_avatar_url = rawState.robot_avatar
    }
  }

  formData.welcomes = stringifyQuestionConfig(formData.welcomes || { content: '', question: [] })
  formData.prompt_struct = JSON.stringify(formData.prompt_struct || {})
  formData.unknown_question_prompt = stringifyQuestionConfig(
    formData.unknown_question_prompt || { content: '', question: [] }
  )
  formData.cache_config = JSON.stringify(formData.cache_config || {})

  const applicationType = Number(rawState.application_type || 2)

  saveRobot(formData, applicationType)
    .then(async (res) => {
      if (res.res != 0) {
        return message.error(res.msg)
      }

      message.success(t('msg_save_success'))
      await refreshSavedRobot(rawState.id)
    })
    .catch(() => {})
}

let loadVersion = 0
const loadRobotContext = async (id) => {
  if (!id) {
    isRobotReady.value = false
    return
  }

  const currentVersion = ++loadVersion
  isRobotReady.value = false

  await refreshRobot(id)
  if (currentVersion !== loadVersion) return

  robotStore.getRobotLists()
  robotStore.getGroupList()
  modelStore.getAllmodelList()
  getWorkFlowRobotList()

  isRobotReady.value = true
}

provide('robotInfo', {
  robotInfo: robotInfo.value,
  updateRobotInfo,
  scrollBoxToBottom,
  getRobot: refreshRobot,
  getWorkFlowRobotList,
  robotList: workFlowRobotList
})

watch(
  () => route.query.menu,
  (menu) => {
    const key = String(menu || '')
    if (validMenuKeys.value.includes(key) && key !== 'model-management') {
      lastExplicitMenuKey.value = key
    }
  },
  { immediate: true }
)

watch(
  [robotId, robotKey, () => route.query.menu],
  () => {
    syncRouteQuery()
  },
  { immediate: true }
)

watch(
  robotId,
  (id) => {
    loadRobotContext(id)
  },
  { immediate: true }
)

onMounted(() => {
  syncRouteQuery()
})
</script>

<style lang="less" scoped>
.clawbot-settings-page {
  width: 100%;
  height: 100%;
  display: flex;
  overflow: hidden;
  background: #fff;
}

.settings-sidebar {
  width: 200px;
  height: 100%;
  flex-shrink: 0;
  background: #fff;
  border-right: 1px solid #e6ebf3;
}

.settings-sidebar-title {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  border-bottom: 1px solid #edf1f6;
  font-size: 16px;
  font-weight: 600;
  color: #101828;
}

.settings-menu {
  padding: 8px;
}

.settings-menu-item {
  height: 40px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: 8px;
  color: #344054;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.2s ease, background-color 0.2s ease;

  &:hover,
  &.active {
    color: #2475fc;
    background: #e6efff;
  }

  & + .settings-menu-item {
    margin-top: 4px;
  }

  .menu-icon {
    font-size: 16px;
    flex-shrink: 0;
  }

  .menu-label {
    min-width: 0;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}

.settings-content {
  flex: 1;
  min-width: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.settings-content-title {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  flex-shrink: 0;
  background: #fff;
  border-bottom: 1px solid #edf1f6;
  font-size: 16px;
  font-weight: 600;
  color: #101828;
}

.settings-content-body {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  background: #fff;
}

.settings-loading {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.basic-config-scroll {
  height: 100%;
  overflow-y: auto;
  padding: 16px 24px 24px;

  .setting-box {
    margin-bottom: 8px;
  }
}

</style>
