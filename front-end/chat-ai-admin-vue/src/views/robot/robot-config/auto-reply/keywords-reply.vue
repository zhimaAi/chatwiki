<template>
  <div class="user-model-page">
    <!-- 自动回复 开关 -->
    <div class="switch-block">
      <span class="switch-title">{{ t('auto_reply') }}</span>
      <a-switch
        @change="keyWordReplySwitchChange"
        :checked="keywordReplyStatus"
        :checked-children="t('switch_on')"
        :un-checked-children="t('switch_off')"
      />
      <span class="switch-desc">{{ t('auto_reply_desc') }}</span>
    </div>
    <a-alert show-icon>
      <template #message>
        <p>{{ t('keyword_reply_tip') }}</p>
      </template>
    </a-alert>
    <div class="search-block">
      <div class="left-block">
        <a-button type="primary" @click="handleAddRule">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('btn_add_rule') }}
        </a-button>
        <!-- 回复类型：下拉选择 图文  文本  图片  小程序 和链接 -->
        <div class="search-item">
          <a-select
            v-model="reply_type"
            :placeholder="t('label_reply_type')"
            allowClear
            :options="replyTypeOptions"
            style="width: 240px;"
            @change="onReplyTypeChange"
          />
        </div>

        <div class="search-item">
          <a-input-search
            v-model:value="search_keyword"
            :placeholder="t('ph_search_keyword')"
            allowClear
            style="width: 240px;"
            @search="onSearch"
          >
            <!-- <template #suffix>
              <SearchOutlined />
            </template> -->
          </a-input-search>
        </div>
      </div>
      <div class="right-block">
        <span>{{ t('continue_reply') }}</span>
        <a-tooltip>
          <template #title>{{ t('continue_reply_tip') }}</template>
          <QuestionCircleOutlined />：
        </a-tooltip>
        <a-switch
          @change="handleChangeSwitch"
          :checked="keywordReplyAiReplyStatus"
          :checked-children="t('switch_on')"
          :un-checked-children="t('switch_off')"
        />
      </div>
    </div>
    <div class="list-box">
      <a-table
        :columns="columns"
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
        <template #bodyCell="{ column, record }">
          <!-- 模糊匹配  -->
          <template v-if="column.key === 'half_keyword'">
            <div class="tags-wra flex">
              <div
                class="tag"
                v-for="(item, index) in record.half_keyword"
                :key="index"
                v-show="index < 3"
              >
                <a-tag>{{ item }}</a-tag>
              </div>
              <a-popover>
                <template #content>
                  <div class="popover-cont flex">
                    <div
                      class="tag"
                      v-for="(item, index) in record.half_keyword"
                      :key="index"
                      v-show="index >= 3"
                    >
                      <a-tag>{{ item }}</a-tag>
                    </div>
                  </div>
                </template>
                <div class="more-tag" v-if="record.half_keyword.length > 3"
                  >+{{ record.half_keyword.length - 3 }}</div
                >
              </a-popover>
            </div>
          </template>
          <!-- 精准匹配  -->
          <template v-if="column.key === 'full_keyword'">
            <div class="tags-wra flex">
              <div
                class="tag"
                v-for="(item, index) in record.full_keyword"
                :key="index"
                v-show="index < 3"
              >
                <a-tag>{{ item }}</a-tag>
              </div>
              <a-popover>
                <template #content>
                  <div class="popover-cont flex">
                    <div
                      class="tag"
                      v-for="(item, index) in record.full_keyword"
                      :key="index"
                      v-show="index >= 3"
                    >
                      <a-tag>{{ item }}</a-tag>
                    </div>
                  </div>
                </template>
                <div class="more-tag" v-if="record.full_keyword.length > 3"
                  >+{{ record.full_keyword.length - 3 }}</div
                >
              </a-popover>
            </div>
          </template>
          <!-- 回复内容 -->
          <template v-if="column.key === 'reply_content'">
            <span style="color:#595959;">{{ summarizeReplyTypes(record.reply_content) || '--' }}</span>
          </template>

          <!-- 启用状态 开关-->
          <template v-if="column.key === 'switch_status'">
            <a-switch
              :checked="record.switch_status"
              :checkedValue="'1'"
              :un-checkedValue="'0'"
              :checked-children="t('switch_on')"
              :un-checked-children="t('switch_off')"
              @change="handleReplySwitchChange(record, $event)"
            />
          </template>


          <!-- 操作 -->
          <template v-if="column.key === 'action'">
            <a-flex :gap="8">
              <a @click="handleEdit(record)">{{ t('btn_edit') }}</a>
              <a @click="handleDelete(record)">{{ t('btn_delete') }}</a>
              <a @click="handleCopy(record)">{{ t('btn_copy') }}</a>
            </a-flex>
          </template>
        </template>
      </a-table>
    </div>
  </div>
  
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { saveRobotAbilitySwitchStatus, getRobotKeywordReplyList, saveRobotAbilityAiReplyStatus, updateRobotKeywordReplySwitchStatus, deleteRobotKeywordReply } from '@/api/explore/index.js'
import { REPLY_TYPE_OPTIONS, REPLY_TYPE_LABEL_MAP } from '@/constants/index'
import { useRobotStore } from '@/stores/modules/robot'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.auto-reply.keywords-reply')
const robotStore = useRobotStore()

