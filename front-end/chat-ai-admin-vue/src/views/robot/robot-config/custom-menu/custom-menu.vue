<template>
  <div class="user-model-page">
    <!-- <div class="page-title">关注后自动回复</div> -->
    <div class="breadcrumb-wrap">
      <svg-icon @click="goBack" name="back" style="font-size: 20px;" />
      <div @click="goBack" class="breadcrumb-title">{{ t('title_custom_menu') }}</div>
      <a-switch
        :checked="abilitySwitchChecked"
        :checked-children="t('btn_switch_on')"
        :un-checked-children="t('btn_switch_off')"
        @change="onAbilitySwitchChange"
      />
      <span class="switch-tip">{{ t('tip_switch_description') }}</span>
    </div>
    <a-alert>
      <template #message>
        <p class="text_color_ed744a">{{ t('alert_tip_1') }}</p>
        <p>{{ t('alert_tip_2') }}</p>
        <p>{{ t('alert_tip_3') }}</p>
        <p>{{ t('alert_tip_4') }}</p>
      </template>
    </a-alert>
    <!-- 公众号列表 -->
    <div class="mp-list-block">
      <div class="mp-list" :class="{ expanded }" ref="mpListRef">
        <div class="mp-card" v-for="mp in (expanded ? mpAccounts : mpAccounts.slice(0, visibleCount))" :key="mp.id"
          :class="{ selected: mp.appid === selectedAppid }" @click="selectMp(mp)">
          <img :src="mp.logo" class="mp-logo" />
          <span class="mp-name">{{ mp.name }}</span>
        </div>
        <a-button v-if="!expanded && mpAccounts.length > visibleCount" type="dashed" class="more-btn"
          @click="expanded = true">
          {{ t('btn_more') }} +{{ mpAccounts.length - visibleCount }}
        </a-button>
      </div>
    </div>
    <div class="search-block">
      <Tooltip :title="t('tip_sync_menu')">
        <a-button @click="syncMenu" :loading="syncLoading">
          {{ t('btn_sync_menu') }}
          <QuestionCircleOutlined style="font-size: 16px;" />
        </a-button>
      </Tooltip>

      <!-- 开关控制 -->
      <div class="flex" style="gap: 4px; align-items: center;">
        <a-switch v-model:checked="appMenuSwitch" :checkedValue="'1'" :un-checkedValue="'0'" :checked-children="t('btn_switch_on')" :un-checked-children="t('btn_switch_off')" @change="onAppSwitchChange" />
        <span style="font-size: 14px; color: #8c8c8c;">{{ t('tip_close_menu_warning') }}</span>
      </div>
    </div>

    <div class="main">
      <div class="main-left">
        <div class="iphone-mock" :style="{ backgroundImage: `url(${iphoneBg})` }">
          <div class="bottom-menu" ref="bottomMenuRef">
            <Draggable :list="menus" item-key="id" class="root-menu-draggable" :animation="200" :ghost-class="'ghost'"
              @start="onRootDragStart" @end="onRootDragEnd">
              <template #item="{ element, index }">
                <div class="root-menu-item" :class="{ active: activeRootIndex === index && activeSubIndex === -1 }"
                  @click="onSelectRoot(index)">
                  <span class="name">{{ element.menu_name || t('label_root_menu') }}</span>
                  <svg-icon class="del-root" name="delete-line" @click.stop="removeRootMenu(index)"
                    style="font-size: 18px;" />
                  <div v-if="activeRootIndex === index" class="submenu-panel" @click.stop>
                    <div class="submenu-card" :style="{ width: rootItemWidth + 'px' }">
                      <Draggable :list="element.sub_menu_list" item-key="'k_'+index" class="sub-menu-draggable"
                        :animation="200" :ghost-class="'ghost'" @start="onSubDragStart(index)" @end="onSubDragEnd(index)">
                        <template #item="{ element: sub, index: si }">
                          <div class="sub-menu-item"
                            :class="{ active: activeRootIndex === index && activeSubIndex === si }"
                            @click.stop="onSelectSub(index, si)">
                            <span class="text">{{ sub.menu_name || t('label_sub_menu') }}</span>
                            <svg-icon class="del" name="delete-line" @click.stop="removeSubMenu(index, si)"
                              style="font-size: 18px;" />
                          </div>
                        </template>
                      </Draggable>
                      <div class="submenu-add" v-if="(element.sub_menu_list || []).length < 5"
                        @click="addSubMenu(index)">
                        <PlusOutlined />
                      </div>
                    </div>
                    <div class="submenu-arrow"></div>
                  </div>
                </div>
              </template>
            </Draggable>
            <div v-if="menus.length < 3" class="root-menu-add" @click="addRootMenu">
              <PlusOutlined />
            </div>
          </div>
        </div>
      </div>
      <div class="main-right">
        <div v-if="editing.type === 'root'" class="editor">
          <div class="form-item">
            <div class="label"><span class="required">*</span>{{ t('label_root_menu') }}</div>
            <a-input ref="rootNameInputRef" v-model:value="menus[activeRootIndex].menu_name" :placeholder="t('ph_input_name')" @input="onRootNameInput" maxLength="5" />
            <div class="emoji-row">
              <a-popover v-model:open="showEmoji" placement="bottomLeft" trigger="click" :getPopupContainer="getPopup">
                <template #content>
                  <Picker :data="emojiIndex" :emojiSize="18" :showPreview="false" set="apple" @select="onEmojiSelect" />
                </template>
                <a-tooltip :title="t('btn_insert_emoji')">😊</a-tooltip>
              </a-popover>
            </div>
            <div class="tip">{{ t('tip_emoji_support') }}</div>
          </div>
          <template v-if="(menus[activeRootIndex].sub_menu_list || []).length === 0">
            <div class="form-item">
              <div class="label">{{ t('label_menu_function') }}</div>
              <MenuActEditor ref="actEditorRef" :key="'root-'+activeRootIndex" v-model:value="menus[activeRootIndex].act" />
            </div>
          </template>
        </div>

        <div v-else-if="editing.type === 'sub'" class="editor">
          <div class="form-item">
            <div class="label"><span class="required">*</span>{{ t('label_sub_menu') }}</div>
            <a-input ref="subNameInputRef" v-model:value="menus[activeRootIndex].sub_menu_list[activeSubIndex].menu_name" :placeholder="t('ph_input_name')"
              @input="onSubNameInput" maxLength="18" />
            <div class="emoji-row">
              <a-popover v-model:open="showSubEmoji" placement="bottomLeft" trigger="click" :getPopupContainer="getPopup">
                <template #content>
                  <Picker :data="emojiIndex" :emojiSize="18" :showPreview="false" set="apple" @select="onSubEmojiSelect" />
                </template>
                <a-tooltip :title="t('btn_insert_emoji')">😊</a-tooltip>
              </a-popover>
            </div>
            <div class="tip">{{ t('tip_emoji_support') }}</div>
          </div>
          <div class="form-item">
            <div class="label">{{ t('label_menu_function') }}</div>
            <MenuActEditor ref="actEditorRef" :key="'sub-'+activeRootIndex+'-'+activeSubIndex" v-model:value="menus[activeRootIndex].sub_menu_list[activeSubIndex].act" />
          </div>
        </div>
      </div>
    </div>
    <div class="footer-save">
      <a-button type="primary" @click="onSave" :loading="saveLoading">{{ t('btn_save_and_apply') }}</a-button>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { getCustomMenuList, saveCustomMenu, syncWxMenuToShow } from '@/api/explore/index.js'
