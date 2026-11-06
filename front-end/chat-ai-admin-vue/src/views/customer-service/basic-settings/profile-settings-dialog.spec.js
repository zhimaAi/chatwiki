import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { reactive, ref } from 'vue'
import BasicInfo from '@/views/robot/robot-config/basic-config/components/basic-config.vue'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_ROBOT_AVATAR_LIST, DEFAULT_WORKFLOW_AVATAR } from '@/constants'

const mocks = vi.hoisted(() => ({ save: vi.fn(), read: vi.fn(), error: vi.fn() }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/api/robot', () => ({ saveRobot: mocks.save }))
vi.mock('@/utils', () => ({ getBase64: mocks.read }))
vi.mock('@/stores/modules/robot', () => ({ useRobotStore: () => ({ robotList: [], robotGroupList: [] }) }))
vi.mock('ant-design-vue', () => ({ message: { error: mocks.error } }))
const container = { template: '<div><slot /></div>' }
const stubs = {
  'svg-icon': true, 'a-tooltip': container, 'a-flex': container,
  'a-button': { template: '<button><slot /></button>' },
  'a-modal': { name: 'AModal', props: ['open', 'confirmLoading', 'closable', 'okButtonProps'], emits: ['ok', 'cancel'], template: '<section v-if="open" role="dialog"><slot /></section>' },
  'a-input': { name: 'AInput', props: ['value', 'status'], emits: ['update:value'], template: '<input :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-textarea': { name: 'ATextarea', props: ['value'], emits: ['update:value'], template: '<textarea :value="value" @input="$emit(\'update:value\', $event.target.value)" />' },
  'a-alert': { name: 'AAlert', template: '<div><slot name="action" /></div>' }
}
function render(inlineSummary = true, applicationType = 2) {
  const robotInfo = reactive({ id: 7, robot_key: 'seven', application_type: applicationType, robot_name: '原名称', robot_intro: '原简介', robot_avatar_url: '/custom.png' })
  const legacySave = vi.fn()
  const state = { robotInfo, savedRevision: ref(0), saveAction: vi.fn(async task => {
    await task()
    const payload = mocks.save.mock.calls.at(-1)[0]
    Object.assign(robotInfo, { robot_name: payload.robot_name, robot_intro: payload.robot_intro, robot_avatar_url: payload.robot_avatar_url || '/uploaded.png' })
    state.savedRevision.value++
    return true
  }) }
  const wrapper = mount(BasicInfo, { props: { inlineSummary }, global: { provide: { customerRobotContext: state, robotInfo: { robotInfo, updateRobotInfo: legacySave } }, stubs } })
  return { wrapper, state, legacySave }
}
const named = (wrapper, name) => wrapper.findComponent({ name })
const open = wrapper => wrapper.find('.edit-box-header button').trigger('click')
const confirm = async wrapper => { named(wrapper, 'AModal').vm.$emit('ok'); await flushPromises() }
const cancel = async wrapper => { named(wrapper, 'AModal').vm.$emit('cancel'); await flushPromises() }
const deferred = () => { let resolve; const promise = new Promise(done => { resolve = done }); return { promise, resolve } }
async function upload(wrapper, file) {
  const input = wrapper.find('input[type="file"]')
  Object.defineProperty(input.element, 'files', { configurable: true, value: [file] })
  await input.trigger('change')
  await flushPromises()
}

describe('workspace profile settings dialog', () => {
  beforeEach(() => { vi.clearAllMocks(); mocks.save.mockResolvedValue({ res: 0 }); mocks.read.mockResolvedValue('data:image/png;base64,avatar') })

  it('opens from either trigger and discards private name/intro/avatar drafts on cancel', async () => {
    const { wrapper, state } = render()
    await open(wrapper)
    await named(wrapper, 'AInput').find('input').setValue('新名称')
    await named(wrapper, 'ATextarea').find('textarea').setValue('新简介')
    expect(named(wrapper, 'AInput').props('value')).toBe('新名称')
    await wrapper.findAll('.default-avatar')[3].trigger('click')
    expect(wrapper.find('.profile-copy').text()).toContain('原名称')
    expect(state.robotInfo.robot_avatar_url).toBe('/custom.png')
    await cancel(wrapper)
    expect(state.saveAction).not.toHaveBeenCalled()
    await wrapper.find('.profile-avatar').trigger('click')
    expect(named(wrapper, 'AInput').props('value')).toBe('原名称')
    expect(named(wrapper, 'ATextarea').props('value')).toBe('原简介')
    expect(wrapper.find('.avatar-button img').attributes('src')).toBe('/custom.png')
    expect(wrapper.find('.robot-form-box').exists()).toBe(false)
    wrapper.unmount()
  })

  it('preserves all six defaults, highlights the selected one and saves through the existing application type', async () => {
    const { wrapper, state } = render(true, 1)
    await open(wrapper)
    const avatars = [DEFAULT_WORKFLOW_AVATAR, ...DEFAULT_ROBOT_AVATAR_LIST]
    expect(wrapper.findAll('.default-avatar img').map(image => image.attributes('src'))).toEqual(avatars)
    for (let index = 0; index < avatars.length; index++) {
      await wrapper.findAll('.default-avatar')[index].trigger('click')
      expect(wrapper.find('.avatar-button img').attributes('src')).toBe(avatars[index])
      expect(wrapper.findAll('.default-avatar.active')).toHaveLength(1)
      expect(wrapper.findAll('.default-avatar')[index].attributes('aria-pressed')).toBe('true')
    }
    expect(state.robotInfo.robot_avatar_url).toBe('/custom.png')
    await confirm(wrapper)
    expect(mocks.save).toHaveBeenCalledWith(expect.objectContaining({ robot_avatar_url: avatars.at(-1), id: 7, robot_key: 'seven' }), 1)
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    expect(wrapper.find('.profile-avatar img').attributes('src')).toBe(avatars.at(-1))
    wrapper.unmount()
  })

  it('keeps failed edits for retry, guards duplicate saves and prevents closing during a pending save', async () => {
    const { wrapper, state } = render()
    state.saveAction.mockResolvedValueOnce(false)
    await open(wrapper)
    await named(wrapper, 'AInput').find('input').setValue('保存名称')
    await confirm(wrapper)
    expect(named(wrapper, 'AAlert').exists()).toBe(true)
    expect(state.robotInfo.robot_name).toBe('原名称')
    const pending = deferred()
    mocks.save.mockReturnValueOnce(pending.promise)
    await named(wrapper, 'AAlert').find('button').trigger('click')
    expect(named(wrapper, 'AModal').props('confirmLoading')).toBe(true)
    await cancel(wrapper)
    await confirm(wrapper)
    expect(named(wrapper, 'AModal').exists()).toBe(true)
    expect(state.saveAction).toHaveBeenCalledTimes(2)
    pending.resolve({ res: 0 })
    await flushPromises()
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    expect(wrapper.find('.profile-copy').text()).toContain('保存名称')
    wrapper.unmount()
  })

  it('validates names, avatar formats and the existing 100 KB upload limit', async () => {
    const { wrapper } = render()
    await open(wrapper)
    await named(wrapper, 'AInput').find('input').setValue('  ')
    await confirm(wrapper)
    expect(named(wrapper, 'AInput').props('status')).toBe('error')
    expect(mocks.save).not.toHaveBeenCalled()
    await upload(wrapper, new File(['gif'], 'avatar.gif', { type: 'image/gif' }))
    await upload(wrapper, new File([new Uint8Array(102401)], 'avatar.png', { type: 'image/png' }))
    expect(mocks.error).toHaveBeenNthCalledWith(1, 'profile_avatar_type_error')
    expect(mocks.error).toHaveBeenNthCalledWith(2, 'msg_avatar_size_limit')
    expect(mocks.read).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('preserves uploaded File objects and excludes the local preview URL from the request', async () => {
    const { wrapper } = render()
    await open(wrapper)
    const file = new File(['image'], 'avatar.webp', { type: 'image/webp' })
    await upload(wrapper, file)
    expect(wrapper.find('.avatar-button img').attributes('src')).toBe('data:image/png;base64,avatar')
    expect(wrapper.findAll('.default-avatar.active')).toHaveLength(0)
    await confirm(wrapper)
    const [payload, type] = mocks.save.mock.calls[0]
    expect(payload.robot_avatar).toBe(file)
    expect(payload).not.toHaveProperty('robot_avatar_url')
    expect(type).toBe(2)
    wrapper.unmount()
  })

  it('ignores late upload previews after selecting a default or cancelling and reopening', async () => {
    const { wrapper } = render()
    await open(wrapper)
    const pending = deferred()
    mocks.read.mockReturnValueOnce(pending.promise)
    await upload(wrapper, new File(['image'], 'avatar.png', { type: 'image/png' }))
    expect(named(wrapper, 'AModal').props('okButtonProps').disabled).toBe(true)
    await wrapper.findAll('.default-avatar')[0].trigger('click')
    pending.resolve('stale-preview')
    await flushPromises()
    expect(wrapper.find('.avatar-button img').attributes('src')).toBe(DEFAULT_ROBOT_AVATAR)
    const next = deferred()
    mocks.read.mockReturnValueOnce(next.promise)
    await upload(wrapper, new File(['image'], 'avatar.png', { type: 'image/png' }))
    await cancel(wrapper)
    await open(wrapper)
    next.resolve('another-stale-preview')
    await flushPromises()
    expect(wrapper.find('.avatar-button img').attributes('src')).toBe('/custom.png')
    wrapper.unmount()
  })

  it('keeps the legacy inline editor and its save handler', async () => {
    const { wrapper, legacySave } = render(false)
    await open(wrapper)
    expect(wrapper.find('.robot-form-box').isVisible()).toBe(true)
    expect(named(wrapper, 'AModal').exists()).toBe(false)
    await wrapper.find('.edit-box-header button').trigger('click')
    expect(legacySave).toHaveBeenCalledWith(expect.objectContaining({ robot_name: '原名称', robot_intro: '原简介' }))
    wrapper.unmount()
  })
})
