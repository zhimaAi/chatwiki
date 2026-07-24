<style lang="less" scoped>
.workflow-page {
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  overflow: hidden;
  position: relative;

  .page-body {
    flex: 1;
    display: flex;
    flex-flow: row nowrap;
    overflow: hidden;
    background: #f0f2f5;

    .page-left {
      height: 100%;
      padding: 8px;
    }

    .page-container {
      flex: 1;
      height: 100%;
    }
  }
}

.logic-flow-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>

<template>
  <div class="workflow-page">
    <PageHeader
      ref="pageHeaderRef"
      @edit="handleClickEdit"
      @save="handleSave"
      @release="handleRelease"
      @getGlobal="getGlobal"
      @getVersionRecord="getVersionRecord"
      :lf="lf"
      :currentVersion="currentVersionData"
      :start_node_params="start_node_params"
      :saveLoading="saveLoading"
      :isEditing="isEditing"
      :isLockedByOther="isLockedByOther"
      :lockRemoteAddr="lockRemoteAddr"
      :lockUserAgent="lockUserAgent"
      :loginUserName="loginUserName"
      :autoSaveEnabled="autoSaveEnabled"
    />
    <div class="page-body">
      <div class="page-left">
        <PageSidebar />
      </div>
      <div class="page-container">
        <WorkflowCanvas
          ref="workflowCanvasRef"
          @selectedNode="handleSelectedNode"
          @onDeleteNode="onDeleteNode"
          @runTest="handleRunTest"
        />
      </div>
    </div>
    <AgentAbnormalModal ref="agentAbnormalModalRef" />
    <AddRobotAlert ref="addRobotAlertRef" />
    <VersionModel ref="versionModelRef" @handleOpenErrorNode="handleOpenErrorNode" />
    <PublishDetail
      ref="publishDetailRef"
      @preview="handlePreviewVersion"
      @setVersion="setVersion"
      :isLockedByOther="isLockedByOther"
    />
  </div>
</template>

<script setup>
import { useWorkflowStore } from '@/stores/modules/workflow'
import {
  getNodeList,
  saveNodes as saveNodesRequest,
  getDraftKey,
  heartbeatDraftKey,
  releaseDraftKey,
  workFlowNextVersion,
  getRobotList
} from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'
import { useUserStore } from '@/stores/modules/user'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { generateUniqueId, duplicateRemoval, removeRepeat } from '@/utils/index'
import { onMounted, ref, onUnmounted, watch, computed, h, provide} from 'vue'
import { onBeforeRouteLeave, useRoute } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import PageSidebar from './components/page-sidebar.vue'
import WorkflowCanvas from './components/workflow-canvas.vue'
import PageHeader from './components/page-header.vue'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'
import { getNodeTypes } from './components/node-list'
import VersionModel from './components/version-model.vue'
import PublishDetail from './components/publish-detail.vue'
import AgentAbnormalModal from './components/agent-abnormal-modal.vue'
import { getModelConfigOption } from '@/api/model/index'
import { getModelOptionsList } from '@/components/model-select/index.js'
import { useModelStore } from '@/stores/modules/model'
import {downloadPlugin, openPlugin} from "@/api/plugins/index.js"
import { useI18n } from '@/hooks/web/useI18n'
import { PATH_URL } from '@/utils/http/axios/service'
import {
  WORKFLOW_AUTO_SAVE_MS,
  createWorkflowHeartbeatController,
  createWorkflowLeaseToken,
  getWorkflowWindowIdentifier
} from './workflow-edit-lock'

const { t } = useI18n('views.workflow.index')

const modelStore = useModelStore()

const route = useRoute()
const query = route.query
const robot_key = ref(route.query.robot_key)

// sessionStorage 按窗口隔离；lease_token 按页面租约实例隔离。
const workflowWindowIdentifier = getWorkflowWindowIdentifier()
const workflowLeaseToken = ref(createWorkflowLeaseToken())
const getUniIdentifier = () => workflowWindowIdentifier
const getLeaseToken = () => workflowLeaseToken.value
const buildUserAgent = () => {
  try {
    const ua = navigator.userAgent || ''
    const platform = navigator.platform || ''
    let os = 'Unknown'
    if (/Windows/i.test(ua)) os = 'Windows'
    else if (/Macintosh|Mac OS X/i.test(ua)) os = 'MacOS'
    else if (/Linux/i.test(ua)) os = 'Linux'
    else if (/Android/i.test(ua)) os = 'Android'
    else if (/iPhone|iPad|iPod/i.test(ua)) os = 'iOS'
    let browser = 'Unknown'
    let version = ''
    const m = ua.match(/Edg\/([\d\.]+)/) || ua.match(/Chrome\/([\d\.]+)/) || ua.match(/Firefox\/([\d\.]+)/) || ua.match(/Version\/([\d\.]+).*Safari/)
    if (m) {
      if (ua.includes('Edg/')) browser = 'Edge'
      else if (ua.includes('Chrome/')) browser = 'Chrome'
      else if (ua.includes('Firefox/')) browser = 'Firefox'
      else if (ua.includes('Safari') && !ua.includes('Chrome')) browser = 'Safari'
      version = m[1]
    }
    return `platform=${platform}; os=${os}; browser=${browser}/${version}; ua=${ua}`
  } catch (e) {
    return 'ua=unknown'
  }
}

const pageUserAgent = buildUserAgent()

