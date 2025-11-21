<template>
  <div class="multi-reply-card">
    <a-tabs v-model:activeKey="activeTab" @change="onTabChange">
      <a-tab-pane key="imageText" :tab="tabLabel('图文链接')">
        <div class="form-grid">
          <a-form :model="imageText" :rules="rules.imageText" ref="imageTextForm" layout="horizontal" :labelCol="labelCol" :wrapperCol="wrapperCol">
            <a-form-item label="链接地址" name="url">
              <a-input v-model:value="imageText.url" placeholder="链接请已http或https开头（备注将展示在客户昵称后）" />
            </a-form-item>
            <a-form-item label="链接标题" name="title">
              <a-input v-model:value="imageText.title" placeholder="链接标题" />
            </a-form-item>
            <a-form-item label="链接描述" name="description">
              <a-textarea v-model:value="imageText.description" auto-size placeholder="链接描述" :maxlength="300" :showCount="true" />
            </a-form-item>
            <a-form-item label="链接图片" name="thumb_url">
        <div class="upload-row">
          <a-upload :show-upload-list="false" :before-upload="handleBeforeUpload"
            :custom-request="(options) => handleUpload(options, 'received_message_images')" accept=".png,.jpg,.jpeg">
            <a-button type="default"><UploadOutlined />上传图片</a-button>
          </a-upload>
          <a-input @paste="onPaste" @click="pasteImage" placeholder="复制粘贴上传图片">复制粘贴上传</a-input>
        </div>
        <div v-if="imageText.thumb_url" class="preview-wrap">
          <img v-viewer :src="imageText.thumb_url" class="preview" />
          <span class="del-icon" @click="removeImage"><CloseOutlined /></span>
        </div>
      </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <a-tab-pane key="text" :tab="tabLabel('文字')">
        <div class="text-area">
          <a-textarea ref="textAreaRef" v-model:value="text.description" :rows="6" placeholder="请输入" :maxlength="300" :showCount="true" />
          <div class="emoji-row">
            <a-popover v-model:open="showEmoji" placement="bottomLeft" trigger="click" :getPopupContainer="getPopup">
              <template #content>
                <Picker :data="emojiIndex" :emojiSize="18" :showPreview="false" set="apple" @select="onEmojiSelect" />
              </template>
              <a-tooltip title="插入表情">😊</a-tooltip>
            </a-popover>
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane key="image" :tab="tabLabel('图片')">
        <div class="form-grid">
          <div class="upload-drop" @paste="onPaste" @drop.prevent="onDrop" @dragover.prevent>
            <div class="upload-actions">
              <div v-if="image.thumb_url" class="preview-wrap">
                <img v-viewer :src="image.thumb_url" class="preview" />
                <span class="del-icon" @click="removeImage"><CloseOutlined /></span>
              </div>
              <a-upload :show-upload-list="false" :before-upload="handleBeforeUpload"
                :custom-request="(options) => handleUpload(options, 'received_message_images')" accept=".png,.jpg,.jpeg">
                <div class="hint">支持点击空白处、拖拽、粘贴图片，上传图片不得超过2M，仅支持png、jpg、jpeg格式</div>
              </a-upload>
            </div>
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane key="url" :tab="tabLabel('网址')">
        <div class="form-grid">
          <a-form :model="urlCard" :rules="rules.url" ref="urlForm" layout="horizontal" :labelCol="labelCol" :wrapperCol="wrapperCol">
            <a-form-item label="链接标题" name="title">
              <a-input v-model:value="urlCard.title" placeholder="请输入" />
            </a-form-item>
            <a-form-item label="打开网址" name="url">
              <a-input v-model:value="urlCard.url" placeholder="链接请已http或https开头（备注将展示在客户昵称后）" />
            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <a-tab-pane key="card" :tab="tabLabel('小程序卡片')">
        <div class="form-grid">
          <a-form :model="miniCard" :rules="rules.card" ref="cardForm" layout="horizontal" :labelCol="cardLabelCol" :wrapperCol="cardWrapperCol">
            <a-alert type="info" style="margin-bottom: 16px;"
              message="小程序卡片仅支持在公众号，微信客服和微信小程序中发送，其他渠道会发送失败公众号内回复必须是关联小程序，微信小程序内回复必须是当前的小程序" />
            <a-form-item label="小程序标题" name="title">
              <a-input v-model:value="miniCard.title" placeholder="请输入小程序卡片标题" />
            </a-form-item>
            <a-form-item label="小程序appID" name="appid">
              <a-input v-model:value="miniCard.appid" placeholder="请输入小程序appID" />
              <div style="margin-top: 2px; color: #FB363F;">小程序右上角三个点>名字>更多资料>Appid;公众号内回复时必须跟小程序是关联关系</div>
            </a-form-item>
            <a-form-item label="小程序路径" name="page_path">
              <a-input v-model:value="miniCard.page_path" placeholder="请输入小程序路径" />
              <div style="margin-top: 2px; color: #8C8C8C;">请联系小程序开发者获取路径,比如/pages/index/index,注意，路径建议以/开头</div>
            </a-form-item>
            <a-form-item label="小程序封面">
          <div class="upload-row">
            <a-upload :show-upload-list="false" :before-upload="handleBeforeUpload"
              :custom-request="(options) => handleUpload(options, 'received_message_images')" accept=".png,.jpg,.jpeg">
              <a-button type="default">上传图片</a-button>
            </a-upload>
            <a-input @paste="onPaste" @click="pasteImage" placeholder="复制粘贴上传图片">复制粘贴上传</a-input>
          </div>
          <div v-if="miniCard.thumb_url" class="preview-wrap">
            <img v-viewer :src="miniCard.thumb_url" class="preview" />
            <span class="del-icon" @click="removeImage"><CloseOutlined /></span>
          </div>
        </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <template #rightExtra>
        <a-button type="text" @click="emitDel">
          <CloseCircleOutlined style="font-size: 16px;" />
        </a-button>
      </template>
    </a-tabs>
  </div>
