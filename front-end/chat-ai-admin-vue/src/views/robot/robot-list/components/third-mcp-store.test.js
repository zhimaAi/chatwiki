import { describe, it, expect, vi, afterEach } from 'vitest'
import { shallowMount, flushPromises } from '@vue/test-utils'
import ThirdMcpStore from './third-mcp-store.vue'
import { saveTMcpProvider, authTMcpProvider } from '@/api/robot/thirdMcp.js'

vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('@/api/robot/thirdMcp.js', () => ({
  saveTMcpProvider: vi.fn(() => Promise.resolve({ data: { provider_id: 42 } })),
  authTMcpProvider: vi.fn(() => Promise.resolve({ data: {} }))
}))
vi.mock('@/utils/http/axios/config.js', () => ({ setShowReqError: vi.fn() }))
vi.mock('ant-design-vue', () => ({ message: {
  loading: vi.fn(() => vi.fn()), success: vi.fn(), warning: vi.fn(), error: vi.fn()
} }))

function mountForm() {
  return shallowMount(ThirdMcpStore, { global: {
    renderStubDefaultSlot: true,
    stubs: {
      'a-input': { name: 'AInput', props: ['value'], template: '<input :value="value" />' },
      'a-input-number': { props: ['value'], template: '<input :value="value" />' },
      'a-modal': { template: '<section><slot/><button data-test="save" @click="$emit(\'ok\')">Save</button></section>' },
      'a-button': { template: '<button><slot/></button>' }
    }
  } })
}
function presetButton(wrapper) {
  return wrapper.findAll('button').find(b => b.text() === 'btn_parallel_search')
}

afterEach(() => { vi.useRealTimers(); vi.clearAllMocks() })

describe('Parallel Search MCP preset', () => {
  it('saves an anonymous preset through the existing form and authorization flow', async () => {
    vi.useFakeTimers()
    const wrapper = mountForm()
    wrapper.vm.show()
    await flushPromises()
    await presetButton(wrapper).trigger('click')
    await wrapper.get('[data-test="save"]').trigger('click')
    await flushPromises()
    expect(saveTMcpProvider).toHaveBeenCalledOnce()
    const saved = saveTMcpProvider.mock.calls[0][0]
    expect(saved).toMatchObject({ name: 'Parallel Search', url: 'https://search.parallel.ai/mcp', request_timeout: '30' })
    expect(JSON.parse(saved.headers)).toEqual({ 'User-Agent': 'ChatWiki/0.1.0' })
    expect(saved.description).toBeTruthy()
    await vi.advanceTimersByTimeAsync(1800)
    await flushPromises()
    expect(authTMcpProvider).toHaveBeenCalledWith({ provider_id: 42 })
    wrapper.unmount()
  })

  it('keeps the manual form empty until explicitly selected and resets after reopening', async () => {
    const wrapper = mountForm()
    wrapper.vm.show()
    await flushPromises()
    expect(wrapper.findAll('input').map(i => i.element.value)).toEqual(['', '', '30'])
    await presetButton(wrapper).trigger('click')
    expect(wrapper.findAll('input').map(i => i.element.value)).toEqual([
      'Parallel Search', 'https://search.parallel.ai/mcp', '30', 'User-Agent', 'ChatWiki/0.1.0'
    ])
    wrapper.vm.show()
    await flushPromises()
    expect(wrapper.findAll('input').map(i => i.element.value)).toEqual(['', '', '30'])
    wrapper.unmount()
  })

  it('does not offer or change presets when editing a saved provider', async () => {
    const wrapper = mountForm()
    const info = { id: 7, name: 'Existing', description: 'Saved server', url: 'https://example.org/mcp', request_timeout: '12', headers: '{"Authorization":"Bearer fixture"}' }
    wrapper.vm.show(info)
    await flushPromises()
    expect(presetButton(wrapper)).toBeUndefined()
    await wrapper.get('[data-test="save"]').trigger('click')
    await flushPromises()
    expect(saveTMcpProvider.mock.calls[0][0]).toMatchObject({ name: info.name, url: info.url, request_timeout: info.request_timeout, provider_id: 7, headers: info.headers })
    expect(authTMcpProvider).not.toHaveBeenCalled()
    wrapper.unmount()
  })
})
