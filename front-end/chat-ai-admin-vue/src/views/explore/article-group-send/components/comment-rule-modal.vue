<template>
  <a-modal v-model:open="open" :title="t('modal_title')" :width="720" @ok="onOk" @cancel="onCancel">
    <div class="rule-modal">
      <div class="toolbar">
        <a-input-search
          v-model:value="query.rule_name"
          :placeholder="t('search_rule_placeholder')"
          @search="onSearch"
          allowClear
        />
        <a-button class="create-btn" @click="onCreate">{{ t('btn_go_create_rule') }}</a-button>
      </div>
      <div class="table-header">
        <div class="th th-name">{{ t('th_rule_name') }}</div>
        <div class="th">{{ t('th_delete_comment') }}</div>
        <div class="th">{{ t('th_reply_comment') }}</div>
        <div class="th">{{ t('th_elect_comment') }}</div>
        <div class="th">{{ t('th_rule_enable') }}</div>
      </div>
      <div v-if="loading" class="loading-box"><a-spin /></div>
        <div v-else>
          <div v-if="list.length === 0" class="empty-box"><a-empty :description="t('empty_no_rule')" /></div>
          <a-radio-group v-else v-model:value="selectedRuleId" class="radio-list">
            <div class="radio-row" v-for="it in list" :key="it.id">
              <div class="name-col">
                <a-radio :value="String(it.id)">{{ it.rule_name }}</a-radio>
                <a-tag v-if="String(it.is_default) === '1'" color="blue" class="default-tag">{{ t('tag_default') }}</a-tag>
              </div>
              <div class="col"><a-switch :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" :checked="String(it.delete_comment_switch) === '1'" disabled /></div>
              <div class="col"><a-switch :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" :checked="String(it.reply_comment_switch) === '1'" disabled /></div>
              <div class="col"><a-switch :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" :checked="String(it.elect_comment_switch) === '1'" disabled /></div>
              <div class="col"><a-switch :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" :checked="String(it.switch) === '1'" disabled /></div>
            </div>
          </a-radio-group>
          <div class="pager-box">
            <a-pagination size="small" :current="query.page" :pageSize="query.size" :total="total" @change="onPageChange" />
          </div>
        </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { getCommentRuleList, setBatchSendTaskCommentRule } from '@/api/robot'
import { useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const router = useRouter()
const { t } = useI18n('views.explore.article-group-send.components.comment-rule-modal')
const open = ref(false)
const list = ref([])
const loading = ref(false)
const total = ref(0)
const query = reactive({ page: 1, size: 10, rule_name: '', is_default: -1 })
const selectedRuleId = ref('')
const taskId = ref('')
const selectOnly = ref(false)

const show = (record) => {
  taskId.value = String(record.id || record.task_id || '')
  selectedRuleId.value = String(record.comment_rule_id || '')
  selectOnly.value = !!record.select_only || !taskId.value
  open.value = true
  query.page = 1
  loadList()
}

const loadList = () => {
  loading.value = true
  getCommentRuleList({ page: query.page, size: query.size, rule_name: query.rule_name, is_default: -1 })
    .then((res) => {
      const data = res?.data || {}
      const arr = Array.isArray(data.list) ? data.list : []
      list.value = arr.map((it) => ({
        ...it,
        id: String(it.id || ''),
        rule_name: it.rule_name || it.name || ''
      }))
      total.value = Number(data.total || arr.length || 0)
      if (!selectedRuleId.value) {
        const def = list.value.find((it) => String(it.is_default) === '1')
        if (def) selectedRuleId.value = String(def.id)
      }
    })
    .finally(() => { loading.value = false })
}

const onSearch = () => { query.page = 1; loadList() }
const onPageChange = (p) => { query.page = p; loadList() }
const onCreate = () => { 
  router.push({ path: '/explore/index/ai-comment-management/create-custom-rule' })
 }

const onOk = async () => {
  if (!selectedRuleId.value) { message.warning(t('error_select_rule')) ; return }
  if (selectOnly.value || !taskId.value) {
    const picked = (list.value || []).find((it) => String(it.id) === String(selectedRuleId.value)) || null
    const rule_name = picked?.rule_name || ''
    const is_default = String(picked?.is_default || '')
    emit('selected', { rule_id: selectedRuleId.value, rule_name, is_default })
    open.value = false
    return
  }
  await setBatchSendTaskCommentRule({ task_id: taskId.value, rule_id: selectedRuleId.value })
  message.success(t('message_set_success'))
  open.value = false
  emit('updated')
}
const onCancel = () => { open.value = false }

const emit = defineEmits(['updated', 'selected'])
defineExpose({ show })
</script>

<style lang="less" scoped>
.rule-modal { padding-top: 4px; }
.toolbar { display: flex; gap: 65px; align-items: center; }
.create-btn { flex: 0 0 auto; }
.table-header { margin-top: 12px; background: #f5f5f5; padding: 6px 12px; border-radius: 6px; color: #595959; display: flex; align-items: center; gap: 12px; }
.table-header .th-name { flex: 1; }
.table-header .th { width: 90px; }
.radio-list { display: flex; margin-top: 8px; flex-direction: column; gap: 10px; }
.radio-row { display: flex; align-items: center; gap: 12px; padding: 6px 12px; }
.radio-row .name-col { flex: 1; display: flex; align-items: center; gap: 8px; }
.radio-row .col { width: 90px; }
.default-tag { margin-left: 4px; }
.pager-box { display: flex; justify-content: flex-end; margin-top: 12px; }
.loading-box { display: flex; align-items: center; justify-content: center; padding: 24px 0; }
.empty-box { display: flex; align-items: center; justify-content: center; padding: 32px 0; }
</style>
