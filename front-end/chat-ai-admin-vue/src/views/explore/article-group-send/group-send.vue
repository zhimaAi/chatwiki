<template>
  <div class="group-send-wrapper">
        <div class="toolbar">
          <a-button type="primary" @click="handleCreateSend">创建群发</a-button>
          <!-- <a-button :disabled="selectedRowKeys.length === 0">批量管理</a-button> -->
        </div>

        <div class="table-box">
          <a-table :columns="columns" :data-source="taskList" :loading="loadingTasks" row-key="id"
            :pagination="{ current: pager.page, pageSize: pager.size, total: pager.total, showSizeChanger: true, pageSizeOptions: ['10', '20', '50', '100'] }"
            @change="onTableChange" :row-class-name="rowClassName">
            <template #headerCell="{ column }">
              <span v-if="typeof column.title === 'string'">{{ column.title }}</span>
            </template>
            <template #expandIcon="{ expanded, onExpand, record }">
              <span class="expand-icon" @click="onExpand(record)">
                <svg-icon v-if="expanded" name="table-hide" size="16px" />
                <svg-icon v-else name="table-expand" size="16px" />
              </span>
            </template>
            <template #expandedRowRender="{ record }">
              <div class="expanded-draft">
                <img v-if="record.thumb_url" class="thumb" :src="record.thumb_url" />
                <img v-else class="thumb" src="@/assets/img/default-cover.png" />
                <div class="info">
                  <div class="title">{{ record.title }}</div>
                  <div class="meta">所属分组：{{ record.group_name }}</div>
                  <div class="digest">{{ record.digest || '暂无摘要' }}</div>
                </div>
                <div class="expanded-meta" @click="openCommentDrawer(record)">
                  <svg-icon name="comment" size="32px" style="color: transparent;" />
                  {{ record.max_comment_id }}
                </div>
              </div>
            </template>
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'task_name'">
                <div class="task-cell">
                  <a-tag v-if="record.is_top == '1'" color="blue">置顶</a-tag>
                  <span class="name">{{ record.task_name }}</span>
                </div>
              </template>
              <template v-else-if="column.key === 'send_time'">
                <div class="send-time-cell">
                  <div class="status-tag" :class="{'status-sending': record.send_status == '1', 'status-success': record.send_status == '2'}">
                    <svg-icon style="color: white;" name="status-fail" v-if="record.send_status == '-1'" />
                    <svg-icon style="color: white;" name="status-fail" v-else-if="record.send_status == '0'" />
                    <svg-icon style="color: white;" name="status-sending" v-else-if="record.send_status == '1'" />
                    <svg-icon style="color: white;" name="status-success" v-else-if="record.send_status == '2'" />
                    <svg-icon style="color: white;" name="status-fail" v-else-if="record.send_status == '3'" />
                    {{ statusTextMap[record.send_status] }}
                    <a-tooltip :title="record.send_res?.errmsg" v-if="record.send_status === '3'">
                      <ExclamationCircleOutlined style="font-size: 16px; color: #FF4D4F;" />
                    </a-tooltip>
                  </div>
                  <div v-if="record.send_time != '0'" class="time">定时群发：{{ formatTime(record.send_time) }}</div>
                  <div v-else class="time">立即发送：{{ formatTime(record.create_time) }}</div>
                </div>
              </template>
              <template v-else-if="column.key === 'receiver'">
                <span style="color: #595959;">{{ toUserTypeMap[record.to_user_type] || '全部粉丝' }}</span>
              </template>
              <template v-else-if="column.key === 'open_status'">
                <a-switch :checked="record.open_status == '1'" checked-children="开" un-checked-children="关"
                  @change="onToggleOpen(record, $event)" />
              </template>
              <template v-else-if="column.key === 'comment_status'">
                <a-switch :checked="record.comment_status == '1'" checked-children="开" un-checked-children="关"
                  @change="onToggleMessageComment(record, $event)" />
              </template>
              <template v-else-if="column.key === 'ai_comment_status'">
                <div class="ai-comment-cell">
                  <a-switch :checked="record.ai_comment_status == '1'" checked-children="开" un-checked-children="关"
                    @change="onToggleAiComment(record, $event)" />
                  <div class="comment-rule-cell">
                    <span class="rule">
                      <span v-if="record.is_default == '1'" class="default-tag">默认</span>
                      {{ record.comment_rule_name || '默认规则' }}
                    </span>
                    <a class="edit-link" @click="editCommentRule(record)">修改</a>
                  </div>
                </div>
              </template>
              <template v-else-if="column.key === 'actions'">
                <a-space style="gap: 16px;">
                  <a @click="editTask(record)">编辑</a>
                  <a-dropdown>
                    <a>更多</a>
                    <template #overlay>
                      <a-menu>
                        <a-menu-item @click="toggleTop(record)">{{ record.is_top == '1' ? '取消置顶' : '置 顶'
                          }}</a-menu-item>
                        <a-menu-item @click="deleteTask(record)">删 除</a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </a-space>
              </template>
            </template>
          </a-table>
        </div>

        
        <CommentDrawer ref="commentDrawerRef" />
        <CommentRuleModal ref="commentRuleModalRef" @updated="getTaskList" />
        <CreateSendModal ref="editSendModalRef" :app-id="appId" :access-key="accessKey" @updated="getTaskList" @created="getTaskList" />
  </div>
