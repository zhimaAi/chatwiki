<template>
  <a-modal
    class="add-model-alert"
    width="800px"
    v-model:open="show"
    :title="currentTitle"
    @ok="handleOk"
    @cancel="handleClose"
  >
    <div class="form-wrapper" v-if="object_type == 1">
      <div class="select">已选择({{ checkedList.length }})</div>
      <a-table
        :data-source="dataList"
        :pagination="false"
        :scroll="{ y: 500 }"
        row-key="id"
        :row-selection="{ selectedRowKeys: checkedList, onChange: onSelectChange }"
      >
        <a-table-column title="机器人" data-index="robot_name" :width="350">
          <template #default="{ record }">
            <div class="tools-wrapper">
              <img class="model-logo" :src="record.robot_avatar" alt="" />
              <div class="item-name">{{ record.robot_name }}</div>
            </div>
          </template>
        </a-table-column>
        <a-table-column title="权限" data-index="operate_rights" :width="350">
          <template #default="{ record }">
            <a-radio-group v-model:value="record.operate_rights">
              <a-radio value="4">管理</a-radio>
              <a-radio value="2">编辑</a-radio>
              <a-radio value="1">查看</a-radio>
            </a-radio-group>
          </template>
        </a-table-column>
      </a-table>
    </div>
    <div class="form-wrapper" v-else-if="object_type == 2">
      <div class="main-tabs-box">
        <a-segmented v-model:value="tabs" :options="tabOption" @change="handleChangeTabs">
          <template #label="{ payload }">
            {{ payload.title }}
          </template>
        </a-segmented>
      </div>
      <template v-if="tabs == 1">
        <div class="select">已选择({{ checkedList.length }})</div>
        <a-table
          :data-source="dataList"
          :pagination="false"
          :scroll="{ y: 500 }"
          row-key="id"
          :row-selection="{ selectedRowKeys: checkedList, onChange: onSelectChange }"
        >
          <a-table-column title="知识库" data-index="library_name" :width="350">
            <template #default="{ record }">
              <div class="tools-wrapper">
                <div class="item-name">{{ record.library_name }}</div>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="权限" data-index="operate_rights" :width="350">
            <template #default="{ record }">
              <a-radio-group v-model:value="record.operate_rights">
                <a-radio value="4">管理</a-radio>
                <a-radio value="2">编辑</a-radio>
                <a-radio value="1">查看</a-radio>
              </a-radio-group>
            </template>
          </a-table-column>
        </a-table>
      </template>
      <template v-else>
        <a-table :data-source="departLists" :pagination="false" :scroll="{ y: 500 }" row-key="id">
          <a-table-column title="知识库" data-index="name" :width="350">
            <template #default="{ record }">
              <div class="tools-wrapper">
                <div class="item-name">{{ record.name }}</div>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="权限" data-index="operate_rights" :width="350">
            <template #default="{ record }">
              <a-radio-group v-model:value="record.operate_rights" disabled>
                <a-radio value="4">管理</a-radio>
                <a-radio value="2">编辑</a-radio>
                <a-radio value="1">查看</a-radio>
              </a-radio-group>
            </template>
          </a-table-column>
        </a-table>
      </template>
    </div>
    <div class="form-wrapper" v-else-if="object_type == 3">
      <div class="select">已选择({{ checkedList.length }})</div>

      <a-table
        :data-source="dataList"
        :pagination="false"
        :scroll="{ y: 500 }"
        row-key="id"
        :row-selection="{ selectedRowKeys: checkedList, onChange: onSelectChange }"
      >
        <a-table-column title="数据库" data-index="name" :width="350">
          <template #default="{ record }">
            <div class="tools-wrapper">
              <div class="item-name">{{ record.name }}</div>
            </div>
          </template>
        </a-table-column>
        <a-table-column title="权限" data-index="operate_rights" :width="350">
          <template #default="{ record }">
            <a-radio-group v-model:value="record.operate_rights">
              <a-radio value="4">管理</a-radio>
              <a-radio value="2">编辑</a-radio>
              <a-radio value="1">查看</a-radio>
            </a-radio-group>
          </template>
        </a-table-column>
      </a-table>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { getPermissionManageList, batchSavePermissionManage } from '@/api/department/index.js'
