<template>
  <div class="user-model-page">
    <a-alert show-icon style="padding-top: 16px">
      <template #message>
        <p>{{ t('alert_message_1') }}</p>
        <p>{{ t('alert_message_2') }}</p>
      </template>
    </a-alert>
    <div class="search-block">
      <div class="left-block">
        <span>{{ t('auto_cluster') }}</span>
        <a-tooltip>
          <template #title>{{ t('auto_cluster_tooltip') }}</template>
          <QuestionCircleOutlined />：
        </a-tooltip>
        <a-switch
          @change="handleChangeSwitch"
          :checked="unknown_summary_status"
          :checked-children="t('switch_on')"
          :un-checked-children="t('switch_off')"
        />
        <a @click="handleOpenClusterModal">{{ t('settings') }}</a>
      </div>
      <div class="right-block">
          <DateSelect datekey="4" @dateChange="onDateChange" />
          <a-button type="primary" @click="onSearch">{{ t('query') }}</a-button>
          <a-dropdown>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="handleDownload('docx')" key="1">{{ t('download_docx') }}</a-menu-item>
                <a-menu-item @click="handleDownload('xlsx')" key="2">{{ t('download_xlsx') }}</a-menu-item>
              </a-menu>
            </template>
            <a-button type="primary">
              {{ t('download') }}
              <DownOutlined />
            </a-button>
          </a-dropdown>
        </div>
    </div>
    <div class="list-box">
      <a-table
        :data-source="tableData"
        :loading="loading"
        :pagination="{
          current: pager.page,
          total: pager.total,
          pageSize: pager.size,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100']
        }"
        @change="onTableChange"
      >
        <a-table-column key="question" :title="t('column_cluster_question')" :width="480">
          <template #default="{ record }">
            <div class="qa-list-box">
              <div class="list-item">
                <div class="list-label">{{ t('column_question') }}</div>
                <div class="list-content">
                  {{ record.question }}
                  <template v-if="record.unknown_total > 0">
                    <a-popover placement="topLeft">
                      <template #content>
                        <div class="classify-scroll-box">
                          <div
                            class="list-item"
                            v-for="(item, index) in buildClusterQuestionList(record)"
                            :key="index"
                          >
                            <a @click="handleOpenChatContext(record, item)">{{ item }}</a>
                          </div>
                        </div>
                      </template>
                      <template #title>
                        <span>{{
                          t('cluster_count', { count: buildClusterQuestionList(record).length })
                        }}</span>
                      </template>
                      <span class="context-link">
                        <a @click.stop="handleOpenChatContext(record, record.question)">{{
                          t('action_view_context')
                        }}</a>
                        <a>（{{ buildClusterQuestionList(record).length }}）</a>
                      </span>
                    </a-popover>
                  </template>
                  <template v-else>
                    <a @click="handleOpenChatContext(record, record.question)">{{
                      t('action_view_context')
                    }}</a>
                  </template>
                </div>
              </div>
              <div class="list-item" v-if="record.answer">
                <div class="list-label">{{ t('column_answer') }}</div>
                <div class="list-content">{{ record.answer }}</div>
              </div>
              <div class="fragment-img" v-viewer>
                <img v-for="(item, index) in record.images" :key="index" :src="item" alt="" />
              </div>
            </div>
          </template>
        </a-table-column>
        <!-- <a-table-column key="unknown_total" title="未知问题条数" :width="140">
            <template #default="{ record }">
              <a-popover placement="topLeft">
                <template #content>
                  <div class="classify-scroll-box">
                    <div
                      class="list-item"
                      v-for="(item, index) in record.unknown_list"
                      :key="index"
                    >
                      {{ item }}
                    </div>
                  </div>
                </template>
                <template #title>
                  <span>共聚类了{{ record.unknown_total }}个未知问题</span>
                </template>
                <a>{{ record.unknown_total }}</a>
              </a-popover>
            </template>
          </a-table-column> -->
        <a-table-column key="show_date" :title="t('column_trigger_date')" :width="120">
          <template #default="{ record }">{{ record.show_date }} </template>
        </a-table-column>
        <a-table-column key="to_library_id" :title="t('column_imported')" :width="160">
          <template #default="{ record }">
            <div class="import-td-box">
              <div class="status-block success" v-if="record.to_library_id > 0">
                <span></span>
                {{ t('imported_yes') }}
              </div>
              <div class="status-block none" v-else>
                <span></span>
                {{ t('imported_no') }}
              </div>
              <a
                :href="`/#/library/details/knowledge-document?id=${record.to_library_id}`"
                target="_blank"
                >{{ record.to_library_name }}</a
              >
            </div>
          </template>
        </a-table-column>
        <a-table-column key="action" :title="t('column_action')" :width="220">
          <template #default="{ record }">
            <a-flex :gap="8">
              <a @click="handleOpenAnswerModal(record)">{{ t('action_set_answer') }}</a>
              <a @click="handleOpenLibrary(record)">{{ t('action_import_library') }}</a>
            </a-flex>
          </template>
        </a-table-column>
      </a-table>
    </div>
    <ImportLibraryModal @ok="getTableData" ref="importLibraryModalRef" />
    <AutoClusterModal @ok="handleSaveCluster" ref="autoClusterModalRef" />
    <SetAnswerModal @ok="getTableData" ref="setAnswerModalRef" />
    <a-drawer v-model:open="contextVisible" :title="t('context_title')" :width="760" destroyOnClose>
      <a-spin :spinning="contextLoading">
        <MessageList
          ref="messageListRef"
          :messages="contextMessages"
          :robot-info="contextRobotInfo"
          :is-empty="contextMessages.length === 0"
          :channel-item="[]"
        />
      </a-spin>
    </a-drawer>
  </div>