const addRobotAlertRef = ref(null)
const agentAbnormalModalRef = ref(null)
const workflowCanvasRef = ref(null)

const workflowStore = useWorkflowStore()
const robotStore = useRobotStore()
const lf = computed(() => {
  return workflowCanvasRef.value?.lfRef
})

// 触发器列表
const triggerList = computed(() => workflowStore.triggerList)

const currentVersion = ref('')
const currentVersionData = ref('')
const pageHeaderRef = ref(null)

const loop_save_canvas_status = computed(()=>{
  return robotStore.robotInfo.loop_save_canvas_status
})

watch(loop_save_canvas_status, (val) => {
  // 在循环节点里面打开 运行测试 保存一下草稿
  if(val > 0){
    handleSave('automatic')
  }
})

// const nodes = []
const nodeTypes = getNodeTypes()

// 协同编辑与自动保存控制
const isEditing = ref(true) // 是否处于可编辑状态
const isLockedByOther = ref(false) // 是否被其他人编辑锁定
const autoSaveEnabled = ref(true) // 是否允许自动保存
let inactivityTimer = null // 无改动计时器（5分钟）
let changeMonitorTimer = null // 变更监控定时器
let lastChangeTs = Date.now() // 最近一次检测到数据变更的时间戳
let lastChangeHash = '' // 最近一次快照hash
const INACTIVITY_MS = 5 * 60 * 1000
// 他人持锁信息
const lockRemoteAddr = ref('')
const lockUserAgent = ref('')
const loginUserName = ref('')
const ownsEditLock = ref(false)
const isPageVisible = ref(document.visibilityState === 'visible')
const userStore = useUserStore()
const localeStore = useLocaleStoreWithOut()
let releaseStarted = false
let visibleLockAttempted = false
let visibleAcquirePromise = null
let draftSaveTimeBeforeHidden = 0

const getWorkflowLockPayload = (targetRobotKey = robot_key.value, targetLeaseToken = getLeaseToken()) => ({
  robot_key: targetRobotKey,
  uni_identifier: getUniIdentifier(),
  lease_token: targetLeaseToken,
  user_agent: pageUserAgent
})

provide('getWorkflowLockPayload', getWorkflowLockPayload)

function applyLockConflict (data = {}) {
  ownsEditLock.value = false
  isLockedByOther.value = true
  isEditing.value = false
  autoSaveEnabled.value = false
  lockRemoteAddr.value = data.remote_addr || ''
  lockUserAgent.value = data.user_agent || ''
  loginUserName.value = data.login_user_name || ''
  robotStore.setIsLockedByOther(true)
  updateAutoSaveTimer()
}

function handleLockLost (data = {}) {
  heartbeatController.stop()
  applyLockConflict(data)
  message.warning(t('msg_switched_to_view_mode'))
}

const heartbeatController = createWorkflowHeartbeatController({
  sendHeartbeat: () => heartbeatDraftKey(getWorkflowLockPayload()),
  onLockLost: handleLockLost
})

async function saveNodes (payload) {
  if (!isPageVisible.value && payload?.draft_save_type === 'automatic') return null
  try {
    return await saveNodesRequest({ ...payload, ...getWorkflowLockPayload() })
  } catch (error) {
    if (error?.data?.lock_conflict == 1) {
      handleLockLost(error.data)
    }
    throw error
  }
}

provide('handleWorkflowLockLost', handleLockLost)

function visibilityHandler () {
  if (document.visibilityState === 'hidden') {
    isPageVisible.value = false
    visibleLockAttempted = false
    draftSaveTimeBeforeHidden = +robotStore.robotInfo.draft_save_time || 0
    heartbeatController.stop()
    updateAutoSaveTimer()
    return
  }

  isPageVisible.value = true
  void tryAcquireVisibleLock()
}

function sendReleaseKeepalive (payload) {
  const apiBase = String(PATH_URL || '').replace(/\/$/, '')
  const body = new URLSearchParams(payload)
  void fetch(`${apiBase}/manage/releaseDraftKey`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
      'X-Requested-With': 'XMLHttpRequest',
      'App-Type': '',
      lang: localeStore.getCurrentLocale.lang,
      token: userStore.getToken || ''
    },
    body,
    keepalive: true,
    credentials: 'same-origin'
  }).catch(() => {})
}

function computeCanvasHash () {
  try {
    const data = getCanvasData()
    return JSON.stringify(data)
  } catch (e) {
    return ''
  }
}

function startChangeMonitor () {
  stopChangeMonitor()
  // 记录进入编辑后的初始快照
  lastChangeHash = computeCanvasHash()
  lastChangeTs = Date.now()
  changeMonitorTimer = setInterval(() => {
    if (!isPageVisible.value || !isEditing.value) return
    const hash = computeCanvasHash()
    if (hash && hash !== lastChangeHash) {
      lastChangeHash = hash
      lastChangeTs = Date.now()
      // 内容发生变化后恢复自动保存；编辑权只由服务端锁决定。
      ensureAutoSaveOnChange()
    }
  }, 2 * 1000)
}

function stopChangeMonitor () {
  changeMonitorTimer && clearInterval(changeMonitorTimer)
  changeMonitorTimer = null
}

function ensureAutoSaveOnChange () {
  if (!isPageVisible.value || !isEditing.value || !ownsEditLock.value) return
  if (!autoSaveEnabled.value) {
    autoSaveEnabled.value = true
  }
  updateAutoSaveTimer()
}

