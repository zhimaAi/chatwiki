<style lang="less" scoped>
.web-app-box {
  display: flex;

  .box-left {
    flex: 1;

    .box-wrapper {
      margin-bottom: 16px;
      &:last-child {
        margin-bottom: 0;
      }
    }
  }
  .web-app-info {
    .web-app-link {
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }
    .link-action {
      margin-top: 8px;
      .action-btn {
        margin-right: 8px;
      }
    }
    .access-restrictions {
      margin-top: 16px;
    }
  }

  .form-box {
    .form-item {
      margin-bottom: 16px;
    }
  }
}
.box-right {
  width: 418px;
  margin: 0 96px 0 48px;

  .demo-box {
    overflow: hidden;
  }
  .preview-img {
    display: block;
    border-radius: 9px;
    box-shadow: 0 4px 32px 0 rgba(0, 0, 0, 0.16);
  }
  iframe {
    display: block;
    transform-origin: top left;
    border-radius: 4px;
    box-shadow: 0 4px 32px 0 rgba(0, 0, 0, 0.16);
  }
}

.window-size-box {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 14px;
}

.iframe-setting-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.iframe-switch-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  color: #595959;
}

.iframe-checkbox {
  display: flex;
  margin-top: 14px;
}

.unit {
  color: var(--app-text-tertiary);
  font-size: 12px;
}

.page-display-group {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px 20px;
  min-height: 32px;
}

