<template>
  <div class="comment-record-page">
    <div class="toolbar">
      <a-select
        v-model:value="filters.check_result"
        allowClear
        :placeholder="t('filter_result_placeholder')"
        style="width: 160px"
        @change="search"
      >
        <a-select-option v-for="(label, val) in typeMap" :key="val" :value="val">{{ label }}</a-select-option>
      </a-select>
      <a-input-search
        v-model:value="filters.comment_text"
        allowClear
        :placeholder="t('search_comment_placeholder')"
        style="width: 300px"
        @search="search"
      />
    </div>

    <div class="table-box">
      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        row-key="id"
        :pagination="{ current: pager.page, pageSize: pager.size, total: pager.total, showSizeChanger: true, pageSizeOptions: ['10','20','50','100'] }"
        @change="onTableChange"
      >
        <template #headerCell="{ column }">
          <span v-if="typeof column.title === 'string'">{{ column.title }}</span>
        </template>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key==='content_text'">
            <a-tooltip :title="getTooltipTitle(record.content_text, record)" placement="top">
              <div class="ellipsis1 title" :class="`titleRef_${record.create_time}`">{{record.content_text}}</div>
            </a-tooltip>
          </template>
          <template v-else-if="column.key==='ai_comment_result_text'">
            <span>{{ record.ai_comment_result_text }}</span>
          </template>
          <template v-else-if="column.key==='reply_comment_text'">
            <a-tooltip :title="getTooltipTitle(record.reply_comment_text, record)" placement="top">
              <div class="ellipsis1 title" :class="`titleRef_${record.create_time}`">{{record.reply_comment_text}}</div>
            </a-tooltip>
          </template>
          <template v-else-if="column.key==='draft_title'">
            <a-tooltip :title="getTooltipTitle(record.draft_title, record)" placement="top">
              <div class="ellipsis1 title" :class="`titleRef_${record.create_time}`">{{record.draft_title}}</div>
            </a-tooltip>
          </template>
          <template v-else-if="column.key==='rule_name'">
            <span class="ellipsis1">{{ record.rule_name }}</span>
          </template>
          <template v-else-if="column.key==='ai_comment_rule_text'">
            <div class="rule-tags">
              <template v-for="(t, idx) in (record.ai_comment_rule_text || []).slice(0, 3)" :key="t + idx">
                <span>{{ t }}</span>
                <span v-if="idx < (record.ai_comment_rule_text || []).length - 1">, </span>
              </template>
              <a-tooltip v-if="(record.ai_comment_rule_text || []).length > 3">
                <template #title>{{ (record.ai_comment_rule_text || []).join(', ') }}</template>
                <span style="cursor: pointer;">...+{{ (record.ai_comment_rule_text || []).length - 3 }}</span>
              </a-tooltip>
            </div>
          </template>
          <template v-else-if="column.key==='ai_exec_time'">
            <span>{{ formatDateTime(record.ai_exec_time) }}</span>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getCommentList } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.explore.ai-comment-management.comment-processing-record')

const filters = reactive({ check_result: undefined, comment_text: '' })
const typeMap = reactive({})
const list = ref([])
const loading = ref(false)
const pager = reactive({ page: 1, size: 10, total: 0 })

const columns = [
  { title: t('column_content_text'), dataIndex: 'content_text', key: 'content_text', width: 220 },
  { title: t('column_result'), dataIndex: 'ai_comment_result_text', key: 'ai_comment_result_text', width: 130 },
  { title: t('column_reply_text'), dataIndex: 'reply_comment_text', key: 'reply_comment_text', width: 220 },
  { title: t('column_draft_title'), dataIndex: 'draft_title', key: 'draft_title', width: 120 },
  { title: t('column_rule_name'), dataIndex: 'rule_name', key: 'rule_name', width: 120 },
  { title: t('column_rule_detail'), dataIndex: 'ai_comment_rule_text', key: 'ai_comment_rule_text', width: 300 },
  { title: t('column_exec_time'), dataIndex: 'ai_exec_time', key: 'ai_exec_time', width: 120 },
]

const loadList = () => {
  loading.value = true
  const params = { page: pager.page, size: pager.size }
  if (filters.check_result) params.check_result = filters.check_result
  if (filters.comment_text) params.comment_text = filters.comment_text
  getCommentList(params).then((res) => {
    const data = res?.data || { list: [], total: 0, type: {} }
    list.value = data.list || []
    pager.total = +data.total || 0
    const t = data.type || {}
    Object.keys(typeMap).forEach(k => delete typeMap[k])
    Object.entries(t).forEach(([k, v]) => { typeMap[k] = v })
  }).finally(() => { loading.value = false })
}

const search = () => { pager.page = 1; loadList() }
const onTableChange = (pagination) => { pager.page = pagination.current; pager.size = pagination.pageSize; loadList() }

const formatDateTime = (ts) => {
  try {
    const d = new Date(Number(ts) * 1000)
    const yy = String(d.getFullYear()).slice(2)
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    return `${yy}-${m}-${dd} ${hh}:${mm}`
  } catch (e) { return '' }
}

// 获取 tooltip 标题
function getTooltipTitle(text, record) {
  if (!text) return null
  
  // 创建临时元素来测量文本宽度
  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')
  // 14px 根据实际字体大小修改
  context.font = '14px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
  
  const textWidth = context.measureText(text).width

  const titleRef = document.querySelector(`.titleRef_${record.create_time}`)
  if (titleRef) {
    record.title_width = titleRef.offsetWidth
  }

  const maxWidth = record?.title_width || 250
  return textWidth > maxWidth ? text : null
}

onMounted(() => { loadList() })
</script>

<style lang="less" scoped>
.comment-record-page {
  padding: 2px 0px 24px;

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
.ellipsis1 { display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; }
.rule-tags { display: flex; flex-wrap: wrap; gap: 4px; }
</style>
