<template>
  <div class="web-app-page">
    <CustomerPageHeader v-model:active-key="activeTab" :title="t('web_app')" :tabs="headerTabs" />

    <div v-if="loadFailed" class="page-loading">
      <a-empty :description="t('error_robot')">
        <a-button type="primary" @click="loadRobotConfig">{{ t('retry') }}</a-button>
      </a-empty>
    </div>

    <div v-else-if="!ready" class="page-loading">
      <a-spin size="large" />
    </div>

    <div v-else-if="activeTab === 'embed'" class="embed-body">
      <EmbedWebsite variant="workspace" />
    </div>

    <div v-else class="web-app-body">
      <div class="config-column" :class="{ 'is-scrolling': configScrolling }" @scroll="onConfigScroll">
        <section class="link-banner">
          <h2>
            <span>{{ t('web_link') }}</span>
            <span class="banner-tip">{{ t('web_link_tip') }}</span>
          </h2>
          <div class="link-row">
            <a-input class="link-input" :value="h5Link" readonly>
              <template #prefix><GlobalOutlined /></template>
            </a-input>
            <a-button @click="copyLink">
              <CopyOutlined />{{ t('copy') }}
            </a-button>
            <a-button @click="previewLink">
              <EyeOutlined />{{ t('preview') }}
            </a-button>
            <a-popover trigger="click" placement="bottom">
              <template #content>
                <img v-if="qrcode" class="qrcode" :src="qrcode" :alt="t('qr_code')" />
              </template>
              <a-button>
                <QrcodeOutlined />{{ t('qr_code') }}
              </a-button>
            </a-popover>
          </div>
        </section>

        <section class="setting-section">
          <h2 class="section-bar">{{ t('display_settings') }}</h2>
          <div class="section-body">
            <div class="setting-row">
              <span class="setting-label">{{ t('page_title') }}</span>
              <div class="setting-inline">
                <a-switch :checked="formState.navbarShow === 1" @change="(checked) => (formState.navbarShow = checked ? 1 : 2)" />
                <PageTitleInput
                  v-if="formState.navbarShow === 1"
                  v-model:value="formState.pageTitle"
                  v-model:avatar="formState.logo"
                  size="middle"
                  class="setting-control"
                />
                <ColorPicker
                  v-if="formState.navbarShow === 1"
                  v-model:value="formState.pageStyle.navbarBackgroundColor"
                  class="setting-color"
                />
              </div>
            </div>

            <div class="setting-row">
              <span class="setting-label">{{ t('bottom_tip') }}</span>
              <div class="setting-inline">
                <a-switch
                  :checked="formState.ai_generated_tip_show === 1"
                  @change="(checked) => (formState.ai_generated_tip_show = checked ? 1 : 2)"
                />
                <a-input
                  v-if="formState.ai_generated_tip_show === 1"
                  v-model:value="formState.ai_generated_tip"
                  class="setting-control"
                  :maxlength="15"
                  :placeholder="t('tip_content_placeholder')"
                />
              </div>
            </div>

            <div class="setting-row">
              <span class="setting-label">{{ t('page_display') }}</span>
              <div class="setting-check-group">
                <a-checkbox
                  :checked="formState.avatarShow === 1"
                  @change="(e) => (formState.avatarShow = e.target.checked ? 1 : 2)"
                >
                  {{ t('session_avatar') }}
                </a-checkbox>
                <a-checkbox
                  :checked="formState.new_session_btn_show === 1"
                  @change="(e) => (formState.new_session_btn_show = e.target.checked ? 1 : 2)"
                >
                  {{ t('history_and_new_session') }}
                </a-checkbox>
              </div>
            </div>

            <div class="setting-row">
              <span class="setting-label">{{ t('access_restriction') }}</span>
              <a-radio-group v-model:value="formState.accessRestrictionsType">
                <a-radio :value="1">{{ t('no_restriction') }}</a-radio>
                <a-radio :value="2">{{ t('login_required') }}</a-radio>
              </a-radio-group>
            </div>
          </div>
        </section>

        <section class="setting-section">
          <h2 class="section-bar">{{ t('language_and_open') }}</h2>
          <div class="section-body">
            <div class="setting-row">
              <span class="setting-label">{{ t('language') }}</span>
              <a-select v-model:value="formState.lang" class="setting-control short">
                <a-select-option value="zh-CN">{{ t('language_zh_cn') }}</a-select-option>
                <a-select-option value="en-US">{{ t('language_en_us') }}</a-select-option>
                <a-select-option value="auto">{{ t('language_auto') }}</a-select-option>
              </a-select>
            </div>
            <div class="setting-row">
              <span class="setting-label">{{ t('open_method') }}</span>
              <a-radio-group v-model:value="formState.open_type">
                <a-radio :value="1">{{ t('open_new_tab') }}</a-radio>
                <a-radio :value="2">{{ t('open_new_window') }}</a-radio>
              </a-radio-group>
            </div>
            <div v-if="formState.open_type === 2" class="setting-row">
              <span class="setting-label">{{ t('window_size') }}</span>
              <a-space>
                <a-input-number v-model:value="formState.window_width" :min="500" :max="2000" />
                <span>×</span>
                <a-input-number v-model:value="formState.window_height" :min="500" :max="2000" />
                PX
              </a-space>
            </div>
          </div>
        </section>

        <QuickInstruction appearance="workspace" :type="robotDetail.app_id" @updata="onQuickCommandUpdated" />
      </div>

      <div class="preview-column">
        <div class="phone-shell">
          <iframe
            id="customer-web-preview"
            :src="previewIframeSrc"
            class="phone-iframe"
            frameborder="0"
            @load="syncPreview"
          ></iframe>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onUnmounted, reactive, ref, toRaw, watch } from 'vue'
