import { ref } from 'vue'
import { defineStore } from 'pinia'
import { getRobotInfo, getRobotList, saveRobot } from '@/api/robot'
import { normalizeCustomerServiceQueryValue } from '@/views/customer-service/navigation'

// 客服支持聊天机器人(0)和智能体(2)；先排除空值/空白，避免转成0误收。
function isCustomerServiceApplication(robot) {
  const type = robot?.application_type
  return (typeof type === 'number' || typeof type === 'string') &&
    String(type).trim() !== '' && [0, 2].includes(Number(type))
}

function stringifyQuestionConfig(value) {
  if (typeof value === 'string') return value
  const config = value || { content: '', question: [] }
  return JSON.stringify({
    ...config,
    question: (config.question || []).map((item) => typeof item === 'string' ? item : item?.content || '')
  })
}

function buildSettingsPayload(source, patch) {
  const payload = JSON.parse(JSON.stringify({ ...source, ...patch }))
  if (source.robot_avatar instanceof File) {
    payload.robot_avatar = new File([source.robot_avatar], source.robot_avatar.name)
    delete payload.robot_avatar_url
  } else if (source.robot_avatar || source.robot_avatar_url) {
    payload.robot_avatar_url = source.robot_avatar || source.robot_avatar_url
  }
  payload.welcomes = stringifyQuestionConfig(payload.welcomes)
  payload.unknown_question_prompt = stringifyQuestionConfig(payload.unknown_question_prompt)
  for (const field of ['prompt_struct', 'cache_config']) {
    if (typeof payload[field] !== 'string') payload[field] = JSON.stringify(payload[field] ?? {})
  }
  // 保存应用的真实类型，数字0是有效值，不能用逻辑或表达式兜成类型2。
  payload.application_type = Number(source.application_type)
  return payload
}

