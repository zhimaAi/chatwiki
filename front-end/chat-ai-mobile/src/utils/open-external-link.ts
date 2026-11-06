export function openExternalLink(url: string, config: { open_type?: number | string; window_width?: number | string; window_height?: number | string }) {
  const mobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
  if (mobile || Number(config.open_type) === 1) return window.open(url, '_blank')
  const width = Number(config.window_width) || 1200
  const height = Number(config.window_height) || 650
  const left = (window.screen.width - width) / 2
  const top = (window.screen.height - height) / 2
  return window.open(url, '_blank', `width=${width},height=${height},left=${left},top=${top},resizable=yes`)
}
