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
      :isLeader="isLeaderFlag"
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
    <AddRobotAlert ref="addRobotAlertRef" />
    <VersionModel ref="versionModelRef" />
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
import { getNodeList, saveNodes, getDraftKey } from '@/api/robot/index'
import { useRobotStore } from '@/stores/modules/robot'
import { generateUniqueId, duplicateRemoval, removeRepeat } from '@/utils/index'
import { onMounted, ref, onUnmounted, watch, computed, h, provide} from 'vue'
import { useRoute } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import PageSidebar from './components/page-sidebar.vue'
import WorkflowCanvas from './components/workflow-canvas.vue'
import PageHeader from './components/page-header.vue'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'
import { getNodeTypes } from './components/node-list'
import VersionModel from './components/version-model.vue'
import PublishDetail from './components/publish-detail.vue'
import { getModelConfigOption } from '@/api/model/index'
import { getModelOptionsList } from '@/components/model-select/index.js'
import { useModelStore } from '@/stores/modules/model'
import {downloadPlugin, openPlugin} from "@/api/plugins/index.js";

const modelStore = useModelStore()

const route = useRoute()
const query = route.query
const robot_key = ref(route.query.robot_key)

// 唯一标识与 user_agent 组装
const UNI_STORAGE_KEY = 'wf_uni_identifier'
const getUniIdentifier = () => {
  try {
    let id = localStorage.getItem(UNI_STORAGE_KEY)
    if (!id) {
      id = `${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
      localStorage.setItem(UNI_STORAGE_KEY, id)
    }
    return id
  } catch (e) {
    return `${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
  }
}
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

const addRobotAlertRef = ref(null)
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
// 领导者状态用于 UI 提示
const isLeaderFlag = ref(false)
// 跨标签页自动保存协调（仅允许一个标签执行自动保存）
const TAB_ID = `${Date.now()}_${Math.random()}`
const LEADER_TTL_MS = 60 * 1000
const HEARTBEAT_MS = 20 * 1000
let heartbeatTimer = null
const LEADER_PREFIX = 'wf_autosave_leader_'
const getLeaderKey = () => `${LEADER_PREFIX}${robot_key.value}`
const getOpenTabsKey = () => `wf_open_tabs_${robot_key.value}`

function incrementOpenTabs () {
  try {
    const n = parseInt(localStorage.getItem(getOpenTabsKey()) || '0', 10) || 0
    localStorage.setItem(getOpenTabsKey(), String(n + 1))
  } catch (e) {
    console.warn('incrementOpenTabs error', e)
  }
}
function decrementOpenTabs () {
  try {
    const n = parseInt(localStorage.getItem(getOpenTabsKey()) || '0', 10) || 0
    const next = Math.max(0, n - 1)
    if (next <= 0) {
      setTimeout(() => {
        localStorage.removeItem(getOpenTabsKey())
        // 所有同工作流页面关闭时，强制释放领导者键
        localStorage.removeItem(getLeaderKey())
      }, 0)
    } else {
      localStorage.setItem(getOpenTabsKey(), String(next))
    }
  } catch (e) {
    console.warn('decrementOpenTabs error', e)
   }
}

function cleanupExpiredLeaderKeys () {
  try {
    const now = Date.now()
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(LEADER_PREFIX)) {
        const raw = localStorage.getItem(key)
        let obj = null
        try { obj = raw ? JSON.parse(raw) : null } catch (e) { obj = null }
        if (!obj || typeof obj.ts !== 'number' || now - obj.ts > LEADER_TTL_MS) {
          localStorage.removeItem(key)
        }
      }
    }
  } catch (e) {
    console.warn('cleanupExpiredLeaderKeys error', e)
   }
}
function getLeader () {
  try {
    const raw = localStorage.getItem(getLeaderKey())
    return raw ? JSON.parse(raw) : null
  } catch (e) {
    return null
  }
}
function becomeLeader () {
  localStorage.setItem(getLeaderKey(), JSON.stringify({ tabId: TAB_ID, ts: Date.now() }))
}
function isLeader () {
  const leader = getLeader()
  if (!leader) {
    becomeLeader()
    return true
  }
  const expired = Date.now() - leader.ts > LEADER_TTL_MS
  if (expired) {
    becomeLeader()
    return true
  }
  return leader.tabId === TAB_ID
}
function releaseLeader () {
  const leader = getLeader()
  if (leader && leader.tabId === TAB_ID) {
    localStorage.removeItem(getLeaderKey())
  }
}
function releaseSpecificLeader (key) {
  try {
    const raw = localStorage.getItem(key)
    const obj = raw ? JSON.parse(raw) : null
    if (obj && obj.tabId === TAB_ID) {
      localStorage.removeItem(key)
    }
  } catch (e) {
    // console.warn('releaseSpecificLeader parse error', e)
  }
}
function storageHandler (e) {
  if (e.key === getLeaderKey()) {
    // 领导者变更时更新自动保存定时器
    updateAutoSaveTimer()
    cleanupExpiredLeaderKeys()
    isLeaderFlag.value = isLeader()
  }
}
function visibilityHandler () {
  // 页面可见时尝试争取领导权，不可见时交由其他标签页
  if (document.visibilityState === 'visible') {
    isLeaderFlag.value = isLeader()
  }
  cleanupExpiredLeaderKeys()
  updateAutoSaveTimer()
}
function startLeaderHeartbeat () {
  stopLeaderHeartbeat()
  // 初始化领导者身份
  cleanupExpiredLeaderKeys()
  isLeaderFlag.value = isLeader()
  heartbeatTimer = setInterval(() => {
    const leaderNow = isLeader()
    isLeaderFlag.value = leaderNow
    if (leaderNow) {
      localStorage.setItem(getLeaderKey(), JSON.stringify({ tabId: TAB_ID, ts: Date.now() }))
    }
  }, HEARTBEAT_MS)
  window.addEventListener('storage', storageHandler)
  document.addEventListener('visibilitychange', visibilityHandler)
  window.addEventListener('beforeunload', () => { releaseLeader(); decrementOpenTabs() })
  window.addEventListener('pagehide', () => { releaseLeader(); decrementOpenTabs() })
}
function stopLeaderHeartbeat () {
  heartbeatTimer && clearInterval(heartbeatTimer)
  heartbeatTimer = null
  window.removeEventListener('storage', storageHandler)
  document.removeEventListener('visibilitychange', visibilityHandler)
  window.removeEventListener('beforeunload', releaseLeader)
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
    if (!isEditing.value) return
    const hash = computeCanvasHash()
    if (hash && hash !== lastChangeHash) {
      lastChangeHash = hash
      lastChangeTs = Date.now()
      // 内容发生变化：若处于非领导者或自动保存已停用，立即恢复
      ensureLeaderAndAutoSaveOnChange()
    }
  }, 2 * 1000)
}