.web-app-box.is-workspace {
  align-items: stretch;
  gap: 0;
  height: 100%;
  min-height: 0;

  .box-left {
    display: flex;
    flex-direction: column;
    height: 100%;
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
  .box-right {
    display: flex;
    flex: none;
    flex-direction: column;
    align-items: center;
    width: 400px;
    height: 100%;
    margin: 0;
    padding: 24px 16px;
    overflow: auto;
  }
  .workspace-form {
    display: contents;
  }
  .embed-banner { order: 0; }
  .style-card { order: 1; }
  .float-card { order: 2; }
  .language-card { order: 3; }
  .quick-card { order: 4; }
}
.embed-banner {
  margin-top: 24px;
  padding: 24px;
  border-radius: 12px;
  background: var(--app-selected);
}
.embed-banner-head {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 20px;

  h2 {
    margin: 0;
    color: var(--app-text-primary);
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }
  p {
    margin: 0;
    color: var(--app-text-secondary);
    font-size: 12px;
    line-height: 20px;
  }
}
.embed-code-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.embed-code-field {
  display: flex;
  flex: 1;
  align-items: center;
  min-width: 0;
  height: 40px;
  padding: 0 12px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: var(--app-bg-container);
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;

  &:hover {
    border-color: var(--app-primary);
  }
  &:focus-visible {
    outline: none;
    border-color: var(--app-primary);
    box-shadow: var(--app-shadow-focus);
  }
}
.embed-code-text {
  overflow: hidden;
  color: var(--app-text-primary);
  font-family: Menlo, ui-monospace, monospace;
  font-size: 13px;
  line-height: 20px;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.embed-code-tip {
  margin: 0 0 16px;
  color: var(--app-text-secondary);
  font-size: 14px;
  line-height: 22px;
}
.embed-code-block {
  max-height: 360px;
  margin: 0;
  padding: 0;
  overflow: auto;
  background: transparent;
  color: var(--app-text-primary);
  font-family: Menlo, ui-monospace, monospace;
  font-size: 13px;
  line-height: 18px;
  white-space: pre-wrap;
  word-break: break-all;
}
.embed-code-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
}
.web-app-box.is-workspace .preview-meta {
  margin-top: 10px;
  color: var(--app-text-tertiary);
  font-size: 12px;
  line-height: 20px;
  text-align: center;
}
.web-app-box.is-workspace :deep(.card-icon) {
  display: none;
}
.web-app-box.is-workspace :deep(.card-box) {
  margin: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  overflow: visible;
}
.web-app-box.is-workspace :deep(.card-box-header) {
  height: 40px;
  margin-top: 16px;
  padding: 0 20px;
  border-radius: 8px;
  background: #f5f5f5;
}
.web-app-box.is-workspace :deep(.card-box-body) {
  padding: 24px 0 0 20px;
}
.web-app-box.is-workspace :deep(.card-box .title-text) {
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 22px;
}
.web-app-box.is-workspace :deep(.page-title-input .avatar) {
  width: 32px;
  height: 32px;
  border-radius: 16px;
  object-fit: cover;
}
.web-app-box.is-workspace .iframe-setting-row :deep(.ant-input-number) {
  width: 72px;
}
.web-app-box.is-workspace .language-card :deep(.ant-select) {
  width: 220px !important;
}
.web-app-box.is-workspace .box-left .box-wrapper {
  margin-bottom: 0;
}
.web-app-box.is-workspace :deep(.ant-form-item) {
  margin-bottom: 28px;
}
.web-app-box.is-workspace .style-card :deep(.ant-form-item:last-child),
.web-app-box.is-workspace .language-card :deep(.ant-form-item:last-child) {
  margin-bottom: 0;
}
.web-app-box.is-workspace :deep(.ant-form-item-row) {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
}
.web-app-box.is-workspace :deep(.ant-form-item-label) {
  flex: 0 0 112px;
  padding: 0 24px 0 0;
  text-align: left;
}
.web-app-box.is-workspace :deep(.ant-form-item-label > label) {
  height: auto;
  min-height: 32px;
  color: var(--app-text-primary);
}
.web-app-box.is-workspace :deep(.no-required-mark .ant-form-item-label > label.ant-form-item-required::before),
.web-app-box.is-workspace :deep(.no-required-mark .ant-form-item-label > label::before),
.web-app-box.is-workspace :deep(.no-required-mark .ant-form-item-label > label::after) {
  display: none;
}
.web-app-box.is-workspace :deep(.ant-form-item-control) {
  flex: 1;
  min-width: 0;
}
.web-app-box.is-workspace .demo-box {
  flex-shrink: 0;
  box-sizing: content-box;
  border: 8px solid #fff;
  border-radius: 30px;
  background: #f0f2f5;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.06), 0 4px 24px rgba(0, 0, 0, 0.1);
}
.web-app-box.is-workspace .box-right iframe {
  border-radius: 0;
  background: #fff;
  box-shadow: none;
}
@media (max-width: 1180px) {
  .web-app-box.is-workspace .box-right {
    width: 340px;
  }
}
@media (max-width: 1023px), (max-height: 620px) {
  .web-app-box.is-workspace {
    flex-direction: column;
    height: auto;
    overflow: visible;
  }
  .web-app-box.is-workspace .box-left {
    width: 100%;
    height: auto;
    overflow: visible;
  }
  .web-app-box.is-workspace .box-right {
    display: flex;
    width: 100%;
    height: auto;
    margin-top: 16px;
    justify-content: center;
  }
}
</style>

<style lang="less">
/* 弹窗外壳被 Teleport 到 body，标题/内边距/关闭按钮只能用全局样式，
   通过 wrap-class-name="embed-code-modal" 限定作用域。 */
.embed-code-modal {
  .ant-modal-content {
    padding: 24px;
    border-radius: 16px;
  }
  .ant-modal-header {
    margin-bottom: 12px;
    background: transparent;
  }
  .ant-modal-title {
    color: var(--app-text-primary);
    font-size: 24px;
    font-weight: 600;
    line-height: 32px;
  }
  .ant-modal-close {
    top: 20px;
    inset-inline-end: 20px;
    color: var(--app-text-tertiary);

    &:hover {
      color: var(--app-text-primary);
    }
  }
}
</style>

<template>
  <div class="web-app-box" :class="{ 'is-workspace': isWorkspace }">
    <div class="box-left" :class="{ 'is-scrolling': leftScrolling }" @scroll="onLeftScroll">
      <section v-if="isWorkspace" class="embed-banner">
        <div class="embed-banner-head">
          <h2>{{ t('embed_code_title') }}</h2>
          <p>{{ t('copy_sdk_code_title') }}</p>
        </div>
        <div class="embed-code-row">
          <div
            class="embed-code-field"
            role="button"
            tabindex="0"
            :aria-label="t('embed_code_title')"
            @click="openCodeModal"
            @keydown.enter.prevent="openCodeModal"
            @keydown.space.prevent="openCodeModal"
          >
            <span class="embed-code-text">{{ sdkCode }}</span>
          </div>
          <a-button size="large" @click="copySDKCode">
            <CopyOutlined />{{ copied ? t('copied_btn') : t('copy_btn') }}
          </a-button>
        </div>
      </section>
      <a-form
        ref="formRef"
        class="workspace-form"
        layout="vertical"
        :model="formState"
        :rules="formRules"
      >
        <div class="box-wrapper style-card">
          <card-box :title="t('style_settings_title')">
            <template #icon>
              <svg-icon name="pc" style="font-size: 16px; color: #262626"></svg-icon>
            </template>
            <template #action>
              <a-button v-if="!isWorkspace" @click="saveForm" size="small" type="primary">
                {{ t('save_btn') }}
              </a-button>
            </template>
            <div class="web-app-style form-box">
              <a-form-item class="form-item no-required-mark" :label="t('page_title_label')" name="headTitle">
                <PageTitleInput
                  v-model:avatar="formState.headImage"
                  v-model:value="formState.headTitle"
                  :size="isWorkspace ? 'middle' : 'large'"
                />
              </a-form-item>
              <a-form-item
                class="form-item"
                :label="isWorkspace ? t('product_intro_label') : t('subtitle_label')"
                name="headSubTitle"
              >
                <a-textarea
                  v-model:value="formState.headSubTitle"
                  :placeholder="t('subtitle_placeholder')"
                  :rows="isWorkspace ? 1 : 3"
                  :auto-size="isWorkspace ? { minRows: 1, maxRows: 8 } : false"
                />
              </a-form-item>
              <a-form-item
                class="form-item no-required-mark"
                :label="isWorkspace ? t('window_background_label') : t('color_label')"
                :name="['pageStyle', 'headBackgroundColor']"
              >
                <GradientColorPicker
                  v-model:value="formState.pageStyle.headBackgroundColor"
                  :compact="isWorkspace"
                />
              </a-form-item>
              <a-form-item class="form-item" :label="isWorkspace ? t('window_size_label') : t('iframe_size_label')">
                <div class="iframe-setting-row">
                  <a-flex align="center" :gap="8">
                    <span>{{ t('iframe_width_label') }}</span>
                    <a-input-number
                      v-model:value="formState.iframe_width"
                      :min="320"
                      :max="2000"
                      :precision="0"
                    />
                    <span v-if="isWorkspace" class="unit">px</span>
                  </a-flex>
                  <a-flex align="center" :gap="8">
                    <span>{{ t('iframe_height_label') }}</span>
                    <a-input-number
                      v-model:value="formState.iframe_height"
                      :min="400"
                      :max="2000"
                      :precision="0"
                    />
                    <span v-if="isWorkspace" class="unit">px</span>
                  </a-flex>
                </div>
                <a-checkbox
                  v-if="isWorkspace"
                  v-model:checked="formState.iframe_resize_enabled"
                  class="iframe-checkbox"
                >
                  {{ t('allow_user_drag') }}
                </a-checkbox>
                <div v-else class="iframe-switch-row">
                  <span>{{ t('iframe_resize_label') }}</span>
                  <a-switch v-model:checked="formState.iframe_resize_enabled" />
                </div>
              </a-form-item>
              <a-form-item class="form-item" :label="isWorkspace ? t('window_position_label') : t('iframe_position_label')">
                <div class="iframe-setting-row">
                  <a-flex align="center" :gap="8">
                    <span>{{ t('iframe_right_label') }}</span>
                    <a-input-number
                      v-model:value="formState.iframe_right"
                      :min="0"
                      :max="2000"
                      :precision="0"
                    />
                    <span v-if="isWorkspace" class="unit">px</span>
                  </a-flex>
                  <a-flex align="center" :gap="8">
                    <span>{{ t('iframe_bottom_label') }}</span>
                    <a-input-number
                      v-model:value="formState.iframe_bottom"
                      :min="0"
                      :max="2000"
                      :precision="0"
                    />
                    <span v-if="isWorkspace" class="unit">px</span>
                  </a-flex>
                </div>
                <a-checkbox
                  v-if="isWorkspace"
                  v-model:checked="formState.iframe_drag_enabled"
                  class="iframe-checkbox"
                >
                  {{ t('allow_user_drag') }}
                </a-checkbox>
                <div v-else class="iframe-switch-row">
                  <span>{{ t('iframe_drag_label') }}</span>
                  <a-switch v-model:checked="formState.iframe_drag_enabled" />
                </div>
              </a-form-item>
              <a-form-item v-if="isWorkspace" class="form-item" :label="t('page_display')">
                <div class="page-display-group">
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
              </a-form-item>
              <template v-if="!isWorkspace">
                <a-form-item class="form-item" :label="t('language_label')" name="lang">
                  <a-select
                    style="width: 180px"
                    v-model:value="formState.lang"
                    :placeholder="t('language_placeholder')"
                  >
                    <a-select-option value="zh-CN">{{ t('language_zh_cn') }}</a-select-option>
                    <a-select-option value="en-US">{{ t('language_en_us') }}</a-select-option>
                    <a-select-option value="auto">自动适配语言</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item class="form-item" :label="t('url_open_type_label')" name="open_type" required>
                  <a-radio-group v-model:value="formState.open_type">
                    <a-radio :value="1">{{ t('open_new_tab') }}</a-radio>
                    <a-radio :value="2">{{ t('open_new_window') }}
                      <a-tooltip :title="t('open_new_window_tooltip')">
                        <QuestionCircleOutlined />
                      </a-tooltip>
                    </a-radio>
                  </a-radio-group>
                  <a-form-item-rest v-if="formState.open_type == 2">
                    <div class="window-size-box">
                      <a-flex align="center" :gap="8">
                        <div>{{ t('window_height_label') }}</div>
                        <a-input-number  v-model:value="formState.window_height" :min="500" :max="2000" />
                        PX
                      </a-flex>
                      <a-flex align="center" :gap="8">
                        <div>{{ t('window_width_label') }}</div>
                        <a-input-number  v-model:value="formState.window_width" :min="500" :max="2000" />
                        PX
                      </a-flex>
                    </div>
                  </a-form-item-rest>
                </a-form-item>
                <a-form-item class="form-item" :label="t('label_show_history_and_new_session_btn')" name="new_session_btn_show" required>
                  <a-radio-group v-model:value="formState.new_session_btn_show" name="new_session_btn_show">
                    <a-radio :value="1">{{ t('label_show') }}</a-radio>
                    <a-radio :value="2">{{ t('label_hide') }}</a-radio>
                  </a-radio-group>
                </a-form-item>
                <a-form-item class="form-item" :label="t('label_show_avatar')" name="avatarShow" required>
                  <a-radio-group v-model:value="formState.avatarShow" name="avatarShow">
                    <a-radio :value="1">{{ t('label_show') }}</a-radio>
                    <a-radio :value="2">{{ t('label_hide') }}</a-radio>
                  </a-radio-group>
                </a-form-item>
              </template>
            </div>
          </card-box>
        </div>

        <div v-if="isWorkspace" class="box-wrapper language-card">
          <card-box :title="t('language_and_open_label')">
            <div class="form-box">
              <a-form-item class="form-item no-required-mark" :label="t('language_label')" name="lang">
                <a-select
                  style="width: 220px"
                  v-model:value="formState.lang"
                  :placeholder="t('language_placeholder')"
                >
                  <a-select-option value="zh-CN">{{ t('language_zh_cn') }}</a-select-option>
                  <a-select-option value="en-US">{{ t('language_en_us') }}</a-select-option>
                  <a-select-option value="auto">自动适配语言</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item class="form-item no-required-mark" :label="t('open_method_label')" name="open_type" required>
                <a-radio-group v-model:value="formState.open_type">
                  <a-radio :value="1">{{ t('open_new_tab') }}</a-radio>
                  <a-radio :value="2">{{ t('open_new_window') }}
                    <a-tooltip :title="t('open_new_window_tooltip')">
                      <QuestionCircleOutlined />
                    </a-tooltip>
                  </a-radio>
                </a-radio-group>
                <a-form-item-rest v-if="formState.open_type == 2">
                  <div class="window-size-box">
                    <a-flex align="center" :gap="8">
                      <div>{{ t('window_height_label') }}</div>
                      <a-input-number  v-model:value="formState.window_height" :min="500" :max="2000" />
                      PX
                    </a-flex>
                    <a-flex align="center" :gap="8">
                      <div>{{ t('window_width_label') }}</div>
                      <a-input-number  v-model:value="formState.window_width" :min="500" :max="2000" />
                      PX
                    </a-flex>
                  </div>
                </a-form-item-rest>
              </a-form-item>
            </div>
          </card-box>
        </div>
      </a-form>

      <div v-if="!isWorkspace" class="box-wrapper sdk-code-card">
        <card-box :title="t('copy_sdk_code_title')">
          <template #icon>
            <svg-icon name="sdk" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <template #action>
            <a-button @click="copySDKCode" size="small">
              {{ copied ? t('copied_btn') : t('copy_btn') }}
            </a-button>
          </template>
          <div class="sdk-code">
            <pre><code>{{ sdkCode }}</code></pre>
          </div>
        </card-box>
      </div>

      <div class="box-wrapper quick-card">
        <QuickInstruction
          :type="robotInfo.app_id_embed"
          :appearance="isWorkspace ? 'workspace' : 'card'"
          @updata="updataQuickComand"
        ></QuickInstruction>
      </div>

      <div class="box-wrapper float-card">
        <FloatIconSetting
          :form="formState"
          :appearance="isWorkspace ? 'workspace' : 'card'"
          @save="handleFloatBtnCongiSave"
        ></FloatIconSetting>
      </div>
    </div>
    <a-modal
      v-if="isWorkspace"
      v-model:open="codeModalOpen"
      :width="640"
      :footer="null"
      centered
      wrap-class-name="embed-code-modal"
    >
      <template #title>{{ t('embed_code_title') }}</template>
      <p class="embed-code-tip">{{ t('copy_sdk_code_title') }}</p>
      <pre class="embed-code-block"><code>{{ sdkCode }}</code></pre>
      <div class="embed-code-actions">
        <a-button @click="codeModalOpen = false">{{ t('close_btn') }}</a-button>
        <a-button type="primary" @click="copySDKCode">
          <CopyOutlined />{{ copied ? t('copied_btn') : t('copy_btn') }}
        </a-button>
      </div>
    </a-modal>
    <div class="box-right">
      <div class="demo-box" :style="previewBoxStyle">
        <iframe
          id="web-preview"
          :src="previewIframeSrc"
          :style="previewIframeStyle"
          frameborder="0"
          @load="handlePreviewLoad"
        ></iframe>
      </div>
      <div v-if="isWorkspace" class="preview-meta">
        {{ previewSize.width }} × {{ previewSize.height }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { getSdkCode } from './sdk-code'
import { nextTick, onBeforeUnmount, onMounted, ref, reactive, toRaw, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useRobotStore } from '@/stores/modules/robot'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'
import { toPageProtocol } from '@/utils/preview-url'
import { editExternalConfig } from '@/api/robot/index'
import CardBox from './card-box.vue'
import PageTitleInput from './page-title-input.vue'
import GradientColorPicker from './gradient-color-picker.vue'
import QuickInstruction from './quick-instruction.vue'
import FloatIconSetting from './float-icon-setting.vue'
import { CopyOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'default'
  }
})
const isWorkspace = computed(() => props.variant === 'workspace')
const emit = defineEmits(['saved'])