</template>

<script setup>
import { reactive, ref, computed, nextTick } from 'vue'
import { QuestionCircleOutlined, DownOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import ImportLibraryModal from './components/import-library-modal.vue'
import AutoClusterModal from './components/auto-cluster-modal.vue'
import SetAnswerModal from './components/set-answer-modal.vue'
import MessageList from '../../session-record/components/message-list.vue'
import { getUnknownIssueSummary, setUnknownIssueSummary, getUnknownIssueChatContext } from '@/api/robot/index.js'
import dayjs from 'dayjs'
import { useRobotStore } from '@/stores/modules/robot'
import DateSelect from '../components/date.vue'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'
import { extractVoiceInfo, removeVoiceFormat } from '@/utils/index'
const { t } = useI18n('views.robot.robot-config.unknown-issue.summarize.index')
const userStore = useUserStore()
const robotStore = useRobotStore()

const { getRobot } = robotStore

const robotInfo = computed(() => {
  return robotStore.robotInfo
})

const unknown_summary_status = computed(() => {
  return robotInfo.value.unknown_summary_status == 1
})

const query = useRoute().query
const router = useRouter()
const activeKey = ref(2)

const requestParams = reactive({
  start_day: '',
  end_day: ''
})

const pager = reactive({
  page: 1,
  size: 10,
  total: 0
})
const tableData = ref([])
const loading = ref(false)
const getTableData = () => {
  let parmas = {
    robot_id: query.id,
    ...requestParams,
    ...pager
  }
  loading.value = true
  getUnknownIssueSummary({
    ...parmas
  })
    .then((res) => {
      tableData.value = res.data.list.map((item) => {
        return {
          ...item,
          images: item.images ? JSON.parse(item.images) : [],
          unknown_list: item.unknown_list ? JSON.parse(item.unknown_list) : []
        }
      })
      pager.total = +res.data.total || 0
    })
    .finally(() => {
      loading.value = false
    })
}
getTableData()

const handleDownload = (type) => {
  let TOKEN = userStore.getToken
  let start_day = requestParams.start_day
  let end_day = requestParams.end_day
  let srcUrl = `/manage/getUnknownIssueSummary?robot_id=${query.id}&token=${TOKEN}&start_day=${start_day}&end_day=${end_day}&export=${type}`
  window.location.href = srcUrl
}
const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getTableData()
}

