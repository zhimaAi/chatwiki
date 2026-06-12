<template>
  <div ref="rootRef" class="chat-model-select">
    <button class="chat-model-trigger" type="button" :disabled="disabled" @click="toggleOpen">
      <span>{{ triggerLabel }}</span>
      <DownOutlined class="trigger-arrow" :class="{ open }" />
    </button>

    <div v-if="open" class="chat-model-panel">
      <div class="panel-search">
        <SearchOutlined class="panel-search-icon" />
        <input
          v-model.trim="keyword"
          class="panel-search-input"
          type="text"
          :placeholder="t('ph_search_model')"
        />
      </div>

      <div class="panel-content">
        <div class="provider-list-wrapper">
          <div class="provider-list">
            <button
              v-for="provider in filteredProviders"
              :key="provider.id"
              type="button"
              class="provider-item"
              :class="{ active: provider.id === selectedProviderId }"
              @click="selectedProviderId = provider.id"
            >
              <div class="provider-info">
                <img v-if="provider.model_icon_url" class="provider-icon" :src="provider.model_icon_url" alt="" />
                <span v-else class="provider-icon fallback">{{ getProviderFallback(provider) }}</span>
                <span class="provider-name">{{ provider.providerName }}</span>
              </div>
              <span class="provider-count">{{ provider.children.length }}</span>
            </button>
          </div>
        </div>

        <div class="model-section">
          <div class="model-list">
            <button
              v-for="model in visibleModels"
              :key="model.key"
              type="button"
              class="model-item"
              :class="{ active: isSelectedModel(model) }"
              @click="handleSelectModel(model)"
            >
              <span class="model-name">{{ model.label || model.show_model_name || model.name }}</span>
              <CheckOutlined v-if="isSelectedModel(model)" class="model-check" />
            </button>

            <div v-if="!visibleModels.length" class="model-empty">{{ t('msg_no_models') }}</div>
          </div>

          <div class="panel-footer">
            <slot name="add-model" :onClick="handleAddModel">
              <button
                class="add-model-btn"
                type="button"
                @click="handleAddModel"
              >
                <PlusOutlined class="add-model-icon" />
                <span>{{ t('btn_add_model') }}</span>
              </button>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { CheckOutlined, DownOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getModelConfigOption } from '@/api/model/index'
import { getModelOptionsList } from '@/components/model-select/index.js'

const { t } = useI18n('views.clawbot.chat.components.chat-model-select')
const emit = defineEmits(['update:modeId', 'update:modeName', 'update:useConfigId', 'change'])

const props = defineProps({
  modelType: {
    type: String,
    default: 'LLM'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  modeId: {
    type: [String, Number],
    default: ''
  },
  modeName: {
    type: String,
    default: ''
  },
  useConfigId: {
    type: [String, Number],
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  }
})

const router = useRouter()
const rootRef = ref(null)
const open = ref(false)
const keyword = ref('')
const providerList = ref([])
const selectedProviderId = ref('')
const innerModeId = ref(props.modeId ? String(props.modeId) : '')
const innerModeName = ref(props.modeName || '')
const innerUseConfigId = ref(props.useConfigId ? String(props.useConfigId) : '')

watch(
  () => [props.modeId, props.modeName, props.useConfigId],
  ([modeId, modeName, useConfigId]) => {
    innerModeId.value = modeId ? String(modeId) : ''
    innerModeName.value = modeName || ''
    innerUseConfigId.value = useConfigId ? String(useConfigId) : ''
  },
  { immediate: true }
)

watch(
  () => props.disabled,
  (disabled) => {
    if (disabled) {
      open.value = false
    }
  }
)

const normalizeProviders = (list = []) => {
  return list.map((item) => ({
    ...item,
    id: String(item.id),
    providerName: item.config_name || item.model_name || '',
    children: (item.children || []).map((child) => ({
      ...child,
      model_config_id: String(child.model_config_id),
      id: child.id ? String(child.id) : '',
      label: child.label || child.show_model_name || child.name || '',
    })),
  }))
}

const filteredProviders = computed(() => {
  const search = keyword.value.trim().toLowerCase()

  if (!search) {
    return providerList.value.filter((item) => item.children.length > 0)
  }

  return providerList.value
    .map((provider) => {
      const providerMatch = provider.providerName.toLowerCase().includes(search)
      const children = provider.children.filter((child) => {
        return (child.label || '').toLowerCase().includes(search)
      })

      return {
        ...provider,
        children: providerMatch ? provider.children : children,
      }
    })
    .filter((provider) => provider.children.length > 0)
})

const selectedProvider = computed(() => {
  return filteredProviders.value.find((item) => item.id === selectedProviderId.value) || filteredProviders.value[0] || null
})

const visibleModels = computed(() => {
  return selectedProvider.value?.children || []
})

const selectedModel = computed(() => {
  for (const provider of providerList.value) {
    const matched = provider.children.find((child) => {
      return child.model_config_id === innerModeId.value && child.name === innerModeName.value
    })

    if (matched) {
      return matched
    }
  }

  return null
})

const triggerLabel = computed(() => {
  return selectedModel.value?.label || props.placeholder || t('ph_select_model')
})

const syncSelectedProvider = () => {
  if (!filteredProviders.value.length) {
    selectedProviderId.value = ''
    return
  }

  const selectedExists = filteredProviders.value.some((item) => item.id === selectedProviderId.value)
  if (selectedExists) {
    return
  }

  const providerBySelection = filteredProviders.value.find((item) =>
    item.children.some((child) => child.model_config_id === innerModeId.value && child.name === innerModeName.value)
  )

  selectedProviderId.value = providerBySelection?.id || filteredProviders.value[0].id
}

watch(filteredProviders, syncSelectedProvider, { immediate: true })
watch(keyword, syncSelectedProvider)

const getProviderFallback = (provider) => {
  return (provider.providerName || '?').slice(0, 1).toUpperCase()
}

const isSelectedModel = (model) => {
  return model.model_config_id === innerModeId.value && model.name === innerModeName.value
}

const toggleOpen = () => {
  if (props.disabled) {
    return
  }

  open.value = !open.value
}

const handleSelectModel = (model) => {
  if (props.disabled) {
    return
  }

  innerModeId.value = model.model_config_id
  innerModeName.value = model.name
  innerUseConfigId.value = model.id || ''

  emit('update:modeId', model.model_config_id)
  emit('update:modeName', model.name)
  emit('update:useConfigId', model.id || '')
  emit('change', model)

  open.value = false
}

const handleAddModel = () => {
  if (props.disabled) {
    return
  }

  open.value = false
  router.push('/user/model')
}

const getModelList = async () => {
  const res = await getModelConfigOption({
    model_type: props.modelType
  })

  let list = res?.data || []
  list = list.sort((a, b) => {
    if (a.weight > 0 && b.weight > 0) return a.weight - b.weight
    if (a.weight > 0 && b.weight <= 0) return -1
    if (a.weight <= 0 && b.weight > 0) return 1
    return 0
  })

  const { newList } = getModelOptionsList(list)
  providerList.value = normalizeProviders(newList)
  syncSelectedProvider()
}

const handleDocumentClick = (event) => {
  if (!rootRef.value?.contains(event.target)) {
    open.value = false
  }
}

onMounted(() => {
  getModelList()
  document.addEventListener('mousedown', handleDocumentClick)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleDocumentClick)
})
</script>

