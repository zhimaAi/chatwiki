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
      <div class="select">{{ t('selected_count', { count: checkedList.length }) }}</div>
      <a-table
        :data-source="dataList"
        :pagination="false"
        :scroll="{ y: 500 }"
        row-key="id"
        :row-selection="{ selectedRowKeys: checkedList, onChange: onSelectChange }"
      >
        <a-table-column :title="t('robot')" data-index="robot_name" :width="350">
          <template #default="{ record }">
            <div class="tools-wrapper">
              <img class="model-logo" :src="record.robot_avatar" alt="" />
              <div class="item-name">{{ record.robot_name }}</div>
            </div>
          </template>
        </a-table-column>
        <a-table-column :title="t('permission')" data-index="operate_rights" :width="350">
          <template #default="{ record }">
            <a-radio-group v-model:value="record.operate_rights">
              <a-radio value="4">{{ t('manage') }}</a-radio>
              <a-radio value="2">{{ t('edit') }}</a-radio>
              <a-radio value="1">{{ t('view') }}</a-radio>
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
        <div class="select">{{ t('selected_count', { count: checkedList.length }) }}</div>
        <a-table
          :data-source="dataList"
          :pagination="false"
          :scroll="{ y: 500 }"
          row-key="id"
          :row-selection="{ selectedRowKeys: checkedList, onChange: onSelectChange }"
        >
          <a-table-column :title="t('knowledge_base')" data-index="library_name" :width="350">
            <template #default="{ record }">
              <div class="tools-wrapper">
                <div class="item-name">{{ record.library_name }}</div>
              </div>
            </template>
          </a-table-column>
          <a-table-column :title="t('permission')" data-index="operate_rights" :width="350">
            <template #default="{ record }">
              <a-radio-group v-model:value="record.operate_rights">
                <a-radio value="4">{{ t('manage') }}</a-radio>
                <a-radio value="2">{{ t('edit') }}</a-radio>
                <a-radio value="1">{{ t('view') }}</a-radio>
              </a-radio-group>
            </template>
          </a-table-column>
        </a-table>
      </template>
      <template v-else>
        <a-table :data-source="departLists" :pagination="false" :scroll="{ y: 500 }" row-key="id">
          <a-table-column :title="t('knowledge_base')" data-index="name" :width="350">
            <template #default="{ record }">
              <div class="tools-wrapper">
                <div class="item-name">{{ record.name }}</div>
              </div>
            </template>
          </a-table-column>
          <a-table-column :title="t('permission')" data-index="operate_rights" :width="350">
            <template #default="{ record }">
              <a-radio-group v-model:value="record.operate_rights" disabled>
                <a-radio value="4">{{ t('manage') }}</a-radio>
                <a-radio value="2">{{ t('edit') }}</a-radio>
                <a-radio value="1">{{ t('view') }}</a-radio>
              </a-radio-group>
            </template>
          </a-table-column>
        </a-table>
      </template>
    </div>
    <div class="form-wrapper" v-else-if="object_type == 3">
      <div class="select">{{ t('selected_count', { count: checkedList.length }) }}</div>

      <a-table
        :data-source="dataList"
        :pagination="false"
        :scroll="{ y: 500 }"
        row-key="id"
        :row-selection="{ selectedRowKeys: checkedList, onChange: onSelectChange }"
      >
        <a-table-column :title="t('database')" data-index="name" :width="350">
          <template #default="{ record }">
            <div class="tools-wrapper">
              <div class="item-name">{{ record.name }}</div>
            </div>
          </template>
        </a-table-column>
        <a-table-column :title="t('permission')" data-index="operate_rights" :width="350">
          <template #default="{ record }">
            <a-radio-group v-model:value="record.operate_rights">
              <a-radio value="4">{{ t('manage') }}</a-radio>
              <a-radio value="2">{{ t('edit') }}</a-radio>
              <a-radio value="1">{{ t('view') }}</a-radio>
            </a-radio-group>
          </template>
        </a-table-column>
      </a-table>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { getPermissionManageList, batchSavePermissionManage } from '@/api/department/index.js'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.manage.components.see-model-alert')
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
const baseOption = computed(() => [
  {
    value: 2,
    payload: {
      title: t('department_knowledge_base')
    }
  },
  {
    value: 1,
    payload: {
      title: t('member_knowledge_base')
    }
  }
])
const tabOption = ref(baseOption.value)

// 监听 baseOption 的变化，当语言切换时更新 tabOption
watch(baseOption, (newVal) => {
  tabOption.value = newVal
}, { immediate: true })

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
    currentTitle.value = t('add_member_robot')
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
      currentTitle.value = t('edit_member_robot')
      checkedList.value = formatId(record.managed_robot_list)
    }
  } else if (key == 2) {
    console.log(record, '==')
    currentTitle.value = t('add_member_knowledge')
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
      currentTitle.value = t('edit_member_knowledge')
      checkedList.value = formatId(record.managed_library_list)
    }
  } else if (key == 3) {
    currentTitle.value = t('add_member_database')

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
      currentTitle.value = t('edit_member_database')
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
    currentTitle.value = t('batch_update_member_robot')
    list = props.robotList
  } else if (key == 2) {
    currentTitle.value = t('batch_update_member_knowledge')
    list = props.libraryList
  } else if (key == 3) {
    currentTitle.value = t('batch_update_member_database')
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

  // 创建响应式的 tab 选项
  const dynamicTabOption = computed(() => [
    {
      value: 2,
      payload: {
        title: identity_type.value == 1 ? t('department_knowledge_base') : t('superior_knowledge_base')
      }
    },
    {
      value: 1,
      payload: {
        title: identity_type.value == 1 ? t('member_knowledge_base') : t('department_knowledge_base')
      }
    }
 ])
  tabOption.value = dynamicTabOption.value
  
  // 监听动态选项的变化
  watch(dynamicTabOption, (newVal) => {
    tabOption.value = newVal
  }, { immediate: true })
  tabs.value = 1
  object_type.value = key

  let list = []
  if (key == 1) {
    currentTitle.value = t('add_member_robot')
    list = props.robotList
  } else if (key == 2) {
    currentTitle.value = t('add_member_knowledge')
    list = props.libraryList
  } else if (key == 3) {
    currentTitle.value = t('add_member_database')
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
    message.success(t('save_success'))
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
