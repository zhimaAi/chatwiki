import { describe, expect, it } from 'vitest'
import {
  DEFAULT_SESSION_EXPIRE_HOURS,
  isValidSessionExpireHours,
  sessionExpireHoursFromSeconds
} from './session-expire'

describe('session expire hours', () => {
  it('accepts the configured range and rejects values outside it', () => {
    expect(isValidSessionExpireHours(0.1)).toBe(true)
    expect(isValidSessionExpireHours(24)).toBe(true)
    expect(isValidSessionExpireHours(48)).toBe(true)
    expect(isValidSessionExpireHours(0.09)).toBe(false)
    expect(isValidSessionExpireHours(48.1)).toBe(false)
  })

  it('converts saved seconds and falls back to 24 hours', () => {
    expect(sessionExpireHoursFromSeconds(86400)).toBe(24)
    expect(sessionExpireHoursFromSeconds(7200)).toBe(2)
    expect(sessionExpireHoursFromSeconds(360)).toBe(0.1)
    expect(sessionExpireHoursFromSeconds(0)).toBe(DEFAULT_SESSION_EXPIRE_HOURS)
  })
})
