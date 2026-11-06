import { ref, computed, watch } from 'vue'
import { defineStore } from 'pinia'
import { message } from 'ant-design-vue'
import { DEFAULT_ROBOT_AVATAR } from '@/constants'
import { deleteRobot, getRobotList, getRobotInfo, robotCopy, saveRobot } from '@/api/robot'
import { saveClawbotConf } from '@/api/clawbot'
import { checkRole, getRobotPermission } from '@/utils/permission'
import { usePermissionStore } from '@/stores/modules/permission'

// 默认助手配置
const DEFAULT_ASSISTANT_NAME = 'Agent'
const DEFAULT_ASSISTANT_INTRO = '智能对话助手，为您解答各类问题，提供专业建议。'
const CLAWBOT_APPLICATION_TYPE = 2 // Agent 类型标识
const STORAGE_KEY = 'clawbot-current-id' // localStorage 持久化 key
let robotInfoRequestSeq = 0

const normalizeQuestion = (item) => {
  if (typeof item === 'string') return item
  return item?.content || ''
}

const stringifyQuestionConfig = (data = {}) => {
  if (typeof data === 'string') {
    return data
  }

  return JSON.stringify({
    ...data,
    question: (data.question || []).map(normalizeQuestion)
  })
}

const stringifyJsonField = (value, fallback) => {
  if (typeof value === 'string') {
    return value
  }

  return JSON.stringify(value ?? fallback)
}

/**
 * 标准化助手数据，确保字段完整且类型正确
 */
const normalizeAssistant = (item = {}) => ({
  ...item,
  id: item.id === undefined || item.id === null ? '' : String(item.id),
  robot_key: item.robot_key || '',
  robot_name: item.robot_name || DEFAULT_ASSISTANT_NAME,
  robot_intro: item.robot_intro || DEFAULT_ASSISTANT_INTRO,
  robot_avatar: item.robot_avatar || item.robot_avatar_url || DEFAULT_ROBOT_AVATAR,
  robot_avatar_url: item.robot_avatar_url || item.robot_avatar || DEFAULT_ROBOT_AVATAR,
  application_type: Number(item.application_type ?? CLAWBOT_APPLICATION_TYPE),
})

/**
 * Clawbot 模块级 Store（助手管理）
 *
 * 职责：管理助手列表、当前选中助手、模块初始化
 * 生命周期：进入 /clawbot 模块时初始化，切换子页面不丢失
 *
 * 模块进入时通过 initModule() 锁定一个机器人：
 * - 拉取助手列表，独立入口默认在列表为空时创建助手，新版入口按 URL 指定应用
 * - 初始化完成后 isReady=true，网关组件据此决定是否渲染子路由
 */
