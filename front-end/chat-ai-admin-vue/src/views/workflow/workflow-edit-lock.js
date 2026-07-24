export const WORKFLOW_AUTO_SAVE_MS = 60 * 1000
export const WORKFLOW_LOCK_HEARTBEAT_MS = 30 * 1000
export const WORKFLOW_LOCK_LEASE_MS = 120 * 1000
export const WORKFLOW_LOCK_RETRY_MS = 5 * 1000

const WINDOW_ID_STORAGE_KEY = 'wf_window_identifier'

const randomIdentifier = () => {
  if (globalThis.crypto?.randomUUID) {
    return globalThis.crypto.randomUUID()
  }
  return `${Date.now()}_${Math.random().toString(36).slice(2, 12)}`
}

export const getWorkflowWindowIdentifier = (storage = globalThis.sessionStorage) => {
  try {
    let identifier = storage.getItem(WINDOW_ID_STORAGE_KEY)
    if (!identifier) {
      identifier = randomIdentifier()
      storage.setItem(WINDOW_ID_STORAGE_KEY, identifier)
    }
    return identifier
  } catch (e) {
    return randomIdentifier()
  }
}

export const createWorkflowLeaseToken = () => randomIdentifier()

export const createWorkflowHeartbeatController = ({
  sendHeartbeat,
  onLockLost,
  now = () => Date.now(),
  setTimer = (callback, delay) => setTimeout(callback, delay),
  clearTimer = (timer) => clearTimeout(timer),
  heartbeatMs = WORKFLOW_LOCK_HEARTBEAT_MS,
  retryMs = WORKFLOW_LOCK_RETRY_MS,
  leaseMs = WORKFLOW_LOCK_LEASE_MS
}) => {
  let timer = null
  let running = false
  let inFlight = false
  let leaseDeadline = 0

  const updateLeaseDeadline = (data = {}) => {
    const ttlSeconds = Number(data.lock_ttl)
    const confirmedLeaseMs = ttlSeconds > 0 ? ttlSeconds * 1000 : leaseMs
    leaseDeadline = now() + confirmedLeaseMs
  }

  const stop = () => {
    running = false
    if (timer) {
      clearTimer(timer)
      timer = null
    }
  }

  const loseLock = (data = {}) => {
    if (!running) return
    stop()
    onLockLost(data)
  }

  const schedule = (delay) => {
    if (!running) return
    if (timer) clearTimer(timer)
    timer = setTimer(runHeartbeat, delay)
  }

  const runHeartbeat = async () => {
    timer = null
    if (!running || inFlight) return
    if (leaseDeadline > 0 && now() >= leaseDeadline) {
      loseLock({ lock_conflict: 1 })
      return
    }

    inFlight = true
    try {
      const response = await sendHeartbeat()
      if (!running) return
      const data = response?.data || {}
      if (!data.lock_res) {
        loseLock(data)
        return
      }
      updateLeaseDeadline(data)
      schedule(heartbeatMs)
    } catch (e) {
      if (!running) return
      if (leaseDeadline > 0 && now() >= leaseDeadline) {
        loseLock(e?.data || { lock_conflict: 1 })
        return
      }
      schedule(retryMs)
    } finally {
      inFlight = false
    }
  }

  const start = (lockData = {}) => {
    stop()
    running = true
    updateLeaseDeadline(lockData)
    schedule(heartbeatMs)
  }

  const trigger = () => {
    if (!running || inFlight) return
    if (timer) {
      clearTimer(timer)
      timer = null
    }
    void runHeartbeat()
  }

  return {
    start,
    stop,
    trigger,
    isRunning: () => running,
    getLeaseDeadline: () => leaseDeadline
  }
}
