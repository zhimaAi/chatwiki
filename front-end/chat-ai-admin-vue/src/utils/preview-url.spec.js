import { describe, expect, it } from 'vitest'
import { toPageProtocol } from './preview-url'

describe('toPageProtocol', () => {
  it('upgrades http previews when the admin is https', () => {
    expect(toPageProtocol('http://h5.example.com/#/chat?robot_key=a', 'https:'))
      .toBe('https://h5.example.com/#/chat?robot_key=a')
  })

  it('keeps http previews on an http admin', () => {
    expect(toPageProtocol('http://h5.example.com/#/chat', 'http:'))
      .toBe('http://h5.example.com/#/chat')
  })

  it('leaves relative, protocol-relative and empty values alone', () => {
    expect(toPageProtocol('//h5.example.com/#/chat', 'https:'))
      .toBe('//h5.example.com/#/chat')
    expect(toPageProtocol('/preview/h5', 'https:')).toBe('/preview/h5')
    expect(toPageProtocol('', 'https:')).toBe('')
  })
})
