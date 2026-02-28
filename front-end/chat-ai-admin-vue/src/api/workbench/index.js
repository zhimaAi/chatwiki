import request from '@/utils/http/axios'

/**
 * 获取工作台配置
 */
export const getWorkbenchConfig = (params = {}) => {
  return request.get({
    url: '/manage/workbench/getConfig',
    params: params
  })
}

/**
 * 工作台保存配置
 * @param {Object} data - 配置数据
 * @param {string} data.default_robot_id - 首页机器人ID
 * @param {string} data.enable_last_app_entry - 默认进入上一次访问的机器人，1开启，0关闭
 */
export const saveWorkbenchConfig = (data = {}) => {
  return request.post({
    url: '/manage/workbench/saveConfig',
    data: data
  })
}

/**
 * 工作台机器人置顶
 * @param {Object} data - 置顶数据
 * @param {string} data.robot_id - 置顶机器人ID
 */
export const topWorkbenchRobot = (data = {}) => {
  return request.post({
    url: '/manage/workbench/topRobot',
    data: data
  })
}

/**
 * 获取机器人分组列表
 */
export const getRobotGroupList = (params = {}) => {
  return request.get({
    url: '/manage/getRobotGroupList',
    params: params
  })
}

/**
 * 工作台团队机器人列表
 */
export const getWorkbenchTeamRobotList = (params = {}) => {
  return request.get({
    url: '/manage/workbench/teamRobotList',
    params: params
  })
}

/**
 * 获取历史访问机器人列表
 */
export const getRobotHistoryVisit = (params = {}) => {
  return request.get({
    url: '/manage/workbench/robotHistoryVisit',
    params: params
  })
}

/**
 * 工作台记录使用机器人
 * @param {Object} data - 记录数据
 * @param {string} data.robot_id - 机器人ID
 */
export const postRecordRobotVisit = (data = {}) => {
  return request.post({
    url: '/manage/workbench/recordRobotVisit',
    data: data
  })
}

