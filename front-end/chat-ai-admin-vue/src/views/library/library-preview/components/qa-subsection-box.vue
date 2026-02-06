<template>
  <div class="subsection-box">
    <!-- <div class="subsection-box-title">
      分段预览
      <span>共{{ props.total }}个分段</span>
    </div> -->
    <a-table :data-source="paragraphLists" :pagination="false">
      <a-table-column key="id" data-index="id" :width="1148">
        <template #title>{{ t('title_qa_total', { count: props.total }) }}</template>
        <template #default="{ record }">
          <div class="qa-list-box">
            <div class="list-item">
              <div class="list-label">{{ t('label_question') }}</div>
              <div class="list-content" v-html="textToHighlight(record.question, props.search)"></div>
            </div>
            <div class="list-item" v-if="record.similar_questions && record.similar_questions.length">
              <div class="list-label">{{ t('label_similar_questions') }}</div>
              <div class="list-content" v-html="textToHighlight(record.similar_questions.join('/'), props.search)"></div>
            </div>
            <div class="list-item">
              <div class="list-label">{{ t('label_answer') }}</div>
              <div class="list-content" v-html="textToHighlight(record.answer, props.search)"></div>
            </div>
            <div class="fragment-img" v-viewer>
              <img v-for="(item, index) in record.images" :key="index" :src="item" alt="" />
            </div>
          </div>
        </template>
      </a-table-column>
      <a-table-column :title="t('title_embed_status')" key="status_text" data-index="status_text" :width="128">
        <template #default="{ record }">
          <span class="status-tag" v-if="record.status == 0"><ClockCircleFilled /> {{ t('status_not_converted') }}</span>
          <span class="status-tag complete" v-if="record.status == 1"
            ><CheckCircleFilled /> {{ t('status_converted') }}</span
          >
          <a-tooltip placement="top" v-if="record.status == 2">
            <template #title>
              <span>{{ record.errmsg }}</span>
            </template>
            <span>
              <span class="status-tag status-error"><CloseCircleFilled /> {{ t('status_convert_error') }}</span>
            </span>
          </a-tooltip>
          <span class="status-tag running" v-if="record.status == 3"
            ><a-spin size="small" /> {{ t('status_converting') }}</span
          >
        </template>
      </a-table-column>
      <a-table-column :title="t('title_operation')" key="action" data-index="action" :width="120">
        <template #default="{ record, index }">
          <div class="right-opration">
            <a-dropdown>
              <div class="hover-btn-box">
                <StarFilled
                  @click="handleSetCategory(record, {})"
                  :style="{ color: getColor(record), 'font-size': '16px' }"
                  v-if="record.category_id > 0"
                />
                <StarOutlined
                  @click="handleSetCategory(record, startLists[0])"
                  style="font-size: 16px"
                  v-else
                />
              </div>
              <template #overlay>
                <a-menu>
                  <a-menu-item v-for="star in startLists" :key="star.id">
                    <div class="start-item" @click="handleSetCategory(record, star)">
                      <StarFilled :style="{ color: colorLists[star.type] }" />
                      <div>{{ star.name || '-' }}</div>
                    </div>
                  </a-menu-item>
                  <a-menu-item>
                    <div class="start-item" @click="handleOpenSetStartModal">
                      <SettingOutlined />
                      <div>{{ t('btn_mark_settings') }}</div>
                    </div>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
            <a-tooltip>
              <template #title>{{ t('tooltip_re_convert') }}</template>
              <div class="hover-btn-box" @click="toReSegmentationPage(record, index)">
                <SyncOutlined />
              </div>
            </a-tooltip>
            <a-dropdown placement="bottomRight">
              <div class="hover-btn-box">
                <MoreOutlined />
              </div>
              <template #overlay>
                <a-menu>
                  <a-menu-item>
                    <div @click.stop="handleOpenEditModal(record)">{{ t('btn_edit') }}</div>
                  </a-menu-item>
                  <a-menu-item>
                    <div @click.stop="hanldleDelete(record.id)">{{ t('btn_delete') }}</div>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </template>
      </a-table-column>
    </a-table>
    <ClassificationMarkModal @ok="getCategoryLists" ref="classificationMarkModalRef" />
  </div>
