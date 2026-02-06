<template>
  <div class="subManage-edit">
    <a-breadcrumb class="subManage-breadcrumb">
      <a-breadcrumb-item><a :href="autoReplyUrl">{{ t('title_smart_menu') }}</a></a-breadcrumb-item>
      <a-breadcrumb-item>{{ t('title_add_rule') }}</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="main">
      <div class="left-pane">
        <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
          <a-form-item name="menu_title" :label="t('label_menu_name')" :rules="[{ required: true, message: t('msg_input_menu_name') }]">
            <a-input v-model:value="form.menu_title" :placeholder="t('ph_input')" />
          </a-form-item>

          <a-form-item name="menu_description" :label="t('label_menu_content')" :rules="[{ required: true, message: t('msg_input_menu_content') }]">
            <div class="menu-content-box">
              <a-textarea ref="menuContentRef" v-model:value="form.menu_description" :rows="6"
                :placeholder="t('ph_menu_content')" />
              <div class="insert-row">
                <span class="insert-label">{{ t('label_insert') }}</span>
                <a-dropdown>
                  <a class="insert-btn" @click.prevent>{{ t('btn_time') }}</a>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item v-for="fmt in timeFormats" :key="fmt">
                        <a-tooltip :title="formatNow(fmt)">
                          <a @click.prevent="insertTime(fmt)">{{ fmt }}</a>
                        </a-tooltip>
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
                <a-popover v-model:open="emojiOpen" placement="bottomLeft" trigger="click"
                  :getPopupContainer="getPopup">
                  <template #content>
                    <Picker :data="emojiIndex" :emojiSize="18" :showPreview="false" set="apple"
                      @select="onEmojiSelect" />
                  </template>
                  <a class="insert-btn">{{ t('btn_emoji') }}</a>
                </a-popover>
              </div>
            </div>
          </a-form-item>

          <a-form-item name="menu_content" :label="t('label_menu_title')" :rules="[{ validator: validateMenuContent }]">
            <a-input v-model:value="hiddenMenuContent" style="display:none" />
            <a-form-item-rest>
              <div class="title-list">
                <div class="title-row" v-for="(it, idx) in titleItems" :key="it.id"
                  :class="{ 'drag-over': dragOverIndex === idx }" draggable="true" @dragstart="onDragStart(idx)"
                  @dragover.prevent @dragenter.prevent="onDragEnter(idx)" @dragleave="onDragLeave(idx)"
                  @drop.prevent="onDrop(idx)" @dragend="onDragEnd">
                  <div class="title-item">
                    <div class="drag-handle">
                      <HolderOutlined />
                    </div>
                    <div v-if="!it.isPure && !it.linkType" class="index-badge" v-show="!it.editingIndex"
                      @click="startEditOrder(idx)">
                      {{ it.orderDisplay }}</div>
                    <input v-if="!it.isPure && !it.linkType && it.editingIndex" class="index-input-edit" type="number" min="1" max="10"
                      inputmode="numeric" autofocus v-model="it.orderDisplayStr" @blur="endEditOrder(idx)"
                      @keyup.enter="endEditOrder(idx)" />
                    <div class="input-box" :class="{ focused: focusedId === it.id }">
                      <a-input v-model:value="it.text" :disabled="it.isNewline === true" :placeholder="t('ph_menu_title')"
                        maxLength="40" @focus="onFocusTitle(idx, $event)" @blur="onBlurTitle(idx)">
                        <template #suffix>
                          <a v-if="!it.isPure && !it.linkType && (!!(it.replyConfig?.rule_id) || !!it.rule_id || !!it.replyConfig)" class="suffix-action" @click="openReplyModal(idx)">{{ t('btn_edit') }}</a>
                          <a-tooltip v-else-if="!it.isPure && !it.linkType && !keywordReplyStatus">
                            <template #title>{{ t('tip_need_auto_reply') }}</template>
                            <a class="suffix-action disabled" @click.prevent>{{ t('btn_add_reply') }}</a>
                          </a-tooltip>
                          <a v-else-if="!it.isPure && !it.linkType" class="suffix-action" @click="openReplyModal(idx)">{{ t('btn_add_reply') }}</a>
                        </template>
                      </a-input>
                      <div v-if="focusedId === it.id" class="float-toolbar" @mousedown.stop="onToolbarMouseDown">
                        <a @mousedown.prevent.stop="openMiniLinkModal(idx)">{{ t('btn_insert_mini_link') }}</a>
                        <a @mousedown.prevent.stop="openUrlLinkModal(idx)">{{ t('btn_insert_url_link') }}</a>
                        <a v-if="!it.isPure && !it.linkType" @mousedown.prevent.stop="convertToPure(idx)">{{ t('btn_to_plain_text') }}</a>
                        <a v-if="!it.isPure && !it.linkType" @mousedown.prevent.stop="insertNewline(idx)">{{ t('btn_newline') }}</a>
                      </div>
                    </div>
                  </div>
                  <div class="row-icons">
                    <svg-icon @click="addTitleAfter(idx)" name="smart-add" size="16"
                      style="color: #8C8C8C; cursor: pointer;"></svg-icon>
                    <svg-icon @click="removeTitle(idx)" name="smart-delete" size="16"
                      style="color: #8C8C8C; cursor: pointer;"></svg-icon>
                  </div>
                </div>
              </div>
            </a-form-item-rest>
          </a-form-item>

          <div class="btn-container">
            <a-button type="primary" @click="onSubmit">{{ t('btn_save') }}</a-button>
          </div>
        </a-form>
      </div>
      <div class="right-pane">
        <img :src="smartMenuRuleImgSrc" alt="smart-menu-preview" class="preview-smart-img" />
      </div>
    </div>

    <a-modal v-model:open="replyModalOpen" :title="`${replyForm.rule_id ? t('title_edit_menu_reply') : t('title_add_menu_reply')}`" :width="720" @ok="onReplyModalOk"
      @cancel="onReplyModalCancel">
      <a-form :model="replyForm" layout="vertical" class="reply-modal-form">
        <a-form-item :label="t('label_rule_name')" name="rule_name" :rules="[{ required: true, message: t('msg_input_rule_name') }]">
          <a-input v-model:value="replyForm.rule_name" />
        </a-form-item>
        <a-form-item :label="t('label_keyword')" name="keyword" :rules="[{ required: true, message: t('msg_input_keyword') }]">
          <a-input v-model:value="replyForm.keyword" :maxlength="30" :placeholder="t('ph_keyword')" @focus="onKeywordFocus" @pressEnter="onKeywordFocus" @blur="onKeywordFocus">
            <template #addonBefore>
              <a-select v-model:value="replyForm.match_type" style="width: 120px">
                <a-select-option value="full">{{ t('select_full_match') }}</a-select-option>
                <a-select-option value="partial">{{ t('select_partial_match') }}</a-select-option>
              </a-select>
            </template>
          </a-input>
          <div class="tip-box" style="margin-top: 4px;">{{ t('tip_keyword_max_length') }}</div>
        </a-form-item>

        <a-form-item name="replyList" :label="t('label_reply_content')" :rules="[{ required: true, message: '' }]">
          <div class="item-box">
            <MultiReply v-for="(it, idx) in replyForm.replyList" :key="idx" ref="replyRefs" v-model:value="replyForm.replyList[idx]"
              :reply_index="idx" :show_smart="true" @change="onReplyContentChange" @del="onReplyDelItem" />
            <a-button type="dashed" style="width: 100%;" :disabled="replyForm.replyList.length >= 5"
              @click="addReplyContent">
              <template #icon>
                <PlusOutlined />
              </template>
              {{ t('btn_add_reply_content') }}({{ replyForm.replyList.length }}/5)
            </a-button>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="miniLinkModalOpen"  :title="t('title_edit_mini_card')" :width="720" @ok="onMiniLinkOk"
      @cancel="onMiniLinkCancel">
      <a-form :model="miniLinkForm" layout="vertical">
        <a-form-item :label="t('label_card_title')" name="title" :rules="[{ required: true, message: t('msg_input_card_title') }]">
          <a-input v-model:value="miniLinkForm.title" :placeholder="t('ph_card_title')" />
        </a-form-item>
        <a-form-item :label="t('label_appid')" name="appid"
          :rules="[{ required: true, message: t('msg_input_appid') }, { validator: appidValidator }]">
          <a-input v-model:value="miniLinkForm.appid" :placeholder="t('ph_appid')" />
          <div class="tips-box">
            <div>{{ t('tip_appid_1') }}<a @click="openDocs('appid')">{{ t('btn_find_appid') }}</a></div>
            <div>{{ t('tip_appid_2') }}<a @click="openDocs('assoc')">{{ t('btn_view_assoc') }}</a></div>
          </div>
        </a-form-item>
        <a-form-item :label="t('label_mini_path')" name="path" :rules="[{ required: true, message: t('msg_input_mini_path') }]">
          <a-input v-model:value="miniLinkForm.path" :placeholder="t('ph_mini_path')" />
          <div class="tips-box">
            <div>{{ t('tip_mini_path') }}</div>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="urlLinkModalOpen" :title="t('title_edit_link')" :width="720" @ok="onUrlLinkOk" @cancel="onUrlLinkCancel">
      <a-form :model="urlLinkForm" :labelCol="labelCol" :wrapperCol="wrapperCol">
        <a-form-item :label="t('label_title')" name="title" :rules="[{ required: true, message: t('msg_input_title') }]">
          <a-input v-model:value="urlLinkForm.title" :placeholder="t('ph_title')" />
        </a-form-item>
        <a-form-item :label="t('label_open_link')" name="url"
          :rules="[{ required: true, message: t('msg_input_open_link') }, { validator: urlHttpValidator }]">
          <a-input v-model:value="urlLinkForm.url" :placeholder="t('ph_open_link')" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script setup>
