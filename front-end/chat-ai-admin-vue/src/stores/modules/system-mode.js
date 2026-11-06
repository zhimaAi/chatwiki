import { defineStore } from 'pinia'
import { getSystemMode, saveSystemMode } from '@/api/user'

export const SYSTEM_MODES = Object.freeze({
  AI_SERVICE: 'ai_service',
  AGENT: 'agent',
  ADVANCED: 'advanced'
})

const isSystemMode = (mode) => Object.values(SYSTEM_MODES).includes(mode)

// 首页和布局导航会同时初始化企业模式，复用同一个请求避免旧响应覆盖。
let inflightLoad = null

export const useSystemModeStore = defineStore('system-mode', {
  state: () => ({
    currentMode: SYSTEM_MODES.ADVANCED,
    modeSelected: true,
    scopeKey: '',
    storageStatus: 'idle',
    storageError: ''
  }),
  getters: {
    needsModeSelection: state => state.storageStatus === 'ready' && !state.modeSelected
  },
  actions: {
    async loadForScope({ companyId }) {
      if (!companyId) {
        this.currentMode = SYSTEM_MODES.ADVANCED
        this.modeSelected = true
        this.scopeKey = ''
        this.storageStatus = 'error'
        this.storageError = 'error_scope'
        return false
      }

      const scopeKey = String(companyId)
      if (inflightLoad?.key === scopeKey) return inflightLoad.promise

      const promise = this.performLoad(scopeKey)
      inflightLoad = { key: scopeKey, promise }
      try {
        return await promise
      } finally {
        if (inflightLoad?.promise === promise) inflightLoad = null
      }
    },
    async performLoad(scopeKey) {
      this.currentMode = SYSTEM_MODES.ADVANCED
      this.modeSelected = true
      this.scopeKey = scopeKey
      this.storageStatus = 'loading'
      this.storageError = ''

      try {
        const res = await getSystemMode()
        if (this.scopeKey !== scopeKey) return false
        this.currentMode = isSystemMode(res?.data?.system_mode)
          ? res.data.system_mode
          : SYSTEM_MODES.ADVANCED
        // 旧服务没有首次选择字段时保持原入口，不能把全部老企业视为新注册企业。
        this.modeSelected = ![0, '0'].includes(res?.data?.system_mode_selected)
        this.storageStatus = 'ready'
        return true
      } catch {
        if (this.scopeKey !== scopeKey) return false
        this.storageStatus = 'error'
        this.storageError = 'error_storage_read'
        return false
      }
    },
    async saveMode(mode) {
      if (!isSystemMode(mode)) {
        throw new Error('Invalid system mode')
      }
      if (this.storageStatus !== 'ready' || !this.scopeKey) {
        throw new Error('System mode storage is not ready')
      }

      const scopeKey = this.scopeKey
      const response = await saveSystemMode({ system_mode: mode })
      if (response?.res !== undefined && response.res !== 0) throw new Error(response.msg || 'Save system mode failed')
      if (this.scopeKey !== scopeKey) {
        throw new Error('System mode scope changed')
      }
      this.currentMode = mode
      this.modeSelected = true
    }
  }
})
