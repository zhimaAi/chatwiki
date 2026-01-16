<template>
  <div>
    <a-modal v-model:open="open" title="鉴权参数" :width="700">
      <div class="mt24">
        <a-alert
          class="zm-alert-info"
          message="选择的参数会自动填充到http节点，鉴权参数在导出CSL或上架模板时会自动清空"
          type="info"
        />
      </div>
      <div class="add-btn-block">
        <a-button type="primary" :icon="h(PlusOutlined)" ghost @click="handleAddKey"
          >添加参数</a-button
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
              <a-input v-model:value="record.auth_key" placeholder="请输入" />
            </template>
            <template v-if="column.dataIndex === 'auth_value'">
              <a-input v-model:value="record.auth_value" placeholder="请输入" />
            </template>
            <template v-if="column.dataIndex === 'auth_value_addto'">
              <a-select v-model:value="record.auth_value_addto" style="width: 100%">
                <a-select-option value="HEADERS">HEADERS</a-select-option>
                <a-select-option value="PARAMS">PARAMS</a-select-option>
                <a-select-option value="BODY">BODY</a-select-option>
              </a-select>
            </template>
            <template v-if="column.dataIndex === 'auth_remark'">
              <a-input :maxLength="10" v-model:value="record.auth_remark" placeholder="请输入" />
            </template>
            <template v-if="column.dataIndex === 'action'">
              <a-popconfirm
                title="确认删除该字段?"
                ok-text="确定"
                cancel-text="取消"
                @confirm="handleDel(record)"
              >
                <a>删除</a>
              </a-popconfirm>
            </template>
          </template>
        </a-table>
      </div>

      <template #footer>
        <a-button @click="open = false">取消</a-button>
        <a-button @click="handleSave">仅保存参数</a-button>
        <a-button type="primary" @click="handleOk">保存并添加</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { getUuid } from '@/utils/index'
import { ref, h, reactive } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getHttpAuthConfig, saveHttpAuthConfig } from '@/api/robot/index'
import { message } from 'ant-design-vue'
const open = ref(false)

const tableData = ref([])

const emit = defineEmits(['ok'])

const columns = [
  {
    title: 'Key',
    dataIndex: 'auth_key'
  },
  {
    title: 'Value',
    dataIndex: 'auth_value'
  },
  {
    title: 'Add To',
    dataIndex: 'auth_value_addto'
  },
  {
    title: '备注',
    dataIndex: 'auth_remark'
  },
  {
    title: '操作',
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
    return message.error('最多添加100条')
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
    message.error(`请填写第${errorIndex.join(',')}行参数`)
    return false
  }
  if (resultList.length == 0) {
    message.error('请添加参数')
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
      message.success('保存成功')
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
    message.success('保存成功')
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