function startInactivityWatcher () {
  stopInactivityWatcher()
  inactivityTimer = setInterval(() => {
    if (isPageVisible.value && isEditing.value && autoSaveEnabled.value) {
      const idle = Date.now() - lastChangeTs
      if (idle >= INACTIVITY_MS) {
        autoSaveEnabled.value = false
        message.warning(t('msg_auto_save_paused'))
        updateAutoSaveTimer()
      }
    }
  }, 30 * 1000)
}

function stopInactivityWatcher () {
  inactivityTimer && clearInterval(inactivityTimer)
  inactivityTimer = null
}

async function acquireEditLock ({ leaseToken = getLeaseToken(), activate = true } = {}) {
  try {
    const res = await getDraftKey(getWorkflowLockPayload(robot_key.value, leaseToken))
    const data = res?.data || {}
    if (data.lock_res) {
      workflowLeaseToken.value = leaseToken
      ownsEditLock.value = true
      releaseStarted = false
      isLockedByOther.value = false
      isEditing.value = activate
      autoSaveEnabled.value = activate
      lockRemoteAddr.value = ''
      lockUserAgent.value = ''
      loginUserName.value = ''
      robotStore.setIsLockedByOther(false)
      if (activate) {
        heartbeatController.start(data)
      }
      updateAutoSaveTimer()
      return data
    }
    applyLockConflict(data)
    return null
  } catch (e) {
    // 获取锁失败必须保持只读，不能让无法确认所有权的窗口继续编辑。
    applyLockConflict(e?.data || {})
    return null
  }
}

async function tryAcquireVisibleLock () {
  if (!isPageVisible.value) return null
  if (visibleAcquirePromise) {
    visibleLockAttempted = true
    return visibleAcquirePromise
  }
  if (visibleLockAttempted) return visibleAcquirePromise

  visibleLockAttempted = true
  const acquirePromise = (async () => {
    const knownDraftSaveTime = draftSaveTimeBeforeHidden
    const nextLeaseToken = createWorkflowLeaseToken()
    isEditing.value = false
    autoSaveEnabled.value = false
    heartbeatController.stop()
    updateAutoSaveTimer()

    const lockData = await acquireEditLock({ leaseToken: nextLeaseToken, activate: false })
    if (!lockData) {
      message.warning(t('msg_switched_to_view_mode'))
      return null
    }
    if (!isPageVisible.value) {
      return lockData
    }

    heartbeatController.start(lockData)
    try {
      const robotRes = await robotStore.getRobot(route.query.id)
      if (!robotRes) {
        throw new Error('refresh workflow draft time failed')
      }
      const latestDraftSaveTime = +robotStore.robotInfo.draft_save_time || 0
      if (latestDraftSaveTime > knownDraftSaveTime) {
        await loadDraftNodes()
      }
      if (!isPageVisible.value) {
        return lockData
      }

      isEditing.value = true
      autoSaveEnabled.value = true
      robotStore.setIsLockedByOther(false)
      lastChangeHash = computeCanvasHash()
      lastChangeTs = Date.now()
      updateAutoSaveTimer()
      return lockData
    } catch (e) {
      heartbeatController.stop()
      await releaseEditLock({ targetLeaseToken: nextLeaseToken })
      applyLockConflict({})
      message.warning(t('msg_switched_to_view_mode'))
      return null
    }
  })()

  visibleAcquirePromise = acquirePromise
  try {
    return await acquirePromise
  } finally {
    if (visibleAcquirePromise === acquirePromise) {
      visibleAcquirePromise = null
    }
  }
}

async function releaseEditLock ({ keepalive = false, targetRobotKey = robot_key.value, targetLeaseToken = getLeaseToken() } = {}) {
  if (!ownsEditLock.value || releaseStarted) return
  const payload = getWorkflowLockPayload(targetRobotKey, targetLeaseToken)
  releaseStarted = true
  ownsEditLock.value = false
  heartbeatController.stop()
  updateAutoSaveTimer()
  if (keepalive) {
    sendReleaseKeepalive(payload)
    return
  }
  try {
    await releaseDraftKey(payload)
  } catch (e) {
    // 主动释放失败时由 Redis TTL 兜底。
  }
}

function handlePageExit () {
  visibleLockAttempted = false
  draftSaveTimeBeforeHidden = +robotStore.robotInfo.draft_save_time || 0
  if (!ownsEditLock.value) return
  void releaseEditLock({ keepalive: true })
}

async function handlePageShow (event) {
  if (!event.persisted) return
  isPageVisible.value = true
  await tryAcquireVisibleLock()
}

function getNode (list) {
  list = list || []
  let nodes = []
  let edges = []

  list.forEach((item) => {
    // 边数据处理
    if (item.node_type == 0) {
      let edge = JSON.parse(item.node_info_json)

      if(edge.type == 'custom-edge'){
        edge.pointsList = [];
        edge.type = 'custom-bezier-edge'
      }

      edges.push(edge)
    } else {
      const type = nodeTypes[item.node_type]

      item.type = type

      // 节点数据处理
      let node = JSON.parse(item.node_info_json)

      if (item.node_type != 0) {
        node.type = nodeTypes[item.node_type]
        node.id = item.node_key || generateUniqueId(node.type)
        node.x = node.x || 0
        node.y = node.y || 0
      }

      if (item.node_type == -1) {
        item.node_params = node.dataRaw
      }

      if (!node.nodeSortKey) {
        item.nodeSortKey = node.id.substring(0, 8) + node.id.substring(node.id.length - 8)
      } else {
        item.nodeSortKey = node.nodeSortKey
      }

      item.dataRaw = node.dataRaw || item.node_params

      // 删除不要的参数
      delete item.node_info_json

      node.loop_parent_key = item.loop_parent_key // 父节点id

      // 设置 properties
      node.properties = item
      node.properties.width = node.width
      node.properties.height = node.height
      nodes.push(node)
    }
  })

  setWorkflowData({ nodes: nodes, edges: edges })
}

