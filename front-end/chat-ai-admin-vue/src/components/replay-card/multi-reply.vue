<template>
  <div class="multi-reply-card">
    <a-tabs v-model:activeKey="activeTab" @change="onTabChange">
      <a-tab-pane key="imageText">
        <template #tab>
          <span class="tab-with-icon"><svg-icon name="reply-content-image-text" size="14" /> {{ t('tab_image_text') }}</span>
        </template>
        <div class="form-grid">
          <a-form :model="imageText" :rules="rules.imageText" ref="imageTextForm" layout="vertical"
            :labelCol="labelCol" :wrapperCol="wrapperCol">
            <a-form-item :label="t('label_link_url')" name="url">
              <a-input v-model:value="imageText.url" :placeholder="t('ph_link_url')" />
            </a-form-item>
            <a-form-item :label="t('label_link_title')" name="title">
              <a-input v-model:value="imageText.title" :placeholder="t('ph_link_title')" />
            </a-form-item>
            <a-form-item :label="t('label_link_desc')" name="description">
              <a-textarea v-model:value="imageText.description" auto-size :placeholder="t('ph_link_desc')" :maxlength="300"
                :showCount="true" />
            </a-form-item>
            <a-form-item :label="t('label_link_image')" name="thumb_url">
              <a-input v-model:value="imageText.thumb_url" style="display:none" />
              <a-form-item-rest>
                <div class="upload-row">
                  <a-upload :show-upload-list="false" :before-upload="handleBeforeUpload"
                    :custom-request="(options) => handleUpload(options, 'received_message_images')"
                    accept=".png,.jpg,.jpeg">
                    <a-button type="default">
                      <UploadOutlined />{{ t('btn_upload_image') }}
                    </a-button>
                  </a-upload>
                  <a-input @paste="onPaste" @click="pasteImage" :placeholder="t('ph_paste_upload')"></a-input>
                </div>
                <div v-if="imageText.thumb_url" class="preview-wrap">
                  <img v-viewer :src="imageText.thumb_url" class="preview" />
                  <span class="del-icon" @click="removeImage">
                    <CloseOutlined />
                  </span>
                </div>
              </a-form-item-rest>
            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <a-tab-pane key="text">
        <template #tab>
          <span class="tab-with-icon"><svg-icon name="reply-content-text" size="14" /> {{ t('tab_text') }}</span>
        </template>
        <div class="text-area">
          <a-textarea ref="textAreaRef" v-model:value="text.description" :rows="6" :placeholder="t('ph_input')" :maxlength="300"
            :showCount="true" />
          <div class="emoji-row">
            <a-popover v-model:open="showEmoji" placement="bottomLeft" trigger="click" :getPopupContainer="getPopup">
              <template #content>
                <Picker :data="emojiIndex" :emojiSize="18" :showPreview="false" set="apple" @select="onEmojiSelect" />
              </template>
              <a-tooltip :title="t('tip_insert_emoji')">😊</a-tooltip>
            </a-popover>
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane key="image">
        <template #tab>
          <span class="tab-with-icon"><svg-icon name="reply-content-image" size="14" /> {{ t('tab_image') }}</span>
        </template>
        <div class="form-grid">
          <div class="upload-drop" @paste="onPaste" @drop.prevent="onDrop" @dragover.prevent>
            <div class="upload-actions">
              <div v-if="image.thumb_url" class="preview-wrap">
                <img v-viewer :src="image.thumb_url" class="preview" />
                <span class="del-icon" @click="removeImage">
                  <CloseOutlined />
                </span>
              </div>
              <a-upload :show-upload-list="false" :before-upload="handleBeforeUpload"
                :custom-request="(options) => handleUpload(options, 'received_message_images')"
                accept=".png,.jpg,.jpeg">
                <div class="hint">{{ t('tip_image_upload') }}</div>
              </a-upload>
            </div>
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane key="url">
        <template #tab>
          <span class="tab-with-icon"><svg-icon name="reply-content-link" size="14" /> {{ t('tab_url') }}</span>
        </template>
        <div class="form-grid">
          <a-form :model="urlCard" :rules="rules.url" ref="urlForm" layout="vertical" :labelCol="labelCol"
            :wrapperCol="wrapperCol">
            <a-form-item :label="t('label_link_title')" name="title">
              <a-input v-model:value="urlCard.title" :placeholder="t('ph_input')" />
            </a-form-item>
            <a-form-item :label="t('label_link_url')" name="url">
              <a-input v-model:value="urlCard.url" :placeholder="t('ph_link_url_2')" />
            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <a-tab-pane key="card">
        <template #tab>
          <span class="tab-with-icon"><svg-icon name="applet" size="14" /> {{ t('tab_mini_card') }}</span>
        </template>
        <div class="form-grid">
          <a-form :model="miniCard" :rules="rules.card" ref="cardForm" layout="vertical" :labelCol="cardLabelCol"
            :wrapperCol="cardWrapperCol">
            <a-alert type="info" style="margin-bottom: 16px;"
              :message="t('tip_mini_card')" />
            <a-form-item :label="t('label_mini_title')" name="title">
              <a-input v-model:value="miniCard.title" :placeholder="t('ph_mini_title')" />
            </a-form-item>
            <a-form-item :label="t('label_mini_appid')" name="appid">
              <a-input v-model:value="miniCard.appid" :placeholder="t('ph_mini_appid')" />
              <div style="margin-top: 2px; color: #FB363F;">{{ t('tip_mini_appid') }}</div>
            </a-form-item>
            <a-form-item :label="t('label_mini_path')" name="page_path">
              <a-input v-model:value="miniCard.page_path" :placeholder="t('ph_mini_path')" />
              <div style="margin-top: 2px; color: #8C8C8C;">{{ t('tip_mini_path') }}</div>
            </a-form-item>
            <a-form-item :label="t('label_mini_cover')" name="thumb_url">
              <a-input v-model:value="miniCard.thumb_url" style="display:none" />
              <a-form-item-rest>
                <div class="upload-row">
                  <a-upload :show-upload-list="false" :before-upload="handleBeforeUpload"
                    :custom-request="(options) => handleUpload(options, 'received_message_images')"
                    accept=".png,.jpg,.jpeg">
                    <a-button type="default">{{ t('btn_upload_image_2') }}</a-button>
                  </a-upload>
                  <a-input @paste="onPaste" @click="pasteImage" :placeholder="t('ph_paste_upload_2')"></a-input>
                </div>
                <div v-if="miniCard.thumb_url" class="preview-wrap">
                  <img v-viewer :src="miniCard.thumb_url" class="preview" />
                  <span class="del-icon" @click="removeImage">
                    <CloseOutlined />
                  </span>
                </div>
              </a-form-item-rest>
            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <a-tab-pane key="smartMenu" v-if="show_smart">
        <template #tab>
          <span class="tab-with-icon"><svg-icon name="reply-content-smart" size="14" /> {{ t('tab_smart_menu') }}</span>
        </template>
        <div class="smart-tab-box">
          <a-button type="dashed" class="choose-smart-btn" @click="openSmartChoose">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_choose_smart') }}
          </a-button>
          <div v-if="smartPreview.id" class="smart-preview">
            <div class="smart-list-card">
              <div class="smart-card-header">
                <a-avatar :src="smartPreview.avatar_url || defaultAvatar" :size="40" shape="square"
                  class="smart-avatar" />
                <div class="smart-card-content">
                  <div class="smart-card-title">{{ smartPreview.menu_description }}</div>
                  <div class="smart-card-text">
                    <template v-for="(line, li) in smartPreview.lines" :key="li">
                      <div class="smart-reply-line" @click="onSmartReplyLineClick($event)">
                        <span class="smart-line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
                        <div v-else-if="line.kind === 'newline'" class="empty-line"></div>
                        <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                        <a v-else-if="line.kind === 'keyword'" href="javascript:;" class="link">
                          <div v-if="line.serial_no">{{ line.serial_no }}</div>
                          {{ line.text }}
                        </a>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
              <div class="smart-card-footer">
                <div class="smart-footer-left">
                  <div class="smart-menu-name">{{ smartPreview.title || t('smart_menu') }}</div>
                  <div class="smart-update-time">{{ t('updated_at') }}{{ formatDateFn(smartPreview.update_time, 'YYYY-MM-DD') }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </a-tab-pane>

      <template #rightExtra>
        <a-button type="text" @click="emitDel">
          <CloseCircleOutlined style="font-size: 16px;" />
        </a-button>
      </template>
    </a-tabs>

    <a-modal v-if="show_smart" v-model:open="smartChooseOpen" :title="t('btn_choose_smart')" :width="860" @ok="onSmartChooseOk"
      @cancel="onSmartChooseCancel">
      <div class="smart-modal">
        <div class="smart-modal-header">
          <a-button type="dashed" @click="openAddSmartMenu">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_add_smart') }}
          </a-button>
          <template v-if="isSubscribeMode">
            <span class="ml8">{{ t('label_select_robot') }}</span>
            <a-select v-model:value="selectedRobotId" style="width: 220px" @change="(val) => loadSmartList(val)">
              <a-select-option v-for="rb in robotOptions" :value="rb.id" :key="rb.id">{{ rb.name }}</a-select-option>
            </a-select>
          </template>
        </div>
        <div v-if="smartList.length === 0" class="smart-empty">
          <a-empty :description="t('empty_smart_menu')" />
        </div>
        <div class="smart-grid" v-else>
          <div class="smart-card" v-for="it in smartList" :key="it.id" :class="{ selected: selectedSmartId === it.id }"
            @click="onTogglePick(it.id)">
            <div class="smart-list-card">
              <div class="smart-card-header">
                <a-avatar :src="it.avatar_url || defaultAvatar" :size="40" shape="square" class="smart-avatar" />
                <div class="smart-card-content">
                  <div class="smart-card-title">{{ it.menu_description }}</div>
                  <div class="smart-card-text">
                    <template v-for="(line, li) in getDisplayLines(it)" :key="li">
                      <div class="smart-reply-line" @click="onSmartReplyLineClick($event)">
                        <span class="smart-line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
                        <div v-else-if="line.kind === 'newline'" class="empty-line"></div>
                        <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                        <a v-else-if="line.kind === 'keyword'" href="javascript:;" class="link">
                          <div v-if="line.serial_no">{{ line.serial_no }}</div>
                          {{ line.text }}
                        </a>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
              <div class="smart-card-footer">
                <div class="smart-footer-left">
                  <div class="smart-menu-name">{{ it.menu_title || it.menu_name || t('smart_menu_2') }}</div>
                  <div class="smart-update-time">{{ t('updated_at_2') }}{{ formatDateFn(it.update_time, 'YYYY-MM-DD') }}</div>
                </div>
                <div class="smart-footer-right" @click.stop>
                  <template v-if="selectedSmartId === it.id">
                    <span class="smart-radio-inline">
                      <div class="picked" @click="onTogglePick(it.id)"></div>
                    </span>
                  </template>
                  <template v-else>
                    <a-dropdown placement="bottomRight">
                      <a class="more-btn" @click.prevent>
                        <EllipsisOutlined />
                      </a>
                      <template #overlay>
                        <a-menu>
                          <a-menu-item key="edit" @click="editSmart(it)">{{ t('btn_edit') }}</a-menu-item>
                          <a-menu-item key="del" class="del" @click="delSmart(it)"
                            style="color: #FF4D4F;">{{ t('btn_delete') }}</a-menu-item>
                        </a-menu>
                      </template>
                    </a-dropdown>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { PlusOutlined, CloseCircleOutlined, UploadOutlined, CloseOutlined, EllipsisOutlined } from '@ant-design/icons-vue'
import emojiDataJson from 'emoji-mart-vue-fast/data/all.json'
import "emoji-mart-vue-fast/css/emoji-mart.css"
import { ref, reactive, watch, nextTick, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { uploadFile } from '@/api/app'
import { Picker, EmojiIndex } from 'emoji-mart-vue-fast/src'
import { useRoute } from 'vue-router'
import { getSmartMenuList, deleteSmartMenu, getSpecifyAbilityConfig, getRobotSpecifyAbilityConfig } from '@/api/explore/index.js'
import { getRobotList } from '@/api/robot/index.js'
import dayjs from 'dayjs'
import { DEFAULT_ROBOT_AVATAR } from '@/constants/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('components.replay-card.multi-reply')
const route = useRoute()

const user_smart_configs = ref(['exploreSubscribeReply', 'exploreCustomMenu'])
// 智能菜单回复是否开启
const user_show_smart = ref(false)
const robot_show_smart = ref(false)

const show_smart = computed(() => {
  if (user_smart_configs.value.includes(route.name)) {
    return user_show_smart.value
  } else {
    return robot_show_smart.value
  }
})

const props = defineProps({
  value: { type: Object, default: () => ({ type: 'text', description: '' }) },
  reply_index: { type: Number, default: 0 }
})
const emit = defineEmits(['update:value', 'change', 'del'])

const activeTab = ref(props.value.type || 'text')

// states per tab
const imageText = reactive({ url: '', title: '', description: '', thumb_url: '' })
const text = reactive({ description: '' })
const image = reactive({ thumb_url: '' })
const urlCard = reactive({ title: '', url: '' })
const miniCard = reactive({ title: '', appid: '', page_path: '', thumb_url: '' })
const smartPreview = reactive({ id: '', title: '', update_time: 0, avatar_url: '', lines: [] })
const smartChooseOpen = ref(false)
const smartList = ref([])
const selectedSmartId = ref('')
const robotOptions = ref([])
const selectedRobotId = ref('')
const defaultAvatar = DEFAULT_ROBOT_AVATAR
const query = useRoute().query
const isSubscribeMode = computed(() => !robot_show_smart.value && user_smart_configs.value.includes(route.name))
const pendingSmartValue = ref(null)
const smartInitLoading = ref(true)

function initRobotShowSmart (robot_id) {
  const rid = String(robot_id || '')
  if (!rid) { robot_show_smart.value = false; return Promise.resolve() }
  return getRobotSpecifyAbilityConfig({ ability_type: 'robot_smart_menu', robot_id: rid }).then((res) => {
    const st = res?.data?.robot_config?.switch_status
    robot_show_smart.value = String(st || '0') == '1'
  })
}

function initUserShowSmart () {
  getSpecifyAbilityConfig({ ability_type: 'robot_smart_menu' }).then((res) => {
    const st = res?.data?.user_config?.switch_status
    user_show_smart.value = String(st || '0') == '1'
  })
}

let isUploading = false
const emojiIndex = new EmojiIndex(emojiDataJson)
const labelCol = { span: 5 }
const wrapperCol = { span: 19 }
const cardLabelCol = { span: 4 }
const cardWrapperCol = { span: 20 }
const imageTextForm = ref(null)
const urlForm = ref(null)
const cardForm = ref(null)

watch(() => text.description, (v) => {
  const s = String(v || '')
  if (s.length > 300) {
    message.warning(t('msg_max_chars'))
    text.description = s.slice(0, 300)
  }
})
watch(() => imageText.description, (v) => {
  const s = String(v || '')
  if (s.length > 300) {
    message.warning(t('msg_max_chars'))
    imageText.description = s.slice(0, 300)
  }
})

const rules = {
  imageText: {
    url: [{ required: true, message: t('msg_input_http_https') }, { validator: httpValidator }],
    title: [{ required: true, message: t('msg_input_title') }],
    description: [{ required: true, message: t('msg_input_desc') }],
    thumb_url: [{ required: true, message: t('msg_upload_link_image') }]
  },
  url: {
    title: [{ required: true, message: t('msg_input_title_2') }],
    url: [{ required: true, message: t('msg_input_http_https_2') }, { validator: httpValidator }]
  },
  card: {
    title: [{ required: true, message: t('msg_input_mini_title') }],
    appid: [{ required: true, message: t('msg_input_mini_appid') }],
    page_path: [{ required: true, message: t('msg_input_mini_path') }],
    thumb_url: [{ required: true, message: t('msg_upload_mini_cover') }]
  },
  smartMenu: {
    id: [{ required: true, message: t('msg_select_smart') }]
  }
}

function httpValidator (_rule, value) {
  if (!value) return Promise.resolve()
  const ok = /^https?:\/\//.test(value)
  return ok ? Promise.resolve() : Promise.reject(t('msg_link_protocol'))
}

const showEmoji = ref(false)
const textAreaRef = ref(null)
function onEmojiSelect (emoji) {
  const char = emoji?.native || ''
  if (!char) return
  // 尝试在光标处插入表情，若不可用则追加到末尾
  let el = textAreaRef.value?.$el?.querySelector('textarea')
  if (el && typeof el.selectionStart === 'number') {
    const start = el.selectionStart
    const end = el.selectionEnd
    const val = text.description || ''
    let nextVal = val.slice(0, start) + char + val.slice(end)
    if ((nextVal || '').length > 300) {
      message.warning(t('msg_max_chars_2'))
      nextVal = (nextVal || '').slice(0, 300)
    }
    text.description = nextVal
    // 更新光标位置
    nextTick(() => {
      el.focus()
      const pos = Math.min(start + char.length, (text.description || '').length)
      el.setSelectionRange(pos, pos)
    })
  } else {
    let nextVal = (text.description || '') + char
    if ((nextVal || '').length > 300) {
      message.warning(t('msg_max_chars_2'))
      nextVal = (nextVal || '').slice(0, 300)
    }
    text.description = nextVal
  }
  showEmoji.value = false
  onChange()
}

function onTabChange () { onChange() }

function emitDel () { emit('del', props.reply_index) }

function handleBeforeUpload (file) {
  const isValidType = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isValidType) {
    message.error(t('msg_image_format'))
    return false
  }
  const isLt2M = file.size / 1024 < 1024 * 2
  if (!isLt2M) {
    message.error(t('msg_image_size'))
    return false
  }
  return true
}

async function handleUpload ({ file, onError, onSuccess }, category) {
  try {
    const res = await uploadFile({
      category,
      file
    })
    const url = res?.data?.link || res?.data?.url || ''
    if (activeTab.value === 'imageText') {
      imageText.thumb_url = url
      nextTick(() => { imageTextForm.value?.validateFields?.(['thumb_url']) })
    } else if (activeTab.value === 'image') {
      image.thumb_url = url
    } else if (activeTab.value === 'card') {
      miniCard.thumb_url = url
      nextTick(() => { cardForm.value?.validateFields?.(['thumb_url']) })
    }
    onSuccess && onSuccess(res)
    onChange()
  } catch (e) {
    message.error(t('msg_upload_failed'))
    onError && onError(e)
  }
}

function onChange () {
  const payload = getPayload()
  emit('update:value', payload)
  emit('change', { ...payload, reply_index: props.reply_index })
}

function removeImage () {
  if (activeTab.value === 'imageText') {
    imageText.thumb_url = ''
    nextTick(() => { imageTextForm.value?.validateFields?.(['thumb_url']) })
  }
  else if (activeTab.value === 'image') {
    image.thumb_url = ''
  }
  else if (activeTab.value === 'card') {
    miniCard.thumb_url = ''
    nextTick(() => { cardForm.value?.validateFields?.(['thumb_url']) })
  }
  onChange()
}

function getPayload () {
  const t = activeTab.value
  if (t === 'imageText') {
    return {
      reply_type: t,
      thumb_url: imageText.thumb_url || '',
      title: imageText.title || '',
      description: imageText.description || '',
      url: imageText.url || '',
      page_path: '',
      appid: '',
      status: '1',
      smart_menu_id: '',
      more_img_text_json: '',
      media_id: '',
      pic: imageText.thumb_url || '',
      type: t
    }
  }
  if (t === 'text') {
    return {
      reply_type: t,
      thumb_url: '',
      title: '',
      description: text.description || '',
      url: '',
      page_path: '',
      appid: '',
      status: '1',
      smart_menu_id: '',
      more_img_text_json: '',
      media_id: '',
      pic: '',
      type: t
    }
  }
  if (t === 'image') {
    return {
      reply_type: t,
      thumb_url: image.thumb_url || '',
      title: '',
      description: '',
      url: '',
      page_path: '',
      appid: '',
      status: '1',
      smart_menu_id: '',
      more_img_text_json: '',
      media_id: '',
      pic: image.thumb_url || '',
      type: t
    }
  }
  if (t === 'url') {
    return {
      reply_type: t,
      thumb_url: '',
      title: urlCard.title || '',
      description: '',
      url: urlCard.url || '',
      page_path: '',
      appid: '',
      status: '1',
      smart_menu_id: '',
      more_img_text_json: '',
      media_id: '',
      pic: '',
      type: t
    }
  }
  if (t === 'card') {
    return {
      reply_type: t,
      thumb_url: miniCard.thumb_url || '',
      title: miniCard.title || '',
      description: '',
      url: '',
      page_path: miniCard.page_path || '',
      appid: miniCard.appid || '',
      status: '1',
      smart_menu_id: '',
      more_img_text_json: '',
      media_id: '',
      pic: miniCard.thumb_url || '',
      type: t
    }
  }
  if (t === 'smartMenu') {
    return {
      reply_type: t,
      thumb_url: '',
      title: smartPreview.title || '',
      description: '',
      url: '',
      page_path: '',
      appid: '',
      status: '1',
      smart_menu_id: smartPreview.id || '',
      more_img_text_json: '',
      media_id: '',
      pic: '',
      type: t
    }
  }
  return {
    reply_type: 'text',
    thumb_url: '',
    title: '',
    description: '',
    url: '',
    page_path: '',
    appid: '',
    status: '1',
    smart_menu_id: '',
    more_img_text_json: '',
    media_id: '',
    pic: '',
    type: 'text'
  }
}

// 规避某些环境下 ant popover 的 getPopupContainer 使用 triggerNode.parentNode 导致的空引用
function getPopup () { return document.body }

function pasteImage () { message.info(t('msg_paste_image')); }
function onPaste (e) {
  if (isUploading) return
  const items = e.clipboardData?.items || []
  for (const it of items) {
    if (it.type.indexOf('image') !== -1) {
      const file = it.getAsFile()
      if (!file) continue
      if (!handleBeforeUpload(file)) return
      isUploading = true
      handleUpload({ file, onError: () => { isUploading = false }, onSuccess: () => { isUploading = false } }, 'received_message_images')
      break
    }
  }
}
function onDrop (e) {
  const file = e.dataTransfer?.files?.[0]
  if (file && handleBeforeUpload(file)) {
    handleUpload({ file, onError: () => { }, onSuccess: () => { } }, 'received_message_images')
  }
}

// initialize from value
watch(() => props.value, (val) => {
  const t = val?.type || 'text'
  if (t === 'smartMenu') {
    if (show_smart.value && !smartInitLoading.value) {
      activeTab.value = 'smartMenu'
    } else {
      activeTab.value = 'text'
    }
  } else {
    activeTab.value = t
  }
  if (t === 'imageText') Object.assign(imageText, val)
  else if (t === 'text') Object.assign(text, { description: val?.description || '' })
  else if (t === 'image') Object.assign(image, { thumb_url: val?.thumb_url || '' })
  else if (t === 'url') Object.assign(urlCard, { title: val?.title || '', url: val?.url || '' })
  else if (t === 'card') Object.assign(miniCard, val)
  else if (t === 'smartMenu') {
    if (show_smart.value) {
      smartPreview.id = String(val?.smart_menu_id || '')
      if (smartPreview.id && val?.smart_menu?.id) {
        const data = val.smart_menu || {}
        smartPreview.title = data.menu_title || data.menu_name || ''
        smartPreview.update_time = Number(data.update_time || data.create_time || 0)
        smartPreview.avatar_url = data.avatar_url || ''
        smartPreview.menu_description = data.menu_description || ''
        smartPreview.lines = buildMenuLines(Array.isArray(data.menu_content) ? data.menu_content : [])
      }
      selectedSmartId.value = String(smartPreview.id || '')
    } else {
      pendingSmartValue.value = val
    }
  }
}, { immediate: true })

// emit changes on input editing
watch(imageText, () => { onChange() }, { deep: true })
watch(text, () => { onChange() }, { deep: true })
watch(image, () => { onChange() }, { deep: true })
watch(urlCard, () => { onChange() }, { deep: true })
watch(miniCard, () => { onChange() }, { deep: true })

function formatDateFn (date, format = 'YYYY-MM-DD') {
  if (!date) return '';
  return dayjs(Number(date) * 1000).format(format)
}

function buildMenuLines (menu_content) {
  const out = [];
  (Array.isArray(menu_content) ? menu_content : []).forEach((mc) => {
    const t = String(mc?.menu_type || '')
    const txt = String(mc?.content || '')
    if (t === '0') {
      if (txt === '') {
        out.push({ kind: 'newline' })
      } else if (/<a[\s\S]*?<\/a>/.test(txt)) {
        const sanitized = /href\s*=\s*['"]\s*#\s*['"]/i.test(txt) ? txt.replace(/href\s*=\s*['"]\s*#\s*['"]/ig, 'href="javascript:;"') : txt.replace(/href=/ig, 'target="_blank" href=');
        out.push({
          kind: 'html',
          html: sanitized
        })
      } else {
        out.push({
          kind: 'text',
          text: txt
        })
      }
    } else if (t === '1') {
      out.push({
        kind: 'keyword',
        text: txt,
        serial_no: mc?.serial_no || ''
      })
    }
  })
  return out.slice(0, 8)
}

function buildPreviewLinesFromReplyContent (list) {
  const lines = [];
  (Array.isArray(list) ? list : []).forEach((rc) => {
    const t = rc?.type || rc?.reply_type
    if (t === 'text') {
      if (rc?.description) {
        lines.push({
          kind: 'text',
          text: rc.description
        })
      }
    } else if (t === 'imageText') {
      const txt = rc?.title || rc?.description || t('label_image_text');
      lines.push({ kind: 'keyword', text: txt })
    } else if (t === 'url') {
      const txt = rc?.title || t('label_link');
      lines.push({ kind: 'keyword', text: txt })
    } else if (t === 'image') {
      lines.push({ kind: 'text', text: t('label_image') })
    } else if (t === 'card') {
      const txt = rc?.title || t('label_mini_card');
      lines.push({ kind: 'text', text: txt })
    } else if (t === 'smartMenu') {
      const txt = rc?.title || t('label_smart_menu');
      lines.push({ kind: 'text', text: txt })
    }
  })
  return lines.slice(0, 8)
}

function getDisplayLines (item) {
  if (Array.isArray(item?.menu_content) && item.menu_content.length) return buildMenuLines(item.menu_content)
  return buildPreviewLinesFromReplyContent(item?.reply_content || [])
}

function openSmartChoose () {
  smartChooseOpen.value = true
  selectedSmartId.value = String(smartPreview.id || '')
  if (isSubscribeMode.value) {
    // 关注后回复：先选择机器人，再拉取其智能菜单
    loadRobotOptions().then(() => {
      const def = selectedRobotId.value || robotOptions.value[0]?.id || ''
      selectedRobotId.value = String(def)
      if (def) loadSmartList(def)
    })
  } else {
    loadSmartList()
  }
}
function openAddSmartMenu () { const url = `/#/robot/ability/smart-menu?id=${query.id}&robot_key=${query.robot_key}`; window.open(url, '_blank') }
function loadSmartList (robot_id) {
  const rid = String(robot_id || query.id || '')
  getSmartMenuList({ robot_id: rid }).then((res) => {
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
  smartList.value = list.map((it) => ({
    id: String(it.id || ''),
    menu_title: it.menu_title || it.menu_name || t('smart_menu_3'),
    update_time: Number(it.update_time || it.create_time || 0),
    avatar_url: it.avatar_url || '',
    menu_description: it.menu_description || '',
    menu_content: Array.isArray(it.menu_content) ? it.menu_content : [],
    reply_content: Array.isArray(it.reply_content) ? it.reply_content : []
  }))
  // 默认不选中任何智能菜单
}).catch(() => { smartList.value = [] })
}

async function loadRobotOptions () {
  try {
    const res = await getRobotList()
    const lists = Array.isArray(res?.data) ? res.data : []
    robotOptions.value = lists.map((r) => ({ id: String(r.id || ''), name: r.robot_name || r.name || r.app_name || r.id }))
  } catch (_) {
    robotOptions.value = []
  }
}

function onSmartChooseOk () {
  if (!selectedSmartId.value) { message.warning(t('msg_select_one_smart')); return }
  const picked = smartList.value.find((it) => String(it.id) === String(selectedSmartId.value))
  if (!picked) { message.error(t('msg_smart_not_found')); return }
  smartPreview.id = String(picked.id || selectedSmartId.value)
  smartPreview.title = picked.menu_title || picked.menu_name || ''
  smartPreview.update_time = Number(picked.update_time || picked.create_time || 0)
  smartPreview.avatar_url = picked.avatar_url || ''
  smartPreview.menu_description = picked.menu_description || ''
  smartPreview.lines = getDisplayLines(picked)
  activeTab.value = 'smartMenu'
  onChange()
  smartChooseOpen.value = false
}
function onSmartChooseCancel () { smartChooseOpen.value = false }

function onTogglePick (id) {
  const sid = String(id || '')
  selectedSmartId.value = (String(selectedSmartId.value || '') === sid) ? '' : sid
}

function editSmart (it) {
  const url = `/#/robot/ability/smart-menu/add-rule?id=${query.id}&robot_key=${query.robot_key}&rule_id=${it.id}`
  window.open(url, '_blank')
}
function delSmart (it) {
  deleteSmartMenu({ id: it.id, robot_id: query.id }).then((res) => {
    if (res && res.res == 0) {
      message.success(t('msg_delete_success'))
      loadSmartList()
    }
  })
}

function onSmartReplyLineClick (e) {
  const a = e.target?.closest?.('a')
  if (!a) return
  const href = String(a.getAttribute('href') || '')
  if (href === '#' || href === 'javascript:;') {
    e.preventDefault();
    e.stopPropagation()
  }
}

async function validate () {
  const t = activeTab.value
  try {
    if (t === 'imageText') {
      try { await imageTextForm.value?.validate(); return true } catch (_e) { message.warning(t('msg_complete_image_text')); return false }
    }
    if (t === 'text') {
      const s = String(text.description || '').trim()
      if (!s) { message.warning(t('msg_input_text')); return false }
      return true
    }
    if (t === 'image') {
      const url = String(image.thumb_url || '').trim()
      if (!url) { message.warning(t('msg_upload_image_2')); return false }
      return true
    }
    if (t === 'url') {
      try { await urlForm.value?.validate(); return true } catch (_e) { message.warning(t('msg_complete_url')); return false }
    }
    if (t === 'card') {
      try { await cardForm.value?.validate(); return true } catch (_e) { message.warning(t('msg_complete_mini_card')); return false }
    }
    if (t === 'smartMenu') {
      const sid = String(selectedSmartId.value || smartPreview.id || '')
      if (!sid) { message.warning(t('msg_select_smart_2')); return false }
      return true
    }
  } catch (_e) {
    return false
  }
  return true
}

onMounted(async () => {
  try {
    if (user_smart_configs.value.includes(route.name)) {
      await initUserShowSmart()
    } else {
      await initRobotShowSmart(route.query.id)
    }
  } finally {
    smartInitLoading.value = false
  }
})

watch(() => show_smart.value, (v) => {
  if (v && pendingSmartValue.value) {
    const val = pendingSmartValue.value
    activeTab.value = 'smartMenu'
    smartPreview.id = String(val?.smart_menu_id || '')
    if (smartPreview.id && val?.smart_menu?.id) {
      const data = val.smart_menu || {}
      smartPreview.title = data.menu_title || data.menu_name || ''
      smartPreview.update_time = Number(data.update_time || data.create_time || 0)
      smartPreview.avatar_url = data.avatar_url || ''
      smartPreview.menu_description = data.menu_description || ''
      smartPreview.lines = buildMenuLines(Array.isArray(data.menu_content) ? data.menu_content : [])
    }
    selectedSmartId.value = String(smartPreview.id || '')
    pendingSmartValue.value = null
  }
})

defineExpose({ validate })
</script>

<style scoped lang="less">
.multi-reply-card {
  width: 694px;
  background: #F2F4F7;
  border-radius: 6px;
  padding: 0px 16px 16px;
  margin-bottom: 8px;
}

.tab-with-icon {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.smart-tab-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.choose-smart-btn {
  align-self: flex-start;
}

.smart-preview {
  margin-top: 4px;
  width: 355px;
  box-sizing: border-box;
}

.smart-modal-header {
  margin-bottom: 8px;
}

.smart-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.smart-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(300px, 1fr));
  gap: 16px;
}

