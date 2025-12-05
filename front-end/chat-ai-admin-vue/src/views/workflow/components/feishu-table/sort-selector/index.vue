<style lang="less" scoped>
.field-list {
  .field-list-row {
    position: relative;
    display: flex;
    flex-wrap: nowrap;
    align-items: center;
    margin-bottom: 4px;
    &:last-child {
      margin-bottom: 0;
    }
  }
  .field-list-col {
    padding: 0 4px;
  }
  .field-value-col {
    width: 208px;
  }
  .field-name-col {
    width: 220px;
    line-height: 22px;
    font-size: 14px;
    color: #595959;
  }
  .field-del-head,
  .field-del-col {
    width: 24px;
    display: flex;
    align-items: center;
  }
  .field-del-col {
    .del-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }
  .field-drag-col {
    width: 24px;
    margin-right: 8px;
    .drag-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 24px;
      height: 24px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
      &:hover {
        border-radius: 6px;
        background: #d9d9d9;
      }
    }
  }
  .field-name-box {
    display: flex;
    .field-type-tag {
      padding: 0 6px;
      margin-left: 8px;
      border-radius: 6px;
      font-size: 12px;
      line-height: 16px;
      font-weight: 400;
      color: rgb(89, 89, 89);
      border: 1px solid rgba(0, 0, 0, 0.15);
      background: #fff;
    }
  }
}
.field-list-row:hover {
  border-radius: 6px;
  background: #e4e6eb;
}
.add-btn-box {
  margin-top: 8px;
}
</style>

<template>
  <div>
    <div class="field-list">
      <div class="field-list-row" v-for="item in conditions" :key="item.name">
        <div class="field-list-col field-drag-col">
          <span class="drag-btn">
            <HolderOutlined />
          </span>
        </div>
        <div class="field-list-col field-name-col">
          <div class="field-name-box">
            <span class="field-name">{{ item.field_name }}</span>
            <span class="field-type-tag">{{ item.ui_type }}</span>
          </div>
        </div>
        <div class="field-list-col field-value-col">
          <a-select v-model:value="item.is_asc" style="width: 100%" @change="changeValue">
            <a-select-option :value="opt.value" v-for="opt in operatorOptions" :key="opt.key">
              {{ opt.label }}
            </a-select-option>
          </a-select>
        </div>
        <div class="field-list-col field-del-col" v-if="props.showDelete">
          <span class="del-btn" @click="removeCondition(index)">
            <svg-icon class="del-icon" name="close-circle"></svg-icon>
          </span>
        </div>
      </div>
    </div>

    <div class="add-btn-box">
      <a-tooltip title="请先选择数据库" style="width: 100%" v-if="disabled">
        <span>
          <a-button class="add-btn" type="dashed" disabled block>
            <PlusOutlined /> 添加排序字段
          </a-button>
        </span>
      </a-tooltip>

      <a-button class="add-btn" type="dashed" block @click="handleAddField" v-else>
        <PlusOutlined /> 添加排序字段
      </a-button>
    </div>

    <FieldSelectAlert :fields="fields" ref="fieldSelectAlertRef" @ok="onChangeSelect" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { PlusOutlined, HolderOutlined } from '@ant-design/icons-vue'
import FieldSelectAlert from './field-select-alert.vue'

const emit = defineEmits(['change'])

const props = defineProps({
  showAdd: {
    type: Boolean,
    default: true
  },
  showDelete: {
    type: Boolean,
    default: true
  },
  fields: {
    type: Array,
    default: () => []
  },
  list: {
    type: Array,
    default: () => []
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const operatorOptions = [
  { label: '升序', value: 1, key: 'asc' },
  { label: '降序', value: 0, key: 'desc' }
]

const fieldSelectAlertRef = ref()
const conditions = ref([])

watch(
  () => props.list,
  (newVal) => {
    conditions.value = newVal.map(item => {
      return {
        ...item,
        is_asc: item.is_asc ? 1 : 0
      }
    })
  }
)

const removeCondition = (index) => {
  conditions.value.splice(index, 1)

  change()
}

const onChangeSelect = (selectedRowKeys, selectedRows) => {
  conditions.value = selectedRows.map((item) => {
    return { ...item }
  })
  change()
}

const changeValue = () => {
  change()
}

const change = () => {
  let data = conditions.value.map(item => {
    return {
      ...item,
      is_asc: item.is_asc == 1 ? true : false
    }
  })
  emit('change', data)
}

const handleAddField = () => {
  fieldSelectAlertRef.value.open({
    selectedRows: props.list,
    selectedRowKeys: props.list.map((item) => {
      return item.field_name
    })
  })
}
</script>