</template>

<script setup>
import { reactive, ref, createVNode, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { getSpecifyAbilityConfig } from '@/api/explore'
import {
  getBatchSendTaskList,
  setBatchSendTaskTopStatus,
  deleteBatchSendTask,
  setBatchSendTaskEnableStatus,
  setBatchSendTaskCommentRuleStatus,
  changeCommentStatus
} from '@/api/robot'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import CommentDrawer from './components/comment-drawer.vue'
import CommentRuleModal from './components/comment-rule-modal.vue'
import CreateSendModal from './components/create-send-modal.vue'
import { addNoReferrerMeta, removeNoReferrerMeta } from '@/utils/index.js'

const router = useRouter()

const props = defineProps({
  appId: { type: String, default: '' },
  accessKey: { type: String, default: '' }
})


const taskList = ref([])
const loadingTasks = ref(false)
const pager = reactive({ page: 1, size: 10, total: 0 })

const statusTextMap = {
  '-1': '已删除',
  '0': '未发送',
  '1': '发送中',
  '2': '已发送',
  '3': '发送失败',
}
const toUserTypeMap = { '0': '全部粉丝' }

const columns = [
  { title: '群发', dataIndex: 'task_name', key: 'task_name' },
  { title: '群发时间', dataIndex: 'send_time', key: 'send_time' },
  { title: '接收者', dataIndex: 'receiver', key: 'receiver' },
  { title: '启用', dataIndex: 'open_status', key: 'open_status' },
  { title: '留言', dataIndex: 'comment_status', key: 'comment_status' },
  { title: 'AI评论', dataIndex: 'ai_comment_status', key: 'ai_comment_status' },
  { title: '操作', dataIndex: 'actions', key: 'actions' },
]

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getTaskList()
}

const rowClassName = (record) => {
  return record.is_top == '1' ? 'row-top' : ''
}

const getTaskList = () => {
  loadingTasks.value = true
  getBatchSendTaskList({ page: pager.page, size: pager.size, app_id: props.appId }).then((res) => {
    const data = res?.data || { list: [], total: 0 }
    data.list.forEach(item => {
      item.send_res = JSON.parse(item.send_res || '{}')
    })
    taskList.value = data.list || []
    pager.total = +data.total || 0
  }).finally(() => { loadingTasks.value = false })
}

const toggleTop = async (record) => {
  // 如果是取消置顶需要二次确认
  if (record.is_top == '1') {
    Modal.confirm({
      title: '确认取消置顶？',
      icon: createVNode(ExclamationCircleOutlined),
      onOk: async () => {
        await setBatchSendTaskTopStatus({ task_id: record.id, is_top: record.is_top == '1' ? 0 : 1 })
        message.success('已取消置顶')
        getTaskList()
      }
    })
    return
  }
  await setBatchSendTaskTopStatus({ task_id: record.id, is_top: record.is_top == '1' ? 0 : 1 })
  message.success(record.is_top == '1' ? '已取消置顶' : '已置顶')
  getTaskList()
}

