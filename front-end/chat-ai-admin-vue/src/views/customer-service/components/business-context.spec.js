import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { h, inject, reactive } from 'vue'
const mocks = vi.hoisted(() => ({ customer: null, robot: null, save: vi.fn(), info: vi.fn() }))
vi.mock('@/api/robot', () => ({ saveRobot: mocks.save, getRobotInfo: mocks.info }))
vi.mock('@/stores/modules/customer-service', () => ({ useCustomerServiceStore: () => mocks.customer }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => mocks.robot }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('ant-design-vue', () => ({ message: { error: vi.fn() } }))
import BusinessContext from './business-context.vue'

const robot = name => ({ id: 7, robot_key: 'seven', application_type: 0, robot_name: name, welcomes: '{"content":"","question":[]}' })
const deferred = () => { let resolve; const promise = new Promise(done => { resolve = done }); return { promise, resolve } }
let context
const Probe = { setup() { context = inject('customerRobotContext'); return () => null } }
const render = () => mount(BusinessContext, { props: { robot: robot('initial') }, slots: { default: () => h(Probe) } })

describe('customer robot configuration adapter', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.customer = { robot: robot('initial'), syncRobot: vi.fn() }
    mocks.robot = { robotInfo: reactive({}), setRobotInfo: vi.fn(data => Object.assign(mocks.robot.robotInfo, data)) }
  })
  it('saves each snapshot in order and keeps newer form edits while an earlier detail refresh completes', async () => {
    const first = deferred()
    const second = deferred()
    mocks.save.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise)
    mocks.info.mockResolvedValueOnce({ data: robot('first') }).mockResolvedValueOnce({ data: robot('second') })
    const wrapper = render()
    const saveOne = context.updateRobotInfo({ robot_name: 'first' })
    const saveTwo = context.updateRobotInfo({ robot_name: 'second' })
    await flushPromises()
    expect(mocks.save).toHaveBeenCalledTimes(1)
    expect(mocks.save.mock.calls[0][0].robot_name).toBe('first')
    first.resolve({ res: 0 })
    await saveOne
    await flushPromises()
    expect(mocks.save.mock.calls[1][0].robot_name).toBe('second')
    expect(context.robotInfo.robot_name).toBe('second')
    second.resolve({ res: 0 })
    await saveTwo
    await context.waitForSave()
    expect(context.robotInfo.robot_name).toBe('second')
    expect(mocks.customer.syncRobot).toHaveBeenLastCalledWith(robot('second'))
    wrapper.unmount()
  })
  it('blocks testing after failure and prevents a departed robot response from changing the new robot', async () => {
    mocks.save.mockRejectedValueOnce(new Error('save failed'))
    const wrapper = render()
    await expect(context.updateRobotInfo({ robot_name: 'failed' })).resolves.toBe(false)
    await expect(context.waitForSave()).rejects.toThrow('save failed')
    const pending = deferred()
    mocks.save.mockReturnValueOnce(pending.promise)
    mocks.info.mockResolvedValueOnce({ data: robot('late') })
    const save = context.updateRobotInfo({ robot_name: 'retry' })
    await flushPromises()
    wrapper.unmount()
    mocks.customer.robot = { id: 8, robot_key: 'eight' }
    Object.assign(mocks.robot.robotInfo, { id: 8, robot_key: 'eight', robot_name: 'new robot' })
    pending.resolve({ res: 0 })
    await save
    expect(mocks.robot.robotInfo.robot_name).toBe('new robot')
    expect(mocks.customer.syncRobot).not.toHaveBeenCalled()
  })
})