</template>

<script setup>
import { CloseCircleOutlined, UploadOutlined, CloseOutlined } from '@ant-design/icons-vue'
import emojiDataJson from 'emoji-mart-vue-fast/data/all.json'
import "emoji-mart-vue-fast/css/emoji-mart.css"
import { ref, reactive, watch, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { uploadFile } from '@/api/app'
import { Picker, EmojiIndex } from 'emoji-mart-vue-fast/src'


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

let isUploading = false
const emojiIndex = new EmojiIndex(emojiDataJson)
const labelCol = { span: 3 }
const wrapperCol = { span: 21 }
const cardLabelCol = { span: 4 }
const cardWrapperCol = { span: 20 }

watch(() => text.description, (v) => {
  const s = String(v || '')
  if (s.length > 300) {
    message.warning('最多输入300个字')
    text.description = s.slice(0, 300)
  }
})
watch(() => imageText.description, (v) => {
  const s = String(v || '')
  if (s.length > 300) {
    message.warning('最多输入300个字')
    imageText.description = s.slice(0, 300)
  }
})

const rules = {
  imageText: {
    url: [{ required: true, message: '请输入http/https地址' }, { validator: httpValidator }],
    title: [{ required: true, message: '请输入标题' }],
    description: [{ required: true, message: '请输入描述' }],
    thumb_url: [{ required: true, message: '请上传链接图片' }]
  },
  url: {
    title: [{ required: true, message: '请输入标题' }],
    url: [{ required: true, message: '请输入http/https地址' }, { validator: httpValidator }]
  },
  card: {
    title: [{ required: true, message: '请输入小程序标题' }],
    appid: [{ required: true, message: '请输入小程序appID' }],
    page_path: [{ required: true, message: '请输入小程序路径' }]
  }
}

function httpValidator (_rule, value) {
  if (!value) return Promise.resolve()
  const ok = /^https?:\/\//.test(value)
  return ok ? Promise.resolve() : Promise.reject('链接需以http或https开头')
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
      message.warning('最多输入300个字')
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
      message.warning('最多输入300个字')
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
    message.error('只支持JPG、PNG格式的图片')
    return false
  }
  const isLt2M = file.size / 1024 < 1024 * 2
  if (!isLt2M) {
    message.error('图片大小不能超过2M')
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
    if (activeTab.value === 'imageText') imageText.thumb_url = url
    else if (activeTab.value === 'image') image.thumb_url = url
    else if (activeTab.value === 'card') miniCard.thumb_url = url
    onSuccess && onSuccess(res)
    onChange()
  } catch (e) {
    message.error('上传失败')
    onError && onError(e)
  }
}

function onChange () {
  const payload = getPayload()
  emit('update:value', payload)
  emit('change', { ...payload, reply_index: props.reply_index })
}

function removeImage () {
  if (activeTab.value === 'imageText') imageText.thumb_url = ''
  else if (activeTab.value === 'image') image.thumb_url = ''
  else if (activeTab.value === 'card') miniCard.thumb_url = ''
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
      auto_menu_id: '',
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
      auto_menu_id: '',
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
      auto_menu_id: '',
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
      auto_menu_id: '',
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
      auto_menu_id: '',
      more_img_text_json: '',
      media_id: '',
      pic: miniCard.thumb_url || '',
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
    auto_menu_id: '',
    more_img_text_json: '',
    media_id: '',
    pic: '',
    type: 'text'
  }
}

// 规避某些环境下 ant popover 的 getPopupContainer 使用 triggerNode.parentNode 导致的空引用
function getPopup () { return document.body }

function pasteImage () { message.info('复制图片后按Ctrl+V可尝试粘贴图片'); }
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
  activeTab.value = t
  if (t === 'imageText') Object.assign(imageText, val)
  else if (t === 'text') Object.assign(text, { description: val?.description || '' })
  else if (t === 'image') Object.assign(image, { thumb_url: val?.thumb_url || '' })
  else if (t === 'url') Object.assign(urlCard, { title: val?.title || '', url: val?.url || '' })
  else if (t === 'card') Object.assign(miniCard, val)
}, { immediate: true })

// emit changes on input editing
watch(imageText, () => { onChange() }, { deep: true })
watch(text, () => { onChange() }, { deep: true })
watch(image, () => { onChange() }, { deep: true })
watch(urlCard, () => { onChange() }, { deep: true })
watch(miniCard, () => { onChange() }, { deep: true })

function tabLabel (txt) { return txt }
</script>

<style scoped lang="less">
.multi-reply-card {
  width: 694px;
  background: #F2F4F7;
  border-radius: 6px;
  padding: 0px 16px 16px;
  margin-bottom: 8px;
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
</style>