const onSearch = () => {
  pager.page = 1
  getTableData()
}

const onDateChange = (data) => {
  requestParams.start_day = data.start_time
  requestParams.end_day = data.end_time
  onSearch()
}

const importLibraryModalRef = ref(null)
const handleOpenLibrary = (record) => {
  importLibraryModalRef.value.showModal({
    id: record.id,
    library_id: record.to_library_id > 0 ? record.to_library_id : '',
    answer: record.answer,
    question: record.question,
    similar_questions: record.unknown_list,
    images: record.images
  })
}

const autoClusterModalRef = ref(null)
const handleOpenClusterModal = () => {
  autoClusterModalRef.value.show({
    unknown_summary_model_config_id: robotInfo.value.unknown_summary_model_config_id,
    unknown_summary_use_model: robotInfo.value.unknown_summary_use_model,
    unknown_summary_similarity: robotInfo.value.unknown_summary_similarity,
    unknown_summary_status: robotInfo.value.unknown_summary_status
  })
}

const handleChangeSwitch = (val) => {
  if (val) {
    autoClusterModalRef.value.show({
      unknown_summary_model_config_id: robotInfo.value.unknown_summary_model_config_id,
      unknown_summary_use_model: robotInfo.value.unknown_summary_use_model,
      unknown_summary_similarity: robotInfo.value.unknown_summary_similarity,
      unknown_summary_status: 1
    })
  } else {
    setUnknownIssueSummary({
      id: query.id,
      unknown_summary_model_config_id: robotInfo.value.unknown_summary_model_config_id,
      unknown_summary_use_model: robotInfo.value.unknown_summary_use_model,
      unknown_summary_similarity: robotInfo.value.unknown_summary_similarity,
      unknown_summary_status: 0
    }).then((res) => {
      message.success(t('save_success'))
      getRobot(query.id)
    })
  }
}

const setAnswerModalRef = ref(null)
const handleOpenAnswerModal = (record) => {
  setAnswerModalRef.value.show({
    id: record.id,
    question: record.question,
    answer: record.answer,
    unknown_list: record.unknown_list.join('\n'),
    images: record.images
  })
}

const handleSaveCluster = () => {
  getRobot(query.id)
}

const contextVisible = ref(false)
const contextLoading = ref(false)
const contextMessages = ref([])
const contextRobotInfo = ref({})
const contextAnchorId = ref(0)
const messageListRef = ref(null)

const parseJsonString = (val) => {
  if (!val || typeof val !== 'string') {
    return val
  }
  try {
    return JSON.parse(val)
  } catch (_e) {
    return val
  }
}

const tryParseJson = (val) => {
  if (typeof val !== 'string') {
    return val
  }
  try {
    return JSON.parse(val)
  } catch (_e) {
    return val
  }
}

const extractMessageText = (node) => {
  if (node === null || node === undefined) {
    return []
  }
  if (typeof node === 'string') {
    return [node]
  }
  if (Array.isArray(node)) {
    return node.flatMap((item) => extractMessageText(item))
  }
  if (typeof node === 'object') {
    const itemType = String(node?.type || node?.reply_type || '')
    if (itemType === 'text') {
      return [String(node?.text || node?.content || node?.description || '')]
    }
    if (itemType === 'image') {
      return ['[image]']
    }
    if (itemType === 'voice') {
      return ['[voice]']
    }
    const direct = [node?.text, node?.content, node?.description].filter((x) => typeof x === 'string')
    if (direct.length > 0) {
      return direct
    }
  }
  return []
}