import { storeToRefs } from 'pinia'
import {
  CopyOutlined, EyeOutlined, GlobalOutlined, QrcodeOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import QRCode from 'qrcode'
import { useI18n } from '@/hooks/web/useI18n'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
import { useRobotStore } from '@/stores/modules/robot'
import { editExternalConfig } from '@/api/robot/index'
import { copyText } from '@/utils/index'
import { toPageProtocol } from '@/utils/preview-url'
import PageTitleInput from '@/views/robot/robot-config/external-service/components/page-title-input.vue'
import ColorPicker from '@/components/color-picker/index.vue'
import QuickInstruction from '@/views/robot/robot-config/external-service/components/quick-instruction.vue'
import EmbedWebsite from '@/views/robot/robot-config/external-service/components/embed-website.vue'
import CustomerPageHeader from '../components/page-header.vue'

const { t } = useI18n('views.customer-service.index')
const customerServiceStore = useCustomerServiceStore()
const robotStore = useRobotStore()
const { robotInfo: robotDetail } = storeToRefs(robotStore)
const { getRobot } = robotStore

const activeTab = ref('webapp')
const headerTabs = computed(() => [
  { key: 'webapp', label: 'webAPP' },
  { key: 'embed', label: t('embed_website') }
])
const configScrolling = ref(false)
let configScrollTimer
const onConfigScroll = () => {
  configScrolling.value = true
  clearTimeout(configScrollTimer)
  configScrollTimer = setTimeout(() => { configScrolling.value = false }, 800)
}
const qrcode = ref('')
const ready = ref(false)
const loadFailed = ref(false)
const customerRobot = computed(() => customerServiceStore.robot)

const defaultConfig = {
  accessRestrictionsType: 1,
  logo: '',
  pageTitle: '',
  avatarShow: 1,
  navbarShow: 1,
  ai_generated_tip_show: 1,
  ai_generated_tip: '内容由AI生成，仅供参考',
  lang: 'zh-CN',
  pageStyle: { navbarBackgroundColor: '#2475FC' },
  open_type: 1,
  window_width: 1200,
  window_height: 650,
  new_session_btn_show: 2
}
const formState = reactive(JSON.parse(JSON.stringify(defaultConfig)))

const h5Link = computed(() => (
  robotDetail.value?.h5_domain
    ? `${robotDetail.value.h5_domain}/#/chat/h5?robot_key=${robotDetail.value.robot_key || ''}`
    : ''
))

const previewIframeSrc = computed(() => toPageProtocol(robotDetail.value?.h5_website))

const syncFormFromStore = () => {
  const source = JSON.parse(JSON.stringify(robotStore.external_config_h5 || {}))
  Object.assign(formState, JSON.parse(JSON.stringify(defaultConfig)), source)
  formState.pageStyle = {
    ...defaultConfig.pageStyle,
    ...(source.pageStyle || {})
  }
}

const generateQrcode = async () => {
  if (!h5Link.value) {
    qrcode.value = ''
    return
  }
  try {
    qrcode.value = await QRCode.toDataURL(h5Link.value)
  } catch {
    qrcode.value = ''
  }
}

const copyLink = () => {
  copyText(h5Link.value)
  message.success(t('copy_success'))
}

const previewLink = () => {
  if (h5Link.value) window.open(h5Link.value, '_blank', 'noopener')
}

let previewTimer = null
const postPreview = () => {
  const iframe = document.getElementById('customer-web-preview')
  if (!iframe?.contentWindow) return
  iframe.contentWindow.postMessage({
    type: 'onPreview',
    data: JSON.parse(JSON.stringify(toRaw(formState)))
  }, '*')
}

const syncPreview = () => {
  if (previewTimer) clearTimeout(previewTimer)
  previewTimer = setTimeout(postPreview, 120)
}

const onQuickCommandUpdated = (data) => {
  const iframe = document.getElementById('customer-web-preview')
  if (!iframe?.contentWindow) return
  iframe.contentWindow.postMessage({
    type: 'updataQuickComand',
    data: JSON.parse(JSON.stringify(data))
  }, '*')
}

let autoSaveTimer = null
let autoSaveEnabled = false
let saveInFlight = false
let saveQueued = false
let tipWarningShown = false

const buildPayload = () => {
  const payload = JSON.parse(JSON.stringify(toRaw(formState)))
  payload.window_width = Number(payload.window_width) || 1200
  payload.window_height = Number(payload.window_height) || 650
  payload.ai_generated_tip = payload.ai_generated_tip?.trim() || ''
  return payload
}

const saveWebApp = async () => {
  const id = robotDetail.value?.id
  if (!ready.value || !id || saveInFlight) {
    if (saveInFlight) saveQueued = true
    return
  }
  if (formState.ai_generated_tip_show === 1 && !formState.ai_generated_tip?.trim()) {
    if (!tipWarningShown) {
      tipWarningShown = true
      message.warning(t('tip_content_required'))
    }
    return
  }
  tipWarningShown = false
  saveInFlight = true
  const payload = buildPayload()
  try {
    await editExternalConfig({ id, external_config_h5: JSON.stringify(payload) })
    if (String(robotDetail.value?.id) === String(id)) {
      Object.assign(robotStore.external_config_h5, payload)
    }
  } catch {
    message.error(t('save_failed'))
  } finally {
    saveInFlight = false
    if (saveQueued) {
      saveQueued = false
      scheduleAutoSave(0)
    }
  }
}

const scheduleAutoSave = (delay = 500) => {
  if (!autoSaveEnabled || !ready.value) return
  if (saveInFlight) {
    saveQueued = true
    return
  }
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    autoSaveTimer = null
    saveWebApp()
  }, delay)
}

