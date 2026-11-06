import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount, shallowMount } from '@vue/test-utils'
import { reactive } from 'vue'
import { message, Modal } from 'ant-design-vue'
import Tooltip from 'ant-design-vue/es/tooltip'

const api = vi.hoisted(() => ({ listLibraries: vi.fn(), upload: vi.fn(), ensure: vi.fn(), relate: vi.fn(), listRobots: vi.fn(), robotInfo: vi.fn(), copy: vi.fn(), remove: vi.fn(), libraryInfo: vi.fn(), workflowRelation: vi.fn(), setRate: vi.fn(), openCreator: vi.fn(), openImporter: vi.fn(), openExport: vi.fn(), exportUrl: vi.fn() }))
vi.mock('@/api/library', () => ({ getLibraryList: api.listLibraries, addLibraryFile: api.upload }))
vi.mock('@/api/robot', () => ({ ensureRobotNormalLibrary: api.ensure, relationLibrary: api.relate, getRobotList: api.listRobots, getRobotInfo: api.robotInfo, robotCopy: api.copy, deleteRobot: api.remove, relationWorkFlow: api.workflowRelation }))
vi.mock('@/stores/modules/library', () => ({ useLibraryStore: () => ({ getLibraryInfo: api.libraryInfo }) }))
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: (key, values) => {
  if (key === 'library_document_count_tip') return `${values.count} 个文档`
  if (key === 'library_storage_size_tip') return `总大小${values.size}`
  if (key === 'library_related_apps_tip') return `已关联 ${values.count} 个`
  return key
} }) }))
vi.mock('@/hooks/web/useOpenUrlWithToken', () => ({ useOpenUrlWithToken: () => ({ openUrlWithToken: api.exportUrl }) }))
vi.mock('vue-router', () => ({ useRouter: () => ({ resolve: location => ({ href: JSON.stringify(location) }) }) }))
vi.mock('ant-design-vue', () => ({ message: { success: vi.fn(), error: vi.fn() }, Modal: { confirm: vi.fn() } }))
vi.mock('@/views/library/library-details/qa-knowledge-document/index.vue', () => ({ default: { props: ['library_id'], template: '<div />' } }))
vi.mock('@/views/library/library-details/knowledge-document.vue', () => ({ default: { props: ['library_id'], template: '<div />' } }))
vi.mock('@/views/library/library-details/knowledge-config.vue', () => ({ default: { props: { libraryId: String, embedded: Boolean }, template: '<div />' } }))
vi.mock('@/views/library/add-library/components/upload-input.vue', () => ({ default: { props: ['value'], emits: ['update:value'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/basic-config/components/associated-knowledge-base/library-select-alert.vue', () => ({ default: { template: '<div />' } }))
vi.mock('@/views/robot/robot-list/components/add-robot-alert.vue', () => ({ default: { props: ['createRedirect'], emits: ['ok', 'editRobot'], setup: (_, { expose }) => { expose({ open: api.openCreator }) }, template: '<div />' } }))
vi.mock('@/views/robot/robot-list/components/import-csl-alert.vue', () => ({ default: { props: ['applicationType'], emits: ['ok'], setup: (_, { expose }) => { expose({ open: api.openImporter }) }, template: '<div />' } }))
vi.mock('@/views/robot/robot-list/components/workflow-export-modal.vue', () => ({ default: { emits: ['ok'], setup: (_, { expose }) => { expose({ handle: api.openExport }) }, template: '<div />' } }))
vi.mock('@/views/robot/robot-config/rate-limit/index.vue', () => ({ default: { props: ['robotKey'], emits: ['status-change'], setup: (_, { expose }) => { expose({ setEnabled: api.setRate }) }, template: '<div />' } }))
vi.mock('@/views/robot/robot-config/invoke-logs/index.vue', () => ({ default: { props: ['robotId', 'robotKey'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/statistical_analysis/hit-statics.vue', () => ({ default: { props: ['robotId', 'robotKey', 'filters'], emits: ['unknown'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/qa-feedback/index.vue', () => ({ default: { props: ['robotId'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/unknown_issue/index.vue', () => ({ default: { props: ['robotId', 'filters'], template: '<div />' } }))
vi.mock('@/views/robot/robot-config/unknown_issue/summarize/index.vue', () => ({ default: { props: ['robotId', 'filters'], template: '<div />' } }))

import Knowledge from './knowledge-base/index.vue'
import Functions from './function-center/index.vue'
import Analytics from './analytics/index.vue'
import CustomerPageHeader from './components/page-header.vue'
import UploadInput from '@/views/library/add-library/components/upload-input.vue'
import QaDocument from '@/views/library/library-details/qa-knowledge-document/index.vue'
import NormalDocument from '@/views/library/library-details/knowledge-document.vue'
import ImportCsl from '@/views/robot/robot-list/components/import-csl-alert.vue'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'
import RateLimit from '@/views/robot/robot-config/rate-limit/index.vue'
import ExportModal from '@/views/robot/robot-list/components/workflow-export-modal.vue'
import InvokeLogs from '@/views/robot/robot-config/invoke-logs/index.vue'
import HitStatistics from '@/views/robot/robot-config/statistical_analysis/hit-statics.vue'
import UnknownStats from '@/views/robot/robot-config/unknown_issue/index.vue'

const container = { template: '<div><slot /><slot name="title" /><slot name="extra" /><slot name="overlay" /></div>' }
const options = context => ({ global: { provide: { customerRobotContext: context }, stubs: { CustomerPageHeader: false, AddRobotAlert: false, ImportCsl: false, ExportModal: false, RateLimit: false, 'a-radio-group': { ...container, name: 'ARadioGroup' }, 'a-radio-button': container, 'a-button': { template: '<button><slot name="icon" /><slot /></button>' }, 'a-spin': container, 'a-card': container, 'a-tooltip': { name: 'ATooltip', props: ['title', 'placement', 'arrow', 'overlayClassName'], template: '<span><slot /></span>' }, 'a-dropdown': { ...container, name: 'ADropdown', props: ['open'], emits: ['openChange'] }, 'a-menu': { ...container, name: 'AMenu' }, 'a-menu-item': { ...container, name: 'AMenuItem' }, 'a-modal': { name: 'AModal', props: ['open', 'title'], emits: ['cancel'], template: '<div v-if="open"><slot /></div>' }, 'a-switch': { name: 'ASwitch', props: ['checked', 'disabled', 'loading'], emits: ['change'], template: '<button role="switch" :disabled="disabled" @click="$emit(\'change\', !checked)" />' } } } })
function context(robotInfo = {}) { return { robotInfo: reactive({ id: 7, robot_key: 'customer-seven', default_library_id: 101, normal_library_id: 0, library_ids: '101', work_flow_ids: '', ...robotInfo }), waitForSave: vi.fn().mockResolvedValue(), refreshRobot: vi.fn().mockResolvedValue(), saveAction: vi.fn(async task => { await task(); return true }) } }
async function changeTab(wrapper, tab) {
  const header = wrapper.findComponent(CustomerPageHeader)
  const index = header.props('tabs').findIndex(item => item.key === tab)
  await header.findAll('[role="tab"]')[index].trigger('click')
  await flushPromises()
}
const button = (wrapper, text) => wrapper.findAll('button').find(item => item.text() === text)

describe('customer workspace business pages', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.listLibraries.mockResolvedValue({ data: [{ id: 101, type: 2 }, { id: 202, type: 0 }] })
    api.ensure.mockResolvedValue({ data: { library_id: 202, created: true } })
    api.upload.mockResolvedValue({ res: 0 })
    api.listRobots.mockResolvedValue({ data: [{ id: 900, robot_key: 'workflow-nine', robot_name: 'Flow', application_type: 1 }] })
    api.robotInfo.mockResolvedValue({ data: { id: 900, robot_key: 'workflow-nine', robot_name: 'Flow', application_type: 1, en_name: 'workflow_nine' } })
  })

  it('browsing and cancelling do not create a normal library; details receive library IDs', async () => {
    const wrapper = shallowMount(Knowledge, options(context()))
    await flushPromises()
    expect(wrapper.findComponent(QaDocument).props('library_id')).toBe('101')
    await changeTab(wrapper, 'normal')
    wrapper.findComponent(UploadInput).vm.$emit('update:value', [new File(['hi'], 'a.txt')])
    await flushPromises()
    await button(wrapper, 'cancel').trigger('click')
    expect(api.ensure).not.toHaveBeenCalled()
    expect(api.upload).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('creates on submit, ignores a double click and reuses the bound library after upload failure', async () => {
    const current = context()
    current.refreshRobot.mockImplementation(async () => { current.robotInfo.normal_library_id = 202; current.robotInfo.library_ids = '101,202' })
    let rejectUpload
    api.upload.mockImplementationOnce(() => new Promise((_, reject) => { rejectUpload = reject }))
    const wrapper = shallowMount(Knowledge, options(current))
    await flushPromises()
    await changeTab(wrapper, 'normal')
    wrapper.findComponent(UploadInput).vm.$emit('update:value', [new File(['hi'], 'a.txt')])
    await flushPromises()
    await button(wrapper, 'upload_documents').trigger('click')
    await button(wrapper, 'upload_documents').trigger('click')
    expect(api.ensure).toHaveBeenCalledTimes(1)
    expect(api.ensure).toHaveBeenCalledWith(7)
    rejectUpload(new Error('upload failed'))
    await flushPromises()
    api.ensure.mockResolvedValue({ data: { library_id: 202, created: false } })
    await button(wrapper, 'upload_documents').trigger('click')
    await flushPromises()
    expect(api.upload.mock.calls[1][0].get('library_id')).toBe('202')
    expect(current.robotInfo.default_library_id).toBe(101)
    expect(wrapper.findComponent(NormalDocument).props('library_id')).toBe('202')
    wrapper.unmount()
  })

  describe('other knowledge base cards', () => {
    beforeEach(() => {
      api.listLibraries.mockResolvedValue({ data: [
        { id: 101, type: 2, library_name: 'Built-in QA' },
        { id: 202, type: 0, library_name: 'Built-in documents' },
        { id: '303', type: '2', library_name: 'Other QA', file_total: '3', file_size: '2048', robot_nums: '2' },
        { id: 404, type: 0, library_name: 'Other documents' },
        { id: 505, type: 1, library_name: 'Public documents' }
      ] })
      api.relate.mockReset().mockResolvedValue({ res: 0 })
    })

    it('shows all other libraries with real metrics and association status, excluding both built-in libraries', async () => {
      const wrapper = shallowMount(Knowledge, options(context({ normal_library_id: '202', library_ids: '101,202,303' })))
      await flushPromises()
      await changeTab(wrapper, 'other')
      const cards = wrapper.findAll('.library-card')
      expect(cards.map(card => card.find('.library-name').text())).toEqual(['Other QA', 'Other documents', 'Public documents'])
      expect(cards[0].classes()).toContain('linked')
      expect(cards[0].find('.linked-badge').text()).toBe('library_linked')
      expect(cards[0].find('.library-metrics').text()).toBe('32KB2')
      expect(cards[1].find('.library-metrics').text()).toBe('00 B0')
      expect(cards[1].find('.linked-badge').exists()).toBe(false)
      expect(api.relate).not.toHaveBeenCalled()
      wrapper.unmount()
    })

    it('links a card once, preserving built-in and freshly loaded associations and refreshing status', async () => {
      const current = context({ normal_library_id: 202, library_ids: '101,202,303' })
      current.refreshRobot.mockImplementation(async () => {
        current.robotInfo.library_ids = api.relate.mock.calls.at(-1)?.[0].library_ids || '101,202,303,707'
      })
      let finish
      api.relate.mockImplementationOnce(() => new Promise(resolve => { finish = resolve }))
      const wrapper = shallowMount(Knowledge, options(current))
      await flushPromises()
      await changeTab(wrapper, 'other')
      await wrapper.findAll('.relation-button')[1].trigger('click')
      await flushPromises()
      expect(api.relate).toHaveBeenCalledTimes(1)
      expect(api.relate).toHaveBeenCalledWith({ id: 7, library_ids: '101,202,303,707,404' })
      expect(wrapper.findAll('.relation-button').every(button => button.attributes('disabled') !== undefined)).toBe(true)
      await wrapper.findAll('.relation-button')[1].trigger('click')
      expect(api.relate).toHaveBeenCalledTimes(1)
      expect(wrapper.findAll('.library-card')[1].classes()).not.toContain('linked')
      api.listLibraries.mockResolvedValueOnce({ data: [{ id: 404, type: 0, library_name: 'Other documents', robot_nums: '1' }] })
      finish({ res: 0 })
      await flushPromises()
      expect(wrapper.find('.library-card').classes()).toContain('linked')
      expect(wrapper.find('.relation-button').attributes('aria-label')).toBe('remove_relation')
      expect(wrapper.find('.relation-button').attributes('aria-pressed')).toBe('true')
      expect(wrapper.find('.library-metrics').text()).toBe('00 B1')
      wrapper.unmount()
    })

    it('requires confirmation to unlink and retains the saved status after failure, allowing retry', async () => {
      const current = context({ normal_library_id: 202, library_ids: '101,202,303,404' })
      current.refreshRobot.mockImplementation(async () => {
        if (api.relate.mock.calls.length > 1) current.robotInfo.library_ids = api.relate.mock.calls.at(-1)[0].library_ids
      })
      api.relate.mockRejectedValueOnce(new Error('denied'))
      const wrapper = shallowMount(Knowledge, options(current))
      await flushPromises()
      await changeTab(wrapper, 'other')
      await wrapper.find('.relation-button').trigger('click')
      expect(api.relate).not.toHaveBeenCalled()
      await Modal.confirm.mock.calls[0][0].onOk()
      await flushPromises()
      expect(message.error).toHaveBeenCalledWith('denied')
      expect(wrapper.find('.library-card').classes()).toContain('linked')
      await wrapper.find('.relation-button').trigger('click')
      await Modal.confirm.mock.calls[1][0].onOk()
      await flushPromises()
      expect(api.relate).toHaveBeenLastCalledWith({ id: 7, library_ids: '101,202,404' })
      expect(wrapper.find('.library-card').classes()).not.toContain('linked')
      expect(wrapper.findAll('.library-card')[1].classes()).toContain('linked')
      wrapper.unmount()
    })

    it('does not submit associations if the application is left while waiting for configuration saves', async () => {
      let finish
      const current = context({ normal_library_id: 202 })
      current.waitForSave.mockImplementation(() => new Promise(resolve => { finish = resolve }))
      const wrapper = shallowMount(Knowledge, options(current))
      await flushPromises()
      await changeTab(wrapper, 'other')
      await wrapper.find('.relation-button').trigger('click')
      wrapper.unmount()
      current.robotInfo.id = 8
      finish()
      await flushPromises()
      expect(api.relate).not.toHaveBeenCalled()
      expect(current.refreshRobot).not.toHaveBeenCalled()
    })

    it('opens card details and the more menu using library identifiers', async () => {
      const opened = vi.spyOn(window, 'open').mockImplementation(() => null)
      const wrapper = shallowMount(Knowledge, options(context({ normal_library_id: 202 })))
      await flushPromises()
      await changeTab(wrapper, 'other')
      await wrapper.find('.library-info').trigger('click')
      expect(opened).toHaveBeenLastCalledWith(JSON.stringify({ path: '/library/details', query: { id: '303' } }), '_blank', 'noopener,noreferrer')
      wrapper.findAllComponents({ name: 'ADropdown' })[2].vm.$emit('openChange', true)
      await flushPromises()
      expect(wrapper.findAll('.library-card')[2].classes()).toContain('actions-open')
      expect(wrapper.findAll('.library-more')[2].attributes('aria-expanded')).toBe('true')
      wrapper.findAllComponents({ name: 'AMenu' })[2].vm.$emit('click', { key: 'view' })
      await flushPromises()
      expect(opened).toHaveBeenLastCalledWith(JSON.stringify({ path: '/public-library/config', query: { library_id: 505 } }), '_blank', 'noopener,noreferrer')
      expect(wrapper.find('.actions-open').exists()).toBe(false)
      wrapper.unmount()
      opened.mockRestore()
    })

    it('filters other libraries by type and association without creating or changing libraries', async () => {
      api.listLibraries.mockResolvedValueOnce({ data: [{ id: 101, type: 2 }, { id: 202, type: 0 }, { id: 303, type: '0', library_name: '普通资料' }, { id: 404, type: '2', library_name: '问答资料' }, { id: 505, type: 1, library_name: '外部资料' }] })
      const wrapper = shallowMount(Knowledge, options(context({ normal_library_id: 202, library_ids: '101,202,404' })))
      await flushPromises()
      await changeTab(wrapper, 'other')
      expect(wrapper.findAll('.library-card')).toHaveLength(3)
      wrapper.findAllComponents({ name: 'ADropdown' })[0].vm.$emit('openChange', true)
      await flushPromises()
      expect(wrapper.find('.actions-open').exists()).toBe(true)
      for (const [index, name] of [[1, '普通资料'], [2, '问答资料'], [3, '外部资料'], [4, '问答资料']]) {
        await wrapper.findAll('.library-filters button')[index].trigger('click')
        expect(wrapper.findAll('.library-name').map(item => item.text())).toEqual([name])
        expect(wrapper.find('.actions-open').exists()).toBe(false)
      }
      expect(wrapper.find('.other-libraries-total').exists()).toBe(true)
      expect(api.listLibraries).toHaveBeenCalledTimes(1)
      expect(api.ensure).not.toHaveBeenCalled()
      expect(api.relate).not.toHaveBeenCalled()
      wrapper.unmount()
    })

    it('uses link state icons and real values for all three metric tooltips', async () => {
      const wrapper = shallowMount(Knowledge, options(context({ normal_library_id: 202, library_ids: '101,202,303' })))
      await flushPromises()
      await changeTab(wrapper, 'other')
      const cards = wrapper.findAll('.library-card')
      expect(cards[0].find('.relation-button').attributes('aria-pressed')).toBe('true')
      expect(cards[0].find('.relation-button img').attributes('src')).toContain('/other-knowledge/unlink.svg')
      expect(cards[1].find('.relation-button').attributes('aria-pressed')).toBe('false')
      expect(cards[1].find('.relation-button img').attributes('src')).toContain('/other-knowledge/link.svg')
      expect(cards[0].find('.library-more img').attributes('src')).toContain('/normal-knowledge/more.svg')
      const tooltipTitles = card => card.findAllComponents({ name: 'ATooltip' }).filter(tooltip => tooltip.props('overlayClassName') === 'other-library-tooltip').map(tooltip => tooltip.vm.$slots.title().filter(node => typeof node.children === 'string').map(node => node.children).join('').trim())
      expect(tooltipTitles(cards[0])).toEqual(['3 个文档', '总大小2 KB', '已关联 2 个'])
      expect(tooltipTitles(cards[1])).toEqual(['0 个文档', '总大小0 B', '已关联 0 个'])
      wrapper.unmount()
    })

    it('renders the real tooltip and its design arrow when a statistic is hovered', async () => {
      const config = options(context({ normal_library_id: 202 }))
      config.attachTo = document.body
      config.global.components = { ATooltip: Tooltip }
      config.global.stubs['a-tooltip'] = false
      config.global.stubs['a-avatar'] = true
      config.global.stubs['a-empty'] = true
      const wrapper = mount(Knowledge, config)
      try {
        await flushPromises()
        await changeTab(wrapper, 'other')
        expect(wrapper.find('.library-info').attributes('title')).toBeUndefined()
        await wrapper.find('.library-metrics [aria-label="3 个文档"]').trigger('mouseenter')
        await vi.waitFor(() => {
          const popup = document.querySelector('.other-library-tooltip .ant-tooltip-inner')
          expect(popup?.textContent).toBe('3 个文档')
          expect(popup?.querySelector('.library-tooltip-beak')?.getAttribute('src')).toContain('/other-knowledge/tooltip-beak.svg')
        })
        expect(document.querySelector('.other-library-tooltip .ant-tooltip-arrow')).toBeNull()
      } finally { wrapper.unmount() }
    })

    it('shows an empty state when only the two built-in libraries are accessible', async () => {
      api.listLibraries.mockResolvedValueOnce({ data: [{ id: 101, type: 2 }, { id: '202', type: 0 }] })
      const wrapper = shallowMount(Knowledge, options(context({ normal_library_id: 202 })))
      await flushPromises()
      await changeTab(wrapper, 'other')
      expect(wrapper.find('.library-card').exists()).toBe(false)
      expect(wrapper.find('a-empty-stub').attributes('description')).toBe('empty_libraries')
      expect(api.ensure).not.toHaveBeenCalled()
      wrapper.unmount()
    })
  })

  it('opens imported workflows and their logs using workflow identifiers', async () => {
    const opened = vi.spyOn(window, 'open').mockImplementation(() => null)
    const wrapper = shallowMount(Functions, options(context()))
    await flushPromises()
    expect(api.listRobots).toHaveBeenCalledWith({ application_type: 1 })
    expect(wrapper.findComponent(ImportCsl).props('applicationType')).toBe(1)
    wrapper.findComponent(ImportCsl).vm.$emit('ok', { id: 901, robot_key: 'imported-flow' })
    await flushPromises()
    expect(opened).toHaveBeenCalledWith(JSON.stringify({ path: '/robot/config/workflow', query: { id: 901, robot_key: 'imported-flow', source: 'customer-service' } }), '_blank', 'noopener,noreferrer')
    wrapper.findAllComponents({ name: 'AMenu' }).at(-1).vm.$emit('click', { key: 'logs' })
    await flushPromises()
    expect(wrapper.findComponent(InvokeLogs).props()).toMatchObject({ robotId: 900, robotKey: 'workflow-nine' })
    wrapper.unmount()
    opened.mockRestore()
  })

  describe('function center', () => {
    it('shows all workflows and rate limiting, opens cards and new workflows in new tabs', async () => {
      api.listRobots.mockResolvedValueOnce({ data: [
        { id: 900, robot_key: 'flow-nine', robot_name: 'Flow nine', application_type: '1' },
        { id: 901, robot_key: 'flow-ten', robot_name: 'Flow ten', application_type: 1 },
        { id: 7, robot_key: 'customer-seven', robot_name: 'Customer', application_type: 0 }
      ] })
      const opened = vi.spyOn(window, 'open').mockImplementation(() => null)
      const wrapper = shallowMount(Functions, options(context()))
      await flushPromises()
      expect(wrapper.findAll('.workflow-card')).toHaveLength(2)
      expect(wrapper.find('.rate-card').exists()).toBe(true)
      expect(wrapper.text()).not.toContain('handoff')
      await wrapper.find('.workflow-card').trigger('click')
      expect(opened).toHaveBeenLastCalledWith(JSON.stringify({ path: '/robot/config/workflow', query: { id: 900, robot_key: 'flow-nine', source: 'customer-service' } }), '_blank', 'noopener,noreferrer')
      const creator = wrapper.findComponent(AddRobotAlert)
      expect(creator.props('createRedirect')).toBe('none')
      await wrapper.findAllComponents({ name: 'AMenuItem' }).find(item => item.text() === 'create_workflow').trigger('click')
      expect(api.openCreator).toHaveBeenCalledWith(1)
      creator.vm.$emit('ok', { id: 902, robot_key: 'created-flow' })
      await flushPromises()
      expect(opened).toHaveBeenLastCalledWith(JSON.stringify({ path: '/robot/config/workflow', query: { id: 902, robot_key: 'created-flow', source: 'customer-service', show_tips: 1 } }), '_blank', 'noopener,noreferrer')
      expect(api.listRobots).toHaveBeenCalledTimes(2)
      wrapper.unmount()
      opened.mockRestore()
    })

    it('persists workflow associations through the save queue, preserving other associations', async () => {
      const current = context({ work_flow_ids: '901' })
      current.refreshRobot.mockImplementation(async () => { current.robotInfo.work_flow_ids = '901,902' })
      let finish
      api.workflowRelation.mockImplementationOnce(() => new Promise(resolve => { finish = resolve }))
      current.saveAction.mockImplementation(async task => {
        await task()
        current.robotInfo.work_flow_ids = api.workflowRelation.mock.calls.at(-1)[0].work_flow_ids
        return true
      })
      const opened = vi.spyOn(window, 'open').mockImplementation(() => null)
      const wrapper = shallowMount(Functions, options(current))
      await flushPromises()
      const workflowSwitch = wrapper.find('.workflow-card').findComponent({ name: 'ASwitch' })
      await workflowSwitch.trigger('click')
      await flushPromises()
      expect(api.workflowRelation).toHaveBeenCalledWith({ id: 7, work_flow_ids: '901,902,900' })
      expect(current.saveAction).toHaveBeenCalledTimes(1)
      expect(workflowSwitch.props('disabled')).toBe(true)
      expect(opened).not.toHaveBeenCalled()
      finish({ res: 0 })
      await flushPromises()
      expect(workflowSwitch.props('checked')).toBe(true)
      expect(wrapper.find('.workflow-card').classes()).toContain('active')
      current.refreshRobot.mockResolvedValue()
      api.workflowRelation.mockResolvedValueOnce({ res: 0 })
      await workflowSwitch.trigger('click')
      await flushPromises()
      expect(api.workflowRelation).toHaveBeenLastCalledWith({ id: 7, work_flow_ids: '901,902' })
      expect(workflowSwitch.props('checked')).toBe(false)
      wrapper.unmount()
      opened.mockRestore()
    })

    it('retains the saved association after failure and ignores a pending toggle after leaving the application', async () => {
      const current = context({ work_flow_ids: '900' })
      current.saveAction.mockResolvedValueOnce(false)
      const wrapper = shallowMount(Functions, options(current))
      await flushPromises()
      const workflowSwitch = wrapper.find('.workflow-card').findComponent({ name: 'ASwitch' })
      await workflowSwitch.trigger('click')
      await flushPromises()
      expect(workflowSwitch.props('checked')).toBe(true)
      let finish
      current.waitForSave.mockImplementationOnce(() => new Promise(resolve => { finish = resolve }))
      await workflowSwitch.trigger('click')
      wrapper.unmount()
      current.robotInfo.id = 8
      finish()
      await flushPromises()
      expect(current.saveAction).toHaveBeenCalledTimes(1)
    })

    it('uses the rate form validation for the card switch and opens its settings when saving fails', async () => {
      api.setRate.mockResolvedValueOnce(true).mockResolvedValueOnce(false)
      const opened = vi.spyOn(window, 'open').mockImplementation(() => null)
      const wrapper = shallowMount(Functions, { ...options(context()), attachTo: document.body })
      await flushPromises()
      const rate = wrapper.findComponent(RateLimit)
      expect(rate.props('robotKey')).toBe('customer-seven')
      rate.vm.$emit('status-change', { ready: true, loading: false, enabled: true, error: '' })
      await flushPromises()
      const rateSwitch = wrapper.find('.rate-card').findComponent({ name: 'ASwitch' })
      await rateSwitch.trigger('click')
      await flushPromises()
      expect(api.setRate).toHaveBeenCalledWith(false)
      expect(wrapper.find('.workflow-toolbar').exists()).toBe(true)
      await rateSwitch.trigger('click')
      await flushPromises()
      expect(wrapper.find('.workflow-toolbar').exists()).toBe(false)
      expect(wrapper.find('.rate-pane').isVisible()).toBe(true)
      expect(opened).not.toHaveBeenCalled()
      await wrapper.find('.back-button').trigger('click')
      expect(wrapper.find('.workflow-toolbar').exists()).toBe(true)
      expect(wrapper.find('.rate-pane').isVisible()).toBe(false)
      wrapper.unmount()
      opened.mockRestore()
    })

    it('edits, copies, deletes and exports the selected workflow without using the current customer ID', async () => {
      const wrapper = shallowMount(Functions, options(context()))
      await flushPromises()
      const workflow = { id: 900, robot_key: 'workflow-nine', robot_name: 'Flow', application_type: 1 }
      const menu = () => wrapper.findAllComponents({ name: 'AMenu' }).at(-1)
      menu().vm.$emit('click', { key: 'edit' })
      await flushPromises()
      expect(api.robotInfo).toHaveBeenCalledWith({ id: 900 })
      expect(api.openCreator).toHaveBeenCalledWith(1, true, { ...workflow, en_name: 'workflow_nine' })
      menu().vm.$emit('click', { key: 'export' })
      expect(api.openExport).toHaveBeenCalledWith(workflow)
      wrapper.findComponent(ExportModal).vm.$emit('ok', workflow, { export_library: 1 })
      expect(api.exportUrl).toHaveBeenCalledWith(expect.stringContaining('id=900'))
      menu().vm.$emit('click', { key: 'copy' })
      expect(api.copy).not.toHaveBeenCalled()
      await Modal.confirm.mock.calls.at(-1)[0].onOk()
      expect(api.copy).toHaveBeenCalledWith({ id: 900 })
      menu().vm.$emit('click', { key: 'delete' })
      expect(api.remove).not.toHaveBeenCalled()
      await Modal.confirm.mock.calls.at(-1)[0].onOk()
      expect(api.remove).toHaveBeenCalledWith({ id: 900 })
      wrapper.unmount()
    })
  })

  it('keeps the current robot and date/channel conditions when opening unknown details', async () => {
    const wrapper = shallowMount(Analytics, options(context()))
    const filters = { start_date: '2026-10-01', end_date: '2026-10-20', app_type: 'web', app_id: -1 }
    wrapper.findComponent(HitStatistics).vm.$emit('unknown', filters)
    await flushPromises()
    expect(wrapper.findComponent(UnknownStats).props()).toMatchObject({ robotId: 7, filters })
    await changeTab(wrapper, 'data')
    expect(wrapper.findComponent(HitStatistics).props('filters')).toEqual(filters)
    wrapper.unmount()
  })
})