import { getWechatAppList } from '@/api/robot'
import { message, Tooltip, Modal } from 'ant-design-vue'
import { getSpecifyAbilityConfig, saveUserAbility, closeWxMenu } from '@/api/explore'
import emojiDataJson from 'emoji-mart-vue-fast/data/all.json'
import 'emoji-mart-vue-fast/css/emoji-mart.css'
import { Picker, EmojiIndex } from 'emoji-mart-vue-fast/src'
import MenuActEditor from './menu-act-editor.vue'
import Draggable from 'vuedraggable'
import iphoneBg from '@/assets/img/iphone-bg.png'
import { useI18n } from '@/hooks/web/useI18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n('views.robot.robot-config.custom-menu.custom-menu')

const mpAccounts = ref([])
const selectedAppid = ref('')
const expanded = ref(false)
const mpListRef = ref(null)
const visibleCount = ref(0)
const CARD_WIDTH = 160
const GAP = 8
const MORE_BTN_WIDTH = 96
const syncLoading = ref(false)
const saveLoading = ref(false)

function calcVisibleCount () {
  const el = mpListRef.value
  if (!el) { visibleCount.value = 0; return }
  const w = el.clientWidth || 0
  const per = CARD_WIDTH + GAP
  const count = Math.floor((w - MORE_BTN_WIDTH) / per)
  visibleCount.value = Math.max(count, 0)
  calcRootItemWidth()
}