import { PlusOutlined, HolderOutlined } from '@ant-design/icons-vue'
import { onMounted, ref, reactive, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import MultiReply from '@/components/replay-card/multi-reply.vue'
import { getSmartMenu, saveSmartMenu, checkKeyWordRepeat, saveRobotKeywordReply, getRobotKeywordReply } from '@/api/explore/index.js'
import dayjs from 'dayjs'
import emojiDataJson from 'emoji-mart-vue-fast/data/all.json'
import { Picker, EmojiIndex } from 'emoji-mart-vue-fast/src'
import smartMenuRuleImgSrc from '@/assets/img/robot/smart-menu-rule.png'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.smart-menu.add-rule')
const labelCol = { span: 3 }
const wrapperCol = { span: 21 }
const query = useRoute().query
const ruleId = ref(+query.rule_id || +query['rule-id'] || 0)
const router = useRouter()
const autoReplyUrl = computed(() => `/#/robot/ability/smart-menu?id=${query.id}&robot_key=${query.robot_key}`)
const robotStore = useRobotStore()
const keywordReplyStatus = computed(() => robotStore.keywordReplySwitchStatus === '1')
const formRef = ref(null)
const form = reactive({
  menu_title: '',
  menu_description: '',
})
const rules = {}
const replyRefs = ref([])
const emojiIndex = new EmojiIndex(emojiDataJson)
const emojiOpen = ref(false)
const menuContentRef = ref(null)
const MAX_TITLES = 10
const timeFormats = [
  'yyyy-MM-dd hh:mm:ss',
  'MM-dd hh:mm:ss',
  'yyyy-MM-dd hh:mm',
  'MM-dd hh:mm',
  'yyyy-MM-dd',
  'MM月dd日'
]
function formatNow (fmt) {
  const map = {
    'yyyy-MM-dd hh:mm:ss': 'YYYY-MM-DD HH:mm:ss',
    'MM-dd hh:mm:ss': 'MM-DD HH:mm:ss',
    'yyyy-MM-dd hh:mm': 'YYYY-MM-DD HH:mm',
    'MM-dd hh:mm': 'MM-DD HH:mm',
    'yyyy-MM-dd': 'YYYY-MM-DD',
    'MM月dd日': 'MM月DD日'
  }
  return dayjs().format(map[fmt] || 'YYYY-MM-DD HH:mm:ss')
}
function insertAtCaretInMenuContent (text) {
  const el = menuContentRef.value?.$el?.querySelector('textarea')
  if (el && typeof el.selectionStart === 'number') {
    const start = el.selectionStart
    const end = el.selectionEnd
    const val = form.menu_description || ''
    form.menu_description = val.slice(0, start) + text + val.slice(end)
    nextTick(() => {
      el.focus()
      const pos = Math.min(start + text.length, (form.menu_description || '').length)
      el.setSelectionRange(pos, pos)
    })
  } else {
    form.menu_description = (form.menu_description || '') + text
  }
}
function insertTime (fmt) {
  insertAtCaretInMenuContent(`{{${fmt}}}`)
}
function onEmojiSelect (emoji) {
  const char = emoji?.native || '';
  if (!char) return;
  insertAtCaretInMenuContent(char);
  emojiOpen.value = false
}
function getPopup () {
  return document.body
}

function urlHttpValidator (_rule, value) {
  if (!value) return Promise.resolve()
  return /^https?:\/\//.test(String(value)) ? Promise.resolve() : Promise.reject(t('msg_link_http_error'))
}

function appidValidator (_rule, value) {
  if (!value) return Promise.resolve()
  return /^wx/.test(String(value)) ? Promise.resolve() : Promise.reject(t('msg_appid_error'))
}

function onKeywordFocus () {
  const kw = String(replyForm.keyword || '').trim()
  if (!kw) return
  checkKeyWordRepeat({ id: replyForm.rule_id || 0, robot_id: query.id, keyword: kw })
    .then((res) => {
      const repeat = res?.data?.is_repeat
      const ruleName = res?.data?.rule_name || ''
      if (repeat) message.error(t('msg_keyword_repeat', { ruleName }))
    })
}

function serializeReplyContent (list) { return (list || []).map((it) => ({ ...it, status: '1' })) }
function serializeReplyTypeCodes (list) { const map = { text: '2', image: '4', card: '3', imageText: '1', url: '5', smartMenu: '6' }; return (list || []).map((it) => map[it.type] || '').filter(Boolean) }

function validateMenuContent (_rule, _value) {
  const ok = Array.isArray(titleItems.value) && titleItems.value.length > 0
  return ok ? Promise.resolve() : Promise.reject(t('msg_add_menu_title'))
}

const titleItems = ref([
  {
    id: Date.now(),
    text: '',
    isPure: false,
    orderDisplay: 1,
    replyConfig: null,
    editingIndex: false,
    orderDisplayStr: ''
  }
])
const hiddenMenuContent = ref('')
const focusedId = ref(null)
const toolbarInteracting = ref(false)
let draggingIndex = -1
const dragOverIndex = ref(-1)
function onDragStart (idx) { draggingIndex = idx }
function onDragEnter (idx) {
  if (draggingIndex === -1 || draggingIndex === idx) {
    dragOverIndex.value = idx; return
  }
  const list = [...titleItems.value]
  const [moved] = list.splice(draggingIndex, 1)
  list.splice(idx, 0, moved)
  titleItems.value = list
  draggingIndex = idx
  dragOverIndex.value = idx
}
function onDragLeave (idx) {
  if (dragOverIndex.value === idx) dragOverIndex.value = -1
}
function onDrop (idx) {
  if (draggingIndex === idx || draggingIndex < 0) return
  const list = [...titleItems.value]
  const [moved] = list.splice(draggingIndex, 1)
  list.splice(idx, 0, moved)
  titleItems.value = list
  draggingIndex = -1
  dragOverIndex.value = -1
}
function onDragEnd () {
  draggingIndex = -1;
  dragOverIndex.value = -1
}
function onFocusTitle (idx) {
  const it = titleItems.value[idx];
  focusedId.value = it?.id ?? null
}
function onBlurTitle (idx) {
  const it = titleItems.value[idx]
  setTimeout(() => {
    if (toolbarInteracting.value) {
      toolbarInteracting.value = false; return
    }
    if ((focusedId.value ?? null) === (it?.id ?? null)) focusedId.value = null
  }, 0)
}
function onToolbarMouseDown () {
  toolbarInteracting.value = true
}
function insertNewline (idx) {
  const it = titleItems.value[idx]
  it.text = t('text_newline')
  it.isPure = true
  it.isNewline = true
  it.replyConfig = null
  focusedId.value = null
}
function convertToPure (idx) {
  const it = titleItems.value[idx];
  it.isPure = true; it.replyConfig = null
}
function addTitleAfter (idx) {
  const count = titleItems.value.length
  if (count >= MAX_TITLES) { message.warning(t('msg_menu_title_max')); return }
  const nextOrder = count + 1
  titleItems.value.splice(idx + 1, 0, { id: Date.now() + Math.random(), text: '', isPure: false, orderDisplay: nextOrder, replyConfig: null, editingIndex: false, orderDisplayStr: '' })
}
function removeTitle (idx) {
  if (titleItems.value.length <= 1) { message.warning(t('msg_menu_title_min')); return }
  titleItems.value.splice(idx, 1)
}
function startEditOrder (idx) {
  const it = titleItems.value[idx];
  it.orderDisplayStr = String(it.orderDisplay || '');
  it.editingIndex = true;
  nextTick(() => {

  })
}
function endEditOrder (idx) {
  const it = titleItems.value[idx];
  const n = parseInt(String(it.orderDisplayStr || '').replace(/\D+/g, ''), 10);
  it.orderDisplay = isNaN(n) ? it.orderDisplay : Math.max(1, Math.min(10, n));
  it.editingIndex = false
}

watch(titleItems, (list) => {
  try {
    const payload = (list || []).map((it) => ({ menu_type: it.isPure ? '0' : '1', content: String(it.text || '') }))
    hiddenMenuContent.value = JSON.stringify(payload)
  } catch (_) {
    hiddenMenuContent.value = String((list || []).length || '')
  }
}, { deep: true, immediate: true })

const replyModalOpen = ref(false)
const replyForm = reactive({
  rule_name: '',
  match_type: 'full',
  keyword: '',
  replyList: [
    {
      type: 'text', description: ''
    }
  ]
})
let editingTitleIndex = -1
function openReplyModal (idx) {
  if (!titleItems.value[idx].text) { message.error(t('msg_input_menu_title')); return }
  editingTitleIndex = idx
  const it = titleItems.value[idx]
  replyForm.rule_name = it.replyConfig?.rule_name || String(it.text || '')
  replyForm.keyword = it.replyConfig?.keyword || String(it.text || '')
  replyForm.match_type = it.replyConfig?.match_type || 'full'
  replyForm.replyList = Array.isArray(it.replyConfig?.replyList) ? JSON.parse(JSON.stringify(it.replyConfig.replyList)) : [{ type: 'text', description: '' }]
  replyForm.rule_id = it.replyConfig?.rule_id || it.rule_id || ''
  if (replyForm.rule_id) {
    getRobotKeywordReply({ id: replyForm.rule_id })
      .then((res) => {
        const data = res?.data || {}
        replyForm.rule_name = data?.name || replyForm.rule_name
        const fk = Array.isArray(data?.full_keyword) ? data.full_keyword : []
        const hk = Array.isArray(data?.half_keyword) ? data.half_keyword : []
        if (fk.length > 0) {
          replyForm.match_type = 'full'
          replyForm.keyword = String(fk[0] || '')
        } else if (hk.length > 0) {
          replyForm.match_type = 'partial'
          replyForm.keyword = String(hk[0] || '')
        }
        const list = Array.isArray(data?.reply_content) ? data.reply_content : []
        replyForm.replyList = list.map((rc) => ({
          type: (rc?.type || rc?.reply_type || 'text'),
          description: rc?.description || '',
          thumb_url: rc?.thumb_url || rc?.pic || '',
          title: rc?.title || '',
          url: rc?.url || '',
          appid: rc?.appid || '',
          page_path: rc?.page_path || '',
          smart_menu_id: rc?.smart_menu_id || '',
          smart_menu: rc?.smart_menu || {},
        }))
      })
      .catch(() => {})
  }
  replyModalOpen.value = true
}
async function onReplyModalOk () {
  if (editingTitleIndex < 0) { replyModalOpen.value = false; return }
  const name = String(replyForm.rule_name || '').trim()
  const kw = String(replyForm.keyword || '').trim()
  if (!name || !kw) { message.warning(t('msg_complete_rule_keyword')); return }
  if (!Array.isArray(replyForm.replyList) || replyForm.replyList.length === 0) { message.warning(t('msg_add_reply_item')); return }
  for (const comp of replyRefs.value) {
    if (comp && comp.validate) {
      const ok = await comp.validate()
      if (!ok) { return }
    }
  }
  checkKeyWordRepeat({ id: replyForm.rule_id || 0, robot_id: query.id, keyword: kw })
    .then((res) => {
      const repeat = res?.data?.is_repeat
      const ruleName = res?.data?.rule_name || ''
      if (repeat) { message.error(t('msg_keyword_repeat', { ruleName })); return Promise.reject('repeat') }
      const payload = {
        robot_id: query.id,
        name,
        full_keyword: replyForm.match_type === 'full' ? [kw] : [],
        half_keyword: replyForm.match_type === 'partial' ? [kw] : [],
        reply_content: JSON.stringify(serializeReplyContent(replyForm.replyList)),
        reply_type: serializeReplyTypeCodes(replyForm.replyList),
        reply_num: '0',
        forced_enable: 1
      }
      if (replyForm.rule_id) payload.id = replyForm.rule_id
      return saveRobotKeywordReply(payload)
    })
    .then((res) => {
      if (!(res && res.res == 0)) { message.error(t('msg_save_failed')); return }
      const newId = (res?.data?.id) || ''
      const cfg = {
        rule_name: replyForm.rule_name,
        keyword: replyForm.keyword,
        match_type: replyForm.match_type,
        replyList: JSON.parse(JSON.stringify(replyForm.replyList)),
        rule_id: String(newId || replyForm.rule_id || '')
      }
      titleItems.value[editingTitleIndex].replyConfig = cfg
      replyForm.rule_id = cfg.rule_id
      replyModalOpen.value = false
      message.success(t('msg_save_success'))
    })
    .catch((e) => { if (e !== 'repeat') { /* ignore */ } })
}
function onReplyModalCancel () { replyModalOpen.value = false }
function addReplyContent () {
  if (replyForm.replyList.length >= 5) return;
  replyForm.replyList.push({ type: 'text', description: '' })
}
function onReplyDelItem (idx) { replyForm.replyList.splice(idx, 1) }
function onReplyContentChange ({ reply_index, ...rest }) {
  if (reply_index >= 0 && reply_index < replyForm.replyList.length) replyForm.replyList[reply_index] = rest
}

const miniLinkModalOpen = ref(false)
const miniLinkForm = reactive({ title: '', appid: '', path: '' })
const urlLinkModalOpen = ref(false)
const urlLinkForm = reactive({ title: '', url: '' })
function openDocs (t) {
  if (t === 'appid') {
    window.open('https://www.kancloud.cn/wikizhima/wikixkf/1011121', '_blank')
  } else {
    window.open('https://www.kancloud.cn/wikizhima/wikixkf/1010959', '_blank')
  }
}
function openMiniLinkModal (idx) {
  editingTitleIndex = idx;
  const it = titleItems.value[idx];
  const raw = String(it.text || '');
  let title = raw;
  let appid = '';
  let path = '';
  if (/<a[\s\S]*?<\/a>/i.test(raw)) {
    const inner = raw.match(/<a[\s\S]*?>([\s\S]*?)<\/a>/i)
    const am = raw.match(/data-miniprogram-appid=['"][\s\S]*?['"]/i)
    const pm = raw.match(/data-miniprogram-path=['"][\s\S]*?['"]/i)
    if (inner && inner[1] != null) title = inner[1]
    if (am) {
      const m = String(am[0] || '').match(/data-miniprogram-appid=['"]\s*([^'"]+)\s*['"]/i)
      if (m && m[1] != null) appid = m[1]
    }
    if (pm) {
      const m = String(pm[0] || '').match(/data-miniprogram-path=['"]\s*([^'"]+)\s*['"]/i)
      if (m && m[1] != null) path = m[1]
    }
  }
  miniLinkForm.title = String(title || '')
  miniLinkForm.appid = String(appid || '')
  miniLinkForm.path = String(path || '')
  miniLinkModalOpen.value = true
}
function onMiniLinkOk () {
  const t = String(miniLinkForm.title || '').trim()
  const a = String(miniLinkForm.appid || '').trim()
  const p = String(miniLinkForm.path || '').trim()
  if (!t || !a || !p) {
    message.warning(t('msg_complete_mini_card')); return
  }
  if (!/^wx/.test(a)) {
    message.warning(t('msg_appid_error')); return
  }
  if (editingTitleIndex < 0) {
    miniLinkModalOpen.value = false; return
  }
  const it = titleItems.value[editingTitleIndex]
  it.text = `<a  href='#' data-miniprogram-path='${p}' data-miniprogram-appid='${a}'  >${t}</a>`
  it.linkType = 'mini'
  it.replyConfig = null
  miniLinkModalOpen.value = false
}
function onMiniLinkCancel () {
  miniLinkModalOpen.value = false
}
function openUrlLinkModal (idx) {
  editingTitleIndex = idx;
  const it = titleItems.value[idx];
  const raw = String(it.text || '')
  let title = raw
  let url = ''
  if (/<a[\s\S]*?<\/a>/i.test(raw)) {
    const inner = raw.match(/<a[\s\S]*?>([\s\S]*?)<\/a>/i)
    const hm = raw.match(/href\s*=\s*['"][\s\S]*?['"]/i)
    if (inner && inner[1] != null) title = inner[1]
    if (hm) {
      const m = String(hm[0] || '').match(/href\s*=\s*['"]\s*([^'"]+)\s*['"]/i)
      if (m && m[1] != null) url = m[1]
    }
    if (url) url = url.replace(/`/g, '').trim()
  }
  urlLinkForm.title = String(title || '')
  urlLinkForm.url = String(url || '')
  urlLinkModalOpen.value = true
}
function onUrlLinkOk () {
  const t = String(urlLinkForm.title || '').trim()
  const u = String(urlLinkForm.url || '').trim()
  if (!t || !u) { message.warning(t('msg_complete_link')); return }
  if (!/^https?:\/\//.test(u)) { message.warning(t('msg_link_http_error')); return }
  if (editingTitleIndex < 0) { urlLinkModalOpen.value = false; return }
  const it = titleItems.value[editingTitleIndex]
  it.text = `<a href="${u}">${t}</a>`
  it.linkType = 'url'
  it.replyConfig = null
  urlLinkModalOpen.value = false
}
function onUrlLinkCancel () {
  urlLinkModalOpen.value = false
}

watch(() => form.duration_type, (val) => {
  if (val === 'week') {
    if (!Array.isArray(form.week_day) || form.week_day.length === 0) {
      form.week_day = ['1']
    }
    form.week_duration = Array.isArray(form.week_day) ? form.week_day : [form.week_day]
  } else if (val === 'time_range') {
    if (Array.isArray(form.date_range) && form.date_range.length === 2 && form.date_range[0] && form.date_range[1]) {
      form.start_day = dayjs(form.date_range[0]).format('YYYY-MM-DD')
      form.end_day = dayjs(form.date_range[1]).format('YYYY-MM-DD')
    } else {
      form.start_day = ''
      form.end_day = ''
    }
  } else {
    form.start_day = ''
    form.end_day = ''
  }
})

watch(() => form.week_day, (v) => {
  if (form.duration_type === 'week') {
    form.week_duration = Array.isArray(v) ? v : [v]
  }
})

watch(() => form.date_range, (v) => {
  if (form.duration_type === 'time_range' && Array.isArray(v) && v.length === 2 && v[0] && v[1]) {
    form.start_day = dayjs(v[0]).format('YYYY-MM-DD')
    form.end_day = dayjs(v[1]).format('YYYY-MM-DD')
  }
})

const onSubmit = () => {
  formRef.value?.validate().then(async () => {
    if (!form.menu_title || !form.menu_description) { message.warning(t('msg_complete_menu_content')); return }
    const menu_content = titleItems.value.map((it) => {
      const hasLink = !!it.linkType
      const isPure = !!it.isPure
      const isNewline = it.isNewline === true
      const text = String(it.text || '')
      if (isNewline || text === '') { return { menu_type: '0', serial_no: '', content: '', rule_id: '' } }
      if (hasLink) { return { menu_type: '0', serial_no: '', content: text, rule_id: '' } }
      if (isPure) { return { menu_type: '0', serial_no: '', content: text, rule_id: '' } }
      const rid = String(it.replyConfig?.id || it.replyConfig?.rule_id || '')
      const sno = String(it.orderDisplay || '')
      return { menu_type: '1', serial_no: sno, content: text, rule_id: rid }
    })
    const payload = {
      robot_id: query.id,
      menu_title: form.menu_title,
      menu_description: form.menu_description,
      menu_content: JSON.stringify(menu_content)
    }
    if (ruleId.value) payload.id = ruleId.value
    try {
      const res = await saveSmartMenu(payload)
      if (res && res.res == 0) {
        message.success(t('msg_save_success'))
        router.push({ path: '/robot/ability/smart-menu', query: { id: query.id, robot_key: query.robot_key } })
      }
    } catch (e) {
      message.error(t('msg_save_failed'))
    }
  })
}

onMounted(async () => {
  const copyId = +(query.copy_id || 0)
  const fetchOne = async (rid) => {
    const res = await getSmartMenu({ id: rid, robot_id: query.id })
    const data = res?.data || {}
    form.menu_title = data.menu_title || ''
    form.menu_description = data.menu_description || data.menu_content || ''
    if (Array.isArray(data.menu_content)) {
      const arr = data.menu_content
      titleItems.value = arr.map((mc, i) => {
        const t = String(mc?.menu_type || '')
        const txt = String(mc?.content || '')
        const isAnchor = /<a[\s\S]*?<\/a>/.test(txt)
        const isMini = isAnchor && /data-miniprogram-appid/i.test(txt)
        return {
          id: Date.now() + i,
          text: (t === '0' && txt === '') ? t('text_newline') : txt,
          isPure: t === '0',
          isNewline: (t === '0' && txt === ''),
          orderDisplay: ((mc && mc.serial_no !== undefined && String(mc.serial_no).trim() !== '') ? Number(mc.serial_no) : (i + 1)),
          editingIndex: false,
          orderDisplayStr: '',
          replyConfig: null,
          linkType: isAnchor ? (isMini ? 'mini' : 'url') : '',
          rule_id: String(mc?.rule_id || '')
        }
      })

      const needLoad = titleItems.value
        .map((it, idx) => ({ idx, rid: String(it.rule_id || '') }))
        .filter((x) => !!x.rid)
      for (const item of needLoad) {
        try {
          const resOne = await getRobotKeywordReply({ id: item.rid })
          const dataOne = resOne?.data || {}
          const list = Array.isArray(dataOne?.reply_content) ? dataOne.reply_content : []
          titleItems.value[item.idx].replyConfig = {
            rule_name: dataOne?.name || String(titleItems.value[item.idx].text || ''),
            keyword: String(titleItems.value[item.idx].text || ''),
            match_type: 'full',
            replyList: list.map((rc) => ({
              type: (rc?.type || rc?.reply_type || 'text'),
              description: rc?.description || '',
              thumb_url: rc?.thumb_url || rc?.pic || '',
              title: rc?.title || '',
              url: rc?.url || '',
              appid: rc?.appid || '',
              page_path: rc?.page_path || '',
              smart_menu_id: rc?.smart_menu_id || '',
              smart_menu: rc?.smart_menu || {},
            })),
            rule_id: item.rid
          }
        } catch (_) {}
      }
    }
  }
  try {
    if (!ruleId.value && copyId) { await fetchOne(copyId); return }
    if (!ruleId.value) return
    await fetchOne(ruleId.value)
  } catch (e) { message.error(t('msg_load_failed')) }
})
</script>
<style lang="less" scoped>
.subManage-edit {
  padding: 16px 24px;
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #fff;
  overflow-x: hidden;
  overflow-y: auto;

  .subManage-breadcrumb {
    display: flex;
    align-items: center;
    color: #000000;
    font-family: "PingFang SC";
    font-size: 14px;
    font-style: normal;
    line-height: 22px;
    padding-bottom: 16px;
  }

  .main {
    padding: 0 8px;
    border-radius: 6px;
    background-color: white;
    padding-bottom: 24px;
    display: flex;
    align-items: flex-start;
    justify-content: space-between;

    .title {
      border-radius: 6px;
      padding: 12px 0 12px 24px;
      align-items: flex-start;
      border-bottom: 1px solid var(--07, #F0F0F0);
      background: #FFF;
      display: flex;
      align-items: center;
      color: #262626;
      font-family: "PingFang SC";
      font-size: 14px;
      font-style: normal;
      font-weight: 600;
      line-height: 22px;
      gap: 8px;
      margin-bottom: 24px;
    }
  }

  .left-pane {
    width: 650px;
  }

  .right-pane {
    width: 240px;
    display: flex;
    justify-content: center;
  }

  .preview-smart-img {
    width: 240px;
    height: 520px;
    object-fit: contain;
    border-radius: 6px;
  }

  .left-pane ::v-deep(.ant-form-item) {
    margin-bottom: 16px;
  }

  .mr-8 {
    margin-right: 8px;
  }

  .mr16 {
    margin-right: 16px;
  }

  .nav-box {
    color: #262626;
    font-size: 14px;
    font-style: normal;
    font-weight: 600;
    line-height: 22px;
    margin-bottom: 4px;
  }

  .flex {
    display: flex;
  }

  .btn-container {
    position: fixed;
    bottom: 0;
    right: 16px;
    display: flex;
    width: calc(100% - 270px);
    padding: 16px 1055px 16px 32px;
    align-items: center;
    border-radius: 0 0 2px 2px;
    background: #FFF;
    box-shadow: 0 -8px 4px 0 #0000000a;
  }

  .flex-center {
    display: flex;
    align-items: center;
  }

  .ml4 {
    margin-left: 4px;
  }

  .tip-box {
    color: #8c8c8c;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
    white-space: wrap;
    max-width: 624px;
  }
}
/* 移除全局 ::v-deep(.ant-form-item-label) 以避免影响子组件（例如 MultiReply） */

.interval-item ::v-deep(.ant-form-item-label) {
  width: 130px;
  min-width: 130px;
  flex: 0 0 130px;
}

.message-type-item {
  margin-bottom: 0;
}

/* 避免影响 MultiReply 内部表单项的间距 */

.menu-content-box {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px;
  background: #F2F4F7;
  border-radius: 6px;
}

.menu-content-box ::v-deep(.ant-input) {
  background: #fff;
  border-color: #D9D9D9;
  border-radius: 6px;
}

.menu-content-box ::v-deep(.ant-input:focus),
.menu-content-box ::v-deep(.ant-input-focused) {
  border-color: #2475FC;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.insert-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.insert-label {
  color: #262626;
}

.insert-btn {
  color: #595959;
  border: 1px solid #D9D9D9;
  background: #fff;
  border-radius: 6px;
  padding: 2px 10px;
  display: inline-block;
}

.insert-btn:hover {
  border-color: #C9CED6;
  background: #fff;
}

.title-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-item {
  box-sizing: border-box;
  height: 30px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  background: #fff;
  flex: 1;
}

.title-row.drag-over .title-item {
  border-color: #2475FC;
  background: rgba(24, 144, 255, 0.04);
}

.drag-handle {
  color: #8c8c8c;
  cursor: grab;
}

.index-input {
  width: 60px;
}

.index-badge {
  width: 20px;
  height: 18px;
  border-radius: 6px;
  background: #D8DDE5;
  color: #242933;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  cursor: pointer;
}

.index-input-edit {
  width: 20px;
  height: 18px;
  border-radius: 6px;
  background: #D8DDE5;
  border: none;
  text-align: center;
  font-size: 12px;
  outline: none;
  padding: 0;
}

.input-box {
  position: relative;
  flex: 1;
}

.float-toolbar {
  position: absolute;
  top: -34px;
  left: 0;
  display: flex;
  gap: 12px;
  background: #fff;
  border: 1px solid #F0F0F0;
  border-radius: 6px;
  padding: 4px 8px;
}

.suffix-action {
  color: #2475FC;
}
.suffix-action.disabled {
  color: #bfbfbf;
  cursor: not-allowed;
}

.actions {
  width: 80px;
  text-align: right;
}

.row-icons {
  display: flex;
  gap: 18px;
}

.input-box ::v-deep(.ant-input-affix-wrapper),
.input-box ::v-deep(.ant-input) {
  height: 22px;
  line-height: 22px;
  border: none;
  box-shadow: none;
  background: transparent;
  padding: 0 4px;
}

.input-box ::v-deep(.ant-input-affix-wrapper:hover),
.input-box ::v-deep(.ant-input:hover),
.input-box ::v-deep(.ant-input-affix-wrapper-focused),
.input-box ::v-deep(.ant-input:focus),
.input-box ::v-deep(.ant-input-focused) {
  border: none;
  box-shadow: none;
}

.tips-box {
  margin-top: 4px;
  display: flex;
  flex-direction: column;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

::v-deep(.ant-form-vertical .ant-form-item-label) {
  width: 100%;
  min-width: 0;
  flex: initial;
  padding-bottom: 4px;
}

.reply-modal-form ::v-deep(.ant-form-item-label) {
  width: 100%;
  min-width: 0;
  flex: initial;
  padding-bottom: 4px;
}

.reply-modal-form ::v-deep(.ant-form-item) {
  margin-bottom: 16px;
}
</style>