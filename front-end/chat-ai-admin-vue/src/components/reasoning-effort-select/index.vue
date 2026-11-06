<template>
  <div class="reasoning-effort-select" @keydown.capture="handleKeydown">
    <a-auto-complete
      :value="inputValue"
      :open="dropdownOpen"
      :options="options"
      :dropdown-match-select-width="120"
      :default-active-first-option="false"
      :status="invalid ? 'error' : undefined"
      :placeholder="t('ph_select_or_input')"
      :allow-clear="false"
      @search="handleSearch"
      @change="draftValue = $event"
      @select="confirmValue"
      @dropdownVisibleChange="handleDropdownVisibleChange"
      @focus="handleFocus"
      @blur="handleBlur"
    />
    <div v-if="invalid" class="error-text">{{ t('error_max_length') }}</div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { REASONING_EFFORTS, isValidReasoningEffort } from '@/utils/reasoning-effort'

const props = defineProps({ value: { type: String, default: 'none' } })
const emit = defineEmits(['update:value', 'change'])
const { t } = useI18n('components.reasoning-effort-select.index')
const draftValue = ref('')
const searchValue = ref('')
const focused = ref(false)
const dropdownOpen = ref(false)
const keyboardNavigated = ref(false)
const invalid = ref(false)
const displayValue = computed(() => props.value)
const inputValue = computed(() => focused.value ? draftValue.value : displayValue.value)
const options = computed(() => {
  const values = REASONING_EFFORTS.includes(props.value)
    ? REASONING_EFFORTS
    : [...REASONING_EFFORTS, props.value]
  const query = searchValue.value.toLowerCase()
  return values.filter(value => value.toLowerCase().includes(query)).map(value => ({ value, label: value }))
})

const resetDraft = () => {
  draftValue.value = displayValue.value
  searchValue.value = ''
  invalid.value = false
  handleDropdownVisibleChange(false)
}

watch(displayValue, resetDraft)

const handleFocus = () => {
  focused.value = true
  resetDraft()
  handleDropdownVisibleChange(true)
}

const handleBlur = () => {
  focused.value = false
  resetDraft()
}

const handleSearch = (value) => {
  searchValue.value = value
  invalid.value = !isValidReasoningEffort(value)
  keyboardNavigated.value = false
}

const handleDropdownVisibleChange = (open) => {
  dropdownOpen.value = open
  if (!open) keyboardNavigated.value = false
}

const handleKeydown = (event) => {
  if (event.isComposing) return
  if (event.key === 'Escape') {
    if (!dropdownOpen.value && draftValue.value === displayValue.value) return
    event.preventDefault()
    event.stopPropagation()
    resetDraft()
    return
  }
  if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
    keyboardNavigated.value = true
    dropdownOpen.value = true
    return
  }
  if (event.key !== 'Enter' || (keyboardNavigated.value && dropdownOpen.value && options.value.length)) return

  event.preventDefault()
  event.stopPropagation()
  confirmValue(searchValue.value.trim())
}

const confirmValue = (value) => {
  if (!value) {
    resetDraft()
    return
  }
  if (!isValidReasoningEffort(value)) {
    invalid.value = true
    return
  }
  resetDraft()
  draftValue.value = value
  if (value === props.value) return
  emit('update:value', value)
  emit('change', value)
}
</script>

<style scoped>
.reasoning-effort-select { min-width: 0; width: 120px; max-width: 100%; }
.reasoning-effort-select :deep(.ant-select) { width: 100%; }
.reasoning-effort-select :deep(.ant-select-selection-search-input) { text-overflow: ellipsis; }
.error-text { color: #ff4d4f; font-size: 12px; }
</style>
