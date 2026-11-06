<template>
  <div class="system-mode-setting" :class="{ onboarding }">
    <div v-if="modeStore.storageStatus === 'loading'" class="mode-loading">
      <a-spin size="small" />
      <span>{{ t('loading') }}</span>
    </div>

    <div v-if="modeStore.storageStatus === 'error'" class="mode-error">
      <a-alert type="error" show-icon :message="t(modeStore.storageError)" />
      <a-button @click="loadMode">{{ t('retry') }}</a-button>
    </div>

    <div class="mode-options" role="radiogroup" :aria-label="t('title')">
      <label
        v-for="option in modeOptions"
        :key="option.value"
        class="mode-card"
        :class="{
          selected: selectedMode === option.value,
          disabled: modeStore.storageStatus !== 'ready' || isSwitching || isConfirmOpen
        }"
      >
        <input
          :checked="selectedMode === option.value"
          class="mode-radio"
          type="radio"
          name="system-mode"
          :value="option.value"
          :disabled="modeStore.storageStatus !== 'ready' || isSwitching || isConfirmOpen"
          @change="selectMode(option.value)"
        />
        <span class="mode-icon"><component :is="option.icon" /></span>
        <CheckCircleFilled v-if="selectedMode === option.value" class="selected-icon" />
        <span class="mode-title">
          {{ t(option.label) }}
          <span v-if="option.value === SYSTEM_MODES.AI_SERVICE" class="recommended-tag">
            {{ t('recommended') }}
          </span>
        </span>
        <span class="mode-description">{{ t(option.description) }}</span>
      </label>
    </div>

    <div v-if="actionError" class="mode-error action-error">
      <a-alert type="error" show-icon :message="t(actionError)" />
      <a-button v-if="actionError === 'error_navigation'" @click="enterHome">
        {{ t('enter_home') }}
      </a-button>
    </div>

    <div v-if="onboarding" class="start-action">
      <a-button type="primary" :loading="isSwitching" :disabled="modeStore.storageStatus !== 'ready'" @click="applyMode">{{ t('start_now') }}</a-button>
    </div>

    <a-modal
      v-model:open="isConfirmOpen"
      :title="t('confirm_title')"
      :ok-text="t('confirm_ok')"
      :cancel-text="t('confirm_cancel')"
      :confirm-loading="isSwitching"
      :cancel-button-props="{ disabled: isSwitching }"
      :closable="!isSwitching"
      :keyboard="!isSwitching"
      :mask-closable="false"
      @ok="applyMode"
      @cancel="cancelMode"
    >
      <p>{{ t('confirm_content', { mode: selectedModeLabel }) }}</p>
    </a-modal>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  CheckCircleFilled,
  CommentOutlined,
  UserOutlined,
  ReadOutlined
} from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useCompanyStore } from '@/stores/modules/company'
import { SYSTEM_MODES, useSystemModeStore } from '@/stores/modules/system-mode'

const { t } = useI18n('views.system-mode')
const props = defineProps({ onboarding: { type: Boolean, default: false } })
const router = useRouter()
const companyStore = useCompanyStore()
const modeStore = useSystemModeStore()
const selectedMode = ref(props.onboarding ? SYSTEM_MODES.AI_SERVICE : SYSTEM_MODES.ADVANCED)
const isConfirmOpen = ref(false)
const isSwitching = ref(false)
const actionError = ref('')
let active = true
let loadSequence = 0

const modeOptions = [
  {
    value: SYSTEM_MODES.AI_SERVICE,
    label: 'ai_service',
    description: 'description_ai_service',
    icon: CommentOutlined
  },
  {
    value: SYSTEM_MODES.AGENT,
    label: 'agent',
    description: 'description_agent',
    icon: UserOutlined
  },
  {
    value: SYSTEM_MODES.ADVANCED,
    label: 'advanced',
    description: 'description_advanced',
    icon: ReadOutlined
  }
]

const selectedModeLabel = computed(() => {
  const selectedOption = modeOptions.find((option) => option.value === selectedMode.value)
  return t(selectedOption?.label || 'advanced')
})

const loadMode = async () => {
  const version = ++loadSequence
  isConfirmOpen.value = false
  actionError.value = ''
  const loaded = await modeStore.loadForScope({ companyId: companyStore.id })
  if (!active || version !== loadSequence) return
  if (props.onboarding && loaded && !modeStore.needsModeSelection) { enterHome(); return }
  selectedMode.value = props.onboarding ? SYSTEM_MODES.AI_SERVICE : modeStore.currentMode
}

