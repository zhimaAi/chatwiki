<template>
  <div>
    <a-modal v-model:open="open" :title="t('title_auth_param')" :width="700">
      <div class="mt24">
        <a-alert
          class="zm-alert-info"
          :message="t('msg_auth_param_tip')"
          type="info"
        />
      </div>
      <div class="add-btn-block">
        <a-button type="primary" :icon="h(PlusOutlined)" ghost @click="handleAddKey"
          >{{ t('btn_add_param') }}</a-button
        >
      </div>
      <div class="list-box">
        <a-table
          :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
          :columns="columns"
          :pagination="false"
          :data-source="tableData"
          row-key="uni_key"
          :scroll="{ y: 550 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'auth_key'">
              <a-input v-model:value="record.auth_key" :placeholder="t('ph_input')" />
            </template>
            <template v-if="column.dataIndex === 'auth_value'">
              <a-input v-model:value="record.auth_value" :placeholder="t('ph_input')" />
            </template>
            <template v-if="column.dataIndex === 'auth_value_addto'">
              <a-select v-model:value="record.auth_value_addto" style="width: 100%">
                <a-select-option value="HEADERS">HEADERS</a-select-option>
                <a-select-option value="PARAMS">PARAMS</a-select-option>
                <a-select-option value="BODY">BODY</a-select-option>
              </a-select>
            </template>
            <template v-if="column.dataIndex === 'auth_remark'">
              <a-input :maxLength="10" v-model:value="record.auth_remark" :placeholder="t('ph_input')" />
            </template>
            <template v-if="column.dataIndex === 'action'">
              <a-popconfirm
                :title="t('msg_confirm_delete')"
                :ok-text="t('btn_confirm')"
                :cancel-text="t('btn_cancel')"
                @confirm="handleDel(record)"
              >
                <a>{{ t('label_delete') }}</a>
              </a-popconfirm>
            </template>
          </template>
        </a-table>
      </div>

      <template #footer>
        <a-button @click="open = false">{{ t('btn_cancel') }}</a-button>
        <a-button @click="handleSave">{{ t('btn_only_save') }}</a-button>
        <a-button type="primary" @click="handleOk">{{ t('btn_save_and_add') }}</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { getUuid } from '@/utils/index'
import { ref, h, reactive } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getHttpAuthConfig, saveHttpAuthConfig } from '@/api/robot/index'
import { message } from 'ant-design-vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.http-node.add-authentication-modal')
const open = ref(false)

const tableData = ref([])

const emit = defineEmits(['ok'])

const columns = [
  {
    title: t('label_key'),
    dataIndex: 'auth_key'
  },
  {
    title: t('label_value'),
    dataIndex: 'auth_value'
  },
  {
    title: t('label_add_to'),
    dataIndex: 'auth_value_addto'
  },
  {
    title: t('label_remark'),
    dataIndex: 'auth_remark'
  },
  {
    title: t('label_operation'),
    dataIndex: 'action',
    width: 80
  }
]

const state = reactive({
  selectedRowKeys: []
})

const onSelectChange = (selectedRowKeys) => {
  console.log('selectedRowKeys changed: ', selectedRowKeys)
  state.selectedRowKeys = selectedRowKeys
}
const show = () => {
  open.value = true
  state.selectedRowKeys = []
  getConfigList()
}

const handleAddKey = () => {
  if (tableData.value.length >= 100) {
    return message.error(t('msg_max_100'))
  }
  tableData.value.push({
    auth_key: '',
    auth_value: '',
    auth_value_addto: 'HEADERS',
    auth_remark: '',
    uni_key: getUuid(16)
  })
}

const handleDel = (record) => {
  tableData.value = tableData.value.filter((item) => item.uni_key !== record.uni_key)
}

const getConfigList = () => {
  getHttpAuthConfig().then((res) => {
    let list = res.data || []
    tableData.value = list.map((item) => {
      return {
        ...item,
        uni_key: getUuid(16)
      }
    })
  })
}
function getParmas() {
  let errorIndex = []
  let resultList = []

  tableData.value.forEach((item, index) => {
    resultList.push({
      auth_key: item.auth_key,
      auth_value: item.auth_value,
      auth_value_addto: item.auth_value_addto,
      auth_remark: item.auth_remark
    })
    if (item.auth_key == '' || item.auth_value == '') {
      errorIndex.push(index + 1)
    }
  })

  if (errorIndex.length > 0) {
    message.error(t('msg_fill_params', { index: errorIndex.join(',') }))
    return false
  }
  if (resultList.length == 0) {
    message.error(t('msg_add_param'))
    return false
  }

  return resultList
}
const handleOk = () => {
  let http_auth_config_list = getParmas()
  if (http_auth_config_list == false) {
    return
  }
  let list = []
  state.selectedRowKeys.forEach((uni_key) => {
    list.push(tableData.value.find((it) => it.uni_key == uni_key))
  })
  emit('ok', list)
  saveHttpAuthConfig({
    http_auth_config_list: JSON.stringify(http_auth_config_list)
  })
    .then((res) => {
      message.success(t('msg_save_success'))
    })
    .finally(() => {
      state.selectedRowKeys = []
      open.value = false
    })
}

const handleSave = () => {
  let http_auth_config_list = getParmas()
  if (http_auth_config_list == false) {
    return
  }
  saveHttpAuthConfig({
    http_auth_config_list: JSON.stringify(http_auth_config_list)
  }).then((res) => {
    message.success(t('msg_save_success'))
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.mt24 {
  margin-top: 24px;
}
.add-btn-block {
  margin-top: 16px;
  margin-bottom: 8px;
}

.list-box {
  margin-bottom: 24px;
}
</style>
