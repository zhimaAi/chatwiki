<template>
  <div class="subsection-box">
    <a-table
      :data-source="paragraphLists"
      row-key="id"
      :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
      :pagination="false"
      :loading="props.isLoading"
      :scroll="{x: 1200}"
      @change="tableChange"
    >
      <a-table-column key="id" data-index="id" :width="1148">
        <template #title>{{ t('title_qa_total', { total: props.total }) }}</template>
        <template #default="{ record }">
          <div class="qa-list-box" @dblclick="handleOpenEditModal(record)">
            <div class="list-item">
              <div class="list-label">{{ t('label_question') }}</div>
              <div class="list-content">
                <a-tooltip placement="top" v-if="record.status == 2">
                  <template #title>
                    <span>{{ record.errmsg }}</span>
                  </template>
                  <span>
                    <span class="status-error"><ExclamationCircleFilled /> </span>
                  </span>
                </a-tooltip>
                <span @click="handleOpenEditModal(record)" v-html="textToHighlight(record.question, props.search)" class="question-text"></span>
                <template v-if="record.similar_questions && record.similar_questions.length">
                  <a-popover placement="topLeft" :overlayInnerStyle="{ 'padding-right': '2px' }">
                    <template #content>
                      <div class="similar-question-list">
                        <div
                          class="similar-question-list-item"
                          v-for="(item, index) in record.similar_questions"
                          :key="index"
                        >
                          {{ item }}
                        </div>
                      </div>
                    </template>
                    <template #title>
                      <span>{{ t('msg_similar_questions', { length: record.similar_questions.length }) }}</span>
                    </template>
                    <a>（{{ record.similar_questions.length }}）</a>
                  </a-popover>
                </template>
              </div>
            </div>
            <div class="list-item list-item-answer">
              <div class="list-label">{{ t('label_answer') }}</div>
              <div class="list-content" v-html="textToHighlight(record.answer, props.search)"></div>
            </div>
            <div class="fragment-img" v-viewer>
              <img v-for="(item, index) in record.images" :key="index" :src="item" alt="" />
            </div>
          </div>
        </template>
      </a-table-column>
      <a-table-column
        :title="t('col_total')"
        key="total_hits"
        data-index="total_hits"
        :width="108"
        :sorter="true"
      >
      </a-table-column>
      <a-table-column
        :title="t('col_yesterday')"
        key="yesterday_hits"
        data-index="yesterday_hits"
        :width="108"
        :sorter="true"
      >
      </a-table-column>
      <a-table-column
        :title="t('col_today')"
        key="today_hits"
        data-index="today_hits"
        :width="108"
        :sorter="true"
      >
      </a-table-column>
      <a-table-column
        v-for="meta in metaData"
        :key="`meta_${meta.key}`"
        :data-index="`meta_${meta.key}`"
        :width="160"
      >
        <template #title>
          <a-tooltip :title="meta.name?.length > 8 ? meta.name : '' ">
            <div class="zm-line1">{{meta.name}}</div>
          </a-tooltip>
        </template>
        <template #default="{record}">
          <a-tooltip :title="record[`meta_${meta.key}`]?.length > 8 ? record[`meta_${meta.key}`] : '' ">
            <div class="zm-line1">{{record[`meta_${meta.key}`]}}</div>
          </a-tooltip>
        </template>
      </a-table-column>
      <a-table-column :title="t('col_action')" key="action" data-index="action" :width="120" fixed="right">
        <template #default="{ record, index }">
          <div class="right-opration" >
            <div class="hover-btn-box" @click="handleSetCategory(record, 0)" v-if="record.category_id > 0">
              <StarFilled style="color: #F59A23;" />
            </div>
            <div class="hover-btn-box" v-else @click="handleSetCategory(record, 4)">
              <StarOutlined />
            </div>
            <div class="hover-btn-box" @click.stop="handleOpenEditModal(record)">
              <EditOutlined />
            </div>
            <div class="hover-btn-box" @click.stop="hanldleDelete(record.id)">
              <DeleteOutlined />
            </div>
          </div>
        </template>
      </a-table-column>
    </a-table>
  </div>
</template>
<script setup>
import { reactive, ref, computed, createVNode } from 'vue'
import { message } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  ExclamationCircleFilled,
  EditOutlined,
  DeleteOutlined,
  StarOutlined,
  StarFilled
} from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import colorLists from '@/utils/starColors.js'
import {
  deleteParagraph,
  editParagraph,
  getCategoryList,
  updateParagraphCategory
} from '@/api/library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.qa-knowledge-document.components.subsection-box')

