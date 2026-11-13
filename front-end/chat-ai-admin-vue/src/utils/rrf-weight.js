export function normalizeRrfWeight(weight) {
  const vector = Number(weight?.vector ?? 0)
  const search = Number(weight?.search ?? 0)
  const sum = vector + search
  const percent = Number.isFinite(sum) && sum > 0 ? Math.round(vector * 100 / sum) : 70
  return { vector: percent, search: 100 - percent }
}