const flushAutoSave = () => {
  if (!autoSaveTimer) return
  clearTimeout(autoSaveTimer)
  autoSaveTimer = null
  saveWebApp()
}

const resetAutoSave = () => {
  autoSaveEnabled = false
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = null
  saveQueued = false
}

const loadRobotConfig = async () => {
  const id = customerRobot.value?.id
  resetAutoSave()
  ready.value = false
  loadFailed.value = false
  if (!id) return
  try {
    const loaded = await getRobot(id)
    if (!loaded || String(robotDetail.value?.id) !== String(id)) {
      loadFailed.value = true
      return
    }
    syncFormFromStore()
    await generateQrcode()
    ready.value = true
    await nextTick()
    autoSaveEnabled = true
    syncPreview()
  } catch {
    loadFailed.value = true
  }
}

watch(
  () => customerRobot.value?.id,
  () => loadRobotConfig(),
  { immediate: true }
)

watch(
  formState,
  () => {
    syncPreview()
    scheduleAutoSave()
  },
  { deep: true }
)
watch(h5Link, generateQrcode)
watch(activeTab, (tab, previous) => {
  if (previous === 'webapp' && tab === 'embed') flushAutoSave()
})

onUnmounted(() => {
  flushAutoSave()
  clearTimeout(configScrollTimer)
  if (previewTimer) clearTimeout(previewTimer)
})
</script>