const syncMenu = async () => {
  if (!selectedAppid.value) { message.error(t('msg_select_account')); return }
  syncLoading.value = true
  syncWxMenuToShow({ appid: selectedAppid.value }).then((res) => {
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    const roots = list.filter((it) => Number(it.menu_level) === 1)
    const subs = list.filter((it) => Number(it.menu_level) === 2)
    if (roots.length) {
      menus.value = roots
        .sort((a, b) => Number(a.seq_id || 0) - Number(b.seq_id || 0))
        .map((r) => {
          const inlineSubs = Array.isArray(r.sub_menu_list) ? r.sub_menu_list : []
          const byParentSubs = subs.filter((s) => Number(s.parent_menu_id || 0) === Number(r.id || 0))
          const children = (inlineSubs.length ? inlineSubs : byParentSubs).slice().sort((a, b) => Number(a.seq_id || 0) - Number(b.seq_id || 0))
          return {
            id: Number(r.id || 0),
            admin_user_id: Number(r.admin_user_id || 0),
            appid: String(r.appid || ''),
            seq_id: Number(r.seq_id || 0),
            menu_name: r.menu_name || t('label_root_menu'),
            menu_level: 1,
            parent_menu_id: 0,
            act: toAct(r),
            sub_menu_list: children.map((sm, si) => ({
              id: Number(sm.id || 0),
              admin_user_id: Number(sm.admin_user_id || 0),
              appid: String(sm.appid || ''),
              seq_id: Number(sm.seq_id || si),
              menu_name: sm.menu_name || t('label_sub_menu'),
              menu_level: 2,
              parent_menu_id: Number(r.id || 0),
              act: toAct(sm)
            }))
          }
        })
      activeRootIndex.value = 0
      activeSubIndex.value = -1
      editing.type = 'root'
    }
    nextTick(calcRootItemWidth)
    message.success(t('msg_sync_success'))
  }).finally(() => {
    syncLoading.value = false
  })
}

const getWechatAppListFn = async () => {
  try {
    const prev = selectedAppid.value
    const res = await getWechatAppList({ app_type: 'official_account', app_name: '' })
    const list = Array.isArray(res?.data) ? res.data : []
    // 只需要account_is_verify为true的公众号
    mpAccounts.value = list
      .filter((it) => it.account_is_verify == 'true')
      .map((it) => ({
        id: it.id,
        appid: it.app_id,
        name: it.app_name,
        logo: it.app_avatar,
        custom_menu_status: String(it.custom_menu_status || '0')
      }))
    if (prev) {
      const hit = mpAccounts.value.find((it) => it.appid === prev)
      selectedAppid.value = hit ? prev : (mpAccounts.value[0]?.appid || '')
    } else {
      selectedAppid.value = mpAccounts.value[0]?.appid || ''
    }
  } catch (_e) {
    mpAccounts.value = []
    selectedAppid.value = ''
  }
}

onMounted(async () => {
  initAbilitySwitch()
  await getWechatAppListFn()
  initAppMenuSwitch()
  nextTick(calcVisibleCount)
  window.addEventListener('resize', calcVisibleCount)
  if (selectedAppid.value) {
    loadMenuList()
  }
})
onUnmounted(() => { window.removeEventListener('resize', calcVisibleCount) })

function selectMp (mp) {
  selectedAppid.value = mp.appid
  expanded.value = true
  initAppMenuSwitch()
  if (selectedAppid.value) {
    loadMenuList()
  }
}

const abilitySwitchChecked = ref(false)
const initAbilitySwitch = () => {
  getSpecifyAbilityConfig({ ability_type: 'official_custom_menu' }).then((res) => {
    const st = res?.data?.user_config?.switch_status
    abilitySwitchChecked.value = String(st || '0') === '1'
  })
}
const onAbilitySwitchChange = (checked) => {
  const newStatus = checked ? '1' : '0'
  if (newStatus === '0') {
    Modal.confirm({
      title: t('confirm_title'),
      content: t('confirm_close_ability'),
      onOk: () => {
        saveUserAbility({ ability_type: 'official_custom_menu', switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            abilitySwitchChecked.value = false
            message.success(t('msg_operation_success'))
            // 刷新公众号列表
            nextTick(async () => {
              await getWechatAppListFn()
              initAppMenuSwitch()
            })
          }
        })
      }
    })
    return
  }
  saveUserAbility({ ability_type: 'official_custom_menu', switch_status: newStatus }).then((res) => {
    if (res && res.res == 0) {
      abilitySwitchChecked.value = true
      message.success(t('msg_operation_success'))
    }
  })
}

const appMenuSwitch = ref('1')
const initAppMenuSwitch = () => {
  if (!selectedAppid.value) { appMenuSwitch.value = '0'; return }
  const cur = mpAccounts.value.find((it) => it.appid === selectedAppid.value)
  appMenuSwitch.value = String(cur?.custom_menu_status || '0')
}
const onAppSwitchChange = (checked) => {
  const switch_status = checked
  if (!selectedAppid.value) { message.error(t('msg_select_account')); return }
  if (switch_status === '0') {
    Modal.confirm({
      title: t('confirm_title'),
      content: abilitySwitchChecked.value ? t('confirm_close_app_menu') : t('confirm_close_app_menu_local'),
      onOk: async () => {
        try {
          await closeWxMenu({ appid: selectedAppid.value })
          message.success(t('msg_operation_success'))
        } finally {
          await getWechatAppListFn()
          initAppMenuSwitch()
        }
      },
      onCancel: () => {
        appMenuSwitch.value = '1'
      }
    })
    return
  }
  // 开启
  Modal.confirm({
    title: t('confirm_title'),
    content: abilitySwitchChecked.value ? t('confirm_open_app_menu') : t('confirm_close_app_menu_local'),
    onOk: async () => {
      const menu_json = toMenuJson()
      saveCustomMenu({ appid: selectedAppid.value, menu_json: JSON.stringify(menu_json) })
        .then(() => { message.success(t('msg_save_success')) })
        .finally(async () => {
          await getWechatAppListFn()
          initAppMenuSwitch()
        })
    },
    onCancel: () => {
      appMenuSwitch.value = '0'
    }
  })
}

