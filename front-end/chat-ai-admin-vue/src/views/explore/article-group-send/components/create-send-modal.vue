<template>
  <div>
    <a-modal class="create-send-modal" v-model:open="open" :title="mode === 'edit' ? t('modal_title_edit') : t('modal_title_create')" :width="720"
      @ok="onOk" @cancel="onCancel">
      <a-form layout="vertical" ref="sendFormRef" :model="sendForm" :rules="sendFormRules" class="create-send-form">
        <a-form-item :label="t('field_task_name')" name="task_name" required>
          <a-input v-model:value="sendForm.task_name" :placeholder="t('placeholder_task_name')" />
        </a-form-item>

        <a-form-item :label="t('field_to_user_type')" name="to_user_type">
          <a-radio-group v-model:value="sendForm.to_user_type">
            <a-radio :value="0">{{ t('to_user_all_fans') }}</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item :label="t('field_send_content')" required>
          <a-button type="dashed" style="margin-bottom: 8px" @click="openSelectDraft">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_select_draft') }}
          </a-button>
          <div class="draft-preview" v-if="selectedDraft">
            <img v-if="selectedDraft.thumb_url" class="thumb" :src="selectedDraft.thumb_url" />
            <img v-else class="thumb" src="@/assets/img/default-cover.png" />
            <div class="info">
              <div class="title">{{ selectedDraft.title }}</div>
              <div class="meta">{{ t('draft_group_label') }}{{ selectedDraft.group_name || (getGroupNameProp ?
                getGroupNameProp(selectedDraft.group_id) :
                '') }}</div>
              <div class="digest">{{ selectedDraft.digest }}</div>
            </div>
          </div>
        </a-form-item>

        <a-form-item :label="t('field_select_ai_rule')">
          <a-button type="dashed" @click="openSelectRule">{{ t('btn_select_rule') }}</a-button>
          <div class="desc" v-if="selectedRuleId">{{ t('selected_rule_desc', { name: selectedRuleName || ('ID ' + selectedRuleId) }) }}</div>
        </a-form-item>

        <a-form-item :label="t('field_comment_status')" name="comment_status">
          <a-switch v-model:checked="sendForm.comment_status" :checkedValue="1" :un-checkedValue="0" />
        </a-form-item>

        <a-form-item :label="t('field_send_time')">
          <a-radio-group v-model:value="sendTimeType">
            <a-radio value="now">{{ t('send_time_now') }}</a-radio>
            <a-radio value="custom">{{ t('send_time_custom') }}</a-radio>
          </a-radio-group>
          <div style="margin-top: 8px" v-if="sendTimeType === 'custom'">
            <a-date-picker v-model:value="sendCustomTime" show-time style="width: 100%" />
          </div>
        </a-form-item>

        <a-form-item :label="t('field_open_status')" name="open_status">
          <a-switch v-model:checked="sendForm.open_status" :checkedValue="1" :un-checkedValue="0" />
          <div class="desc">{{ t('desc_open_status') }}</div>
        </a-form-item>
      </a-form>
    </a-modal>
    <CommentRuleModal ref="commentRuleModalRef" @updated="onRuleUpdated" @selected="onRuleSelected" />
    <a-modal v-model:open="selectDraftOpen" :title="t('modal_select_draft_title')" :width="720" @ok="onConfirmSelectDraft">
      <div class="select-draft-modal">
        <div v-if="selectDraftLoading" class="loading-box"><a-spin /></div>
        <template v-else>
          <a-radio-group v-model:value="selectDraftId" style="width: 100%">
            <div class="draft-list">
              <label class="draft-item" v-for="it in selectDraftList" :key="it.id">
                <a-radio :value="it.id" />
                <img v-if="it.thumb_url" class="thumb" :src="it.thumb_url" />
                <img v-else class="thumb" src="@/assets/img/default-cover.png" />
                <div class="info">
                  <div class="title">{{ it.title }}</div>
                  <div class="meta">{{ t('draft_group_label') }}{{ getGroupNameFn(it.group_id) }}</div>
                  <div class="digest">{{ it.digest }}</div>
                </div>
              </label>
            </div>
          </a-radio-group>
          <div class="pagination-box">
            <a-pagination v-model:current="selectPaginations.page" v-model:page-size="selectPaginations.size"
              :total="selectTotal" :pageSizeOptions="['10', '20', '50']" show-size-changer @change="onSelectPageChange" />
          </div>
        </template>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { message, Modal } from 'ant-design-vue'
