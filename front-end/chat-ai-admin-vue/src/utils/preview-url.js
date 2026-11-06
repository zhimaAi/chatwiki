export const toPageProtocol = (url, protocol = window.location.protocol) => {
  if (!url || typeof url !== 'string') return ''
  if (protocol !== 'https:') return url
  return url.replace(/^http:\/\//i, 'https://')
}