function toAct (item) {
    const code = Number(item?.choose_act_item || 0)
    const p = item?.act_params || {}
    const ap = {}
    if (code === 1) {
      const replies = Array.isArray(p.reply_content) ? p.reply_content : []
      ap.replyList = replies.map((rc) => ({
        type: rc?.type || rc?.reply_type || 'text',
        description: rc?.description || '',
        thumb_url: rc?.thumb_url || rc?.pic || '',
        title: rc?.title || '',
        url: rc?.url || '',
        appid: rc?.appid || '',
        page_path: rc?.page_path || '',
        standbyUrl: rc?.standbyUrl || rc?.standby_url || '',
        smart_menu_id: rc?.smart_menu_id || '',
        smart_menu: rc?.smart_menu || {},
      }))
      ap.reply_num = Number(item.reply_num || p.reply_num || 0)
    } else if (code === 2) { ap.linkUrl = String(p.linkUrl || '') }
    else if (code === 3) { ap.appid = String(p.appid || ''); ap.page_path = String(p.pagepath || p.page_path || ''); ap.standbyUrl = String(p.standbyUrl || p.standby_url || '') }
    else if (code === 5) { ap.key = String(p.key || '') }
    return { choose_act_item: code, act_params: ap }
}

async function loadMenuList () {
  try {
    const res = await getCustomMenuList({ appid: selectedAppid.value })
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    const roots = list.filter((it) => Number(it.menu_level) === 1)
    const subs = list.filter((it) => Number(it.menu_level) === 2)
    menus.value = roots
      .sort((a, b) => Number(a.seq_id || 0) - Number(b.seq_id || 0))
      .map((r) => {
        const inlineSubs = Array.isArray(r.sub_menu_list) ? r.sub_menu_list : []
        const byParentSubs = subs.filter((s) => Number(s.parent_menu_id || 0) === Number(r.id || 0))
        const children = (inlineSubs.length ? inlineSubs : byParentSubs).slice().sort((a, b) => Number(a.seq_id || 0) - Number(b.seq_id || 0))
        return {
          id: Number(r.id || 0),
          admin_user_id: Number(r.admin_user_id || 0),
          appid: String(r.appid || ''),
          seq_id: Number(r.seq_id || 0),
          menu_name: r.menu_name || t('label_root_menu'),
          menu_level: 1,
          parent_menu_id: 0,
          act: toAct(r),
          sub_menu_list: children.map((sm, si) => ({
            id: Number(sm.id || 0),
            admin_user_id: Number(sm.admin_user_id || 0),
            appid: String(sm.appid || ''),
            seq_id: Number(sm.seq_id || si),
            menu_name: sm.menu_name || t('label_sub_menu'),
            menu_level: 2,
            parent_menu_id: Number(r.id || 0),
            act: toAct(sm)
          }))
        }
      })
    nextTick(calcRootItemWidth)
  } catch (_) {
    // 保持默认
  }
}

const goBack = () => {
  if (route.query.id && route.query.robot_key) {
    router.push({ path: '/robot/config/function-center', query: { id: route.query.id, robot_key: route.query.robot_key } })
  } else {
    router.push({ path: '/explore/index' })
  }
}

// ====== 自定义菜单编辑区 ======
const activeRootIndex = ref(0)
const activeSubIndex = ref(-1)
const editing = reactive({ type: 'root' })
const actEditorRef = ref(null)
const rootNameInputRef = ref(null)
const showEmoji = ref(false)
const emojiIndex = new EmojiIndex(emojiDataJson)
const subNameInputRef = ref(null)
const showSubEmoji = ref(false)
function getPopup () { return document.body }
function onEmojiSelect (emoji) {
  const char = emoji?.native || ''
  if (!char) return
  const el = rootNameInputRef.value?.$el?.querySelector('input')
  const val = String(menus.value[activeRootIndex.value]?.menu_name || '')
  if (el && typeof el.selectionStart === 'number') {
    const start = el.selectionStart
    const end = el.selectionEnd
    const nextVal = val.slice(0, start) + char + val.slice(end)
    const clamped = clampNameLen(nextVal, 5)
    if (clamped !== nextVal) { message.warning(t('msg_root_menu_max_5_chars')) }
    menus.value[activeRootIndex.value].menu_name = clamped
    nextTick(() => {
      el.focus()
      const pos = (menus.value[activeRootIndex.value].menu_name || '').length
      el.setSelectionRange(pos, pos)
    })
  } else {
    const nextVal = val + char
    const clamped = clampNameLen(nextVal, 5)
    if (clamped !== nextVal) { message.warning(t('msg_root_menu_max_5_chars')) }
    menus.value[activeRootIndex.value].menu_name = clamped
  }
  showEmoji.value = false
}