// 来自左侧菜单的能力开关（关键词回复）
const keywordReplyStatus = computed(() => robotStore.keywordReplySwitchStatus === '1')
const keywordReplyAiReplyStatus = computed(() => robotStore.keywordReplyAiReplyStatus === '1')

const query = useRoute().query
const router = useRouter()

const columns = ref([
  {
    title: t('col_rule_name'),
    dataIndex: 'name',
    key: 'name',
    width: 120
  },
  {
    title: t('col_fuzzy_match'),
    dataIndex: 'half_keyword',
    key: 'half_keyword',
    width: 220
  },
  {
    title: t('col_exact_match'),
    dataIndex: 'full_keyword',
    key: 'full_keyword',
    width: 220
  },
  {
    title: t('col_reply_content'),
    dataIndex: 'reply_content',
    key: 'reply_content',
    width: 120
  },
  {
    title: t('col_enabled_status'),
    dataIndex: 'switch_status',
    key: 'switch_status',
    width: 120
  },
  {
    title: t('col_action'),
    dataIndex: 'action',
    key: 'action',
    width: 120
  }
])

const pager = reactive({
  page: 1,
  size: 10,
  total: 0
})
const replyTypeOptions = REPLY_TYPE_OPTIONS()
const tableData = ref([])
const loading = ref(false)
const reply_type = ref('')
const search_keyword = ref('')
const getTableData = () => {
  console.log('search_keyword.value', search_keyword.value)
  console.log('reply_type.value', reply_type.value)
  const parmas = {
    robot_id: query.id,
    keyword: search_keyword.value || '',
    reply_type: reply_type.value || '',
    page: pager.page,
    size: pager.size
  }
  loading.value = true
  getRobotKeywordReplyList({
    ...parmas
  })
    .then((res) => {
      const data = res?.data || { list: [], total: 0, page: pager.page, size: pager.size }
      tableData.value = (data.list || []).map((item) => ({
        ...item,
        // 确保这些字段存在
        full_keyword: Array.isArray(item.full_keyword) ? item.full_keyword : [],
        half_keyword: Array.isArray(item.half_keyword) ? item.half_keyword : [],
        reply_content: Array.isArray(item.reply_content) ? item.reply_content : [],
        switch_status: String(item.switch_status ?? '0')
      }))
      pager.total = +data.total || 0
    })
    .finally(() => {
      loading.value = false
    })
}
getTableData()

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getTableData()
}

const onSearch = () => {
  pager.page = 1
  getTableData()
}