.smart-card {
  position: relative;
  cursor: pointer;
}

.smart-card.selected .smart-list-card {
  border-color: #2475FC;
  box-shadow: 0 0 0 2px rgba(36, 117, 252, 0.2);
}

.smart-footer-right {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.smart-action {
  color: #595959;
  cursor: pointer;
}

.smart-action.del {
  color: #FF4D4F;
}

.smart-radio-inline {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 1px solid #2475FC;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.smart-radio-inline .picked {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #2475FC;
}

.form-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.upload-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-row {
  margin-top: 8px;
}

.preview {
  cursor: pointer;
  margin-top: 8px;
  width: 120px;
  height: 120px;
  object-fit: contain;
  border-radius: 6px;
  border: 1px solid #d9d9d9;
}

.preview-box {
  display: flex;
  align-items: center;
}

.preview-wrap {
  width: 120px;
  position: relative;
  display: inline-block;
}

.del-icon {
  position: absolute;
  top: 4px;
  right: -10px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  transition: opacity .2s ease;
}

.preview-wrap:hover .del-icon {
  opacity: 1;
}

.text-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.emoji-row {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}

.emoji-grid {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  padding: 6px;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 6px;
}

.emoji-item {
  cursor: pointer;
}

.upload-drop {
  border-radius: 6px;

  .upload-actions {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .hint {
      cursor: pointer;
      display: flex;
      padding: 12px 0;
      justify-content: center;
      align-items: center;
      gap: 10px;
      background: #F2F4F7;
      color: #3a4559;
      text-align: center;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
  }
}

.ml8 {
  margin-left: 8px;
}

@media (max-width:768px) {
  .preview {
    width: 96px;
    height: 96px;
  }
}

.smart-list-card {
  border: 1px solid #e6e8ec;
  border-radius: 12px;
  background: #fff;
  padding: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04);
  transition: box-shadow .2s ease, transform .2s ease;
}

.smart-list-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.smart-card-header {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px 60px 12px 12px;
  background: #F2F4F7;
  border-radius: 12px 12px 0 0;
  height: 200px;
  box-shadow: 0 -4px 8px 0 #00000014 inset;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #d9d9d9 transparent;
}

.smart-card-header::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.smart-card-header::-webkit-scrollbar-track {
  background: transparent;
  border-radius: 8px;
}

.smart-card-header::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 8px;
}