const { t } = useI18n('views.robot.robot-config.external-service.components.embed-website')
const robotStore = useRobotStore()
const { robotInfo, external_config_pc } = storeToRefs(robotStore)
const { getRobot } = robotStore

const sdkCode = computed(() => getSdkCode(robotInfo.value))
const copied = ref(false)
let copyTimer = null
const codeModalOpen = ref(false)
const openCodeModal = () => {
  codeModalOpen.value = true
}

const copySDKCode = async () => {
  try {
    if (navigator.clipboard?.writeText) await navigator.clipboard.writeText(sdkCode.value)
    else copyText(sdkCode.value)
  } catch {
    copyText(sdkCode.value)
  }
  copied.value = true
  message.success(t('copy_success'))
  if (copyTimer) clearTimeout(copyTimer)
  copyTimer = setTimeout(() => {
    copied.value = false
  }, 1500)
}

const formRef = ref()
const formState = reactive({
  headTitle: external_config_pc.value.headTitle,
  headSubTitle: external_config_pc.value.headSubTitle,
  headImage: external_config_pc.value.headImage,
  lang: external_config_pc.value.lang,
  pageStyle: external_config_pc.value.pageStyle,
  floatBtn: external_config_pc.value.floatBtn,
  open_type: external_config_pc.value.open_type,
  window_width: external_config_pc.value.window_width,
  window_height: external_config_pc.value.window_height,
  iframe_width: external_config_pc.value.iframe_width,
  iframe_height: external_config_pc.value.iframe_height,
  iframe_right: external_config_pc.value.iframe_right,
  iframe_bottom: external_config_pc.value.iframe_bottom,
  iframe_resize_enabled: external_config_pc.value.iframe_resize_enabled,
  iframe_drag_enabled: external_config_pc.value.iframe_drag_enabled,
  new_session_btn_show: external_config_pc.value.new_session_btn_show,
  avatarShow: external_config_pc.value.avatarShow
})

