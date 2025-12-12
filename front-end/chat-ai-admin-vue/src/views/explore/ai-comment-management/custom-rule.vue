<template>
  <div class="custom-rule-page">
    <div class="toolbar">
      <a-button type="primary" @click="goCreate">添加规则</a-button>
      <a-input-search v-model:value="searchKey" allowClear placeholder="请输入规则名称搜索" style="width: 300px" @search="search" @pressEnter="search" />
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
              <span>{{ (record.task_info || []).slice(0, 2).map(it => it.task_name).join('，') }}</span>
              <a-popover v-if="(record.task_info || []).length > 2" trigger="hover" placement="top">
                <template #content>
                  <div style="max-width: 360px; white-space: normal;">
                    {{ (record.task_info || []).map(it => it.task_name).join('，') }}
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
              <a @click="editRule(record)">编辑</a>
              <span v-if="record.is_default == '1'" class="actions-disabled">删除</span>
              <a v-else @click="delRule(record)">删除</a>
              <a @click="copyRule(record)">复制</a>
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

const router = useRouter()

const searchKey = ref('')
const ruleList = ref([])
const loading = ref(false)
const pager = reactive({ page: 1, size: 10, total: 0 })

const columns = [
  { title: '规则名', dataIndex: 'rule_name', key: 'rule_name' },
  { title: '自动删评', dataIndex: 'delete_comment_switch', key: 'delete_comment_switch' },
  { title: '自动回复评论', dataIndex: 'reply_comment_switch', key: 'reply_comment_switch' },
  { title: '评论精选', dataIndex: 'elect_comment_switch', key: 'elect_comment_switch' },
  { title: '生效群发文章', dataIndex: 'task_info', key: 'task_info' },
  { title: '规则启用', dataIndex: 'switch', key: 'switch' },
  { title: '创建时间', dataIndex: 'create_time', key: 'create_time' },
  { title: '操作', dataIndex: 'actions', key: 'actions' },
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
  delete_comment_switch: '自动删评',
  reply_comment_switch: '自动回复评论',
  elect_comment_switch: '评论精选',
  switch: '规则启用',
}

const onToggle = async (record, field, checked) => {
  if (!checked) {
    Modal.confirm({
      title: '提示',
      content: `确认关闭${record.rule_name}中的${switchMap[field]}？`,
      icon: createVNode(ExclamationCircleOutlined),
      onOk: async () => {
        const val = checked ? 1 : 0
        await changeCommentRuleStatus({ id: record.id, change_fields: field, switch_status: val })
        record[field] = String(val)
        message.success('操作成功')
      }
    })
    return
  }
  const val = checked ? 1 : 0
  await changeCommentRuleStatus({ id: record.id, change_fields: field, switch_status: val })
  record[field] = String(val)
  message.success('操作成功')
}

const editRule = (record) => {
  router.push({ path: '/explore/index/ai-comment-management/create-custom-rule', query: { id: record.id } })
}
const copyRule = (record) => {
  router.push({ path: '/explore/index/ai-comment-management/create-custom-rule', query: { id: record.id, copy: 1 } })
}
const delRule = (record) => {
  Modal.confirm({
    title: '确认删除该规则，删除后，使用该规则的群发，会使用默认规则？',
    icon: createVNode(ExclamationCircleOutlined),
    onOk: async () => {
      await deleteCommentRule({ id: record.id })
      message.success('删除成功')
      loadList()
    }
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