import { createBatchSendTask, getOfficialDraftList } from '@/api/robot'
import { PlusOutlined } from '@ant-design/icons-vue'
import CommentRuleModal from '@/views/explore/article-group-send/components/comment-rule-modal.vue'
import { addNoReferrerMeta, removeNoReferrerMeta } from '@/utils/index.js'
import { useI18n } from '@/hooks/web/useI18n'

const props = defineProps({
  appId: { type: String, default: '' },
  accessKey: { type: String, default: '' },
  getGroupName: { type: Function, default: null }
})

const router = useRouter()
const { t } = useI18n('views.explore.article-group-send.components.create-send-modal')

const open = ref(false)
const mode = ref('create')
const currentTask = ref(null)
const selectedDraft = ref(null)
const sendForm = reactive({ task_name: '', to_user_type: 0, comment_status: 1, open_status: 1 })
const sendFormRules = { task_name: [{ required: true, message: t('error_enter_task_name') }] }
const sendFormRef = ref(null)
const sendTimeType = ref('now')
const sendCustomTime = ref(null)
const getGroupNameProp = props.getGroupName
const getGroupNameFn = (gid) => {
  if (typeof getGroupNameProp === 'function') return getGroupNameProp(gid)
  if (gid == 0) return t('group_unassigned')
  return ''
}

const selectDraftOpen = ref(false)
const selectDraftList = ref([])
const selectDraftLoading = ref(false)
const selectPaginations = reactive({ page: 1, size: 10 })
const selectTotal = ref(0)
const selectDraftId = ref(undefined)
const commentRuleModalRef = ref(null)
const selectedRuleId = ref('')
const selectedRuleName = ref('')

const openSelectDraft = () => {
  selectDraftOpen.value = true
  selectDraftId.value = selectedDraft.value?.id
  loadSelectDrafts()
}

const loadSelectDrafts = async () => {
  selectDraftLoading.value = true
  const rsp = await getOfficialDraftList({
    page: selectPaginations.page,
    size: selectPaginations.size,
    app_id: props.appId,
  })
  const data = rsp?.data || {}
  selectDraftList.value = data.list || []
  selectTotal.value = data.total || 0
  selectDraftLoading.value = false
}

const onSelectPageChange = (page, size) => {
  selectPaginations.page = page
  selectPaginations.size = size
  loadSelectDrafts()
}

const onConfirmSelectDraft = () => {
  if (!selectDraftId.value) { message.error(t('error_select_draft')) ; return }
  const it = selectDraftList.value.find(d => d.id === selectDraftId.value)
  if (it) {
    selectedDraft.value = it
    const curName = String(sendForm.task_name || '').trim()
    if (!curName) {
      sendForm.task_name = it.title || ''
    }
  }
  selectDraftOpen.value = false
}

const show = (payload = {}) => {
  if (payload.task) {
    mode.value = 'edit'
    currentTask.value = payload.task
    selectedDraft.value = {
      thumb_url: payload.task.thumb_url,
      title: payload.task.title,
      digest: payload.task.digest,
      group_id: payload.task.group_id,
      group_name: payload.task.group_name,
      id: payload.task.draft_id,
    }
    sendForm.task_name = payload.task.task_name || ''
    sendForm.to_user_type = Number(payload.task.to_user_type ?? 0)
    sendForm.comment_status = Number(payload.task.comment_status || 0)
    sendForm.open_status = Number(payload.task.open_status || 0)
    selectedRuleId.value = String(payload.task.comment_rule_id || '')
    selectedRuleName.value = payload.task.comment_rule_name || ''
    const st = Number(payload.task.send_time || 0)
    if (st > 0) {
      sendTimeType.value = 'custom'
      sendCustomTime.value = dayjs(st * 1000)
    } else {
      sendTimeType.value = 'now'
      sendCustomTime.value = null
    }
    open.value = true
    return
  }
  mode.value = 'create'
  const draft = payload.draft || null
  selectedDraft.value = draft
  sendForm.task_name = draft?.title || ''
  sendForm.to_user_type = 0
  sendForm.comment_status = 1
  sendForm.open_status = 1
  sendTimeType.value = 'now'
  sendCustomTime.value = null
  selectedRuleId.value = ''
  selectedRuleName.value = ''
  open.value = true
}

