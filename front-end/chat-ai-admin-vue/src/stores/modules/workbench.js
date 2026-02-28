import { defineStore } from 'pinia'
import { store } from '../index'
import { getWorkbenchConfig, getRobotHistoryVisit, postRecordRobotVisit } from '@/api/workbench'
import { getRobotInfo } from '@/api/robot/index'

export const useWorkbenchStore = defineStore('workbench', {
  state: () => {
    return {
      // 工作台相关状态
      currentTab: '',
      sidebarCollapsed: false,
      // 当前机器人 robot_key
      currentRobotKey: '',
      // 工作台配置数据（config, home）
      workbenchConfig: null,
      // 最近使用的机器人列表（预留）
      recentRobots: [],
      // iframe 基础地址
      iframeBaseUrl: '',
      // 选中状态：'home' 表示首页，robotId 表示最近使用中的某项
      selectedMode: '',
      isNewChat: true,
      // 是否正在初始化
      isInitializing: false,
      // 空状态
      showEmptyState: false,
      // iframe 相关状态
      iframeLoading: false,
      iframeKey: 0,
    }
  },
  getters: {
    getCurrentTab() {
      return this.currentTab
    },
    isSidebarCollapsed() {
      return this.sidebarCollapsed
    },
    getCurrentRobotKey() {
      return this.currentRobotKey
    },
    getWorkbenchConfig() {
      return this.workbenchConfig
    },
    // iframe 地址
    iframeSrc() {
      // 1. 确定基础地址来源（开发环境用本地，生产环境用同步到的域名）
      const baseUrl = import.meta.env.DEV
        ? import.meta.env.VITE_WORKBENCH_CHAT_IFRAME_BASE_URL
        : this.iframeBaseUrl

      // 2. 统一校验：确保 Key 和 BaseUrl 都存在
      if (!this.currentRobotKey || !baseUrl) {
        return ''
      }

      // 3. 拼接最终地址
      return `${baseUrl}/#/pc?robot_key=${this.currentRobotKey}&isForceNewChat=${Number(this.isNewChat)}`
    }
  },
  actions: {
    setCurrentTab(tab) {
      this.currentTab = tab
    },
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },
    setSidebarCollapsed(collapsed) {
      this.sidebarCollapsed = collapsed
    },
    /**
     * 获取工作台配置
     */
    async fetchWorkbenchConfig() {
      try {
        const res = await getWorkbenchConfig()
        const data = res?.data || {}

        if (Object.keys(data).length) {
          this.workbenchConfig = data
          
          return data
        }

        // 无数据时返回当前缓存，避免状态被清空
        return this.workbenchConfig || {}
      } catch (error) {
        console.error('获取工作台配置失败:', error)
        // 异常时保留现有状态
        return this.workbenchConfig || {}
      }
    },

    async initConfig() {
      this.isInitializing = true
      this.showEmptyState = false

      try {
        await this.fetchWorkbenchConfig()

        const { config, home } = this.workbenchConfig || {}

        // 检查 home 和必要的字段是否存在
        if (!home || !home.robot_id || !home.robot_key) {
          this.showEmptyState = true
          this.selectHome(null, null, true)
          return
        }

        const robotId = home.robot_id
        const robotKey = home.robot_key

        // 获取机器人信息，同步 iframe 基础地址
        await this.fetchRobotInfo(robotId)

        // 根据配置选择模式
        if (config?.enable_last_app_entry == 1) {
          this.selectRecent(robotKey, robotId, false)
        } else {
          this.selectHome(robotKey, robotId, true)
        }
      } finally {
        this.isInitializing = false
      }
    },

    /**
     * 刷新首页配置并返回 home 关键字段
     */
    async refreshHomeConfig() {
      const {home, config} = await this.fetchWorkbenchConfig()
      
      let robotKey = config.default_robot_key
      let robotId = config.default_robot_id

      if(!robotKey || !robotId){
        robotKey = home.robot_key
        robotId = home.robot_id
      } 

      this.selectHome(robotKey, robotId, true)
      return { robot_key: robotKey, robot_id: robotId}
    },

    /**
     * 设置选中模式为首页
     * @param {string} robotKey - 机器人key
     * @param {string} robotId - 机器人ID
     * @param {boolean} isNewChat - 是否为新聊天
     */
    selectHome(robotKey, robotId, isNewChat = true) {
      this.isNewChat = isNewChat
      this.currentRobotKey = robotKey
      this.selectedMode = 'home'
      if(robotId && robotKey){
        this.recordRobotVisit(robotId, true)
      }
    },

    /**
     * 设置选中模式为最近使用项
     * @param {string} robotKey - 机器人key
     * @param {string} robotId - 机器人ID
     */
    selectRecent(robotKey, robotId, isNewChat = false) {
      this.isNewChat = isNewChat
      this.currentRobotKey = robotKey
      this.selectedMode = robotId
    },
    /**
     * 更新最近使用机器人列表（预留）
     */
    updateRecentRobots(robot) {
      // TODO: 实现最近使用记录更新逻辑
      // 当后端接口实现后，调用接口更新最近使用列表
      console.log('更新最近使用机器人:', robot)
    },
    /**
     * 获取历史访问机器人列表
     */
    async fetchRobotHistoryVisit() {
      try {
        const res = await getRobotHistoryVisit()
        if (res && res.data) {
          this.recentRobots = res.data
          return res.data
        }
      } catch (error) {
        console.error('获取历史访问列表失败:', error)
      }
    },
    /**
     * 记录使用机器人
     * @param {string} robotId - 机器人ID
     * @param {boolean} refreshList - 是否刷新最近使用列表，默认true
     */
    async recordRobotVisit(robotId, refreshList = true) {
      try {
        const res = await postRecordRobotVisit({ robot_id: robotId })
        if (res) {
          // 更新工作台配置
          // this.workbenchConfig = res.data

          // this.fetchWorkbenchConfig()

          // 记录成功后刷新最近使用列表
          if (refreshList) {
            this.fetchRobotHistoryVisit()
          }
        }
      } catch (error) {
        console.error('记录使用机器人失败:', error)
      }
    },
    /**
     * 获取机器人信息并同步 iframe 基础地址
     * @param {string} id - 机器人ID
     */
    async fetchRobotInfo(id) {
      try {
        const res = await getRobotInfo({ id })
        if (res) {
          if (res.data.h5_domain) {
            this.iframeBaseUrl = res.data.h5_domain
          }
          return res.data
        }
      } catch (error) {
        console.error('获取机器人信息失败:', error)
      }
    },
    /**
     * 刷新 iframe（强制重新加载）
     */
    refreshIframe() {
      this.iframeKey += 1
      this.iframeLoading = true
    },
    /**
     * iframe 加载完成处理
     */
    handleIframeLoad() {
      this.iframeLoading = false
    }
  },
  persist: {
    key: 'workbench-store',
    paths: ['currentRobotKey', 'iframeBaseUrl', 'sidebarCollapsed']
  }
})

export const useWorkbenchStoreWithOut = () => {
  return useWorkbenchStore(store)
}