.smart-card-header::-webkit-scrollbar-thumb:hover {
  background: #A1A7B3;
}

.smart-card-content {
  width: 100%;
  flex: 1;
  padding: 12px;
  background: #fff;
  border-radius: 2px 12px 12px 12px;
}

.smart-card-title {
  white-space: wrap;
  align-self: stretch;
  color: #1a1a1a;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  margin-bottom: 14px;
}

.smart-card-text {
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #3a4559;
  font-size: 14px;
  line-height: 22px;
  overflow: visible;

  .smart-reply-line {
    line-height: 22px;

    .smart-line-text {
      color: #3a4559;
    }

    .empty-line {
      height: 22px;
    }

    .link {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }
}

.smart-reply-line .link:hover {
  opacity: 0.7;
}

.smart-card-footer {
  box-sizing: border-box;
  margin-top: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-top: 1px solid #edeff2;
}

.smart-footer-left {
  display: flex;
  flex-direction: column;
  align-items: self-start;
}

.smart-menu-name {
  align-self: stretch;
  color: #262626;
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.smart-update-time {
  align-self: stretch;
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
}

.smart-avatar {
  border-radius: 4px !important;
}

/* 防止父级 vertical 表单影响内部横向表单布局 */
.multi-reply-card ::v-deep(.ant-form-item-row) {
  flex-direction: row;
}

.multi-reply-card ::v-deep(.ant-form-item) {
  margin-bottom: 16px;
}

.more-btn {
  color: #595959;
  width: 24px;
  height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
}

.more-btn:hover {
  color: #595959;
  background: #E4E6EB;
}
</style>
