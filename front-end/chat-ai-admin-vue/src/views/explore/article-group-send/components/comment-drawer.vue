<template>
  <a-drawer v-model:open="open" :title="t('drawer_title')" :width="420" :destroyOnClose="true">
    <div class="article-card">
      <img v-if="article.thumb_url" class="thumb" :src="article.thumb_url" />
      <img v-else class="thumb" src="@/assets/img/default-cover.png" />
      <div class="info">
        <div class="title">{{ article.title }}</div>
        <div class="meta">{{ t('article_group_label') }}{{ article.group_name }}</div>
        <div class="digest">{{ article.digest }}</div>
      </div>
    </div>
    <div class="comment-header">{{ t('comment_header', { total: pagination.total }) }}</div>
    <div class="list-box" v-if="!isLoading" ref="bodyAnchorRef">
      <div class="comment-item" v-for="it in list" :key="it.id">
        <a-avatar :size="40" :src="userAvatar" />
        <div class="comment-content">
          <div class="top-line">
            <span class="name">{{ formatUserName(it.open_id) }}</span>
            <a-tag v-if="isSelected(it)" color="#FF4D4F" style="margin-right: 0;">{{ t('tag_selected') }}</a-tag>
            <span class="time">{{ formatDisplayChatTime(it.comment_create_time) }}</span>
            <div class="delete-status" v-if="it.delete_status == '1'">{{ t('deleted_status') }}</div>
            <span class="actions" v-else>
              <DeleteOutlined class="delete-icon" @click="handleDelete(it)" />
            </span>
          </div>
          <div class="text">{{ it.content_text }}</div>
          <div class="ai-reply" v-if="it.reply_comment_text">
            <div class="top-line ai-reply-content">
              <a-avatar :size="32" :src="defaultAvatar" style="flex-basis: 32px;" />
              <div class="ai-replay-content-right">
                <div class="ai-replay-content-right-top">
                  <span class="name ai">{{ t('ai_bot_name') }}</span>
                  <span class="time">{{ formatDisplayChatTime(it.reply_create_time) }}</span>
                  <div class="delete-status" v-if="it.delete_status == '1'">{{ t('deleted_status') }}</div>
                  <span class="actions" v-else>
                    <DeleteOutlined class="delete-icon" @click="handleDeleteReply(it)" />
                  </span>
                </div>
                <div class="text">{{ it.reply_comment_text }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="empty" v-if="list.length === 0"><a-empty :description="t('empty_no_comment')" /></div>
      <div class="load-more-tips" v-else>
        <span v-if="loadingMore">{{ t('load_more_loading') }}</span>
        <span v-else-if="!hasMore">{{ t('load_more_no_more') }}</span>
      </div>
    </div>
    <div class="loading-box" v-else><a-spin /></div>
  </a-drawer>
</template>

<script setup>
import { ref, reactive, createVNode, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { DeleteOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getCommentList, deleteComment, deleteCommentReply } from '@/api/robot'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_USER_AVATAR } from '@/constants/index'
import { formatDisplayChatTime, addNoReferrerMeta, removeNoReferrerMeta } from '@/utils/index.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.explore.article-group-send.components.comment-drawer')

const defaultAvatar = DEFAULT_ROBOT_AVATAR
const userAvatar = DEFAULT_USER_AVATAR
const open = ref(false)
const article = ref({})
const isLoading = ref(false)
const pagination = reactive({ page: 1, size: 20, total: 0 })
const list = ref([])
const hasMore = ref(true)
const loadingMore = ref(false)
const bodyAnchorRef = ref(null)
let bodyEl = null

const show = (record) => {
  article.value = { ...record }
  open.value = true
  pagination.page = 1
  loadComments()
}

const loadComments = () => {
  isLoading.value = true
  getCommentList({ task_id: article.value.id, page: pagination.page, size: pagination.size })
    .then((res) => {
      const data = res?.data || {}
      list.value = Array.isArray(data.list) ? data.list : []
      pagination.page = +data.page || pagination.page
      pagination.size = +data.size || pagination.size
      pagination.total = +data.total || list.value.length || 0
      hasMore.value = list.value.length < pagination.total
    })
    .finally(() => { isLoading.value = false })
}

const handleDelete = (item) => {
  Modal.confirm({
    title: t('modal_delete_comment_title'),
    icon: createVNode(ExclamationCircleOutlined),
    okText: t('btn_delete'),
    cancelText: t('btn_cancel'),
    onOk: async () => {
      await deleteComment({ msg_id: item.msg_data_id, comment_id: item.user_comment_id, access_key: article.value.access_key })
      message.success(t('message_delete_success'))
      loadComments()
    }
  })
}