const enterHome = async () => {
  if (!active) return
  isSwitching.value = true
  try {
    if (props.onboarding) await router.replace('/')
    else window.location.replace(router.resolve('/').href)
  } catch {
    // 本地值已保存；导航失败不能提示成保存失败，也不应回滚生效模式。
    actionError.value = 'error_navigation'
    isSwitching.value = false
  }
}

const selectMode = (mode) => {
  if (modeStore.storageStatus !== 'ready' || isSwitching.value || isConfirmOpen.value) return
  if (!props.onboarding && mode === modeStore.currentMode) return

  selectedMode.value = mode
  actionError.value = ''
  if (!props.onboarding) isConfirmOpen.value = true
}

const cancelMode = () => {
  if (isSwitching.value) return
  isConfirmOpen.value = false
  selectedMode.value = modeStore.currentMode
}

const applyMode = async () => {
  if ((!props.onboarding && !isConfirmOpen.value) || isSwitching.value) return
  if (modeStore.storageStatus !== 'ready' || (!props.onboarding && selectedMode.value === modeStore.currentMode)) {
    isConfirmOpen.value = false
    if (!props.onboarding) selectedMode.value = modeStore.currentMode
    return
  }

  isSwitching.value = true
  actionError.value = ''
  try {
    await modeStore.saveMode(selectedMode.value)
  } catch {
    if (!active) return
    actionError.value = 'error_storage_write'
    isConfirmOpen.value = false
    if (!props.onboarding) selectedMode.value = modeStore.currentMode
    isSwitching.value = false
    return
  }
  isConfirmOpen.value = false
  enterHome()
}

watch(() => companyStore.id, loadMode, { immediate: true })
onBeforeUnmount(() => { active = false })
</script>

<style lang="less" scoped>
.system-mode-setting {
  min-height: 240px;
  padding-top: 8px;
}
.start-action { margin-top: 40px; text-align: center; .ant-btn { min-width: 196px; height: 42px; border-radius: 999px; font-size: 16px; font-weight: 600; } }
.onboarding .mode-card { min-height: 266px; padding: 30px; background: var(--app-bg-container); }
.onboarding .mode-card.selected { box-shadow: var(--app-shadow-card); }
.mode-options {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
}
.mode-card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 220px;
  padding: 26px;
  border: 2px solid var(--app-border);
  border-radius: 20px;
  background: var(--app-bg-layout);
  cursor: pointer;
  transition: border-color 200ms, background-color 200ms, box-shadow 200ms;

  &:not(.disabled):hover {
    border-color: var(--app-primary);
  }
  &.selected {
    border-color: var(--app-primary);
    background: var(--app-bg-container);
  }
  &:focus-within {
    outline: 2px solid var(--app-primary);
    outline-offset: 3px;
  }
  &.disabled {
    cursor: not-allowed;
    opacity: .6;
  }
}
.mode-radio {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
}
.mode-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  margin-bottom: 16px;
  color: var(--app-primary);
  background: var(--app-selected);
  border-radius: 16px;
  font-size: 27px;
}
.selected-icon {
  position: absolute;
  top: 20px;
  right: 20px;
  color: var(--app-primary);
  font-size: 26px;
}
.mode-title {
  margin-bottom: 12px;
  color: var(--app-text-primary);
  font-size: 18px;
  font-weight: 600;
}
.recommended-tag {
  display: inline-block;
  margin-left: 6px;
  padding: 1px 8px;
  border-radius: 5px;
  color: var(--app-text-on-primary);
  background: var(--app-primary-active);
  font-size: 12px;
  vertical-align: middle;
}
.mode-description {
  color: var(--app-text-secondary);
  font-size: 14px;
  line-height: 24px;
}
.mode-loading {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  color: var(--app-text-secondary);
  font-size: 14px;
}
.mode-error {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;

  .ant-alert {
    flex: 1;
  }
}
.action-error {
  margin: 20px 0 0;
}
.system-mode-setting:not(.onboarding) {
  .mode-options {
    gap: 16px;
  }
  .mode-card {
    min-height: 180px;
    padding: 20px;
  }
  .mode-icon {
    width: 48px;
    height: 48px;
    margin-bottom: 14px;
    border-radius: 14px;
    font-size: 24px;
  }
  .selected-icon {
    top: 16px;
    right: 16px;
    font-size: 22px;
  }
  .mode-title {
    margin-bottom: 8px;
    font-size: 16px;
  }
  .mode-description {
    font-size: 13px;
    line-height: 22px;
  }
}
@media (max-width: 1200px) {
  .mode-options {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
@media (max-width: 800px) {
  .mode-options {
    grid-template-columns: 1fr;
  }
  .mode-card {
    min-height: 0;
    padding: 24px;
  }
  .onboarding .mode-card {
    min-height: 0;
    padding: 24px;
  }
}
</style>