function onSubEmojiSelect (emoji) {
  const char = emoji?.native || ''
  if (!char) return
  const el = subNameInputRef.value?.$el?.querySelector('input')
  const val = String(menus.value[activeRootIndex.value]?.sub_menu_list?.[activeSubIndex.value]?.menu_name || '')
  const max = 18
  if (el && typeof el.selectionStart === 'number') {
    const start = el.selectionStart
    const end = el.selectionEnd
    const nextVal = val.slice(0, start) + char + val.slice(end)
    const clamped = clampNameLen(nextVal, max)
    if (clamped !== nextVal) { message.warning(t('msg_sub_menu_max_18_chars')) }
    menus.value[activeRootIndex.value].sub_menu_list[activeSubIndex.value].menu_name = clamped
    nextTick(() => {
      el.focus()
      const pos = (menus.value[activeRootIndex.value].sub_menu_list[activeSubIndex.value].menu_name || '').length
      el.setSelectionRange(pos, pos)
    })
  } else {
    const nextVal = val + char
    const clamped = clampNameLen(nextVal, max)
    if (clamped !== nextVal) { message.warning(t('msg_sub_menu_max_18_chars')) }
    menus.value[activeRootIndex.value].sub_menu_list[activeSubIndex.value].menu_name = clamped
  }
  showSubEmoji.value = false
}

function onSubNameInput () {
  const v = String(menus.value[activeRootIndex.value]?.sub_menu_list?.[activeSubIndex.value]?.menu_name || '')
  const cl = clampNameLen(v, 18)
  if (cl !== v) {
    menus.value[activeRootIndex.value].sub_menu_list[activeSubIndex.value].menu_name = cl
    message.warning(t('msg_sub_menu_max_18_chars'))
  }
}

const segmentGraphemes = (s) => {
  try {
    const seg = new Intl.Segmenter('zh', { granularity: 'grapheme' })
    return Array.from(seg.segment(String(s))).map((it) => it.segment)
  } catch (_) {
    return Array.from(String(s))
  }
}
const clampNameLen = (s, n = 5) => {
  const segs = segmentGraphemes(s)
  return segs.slice(0, n).join('')
}
const onRootNameInput = () => {
  const v = String(menus.value[activeRootIndex.value]?.menu_name || '')
  const cl = clampNameLen(v, 5)
  if (cl !== v) {
    menus.value[activeRootIndex.value].menu_name = cl
    message.warning(t('msg_root_menu_max_5_chars'))
  }
}

function makeEmptyAct () {
  return {
    choose_act_item: 1,
    act_params: {
      replyList: [{
        type: 'text',
        description: ''
      }],
      reply_num: 0
    }
  }
}

const menus = ref([
  {
    id: 0,
    admin_user_id: 0,
    appid: '',
    seq_id: 0,
    menu_name: t('label_root_menu'),
    menu_level: 1,
    parent_menu_id: 0,
    choose_act_item: 0,
    act_params: {},
    act: makeEmptyAct(),
    sub_menu_list: []
  }
])

const bottomMenuRef = ref(null)
const rootItemWidth = ref(64)
function calcRootItemWidth () {
  const wrap = bottomMenuRef.value
  if (!wrap) return
  const total = wrap.clientWidth || 0
  const hasAdd = menus.value.length < 3
  const addWidth = hasAdd ? 32 : 0
  const gap = 8
  const n = Math.max(menus.value.length, 1)
  const avail = total - addWidth - gap * (n - 1)
  const w = Math.max(Math.floor(avail / n), 48)
  rootItemWidth.value = w
}