const toAddRobot = (val) => {
  // router.push({ name: 'addRobot' })
  addRobotAlertRef.value.open(val, true)
}

const getParsedNodeInfo = (item) => {
  if (!item?.node_info_json) return null
  if (typeof item.node_info_json === 'string') {
    try {
      return JSON.parse(item.node_info_json)
    } catch (e) {
      return null
    }
  }

  return item.node_info_json
}

const getAgentParams = (item, nodeInfo) => {
  const dataRaw = nodeInfo?.dataRaw || item?.node_params
  if (!dataRaw) return null
  if (typeof dataRaw === 'string') {
    try {
      return JSON.parse(dataRaw)
    } catch (e) {
      return null
    }
  }

  return dataRaw
}

const getAbnormalAgentNodes = (list) => {
  return (Array.isArray(list) ? list : []).reduce((nodes, item) => {
    if (item.node_type != 52) return nodes
    const nodeInfo = getParsedNodeInfo(item)
    const nodeParams = getAgentParams(item, nodeInfo)
    if (nodeParams?.agent?.robot_id == 0) {
      nodes.push({ item, nodeInfo, nodeParams })
    }
    return nodes
  }, [])
}

const checkAgentAbnormalNodes = async (list) => {
  const abnormalNodes = getAbnormalAgentNodes(list)
  if (!abnormalNodes.length) return list

  const res = await getRobotList({ application_type: 2 })
  const agentRobotOptions = (res.data || []).filter(item => item.has_published == 1)

  if (!agentRobotOptions.length) {
    message.warning('暂无可选择的Agent机器人')
    return list
  }

  const selectedRobot = await agentAbnormalModalRef.value?.open(agentRobotOptions)
  if (!selectedRobot) return list

  abnormalNodes.forEach(({ item, nodeInfo, nodeParams }) => {
    nodeParams.agent.robot_id = Number(selectedRobot.id)
    nodeParams.agent.robot_info = selectedRobot
    item.node_params = JSON.stringify(nodeParams)
    if (nodeInfo) {
      nodeInfo.dataRaw = item.node_params
      item.node_info_json = JSON.stringify(nodeInfo)
    }
    item.node_name = selectedRobot.robot_name || item.node_name
    item.node_icon = selectedRobot.robot_avatar || item.node_icon
  })

  return list
}

const setWorkflowData = (data) => {
  workflowCanvasRef.value.setData(data)
}

const getCanvasData = () => {
  let data = workflowCanvasRef.value.getData()

  let list = []
  let edgeMap = {}
  // 先处理边数据
  data.edges.forEach((item) => {
    let obj = {
      node_key: item.id,
      node_name: 'edge',
      node_type: 0
    }

    let node_info_json = {
      ...item
    }

    obj.node_info_json = node_info_json

    list.push(obj)
    if (item.sourceAnchorId) {
      edgeMap[item.sourceAnchorId] = item.targetNodeId
    }

    if (item.targetAnchorId) {
      edgeMap[item.targetAnchorId] = item.sourceNodeId
    }
  })

  data.nodes.forEach((item) => {
    let obj = {
      ...item.properties,
      node_type: +item.properties.node_type,
      node_key: item.id
    }

    obj.node_info_json = {
      type: item.type,
      x: item.x,
      y: item.y,
      width: item.properties.width,
      height: item.properties.height,
      id: item.id,
      nodeSortKey: obj.nodeSortKey,
      dataRaw: item.properties.node_params,
    }

    // 关联next_node_key
    obj.next_node_key = edgeMap[obj.nodeSortKey + '-anchor_right'] || ''
    obj.prev_node_key = edgeMap[obj.nodeSortKey + '-anchor_left'] || ''

    let node_params = JSON.parse(obj.node_params)

    // 开始节点
    if (obj.node_type == 1) {

    }

    if (obj.node_type == 2) {
      // 判断分支
      if (node_params.term && node_params.term.length > 0) {
        node_params.term.forEach((msg, index) => {
          let key = obj.nodeSortKey + '-anchor_' + index
          msg.next_node_key = edgeMap[key] || ''
        })
      }
    }

    if (obj.node_type == 3) {
      // 问题分类
      if (node_params.cate && node_params.cate.categorys && node_params.cate.categorys.length > 0) {
        node_params.cate.categorys.forEach((msg, index) => {
          let key = obj.nodeSortKey + '-anchor_' + index
          msg.next_node_key = edgeMap[key] || ''
        })
      }
    }

    if (obj.node_type == 43) {
      // 问答
      if(node_params.question.answer_type == 'menu'){
        if (node_params.question && node_params.question.reply_content_list && node_params.question.reply_content_list.length > 0) {
          let menu_content = node_params.question.reply_content_list[0].smart_menu?.menu_content || []
          menu_content.forEach((msg, index) => {
            let key = obj.nodeSortKey + '-anchor_' + index
            msg.next_node_key = edgeMap[key] || ''
          })
        }
      }

    }

    if (obj.node_type == 17) {
      // 代码运行
      let exception = edgeMap[obj.nodeSortKey + '-anchor_right_exception']
      if (exception) {
        node_params.code_run.exception = exception
      }
    }

    if (obj.node_type == 4) {
      // http请求
      let exception = edgeMap[obj.nodeSortKey + '-anchor_right_exception']
      if (exception) {
        node_params.curl.exception = exception
      }
    }

    if (obj.node_type == 41) {
      // 工作流
      let exception = edgeMap[obj.nodeSortKey + '-anchor_right_exception']
      if (exception) {
        node_params.workflow.exception = exception
      }
    }

    obj.node_params = node_params

    // 删除无用字段
    delete obj.dataRaw

    list.push(obj)
  })

  return list
}