const emit = defineEmits([
  'handleDelParagraph',
  'handleScrollTargetPage',
  'openEditSubscription',
  'handleSort',
  'getList'
])
const props = defineProps({
  paragraphLists: {
    type: Array,
    default: () => []
  },
  total: {
    type: [Number, String],
    default: 0
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  search: {
    type: String,
    default: ''
  }
})

const metaData = computed(() => {
  return props.paragraphLists?.[0]?.meta_list || []
})

const state = reactive({
  selectedRowKeys: []
})
const onSelectChange = (selectedRowKeys) => {
  state.selectedRowKeys = selectedRowKeys
}

const resetSelect = () => {
  state.selectedRowKeys = []
}

const handleOpenEditModal = (item) => {
  emit('openEditSubscription', item)
}
const hanldleDelete = (id) => {
  Modal.confirm({
    title: t('title_hint'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete'),
    onOk() {
      return new Promise((resolve, reject) => {
        deleteParagraph({ id })
          .then((res) => {
            message.success(t('msg_delete_success'))
            emit('handleDelParagraph', id)
            resolve()
          })
          .catch(() => {
            resolve()
          })
      })
    },
    onCancel() {}
  })
}

const handleSetCategory = (item, category_id) => {
  updateParagraphCategory({
    id: item.id,
    category_id
  }).then((res) => {
    message.success(t('msg_set_success'))
    emit('getList')
  })
}

const tableChange = (a, b, sort) => {
  let sort_field = sort.field
  let sort_type = sort.order
  if (sort_type) {
    sort_type = sort_type === 'ascend' ? 'asc' : 'desc'
  } else {
    sort_field = ''
  }
  emit('handleSort', {
    sort_field,
    sort_type
  })
}

function textToHighlight(fullText, highlightText, options = {}) {
  if (!highlightText || !fullText) return fullText;

  const {
    highlightClass = 'highlight',
    caseSensitive = false,
    wholeWord = false
  } = options;

  const flags = caseSensitive ? 'g' : 'gi';
  let regexPattern;

  if (wholeWord) {
    // 使用单词边界匹配完整单词
    regexPattern = new RegExp(`\\b${escapeRegExp(highlightText)}\\b`, flags);
  } else {
    regexPattern = new RegExp(escapeRegExp(highlightText), flags);
  }

  return fullText.replace(regexPattern, match =>
    `<span class="${highlightClass}">${match}</span>`
  );
}

/**
 * 转义正则表达式特殊字符
 * @param {string} string
 * @returns {string}
 */
function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}


defineExpose({ handleOpenEditModal, state, resetSelect })
</script>
<style lang="less" scoped>
.subsection-box {
  width: 100%;
  .subsection-box-title {
    display: flex;
    align-items: center;
    font-size: 14px;
    line-height: 22px;
    font-weight: 600;
    color: #242933;
    span {
      color: #7a8699;
      font-weight: 400;
      margin-left: 8px;
    }
  }
}

.qa-list-box {
  .question-text {
    cursor: pointer;
    &:hover {
      color: #3475fc;
    }
  }
  .list-item {
    display: flex;
    flex-wrap: wrap;
    line-height: 22px;
    font-size: 14px;
    color: #262626;
    margin-bottom: 6px;
    .list-label {
      margin-right: 12px;
    }
    .list-content {
      flex: 1;
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
      word-break: break-all;
      &::v-deep(.highlight) {
        background-color: #FFEB3B; /* 黄色高亮 */
        padding: 0 2px;
        border-radius: 2px;
        box-shadow: 0 1px 1px rgba(0,0,0,0.1);
      }
    }
  }
  .list-item-answer {
    color: #8c8c8c;
  }
  .fragment-img {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding-left: 40px;
    img {
      width: 80px;
      height: 80px;
      border-radius: 6px;
      cursor: pointer;
    }
  }
}

.status-error {
  cursor: pointer;
  color: #fb363f;
}

.right-opration {
  display: flex;
  align-items: center;
  gap: 8px;
}
.hover-btn-box {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
  &:hover {
    background: #e4e6eb;
    border-radius: 6px;
  }
}

.similar-question-list {
  padding-right: 16px;
  width: 260px;
  max-height: 280px;
  overflow-y: auto;
  padding-top: 4px;
  .similar-question-list-item {
    margin-bottom: 8px;
    color: #000;
    font-size: 14px;
    line-height: 22px;
  }
  &::-webkit-scrollbar {
    width: 6px; /* 垂直滚动条宽度 */
    height: 6px; /* 水平滚动条高度 */
  }

  /* 滚动条轨道 */
  &::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 10px;
  }

  /* 滚动条滑块 */
  &::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 10px;
    transition: background 0.3s ease;
  }

  /* 滚动条滑块悬停状态 */
  &::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  /* 滚动条角落 */
  &::-webkit-scrollbar-corner {
    background: #f1f1f1;
  }
}
</style>