<style lang="less" scoped>
.web-app-page {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  margin: 8px 8px 8px 0;
  overflow: hidden;
  border-radius: 12px;
  background: #fff;
}
.page-loading {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
}
.embed-body,
.web-app-body {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
.web-app-body {
  display: flex;
}
.config-column {
  flex: 1;
  min-width: 0;
  padding: 0 28px 28px;
  overflow-x: hidden;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;

  &::-webkit-scrollbar { width: 6px; height: 6px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb { background: transparent; border-radius: 4px; }
  &::-webkit-scrollbar-button { display: none; width: 0; height: 0; }
  &.is-scrolling {
    scrollbar-color: #d1d5db transparent;
    &::-webkit-scrollbar-thumb { background: #d1d5db; }
  }
  @supports selector(::-webkit-scrollbar) {
    scrollbar-width: auto;
    scrollbar-color: auto;
    &.is-scrolling { scrollbar-color: auto; }
  }
}
.preview-column {
  display: flex;
  flex-shrink: 0;
  align-items: flex-start;
  justify-content: center;
  width: 400px;
  padding: 24px 24px 24px 0;
  overflow: auto;
}
.link-banner {
  margin-top: 24px;
  padding: 20px;
  border-radius: 8px;
  background: var(--app-selected);

  h2 {
    display: flex;
    align-items: baseline;
    gap: 12px;
    margin: 0 0 12px;
    color: var(--app-text-primary);
    font-size: 16px;
    font-weight: 500;
    line-height: 24px;
  }
}
.banner-tip {
  color: var(--app-text-secondary);
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
}
.link-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.link-input {
  flex: 1;
  min-width: 220px;
}
.qrcode {
  display: block;
  width: 180px;
  height: 180px;
}
.section-bar {
  display: flex;
  align-items: center;
  height: 40px;
  margin: 32px 0 0;
  padding: 0 20px;
  border-radius: 8px;
  background: #f5f5f5;
  color: var(--app-text-primary);
  font-size: 14px;
  font-weight: 500;
  line-height: 22px;
}
.section-body {
  padding: 24px 0 0 20px;
}
.setting-row {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 28px;

  &:last-child {
    margin-bottom: 0;
  }

}
.setting-label {
  flex: 0 0 112px;
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 22px;
}
.setting-control {
  flex: 1;
  min-width: 0;
}
.setting-control.short {
  flex: 0 1 220px;
  max-width: 220px;
}
.setting-inline {
  display: flex;
  flex: 1;
  align-items: center;
  gap: 12px;
  min-width: 0;
}
.setting-inline :deep(.page-title-input) {
  min-width: 0;
}
.setting-inline :deep(.page-title-input .avatar) {
  width: 32px;
  height: 32px;
  border-radius: 16px;
  object-fit: cover;
}
.setting-color {
  flex: 0 0 140px;
  width: 140px;
}
.setting-color :deep(.color-picker-wrapper) {
  display: block;
  width: 100%;
}
.setting-color :deep(.color-picker-wrapper .ant-input) {
  padding-right: 32px;
}
.setting-check-group {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px 20px;
  min-height: 32px;
  min-width: 0;
}
.phone-shell {
  display: flex;
  flex-direction: column;
  width: 352px;
  height: min(720px, calc(100vh - 140px));
  height: min(720px, calc(100dvh - 140px));
  border: 8px solid #fff;
  border-radius: 30px;
  background: #f0f2f5;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.06), 0 4px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
.phone-iframe {
  flex: 1;
  width: 100%;
  min-height: 0;
  border: 0;
  border-radius: 22px;
  background: #fff;
}
@media (max-width: 1180px) {
  .preview-column {
    width: 340px;
    padding-right: 16px;
  }
  .phone-shell {
    width: 320px;
  }
}
@media (max-width: 1023px), (max-height: 620px) {
  .web-app-body {
    flex-direction: column;
    overflow-y: auto;
  }
  .config-column,
  .embed-body {
    overflow: visible;
  }
  .preview-column {
    width: 100%;
    padding: 0 0 24px;
  }
  .phone-shell {
    height: min(560px, calc(100vh - 160px));
    height: min(560px, calc(100dvh - 160px));
  }
}
</style>
