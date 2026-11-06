export function resolveSdkAvatarUrl(icon: string, origin: string, fallback: string): string {
  try {
    return new URL(icon || fallback, origin).href
  } catch {
    return fallback
  }
}
