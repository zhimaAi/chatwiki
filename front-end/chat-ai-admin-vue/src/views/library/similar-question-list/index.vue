<template>
  <div class="similar-question-list-page">
    <!-- 页面标题 -->
    <div class="page-title">{{ t('page_title', { total }) }}</div>

    <a-alert
      class="page-alert"
      :message="t('page_alert')"
      type="info"
      show-icon
    />

    <div class="batch-actions">
      <a-button @click="handleSettings">
        <template #icon><SettingOutlined /></template>
        {{ t('settings') }}
      </a-button>
      <a-button @click="handleBatchMerge">{{ t('batch_merge') }}</a-button>
      <a-button @click="handleBatchDelete">{{ t('batch_delete') }}</a-button>
    </div>

    <!-- 列表表格 -->
    <div class="table-wrapper">
      <a-table
        :data-source="questionList"
        row-key="id"
        :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }"
        :pagination="false"
        :loading="loading"
      >
        <a-table-column key="content" data-index="content">
          <template #title>{{ t('qa_title', { total }) }}</template>
          <template #default="{ record }">
            <div class="qa-list-box">
              <div class="list-item">
                <div class="list-label">{{ t('label_question') }}</div>
                <div class="list-content">{{ record.question }}</div>
              </div>
              <div class="list-item list-item-answer">
                <div class="list-label">{{ t('label_answer') }}</div>
                <div class="list-content">{{ record.answer }}</div>
              </div>
              <div class="fragment-img" v-if="record.images && record.images.length" v-viewer>
                <img v-for="(img, index) in record.images" :key="index" :src="img" alt="" />
              </div>
            </div>
          </template>
        </a-table-column>
        <a-table-column key="action" data-index="action" :width="120" align="center">
          <template #title>{{ t('col_action') }}</template>
          <template #default="{ record }">
            <div class="action-buttons">
              <a-button type="link" size="small" @click="handleMerge(record)">{{ t('merge') }}</a-button>
              <a-button type="link" size="small" danger @click="handleDelete(record)">{{ t('delete') }}</a-button>
            </div>
          </template>
        </a-table-column>
      </a-table>
    </div>

    <!-- 合并弹窗 -->
    <MergeModal ref="mergeModalRef" @confirm="handleMergeConfirm" @cancel="handleMergeCancel" />
    <!-- 设置弹窗 -->
    <SettingsModal ref="settingsModalRef" @confirm="handleSettingsConfirm" @cancel="handleSettingsCancel" />
  </div>
</template>