const handleSave = async (type) => {
  if (!isPageVisible.value && type === 'automatic') return
  if (!ownsEditLock.value || isLockedByOther.value) {
    // 手动保存时给出提示，自动保存静默跳过
    if (type === 'handle') {
      message.warning(t('msg_save_locked_by_other'))
    }
    return
  }
  if (!isEditing.value) return
  if (!autoSaveEnabled.value && type === 'automatic') return

  try {
    await workflowCanvasRef.value.layoutAllGroups()
  } catch (e) {
    // 分组整理失败不阻塞保存
  }

  let list = getCanvasData()

  if (type === 'handle') {
    // 手动保存增强：传入服务端最近草稿时间戳与覆盖标识
    // await robotStore.getRobot(query.id)
    const clientDraftTs = +robotStore.robotInfo.draft_save_time || 0
    const basePayload = {
      robot_key: robot_key.value,
      data_type: 1,
      node_list: JSON.stringify(list),
      draft_save_type: 'handle',
      draft_save_time: clientDraftTs,
      re_cover_save: 0,
      uni_identifier: getUniIdentifier(),
      lease_token: getLeaseToken(),
      user_agent: pageUserAgent
    }

    const result = await confirmOverrideAndSave(basePayload, false)
  
    if (result.saved) {
      message.success(t('msg_save_success'))
    } else if (!result.behind) {
      // message.error('保存失败，请稍后重试')
    }
    return
  }

  // 自动保存沿用原逻辑，但在成功后刷新 robotInfo 的草稿时间戳
  try {
    await saveNodes({
      robot_key: robot_key.value,
      data_type: 1,
      node_list: JSON.stringify(list),
      draft_save_type: 'automatic',
      draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
      uni_identifier: getUniIdentifier(),
      lease_token: getLeaseToken(),
      user_agent: pageUserAgent
    })
    await robotStore.getRobot(query.id)
    const ts = +robotStore.robotInfo.draft_save_time || dayjs().unix()
    robotStore.setDrafSaveTime({
      draft_save_type: 'automatic',
      draft_save_time: ts
    })
  } catch (e) {
    // 自动保存失败不提示
  }
}

let timer = null
let autoBehindConfirming = false
// behind 弹窗节流与封装
let autoBehindPromptShown = false // 本会话内已提示过自动保存覆盖
let lastAutoBehindPromptTs = 0 // 上次提示时间戳
const AUTO_PROMPT_COOLDOWN_MS = 2 * 60 * 1000 // 自动保存覆盖弹窗最短提示间隔
const BEHIND_MODAL_TEXT = computed(() => t('msg_confirm_override_draft'))

