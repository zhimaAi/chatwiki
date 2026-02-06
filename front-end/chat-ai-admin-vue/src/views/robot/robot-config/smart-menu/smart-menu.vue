<template>
  <div class="user-model-page">
    <!-- <div class="page-title">智能菜单</div> -->
    <div class="switch-block">
      <span class="switch-title">{{ t('title_smart_menu') }}</span>
      <a-switch @change="smartMenuSwitchChange" :checked="smartMenuStatus" :checked-children="t('switch_on')"
        :un-checked-children="t('switch_off')" />
      <span class="switch-desc">
        {{ t('switch_desc') }}
      </span>
    </div>
    <a-alert>
      <template #message>
        <p>{{ t('alert_description') }}</p>
        <p>{{ t('alert_usage') }}</p>
      </template>
    </a-alert>
    <div class="search-block">
      <div class="left-block">
        <a-button type="primary" @click="handleAddReply">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('btn_add_smart_menu') }}
        </a-button>
      </div>
    </div>
    <div class="list-box">
      <!-- 空状态 -->
      <div v-if="smartList.length === 0" class="list-empty">
        <ListEmpty size="200" :text="t('text_no_data')" />
      </div>
      <div v-else class="list-grid">
        <div class="list-card" v-for="it in smartList" :key="it.id">
          <div class="card-header">
            <a-avatar :src="it.avatar_url || defaultAvatar" :size="40" shape="square" class="sm-avatar" />
            <div class="card-content">
              <div class="card-title">{{ it.menu_description }}</div>
              <div class="card-text">
                <template v-for="(line, li) in getDisplayLines(it)" :key="li">
                  <div class="reply-line" @click="onReplyLineClick($event)">
                    <span class="line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
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
          <div class="card-footer">
            <div class="footer-left">
              <div class="menu-name">{{ it.menu_title }}</div>
              <div class="update-time">{{ t('text_updated_on') }}{{ formatDateFn(it.update_time, 'YYYY-MM-DD') }}</div>
            </div>
            <a-dropdown placement="bottomRight">
              <a class="more-btn" @click.prevent>
                <EllipsisOutlined />
              </a>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="edit" @click="handleEdit(it)">{{ t('btn_edit') }}</a-menu-item>
                  <a-menu-item key="del" @click="handleDelete(it)" style="color: #FF4D4F;">{{ t('btn_delete') }}</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { PlusOutlined, EllipsisOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { saveRobotAbilitySwitchStatus, getSmartMenuList, deleteSmartMenu } from '@/api/explore/index.js'
import { useRobotStore } from '@/stores/modules/robot'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import { DEFAULT_ROBOT_AVATAR } from '@/constants/index'
import ListEmpty from '@/views/robot/robot-config/function-center/components/list-empty.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.smart-menu.smart-menu')

const robotStore = useRobotStore()
// 来自左侧菜单的能力开关（关键词回复）
const smartMenuStatus = computed(() => robotStore.smartMenuSwitchStatus === '1')
const query = useRoute().query
const router = useRouter()
const loading = ref(false)
const smartList = ref([])
const defaultAvatar = DEFAULT_ROBOT_AVATAR
const getTableData = () => {
  const parmas = { robot_id: query.id }
  loading.value = true
  getSmartMenuList({ ...parmas })
    .then((res) => {
      const data = res?.data || { list: [] }
      smartList.value = (data.list || []).map((item) => ({
        id: Number(item.id || 0),
        menu_title: item.menu_title || '',
        update_time: Number(item.update_time || item.create_time || 0),
        avatar_url: item.avatar_url || '',
        switch_status: String(item.switch_status ?? '0'),
        priority_num: Number(item.priority_num || 0),
        reply_num: Number(item.reply_num || 0),
        reply_content: Array.isArray(item.reply_content) ? item.reply_content : [],
        menu_content: Array.isArray(item.menu_content) ? item.menu_content : [],
        menu_description: item.menu_description || ''
      }))
    })
    .catch(() => {
    })
    .finally(() => { loading.value = false })
}

function buildPreviewLines (list) {
  const lines = []
    ; (Array.isArray(list) ? list : []).forEach((rc) => {
      const t = rc?.type || rc?.reply_type
      if (t === 'text') {
        if (rc?.description) lines.push({ text: rc.description })
      } else if (t === 'imageText') {
        const txt = rc?.title || rc?.description || t('text_image_text_link')
        lines.push({ text: txt, url: rc?.url })
      } else if (t === 'url') {
        const txt = rc?.title || t('text_link')
        lines.push({ text: txt, url: rc?.url })
      } else if (t === 'image') {
        lines.push({ text: t('text_image') })
      } else if (t === 'card') {
        const txt = rc?.title || t('text_mini_program_card')
        lines.push({ text: txt })
      }
    })
  return lines.slice(0, 8)
}