function stopChangeMonitor () {
  changeMonitorTimer && clearInterval(changeMonitorTimer)
  changeMonitorTimer = null
}

// 当检测到内容变更且提示“可能在其他页面编辑”时，立即设为领导者并恢复自动保存
function ensureLeaderAndAutoSaveOnChange () {
  try {
    if (!isEditing.value) return
    if (isLockedByOther.value) return
    const needRecover = (!autoSaveEnabled.value || !isLeader())
    if (needRecover) {
      becomeLeader()
      autoSaveEnabled.value = true
      updateAutoSaveTimer()
    } else {
      // 已是领导者且自动保存开启，也确保心跳与定时器正常
      updateAutoSaveTimer()
    }
  } catch (e) {
    // console.warn('ensureLeaderAndAutoSaveOnChange error', e)
  }
}

function startInactivityWatcher () {
  stopInactivityWatcher()
  inactivityTimer = setInterval(() => {
    if (isEditing.value && autoSaveEnabled.value) {
      const idle = Date.now() - lastChangeTs
      if (idle >= INACTIVITY_MS) {
        autoSaveEnabled.value = false
        message.warning('已超过5分钟无改动，自动保存已暂停')
        updateAutoSaveTimer()
      }
    }
  }, 30 * 1000)
}

function stopInactivityWatcher () {
  inactivityTimer && clearInterval(inactivityTimer)
  inactivityTimer = null
}

// 编辑锁，根据规范调用后端接口
async function checkEditLock () {
  try {
    const res = await getDraftKey({ robot_key: robot_key.value, uni_identifier: getUniIdentifier(), user_agent: buildUserAgent() })
    // const res = {"msg":"success","res":0,"data":{"is_self":true,"lock_res":true,"lock_ttl":955,"remote_addr":"171.83.17.34","robot_key":"yw5BnxX80G","staff_id":3432,"user_agent":"Mozilla\/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/142.0.0.0 Safari\/537.36"}}
    const data = res?.data || {}
    // lock_res 表示是否成功获取到编辑锁；is_self 表示是否为自己
    // 若未获取到锁，则视为他人占用
    isLockedByOther.value = !data.lock_res
    lockRemoteAddr.value = data.remote_addr || ''
    lockUserAgent.value = data.user_agent || ''
    loginUserName.value = data.login_user_name || ''
    robotStore.setIsLockedByOther(isLockedByOther.value)
    updateAutoSaveTimer()
  } catch (e) {
    // 接口异常时默认不锁定，避免阻断编辑
    isLockedByOther.value = false
    robotStore.setIsLockedByOther(isLockedByOther.value)
  }
}

