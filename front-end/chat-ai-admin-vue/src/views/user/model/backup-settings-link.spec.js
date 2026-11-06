import { afterEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, shallowMount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { createPinia } from 'pinia'
import ModelPage from './index.vue'
import { getModelConfigList } from '@/api/model/index'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/api/model/index', () => {
  const list = () => Promise.resolve({ data: [] })
  return {
    getModelConfigList: vi.fn(list), delModelConfig: list, getSelfModelBuylog: list,
    getSelfModelConfigs: list, showModelConfigList: list, delUseModelConfig: list,
    setModelConfigWeight: list
  }
})
vi.mock('./components/page-tabs.vue', () => ({ default: { name: 'PageTabs', props: ['value'], template: '<div />' } }))
vi.mock('./components/model-exception.vue', () => ({ default: { name: 'ModelException', template: '<div>backup-settings</div>' } }))
vi.mock('./components/set-model-alert.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/see-model-alert.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/buy-record.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/kefu-modal.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/has-modal-list.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/add-model-new.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/add-confing-modal.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./components/model-list.vue', () => ({ default: { template: '<div />' } }))
vi.mock('./model-config', () => ({ getModelTableConfig: () => [] }))
vi.mock('@/utils/index', () => ({ formatPriceWithCommas: value => value }))
vi.mock('@/stores/modules/user', () => ({ useUserStore: () => ({ getToken: '' }) }))

const wrappers = []
afterEach(() => wrappers.splice(0).forEach(wrapper => wrapper.unmount()))
async function render(path) {
  const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/user/model', component: ModelPage }] })
  await router.push(path)
  const wrapper = shallowMount(ModelPage, { global: {
    plugins: [createPinia(), router],
    stubs: { ASelect: { inheritAttrs: false, props: ['children'], template: '<div />' } }
  } })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}

describe('model settings deep link', () => {
  it.each([
    ['no model configured', []],
    ['an existing primary model', [{ model_define: 'openai', config_info: { id: 7 } }]]
  ])('opens the backup settings tab with %s', async (_name, models) => {
    getModelConfigList.mockResolvedValueOnce({ data: models })
    const wrapper = await render('/user/model?activeTab=3')
    expect(wrapper.findComponent({ name: 'PageTabs' }).props('value')).toBe(3)
    expect(wrapper.findComponent({ name: 'ModelException' }).exists()).toBe(true)
  })
})
