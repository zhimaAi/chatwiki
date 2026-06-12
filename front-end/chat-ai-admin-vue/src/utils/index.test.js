import { describe, expect, it } from 'vitest'
import { objToFormData } from './index'

describe('objToFormData', () => {
  it('passes FormData through unchanged', () => {
    const formData = new FormData()
    formData.append('name', 'test')

    expect(objToFormData(formData)).toBe(formData)
  })

  it('expands bracket-array fields into repeated form entries', () => {
    const formData = objToFormData({
      id: 1,
      'keywords[]': ['foo', 'bar'],
      description: 'desc'
    })

    expect(formData.get('id')).toBe('1')
    expect(formData.get('description')).toBe('desc')
    expect(formData.getAll('keywords[]')).toEqual(['foo', 'bar'])
  })
})