import { message } from 'ant-design-vue'
const emit = defineEmits(['save', 'ok'])
const props = defineProps({
  robotList: {
    type: Array,
    default: () => []
  },
  libraryList: {
    type: Array,
    default: () => []
  },
  formList: {
    type: Array,
    default: () => []
  }
})
const checkedList = ref([])
const currentTitle = ref('')
const show = ref(false)
const dataList = ref([])
const departLists = ref([])

const formatId = (arr) => {
  const newArr = []
  arr.map((item) => {
    newArr.push(item.id)
  })
  return newArr
}

const onSelectChange = (selectedRowKeys) => {
  checkedList.value = selectedRowKeys
}

const tabs = ref(1)
const baseOption = [
  {
    value: 2,
    payload: {
      title: '部门知识库'
    }
  },
  {
    value: 1,
    payload: {
      title: '成员知识库'
    }
  }
]
const tabOption = ref(baseOption)

let currentRecord = null
const object_type = ref(1) // 1:机器人 2:知识库 3:数据库
const identity_type = ref(1) //	1:用户 2:部门
const identity_id = ref()

let depratmentIds = []
const open = async (key, status, record) => {
  checkedList.value = []
  object_type.value = key
  depratmentIds = []
  identity_type.value = 1
  tabs.value = 1
  currentRecord = JSON.parse(JSON.stringify(record))
  identity_id.value = currentRecord.id
  record = JSON.parse(JSON.stringify(record))
  tabOption.value = baseOption
  record.managed_robot_list = record.managed_robot_list || []
  record.managed_library_list = record.managed_library_list || []
  record.managed_form_list = record.managed_form_list || []
  let list = []
  if (key == 1) {
    currentTitle.value = '添加成员管理的机器人'
    list = props.robotList
    list.forEach((item) => {
      let filterItem = record.managed_robot_list.filter((it) => it.id == item.id)
      if (filterItem.length) {
        item.operate_rights = filterItem[0].operate_rights
      } else {
        item.operate_rights = '4'
      }
    })
    // 编辑
    if (status && status == 'edit') {
      currentTitle.value = '编辑成员管理的机器人'
      checkedList.value = formatId(record.managed_robot_list)
    }
  } else if (key == 2) {
    console.log(record, '==')
    currentTitle.value = '添加成员管理的知识库'
    // if (record.departments && record.departments.length) {
    //   depratmentIds = record.departments.map((item) => item.id)
    //   getDepartmentList()
    // }
    list = props.libraryList
    list.forEach((item) => {
      let filterItem = record.managed_library_list.filter((it) => it.id == item.id)
      if (filterItem.length) {
        item.operate_rights = filterItem[0].operate_rights
      } else {
        item.operate_rights = '4'
      }
    })
    // 编辑
    if (status && status == 'edit') {
      currentTitle.value = '编辑成员管理的知识库'
      checkedList.value = formatId(record.managed_library_list)
    }
  } else if (key == 3) {
    currentTitle.value = '添加成员管理的数据库'

    list = props.formList
    list.forEach((item) => {
      let filterItem = record.managed_form_list.filter((it) => it.id == item.id)
      if (filterItem.length) {
        item.operate_rights = filterItem[0].operate_rights
      } else {
        item.operate_rights = '4'
      }
    })

    // 编辑
    if (status && status == 'edit') {
      currentTitle.value = '编辑成员管理的数据库'
      checkedList.value = formatId(record.managed_form_list)
    }
  }
  dataList.value = list
  show.value = true
}

const getDepartmentList = () => {
  getPermissionManageList({
    object_type: object_type.value,
    identity_type: identity_type.value,
    identity_id: identity_id.value,
    tab: identity_type.value
  }).then((res) => {
    departLists.value = res.data || []
  })
}