async function acquireEditLock () {
  try {
    const res = await getDraftKey({ robot_key: robot_key.value, uni_identifier: getUniIdentifier(), user_agent: buildUserAgent() })
    // const res = {"msg":"success","res":0,"data":{"is_self":true,"lock_res":true,"lock_ttl":955,"remote_addr":"171.83.17.34","robot_key":"yw5BnxX80G","staff_id":3432,"user_agent":"Mozilla\/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/142.0.0.0 Safari\/537.36"}}
    const data = res?.data || {}
    if (data.lock_res) {
      isLockedByOther.value = false
      lockRemoteAddr.value = ''
      lockUserAgent.value = ''
      loginUserName.value = ''
      robotStore.setIsLockedByOther(isLockedByOther.value)
      return true
    }
    // 他人占用
    isLockedByOther.value = !data.lock_res
    lockRemoteAddr.value = data.remote_addr || ''
    lockUserAgent.value = data.user_agent || ''
    loginUserName.value = data.login_user_name || ''
    robotStore.setIsLockedByOther(isLockedByOther.value)
    return false
  } catch (e) {
    // 异常时默认允许进入编辑（可根据需要改为严格模式）
    isLockedByOther.value = false
    robotStore.setIsLockedByOther(isLockedByOther.value)
    return true
  }
}

