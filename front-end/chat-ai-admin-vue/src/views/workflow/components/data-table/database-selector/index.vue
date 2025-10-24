<style lang="less" scoped>
.select-data-btn {
  width: 100%;
}
.data-table-info {
  position: relative;
  padding: 14px 16px;
  border-radius: 6px;
  border: 1px solid #d8dde5;
  background: #fff;

  .data-table-name {
    font-size: 14px;
    line-height: 22px;
    font-weight: 600;
    color: #262626;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .data-table-desc {
    font-size: 12px;
    line-height: 20px;
    font-weight: 400;
    color: #8c8c8c;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .del-btn {
    position: absolute;
    right: 16px;
    top: 28px;
    font-size: 16px;
    color: #595959;
    cursor: pointer;
  }
}
</style>

<template>
  <div class="">
    <div class="data-table-info" v-if="props.value" >
      <div class="data-table-name">{{ selectedRow ? selectedRow.name : props.title }}</div>
      <div class="data-table-desc">{{ selectedRow ? selectedRow.description : props.description }}</div>
      <svg-icon class="del-btn" name="close-circle" @click="handleClear"></svg-icon>
    </div>

    <a-button class="select-data-btn" type="primary" ghost @click="handleSelectData" v-else><PlusOutlined /> 选择数据表</a-button>

    <TableSelectAlert ref="tableSelectAlertRef" @ok="onSelectOk" :options="options" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import TableSelectAlert from './table-select-alert.vue'

const emit = defineEmits(['update:value', 'change'])

const props = defineProps({
  value: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    default: () => []
  },
  description: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  }
})

const tableSelectAlertRef = ref(null)

const selectedRow = computed(() => {
  let arr = props.options.filter(item => item.id === props.value)

  if (arr.length) {
    return arr[0]
  }

  return null
})


const handleSelectData = () => {
  tableSelectAlertRef.value.open()
}

const onSelectOk = (value, record ) => {
  emit('update:value', value)
  emit('change', value, record)
}

const handleClear = () => {
  emit('update:value', '')
  emit('change', '')
}

</script>
