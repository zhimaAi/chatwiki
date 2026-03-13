<template>
  <div class="chatclaw-token-log-page">
    <div class="page-title">{{ t('title') }}</div>
    <div v-if="loading" class="loading-box">
      <a-spin />
    </div>
    <div v-else class="table-wrapper">
      <a-table
        :columns="columns"
        :data-source="list"
        :pagination="pagination"
        row-key="id"
        size="middle"
        @change="onTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'id'">
            {{ record.id }}
          </template>
          <template v-else-if="column.key === 'token'">
            <a-tooltip :title="record.token">
              <span class="token-cell">{{ record.token || '--' }}</span>
            </a-tooltip>
          </template>
          <template v-else-if="column.key === 'create_time'">
            {{ formatTime(record.create_time) }}
          </template>
          <template v-else-if="column.key === 'expired_at'">
            {{ formatExpiredAt(record.expired_at) }}
          </template>
          <template v-else-if="column.key === 'revoke_time'">
            {{ formatTime(record.revoke_time) }}
          </template>
          <template v-else-if="column.key === 'status'">
            <span :class="getStatusClass(record)">{{ getStatusText(record) }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-button
              type="link"
              size="small"
              :disabled="isRevoked(record)"
              @click="handleForceOffline(record)"
            >
              {{ t('action_force_offline') }}
            </a-button>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getChatClawTokenLogListApi, forceOfflineChatClawTokenApi } from '@/api/chatclaw'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.chatclaw-token-log')

const list = ref([])
const total = ref(0)
const loading = ref(true)
const paginationState = ref({ current: 1, pageSize: 10 })

const pagination = computed(() => ({
  current: paginationState.value.current,
  pageSize: paginationState.value.pageSize,
  total: total.value,
  showSizeChanger: true,
  showTotal: (totalNum) => `共 ${totalNum} 条`,
  pageSizeOptions: ['10', '20', '50']
}))

const columns = computed(() => [
  { title: t('col_id'), key: 'id', width: 72, align: 'center' },
  { title: 'Token', dataIndex: 'token', key: 'token', width: 200, ellipsis: true },
  { title: t('col_os'), dataIndex: 'os_type', key: 'os_type', width: 100 },
  { title: t('col_os_version'), dataIndex: 'os_version', key: 'os_version', width: 110 },
  { title: 'IP', dataIndex: 'client_ip', key: 'client_ip', width: 130 },
  { title: t('col_issue_time'), key: 'create_time', width: 130 },
  { title: t('col_expire_time'), key: 'expired_at', width: 130 },
  { title: t('col_revoke_time'), key: 'revoke_time', width: 130 },
  { title: t('col_status'), key: 'status', width: 88 },
  { title: t('col_action'), key: 'action', width: 100 }
])

function formatTime(ts) {
  const n = Number(ts)
  if (!n || isNaN(n)) return '--'
  const d = new Date(n * 1000)
  const y = String(d.getFullYear()).slice(-2)
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${h}:${min}`
}

function formatExpiredAt(ts) {
  if (!ts) return '--'
  return formatTime(ts)
}

function isRevoked(record) {
  return Number(record?.status) === 2
}

function isNaturalExpired(record) {
  const n = Number(record?.expired_at)
  return !!n && !isNaN(n) && Date.now() / 1000 > n
}

function getStatusText(record) {
  if (isRevoked(record)) return t('status_revoked')
  if (isNaturalExpired(record)) return t('status_expired')
  return t('status_active')
}

function getStatusClass(record) {
  if (isRevoked(record)) return 'status-revoked'
  if (isNaturalExpired(record)) return 'status-expired'
  return 'status-active'
}

function fetchList() {
  loading.value = true
  const { current, pageSize } = paginationState.value
  getChatClawTokenLogListApi({ page: current, size: pageSize })
    .then((res) => {
      const data = res.data || {}
      list.value = data.list || []
      total.value = Number(data.total) || 0
    })
    .finally(() => {
      loading.value = false
    })
}

function onTableChange(pag) {
  paginationState.value.current = pag.current
  paginationState.value.pageSize = pag.pageSize
  fetchList()
}

function handleForceOffline(record) {
  Modal.confirm({
    title: t('confirm_force_offline_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_force_offline_content'),
    okText: t('action_force_offline'),
    cancelText: t('cancel'),
    okType: 'danger',
    onOk() {
      return forceOfflineChatClawTokenApi({ id: record.id }).then(() => {
        message.success(t('msg_force_offline_success'))
        fetchList()
      })
    }
  })
}

onMounted(() => {
  fetchList()
})
</script>

<style lang="less" scoped>
.chatclaw-token-log-page {
  padding: 24px;
  height: 100%;
  width: 100%;
  overflow-y: auto;
  background-color: #fff;

  .page-title {
    color: #000000;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 24px;
    padding: 0;
  }

  .loading-box {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 200px;
  }

  .table-wrapper {
    background: #fff;
    padding: 16px;
    border-radius: 4px;
  }

  .status-active {
    color: #52c41a;
  }

  .status-expired {
    color: #999;
  }

  .status-revoked {
    color: #ff4d4f;
  }

  .token-cell {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: bottom;
  }
}
</style>
