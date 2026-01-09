<template>
  <div class="content-box">
    <div class="search-block">
      <a-input
        v-model:value="searchState.search"
        :placeholder="t('quota.search_placeholder')"
        style="width: 220px"
        @change="onSearch"
      >
        <template #suffix>
          <SearchOutlined @click="onSearch" />
        </template>
      </a-input>
    </div>
    <div class="list-box">
      <a-table :data-source="tableData" :loading="loading" :pagination="false" sticky>
        <a-table-column :title="t('quota.type')" data-index="token_app_type" :width="120">
          <template #default="{ record }">
            <a-flex :gap="4" align="center">
              <a-avatar
                v-if="record.token_app_type == 'chatwiki_robot'"
                :size="24"
                :src="DEFAULT_ROBOT_AVATAR"
              />
              <a-avatar
                v-if="record.token_app_type == 'workflow'"
                :size="24"
                :src="DEFAULT_WORKFLOW_AVATAR"
              />
              <a-avatar v-if="record.token_app_type == 'other'" :size="24">
                <template #icon><AppstoreOutlined style="font-size: 16px" /></template>
              </a-avatar>
              {{ token_app_type_map[record.token_app_type] }}
            </a-flex>
          </template>
        </a-table-column>
        <a-table-column :title="t('quota.app_name')" data-index="robot_name" :width="160"> </a-table-column>
        <a-table-column :title="t('quota.today_consumption')" data-index="today_use_token" :width="130">
          <template #default="{ record }">{{ formatNum(record.today_use_token) }}</template>
        </a-table-column>
        <a-table-column :title="t('quota.token_limit')" data-index="amount" :width="120">
          <template #default="{ record }">{{
            record.is_config == 1 && record.switch_status == 1 ? formatNum(record.max_token) : '-'
          }}</template>
        </a-table-column>
        <a-table-column :title="t('quota.consumed')" data-index="amount2" :width="120">
          <template #default="{ record }">{{
            record.switch_status == 1 ? formatNum(record.use_token) : '-'
          }}</template>
        </a-table-column>
        <a-table-column :title="t('quota.remaining')" data-index="amount3" :width="120">
          <template #default="{ record }">{{
            record.switch_status == 1 ? formatNum(record.max_token - record.use_token) : '-'
          }}</template>
        </a-table-column>
        <a-table-column :title="t('quota.remark')" data-index="description" :width="200" :ellipsis="true">
          <template #default="{ record }">{{ record.description || '--' }}</template>
        </a-table-column>
        <a-table-column :title="t('quota.action')" data-index="action" :width="120">
          <template #default="{ record }">
            <a-flex align="center" :gap="8">
              <a-switch
                @change="handleChangeSwitch(record)"
                :checked="record.switch_status == 1"
                :checked-children="t('quota.switch_on')"
                :un-checked-children="t('quota.switch_off')"
              />
              <a v-if="record.switch_status == 1" @click="handleOpenModal(record, true)">{{ t('quota.setting') }}</a>
            </a-flex>
          </template>
        </a-table-column>
      </a-table>
    </div>
    <TokenQuotaModal ref="tokenQuotaModalRef" @ok="handleSetOk" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { SearchOutlined, AppstoreOutlined } from '@ant-design/icons-vue'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_WORKFLOW_AVATAR } from '@/constants/index.js'
import TokenQuotaModal from './components/token-quota-modal.vue'
import { tokenLimitList, tokenLimitCreate, tokenLimitSwitch } from '@/api/manage/index.js'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.usetoken')

const searchState = reactive({
  search: ''
})

const tableData = ref([])

const token_app_type_map = {
  chatwiki_robot: t('robot'),
  workflow: t('workflow'),
  other: t('other')
}

const loading = ref(false)
const onSearch = () => {
  getData()
}

const getData = () => {
  loading.value = true
  tokenLimitList({
    ...searchState
  })
    .then((res) => {
      tableData.value = res.data.list || []
    })
    .finally(() => {
      loading.value = false
    })
}

let tempRecord = null
const handleChangeSwitch = (record) => {
  tempRecord = record
  if (tempRecord.switch_status == 1) {
    handleSetSwitch()
    return
  }
  handleOpenModal(record)
}

const handleSetOk = () => {
  if (tempRecord == null) {
    getData()
  } else {
    handleSetSwitch()
  }
}
const handleSetSwitch = () => {
  tokenLimitSwitch({
    robot_id: tempRecord.robot_id,
    token_app_type: tempRecord.token_app_type,
    switch_status: tempRecord.switch_status == 1 ? 0 : 1
  }).then((res) => {
    message.success(`${tempRecord.switch_status == 1 ? t('quota.close_success') : t('quota.open_success')}`)
    getData()
  })
}

function formatNum(num) {
  if (num <= 0 || !num) {
    return 0
  }
  return (num / 1000).toFixed(3)
}

const tokenQuotaModalRef = ref(null)
const handleOpenModal = (record, isRowEdit) => {
  if (isRowEdit) {
    tempRecord = null
  }
  tokenQuotaModalRef.value.show(JSON.parse(JSON.stringify(record)))
}

onMounted(() => {
  getData()
})
</script>

<style lang="less" scoped>
.content-box {
  padding: 0 16px;
  .search-block {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .list-box {
    margin-top: 16px;
  }
}
</style>