async function confirmOverrideAndSave (basePayload, isAuto) {
  try {
    const payload = { ...basePayload, ...getWorkflowLockPayload() }
    const res = await saveNodes(payload)
    const behind = res?.data?.behind_draft == 1 || res?.behind_draft == 1
    if (!behind) {
      // 正常保存成功后刷新服务端草稿时间戳
      await robotStore.getRobot(query.id)
      const ts = +robotStore.robotInfo.draft_save_time || dayjs().unix()
      robotStore.setDrafSaveTime({
        draft_save_type: isAuto ? 'automatic' : 'handle',
        draft_save_time: ts
      })
      return { saved: true, behind: false }
    }

    // 版本落后：自动保存节流，手动保存每次提示
    if (isAuto) {
      const now = Date.now()
      if (autoBehindConfirming || (autoBehindPromptShown && (now - lastAutoBehindPromptTs < AUTO_PROMPT_COOLDOWN_MS))) {
        // 自动保存在提示后的一段时间内不再弹窗，减轻干扰；暂停自动保存
        autoSaveEnabled.value = false
        updateAutoSaveTimer()
        return { saved: false, behind: true }
      }
      autoBehindConfirming = true
      autoBehindPromptShown = true
      lastAutoBehindPromptTs = now
    }

    return await new Promise((resolve) => {
      Modal.confirm({
        title: t('title_prompt'),
        icon: null,
        content: BEHIND_MODAL_TEXT.value,
        okText: t('btn_confirm'),
        cancelText: t('btn_cancel'),
        async onOk () {
          const forcePayload = { ...basePayload, re_cover_save: 1, ...getWorkflowLockPayload() }
          const res2 = await saveNodes(forcePayload)
          if (res2 && res2.res == 0 || !res2?.data?.behind_draft) {
            await robotStore.getRobot(query.id)
            const ts = +robotStore.robotInfo.draft_save_time || dayjs().unix()
            robotStore.setDrafSaveTime({
              draft_save_type: isAuto ? 'automatic' : 'handle',
              draft_save_time: ts
            })
          }
          if (isAuto) {
            autoSaveEnabled.value = true
            updateAutoSaveTimer()
            autoBehindConfirming = false
            autoBehindPromptShown = false // 覆盖后重置，允许未来再次提示
            lastAutoBehindPromptTs = Date.now()
          }
          resolve({ saved: true, behind: true })
        },
        onCancel () {
          if (isAuto) {
            autoSaveEnabled.value = false // 取消后暂停自动保存，避免频繁提示
            updateAutoSaveTimer()
            autoBehindConfirming = false
            // 保留 autoBehindPromptShown=true；在冷却时间后允许再次提醒
            lastAutoBehindPromptTs = Date.now()
          }
          resolve({ saved: false, behind: true })
        }
      })
    })
  } catch (e) {
    // 若后端以异常形式返回 behind_draft，可视作落后处理；否则静默失败
    const behind = e?.data?.behind_draft == 1 || e?.behind_draft == 1
    if (!behind) return { saved: false, behind: false }
    if (isAuto) {
      const now = Date.now()
      if (autoBehindConfirming || (autoBehindPromptShown && (now - lastAutoBehindPromptTs < AUTO_PROMPT_COOLDOWN_MS))) {
        autoSaveEnabled.value = false
        updateAutoSaveTimer()
        return { saved: false, behind: true }
      }
      autoBehindConfirming = true
      autoBehindPromptShown = true
      lastAutoBehindPromptTs = now
    }
    return await new Promise((resolve) => {
      Modal.confirm({
        title: t('title_prompt'),
        icon: null,
        content: BEHIND_MODAL_TEXT.value,
        okText: t('btn_confirm'),
        cancelText: t('btn_cancel'),
        async onOk () {
          const forcePayload = { ...basePayload, re_cover_save: 1, ...getWorkflowLockPayload() }
          const res2 = await saveNodes(forcePayload)
          if (res2 && res2.res == 0 || !res2?.data?.behind_draft) {
            await robotStore.getRobot(query.id)
            const ts = +robotStore.robotInfo.draft_save_time || dayjs().unix()
            robotStore.setDrafSaveTime({
              draft_save_type: isAuto ? 'automatic' : 'handle',
              draft_save_time: ts
            })
          }
          if (isAuto) {
            autoSaveEnabled.value = true
            updateAutoSaveTimer()
            autoBehindConfirming = false
            autoBehindPromptShown = false
            lastAutoBehindPromptTs = Date.now()
          }
          resolve({ saved: true, behind: true })
        },
        onCancel () {
          if (isAuto) {
            autoSaveEnabled.value = false
            updateAutoSaveTimer()
            autoBehindConfirming = false
            lastAutoBehindPromptTs = Date.now()
          }
          resolve({ saved: false, behind: true })
        }
      })
    })
  }
}

// 自动保存（含版本落后弹窗确认）
async function handleAutoSaveWithConflictCheck () {
  if (!isPageVisible.value) return
  if (!ownsEditLock.value || isLockedByOther.value) return
  if (!isEditing.value) return
  if (!autoSaveEnabled.value) return

  const list = getCanvasData()
  const clientDraftTs = +robotStore.robotInfo.draft_save_time || 0
  const basePayload = {
    robot_key: robot_key.value,
    data_type: 1,
    node_list: JSON.stringify(list),
    draft_save_type: 'automatic',
    draft_save_time: clientDraftTs,
    re_cover_save: 0,
    uni_identifier: getUniIdentifier(),
    lease_token: getLeaseToken(),
    user_agent: pageUserAgent
  }
  await confirmOverrideAndSave(basePayload, true)
}
function updateAutoSaveTimer () {
  timer && clearInterval(timer)
  timer = null
  if (import.meta.env.PROD && isPageVisible.value && ownsEditLock.value && isEditing.value && autoSaveEnabled.value && !isLockedByOther.value) {
    timer = setInterval(() => {
      handleAutoSaveWithConflictCheck()
    }, WORKFLOW_AUTO_SAVE_MS)
  }
}

onUnmounted(() => {
  window.removeEventListener('pagehide', handlePageExit)
  window.removeEventListener('pageshow', handlePageShow)
  document.removeEventListener('visibilitychange', visibilityHandler)
  timer && clearInterval(timer)
  stopInactivityWatcher()
  stopChangeMonitor()
  heartbeatController.stop()
  void releaseEditLock()
})

onBeforeRouteLeave(async () => {
  await releaseEditLock()
})


const versionModelRef = ref(null)

const saveLoading = ref(false)