// 管理机器人基础数据及正式模型/人设配置；URL标识、默认选择和router跳转由客服布局负责。
export const useCustomerServiceStore = defineStore('customer-service', () => {
  // robots为类型0/2的应用列表，robot为已校验详情；初始及重置时分别为空数组、null。
  // 列表状态为idle/loading/ready/error（空列表也ready）；详情为empty/loading/ready/error。
  // robotError是详情提示key：空串表示无错误，invalid_robot为标识/类型不符，error_robot为请求失败。
  const robots = ref([])
  const robot = ref(null)
  const robotListStatus = ref('idle')
  const robotStatus = ref('empty')
  const robotError = ref('')
  // 列表、详情分别计数，互不阻挡；新请求或reset后，旧成功/失败结果都不能覆盖当前状态。
  let listRequestSequence = 0
  let detailRequestSequence = 0
  let configRequestSequence = 0
  let pendingConfigRequest = 0
  let savingSettings = false

  // 只加载候选列表，不改当前详情；成功返回过滤后的数组（含空数组），失败/过期返回null。
  const loadRobots = async () => {
    const requestId = ++listRequestSequence
    robots.value = []
    robotListStatus.value = 'loading'

    try {
      // 接口默认排除智能体，显式关闭该过滤后再排除客服不支持的工作流。
      const response = await getRobotList({ filter_claw: false })
      if (requestId !== listRequestSequence) return null
      if (!Array.isArray(response?.data)) throw new Error('Invalid robot list response')

      robots.value = response.data.filter(isCustomerServiceApplication)
      robotListStatus.value = 'ready'
      return robots.value
    } catch {
      if (requestId !== listRequestSequence) return null
      robotListStatus.value = 'error'
      return null
    }
  }

  // 按显式标识独立加载详情，不依赖列表成功；成功返回详情，无参/无效/失败/过期均返回null。
  const loadRobot = async ({ id, robotKey } = {}) => {
    const requestId = ++detailRequestSequence
    configRequestSequence++
    pendingConfigRequest = 0
    const targetId = normalizeCustomerServiceQueryValue(id)
    const targetKey = normalizeCustomerServiceQueryValue(robotKey)
    // 先清旧详情和错误，避免切换期间仍把上一应用当成当前应用展示。
    robot.value = null
    robotError.value = ''

    // 完全无参不请求；是否从列表补默认项由布局处理。
    if (!targetId && !targetKey) {
      robotStatus.value = 'empty'
      return null
    }
    // 残缺标识只提示链接无效，不在Store内默认换应用。
    if (!targetId || !targetKey) {
      robotStatus.value = 'error'
      robotError.value = 'invalid_robot'
      return null
    }

    robotStatus.value = 'loading'
    try {
      const response = await getRobotInfo({ id: targetId })
      if (requestId !== detailRequestSequence) return null

      const robotInfo = response?.data
      const responseId = normalizeCustomerServiceQueryValue(robotInfo?.id)
      // 缺失ID不能转成字面量undefined/null后误匹配；残缺或不一致的链接不改选机器人。
      if (
        !isCustomerServiceApplication(robotInfo) || !responseId || responseId !== targetId ||
        !robotInfo.robot_key || robotInfo.robot_key !== targetKey
      ) {
        robotStatus.value = 'error'
        robotError.value = 'invalid_robot'
        return null
      }

      robot.value = robotInfo
      robotStatus.value = 'ready'
      return robot.value
    } catch {
      if (requestId !== detailRequestSequence) return null
      robotStatus.value = 'error'
      robotError.value = 'error_robot'
      return null
    }
  }

  const isCurrentRobot = (snapshot) => {
    return robotStatus.value === 'ready' && !!snapshot && robot.value === snapshot &&
      !!normalizeCustomerServiceQueryValue(snapshot.id) && !!snapshot.robot_key && isCustomerServiceApplication(snapshot)
  }

  /**
   * expectedIdentity必须是请求前的原robot对象，不是新造的{id,key}。
   * 每次welcome重新捕获；配置更新统一替换对象，旧welcome只能丢弃配置回填，不能丢弃会话结果。
   */
  const mergeCurrentRobotInfo = (expectedIdentity, partialRobot) => {
    if (!isCurrentRobot(expectedIdentity) || savingSettings || pendingConfigRequest ||
      !partialRobot || typeof partialRobot !== 'object' || Array.isArray(partialRobot)) return false
    if ('id' in partialRobot && normalizeCustomerServiceQueryValue(partialRobot.id) !==
      normalizeCustomerServiceQueryValue(expectedIdentity.id)) return false
    if ('robot_key' in partialRobot && partialRobot.robot_key !== expectedIdentity.robot_key) return false
    // 同一身份的配置回填不能将聊天机器人与智能体互相转换。
    if ('application_type' in partialRobot && (
      !isCustomerServiceApplication(partialRobot) ||
      Number(partialRobot.application_type) !== Number(expectedIdentity.application_type)
    )) return false
    robot.value = {
      ...expectedIdentity,
      ...partialRobot,
      id: expectedIdentity.id,
      robot_key: expectedIdentity.robot_key,
      application_type: expectedIdentity.application_type
    }
    configRequestSequence++
    return true
  }

  /** 传入本次刷新前的原robot对象；失败/过期返回null，不清robot或ready，不重建聊天。 */
  const refreshCurrentRobot = async (expectedIdentity) => {
    if (!isCurrentRobot(expectedIdentity)) return null
    const detailSequence = detailRequestSequence
    const requestId = ++configRequestSequence
    pendingConfigRequest = requestId
    // 保留配置值，仅更新快照边界，使刷新前发出的welcome不能随后覆盖新配置。
    robot.value = { ...expectedIdentity }
    const refreshSnapshot = robot.value
    try {
      const response = await getRobotInfo({ id: normalizeCustomerServiceQueryValue(expectedIdentity.id) })
      if (detailSequence !== detailRequestSequence || requestId !== configRequestSequence ||
        !isCurrentRobot(refreshSnapshot)) return null
      const nextRobot = response?.data
      if (!isCustomerServiceApplication(nextRobot) || !normalizeCustomerServiceQueryValue(nextRobot?.id) ||
        normalizeCustomerServiceQueryValue(nextRobot.id) !== normalizeCustomerServiceQueryValue(expectedIdentity.id) ||
        Number(nextRobot.application_type) !== Number(expectedIdentity.application_type) ||
        !nextRobot.robot_key || nextRobot.robot_key !== expectedIdentity.robot_key) return null
      robot.value = nextRobot
      return robot.value
    } catch {
      return null
    } finally {
      if (pendingConfigRequest === requestId) pendingConfigRequest = 0
    }
  }

  const saveRobotSettings = async (patch) => {
    if (savingSettings) throw new Error('settings_save_in_progress')
    const source = robot.value
    if (!isCurrentRobot(source)) throw new Error('robot_not_ready')
    if (!patch || typeof patch !== 'object' || Array.isArray(patch)) throw new Error('invalid_settings_patch')
    const keys = Object.keys(patch)
    const isPromptPatch = keys.length === 1 && keys[0] === 'prompt' && typeof patch.prompt === 'string'
    const isModelPatch = keys.length === 2 && keys.includes('model_config_id') && keys.includes('use_model') &&
      patch.model_config_id !== undefined && patch.model_config_id !== null && typeof patch.use_model === 'string'
    if (!isPromptPatch && !isModelPatch) throw new Error('invalid_settings_patch')

    const detailSequence = detailRequestSequence
    const payload = buildSettingsPayload(source, patch)
    savingSettings = true
    configRequestSequence++
    pendingConfigRequest = 0
    // 不乐观改模型/人设；更换快照只隔离保存前的welcome及配置刷新。
    robot.value = { ...source }
    try {
      await saveRobot(payload, payload.application_type)
      // 服务端保存已经成功；换身份或刷新失败只能报告部分成功，不能伪装为保存失败。
      if (detailSequence !== detailRequestSequence || robotStatus.value !== 'ready' ||
        normalizeCustomerServiceQueryValue(robot.value?.id) !== normalizeCustomerServiceQueryValue(source.id) ||
        robot.value?.robot_key !== source.robot_key) return { saved: true, refreshed: false }
      const refreshed = await refreshCurrentRobot(robot.value)
      return { saved: true, refreshed: !!refreshed }
    } finally {
      // reset不提前解此锁：离开前已到服务端的保存无法撤销，不能并发覆盖同一正式配置。
      savingSettings = false
    }
  }

  const reset = () => {
    // 布局重新进入或卸载时，让在途结果失效，不取消其他模块的请求。
    // 序号必须继续递增，不能归零：后续新请求可能与旧请求撞号，导致旧结果重新被接纳。
    listRequestSequence++
    detailRequestSequence++
    configRequestSequence++
    pendingConfigRequest = 0
    robots.value = []
    robot.value = null
    robotListStatus.value = 'idle'
    robotStatus.value = 'empty'
    robotError.value = ''
  }

  const syncRobot = data => {
    if (!isCurrentRobot(robot.value) || !isCustomerServiceApplication(data) ||
      Number(data.application_type) !== Number(robot.value.application_type)) return
    if (String(data?.id) !== String(robot.value?.id) || data.robot_key !== robot.value?.robot_key) return
    robot.value = data
    const index = robots.value.findIndex(item => String(item.id) === String(data.id))
    if (index >= 0) robots.value[index] = { ...robots.value[index], ...data }
  }

  return {
    syncRobot,
    robots,
    robot,
    robotListStatus,
    robotStatus,
    robotError,
    loadRobots,
    loadRobot,
    mergeCurrentRobotInfo,
    refreshCurrentRobot,
    saveRobotSettings,
    reset
  }
})