const onReplyTypeChange = (val) => {
  reply_type.value = val
  onSearch()
}

const handleAddRule = () => {
  router.push({
    path: '/robot/ability/auto-reply/add-rule',
    query: {
      id: query.id,
      robot_key: query.robot_key
    }
  })
}

  const handleEdit = (record) => {
    router.push({
      path: '/robot/ability/auto-reply/add-rule',
      query: {
        id: query.id,
        robot_key: query.robot_key,
        rule_id: record.id
      }
    })
  }

  const handleCopy = (record) => {
    router.push({
      path: '/robot/ability/auto-reply/add-rule',
      query: {
        id: query.id,
        robot_key: query.robot_key,
        copy_id: record.id
      }
    })
  }

const keyWordReplySwitchChange = (checked) => {
  const switch_status = checked ? '1' : '0'
  if (switch_status === '0') {
    Modal.confirm({
      title: t('title_tip'),
      content: t('msg_close_warning'),
      onOk: () => {
        saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_auto_reply', switch_status }).then((res) => {
          if (res && res.res == 0) {
            robotStore.setKeywordReplySwitchStatus(switch_status)
            message.success(t('msg_operation_success'))
            window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
          }
        })
      }
    })
    return
  }
  saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_auto_reply', switch_status }).then((res) => {
    if (res && res.res == 0) {
      robotStore.setKeywordReplySwitchStatus(switch_status)
      message.success(t('msg_operation_success'))
      window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
    }
  })
}

const handleReplySwitchChange = (record, checked) => {
  const switch_status = checked
  updateRobotKeywordReplySwitchStatus({ id: record.id, robot_id: query.id, switch_status }).then((res) => {
    if (res && res.res == 0) {
      record.switch_status = switch_status
      message.success(t('msg_operation_success'))
    }
  })
}

const handleChangeSwitch = (checked) => {
  const ai_reply_status = checked ? '1' : '0'
  saveRobotAbilityAiReplyStatus({ robot_id: query.id, ability_type: 'robot_auto_reply', ai_reply_status }).then((res) => {
    if (res && res.res == 0) {
      robotStore.setKeywordReplyAiReplyStatus(ai_reply_status)
      message.success(t('msg_operation_success'))
    }
  })
}

const handleDelete = (record) => {
  // 确认删除
  Modal.confirm({
    title: t('title_confirm_delete'),
    okText: t('btn_confirm'),
    onOk: () => {
      deleteRobotKeywordReply({ id: record.id, robot_id: query.id }).then((res) => {
        if (res && res.res == 0) {
          message.success(t('msg_delete_success'))
          getTableData()
        }
      })
    }
  })
}

function mapReplyTypeLabel (t) {
  return REPLY_TYPE_LABEL_MAP()[t] || ''
}

function summarizeReplyTypes (list) {
  if (!Array.isArray(list)) return ''
  const labels = list
    .map((rc) => mapReplyTypeLabel(rc?.type))
    .filter((s) => !!s)
  // 去重并使用/连接
  const uniq = [...new Set(labels)]
  return uniq.join('/')
}
</script>

<style lang="less" scoped>
.user-model-page {
  height: 100%;
  width: 100%;
  .search-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 16px;
    .left-block {
      display: flex;
      align-items: center;
      gap: 16px;
    }
    .right-block {
      display: flex;
      align-items: center;
      gap: 2px;
    }
  }
  .list-box {
    margin-top: 8px;
  }
  ::v-deep(.ant-alert) {
    align-items: baseline;
  }
}

.switch-block {
  display: flex;
  align-items: center;
  margin-bottom: 16px;

  .switch-title {
    margin-right: 12px;
    color: #262626;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
}
.switch-desc {
  margin-left: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.tags-wra {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;

  .ant-tag {
    margin: 0;
  }
}

.more-tag {
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
  cursor: pointer;
}

.flex {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.popover-cont {
  max-width: 560px;
}
</style>