const deleteTask = (record) => {
  Modal.confirm({
    title: `确认删除群发任务：${record.task_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk: async () => {
      await deleteBatchSendTask({ task_id: record.id })
      message.success('删除成功')
      getTaskList()
    }
  })
}

const onToggleOpen = async (record, checked) => {
  const open_status = checked ? 1 : 0
  // 如果是关闭，需要二次确认
  if (!checked) {
    Modal.confirm({
      title: '确认关闭群发任务？',
      icon: createVNode(ExclamationCircleOutlined),
      onOk: async () => {
        await setBatchSendTaskEnableStatus({ task_id: record.id, open_status })
        record.open_status = String(open_status)
        message.success('已关闭群发任务')
        getTaskList()
      }
    })
    return
  }
  await setBatchSendTaskEnableStatus({ task_id: record.id, open_status })
  record.open_status = String(open_status)
  message.success('操作成功')
}

const onToggleAiComment = async (record, checked) => {
  if (checked) {
    try {
      const res = await getSpecifyAbilityConfig({ ability_type: 'official_account_ai_comment' })
      const status = res?.data?.user_config?.switch_status
      if (String(status || '0') !== '1') {
        Modal.confirm({
          title: 'AI评论',
          content: createVNode('div', null, [
            createVNode('span', { style: 'color: #ff4d4f;' }, 'AI评论功能暂未开启'),
            createVNode('span', null, '，您可到探索>功能中 去开启')
          ]),
          okText: '去开启',
          cancelText: '取消',
          onOk: () => {
            router.push('/explore/index/ai-comment-management')
          }
        })
        return
      }
    } catch (e) {
      console.error(e)
    }
  }

  const ai_comment_status = checked ? 1 : 0
  // 关闭需要二次确认
  if (!checked) {
    Modal.confirm({
      title: '确认关闭AI评论规则？',
      icon: createVNode(ExclamationCircleOutlined),
      onOk: async () => {
        await setBatchSendTaskCommentRuleStatus({ task_id: record.id, ai_comment_status })
        record.ai_comment_status = String(ai_comment_status)
        message.success('已关闭AI评论规则')
        getTaskList()
      }
    })
    return
  }
  await setBatchSendTaskCommentRuleStatus({ task_id: record.id, ai_comment_status })
  record.ai_comment_status = String(ai_comment_status)
  // message.success('操作成功')
  commentRuleModalRef.value && commentRuleModalRef.value.show({ ...record })
}

const onToggleMessageComment = async (record, checked) => {
  const comment_status = checked ? 1 : 0
  // 关闭需要二次确认
  if (!checked) {
    Modal.confirm({
      title: '确认关闭留言？',
      icon: createVNode(ExclamationCircleOutlined),
      onOk: async () => {
        await changeCommentStatus({ task_id: record.id, msg_id: record.msg_data_id, access_key: props.accessKey, comment_status })
        record.comment_status = String(comment_status)
        message.success('已关闭留言')
        getTaskList()
      }
    })
    return
  }
  await changeCommentStatus({ task_id: record.id, msg_id: record.msg_data_id, access_key: props.accessKey, comment_status })
  record.comment_status = String(comment_status)
  message.success('操作成功')
}

const editSendModalRef = ref(null)
const editTask = (record) => { editSendModalRef.value && editSendModalRef.value.show({ task: record }) }
const handleCreateSend = () => { editSendModalRef.value && editSendModalRef.value.show({}) }
const commentRuleModalRef = ref(null)
const editCommentRule = (record) => { commentRuleModalRef.value && commentRuleModalRef.value.show({ ...record }) }


const formatTime = (ts) => {
  if (!ts) return ''
  const d = new Date(Number(ts) * 1000)
  const yy = String(d.getFullYear()).slice(2)
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${yy}-${m}-${dd} ${hh}:${mm}`
}

watch(() => props.appId, () => {
  pager.page = 1
  if (props.appId) getTaskList()
}, { immediate: true })

onMounted(() => { addNoReferrerMeta() })
onUnmounted(() => { removeNoReferrerMeta() })

const commentDrawerRef = ref(null)
const openCommentDrawer = (record) => { commentDrawerRef.value && commentDrawerRef.value.show({ ...record }) }
</script>

<style lang="less" scoped>


.line-box {
  height: 1px;
  background: #F0F0F0;
  margin-top: 16px;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.table-box {
  margin-top: 12px;
}

.task-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
}

::v-deep(.ant-table-tbody > tr.row-top > td) {
  background-color: #F5F5F5 !important;
}

.send-time-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-direction: column;
  align-items: flex-start;
  color: #595959;

  .status-tag {
    display: flex;
    padding: 0 6px;
    align-items: center;
    gap: 2px;
    border-radius: 6px;
    background: #FAE4DC;
    color: #ED744A;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }

  .status-sending {
    background: #D4E3FC;
    color: #2475FC;
  }

  .status-success {
    background: #CAFCE4;
    color: #21A665;
  }
}

.expanded-draft {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 12px;
  border: 1px solid #edeff2;
  border-radius: 8px;
  margin-left: 48px;
}

.expand-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.expanded-draft .thumb {
  width: 146px;
  height: 96px;
  border-radius: 6px;
  object-fit: cover;
}

.expanded-draft .info {
  flex: 1;
  overflow: hidden;
}

.expanded-draft .info .title {
  font-weight: 600;
  color: #262626;
}

.expanded-draft .info .meta {
  color: #8c8c8c;
  margin: 4px 0;
}

.expanded-draft .info .digest {
  color: #595959;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.expanded-draft .expanded-meta {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  color: #8c8c8c;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  margin-left: 4px;
  cursor: pointer;
}

.ai-comment-cell {
  .comment-rule-cell {
    display: block;
    margin-top: 4px;
    color: #8C8C8C;
    font-size: 12px;

    .edit-link {
      margin-left: 4px;
    }

    .default-tag {
      color: #1677ff;
      background: #e6f4ff;
      padding: 0px 2px;
      border-radius: 4px;
      font-size: 12px;
      border: 1px solid #91caff;
    }
  }
}

 
</style>