const normalizeReceivedText = (content) => {
  if (!content || typeof content !== 'string') {
    return content
  }
  let parsed = tryParseJson(content)
  if (typeof parsed === 'string') {
    parsed = tryParseJson(parsed)
  }
  if (typeof parsed === 'string') {
    return content
  }
  const textList = extractMessageText(parsed)
    .map((item) => String(item || '').trim())
    .filter((item) => item !== '')
  return textList.length > 0 ? textList.join('\n') : content
}

const normalizeContextMessage = (item, robot, customer) => {
  const msg = {
    ...item,
    uid: item.id,
    loading: false,
    name: item.name || item.nickname
  }
  if (msg.is_customer == 1) {
    msg.name = msg.name || customer?.name || ''
    msg.avatar = msg.avatar || customer?.avatar || ''
  } else {
    msg.name = msg.name || robot?.robot_name || ''
    msg.robot_avatar = msg.avatar || robot?.robot_avatar || ''
    msg.avatar = msg.robot_avatar
  }
  msg.menu_json = parseJsonString(msg.menu_json)
  msg.quote_file = parseJsonString(msg.quote_file)
  msg.reply_content_list = parseJsonString(msg.reply_content_list) || []
  if (msg.is_customer == 1) {
    msg.content = normalizeReceivedText(msg.content)
  }
  msg.voice_content = extractVoiceInfo(msg.content)
  msg.content = removeVoiceFormat(msg.content)
  return msg
}

const buildClusterQuestionList = (record) => {
  const list = []
  if (record.question) {
    list.push(record.question)
  }
  if (Array.isArray(record.unknown_list)) {
    record.unknown_list.forEach((item) => {
      if (item && item !== record.question) {
        list.push(item)
      }
    })
  }
  return list
}

const handleOpenChatContext = async (record, question) => {
  contextLoading.value = true
  try {
    const res = await getUnknownIssueChatContext({
      robot_id: query.id,
      trigger_day: record.trigger_day,
      question: question
    })
    const data = res.data || {}
    const list = Array.isArray(data.list) ? data.list : []
    contextRobotInfo.value = data.robot || {}
    contextAnchorId.value = Number(data.anchor_message_id || 0)
    contextMessages.value = list
      .sort((a, b) => Number(a.id || 0) - Number(b.id || 0))
      .map((item) => normalizeContextMessage(item, data.robot, data.customer))
    if (contextMessages.value.length === 0) {
      message.warning(t('context_not_found'))
      return
    }
    contextVisible.value = true
    await nextTick()
    if (messageListRef.value && contextAnchorId.value > 0) {
      messageListRef.value.scrollToMessage(contextAnchorId.value, 'top', 72)
    } else if (messageListRef.value) {
      messageListRef.value.scrollToBottom()
    }
  } finally {
    contextLoading.value = false
  }
}

const handleChangeTab = () => {
  router.push({
    path: '/robot/config/unknown_issue',
    query: {
      ...query
    }
  })
}
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  padding: 16px 24px;
  .search-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 16px;
    .left-block {
      display: flex;
      align-items: center;
      gap: 4px;
      .anticon {
        cursor: pointer;
      }
    }
    .right-block {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
  .list-box {
    margin-top: 16px;
  }
  ::v-deep(.ant-alert) {
    align-items: baseline;
  }
}
.import-td-box {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  .status-block {
    display: flex;
    align-items: center;
    gap: 8px;
    &.success {
      color: #00a854;
    }
    &.none {
      color: #8c8c8c;
      span {
        background-color: #8c8c8c;
      }
    }
    span {
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background-color: #00a854;
    }
  }
}
.fragment-img {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
  padding-left: 40px;
  img {
    width: 80px;
    height: 80px;
    border-radius: 6px;
    cursor: pointer;
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
  }
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
.classify-scroll-box {
  max-height: 400px;
  min-height: 180px;
  margin-top: 4px;
  overflow: hidden;
  overflow-y: auto;
  .list-item {
    font-size: 14px;
    line-height: 22px;
    color: #262626;
    margin-bottom: 4px;
  }
  /* 整个页面的滚动条 */
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
