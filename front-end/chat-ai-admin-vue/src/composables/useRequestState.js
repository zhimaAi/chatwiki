import { onBeforeUnmount, ref } from 'vue'

export function useRequestState() {
  const loading = ref(false)
  const error = ref('')
  let sequence = 0
  let alive = true
  async function run(request, apply) {
    const version = ++sequence
    loading.value = true
    error.value = ''
    try {
      const result = await request()
      if (!alive || version !== sequence) return null
      apply?.(result)
      return result
    } catch (err) {
      if (alive && version === sequence) error.value = err?.message || err?.msg || String(err)
      return null
    } finally {
      if (alive && version === sequence) loading.value = false
    }
  }
  onBeforeUnmount(() => { alive = false; sequence++ })
  return { loading, error, run }
}