const previewIframeSrc = computed(() => {
  let { pc_domain, robot_key } = robotInfo.value
  return toPageProtocol(`${pc_domain}/web/#/chat?robot_key=${robot_key}`)
})

const toIntegerInRange = (value, fallback, min, max) => {
  const number = Number(value)
  if (!Number.isFinite(number)) {
    return fallback
  }
  return Math.min(max, Math.max(min, Math.round(number)))
}

const previewSize = computed(() => {
  const width = toIntegerInRange(formState.iframe_width, 418, 320, 2000)
  const height = toIntegerInRange(formState.iframe_height, 680, 400, 2000)
  const maxWidth = isWorkspace.value ? 344 : 418
  const maxHeight = isWorkspace.value ? 584 : 720
  const scale = Math.min(1, maxWidth / width, maxHeight / height)

  return { width, height, scale }
})

const previewBoxStyle = computed(() => ({
  width: `${previewSize.value.width * previewSize.value.scale}px`,
  height: `${previewSize.value.height * previewSize.value.scale}px`
}))

const previewIframeStyle = computed(() => ({
  width: `${previewSize.value.width}px`,
  height: `${previewSize.value.height}px`,
  transform: `scale(${previewSize.value.scale})`
}))

let autoSaveTimer = null
let autoSaveEnabled = false
let autoSaveRunning = false
let autoSaveQueued = false

