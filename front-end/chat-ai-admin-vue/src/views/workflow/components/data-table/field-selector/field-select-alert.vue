<style lang="less" scoped>
.library-checkbox-box {
  padding-top: 16px;

  .list-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
    margin-bottom: 8px;
  }

  .list-box {
    display: flex;
    flex-flow: row wrap;
    height: 388px;
    width: 100%;
    overflow-y: auto;
    align-content: flex-start;
    margin: 0 -8px;

    .list-item-wraapper {
      padding: 8px;
      width: 50%;
    }

    .list-item {
      width: 100%;
      padding: 14px 12px;
      border: 1px solid #f0f0f0;
      border-radius: 2px;

      &:hover {
        cursor: pointer;
        box-shadow: 0 4px 16px 0 #1b3a6929;
      }

      .library-name {
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .library-desc {
        line-height: 20px;
        margin-top: 2px;
        font-size: 12px;
        font-weight: 400;
        color: #8c8c8c;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }

    .list-item :deep(span:last-child) {
      flex: 1;
      overflow: hidden;
    }
  }
}
.empty-box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding-top: 40px;
  padding-bottom: 40px;
  color: #8c8c8c;
  img {
    width: 150px;
    height: 150px;
  }
}
</style>

<template>
  <a-modal width="746px" v-model:open="state.show" title="关联数据表" @ok="handleOk" ok-text="确定" cancel-text="取消">
    <div class="field-select-box">
      <div class="field-list-box" v-if="dataSource.length">
        <a-table rowKey="name" :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }" :dataSource="dataSource" :columns="columns" :loading="loading" :pagination="false" />
      </div>
      <div class="empty-box" v-if="!dataSource.length">
        <img src="@/assets/img/library/preview/empty.png" alt="" />
        <div>
          暂无字段, 请先去数据表添加
          <a @click="openAddLibrary"> 去添加</a>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useDataTableStore } from '@/stores/modules/data-table'

const emit = defineEmits(['ok'])

const dataTableStore = useDataTableStore()

const state = reactive({
  show: false,
  formId: '',
  selectedRowKeys: [],
  selectedRows: []
})

const columns = ref([
  {
    title: '字段名称',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '数据类型',
    dataIndex: 'type',
    key: 'type'
  },
  {
    title: '字段描述',
    dataIndex: 'description',
    key: 'description'
  }
])

const dataSource = ref([])

const onSelectChange = (selectedRowKeys, selectedRows) => {
  state.selectedRowKeys = selectedRowKeys;
  state.selectedRows = selectedRows
};

const open = ({
  formId,
  selectedRowKeys,
  selectedRows
}) => {
  state.formId = formId
  state.selectedRowKeys = selectedRowKeys
  state.selectedRows = selectedRows

  getList(selectedRowKeys)
  
  state.show = true
}

const openAddLibrary = () => {
  window.open('/#/database/details/field-manage?form_id='+state.formId)
}

const handleOk = () => {
  emit('ok', state.selectedRowKeys, state.selectedRows)
  state.show = false;
}

const loading = ref(false)
const getList = (selectedRowKeys) => {
  loading.value = true
  dataTableStore.getFormFieldList({ form_id: state.formId})
    .then((list) => {
      loading.value = false

      // list = list.filter(item => !selectedRowKeys.includes(item.name))

      list.forEach(item => {
        item.value = ''
        item.id = item.id * 1
      })

      dataSource.value = list
    })
    .catch(() => {
      loading.value = false
    })
}

defineExpose({
  open
})
</script>