const openVersionModel = (node_list) => {
  versionModelRef.value.show(node_list)
}
// 发布机器人
const handleRelease = async () => {
  if (!ownsEditLock.value || isLockedByOther.value) {
    message.warning(t('msg_publish_locked_by_other'))
    return
  }
  let list = getCanvasData()

  let errorNodes = []
  for (let i = 0; i < list.length; i++) {
    let node = list[i]
    // 跳过边节点
    if (node.node_type == 0 || node.node_type == -1 || node.node_type == 43) {
      // 跳过
      continue
    }
    if(!node.loop_parent_key){
      // 分组里面的不用校验
      if (node.node_type == 1) {
        if (node.next_node_key == '') {
          errorNodes.push(node)
        }
      } else if (node.node_type == 7 || node.node_type == 3) {
        if (node.prev_node_key == '') {
          errorNodes.push(node)
        }
      } else {
        if (node.next_node_key == '' || node.prev_node_key == '') {
          errorNodes.push(node)
        }
      }

    }

  }

  if (errorNodes.length > 0) {
    // message.error(t('msg_unconnected_nodes'))
    // return
  }

  // 先保存草稿，再发布
  // saveLoading.value = true;
  try {
    // message.loading('保存中...')
    await saveNodes({
      robot_key: robot_key.value,
      data_type: 1,
      node_list: JSON.stringify(list),
      draft_save_type: 'handle',
      draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
      uni_identifier: getUniIdentifier(),
      lease_token: getLeaseToken(),
      user_agent: pageUserAgent
    })
    // 保存后刷新服务端草稿时间戳
    await robotStore.getRobot(query.id)
    robotStore.setDrafSaveTime({
      draft_save_type: 'handle',
      draft_save_time: +robotStore.robotInfo.draft_save_time || 0
    })
    openVersionModel(JSON.stringify(list))
    // const res = await saveNodes({
    //   robot_key: robot_key.value,
    //   data_type: 2,
    //   node_list: JSON.stringify(list)
    // })
    // setTimeout(()=>{
    //   message.destroy()
    //   saveLoading.value = false;
    //   if (res && res.res == 0) {
    //     message.success('发布成功')
    //   }
    // },400)
  } catch (e) {
    saveLoading.value = false;
    // message.success('发布失败，请重试')
  }
}

// 选择节点
let selectedNode = ref(null)
const handleSelectedNode = (data) => {
  selectedNode.value = data
  // 结束节点不支持编辑
  if (data.properties.node_sub_type == 51) {
    return
  }

  // 内容发生变更（节点选择通常伴随编辑行为），尝试恢复自动保存并设置为领导者
  ensureAutoSaveOnChange()

}

// 删除节点
const onDeleteNode = () => {
  // 删除节点属于内容变更，立即恢复自动保存并设置为领导者
  ensureAutoSaveOnChange()
}

const getModelList = async () => {
  return getModelConfigOption({
    model_type: 'LLM'
  }).then((res) => {
    let list = res.data || []
    let { newList  } = getModelOptionsList(list)
    robotStore.setModelList(newList)
  })
}

const start_node_params = ref({
  diy_global: [],
  sys_global: [],
})
const getGlobal = () => {
  let list = getCanvasData()
  let start_node = list.filter(item => item.node_type == 1)
  if (start_node.length > 0) {
    start_node_params.value = start_node[0].node_params.start
  }
}

const publishDetailRef = ref(null)
const getVersionRecord = () => {
  // 历史记录打开时，禁止自动保存
  autoSaveEnabled.value = false
  updateAutoSaveTimer()
  publishDetailRef.value.showDrawer()
}

const setVersion = (data) => {
  if (!ownsEditLock.value || isLockedByOther.value) {
    message.warning(t('msg_restore_locked_by_other'))
    return
  }
  // 版本恢复前，先保存当前草稿一次
  const currentDraft = getCanvasData()
  saveNodes({
    robot_key: robot_key.value,
    data_type: 1,
    node_list: JSON.stringify(currentDraft),
    draft_save_type: 'handle',
    draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
    uni_identifier: getUniIdentifier(),
    lease_token: getLeaseToken(),
    user_agent: pageUserAgent
  }).then(() => {
    // 刷新草稿时间戳
    robotStore.getRobot(query.id).then(() => {
      robotStore.setDrafSaveTime({
        draft_save_type: 'handle',
        draft_save_time: +robotStore.robotInfo.draft_save_time || 0
      })
    })
    // 设置为当前的版本（切换到恢复的版本数据）
    currentVersion.value = ''
    currentVersionData.value = null
    getNode(data)
    message.success(t('msg_draft_saved_and_restored'))
    // 恢复后重新启用自动保存（回到编辑态）
    isEditing.value = true
    autoSaveEnabled.value = true
    updateAutoSaveTimer()
  })
}

const handlePreviewVersion = async (data, version) => {
  // 在发布详情中切换/预览前，先保存当前草稿，避免丢失
  if (ownsEditLock.value && !isLockedByOther.value) {
    const currentDraft = getCanvasData()
    try {
      await saveNodes({
        robot_key: robot_key.value,
        data_type: 1,
        node_list: JSON.stringify(currentDraft),
        draft_save_type: 'automatic',
        draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
        uni_identifier: getUniIdentifier(),
        lease_token: getLeaseToken(),
        user_agent: pageUserAgent
      })
      // 更新本地草稿时间，保持发布详情“最近保存于”显示
      await robotStore.getRobot(query.id)
      robotStore.setDrafSaveTime({
        draft_save_type: 'automatic',
        draft_save_time: +robotStore.robotInfo.draft_save_time || 0
      })
    } catch (e) {
      // 失败不阻断切换，但建议后端容错
    }
  }
  // currentVersion.value = version.version_id || ''
  currentVersionData.value = version
  clearInterval(timer)
  getNode(data)
}
// 运行测试
const handleRunTest = async () => {
  if (!ownsEditLock.value || isLockedByOther.value) {
    message.warning(t('msg_save_locked_by_other'))
    return
  }
  if (!isEditing.value) return

  let list = getCanvasData()

  await saveNodes({
    robot_key: robot_key.value,
    data_type: 1,
    node_list: JSON.stringify(list),
    draft_save_type: 'automatic',
    draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
    uni_identifier: getUniIdentifier(),
    lease_token: getLeaseToken(),
    user_agent: pageUserAgent
  })
  await robotStore.getRobot(query.id)
  workFlowNextVersion({
    robot_key: robot_key.value,
  }).then(res => {
    pageHeaderRef.value.openRunTest()
  }).catch((res) => {
    if(res.data && res.data.err_node_key){
      handleOpenErrorNode(res.data.err_node_key)
    }
  })

}