const autoSaveForm = () => {
  if (!isWorkspace.value || !autoSaveEnabled) return
  if (autoSaveRunning) {
    autoSaveQueued = true
    return
  }
  formRef.value
    ?.validate()
    .then(async () => {
      autoSaveRunning = true
      try {
        await saveWebSiteInfo({ silent: true })
      } finally {
        autoSaveRunning = false
        if (autoSaveQueued) {
          autoSaveQueued = false
          scheduleAutoSave(0)
        }
      }
    })
    .catch(() => {})
}

const scheduleAutoSave = (delay = 500) => {
  if (!isWorkspace.value || !autoSaveEnabled) return
  if (autoSaveRunning) {
    autoSaveQueued = true
    return
  }
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    autoSaveTimer = null
    autoSaveForm()
  }, delay)
}

const flushAutoSave = () => {
  if (!autoSaveTimer) return
  clearTimeout(autoSaveTimer)
  autoSaveTimer = null
  autoSaveForm()
}

watch(formState, (val) => {
  updatePreview(val)
  scheduleAutoSave()
})

onMounted(() => {
  nextTick(() => {
    autoSaveEnabled = true
  })
})

onBeforeUnmount(() => {
  flushAutoSave()
  clearTimeout(leftScrollTimer)
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = null
  if (copyTimer) clearTimeout(copyTimer)
  copyTimer = null
})