const openSelectRule = () => {
  commentRuleModalRef.value && commentRuleModalRef.value.show({ select_only: true, task_id: currentTask.value?.id, comment_rule_id: selectedRuleId.value })
}

const onRuleSelected = (payload) => {
  selectedRuleId.value = String(payload?.rule_id || '')
  selectedRuleName.value = payload?.rule_name || ''
}

const onRuleUpdated = () => { emit('updated') }

const onOk = async () => {
  try {
    await sendFormRef.value?.validate()
  } catch (e) { return }
  if (!selectedDraft.value) { message.error(t('error_no_draft_selected')) ; return }
  let ts = 0
  if (sendTimeType.value === 'custom' && sendCustomTime.value) ts = Math.floor(sendCustomTime.value.valueOf() / 1000)
  const params = {
    task_name: sendForm.task_name,
    app_id: props.appId,
    access_key: props.accessKey,
    send_time: ts,
    comment_status: sendForm.comment_status,
    open_status: sendForm.open_status,
    to_user_type: sendForm.to_user_type,
  }
  if (selectedRuleId.value) params.comment_rule_id = selectedRuleId.value
  if (selectedDraft.value?.id) params.draft_id = selectedDraft.value.id
  if (mode.value === 'edit') {
    if (!currentTask.value) { message.error(t('error_invalid_task')) ; return }
    params.task_id = currentTask.value.id
    await createBatchSendTask(params)
    emit('updated')
    open.value = false
    return
  }
  await createBatchSendTask(params)

  // 如果当前页面不是文章群发页面，提示用户去群发管理查看进度
  if (router.currentRoute.value.path !== '/explore/index/article-group-send/group-send') {
    Modal.confirm({
      title: t('modal_task_created_title'),
      content: t('modal_task_created_content'),
      okText: t('modal_task_created_ok'),
      cancelText: t('btn_cancel'),
      onOk: () => {
        router.push({ path: '/explore/index/article-group-send/group-send' })
      }
    })
  } else {
    message.success(t('message_task_created'))
  }
  emit('created')
  open.value = false
}
const onCancel = () => { open.value = false }

const emit = defineEmits(['created', 'updated'])
defineExpose({ show })

onMounted(() => { addNoReferrerMeta() })
onUnmounted(() => { removeNoReferrerMeta() })
</script>

<style scoped lang="less">
.create-send-modal {
  .form-row {
    display: flex;
    margin-bottom: 16px;

    .label {
      width: 96px;
      color: #262626;
      font-weight: 600;
      line-height: 24px;
    }

    .content {
      flex: 1;
    }

    .desc {
      margin-top: 8px;
      color: #8c8c8c;
    }
  }

  .draft-preview {
    display: flex;
    gap: 10px;
    align-items: center;
    height: 132px;
    box-sizing: border-box;
    padding: 16px;
    border: 1px solid #D9D9D9;
    border-radius: 6px;

    .thumb {
      width: 146px;
      height: 96px;
      border-radius: 6px;
      object-fit: cover;
    }

    .info {
      flex: 1;
      overflow: hidden;
    }

    .info .title {
      font-weight: 600;
      font-size: 16px;
      color: #262626;
    }

    .info .meta {
      color: #8c8c8c;
      font-size: 14px;
      margin: 4px 0;
    }

    .info .digest {
      color: #595959;
      font-size: 14px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }
}

.create-send-form {
  .desc {
    margin-top: 8px;
    color: #8c8c8c;
  }
}

.select-draft-modal {
  .draft-item {
    cursor: pointer;
  }

  .pagination-box {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    margin-top: 16px;
    margin-bottom: 24px;
  }
}

.draft-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.draft-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 16px;
  border-bottom: 1px solid #D9D9D9;
}

.draft-item .thumb {
  width: 146px;
  height: 96px;
  border-radius: 4px;
  object-fit: cover;
}

.draft-item .info {
  flex: 1;
  overflow: hidden;
}

.draft-item .info .title {
  align-self: stretch;
  color: #262626;
  font-size: 16px;
  font-style: normal;
  font-weight: 600;
  line-height: 24px;
}

.draft-item .info .meta {
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
  margin: 4px 0 6px;
}

.draft-item .info .digest {
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

.draft-item .icon:hover {
  opacity: 0.7;
}

.loading-box {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
}
</style>
