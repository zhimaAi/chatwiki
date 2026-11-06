export const SESSION_EXPIRE_MIN_HOURS = 0.1
export const SESSION_EXPIRE_MAX_HOURS = 48
export const DEFAULT_SESSION_EXPIRE_HOURS = 24

export const isValidSessionExpireHours = (hours) => {
  const value = Number(hours)
  return Number.isFinite(value) && value >= SESSION_EXPIRE_MIN_HOURS && value <= SESSION_EXPIRE_MAX_HOURS
}

export const sessionExpireHoursFromSeconds = (seconds) => {
  const hours = Number(seconds) / 3600
  if (!isValidSessionExpireHours(hours)) {
    return DEFAULT_SESSION_EXPIRE_HOURS
  }
  return Math.round(hours * 10) / 10
}
