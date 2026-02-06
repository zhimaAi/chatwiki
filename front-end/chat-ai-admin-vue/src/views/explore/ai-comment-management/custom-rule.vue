<template>
  <div class="custom-rule-page">
    <div class="toolbar">
      <a-button type="primary" @click="goCreate">{{ t('btn_add_rule') }}</a-button>
      <a-input-search
        v-model:value="searchKey"
        allowClear
        :placeholder="t('search_rule_name_placeholder')"
        style="width: 300px"
        @search="search"
        @pressEnter="search"
      />
    </div>

    <div class="table-box">
      <a-table
        :columns="columns"
        :data-source="ruleList"
        :loading="loading"
        row-key="id"
        :pagination="{ current: pager.page, pageSize: pager.size, total: pager.total, showSizeChanger: true, pageSizeOptions: ['10','20','50','100'] }"
        @change="onTableChange"
      >
        <template #headerCell="{ column }">
          <span v-if="typeof column.title === 'string'">{{ column.title }}</span>
        </template>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key==='delete_comment_switch'">
            <a-switch :checked="Number(record.delete_comment_switch) === 1" @change="(checked) => onToggle(record, 'delete_comment_switch', checked)" />
          </template>
          <template v-else-if="column.key==='reply_comment_switch'">
            <a-switch :checked="Number(record.reply_comment_switch) === 1" @change="(checked) => onToggle(record, 'reply_comment_switch', checked)" />
          </template>
          <template v-else-if="column.key==='elect_comment_switch'">
            <a-switch :checked="Number(record.elect_comment_switch) === 1" @change="(checked) => onToggle(record, 'elect_comment_switch', checked)" />
          </template>
          <template v-else-if="column.key==='switch'">
            <a-switch :checked="Number(record.switch) === 1" @change="(checked) => onToggle(record, 'switch', checked)" />
          </template>
          <template v-else-if="column.key==='task_info'">
            <span v-if="(record.task_info || []).length === 0">-</span>
            <template v-else>
              <span>{{ (record.task_info || []).slice(0, 2).map(it => it.task_name).join(', ') }}</span>
              <a-popover v-if="(record.task_info || []).length > 2" trigger="hover" placement="top">
                <template #content>
                  <div style="max-width: 360px; white-space: normal;">
                    {{ (record.task_info || []).map(it => it.task_name).join(', ') }}
                  </div>
                </template>
                <a style="margin-left: 4px">+{{ (record.task_info || []).length - 2 }}</a>
              </a-popover>
            </template>
          </template>
          <template v-else-if="column.key==='create_time'">
            <span>{{ formatDate(record.create_time) }}</span>
          </template>
          <template v-else-if="column.key==='actions'">
            <a-space>
              <a @click="editRule(record)">{{ t('action_edit') }}</a>
              <span v-if="record.is_default == '1'" class="actions-disabled">{{ t('action_delete') }}</span>
              <a v-else @click="delRule(record)">{{ t('action_delete') }}</a>
              <a @click="copyRule(record)">{{ t('action_copy') }}</a>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { createVNode, ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getCommentRuleList, deleteCommentRule, changeCommentRuleStatus } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'

const router = useRouter()
const { t } = useI18n('views.explore.ai-comment-management.custom-rule')

const searchKey = ref('')
const ruleList = ref([])
const loading = ref(false)
const pager = reactive({ page: 1, size: 10, total: 0 })

const columns = [
  { title: t('column_rule_name'), dataIndex: 'rule_name', key: 'rule_name' },
  { title: t('column_delete_comment'), dataIndex: 'delete_comment_switch', key: 'delete_comment_switch' },
  { title: t('column_reply_comment'), dataIndex: 'reply_comment_switch', key: 'reply_comment_switch' },
  { title: t('column_elect_comment'), dataIndex: 'elect_comment_switch', key: 'elect_comment_switch' },
  { title: t('column_task_info'), dataIndex: 'task_info', key: 'task_info' },
  { title: t('column_switch'), dataIndex: 'switch', key: 'switch' },
  { title: t('column_create_time'), dataIndex: 'create_time', key: 'create_time' },
  { title: t('column_actions'), dataIndex: 'actions', key: 'actions' },
]

const loadList = () => {
  loading.value = true
  getCommentRuleList({ page: pager.page, size: pager.size, rule_name: searchKey.value, is_default: 0 }).then((res) => {
    const data = res?.data || { list: [], total: 0 }
    ruleList.value = data.list || []
    pager.total = +data.total || 0
  }).finally(() => { loading.value = false })
}

const search = () => { pager.page = 1; loadList() }
const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  loadList()
}

const switchMap = {
  delete_comment_switch: t('switch_delete_comment'),
  reply_comment_switch: t('switch_reply_comment'),
  elect_comment_switch: t('switch_elect_comment'),
  switch: t('switch_rule_enable'),
}

const onToggle = async (record, field, checked) => {
  if (!checked) {
    Modal.confirm({
      title: t('confirm_modal_title'),
      content: t('confirm_close_switch', { rule_name: record.rule_name, switch_name: switchMap[field] }),
      icon: createVNode(ExclamationCircleOutlined),
      onOk: async () => {
        const val = checked ? 1 : 0
        await changeCommentRuleStatus({ id: record.id, change_fields: field, switch_status: val })
        record[field] = String(val)
        message.success(t('message_operation_success'))
      }
    })
    return
  }
  const val = checked ? 1 : 0
  await changeCommentRuleStatus({ id: record.id, change_fields: field, switch_status: val })
  record[field] = String(val)
  message.success(t('message_operation_success'))
}

const editRule = (record) => {
  router.push({ path: '/explore/index/ai-comment-management/create-custom-rule', query: { id: record.id } })
}
const copyRule = (record) => {
  router.push({ path: '/explore/index/ai-comment-management/create-custom-rule', query: { id: record.id, copy: 1 } })
}
const delRule = (record) => {
  Modal.confirm({
    title: t('confirm_delete_rule_title'),
    icon: createVNode(ExclamationCircleOutlined),
    onOk: async () => {
      await deleteCommentRule({ id: record.id })
      message.success(t('message_delete_success'))
      loadList()
    },
    content: t('confirm_delete_rule_content'),
  })
}

const goCreate = () => { router.push({ path: '/explore/index/ai-comment-management/create-custom-rule' }) }

const formatDate = (ts) => {
  try {
    const d = new Date(Number(ts) * 1000)
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${dd}`
  } catch (e) { return '' }
}

onMounted(() => { loadList() })
</script>

<style lang="less" scoped>
.custom-rule-page {
  ::v-deep .ant-table-thead >tr>th {
    color: #262626;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }

  ::v-deep .ant-table-row {
    color: #595959;
  }
}
.toolbar { display: flex; align-items: center; gap: 12px; }
.table-box { margin-top: 16px; }
.actions-disabled {
  color: #8c8c8c;
  cursor: not-allowed;
}
</style>