// 头部“编辑”按钮
const handleClickEdit = async () => {
  if (isLockedByOther.value) {
    message.warning(t('msg_edit_locked_by_other'))
    return
  }
  const ok = await acquireEditLock()
  if (ok) {
    isEditing.value = true
    autoSaveEnabled.value = true
    // 重新记录快照并启动无改动监控
    lastChangeHash = computeCanvasHash()
    lastChangeTs = Date.now()
    startChangeMonitor()
    updateAutoSaveTimer()
    message.success(t('msg_enter_edit_mode'))
    toAddRobot(1)
  }
}

const checkNodePluginStatus = (nds) => {
  nds = Array.isArray(nds) ? nds : []
  let closeNds = []
  let updateNds = []
  nds.forEach(nd => {
    if (nd.node_type == 21) {
      if (nd.has_loaded === "" || nd.plugin_version !== nd.remote_plugin_version) {
        updateNds.push(nd)
      } else if (nd.has_loaded === "false") {
        closeNds.push(nd)
      }
    }
  })
  if (closeNds.length || updateNds.length) {
    let names = [...new Set([...closeNds.map(i => i.node_name), ...updateNds.map(i => i.node_name)])]
    const titleStyle = {
      'font-weight': 500,
      'margin-top': '8px',
      'margin-bottom': '4px',
    }
    const nameStyle = {'color': '#8c8c8c'}
    Modal.confirm({
      title: t('title_workflow_fix'),
      content: h('div', {}, [
        h('div', {style: titleStyle}, t('msg_plugin_abnormal_prompt')),
        h('div', {style: nameStyle}, names.join('、'))
      ]),
      onOk: async () => {
        const tasks = []
        if (closeNds.length) tasks.push(openPlugin({name: closeNds.map(i => i.plugin_name).toString()}))
        if (updateNds.length) {
          tasks.push(downloadPlugin({
            download_data: JSON.stringify(updateNds.map(i => ({
              url: i.latest_version_detail_url,
              version_id: i.latest_version_detail_id
            })))
          }))
        }
        await Promise.all(tasks)
        message.success(t('msg_operation_reload'))
        setTimeout(() => {
          window.location.reload()
        }, 1200)
      }
    })
  }
}

const loadDraftNodes = async ({ checkPluginStatus = false } = {}) => {
  const res = await getNodeList({
    robot_key: robot_key.value,
    data_type: 1
  })

  const nodeList = await checkAgentAbnormalNodes(res.data)
  getNode(nodeList)
  if (checkPluginStatus) {
    checkNodePluginStatus(nodeList)
  }
}

const init = async () => {
  await getModelList()
  await workflowStore.getTriggerList(robot_key.value);
  workflowStore.getTriggerOfficialMsg(robot_key.value)
  await modelStore.getAllmodelList()
  workflowStore.getAllLibraryList();

  await loadDraftNodes({ checkPluginStatus: true })
}

const handleAutoSaveDraft = async (type = 'automatic') => {
  await handleSave(type)
}

provide('handleAutoSaveDraft', handleAutoSaveDraft)

const handleOpenErrorNode = (data) => {
  workflowCanvasRef.value.focusOnNode(data)
}

onMounted(async () => {
  window.addEventListener('pagehide', handlePageExit)
  window.addEventListener('pageshow', handlePageShow)
  document.addEventListener('visibilitychange', visibilityHandler)

  await init()

  // 初次进入仅在页面可见时抢锁；后续每次回到前台再尝试一次。
  if (isPageVisible.value) {
    visibleLockAttempted = true
    const acquired = await acquireEditLock()
    if (!acquired) {
      message.warning(t('msg_switched_to_view_mode'))
    }
  } else {
    isEditing.value = false
    autoSaveEnabled.value = false
  }
  updateAutoSaveTimer()

  // 复用当前组件切换工作流时，先释放旧锁，再用新租约抢占新工作流。
  watch(() => route.query.robot_key, async (newKey) => {
    if (!newKey || newKey === robot_key.value) return
    const oldRobotKey = robot_key.value
    const oldLeaseToken = getLeaseToken()
    await releaseEditLock({ targetRobotKey: oldRobotKey, targetLeaseToken: oldLeaseToken })
    robot_key.value = newKey
    workflowLeaseToken.value = createWorkflowLeaseToken()
    releaseStarted = false
    applyLockConflict({})
    await init()
    visibleLockAttempted = isPageVisible.value
    if (isPageVisible.value) {
      const nextAcquired = await acquireEditLock()
      if (!nextAcquired) {
        message.warning(t('msg_switched_to_view_mode'))
      }
    }
  })

  // 活动监听与5分钟无操作停止
  startInactivityWatcher()
  startChangeMonitor()
  if (route.query.show_tips) {
    message.info(t('msg_canvas_scroll_tips'), 6)
  }
})

</script>