</template>
<script setup>
import { reactive, ref, computed, createVNode } from 'vue'
import { message } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  MoreOutlined,
  SyncOutlined,
  ClockCircleFilled,
  LoadingOutlined,
  StarOutlined,
  StarFilled,
  SettingOutlined
} from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import ClassificationMarkModal from './classification-mark-modal.vue'
import colorLists from '@/utils/starColors.js'
import {
  deleteParagraph,
  editParagraph,
  getCategoryList,
  updateParagraphCategory
} from '@/api/library'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-preview.components.qa-subsection-box')

const route = useRoute()
const query = route.query
const emit = defineEmits([
  'handleDelParagraph',
  'handleScrollTargetPage',
  'openEditSubscription',
  'handleConvert',
  'getStatrList'
])
const props = defineProps({
  paragraphLists: {
    type: Array,
    default: () => []
  },
  search: {
    type: String,
    default: ''
  },
  total: {
    type: [Number, String],
    default: 0
  }
})

const toReSegmentationPage = (item, index) => {
  let { id, title, content, question, answer, images, category_id } = item
  let similar_questions = item.similar_questions || []
  Modal.confirm({
    title: t('title_confirm_re_convert'),
    icon: null,
    content: t('msg_confirm_re_convert', { index: index + 1 }),
    onOk() {
      return new Promise((resolve, reject) => {
        editParagraph({ id, title, content, question, answer, images, category_id, similar_questions: JSON.stringify(similar_questions) })
          .then((res) => {
            emit('handleConvert', item)
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

const handleOpenEditModal = (item) => {
  emit('openEditSubscription', item)
}
const hanldleDelete = (id) => {
  Modal.confirm({
    title: t('title_tip'),
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

const startLists = ref([])
const getCategoryLists = () => {
  let params = {}
  if (query.id) {
    params.file_id = query.id
  }
  getCategoryList(params).then((res) => {
    startLists.value = res.data || []
    emit('getStatrList', res.data || [])
  })
}
getCategoryLists()

const classificationMarkModalRef = ref(null)
const handleOpenSetStartModal = () => {
  classificationMarkModalRef.value.show()
}

const getColor = (data) => {
  let type = startLists.value.filter((item) => item.id == data.category_id)[0]?.type
  if (type) {
    return colorLists[type]
  } else {
    return '#F4EA2A'
  }
}
const handleSetCategory = (item, star = {}) => {
  updateParagraphCategory({
    id: item.id,
    category_id: star.id || 0
  }).then((res) => {
    message.success(t('msg_set_success'))
    item.category_id = star.id
    getCategoryLists()
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

defineExpose({ handleOpenEditModal })
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
  .list-item {
    display: flex;
    flex-wrap: wrap;
    line-height: 22px;
    font-size: 14px;
    color: #595959;
    margin-bottom: 6px;
    .list-label {
      width: 65px;
      font-weight: 600;
      color: #262626;
    }
    .list-content {
      flex: 1;
      word-break: break-all;
      &::v-deep(.highlight) {
        background-color: #FFEB3B; /* 黄色高亮 */
        padding: 0 2px;
        border-radius: 2px;
        box-shadow: 0 1px 1px rgba(0,0,0,0.1);
      }
    }
  }
  .fragment-img {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding-left: 65px;
    img {
      width: 80px;
      height: 80px;
      border-radius: 6px;
      cursor: pointer;
    }
  }
}

.status-tag {
  display: inline-block;
  height: 24px;
  line-height: 24px;
  padding: 0 6px;
  border-radius: 2px;
  font-size: 14px;
  font-weight: 500;
  text-align: center;
  color: #3a4559;
  background-color: #edeff2;

  &.running {
    color: #2475fc;
    background-color: #e8effc;
  }

  &.complete {
    color: #21a665;
    background: #e8fcf3;
  }

  &.error {
    cursor: pointer;
    color: #fb363f;
    background-color: #f5c6c8;
  }

  &.warning {
    cursor: pointer;
    background: #faebe6;
    color: #ed744a;
  }

  &.status-error {
    cursor: pointer;
    color: #fb363f;
  }
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

.start-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
  font-size: 14px;
  .anticon {
    font-size: 16px;
  }
}
</style>
