<template>
  <div class="subsection-box">
    <a-table
      :data-source="paragraphLists"
      row-key="id"
      :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
      :pagination="false"
      :loading="props.isLoading"
    >
      <a-table-column key="id" data-index="id" :width="1148">
        <template #title>{{ t('title_qa', { total: props.total }) }}</template>
        <template #default="{ record }">
          <div class="qa-list-box">
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

                {{ record.question }}
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
                      <span>{{ t('msg_similar_questions', { count: record.similar_questions.length }) }}</span>
                    </template>
                    <a>（{{ record.similar_questions.length }}）</a>
                  </a-popover>
                </template>
              </div>
            </div>
            <div class="list-item list-item-answer">
              <div class="list-label">{{ t('label_answer') }}</div>
              <div class="list-content">{{ record.answer }}</div>
            </div>
            <div class="fragment-img" v-viewer>
              <img v-for="(item, index) in record.images" :key="index" :src="item" alt="" />
            </div>
          </div>
        </template>
      </a-table-column>
      <a-table-column :title="t('title_import_to_knowledge')" key="is_import" data-index="is_import" :width="180">
        <template #default="{ record }">
          <div v-if="record.is_import == 1" class="status-block status-green">
            <CheckCircleFilled />{{ t('status_imported') }}
          </div>
          <div v-if="record.is_import == 0" class="status-block status-gray">
            <ExclamationCircleFilled />{{ t('status_not_imported') }}
          </div>
          <div>
            <a
              :href="`/#/library/details/knowledge-document?id=${record.library_id}`"
              target="_blank"
              >{{ record.library_name }}</a
            >
          </div>
        </template>
      </a-table-column>
      <a-table-column :title="t('title_action')" key="action" data-index="action" :width="120">
        <template #default="{ record, index }">
          <div class="right-opration">
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
import { reactive, ref, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  CheckCircleFilled,
  ExclamationCircleFilled,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import { deleteFAQFileQA } from '@/api/library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.ai-extract-faq.detail.components.subsection-box')

const emit = defineEmits(['handleDelParagraph', 'openEditSubscription'])
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
  }
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
        deleteFAQFileQA({ ids: id })
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

.status-block {
  display: flex;
  align-items: center;
  width: fit-content;
  gap: 3px;
  padding: 0 6px;
  height: 22px;
  border-radius: 6px;
  line-height: 22px;
  font-size: 14px;
  &.status-bule {
    background: #d4e3fc;
    color: #2475fc;
  }
  &.status-green {
    background: #cafce4;
    color: #21a665;
  }
  &.status-gray {
    background: #ebeff5;
    color: #3a4559;
  }
}

.qa-list-box {
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