const handleChangeTabs = () => {
  if (tabs.value == 2) {
    getDepartmentList()
  }
}

const handleBatchOpen = (key, ids) => {
  checkedList.value = []
  depratmentIds = []
  identity_type.value = 1
  tabs.value = 1
  object_type.value = key
  identity_id.value = ids
  let list = []
  if (key == 1) {
    currentTitle.value = '批量修改成员管理的机器人'
    list = props.robotList
  } else if (key == 2) {
    currentTitle.value = '批量修改成员管理的知识库'
    list = props.libraryList
  } else if (key == 3) {
    currentTitle.value = '批量修改成员管理的数据库'
    list = props.formList
  }
  list.forEach((item) => {
    item.operate_rights = '4'
  })
  dataList.value = list
  show.value = true
}

const handleDepartmentOpen = (key, treeParmas, datas) => {
  // 添加部门
  checkedList.value = []
  depratmentIds = []

  if (treeParmas.is_user) {
    identity_type.value = 1
  } else {
    identity_type.value = 2
  }

  identity_id.value = treeParmas.id

  tabOption.value = [
    {
      value: 2,
      payload: {
        title: identity_type.value == 1 ? '部门知识库' : '上级知识库'
      }
    },
    {
      value: 1,
      payload: {
        title: identity_type.value == 1 ? '成员知识库' : '部门知识库'
      }
    }
  ]
  tabs.value = 1
  object_type.value = key

  let list = []
  if (key == 1) {
    currentTitle.value = `添加${identity_type == 1 ? '成员' : '部门'}管理的机器人`
    list = props.robotList
  } else if (key == 2) {
    currentTitle.value = `添加${identity_type == 1 ? '成员' : '部门'}管理的知识库`
    list = props.libraryList
  } else if (key == 3) {
    currentTitle.value = `添加${identity_type == 1 ? '成员' : '部门'}管理的数据库`
    list = props.formList
  }
  list.forEach((item) => {
    let filterItem = datas.filter((it) => it.object_id == item.id)
    if (filterItem.length) {
      item.operate_rights = filterItem[0].operate_rights
    } else {
      item.operate_rights = '4'
    }
  })
  checkedList.value = datas.map((item) => item.object_id)
  dataList.value = list
  show.value = true
}

const handleOk = () => {
  let parmas = {
    identity_ids: identity_id.value,
    identity_type: identity_type.value,
    object_type: object_type.value
  }

  let object_array = []
  dataList.value.forEach((item) => {
    if (checkedList.value.includes(item.id)) {
      object_array.push({
        object_id: +item.id,
        object_type: +object_type.value,
        operate_rights: +item.operate_rights
      })
    }
  })
  parmas.object_array = JSON.stringify(object_array)
  batchSavePermissionManage(parmas).then((res) => {
    show.value = false
    message.success('保存成功')
    emit('ok')
  })
}

const handleClose = () => {}

// watch(
//   () => checkedList.value,
//   (val) => {
//   },
// );

defineExpose({
  open,
  handleBatchOpen,
  handleDepartmentOpen
})
</script>

<style lang="less" scoped>
.add-model-alert {
  .form-wrapper {
    display: flex;
    flex-direction: column;
    margin-top: 14px;
  }

  .select {
    padding-bottom: 8px;
  }

  .model-logo {
    display: block;
    height: 26px;
  }
  .list-item {
    width: 100%;
    border-bottom: 1px solid #d9d9d9;
  }

  .tools-wrapper {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    gap: 8px;
  }
}

.main-tabs-box {
  display: flex;
  margin-bottom: 16px;
  ::v-deep(.ant-segmented .ant-segmented-item-selected) {
    color: #2475fc;
  }
  ::v-deep(.ant-segmented .ant-segmented-item-label) {
    padding: 0 16px;
  }
}
</style>