function releaseEditLock () {
  // 规范未提供释放接口，此处占位以便未来对接后端释放
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
  if (isLockedByOther.value) {
    // 手动保存时给出提示，自动保存静默跳过
    if (type === 'handle') {
      message.warning('当前已有其他用户在编辑中，无法保存')
    }
    return
  }
  if (!isEditing.value) return
  if (!autoSaveEnabled.value && type === 'automatic') return
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
      user_agent: buildUserAgent()
    }

    const result = await confirmOverrideAndSave(basePayload, false)
    console.log('result', 2)
    if (result.saved) {
      message.success('保存成功')
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
      user_agent: buildUserAgent()
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
const BEHIND_MODAL_TEXT = '当前草稿已经不是最新的草稿版本，保存后会覆盖，确认保存么？'

async function confirmOverrideAndSave (basePayload, isAuto) {
  try {
     const payload = { ...basePayload, uni_identifier: getUniIdentifier(), user_agent: buildUserAgent() }
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
        title: '提示',
        icon: null,
        content: BEHIND_MODAL_TEXT,
        okText: '确 定',
        cancelText: '取 消',
        async onOk () {
          const forcePayload = { ...basePayload, re_cover_save: 1, uni_identifier: getUniIdentifier(), user_agent: buildUserAgent() }
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
        title: '提示',
        icon: null,
        content: BEHIND_MODAL_TEXT,
        okText: '确 定',
        cancelText: '取 消',
        async onOk () {
          const forcePayload = { ...basePayload, re_cover_save: 1, uni_identifier: getUniIdentifier(), user_agent: buildUserAgent() }
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
  if (isLockedByOther.value) return
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
    uni_identifier: getUniIdentifier()
  }
  await confirmOverrideAndSave(basePayload, true)
}
function updateAutoSaveTimer () {
  timer && clearInterval(timer)
  timer = null
  // 同步领导者状态以便 UI 及时展示
  isLeaderFlag.value = isLeader()
  if (import.meta.env.PROD && isEditing.value && autoSaveEnabled.value && !isLockedByOther.value && isLeader()) {
    timer = setInterval(() => {
      handleAutoSaveWithConflictCheck()
    }, 1 * 60 * 1000)
  }
}

onUnmounted(() => {
  timer && clearInterval(timer)
  stopInactivityWatcher()
  stopChangeMonitor()
  stopLeaderHeartbeat()
  releaseEditLock()
  releaseLeader()
  decrementOpenTabs()
})


const versionModelRef = ref(null)

const saveLoading = ref(false)

const openVersionModel = (node_list) => {
  versionModelRef.value.show(node_list)
}
// 发布机器人
const handleRelease = async () => {
  if (isLockedByOther.value) {
    message.warning('当前已有其他用户在编辑中，无法发布')
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
    message.error('存在未关联的节点，请先关联')
    return
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
      user_agent: buildUserAgent()
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
  ensureLeaderAndAutoSaveOnChange()

}

// 删除节点
const onDeleteNode = () => {
  // 删除节点属于内容变更，立即恢复自动保存并设置为领导者
  ensureLeaderAndAutoSaveOnChange()
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
  if (isLockedByOther.value) {
    message.warning('当前已有其他用户在编辑中，无法恢复')
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
    user_agent: buildUserAgent()
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
    message.success('已保存当前草稿并恢复版本')
    // 恢复后重新启用自动保存（回到编辑态）
    isEditing.value = true
    autoSaveEnabled.value = true
    updateAutoSaveTimer()
  })
}

const handlePreviewVersion = async (data, version) => {
  // 在发布详情中切换/预览前，先保存当前草稿，避免丢失
  if (!isLockedByOther.value) {
    const currentDraft = getCanvasData()
    try {
      await saveNodes({
        robot_key: robot_key.value,
        data_type: 1,
        node_list: JSON.stringify(currentDraft),
        draft_save_type: 'automatic',
        draft_save_time: +robotStore.robotInfo.draft_save_time || 0,
        uni_identifier: getUniIdentifier(),
        user_agent: buildUserAgent()
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
const handleRunTest = () => {
  pageHeaderRef.value.openRunTest()
}

// 头部“编辑”按钮
const handleClickEdit = async () => {
  if (isLockedByOther.value) {
    message.warning('当前已有其他用户在编辑中，无法编辑')
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
    message.success('已进入编辑模式')
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
      title: '工作流异常修复',
      content: h('div', {}, [
        h('div', {style: titleStyle}, '以下插件状态异常，是否自动安装插件并开启？'),
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
        message.success('操作完成，即将重新加载...')
        setTimeout(() => {
          window.location.reload()
        }, 1200)
      }
    })
  }
}

const init = async () => {
  await getModelList()
  await workflowStore.getTriggerList(robot_key.value);
  workflowStore.getTriggerOfficialMsg(robot_key.value)
  await modelStore.getAllmodelList()
  workflowStore.getAllLibraryList();

  const res = await getNodeList({
    robot_key: robot_key.value,
    data_type: 1
  })

  getNode(res.data)
  checkNodePluginStatus(res.data)
}

const handleAutoSaveDraft = async (type = 'automatic') => {
  await handleSave(type)
}

provide('handleAutoSaveDraft', handleAutoSaveDraft)

onMounted(async () => {
  // 记录当前工作流的打开页面计数
  incrementOpenTabs()

  await init()

  // 初次进入：若存在草稿记录则默认查看模式，需要点击编辑按钮进入编辑
  await checkEditLock()
  if (isLockedByOther.value) {
    isEditing.value = false
    autoSaveEnabled.value = false
    message.warning('其他人在编辑，已切换为查看模式')
  } else {
    // 无他人占用，直接进入编辑模式（无需手动点击编辑）
    const ok = await acquireEditLock()
    if (ok) {
      isEditing.value = true
      autoSaveEnabled.value = true
    } else {
      isEditing.value = false
      autoSaveEnabled.value = false
    }
  }
  // 启动领导者协调与心跳，仅允许一个标签页自动保存
  startLeaderHeartbeat()
  updateAutoSaveTimer()

  // 监听路由中 robot_key 变化，切换工作流时释放旧领导者并重启心跳
  let currentLeaderKey = getLeaderKey()
  watch(() => route.query.robot_key, (newKey) => {
    if (!newKey || newKey === robot_key.value) return
    // 释放旧工作流的领导者占用（仅当本标签页是旧键的领导者时）
    releaseSpecificLeader(currentLeaderKey)
    // 更新旧工作流的打开页面计数
    decrementOpenTabs()
    // 切换到新的工作流键
    robot_key.value = newKey
    currentLeaderKey = getLeaderKey()
    // 新工作流计数+1
    incrementOpenTabs()
    // 重新争取领导者并重启自动保存评估
    startLeaderHeartbeat()
    updateAutoSaveTimer()
  })

  // 活动监听与5分钟无操作停止
  startInactivityWatcher()
  startChangeMonitor()
  if (route.query.show_tips) {
    message.info('按住Shift 滚动鼠标可左右移动画布，按住Ctrl 滚动鼠标可放大缩小画布', 6)
  }
})

</script>
