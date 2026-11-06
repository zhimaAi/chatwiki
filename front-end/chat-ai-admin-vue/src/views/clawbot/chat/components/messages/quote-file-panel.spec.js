import { describe, expect, it, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
vi.mock('@/hooks/web/useI18n', () => ({ useI18n: () => ({ t: key => key }) }))
vi.mock('vue-router', () => ({ useRouter: () => ({ resolve: ({ path, query }) => ({ href: `#${path}?id=${query.id}` }) }) }))
import QuoteFilePanel from './quote-file-panel.vue'

describe('knowledge source links', () => {
  it('opens documents by file ID and featured knowledge by library ID', () => {
    const wrapper = shallowMount(QuoteFilePanel, {
      props: { item: { show_quote_file: true, quote_file: [
        { id: 11, file_name: 'qa.csv', library_id: 501 },
        { id: 12, file_name: 'source.pdf', library_id: 502 },
        { id: 13, library_id: 503, library_name: 'Featured' }
      ] } }
    })
    const links = wrapper.findAll('a')
    expect(links.map(link => link.attributes('href'))).toEqual(['#/library/preview?id=11', '#/library/preview?id=12', '#/library/details/categary-manage?id=503'])
    expect(links.every(link => link.attributes('target') === '_blank' && link.attributes('rel').includes('noopener'))).toBe(true)
    wrapper.unmount()
  })
})