<script setup>
import { ref, computed, createVNode, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { getSimilarQuestions, deleteParagraph, mergeQAParagraph } from '@/api/library'
import MergeModal from './components/merge-modal.vue'
import SettingsModal from './components/settings-modal.vue'

const route = useRoute()
const { t } = useI18n('views.library.similar-question-list.index')
const loading = ref(false)

// 选中的行
const selectedRowKeys = ref([])

// 数据
const questionList = ref([])

// 总数
const total = computed(() => questionList.value.length)

// 合并弹窗相关
const mergeModalRef = ref(null)
const source = ref(null)

// 设置弹窗相关
const settingsModalRef = ref(null)
const similarityThreshold = ref(0.8)

// 选择变化
const onSelectChange = (keys) => {
  selectedRowKeys.value = keys
}

// 合并单条
const handleMerge = (record) => {
  const options = [{...source.value}, record]

  mergeModalRef.value?.open({
    options: [...options],
    defaultSelected: options[0].id
  })
}

// 删除单条
const handleDelete = (record) => {
  Modal.confirm({
    title: t('confirm_delete_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_delete_content'),
    onOk() {
      return new Promise((resolve) => {
        deleteParagraph({ id: record.id })
          .then(() => {
            message.success(t('delete_success'))
            getSimilarQuestionList()
            // 跨窗口通知父页面刷新列表
            window.opener?.postMessage({ type: 'qa-merged', libraryId: route.query.library_id }, '*')
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

// 批量合并
const handleBatchMerge = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning(t('please_select_to_merge'))
    return
  }
  const selectedRecords = questionList.value.filter(item => selectedRowKeys.value.includes(item.id))
  selectedRecords.unshift(source.value)
  mergeModalRef.value?.open({
    options: [...selectedRecords],
    defaultSelected: selectedRecords[0].id
  })
}

// 批量删除
const handleBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning(t('please_select_to_delete'))
    return
  }
  Modal.confirm({
    title: t('confirm_batch_delete_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_batch_delete_content', { count: selectedRowKeys.value.length }),
    onOk() {
      return new Promise((resolve) => {
        deleteParagraph({ id: selectedRowKeys.value.join(',') })
          .then(() => {
            message.success(t('batch_delete_success'))
            selectedRowKeys.value = []
            getSimilarQuestionList()
            // 跨窗口通知父页面刷新列表
            window.opener?.postMessage({ type: 'qa-merged', libraryId: route.query.library_id }, '*')
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

// 处理合并确认
const handleMergeConfirm = (selectedId, options) => {
  if (!selectedId) {
    message.warning(t('please_select_merge_result'))
    return
  }

  // 从 options 中获取所有需要删除的ID（除了选中的目标ID）
  const allIds = options.map(item => item.id)
  const deleteIds = allIds.filter(id => id !== selectedId)

  // 调用合并接口
  mergeQAParagraph({
    source_data_id: source.value.id,
    target_data_id: selectedId,
    delete_data_ids: deleteIds.join(',')
  })
    .then(() => {
      message.success(t('merge_success'))
      selectedRowKeys.value = []
      getSimilarQuestionList()
      // 跨窗口通知父页面刷新列表
      window.opener?.postMessage({ type: 'qa-merged', libraryId: route.query.library_id }, '*')
    })
    .catch(() => {})
}

// 处理合并取消
const handleMergeCancel = () => {
}

// 打开设置弹窗
const handleSettings = () => {
  settingsModalRef.value?.open({
    defaultValue: similarityThreshold.value
  })
}

// 处理设置确认
const handleSettingsConfirm = (value) => {
  similarityThreshold.value = value
  message.success(t('settings_success'))
  // 刷新列表
  getSimilarQuestionList()
}

// 处理设置取消
const handleSettingsCancel = () => {
}

const getSimilarQuestionList = () => {
  loading.value = true
  getSimilarQuestions({
    library_id: route.query.library_id,
    data_id: route.query.question_id,
    type: 'list'
  }).then((res) => {
    loading.value = false
    questionList.value = res.data.list || []
    source.value = res.data.source_data || {}
  })
}

// 组件挂载时获取相似问答列表
onMounted(() => {
  getSimilarQuestionList()
})
</script>

<style lang="less" scoped>
.similar-question-list-page {
  padding: 24px;
  background: #fff;
  min-height: 100%;
  max-width: 1200px;
  margin: 0 auto;

  .page-title {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    line-height: 24px;
    margin-bottom: 16px;
  }

  .page-alert {
    margin-bottom: 16px;
  }

  .batch-actions {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
  }

  .table-wrapper {
    :deep(.ant-table) {
      .ant-table-thead > tr > th {
        background: #fafafa;
        font-weight: 600;
      }
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
        color: #8c8c8c;
        flex-shrink: 0;
      }

      .list-content {
        flex: 1;
        word-break: break-all;
      }
    }

    .list-item-answer {
      .list-content {
        color: #8c8c8c;
      }
    }

    .fragment-img {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      padding-left: 40px;
      margin-top: 8px;

      img {
        width: 80px;
        height: 80px;
        border-radius: 6px;
        cursor: pointer;
        object-fit: cover;
      }
    }
  }

  .action-buttons {
    display: flex;
    justify-content: center;
    gap: 8px;
  }
}
</style>