const handleDeleteReply = (item) => {
  Modal.confirm({
    title: t('modal_delete_reply_title'),
    icon: createVNode(ExclamationCircleOutlined),
    okText: t('btn_delete'),
    cancelText: t('btn_cancel'),
    onOk: async () => {
      await deleteCommentReply({ msg_id: item.msg_data_id, comment_id: item.user_comment_id, access_key: article.value.access_key })
      message.success(t('message_delete_success'))
      loadComments()
    }
  })
}

defineExpose({ show })

const formatUserName = (open_id) => {
  if (!open_id) return t('default_user_name')
  const suf = String(open_id).slice(-4)
  return t('user_name_with_suffix', { suf })
}

const isSelected = (it) => {
  if (String(it.comment_type) === '1') return true
  const text = String(it.ai_comment_result_text || '')
  return text.includes(t('auto_feature_keyword'))
}

const tryAttach = () => {
  if (bodyAnchorRef.value) {
    bodyEl = bodyAnchorRef.value.closest('.ant-drawer-body')
    if (bodyEl) bodyEl.addEventListener('scroll', onBodyScroll, { passive: true })
  }
}

const detach = () => {
  if (bodyEl) {
    bodyEl.removeEventListener('scroll', onBodyScroll)
    bodyEl = null
  }
}

const onBodyScroll = () => {
  if (!hasMore.value || loadingMore.value || !bodyEl) return
  const threshold = 80
  const nearBottom = bodyEl.scrollTop + bodyEl.clientHeight >= bodyEl.scrollHeight - threshold
  if (nearBottom) {
    loadingMore.value = true
    const nextPage = pagination.page + 1
    getCommentList({ task_id: article.value.id, page: nextPage, size: pagination.size })
      .then((res) => {
        const data = res?.data || {}
        const nextList = Array.isArray(data.list) ? data.list : []
        pagination.page = +data.page || nextPage
        pagination.size = +data.size || pagination.size
        pagination.total = +data.total || pagination.total
        list.value = list.value.concat(nextList)
        hasMore.value = list.value.length < pagination.total
      })
      .finally(() => { loadingMore.value = false })
  }
}

watch(open, (val) => {
  if (val) {
    hasMore.value = true
    loadingMore.value = false
    nextTick(tryAttach)
  } else {
    detach()
  }
})

onMounted(() => { addNoReferrerMeta() })
onUnmounted(() => { removeNoReferrerMeta() })
</script>

<style lang="less" scoped>
.article-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid #edeff2;
  border-radius: 8px;
  background: #fff;
}
.article-card .thumb {
  width: 146px;
  height: 96px;
  border-radius: 6px;
  object-fit: cover;
}
.article-card .info {
  flex: 1;
  overflow: hidden;
}
.article-card .info .source {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #8c8c8c;
}
.article-card .info .source .brand {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
  height: 20px;
  border: 1px solid #00000026;
  background: #0000000a;
  border-radius: 6px;
  font-size: 12px;
  color: #bfbfbf;
}
.article-card .info .title {
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
  color: #262626;
}
.article-card .info .meta {
  margin-top: 2px;
  font-size: 12px;
  line-height: 20px;
  color: #8c8c8c;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-card .info .digest {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  flex: 1 0 0;
  overflow: hidden;
  color: #595959;
  text-overflow: ellipsis;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.comment-header {
  margin: 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.comment-item {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.comment-content {
  flex: 1;
}

.top-line {
  display: flex;
  align-items: center;
  gap: 4px;
}

.name {
  font-size: 14px;
  color: #262626;
}

.time {
  font-size: 12px;
  color: #8c8c8c;
}

.delete-status {
  font-size: 12px;
  color: #FF4D4F;
  margin-left: auto;
}

.actions {
  margin-left: auto;
  opacity: 0;
  transition: opacity 0.2s;
  display: flex;
  padding: 4px;
  justify-content: center;
  align-items: center;
  gap: 4px;
  border-radius: 6px;

  &:hover {
    background: #E4E6EB;
  }
}

.delete-icon {
  font-size: 16px;
  color: #8c8c8c;
  cursor: pointer;
}

.comment-item:hover .actions {
  opacity: 1;
}

.text {
  margin-top: 4px;
  font-size: 14px;
  color: #595959;
  white-space: pre-wrap;
}

.ai-reply {
  margin-top: 16px;

  .ai-reply-content {
    display: flex; 
    align-items: self-start; 
    gap: 8px;

    .ai-replay-content-right {
      flex: 1;

      .ai-replay-content-right-top {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
  }
}

.loading-box {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.empty {
  padding: 24px 0;
}

.load-more-tips {
  padding: 8px 0;
  text-align: center;
  font-size: 12px;
  color: #8c8c8c;
}
</style>
