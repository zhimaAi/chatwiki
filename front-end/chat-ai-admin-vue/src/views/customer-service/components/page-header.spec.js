import { describe, expect, it } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, ref } from 'vue'
import CustomerPageHeader from './page-header.vue'

const tabs = [{ key: 'qa', label: '问答知识库' }, { key: 'normal', label: '普通知识库' }, { key: 'other', label: '其他知识库' }]
function render(options = {}, mountOptions = {}) {
  return mount(CustomerPageHeader, { props: { title: '知识库', tabs, activeKey: 'qa', ...options }, ...mountOptions })
}

describe('customer page header', () => {
  it('keeps title-only headers, save status, back navigation and extra actions in the same component', async () => {
    const wrapper = mount(CustomerPageHeader, {
      props: { title: '基础设置', back: true, backLabel: '返回功能中心' },
      slots: { default: '<span>保存中</span>', extra: '<button>操作</button>' }
    })
    expect(wrapper.find('h1').text()).toBe('基础设置')
    expect(wrapper.find('[role="tablist"]').exists()).toBe(false)
    expect(wrapper.find('.header-status').text()).toBe('保存中')
    expect(wrapper.find('.header-extra').text()).toBe('操作')
    await wrapper.find('.back-button').trigger('click')
    expect(wrapper.emitted('back')).toHaveLength(1)
    wrapper.unmount()
  })

  it('updates the active key before notifying the page and does not reload an already active tab', async () => {
    const host = mount(defineComponent({
      components: { CustomerPageHeader },
      setup() { return { active: ref('qa'), observed: ref(''), tabs } },
      template: '<CustomerPageHeader title="知识库" v-model:active-key="active" :tabs="tabs" @change="observed = active" /><output>{{ observed }}</output>'
    }))
    const header = host.findComponent(CustomerPageHeader)
    await header.findAll('[role="tab"]')[0].trigger('click')
    expect(header.emitted('change')).toBeUndefined()
    await header.findAll('[role="tab"]')[1].trigger('click')
    expect(host.find('output').text()).toBe('normal')
    expect(header.findAll('[role="tab"]')[1].attributes('aria-selected')).toBe('true')
    expect(header.emitted('change')).toEqual([['normal']])
    host.unmount()
  })

  it('supports keyboard navigation, wraps at the ends and skips disabled tabs', async () => {
    const wrapper = render({ tabs: [tabs[0], { ...tabs[1], disabled: true }, tabs[2]] }, { attachTo: document.body })
    await wrapper.findAll('[role="tab"]')[0].trigger('keydown', { key: 'ArrowRight' })
    expect(wrapper.emitted('update:activeKey').at(-1)).toEqual(['other'])
    expect(document.activeElement).toBe(wrapper.findAll('[role="tab"]')[2].element)
    await wrapper.setProps({ activeKey: 'other' })
    await wrapper.findAll('[role="tab"]')[2].trigger('keydown', { key: 'ArrowRight' })
    expect(wrapper.emitted('update:activeKey').at(-1)).toEqual(['qa'])
    await wrapper.setProps({ activeKey: 'qa' })
    await wrapper.findAll('[role="tab"]')[0].trigger('keydown', { key: 'ArrowLeft' })
    expect(wrapper.emitted('update:activeKey').at(-1)).toEqual(['other'])
    await wrapper.findAll('[role="tab"]')[2].trigger('keydown', { key: 'Home' })
    await flushPromises()
    expect(document.activeElement).toBe(wrapper.findAll('[role="tab"]')[0].element)
    expect(wrapper.findAll('[role="tab"]')[1].attributes('disabled')).toBeDefined()
    await wrapper.findAll('[role="tab"]')[0].trigger('keydown', { key: 'End' })
    expect(wrapper.emitted('update:activeKey').at(-1)).toEqual(['other'])
    wrapper.unmount()
  })

  it('keeps numeric keys and translated labels reactive', async () => {
    const wrapper = render({ tabs: [{ key: 0, label: '统计' }, { key: 1, label: '反馈' }], activeKey: 0 })
    expect(wrapper.findAll('[role="tab"]')[0].attributes('tabindex')).toBe('0')
    await wrapper.findAll('[role="tab"]')[1].trigger('click')
    expect(wrapper.emitted('change')).toEqual([[1]])
    await wrapper.setProps({ title: 'Analytics', tabs: [{ key: 0, label: 'Statistics' }, { key: 1, label: 'Feedback' }] })
    expect(wrapper.find('h1').text()).toBe('Analytics')
    expect(wrapper.find('[role="tab"]').text()).toBe('Statistics')
    wrapper.unmount()
  })
})