export const useClawbotStore = defineStore('clawbot', () => {
  // ===== 模块初始化状态 =====
  const isReady = ref(false)           // 模块是否就绪（所有子页面依赖此状态）
  const isInitializing = ref(false)    // 是否正在执行 initModule
  const loading = ref(false)           // 助手列表加载中
  const creating = ref(false)          // 助手创建中
  let moduleInitialization = null
  // 每次请求捕获发起时的判断函数，避免布局切换后使用新上下文接收旧请求结果。
  let contextIsCurrent = () => true

  // ===== 助手管理 =====
  const assistantList = ref([])        // 当前用户的助手列表
  // 从 localStorage 恢复上次选中的助手 ID
  const currentAssistantId = ref(localStorage.getItem(STORAGE_KEY) || '')

  // ===== 机器人详细配置 =====
  const robotInfo = ref(null)

  // 自动同步 currentAssistantId 到 localStorage，并在就绪后重新拉取 robotInfo
  watch(currentAssistantId, (val) => {
    if (val) {
      localStorage.setItem(STORAGE_KEY, val)
    } else {
      localStorage.removeItem(STORAGE_KEY)
      robotInfo.value = null
    }
    // 模块就绪后切换助手时重新获取详细配置
    if (isReady.value && !isInitializing.value && val) {
      fetchRobotInfo(val)
    }
  })

  // 当前选中的助手（派生计算属性）
  const currentAssistant = computed(() => {
    return assistantList.value.find((item) => item.id === currentAssistantId.value) || null
  })

  /**
   * 应用助手列表：标准化数据 + 校正当前选中项
   * 默认在当前助手失效时回退到第一个；新版初始化保留选中项，交由 URL 校验决定是否进入。
   */
  function applyAssistantList(list = [], { preserveSelection = false } = {}) {
    const normalizedList = (list || []).map(normalizeAssistant).filter((item) => item.id)
    assistantList.value = normalizedList

    if (preserveSelection) return

    if (!normalizedList.length) {
      currentAssistantId.value = ''
      return
    }

    const hasCurrent = normalizedList.some((item) => item.id === currentAssistantId.value)
    if (!hasCurrent) {
      currentAssistantId.value = normalizedList[0].id
    }
  }

  function syncAssistantInfo(data) {
    if (!data?.id) return

    const normalized = normalizeAssistant(data)
    const index = assistantList.value.findIndex((item) => item.id === normalized.id)

    if (index >= 0) {
      assistantList.value.splice(index, 1, {
        ...assistantList.value[index],
        ...normalized
      })
      return
    }

    if (normalized.id === currentAssistantId.value) {
      assistantList.value.unshift(normalized)
    }
  }

  function replaceRobotInfo(data) {
    robotInfo.value = data
    syncAssistantInfo(data)
  }

  function buildAssistantSavePayload(rawState = {}, patch = {}) {
    const formData = JSON.parse(JSON.stringify({
      ...rawState,
      ...patch
    }))

    if (rawState.robot_avatar instanceof File) {
      formData.robot_avatar = new File([rawState.robot_avatar], rawState.robot_avatar.name)
      delete formData.robot_avatar_url
    } else if (rawState.robot_avatar) {
      formData.robot_avatar_url = rawState.robot_avatar
    } else if (rawState.robot_avatar_url) {
      formData.robot_avatar_url = rawState.robot_avatar_url
    }

    formData.welcomes = stringifyQuestionConfig(formData.welcomes || { content: '', question: [] })
    formData.prompt_struct = stringifyJsonField(formData.prompt_struct, {})
    formData.unknown_question_prompt = stringifyQuestionConfig(
      formData.unknown_question_prompt || { content: '', question: [] }
    )
    formData.cache_config = stringifyJsonField(formData.cache_config, {})

    return {
      formData,
      applicationType: Number(rawState.application_type || currentAssistant.value?.application_type || 2)
    }
  }

  /** 合并更新 robotInfo（供 chat store 从 chatWelcome 响应同步部分 robot 数据） */
  function updateRobotInfo(data) {
    if (!data) return
    robotInfo.value = {
      ...(robotInfo.value || {}),
      ...data
    }
    syncAssistantInfo(data)
  }

  /** 获取当前助手的详细配置 */
  async function fetchRobotInfo(targetId = currentAssistantId.value) {
    const id = targetId ? String(targetId) : ''
    const requestSeq = ++robotInfoRequestSeq
    const isCurrent = contextIsCurrent

    if (!id) {
      robotInfo.value = null
      return
    }

    try {
      const res = await getRobotInfo({ id })
      if (!isCurrent() || requestSeq !== robotInfoRequestSeq || currentAssistantId.value !== id) {
        return res
      }
      replaceRobotInfo(res?.data || null)
      return res
    } catch {
      if (isCurrent() && requestSeq === robotInfoRequestSeq && currentAssistantId.value === id) {
        robotInfo.value = null
      }
    }
  }

  /** 更新 Clawbot 配置（开关等），保存后刷新 robotInfo */
  async function updateClawbotConf(data) {
    await saveClawbotConf(data)
    await fetchRobotInfo()
  }

  /** 保存当前助手配置，支持局部字段更新和可选乐观更新 */
  async function saveAssistant(partialData = {}, options = {}) {
    const {
      optimistic = false,
      successMessage = '',
      refreshAfterSave = true,
      rollbackOnError = false
    } = options

    const sourceRobotInfo = robotInfo.value ? { ...robotInfo.value } : null
    const isCurrent = contextIsCurrent
    const canApply = () => isCurrent() && currentAssistantId.value === String(sourceRobotInfo?.id || '')
    if (!sourceRobotInfo?.id) {
      message.error('当前助手信息未加载完成，请稍后重试')
      throw new Error('clawbot assistant info missing')
    }

    const rollbackState = rollbackOnError
      ? Object.keys(partialData).reduce((result, key) => {
          result[key] = sourceRobotInfo[key]
          return result
        }, {})
      : null

    if (optimistic) {
      updateRobotInfo(partialData)
    }

    try {
      const { formData, applicationType } = buildAssistantSavePayload(sourceRobotInfo, partialData)
      const res = await saveRobot(formData, applicationType)

      if (refreshAfterSave) {
        await fetchRobotInfo(String(sourceRobotInfo.id))
      } else if (!optimistic && canApply()) {
        updateRobotInfo(partialData)
      }

      if (successMessage) {
        message.success(successMessage)
      }

      return res
    } catch (err) {
      if (optimistic && rollbackOnError && rollbackState && canApply()) {
        updateRobotInfo(rollbackState)
      }
      throw err
    }
  }

  /** 拉取助手列表；上下文已失效时返回 null，调用方不能继续选中旧请求中的助手。 */
  async function fetchAssistants({ isCurrent = contextIsCurrent, preserveSelection = false } = {}) {
    loading.value = true
    try {
      const res = await getRobotList({ application_type: CLAWBOT_APPLICATION_TYPE })
      await usePermissionStore().checkPermission()
      if (!isCurrent()) return null
      const list = (res?.data || []).filter((item) => getRobotPermission(item.id) > 0)
      applyAssistantList(list, { preserveSelection })
      return assistantList.value
    } finally {
      loading.value = false
    }
  }

  /** 创建默认助手，创建后自动刷新列表并选中 */
  async function createAssistant() {
    if (!checkRole(['ClawBotCreate'])) return null
    const isCurrent = contextIsCurrent
    creating.value = true
    try {
      const res = await saveRobot(
        {
          robot_name: DEFAULT_ASSISTANT_NAME,
          robot_intro: DEFAULT_ASSISTANT_INTRO,
          robot_avatar_url: DEFAULT_ROBOT_AVATAR,
          group_id: '0',
        },
        CLAWBOT_APPLICATION_TYPE
      )
      const createdId = String(res?.data?.id || '')
      if (!isCurrent()) return null
      await fetchAssistants({ isCurrent })
      if (!isCurrent()) return null
      if (createdId) {
        selectAssistant(createdId)
      } else if (assistantList.value.length) {
        currentAssistantId.value = assistantList.value[0].id
      }
      return currentAssistant.value
    } finally {
      creating.value = false
    }
  }

  /** 删除后切换到剩余助手；独立入口默认在列表为空时创建，新版入口可关闭自动创建。 */
  async function deleteAssistant(id, { autoCreate = true } = {}) {
    const normalizedId = String(id)
    const isCurrent = contextIsCurrent
    await deleteRobot({ id: normalizedId })
    if (!isCurrent()) return
    const wasCurrent = currentAssistantId.value === normalizedId
    await fetchAssistants({ isCurrent })
    if (!isCurrent()) return
    if (!assistantList.value.length) {
      if (autoCreate) await createAssistant()
      return
    }
    if (wasCurrent && assistantList.value.length) {
      currentAssistantId.value = assistantList.value[0].id
    }
  }

  /** 复制助手 */
  async function copyAssistant(id) {
    const isCurrent = contextIsCurrent
    await robotCopy({ id: String(id) })
    if (isCurrent()) await fetchAssistants({ isCurrent })
  }

  const getAssistantByQuery = (targetId, targetRobotKey = '') => {
    const normalizedId = targetId === undefined || targetId === null ? '' : String(targetId)
    const normalizedRobotKey = targetRobotKey === undefined || targetRobotKey === null ? '' : String(targetRobotKey)

    if (!normalizedId) {
      return null
    }

    const assistant = assistantList.value.find((item) => item.id === normalizedId)
    if (!assistant) {
      return null
    }

    if (normalizedRobotKey && String(assistant.robot_key || '') !== normalizedRobotKey) {
      return null
    }

    return assistant
  }

  /** 选中指定助手（ID 必须在列表中存在） */
  function selectAssistant(id) {
    const normalizedId = String(id)
    const exists = assistantList.value.some((item) => item.id === normalizedId)
    if (exists) {
      currentAssistantId.value = normalizedId
    }
    return currentAssistant.value
  }

  function selectAssistantByQuery(targetId, targetRobotKey = '') {
    const assistant = getAssistantByQuery(targetId, targetRobotKey)
    if (!assistant) {
      return null
    }

    currentAssistantId.value = assistant.id
    return assistant
  }

  /**
   * 模块初始化入口（幂等）
   *
   * 由独立布局或新版布局调用，串行处理两种入口对同一 store 的初始化。
   * 默认保留独立入口的自动创建行为；新版传入 autoCreate=false，严格校验 URL 中的应用。
   * isCurrent 用于丢弃切换应用、重试或卸载后返回的结果；新版还会单独核对就绪身份。
   */
  async function initModule(options = {}) {
    // 同一模块的两种布局可连续进入；后一入口必须等待并按自己的 URL 重新初始化。
    while (moduleInitialization) {
      await moduleInitialization.catch(() => null)
    }
    if (options.isCurrent && !options.isCurrent()) return null
    moduleInitialization = initializeModule(options)
    try {
      return await moduleInitialization
    } finally {
      moduleInitialization = null
    }
  }

  async function initializeModule({
    forceRefresh = false,
    targetId = '',
    targetRobotKey = '',
    autoCreate = true,
    isCurrent = () => true
  } = {}) {
    contextIsCurrent = isCurrent

    if (isReady.value && !forceRefresh) {
      if (targetId) {
        selectAssistantByQuery(targetId, targetRobotKey)
      }
      return currentAssistant.value
    }

    isInitializing.value = true
    try {
      await fetchAssistants({ isCurrent, preserveSelection: !autoCreate })
      if (!isCurrent()) return null
      // 新版只能打开明确指定的应用，不能把失效链接兜底为其他 Agent 或触发创建。
      if (!autoCreate && (!targetId || !getAssistantByQuery(targetId, targetRobotKey))) return null
      if (!assistantList.value.length) {
        if (autoCreate) await createAssistant()
      } else if (targetId) {
        selectAssistantByQuery(targetId, targetRobotKey)
      }
      // 拉取当前助手的详细配置
      await fetchRobotInfo(currentAssistantId.value)
      if (!isCurrent()) return null
      isReady.value = true
      return currentAssistant.value
    } finally {
      isInitializing.value = false
    }
  }

  return {
    // 模块状态
    isReady,
    isInitializing,
    loading,
    creating,

    // 助手管理
    assistantList,
    currentAssistantId,
    currentAssistant,
    initModule,
    fetchAssistants,
    createAssistant,
    deleteAssistant,
    copyAssistant,
    selectAssistant,
    selectAssistantByQuery,

    // 机器人详细配置
    robotInfo,
    fetchRobotInfo,
    saveAssistant,
    updateRobotInfo,
    updateClawbotConf,
  }
})