const leftScrolling = ref(false)
let leftScrollTimer
const onLeftScroll = () => {
  if (!isWorkspace.value) return
  leftScrolling.value = true
  clearTimeout(leftScrollTimer)
  leftScrollTimer = setTimeout(() => { leftScrolling.value = false }, 800)
}

const updataQuickComand = (data) => {
  updatePreview(data, 'updataQuickComand')
}

const handlePreviewLoad = () => {
  updatePreview(formState)
}

const updatePreview = (data, type) => {
  const iframe = document.getElementById('web-preview')
  if (!iframe?.contentWindow) return
  iframe.contentWindow.postMessage(
    {
      type: type || 'onPreview',
      data: JSON.parse(JSON.stringify(data))
    },
    '*'
  )
}

const formRules = {
  lang: [
    {
      required: true,
      message: t('please_select_language'),
      trigger: 'change'
    }
  ],
  headTitle: [
    {
      required: true,
      message: t('please_input_title'),
      trigger: 'input'
    },
    {
      trigger: 'input',
      validator: () => {
        if (!formState.headImage) {
          return Promise.reject(t('please_upload_logo'))
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  pageStyle: {
    headBackgroundColor: [
      {
        required: true,
        message: t('please_select_color'),
        trigger: 'change'
      }
    ]
  }
}

// 保存样式设置
const saveWebSiteInfo = async ({ silent = false } = {}) => {
  const { id } = robotInfo.value
  formState.window_width = +formState.window_width || 1200
  formState.window_height = +formState.window_height || 650
  formState.iframe_width = toIntegerInRange(formState.iframe_width, 418, 320, 2000)
  formState.iframe_height = toIntegerInRange(formState.iframe_height, 680, 400, 2000)
  formState.iframe_right = toIntegerInRange(formState.iframe_right, 50, 0, 2000)
  formState.iframe_bottom = toIntegerInRange(formState.iframe_bottom, 50, 0, 2000)
  const formData = JSON.parse(JSON.stringify(toRaw(formState)))

  try {
    await editExternalConfig({
      id: id,
      external_config_pc: JSON.stringify(formData)
    })
    if (isWorkspace.value) {
      Object.assign(robotStore.external_config_pc, formData)
    } else {
      getRobot(id)
    }
    if (!silent) message.success(t('save_success'))
    emit('saved')
    return true
  } catch {
    message.error(t('save_failed'))
    return false
  }
}

const saveForm = () => {
  formRef.value
    .validate()
    .then(() => {
      saveWebSiteInfo()
    })
    .catch((error) => {
      console.log('error', error)
    })
}

// 保存浮标设置
const handleFloatBtnCongiSave = (data) => {
  formState.floatBtn = { ...data }
  if (!isWorkspace.value) saveForm()
}
</script>
