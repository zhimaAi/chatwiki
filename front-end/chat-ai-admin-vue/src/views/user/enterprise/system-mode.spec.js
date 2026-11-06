import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const api = vi.hoisted(() => ({
  getSystemMode: vi.fn(),
  saveSystemMode: vi.fn()
}))

vi.mock('@/api/user', () => api)

import { SYSTEM_MODES, useSystemModeStore } from '@/stores/modules/system-mode'

describe('system mode store', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('loads the enterprise mode from the server', async () => {
    api.getSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.AI_SERVICE } })
    const store = useSystemModeStore()

    await expect(store.loadForScope({ companyId: 1 })).resolves.toBe(true)
    expect(store.currentMode).toBe(SYSTEM_MODES.AI_SERVICE)
    expect(store.storageStatus).toBe('ready')
  })

  it('keeps advanced for old enterprises', async () => {
    api.getSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.ADVANCED } })
    const store = useSystemModeStore()

    await store.loadForScope({ companyId: 1 })
    expect(store.currentMode).toBe(SYSTEM_MODES.ADVANCED)
    expect(store.needsModeSelection).toBe(false)
  })

  it.each([0, '0'])('requires the first selection only when the server explicitly returns %s', async selected => {
    api.getSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.AI_SERVICE, system_mode_selected: selected } })
    const store = useSystemModeStore()
    await store.loadForScope({ companyId: 1 })
    expect(store.needsModeSelection).toBe(true)
    api.saveSystemMode.mockResolvedValue({ res: 0 })
    await store.saveMode(SYSTEM_MODES.AI_SERVICE)
    expect(store.needsModeSelection).toBe(false)
    expect(api.saveSystemMode).toHaveBeenCalledWith({ system_mode: SYSTEM_MODES.AI_SERVICE })
  })

  it('keeps the onboarding requirement after a failed save', async () => {
    api.getSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.AI_SERVICE, system_mode_selected: 0 } })
    api.saveSystemMode.mockResolvedValue({ res: 1, msg: 'denied' })
    const store = useSystemModeStore()
    await store.loadForScope({ companyId: 1 })
    await expect(store.saveMode(SYSTEM_MODES.AGENT)).rejects.toThrow('denied')
    expect(store.currentMode).toBe(SYSTEM_MODES.AI_SERVICE)
    expect(store.needsModeSelection).toBe(true)
  })

  it('saves to the server before applying the new mode', async () => {
    api.getSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.ADVANCED } })
    api.saveSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.AI_SERVICE } })
    const store = useSystemModeStore()
    await store.loadForScope({ companyId: 1 })

    await store.saveMode(SYSTEM_MODES.AI_SERVICE)
    expect(api.saveSystemMode).toHaveBeenCalledWith({ system_mode: SYSTEM_MODES.AI_SERVICE })
    expect(store.currentMode).toBe(SYSTEM_MODES.AI_SERVICE)
  })

  it('keeps the current mode when saving fails', async () => {
    api.getSystemMode.mockResolvedValue({ data: { system_mode: SYSTEM_MODES.ADVANCED } })
    api.saveSystemMode.mockRejectedValue(new Error('network'))
    const store = useSystemModeStore()
    await store.loadForScope({ companyId: 1 })

    await expect(store.saveMode(SYSTEM_MODES.AGENT)).rejects.toThrow('network')
    expect(store.currentMode).toBe(SYSTEM_MODES.ADVANCED)
  })

  it('coalesces concurrent loads for the same enterprise', async () => {
    let resolveRequest
    api.getSystemMode.mockImplementation(
      () =>
        new Promise((resolve) => {
          resolveRequest = resolve
        })
    )
    const store = useSystemModeStore()

    const first = store.loadForScope({ companyId: 1 })
    const second = store.loadForScope({ companyId: 1 })
    expect(api.getSystemMode).toHaveBeenCalledTimes(1)

    resolveRequest({ data: { system_mode: SYSTEM_MODES.AGENT } })
    await expect(Promise.all([first, second])).resolves.toEqual([true, true])
    expect(store.currentMode).toBe(SYSTEM_MODES.AGENT)
  })

  it('ignores a stale response after the enterprise changes', async () => {
    let resolveFirst
    let resolveSecond
    api.getSystemMode
      .mockImplementationOnce(
        () =>
          new Promise((resolve) => {
            resolveFirst = resolve
          })
      )
      .mockImplementationOnce(
        () =>
          new Promise((resolve) => {
            resolveSecond = resolve
          })
      )
    const store = useSystemModeStore()

    const first = store.loadForScope({ companyId: 1 })
    const second = store.loadForScope({ companyId: 2 })
    resolveSecond({ data: { system_mode: SYSTEM_MODES.AGENT } })
    await second
    resolveFirst({ data: { system_mode: SYSTEM_MODES.AI_SERVICE } })
    await expect(first).resolves.toBe(false)
    expect(store.currentMode).toBe(SYSTEM_MODES.AGENT)
    expect(store.scopeKey).toBe('2')
  })

  it('reports a server read failure and keeps advanced mode', async () => {
    api.getSystemMode.mockRejectedValue(new Error('network'))
    const store = useSystemModeStore()

    await expect(store.loadForScope({ companyId: 1 })).resolves.toBe(false)
    expect(store.storageStatus).toBe('error')
    expect(store.storageError).toBe('error_storage_read')
    expect(store.currentMode).toBe(SYSTEM_MODES.ADVANCED)
  })
})
