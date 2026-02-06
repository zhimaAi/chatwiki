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
  <a-modal width="746px" v-model:open="state.show" :title="t('title_relate_data_table')" @ok="handleOk" :ok-text="t('btn_confirm')" :cancel-text="t('btn_cancel')">
    <div class="field-select-box">
      <div class="field-list-box" v-if="fields.length">
        <a-table
          rowKey="field_name"
          :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
          :dataSource="fields"
          :columns="columns"
          :pagination="false"/>
      </div>
      <div class="empty-box" v-if="!fields.length">
        <img src="@/assets/img/library/preview/empty.png" alt=""/>
        <div>{{ t('msg_no_fields_add') }}</div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.feishu-table.sort-selector.field-select-alert')

const emit = defineEmits(['ok'])
const props = defineProps({
  fields: {
    type: Array,
    default: () => ([])
  }
})

const state = reactive({
  show: false,
  selectedRowKeys: [],
  selectedRows: []
})

const columns = ref([
  {
    title: t('label_field_name'),
    dataIndex: 'field_name',
  },
  {
    title: t('label_data_type'),
    dataIndex: 'ui_type',
  },
])

const onSelectChange = (selectedRowKeys, selectedRows) => {
  state.selectedRowKeys = selectedRowKeys;
  state.selectedRows = selectedRows
};

const open = ({selectedRowKeys, selectedRows}) => {
  state.selectedRowKeys = selectedRowKeys
  state.selectedRows = selectedRows

  state.show = true
}


const handleOk = () => {
  emit('ok', state.selectedRowKeys, state.selectedRows)
  state.show = false;
}

defineExpose({
  open
})
</script>