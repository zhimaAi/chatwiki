import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import { flushPromises } from '@vue/test-utils'
const api = vi.hoisted(() => ({ list: vi.fn(), info: vi.fn() }))
vi.mock('@/api/robot', () => ({ getRobotList: api.list, getRobotInfo: api.info, deleteRobot: vi.fn(), robotCopy: vi.fn(), saveRobot: vi.fn() }))
vi.mock('@/api/clawbot', () => ({ saveClawbotConf: vi.fn() }))
vi.mock('@/utils/permission', () => ({ checkRole: () => true, getRobotPermission: () => 1 }))
vi.mock('@/stores/modules/permission', () => ({ usePermissionStore: () => ({ checkPermission: vi.fn().mockResolvedValue() }) }))
vi.mock('ant-design-vue', () => ({ message: { error: vi.fn() } }))
import { useClawbotStore } from './clawbot'

describe('Agent detail request ownership', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
    api.list.mockResolvedValue({ data: [{ id: '1', robot_key: 'one' }, { id: '2', robot_key: 'two' }] })
    api.info.mockImplementation(async ({ id }) => ({ data: { id, robot_key: id === '1' ? 'one' : 'two' } }))
  })
  it('waits for the only detail request when switching an already initialized Agent', async () => {
    const store = useClawbotStore()
    await store.initModule({ targetId: '1', targetRobotKey: 'one', autoCreate: false })
    let finish
    api.info.mockImplementationOnce(() => new Promise(resolve => { finish = resolve }))
    const switching = store.initModule({ targetId: '2', targetRobotKey: 'two', autoCreate: false, forceRefresh: true })
    await flushPromises()
    expect(api.info).toHaveBeenCalledTimes(2)
    finish({ data: { id: '2', robot_key: 'two' } })
    await switching
    expect(store.robotInfo.id).toBe('2')
    expect(store.isInitializing).toBe(false)
    expect(api.info).toHaveBeenCalledTimes(2)
  })
  it('does not accept an old Agent response after switching back', async () => {
    const store = useClawbotStore()
    await store.initModule({ targetId: '1', targetRobotKey: 'one', autoCreate: false })
    let finishOld
    api.info.mockImplementationOnce(() => new Promise(resolve => { finishOld = resolve }))
    store.currentAssistantId = '2'
    await nextTick()
    store.currentAssistantId = '1'
    await nextTick()
    await nextTick()
    finishOld({ data: { id: '2', robot_key: 'two' } })
    await nextTick()
    await nextTick()
    expect(store.robotInfo.id).toBe('1')
  })
})