function buildMenuLines (menu_content) {
  const out = []
    ; (Array.isArray(menu_content) ? menu_content : []).forEach((mc) => {
      const t = String(mc?.menu_type || '')
      const txt = String(mc?.content || '')
      if (t === '0') {
        if (txt === '') { out.push({ kind: 'newline' }) }
        else if (/<a[\s\S]*?<\/a>/.test(txt)) {
          const sanitized = /href\s*=\s*['"]\s*#\s*['"]/i.test(txt) ? txt.replace(/href\s*=\s*['"]\s*#\s*['"]/ig, 'href="javascript:;"') : txt.replace(/href=/ig, 'target="_blank" href=');
          out.push({
            kind: 'html',
            html: sanitized
          })
        }
        else {
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
  return out.slice(0, 20)
}

function getDisplayLines (item) {
  if (Array.isArray(item?.menu_content) && item.menu_content.length) return buildMenuLines(item.menu_content)
  return buildPreviewLines(item?.reply_content || [])
}

function onReplyLineClick (e) {
  const a = e.target?.closest?.('a')
  if (!a) return
  const href = String(a.getAttribute('href') || '')
  if (href === '#' || href === 'javascript:;') { e.preventDefault(); e.stopPropagation() }
}

const handleAddReply = () => {
  router.push({
    path: '/robot/ability/smart-menu/add-rule',
    query: {
      id: query.id,
      robot_key: query.robot_key
    }
  })
}

const handleEdit = (record) => {
  router.push({
    path: '/robot/ability/smart-menu/add-rule',
    query: {
      id: query.id,
      robot_key: query.robot_key,
      rule_id: record.id
    }
  })
}

const handleCopy = (record) => {
  router.push({
    path: '/robot/ability/smart-menu/add-rule',
    query: {
      id: query.id,
      robot_key: query.robot_key,
      copy_id: record.id
    }
  })
}

const handleDelete = (record) => {
  Modal.confirm({
    title: t('msg_confirm_delete'),
    okText: t('btn_confirm'),
    onOk: () => {
      deleteSmartMenu({ id: record.id, robot_id: query.id }).then((res) => {
        if (res && res.res == 0) {
          message.success(t('msg_delete_success'))
          getTableData()
        }
      })
    }
  })
}

const smartMenuSwitchChange = (checked) => {
  const switch_status = checked ? '1' : '0'
  if (!checked) {
    Modal.confirm({
      title: t('msg_confirm_close'),
      okText: t('btn_confirm'),
      cancelText: t('btn_cancel'),
      onOk() {
        saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_smart_menu', switch_status }).then((res) => {
          if (res && res.res == 0) {
            robotStore.setSmartMenuSwitchStatus(switch_status)
            message.success(t('msg_operation_success'))
            window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
          }
        })
      }
    })
    return
  }
  saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_smart_menu', switch_status }).then((res) => {
    if (res && res.res == 0) {
      robotStore.setSmartMenuSwitchStatus(switch_status)
      message.success(t('msg_operation_success'))
      window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
    }
  })
}

function formatDateFn (date, format = 'YYYY-MM-DD') {
  if (!date) return ''
  return dayjs(date * 1000).format(format)
}

onMounted(async () => {
  getTableData()
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
    font-weight: 600;
    line-height: 24px;
  }

  .search-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 16px;

    .left-block {
      display: flex;
      align-items: center;
      gap: 16px;
    }
  }

  .list-box {
    margin-top: 8px;
  }

  ::v-deep(.ant-alert) {
    align-items: baseline;
  }
}

.switch-block {
  display: flex;
  align-items: center;
  margin-bottom: 16px;

  .switch-title {
    margin-right: 12px;
    color: #262626;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
}

.switch-desc {
  margin-left: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.flex {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.more-btn {
  flex: 0 0 auto;
}

.list-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(300px, 1fr));
  gap: 16px;
  padding: 8px 4px;
}

.list-card {
  border: 1px solid #e6e8ec;
  border-radius: 12px;
  background: #fff;
  padding: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04);
  transition: box-shadow .2s ease, transform .2s ease;
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px 60px 12px 12px;
  background: #F2F4F7;
  border-radius: 12px 12px 0 0;
  height: 180px;
  box-shadow: 0 -4px 8px 0 #00000014 inset;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #d9d9d9 transparent;
}

.card-header::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.card-header::-webkit-scrollbar-track {
  background: transparent;
  border-radius: 8px;
}

.card-header::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 8px;
}

.card-header::-webkit-scrollbar-thumb:hover {
  background: #A1A7B3;
}

.card-content {
  width: 100%;
  flex: 1;
  padding: 12px;
  background: #fff;
  border-radius: 2px 12px 12px 12px;
}

.card-title {
  white-space: pre-wrap;
  align-self: stretch;
  color: #1a1a1a;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  margin-bottom: 14px;
}

.card-text {
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #3a4559;
  font-size: 14px;
  line-height: 22px;
  overflow: visible;
  .reply-line {
    line-height: 22px;
    .line-text {
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

.empty-line {
  height: 22px;
}

.card-footer {
  box-sizing: border-box;
  margin-top: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-top: 1px solid #edeff2;
}

.footer-left {
  display: flex;
  flex-direction: column;
  align-items: self-start;
}

.menu-name {
  align-self: stretch;
  color: #262626;
  font-size: 14px;
  font-style: normal;
  font-weight: 600;
  line-height: 22px;
}

.update-time {
  align-self: stretch;
  color: #8c8c8c;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: 16px;
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

@media (min-width: 1920px) {
  .list-grid {
    grid-template-columns: repeat(4, minmax(300px, 1fr));
  }
}

.sm-avatar {
  border-radius: 4px !important;
}

.list-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
}
</style>