function addRootMenu () {
  if (menus.value.length >= 3) return
  menus.value.push({
    id: 0,
    admin_user_id: 0,
    appid: selectedAppid.value,
    seq_id: 0,
    menu_name: t('label_root_menu'),
    menu_level: 1,
    parent_menu_id: 0,
    choose_act_item: 0,
    act_params: {},
    act: makeEmptyAct(),
    sub_menu_list: []
  })
  activeRootIndex.value = menus.value.length - 1
  activeSubIndex.value = -1
  editing.type = 'root'
  nextTick(calcRootItemWidth)
}
function onSelectRoot (idx) {
  activeRootIndex.value = idx;
  activeSubIndex.value = -1;
  editing.type = 'root'
}
function removeRootMenu (idx) {
  if (menus.value.length <= 1) {
    message.warning(t('msg_keep_at_least_one_root_menu'));
    return
  }
  Modal.confirm({
    title: t('confirm_title'),
    content: t('confirm_delete_root_menu'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk() {
      menus.value.splice(idx, 1)
      if (activeRootIndex.value >= menus.value.length) {
        activeRootIndex.value = menus.value.length - 1
      } else if (activeRootIndex.value > idx) {
        activeRootIndex.value = activeRootIndex.value - 1
      }
      activeSubIndex.value = -1
      editing.type = 'root'
      nextTick(calcRootItemWidth)
    }
  })
}
function addSubMenu (rootIdx) {
  const root = menus.value[rootIdx]
  if (!root) return
  root.sub_menu_list = root.sub_menu_list || []
  if (root.sub_menu_list.length >= 5) return
  const doAdd = () => {
    root.sub_menu_list.push({
      id: 0,
      admin_user_id: 0,
      appid: selectedAppid.value,
      seq_id: root.sub_menu_list.length,
      menu_name: t('label_sub_menu'),
      menu_level: 2,
      parent_menu_id: root.id,
      choose_act_item: 0,
      act_params: {},
      act: makeEmptyAct()
    })
    activeRootIndex.value = rootIdx
    activeSubIndex.value = (root.sub_menu_list || []).length - 1
    editing.type = 'sub'
  }
  if (root.sub_menu_list.length === 0) {
    Modal.confirm({
      title: t('confirm_title'),
      content: t('confirm_add_sub_menu'),
      okText: t('btn_confirm'),
      cancelText: t('btn_cancel'),
      onOk() {
        root.act = makeEmptyAct()
        doAdd()
      }
    })
  } else {
    doAdd()
  }
}
function onSelectSub (rootIdx, subIdx) {
  activeRootIndex.value = rootIdx;
  activeSubIndex.value = subIdx;
  editing.type = 'sub'
}
function removeSubMenu (rootIdx, subIdx) {
  const root = menus.value[rootIdx]
  if (!root) return
  Modal.confirm({
    title: t('confirm_title'),
    content: t('confirm_delete_sub_menu'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk() {
      root.sub_menu_list.splice(subIdx, 1)
      const len = (root.sub_menu_list || []).length
      if (activeRootIndex.value === rootIdx) {
        if (len === 0) {
          activeSubIndex.value = -1
          editing.type = 'root'
        } else if (activeSubIndex.value >= len) {
          activeSubIndex.value = len - 1
          editing.type = 'sub'
        }
      }
    }
  })
}
let draggingRootSelected = null
function onRootDragStart () {
  draggingRootSelected = menus.value[activeRootIndex.value] || null
}
function onRootDragEnd () {
  menus.value.forEach((m, i) => { m.seq_id = i })
  if (draggingRootSelected) {
    const newIdx = menus.value.findIndex((m) => m === draggingRootSelected)
    if (newIdx >= 0) {
      activeRootIndex.value = newIdx
    }
  }
  draggingRootSelected = null
}
let draggingSubSelected = null
function onSubDragStart (rootIdx) {
  const root = menus.value[rootIdx]
  if (!root) { draggingSubSelected = null; return }
  if (activeRootIndex.value === rootIdx && activeSubIndex.value >= 0) {
    draggingSubSelected = (root.sub_menu_list || [])[activeSubIndex.value] || null
  } else {
    draggingSubSelected = null
  }
}
function onSubDragEnd (rootIdx) {
  const root = menus.value[rootIdx]
  if (!root) return
  (root.sub_menu_list || []).forEach((sm, i) => { sm.seq_id = i })
  if (activeRootIndex.value === rootIdx && draggingSubSelected) {
    const newIdx = (root.sub_menu_list || []).findIndex((sm) => sm === draggingSubSelected)
    if (newIdx >= 0) {
      activeSubIndex.value = newIdx
    }
    editing.type = (activeSubIndex.value >= 0) ? 'sub' : 'root'
  }
  draggingSubSelected = null
}

function serializeReplyContent (list) {
  return (list || []).map((it) => ({
    ...it,
    status: '1'
  }))
}
function serializeReplyTypeCodes (list) {
  const map = {
    text: '2',
    image: '4',
    card: '3',
    imageText: '1',
    url: '5',
    smartMenu: '6'
  }; return (list || []).map((it) => map[it.type] || '').filter(Boolean)
}

function buildActPayload (act) {
  const item = Number(act?.choose_act_item || 0)
  const p = act?.act_params || {}
  if (item === 1) {
    return {
      item,
      reply_content: serializeReplyContent(p.replyList || []),
      reply_type: serializeReplyTypeCodes(p.replyList || []),
      reply_num: Number(p.reply_num || 0)
    }
  }
  if (item === 2) {
    return {
      item,
      linkUrl: String(p.linkUrl || '')
    }
  }
  if (item === 3) {
    return {
      item,
      appid: String(p.appid || ''),
      pagepath: String(p.page_path || ''),
      standbyUrl: String(p.standbyUrl || '')
    }
  }
  if (item === 5) {
    return {
      item,
      key: String(p.key || '')
    }
  }
  return {
    item: 0
  }
}

function httpOk (s) {
  return /^https?:\/\//.test(String(s || ''))
}
function validateReplyList (list) {
  const arr = Array.isArray(list) ? list : []
  if (!arr.length) return false
  for (const rc of arr) {
    const t = rc?.type || rc?.reply_type
    if (t === 'text') {
      if (!String(rc?.description || '').trim()) return false
    } else if (t === 'image') {
      if (!String(rc?.thumb_url || '').trim()) return false
    } else if (t === 'imageText') {
      if (!httpOk(rc?.url) || !String(rc?.title || '').trim() || !String(rc?.description || '').trim() || !String(rc?.thumb_url || '').trim()) return false
    } else if (t === 'url') {
      if (!httpOk(rc?.url) || !String(rc?.title || '').trim()) return false
    } else if (t === 'card') {
      if (!String(rc?.title || '').trim() || !String(rc?.appid || '').trim() || !String(rc?.page_path || '').trim() || !String(rc?.thumb_url || '').trim()) return false
    } else if (t === 'smartMenu') {
      const sid = String(rc?.smart_menu_id || rc?.smart_menu?.id || '').trim();
      if (!sid) return false
    }
  }
  return true
}
function validateActDeep (act) {
  const item = Number(act?.choose_act_item || 0)
  const p = act?.act_params || {}
  if (item === 0) return true
  if (item === 1) return validateReplyList(p.replyList)
  if (item === 2) return httpOk(p.linkUrl)
  if (item === 3) return String(p.appid || '').trim().length > 0 && String(p.page_path || '').trim().length > 0 && httpOk(p.standbyUrl)
  if (item === 5) return String(p.key || '').trim().length > 0
  return true
}

function toMenuJson () {
  return menus.value.map((root, ri) => {
    const node = {
      admin_user_id: 1,
      appid: selectedAppid.value,
      seq_id: ri,
      menu_name: root.menu_name || t('label_root_menu'),
      menu_level: 1,
      parent_menu_id: 0,
      choose_act_item: (root.sub_menu_list || []).length > 0 ? 0 : Number(root.act?.choose_act_item || 0),
      act_params: (root.sub_menu_list || []).length > 0 ? { item: 0 } : buildActPayload(root.act),
      oper_user_id: 1,
      sub_menu_list: (root.sub_menu_list || []).map((sm, si) => {
        const child = {
          admin_user_id: 1,
          appid: selectedAppid.value,
          seq_id: si,
          menu_name: sm.menu_name || t('label_sub_menu'),
          menu_level: 2,
          parent_menu_id: Number(root.id || 0),
          choose_act_item: Number(sm.act?.choose_act_item || 0),
          act_params: buildActPayload(sm.act),
          oper_user_id: 1,
          create_time: toolTime(),
          update_time: toolTime()
        }
        const sid = Number(sm.id || 0)
        if (sid > 0) child.id = sid
        return child
      }),
      create_time: toolTime(),
      update_time: toolTime()
    }
    const rid = Number(root.id || 0)
    if (rid > 0) node.id = rid
    return node
  })
}

function toolTime () {
  return Math.floor(Date.now() / 1000)
}

async function onSave () {
  if (saveLoading.value) {
    return
  }
  if (!selectedAppid.value) {
    message.error(t('msg_select_account'));
    return
  }
  if (actEditorRef.value && actEditorRef.value.getValue) {
    const val = actEditorRef.value.getValue()
    if (editing.type === 'root') {
      menus.value[activeRootIndex.value].act = val
    } else if (editing.type === 'sub') {
      const root = menus.value[activeRootIndex.value]
      if (root && root.sub_menu_list && root.sub_menu_list[activeSubIndex.value]) {
        root.sub_menu_list[activeSubIndex.value].act = val
      }
    }
  }
  if (actEditorRef.value && actEditorRef.value.validate) {
    const ok = await actEditorRef.value.validate()
    if (!ok) { return }
  }
  // 简单校验
  for (let ri = 0; ri < menus.value.length; ri++) {
    const root = menus.value[ri]
    const rootName = String(root.menu_name || t('label_root_menu'))
    if ((root.sub_menu_list || []).length === 0) {
      if (!validateActDeep(root.act)) { message.error(t('msg_complete_root_menu_config', { name: rootName })); return }
    }
    for (let si = 0; si < (root.sub_menu_list || []).length; si++) {
      const sm = root.sub_menu_list[si]
      const subName = String(sm.menu_name || t('label_sub_menu'))
      if (!validateActDeep(sm.act)) { message.error(t('msg_complete_sub_menu_config', { rootName, subName })); return }
    }
  }
  const menu_json = toMenuJson()
  Modal.confirm({
    title: t('confirm_save_and_apply'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: async () => {
      saveLoading.value = true
      try {
        await saveCustomMenu({
          appid: selectedAppid.value,
          menu_json: JSON.stringify(menu_json)
        })
        message.success(t('msg_save_success'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

// 动态宽度监听
watch(() => menus.value.length, () => {
  nextTick(calcRootItemWidth)
})
watch(() => activeRootIndex.value, () => {
  nextTick(calcRootItemWidth)
})
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;

  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding-bottom: 16px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 400;
    line-height: 24px;
  }

  .search-block {
    display: flex;
    align-items: center;
    gap: 16px;
    margin: 24px 0 12px;
  }

  .list-box {
    margin-top: 8px;
  }

  ::v-deep(.ant-alert) {
    align-items: baseline;
  }
}

.flex {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.mp-list-block {
  margin: 16px 0 8px 0;
}

.mp-list {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
}

.mp-list.expanded {
  flex-wrap: wrap;
}

.mp-card {
  cursor: pointer;
  min-width: 160px;
  padding: 8px 12px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #edeff2;
  display: inline-flex;
  align-items: center;
  gap: 8px;

  &:hover {
    box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.08);
  }
}

.selected {
  border-color: #1890ff;
  background-color: rgba(24, 144, 255, 0.04);
}

.mp-logo {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.mp-name {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
}

.more-btn {
  flex: 0 0 auto;
}

.breadcrumb-wrap {
  width: fit-content;
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-bottom: 16px;
}

.breadcrumb-title {
  margin: 0 12px 0 2px;
  color: #262626;
  font-size: 16px;
  font-style: normal;
  font-weight: 600;
  line-height: 24px;
}

.switch-tip {
  margin-left: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.main {
  display: flex;
  gap: 80px;
}

.main-left {
  flex: 0 0 auto;
}

.main-right {
  max-width: 694px;
  flex: 1;
}

.iphone-mock {
  width: 264px;
  height: 544px;
  background-repeat: no-repeat;
  background-size: contain;
  position: relative;
}

.bottom-menu {
  position: absolute;
  height: 32px;
  width: calc(100% - 60px);
  bottom: 42px;
  right: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.root-menu-draggable {
  display: flex;
  flex: 1;
}

.root-menu-item {
  box-sizing: border-box;
  line-height: 32px;
  height: 32px;
  flex: 1;
  text-align: center;
  cursor: pointer;
  position: relative;
  border: 1px solid transparent;
  font-size: 12px;
  color: #262626;
  border-right: 1px solid #D9D9D9;

  .name {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.del-root {
  position: absolute;
  top: -2px;
  right: -2px;
  cursor: pointer;
  opacity: 0;
  transition: opacity .2s ease;
}

.root-menu-item:hover .del-root {
  opacity: 1;
}

.root-menu-add {
  width: 32px;
  height: 32px;
  border: 1px solid transparent;
  color: #262626;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.submenu-panel {
  position: absolute;
  bottom: 42px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.submenu-card {
  width: 64px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.sub-menu-item {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0px 2px;
  border: 1px solid transparent;
  cursor: pointer;
  position: relative;
  border-bottom: 1px solid #F0F0F0;

  .text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .del {
    position: absolute;
    top: -2px;
    right: -2px;
  }
}

.sub-menu-item .del {
  color: #595959;
  display: none;
  font-size: 18px;
  cursor: pointer;
}

.sub-menu-item:hover .del {
  display: inline;
}

.root-menu-item.active,
.sub-menu-item.active {
  border-color: #2475FC;
}

.submenu-add {
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
  border-top: 1px solid #f0f0f0;
  cursor: pointer;
  color: #262626;
  font-size: 14px;
}

.submenu-arrow {
  width: 0;
  height: 0;
  border-left: 15px solid transparent;
  border-right: 15px solid transparent;
  border-top: 8px solid #fff;
  filter: drop-shadow(0 2px 0 #edeff2);
}

.form-item {
  margin-bottom: 24px;
}

.form-item .label {
  color: #262626;
  font-weight: 400;
  margin-bottom: 6px;
}

.form-item .tip {
  color: #8c8c8c;
  font-size: 12px;
  margin-top: 4px;
}

.required {
  color: #FB363F;
  margin-right: 4px;
}

.emoji-row {
  width: 100%;
  display: flex;
  align-items: flex-end;
  cursor: pointer;
  gap: 8px;
  margin-top: 4px;
}

.footer-save {
  position: fixed;
  bottom: 0;
  right: 16px;
  display: flex;
  width: 100%;
  padding: 16px 1055px 16px 406px;
  align-items: center;
  border-radius: 0 0 2px 2px;
  background: #FFF;
  box-shadow: 0 -8px 4px 0 #0000000a;
}
</style>