<style lang="less" scoped>
.chat-model-select {
  position: relative;
}

.chat-model-trigger {
  height: 32px;
  padding: 0 12px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: #262626;
  font-size: 14px;
  line-height: 22px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    background: #f2f4f7;
  }

  &:disabled {
    cursor: not-allowed;
    color: #8c8c8c;
    background: transparent;
  }
}

.trigger-arrow {
  font-size: 12px;
  color: #595959;
  transition: transform 0.2s ease;

  &.open {
    transform: rotate(180deg);
  }
}

.chat-model-panel {
  position: absolute;
  left: 0;
  bottom: calc(100% + 12px);
  z-index: 20;
  width: 520px;
  height: 280px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  background: #fff;
  box-shadow: 0 12px 40px -8px rgba(15, 23, 42, 0.18);
  box-sizing: border-box;
}

.panel-search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px 11px;
  border-bottom: 1px solid #f0f0f0;
}

.panel-search-icon {
  flex-shrink: 0;
  font-size: 14px;
  color: #8c8c8c;
}

.panel-search-input {
  flex: 1;
  min-width: 0;
  border: 0;
  outline: none;
  background: transparent;
  color: #262626;
  font-size: 14px;
  line-height: 22px;

  &::placeholder {
    color: #8c8c8c;
  }
}

.panel-content {
  display: flex;
  align-items: stretch;
  height: calc(100% - 43px);
}

.provider-list-wrapper {
  width: 152px;
  height: 100%;
  overflow-y: auto;
  background: #f2f4f7;
  box-sizing: border-box;
  scrollbar-width: thin;
  scrollbar-color: #fff;

  &::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    width: 1px;
    background: #f1f5f9;
    pointer-events: none;
  }
  
  &::-webkit-scrollbar {
    width: 4px;
    background: #ffffff;
  }

  &::-webkit-scrollbar-track {
    background: #ffffff;
  }

  &::-webkit-scrollbar-thumb {
    border-radius: 999px;
    background: transparent;
  }

  &:hover {
    scrollbar-color: #FFF;
  }

  &:hover::-webkit-scrollbar-thumb {
    background: #d0d7e2;
  }

  &:hover::-webkit-scrollbar-thumb:hover {
    background: #bcc6d4;
  }
}

.provider-list {
  position: relative;
  min-height: 100%;
  padding: 6px 0;
  box-sizing: border-box;
}

.provider-item {
  position: relative;
  width: 100%;
  padding: 8px 12px;
  border: 0;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.2s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.72);
  }

  &.active {
    background: #fff;
  }

  &.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 8px;
    width: 2px;
    height: 20px;
    border-radius: 0 4px 4px 0;
    background: #155dfc;
  }
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.provider-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.provider-icon.fallback {
  border-radius: 50%;
  background: #dbe6ff;
  color: #2475fc;
  font-size: 12px;
  line-height: 20px;
  text-align: center;
}

.provider-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

.provider-item.active .provider-name {
  color: #262626;
}

.provider-count {
  flex-shrink: 0;
  color: #7a8699;
  font-size: 12px;
  line-height: 20px;
}

.model-section {
  flex: 1;
  min-width: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.model-list {
  flex: 1;
  overflow-y: auto;
  padding: 6px 4px 0;
  box-sizing: border-box;
  scrollbar-width: thin;
}

.model-item {
  width: calc(100% - 8px);
  margin: 2px 4px;
  padding: 8px 12px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    background: #f6f9ff;
  }

  &.active {
    background: #e5efff;
  }
}

.model-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #1d293d;
  font-size: 14px;
  line-height: 22px;
}

.model-item.active .model-name,
.model-check {
  color: #2475fc;
}

.model-check {
  flex-shrink: 0;
  font-size: 14px;
}

.model-empty {
  padding: 12px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.panel-footer {
  border-top: 1px solid #f1f5f9;
  padding: 5px 4px 4px;
}

.add-model-btn {
  width: 100%;
  height: 36px;
  padding: 0 12px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #155dfc;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    opacity: 0.8;
  }
}

.add-model-icon {
  font-size: 14px;
}
</style>
